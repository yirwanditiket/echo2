package configs

import (
	"net/http"
	"testing"
)

func TestRoute_GetMethod(t *testing.T) {
	tests := []struct {
		name     string
		route    Route
		expected string
	}{
		{
			name:     "empty method should default to GET",
			route:    Route{Method: ""},
			expected: http.MethodGet,
		},
		{
			name:     "explicit GET method",
			route:    Route{Method: "GET"},
			expected: "GET",
		},
		{
			name:     "POST method",
			route:    Route{Method: "POST"},
			expected: "POST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.route.GetMethod(); got != tt.expected {
				t.Errorf("Route.GetMethod() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRoute_GetResponseBody(t *testing.T) {
	tests := []struct {
		name     string
		route    Route
		expected string
	}{
		{
			name:     "empty response body",
			route:    Route{ResponseBody: ""},
			expected: "",
		},
		{
			name:     "non-empty response body",
			route:    Route{ResponseBody: "Hello World"},
			expected: "Hello World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.route.GetResponseBody(); got != tt.expected {
				t.Errorf("Route.GetResponseBody() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRoute_GetResponseHeaders(t *testing.T) {
	tests := []struct {
		name     string
		route    Route
		expected map[string]string
	}{
		{
			name:     "nil headers should return empty map",
			route:    Route{ResponseHeader: nil},
			expected: make(map[string]string),
		},
		{
			name:     "empty headers map",
			route:    Route{ResponseHeader: make(map[string]string)},
			expected: make(map[string]string),
		},
		{
			name: "headers with content",
			route: Route{ResponseHeader: map[string]string{
				"Content-Type": "application/json",
				"X-Custom":     "value",
			}},
			expected: map[string]string{
				"Content-Type": "application/json",
				"X-Custom":     "value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.route.GetResponseHeaders()
			if len(got) != len(tt.expected) {
				t.Errorf("Route.GetResponseHeaders() length = %v, want %v", len(got), len(tt.expected))
				return
			}

			for key, expectedValue := range tt.expected {
				if gotValue, exists := got[key]; !exists || gotValue != expectedValue {
					t.Errorf("Route.GetResponseHeaders()[%s] = %v, want %v", key, gotValue, expectedValue)
				}
			}
		})
	}
}

func TestRoute_GetResponseStatus(t *testing.T) {
	tests := []struct {
		name     string
		route    Route
		expected int
	}{
		{
			name:     "zero status should default to 200",
			route:    Route{ResponseStatus: 0},
			expected: 200,
		},
		{
			name:     "explicit 404 status",
			route:    Route{ResponseStatus: 404},
			expected: 404,
		},
		{
			name:     "500 status",
			route:    Route{ResponseStatus: 500},
			expected: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.route.GetResponseStatus(); got != tt.expected {
				t.Errorf("Route.GetResponseStatus() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRouteCondition_GetResponseBody(t *testing.T) {
	tests := []struct {
		name      string
		condition RouteCondition
		expected  string
	}{
		{
			name:      "empty response body",
			condition: RouteCondition{ResponseBody: ""},
			expected:  "",
		},
		{
			name:      "non-empty response body",
			condition: RouteCondition{ResponseBody: "Authorized"},
			expected:  "Authorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.condition.GetResponseBody(); got != tt.expected {
				t.Errorf("RouteCondition.GetResponseBody() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRouteCondition_GetResponseStatus(t *testing.T) {
	tests := []struct {
		name      string
		condition RouteCondition
		expected  int
	}{
		{
			name:      "zero status should default to 200",
			condition: RouteCondition{ResponseStatus: 0},
			expected:  200,
		},
		{
			name:      "explicit 201 status",
			condition: RouteCondition{ResponseStatus: 201},
			expected:  201,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.condition.GetResponseStatus(); got != tt.expected {
				t.Errorf("RouteCondition.GetResponseStatus() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRouteCondition_MatchesHeaders(t *testing.T) {
	tests := []struct {
		name           string
		condition      RouteCondition
		requestHeaders map[string]string
		expected       bool
	}{
		{
			name: "single header match",
			condition: RouteCondition{
				HeaderMatch: map[string]string{
					"Authorization": "Bearer token123",
				},
			},
			requestHeaders: map[string]string{
				"Authorization": "Bearer token123",
				"Content-Type":  "application/json",
			},
			expected: true,
		},
		{
			name: "single header mismatch",
			condition: RouteCondition{
				HeaderMatch: map[string]string{
					"Authorization": "Bearer token123",
				},
			},
			requestHeaders: map[string]string{
				"Authorization": "Bearer different-token",
				"Content-Type":  "application/json",
			},
			expected: false,
		},
		{
			name: "missing required header",
			condition: RouteCondition{
				HeaderMatch: map[string]string{
					"X-API-Key": "secret",
				},
			},
			requestHeaders: map[string]string{
				"Authorization": "Bearer token123",
			},
			expected: false,
		},
		{
			name: "multiple headers match",
			condition: RouteCondition{
				HeaderMatch: map[string]string{
					"Authorization": "Bearer token123",
					"X-API-Key":     "secret",
				},
			},
			requestHeaders: map[string]string{
				"Authorization": "Bearer token123",
				"X-API-Key":     "secret",
				"Content-Type":  "application/json",
			},
			expected: true,
		},
		{
			name: "multiple headers partial match",
			condition: RouteCondition{
				HeaderMatch: map[string]string{
					"Authorization": "Bearer token123",
					"X-API-Key":     "secret",
				},
			},
			requestHeaders: map[string]string{
				"Authorization": "Bearer token123",
				"X-API-Key":     "different-secret",
			},
			expected: false,
		},
		{
			name: "empty condition matches any headers",
			condition: RouteCondition{
				HeaderMatch: map[string]string{},
			},
			requestHeaders: map[string]string{
				"Authorization": "Bearer token123",
			},
			expected: true,
		},
		{
			name: "case insensitive header name match",
			condition: RouteCondition{
				HeaderMatch: map[string]string{
					"X-API-Key": "secret",
				},
			},
			requestHeaders: map[string]string{
				"x-api-key": "secret", // lowercase header name
			},
			expected: true,
		},
		{
			name: "case insensitive header name mismatch (value case sensitive)",
			condition: RouteCondition{
				HeaderMatch: map[string]string{
					"X-API-Key": "Secret", // capital S
				},
			},
			requestHeaders: map[string]string{
				"x-api-key": "secret", // lowercase s
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.condition.MatchesHeaders(tt.requestHeaders); got != tt.expected {
				t.Errorf("RouteCondition.MatchesHeaders() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestServerConfig_GetLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		config   ServerConfig
		expected string
	}{
		{
			name:     "empty log level should default to info",
			config:   ServerConfig{LogLevel: ""},
			expected: "info",
		},
		{
			name:     "debug log level",
			config:   ServerConfig{LogLevel: "debug"},
			expected: "debug",
		},
		{
			name:     "info log level",
			config:   ServerConfig{LogLevel: "info"},
			expected: "info",
		},
		{
			name:     "warn log level",
			config:   ServerConfig{LogLevel: "warn"},
			expected: "warn",
		},
		{
			name:     "error log level",
			config:   ServerConfig{LogLevel: "error"},
			expected: "error",
		},
		{
			name:     "uppercase log level should be converted to lowercase",
			config:   ServerConfig{LogLevel: "DEBUG"},
			expected: "debug",
		},
		{
			name:     "mixed case log level should be converted to lowercase",
			config:   ServerConfig{LogLevel: "WaRn"},
			expected: "warn",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.GetLogLevel(); got != tt.expected {
				t.Errorf("ServerConfig.GetLogLevel() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRoute_GetResponseDump(t *testing.T) {
	tests := []struct {
		name     string
		route    Route
		expected bool
	}{
		{
			name:     "default response dump should be false",
			route:    Route{ResponseDump: false},
			expected: false,
		},
		{
			name:     "explicit false response dump",
			route:    Route{ResponseDump: false},
			expected: false,
		},
		{
			name:     "explicit true response dump",
			route:    Route{ResponseDump: true},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.route.GetResponseDump(); got != tt.expected {
				t.Errorf("Route.GetResponseDump() = %v, want %v", got, tt.expected)
			}
		})
	}
}
