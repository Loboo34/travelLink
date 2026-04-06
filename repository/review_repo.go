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

type ReviewRepo struct {
	db *mongo.Database
}

func NewReviewRepo(db *mongo.Database) *ReviewRepo {
	return &ReviewRepo{db: db}
}

func (r *ReviewRepo) Create(ctx context.Context, review *model.Review) error {
	_, err := r.db.Collection("reviews").InsertOne(ctx, review)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return &model.ConflictError{Message: "you have already reviewed this"}
		}
		return fmt.Errorf("creating review: %w", err)
	}
	return nil
}

func (r *ReviewRepo) GetByID(ctx context.Context, reviewID primitive.ObjectID) (*model.Review, error) {
	var review model.Review
	err := r.db.Collection("reviews").FindOne(ctx, bson.M{"_id": reviewID}).Decode(&review)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &model.NotFoundError{Resource: "review", ID: reviewID.Hex()}
		}
		return nil, fmt.Errorf("getting review: %w", err)
	}
	return &review, nil
}

// GetByTarget returns all reviews for a given target (accommodation, activity, etc.)
func (r *ReviewRepo) GetByTarget(ctx context.Context, referenceID primitive.ObjectID, reviewFor model.ReviewFor) ([]model.Review, error) {
	filter := bson.M{
		"referenceID": referenceID,
		"reviewFor":   reviewFor,
	}
	cursor, err := r.db.Collection("reviews").Find(ctx, filter, options.Find().SetSort(bson.M{"createdAt": -1}))
	if err != nil {
		return nil, fmt.Errorf("fetching reviews: %w", err)
	}
	defer cursor.Close(ctx)

	var reviews []model.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, fmt.Errorf("decoding reviews: %w", err)
	}
	return reviews, nil
}

func (r *ReviewRepo) Update(ctx context.Context, reviewID primitive.ObjectID, rating int, content string) error {
	update := bson.M{
		"$set": bson.M{
			"rating":    rating,
			"content":   content,
			"updatedAt": time.Now().UTC(),
		},
	}
	res, err := r.db.Collection("reviews").UpdateOne(ctx, bson.M{"_id": reviewID}, update)
	if err != nil {
		return fmt.Errorf("updating review: %w", err)
	}
	if res.MatchedCount == 0 {
		return &model.NotFoundError{Resource: "review", ID: reviewID.Hex()}
	}
	return nil
}

func (r *ReviewRepo) Delete(ctx context.Context, reviewID primitive.ObjectID) error {
	res, err := r.db.Collection("reviews").DeleteOne(ctx, bson.M{"_id": reviewID})
	if err != nil {
		return fmt.Errorf("deleting review: %w", err)
	}
	if res.DeletedCount == 0 {
		return &model.NotFoundError{Resource: "review", ID: reviewID.Hex()}
	}
	return nil
}

// UpdateTargetRating atomically recalculates and sets rating on the parent document.
// collection is e.g. "accommodations", "activities", "flights", "packages"
func (r *ReviewRepo) UpdateTargetRating(ctx context.Context, collection string, referenceID primitive.ObjectID, newRating int, delta int) error {
	update := bson.M{
		"$inc": bson.M{
			"ratingSum":   newRating * delta, // delta: +1 on create, -1 on delete
			"reviewCount": delta,
		},
	}
	_, err := r.db.Collection(collection).UpdateOne(ctx, bson.M{"_id": referenceID}, update)
	if err != nil {
		return fmt.Errorf("updating target rating: %w", err)
	}
	return nil
}