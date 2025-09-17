package configs

import (
	"net/http"
	"strings"
)

// ServerConfig contains server configuration
type ServerConfig struct {
	Address  string  `yaml:"address" default:":12330"`
	LogLevel string  `yaml:"log_level" default:"info"`
	Routes   []Route `yaml:"routes"`
}

// Route represents a single route configuration
type Route struct {
	Path           string            `yaml:"path"`
	Method         string            `yaml:"method,omitempty"`
	ResponseBody   string            `yaml:"response_body,omitempty"`
	ResponseHeader map[string]string `yaml:"response_header,omitempty"`
	ResponseStatus int               `yaml:"response_status,omitempty"`
	ResponseDump   bool              `yaml:"response_dump,omitempty"`
	Conditions     []RouteCondition  `yaml:"conditions,omitempty"`
}

// RouteCondition represents a conditional response based on header matching
type RouteCondition struct {
	HeaderMatch    map[string]string `yaml:"header_match"`
	ResponseBody   string            `yaml:"response_body,omitempty"`
	ResponseHeader map[string]string `yaml:"response_header,omitempty"`
	ResponseStatus int               `yaml:"response_status,omitempty"`
}

// GetMethod returns the HTTP method for the route, defaulting to GET
func (r *Route) GetMethod() string {
	if r.Method == "" {
		return http.MethodGet
	}
	return r.Method
}

// GetResponseBody returns the response body, defaulting to empty string
func (r *Route) GetResponseBody() string {
	return r.ResponseBody
}

// GetResponseHeaders returns the response headers, defaulting to empty map
func (r *Route) GetResponseHeaders() map[string]string {
	if r.ResponseHeader == nil {
		return make(map[string]string)
	}
	return r.ResponseHeader
}

// GetResponseStatus returns the response status code, defaulting to 200
func (r *Route) GetResponseStatus() int {
	if r.ResponseStatus == 0 {
		return 200
	}
	return r.ResponseStatus
}

// GetResponseDump returns whether request headers and query parameters should be
// included in the response body for debugging purposes. Defaults to false.
func (r *Route) GetResponseDump() bool {
	return r.ResponseDump
}

// GetResponseBody returns the response body for a condition, defaulting to empty string
func (c *RouteCondition) GetResponseBody() string {
	return c.ResponseBody
}

// GetResponseHeaders returns the response headers for a condition, defaulting to empty map
func (c *RouteCondition) GetResponseHeaders() map[string]string {
	if c.ResponseHeader == nil {
		return make(map[string]string)
	}
	return c.ResponseHeader
}

// GetResponseStatus returns the response status code for a condition, defaulting to 200
func (c *RouteCondition) GetResponseStatus() int {
	if c.ResponseStatus == 0 {
		return 200
	}
	return c.ResponseStatus
}

// MatchesHeaders checks if the condition's header requirements match the request headers
func (c *RouteCondition) MatchesHeaders(requestHeaders map[string]string) bool {
	for expectedKey, expectedValue := range c.HeaderMatch {
		found := false
		for actualKey, actualValue := range requestHeaders {
			// Case-insensitive header name comparison
			if strings.EqualFold(expectedKey, actualKey) && actualValue == expectedValue {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// GetLogLevel returns the log level, defaulting to "info"
func (s *ServerConfig) GetLogLevel() string {
	if s.LogLevel == "" {
		return "info"
	}
	return strings.ToLower(s.LogLevel)
}
