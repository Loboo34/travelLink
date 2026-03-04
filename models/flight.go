package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Flight struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OriginID      primitive.ObjectID `bson:"originID" json:"originID"`
	DestinationID primitive.ObjectID `bson:"destinationID" json:"destinationID"`
	DepartureTime time.Time          `bson:"departureTime" json:"departureTime"`
	ArrivalTime   time.Time          `bson:"arrivalTime" json:"arrivalTime"`
	AirlineID     primitive.ObjectID `bson:"airlineID" json:"airlineID"`
	FlightNumber  string             `bson:"flightNumber" json:"flightNumber"`
	CabinClass    []FlightCabinClass `bson:"cabinClass" json:"cabinClass"`
	Stops         int                `bson:"stops" json:"stops"`
	PlaneID       primitive.ObjectID `bson:"planeID" json:"planeID"`
	Status        string             `bson:"status" json:"status"`
}

type FlightOffer struct {
	ID                primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	FlightID          primitive.ObjectID   `bson:"flightID" json:"flightID"`
	ProviderReference string               `bson:"providerReference" json:"providerReference"`
	Provider          string               `bson:"provider" json:"provider"`
	OneWay            bool                 `bson:"oneway" json:"oneway"`
	Segments          []primitive.ObjectID `bson:"flight" json:"flight"`
	PriceTotal        int64                `bson:"priceTotal" json:"priceTotal"`
	BaggageAllowance  string               `bson:"baggageAllowance,omitempty" json:"baggageAllowance,omitempty"`
	LastTicketingDate *time.Time           `bson:"lastTicketingDate,omitempty" json:"lastTicketingDate,omitempty"`
	BookableSeats     int                  `bson:"bookable_seats" json:"bookable_seats"`
	CachedAt          time.Time            `bson:"cached_at" json:"cachedAt"`
	ExpiresAt         *time.Time           `bson:"expiresAt,omitempty" json:"expiresAt,omitempty"`
	IsActive          bool                 `bson:"isActive" json:"isActive"`
}

type PlaneModel struct {
	ID           primitive.ObjectID
	Name         string
	Model        string
	CabinClasses []CabinClass
}

type CabinClass struct {
	Name         string
	SeatCapacity int
	AmenityLevel string
}

type FlightCabinClass struct {
	ID            primitive.ObjectID
	ClassType     string
	TotalSeats    int
	ReservedSeats int
	Price         int64
}

type Airline struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Name string             `bson:"name" json:"name"`
	Code string             `bson:"code" json:"code"`
}

type Airport struct {
	ID        primitive.ObjectID
	Code      string // IATA  NBO
	Name      string
	City      string
	Country   string
	Latitude  float64
	Longitude float64
	Timezone  string
}

type Route struct {
	ID                   primitive.ObjectID `bson:"_id" json:"id"`
	OriginAirportID      primitive.ObjectID `bson:"originAirportID" json:"originAirportID"`
	DestinationAirportID primitive.ObjectID `bson:"destinationAirportID" json:"destinationAirportID"`
	EstimatedDurationMin int                `bson:"estimatedDurationMin" json:"estimatedDurationMin"`
}
