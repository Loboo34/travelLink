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

type AccommodationRepo struct {
	db *mongo.Database
}

func NewAccommodationRepo(db *mongo.Database) *AccommodationRepo {
	return &AccommodationRepo{db: db}
}

type AccommodationFilter struct {
	Location     model.LocationSearch
	CheckInDate  time.Time
	CheckOutDate time.Time
	Guests       model.GuestCount
	PropertyType model.PropertyType
	Rooms        int
	SortBy       model.AccommodationSortOption
	Page         int
	PageSize     int
}

func (r *AccommodationRepo) SearchAccommodationAvailability(ctx context.Context, a *AccommodationFilter) ([]model.AccommodationSearchResult, error) {

	nights := int(a.CheckOutDate.Sub(a.CheckInDate).Hours() / 24)
	totalGuests := a.Guests.Adults + a.Guests.Children + a.Guests.Infants
	pipeline := mongo.Pipeline{
		stageMatchAvailability(a.Rooms, a.CheckInDate, a.CheckOutDate),
		stageLookUpProperties(),
		stageUnwindAccommodations(),
		stageMatchesProperties(a.Location, a.PropertyType, totalGuests),
		stageGroupByAccommodation(nights),
		stageSortAvailability(a.SortBy),
		bson.D{{Key: "$match", Value: bson.M{"availableNights": nights}}},
		bson.D{{Key: "$skip", Value: (a.Page - 1) * a.PageSize}},
		bson.D{{Key: "$limit", Value: a.PageSize}},
	}

	cursor, err := r.db.Collection("accommodation_availability").Aggregate(ctx, pipeline)
	if err != nil {
		utils.Logger.Warn("Error aggragating accommodations")
		return nil, fmt.Errorf("Failed accommodation aggragation: %w", err)
	}

	defer cursor.Close(ctx)

	var accommodations []model.AccommodationSearchResult
	if err := cursor.All(ctx, &accommodations); err != nil {
		utils.Logger.Warn("Failed to decode accommodation availability")
		return nil, fmt.Errorf("%w", err)
	}

	return accommodations, nil

}

func stageMatchAvailability(rooms int, checkIn, checkOut time.Time) bson.D {
	startOfDay := time.Date(checkIn.Year(), checkIn.Month(), checkIn.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(checkOut.Year(), checkOut.Month(), checkOut.Day(), 0, 0, 0, 0, time.UTC)

	return bson.D{{Key: "$match", Value: bson.M{
		"isActive": true,
		"$expr": bson.M{"$gte": []interface{}{
			bson.M{"$subtract": []interface{}{"$totalRooms", "$reservedRooms"}},
			rooms,
		}},
		"date": bson.M{"$gte": startOfDay, "$lt": endOfDay},
	}}}
}

func stageLookUpProperties() bson.D {
	return bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "accommodations",
		"localField":   "accommodationID",
		"foreignField": "_id",
		"as":           "accommodationDocs",
	}}}
}

func stageUnwindAccommodations() bson.D {
	return bson.D{{Key: "$unwind", Value: "$accommodationDocs"}}
}

func stageMatchesProperties(location model.LocationSearch, property model.PropertyType, totalGuests int) bson.D {
	filter := bson.M{
		"accommodationDoc.isActive": true,
	}

	if location.Longitude != 0 || location.Latitude != 0 {
		filter["accommodationDoc.location"] = bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{location.Longitude, location.Latitude},
				},
				"$maxDistance": location.RadiusKm * 1000,
			},
		}
	} else if location.City != "" {

		filter["accommodationDoc.address.city"] = bson.M{
			"$regex": location.City, "$options": "i",
		}
		if location.Country != "" {
			filter["accommodationDoc.address.country"] = bson.M{
				"$regex": location.Country, "$options": "i",
			}
		}

	}

	if property != "" {
		filter["accommodationDoc.type"] = property
	}

	if totalGuests > 0 {
		filter["accommodationDoc.roomType"] = bson.M{
			"$elemMatch": bson.M{
				"maxGuests": bson.M{"$gte": totalGuests},
			},
		}
	}

	return bson.D{{Key: "$match", Value: filter}}
}

func stageGroupByAccommodation(nights int) bson.D {
	return bson.D{{Key: "$group", Value: bson.M{
		"_id":              "$accommodationID",
		"availableNights":  bson.M{"$sum": 1},
		"totalPrice":       bson.M{"$sum": "$pricePerNight"},
		"accommodationDoc": bson.M{"$first": "$accommodationDoc"},
		"minAvailableRooms": bson.M{
			"$min": bson.M{
				"$subtract": []interface{}{"$totalRooms", "$reservedRooms"},
			},
		},
		"currency": bson.M{"$first": "$currency"},
	}}}
}

func stageSortAvailability(sort model.AccommodationSortOption) bson.D {
	switch sort {
	case model.SortAccommodationByRating:
		return bson.D{{Key: "$sort", Value: bson.D{
			{Key: "accommodationDoc.rating", Value: -1},
			{Key: "totalPrice", Value: 1},
		}}}
	default:
		return bson.D{{Key: "$sort", Value: bson.D{
			{Key: "totalPrice", Value: 1},
			{Key: "accommodationDoc.rating", Value: -1},
		}}}
	}
}
