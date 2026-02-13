package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Package struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name                string             `bson:"name" json:"name"` 
    Description         string             `bson:"description" json:"description"`
    Destination         string             `bson:"destination" json:"destination"`
    DurationDays        int                `bson:"duration_days" json:"duration_days"`
    StartDateFrom       time.Time         `bson:"startDateFrom,omitempty" json:"startDateFrom,omitempty"` 
    StartDateTo         time.Time         `bson:"startDate_to,omitempty" json:"startDateTo,omitempty"`
    BasePrice           float64            `bson:"basePrice" json:"basePrice"`
    Currency            string             `bson:"currency" json:"currency"`
    IncludedComponents  []ComponentSummary `bson:"includedComponents" json:"includedComponents"` 
    Tags                []string           `bson:"tags" json:"tags,omitempty"` 
    Images              []string           `bson:"images" json:"images,omitempty"`
    Rating              float64            `bson:"rating" json:"rating,omitempty"`
    ReviewCount         int                `bson:"reviewCount" json:"reviewCount,omitempty"`
    CachedAt            time.Time          `bson:"cached_at" json:"cached_at"`
    ExpiresAt           time.Time         `bson:"expiresAt,omitempty" json:"expiresAt,omitempty"`
}

type ComponentSummary struct {
    Type        string  `bson:"type" json:"type"` // "flight", "accommodation", "activity"
    Title       string  `bson:"title" json:"title"` 
    Count       int     `bson:"count" json:"count"` 
    ApproxPrice float64 `bson:"approxPrice" json:"approxPrice,omitempty"`
}
