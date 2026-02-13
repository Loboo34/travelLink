package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Flight struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DepartureAirport string             `bson:"departureAirport" json:"departureAirport"`
	ArrivalAirport   string             `bson:"arrivalAirport" json:"arrivalAirport"`
	DepartureTime    time.Time          `bson:"departureTime" json:"departureTime"`
	ArrivalTime      time.Time          `bson:"arrivalTime" json:"arrivalTime"`
	Airline          string             `bson:"airline" json:"airline"`
	FlightNumber     string             `bson:"flightNumber" json:"flightNumber"`
	CabinClass       string             `bson:"cabinClass" json:"cabinClass"`
	Stops            string             `bson:"stops" json:"stops"`
	PlaneType        string             `bson:"planeType" json:"planeType"`
}

type FlightOffer struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProviderReference string             `bson:"providerReference" json:"providerReference"` 
	Provider          string             `bson:"provider" json:"provider"`                     
	OneWay            bool               `bson:"one_way" json:"one_way"`
	Segments          []Flight           `bson:"flight" json:"flight"`
	PriceTotal        float64            `bson:"priceTotal" json:"priceTotal"`
	Currency          string             `bson:"currency" json:"currency"` 
	BaggageAllowance  string             `bson:"baggageAllowance,omitempty" json:"baggageAllowance,omitempty"`
	LastTicketingDate *time.Time         `bson:"lastTicketingDate,omitempty" json:"lastTicketingDate,omitempty"`
	BookableSeats     int                `bson:"bookable_seats" json:"bookable_seats"`
	CachedAt          time.Time          `bson:"cached_at" json:"cachedAt"` 
	ExpiresAt         *time.Time         `bson:"expiresAt,omitempty" json:"expiresAt,omitempty"`
}
