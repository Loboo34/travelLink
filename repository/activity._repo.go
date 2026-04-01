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

type ActivityRepo struct {
	db *mongo.Database
}

func NewActivityRepo(db *mongo.Database) *ActivityRepo {
	return &ActivityRepo{db: db}
}

func (r *ActivityRepo) Add(ctx context.Context, activity *model.Activity) error {
	_, err := r.db.Collection("activities").InsertOne(ctx, activity)
	if err != nil {
		return fmt.Errorf("adding activity")
	}

	return nil
}

func (r *ActivityRepo) Update(ctx context.Context, activityID primitive.ObjectID, title string, duration int, inclusion, exclusion []string, point *model.MeetingPoint) error {
	var activity model.Activity
	if err := r.db.Collection("activities").FindOne(ctx, bson.M{"_id": activityID}).Decode(&activity); err != nil {
		return mongo.ErrNoDocuments
	}

	update := bson.M{
		"$set": bson.M{
			"title":           title,
			"meetingPoint":    point,
			"durationMinutes": duration,
			"inclusions":      inclusion,
			"exclusion":       exclusion,
			"updatedAt":       time.Now(),
		},
	}

	_, err := r.db.Collection("activities").UpdateOne(ctx, bson.M{"_id": activityID}, update)
	if err != nil {
		return fmt.Errorf("updating activity: %w", err)
	}

	return nil
}

func (r *ActivityRepo) SetActive(ctx context.Context, activityID primitive.ObjectID, isActive bool) error {
	var activity model.Activity
	if err := r.db.Collection("activities").FindOne(ctx, bson.M{"_id": activityID}).Decode(&activity); err != nil {
		return mongo.ErrNoDocuments
	}

	update := bson.M{"$set": bson.M{"isActive": isActive, "updatedAt": time.Now()}}
	_, err := r.db.Collection("activities").UpdateOne(ctx, bson.M{"_id": activityID}, update)
	if err != nil {
		return fmt.Errorf("updating activity active flag: %w", err)
	}

	return nil
}

func (r *ActivityRepo) Delete(ctx context.Context, activityID primitive.ObjectID) error {
	result, err := r.db.Collection("activities").DeleteOne(ctx, bson.M{"_id": activityID})
	if err != nil {
		return fmt.Errorf("deleting activity: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("not found: %w", err)
	}

	return nil
}

func (r *ActivityRepo) Timeslot(ctx context.Context, timeSlot *model.ActivityTimeslot) error {
	_, err := r.db.Collection("activity_timeslot").InsertOne(ctx, timeSlot)
	if err != nil {
		return fmt.Errorf("adding activity slot")
	}

	return nil
}

func (r *ActivityRepo) IsActive(ctx context.Context, timeslotID primitive.ObjectID, isActive bool) error {
	var activity model.ActivityTimeslot
	if err := r.db.Collection("activity_timeslot").FindOne(ctx, bson.M{"_id": timeslotID}).Decode(&activity); err != nil {
		return mongo.ErrNoDocuments
	}

	update := bson.M{
		"$set": bson.M{
			"isActive":  isActive,
			"updatedAt": time.Now(),
		},
	}

	_, err := r.db.Collection("activities").UpdateOne(ctx, bson.M{"_id": timeslotID}, update)
	if err != nil {
		return fmt.Errorf("updating activity: %w", err)
	}

	return nil
}

func (r *ActivityRepo) UpdateTimeslot(ctx context.Context, timeSlotID primitive.ObjectID, startTime time.Time, duration, totalSlots, groupSize int, price int64) error {
	var activity model.ActivityTimeslot
	if err := r.db.Collection("activity_timeslot").FindOne(ctx, bson.M{"_id": timeSlotID}).Decode(&activity); err != nil {
		return mongo.ErrNoDocuments
	}

	update := bson.M{
		"$set": bson.M{
			"startTime":       startTime,
			"durationMinutes": duration,
			"totalSlots":      totalSlots,
			"pricePerPerson":  price,
			"groupSizeMax":    groupSize,
			"updatedAt":       time.Now(),
		},
	}

	_, err := r.db.Collection("activity_timeslot").UpdateOne(ctx, bson.M{"_id": timeSlotID}, update)
	if err != nil {
		return fmt.Errorf("updating activity: %w", err)
	}

	return nil
}
