package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Car struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OriginID primitive.ObjectID `bson:"originID" json:"originID"`
	
}

type CarAvailability struct {
}
