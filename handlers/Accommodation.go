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
// Create accommodaion
func AddAccommodation(w http.ResponseWriter, r *http.Request) {
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
		PropertyType string   `json:"propertyType"`
		Name         string   `json:"name"`
		Address      string   `json:"address"`
		Ammenities   []string `json:"ammenities"`
		Description  string   `json:"descripton"`
		Images       []string `json:"images"`
		Location     string   `json:"location"`
		Fee          float64  `json:"fee"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	accommodationCollection := database.DB.Collection("accomodations")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	accommodation := model.Accommodation{
		ID:           primitive.NewObjectID(),
		PropertyType: req.PropertyType,
		Name:         req.Name,
		Address:      req.Address,
		Amenities:    req.Ammenities,
		Description:  req.Description,
		Images:       req.Images,
		Location:     req.Location,
		Fee:          req.Fee,
		Rating:       0,
	}

	_, err = accommodationCollection.InsertOne(ctx, accommodation)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating accommodation")
		return
	}

	utils.Logger.Info("Successfully created Accommodation")
	utils.RespondWithJson(w, http.StatusCreated, "Created acommodation", map[string]interface{}{})

}

// update accommodation
func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PATCH allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	accommodationIDStr := vars["accommodationID"]
	if accommodationIDStr == "" {
		utils.RespondWithError(w, http.StatusNotFound, "Missing flight ID")
		return
	}

	accommodaionID, err := primitive.ObjectIDFromHex(accommodationIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid flight ID")
		return
	}

	var req struct {
		Fee         string   `json:"fee"`
		Description string   `json:"description"`
		Amenities   string   `json:"amenities"`
		Images      []string `json:"imags"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	accommodationCollection := database.DB.Collection("accommodation")
	var accommodation model.Accommodation

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = accommodationCollection.FindOne(ctx, bson.M{"_id": accommodaionID}).Decode(&accommodation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Accommodation not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding accommodation")
		}
		return
	}

	update := bson.M{
		"$set": bson.M{
			"fee":         req.Fee,
			"description": req.Description,
			"amenities":   req.Amenities,
			"images":      req.Images,
		},
	}

	result, err := accommodationCollection.UpdateOne(ctx, bson.M{"_id": accommodaionID}, update)
	if err != nil {
		utils.Logger.Warn("Failed to update accommodation")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating accommodation")
		return 
	}

	if result.MatchedCount == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Accommodation not found")
		return 
	}


	utils.Logger.Info("Accommodation updated")
	utils.RespondWithJson(w, http.StatusOK, "accommodation updated successfully", map[string]interface{}{})

}

//Delete accommodation
//booking stats

//user
//get accommodations
//get accommodationID
//search accommodations
//get available rooms
//get accommodation availability
//get by location
//get reviews
//leave review
//book accommodation
