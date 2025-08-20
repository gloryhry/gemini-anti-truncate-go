#!/bin/bash

# Test runner script for the Gemini Anti-Truncate Go project

# Exit on any error
set -e

# Get the project directory
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Change to the project directory
cd "$PROJECT_DIR"

echo "Running tests for Gemini Anti-Truncate Go project..."

# 1. Run unit tests with coverage
echo "1. Running unit tests with coverage..."
go test -v -coverprofile=coverage.out -coverpkg=./... ./internal/...

# 2. Run integration tests
echo "2. Running integration tests..."
go test -v ./test/...

# 3. Generate coverage report
echo "3. Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

# 4. Run benchmark tests
echo "4. Running benchmark tests..."
go test -bench=. -benchmem ./test/...

# 5. Run race condition tests
echo "5. Running race condition tests..."
go test -race ./internal/...

echo "All tests completed successfully!"

# Show coverage summary
echo "Coverage summary:"
go tool cover -func=coverage.out

echo "Coverage report saved to coverage.html"