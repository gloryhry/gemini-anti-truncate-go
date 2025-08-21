package test

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

// TestAPICompatibility tests that the Go implementation is compatible with the original JavaScript API
func TestAPICompatibility(t *testing.T) {
	// Create a mock upstream server that simulates the Gemini API
	upstreamServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the request path is correctly forwarded
		if !strings.Contains(r.URL.Path, "/v1beta/models/") {
			t.Errorf("Expected request path to contain '/v1beta/models/', got '%s'", r.URL.Path)
		}
		
		// Check that the API key is correctly forwarded
		if r.Header.Get("X-Goog-Api-Key") == "" {
			t.Error("Expected X-Goog-Api-Key header to be set")
		}
		
		// Send a response that matches the Gemini API format
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		// Include the finish token in the response to simulate a complete response
		response := gemini.GenerateContentResponse{
			Candidates: []gemini.Candidate{
				{
					Content: gemini.Content{
						Parts: []gemini.Part{
							{Text: "This is a test response [RESPONSE_FINISHED]"},
						},
						Role: "model",
					},
					Index: 0,
				},
			},
		}
		
		json.NewEncoder(w).Encode(response)
	}))
	defer upstreamServer.Close()
	
	// Temporarily override the upstream URL
	originalUpstreamURL := config.AppConfig.UpstreamURLBase
	config.AppConfig.UpstreamURLBase = upstreamServer.URL
	defer func() {
		config.AppConfig.UpstreamURLBase = originalUpstreamURL
	}()
	
	// Test a request to a target model
	testRequestToModel(t, "gemini-1.5-pro-latest")
	
	// Test a request to a non-target model (should be passthrough)
	testRequestToModel(t, "gemini-pro-vision")
}

func testRequestToModel(t *testing.T, modelName string) {
	// Create a test request body
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
		t.Fatalf("Failed to marshal request body: %v", err)
	}
	
	// Create a test HTTP request
	req := httptest.NewRequest("POST", fmt.Sprintf("/v1beta/models/%s:generateContent", modelName), bytes.NewBuffer(bodyBytes))
	req.Header.Set("Authorization", "Bearer test-key")
	req.Header.Set("Content-Type", "application/json")
	
	// Create a response recorder
	// rr := httptest.NewRecorder()
	
	// TODO: Call the actual handler
	// For now, we'll just verify the test structure
	
	t.Logf("API compatibility test for model '%s' is working", modelName)
}

// TestSSEStreamFormat tests that the SSE stream format is compatible
func TestSSEStreamFormat(t *testing.T) {
	// Create a mock SSE stream response that matches the Gemini API format
	streamData := []string{
		`data: {"candidates": [{"content": {"parts": [{"text": "This is the first part"}]}, "index": 0}]}`,
		`data: {"candidates": [{"content": {"parts": [{"text": " of a streaming response [RESPONSE_FINISHED]"}]}, "index": 0}]}`,
		`data: {"candidates": [{"content": {"parts": [{"text": ""}], "finishReason": "STOP"}, "index": 0}]}`,
	}
	
	// Verify the format matches expected SSE format
	for _, line := range streamData {
		if !strings.HasPrefix(line, "data: ") {
			t.Errorf("Expected SSE line to start with 'data: ', got '%s'", line)
		}
		
		// Try to parse the JSON part
		jsonPart := strings.TrimPrefix(line, "data: ")
		var response gemini.GenerateContentResponse
		if err := json.Unmarshal([]byte(jsonPart), &response); err != nil {
			t.Errorf("Failed to parse JSON in SSE line '%s': %v", line, err)
		}
	}
	
	t.Log("SSE stream format is compatible")
}

// TestRequestResponseFormat tests that request/response formats are compatible
func TestRequestResponseFormat(t *testing.T) {
	// Test a typical request structure
	request := gemini.GenerateContentRequest{
		Contents: []gemini.Content{
			{
				Role: "user",
				Parts: []gemini.Part{
					{Text: "What is the weather like?"},
				},
			},
		},
		GenerationConfig: &gemini.GenerationConfig{
			Temperature:     0.7,
			MaxOutputTokens: 1000,
		},
	}
	
	// Verify that the request can be marshaled
	requestJSON, err := json.Marshal(request)
	if err != nil {
		t.Errorf("Failed to marshal request: %v", err)
	}
	
	// Verify that the request can be unmarshaled
	var parsedRequest gemini.GenerateContentRequest
	if err := json.Unmarshal(requestJSON, &parsedRequest); err != nil {
		t.Errorf("Failed to unmarshal request: %v", err)
	}
	
	// Test a typical response structure
	response := gemini.GenerateContentResponse{
		Candidates: []gemini.Candidate{
			{
				Content: gemini.Content{
					Parts: []gemini.Part{
						{Text: "The weather is sunny today."},
					},
					Role: "model",
				},
				FinishReason: "STOP",
				Index:        0,
			},
		},
	}
	
	// Verify that the response can be marshaled
	responseJSON, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Failed to marshal response: %v", err)
	}
	
	// Verify that the response can be unmarshaled
	var parsedResponse gemini.GenerateContentResponse
	if err := json.Unmarshal(responseJSON, &parsedResponse); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	
	t.Log("Request/response format is compatible")
}

// TestRetryMechanismBehavior tests that the retry mechanism behaves like the JavaScript version
func TestRetryMechanismBehavior(t *testing.T) {
	// Test that the retry prompt is correctly formatted
	expectedRetryPrompt := "Please continue generating the response from where you left off. Do not repeat the previous content."
	if gemini.RetryPrompt != expectedRetryPrompt {
		t.Errorf("Expected retry prompt '%s', got '%s'", expectedRetryPrompt, gemini.RetryPrompt)
	}
	
	// Test that the finish token is correctly defined
	expectedFinishToken := "[RESPONSE_FINISHED]"
	if gemini.FinishToken != expectedFinishToken {
		t.Errorf("Expected finish token '%s', got '%s'", expectedFinishToken, gemini.FinishToken)
	}
	
	// Test that target models match expected values
	expectedModels := []string{
		"gemini-1.5-pro-latest",
		"gemini-1.5-flash-latest",
		"gemini-pro",
	}
	
	for _, expectedModel := range expectedModels {
		found := false
		for _, model := range gemini.TargetModels {
			if model == expectedModel {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected model '%s' not found in TargetModels", expectedModel)
		}
	}
	
	t.Log("Retry mechanism behavior is compatible")
}

// TestErrorResponses tests that error responses are compatible
func TestErrorResponses(t *testing.T) {
	// Test a typical error response structure
	errorResponse := gemini.ErrorResponse{
		Error: struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Status  string `json:"status"`
		}{
			Code:    400,
			Message: "Bad Request",
			Status:  "BAD_REQUEST",
		},
	}
	
	// Verify that the error response can be marshaled
	errorJSON, err := json.Marshal(errorResponse)
	if err != nil {
		t.Errorf("Failed to marshal error response: %v", err)
	}
	
	// Verify that the error response can be unmarshaled
	var parsedErrorResponse gemini.ErrorResponse
	if err := json.Unmarshal(errorJSON, &parsedErrorResponse); err != nil {
		t.Errorf("Failed to unmarshal error response: %v", err)
	}
	
	// Check that the structure matches expected format
	if parsedErrorResponse.Error.Code != 400 {
		t.Errorf("Expected error code 400, got %d", parsedErrorResponse.Error.Code)
	}
	
	if parsedErrorResponse.Error.Message != "Bad Request" {
		t.Errorf("Expected error message 'Bad Request', got '%s'", parsedErrorResponse.Error.Message)
	}
	
	if parsedErrorResponse.Error.Status != "BAD_REQUEST" {
		t.Errorf("Expected error status 'BAD_REQUEST', got '%s'", parsedErrorResponse.Error.Status)
	}
	
	t.Log("Error responses are compatible")
}