package service

import (
	"context"
	"fmt"
	"time"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccommodationService struct {
	accommodationRepo *repository.AccommodationRepo
}

func NewAccommodationService(accommodationRepo *repository.AccommodationRepo) *AccommodationService {
	return &AccommodationService{accommodationRepo: accommodationRepo}
}

type AccommodationRequest struct {
	HostID       *primitive.ObjectID `json:"hostID"`
	PropertyType model.PropertyType  `json:"propertyType"`
	Name         string              `json:"name"`
	Address      model.Address       `json:"address"`
	Amenities    []model.Amenity     `json:"amenities"`
	Description  string              `json:"description"`
	Images       []string            `json:"images"`
	Location     model.GeoLocation   `json:"location"`
	RoomType     []model.RoomType    `json:"roomType"`
}

type AccommodationResult struct {
	AccommodationID primitive.ObjectID `json:"accommodationID"`
	PropertyType    model.PropertyType `json:"propertyType"`
	Name            string             `json:"name"`
	Adress          model.Address      `json:"address"`
}

type AccommodationUpdate struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Amenities   []string `json:"amenities"`
	Images      []string `json:"images"`
}

type AvailabilityRequest struct {
	AccommodationID primitive.ObjectID `json:"accommodationID"`
	RoomTypeID      primitive.ObjectID `json:"roomTypeID"`
	Date            time.Time          `json:"date"`
	TotalRooms      int                `json:"totalRooms"`
	PricePerNight   int64              `json:"pricePerNight"`
}

type AvailabilityResult struct {
	AccommodationID primitive.ObjectID `json:"accommodationID"`
	RoomTypeID      primitive.ObjectID     `json:"roomTypeID"`
	TotalRooms      int                `json:"totalRooms"`
}

func (s *AccommodationService) Add(ctx context.Context, req AccommodationRequest) (*AccommodationResult, error) {

	accommodationID := primitive.NewObjectID()

	accommodation := &model.Accommodation{
		ID:           accommodationID,
		HostID:       req.HostID,
		PropertyType: req.PropertyType,
		Name:         req.Name,
		Address:      req.Address,
		Description:  req.Description,
		Amenities:    req.Amenities,
		Images:       req.Images,
		Location:     req.Location,
		RoomType:     req.RoomType,
		Rating:       0,
		ReviewCount:  0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.accommodationRepo.Add(ctx, accommodation); err != nil {
		return nil, fmt.Errorf("acceing accommodation: %w", err)
	}

	return &AccommodationResult{
		AccommodationID: accommodationID,
		PropertyType:    req.PropertyType,
		Name:            req.Name,
		Adress:          req.Address,
	}, nil

}

func (s *AccommodationService) Update(ctx context.Context, accommodationId primitive.ObjectID, req AccommodationUpdate) (*AccommodationResult, error) {
	if err := s.accommodationRepo.Update(ctx, accommodationId, req.Name, req.Description, req.Amenities, req.Images); err != nil {
		return nil, fmt.Errorf("updating accommodtion: %w ", err)
	}

	return &AccommodationResult{

		AccommodationID: accommodationId,
		Name:            req.Name,
	}, nil
}

func (s *AccommodationService) Delete(ctx context.Context, accommodationID primitive.ObjectID) error {
	if err := s.accommodationRepo.Delete(ctx, accommodationID); err != nil {
		return fmt.Errorf("deleting accommodation: %w", err)
	}

	return nil
}

func (s *AccommodationService) Availability(ctx context.Context, req AvailabilityRequest) (*AvailabilityResult, error) {

	availabilityID := primitive.NewObjectID()
	availability := &model.AccommodationAvailability{
		ID:              availabilityID,
		AccommodationID: req.AccommodationID,
		RoomTypeID:      req.RoomTypeID,
		Date:            req.Date,
		TotalRooms:      req.TotalRooms,
		ReservedRooms:   0,
		PricePerNight:   req.PricePerNight,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.accommodationRepo.Availability(ctx, availability); err != nil {
		return nil, fmt.Errorf("creating availability: %w", err)
	}

	return &AvailabilityResult{
		AccommodationID: req.AccommodationID,
		RoomTypeID:      req.RoomTypeID,
		TotalRooms:      req.TotalRooms,
	}, nil
}

func (s *AccommodationService) IsActive(ctx context.Context, availabilityID primitive.ObjectID, status bool) error {
	if err := s.accommodationRepo.IsActive(ctx, availabilityID, status); err != nil {
		return fmt.Errorf("updating status: %w", err)
	}

	return nil

}

func (s *AccommodationService) Remove(ctx context.Context, accommodationID primitive.ObjectID) error {
	if err := s.accommodationRepo.Remove(ctx, accommodationID); err != nil {
		return fmt.Errorf("deleting accommodation: %w", err)
	}

	return nil
}


func (s *AccommodationService) GetAccomodations(ctx context.Context) (*[]model.Accommodation, error) {
	accommodations, err := s.accommodationRepo.GetAccomodations(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting accommodations: %w", err)
	}

	return &accommodations, nil
}

func (s *AccommodationService) GetAccommodation(ctx context.Context, accommodationID primitive.ObjectID) (*model.Accommodation, error){
 accommodation, err := s.accommodationRepo.GetAccomodation(ctx, accommodationID)
 if err != nil{
	return nil, fmt.Errorf("getting accommodation: %w", err)
 }

 return accommodation, nil
}


func (s *AccommodationService) GetAvailabilities(ctx context.Context) (*[]model.AccommodationAvailability, error){
	availability, err := s.accommodationRepo.GetAvailabilities(ctx)
	if err != nil{
		return nil, fmt.Errorf("getting available accommodations: %w", err)
	}

	return &availability, nil
}

func (s *AccommodationService) GetAvailability(ctx context.Context, availabilityID primitive.ObjectID) (*model.AccommodationAvailability, error){
	accom, err := s.accommodationRepo.GetAvailability(ctx, availabilityID)
	if err != nil{
		return nil, fmt.Errorf("getting available accommodation: %w", err)
	}

	return accom, nil
}