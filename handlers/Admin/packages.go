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
// create package
func CreatePackage(w http.ResponseWriter, r *http.Request) {
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
		Name               string                   `json:"name"`
		Description        string                   `json:"description"`
		Destination        string                   `json:"destination"`
		DurationDays       int                      `json:"durationDays"`
		StartDateFrom      time.Time                `json:"startDateFrom"`
		StartDateTo        time.Time                `json:"startDateTo"`
		BasePrice          float64                  `json:"basePrice"`
		IncludedComponents []model.ComponentSummary `json:"includedComponents"`
		Tags               []string                 `json:"tags"`
		Images             []string                 `json:"images"`
		ExpiresAt          time.Time                `json:"expiresAt"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	packageCollection := database.DB.Collection("packags")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var expiresAt *time.Time
	if !req.ExpiresAt.IsZero() {
		expiresAt = &req.ExpiresAt
	}

	create := model.Package{
		ID:                 primitive.NewObjectID(),
		Name:               req.Name,
		Description:        req.Description,
		Destination:        req.Destination,
		DurationDays:       req.DurationDays,
		StartDateFrom:      req.StartDateFrom,
		StartDateTo:        req.StartDateTo,
		BasePrice:          req.BasePrice,
		IncludedComponents: req.IncludedComponents,
		Tags:               req.Tags,
		Images:             req.Images,
		ExpiresAt:          *expiresAt,
		IsActive:           true,
	}

	_, err = packageCollection.InsertOne(ctx, create)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating package")
		utils.Logger.Warn("Failed to create package")
	}

	utils.RespondWithJson(w, http.StatusCreated, "Package created", map[string]interface{}{})
	utils.Logger.Info("Created package successfully")
}

// update package
func UpadatePackage(w http.ResponseWriter, r *http.Request) {
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
	packageIDStr := vars["packageID"]
	if packageIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing paackage ID")
		return
	}

	packageID, err := primitive.ObjectIDFromHex(packageIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid package ID")
		return
	}

	var req struct {
		Name               string                   `json:"name"`
		Description        string                   `json:"desription"`
		Destination        string                   `json:"desination"`
		DurationDays       int                      `json:"durationDays"`
		StartDateFrom      time.Time                `json:"startDateFrom"`
		StartDateTo        time.Time                `json:"startDateTo"`
		BasePrice          float64                  `json:"basePrice"`
		IncludedComponents []model.ComponentSummary `json:"includedComponents"`
		Tags               []string                 `json:"tags"`
		Images             []string                 `json:"images"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	packageCollection := database.DB.Collection("packags")
	var pack model.Package

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = packageCollection.FindOne(ctx, bson.M{"_id": packageID}).Decode(&pack)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "package not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding package")
		}
		return
	}

	update := bson.M{
		"$set": bson.M{
			"name":               req.Name,
			"description":        req.Description,
			"Destination":        req.Destination,
			"durationDays":       req.DurationDays,
			"StartDateFrom":      req.StartDateFrom,
			"startDateTo":        req.StartDateTo,
			"basePrice":          req.BasePrice,
			"includedComponents": req.IncludedComponents,
			"tags":               req.Tags,
			"images":             req.Images,
		},
	}

	_, err = packageCollection.UpdateOne(ctx, bson.M{"_id": packageID}, update)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating package")
		utils.Logger.Warn("Failed to update package")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, "Package Updated", map[string]interface{}{})
	utils.Logger.Info("Updated package successfully")

}

// publish/unpublish package
func Active(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PATCH allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing Admin ID")
		return
	}

	vars := mux.Vars(r)
	packageIDStr := vars["packageID"]
	if packageIDStr == "" {
		utils.RespondWithError(w, http.StatusNotFound, "Missing package ID")
		return
	}

	packageID, err := primitive.ObjectIDFromHex(packageIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid package ID")
		return
	}

	var req struct {
		IsActive bool `json:"isActive"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	packageCollection := database.DB.Collection("packages")
	//var package model.FlightOffer

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"isActive": req.IsActive,
		},
	}

	result, err := packageCollection.UpdateOne(ctx, bson.M{"_id": packageID}, update)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating package status")
		utils.Logger.Warn("Failed to update package status")
		return
	}

	if result.MatchedCount == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Package not found")
		return
	}

	utils.Logger.Info("Update successful")
	utils.RespondWithJson(w, http.StatusOK, "Package status updated", map[string]interface{}{
		"isActive": req.IsActive,
	})
}

// delete
func DeletePackage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only DELETE allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing package ID")
		return
	}

	vars := mux.Vars(r)
	packageIDStr := vars["packagID"]
	if packageIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing package ID")
		return
	}

	PackageID, err := primitive.ObjectIDFromHex(packageIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid package ID")
		return
	}

	packageCollection := database.DB.Collection("packages")
	var pack model.Package

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = packageCollection.FindOne(ctx, bson.M{"_id": PackageID}).Decode(&pack)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Package not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding package")
			utils.Logger.Warn("Failed to find package")
		}
		return
	}

	_, err = packageCollection.DeleteOne(ctx, bson.M{"_id": PackageID})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting package")
		utils.Logger.Warn("Failed to delete package ")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, "Deleted package successfully", map[string]interface{}{})
	utils.Logger.Info("Package deleted")
}

// package stats

