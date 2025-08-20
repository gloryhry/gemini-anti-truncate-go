package util

import (
	"gemini-anti-truncate-go/internal/config"
	"log"
)

// Infof logs an informational message.
func Infof(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

// Debugf logs a debug message only if DebugMode is enabled.
func Debugf(format string, v ...interface{}) {
	if config.AppConfig != nil && config.AppConfig.DebugMode {
		log.Printf("[DEBUG] "+format, v...)
	}
}

// Errorf logs an error message.
func Errorf(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}
