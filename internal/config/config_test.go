package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("UPSTREAM_URL_BASE", "https://test.googleapis.com")
	os.Setenv("MAX_RETRIES", "10")
	os.Setenv("DEBUG_MODE", "true")
	os.Setenv("HTTP_PORT", "9090")
	
	// Load configuration
	Load()
	
	// Test values
	if AppConfig.UpstreamURLBase != "https://test.googleapis.com" {
		t.Errorf("Expected UpstreamURLBase to be 'https://test.googleapis.com', got '%s'", AppConfig.UpstreamURLBase)
	}
	
	if AppConfig.MaxRetries != 10 {
		t.Errorf("Expected MaxRetries to be 10, got %d", AppConfig.MaxRetries)
	}
	
	if !AppConfig.DebugMode {
		t.Error("Expected DebugMode to be true")
	}
	
	if AppConfig.Port != 9090 {
		t.Errorf("Expected Port to be 9090, got %d", AppConfig.Port)
	}
	
	// Clean up environment variables
	os.Unsetenv("UPSTREAM_URL_BASE")
	os.Unsetenv("MAX_RETRIES")
	os.Unsetenv("DEBUG_MODE")
	os.Unsetenv("HTTP_PORT")
}

func TestGetEnv(t *testing.T) {
	// Set a test environment variable
	os.Setenv("TEST_VAR", "test_value")
	
	// Test getting the variable
	value := getEnv("TEST_VAR", "default_value")
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", value)
	}
	
	// Test getting a non-existent variable (should return default)
	defaultValue := getEnv("NON_EXISTENT_VAR", "default_value")
	if defaultValue != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", defaultValue)
	}
	
	// Clean up
	os.Unsetenv("TEST_VAR")
}

func TestGetEnvAsInt(t *testing.T) {
	// Set a test environment variable
	os.Setenv("TEST_INT", "42")
	
	// Test getting the variable
	value := getEnvAsInt("TEST_INT", 10)
	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}
	
	// Test getting a non-existent variable (should return default)
	defaultValue := getEnvAsInt("NON_EXISTENT_INT", 10)
	if defaultValue != 10 {
		t.Errorf("Expected 10, got %d", defaultValue)
	}
	
	// Test getting an invalid integer (should return default)
	os.Setenv("INVALID_INT", "not_a_number")
	invalidValue := getEnvAsInt("INVALID_INT", 10)
	if invalidValue != 10 {
		t.Errorf("Expected 10, got %d", invalidValue)
	}
	
	// Clean up
	os.Unsetenv("TEST_INT")
	os.Unsetenv("INVALID_INT")
}

func TestGetEnvAsBool(t *testing.T) {
	// Test true values
	os.Setenv("TEST_BOOL_TRUE1", "true")
	os.Setenv("TEST_BOOL_TRUE2", "1")
	
	// Test false values
	os.Setenv("TEST_BOOL_FALSE1", "false")
	os.Setenv("TEST_BOOL_FALSE2", "0")
	
	// Test invalid value (should return default)
	os.Setenv("TEST_BOOL_INVALID", "not_a_bool")
	
	// Test true values
	if !getEnvAsBool("TEST_BOOL_TRUE1", false) {
		t.Error("Expected true for 'true'")
	}
	
	if !getEnvAsBool("TEST_BOOL_TRUE2", false) {
		t.Error("Expected true for '1'")
	}
	
	// Test false values
	if getEnvAsBool("TEST_BOOL_FALSE1", true) {
		t.Error("Expected false for 'false'")
	}
	
	if getEnvAsBool("TEST_BOOL_FALSE2", true) {
		t.Error("Expected false for '0'")
	}
	
	// Test invalid value
	if getEnvAsBool("TEST_BOOL_INVALID", true) != true {
		t.Error("Expected default value true for invalid bool")
	}
	
	// Test non-existent variable
	if getEnvAsBool("NON_EXISTENT_BOOL", true) != true {
		t.Error("Expected default value true for non-existent variable")
	}
	
	// Clean up
	os.Unsetenv("TEST_BOOL_TRUE1")
	os.Unsetenv("TEST_BOOL_TRUE2")
	os.Unsetenv("TEST_BOOL_FALSE1")
	os.Unsetenv("TEST_BOOL_FALSE2")
	os.Unsetenv("TEST_BOOL_INVALID")
}