package handlers_admin

import (
	"encoding/json"
	"net/http"

	"github.com/Loboo34/travel/handlers"
	"github.com/Loboo34/travel/service"
	"github.com/Loboo34/travel/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivityHandler struct {
	activityService *service.ActivityService
}

func NewActivityHandler(activityService *service.ActivityService) *ActivityHandler {
	return &ActivityHandler{activityService: activityService}
}

func (h *ActivityHandler) CreateActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Missing admin ID")
		return
	}

	var req service.ActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	result, err := h.activityService.Create(r.Context(), req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating activity")
		return
	}

	utils.Logger.Info("Created activity")
	utils.RespondWithJson(w, http.StatusCreated, result)
}

func (h *ActivityHandler) UpdateActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PUT/PATCH allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	activityIDStr := vars["activityID"]
	if activityIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing activity ID")
		return
	}

	activityID, err := primitive.ObjectIDFromHex(activityIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid activity ID")
		return
	}

	var req service.ActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	result, err := h.activityService.Update(r.Context(), activityID, req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating activity")
		return
	}

	utils.Logger.Info("Updated activity")
	utils.RespondWithJson(w, http.StatusOK, result)
}

func (h *ActivityHandler) DeleteActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only DELETE allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	activityIDStr := vars["activityID"]
	if activityIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing activity ID")
		return
	}

	activityID, err := primitive.ObjectIDFromHex(activityIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid activity ID")
		return
	}

	if err := h.activityService.Delete(r.Context(), activityID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting activity")
		return
	}

	utils.Logger.Info("Deleted activity")
	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{})
}

func (h *ActivityHandler) CreateTimeSlot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	activityIDStr := vars["activityID"]
	if activityIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing activity ID")
		return
	}

	activityID, err := primitive.ObjectIDFromHex(activityIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid activity ID")
		return
	}

	var req service.ActivityTimeslotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	result, err := h.activityService.CreateTimeSlot(r.Context(), activityID, req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating time slot")
		return
	}

	utils.Logger.Info("Created time slot")
	utils.RespondWithJson(w, http.StatusCreated, result)
}

func (h *ActivityHandler) GetActivities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET")
		return
	}

	results, err := h.activityService.GetActivities(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error decoding activities")
		utils.Logger.Warn("Failed to decode activities")
	}

	utils.RespondWithJson(w, http.StatusOK, results)
}

func (h *ActivityHandler) GetActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	vars := mux.Vars(r)
	activityIDStr := vars["activityID"]
	if activityIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing activity ID")
		return
	}

	activityID, err := primitive.ObjectIDFromHex(activityIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid activity ID")
		return
	}

	result, err := h.activityService.GetActivity(r.Context(), activityID)
	if err != nil {
		handlers.HandleServiceError(w, err, "failed getting activity")
	}

	utils.RespondWithJson(w, http.StatusOK, result)
}

func (h *ActivityHandler) GetTimeslots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET")
		return
	}

	results, err := h.activityService.GetTimeSlots(r.Context())
	if err != nil {
		handlers.HandleServiceError(w, err, "failed getting activity timeslots")
		utils.Logger.Warn("Failed to decode activities")
	}

	utils.RespondWithJson(w, http.StatusOK, results)
}

func (h *ActivityHandler) GetTimeslot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	vars := mux.Vars(r)
	activityIDStr := vars["timeslotID"]
	if activityIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing activity ID")
		return
	}

	activityID, err := primitive.ObjectIDFromHex(activityIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid activity ID")
		return
	}

	result, err := h.activityService.GetActivity(r.Context(), activityID)
	if err != nil {
		handlers.HandleServiceError(w, err, "failed getting activity timeslot")
	}

	utils.RespondWithJson(w, http.StatusOK, result)
}
