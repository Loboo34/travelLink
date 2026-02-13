package handlers

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

// admin
// add flight
func AddFlight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	var req struct {
		DepartureAirport string    `json:"departureAirport"`
		ArrivalAirport   string    `json:"arrivalAirport"`
		DepartureTime    time.Time `json:"departureTime"`
		ArrivalTime      time.Time `json:"arrivalTime"`
		Airline          string    `json:"airline"`
		FlightNumber     string    `json:"flightNumber"`
		Stops            string    `json:"stops"`
		PlaneType        string    `json:"planeType"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Json")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	flightCollection := database.DB.Collection("flights")

	flight := model.Flight{
		ID:               primitive.NewObjectID(),
		DepartureAirport: req.DepartureAirport,
		ArrivalAirport:   req.ArrivalAirport,
		DepartureTime:    req.DepartureTime,
		ArrivalTime:      req.ArrivalTime,
		Airline:          req.Airline,
		FlightNumber:     req.FlightNumber,
		Stops:            req.Stops,
		PlaneType:        req.PlaneType,
	}

	_, err = flightCollection.InsertOne(ctx, flight)
	if err != nil {
		utils.Logger.Warn("Failed to create flight")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating flight")
		return
	}

	utils.Logger.Info("Created Fligh successfully")
	utils.RespondWithJson(w, http.StatusCreated, "Flight created successfully", map[string]interface{}{
		"flightID": flight.ID.Hex(),
	})

}

//update flight
//dlete flight

//user
//get flight/flights
//search flights
//book flight
//get flight bookings
//get flight by routes
//get flight availability/seats
