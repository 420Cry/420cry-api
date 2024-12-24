# Makefile for Go project

# Go binary name
BINARY_NAME=420cry-api

# Build the Go application
build:
	go build -o $(BINARY_NAME) cmd/server/main.go

# Run the Go application
run:
	go run cmd/server/main.go

# Clean the project (remove binaries)
clean:
	rm -f $(BINARY_NAME)

# Install the Go dependencies
install:
	go mod tidy

# Test the Go application
test:
	go test ./...

# Run the server with the binary
dev: build
	./$(BINARY_NAME)
