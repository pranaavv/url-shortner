package views

import (
	"encoding/json"
	"net/http"
	"time"
)

// ErrorResponse defines the standardized JSON error format.
type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// ShortenResponse defines the response returned upon successful URL shortening.
type ShortenResponse struct {
	ID          string    `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
}

// MessageResponse defines a generic message response.
type MessageResponse struct {
	Message string `json:"message"`
}

// RenderJSON writes a JSON response with the provided status code.
func RenderJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

// RenderError writes a standardized JSON error response.
func RenderError(w http.ResponseWriter, statusCode int, message string) {
	RenderJSON(w, statusCode, ErrorResponse{
		Error:  message,
		Status: statusCode,
	})
}
