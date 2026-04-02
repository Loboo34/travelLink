package handlers_admin

import (
	"encoding/json"
	"net/http"

	"github.com/Loboo34/travel/handlers"
	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/service"
	"github.com/Loboo34/travel/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type FlightHandler struct {
	flightService *service.FlightService
}

func NewFlightHandler(flightService *service.FlightService) *FlightHandler {
	return &FlightHandler{flightService: flightService}
}

// add flight
func (h *FlightHandler) AddFlight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	adminID, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.Logger.Error("Failed to get admin ID", zap.Error(err))
		utils.RespondWithError(w, http.StatusUnauthorized, "Missing admin ID")
		return
	}

	utils.Logger.Info("Admin creating flight", zap.String("adminID", adminID.Hex()))

	var req service.FlightRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.flightService.AddFlight(r.Context(), req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error creating flight")
		return
	}

	utils.Logger.Info("Created Fligh successfully")
	utils.RespondWithJson(w, http.StatusCreated, result)

}

// update flight
func (h *FlightHandler) UpdateFight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PUT Allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	flightIDStr := vars["flightID"]
	if flightIDStr == "" {
		utils.RespondWithError(w, http.StatusNotFound, "Missing flight ID")
		return
	}

	flightID, err := primitive.ObjectIDFromHex(flightIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid flight ID")
		return
	}

	var req service.UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	result, err := h.flightService.Update(r.Context(), flightID, req)
	if err != nil {
		utils.Logger.Warn("Error while updating flight")
		// utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update flight")
		handlers.HandleServiceError(w, err, "error updating flight")
		return
	}

	utils.Logger.Info("Updated flight successfully")
	utils.RespondWithJson(w, http.StatusOK, result)

}

// flight status
func (h *FlightHandler) UpdateFlightStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PATCH allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing andmin ID")
		return
	}
	vars := mux.Vars(r)
	flightIDStr := vars["flightID"]
	if flightIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing flight ID")
		return
	}

	flightID, err := primitive.ObjectIDFromHex(flightIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid flight ID")
		return
	}
	var req struct {
		Status model.FlightStatus `json:"status"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	err = h.flightService.UpdateStatus(r.Context(), flightID, &req.Status)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating flight status")
		utils.Logger.Warn("Failed to update flight status")
		return
	}
	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{"status": req.Status})
	utils.Logger.Info("Status updated successfully")
}

// delete flight
func (h *FlightHandler) DeleteFlight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Delete allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin iD")
		return
	}

	vars := mux.Vars(r)
	flightIDStr := vars["flightID"]
	if flightIDStr == "" {
		utils.RespondWithError(w, http.StatusNotFound, "Missing flight ID")
		return
	}

	flightID, err := primitive.ObjectIDFromHex(flightIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid flight ID")
		return
	}

	if err := h.flightService.DeleteFlight(r.Context(), flightID); err != nil {
		// utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete flight")
		handlers.HandleServiceError(w, err, "failed to delete flight")
		return
	}

	utils.Logger.Info("Delete successful")
	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{})

}

// flight offer
func (h *FlightHandler) FlightOffer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Ony POST allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	var req service.Offer

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	resilt, err := h.flightService.CreateOffer(r.Context(), req)
	if err != nil {
		utils.Logger.Warn("Failed to create offer")
		handlers.HandleServiceError(w, err, "failed to create offer")
		return
	}

	utils.Logger.Info("Successfully created offer")
	utils.RespondWithJson(w, http.StatusCreated, resilt)

}

// update offer
func (h *FlightHandler) UpdateOffer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PUT allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	offerIDStr := vars["offerID"]
	if offerIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing offer ID")
		return
	}

	offerID, err := primitive.ObjectIDFromHex(offerIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid offer ID")
		return
	}

	var req service.OfferUpdate

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	result, err := h.flightService.UpdateOffer(r.Context(), offerID, req)
	if err != nil {
		utils.Logger.Warn("Failed to update offer")
		// utils.RespondWithError(w, http.StatusInternalServerError, "Error updating offer")
		handlers.HandleServiceError(w, err, "failed to update offer")
		return
	}

	utils.Logger.Info("Offer updated Successfully")
	utils.RespondWithJson(w, http.StatusOK, result)
}

// offer status
func (h *FlightHandler) IsActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PATCH allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing Admin ID")
		return
	}

	vars := mux.Vars(r)
	offerIDStr := vars["offerID"]
	if offerIDStr == "" {
		utils.RespondWithError(w, http.StatusNotFound, "Missing offer ID")
		return
	}

	offerID, err := primitive.ObjectIDFromHex(offerIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid offer ID")
		return
	}

	var req struct {
		IsActive bool `json:"isActive"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if err := h.flightService.IsActive(r.Context(), offerID, req.IsActive); err != nil {
		// utils.RespondWithError(w, http.StatusInternalServerError, "Error updating offer status")
		handlers.HandleServiceError(w, err, "failed to change offer status")
		utils.Logger.Warn("Failed to update offer status")
		return
	}

	utils.Logger.Info("Update successful")
	utils.RespondWithJson(w, http.StatusOK, map[string]any{
		"isActive": req.IsActive,
	})
}

// delete offer
func (h *FlightHandler) DeleteOffer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only DELETE allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing Admin ID")
		return
	}

	vars := mux.Vars(r)
	offerIDStr := vars["offerID"]
	if offerIDStr == "" {
		utils.RespondWithError(w, http.StatusNotFound, "Missing flight ID")
		return
	}

	offerID, err := primitive.ObjectIDFromHex(offerIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid flight ID")
		return
	}

	if err := h.flightService.DeleteOffer(r.Context(), offerID); err != nil {
		utils.Logger.Warn("Failed to delete offer")
		// utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting offer")
		handlers.HandleServiceError(w, err, "failed to delete offer")
		return
	}

	utils.Logger.Info("Offer deleted successfully")
	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{})
}

func (h *FlightHandler) GetFlights(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	flights, err := h.flightService.GetFlights(r.Context())
	if err != nil{
		handlers.HandleServiceError(w, err, "fetching flights")
		return
	}

	utils.RespondWithJson(w, http.StatusOK,  flights)
}


func (h *FlightHandler)  GetFlight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	vars := mux.Vars(r)
	flightIDStr := vars["flightID"]
	if flightIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing flight ID")
		return
	}

	flightID, err := primitive.ObjectIDFromHex(flightIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid flight ID")
		return
	}

	result, err := h.flightService.GetFlight(r.Context(), flightID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error finding flight")
		utils.Logger.Warn("Failed to find flight")
		return
	}

	utils.Logger.Info("Fetched flight")
	utils.RespondWithJson(w, http.StatusOK, result)
}

func (h *FlightHandler) GetOffers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	offers, err := h.flightService.GetOffers(r.Context())
	if err != nil{
		handlers.HandleServiceError(w, err, "fetching offers")
		return
	}

	utils.RespondWithJson(w, http.StatusOK,  offers)
}

func (h *FlightHandler)  GetOffer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	vars := mux.Vars(r)
	flightIDStr := vars["offerID"]
	if flightIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing offer ID")
		return
	}

	flightID, err := primitive.ObjectIDFromHex(flightIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid flight ID")
		return
	}

	result, err := h.flightService.GetOffer(r.Context(), flightID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error finding offer")
		utils.Logger.Warn("Failed to find offer")
		return
	}

	utils.Logger.Info("Fetched flight")
	utils.RespondWithJson(w, http.StatusOK, result)
}