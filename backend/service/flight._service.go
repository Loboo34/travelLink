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
	OriginID      string                   `json:"originID"`
	DestinationID string                   `json:"destinationID"`
	DepartureTime string                   `json:"departureTime"`
	ArrivalTime   string                   `json:"arrivalTime"`
	AirlineID     string                   `json:"airlineID"`
	FlightNumber  string                   `json:"flightNumber"`
	CabinClass    []model.FlightCabinClass `json:"cabinClass"`
	Segments      []string                 `json:"segments"`
	PlaneID       string                   `json:"planeID"`
}

func (r *FlightRequest) Validate() error {
	if r.OriginID == "" {
		return errors.New("origin ID required")
	}

	if r.DestinationID == "" {
		return fmt.Errorf("destination ID required")
	}

	if r.PlaneID == "" {
		return fmt.Errorf("plane ID required")
	}

	if r.AirlineID == "" {
		return fmt.Errorf("airline ID required")
	}

	if r.FlightNumber == "" {
		return fmt.Errorf("flight number is required")
	}

	if r.DepartureTime == "" {
		return fmt.Errorf("departure time is required")
	}

	if r.ArrivalTime == "" {
		return fmt.Errorf("arrival time is required")
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

	// Parse departure time
	deptTime, err := time.Parse(time.RFC3339, req.DepartureTime)
	if err != nil {
		return nil, &model.ValidationError{Message: "invalid departure time format"}
	}

	// Parse arrival time
	arrTime, err := time.Parse(time.RFC3339, req.ArrivalTime)
	if err != nil {
		return nil, &model.ValidationError{Message: "invalid arrival time format"}
	}

	// Convert string IDs to ObjectIDs
	originID, err := primitive.ObjectIDFromHex(req.OriginID)
	if err != nil {
		return nil, &model.ValidationError{Message: "invalid origin ID"}
	}

	destID, err := primitive.ObjectIDFromHex(req.DestinationID)
	if err != nil {
		return nil, &model.ValidationError{Message: "invalid destination ID"}
	}

	airlineID, err := primitive.ObjectIDFromHex(req.AirlineID)
	if err != nil {
		return nil, &model.ValidationError{Message: "invalid airline ID"}
	}

	planeID, err := primitive.ObjectIDFromHex(req.PlaneID)
	if err != nil {
		return nil, &model.ValidationError{Message: "invalid plane ID"}
	}

	// Convert segment strings to ObjectIDs
	segments := make([]primitive.ObjectID, 0)
	for _, segStr := range req.Segments {
		if segStr == "" {
			continue
		}
		segID, err := primitive.ObjectIDFromHex(segStr)
		if err != nil {
			return nil, &model.ValidationError{Message: "invalid segment ID"}
		}
		segments = append(segments, segID)
	}

	flightID := primitive.NewObjectID()
	now := time.Now()

	flight := &model.Flight{
		ID:            flightID,
		OriginID:      originID,
		DestinationID: destID,
		DepartureTime: deptTime,
		ArrivalTime:   arrTime,
		AirlineID:     airlineID,
		FlightNumber:  req.FlightNumber,
		CabinClass:    req.CabinClass,
		Segments:      segments,
		Stops:         len(segments) - 1,
		PlaneID:       planeID,
		Status:        model.FlightStatusScheduled,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.FlightRepo.Create(ctx, flight); err != nil {
		return nil, fmt.Errorf("creating flight: %w", err)
	}

	return flight, nil

}

type UpdateReq struct {
	DepartureTime time.Time `json:"departureTime"`
	ArrivalTime   time.Time `json:"arrivalTime"`
	Stops         int       `json:"stops"`
}

func (s *FlightService) Update(ctx context.Context, flightID primitive.ObjectID, req UpdateReq) (*UpdateReq, error) {
	_, err := s.FlightRepo.GetFlight(ctx, flightID)
	if err != nil {
		return nil, err
	}
	if err := s.FlightRepo.Update(ctx, flightID, req.ArrivalTime, req.DepartureTime, req.Stops); err != nil {
		return nil, fmt.Errorf("updating flight: %w", err)
	}

	return &UpdateReq{}, nil
}

func (s *FlightService) UpdateStatus(ctx context.Context, flightID primitive.ObjectID, status *model.FlightStatus) error {
	_, err := s.FlightRepo.GetFlight(ctx, flightID)
	if err != nil {
		return err
	}
	if err := s.FlightRepo.UpdateStatus(ctx, flightID, status); err != nil {
		return fmt.Errorf("updating status: %w", err)
	}

	return nil
}

func (s *FlightService) DeleteFlight(ctx context.Context, flightID primitive.ObjectID) error {
	_, err := s.FlightRepo.GetFlight(ctx, flightID)
	if err != nil {
		return err
	}
	if err := s.FlightRepo.Delete(ctx, flightID); err != nil {
		return fmt.Errorf("deleting flight: %w", err)
	}

	return nil
}

type Offer struct {
	FlightID          primitive.ObjectID       `json:"flightID"`
	ProviderReference string                   `json:"providerReference"`
	Provider          string                   `json:"provider"`
	OneWay            bool                     `json:"oneway"`
	Segments          []primitive.ObjectID     `json:"segments"`
	PriceTotal        int64                    `json:"priceTotal"`
	BaggageAllowance  []model.BaggageAllowance `json:"baggageAllowance"`
	LastTicketingDate *time.Time               `json:"lastTicketingDate"`
	BookableSeats     int                      `json:"bookableSeats"`
	ExpiresAt         *time.Time               `json:"expiresAt"`
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
	_, err := s.FlightRepo.GetOffer(ctx, offerID)
	if err != nil {
		return nil, err
	}
	if err := s.FlightRepo.UpdateOffer(ctx, offerID, req.Price, req.BookableSeats, req.OneWay); err != nil {
		return nil, fmt.Errorf("updating offer: %w", err)
	}

	return &model.FlightOffer{}, nil
}

func (s *FlightService) IsActive(ctx context.Context, offerID primitive.ObjectID, isActive bool) error {
	_, err := s.FlightRepo.GetOffer(ctx, offerID)
	if err != nil {
		return err
	}

	if err := s.FlightRepo.IsActive(ctx, offerID, isActive); err != nil {
		return fmt.Errorf("updating status: %w", err)
	}

	return nil
}

func (s *FlightService) DeleteOffer(ctx context.Context, offerID primitive.ObjectID) error {
	_, err := s.FlightRepo.GetOffer(ctx, offerID)
	if err != nil {
		return err
	}
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

func (s *FlightService) GetFlight(ctx context.Context, flightID primitive.ObjectID) (*model.Flight, error) {
	flight, err := s.FlightRepo.GetFlight(ctx, flightID)
	if err != nil {
		return nil, fmt.Errorf("getting flight: %w", err)
	}

	return flight, nil
}

func (s *FlightService) GetOffer(ctx context.Context, offerID primitive.ObjectID) (*model.FlightOffer, error) {
	offer, err := s.FlightRepo.GetOffer(ctx, offerID)
	if err != nil {
		return nil, fmt.Errorf("getting offer: %w", err)
	}

	return offer, nil
}

func (s *FlightService) GetOffers(ctx context.Context) (*[]model.FlightOffer, error) {
	offers, err := s.FlightRepo.GetOffers(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting offers: %w", err)
	}

	return &offers, nil
}
