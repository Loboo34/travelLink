package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Accommodation struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PropertyType string             `bson:"type" json:"type"`
	Name         string             `bson:"name" json:"name"`
	Address      string             `bson:"address" json:"address"`
	Amenities    []string           `bson:"amenities" json:"amenities"`
	Description  string             `bson:"description" json:"description"`
	Images       []string           `bson:"images" json:"images"`
	Location     string             `bson:"location" json:"location"`
	Fee          float64            `bson:"fee" json:"fee"`
	Rating       float64            `bson:"rating" json:"rating"`
	Reviews      []string           `bson:"reviews" json:"reviews"`
	CachedAt     time.Time          `bson:"cached_at" json:"cachedAt"`

	MaxGuests     int                `bson:"maxGuests,omitempty" json:"maxGuests,omitempty"`
	Bedrooms      int                `bson:"bedrooms,omitempty" json:"bedrooms,omitempty"`
	Bathrooms     int                `bson:"bathrooms,omitempty" json:"bathrooms,omitempty"`
	IsEntirePlace bool               `bson:"isEntirePlace" json:"isEntirePlace"`       // true for cottage/BnB, false for hotel rooms
	HostID        primitive.ObjectID `bson:"hostID,omitempty" json:"hostID,omitempty"` // for BnB hosts
	CheckInTime   string             `bson:"checkInTime,omitempty" json:"checkInTime,omitempty"`
	CheckOutTime  string             `bson:"checkOutTime,omitempty" json:"checkOutTime,omitempty"`
	IsActive      bool               `bson:"isActive" json:"isActive"`
}

type AccommodationAvailability struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AccommodationID primitive.ObjectID `bson:"accommodationID" json:"accommodation"`
	Date            time.Time          `bson:"date" json:"date"`
	AvailableRooms  int                `bson:"availableRooms" json:"availableRooms"`
	PricePerNight   float64            `bson:"pricePerNight" json:"pricePerNight"`
	Currency        string             `bson:"currency" json:"currency"`
	RoomType        string             `bson:"roomType" json:"roomType"`
}
