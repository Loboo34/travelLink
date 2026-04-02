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

type AccommodationRepo struct {
	db *mongo.Database
}

func NewAccommodationRepo(db *mongo.Database) *AccommodationRepo {
	return &AccommodationRepo{db: db}
}

func (r *AccommodationRepo) Add(ctx context.Context, accommodation *model.Accommodation) error {
	_, err := r.db.Collection("accommodations").InsertOne(ctx, accommodation)
	if err != nil {
		return fmt.Errorf("adding accommodation: %w", err)
	}

	return nil
}

func (r *AccommodationRepo) Update(ctx context.Context, accommodationID primitive.ObjectID, name, description string, amenities, images []string) error {
	var accom model.Accommodation
	if err := r.db.Collection("accommodations").FindOne(ctx, bson.M{"_id": accommodationID}).Decode(&accom); err != nil {
		return mongo.ErrNoDocuments
	}

	update := bson.M{
		"$set": bson.M{
			"name":        name,
			"description": description,
			"amenities":   amenities,
			"images":      images,
			"updatedAt":   time.Now(),
		},
	}

	_, err := r.db.Collection("accommodations").UpdateOne(ctx, bson.M{"_id": accommodationID}, update)
	if err != nil {
		return fmt.Errorf("updating accommodation: %w", err)
	}

	return nil
}

func (r *AccommodationRepo) Delete(ctx context.Context, accommodationID primitive.ObjectID) error {
	result, err := r.db.Collection("accommodations").DeleteOne(ctx, bson.M{"_id": accommodationID})
	if err != nil {
		return fmt.Errorf("deleting accommodation: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("not found: %w", err)
	}

	return nil
}

func (r *AccommodationRepo) Availability(ctx context.Context, availability *model.AccommodationAvailability) error {
	_, err := r.db.Collection("accommodation_availability").InsertOne(ctx, availability)
	if err != nil {
		return fmt.Errorf("creating availability: %w", err)
	}

	return nil
}

func (r *AccommodationRepo) IsActive(ctx context.Context, availabilityID primitive.ObjectID, isActive bool) error {

	_, err := r.db.Collection("accommodation_availability").UpdateOne(ctx, bson.M{
		"_id": availabilityID,
	}, bson.M{
		"$set": bson.M{
			"isActive": isActive,
		},
	})
	if err != nil {
		return fmt.Errorf("updating status: %w", err)
	}

	return nil
}

func (r *AccommodationRepo) Remove(ctx context.Context, availabilityID primitive.ObjectID) error {
	result, err := r.db.Collection("accommodation_availabiity").DeleteOne(ctx, bson.M{"_id": availabilityID})
	if err != nil {
		return fmt.Errorf("removing availability: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("not found: %w", err)
	}

	return nil
}

func (r *AccommodationRepo) GetAvailabilities(ctx context.Context) ([]model.AccommodationAvailability, error) {
 cursor, err := r.db.Collection("accommodation_Availability").Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("getting available accommodations: %w", err)
	}

	defer cursor.Close(ctx)

	var accoms []model.AccommodationAvailability
	if err := cursor.All(ctx, &accoms); err != nil{
		return nil, fmt.Errorf("decoding available accommodations: %w", err)
	}

	return accoms , nil
}

func (r *AccommodationRepo) GetAvailability(ctx context.Context, availabilityID primitive.ObjectID) (*model.AccommodationAvailability, error) {
	var accoms model.AccommodationAvailability
	if err := r.db.Collection("accommodation_Availability").FindOne(ctx, bson.M{"_id": availabilityID}).Decode(&accoms); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("accommodation not available: %w", err)
		}
		return nil, fmt.Errorf("getting available accommodations")

	}

	return &accoms, nil
}

func (r *AccommodationRepo) GetAccomodation(ctx context.Context, accommodationID primitive.ObjectID) (*model.Accommodation, error) {
	var accommodation model.Accommodation
	_, err := r.db.Collection("accommodations").Find(ctx, bson.M{"_id": accommodationID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("accommodation not available: %w", err)
		}
		return nil, fmt.Errorf("getting  accommodations")
	}

	return &accommodation, nil
}

func (r *AccommodationRepo) GetAccomodations(ctx context.Context) ([]model.Accommodation, error) {
	cursor, err := r.db.Collection("accommodations").Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("fetching accommodations: %w", err)
	}
	defer cursor.Close(ctx)

	var accoms []model.Accommodation
	if err := cursor.All(ctx, &accoms); err != nil{
		return nil, fmt.Errorf("decoding flights: %w", err)
	}

	return accoms, nil
}




