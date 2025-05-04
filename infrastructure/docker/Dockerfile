# Use an official Golang runtime as a parent image
FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod tidy

# Copy the entire project to the working directory
COPY . .

# Build the Go app
RUN go build -o 420cry-api app/server/server.go

# Start a new stage from a smaller image to minimize the size
FROM alpine:latest

# Install necessary dependencies (ca-certificates for HTTPS requests)
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder image
COPY --from=builder /app/420cry-api .

# Expose the port the app runs on (e.g., 8080)
EXPOSE 8080

# Command to run the executable
CMD ["./420cry-api"]
