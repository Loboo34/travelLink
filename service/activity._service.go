package service

import (
	"context"
	"fmt"
	"time"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivityService struct {
	activityRepo *repository.ActivityRepo
}

func NewActivityService(activityRepo *repository.ActivityRepo) *ActivityService {
	return &ActivityService{activityRepo: activityRepo}
}

type ActivityRequest struct {
	Title           string                   `json:"title"`
	Description     string                   `json:"description"`
	City            string                   `json:"city"`
	Country         string                   `json:"country"`
	Location        model.GeoLocation        `json:"location"`
	MeetingPoint    model.MeetingPoint       `json:"meetingPoint"`
	Categories      []model.ActivityCategory `json:"categories"`
	DurationMinutes int                      `json:"durationMinutes"`
	Inclusions      []string                 `json:"inclusions"`
	Exclusions      []string                 `json:"exclusions"`
	Images          []string                 `json:"images"`
}

type ActivityResult struct {
	ActivityID  primitive.ObjectID `json:"activityID"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	City        string             `json:"city"`
	Country     string             `json:"country"`
}

type ActivityTimeslotRequest struct {
	StartTime       time.Time `json:"startTime"`
	DurationMinutes int       `json:"durationMinutes"`
	TotalSlots      int       `json:"totalSlots"`
	PricePerPerson  int64     `json:"pricePerPerson"`
	GroupSizeMax    int       `json:"groupSizeMax"`
}

type ActivityTimeslotResult struct {
	TimeSlotID primitive.ObjectID `json:"timeSlotID"`
	ActivityID primitive.ObjectID `json:"activityID"`
}

func (s *ActivityService) Create(ctx context.Context, req ActivityRequest) (*ActivityResult, error) {

	activityID := primitive.NewObjectID()
	now := time.Now()

	activity := &model.Activity{
		ID:              activityID,
		Title:           req.Title,
		Description:     req.Description,
		City:            req.City,
		Country:         req.Country,
		Location:        req.Location,
		MeetingPoint:    req.MeetingPoint,
		Categories:      req.Categories,
		DurationMinutes: req.DurationMinutes,
		Inclusions:      req.Inclusions,
		Exclusions:      req.Exclusions,
		Images:          req.Images,
		Rating:          0,
		ReviewCount:     0,
		CreatedAt:       now,
		UpdatedAt:       now,
		IsActive:        true,
	}

	if err := s.activityRepo.Add(ctx, activity); err != nil {
		return nil, fmt.Errorf("creating activity: %w", err)
	}

	return &ActivityResult{
		ActivityID:  activityID,
		Title:       req.Title,
		Description: req.Description,
		City:        req.City,
		Country:     req.Country,
	}, nil
}

func (s *ActivityService) Update(ctx context.Context, activityID primitive.ObjectID, req ActivityRequest) (*ActivityResult, error) {
	if err := s.activityRepo.Update(ctx, activityID, req.Title, req.DurationMinutes, req.Inclusions, req.Exclusions, &req.MeetingPoint); err != nil {
		return nil, fmt.Errorf("updating activity: %w", err)
	}

	return &ActivityResult{
		ActivityID:  activityID,
		Title:       req.Title,
		Description: req.Description,
		City:        req.City,
		Country:     req.Country,
	}, nil
}

func (s *ActivityService) Delete(ctx context.Context, activityID primitive.ObjectID) error {
	if err := s.activityRepo.Delete(ctx, activityID); err != nil {
		return fmt.Errorf("deleting activity: %w", err)
	}
	return nil
}

func (s *ActivityService) SetActivityActive(ctx context.Context, activityID primitive.ObjectID, active bool) error {
	if err := s.activityRepo.SetActive(ctx, activityID, active); err != nil {
		return fmt.Errorf("setting activity active state: %w", err)
	}
	return nil
}

func (s *ActivityService) CreateTimeSlot(ctx context.Context, activityID primitive.ObjectID, req ActivityTimeslotRequest) (*ActivityTimeslotResult, error) {
	if activityID.IsZero() {
		return nil, fmt.Errorf("activityID is required")
	}

	timeSlotID := primitive.NewObjectID()
	now := time.Now()

	timeSlot := &model.ActivityTimeslot{
		ID:              timeSlotID,
		ActivityID:      activityID,
		StartTime:       req.StartTime,
		DurationMinutes: req.DurationMinutes,
		TotalSlots:      req.TotalSlots,
		ReservedSlots:   0,
		PricePerPerson:  req.PricePerPerson,
		GroupSizeMax:    req.GroupSizeMax,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.activityRepo.Timeslot(ctx, timeSlot); err != nil {
		return nil, fmt.Errorf("creating timeslot: %w", err)
	}

	return &ActivityTimeslotResult{
		TimeSlotID: timeSlotID,
		ActivityID: activityID,
	}, nil
}

func (s *ActivityService) UpdateTimeSlot(ctx context.Context, timeSlotID primitive.ObjectID, req ActivityTimeslotRequest) error {
	if err := s.activityRepo.UpdateTimeslot(ctx, timeSlotID, req.StartTime, req.DurationMinutes, req.TotalSlots, req.GroupSizeMax, req.PricePerPerson); err != nil {
		return fmt.Errorf("updating timeslot: %w", err)
	}
	return nil
}

func (s *ActivityService) SetTimeSlotActive(ctx context.Context, timeSlotID primitive.ObjectID, active bool) error {
	if err := s.activityRepo.IsActive(ctx, timeSlotID, active); err != nil {
		return fmt.Errorf("setting timeslot active state: %w", err)
	}
	return nil
}

func (s *ActivityService) GetActivity(ctx context.Context, activityID primitive.ObjectID) (*model.Activity, error){
	activity, err := s.activityRepo.GetActivity(ctx, activityID)
	if err != nil{
		return nil, fmt.Errorf("getting activity: %w", err)
	}

	return activity, nil
}

func (s *ActivityService) GetActivities(ctx context.Context)(*[]model.Activity, error){
	activities, err := s.activityRepo.GetActivities(ctx)
	if err != nil{
		return nil, fmt.Errorf("getting activities: %w", err)
	}

	return &activities, nil 
}

func (s *ActivityService) GetTimeslot(ctx context.Context, timeslotID primitive.ObjectID) (*model.ActivityTimeslot, error){
	timeSlot, err := s.activityRepo.GetTmeslot(ctx, timeslotID)
	if err != nil{
		return nil, fmt.Errorf("getting activity timeSlot: %w", err)
	}

	return timeSlot, nil
}

func (s *ActivityService) GetTimeSlots(ctx context.Context)(*[]model.ActivityTimeslot, error){
	timeslots, err := s.activityRepo.GetTimeSlots(ctx)
	if err != nil{
		return nil, fmt.Errorf("getting activity timeslots: %w", err)
	}

	return &timeslots, nil 
}