package configs

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads server configuration from a YAML file
func LoadConfig(filePath string) (*ServerConfig, error) {
	// Read the YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse the YAML content
	var config ServerConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate the configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// validateConfig validates the server configuration
func validateConfig(config *ServerConfig) error {
	if config.Address == "" {
		config.Address = ":12330" // Set default address
	}

	// Validate routes
	for i, route := range config.Routes {
		if route.Path == "" {
			return fmt.Errorf("route %d: path cannot be empty", i)
		}

		// Validate HTTP method if provided
		if route.Method != "" {
			validMethods := map[string]bool{
				"GET": true, "POST": true, "PUT": true, "DELETE": true,
				"PATCH": true, "HEAD": true, "OPTIONS": true,
			}
			if !validMethods[route.Method] {
				return fmt.Errorf("route %d: invalid HTTP method '%s'", i, route.Method)
			}
		}
	}

	return nil
}
