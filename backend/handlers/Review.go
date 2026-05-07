package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Loboo34/travel/service"
	"github.com/Loboo34/travel/utils"
	model "github.com/Loboo34/travel/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}


func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "missing user ID")
		return
	}

	var req service.ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	result, err := h.reviewService.Create(r.Context(), userID, req)
	if err != nil {
		HandleServiceError(w, err, "failed to create review")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, result)
}


func (h *ReviewHandler) GetByTarget(w http.ResponseWriter, r *http.Request) {
	refIDStr := r.URL.Query().Get("referenceID")
	reviewFor := r.URL.Query().Get("reviewFor")

	if refIDStr == "" || reviewFor == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "referenceID and reviewFor are required")
		return
	}

	refID, err := primitive.ObjectIDFromHex(refIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid referenceID")
		return
	}

	reviews, err := h.reviewService.GetByTarget(r.Context(), refID, model.ReviewFor(reviewFor))
	if err != nil {
		HandleServiceError(w, err, "failed to fetch reviews")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, reviews)
}


func (h *ReviewHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "missing user ID")
		return
	}

	reviewID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid review ID")
		return
	}

	var req service.UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	result, err := h.reviewService.Update(r.Context(), reviewID, userID, req)
	if err != nil {
		HandleServiceError(w, err, "failed to update review")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, result)
}


func (h *ReviewHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "missing user ID")
		return
	}

	reviewID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid review ID")
		return
	}

	// admins can delete any review
	_, isAdmin := utils.GetAdminID(r.Context())

	if err := h.reviewService.Delete(r.Context(), reviewID, userID, isAdmin == nil); err != nil {
		HandleServiceError(w, err, "failed to delete review")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]string{"message": "review deleted"})
}