package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateFlightIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "isActive", Value: 1},
				{Key: "bookableSeats", Value: 1},
				{Key: "expiresAt", Value: 1},
			},
		},
	}

	if _, err := db.Collection("flight_offers").Indexes().CreateMany(ctx, indexes); err != nil {
		return fmt.Errorf("creating flight_offers indexes: %w", err)
	}

	return nil
}

func CreateAccommodationIndexes(ctx context.Context, db *mongo.Database) error {
	accomIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "location", Value: "2dsphere"}}},
		{Keys: bson.D{
			{Key: "address.city", Value: 1},
			{Key: "address.country", Value: 1},
			{Key: "type", Value: 1},
		}},
	}

	if _, err := db.Collection("accommodations").Indexes().CreateMany(ctx, accomIndexes); err != nil {
		return fmt.Errorf("creating accommodations indexes: %w", err)
	}

	availabilityIndexes := []mongo.IndexModel{
		{Keys: bson.D{
			{Key: "isActive", Value: 1},
			{Key: "date", Value: 1},
			{Key: "accommodationID", Value: 1},
		}},
	}

	if _, err := db.Collection("accommodation_availability").Indexes().CreateMany(ctx, availabilityIndexes); err != nil {
		return fmt.Errorf("creating accommodation_availability indexes: %w", err)
	}

	return nil
}

func CreateActivityIndexes(ctx context.Context, db *mongo.Database) error {
	activityIndexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "location", Value: "2dsphere"}}},
		{Keys: bson.D{
			{Key: "city", Value: 1},
			{Key: "country", Value: 1},
			{Key: "categories", Value: 1},
		}},
	}

	if _, err := db.Collection("activities").Indexes().CreateMany(ctx, activityIndexes); err != nil {
		return fmt.Errorf("creating activities indexes: %w", err)
	}

	timeslotIndexes := []mongo.IndexModel{
		{Keys: bson.D{
			{Key: "isActive", Value: 1},
			{Key: "date", Value: 1},
			{Key: "activityID", Value: 1},
		}},
	}

	if _, err := db.Collection("activity_timeslots").Indexes().CreateMany(ctx, timeslotIndexes); err != nil {
		return fmt.Errorf("creating activity_timeslots indexes: %w", err)
	}

	return nil
}

func CreateBookingIndexes(ctx context.Context, db *mongo.Database) error {
	collections := []string{"flight_bookings", "accommodation_bookings", "activity_bookings", "package_bookings"}
	index := mongo.IndexModel{Keys: bson.D{{Key: "userID", Value: 1}, {Key: "status", Value: 1}}}

	for _, coll := range collections {
		if _, err := db.Collection(coll).Indexes().CreateMany(ctx, []mongo.IndexModel{index}); err != nil {
			return fmt.Errorf("creating %s indexes: %w", coll, err)
		}
	}

	return nil
}

func CreateReviewIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "targetID", Value: 1}, {Key: "targetType", Value: 1}}},
		{Keys: bson.D{{Key: "userID", Value: 1}, {Key: "targetID", Value: 1}, {Key: "targetType", Value: 1}}, Options: options.Index().SetUnique(true)},
	}

	if _, err := db.Collection("reviews").Indexes().CreateMany(ctx, indexes); err != nil {
		return fmt.Errorf("creating reviews indexes: %w", err)
	}

	return nil
}

func CreateUserIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "role", Value: 1}}},
	}

	if _, err := db.Collection("users").Indexes().CreateMany(ctx, indexes); err != nil {
		return fmt.Errorf("creating users indexes: %w", err)
	}

	return nil
}
