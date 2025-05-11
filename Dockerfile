# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
# -o /app/main: output binary to /app/main
# CGO_ENABLED=0: build a statically linked binary
# GOOS=linux: build for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/api/main.go

# Stage 2: Create the final lightweight image
FROM alpine:latest

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose port (must match the port the app listens on)
EXPOSE 8080

# Command to run the application
# The DSN and APP_PORT will be passed as environment variables from docker-compose.yml
CMD ["./main"]