package model

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlightRequest struct {
	FlightID         primitive.ObjectID `json:"flightID"`
	PassengerDetails []UserDetails      `json:"passengerDetails"`
	SelectedSeats    []string           `json:"selectedSeates"`
	PaymentMethod    PaymentMethod      `json:"paymentMethod"`
	Currency         string             `json:"currency"`
}

func (r *FlightRequest) Validate() error {
	// offer
	if r.FlightID.IsZero() {
		return errors.New("flightOfferID is required")
	}

	// passengers
	if len(r.PassengerDetails) < 1 {
		return errors.New("at least one passenger is required")
	}
	for i, p := range r.PassengerDetails {
		if p.FirstName == "" {
			return fmt.Errorf("passenger %d: firstName is required", i+1)
		}
		if p.LastName == "" {
			return fmt.Errorf("passenger %d: lastName is required", i+1)
		}
		if p.Passport == "" {
			return fmt.Errorf("passenger %d: passportNo is required", i+1)
		}
		if p.Nationality == "" {
			return fmt.Errorf("passenger %d: nationality is required", i+1)
		}
		if p.DateOfBirth.IsZero() {
			return fmt.Errorf("passenger %d: dateOfBirth is required", i+1)
		}
	}

	// payment
	switch r.PaymentMethod {
	case PaymentMethodCard, PaymentMethodMpesa, PaymentMethodBank:
	default:
		return errors.New("invalid payment method")
	}

	if r.Currency == "" {
		r.Currency = "USD" // sensible default
	}

	return nil
}

type AccommodationBookingRequest struct {
	AccommodationID primitive.ObjectID `json:"accommodationID"`
	GuestDetails    []UserDetails      `json:"guestDetails"`
	RoomTypeID      primitive.ObjectID `json:"roomTypeID"`
	Rooms           int                `json:"rooms"`
	CheckIn         time.Time          `json:"checkIn"`
	CheckOut        time.Time          `json:"checkOut"`
	PaymentMethod   PaymentMethod      `json:"paymentMethod"`
	Currency        string             `json:"currency"`
}

func (r *AccommodationBookingRequest) Validate() error {
	if r.AccommodationID.IsZero() {
		return errors.New("Accommodation ID is Required")
	}

	if r.RoomTypeID.IsZero() {
		return errors.New("Room type ID is required")
	}

	if len(r.GuestDetails) < 1 {
		return errors.New("A guest is required for booking")
	}

	for i, G := range r.GuestDetails {
		if G.FirstName == "" {
			return fmt.Errorf("passenger %d: firstName is required", i+1)
		}
		if G.LastName == "" {
			return fmt.Errorf("passenger %d: lastName is required", i+1)
		}
	}

	if r.Rooms < 1 {
		r.Rooms = 1
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	if r.CheckIn.UTC().Before(today) {
		return errors.New("checkIn date cannot be in the past")
	}
	if !r.CheckOut.After(r.CheckIn) {
		return errors.New("checkOut date must be after check-in date")
	}

	// payment
	switch r.PaymentMethod {
	case PaymentMethodCard, PaymentMethodMpesa, PaymentMethodBank:
	default:
		return errors.New("invalid payment method")
	}

	if r.Currency == "" {
		r.Currency = "USD" // sensible default
	}

	return nil

}

type ActivityBookingRequest struct {
	ActivityID         primitive.ObjectID `json:"activityID"`
	ParticipantDetails []UserDetails      `json:"participantDetails"`
	TimeSlotID         primitive.ObjectID `json:"timeSlotID"`
	Participants       int                `json:"participants"`
	PaymentMethod      PaymentMethod      `json:"paymentMethod"`
	Currency           string             `json:"currency"`
}

func (r *ActivityBookingRequest) Validate() error {
	if r.ActivityID.IsZero() {
		return errors.New("activityID is required")
	}

	if r.TimeSlotID.IsZero() {
		return errors.New("timeslotID is required")
	}

	if len(r.ParticipantDetails) < 1 {
		return errors.New("at least one passenger is required")
	}
	for i, p := range r.ParticipantDetails {
		if p.FirstName == "" {
			return fmt.Errorf("passenger %d: firstName is required", i+1)
		}
		if p.LastName == "" {
			return fmt.Errorf("passenger %d: lastName is required", i+1)
		}
	}

	if r.Participants < 0 {
		r.Participants = 1
	}

	switch r.PaymentMethod {
	case PaymentMethodCard, PaymentMethodMpesa, PaymentMethodBank:
	default:
		return errors.New("invalid payment method")
	}

	if r.Currency == "" {
		r.Currency = "USD" // sensible default
	}

	return nil
}

type PackageBookingRequest struct {
	PackageID          primitive.ObjectID   `json:"packageID"`
	TravelersDetails   []UserDetails        `json:"travelersDetails"`
	ComponentSelection []ComponentSelection `json:"componentSelection"`
	PaymentMethod      PaymentMethod        `json:"paymentMethod"`
	Currency           string               `json:"currency"`
}

type ComponentSelection struct {
	ComponentType Component           `json:"componentType"`
	ReferenceID   primitive.ObjectID  `json:"referenceID"`
	RoomTypeID    *primitive.ObjectID `json:"roomTypeID"`
	CheckIn       *time.Time          `json:"checkin"`
	CheckOut      *time.Time          `json:"checkout"`
	Rooms         int                 `json:"rooms"`
	Participants  int                 `json:"participants"`
	TimeslotID    *primitive.ObjectID `json:"timeslotID"`
}

func (r *PackageBookingRequest) Validate() error {
	if r.PackageID.IsZero() {
		return errors.New("activityID is required")
	}

	if len(r.TravelersDetails) < 1 {
		return errors.New("at least one passenger is required")
	}
	for i, p := range r.TravelersDetails {
		if p.FirstName == "" {
			return fmt.Errorf("passenger %d: firstName is required", i+1)
		}
		if p.LastName == "" {
			return fmt.Errorf("passenger %d: lastName is required", i+1)
		}
	}

	switch r.PaymentMethod {
	case PaymentMethodCard, PaymentMethodMpesa, PaymentMethodBank:
	default:
		return errors.New("invalid payment method")
	}

	if r.Currency == "" {
		r.Currency = "USD" // sensible default
	}

	return nil
}
