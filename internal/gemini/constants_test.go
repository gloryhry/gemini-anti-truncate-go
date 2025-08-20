package gemini

import (
	"testing"
)

func TestConstants(t *testing.T) {
	// Test that constants are correctly defined
	expectedFinishToken := "[RESPONSE_FINISHED]"
	if FinishToken != expectedFinishToken {
		t.Errorf("Expected FinishToken to be '%s', got '%s'", expectedFinishToken, FinishToken)
	}
	
	expectedUserPromptSuffix := "\n\n(Note: If you are done, please end your response with [RESPONSE_FINISHED])"
	if UserPromptSuffix != expectedUserPromptSuffix {
		t.Errorf("Expected UserPromptSuffix to be '%s', got '%s'", expectedUserPromptSuffix, UserPromptSuffix)
	}
	
	expectedRetryPrompt := "Please continue generating the response from where you left off. Do not repeat the previous content."
	if RetryPrompt != expectedRetryPrompt {
		t.Errorf("Expected RetryPrompt to be '%s', got '%s'", expectedRetryPrompt, RetryPrompt)
	}
	
	expectedDefaultUpstreamURL := "https://generativelanguage.googleapis.com"
	if DefaultUpstreamURL != expectedDefaultUpstreamURL {
		t.Errorf("Expected DefaultUpstreamURL to be '%s', got '%s'", expectedDefaultUpstreamURL, DefaultUpstreamURL)
	}
	
	if DefaultMaxRetries != 20 {
		t.Errorf("Expected DefaultMaxRetries to be 20, got %d", DefaultMaxRetries)
	}
	
	if DefaultHTTPPort != 8080 {
		t.Errorf("Expected DefaultHTTPPort to be 8080, got %d", DefaultHTTPPort)
	}
	
	// Test that TargetModels contains expected values
	expectedModels := []string{
		"gemini-1.5-pro-latest",
		"gemini-1.5-flash-latest",
		"gemini-pro",
	}
	
	if len(TargetModels) != len(expectedModels) {
		t.Errorf("Expected %d target models, got %d", len(expectedModels), len(TargetModels))
	}
	
	// Check that all expected models are present
	for _, expectedModel := range expectedModels {
		found := false
		for _, model := range TargetModels {
			if model == expectedModel {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected model '%s' not found in TargetModels", expectedModel)
		}
	}
	
	// Test retryable status codes
	expectedRetryableStatus := []int{503, 403, 429}
	if len(RetryableStatus) != len(expectedRetryableStatus) {
		t.Errorf("Expected %d retryable status codes, got %d", len(expectedRetryableStatus), len(RetryableStatus))
	}
	
	// Test fatal status codes
	expectedFatalStatus := []int{500}
	if len(FatalStatus) != len(expectedFatalStatus) {
		t.Errorf("Expected %d fatal status codes, got %d", len(expectedFatalStatus), len(FatalStatus))
	}
}

func TestGenerateContentRequest_GetSystemInstruction(t *testing.T) {
	// Test when SystemInstruction_ is set (snake_case)
	sysInst := &SystemInstruction{
		Role: "system",
		Parts: []Part{{Text: "Test system instruction"}},
	}
	
	req := &GenerateContentRequest{
		SystemInstruction_: sysInst,
	}
	
	result := req.GetSystemInstruction()
	if result != sysInst {
		t.Error("Expected SystemInstruction_ to be returned when set")
	}
	
	// Test when SystemInstruction is set (camelCase)
	req = &GenerateContentRequest{
		SystemInstruction: sysInst,
	}
	
	result = req.GetSystemInstruction()
	if result != sysInst {
		t.Error("Expected SystemInstruction to be returned when set")
	}
	
	// Test when both are set (should prefer snake_case)
	sysInst2 := &SystemInstruction{
		Role: "system",
		Parts: []Part{{Text: "Test system instruction 2"}},
	}
	
	req = &GenerateContentRequest{
		SystemInstruction:  sysInst,   // camelCase
		SystemInstruction_: sysInst2,  // snake_case (should be preferred)
	}
	
	result = req.GetSystemInstruction()
	if result != sysInst2 {
		t.Error("Expected SystemInstruction_ to be preferred when both are set")
	}
	
	// Test when neither is set
	req = &GenerateContentRequest{}
	
	result = req.GetSystemInstruction()
	if result != nil {
		t.Error("Expected nil when neither SystemInstruction nor SystemInstruction_ is set")
	}
}

func TestGenerateContentRequest_SetSystemInstruction(t *testing.T) {
	sysInst := &SystemInstruction{
		Role: "system",
		Parts: []Part{{Text: "Test system instruction"}},
	}
	
	req := &GenerateContentRequest{
		SystemInstruction: &SystemInstruction{
			Role: "system",
			Parts: []Part{{Text: "Old instruction"}},
		},
	}
	
	req.SetSystemInstruction(sysInst)
	
	// Check that SystemInstruction_ is set
	if req.SystemInstruction_ != sysInst {
		t.Error("Expected SystemInstruction_ to be set")
	}
	
	// Check that SystemInstruction is cleared
	if req.SystemInstruction != nil {
		t.Error("Expected SystemInstruction to be cleared")
	}
}