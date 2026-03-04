package utils

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateFlightIndexes(ctx context.Context, collection *mongo.Collection) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "originID", Value: 1},
			{Key: "destinationID", Value: 1},
			{Key: "departureTime", Value: 1},
		},
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}