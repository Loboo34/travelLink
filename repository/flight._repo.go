package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	model "github.com/Loboo34/travel/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FlightRepo struct {
	db *mongo.Database
}

func NewFlightRepo(db *mongo.Database) *FlightRepo {
	return &FlightRepo{db: db}
}

func (r *FlightRepo) Create(ctx context.Context, flight *model.Flight) error {
	_, err := r.db.Collection("flights").InsertOne(ctx, flight)
	if err != nil {
		return fmt.Errorf("creating flight: %w", err)
	}

	return nil
}

func (r *FlightRepo) Update(ctx context.Context, flightID primitive.ObjectID, depTime, arrTime time.Time, stops int) error {
	var flight model.Flight

	if err := r.db.Collection("flights").FindOne(ctx, bson.M{"_id": flightID}).Decode(&flight); err != nil {
		return mongo.ErrNoDocuments
	}

	update := bson.M{
		"$set": bson.M{
			"departureTime": depTime,
			"arrivalTime":   arrTime,
			"stops":         stops,
			"updatedAt":     time.Now(),
		},
	}

	_, err := r.db.Collection("flghts").UpdateOne(ctx, bson.M{"_id": flightID}, update)
	if err != nil {
		return fmt.Errorf("updating flight: %w", err)
	}

	// if result.MatchedCount == 0 {
	// 	return fmt.Errorf("flight  not found: %s", flightID.Hex())
	// }

	return nil

}

func (r *FlightRepo) UpdateStatus(ctx context.Context, flightID primitive.ObjectID, status *model.FlightStatus) error {

	var flight model.Flight

	if err := r.db.Collection("flights").FindOne(ctx, bson.M{"_id": flightID}).Decode(&flight); err != nil {
		return mongo.ErrNoDocuments
	}

	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	_, err := r.db.Collection("flights").UpdateOne(ctx, bson.M{"_id": flightID}, update)
	if err != nil {
		return fmt.Errorf("updating flight status: %w", err)
	}

	return nil

}

func (r *FlightRepo) Delete(ctx context.Context, flightID primitive.ObjectID) error {
	var flight model.Flight

	err := r.db.Collection("flights").FindOneAndDelete(ctx, bson.M{"_id": flightID}).Decode(&flight)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("flight not found")
		}
		return fmt.Errorf("deliting flight: %w", err)
	}

	return nil
}

func (r *FlightRepo) CreateOffer(ctx context.Context, offer *model.FlightOffer) error {
	_, err := r.db.Collection("flight_offers").InsertOne(ctx, offer)
	if err != nil {
		return fmt.Errorf("creating flight offer: %w", err)
	}

	return nil
}

func (r *FlightRepo) UpdateOffer(ctx context.Context, offerID primitive.ObjectID, price, seats int, oneWay bool) error {

	update := bson.M{
		"$set": bson.M{
			"priceTotal":        price,
			"oneway":            oneWay,
			"bookableSeats":     seats,
			"lastTicketingDate": time.Now(),
			"expiresAt":         time.Now(),
		},
	}

	result, err := r.db.Collection("flight_offers").UpdateOne(ctx, bson.M{"_id": offerID}, update)
	if err != nil {
		return fmt.Errorf("updting offer")
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("flight offer not found: %s", offerID.Hex())
	}

	return nil
}

func (r *FlightRepo) IsActive(ctx context.Context, offerID primitive.ObjectID, isActive bool) error {
	update := bson.M{
		"$set": bson.M{
			"isActive":  isActive,
			"updatedAt": time.Now(),
		},
	}

	result, err := r.db.Collection("flight_offers").UpdateOne(
		ctx,
		bson.M{"_id": offerID},
		update,
	)
	if err != nil {
		return fmt.Errorf("database error updating offer status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("flight offer not found: %s", offerID.Hex())
	}

	if result.ModifiedCount == 0 && isActive {

	}

	return nil
}
func (r *FlightRepo) DeleteOffer(ctx context.Context, offerID primitive.ObjectID) error {
	result, err := r.db.Collection("flight_offers").DeleteOne(ctx, bson.M{"_id": offerID})
	if err != nil {
		return fmt.Errorf("hard delete offer: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("not found: %w", err)
	}

	return nil
}
