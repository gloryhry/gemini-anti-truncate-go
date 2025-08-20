package config

import (
	"gemini-anti-truncate-go/internal/gemini"
	"os"
	"strconv"
)

// Config holds all configuration for the application.
type Config struct {
	UpstreamURLBase string
	MaxRetries      int
	DebugMode       bool
	Port            int
}

// AppConfig is a global variable holding the application's configuration.
var AppConfig *Config

// Load loads configuration from environment variables and populates the AppConfig global variable.
func Load() {
	AppConfig = &Config{
		UpstreamURLBase: getEnv("UPSTREAM_URL_BASE", gemini.DefaultUpstreamURL),
		MaxRetries:      getEnvAsInt("MAX_RETRIES", gemini.DefaultMaxRetries),
		DebugMode:       getEnvAsBool("DEBUG_MODE", false),
		Port:            getEnvAsInt("HTTP_PORT", gemini.DefaultHTTPPort),
	}
}

// getEnv retrieves a string value from an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an integer value from an environment variable or returns a default value.
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsBool retrieves a boolean value from an environment variable or returns a default value.
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
