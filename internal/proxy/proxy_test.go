package proxy

import (
	"gemini-anti-truncate-go/internal/gemini"
	"testing"
)

func TestInjectFinishToken(t *testing.T) {
	// Create a test request
	req := &gemini.GenerateContentRequest{
		Contents: []gemini.Content{
			{
				Role: "user",
				Parts: []gemini.Part{
					{Text: "Hello, world!"},
				},
			},
		},
	}
	
	// Apply the injection
	result := InjectFinishToken(req)
	
	// Check that the user prompt suffix was added
	expectedSuffix := "\n\n(Note: If you are done, please end your response with [RESPONSE_FINISHED])"
	if len(result.Contents) == 0 || len(result.Contents[0].Parts) == 0 {
		t.Fatal("Expected contents with parts")
	}
	
	if !contains(result.Contents[0].Parts[0].Text, expectedSuffix) {
		t.Errorf("Expected user prompt to contain suffix '%s', got '%s'", expectedSuffix, result.Contents[0].Parts[0].Text)
	}
	
	// Check that a system instruction was added
	sysInst := result.GetSystemInstruction()
	if sysInst == nil {
		t.Error("Expected system instruction to be added")
	} else {
		expectedInstruction := "Please ensure your response ends with [RESPONSE_FINISHED]"
		if len(sysInst.Parts) == 0 || !contains(sysInst.Parts[0].Text, expectedInstruction) {
			t.Errorf("Expected system instruction to contain '%s', got '%s'", expectedInstruction, sysInst.Parts[0].Text)
		}
	}
	
	// Test with existing system instruction
	req2 := &gemini.GenerateContentRequest{
		Contents: []gemini.Content{
			{
				Role: "user",
				Parts: []gemini.Part{
					{Text: "Hello, world!"},
				},
			},
		},
		SystemInstruction_: &gemini.SystemInstruction{
			Role: "system",
			Parts: []gemini.Part{
				{Text: "Existing instruction"},
			},
		},
	}
	
	result2 := InjectFinishToken(req2)
	
	// Check that the existing instruction was modified
	sysInst2 := result2.GetSystemInstruction()
	if sysInst2 == nil {
		t.Error("Expected system instruction to exist")
	} else {
		expectedAppend := "\n\nPlease ensure your response ends with [RESPONSE_FINISHED]"
		if len(sysInst2.Parts) == 0 || !contains(sysInst2.Parts[0].Text, expectedAppend) {
			t.Errorf("Expected system instruction to contain appended text '%s', got '%s'", expectedAppend, sysInst2.Parts[0].Text)
		}
	}
}

func TestBuildRetryRequest(t *testing.T) {
	// Create an original request
	originalReq := &gemini.GenerateContentRequest{
		Contents: []gemini.Content{
			{
				Role: "user",
				Parts: []gemini.Part{
					{Text: "Hello, world!"},
				},
			},
		},
		SystemInstruction_: &gemini.SystemInstruction{
			Role: "system",
			Parts: []gemini.Part{
				{Text: "Test instruction"},
			},
		},
		GenerationConfig: &gemini.GenerationConfig{
			Temperature: 0.7,
		},
	}
	
	partialText := "This is a partial response"
	
	// Build the retry request
	retryReq := BuildRetryRequest(originalReq, partialText)
	
	// Check that the retry request has the correct structure
	if len(retryReq.Contents) != 3 { // original + model response + retry prompt
		t.Errorf("Expected 3 contents, got %d", len(retryReq.Contents))
	}
	
	// Check that the original content is preserved
	if retryReq.Contents[0].Role != "user" || retryReq.Contents[0].Parts[0].Text != "Hello, world!" {
		t.Error("Expected original user content to be preserved")
	}
	
	// Check that the model response was added
	if retryReq.Contents[1].Role != "model" || retryReq.Contents[1].Parts[0].Text != partialText {
		t.Error("Expected model response to be added")
	}
	
	// Check that the retry prompt was added
	if retryReq.Contents[2].Role != "user" || retryReq.Contents[2].Parts[0].Text != gemini.RetryPrompt {
		t.Error("Expected retry prompt to be added")
	}
	
	// Check that system instruction is preserved
	sysInst := retryReq.GetSystemInstruction()
	if sysInst == nil || len(sysInst.Parts) == 0 || sysInst.Parts[0].Text != "Test instruction" {
		t.Error("Expected system instruction to be preserved")
	}
	
	// Check that generation config is preserved
	if retryReq.GenerationConfig == nil || retryReq.GenerationConfig.Temperature != 0.7 {
		t.Error("Expected generation config to be preserved")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}