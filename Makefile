# Makefile for Peripage Printer API

.PHONY: help build run test test-cover test-race test-integration clean docker-build docker-up docker-down swagger deps

# Default target
help:
	@echo "Available targets:"
	@echo "  make deps             - Download Go dependencies"
	@echo "  make swagger          - Generate Swagger documentation"
	@echo "  make build            - Build the binary"
	@echo "  make run              - Run the application (mock printer)"
	@echo "  make run-ble          - Run the application (BLE printer)"
	@echo "  make test             - Run tests"
	@echo "  make test-cover       - Run tests with coverage report"
	@echo "  make test-race        - Run tests with race detector"
	@echo "  make test-integration - Run integration tests (requires hardware)"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-up-dev    - Start with Docker Compose (dev/mock)"
	@echo "  make docker-up        - Start with Docker Compose (production/BLE)"
	@echo "  make docker-down      - Stop Docker Compose"

# Install dependencies
deps:
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go -o internal/adapters/docs

# Build the binary
build: swagger
	go build -o bin/peripage-server ./cmd/server

# Run with mock printer
run: swagger
	PRINTER_TYPE=mock go run cmd/server/main.go

# Run with BLE printer
run-ble: swagger
	PRINTER_TYPE=ble go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detector
test-race:
	go test ./... -race

# Run integration tests (requires hardware)
test-integration:
	go test -tags=integration ./test/integration/... -v

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f internal/adapters/docs/docs.go
	rm -f internal/adapters/docs/swagger.json
	rm -f internal/adapters/docs/swagger.yaml
	rm -f coverage.out
	rm -f coverage.html

# Build Docker image
docker-build:
	docker-compose build

# Start with Docker Compose (development)
docker-up-dev:
	docker-compose -f docker-compose.dev.yml up --build

# Start with Docker Compose (production)
docker-up:
	docker-compose up --build

# Stop Docker Compose
docker-down:
	docker-compose down
	docker-compose -f docker-compose.dev.yml down
