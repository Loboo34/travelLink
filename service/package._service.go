package service

import (
	"context"
	"fmt"
	"time"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackageService struct {
	packageRepo *repository.PackageRepo
}

func NewPackageService(packageRepo *repository.PackageRepo) *PackageService {
	return &PackageService{packageRepo: packageRepo}
}

type PackageRequest struct {
	Name               string                   `json:"name"`
	Description        string                   `json:"description"`
	Destination        string                   `json:"destination"`
	DurationDays       int                      `json:"durationDays"`
	StartDateFrom      *time.Time               `json:"startDateFrom,omitempty"`
	StartDateTo        *time.Time               `json:"startDateTo,omitempty"`
	MaxTravelers       int                      `json:"maxTravelers"`
	BasePrice          int64                    `json:"basePrice"`
	Currency           string                   `json:"currency"`
	IncludedComponents []model.PackageComponent `json:"includedComponents"`
	Tags               []model.PackageTag       `json:"tags,omitempty"`
	Images             []string                 `json:"images,omitempty"`
	ExpiresAt          *time.Time               `json:"expiresAt,omitempty"`
}

type PackageResult struct {
	PackageID   primitive.ObjectID `json:"packageID"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Destination string             `json:"destination"`
	BasePrice   int64              `json:"basePrice"`
	Currency    string             `json:"currency"`
}

type PackageUpdateRequest struct {
	Name               string                   `json:"name"`
	Description        string                   `json:"description"`
	Destination        string                   `json:"destination"`
	DurationDays       int                      `json:"durationDays"`
	MaxTravelers       int                      `json:"maxTravelers"`
	BasePrice          int64                    `json:"basePrice"`
	IncludedComponents []model.PackageComponent `json:"includedComponents"`
	Tags               []model.PackageTag       `json:"tags,omitempty"`
	Images             []string                 `json:"images,omitempty"`
}

func (s *PackageService) Create(ctx context.Context, req PackageRequest) (*PackageResult, error) {
	packageID := primitive.NewObjectID()
	now := time.Now()

	pkg := &model.Package{
		ID:                 packageID,
		Name:               req.Name,
		Description:        req.Description,
		Destination:        req.Destination,
		DurationDays:       req.DurationDays,
		StartDateFrom:      req.StartDateFrom,
		StartDateTo:        req.StartDateTo,
		MaxTravelers:       req.MaxTravelers,
		BasePrice:          req.BasePrice,
		Currency:           req.Currency,
		IncludedComponents: req.IncludedComponents,
		Tags:               req.Tags,
		Images:             req.Images,
		ExpiresAt:          req.ExpiresAt,
		Rating:             0,
		ReviewCount:        0,
		IsActive:           true,
		CreatedAt:          now,
		UpdatedAt:          now,
		CachedAt:           now,
	}

	if err := s.packageRepo.Add(ctx, pkg); err != nil {
		return nil, fmt.Errorf("creating package: %w", err)
	}

	return &PackageResult{
		PackageID:   packageID,
		Name:        req.Name,
		Description: req.Description,
		Destination: req.Destination,
		BasePrice:   req.BasePrice,
		Currency:    req.Currency,
	}, nil
}

func (s *PackageService) Update(ctx context.Context, packageID primitive.ObjectID, req PackageUpdateRequest) (*PackageResult, error) {
	if err := s.packageRepo.Update(ctx, packageID, req.Name, req.Description, req.Destination, req.DurationDays, req.MaxTravelers, req.BasePrice, req.IncludedComponents, req.Tags, req.Images); err != nil {
		return nil, fmt.Errorf("updating package: %w", err)
	}

	return &PackageResult{
		PackageID:   packageID,
		Name:        req.Name,
		Description: req.Description,
		Destination: req.Destination,
		BasePrice:   req.BasePrice,
	}, nil
}

func (s *PackageService) SetActive(ctx context.Context, packageID primitive.ObjectID, active bool) error {
	if err := s.packageRepo.SetActive(ctx, packageID, active); err != nil {
		return fmt.Errorf("setting package active state: %w", err)
	}
	return nil
}

func (s *PackageService) Delete(ctx context.Context, packageID primitive.ObjectID) error {
	if err := s.packageRepo.Delete(ctx, packageID); err != nil {
		return fmt.Errorf("deleting package: %w", err)
	}
	return nil
}

func (s *PackageService) GetByID(ctx context.Context, packageID primitive.ObjectID) (*model.Package, error) {
	pkg, err := s.packageRepo.GetByID(ctx, packageID)
	if err != nil {
		return nil, fmt.Errorf("fetching package: %w", err)
	}
	return pkg, nil
}
