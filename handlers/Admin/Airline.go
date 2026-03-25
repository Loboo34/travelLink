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

// airline
func CreateAirline(w http.ResponseWriter, r *http.Request) {
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
		Name string `json:"name"`
		Code string `json:"code"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	airlineCollection := database.DB.Collection("airlines")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	create := model.Airline{
		ID:   primitive.NewObjectID(),
		Name: req.Name,
		Code: req.Code,
	}

	_, err = airlineCollection.InsertOne(ctx, create)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating Airline")
		utils.Logger.Warn("Failed to create Airline")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, map[string]interface{}{})
	utils.Logger.Info("Airline created")
}

func UpdateAirline(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PUT allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	airlineIDStr := vars["airlineID"]
	if airlineIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing airline ID")
		return
	}

	airlineID, err := primitive.ObjectIDFromHex(airlineIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid airline ID")
		return
	}

	var req struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	airlineCollection := database.DB.Collection("airlines")
	var airline model.Airline

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = airlineCollection.FindOne(ctx, bson.M{"_id": airlineID}).Decode(&airline)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Airline not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding airline")
		}
		return
	}

	update := bson.M{
		"$set": bson.M{
			"name": req.Name,
			"code": req.Code,
		},
	}

	_, err = airlineCollection.UpdateOne(ctx, bson.M{"_id": airlineID}, update)
	if err != nil {
		utils.Logger.Warn("Failed to update airline")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating airline")
		return
	}

	utils.RespondWithJson(w, http.StatusOK,  map[string]interface{}{})

}

func DeleteAirline(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only DELETE allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	airlineIDStr := vars["airlineID"]
	if airlineIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing airline ID")
		return
	}

	airlineID, err := primitive.ObjectIDFromHex(airlineIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid airline ID")
		return
	}

	airlineCollection := database.DB.Collection("airlines")
	var airline model.Airline

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = airlineCollection.FindOne(ctx, bson.M{"_id": airlineID}).Decode(&airline)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Airline not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding airline")
		}
		return
	}

	_, err = airlineCollection.DeleteOne(ctx, bson.M{"_id": airlineID})
	if err != nil {
		utils.Logger.Warn("Failed to Delete airline")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error Delete airline")
		return
	}

	utils.RespondWithJson(w, http.StatusOK,  map[string]interface{}{})

}

func GetAirlines(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Get Allowed")
		return
	}

	airlineCollection := database.DB.Collection("airlines")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	cursor, err := airlineCollection.Find(ctx, bson.M{})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error finding airlines")
		utils.Logger.Warn("Failed to fetch airlines")
		return
	}

	defer cursor.Close(ctx)

	var airlines []model.Airline
	if err := cursor.All(ctx, &airlines); err != nil {
		utils.Logger.Warn("Failed to decode airlines")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error decodeing airlines")
		return
	}

	utils.RespondWithJson(w, http.StatusOK,  airlines)
}

func GetAirline(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	vars := mux.Vars(r)
	airlineIDStr := vars["airlineID"]
	if airlineIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing airline ID")
		return
	}

	airlineID, err := primitive.ObjectIDFromHex(airlineIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid airline ID")
		return
	}

	airlineCollection := database.DB.Collection("airlines")
	var airline model.Airline

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = airlineCollection.FindOne(ctx, bson.M{"_id": airlineID}).Decode(&airline)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Airline not found")

		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding airline")
			utils.Logger.Warn("Failed to find airline")
		}
	}

	utils.RespondWithJson(w, http.StatusOK,  map[string]interface{}{})
}

// airport
func CreateAirport(w http.ResponseWriter, r *http.Request) {
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
		Name      string  `json:"name"`
		Code      string  `json:"code"`
		City      string  `json:"city"`
		Country   string  `json:"country"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Timezone  string  `json:"timeZone"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	airportCollection := database.DB.Collection("airports")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	create := model.Airport{
		ID:        primitive.NewObjectID(),
		Name:      req.Name,
		Code:      req.Code,
		City:      req.City,
		Country:   req.Country,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Timezone:  req.Timezone,
	}

	_, err = airportCollection.InsertOne(ctx, create)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating Airport")
		utils.Logger.Warn("Failed to create Airport")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, map[string]interface{}{})
	utils.Logger.Info("Airport created")
}

func UpdateAirport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PUT allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	airportIDStr := vars["airportID"]
	if airportIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing airport ID")
		return
	}

	airportID, err := primitive.ObjectIDFromHex(airportIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid airport ID")
		return
	}

	var req struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	airportCollection := database.DB.Collection("airports")
	var airport model.Airport

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = airportCollection.FindOne(ctx, bson.M{"_id": airportID}).Decode(&airport)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Airport not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding airport")
		}
		return
	}

	update := bson.M{
		"$set": bson.M{
			"name": req.Name,
			"code": req.Code,
		},
	}

	_, err = airportCollection.UpdateOne(ctx, bson.M{"_id": airportID}, update)
	if err != nil {
		utils.Logger.Warn("Failed to update airport")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating airport")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{})

}

func DeleteAirport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only DELETE allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	airportIDStr := vars["airportID"]
	if airportIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing airport ID")
		return
	}

	airportID, err := primitive.ObjectIDFromHex(airportIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid airport ID")
		return
	}

	airportCollection := database.DB.Collection("airports")
	var airport model.Airport

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = airportCollection.FindOne(ctx, bson.M{"_id": airportID}).Decode(&airport)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Airport not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding airport")
		}
		return
	}

	_, err = airportCollection.DeleteOne(ctx, bson.M{"_id": airportID})
	if err != nil {
		utils.Logger.Warn("Failed to Delete airport")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error Delete airport")
		return
	}

	utils.RespondWithJson(w, http.StatusOK,  map[string]interface{}{})

}

func GetAirports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Get Allowed")
		return
	}

	airportCollection := database.DB.Collection("airports")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	cursor, err := airportCollection.Find(ctx, bson.M{})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error finding airports")
		utils.Logger.Warn("Failed to fetch airports")
		return
	}

	defer cursor.Close(ctx)

	var airports []model.Airport
	if err := cursor.All(ctx, &airports); err != nil {
		utils.Logger.Warn("Failed to decode airports")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error decodeing airports")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, airports)
}

func GetAirport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	vars := mux.Vars(r)
	airportIDStr := vars["airportID"]
	if airportIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing airport ID")
		return
	}

	airportID, err := primitive.ObjectIDFromHex(airportIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid airport ID")
		return
	}

	airportCollection := database.DB.Collection("airlines")
	var airport model.Airline

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = airportCollection.FindOne(ctx, bson.M{"_id": airportID}).Decode(&airport)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "iArport not found")

		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding airport")
			utils.Logger.Warn("Failed to find airport")
		}
	}

	utils.RespondWithJson(w, http.StatusOK,  map[string]interface{}{})
}

// routes
func CreateRoute(w http.ResponseWriter, r *http.Request) {
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
		OriginAirportID      primitive.ObjectID `json:"originAirportID"`
		DestinationAirportID primitive.ObjectID `json:"destinationAirportID"`
		EstimatedDurationMin int                `json:"estimatedDurationMin"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	routesCollection := database.DB.Collection("routes")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	create := model.Route{
		ID:                   primitive.NewObjectID(),
		OriginAirportID:      req.OriginAirportID,
		DestinationAirportID: req.DestinationAirportID,
		EstimatedDurationMin: req.EstimatedDurationMin,
	}

	_, err = routesCollection.InsertOne(ctx, create)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating Route")
		utils.Logger.Warn("Failed to create Route")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated,  map[string]interface{}{})
	utils.Logger.Info("Route created")
}

func UpdateRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PUT allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	routeIDStr := vars["routeID"]
	if routeIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing route ID")
		return
	}

	routeID, err := primitive.ObjectIDFromHex(routeIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid route ID")
		return
	}

	var req struct {
		OriginAirportID      primitive.ObjectID `json:"originAirportID"`
		DestinationAirportID primitive.ObjectID `json:"destinationAirportID"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	routeCollection := database.DB.Collection("routes")
	var route model.Route

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = routeCollection.FindOne(ctx, bson.M{"_id": routeID}).Decode(&route)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Route not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding route")
		}
		return
	}

	update := bson.M{
		"$set": bson.M{
			"originAirportID":      req.OriginAirportID,
			"destinationAirportID": req.DestinationAirportID,
		},
	}

	_, err = routeCollection.UpdateOne(ctx, bson.M{"_id": routeID}, update)
	if err != nil {
		utils.Logger.Warn("Failed to update route")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating route")
		return
	}

	utils.RespondWithJson(w, http.StatusOK,  map[string]interface{}{})

}

func DeleteRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only DELETE allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	routeIDStr := vars["routeID"]
	if routeIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing route ID")
		return
	}

	routeID, err := primitive.ObjectIDFromHex(routeIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid route ID")
		return
	}

	routeCollection := database.DB.Collection("routes")
	var route model.Route

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = routeCollection.FindOne(ctx, bson.M{"_id": routeID}).Decode(&route)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Route not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding route")
		}
		return
	}

	_, err = routeCollection.DeleteOne(ctx, bson.M{"_id": routeID})
	if err != nil {
		utils.Logger.Warn("Failed to Delete route")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error Delete route")
		return
	}

	utils.RespondWithJson(w, http.StatusOK,  map[string]interface{}{})

}

func GetRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Get Allowed")
		return
	}

	routeCollection := database.DB.Collection("routes")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	cursor, err := routeCollection.Find(ctx, bson.M{})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error finding routes")
		utils.Logger.Warn("Failed to fetch routes")
		return
	}

	defer cursor.Close(ctx)

	var routes []model.Route
	if err := cursor.All(ctx, &routes); err != nil {
		utils.Logger.Warn("Failed to decode routes")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error decodeing routes")
		return
	}

	utils.RespondWithJson(w, http.StatusOK,  routes)
}

func GetRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	vars := mux.Vars(r)
	routeIDStr := vars["routeID"]
	if routeIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing route ID")
		return
	}

	routeID, err := primitive.ObjectIDFromHex(routeIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid route ID")
		return
	}

	routeCollection := database.DB.Collection("routes")
	var route model.Route

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = routeCollection.FindOne(ctx, bson.M{"_id": routeID}).Decode(&route)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Route not found")

		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding route")
			utils.Logger.Warn("Failed to find route")
		}
	}

	utils.RespondWithJson(w, http.StatusOK,  map[string]interface{}{})
}
