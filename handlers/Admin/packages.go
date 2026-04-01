package handlers_admin

import (
	"encoding/json"
	"net/http"

	"github.com/Loboo34/travel/service"
	"github.com/Loboo34/travel/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackageHandler struct {
	packageService *service.PackageService
}

func NewPackageHandler(packageService *service.PackageService) *PackageHandler {
	return &PackageHandler{packageService: packageService}
}

func (h *PackageHandler) CreatePackage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Missing admin ID")
		return
	}

	var req service.PackageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	result, err := h.packageService.Create(r.Context(), req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating package")
		return
	}

	utils.Logger.Info("Created package")
	utils.RespondWithJson(w, http.StatusCreated, result)
}

func (h *PackageHandler) UpdatePackage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only PUT/PATCH allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	packageIDStr := vars["packageID"]
	if packageIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing package ID")
		return
	}

	packageID, err := primitive.ObjectIDFromHex(packageIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid package ID")
		return
	}

	var req service.PackageUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	result, err := h.packageService.Update(r.Context(), packageID, req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating package")
		return
	}

	utils.Logger.Info("Updated package")
	utils.RespondWithJson(w, http.StatusOK, result)
}

func (h *PackageHandler) SetActivePackage(w http.ResponseWriter, r *http.Request) {
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
	packageIDStr := vars["packageID"]
	if packageIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing package ID")
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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if err := h.packageService.SetActive(r.Context(), packageID, req.IsActive); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating package status")
		return
	}

	utils.Logger.Info("Updated package status")
	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{
		"isActive": req.IsActive,
	})
}

func (h *PackageHandler) DeletePackage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only DELETE allowed")
		return
	}

	_, err := utils.GetAdminID(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Missing admin ID")
		return
	}

	vars := mux.Vars(r)
	packageIDStr := vars["packageID"]
	if packageIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing package ID")
		return
	}

	packageID, err := primitive.ObjectIDFromHex(packageIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid package ID")
		return
	}

	if err := h.packageService.Delete(r.Context(), packageID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting package")
		return
	}

	utils.Logger.Info("Deleted package")
	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{})
}

// package stats
