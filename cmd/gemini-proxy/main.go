package main

import (
	"fmt"
	"gemini-anti-truncate-go/internal/config"
	"gemini-anti-truncate-go/internal/handler"
	"gemini-anti-truncate-go/internal/util"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Load application configuration from environment variables
	config.Load()

	// Initialize the router
	r := mux.NewRouter()

	// The primary route that captures all relevant Gemini API paths.
	// This single route will handle both stream and non-stream requests,
	// which are then differentiated within the ProxyHandler.
	r.HandleFunc("/v1beta/models/{model:.+}", handler.ProxyHandler).Methods("POST")

	// Define the server address
	addr := fmt.Sprintf(":%d", config.AppConfig.Port)

	util.Infof("Starting Gemini Anti-Truncate Proxy Server on %s", addr)
	util.Infof("Debug mode is %t", config.AppConfig.DebugMode)

	// Start the HTTP server
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
