package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gemini-anti-truncate-go/internal/config"
	"gemini-anti-truncate-go/internal/gemini"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Initialize config for tests
func init() {
	config.Load()
}

func TestHeaderSuppressingWriter(t *testing.T) {
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Wrap it with our header suppressing writer
	hsw := &headerSuppressingWriter{ResponseWriter: rr}
	
	// Set a header
	hsw.Header().Set("Content-Type", "text/event-stream")
	
	// Write headers
	hsw.WriteHeader(http.StatusOK)
	
	// Try to write headers again (should be suppressed)
	hsw.WriteHeader(http.StatusBadRequest)
	
	// Write some data
	hsw.Write([]byte("test data"))
	
	// Check that the status is OK (not bad request)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}
	
	// Check that the header was set
	if rr.Header().Get("Content-Type") != "text/event-stream" {
		t.Errorf("Expected Content-Type 'text/event-stream', got '%s'", rr.Header().Get("Content-Type"))
	}
	
	// Check that the data was written
	if rr.Body.String() != "test data" {
		t.Errorf("Expected body 'test data', got '%s'", rr.Body.String())
	}
}

func TestProxyHandler_MethodNotAllowed(t *testing.T) {
	// Create a GET request (should be rejected)
	req := httptest.NewRequest("GET", "/v1beta/models/gemini-pro", nil)
	rr := httptest.NewRecorder()
	
	// Call the handler
	ProxyHandler(rr, req)
	
	// Check the response
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
	
	// Check that the response contains an error message
	if !strings.Contains(rr.Body.String(), "Method not allowed") {
		t.Errorf("Expected error message about method not allowed, got '%s'", rr.Body.String())
	}
}

func TestProxyHandler_MissingAPIKey(t *testing.T) {
	// Create a POST request without API key
	req := httptest.NewRequest("POST", "/v1beta/models/gemini-pro", nil)
	rr := httptest.NewRecorder()
	
	// Call the handler
	ProxyHandler(rr, req)
	
	// Check the response
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
	}
	
	// Check that the response contains an error message
	if !strings.Contains(rr.Body.String(), "API key is missing") {
		t.Errorf("Expected error message about missing API key, got '%s'", rr.Body.String())
	}
}

func TestProxyHandler_InvalidJSON(t *testing.T) {
	// Create a POST request with invalid JSON
	req := httptest.NewRequest("POST", "/v1beta/models/gemini-pro", strings.NewReader("invalid json"))
	req.Header.Set("Authorization", "Bearer test-key")
	rr := httptest.NewRecorder()
	
	// Call the handler
	ProxyHandler(rr, req)
	
	// Check the response
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}
	
	// Check that the response contains an error message
	if !strings.Contains(rr.Body.String(), "Invalid JSON") {
		t.Errorf("Expected error message about invalid JSON, got '%s'", rr.Body.String())
	}
}

func TestPassthroughRequest(t *testing.T) {
	// Create a test request
	reqBody := gemini.GenerateContentRequest{
		Contents: []gemini.Content{
			{
				Role: "user",
				Parts: []gemini.Part{
					{Text: "Hello, world!"},
				},
			},
		},
	}
	
	// Marshal the request body
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal("Failed to marshal request body")
	}
	
	// Create a test HTTP request
	req := httptest.NewRequest("POST", "/v1beta/models/gemini-pro", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer test-key")
	req.Header.Set("Content-Type", "application/json")
	
	// Create a response recorder
	rr := httptest.NewRecorder()
	
	// Create a mock upstream server
	upstreamServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request was forwarded correctly
		if r.Header.Get("X-Goog-Api-Key") != "test-key" {
			t.Error("Expected X-Goog-Api-Key header to be set")
		}
		
		// Send a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"candidates": [{"content": {"parts": [{"text": "Hello, world!"}]}}]}`)
	}))
	defer upstreamServer.Close()
	
	// Temporarily override the upstream URL
	originalUpstreamURL := config.AppConfig.UpstreamURLBase
	config.AppConfig.UpstreamURLBase = upstreamServer.URL
	defer func() {
		config.AppConfig.UpstreamURLBase = originalUpstreamURL
	}()
	
	// Call the passthrough function
	passthroughRequest(rr, req, "test-key", &reqBody)
	
	// Check the response
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}
	
	// Check that the response contains the expected content
	if !strings.Contains(rr.Body.String(), "Hello, world!") {
		t.Errorf("Expected response to contain 'Hello, world!', got '%s'", rr.Body.String())
	}
}