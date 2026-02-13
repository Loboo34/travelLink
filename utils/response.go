package utils

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func RespondWithJson(w http.ResponseWriter, code int, message string, payload interface{}) {

	response := ApiResponse{
		Success: true,
		Message: message,
		Data:    payload,
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {

	response := ApiResponse{
		Success: false,
		Error:   message,
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)

}
