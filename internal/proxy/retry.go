package proxy

import "gemini-anti-truncate-go/internal/gemini"

// BuildRetryRequest creates a new request to continue a truncated generation.
// It takes the original request and the partial response text, and constructs a new request
// that instructs the model to continue from where it left off.
func BuildRetryRequest(originalReq *gemini.GenerateContentRequest, partialResponseText string) *gemini.GenerateContentRequest {
	// Create a new request object to avoid modifying the original.
	retryReq := &gemini.GenerateContentRequest{
		// Copy essential fields from the original request.
		SystemInstruction:  originalReq.SystemInstruction,
		SystemInstruction_: originalReq.SystemInstruction_,
		GenerationConfig:   originalReq.GenerationConfig,
		SafetySettings:     originalReq.SafetySettings,
		Tools:              originalReq.Tools,
		// Deep copy the contents to avoid slice modification issues.
		Contents: make([]gemini.Content, len(originalReq.Contents)),
	}
	copy(retryReq.Contents, originalReq.Contents)

	// 1. Add the model's partial response to the conversation history.
	// This provides context for the continuation.
	modelContent := gemini.Content{
		Role: "model",
		Parts: []gemini.Part{
			{Text: partialResponseText},
		},
	}
	retryReq.Contents = append(retryReq.Contents, modelContent)

	// 2. Add a new user prompt instructing the model to continue.
	userPrompt := gemini.Content{
		Role: "user",
		Parts: []gemini.Part{
			{Text: gemini.RetryPrompt},
		},
	}
	retryReq.Contents = append(retryReq.Contents, userPrompt)

	return retryReq
}
