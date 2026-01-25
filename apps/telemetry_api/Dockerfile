# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Generate Swagger documentation
RUN go run github.com/swaggo/swag/cmd/swag@latest init

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o telemetry_api

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/telemetry_api .

# Copy Swagger docs
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE 8083

# Command to run
CMD ["./telemetry_api"]
