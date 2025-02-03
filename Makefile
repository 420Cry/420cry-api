# Go binary name
BINARY_NAME=cry-api

# Build the Go application
build:
	go build -o $(BINARY_NAME) app/server/server.go

# Run the Go application
run:
	go run app/server/server.go

# Clean the project (remove binaries)
clean:
	rm -f $(BINARY_NAME)

# Install the Go dependencies
install:
	go mod tidy

# Lint
lint:
	golangci-lint run

# Test the Go application
test:
	go test ./tests/...

# Run the server with the binary
dev: build
	./$(BINARY_NAME)

# Run the migration script
migrate:
	go run app/migration/migration.go
