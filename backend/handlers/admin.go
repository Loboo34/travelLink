package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Loboo34/travel/auth"
	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/repository"
	"github.com/Loboo34/travel/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type AdminHandler struct {
	userRepo              *auth.UserRepo
	flightBookingRepo     *repository.FlightBookingRepo
	accommodationBookingRepo *repository.AccommodationBookingRepo
	activityBookingRepo   *repository.ActivityBookingRepo
	packageBookingRepo    *repository.PackageBookingRepo
}

func NewAdminHandler(
	userRepo *auth.UserRepo,
	flightBookingRepo *repository.FlightBookingRepo,
	accommodationBookingRepo *repository.AccommodationBookingRepo,
	activityBookingRepo *repository.ActivityBookingRepo,
	packageBookingRepo *repository.PackageBookingRepo,
) *AdminHandler {
	return &AdminHandler{
		userRepo:              userRepo,
		flightBookingRepo:     flightBookingRepo,
		accommodationBookingRepo: accommodationBookingRepo,
		activityBookingRepo:   activityBookingRepo,
		packageBookingRepo:    packageBookingRepo,
	}
}

// GetUsers returns all users
func (h *AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	users, err := h.userRepo.GetAllUsers(r.Context())
	if err != nil {
		utils.Logger.Error("Failed to fetch users", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, users)
}

// UpdateUserRole updates a user's role
func (h *AdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PUT allowed")
		return
	}

	vars := mux.Vars(r)
	userIDStr := vars["userId"]
	if userIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate role
	if req.Role != "admin" && req.Role != "user" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid role")
		return
	}

	// Get user first
	user, err := h.userRepo.GetUser(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Update user role in database
	if err := h.userRepo.UpdateUserRole(r.Context(), userID, model.UserRole(req.Role)); err != nil {
		utils.Logger.Error("Failed to update user role", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update user role")
		return
	}

	// Return updated user
	user.Role = model.UserRole(req.Role)
	utils.RespondWithJson(w, http.StatusOK, user)
}

// ToggleUserStatus toggles user active/inactive status
func (h *AdminHandler) ToggleUserStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PUT allowed")
		return
	}

	vars := mux.Vars(r)
	userIDStr := vars["userId"]
	if userIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userRepo.GetUser(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Toggle and update status
	newStatus := !user.IsActive
	if err := h.userRepo.UpdateUserStatus(r.Context(), userID, newStatus); err != nil {
		utils.Logger.Error("Failed to update user status", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update user status")
		return
	}

	user.IsActive = newStatus
	utils.RespondWithJson(w, http.StatusOK, user)
}

// GetBookings returns all bookings across all types
func (h *AdminHandler) GetBookings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	// Fetch all bookings from all repositories
	flightBookings, _ := h.flightBookingRepo.Bookings(r.Context())
	accommodationBookings, _ := h.accommodationBookingRepo.Bookings(r.Context())
	activityBookings, _ := h.activityBookingRepo.Bookings(r.Context())

	// Convert to unified booking format
	bookings := []map[string]interface{}{}

	for _, fb := range flightBookings {
		bookings = append(bookings, map[string]interface{}{
			"id":        fb.ID.Hex(),
			"type":      "flight",
			"status":    string(fb.Status),
			"createdAt": fb.CreatedAt,
			"totalPrice": fb.AmountPaid,
			"currency":  fb.Payment.Currency,
		})
	}

	for _, ab := range accommodationBookings {
		bookings = append(bookings, map[string]interface{}{
			"id":        ab.ID.Hex(),
			"type":      "accommodation",
			"status":    string(ab.Status),
			"createdAt": ab.CreatedAt,
			"totalPrice": ab.AmountPaid,
			"currency":  ab.Payment.Currency,
		})
	}

	for _, acb := range activityBookings {
		bookings = append(bookings, map[string]interface{}{
			"id":        acb.ID.Hex(),
			"type":      "activity",
			"status":    string(acb.Status),
			"createdAt": acb.CreatedAt,
			"totalPrice": acb.AmountPaid,
			"currency":  acb.Payment.Currency,
		})
	}

	utils.RespondWithJson(w, http.StatusOK, bookings)
}

// UpdateBookingStatus updates a booking's status
func (h *AdminHandler) UpdateBookingStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PUT allowed")
		return
	}

	vars := mux.Vars(r)
	bookingIDStr := vars["bookingId"]
	if bookingIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing booking ID")
		return
	}

	bookingID, err := primitive.ObjectIDFromHex(bookingIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid booking ID")
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"confirmed": true,
		"pending":   true,
		"cancelled": true,
		"completed": true,
		"failed":    true,
	}
	if !validStatuses[req.Status] {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid status")
		return
	}

	bookingStatus := model.BookingStatus(req.Status)

	// Try to update in each repository (order matters for performance)
	// In a real app, you'd store booking type metadata to know which repo to query
	err = h.flightBookingRepo.UpdateBooking(r.Context(), bookingID, bookingStatus, nil)
	if err == nil {
		utils.Logger.Info("Updated flight booking status", zap.String("bookingID", bookingID.Hex()))
		utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{
			"id":     bookingID.Hex(),
			"status": req.Status,
		})
		return
	}

	err = h.accommodationBookingRepo.UpdateBooking(r.Context(), bookingID, bookingStatus, nil)
	if err == nil {
		utils.Logger.Info("Updated accommodation booking status", zap.String("bookingID", bookingID.Hex()))
		utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{
			"id":     bookingID.Hex(),
			"status": req.Status,
		})
		return
	}

	err = h.activityBookingRepo.UpdateBooking(r.Context(), bookingID, bookingStatus, nil)
	if err == nil {
		utils.Logger.Info("Updated activity booking status", zap.String("bookingID", bookingID.Hex()))
		utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{
			"id":     bookingID.Hex(),
			"status": req.Status,
		})
		return
	}

	utils.RespondWithError(w, http.StatusNotFound, "Booking not found")
}

// GetStats returns admin dashboard statistics
func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	// This would require repository methods to count entities
	// For now, returning empty stats
	stats := map[string]interface{}{
		"totalFlights":       0,
		"totalAccommodations": 0,
		"totalActivities":    0,
		"totalPackages":      0,
		"totalOffers":        0,
		"totalBookings":      0,
		"totalUsers":         0,
		"revenue":           0,
	}

	utils.RespondWithJson(w, http.StatusOK, stats)
}
