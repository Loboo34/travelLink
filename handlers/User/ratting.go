package user

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Loboo34/travel/database"
	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Rate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	userIDStr, err := utils.GetUserID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missig user ID")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return 
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing ID")
		return
	}

	var req struct{
		Rating int64 `json:"rating"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return 
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ratingCollection := database.DB.Collection("rating")

	rate := model.Rate{
		RatingFor: model.ReviewFor(id),
		Rating: req.Rating,
		UserID: userID,
	}

	_,err = ratingCollection.InsertOne(ctx, rate)
	if err != nil{
		utils.Logger.Warn("Failed to rate ...")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error rating ")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, "", "")
}
