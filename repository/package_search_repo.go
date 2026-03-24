package repository

import (
	"context"
	"fmt"
	"time"

	model "github.com/Loboo34/travel/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PackageRepo struct {
	db *mongo.Database
}

func NewPackageRepo(db *mongo.Database) *PackageRepo {
	return &PackageRepo{db: db}
}

type PackageFilter struct {
	Destination string
	StartDate   time.Time
	EndDate     time.Time
	Travelers   model.TravelerCount
	MaxBudget   int64
	Tags        []model.PackageTag
	Components  []model.Component
	SortBy      model.PackageSortOption
	Page        int
	PageSize    int
}

func (r *PackageRepo) FindCandidates(
	ctx context.Context,
	f *PackageFilter,
) ([]model.Package, error) {

	filter := bson.M{
		"isActive":     true,
		"maxTravelers": bson.M{"$gte": f.Travelers},
		"destination":  bson.M{"$regex": f.Destination, "$options": "i"},
		// package window must overlap with requested dates
		// overlap condition: startDateFrom <= endDate AND startDateTo >= startDate
		"startDateFrom": bson.M{"$lte": f.EndDate},
		"startDateTo":   bson.M{"$gte": f.StartDate},
	}

	// only apply budget filter if a ceiling was set
	if f.MaxBudget > 0 {
		filter["basePrice"] = bson.M{"$lte": f.MaxBudget}
	}

	// package must have ALL requested tags
	if len(f.Tags) > 0 {
		filter["tags"] = bson.M{"$all": f.Tags}
	}

	// package must contain AT LEAST ONE of the requested component types
	if len(f.Components) > 0 {
		filter["includedComponents"] = bson.M{
			"$elemMatch": bson.M{
				"componentType": bson.M{"$in": f.Components},
			},
		}
	}

	opts := options.Find().
		SetSort(sortPackageFilter(f.SortBy)).
		SetSkip(int64((f.Page - 1) * f.PageSize)).
		SetLimit(int64(f.PageSize))

	cursor, err := r.db.Collection("packages").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("package candidate query: %w", err)
	}
	defer cursor.Close(ctx)

	var packages []model.Package
	if err := cursor.All(ctx, &packages); err != nil {
		return nil, fmt.Errorf("decoding packages: %w", err)
	}

	return packages, nil
}

func (r *PackageRepo) GetAvailability(
	ctx context.Context,
	packageID primitive.ObjectID,
) (*model.PackageAvailability, error) {

	var availability model.PackageAvailability

	err := r.db.Collection("package_availability").FindOne(ctx, bson.M{
		"packageID": packageID,
		"isActive":  true,
	}).Decode(&availability)

	if err == mongo.ErrNoDocuments {
		return nil, nil // no availability record — treat as unavailable
	}
	if err != nil {
		return nil, fmt.Errorf("package availability lookup: %w", err)
	}

	return &availability, nil
}

func (r *FlightRepo) FindActiveOffer(
	ctx context.Context,
	flightID primitive.ObjectID,
	minSeats int,
) (*model.FlightOffer, error) {

	var offer model.FlightOffer

	err := r.db.Collection("flight_offers").FindOne(ctx, bson.M{
		"flightID":      flightID,
		"isActive":      true,
		"bookableSeats": bson.M{"$gte": minSeats},
		"expiresAt":     bson.M{"$gt": time.Now()},
	}).Decode(&offer)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("flight offer lookup: %w", err)
	}

	return &offer, nil
}

func (r *AccommodationRepo) CheckAvailability(
	ctx context.Context,
	accommodationID primitive.ObjectID,
	checkIn time.Time,
	checkOut time.Time,
	minRooms int,
) (bool, int64, error) {

	nights := int(checkOut.Sub(checkIn).Hours() / 24)

	// count nights that have enough available rooms
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"accommodationID": accommodationID,
			"isActive":        true,
			"date":            bson.M{"$gte": checkIn, "$lt": checkOut},
			"$expr": bson.M{
				"$gte": []interface{}{
					bson.M{"$subtract": []interface{}{"$totalRooms", "$reservedRooms"}},
					minRooms,
				},
			},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":        nil,
			"nightCount": bson.M{"$sum": 1},
			"totalPrice": bson.M{"$sum": "$pricePerNight"},
		}}},
	}

	cursor, err := r.db.Collection("accommodation_availability").Aggregate(ctx, pipeline)
	if err != nil {
		return false, 0, fmt.Errorf("accommodation availability check: %w", err)
	}
	defer cursor.Close(ctx)

	var result struct {
		NightCount int   `bson:"nightCount"`
		TotalPrice int64 `bson:"totalPrice"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return false, 0, fmt.Errorf("decoding availability result: %w", err)
		}
	}

	// all nights must have availability
	if result.NightCount < nights {
		return false, 0, nil
	}

	return true, result.TotalPrice, nil
}

func (r *ActivityRepo) FindAvailableTimeslot(
	ctx context.Context,
	activityID primitive.ObjectID,
	date time.Time,
	participants int,
) (*model.ActivityTimeslot, error) {

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	var slot model.ActivityTimeslot

	err := r.db.Collection("activity_timeslots").FindOne(ctx, bson.M{
		"activityID":   activityID,
		"isActive":     true,
		"startTime":    bson.M{"$gte": startOfDay, "$lt": endOfDay},
		"groupSizeMax": bson.M{"$gte": participants},
		"$expr": bson.M{
			"$gt": []interface{}{
				bson.M{"$subtract": []interface{}{"$totalSlots", "$reservedSlots"}},
				0,
			},
		},
	}).Decode(&slot)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("activity timeslot lookup: %w", err)
	}

	return &slot, nil
}

func sortPackageFilter(s model.PackageSortOption) bson.D {
	switch s {
	case model.SortPackageByRating:
		return bson.D{{Key: "ratingSum", Value: -1}}
	case model.SortPackageByDuration:
		return bson.D{{Key: "durationDays", Value: 1}}
	default: // SortPackageByPrice
		return bson.D{{Key: "basePrice", Value: 1}}
	}
}
