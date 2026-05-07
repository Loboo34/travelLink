package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivityCategory string

const (
	ActivityCategoryAdventure   ActivityCategory = "adventure"
	ActivityCategoryWellness    ActivityCategory = "wellness"
	ActivityCategoryCultural    ActivityCategory = "cultural"
	ActivityCategoryFood        ActivityCategory = "food_and_drink"
	ActivityCategoryNature      ActivityCategory = "nature"
	ActivityCategorySightseeing ActivityCategory = "sightseeing"
	ActivityCategoryWater       ActivityCategory = "water_sports"
	ActivityCategoryNightlife   ActivityCategory = "nightlife"
)

type MeetingPoint struct {
	Label    string      `bson:"label" json:"label"`
	Location GeoLocation `bson:"location" json:"location"`
}

type Activity struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	City            string             `bson:"city" json:"city"`
	Country         string             `bson:"country" json:"country"`
	Location        GeoLocation        `bson:"location" json:"location"`
	MeetingPoint    MeetingPoint       `bson:"meetingPoint" json:"meetingPoint"`
	Categories      []ActivityCategory `bson:"categories" json:"categories"`
	DurationMinutes int                `bson:"durationMinutes" json:"durationMinutes"`
	Inclusions      []string           `bson:"inclusions" json:"inclusions,omitempty"`
	Exclusions      []string           `bson:"exclusions" json:"exclusions,omitempty"`
	Images          []string           `bson:"images" json:"images,omitempty"`
	Rating          float64            `bson:"rating" json:"rating,omitempty"`
	ReviewCount     int                `bson:"reviewCount" json:"reviewCount,omitempty"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
	CachedAt        *time.Time          `bson:"cachedAt" json:"cachedAt"`
	IsActive        bool               `bson:"isActive" json:"isActive"`
}

type ActivityTimeslot struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ActivityID      primitive.ObjectID `bson:"activityID" json:"activityID"`
	StartTime       time.Time          `bson:"startTime" json:"startTime"`
	DurationMinutes int                `bson:"durationMinutes" json:"durationMinutes"`
	TotalSlots      int                `bson:"totalSlots" json:"totalSlots"`
	ReservedSlots   int                `bson:"reservedSlots" jso:"reservedSlots"`
	PricePerPerson  int64              `bson:"pricePerPerson" json:"pricePerPerson"`
	GroupSizeMax    int                `bson:"groupSizeMax,omitempty" json:"groupSizeMax"`
	IsActive        bool               `bson:"isActive" json:"isActive"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func (pa *PackageAvailability) ActivityTimeslot() int {
	return pa.TotalSlots - pa.ReservedSlots
}

type ActivitySearchResult struct {
    ActivityID    primitive.ObjectID `bson:"_id" json:"activityID"`
    Activity      Activity     `bson:"activityDoc" json:"activity"`
    PricePerPerson int64             `bson:"pricePerPerson" json:"pricePerPerson"`
    AvailableSlots int               `bson:"availableSlots" json:"availableSlots"`
    StartTime     time.Time          `bson:"startTime" json:"startTime"`
    Currency      string             `bson:"currency" json:"currency"`
}
