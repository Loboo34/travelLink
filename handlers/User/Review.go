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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		Content   string          `json:"content"`
		Note      string          `json:"note"`
		ReviewFor model.ReviewFor `json:"reviewFor"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid json format")
		return
	}

	reviewCollection := database.DB.Collection("reviews")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	create := model.Review{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Content:   req.Content,
		Note:      req.Note,
		ReviewFor: req.ReviewFor,
		CreatedAt: time.Now(),
	}

	_, err = reviewCollection.InsertOne(ctx, create)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error leaving review")
		utils.Logger.Warn("Failed to leave review")
		return
	}

	utils.Logger.Info("Review created")
	utils.RespondWithJson(w, http.StatusCreated, "Review left", map[string]interface{}{})
}

func UpdateReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PATCH allowed")
		return
	}

	userIDStr, err := utils.GetUserID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	var req struct {
		Content string `json:"content"`
		Note    string `json:"note"`
	}

	vars := mux.Vars(r)
	reviewIDStr := vars["reviewID"]
	if reviewIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing review ID")
		return
	}

	reviewID, err := primitive.ObjectIDFromHex(reviewIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid review ID")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	userCollection := database.DB.Collection("users")
	var user model.User

	reviewCollection := database.DB.Collection("reviews")
	var review model.Review

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "User not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding user")
			utils.Logger.Warn("Failed to find user")
		}
		return
	}

	err = reviewCollection.FindOne(ctx, bson.M{"_id": reviewID}).Decode(&review)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Review not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error Finding review")
		}
		return
	}

	update := bson.M{
		"$set": bson.M{
			"content":  req.Content,
			"note":     req.Note,
			"updateAt": time.Now(),
		},
	}

	_, err = reviewCollection.UpdateOne(ctx, bson.M{"_id": reviewID}, update)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating review")
		utils.Logger.Warn("Failed to update review")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, "Review updated", "")

}


