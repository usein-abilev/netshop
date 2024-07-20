package tools

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	Status int         `json:"status"`
	Error  ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

func RespondWithError(w http.ResponseWriter, message string, status int) {
	response := ErrorResponse{
		Status: status,
		Error: ErrorDetail{
			Message: message,
			Details: nil,
		},
	}
	respondWithJSON(w, status, response)
}

func RespondWithSuccess(w http.ResponseWriter, data interface{}) {
	response := SuccessResponse{
		Status: 200,
		Data:   data,
	}
	respondWithJSON(w, 200, response)
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;encoding=utf-8")
	w.WriteHeader(status)
	w.Write(bytes)
}
