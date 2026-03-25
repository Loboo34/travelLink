package utils

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Success bool      `json:"success"`
	Error   *ApiError `json:"error,omitempty"`
	Data    any       `json:"data,omitempty"`
}

type ApiError struct {
	Code    int    `json:"coode"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func RespondWithJson(w http.ResponseWriter, code int,  payload any) {

	response := ApiResponse{
		Success: true,
		Data:    payload,
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {

	response := ApiResponse{
		Success: false,
		Error:   &ApiError{
			Code: code,
			Status: http.StatusText(code),
			Message: message,
		},
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)

}
