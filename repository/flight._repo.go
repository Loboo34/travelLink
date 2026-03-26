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

func (r *FlightRepo) Update(ctx context.Context, flightID primitive.ObjectID) (*model.Flight, error) {
	var flight model.Flight

	if err := r.db.Collection("flights").FindOne(ctx, bson.M{"_id": flightID}).Decode(&flight); err != nil {
		return nil, fmt.Errorf("finding flight: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"departureTime": flight.DepartureTime,
			"arrivalTime":   flight.ArrivalTime,
			"stops":         flight.Stops,
			"updatedAt":     time.Now(),
		},
	}

	_, err := r.db.Collection("flghts").UpdateOne(ctx, bson.M{"_id": flightID}, update)
	if err != nil {
		return nil, fmt.Errorf("updating flight: %w", err)
	}

	return &flight, nil

}

func (r *FlightRepo) UpdateStatus(ctx context.Context, flightID primitive.ObjectID, status string) (*model.Flight, error) {
	var flight model.Flight

	if err := r.db.Collection("flights").FindOne(ctx, bson.M{"_id": flightID}).Decode(&flight); err != nil {
		return nil, fmt.Errorf("finding flight: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	_, err := r.db.Collection("flights").UpdateOne(ctx, bson.M{"_id": flightID}, update)
	if err != nil {
		return nil, fmt.Errorf("updating flight status: %w", err)
	}

	return &flight, nil

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

func (r *FlightRepo) UpdateOffer(ctx context.Context, offerID primitive.ObjectID) (*model.FlightOffer, error) {
	var offer model.FlightOffer

	update := bson.M{
		"$set": bson.M{
			"priceTotal":        offer.PriceTotal,
			"oneway":            offer.OneWay,
			"bookableSeats":     offer.BookableSeats,
			"lastTicketingDate": time.Now(),
			"expiresAt":         time.Now(),
		},
	}

	_, err := r.db.Collection("flight_offers").UpdateOne(ctx, bson.M{"_id": offerID}, update)
	if err != nil {
		return nil, fmt.Errorf("updating offer: %w", err)
	}

	return &offer, nil
}

func (r *FlightRepo) IsActive(ctx context.Context, offerID primitive.ObjectID, isActive bool) error {
	update := bson.M{
		"$set": bson.M{
			"isActive": isActive,
		},
	}

	_, err := r.db.Collection("flight_offers").UpdateOne(ctx, bson.M{"_id": offerID}, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("offer not found")
		}
		return fmt.Errorf("updating status: %w", err)
	}

	return nil

}

func (r *FlightRepo) DeleteOffer(ctx context.Context, offerID primitive.ObjectID) error {
	var offer model.FlightOffer

	err := r.db.Collection("flights").FindOneAndDelete(ctx, bson.M{"_id": offerID}).Decode(&offer)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("offer not found")
		}
		return fmt.Errorf("deliting offer: %w", err)
	}

	return nil
}