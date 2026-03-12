package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson"userID" json:"userID"`
	Content   string             `bson:"content" json:"content"`
	Note      string             `bson:"note" json:"note"`
	ReviewFor ReviewFor          `bson:"reviewFor" json:"reviewFor"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type ReviewFor string

const (
	AccommodationReview ReviewFor = "Accommodation"
	ActivityReview      ReviewFor = "Activity"
)
