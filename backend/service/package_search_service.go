package service

import (
	"context"
	"fmt"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/repository"
	"github.com/Loboo34/travel/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackageSearchService struct {
	packageRepo       *repository.PackageSearchRepo
	activityRepo      *repository.ActivityRepo
	accommodationRepo *repository.AccommodationSearchRepo
	flightRepo        *repository.FlightSearchRepo
}

type PackageSearchResult struct {
	Package          model.Package     `json:"package"`
	TotalPrice       int64             `json:"totalPrice"`
	AvailableSlots   int               `json:"availableSlots"`
	ComponentDetails []ComponentDetail `json:"componentDetails"`
}

type ComponentDetail struct {
	Type        model.Component    `json:"type"`
	ReferenceID primitive.ObjectID `json:"referenceID"`
	Title       string             `json:"title"`
	Price       int64              `json:"price"`
	Required    bool               `json:"required"`
}

type PackageSearchResponse struct {
	Results  []PackageSearchResult `json:"results"`
	Total    int                   `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"pageSize"`
}

func NewPackageSearchService(
	packageRepo *repository.PackageSearchRepo,
	activityRepo *repository.ActivityRepo,
	accommodationRepo *repository.AccommodationSearchRepo,
	flightRepo *repository.FlightSearchRepo,
) *PackageSearchService {
	return &PackageSearchService{
		packageRepo:       packageRepo,
		activityRepo:      activityRepo,
		accommodationRepo: accommodationRepo,
		flightRepo:        flightRepo,
	}
}

func (s *PackageSearchService) Search(
	ctx context.Context,
	params model.PackageSearchParams,
) (*PackageSearchResponse, error) {

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("invalid search params: %w", err)
	}

	filter := &repository.PackageFilter{
		Destination: params.Destination,
		StartDate:   params.StartDate,
		EndDate:     params.EndDate,
		Travelers:   params.Travelers,
		MaxBudget:   params.MaxBudget,
		Tags:        params.Tags,
		Components:  params.Components,
		SortBy:      params.SortBy,
		Page:        params.Page,
		PageSize:    params.PageSize,
	}

	candidates, err := s.packageRepo.FindCandidates(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("finding package candidates: %w", err)
	}

	// return early if nothing matched the catalog query
	if len(candidates) == 0 {
		return &PackageSearchResponse{
			Results:  []PackageSearchResult{},
			Page:     params.Page,
			PageSize: params.PageSize,
		}, nil
	}

	// step 3 — check component availability concurrently for each candidate
	type availResult struct {
		result *PackageSearchResult
		err    error
	}

	resultCh := make(chan availResult, len(candidates))

	for _, pkg := range candidates {
		pkg := pkg // capture loop variable
		go func() {
			result, err := s.checkPackageAvailability(ctx, pkg, params)
			resultCh <- availResult{result, err}
		}()
	}

	// collect — only keep packages where all required components are available
	var results []PackageSearchResult
	for range candidates {
		r := <-resultCh
		if r.err != nil {
			// log but don't fail the entire search
			utils.Logger.Warn("package availability check failed")
			continue
		}
		if r.result != nil {
			results = append(results, *r.result)
		}
	}

	return &PackageSearchResponse{
		Results:  results,
		Total:    len(results),
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

// checkPackageAvailability verifies package slots and all component availability
// returns nil, nil if the package should be excluded from results
func (s *PackageSearchService) checkPackageAvailability(
	ctx context.Context,
	pkg model.Package,
	params model.PackageSearchParams,
) (*PackageSearchResult, error) {

	// check package-level slot availability first
	avail, err := s.packageRepo.GetAvailability(ctx, pkg.ID)
	if err != nil {
		return nil, fmt.Errorf("package availability lookup: %w", err)
	}
	// no availability record or not enough slots
	if avail == nil || avail.AvailableSlots() < params.Travelers.Total() {
		return nil, nil
	}

	// check each component concurrently
	type componentResult struct {
		detail   *ComponentDetail
		required bool
		err      error
	}

	componentCh := make(chan componentResult, len(pkg.IncludedComponents))

	for _, comp := range pkg.IncludedComponents {
		comp := comp
		go func() {
			detail, err := s.checkComponent(ctx, comp, params)
			componentCh <- componentResult{detail, comp.Required, err}
		}()
	}

	var details []ComponentDetail
	var totalPrice int64

	for range pkg.IncludedComponents {
		r := <-componentCh
		if r.err != nil {
			return nil, r.err
		}
		if r.detail == nil && r.required {
			// required component unavailable — exclude entire package
			return nil, nil
		}
		if r.detail != nil {
			totalPrice += r.detail.Price
			details = append(details, *r.detail)
		}
	}

	// must have resolved at least one component
	if len(details) == 0 {
		return nil, nil
	}

	return &PackageSearchResult{
		Package:          pkg,
		TotalPrice:       totalPrice,
		AvailableSlots:   avail.AvailableSlots(),
		ComponentDetails: details,
	}, nil
}

// checkComponent routes to the correct repo based on component type
func (s *PackageSearchService) checkComponent(
	ctx context.Context,
	comp model.PackageComponent,
	params model.PackageSearchParams,
) (*ComponentDetail, error) {

	switch comp.ComponentType {
	case model.ComponentFlight:
		return s.checkFlight(ctx, comp, params)
	case model.ComponentAccommodation:
		return s.checkAccommodation(ctx, comp, params)
	case model.ComponentActivity:
		return s.checkActivity(ctx, comp, params)
	default:
		return nil, fmt.Errorf("unknown component type: %q", comp.ComponentType)
	}
}

func (s *PackageSearchService) checkFlight(
	ctx context.Context,
	comp model.PackageComponent,
	params model.PackageSearchParams,
) (*ComponentDetail, error) {

	offer, err := s.flightRepo.FindActiveOffer(ctx, comp.ReferenceID, params.Travelers.Total())
	if err != nil {
		return nil, fmt.Errorf("flight availability: %w", err)
	}
	if offer == nil {
		return nil, nil // unavailable
	}
	return &ComponentDetail{
		Type:        model.ComponentFlight,
		ReferenceID: comp.ReferenceID,
		Title:       offer.ProviderReference,
		Price:       offer.PriceTotal,
		Required:    comp.Required,
	}, nil
}

func (s *PackageSearchService) checkAccommodation(
	ctx context.Context,
	comp model.PackageComponent,
	params model.PackageSearchParams,
) (*ComponentDetail, error) {

	available, totalPrice, err := s.accommodationRepo.CheckAvailability(
		ctx,
		comp.ReferenceID,
		params.StartDate,
		params.EndDate,
		params.Travelers.Total(),
	)
	if err != nil {
		return nil, fmt.Errorf("accommodation availability: %w", err)
	}
	if !available {
		return nil, nil
	}
	return &ComponentDetail{
		Type:        model.ComponentAccommodation,
		ReferenceID: comp.ReferenceID,
		Price:       totalPrice,
		Required:    comp.Required,
	}, nil
}

func (s *PackageSearchService) checkActivity(
	ctx context.Context,
	comp model.PackageComponent,
	params model.PackageSearchParams,
) (*ComponentDetail, error) {

	slot, err := s.activityRepo.FindAvailableTimeslot(
		ctx,
		comp.ReferenceID,
		params.StartDate,
		params.Travelers.Total(),
	)
	if err != nil {
		return nil, fmt.Errorf("activity availability: %w", err)
	}
	if slot == nil {
		return nil, nil
	}
	return &ComponentDetail{
		Type:        model.ComponentActivity,
		ReferenceID: comp.ReferenceID,
		Price:       slot.PricePerPerson * int64(params.Travelers.Total()),
		Required:    comp.Required,
	}, nil
}
