package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Loboo34/travel/auth"
	"github.com/Loboo34/travel/utils"
)

type UserHandler struct {
	userService *auth.UserService
}

func NewUserHandler(userService *auth.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "only POST allowed")
		return
	}

	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	defer r.Body.Close()

	result, err := h.userService.Register(r.Context(), req)
	if err != nil {
		HandleServiceError(w, err, "failed registering user")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, result)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "only POST allowed")
		return
	}

	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	defer r.Body.Close()

	result, err := h.userService.Login(r.Context(), req)
	if err != nil {
		HandleServiceError(w, err, "failed to login")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, result)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "only GET allowed")
		return
	}

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "missing user ID")
		return
	}

	result, err := h.userService.GetProfile(r.Context(), userID)
	if err != nil {
		HandleServiceError(w, err, "failed to get profile")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, result)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "only PATCH allowed")
		return
	}

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "missing user ID")
		return
	}

	var req auth.UpdateProfile

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid json format")
		return
	}

	result, err := h.userService.Update(r.Context(), userID, req)
	if err != nil {
		HandleServiceError(w, err, "failed to update profile")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, result)
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "only PATCH allowed")
		return
	}

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "missing user ID")
		return
	}

	var req auth.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	if err := h.userService.ChangePassword(r.Context(), userID, req); err != nil {
		HandleServiceError(w, err, "failed to change password")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]string{"message": "password updated successfully"})
}

func (h *UserHandler) ProfilePic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "only PATCH allowed")
		return
	}

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "missing user ID")
		return
	}

	// Parse multipart form with max 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	files := r.MultipartForm.File["image"]
	if len(files) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "no image provided")
		return
	}
	if len(files) > 1 {
		utils.RespondWithError(w, http.StatusBadRequest, "maximum 1 image allowed")
		return
	}

	fileHeader := files[0]
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	file, err := fileHeader.Open()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to read image file")
		return
	}
	defer file.Close()

	url, err := utils.UploadImage(ctx, file, "profile")
	if err != nil {
		utils.Logger.Warn("cloudinary upload failed")
		utils.RespondWithError(w, http.StatusInternalServerError, "image upload failed")
		return
	}

	if err := h.userService.ProfilePic(r.Context(), userID, url); err != nil {
		HandleServiceError(w, err, "failed to update profile picture")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]string{"message": "profile picture updated successfully", "image_url": url})
}



