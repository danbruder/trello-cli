# Build stage
FROM golang:1.23-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o trello-cli .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S trello && \
    adduser -u 1001 -S trello -G trello

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/trello-cli .

# Change ownership to non-root user
RUN chown trello:trello trello-cli

# Switch to non-root user
USER trello

# Expose port (if needed for future web features)
EXPOSE 8080

# Set the binary as entrypoint
ENTRYPOINT ["./trello-cli"]
