package user

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
// get all activities/activity
func GetActivities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET")
		return
	}

	activityCollection := database.DB.Collection("activities")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	cursor, err := activityCollection.Find(ctx, bson.M{})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error finding activities")
		utils.Logger.Warn("Failed to find activities")
		return
	}

	defer cursor.Close(ctx)

	var activities []model.Activity
	if err := cursor.All(ctx, &activities); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error decoding activities")
		utils.Logger.Warn("Failed to decode activities")
	}

	utils.RespondWithJson(w, http.StatusOK, "Fetched activities", activities)
}

func GetActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	vars := mux.Vars(r)
	activityIDStr := vars["activityID"]
	if activityIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing activity ID")
		return
	}

	activityID, err := primitive.ObjectIDFromHex(activityIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid activity ID")
		return
	}

	activityCollection := database.DB.Collection("activities")
	var activity model.Activity

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = activityCollection.FindOne(ctx, bson.M{"_id": activityID}).Decode(&activity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Activity not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding activity")
			utils.Logger.Warn("Failed to fetch activity")
		}
		return
	}

	utils.RespondWithJson(w, http.StatusOK, "Fetched activity", activity)
}

// search activity
type ActivityHandler struct {
	activityService *service.ActivityService
}

func NewActivityHandler(activityService *service.ActivityService) *ActivityHandler {
	return &ActivityHandler{activityService: activityService}
}

func (h *ActivityHandler) ActivitySearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	params, err := ParceAivityParams(r.URL.Query())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Search param errors")
		return
	}

	result, err := h.activityService.Search(r.Context(), params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error making search")
		utils.Logger.Warn("Failed to make search")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, "Results:", result)
}

func ParceAivityParams(q url.Values) (model.ActivitySearch, error) {
	var params model.ActivitySearch

	params.Location = model.LocationSearch{
		City:    q.Get("city"),
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

	date, err := time.Parse("2006-01-02", q.Get("date"))
	if err != nil {
		return params, errors.New("wrong date format")
	}
	params.Date = date
	adults, err := strconv.Atoi(q.Get("adults"))
	if err != nil || adults < 1 {
		return params, errors.New("adults must be a number greater than 0")
	}
	params.Participants.Adults = adults

	if c := q.Get("children"); c != "" {
		children, err := strconv.Atoi(c)
		if err != nil {
			return params, errors.New("children must be a valid number")
		}
		params.Participants.Children = children
	}

	if i := q.Get("infants"); i != "" {
		infants, err := strconv.Atoi(i)
		if err != nil {
			return params, errors.New("infants must be a valid number")
		}
		params.Participants.Infants = infants
	}

	// optional fields
	if d := q.Get("maxDurationMinutes"); d != "" {
		duration, err := strconv.Atoi(d)
		if err != nil {
			return params, errors.New("maxDurationMinutes must be a valid number")
		}
		params.MaxDurationMinutes = duration
	}

	params.Category = model.ActivityCategory(q.Get("category"))
	params.SortBy = model.ActivitySortOptions(q.Get("sortBy"))
	page, _ := strconv.Atoi(q.Get("page"))
	params.Page = page

	pageSize, _ := strconv.Atoi(q.Get("pageSize"))
	params.PageSize = pageSize

	return params, nil

}

//get activity time slots
//check availability
//get reviews
//get by location
//get by categories
