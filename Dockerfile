# Build stage
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and ca-certificates for dependency downloads
RUN apk add --no-cache git ca-certificates

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary, GOOS=linux for Linux target
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o echo-server ./cmd/server

# Final stage - distroless
FROM gcr.io/distroless/static-debian12:nonroot

# Copy ca-certificates from builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary from builder stage
COPY --from=builder /app/echo-server /echo-server

# Use nonroot user for security
USER nonroot:nonroot

# Expose the port the app runs on
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["/echo-server"]

