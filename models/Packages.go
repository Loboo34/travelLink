package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Component string

const (
	ComponentFlight        Component = "Flight"
	ComponentAccommodation Component = "Accommodation"
	ComponentActivity      Component = "Activity"
)

type PackageTag string

const (
	PackageTagHoneymoon PackageTag = "honeymoon"
	PackageTagAdventure PackageTag = "adventure"
	PackageTagFamily    PackageTag = "family"
	PackageTagLuxury    PackageTag = "luxury"
	PackageTagBudget    PackageTag = "budget"
	PackageTagCultural  PackageTag = "cultural"
	PackageTagBeach     PackageTag = "beach"
)

type Package struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name               string             `bson:"name" json:"name"`
	Description        string             `bson:"description" json:"description"`
	Destination        string             `bson:"destination" json:"destination"`
	DurationDays       int                `bson:"durationDays" json:"durationDays"`
	StartDateFrom      *time.Time          `bson:"startDateFrom,omitempty" json:"startDateFrom,omitempty"`
	StartDateTo        *time.Time          `bson:"startDate_to,omitempty" json:"startDateTo,omitempty"`
	MaxTravelers       int                `bson:"maxTravelers" json:"maxTravelers"`
	BasePrice          int64              `bson:"basePrice" json:"basePrice"`
	Currency           string             `bson:"currency" json:"currency"`
	IncludedComponents []PackageComponent `bson:"includedComponents" json:"includedComponents"`
	Tags               []PackageTag       `bson:"tags" json:"tags,omitempty"`
	Images             []string           `bson:"images" json:"images,omitempty"`
	Rating             float64            `bson:"rating" json:"rating,omitempty"`
	ReviewCount        int                `bson:"reviewCount" json:"reviewCount,omitempty"`
	CachedAt           time.Time          `bson:"cached_at" json:"cached_at"`
	ExpiresAt          time.Time          `bson:"expiresAt,omitempty" json:"expiresAt,omitempty"`
	IsActive           bool               `bson:"isAcive" json:"isActive"`
}

type PackageComponent struct {
	ComponentType Component          `bson:"componentType" json:"componentType"`
	ReferenceID   primitive.ObjectID `bson:"referenceID" json:"referenceID"`
	Quantity      int                `bson:"quantity" json:"quantity"`
	Notes         string             `bson:"notes,omitempty" json:"notes,omitempty"`
}

type ComponentSummary struct {
	Type        Component `bson:"type" json:"type"` // "flight", "accommodation", "activity"
	Title       string    `bson:"title" json:"title"`
	Count       int       `bson:"count" json:"count"`
	ApproxPrice int64     `bson:"approxPrice" json:"approxPrice,omitempty"`
}

type PackageAvailability struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PackageID     primitive.ObjectID `bson:"packageID" json:"packageID"`
	TotalSlots    int                `bson:"totalSlots" json:"totalSlots"`
	ReservedSlots int                `bson:"reservedSlots" json:"reservedSlots"`
	IsActive      bool               `bson:"isActive" json:"isActive"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func (pa *PackageAvailability) AvailableSlots() int {
    return pa.TotalSlots - pa.ReservedSlots
}
