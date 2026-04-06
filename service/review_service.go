package service

import (
	"context"
	"fmt"
	"time"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewService struct {
	reviewRepo *repository.ReviewRepo
}

func NewReviewService(reviewRepo *repository.ReviewRepo) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo}
}

type ReviewRequest struct {
	ReferenceID primitive.ObjectID `json:"referenceID"`
	ReviewFor   model.ReviewFor    `json:"reviewFor"`
	Rating      int                `json:"rating"`
	Content     string             `json:"content"`
}

type UpdateReviewRequest struct {
	Rating  int    `json:"rating"`
	Content string `json:"content"`
}

func (s *ReviewService) Create(ctx context.Context, userID primitive.ObjectID, req ReviewRequest) (*model.Review, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, &model.ValidationError{Message: "rating must be between 1 and 5"}
	}
	if req.Content == "" {
		return nil, &model.ValidationError{Message: "review content is required"}
	}
	switch req.ReviewFor {
	case model.ReviewForAccommodation, model.ReviewForActivity, model.ReviewForFlight, model.ReviewForPackage:
	default:
		return nil, &model.ValidationError{Message: "invalid reviewFor value"}
	}

	now := time.Now().UTC()
	review := &model.Review{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		ReferenceID: req.ReferenceID,
		ReviewFor:   req.ReviewFor,
		Rating:      req.Rating,
		Content:     req.Content,
		IsVerified:  false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	col := targetCollection(req.ReviewFor)
	if err := s.reviewRepo.UpdateTargetRating(ctx, col, req.ReferenceID, req.Rating, 1); err != nil {
		fmt.Printf("warn: updating target rating after create: %v\n", err)
	}

	return review, nil
}

func (s *ReviewService) GetByTarget(ctx context.Context, referenceID primitive.ObjectID, reviewFor model.ReviewFor) ([]model.Review, error) {
	reviews, err := s.reviewRepo.GetByTarget(ctx, referenceID, reviewFor)
	if err != nil {
		return nil, fmt.Errorf("fetching reviews: %w", err)
	}
	return reviews, nil
}

func (s *ReviewService) Update(ctx context.Context, reviewID primitive.ObjectID, userID primitive.ObjectID, req UpdateReviewRequest) (*model.Review, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, &model.ValidationError{Message: "rating must be between 1 and 5"}
	}
	if req.Content == "" {
		return nil, &model.ValidationError{Message: "review content is required"}
	}

	existing, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return nil, err
	}
	if existing.UserID != userID {
		return nil, &model.AuthError{Message: "you can only edit your own reviews"}
	}

	if err := s.reviewRepo.Update(ctx, reviewID, req.Rating, req.Content); err != nil {
		return nil, err
	}

	
	col := targetCollection(existing.ReviewFor)
	ratingDiff := req.Rating - existing.Rating
	if ratingDiff != 0 {
		if err := s.reviewRepo.UpdateTargetRating(ctx, col, existing.ReferenceID, ratingDiff, 1); err != nil {
			fmt.Printf(" updating target rating after update: %v\n", err)
		}
	}

	existing.Rating = req.Rating
	existing.Content = req.Content
	existing.UpdatedAt = time.Now().UTC()
	return existing, nil
}

func (s *ReviewService) Delete(ctx context.Context, reviewID primitive.ObjectID, userID primitive.ObjectID, isAdmin bool) error {
	existing, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return err
	}
	if !isAdmin && existing.UserID != userID {
		return &model.AuthError{Message: "you can only delete your own reviews"}
	}

	if err := s.reviewRepo.Delete(ctx, reviewID); err != nil {
		return err
	}

	col := targetCollection(existing.ReviewFor)
	if err := s.reviewRepo.UpdateTargetRating(ctx, col, existing.ReferenceID, existing.Rating, -1); err != nil {
		fmt.Printf("warn: updating target rating after delete: %v\n", err)
	}

	return nil
}

func targetCollection(rf model.ReviewFor) string {
	switch rf {
	case model.ReviewForAccommodation:
		return "accommodations"
	case model.ReviewForActivity:
		return "activities"
	case model.ReviewForFlight:
		return "flights"
	case model.ReviewForPackage:
		return "packages"
	default:
		return ""
	}
}