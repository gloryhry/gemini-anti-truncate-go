package util

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	// Test with Authorization header (Bearer)
	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("Authorization", "Bearer test-api-key")
	
	apiKey := GetAPIKey(req)
	if apiKey != "test-api-key" {
		t.Errorf("Expected 'test-api-key', got '%s'", apiKey)
	}
	
	// Test with X-Goog-Api-Key header
	req = httptest.NewRequest("POST", "/", nil)
	req.Header.Set("X-Goog-Api-Key", "test-api-key-2")
	
	apiKey = GetAPIKey(req)
	if apiKey != "test-api-key-2" {
		t.Errorf("Expected 'test-api-key-2', got '%s'", apiKey)
	}
	
	// Test with Authorization header taking precedence
	req = httptest.NewRequest("POST", "/", nil)
	req.Header.Set("Authorization", "Bearer test-api-key-3")
	req.Header.Set("X-Goog-Api-Key", "test-api-key-4")
	
	apiKey = GetAPIKey(req)
	if apiKey != "test-api-key-3" {
		t.Errorf("Expected 'test-api-key-3', got '%s'", apiKey)
	}
	
	// Test with no API key headers
	req = httptest.NewRequest("POST", "/", nil)
	
	apiKey = GetAPIKey(req)
	if apiKey != "" {
		t.Errorf("Expected empty string, got '%s'", apiKey)
	}
	
	// Test with malformed Authorization header
	req = httptest.NewRequest("POST", "/", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	
	apiKey = GetAPIKey(req)
	if apiKey != "" {
		t.Errorf("Expected empty string for non-Bearer Authorization, got '%s'", apiKey)
	}
}

func TestSendJSONError(t *testing.T) {
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Send a JSON error
	SendJSONError(rr, "Test error message", http.StatusBadRequest)
	
	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}
	
	// Check the content type
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
	
	// Check the response body contains the error message
	expected := `"message":"Test error message"`
	if body := rr.Body.String(); !contains(body, expected) {
		t.Errorf("Expected response body to contain '%s', got '%s'", expected, body)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		(s == substr || 
			(len(s) > len(substr) && 
				(s[:len(substr)] == substr || 
					s[len(s)-len(substr):] == substr || 
					findSubstring(s, substr) != -1)))
}

// Simple substring search
func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}