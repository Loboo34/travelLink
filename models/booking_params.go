package model

import (

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlightRequest struct {
	FlightID         primitive.ObjectID `json:"flightID"`
	PassengerDetails []UserDetails `json:"passengerDetails"`
	SelectedSeats    []string           `json:"selectedSeates"`
	PaymentMethod    PaymentMethod      `json:"paymentMethod"`
}


