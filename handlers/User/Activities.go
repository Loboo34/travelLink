package user

import (
	"context"
	"net/http"
	"time"

	"github.com/Loboo34/travel/database"
	model "github.com/Loboo34/travel/models"
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

//search activity
//get activity time slots
//check availability
//get reviews
//get by location
//get by categories
