package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Activity struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title           string             `bson:"title" json:"title"`
	Price           float64            `bson:"price" json:"price"`
	Description     string             `bson:"description" json:"description"`
	Location        string             `bson:"location" json:"location"`
	Categories      []string           `bson:"categories" json:"categories"`
	DurationMinutes int                `bson:"duration_minutes" json:"duration_minutes"`
	Inclusions      []string           `bson:"inclusions" json:"inclusions,omitempty"`
	Exclusions      []string           `bson:"exclusions" json:"exclusions,omitempty"`
	Images          []string           `bson:"images" json:"images,omitempty"`
	Rating          float64            `bson:"rating" json:"rating,omitempty"`
	ReviewCount     int                `bson:"review_count" json:"review_count,omitempty"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	CachedAt          time.Time          `bson:"cached_at" json:"cachedAt"` 
}

type ActivityTimeslot struct {
    ID                primitive.ObjectID `bson:"_id,omitempty"`
    ActivityID        primitive.ObjectID `bson:"activity_id"`
    StartTime         time.Time          `bson:"start_time"`
    DurationMinutes   int                `bson:"duration_minutes"`
    AvailableSpots    int                `bson:"availableSpots"`
    PricePerPerson    float64            `bson:"priceperPerson"`
    Currency          string             `bson:"currency"`
    GroupSizeMax      int                `bson:"groupSizemax,omitempty"`
}
