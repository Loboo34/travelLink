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

type PackageBookingRepo struct {
	db *mongo.Database
}

func NewPackageBookingRepo(db *mongo.Database) *PackageBookingRepo {
	return &PackageBookingRepo{db: db}
}

func (r *PackageBookingRepo) ReserveSlot(ctx context.Context, packageID primitive.ObjectID, travelers int) error {
	err := r.db.Collection("package_availability").FindOneAndUpdate(ctx, bson.M{
		"_id":      packageID,
		"isActive": true,
		"$expr": bson.M{
			"$gte": []interface{}{
				bson.M{"$subtract": []interface{}{"$totalSlots", "$reservedSlots"}},
				travelers,
			},
		},
	},
		bson.M{"$inc": bson.M{"reservedSlots": travelers}},
	).Err()

	if err != nil {
		return fmt.Errorf("failed to reserve slot")
	}

	return nil
}

func (r *PackageBookingRepo) ReleaseSlot(ctx context.Context, packageID primitive.ObjectID, travelers int) error {
	_, err := r.db.Collection("package_availability").UpdateOne(ctx, bson.M{
		"packageID": packageID,
	}, bson.M{"$inc": bson.M{"reservedSlots": -travelers}})
	if err != nil {
		return fmt.Errorf("releasing package: %w", err)
	}

	return nil
}

func (r *PackageBookingRepo) CreateBooking(
	ctx context.Context,
	booking *model.PackageBooking,
) error {
	_, err := r.db.Collection("package_bookings").InsertOne(ctx, booking)
	if err != nil {
		return fmt.Errorf("creating package booking: %w", err)
	}
	return nil
}

func (r *PackageBookingRepo) UpdateBookingStatus(
	ctx context.Context,
	bookingID primitive.ObjectID,
	status model.BookingStatus,
	payment *model.Payment,
) error {
	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": time.Now().UTC(),
		},
	}
	if payment != nil {
		update["$set"].(bson.M)["payment"] = payment
		update["$set"].(bson.M)["amountPaid"] = payment.TotalAmount
	}
	_, err := r.db.Collection("package_bookings").UpdateOne(
		ctx,
		bson.M{"_id": bookingID},
		update,
	)
	if err != nil {
		return fmt.Errorf("updating package booking status: %w", err)
	}
	return nil
}

func (r *PackageBookingRepo) GetPackage(
	ctx context.Context,
	packageID primitive.ObjectID,
) (*model.Package, error) {
	var pkg model.Package

	err := r.db.Collection("packages").FindOne(ctx, bson.M{
		"_id":      packageID,
		"isActive": true,
	}).Decode(&pkg)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("fetching package: %w", err)
	}

	return &pkg, nil
}

func (r *PackageBookingRepo) GetBooking(ctx context.Context, bookingID primitive.ObjectID) (*model.PackageBooking, error) {
	var booking model.PackageBooking

	err := r.db.Collection("flight_booking").FindOne(ctx, bson.M{"_id": bookingID}).Decode(&booking)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed fetching booking: %w", err)
	}

	return &booking, nil
}

func (r *PackageBookingRepo) Cancel(ctx context.Context, bookingID primitive.ObjectID, reason string) error {

	update := bson.M{
		"$set": bson.M{
			"status": model.BookingStatusCanceled,
			"cancellationReason": reason,
			"refundStautus":      model.BookingStatusCanceled,
			"cancellationDate":   time.Now(),
			"updatedAt": time.Now(),
		},
	}
	_, err := r.db.Collection("package_bookings").UpdateOne(ctx, bson.M{"_id": bookingID}, update)
	if err != nil {
		return fmt.Errorf("error cancelling flight: %w", err)
	}
	return nil

}