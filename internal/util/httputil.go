package util

import (
	"encoding/json"
	"gemini-anti-truncate-go/internal/gemini"
	"net/http"
	"strings"
)

// GetAPIKey extracts the Gemini API key from the request headers.
// It checks for "Authorization: Bearer <key>" and "X-Goog-Api-Key: <key>".
func GetAPIKey(r *http.Request) string {
	// Check Authorization header first
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	// Fallback to X-Goog-Api-Key
	if apiKey := r.Header.Get("X-Goog-Api-Key"); apiKey != "" {
		return apiKey
	}

	return ""
}

// SendJSONError sends a standardized JSON error response to the client.
func SendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := gemini.ErrorResponse{
		Error: struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Status  string `json:"status"`
		}{
			Code:    statusCode,
			Message: message,
			Status:  http.StatusText(statusCode),
		},
	}

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		// If encoding fails, log it and fall back to a plain text response
		Errorf("Failed to encode JSON error response: %v", err)
		http.Error(w, `{"error": "Failed to serialize error message."}`, http.StatusInternalServerError)
	}
}
