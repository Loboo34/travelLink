package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FlightRepo struct {
	db *mongo.Database
}

func NewFligthRepo(db *mongo.Database) *FlightRepo {
	return &FlightRepo{db: db}
}

type FlightFilter struct {
	OriginID      primitive.ObjectID
	DestinationID primitive.ObjectID
	DepartureTime time.Time
	CabinClass    model.CabinClassType
	MinSeats      int
	SortBy        model.SortOptions
	Page          int
	PageSize      int
}

func (r *FlightRepo) SearchOffers(ctx context.Context, f FlightFilter) ([]model.FlightOffer, error) {
	startOfDay := time.Date(
		f.DepartureTime.Year(), f.DepartureTime.Month(), f.DepartureTime.Day(),
		0, 0, 0, 0, time.UTC,
	)

	endOfDay := startOfDay.Add(24 * time.Hour)

	pipeline := mongo.Pipeline{
		stageMatchOffers(f.MinSeats),
		stageLookupSegments(),
		stageMatchSegments(f.OriginID, f.DestinationID, startOfDay, endOfDay, f.CabinClass, f.MinSeats),
		stageSortOffers(f.SortBy),
		bson.D{{Key: "$skip", Value: (f.Page - 1) * f.PageSize}},
		bson.D{{Key: "$limit", Value: f.PageSize}},
	}

	cursor, err := r.db.Collection("flightOffers").Aggregate(ctx, pipeline)
	if err != nil {
		utils.Logger.Warn("Error with flight aggragation")
		return nil, fmt.Errorf("Failed flight search aggragation: %w", err)
	}

	defer cursor.Close(ctx)

	var flightOffers []model.FlightOffer
	if err := cursor.All(ctx, &flightOffers); err != nil {
		utils.Logger.Warn("Failed to decode flight offers")
		return nil, fmt.Errorf("Error decoding offers: %w", err)
	}

	return flightOffers, nil
}

func stageMatchOffers(minSeats int) bson.D {
	return bson.D{{Key: "$match", Value: bson.M{
		"isActive":      true,
		"expiresAt":     bson.M{"$gt": time.Now()},
		"bookableSeats": bson.M{"$gte": minSeats},
	}}}
}
func stageLookupSegments() bson.D {
	return bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "flights",
		"localField":   "segments",
		"foreignField": "_id",
		"as":           "segmentDocs",
	}}}
}
func stageMatchSegments(originID primitive.ObjectID,
	destinationID primitive.ObjectID,
	startOfDay time.Time,
	endOfDay time.Time,
	cabinClass model.CabinClassType,
	minSeats int) bson.D {
	return bson.D{{Key: "$match", Value: bson.M{
		"segmentDocs": bson.M{"$elemMatch": bson.M{
			"originID":      originID,
			"departureTime": bson.M{"$gte": startOfDay, "$lt": endOfDay},
			"status":        bson.M{"$ne": string(model.FlightStatusCancelled)},
			"cabinClasses": bson.M{"$elemMatch": bson.M{
				"classType": cabinClass,
				"$expr": bson.M{
					"$gte": []interface{}{
						bson.M{"$subtract": []interface{}{"$totalSeats", "$reservedSeats"}},
						minSeats,
					},
				},
			}},
		}},

		"segmentDocs.destinationID": destinationID,
	}}}
}
func stageSortOffers(sort model.SortOptions) bson.D {
	switch sort {
	case model.SortByStops:
		return bson.D{{Key: "$sort", Value: bson.D{
			{Key: "stops", Value: 1},
			{Key: "priceTotal", Value: 1},
		}}}
	case model.SortByAirline:
		return bson.D{{Key: "$sort", Value: bson.D{
			{Key: "airline", Value: 1},
			{Key: "priceTotal", Value: 1},
		}}}
	default:
		return bson.D{{Key: "$sort", Value: bson.D{
			{Key: "priceTotal", Value: 1},
			{Key: "stops", Value: 1},
		}}}
	}
}

type AirportRepository struct {
	db *mongo.Database
}

func NewAirportRepository(db *mongo.Database) *AirportRepository {
	return &AirportRepository{db: db}
}

func (r *AirportRepository) FindIDByCode(ctx context.Context, code string) (primitive.ObjectID, error) {
	var airport model.Airport

	err := r.db.Collection("airports").FindOne(ctx, bson.M{
		"code": strings.ToUpper(code),
	}).Decode(&airport)

	if err == mongo.ErrNoDocuments {
		return primitive.NilObjectID, fmt.Errorf("airport with code %q not found", code)
	}
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("airport lookup failed: %w", err)
	}

	return airport.ID, nil
}
