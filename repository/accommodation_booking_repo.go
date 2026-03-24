package repository

import (
	"context"
	"fmt"
	"time"

	model "github.com/Loboo34/travel/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccommodationBookingRepo struct {
	db *mongo.Database
}

func NewAccommodationBookigRepo(db *mongo.Database) *AccommodationBookingRepo {
	return &AccommodationBookingRepo{db: db}
}

func (r *AccommodationBookingRepo) CheckAndReserv(ctx context.Context, accommodationID, roomTypeID primitive.ObjectID, checkIn, checkOut time.Time, rooms int) error {
	

	startOfDay := time.Date(checkIn.Year(), checkIn.Month(), checkIn.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(checkOut.Year(), checkOut.Month(), checkOut.Day(), 0, 0, 0, 0, time.UTC)
	nights := int(checkOut.Sub(checkIn).Hours() / 24)
	result, err := r.db.Collection("accommodation_availability").UpdateMany(ctx, bson.M{
		"_id":        accommodationID,
		"isActive":   true,
		"roomTypeID": roomTypeID,

		"$expr": bson.M{
			"$gte": []interface{}{
				bson.M{"$subrtact": []interface{}{"$totalRooms", "$reservedRooms"}},
			},
		},
		"date": bson.M{"$gte": startOfDay, "$lt": endOfDay},
	},
		bson.M{"$inc": bson.M{"reservedRooms": rooms}},
	)

	// if err == mongo.ErrNoDocuments {
	// 	return nil, errors.New("")
	// }

	if err != nil {
		return fmt.Errorf("Room reservation failed: %w", err)
	}

	if int(result.ModifiedCount) != nights {
		if result.MatchedCount > 0 {
			_ = r.ReleaseReservation(ctx, accommodationID, roomTypeID, checkIn, checkOut, rooms)
		}
		return fmt.Errorf("No rooms available")
	}

	return nil
}

func (r *AccommodationBookingRepo) ReleaseReservation(ctx context.Context, accommodationID, roomTypeID primitive.ObjectID, checkIn, checkOut time.Time, rooms int) error {
	startOfDay := time.Date(checkIn.Year(), checkIn.Month(), checkIn.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(checkOut.Year(), checkOut.Month(), checkOut.Day(), 0, 0, 0, 0, time.UTC)
	_, err := r.db.Collection("accommodation_availability").UpdateMany(ctx, bson.M{
		"accommodationID":        accommodationID,
		"roomTypeID": roomTypeID,
		"date":       bson.M{"$gte": startOfDay, "$lt": endOfDay},
	},
		bson.M{"$inc": bson.M{"reservedRooms": -rooms}})


		if err != nil {
			return fmt.Errorf("Error releasing rooms: %w", err)
		}

		return nil
}

func (r *AccommodationBookingRepo) CreateBooking(ctx context.Context, booking *model.AccommodationBooking) error{
	_, err := r.db.Collection("accommodation_booking").InsertOne(ctx, booking)
	if err != nil{
		return fmt.Errorf("error creating booking")
	}

	return nil
}

func (r *AccommodationBookingRepo) UpdateBooking(ctx context.Context, bookingID primitive.ObjectID, status model.BookingStatus, payment *model.Payment) error{
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

	_, err := r.db.Collection("accommodation_booking").UpdateOne(ctx, bson.M{"_id": bookingID}, update)
	if err != nil{
		return fmt.Errorf("error updating booking status")
	}

	return nil
}

func (r *AccommodationBookingRepo) GetBooking(ctx context.Context, bookingID primitive.ObjectID) (*model.AccommodationBooking, error) {
	var booking model.AccommodationBooking

	err := r.db.Collection("accommodation_booking").FindOne(ctx, bson.M{"_id": bookingID}).Decode(&booking)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed fetching booking: %w", err)
	}

	return &booking, nil
}


func (r *AccommodationBookingRepo) Cancel(ctx context.Context, bookingID primitive.ObjectID, reason string) error {

	update := bson.M{
		"$set": bson.M{
			"status": model.BookingStatusCanceled,
			"cancellationReason": reason,
			"refundStautus":      model.BookingStatusCanceled,
			"cancellationDate":   time.Now(),
			"updatedAt": time.Now(),
		},
	}
	_, err := r.db.Collection("accommodation_booking").UpdateOne(ctx, bson.M{"_id": bookingID}, update)
	if err != nil {
		return fmt.Errorf("error cancelling accommodation: %w", err)
	}
	return nil

}