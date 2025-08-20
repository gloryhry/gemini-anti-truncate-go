package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gemini-anti-truncate-go/internal/config"
	"gemini-anti-truncate-go/internal/gemini"
	"gemini-anti-truncate-go/internal/proxy"
	"gemini-anti-truncate-go/internal/util"
	"io"
	"net/http"
)

// HandleNonStream manages non-streaming requests, including the retry logic for truncated responses.
func HandleNonStream(w http.ResponseWriter, r *http.Request, initialReq *gemini.GenerateContentRequest, apiKey string) {
	currentReq := initialReq
	httpClient := &http.Client{}

	for i := 0; i < config.AppConfig.MaxRetries; i++ {
		util.Debugf("Non-stream attempt %d/%d", i+1, config.AppConfig.MaxRetries)

		// 1. Prepare the upstream request
		reqBodyBytes, err := json.Marshal(currentReq)
		if err != nil {
			util.SendJSONError(w, "Failed to marshal request body", http.StatusInternalServerError)
			return
		}

		upstreamURL := fmt.Sprintf("%s/%s", config.AppConfig.UpstreamURLBase, r.URL.Path)
		upstreamReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, upstreamURL, bytes.NewBuffer(reqBodyBytes))
		if err != nil {
			util.SendJSONError(w, "Failed to create upstream request", http.StatusInternalServerError)
			return
		}

		// Copy essential headers
		upstreamReq.Header.Set("Content-Type", "application/json")
		upstreamReq.Header.Set("X-Goog-Api-Key", apiKey)
		if auth := r.Header.Get("Authorization"); auth != "" {
			upstreamReq.Header.Set("Authorization", auth)
		}

		// 2. Execute the request
		upstreamResp, err := httpClient.Do(upstreamReq)
		if err != nil {
			util.SendJSONError(w, fmt.Sprintf("Upstream request failed: %v", err), http.StatusBadGateway)
			return
		}
		defer upstreamResp.Body.Close()

		// 3. Read the response body
		respBodyBytes, err := io.ReadAll(upstreamResp.Body)
		if err != nil {
			util.SendJSONError(w, "Failed to read upstream response body", http.StatusBadGateway)
			return
		}

		// 4. Check the status and process the response
		if upstreamResp.StatusCode != http.StatusOK {
			// Check for retryable status codes
			isRetryable := false
			for _, code := range gemini.RetryableStatus {
				if upstreamResp.StatusCode == code {
					isRetryable = true
					break
				}
			}
			if isRetryable {
				util.Debugf("Received retryable status %d, retrying...", upstreamResp.StatusCode)
				continue // Go to the next attempt
			}

			// For non-retryable errors, forward the response to the client
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(upstreamResp.StatusCode)
			w.Write(respBodyBytes)
			return
		}

		// 5. Process the successful response
		result, err := proxy.ProcessNonStream(respBodyBytes)
		if err != nil {
			if pErr, ok := err.(*gemini.ProxyError); ok {
				util.SendJSONError(w, pErr.Message, pErr.StatusCode)
			} else {
				util.SendJSONError(w, "Failed to process non-stream response", http.StatusInternalServerError)
			}
			return
		}

		if result.IsComplete || result.HasFunctionCall {
			util.Debugf("Non-stream response is complete or has function call. Finishing.")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(result.FinalResponseJSON))
			return // Success
		}

		// 6. If not complete, prepare for retry
		util.Debugf("Response incomplete, preparing for retry...")
		currentReq = proxy.BuildRetryRequest(initialReq, result.AccumulatedText)
	}

	// If the loop finishes, we've exceeded max retries
	util.Errorf("Non-stream request failed after %d attempts", config.AppConfig.MaxRetries)
	util.SendJSONError(w, "Request failed after maximum retries", http.StatusGatewayTimeout)
}
