package proxy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"gemini-anti-truncate-go/internal/gemini"
	"gemini-anti-truncate-go/internal/util"
	"io"
	"net/http"
	"strings"
)

// StreamProcessingResult holds the outcome of processing a stream.
type StreamProcessingResult struct {
	IsComplete        bool
	HasFunctionCall   bool
	AccumulatedText   string
	FinalResponseJSON string // Used for non-stream handler to get the full JSON
}

// ProcessStream handles the server-sent event (SSE) stream from the upstream API.
// It forwards events to the client, while checking for the finish token to determine if the response is complete.
func ProcessStream(w http.ResponseWriter, upstreamResp *http.Response) (*StreamProcessingResult, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, &gemini.ProxyError{Message: "Streaming unsupported", StatusCode: http.StatusInternalServerError}
	}

	// Set client headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	scanner := bufio.NewScanner(upstreamResp.Body)
	var textBuffer, lookbehindBuffer bytes.Buffer
	var isComplete, hasFunctionCall, inPassthrough bool

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "data:") {
			jsonData := strings.TrimPrefix(line, "data: ")

			if inPassthrough {
				// Once a function call is seen, we just forward everything.
				fmt.Fprintf(w, "%s\n", line)
				flusher.Flush()
				continue
			}

			var streamChunk gemini.GenerateContentResponse
			if err := json.Unmarshal([]byte(jsonData), &streamChunk); err != nil {
				util.Debugf("Error unmarshalling stream chunk: %v. Data: %s", err, jsonData)
				// Forward malformed data as-is
				fmt.Fprintf(w, "%s\n", line)
				flusher.Flush()
				continue
			}

			var currentText, currentThought string
			isThoughtChunk := false

			// Process parts to extract text, thoughts, and function calls
			if len(streamChunk.Candidates) > 0 {
				for _, part := range streamChunk.Candidates[0].Content.Parts {
					if part.Thought {
						isThoughtChunk = true
						currentThought += part.Text
					} else if part.FunctionCall != nil {
						hasFunctionCall = true
						inPassthrough = true // Enter passthrough mode
						currentText += part.Text
					} else {
						currentText += part.Text
					}
				}
			}

			if isThoughtChunk {
				// Forward thought chunks immediately without affecting the main text buffer
				// but clean the finish token just in case it appears there.
				cleanedThoughtLine := strings.Replace(line, gemini.FinishToken, "", -1)
				fmt.Fprintf(w, "%s\n", cleanedThoughtLine)
				flusher.Flush()
				continue // Skip buffer processing for thoughts
			}

			// Append text to main buffer and lookbehind buffer
			textBuffer.WriteString(currentText)
			if lookbehindBuffer.Len() > gemini.TokenLookbehindChars {
				// Slide the lookbehind window
				lookbehindBuffer.Next(lookbehindBuffer.Len() - gemini.TokenLookbehindChars)
			}
			lookbehindBuffer.WriteString(currentText)

			// Check for finish token
			if strings.Contains(lookbehindBuffer.String(), gemini.FinishToken) {
				isComplete = true
				// Clean the token from the output
				line = strings.Replace(line, gemini.FinishToken, "", -1)
			}
		}

		// Forward the (potentially modified) line to the client
		fmt.Fprintf(w, "%s\n", line)

		// Write an empty line to signify the end of the event
		if len(line) > 0 {
			fmt.Fprintf(w, "\n")
		}

		flusher.Flush()
	}

	if err := scanner.Err(); err != nil {
		util.Errorf("Error reading stream from upstream: %v", err)
		return nil, err
	}

	return &StreamProcessingResult{
		IsComplete:      isComplete,
		HasFunctionCall: hasFunctionCall,
		AccumulatedText: textBuffer.String(),
	}, nil
}

// ProcessNonStream checks a complete non-streaming response for the finish token.
func ProcessNonStream(body []byte) (*StreamProcessingResult, error) {
	var response gemini.GenerateContentResponse
	if err := json.Unmarshal(body, &response); err != nil {
		util.Errorf("Error unmarshalling non-stream response: %v", err)
		return nil, &gemini.ProxyError{Message: "Failed to parse upstream response", StatusCode: http.StatusBadGateway}
	}

	var accumulatedText string
	hasFunctionCall := false
	if len(response.Candidates) > 0 {
		for _, part := range response.Candidates[0].Content.Parts {
			if part.FunctionCall != nil {
				hasFunctionCall = true
			}
			accumulatedText += part.Text
		}
	}

	isComplete := strings.Contains(accumulatedText, gemini.FinishToken)
	
	// Clean the finish token from the final text for the client
	finalText := strings.Replace(accumulatedText, gemini.FinishToken, "", -1)

	// Re-assemble the response with the cleaned text
	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		// This is a simplification. It assumes the text is in the first part.
		// A more robust implementation might need to handle multiple text parts.
		response.Candidates[0].Content.Parts[0].Text = finalText
	}

	finalJSON, err := json.Marshal(response)
	if err != nil {
		util.Errorf("Error re-marshalling cleaned response: %v", err)
		return nil, &gemini.ProxyError{Message: "Failed to construct final response", StatusCode: http.StatusInternalServerError}
	}


	return &StreamProcessingResult{
		IsComplete:        isComplete,
		HasFunctionCall:   hasFunctionCall,
		AccumulatedText:   accumulatedText, // The original text with the token
		FinalResponseJSON: string(finalJSON),
	}, nil
}

// Custom error for proxy-specific issues
type ProxyError struct {
	Message    string
	StatusCode int
}

func (e *ProxyError) Error() string {
	return e.Message
}
