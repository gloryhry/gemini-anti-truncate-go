package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gemini-anti-truncate-go/internal/config"
	"gemini-anti-truncate-go/internal/gemini"
	"gemini-anti-truncate-go/internal/proxy"
	"gemini-anti-truncate-go/internal/util"
	"net/http"
)

// headerSuppressingWriter is a wrapper around http.ResponseWriter that suppresses
// subsequent attempts to write headers after the first write. This is crucial for
// the streaming retry logic, where multiple upstream responses are stitched together
// into a single client response stream.
type headerSuppressingWriter struct {
	http.ResponseWriter
	headersSent bool
}

// Header returns the original ResponseWriter's header map if headers have not been sent,
// otherwise it returns a dummy map to prevent panics on subsequent writes.
func (hsw *headerSuppressingWriter) Header() http.Header {
	if hsw.headersSent {
		return http.Header{}
	}
	return hsw.ResponseWriter.Header()
}

// Write delegates to the original ResponseWriter's Write method and marks headers as sent.
func (hsw *headerSuppressingWriter) Write(p []byte) (int, error) {
	hsw.headersSent = true
	return hsw.ResponseWriter.Write(p)
}

// WriteHeader delegates to the original ResponseWriter's WriteHeader method only if
// headers have not been sent yet.
func (hsw *headerSuppressingWriter) WriteHeader(statusCode int) {
	if hsw.headersSent {
		return
	}
	hsw.headersSent = true
	hsw.ResponseWriter.WriteHeader(statusCode)
}

// Flush ensures that the underlying writer is flushed if it implements http.Flusher.
func (hsw *headerSuppressingWriter) Flush() {
	if flusher, ok := hsw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// HandleStream manages streaming requests, including the retry logic for truncated streams.
func HandleStream(w http.ResponseWriter, r *http.Request, initialReq *gemini.GenerateContentRequest, apiKey string) {
	currentReq := initialReq
	httpClient := &http.Client{}
	var accumulatedText string

	// Wrap the original response writer to handle headers correctly across multiple retries.
	wrappedWriter := &headerSuppressingWriter{ResponseWriter: w}

	for i := 0; i < config.AppConfig.MaxRetries; i++ {
		util.Debugf("Stream attempt %d/%d", i+1, config.AppConfig.MaxRetries)

		reqBodyBytes, err := json.Marshal(currentReq)
		if err != nil {
			// Cannot send JSON error if stream has started, so just log and exit.
			if !wrappedWriter.headersSent {
				util.SendJSONError(w, "Failed to marshal request body", http.StatusInternalServerError)
			}
			util.Errorf("Failed to marshal request body: %v", err)
			return
		}

		upstreamURL := fmt.Sprintf("%s/%s", config.AppConfig.UpstreamURLBase, r.URL.Path)
		upstreamReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, upstreamURL, bytes.NewBuffer(reqBodyBytes))
		if err != nil {
			if !wrappedWriter.headersSent {
				util.SendJSONError(w, "Failed to create upstream request", http.StatusInternalServerError)
			}
			util.Errorf("Failed to create upstream request: %v", err)
			return
		}

		upstreamReq.Header.Set("Content-Type", "application/json")
		upstreamReq.Header.Set("X-Goog-Api-Key", apiKey)
		if auth := r.Header.Get("Authorization"); auth != "" {
			upstreamReq.Header.Set("Authorization", auth)
		}

		upstreamResp, err := httpClient.Do(upstreamReq)
		if err != nil {
			if !wrappedWriter.headersSent {
				util.SendJSONError(w, fmt.Sprintf("Upstream request failed: %v", err), http.StatusBadGateway)
			}
			util.Errorf("Upstream request failed: %v", err)
			return
		}
		defer upstreamResp.Body.Close()

		if upstreamResp.StatusCode != http.StatusOK {
			// Handle retryable status codes
			isRetryable := false
			for _, code := range gemini.RetryableStatus {
				if upstreamResp.StatusCode == code {
					isRetryable = true
					break
				}
			}
			if isRetryable {
				util.Debugf("Received retryable status %d, retrying stream...", upstreamResp.StatusCode)
				continue // Go to the next attempt
			}
			
			// For non-retryable errors, forward if possible
			if !wrappedWriter.headersSent {
				util.SendJSONError(w, "Upstream returned non-200 status", upstreamResp.StatusCode)
			}
			util.Errorf("Upstream returned non-200 status: %d", upstreamResp.StatusCode)
			return
		}

		// Process the stream. The wrappedWriter ensures headers are only sent once.
		result, err := proxy.ProcessStream(wrappedWriter, upstreamResp)
		if err != nil {
			util.Errorf("Error processing stream: %v", err)
			return // The connection is likely broken
		}

		// Append the text from this chunk to the total accumulated text
		accumulatedText += result.AccumulatedText

		if result.IsComplete || result.HasFunctionCall {
			util.Debugf("Stream is complete or has function call. Finishing.")
			return // Success
		}

		util.Debugf("Stream incomplete, preparing for retry...")
		currentReq = proxy.BuildRetryRequest(initialReq, accumulatedText)
	}

	util.Errorf("Stream request failed after %d attempts", config.AppConfig.MaxRetries)
	// We cannot send a final error message as the stream is already in progress.
	// The client will experience this as a stream that simply stops.
}
