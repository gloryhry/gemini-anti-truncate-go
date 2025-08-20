# Testing Documentation

This document provides comprehensive information about testing the Gemini Anti-Truncate Go project.

## Table of Contents

1. [Overview](#overview)
2. [Test Structure](#test-structure)
3. [Running Tests](#running-tests)
4. [Test Types](#test-types)
5. [Coverage](#coverage)
6. [CI/CD Integration](#cicd-integration)
7. [Performance Testing](#performance-testing)
8. [Docker Testing](#docker-testing)

## Overview

The Gemini Anti-Truncate Go project includes a comprehensive test suite covering unit tests, integration tests, performance tests, and Docker deployment tests. The tests ensure that the Go implementation is functionally equivalent to the original JavaScript version and meets all performance and reliability requirements.

## Test Structure

The test suite is organized as follows:

```
test/
├── integration_test.go       # Integration tests
├── api_compatibility_test.go # API compatibility tests
├── performance_test.go       # Performance tests
├── docker_test.go            # Docker deployment tests
├── data/                     # Test data files
│   ├── sample_request.json
│   ├── sample_response_complete.json
│   ├── sample_response_incomplete.json
│   ├── sample_stream_complete.txt
│   └── sample_stream_incomplete.txt
├── .env.test                 # Test environment variables
├── run_tests.sh              # Test runner script (Linux/Mac)
├── run_tests.bat             # Test runner script (Windows)
├── coverage.yaml             # Coverage configuration
└── README.md                 # This file

internal/
├── config/
│   └── config_test.go        # Configuration tests
├── gemini/
│   └── constants_test.go     # Constants tests
├── util/
│   └── util_test.go          # Utility tests
├── proxy/
│   └── proxy_test.go         # Proxy logic tests
└── handler/
    └── handler_test.go       # HTTP handler tests
```

## Running Tests

### Using Makefile (Recommended)

```bash
# Run all tests
make test-all

# Run unit tests only
make test

# Run integration tests
make test-integration

# Run tests with coverage
make test-coverage

# Run benchmark tests
make test-bench

# Run race condition tests
make test-race
```

### Using Go Commands

```bash
# Run unit tests
go test -v ./internal/...

# Run integration tests
go test -v ./test/...

# Run tests with coverage
go test -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out -o coverage.html

# Run benchmark tests
go test -bench=. ./test/...

# Run race condition tests
go test -race ./internal/...
```

### Using Test Scripts

```bash
# Linux/Mac
./test/run_tests.sh

# Windows
test\run_tests.bat
```

## Test Types

### Unit Tests

Unit tests verify the correctness of individual functions and components:

- **Configuration tests** (`internal/config/config_test.go`): Test environment variable loading and configuration parsing
- **Constants tests** (`internal/gemini/constants_test.go`): Verify that constants are correctly defined
- **Utility tests** (`internal/util/util_test.go`): Test utility functions like API key extraction and error handling
- **Proxy tests** (`internal/proxy/proxy_test.go`): Test core logic functions like token injection and retry request building
- **Handler tests** (`internal/handler/handler_test.go`): Test HTTP request handling and routing

### Integration Tests

Integration tests verify that components work together correctly:

- **HTTP handler tests** (`test/integration_test.go`): Test end-to-end HTTP request processing
- **API compatibility tests** (`test/api_compatibility_test.go`): Ensure compatibility with the original JavaScript API
- **Stream processing tests**: Verify correct handling of SSE streams
- **Non-stream processing tests**: Verify correct handling of regular JSON responses

### Performance Tests

Performance tests ensure the application meets performance requirements:

- **Benchmark tests** (`test/performance_test.go`): Measure performance of critical functions
- **Concurrency tests**: Verify correct handling of concurrent requests
- **Memory usage tests**: Ensure memory usage remains within acceptable limits
- **Response time tests**: Verify that response times meet requirements

### Docker Tests

Docker tests verify that the application can be correctly containerized and deployed:

- **Docker build tests** (`test/docker_test.go`): Verify that Docker images can be built
- **Docker run tests**: Verify that containers can be started
- **Environment variable tests**: Verify that environment variables are correctly configured
- **Docker Compose tests**: Verify that docker-compose works correctly

## Coverage

The project aims for comprehensive test coverage:

- **Unit test coverage**: 90%+ for critical business logic
- **Integration test coverage**: 80%+ for API endpoints
- **Overall coverage**: 70%+ for the entire codebase

To generate and view coverage reports:

```bash
# Generate coverage report
make test-coverage

# Or manually:
go test -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out -o coverage.html
```

Coverage reports are generated in HTML format and can be viewed in a web browser.

## CI/CD Integration

The project includes GitHub Actions workflows for automated testing:

- **Test workflow** (`.github/workflows/test.yml`): Runs all tests on every push and pull request
- **Docker workflow**: Builds and tests Docker images

The workflows test:

1. Unit tests
2. Integration tests
3. Coverage analysis
4. Benchmark tests
5. Race condition detection
6. Docker image building

## Performance Testing

Performance tests ensure the application meets performance requirements:

### Benchmark Tests

Benchmark tests measure the performance of critical functions:

- `BenchmarkInjectFinishToken`: Measures token injection performance
- `BenchmarkBuildRetryRequest`: Measures retry request building performance
- `BenchmarkJSONMarshal`: Measures JSON marshaling performance
- `BenchmarkJSONUnmarshal`: Measures JSON unmarshaling performance

### Concurrency Tests

Concurrency tests verify that the application can handle multiple simultaneous requests without issues.

### Memory Usage Tests

Memory usage tests ensure that the application doesn't have memory leaks and uses memory efficiently.

## Docker Testing

Docker tests verify that the application can be containerized and deployed correctly:

### Docker Build Tests

Verify that Docker images can be built successfully using the provided Dockerfile.

### Docker Run Tests

Verify that containers can be started and run correctly with the proper configuration.

### Environment Variable Tests

Verify that environment variables are correctly passed to the container and used by the application.

### Docker Compose Tests

Verify that the application can be deployed using docker-compose with the provided configuration.

## Troubleshooting

### Common Issues

1. **Tests failing due to missing dependencies**: Run `go mod tidy` to ensure all dependencies are correctly resolved.

2. **Coverage below threshold**: Review uncovered code paths and add additional tests.

3. **Docker tests failing**: Ensure Docker is installed and running on your system.

4. **Performance test failures**: Check system resources and optimize code if necessary.

### Debugging Tests

To run tests with verbose output:

```bash
go test -v ./...
```

To run a specific test:

```bash
go test -v -run TestName ./path/to/package
```

To enable debugging output during tests, set the DEBUG_MODE environment variable:

```bash
DEBUG_MODE=true go test -v ./...
```

## Best Practices

1. **Write tests before implementing new features** (Test-Driven Development)
2. **Keep tests focused and specific** - each test should verify one behavior
3. **Use descriptive test names** that clearly indicate what is being tested
4. **Test edge cases** and error conditions
5. **Maintain test data** in separate files for easy maintenance
6. **Regularly update tests** when modifying code
7. **Monitor coverage** and strive to improve it over time
8. **Run tests frequently** during development
9. **Regularly scan dependencies** for security vulnerabilities
10. **Verify dependency integrity** before committing changes

## Dependency Management and Security

Regular dependency management and security scanning are crucial for maintaining a healthy codebase. See [DEPENDENCY.md](../DEPENDENCY.md) for detailed information about:

- Dependency conflict resolution
- Security scanning procedures
- Dependency upgrade strategies
- Docker dependency management

### Security Scanning

To scan for vulnerabilities in dependencies:
```bash
# Using make
make security-scan

# Or directly
govulncheck ./...
```

### Dependency Verification

To verify dependency integrity:
```bash
# Using make
make deps-verify

# Or directly
go mod verify
go mod tidy
git diff --exit-code go.mod go.sum
```

## Conclusion

The comprehensive test suite ensures that the Gemini Anti-Truncate Go project is reliable, performant, and compatible with the original JavaScript implementation. Regular testing helps maintain code quality and prevents regressions.