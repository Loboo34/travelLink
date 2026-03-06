package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlightStatus string

const (
	FlightStatusScheduled FlightStatus = "scheduled"
	FlightStatusDelayed   FlightStatus = "delayed"
	FlightStatusCancelled FlightStatus = "cancelled"
	FlightStatusDeparted  FlightStatus = "departed"
	FlightStatusArrived   FlightStatus = "arrived"
)

type Flight struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	OriginID      primitive.ObjectID   `bson:"originID" json:"originID"`
	DestinationID primitive.ObjectID   `bson:"destinationID" json:"destinationID"`
	DepartureTime time.Time            `bson:"departureTime" json:"departureTime"`
	ArrivalTime   time.Time            `bson:"arrivalTime" json:"arrivalTime"`
	AirlineID     primitive.ObjectID   `bson:"airlineID" json:"airlineID"`
	FlightNumber  string               `bson:"flightNumber" json:"flightNumber"`
	CabinClass    []FlightCabinClass   `bson:"cabinClass" json:"cabinClass"`
	Stops         int                  `bson:"stops" json:"stops"` //check how stops are handled
	Segments      []primitive.ObjectID `bson:"segments" json:"segments"`
	PlaneID       primitive.ObjectID   `bson:"planeID" json:"planeID"`
	Status        FlightStatus         `bson:"status" json:"status"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt" json:"UpdatedAt"`
}

type FlightOffer struct {
	ID                primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	FlightID          primitive.ObjectID   `bson:"flightID" json:"flightID"`
	ProviderReference string               `bson:"providerReference" json:"providerReference"`
	Provider          string               `bson:"provider" json:"provider"`
	OneWay            bool                 `bson:"oneway" json:"oneway"`
	Segments          []primitive.ObjectID `bson:"segments" json:"segments"`
	PriceTotal        int64                `bson:"priceTotal" json:"priceTotal"`
	BaggageAllowance  BaggageAllowance     `bson:"baggageAllowance,omitempty" json:"baggageAllowance,omitempty"`
	LastTicketingDate *time.Time           `bson:"lastTicketingDate,omitempty" json:"lastTicketingDate,omitempty"`
	BookableSeats     int                  `bson:"bookable_seats" json:"bookable_seats"`
	CachedAt          time.Time            `bson:"cached_at" json:"cachedAt"`
	ExpiresAt         *time.Time           `bson:"expiresAt,omitempty" json:"expiresAt,omitempty"`
	IsActive          bool                 `bson:"isActive" json:"isActive"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt" json:"UpdatedAt"`
}

type PlaneModel struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Model        string             `bson:"model" json:"model"`
	CabinClasses []CabinClass       `bson:"cabinClasses" json:"cabinClasses"`
}

type CabinClassType string

const (
	CabinClassEconomy  CabinClassType = "economy"
	CabinClassBusiness CabinClassType = "business"
	CabinClassFirst    CabinClassType = "first"
)

type CabinClass struct {
	Name         string         `bson:"name" json:"name"`
	Type         CabinClassType `bson:"type" json:"type"`
	SeatCapacity int            `bson:"seatCapacity" json:"seatCapacity"`
	AmenityLevel string         `bson:"amenityLevel" json:"amenityLevel"`
}

type FlightCabinClass struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClassType     CabinClassType     `bson:"classType" json:"classType"`
	TotalSeats    int                `bson:"totalSeats" json:"totalSeats"`
	ReservedSeats int                `bson:"reservedSeats" json:"reservedSeats"`
	Price         int64              `bson:"price" json:"price"` // in smallest currency unit
}

type BaggageAllowance struct {
	Pieces   int `bson:"pieces" json:"pieces"`
	WeightKg int `bson:"weightKg" json:"weightKg"`
}

type FlightSegment struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FlightID             primitive.ObjectID `bson:"flightID" json:"flightID"`
	OriginAirportID      primitive.ObjectID `bson:"originAiportID" json:"originAirportID"`
	DestinationAirportID primitive.ObjectID `bson:"destinationAirportID" json:"destinationAirportID"`
	ArrivalTime          time.Time          `bson:"arrivalTime" json:"arrivalTime"`
	DepartureTime        time.Time          `bson:"departureTime" json:"departureTime"`
	AirlineID            primitive.ObjectID `bson:"airlineID" json:"airlineID"`
	FlightNumber         string             `bson:"flightNumber" json:"flightNumber"`
}

type Airline struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name" json:"name"`
	Code string             `bson:"code" json:"code"`
}

type Airport struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code      string             `bson:"code" json:"code"`
	Name      string             `bson:"name" json:"name"`
	City      string             `bson:"city" json:"city"`
	Country   string             `bson:"country" json:"country"`
	Latitude  float64            `bson:"latitude" json:"latitude"`
	Longitude float64            `bson:"longitude" json:"longitude"`
	Timezone  string             `bson:"timezone" json:"timezone"`
}

type Route struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OriginAirportID      primitive.ObjectID `bson:"originAirportID" json:"originAirportID"`
	DestinationAirportID primitive.ObjectID `bson:"destinationAirportID" json:"destinationAirportID"`
	EstimatedDurationMin int                `bson:"estimatedDurationMin" json:"estimatedDurationMin"`
}
