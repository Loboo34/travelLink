package service

import (
	"context"
	"fmt"
	"time"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/repository"
)

type AccommodationService struct {
	AccommodationRepo *repository.AccommodationRepo
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

func NewAccommodationService(repo *repository.AccommodationRepo) *AccommodationService {
	return &AccommodationService{
		AccommodationRepo: repo,
	}
}

func (a *AccommodationService) Search(ctx context.Context, params model.AccommodationSearch) (*AccommodationSearchResponse, error) {
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

	results, err := a.AccommodationRepo.SearchAccommodationAvailability(ctx, &filter)
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
