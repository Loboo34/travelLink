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

// flight offer
func FlightOffer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Ony POST allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	var req struct {
		FlightID          primitive.ObjectID `json:"flightID"`
		Provider          string             `json:"provider"`
		ProviderReference string             `json:"providerReference"`
		OneWay            bool               `json:"oneWay"`
		Price             float64            `json:"price"`
		BaggageAllowance  string             `json:"baggageAllowance"`
		LastTicketingDate time.Time          `json:"lastTicketingDate"`
		BookableSeats     int                `json:"bookableSeats"`
		ExpiresAt         time.Time          `json:"expiresAt"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	offersCollection := database.DB.Collection("flight-offers")
	//var offer model.FlightOffer

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	create := model.FlightOffer{
		ID:                primitive.NewObjectID(),
		FlightID:          req.FlightID,
		ProviderReference: req.ProviderReference,
		Provider:          req.Provider,
		OneWay:            req.OneWay,
		PriceTotal:        req.Price,
		BaggageAllowance:  req.BaggageAllowance,
		LastTicketingDate: &req.LastTicketingDate,
		BookableSeats:     req.BookableSeats,
		ExpiresAt:         &req.ExpiresAt,
	}

	_,err = offersCollection.InsertOne(ctx, create)
	if err != nil{
		utils.Logger.Warn("Failed to create offer")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating offer")
		return
	}

	utils.Logger.Info("Successfully created offer")
	utils.RespondWithJson(w, http.StatusCreated, "Created offer", map[string]interface{}{"flightID": req.FlightID})

}

// dlete flight
func DeleteFlight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Delete allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin iD")
		return
	}

	vars := mux.Vars(r)
	flightID := vars["flightID"]
	if flightID == "" {
		utils.RespondWithError(w, http.StatusNotFound, "Missing flight ID")
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

	result, err := flightCollection.DeleteOne(ctx, bson.M{"_id": flightID})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete flight")
		return
	}

	if result.DeletedCount == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Flight not found")
		return
	}

	utils.Logger.Info("Delete successful")
	utils.RespondWithJson(w, http.StatusOK, "Deleted flight successfully", map[string]interface{}{})

}

//delete offer
//update offer

//user
//get flight/flights
//search flights
//book flight
//get flight bookings
//get flight by routes
//get flight availability/seats
