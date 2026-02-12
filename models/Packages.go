package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Package struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Accommodation string             `bson:"accommodation" json:"accommodation"`
	Flight        string             `bson:"flight" json:"flight"`
	Tour          string             `bson:"tour" json:"tour"`
	Price         float64               `bson:"price" json:"price"`
	Duration      string             `bson:"duration" json:"duration"`
	User          []string           `bson:"user" json:"users"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
}
