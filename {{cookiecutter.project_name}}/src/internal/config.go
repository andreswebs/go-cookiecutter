// Package config loads application configuration from environment variables.
package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration.
type Config struct {
	// Port is the TCP port the HTTP server listens on.
	Port string
}

// Load reads environment variables and returns a validated Config.
func Load() (Config, error) {
	return Config{
		Port: optionalString("PORT", "8080"),
	}, nil
}

// optionalString returns the value of the environment variable named by key,
// or defaultVal if the variable is unset or empty.
func optionalString(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// requiredString returns the value of the environment variable named by key.
// It returns an error if the variable is unset or empty.
func requiredString(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("required environment variable %s is not set", key)
	}
	return v, nil
}
