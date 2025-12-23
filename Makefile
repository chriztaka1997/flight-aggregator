.PHONY: run build test clean install-deps

# Run the application
run:
	go run cmd/server/main.go

# Build the application
build:
	mkdir -p bin
	go build -o bin/flight-aggregator cmd/server/main.go

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Install dependencies
install-deps:
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Run linter (requires golangci-lint to be installed)
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...