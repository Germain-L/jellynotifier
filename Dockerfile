# Multi-stage build for Go application
# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies for CGO (required for SQLite)
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Set working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o jellynotifier .

# Final stage - minimal runtime image
FROM alpine:latest

# Install ca-certificates and sqlite for runtime
RUN apk --no-cache add ca-certificates sqlite

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Create data directory for SQLite database
RUN mkdir -p /app/data && chown appuser:appgroup /app/data

# Copy the binary from builder stage
COPY --from=builder /app/jellynotifier .

# Change ownership to non-root user
RUN chown appuser:appgroup jellynotifier

# Switch to non-root user
USER appuser

# Expose port 8080
EXPOSE 8080

# Health check to ensure the server is running
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the binary
CMD ["./jellynotifier"]