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
		DepartureAirport string                   `json:"departureAirport"`
		ArrivalAirport   string                   `json:"arrivalAirport"`
		DepartureTime    time.Time                `json:"departureTime"`
		ArrivalTime      time.Time                `json:"arrivalTime"`
		Airline          string                   `json:"airline"`
		FlightNumber     string                   `json:"flightNumber"`
		CabinClass       []model.FlightCabinClass `json:"cabinClass"`
		Stops            string                   `json:"stops"`
		PlaneType        string                   `json:"planeType"`
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
		CabinClass:       req.CabinClass,
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

// update flight
func UpdateFight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PUT Allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	flightID := vars["flightId"]
	if flightID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing flight ID")
		return
	}

	var req struct {
		DepartureTime time.Time `json:"departureTime"`
		ArrivalTime   time.Time `json:"arrivalTime"`
		Stops         string    `json:"stops"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	flightCollection := database.DB.Collection("flights")
	var flight model.Flight

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = flightCollection.FindOne(ctx, bson.M{"_id": flightID}).Decode(&flight)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Flight not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding flight")
		}
		return
	}

	update := bson.M{
		"$set": bson.M{
			"departureTime": req.DepartureTime,
			"arrivalTime":   req.ArrivalTime,
			"stops":         req.Stops,
		},
	}

	result, err := flightCollection.UpdateOne(ctx, bson.M{"_id": flightID}, update)
	if err != nil {
		utils.Logger.Warn("Error while updating flight")
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update flight")
		return
	}

	if result.MatchedCount == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Flight not found")
		return
	}

	utils.Logger.Info("Updated flight successfully")
	utils.RespondWithJson(w, http.StatusOK, "Flight Updated", map[string]interface{}{"flight": update})

}

//dlete flight

//user
//get flight/flights
//search flights
//book flight
//get flight bookings
//get flight by routes
//get flight availability/seats
