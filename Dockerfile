# Use the official Golang image as the base image for building
FROM golang:1.25.0-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
ARG DEP_UPDATE=false
ARG DEP_CLEAN_CACHE=false

# Update dependencies if requested
RUN if [ "$DEP_UPDATE" = "true" ]; then \
        go get -u ./... && \
        go mod tidy; \
    else \
        go mod download; \
    fi

# Clean Go module cache if requested
RUN if [ "$DEP_CLEAN_CACHE" = "true" ]; then \
        go clean -modcache; \
    fi

# Copy the source code into the container
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gemini-proxy cmd/gemini-proxy/main.go

# Use a minimal base image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/gemini-proxy .

# Change ownership of the binary to the non-root user
RUN chown appuser:appuser gemini-proxy

# Switch to the non-root user
USER appuser

# Expose port 8080
EXPOSE 8080

# Command to run the binary
CMD ["./gemini-proxy"]