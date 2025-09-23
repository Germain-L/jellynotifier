# Multi-stage build for Go application
# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary, GOOS=linux for Linux compatibility
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o jellynotifier .

# Final stage - minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS requests (if needed)
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /root/

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
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/webhook || exit 1

# Run the binary
CMD ["./jellynotifier"]