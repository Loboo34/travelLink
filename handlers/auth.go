package handlers

import (
	"encoding/json"
	"net/http"

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

	utils.RespondWithJson(w, http.StatusCreated, result)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "only GET allowed")
		return 
	}

	userID, err := utils.GetUserID(r.Context())
	if err != nil{
		utils.RespondWithError(w, http.StatusUnauthorized, "missing user ID")
		return 
	}

	result, err := h.userService.GetProfile(r.Context(), userID)
	if err != nil{
		HandleServiceError(w, err, "failed to get profile")
		return 
	}

	utils.RespondWithJson(w, http.StatusOK, result)
}
