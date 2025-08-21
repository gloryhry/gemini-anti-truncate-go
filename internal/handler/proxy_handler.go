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
	"strings"

	"github.com/gorilla/mux"
)

// extractBaseModelName extracts the base model name from a potentially suffixed model path.
// For example, "gemini-2.5-pro:streamGenerateContent" becomes "gemini-2.5-pro".
func extractBaseModelName(modelPath string) string {
	// Split on ':' and take the first part as the base model name
	parts := strings.Split(modelPath, ":")
	if len(parts) > 0 {
		return parts[0]
	}
	return modelPath
}

// ProxyHandler is the main entry point for all incoming API requests.
// It validates the request, decides whether to apply anti-truncate logic,
// and then dispatches to the appropriate stream or non-stream handler.
func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	util.Debugf("Received request for: %s", r.URL.Path)

	// 1. Basic request validation
	if r.Method != http.MethodPost {
		util.SendJSONError(w, "Method not allowed. Please use POST.", http.StatusMethodNotAllowed)
		return
	}

	apiKey := util.GetAPIKey(r)
	if apiKey == "" {
		util.SendJSONError(w, "API key is missing. Please provide it in 'Authorization: Bearer <key>' or 'X-Goog-Api-Key: <key>' header.", http.StatusUnauthorized)
		return
	}

	// 2. Decode the request body
	var req gemini.GenerateContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendJSONError(w, fmt.Sprintf("Invalid JSON in request body: %v", err), http.StatusBadRequest)
		return
	}

	// 3. Determine if this request should be a simple passthrough
	vars := mux.Vars(r)
	modelPath := vars["model"]
	model := extractBaseModelName(modelPath)
	isTargetModel := false
	for _, m := range gemini.TargetModels {
		if m == model {
			isTargetModel = true
			break
		}
	}

	hasSchema := req.GenerationConfig != nil && req.GenerationConfig.ResponseSchema != nil
	
	// Passthrough if not a target model or if it's a structured output request
	if !isTargetModel || hasSchema {
		util.Debugf("Passthrough request for model '%s' (isTarget: %t, hasSchema: %t)", model, isTargetModel, hasSchema)
		passthroughRequest(w, r, apiKey, &req)
		return
	}

	// 4. Inject the finish token prompt
	modifiedReq := proxy.InjectFinishToken(&req)

	// 5. Dispatch to the appropriate handler
	isStream := strings.Contains(r.URL.Path, ":streamGenerateContent")
	if isStream {
		util.Debugf("Dispatching to stream handler")
		HandleStream(w, r, modifiedReq, apiKey)
	} else {
		util.Debugf("Dispatching to non-stream handler")
		HandleNonStream(w, r, modifiedReq, apiKey)
	}
}

// passthroughRequest forwards the request directly to the upstream without modification.
func passthroughRequest(w http.ResponseWriter, r *http.Request, apiKey string, req *gemini.GenerateContentRequest) {
	httpClient := &http.Client{}

	reqBodyBytes, err := json.Marshal(req)
	if err != nil {
		util.SendJSONError(w, "Failed to re-marshal request body for passthrough", http.StatusInternalServerError)
		return
	}

	upstreamURL := fmt.Sprintf("%s/%s", config.AppConfig.UpstreamURLBase, r.URL.Path)
	upstreamReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, upstreamURL, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		util.SendJSONError(w, "Failed to create passthrough upstream request", http.StatusInternalServerError)
		return
	}

	// Copy headers
	upstreamReq.Header.Set("Content-Type", "application/json")
	upstreamReq.Header.Set("X-Goog-Api-Key", apiKey)
	if auth := r.Header.Get("Authorization"); auth != "" {
		upstreamReq.Header.Set("Authorization", auth)
	}

	upstreamResp, err := httpClient.Do(upstreamReq)
	if err != nil {
		util.SendJSONError(w, fmt.Sprintf("Passthrough request failed: %v", err), http.StatusBadGateway)
		return
	}
	defer upstreamResp.Body.Close()

	// Copy upstream response headers to the client
	for key, values := range upstreamResp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(upstreamResp.StatusCode)

	// Stream the response body directly to the client
	if _, err := io.Copy(w, upstreamResp.Body); err != nil {
		util.Errorf("Error streaming passthrough response: %v", err)
	}
}
