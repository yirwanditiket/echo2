# Echo HTTP Server

A high-performance HTTP server built with [Valyala's fasthttp](https://github.com/valyala/fasthttp) framework that serves configurable routes from a YAML configuration file.

## Features

- **High Performance**: Built on top of fasthttp for maximum performance with fasthttp/router for efficient routing
- **YAML Configuration**: Define routes, methods, responses, and headers in a simple YAML file
- **Structured Logging**: Uses Go's `log/slog` with configurable log levels (debug, info, warn, error)
- **Advanced Routing**: Uses fasthttp/router for efficient HTTP method and path-based routing with proper status codes (405 for wrong methods, 404 for missing paths)
- **Flexible Route Configuration**: Support for custom HTTP methods, response bodies, and headers
- **Header-Based Conditional Responses**: Return different responses based on request headers
- **Response Delay Parameter**: Add artificial delays to responses using `?delay=10ms` for testing scenarios with shutdown-aware cancellation support
- **Response Dump**: Include request headers and query parameters in JSON format within the response body for debugging purposes
- **Default Values**: Sensible defaults for method (GET), response body (empty), and headers (empty)
- **Graceful Shutdown**: Properly handles SIGINT and SIGTERM signals with 30-second timeout
- **Comprehensive Testing**: Full unit test coverage for all components
- **Easy to Use**: Simple command-line interface with configurable config file path

## Quick Start

1. **Build the server**:
   ```bash
   go build -o echo-server ./cmd/server
   ```

2. **Create a configuration file** (see [Configuration](#configuration) section below)

3. **Run the server**:
   ```bash
   ./echo-server -config config.yaml
   ```

## Docker

The application includes a multi-stage Dockerfile that builds a minimal, secure container using Google's distroless base image.

### Building the Docker Image

```bash
# Build the Docker image
docker build -t echo-server .
```

### Running with Docker

```bash
# Run with default configuration
docker run -p 8080:8080 echo-server

# Run with custom configuration
# First, create your custom config file, then mount it:
docker run -p 8080:8080 -v $(pwd)/my-config.yaml:/config.yaml echo-server

# Run in background (detached)
docker run -d -p 8080:8080 --name my-echo-server echo-server

# View logs
docker logs my-echo-server

# Stop the container
docker stop my-echo-server
```

### Docker Features

- **Multi-stage build**: Optimized build process with separate build and runtime stages
- **Distroless base image**: Minimal attack surface using `gcr.io/distroless/static:nonroot`
- **Non-root user**: Runs as non-root user for enhanced security
- **Static binary**: CGO-disabled build for maximum compatibility
- **Health check**: Built-in health checking capability
- **Small image size**: Optimized final image size with only necessary components

### Docker Compose (Optional)

Create a `docker-compose.yml` for easier management:

```yaml
version: '3.8'
services:
  echo-server:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - "./config.yaml:/config.yaml:ro"
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "/echo-server", "--help"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
```

Then run:
```bash
# Start the service
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the service
docker-compose down
```

## Configuration

The server reads its configuration from a YAML file. Here's the structure:

### Basic Structure

```yaml
address: ":8080"        # Server address (default: ":12330")
log_level: "info"       # Log level (default: "info")
routes:                 # Array of route configurations
  - path: "/health"
    method: "GET"                    # Optional, defaults to "GET"
    response_body: "OK"              # Optional, defaults to empty string
    response_header:                 # Optional, defaults to empty
      Content-Type: "text/plain"
    response_dump: false             # Optional, enable request dump for this route
```

### Route Configuration

Each route supports the following fields:

- **`path`** (required): The URL path to match
- **`method`** (optional): HTTP method (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS)
  - Default: "GET"
- **`response_body`** (optional): The response body to return
  - Default: empty string
- **`response_header`** (optional): Map of response headers to set
  - Default: empty map
- **`response_status`** (optional): HTTP status code to return
  - Default: 200
- **`response_dump`** (optional): Enable request headers and query parameters dump in JSON format
  - Default: false
- **`conditions`** (optional): Array of conditional responses based on request headers
  - Default: empty array

### Logging Configuration

The server uses Go's structured logging (`log/slog`) with configurable log levels:

- **`log_level`** (optional): Controls the minimum log level that will be output
  - **Available levels**: `debug`, `info`, `warn`/`warning`, `error`
  - **Default**: `info`
  - **Case insensitive**: `DEBUG`, `Info`, `WaRn` are all valid

#### Log Level Behavior

- **`debug`**: Shows all log messages including request details, route matching, and response information
- **`info`**: Shows server startup, shutdown, and request handling information (default)
- **`warn`**: Shows warnings and errors only
- **`error`**: Shows only error messages

#### Example Log Output

```yaml
# With log_level: "debug"
log_level: "debug"
address: ":8080"
routes:
  - path: "/health"
    response_body: "OK"
```

Debug level will output detailed access logs for each request:
```
time=2025-09-17T12:46:59.000Z level=DEBUG msg="Received request" method=GET path=/health
time=2025-09-17T12:46:59.000Z level=INFO msg="Request handled" method=GET path=/health status=200 response_bytes=2
```

### Conditional Responses

Routes can have conditional responses based on request headers. The server checks conditions in order and uses the first matching condition. If no conditions match, it uses the default route response.

Each condition supports:
- **`header_match`** (required): Map of header key-value pairs that must all match
- **`response_body`** (optional): Response body for this condition
- **`response_status`** (optional): HTTP status code for this condition (default: 200)
- **`response_header`** (optional): Response headers for this condition

### Example Configuration

```yaml
address: ":8080"
routes:
  # Simple health check
  - path: "/health"
    method: "GET"
    response_body: "OK"
  
  # JSON API endpoint
  - path: "/api/users"
    method: "GET"
    response_body: '{"users": [{"id": 1, "name": "John Doe"}]}'
    response_header:
      Content-Type: "application/json"
      X-API-Version: "v1.0"
  
  # POST endpoint
  - path: "/api/users"
    method: "POST"
    response_body: '{"message": "User created successfully"}'
    response_header:
      Content-Type: "application/json"
      Location: "/api/users/123"
  
  # HTML response
  - path: "/welcome"
    response_body: |
      <html>
      <head><title>Welcome</title></head>
      <body><h1>Welcome!</h1></body>
      </html>
    response_header:
      Content-Type: "text/html"
  
  # Route using defaults (GET method, empty response)
  - path: "/ping"
  
  # Header-based conditional responses
  - path: "/api/secure"
    method: "GET"
    # Default response (unauthorized)
    response_body: '{"error": "Unauthorized"}'
    response_status: 401
    response_header:
      Content-Type: "application/json"
    
    # Conditional responses based on headers
    conditions:
      # Valid bearer token
      - header_match:
          Authorization: "Bearer secret-token-123"
        response_body: '{"data": "secret information"}'
        response_status: 200
        response_header:
          Content-Type: "application/json"
          X-User-Role: "user"
      
      # Admin API key
      - header_match:
          X-API-Key: "admin-key-456"
        response_body: '{"data": "admin information"}'
        response_status: 200
        response_header:
          Content-Type: "application/json"
          X-User-Role: "admin"
      
      # Multiple headers required (both must match)
      - header_match:
          Authorization: "Bearer secret-token-123"
          X-Client-ID: "mobile-app"
        response_body: '{"data": "mobile-specific content"}'
        response_status: 200
        response_header:
          Content-Type: "application/json"
          X-Client-Type: "mobile"
```

## Command Line Options

- **`-config`**: Path to the YAML configuration file (default: "config.yaml")

```bash
# Use default config file (config.yaml)
./echo-server

# Use custom config file
./echo-server -config /path/to/my-config.yaml
```

## Architecture

### Project Structure

```
├── cmd/
│   └── server/
│       ├── main.go      # Main server implementation
│       └── main_test.go # Server tests
├── configs/
│   ├── types.go         # Configuration types and methods
│   ├── types_test.go    # Types tests
│   ├── loader.go        # Configuration loading logic
│   └── loader_test.go   # Loader tests
├── config.yaml          # Example configuration
├── go.mod              # Go module dependencies
└── README.md           # This file
```

### Components

#### Server (`cmd/server/main.go`)
- **Server**: Main server struct that holds configuration and fasthttp/router instance
- **initializeRouter**: Initializes the fasthttp/router with all configured routes and HTTP methods (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS)
- **handleRouteRequest**: FastHTTP request handler for individual routes (called by router)
- **handleRoute**: Route response handling with condition matching and response generation

#### Configuration (`configs/`)
- **ServerConfig**: Main configuration structure
- **Route**: Individual route configuration
- **LoadConfig**: YAML configuration loader with validation

## Testing

Run all tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run specific package tests:

```bash
# Test configuration package
go test ./configs

# Test server package
go test ./cmd/server
```

### Test Coverage

The project includes comprehensive unit tests covering:

- Configuration loading and validation
- Route matching logic
- Request handling
- Default value application
- Error conditions
- Edge cases

## Development

### Prerequisites

- Go 1.25.0 or later
- Dependencies are managed via go.mod

### Adding New Features

1. Update configuration types in `configs/types.go` if needed
2. Add validation logic in `configs/loader.go`
3. Implement server logic in `cmd/server/main.go`
4. Add comprehensive tests
5. Update this README

### Dependencies

- **[fasthttp](https://github.com/valyala/fasthttp)**: High-performance HTTP server framework
- **[yaml.v3](https://gopkg.in/yaml.v3)**: YAML parsing library

## Performance

The server is built on top of fasthttp, which provides:

- Up to 10x faster than net/http for high-load scenarios
- Zero memory allocations in hot paths
- Efficient HTTP parsing and response generation
- Built-in support for HTTP/1.1 pipelining

## Examples

### Basic Usage

1. Create `config.yaml`:
   ```yaml
   address: ":8080"
   routes:
     - path: "/hello"
       response_body: "Hello, World!"
   ```

2. Start server:
   ```bash
   go run ./cmd/server -config config.yaml
   ```

3. Test:
   ```bash
   curl http://localhost:8080/hello
   # Output: Hello, World!
   ```

### Testing Header-Based Responses

```bash
# Test the secure endpoint with different headers

# Without authorization (default response)
curl http://localhost:8080/api/secure
# Output: {"error": "Unauthorized", "message": "Valid authorization required"}
# Status: 401

# With valid bearer token
curl -H "Authorization: Bearer secret-token-123" http://localhost:8080/api/secure
# Output: {"data": "This is secret information", "user": "authenticated"}
# Status: 200

# With admin API key
curl -H "X-API-Key: admin-key-456" http://localhost:8080/api/secure
# Output: {"data": "Admin-level information", "permissions": ["read", "write", "admin"]}
# Status: 200

# With multiple headers (both required)
curl -H "Authorization: Bearer secret-token-123" -H "X-Client-ID: mobile-app" http://localhost:8080/api/secure
# Output: {"data": "Mobile-specific content", "features": ["offline", "push"]}
# Status: 200
```

### Advanced Configuration

```yaml
address: ":3000"
routes:
  # REST API simulation
  - path: "/api/v1/health"
    method: "GET"
    response_body: '{"status": "healthy", "timestamp": "2024-01-01T00:00:00Z"}'
    response_header:
      Content-Type: "application/json"
      X-Service: "echo-server"
  
  # Form submission endpoint
  - path: "/submit"
    method: "POST"
    response_body: '{"success": true, "message": "Form submitted successfully"}'
    response_header:
      Content-Type: "application/json"
      Access-Control-Allow-Origin: "*"
  
  # Static content
  - path: "/favicon.ico"
    method: "GET"
    response_body: ""
    response_header:
      Content-Type: "image/x-icon"
      Cache-Control: "public, max-age=3600"
```

### Response Delay Parameter

The server supports adding artificial delays to responses using the `delay` query parameter. This is useful for testing timeout scenarios, simulating slow network conditions, or testing client-side loading states.

#### Supported Delay Formats

- **Duration strings**: `10ms`, `1s`, `500us`, `2m` (supports nanoseconds, microseconds, milliseconds, seconds, minutes, hours)
- **Integer milliseconds**: `100` (interpreted as 100 milliseconds for backward compatibility)

#### Usage Examples

```bash
# Add a 10 millisecond delay
curl "http://localhost:8080/hello?delay=10ms"

# Add a 2 second delay
curl "http://localhost:8080/api/users?delay=2s"

# Add a 500 microsecond delay
curl "http://localhost:8080/health?delay=500us"

# Integer format (backward compatibility) - 250 milliseconds
curl "http://localhost:8080/status?delay=250"

# Combined with other query parameters
curl "http://localhost:8080/search?q=test&delay=1s"
```

#### Error Handling

If an invalid delay parameter is provided, the server returns a `400 Bad Request` response:

```bash
# Invalid delay format
curl "http://localhost:8080/hello?delay=invalid"
# Output: Invalid delay parameter: time: invalid duration "invalid"
# Status: 400
```

#### Implementation Notes

- Delays are applied **before** the response is sent
- Negative delays are technically supported but not recommended
- The delay parameter is logged for monitoring and debugging
- Very long delays may cause client timeouts
- **Shutdown-aware cancellation**: If the server receives a shutdown signal (SIGINT or SIGTERM) during a delay, the request will return early without sending a response, enabling graceful shutdown even with long delays
- Uses Go's `context.Context` with `ctx.Done()` to detect server shutdown and cancel ongoing delay operations

### Router Implementation

The server uses [fasthttp/router](https://github.com/fasthttp/router) for efficient HTTP routing, providing significant performance improvements over manual route matching.

#### Implementation Notes

- **Method-specific routing**: Each route is registered with its specific HTTP method (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS)
- **Proper HTTP status codes**: Returns 405 Method Not Allowed for existing paths with wrong methods, 404 Not Found for non-existent paths
- **Unknown methods**: Routes with unrecognized HTTP methods are registered using the router's `ANY` method and logged as warnings
- **Route closure**: Each route handler is created as a closure that captures the specific route configuration
- **Custom NotFound handler**: Provides consistent 404 responses for unmatched routes

### Response Dump

The server supports dumping request headers and query parameters in JSON format as the response body. This feature is configured per route and is useful for debugging, testing, and understanding what data the server receives from clients.

#### Configuration

Enable response dumping by setting `response_dump: true` on individual routes:

```yaml
address: ":8080"
log_level: "debug"
routes:
  - path: "/debug"
    method: "GET"
    response_body: "This will be replaced by JSON dump"
    response_dump: true    # Enable request dump for this route
  - path: "/normal"
    method: "GET"
    response_body: "Normal response"
    response_dump: false   # Normal response (default)
```

#### Behavior

When `response_dump: true` is set on a route:

1. **Replaces response body**: The configured response body is completely replaced with JSON dump
2. **Sets Content-Type**: Automatically sets Content-Type to `application/json`
3. **Pure JSON format**: Request data is formatted as pretty-printed JSON without any additional text

#### Usage Examples

```bash
# Test with headers and query parameters
curl -H "Authorization: Bearer token123" \
     -H "User-Agent: MyApp/1.0" \
     "http://localhost:8080/debug?param1=value1&param2=value2"
```

**Example response (pure JSON):**
```json
{
  "headers": {
    "Authorization": "Bearer token123",
    "User-Agent": "MyApp/1.0",
    "Host": "localhost:8080",
    "Accept": "*/*"
  },
  "query_parameters": {
    "param1": "value1",
    "param2": "value2"
  }
}
```

#### Use Cases

- **API debugging**: Understanding what headers and parameters clients are sending
- **Testing**: Verifying request data during integration tests
- **Development**: Quick inspection of request details without checking server logs
- **Client troubleshooting**: Help clients understand what data is being transmitted

#### Implementation Notes

- **Route-level configuration**: Each route can independently enable/disable response dumping
- **Performance impact**: Minimal overhead when disabled (default), only affects routes with `response_dump: true`
- **Security consideration**: May expose sensitive headers - use with caution in production
- **Response replacement**: When enabled, completely replaces the configured response body with JSON dump
- **Content-Type**: Automatically sets `Content-Type: application/json` when dumping
- **JSON marshaling**: Request data is marshaled using Go's `json.MarshalIndent` with 2-space indentation
- **Error handling**: If JSON marshaling fails, an error is logged and the original response is returned
- **Default value**: `response_dump` defaults to `false` for each route for security and performance reasons

## License

This project is open source. Feel free to use, modify, and distribute.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## Graceful Shutdown

The server supports graceful shutdown when receiving SIGINT (Ctrl+C) or SIGTERM signals:

- **Immediate**: Stops accepting new connections
- **Graceful**: Allows existing requests to complete (up to 30 seconds)
- **Clean Exit**: Properly closes all resources

### Shutdown Process

1. Signal received (SIGINT/SIGTERM)
2. Server stops accepting new connections
3. Existing requests are allowed to finish (30s timeout)
4. Server exits cleanly

```bash
# Normal operation
./echo-server -config config.yaml
# ... server running ...

# Graceful shutdown (Ctrl+C or kill -TERM <pid>)
^C2025/09/16 18:27:17 Received shutdown signal, shutting down gracefully...
2025/09/16 18:27:17 Server exited gracefully
```

## Troubleshooting

### Common Issues

1. **Server won't start**: Check if the port is already in use
2. **Config file not found**: Verify the path to your config file
3. **Invalid YAML**: Validate your YAML syntax
4. **Route not matching**: Ensure exact path matching (case-sensitive)
5. **Graceful shutdown timeout**: Increase timeout or optimize slow endpoints

### Debugging

The server logs all requests and responses to stdout:

```
2024/01/01 12:00:00 Starting server on :8080
2024/01/01 12:00:00 Loaded 3 routes
2024/01/01 12:00:01 Received GET request for /health
2024/01/01 12:00:01 Handled GET /health -> 200 status, 2 bytes
```
