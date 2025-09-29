# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger documentation
RUN go run github.com/swaggo/swag/cmd/swag@latest init

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o devices_command_api

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/devices_command_api .

# Copy Swagger docs
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE 8082

# Run the application
CMD ["./devices_command_api"]
