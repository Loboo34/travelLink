package handlers_admin

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Loboo34/travel/service"
	"github.com/Loboo34/travel/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccommodationHandler struct {
	accommodationService *service.AccommodationService
}

func NewAccommodationHandler(accommodationService *service.AccommodationService) *AccommodationHandler {
	return &AccommodationHandler{accommodationService: accommodationService}
}

// Create accommodaion
func (h *AccommodationHandler) AddAccommodation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Missing admin ID")
		return
	}

	var req service.AccommodationRequest

	if err := r.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	if r.MultipartForm != nil {
		if v := r.FormValue("name"); v != "" {
			req.Name = v
		}

		if v := r.FormValue("description"); v != "" {
			req.Description = v
		}

		if v := r.FormValue("propertyType"); v != "" {
			if err := json.Unmarshal([]byte(v), &req.PropertyType); err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "invalid roomType format")
				return
			}
		}

		if v := r.FormValue("address"); v != "" {
			if err := json.Unmarshal([]byte(v), &req.Address); err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "invalid address format")
				return
			}

		}

		if v := r.FormValue("location"); v != "" {
			if err := json.Unmarshal([]byte(v), &req.Location); err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "invalid location format")
				return
			}
		}
		if v := r.FormValue("amenities"); v != "" {
			if err := json.Unmarshal([]byte(v), &req.Amenities); err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "invalid amenities format")
				return
			}
		}
		if v := r.FormValue("roomType"); v != "" {
			if err := json.Unmarshal([]byte(v), &req.RoomType); err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "invalid roomType format")
				return
			}
		}
		if v := r.FormValue("hostID"); v != "" {
			id, err := primitive.ObjectIDFromHex(v)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "invalid hostID format")
				return
			}
			req.HostID = &id
		}

		var imageURLs []string
		if r.MultipartForm.File != nil && r.MultipartForm.File["images"] != nil {
			files := r.MultipartForm.File["images"]

			if len(files) > 10 {
				utils.RespondWithError(w, http.StatusBadRequest, "maximum 10 images allowed")
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

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

			req.Images = imageURLs
		}

	} else {
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
	}


	result, err := h.accommodationService.Add(r.Context(), req)
	if err != nil {

		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating accommodation")
		return
	}

	utils.Logger.Info("Successfully created Accommodation")
	utils.RespondWithJson(w, http.StatusCreated, result)

}

// update accommodation
func (h *AccommodationHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PATCH allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	accommodationIDStr := vars["accommodationID"]
	if accommodationIDStr == "" {
		utils.RespondWithError(w, http.StatusNotFound, "Missing accommodation ID")
		return
	}

	accommodaionID, err := primitive.ObjectIDFromHex(accommodationIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid flight ID")
		return
	}

	var req service.AccommodationUpdate

	if err := r.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid multipart form")
		return
	}

	if r.MultipartForm != nil {
		// Support form-based update fields
		if v := r.FormValue("name"); v != "" {
			req.Name = v
		}
		if v := r.FormValue("description"); v != "" {
			req.Description = v
		}
		if v := r.FormValue("amenities"); v != "" {
			if err := json.Unmarshal([]byte(v), &req.Amenities); err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "invalid amenities format")
				return
			}
		}

		var imageURLs []string
		if r.MultipartForm.File != nil && r.MultipartForm.File["images"] != nil {
			files := r.MultipartForm.File["images"]

			if len(files) > 10 {
				utils.RespondWithError(w, http.StatusBadRequest, "maximum 10 images allowed")
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

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

			req.Images = imageURLs
		}
	} else {
		// Fallback to JSON body updates
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
	}

	result, err := h.accommodationService.Update(r.Context(), accommodaionID, req)
	if err != nil {
		utils.Logger.Warn("Failed to update accommodation")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating accommodation")
		return
	}

	utils.Logger.Info("Accommodation updated")
	utils.RespondWithJson(w, http.StatusOK, result)

}

func (h *AccommodationHandler) Availability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Missing admin ID")
		return
	}

	var req service.AvailabilityRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	result, err := h.accommodationService.Availability(r.Context(), req)
	if err != nil {
		utils.Logger.Warn("Failed to update accommodation")
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating accommodation")
		return
	}

	utils.Logger.Info("Accommodation availability created")
	utils.RespondWithJson(w, http.StatusOK, result)

}

func (h *AccommodationHandler) IsActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PATCH allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	accommodationIDStr := vars["availabilityID"]
	if accommodationIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing accommodation ID")
		return
	}

	accommodationID, err := primitive.ObjectIDFromHex(accommodationIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid accommodation ID")
		return
	}

	var req struct {
		IsActive bool `json:"isActive"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	err = h.accommodationService.IsActive(r.Context(), accommodationID, req.IsActive)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting accommodation")
		utils.Logger.Warn("Failed to delete accommodation")
		return
	}

	utils.Logger.Info("Accommodation deleted successfully")
	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{})
}

// Delete accommodation
func (h *AccommodationHandler) DeleteAccommodation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only DELETE allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	accommodationIDStr := vars["accommodationID"]
	if accommodationIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing accommodation ID")
		return
	}

	accommodationID, err := primitive.ObjectIDFromHex(accommodationIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid accommodation ID")
		return
	}

	err = h.accommodationService.Delete(r.Context(), accommodationID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting accommodation")
		utils.Logger.Warn("Failed to delete accommodation")
		return
	}

	utils.Logger.Info("Accommodation deleted successfully")
	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{})
}

func (h *AccommodationHandler) RemoveAvailability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only DELETE allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	accommodationIDStr := vars["accommodationID"]
	if accommodationIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing accommodation ID")
		return
	}

	accommodationID, err := primitive.ObjectIDFromHex(accommodationIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid accommodation ID")
		return
	}

	err = h.accommodationService.Remove(r.Context(), accommodationID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting accommodation")
		utils.Logger.Warn("Failed to delete accommodation")
		return
	}

	utils.Logger.Info("Accommodation deleted successfully")
	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{})
}

//booking stats
