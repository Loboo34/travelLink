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

//user
//get accommodations
func GetAcommodations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET")
		return
	}

	accommodationCollection := database.DB.Collection("acommodations")

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

	utils.RespondWithJson(w, http.StatusOK, "Fetched accommodations", accommodations)
}
//get accommodationID
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

	utils.RespondWithJson(w, http.StatusOK, "Fetched accommodation", accommodation)
}
//search accommodations
//get available rooms
//get accommodation availability
//get by location
//get reviews
//leave review