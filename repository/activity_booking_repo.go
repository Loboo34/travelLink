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

type ActivityBookingRepo struct {
	db *mongo.Database
}

func NewActivityBookingRepo(db *mongo.Database) *ActivityBookingRepo {
	return &ActivityBookingRepo{db: db}
}

func (r *ActivityBookingRepo) CheckAndReserv(ctx context.Context, activityID, timeSlotID primitive.ObjectID, participants int) (*model.ActivityTimeslot, error) {
	var slot model.ActivityTimeslot

	err := r.db.Collection("activity_timeslots").FindOneAndUpdate(ctx, bson.M{
		"activityID": activityID,
		"_id":        timeSlotID,
		"isActive":   true,
		//"expiresAt":  bson.M{"$gt": time.Now()},
		"$expr": bson.M{"$gte": []interface{}{
			bson.M{"$subtract": []interface{}{"$totalSlots", "$reservedSlots"}},
			participants,
		}},
	}, bson.M{"$inc": bson.M{"reservedSlots": participants}},
		options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&slot)

	if err == mongo.ErrNoDocuments {
		return nil, errors.New("slots available")
	}

	if err != nil {
		return nil, fmt.Errorf(" reservation failed: %w", err)
	}

	return &slot, nil
}

func (r *ActivityBookingRepo) ReleaseReservation(ctx context.Context, activityID, timeSlotID primitive.ObjectID, participants int) error {
	_, err := r.db.Collection("activity_timeslots").UpdateOne(ctx, bson.M{
		"_id": timeSlotID,
	},
		bson.M{"$inc": bson.M{"reservedSlots": -participants}})

	if err != nil {
		return fmt.Errorf("error releasing slots: %w", err)
	}

	return nil
}

func (r *ActivityBookingRepo) CreateBooking(ctx context.Context, booking *model.ActivityBooking) error {
	_, err := r.db.Collection("activity_booking").InsertOne(ctx, booking)
	if err != nil {
		return fmt.Errorf("error creating booking: %w", err)
	}

	return nil
}

func (r *ActivityBookingRepo) UpdateBooking(ctx context.Context, bookingID primitive.ObjectID, status model.BookingStatus, payment *model.Payment) error {

	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": time.Now(),
		},
	}

	if payment != nil {
		update["$set"].(bson.M)["payment"] = payment
		update["$set"].(bson.M)["amountPaid"] = payment.TotalAmount
	}

	_, err := r.db.Collection("activity_booking").UpdateOne(ctx, bson.M{"_id": bookingID}, update)
	if err != nil {
		return fmt.Errorf("error updating status")
	}
	return nil
}
