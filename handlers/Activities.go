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

	var req struct {
		Title           string   `json:"title"`
		Price           float64  `json:"price"`
		DurationMinutes int      `json:"durationMinutes"`
		Inclusions      []string `json:"inclusions"`
		Exclusions      []string `json:"exclusions"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
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
			"title":           req.Title,
			"price":           req.Price,
			"durationMinutes": req.DurationMinutes,
			"inclusions":      req.Inclusions,
			"Exclusion":       req.Exclusions,
		},
	}

	_, err = activityCollection.UpdateOne(ctx, bson.M{"_id": activityID}, update)
	if err != nil {
		utils.Logger.Warn("Failed to update activity")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating activity")
		return
	}

	utils.Logger.Info("Updated activity")
	utils.RespondWithJson(w, http.StatusOK, "Activity updated successfully", map[string]interface{}{})

}

// delete activity
func DeleteActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "ONly DELETE allowed")
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

	activityCollection := database.DB.Collection("activities")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := activityCollection.DeleteOne(ctx, bson.M{"_id": activityID})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting activity")
		utils.Logger.Warn("Failed to delete activity ")
		return
	}

	if result.DeletedCount == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Activity not found")
	}

	utils.Logger.Info("Deleted activity")
	utils.RespondWithJson(w, http.StatusOK, "Deleted activity successfully", map[string]interface{}{})

}

// create timeslot
func TimeSlot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Post allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
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

	var req struct {
		ActivityID      primitive.ObjectID `json:"activityID"`
		StartTime       time.Time          `json:"startTime"`
		DurationMinutes int                `json:"durationMinutes"`
		AvailableSpots  int                `json:"availableSpots"`
		PricePerPerson  float64            `json:"pricePerPerson"`
		GroupSizeMax    int                `json:"groupSiveMax"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	timeSlotCollection := database.DB.Collection("activity-timeslot")

	activityCollection := database.DB.Collection("activities")
	var activity model.Activity

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = activityCollection.FindOne(ctx, bson.M{"_id": activityID}).Decode(&activity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Accommodation not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding accommodation")
		}
		return
	}

	create := model.ActivityTimeslot{
		ID:              primitive.NewObjectID(),
		ActivityID:      req.ActivityID,
		StartTime:       req.StartTime,
		DurationMinutes: req.DurationMinutes,
		AvailableSpots:  req.AvailableSpots,
		PricePerPerson:  req.PricePerPerson,
		GroupSizeMax:    req.GroupSizeMax,
	}

	_, err = timeSlotCollection.InsertOne(ctx, create)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating time slot")
		utils.Logger.Warn("Failed to create a time slot")
		return
	}

	utils.Logger.Info("Time slot created")
	utils.RespondWithJson(w, http.StatusOK, "Time slot created successfully", map[string]interface{}{})

}

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
