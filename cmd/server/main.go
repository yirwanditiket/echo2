package main

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/yirwanditiket/echo2/configs"
)

// Global logger instance
var logger *slog.Logger

// Global shutdown channel to signal when server is shutting down
var shutdownChan = make(chan struct{})

// setupLogger configures the slog logger with the specified level
func setupLogger(level string) {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn", "warning":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger = slog.New(handler)
	slog.SetDefault(logger)
}

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	config, err := configs.LoadConfig(*configPath)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Setup logger with configured level
	setupLogger(config.GetLogLevel())

	slog.Info("Starting server", "address", config.Address)
	slog.Info("Loaded routes", "count", len(config.Routes))

	// Create the server
	appServer := &Server{config: config}

	// Initialize router with configured routes
	appServer.initializeRouter()

	// Create fasthttp server with router handler
	httpServer := &fasthttp.Server{
		Handler: appServer.router.Handler,
		Name:    "echo-server",
	}

	// Start server in a goroutine
	go func() {
		if err := httpServer.ListenAndServe(config.Address); err != nil {
			slog.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	slog.Info("Received shutdown signal, shutting down gracefully...")

	// Signal all ongoing operations that server is shutting down
	close(shutdownChan)

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := httpServer.ShutdownWithContext(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	} else {
		slog.Info("Server exited gracefully")
	}
}

// Server holds the server configuration and handles requests.
// The server uses fasthttp/router for efficient HTTP routing instead of manual path matching.
// This provides better performance and proper HTTP status code handling.
type Server struct {
	config *configs.ServerConfig // Server configuration loaded from YAML
	router *router.Router        // FastHTTP router for efficient request routing
}

// RequestDump represents the structure for request dump data that is included
// in response bodies when response_dump is enabled in the server configuration.
// This is useful for debugging and understanding what headers and query parameters
// the server receives from clients.
type RequestDump struct {
	Headers         map[string]string `json:"headers"`          // All request headers as key-value pairs
	QueryParameters map[string]string `json:"query_parameters"` // All query parameters as key-value pairs
}

// initializeRouter sets up the fasthttp/router with all configured routes.
// This method creates a new router instance and registers each configured route
// with its specific HTTP method. It provides several benefits over manual routing:
// - Efficient path matching using radix trees
// - Proper HTTP status codes (405 for wrong methods, 404 for missing paths)
// - Support for all HTTP methods including custom ones
// - Better performance for high-traffic scenarios
func (s *Server) initializeRouter() {
	s.router = router.New()

	// Add all configured routes to the router
	for _, route := range s.config.Routes {
		method := strings.ToUpper(route.GetMethod())
		path := route.Path

		// Create a closure to capture the route configuration
		routeConfig := route
		handler := func(ctx *fasthttp.RequestCtx) {
			s.handleRouteRequest(ctx, routeConfig)
		}

		// Register the route with the appropriate HTTP method
		switch method {
		case "GET":
			s.router.GET(path, handler)
		case "POST":
			s.router.POST(path, handler)
		case "PUT":
			s.router.PUT(path, handler)
		case "DELETE":
			s.router.DELETE(path, handler)
		case "PATCH":
			s.router.PATCH(path, handler)
		case "HEAD":
			s.router.HEAD(path, handler)
		case "OPTIONS":
			s.router.OPTIONS(path, handler)
		default:
			// For any other methods, use the ANY method (supports all HTTP methods)
			slog.Warn("Unknown HTTP method, registering as ANY", "method", method, "path", path)
			s.router.ANY(path, handler)
		}

		slog.Debug("Registered route", "method", method, "path", path)
	}

	// Add a catch-all route for 404 handling
	s.router.NotFound = func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetContentType("text/plain")
		ctx.WriteString("404 Not Found")
	}
}

// handleRouteRequest processes a specific route request (used by router).
// This method is called by the fasthttp/router when a route matches an incoming request.
// It serves as an adapter between the router and the existing route processing logic,
// providing access logging and delegating actual response handling to handleRoute.
//
// Parameters:
// - ctx: The fasthttp request context containing request/response data
// - route: The matched route configuration with response details
func (s *Server) handleRouteRequest(ctx *fasthttp.RequestCtx, route configs.Route) {
	path := string(ctx.Path())
	method := string(ctx.Method())

	// Access log at debug level
	slog.Debug("Received request", "method", method, "path", path)

	// Process the matched route
	s.handleRoute(ctx, route)
}

// parseDelayParam extracts and parses the delay parameter from query string
func (s *Server) parseDelayParam(ctx *fasthttp.RequestCtx) (time.Duration, error) {
	delayStr := string(ctx.QueryArgs().Peek("delay"))
	if delayStr == "" {
		return 0, nil
	}

	// Try to parse as duration (e.g., "10ms", "1s", "500us")
	delay, err := time.ParseDuration(delayStr)
	if err != nil {
		// If that fails, try to parse as milliseconds integer (for backward compatibility)
		if ms, parseErr := strconv.Atoi(delayStr); parseErr == nil {
			delay = time.Duration(ms) * time.Millisecond
			return delay, nil
		}
		return 0, err
	}

	return delay, nil
}

// sleepWithCancellation sleeps for the specified duration while checking for shutdown cancellation
// Returns true if the sleep completed normally, false if cancelled due to server shutdown
func (s *Server) sleepWithCancellation(delay time.Duration) bool {
	if delay <= 0 {
		return true
	}

	// Create a timer for the delay
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		// Delay completed successfully
		return true
	case <-shutdownChan:
		// Server is shutting down, return early
		slog.Debug("Request delay cancelled due to server shutdown",
			"remaining_delay", delay.String())
		return false
	}
}

// handleRoute processes a matched route and sends the configured response
func (s *Server) handleRoute(ctx *fasthttp.RequestCtx, route configs.Route) {
	// Parse and apply delay parameter if present
	if delay, err := s.parseDelayParam(ctx); err != nil {
		// Invalid delay parameter, return 400 Bad Request
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetContentType("text/plain")
		ctx.WriteString("Invalid delay parameter: " + err.Error())
		return
	} else if delay > 0 {
		// Apply delay with shutdown cancellation support
		if !s.sleepWithCancellation(delay) {
			// Server is shutting down, return early without sending response
			return
		}
	}

	// Check if any conditions match the request headers
	requestHeaders := s.extractHeaders(ctx)

	var responseBody string
	var responseHeaders map[string]string
	var responseStatus int
	conditionMatched := false

	// Check conditions first
	for _, condition := range route.Conditions {
		if condition.MatchesHeaders(requestHeaders) {
			responseBody = condition.GetResponseBody()
			responseHeaders = condition.GetResponseHeaders()
			responseStatus = condition.GetResponseStatus()
			conditionMatched = true
			slog.Debug("Condition matched", "method", route.GetMethod(), "path", route.Path)
			break
		}
	}

	// If no condition matched, use default route response
	if !conditionMatched {
		responseBody = route.GetResponseBody()
		responseHeaders = route.GetResponseHeaders()
		responseStatus = route.GetResponseStatus()
	}

	// Set response status code
	ctx.SetStatusCode(responseStatus)

	// Set response headers
	for key, value := range responseHeaders {
		ctx.Response.Header.Set(key, value)
	}

	// Set content type if not already set
	if len(ctx.Response.Header.Peek("Content-Type")) == 0 {
		ctx.SetContentType("text/plain")
	}

	// Handle response dump if enabled for this route
	finalResponseBody := responseBody
	if route.GetResponseDump() {
		requestHeaders := s.extractHeaders(ctx)
		queryParams := s.extractQueryParameters(ctx)

		dump := RequestDump{
			Headers:         requestHeaders,
			QueryParameters: queryParams,
		}

		dumpJSON, err := json.MarshalIndent(dump, "", "  ")
		if err != nil {
			slog.Error("Failed to marshal request dump", "error", err)
		} else {
			// Replace response body with JSON dump
			finalResponseBody = string(dumpJSON)
			// Set content type to JSON when dumping
			ctx.Response.Header.Set("Content-Type", "application/json")
		}
	}

	// Set response body
	ctx.WriteString(finalResponseBody)

	slog.Debug("Request handled",
		"method", route.GetMethod(),
		"path", route.Path,
		"status", responseStatus,
		"response_bytes", len(finalResponseBody))
}

// extractHeaders extracts request headers into a map for condition matching
func (s *Server) extractHeaders(ctx *fasthttp.RequestCtx) map[string]string {
	headers := make(map[string]string)

	ctx.Request.Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})

	return headers
}

// extractQueryParameters extracts query parameters into a map for response dumping
func (s *Server) extractQueryParameters(ctx *fasthttp.RequestCtx) map[string]string {
	queryParams := make(map[string]string)

	ctx.QueryArgs().VisitAll(func(key, value []byte) {
		queryParams[string(key)] = string(value)
	})

	return queryParams
}
