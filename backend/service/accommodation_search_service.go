package service

import (
	"context"
	"fmt"
	"time"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/repository"
)



type AccommodationSearchService struct {
	AccommodationRepo *repository.AccommodationSearchRepo
}

type AccommodationSearchResponse struct {
	Results      []model.AccommodationSearchResult `json:"results"`
	Total        int64                             `json:"total"`
	CheckInTime  time.Time                         `json:"checkInTime"`
	CheckOutTime time.Time                         `json:"checkOutTime"`
	Nights       int                               `json:"nights"`
	Page         int                               `json:"page"`
	PageSize     int                               `json:"pageSize"`
}

func NewAccommodationSearchService(repo *repository.AccommodationSearchRepo) *AccommodationSearchService {
	return &AccommodationSearchService{
		AccommodationRepo: repo,
	}
}
func (s *AccommodationSearchService) SearchWithFilter(ctx context.Context, f *repository.AccommodationFilter) ([]model.AccommodationSearchResult, error) {

	if f.Location.Latitude != 0 || f.Location.Longitude != 0 {
		ids, err := s.AccommodationRepo.FindIDsNearLocation(ctx, f.Location.Latitude, f.Location.Longitude, f.Location.RadiusKm)
		if err != nil {
			return nil, err
		}
		if len(ids) == 0 {
			return []model.AccommodationSearchResult{}, nil
		}
		f.CandidateIDs = ids
	}

	return s.AccommodationRepo.SearchAccommodationAvailability(ctx, f)
}

func (s *AccommodationSearchService) Search(ctx context.Context, params model.AccommodationSearch) (*AccommodationSearchResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("Invalid search params: %w", err)
	}

	nights := int(params.CheckOutDate.Sub(params.ChecKInDate).Hours() / 24)

	filter := repository.AccommodationFilter{
		Location:     params.Location,
		CheckInDate:  params.ChecKInDate,
		CheckOutDate: params.CheckOutDate,
		Guests:       params.Guests,
		PropertyType: params.PropertyType,
		Rooms:        params.TotalRooms,
		SortBy:       params.SortBy,
		Page:         params.Page,
		PageSize:     params.PageSize,
	}

	results, err := s.SearchWithFilter(ctx, &filter)
	if err != nil {
		return nil, fmt.Errorf("accommodation search: %w", err)
	}

	//  for i := range results {
	//     if nights > 0 {
	//         results[i].PricePerNight = results[i].TotalPrice / int64(nights)
	//     }
	// }

	return &AccommodationSearchResponse{
		Results:      results,
		Total:        int64(len(results)),
		Page:         params.Page,
		PageSize:     params.PageSize,
		Nights:       nights,
		CheckInTime:  params.ChecKInDate,
		CheckOutTime: params.CheckOutDate,
	}, nil
}
