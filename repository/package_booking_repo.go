package repository

import (
	"context"
	"fmt"

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
		"_id": packageID,
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

func (r *PackageBookingRepo) ReleaseSlot(ctx context.Context, packageID primitive.ObjectID, travelers int) error{
	_, err := r.db.Collection("package_availability").UpdateOne(ctx, bson.M{
		"packageID": packageID,
	}, bson.M{"$inc": bson.M{"reservedSlots": -travelers}},)
	if err != nil{
		return fmt.Errorf("releasing package: %w", err)
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
