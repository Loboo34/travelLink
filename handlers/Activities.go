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

//admin
//Create activity
func CreateActivity(w http.ResponseWriter, r *http.Request){
    if r.Method != http.MethodPost{
        utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST allowed")
        return 
    }

    _,err := utils.GetAdminID()
    if err != nil{
        utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
        return 
    }

    var req struct {
        Title string `json:"title"`
        Price float64 `json:"price"`
        Description string `json:"description"`
        Location string `json:"location"`
        Categories []string `json:"categories"`
        DurationMinutes int `json:"durationMinutes"`
        Inclusions []string `json:"inclusions"`
        Exclusions []string `json:"exclusions"`
        Images []string `json:"images"`
    }

    if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
        return 
    }

    activityCollection := database.DB.Collection("activities")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    create := model.Activity{
        ID: primitive.NewObjectID(),
        Title: req.Title,
        Price: req.Price,
        Description: req.Description,
        Location: req.Location,
        Categories: req.Categories,
        DurationMinutes: req.DurationMinutes,
        Inclusions: req.Inclusions,
        Exclusions: req.Exclusions,
        Images: req.Images,
    }

    _,err = activityCollection.InsertOne(ctx, create)
    if err != nil {
        utils.Logger.Warn("Failed to create activity")
         utils.RespondWithError(w, http.StatusInternalServerError, "Error creating activity")
    }

    utils.Logger.Info("Created activity")
    utils.RespondWithJson(w, http.StatusCreated, "Activity created successfully", map[string]interface{}{})
}
//Update activity
//delete activity
//create timeslot
//booking stats


   //user
//get all activities/activity
//search activity
//get activity time slots
//check availability
//get reviews
//get by location
//get by categories
//book activity