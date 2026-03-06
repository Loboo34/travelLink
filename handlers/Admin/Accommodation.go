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
		HostID       *primitive.ObjectID `json:"hostID"`
		PropertyType model.PropertyType  `json:"propertyType"`
		Name         string              `json:"name"`
		Address      model.Address       `json:"address"`
		Amenities    []model.Amenity     `json:"amenities"`
		Description  string              `json:"description"`
		Images       []string            `json:"images"`
		Location     model.GeoLocation   `json:"location"`
		RoomType     []model.RoomType    `json:"roomType"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	accommodationCollection := database.DB.Collection("accommodations")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	accommodation := model.Accommodation{
		ID:           primitive.NewObjectID(),
		HostID:       req.HostID,
		PropertyType: req.PropertyType,
		Name:         req.Name,
		Address:      req.Address,
		Description:  req.Description,
		Amenities:    req.Amenities,
		Images:       req.Images,
		Location:     req.Location,
		RoomType:     req.RoomType,
		Rating:       0,
		ReviewCount:  0,
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
		utils.RespondWithError(w, http.StatusNotFound, "Missing accommodation ID")
		return
	}

	accommodaionID, err := primitive.ObjectIDFromHex(accommodationIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid flight ID")
		return
	}

	var req struct {
		Name        float64  `json:"name"`
		Description string   `json:"description"`
		Amenities   []string `json:"amenities"`
		Images      []string `json:"images"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	accommodationCollection := database.DB.Collection("accommodations")
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
			"name":        req.Name,
			"description": req.Description,
			"amenities":   req.Amenities,
			"images":      req.Images,
		},
	}

	_, err = accommodationCollection.UpdateOne(ctx, bson.M{"_id": accommodaionID}, update)
	if err != nil {
		utils.Logger.Warn("Failed to update accommodation")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating accommodation")
		return
	}

	// if result.MatchedCount == 0 {
	// 	utils.RespondWithError(w, http.StatusNotFound, "Accommodation not found")
	// 	return
	// }

	utils.Logger.Info("Accommodation updated")
	utils.RespondWithJson(w, http.StatusOK, "accommodation updated successfully", map[string]interface{}{})

}

func Availability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST allowed")
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
		utils.RespondWithError(w, http.StatusNotFound, "Missing accommodation ID")
		return
	}

	accommodaionID, err := primitive.ObjectIDFromHex(accommodationIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid flight ID")
		return
	}

	var req struct {
		RoomTypeID      primitive.ObjectID `json:"roomType"`
		AccommodationID string             `json:"accommodationID"`
		Date            time.Time          `json:"date"`
		TotalRooms      int                `json:"totalRooms"`
		ReservedRooms   int                `json:"reserveRooms"`
		PricePerNight   int64              `json:"pricePerNight"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	availabilityCollection := database.DB.Collection("accommodations-availability")
	var accommodation model.Accommodation

	accommodationCollection := database.DB.Collection("accommodations")

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

	create := model.AccommodationAvailability{
		ID:              primitive.NewObjectID(),
		AccommodationID: accommodaionID,
		RoomTypeID:      req.RoomTypeID,
		Date:            req.Date,
		TotalRooms:      req.TotalRooms,
		ReservedRooms:   0,
		PricePerNight:   req.PricePerNight,
		IsActive:        true,
	}

	_, err = availabilityCollection.InsertOne(ctx, create)
	if err != nil {
		utils.Logger.Warn("Failed to update accommodation")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating accommodation")
		return
	}

	utils.Logger.Info("Accommodation updated")
	utils.RespondWithJson(w, http.StatusOK, "accommodation updated successfully", map[string]interface{}{})

}

// Delete accommodation
func DeleteAccommodation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only DELETE allowed")
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
		utils.RespondWithError(w, http.StatusBadRequest, "Missing accommodation ID")
		return
	}

	accommodationID, err := primitive.ObjectIDFromHex(accommodationIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid accommodation ID")
		return
	}

	accommodationCollection := database.DB.Collection("accommodations")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := accommodationCollection.DeleteOne(ctx, bson.M{"_id": accommodationID})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting accommodation")
		utils.Logger.Warn("Failed to delete accommodation")
		return
	}

	if result.DeletedCount == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Accommodation not found")
		return
	}

	utils.Logger.Info("Accommodation deleted successfully")
	utils.RespondWithJson(w, http.StatusOK, "Accommodation deleted", map[string]interface{}{})
}

//booking stats
