package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Loboo34/travel/database"
	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/service"
	"github.com/Loboo34/travel/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// user
// get accommodations
func GetAcommodations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET")
		return
	}

	accommodationCollection := database.DB.Collection("accommodations")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	cursor, err := accommodationCollection.Find(ctx, bson.M{})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error finding accommodation")
		utils.Logger.Warn("Failed to find accommodation")
		return
	}

	defer cursor.Close(ctx)

	var accommodations []model.Activity
	if err := cursor.All(ctx, &accommodations); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error decoding accommodations")
		utils.Logger.Warn("Failed to decode accommodations")
	}

	utils.RespondWithJson(w, http.StatusOK,  accommodations)
}

// get accommodationID
func GetAccommodation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	vars := mux.Vars(r)
	accommodationIDStr := vars["accommodationID"]
	if accommodationIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing accommodation ID")
		return
	}

	accommodationID, err := primitive.ObjectIDFromHex(accommodationIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid accommodation ID")
		return
	}

	accommodationCollection := database.DB.Collection("accommodations")
	var accommodation model.Accommodation

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = accommodationCollection.FindOne(ctx, bson.M{"_id": accommodationID}).Decode(&accommodation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Accommodation not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding accommodation")
			utils.Logger.Warn("Failed to fetch accommodation")
		}
		return
	}

	utils.RespondWithJson(w, http.StatusOK,  accommodation)
}

// search accommodations
type AccommodationHandler struct {
	accommodationService *service.AccommodationSearchService
}

func NewAccommodationHandler(accommodationService *service.AccommodationSearchService) *AccommodationHandler {
	return &AccommodationHandler{accommodationService: accommodationService}
}

func (h *AccommodationHandler) AccommodationSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	params, err := parseAccomSearchParams(r.URL.Query())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Search param errors")
		return
	}

	results, err := h.accommodationService.Search(r.Context(), params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error making search")
		utils.Logger.Warn("Failed to make search")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, results)

}

func parseAccomSearchParams(q url.Values) (model.AccommodationSearch, error) {
	var params model.AccommodationSearch

	params.Location = model.LocationSearch{
		City: q.Get("city"),
		Country: q.Get("country"),
	}

	 if lat := q.Get("latitude"); lat != "" {
        parsed, err := strconv.ParseFloat(lat, 64)
        if err != nil {
            return params, errors.New("latitude must be a valid number")
        }
        params.Location.Latitude = parsed
    }
    if lng := q.Get("longitude"); lng != "" {
        parsed, err := strconv.ParseFloat(lng, 64)
        if err != nil {
            return params, errors.New("longitude must be a valid number")
        }
        params.Location.Longitude = parsed
    }
    if radius := q.Get("radiusKm"); radius != "" {
        parsed, err := strconv.ParseFloat(radius, 64)
        if err != nil {
            return params, errors.New("radiusKm must be a valid number")
        }
        params.Location.RadiusKm = parsed
    }

	checkIn, err := time.Parse("2006-01-02", q.Get("checkInDate"))
	if err != nil {
		return params, errors.New("wrong date format")
	}
	params.ChecKInDate = checkIn

	checkOut, err := time.Parse("2006-01-02", q.Get("checkOutDate"))
	if err != nil {
		return params, errors.New("wrong date format")
	}
	params.CheckOutDate = checkOut

	 adults, err := strconv.Atoi(q.Get("adults"))
    if err != nil || adults < 1 {
        return params, errors.New("adults must be a number greater than 0")
    }
    params.Guests.Adults = adults

    if c := q.Get("children"); c != "" {
        children, err := strconv.Atoi(c)
        if err != nil {
            return params, errors.New("children must be a valid number")
        }
        params.Guests.Children = children
    }

    if i := q.Get("infants"); i != "" {
        infants, err := strconv.Atoi(i)
        if err != nil {
            return params, errors.New("infants must be a valid number")
        }
        params.Guests.Infants = infants
    }

	  if r := q.Get("totalRooms"); r != "" {
        rooms, err := strconv.Atoi(r)
        if err != nil {
            return params, errors.New("totalRooms must be a valid number")
        }
        params.TotalRooms = rooms
    }

	
	
	params.PropertyType = model.PropertyType(q.Get("propertyType"))
	params.SortBy = model.AccommodationSortOption(q.Get("sortBy"))
	page, _ := strconv.Atoi(q.Get("page"))
	params.Page = page

	pageSize, _ := strconv.Atoi(q.Get("pageSize"))
	params.PageSize = pageSize

	return params, nil

}

//get available rooms
//get accommodation availability
//get by location
//get reviews
//leave review
