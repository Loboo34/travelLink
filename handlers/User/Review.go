package user

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

func LeaveReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Post allowed")
		return
	}

	userIDStr, err := utils.GetUserID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil{
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		Content   string `json:"content"`
		Note      string `json:"note"`
		ReviewFor string `json:"reviewFor"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil{
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid json format")
		return 
	}

	reviewCollection := database.DB.Collection("reviews")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	create := model.Review{
		ID: primitive.NewObjectID(),
		UserID: userID,
		Content: req.Content,
		Note: req.Note,
		ReviewFor: req.ReviewFor,
		CreatedAt: time.Now(),
	}

	_,err = reviewCollection.InsertOne(ctx, create)
	if err != nil{
		utils.RespondWithError(w, http.StatusInternalServerError, "Error leaving review")
		utils.Logger.Warn("Failed to leave review")
		return 
	}
}
