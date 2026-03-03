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

func CreateAirline(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return 
	}

	_,err := utils.GetAdminID()
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
		ID: primitive.NewObjectID(),
		Name : req.Name,
		Code: req.Code,
	}

	_,err = airlineCollection.InsertOne(ctx, create )
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating Airline")
		utils.Logger.Warn("Failed to create Airline")
		return 
	}

	utils.RespondWithJson(w, http.StatusCreated, "Created airline", map[string]interface{}{})
	utils.Logger.Info("Airline created")
}

func UpdateAirline(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPut{
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
	if airlineIDStr == ""{
		utils.RespondWithError(w, http.StatusBadRequest, "Missing airline ID")
		return 
	}

	airlineID, err := primitive.ObjectIDFromHex(airlineIDStr)
	if err != nil{
		utils.RespondWithError(w, http.StatusBadRequest, "invalid airline ID")
		return 
	}

	var req struct{
		Name string `json:"name"`
		Code string `json:"code"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil{
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return 
	}

	airlineCollection := database.DB.Collection("airlines")
	var airline model.Airline

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = airlineCollection.FindOne(ctx, bson.M{"_id": airlineID}).Decode(&airline)
	if err != nil {
		if err == mongo.ErrNoDocuments{
			utils.RespondWithError(w, http.StatusNotFound, "Airline not found")
		}else{
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

	_,err = airlineCollection.UpdateOne(ctx, bson.M{"_id": airlineID}, update)
	if err != nil {
		utils.Logger.Warn("Failed to update airline")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating airline")
		return 
	}

	utils.RespondWithJson(w, http.StatusOK, "Airline updated", map[string]interface{}{})

}

func DeleteAirline(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPut{
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
	if airlineIDStr == ""{
		utils.RespondWithError(w, http.StatusBadRequest, "Missing airline ID")
		return 
	}

	airlineID, err := primitive.ObjectIDFromHex(airlineIDStr)
	if err != nil{
		utils.RespondWithError(w, http.StatusBadRequest, "invalid airline ID")
		return 
	}


	airlineCollection := database.DB.Collection("airlines")
	var airline model.Airline

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = airlineCollection.FindOne(ctx, bson.M{"_id": airlineID}).Decode(&airline)
	if err != nil {
		if err == mongo.ErrNoDocuments{
			utils.RespondWithError(w, http.StatusNotFound, "Airline not found")
		}else{
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding airline")
		}
		return 
	}

	

	_,err = airlineCollection.DeleteOne(ctx, bson.M{"_id": airlineID})
	if err != nil {
		utils.Logger.Warn("Failed to Delete airline")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error Delete airline")
		return 
	}

	utils.RespondWithJson(w, http.StatusOK, "Airline Deleted", map[string]interface{}{})

}