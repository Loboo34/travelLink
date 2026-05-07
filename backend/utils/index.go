package utils

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateFlightIndexes(ctx context.Context, db *mongo.Database) error {

	_, err := db.Collection("flight_offers").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "isActive", Value: 1},
				{Key: "expiresAt", Value: 1},
				{Key: "bookableSeats", Value: 1},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("flight_offers index: %w", err)
	}
	_, err = db.Collection("flights").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "originID", Value: 1},
				{Key: "destinationID", Value: 1},
				{Key: "departureTime", Value: 1},
				{Key: "status", Value: 1},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("flights index: %w", err)
	}

	_, err = db.Collection("airports").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "code", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("airports index: %w", err)
	}

	return nil
}
