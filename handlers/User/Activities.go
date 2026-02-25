package user

import (
	"context"
	"net/http"
	"time"

	"github.com/Loboo34/travel/database"
	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/utils"
	"go.mongodb.org/mongo-driver/bson"
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

//search activity
//get activity time slots
//check availability
//get reviews
//get by location
//get by categories
