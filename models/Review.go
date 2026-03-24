package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewFor string

const (
    ReviewForAccommodation ReviewFor = "accommodation"
    ReviewForActivity      ReviewFor = "activity"
    ReviewForFlight        ReviewFor = "flight"
    ReviewForPackage       ReviewFor = "package"
)

type Review struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserID      primitive.ObjectID `bson:"userID" json:"userID"`
    ReferenceID primitive.ObjectID `bson:"referenceID" json:"referenceID"` 
    ReviewFor   ReviewFor          `bson:"reviewFor" json:"reviewFor"`      
    Rating      int                `bson:"rating" json:"rating"`            
    Content     string             `bson:"content" json:"content"`
    IsVerified  bool               `bson:"isVerified" json:"isVerified"`    
    CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
    UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}