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

type PackageRepo struct {
	db *mongo.Database
}

func NewPackageRepo(db *mongo.Database) *PackageRepo {
	return &PackageRepo{db: db}
}

// CRUD Operations

func (r *PackageRepo) Add(ctx context.Context, pkg *model.Package) error {
	_, err := r.db.Collection("packages").InsertOne(ctx, pkg)
	if err != nil {
		return fmt.Errorf("adding package: %w", err)
	}
	return nil
}

func (r *PackageRepo) Update(ctx context.Context, packageID primitive.ObjectID, name, description, destination string, durationDays, maxTravelers int, basePrice int64, components []model.PackageComponent, tags []model.PackageTag, images []string) error {
	var pkg model.Package
	if err := r.db.Collection("packages").FindOne(ctx, bson.M{"_id": packageID}).Decode(&pkg); err != nil {
		return mongo.ErrNoDocuments
	}

	update := bson.M{
		"$set": bson.M{
			"name":               name,
			"description":        description,
			"destination":        destination,
			"durationDays":       durationDays,
			"maxTravelers":       maxTravelers,
			"basePrice":          basePrice,
			"includedComponents": components,
			"tags":               tags,
			"images":             images,
			"updatedAt":          time.Now(),
		},
	}

	_, err := r.db.Collection("packages").UpdateOne(ctx, bson.M{"_id": packageID}, update)
	if err != nil {
		return fmt.Errorf("updating package: %w", err)
	}

	return nil
}

func (r *PackageRepo) SetActive(ctx context.Context, packageID primitive.ObjectID, isActive bool) error {
	var pkg model.Package
	if err := r.db.Collection("packages").FindOne(ctx, bson.M{"_id": packageID}).Decode(&pkg); err != nil {
		return mongo.ErrNoDocuments
	}

	update := bson.M{
		"$set": bson.M{
			"isActive":  isActive,
			"updatedAt": time.Now(),
		},
	}

	_, err := r.db.Collection("packages").UpdateOne(ctx, bson.M{"_id": packageID}, update)
	if err != nil {
		return fmt.Errorf("updating package active state: %w", err)
	}

	return nil
}

func (r *PackageRepo) Delete(ctx context.Context, packageID primitive.ObjectID) error {
	result, err := r.db.Collection("packages").DeleteOne(ctx, bson.M{"_id": packageID})
	if err != nil {
		return fmt.Errorf("deleting package: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("package not found")
	}

	return nil
}

func (r *PackageRepo) GetByID(ctx context.Context, packageID primitive.ObjectID) (*model.Package, error) {
	var pkg model.Package
	err := r.db.Collection("packages").FindOne(ctx, bson.M{"_id": packageID}).Decode(&pkg)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("package not found")
		}
		return nil, fmt.Errorf("fetching package: %w", err)
	}

	return &pkg, nil
}
