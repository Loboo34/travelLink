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

// Create activity
func CreateActivity(w http.ResponseWriter, r *http.Request) {
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
		Title           string                   `json:"title"`
		Description     string                   `json:"description"`
		City            string                   `json:"city"`
		Country         string                   `json:"country"`
		Location        model.GeoLocation        `json:"location"`
		MeetingPoint    model.MeetingPoint       `json:"meetingPoint"`
		Categories      []model.ActivityCategory `json:"categories"`
		DurationMinutes int                      `json:"durationMinutes"`
		Inclusions      []string                 `json:"inclusions"`
		Exclusions      []string                 `json:"exclusions"`
		Images          []string                 `json:"images"`
	}

	req.Title = r.FormValue("title")
	if req.Title == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing activity title")
		return
	}

	req.Description = r.FormValue("description")
	if req.Description == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing activity description")
		return
	}

	if v := r.FormValue("city"); v != "" {
		if err := json.Unmarshal([]byte(v), &req.City); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid city format")
			return
		}
	}

	if v := r.FormValue("country"); v != "" {
		if err := json.Unmarshal([]byte(v), &req.Country); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid Country format")
			return
		}
	}

	if v := r.FormValue("location"); v != "" {
		if err := json.Unmarshal([]byte(v), &req.Location); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid location format")
			return
		}
	}

		if v := r.FormValue("meetingPoint"); v != "" {
		if err := json.Unmarshal([]byte(v), &req.MeetingPoint); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid meeting point format")
			return
		}
	}

		if v := r.FormValue("categories"); v != "" {
		if err := json.Unmarshal([]byte(v), &req.Categories); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid category format")
			return
		}
	}
	if v := r.FormValue("durationMinutes"); v != "" {
		if err := json.Unmarshal([]byte(v), &req.DurationMinutes); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid duration")
			return
		}
	}

		if v := r.FormValue("inclusions"); v != "" {
		if err := json.Unmarshal([]byte(v), &req.Inclusions); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid inclusion format")
			return
		}
	}

		if v := r.FormValue("Exclusion"); v != "" {
		if err := json.Unmarshal([]byte(v), &req.Exclusions); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid exclusion formart")
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	  var imageURLs []string

    if r.MultipartForm != nil && r.MultipartForm.File["images"] != nil {
        files := r.MultipartForm.File["images"]

        if len(files) > 10 {
            utils.RespondWithError(w, http.StatusBadRequest, "maximum 10 images allowed")
            return
        }

        for _, fileHeader := range files {
            file, err := fileHeader.Open()
            if err != nil {
                utils.RespondWithError(w, http.StatusInternalServerError, "failed to read image file")
                return
            }
            defer file.Close()

            url, err := utils.UploadImage(ctx, file, "accommodations")
            if err != nil {
                utils.Logger.Warn("cloudinary upload failed")
                utils.RespondWithError(w, http.StatusInternalServerError, "image upload failed")
                return
            }

            imageURLs = append(imageURLs, url)
        }
    }

	activityCollection := database.DB.Collection("activities")


	create := model.Activity{
		ID:              primitive.NewObjectID(),
		Title:           req.Title,
		Description:     req.Description,
		City:            req.City,
		Country:         req.Country,
		Location:        req.Location,
		MeetingPoint:    req.MeetingPoint,
		Categories:      req.Categories,
		DurationMinutes: req.DurationMinutes,
		Inclusions:      req.Inclusions,
		Exclusions:      req.Exclusions,
		Images:          imageURLs,
		ReviewCount:     0,
		Rating:          0.0,
		IsActive:        true,
		CreatedAt:       time.Now(),
	}

	_, err = activityCollection.InsertOne(ctx, create)
	if err != nil {
		utils.Logger.Warn("Failed to create activity")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating activity")
	}

	utils.Logger.Info("Created activity")
	utils.RespondWithJson(w, http.StatusCreated,  map[string]interface{}{})
}

// Update activity
func UpdateActivity(w http.ResponseWriter, r *http.Request) {
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
	activityIDStr := vars["activityID"]
	if activityIDStr == "" {
		utils.RespondWithError(w, http.StatusNotFound, "Missing flight ID")
		return
	}

	activityID, err := primitive.ObjectIDFromHex(activityIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid flight ID")
		return
	}

	var req struct {
		Title           string             `json:"title"`
		MeetingPoint    model.MeetingPoint `json:"price"`
		DurationMinutes int                `json:"durationMinutes"`
		Inclusions      []string           `json:"inclusions"`
		Exclusions      []string           `json:"exclusions"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	activityCollection := database.DB.Collection("activities")
	var activity model.Activity

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = activityCollection.FindOne(ctx, bson.M{"_id": activityID}).Decode(&activity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "activity not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding")
		}
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":           req.Title,
			"meetingPoint":    req.MeetingPoint,
			"durationMinutes": req.DurationMinutes,
			"inclusions":      req.Inclusions,
			"exclusion":       req.Exclusions,
			"updatedAt":       time.Now(),
		},
	}

	_, err = activityCollection.UpdateOne(ctx, bson.M{"_id": activityID}, update)
	if err != nil {
		utils.Logger.Warn("Failed to update activity")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating activity")
		return
	}

	utils.Logger.Info("Updated activity")
	utils.RespondWithJson(w, http.StatusOK,  map[string]interface{}{})

}

// delete activity
func DeleteActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "ONly DELETE allowed")
		return
	}
	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing Admin ID")
		return
	}

	vars := mux.Vars(r)
	activityIDStr := vars["activityID"]
	if activityIDStr == "" {
		utils.RespondWithError(w, http.StatusNotFound, "Missing flight ID")
		return
	}

	activityID, err := primitive.ObjectIDFromHex(activityIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid flight ID")
		return
	}

	activityCollection := database.DB.Collection("activities")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := activityCollection.DeleteOne(ctx, bson.M{"_id": activityID})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting activity")
		utils.Logger.Warn("Failed to delete activity ")
		return
	}

	if result.DeletedCount == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Activity not found")
	}

	utils.Logger.Info("Deleted activity")
	utils.RespondWithJson(w, http.StatusOK,  map[string]interface{}{})

}

// create timeslot
func TimeSlot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only Post allowed")
		return
	}

	_, err := utils.GetAdminID()
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	activityIDStr := vars["activityID"]
	if activityIDStr == "" {
		utils.RespondWithError(w, http.StatusNotFound, "Missing flight ID")
		return
	}

	activityID, err := primitive.ObjectIDFromHex(activityIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid flight ID")
		return
	}

	var req struct {
		StartTime       time.Time `json:"startTime"`
		DurationMinutes int       `json:"durationMinutes"`
		TotalSlots      int       `json:"totalSlots"`
		PricePerPerson  int64     `json:"pricePerPerson"`
		GroupSizeMax    int       `json:"groupSizeMax"`
	}

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	timeSlotCollection := database.DB.Collection("activity_timeslot")

	activityCollection := database.DB.Collection("activities")
	var activity model.Activity

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = activityCollection.FindOne(ctx, bson.M{"_id": activityID}).Decode(&activity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.RespondWithError(w, http.StatusNotFound, "Accommodation not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error finding accommodation")
		}
		return
	}

	create := model.ActivityTimeslot{
		ID:              primitive.NewObjectID(),
		ActivityID:      activityID,
		StartTime:       req.StartTime,
		DurationMinutes: req.DurationMinutes,
		TotalSlots:      req.TotalSlots,
		ReservedSlots:   0,
		PricePerPerson:  req.PricePerPerson,
		GroupSizeMax:    req.GroupSizeMax,
		IsActive:        true,
		CreatedAt:       time.Now(),
	}

	_, err = timeSlotCollection.InsertOne(ctx, create)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating time slot")
		utils.Logger.Warn("Failed to create a time slot")
		return
	}

	utils.Logger.Info("Time slot created")
	utils.RespondWithJson(w, http.StatusOK,  map[string]interface{}{})

}

//booking stats
