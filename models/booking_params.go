package model

import (
	"errors"
	"fmt"

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
