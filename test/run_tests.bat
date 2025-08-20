@echo off

REM Test runner script for the Gemini Anti-Truncate Go project (Windows)

REM Get the project directory
set PROJECT_DIR=%~dp0..

REM Change to the project directory
cd /d "%PROJECT_DIR%"

echo Running tests for Gemini Anti-Truncate Go project...

REM 1. Run unit tests with coverage
echo 1. Running unit tests with coverage...
go test -v -coverprofile=coverage.out -coverpkg=./... ./internal/...

REM 2. Run integration tests
echo 2. Running integration tests...
go test -v ./test/...

REM 3. Generate coverage report
echo 3. Generating coverage report...
go tool cover -html=coverage.out -o coverage.html

REM 4. Run benchmark tests
echo 4. Running benchmark tests...
go test -bench=. -benchmem ./test/...

REM 5. Run race condition tests
echo 5. Running race condition tests...
go test -race ./internal/...

echo All tests completed successfully!

REM Show coverage summary
echo Coverage summary:
go tool cover -func=coverage.out

echo Coverage report saved to coverage.html

pause