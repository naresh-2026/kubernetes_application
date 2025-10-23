# syntax=docker/dockerfile:1
FROM golang:1.21 AS builder

# Set working directory inside container
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download
# Copy source files
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go

# Use a minimal base image for running the app
FROM alpine:latest

# Install required certificates
RUN apk --no-cache add ca-certificates

# Set timezone to India (IST)
ENV TZ=Asia/Kolkata
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Set working directory
WORKDIR /myapp1/

# Copy the compiled binary from builder stage
COPY --from=builder /app/main .
COPY static ./static
# Ensure the binary has execution permissions
RUN chmod +x /myapp1/main
# Expose port 8082
EXPOSE 8082
ENV WS_HOST="0.0.0.0"

# Command to run the application
CMD ["./main"]

