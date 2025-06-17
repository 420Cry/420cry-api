# Go binary name
BINARY_NAME=cry-api

.PHONY: build run clean install lint test dev migrate lint-fix

# Build the Go application
build:
	go build -o $(BINARY_NAME) app/cmd/main.go

# Run the Go application
run:
	go run app/cmd/main.go

# Clean the project (remove binaries)
clean:
	rm -f $(BINARY_NAME)

# Install the Go dependencies
install:
	go mod tidy

# Lint
lint:
	golangci-lint run

# Lint fix
lint-fix:
	gofumpt -w .
	goimports -w .

# Test the Go application
test:
	go test ./tests/...

# Run the server with the binary
dev: build
	./$(BINARY_NAME)

# Run the migration script
migrate:
	go run app/migration/migration.go