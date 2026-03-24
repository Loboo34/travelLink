package repository

import (
	"context"
	"fmt"
	"time"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ActivityRepo struct {
	db *mongo.Database
}

func NewActivityRepo(db *mongo.Database) *ActivityRepo {
	return &ActivityRepo{db: db}
}

type ActivityFilter struct {
	Location     model.LocationSearch
	Participants model.TravelerCount
	ForAllAges   bool
	Category     model.ActivityCategory
	Date         time.Time
	Duration     int
	SortBy       model.ActivitySortOptions
	Page         int
	PageSize     int
}

func (r *ActivityRepo) SearchActivity(ctx context.Context, a *ActivityFilter) ([]model.ActivitySearchResult, error) {

	totalParticipants := a.Participants.Adults + a.Participants.Children + a.Participants.Infants
	pipeline := mongo.Pipeline{
		stageMatchTimeSlot(a.Date, totalParticipants),
		stageLookUpActivities(),
		stageUnwindActivity(),
		stageMatchCategories(a.Location, a.Category, a.ForAllAges, a.Duration),
		stageSortTimeslot(a.SortBy),
		bson.D{{Key: "$skip", Value: (a.Page - 1) * a.PageSize}},
		bson.D{{Key: "$limit", Value: a.PageSize}},
	}

	cursor, err := r.db.Collection("activity_timeslot").Aggregate(ctx, pipeline)
	if err != nil {
		utils.Logger.Warn("Error aggragating activities")
		return nil, fmt.Errorf("Failed activity aggragation: %w", err)
	}

	defer cursor.Close(ctx)

	var activities []model.ActivitySearchResult
	if err := cursor.All(ctx, &activities); err != nil {
		utils.Logger.Warn("Failed to decode activities")
		return nil, fmt.Errorf("%w", err)
	}

	return activities, nil
}

func stageMatchTimeSlot(date time.Time, participants int) bson.D {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)
	return bson.D{{Key: "$match", Value: bson.M{
		"isActive": true,
		"date":     bson.M{"$gte": startOfDay, "$lt": endOfDay},
		"$expr": bson.M{
			"$sgt": []interface{}{
				bson.M{
					"$subtract": []interface{}{"totalSlots", "reservedSlots"}},
				0,
			},
		},
	}}}
}

func stageLookUpActivities() bson.D {
	return bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "activities",
		"localField":   "activityID",
		"foreignField": "_id",
		"as":           "activityDoc",
	}}}
}

func stageUnwindActivity() bson.D {
	return bson.D{{Key: "$unwind", Value: "$activityDoc"}}
}

func stageMatchCategories(location model.LocationSearch, category model.ActivityCategory, forAllAges bool,
	maxDuration int) bson.D {
	filter := bson.M{
		"activityDoc.isActive": true,
	}

	if location.Longitude != 0 || location.Latitude != 0 {
		filter["activityDoc.location"] = bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{location.Longitude, location.Latitude},
				},
				"$maxDistance": location.RadiusKm * 1000,
			},
		}
	} else if location.City != "" {
		filter["activityDoc.city"] = bson.M{"$regex": location.City, "$options": "i"}
		if location.Country != "" {
			filter["activityDoc.country"] = bson.M{"$regex": location.Country, "$options": "i"}
		}

	}

	 if category != "" {
        filter["activityDoc.categories"] = bson.M{"$in": []model.ActivityCategory{category}}
    }

	if maxDuration > 0 {
        filter["activityDoc.durationMinutes"] = bson.M{"$lte": maxDuration}
    }


	return bson.D{{Key: "$match", Value: filter}}
}

func stageSortTimeslot(sort model.ActivitySortOptions) bson.D {
	switch sort {
	case model.SortActivityByRating:
		return bson.D{{Key: "$sort", Value: bson.D{
			{Key: "activityDoc.rating", Value: -1},
			{Key: "totalPrice", Value: 1},
		}}}
	default:
		return bson.D{{Key: "$sort", Value: bson.D{
			{Key: "totalPrice", Value: 1},
			{Key: "activityDoc.rating", Value: -1},
		}}}
	}
}
