package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	UserID        string             `bson:"userId" json:"userId"`
	UserPhone     string             `bson:"userPhone" json:"userPhone"`
	Amount        float64            `bson:"ammount" json:"ammount"`
	PaymentMethod string             `bson:"paymentMethod" json:"paymentMethod"`
	Currency      string             `bson:"currency" json:"currency"`
	PaymentStatus string             `bson:"paymentStatus" json:"paymentStatus"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
}
