package service

import (
	"context"
	"fmt"
	"time"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/repository"
)

type ActivitySearchService struct {
	ActivityRepo *repository.ActivitySearchRepo
}

type ActivitySearchResponse struct {
	Results      []model.ActivitySearchResult `json:"results"`
	Date         time.Time                    `json:"date"`
	Total        int64                          `json:"total"`
	Participants model.TravelerCount          `json:"participants"`
	Page         int                          `json:"page"`
	PageSize     int                          `json:"pageSize"`
}

func NewActivitySearchService(repo *repository.ActivitySearchRepo) *ActivitySearchService{
	return &ActivitySearchService{
		ActivityRepo: repo,
	}
}


func (s *ActivitySearchService) Search(ctx context.Context, params model.ActivitySearch)(*ActivitySearchResponse, error){
	if err := params.Validate(); err != nil{
		return nil, fmt.Errorf("Invalid search params: %w", err)
	}

	filter := repository.ActivityFilter{
		Location: params.Location,
		Participants: params.Participants,
		ForAllAges: params.ForAllAges,
		Category: params.Category,
		Date: params.Date,
		Duration: params.MaxDurationMinutes,
		SortBy: params.SortBy,
		Page: params.Page,
		PageSize: params.PageSize,
	}

	result, err := s.ActivityRepo.SearchActivity(ctx, &filter)
	if err != nil{
		return nil, fmt.Errorf("accommodation search error: %w", err)
	}

	return &ActivitySearchResponse{
		Results: result,
		Total: int64(len(result)),
		Participants: params.Participants,
		Page: params.Page,
		Date: params.Date,
		PageSize: params.PageSize,

	}, nil
}
