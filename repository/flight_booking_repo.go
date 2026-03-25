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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FlightBookingRepo struct {
	db *mongo.Database
}

func NewFlightBookingRepo(db *mongo.Database) *FlightBookingRepo {
	return &FlightBookingRepo{db: db}
}

func (r *FlightBookingRepo) CheckAndReserv(ctx context.Context, flightID primitive.ObjectID, seats int) (*model.FlightOffer, error) {
	var offer model.FlightOffer

	err := r.db.Collection("flight_offers").FindOneAndUpdate(ctx, bson.M{
		"_id":           flightID,
		"isActive":      true,
		"expiresAt":     bson.M{"$gt": time.Now()},
		"bookableSeats": bson.M{"$gte": seats},
	},

		bson.M{"$inc": bson.M{"bookableSeats": -seats}},
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&offer)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("No seats available for booking")
	}

	if err != nil {
		return nil, fmt.Errorf("flight reservation failed: %w", err)
	}

	return &offer, nil
}

func (r *FlightBookingRepo) ReleaseReservation(ctx context.Context, flightID primitive.ObjectID, seats int) error {
	_, err := r.db.Collection("flight_offers").UpdateOne(ctx,
		bson.M{"_id": flightID},
		bson.M{"$inc": bson.M{"bookableSeats": seats}},
	)

	if err != nil {
		return fmt.Errorf("Error releasing reserved seats: %w", err)
	}

	return nil
}

func (r *FlightBookingRepo) CreateBooking(ctx context.Context, booking *model.FlightBooking) error{
	_, err := r.db.Collection("flight_booking").InsertOne(ctx, booking)
	if err != nil{
		return fmt.Errorf("error creating booking: %w", err)
	}
	return  nil
}

func (r *FlightBookingRepo) UpdateBooking(ctx context.Context, bookingID primitive.ObjectID, status model.BookingStatus, payment *model.Payment ) error {
	update := bson.M{
		"$set": bson.M{
			"status": status,
			"updatedAt": time.Now(),
		},
	}

	if payment != nil{
		update["$set"].(bson.M)["payment"] = payment
        update["$set"].(bson.M)["amountPaid"] = payment.TotalAmount
	}

	_, err := r.db.Collection("flight_booking").UpdateOne(ctx, bson.M{"bookingID": bookingID}, update)
	if err != nil {
		return fmt.Errorf("Error updating booking status: %w", err)
	}

	return nil 
}

func (r *FlightBookingRepo) GetBooking(ctx context.Context, bookingID primitive.ObjectID) (*model.FlightBooking, error) {
	var booking model.FlightBooking

	err := r.db.Collection("flight_booking").FindOne(ctx, bson.M{"_id": bookingID}).Decode(&booking)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed fetching booking: %w", err)
	}

	return &booking, nil
}

func (r *FlightBookingRepo) Cancel(ctx context.Context, bookingID primitive.ObjectID, reason string) error {

	update := bson.M{
		"$set": bson.M{
			"status": model.BookingStatusCanceled,
			"cancellationReason": reason,
			"refundStautus":      model.BookingStatusCanceled,
			"cancellationDate":   time.Now(),
			"updatedAt": time.Now(),
		},
	}
	_, err := r.db.Collection("flight_booking").UpdateOne(ctx, bson.M{"_id": bookingID}, update)
	if err != nil {
		return fmt.Errorf("error cancelling flight: %w", err)
	}
	return nil

}
