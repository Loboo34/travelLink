package handlers

import (
	"encoding/json"
	"net/http"

	model "github.com/Loboo34/travel/models"
	"github.com/Loboo34/travel/service"
	"github.com/Loboo34/travel/utils"
)

type FlightBookingHandler struct {
	bookingService *service.FlightBookingService
}

func NewFlightBookingHandler(bookingService *service.FlightBookingService) *FlightBookingHandler {
	return &FlightBookingHandler{bookingService: bookingService}
}

// book flight
func (h *FlightBookingHandler) FLightBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	userID, err := utils.GetUserID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "missing user ID")
		return
	}

	var req model.FlightRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	defer r.Body.Close()

	result, err := h.bookingService.Book(r.Context(), userID, req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error booking")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, "Booking complete", result)
}

// book accommodation
type AccommodationBookingHandler struct {
	bookingService *service.AccommodationBookingService
}

func NewAccommodationBookingHandler(bookingService *service.AccommodationBookingService) *AccommodationBookingHandler {
	return &AccommodationBookingHandler{bookingService: bookingService}
}

func (h *AccommodationBookingHandler) AccommodationBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	userID, err := utils.GetUserID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	var req model.AccommodationBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	defer r.Body.Close()

	result, err := h.bookingService.Book(r.Context(), userID, req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error booking accommodations")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, "Booking complete", result)
}

// book activities
type ActivityBookingHandler struct {
	bookingService *service.ActivityBookingService
}

func NewActivityBookingHandler(bookingService *service.ActivityBookingService) *ActivityBookingHandler {
	return &ActivityBookingHandler{bookingService: bookingService}
}

func (h *ActivityBookingHandler) ActivityBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	userID, err := utils.GetUserID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	var req model.ActivityBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	defer r.Body.Close()

	result, err := h.bookingService.Book(r.Context(), userID, req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error booking activity")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, "Booking complete", result)

}

// book package
type PackageBookingHandler struct {
	bookingService *service.PackageBookingService
}

func NewPackageHandler(bookingService *service.PackageBookingService) *PackageBookingHandler {
	return &PackageBookingHandler{bookingService: bookingService}
}

func (h *PackageBookingHandler) PackageBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "only POST allowed")
		return
	}

	userID, err := utils.GetUserID()
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "missing user ID")
		return
	}

	var req model.PackageBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return 
	}

	defer r.Body.Close()

	result, err := h.bookingService.Book(r.Context(), userID, req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error booking activity")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, "Booking complete", result)


}

//get users bookings
//get booking by id
//update booking
//cancel booking
//get ooking history
//createitinerary
//get itinerary
