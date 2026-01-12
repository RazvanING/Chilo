package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponseBody struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	JSONResponse(w, statusCode, ErrorResponseBody{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

func SuccessResponseWithData(w http.ResponseWriter, data interface{}) {
	JSONResponse(w, http.StatusOK, SuccessResponse{Data: data})
}

func SuccessResponseWithMessage(w http.ResponseWriter, message string) {
	JSONResponse(w, http.StatusOK, SuccessResponse{Message: message})
}
