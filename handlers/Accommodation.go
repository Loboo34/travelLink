package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Loboo34/travel/database"
	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
        ID : primitive.NewObjectID(),
        PropertyType: req.PropertyType,
        Name: req.Name,
        Address: req.Address,
        Amenities: req.Ammenities,
        Description: req.Description,
        Images: req.Images,
        Location: req.Location,
        Fee: req.Fee,
        Rating: 0,
    }

    _, err = accommodationCollection.InsertOne(ctx, accommodation)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "Error creating accommodation")
        return
    }

utils.Logger.Info("Successfully created Accommodation")
utils.RespondWithJson(w, http.StatusCreated, "Created acommodation", map[string]interface{}{})

}

//update accommodation
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
