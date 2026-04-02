package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlightService struct {
	FlightRepo *repository.FlightRepo
}

func NewFlightService(flightRepo *repository.FlightRepo) *FlightService {
	return &FlightService{FlightRepo: flightRepo}
}

type FlightRequest struct {
	OriginID      primitive.ObjectID       `bson:"originID" json:"originID"`
	DestinationID primitive.ObjectID       `bson:"destinationID" json:"destinationID"`
	DepartureTime time.Time                `bson:"departureTime" json:"departureTime"`
	ArrivalTime   time.Time                `bson:"arrivalTime" json:"arrivalTime"`
	AirlineID     primitive.ObjectID       `bson:"airlineID" json:"airlineID"`
	FlightNumber  string                   `bson:"flightNumber" json:"flightNumber"`
	CabinClass    []model.FlightCabinClass `bson:"cabinClass" json:"cabinClass"`
	Segments      []primitive.ObjectID     `bson:"segments" json:"segments"`
	PlaneID       primitive.ObjectID       `bson:"palneID" json:"planeID"`
	//Status        string                   `json:"status"`
}

func (r *FlightRequest) Validate() error {
	if r.OriginID.IsZero() {
		return errors.New("origin ID required")
	}

	if r.DestinationID.IsZero() {
		return fmt.Errorf("destination ID required")
	}

	if r.PlaneID.IsZero() {
		return fmt.Errorf("plane ID required")
	}

	if len(r.Segments) == 0 {
		return fmt.Errorf("segment ID required")
	}

	if r.AirlineID.IsZero() {
		return fmt.Errorf("destination ID required")
	}

	if r.FlightNumber == "" {
		return fmt.Errorf("number is required")
	}

	return nil
}

type FlightResult struct {
	FlightID      primitive.ObjectID `json:"flightID"`
	OriginID      primitive.ObjectID `json:"originID"`
	DestinationID primitive.ObjectID `json:"destinationID"`
	DepartureTime time.Time          `json:"departureTime"`
	ArrivalTime   time.Time          `json:"arrivalTime"`
	AirlineID     primitive.ObjectID `json:"airline"`
	FlightNumber  string             `json:"flightNumber"`
}

func (s *FlightService) AddFlight(ctx context.Context, req FlightRequest) (*model.Flight, error) {
	if err := req.Validate(); err != nil {
		return nil, &model.ValidationError{Message: err.Error()}
	}

	flightID := primitive.NewObjectID()
	now := time.Now()

	flight := &model.Flight{
		ID:            flightID,
		OriginID:      req.OriginID,
		DestinationID: req.DestinationID,
		DepartureTime: req.DepartureTime,
		ArrivalTime:   req.ArrivalTime,
		AirlineID:     req.AirlineID,
		FlightNumber:  req.FlightNumber,
		CabinClass:    req.CabinClass,
		Segments:      req.Segments,
		Stops:         len(req.Segments) - 1,
		PlaneID:       req.PlaneID,
		Status:        model.FlightStatusScheduled,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.FlightRepo.Create(ctx, flight); err != nil {
		return nil, fmt.Errorf("creating flight: %w", err)
	}

	return &model.Flight{}, nil

}

type UpdateReq struct {
	DepartureTime time.Time `json:"departureTime"`
	ArrivalTime   time.Time `json:"arrivalTime"`
	Stops         int       `json:"stops"`
}

func (s *FlightService) Update(ctx context.Context, flightID primitive.ObjectID, req UpdateReq) (*UpdateReq, error) {
	err := s.FlightRepo.Update(ctx, flightID, req.ArrivalTime, req.DepartureTime, req.Stops)
	if err != nil {
		return nil, fmt.Errorf("updating flight: %w", err)
	}

	return &UpdateReq{}, nil
}

func (s *FlightService) UpdateStatus(ctx context.Context, flightId primitive.ObjectID, status *model.FlightStatus) error {
	if err := s.FlightRepo.UpdateStatus(ctx, flightId, status); err != nil {
		return fmt.Errorf("updating status: %w", err)
	}

	return nil
}

func (s *FlightService) DeleteFlight(ctx context.Context, flighID primitive.ObjectID) error {
	if err := s.FlightRepo.Delete(ctx, flighID); err != nil {
		return fmt.Errorf("deleting flight: %w", err)
	}

	return nil
}

type Offer struct {
	FlightID          primitive.ObjectID     `json:"flightID"`
	ProviderReference string                 `json:"providerReference"`
	Provider          string                 `json:"provider"`
	OneWay            bool                   `json:"oneway"`
	Segments          []primitive.ObjectID   `json:"segments"`
	PriceTotal        int64                  `json:"priceTotal"`
	BaggageAllowance  []model.BaggageAllowance `json:"baggageAllowance"`
	LastTicketingDate *time.Time             `json:"lastTicketingDate"`
	BookableSeats     int                    `json:"bookableSeats"`
	ExpiresAt         *time.Time             `json:"expiresAt"`
}

func (s *FlightService) CreateOffer(ctx context.Context, req Offer) (*model.FlightOffer, error) {

	offerID := primitive.NewObjectID()
	now := time.Now()

	Offer := &model.FlightOffer{
		ID:                offerID,
		FlightID:          req.FlightID,
		ProviderReference: req.ProviderReference,
		Provider:          req.Provider,
		OneWay:            req.OneWay,
		Segments:          req.Segments,
		PriceTotal:        req.PriceTotal,
		BaggageAllowance:  req.BaggageAllowance,
		LastTicketingDate: req.LastTicketingDate,
		BookableSeats:     req.BookableSeats,
		CachedAt:          now,
		ExpiresAt:         req.ExpiresAt,
		IsActive:          true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	err := s.FlightRepo.CreateOffer(ctx, Offer)
	if err != nil {
		return nil, fmt.Errorf("creating offer: %w", err)
	}

	return &model.FlightOffer{}, nil
}

type OfferUpdate struct {
	Price         int64 `json:"priceTotal"`
	OneWay        bool  `json:"oneway"`
	BookableSeats int   `json:"bookableSeats"`
}

func (s *FlightService) UpdateOffer(ctx context.Context, offerID primitive.ObjectID, req OfferUpdate) (*model.FlightOffer, error) {
	if err := s.FlightRepo.UpdateOffer(ctx, offerID,  req.Price, req.BookableSeats, req.OneWay); err != nil {
		return nil, fmt.Errorf("updating offer: %w", err)
	}

	return &model.FlightOffer{}, nil
}

func (s *FlightService) IsActive(ctx context.Context, offerID primitive.ObjectID, isActive bool) error {
	err := s.FlightRepo.IsActive(ctx, offerID, isActive)
	if err != nil {
		return fmt.Errorf("updating status: %w", err)
	}

	return nil
}

func (s *FlightService) DeleteOffer(ctx context.Context, offerID primitive.ObjectID) error {
	if err := s.FlightRepo.DeleteOffer(ctx, offerID); err != nil {
		return fmt.Errorf("deleting offer :%w", err)
	}

	return nil

}

func (s *FlightService) GetFlights(ctx context.Context) (*[]model.Flight, error) {
	flights, err := s.FlightRepo.GetFlights(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting flights: %w", err)
	}

	return &flights, nil
}

func (s *FlightService) GetFlight(ctx context.Context, flightID primitive.ObjectID) (*model.Flight, error){
 flight, err := s.FlightRepo.GetFlight(ctx, flightID)
 if err != nil{
	return nil, fmt.Errorf("getting flight: %w", err)
 }

 return flight, nil
}

func (s *FlightService) GetOffer(ctx context.Context, offerID primitive.ObjectID)(*model.FlightOffer, error){
	offer, err := s.FlightRepo.GetOffer(ctx, offerID)
	if err != nil{
		return nil, fmt.Errorf("getting offer: %w", err)
	}

	return offer, nil 
}

func (s *FlightService) GetOffers(ctx context.Context)(*[]model.FlightOffer, error){
	offers, err := s.FlightRepo.GetOffers(ctx)
	if err != nil{
		return nil, fmt.Errorf("getting offers: %w", err)
	}

	return &offers, nil
}

