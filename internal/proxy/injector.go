package proxy

import "gemini-anti-truncate-go/internal/gemini"

// InjectFinishToken modifies the request to include instructions for the model
// to append a finish token at the end of its response.
func InjectFinishToken(req *gemini.GenerateContentRequest) *gemini.GenerateContentRequest {
	// 1. Handle System Instruction
	systemInstruction := req.GetSystemInstruction()
	if systemInstruction == nil {
		// If no system instruction exists, create one.
		systemInstruction = &gemini.SystemInstruction{
			Role: "system",
			Parts: []gemini.Part{
				{Text: "You are a helpful assistant. Please ensure your response ends with " + gemini.FinishToken},
			},
		}
	} else {
		// If it exists, append the instruction.
		// We assume the first part is the main text.
		if len(systemInstruction.Parts) > 0 {
			systemInstruction.Parts[0].Text += "\n\nPlease ensure your response ends with " + gemini.FinishToken
		} else {
			// If parts are empty, add a new part.
			systemInstruction.Parts = append(systemInstruction.Parts, gemini.Part{Text: "Please ensure your response ends with " + gemini.FinishToken})
		}
	}
	req.SetSystemInstruction(systemInstruction)

	// 2. Handle User Prompt Suffix
	// Find the last user content and append the reminder.
	if len(req.Contents) > 0 {
		lastContentIndex := len(req.Contents) - 1
		lastContent := &req.Contents[lastContentIndex]

		if lastContent.Role == "user" && len(lastContent.Parts) > 0 {
			lastPartIndex := len(lastContent.Parts) - 1
			// Append to the last part of the last user message.
			lastContent.Parts[lastPartIndex].Text += gemini.UserPromptSuffix
		}
	}

	return req
}
