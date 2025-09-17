package configs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("valid config file", func(t *testing.T) {
		configContent := `address: ":8080"
routes:
  - path: "/health"
    method: "GET"
    response_body: "OK"
  - path: "/api/users"
    method: "POST"
    response_body: '{"message": "created"}'
    response_header:
      Content-Type: "application/json"
`
		configFile := filepath.Join(tempDir, "valid_config.yaml")
		if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		config, err := LoadConfig(configFile)
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		if config.Address != ":8080" {
			t.Errorf("Expected address :8080, got %s", config.Address)
		}

		if len(config.Routes) != 2 {
			t.Errorf("Expected 2 routes, got %d", len(config.Routes))
		}

		// Check first route
		route1 := config.Routes[0]
		if route1.Path != "/health" {
			t.Errorf("Expected path /health, got %s", route1.Path)
		}
		if route1.Method != "GET" {
			t.Errorf("Expected method GET, got %s", route1.Method)
		}
		if route1.ResponseBody != "OK" {
			t.Errorf("Expected response body OK, got %s", route1.ResponseBody)
		}

		// Check second route
		route2 := config.Routes[1]
		if route2.Path != "/api/users" {
			t.Errorf("Expected path /api/users, got %s", route2.Path)
		}
		if route2.Method != "POST" {
			t.Errorf("Expected method POST, got %s", route2.Method)
		}
		if route2.ResponseHeader["Content-Type"] != "application/json" {
			t.Errorf("Expected Content-Type header application/json, got %s", route2.ResponseHeader["Content-Type"])
		}
	})

	t.Run("config with defaults", func(t *testing.T) {
		configContent := `routes:
  - path: "/test"
`
		configFile := filepath.Join(tempDir, "default_config.yaml")
		if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		config, err := LoadConfig(configFile)
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		// Should use default address
		if config.Address != ":12330" {
			t.Errorf("Expected default address :12330, got %s", config.Address)
		}

		// Route should have defaults applied
		route := config.Routes[0]
		if route.GetMethod() != "GET" {
			t.Errorf("Expected default method GET, got %s", route.GetMethod())
		}
		if route.GetResponseBody() != "" {
			t.Errorf("Expected empty response body, got %s", route.GetResponseBody())
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		_, err := LoadConfig("/non/existent/file.yaml")
		if err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}
	})

	t.Run("invalid YAML", func(t *testing.T) {
		configContent := `invalid: yaml: content: [
`
		configFile := filepath.Join(tempDir, "invalid_config.yaml")
		if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		_, err := LoadConfig(configFile)
		if err == nil {
			t.Error("Expected error for invalid YAML, got nil")
		}
	})

	t.Run("empty path validation", func(t *testing.T) {
		configContent := `routes:
  - path: ""
    method: "GET"
`
		configFile := filepath.Join(tempDir, "empty_path_config.yaml")
		if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		_, err := LoadConfig(configFile)
		if err == nil {
			t.Error("Expected error for empty path, got nil")
		}
	})

	t.Run("invalid HTTP method", func(t *testing.T) {
		configContent := `routes:
  - path: "/test"
    method: "INVALID"
`
		configFile := filepath.Join(tempDir, "invalid_method_config.yaml")
		if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		_, err := LoadConfig(configFile)
		if err == nil {
			t.Error("Expected error for invalid HTTP method, got nil")
		}
	})
}
