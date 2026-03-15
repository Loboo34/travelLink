package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PropertyType string

const (
	PropertyTypeHotel      PropertyType = "Hotel"
	PropertyTypeAirBnb     PropertyType = "AirBNB"
	PropertyTypeVilla      PropertyType = "Villa"
	PropertyTypeResort     PropertyType = "Resort"
	PropertyTypeGuesthouse PropertyType = "Guesthouse"
)

type Amenity string

const (
	AmenityWifi        Amenity = "wifi"
	AmenityParking     Amenity = "parking"
	AmenityPool        Amenity = "pool"
	AmenityGym         Amenity = "gym"
	AmenityAirCon      Amenity = "airConditioning"
	AmenityBreakfast   Amenity = "breakfast"
	AmenityPetFriendly Amenity = "petFriendly"
	AmenityKitchen     Amenity = "kitchen"
)

type Address struct {
	Street  string `bson:"street" json:"street"`
	City    string `bson:"city" json:"city"`
	Country string `bson:"country" json:"country"`
	ZipCode string `bson:"zipCode" json:"zipCode"`
}

type GeoLocation struct {
	Type        string    `bson:"type" json:"type"`               // always "Point"
	Coordinates []float64 `bson:"coordinates" json:"coordinates"` // [longitude, latitude]
}

type RoomType struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Description   string             `bson:"description" json:"description"`
	Bedrooms      int                `bson:"bedrooms" json:"bedrooms"`
	Bathrooms     int                `bson:"bathrooms" json:"bathrooms"`
	MaxGuests     int                `bson:"maxGuests" json:"maxGuests"`
	IsEntirePlace bool               `bson:"isEntirePlace" json:"isEntirePlace"`
	Amenities     []Amenity          `bson:"amenities" json:"amenities"`
	Images        []string           `bson:"images" json:"images"`
	BasePrice     int64              `bson:"basePrice" json:"basePrice"`
}
type Accommodation struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	HostID       *primitive.ObjectID `bson:"hostID,omitempty" json:"hostID,omitempty"` // nil for hotels
	PropertyType PropertyType        `bson:"type" json:"type"`
	Name         string              `bson:"name" json:"name"`
	Address      Address             `bson:"address" json:"address"`
	Description  string              `bson:"description" json:"description"`
	Amenities    []Amenity           `bson:"amenities" json:"amenities"`
	Images       []string            `bson:"images" json:"images"`
	Location     GeoLocation         `bson:"location" json:"location"`
	RoomType     []RoomType          `bson:"roomType" json:"roomType"`
	Rating       float64             `bson:"rating" json:"rating"`
	ReviewCount  int                 `bson:"reviews" json:"reviews"`
	CachedAt     time.Time           `bson:"cachedAt" json:"cachedAt"`
	CheckInTime  int                 `bson:"checkInTime,omitempty" json:"checkInTime,omitempty"` //how to deal with checkin/check out-
	CheckOutTime int                 `bson:"checkOutTime,omitempty" json:"checkOutTime,omitempty"`
	IsActive     bool                `bson:"isActive" json:"isActive"`
	CreatedAt    time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time           `bson:"updatedAt" json:"updatedAt"`
}

type AccommodationAvailability struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AccommodationID primitive.ObjectID `bson:"accommodationID" json:"accommodationID"`
	RoomTypeID      primitive.ObjectID `bson:"roomTypeID" json:"roomTypeID"`
	Date            time.Time          `bson:"date" json:"date"`
	TotalRooms      int                `bson:"totalRooms" json:"totalRooms"`
	ReservedRooms   int                `bson:"reservedRooms" json:"reservedRooms"`
	PricePerNight   int64              `bson:"pricePerNight" json:"pricePerNight"`
	Currency        string             `bson:"currency" json:"currency"`
	IsActive        bool               `bson:"isActive" json:"isActive"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type AccommodationSearchResult struct {
	AccommodationID   primitive.ObjectID `bson:"_id" json:"accommodationID"`
	Accommodation     Accommodation      `bson:"accommodationDoc" json:"accommodation"`
	PricePerNight     int64              `bson:"pricePerNight" json:"pricePerNight"`
	MinAvailableRooms int                `bson:"minAvailableRooms" json:"minAvailableRooms"`
	AvailableNights   int                `bson:"availableNights" json:"availableNights"`
	Currency          string             `bson:"currency" json:"currency"`
}
