# Makefile for the Gemini Anti-Truncate Go project

# Variables
PROJECT_NAME = gemini-proxy
BINARY_NAME = $(PROJECT_NAME)
PROJECT_DIR = $(PWD)
TEST_DIR = $(PROJECT_DIR)/test

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOVET = $(GOCMD) vet
GOFMT = $(GOCMD) fmt

# Default target
all: build

# Build the project
build:
	cd cmd/$(PROJECT_NAME) && $(GOBUILD) -o ../../$(BINARY_NAME) .

# Install dependencies
deps:
	$(GOGET) -v ./...

# Run tests
test:
	$(GOTEST) -v ./internal/...

# Run integration tests
test-integration:
	$(GOTEST) -v ./test/...

# Run tests with coverage
test-coverage:
	$(GOTEST) -coverprofile=coverage.out ./internal/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	$(GOCMD) tool cover -func=coverage.out

# Run benchmark tests
test-bench:
	$(GOTEST) -bench=. -benchmem ./test/...

# Run race condition tests
test-race:
	$(GOTEST) -race ./internal/...

# Run all tests
test-all:
	$(GOTEST) -v -coverprofile=coverage.out ./internal/...
	$(GOTEST) -v ./test/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	$(GOCMD) tool cover -func=coverage.out

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	rm -f coverage.html

# Format code
fmt:
	$(GOFMT) ./...

# Vet code
vet:
	$(GOVET) ./...

# Run Docker build
docker-build:
	docker build -t $(PROJECT_NAME) .

# Run Docker container
docker-run:
	docker run -p 8080:8080 --env-file .env $(PROJECT_NAME)

# Run Docker Compose
docker-compose-up:
	docker-compose up --build

# Run Docker Compose in detached mode
docker-compose-up-detached:
	docker-compose up -d --build

# Stop Docker Compose
docker-compose-down:
	docker-compose down

# Help
help:
	@echo "Available targets:"
	@echo "  all                  - Build the project (default)"
	@echo "  build                - Build the project"
	@echo "  deps                 - Install dependencies"
	@echo "  test                 - Run unit tests"
	@echo "  test-integration     - Run integration tests"
	@echo "  test-coverage        - Run tests with coverage"
	@echo "  test-bench           - Run benchmark tests"
	@echo "  test-race            - Run race condition tests"
	@echo "  test-all             - Run all tests"
	@echo "  clean                - Clean build artifacts"
	@echo "  fmt                  - Format code"
	@echo "  vet                  - Vet code"
	@echo "  docker-build         - Build Docker image"
	@echo "  docker-run           - Run Docker container"
	@echo "  docker-compose-up    - Run Docker Compose"
	@echo "  docker-compose-up-detached - Run Docker Compose in detached mode"
	@echo "  docker-compose-down  - Stop Docker Compose"
	@echo "  help                 - Show this help message"

.PHONY: all build deps test test-integration test-coverage test-bench test-race test-all clean fmt vet docker-build docker-run docker-compose-up docker-compose-up-detached docker-compose-down help