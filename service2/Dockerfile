# syntax=docker/dockerfile:1
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Set working directory
WORKDIR /myapp2/

# Install necessary dependencies
RUN apk add --no-cache ca-certificates

# Copy the built server binary from the builder stage
COPY --from=builder /app/main .

# Copy static files (HTML, JS, CSS)
COPY static ./static
RUN chmod +x /myapp2/main
# Expose port 8083 for Kubernetes service
EXPOSE 8083

# Set environment variable for WebSocket to listen on all interfaces
ENV WS_HOST="0.0.0.0"

# Start the server
CMD ["./main"]


# FROM golang:1.23.7 AS builder

# WORKDIR /app

# COPY go.mod ./
# RUN go mod download
# COPY . .

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# FROM alpine:latest

# RUN apk --no-cache add ca-certificates

# WORKDIR /myapp/
# COPY --from=builder /app/main .

# RUN chmod +x /myapp/main

# EXPOSE 8083

# CMD ["./main"]




# Stage 1: Build the Go application