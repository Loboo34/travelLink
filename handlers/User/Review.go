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

func LeaveReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Post allowed")
		return
	}

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	var req struct {
		Content   string          `json:"content"`
		Rating    int             `json:"rating"`
		ReviewFor model.ReviewFor `json:"reviewFor"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid json format")
		return
	}

	vars := mux.Vars(r)

	referenceIDStr := vars["referenceID"]
	if referenceIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing refference ID")
		return
	}

	referenceID, err := primitive.ObjectIDFromHex(referenceIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid reference ID")
		return
	}

	reviewCollection := database.DB.Collection("reviews")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	create := model.Review{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		ReferenceID: referenceID,
		Content:     req.Content,
		Rating:      req.Rating,
		IsVerified:  true,
		ReviewFor:   req.ReviewFor,
		CreatedAt:   time.Now(),
	}

	_, err = reviewCollection.InsertOne(ctx, create)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error leaving review")
		utils.Logger.Warn("Failed to leave review")
		return
	}

	utils.Logger.Info("Review created")
	utils.RespondWithJson(w, http.StatusCreated, map[string]interface{}{})
}

func UpdateReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PATCH allowed")
		return
	}

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	var req struct {
		Content string `json:"content"`
		Rating  int    `json:"rating"`
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
			"rating":   req.Rating,
			"updateAt": time.Now(),
		},
	}

	_, err = reviewCollection.UpdateOne(ctx, bson.M{"_id": reviewID}, update)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating review")
		utils.Logger.Warn("Failed to update review")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, "Review updated")

}

func DeleteReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only DELETE allowed")
		return
	}

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "missing user ID")
		return
	}

	vars := mux.Vars(r)
	reviewIDStr := vars["reviewID"]
	if reviewIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing review ID")
		return
	}

	reviewID, err := primitive.ObjectIDFromHex(reviewIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid review ID")
		return
	}
	reviewCollection := database.DB.Collection("reviews")
	var review model.Review

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = reviewCollection.FindOne(ctx, bson.M{"_id": reviewID, "userID": userID}).Decode(&review)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Review not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding ")
			utils.Logger.Warn("Failed to find the review")
		}
		return
	}

	_, err = reviewCollection.DeleteOne(ctx, bson.M{"_id": reviewID})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting review")
		utils.Logger.Warn("Failed to delete review")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, "Review Deleted")

}

func GetReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	reviewCollection := database.DB.Collection("reviews")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	cursor, err := reviewCollection.Find(ctx, bson.M{})
	if err != nil {
		utils.Logger.Warn("Failed to find revies")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error finding reviews")
	}

	defer cursor.Close(ctx)

	var reviews []model.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error decoding reviews")
		utils.Logger.Warn("Failed to decode reviews")
	}

	utils.RespondWithJson(w, http.StatusOK, reviews)
}
