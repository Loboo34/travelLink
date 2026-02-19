package handlers

import (
	"context"
	"encoding/json"
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

// admin
// Create activity
func CreateActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	var req struct {
		Title           string   `json:"title"`
		Price           float64  `json:"price"`
		Description     string   `json:"description"`
		Location        string   `json:"location"`
		Categories      []string `json:"categories"`
		DurationMinutes int      `json:"durationMinutes"`
		Inclusions      []string `json:"inclusions"`
		Exclusions      []string `json:"exclusions"`
		Images          []string `json:"images"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	activityCollection := database.DB.Collection("activities")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	create := model.Activity{
		ID:              primitive.NewObjectID(),
		Title:           req.Title,
		Price:           req.Price,
		Description:     req.Description,
		Location:        req.Location,
		Categories:      req.Categories,
		DurationMinutes: req.DurationMinutes,
		Inclusions:      req.Inclusions,
		Exclusions:      req.Exclusions,
		Images:          req.Images,
	}

	_, err = activityCollection.InsertOne(ctx, create)
	if err != nil {
		utils.Logger.Warn("Failed to create activity")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating activity")
	}

	utils.Logger.Info("Created activity")
	utils.RespondWithJson(w, http.StatusCreated, "Activity created successfully", map[string]interface{}{})
}

// Update activity
func UpdateActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PATCH allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing Admin ID")
		return
	}

	vars := mux.Vars(r)
	activityIDStr := vars["activityID"]
	if activityIDStr == "" {
		utils.RespondWithError(w, http.StatusNotFound, "Missing flight ID")
		return
	}

	activityID, err := primitive.ObjectIDFromHex(activityIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid flight ID")
		return
	}

    var req struct{
        Title           string   `json:"title"`
		Price           float64  `json:"price"`
        DurationMinutes int      `json:"durationMinutes"`
		Inclusions      []string `json:"inclusions"`
		Exclusions      []string `json:"exclusions"`
    }

	activityCollection := database.DB.Collection("activities")
	var activity model.Activity

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = activityCollection.FindOne(ctx, bson.M{"_id": activityID}).Decode(&activity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "activity not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding")
		}
		return
	}

	update := bson.M{
		"$set": bson.M{
            "title": req.Title,
            "price": req.Price,
            "durationMinutes": req.DurationMinutes,
            "inclusions": req.Inclusions,
            "Exclusion": req.Exclusions,
        },
	}

	_, err = activityCollection.UpdateOne(ctx, bson.M{"_id": activityID}, update)
	if err != nil {
		utils.Logger.Warn("Failed to update activity")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating activity")
		return
	}


    utils.Logger.Info("Updated activity")
    utils.RespondWithJson(w, http.StatusOK,  "Activity updated successfully", map[string]interface{}{})

}

//delete activity
//create timeslot
//booking stats

//user
//get all activities/activity
//search activity
//get activity time slots
//check availability
//get reviews
//get by location
//get by categories
//book activity
