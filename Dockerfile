# ---- Builder Stage ----
# This stage uses the official Go image to build the application binary.
FROM golang:1.22-alpine AS builder

# Build arguments for dependency management
ARG DEP_UPDATE=false
ARG DEP_CLEAN_CACHE=false

# Set the working directory inside the container
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN if [ "$DEP_UPDATE" = "true" ]; then \
        go get -u ./... && \
        go mod tidy; \
    fi && \
    go mod download && \
    if [ "$DEP_CLEAN_CACHE" = "true" ]; then \
        rm -rf /go/pkg/mod/cache; \
    fi && \
    go mod verify

# Copy the source code into the container (only necessary files)
COPY cmd/ cmd/
COPY internal/ internal/

# Build the Go application, creating a statically linked binary
# CGO_ENABLED=0 is important for creating a static binary without C dependencies
# GOOS=linux ensures the binary is built for a Linux environment
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gemini-proxy ./cmd/gemini-proxy

# ---- Final Stage ----
# This stage uses a minimal Alpine image to create a small and secure final image.
FROM alpine:latest

# Alpine Linux comes with a minimal package set, so we add root certificates
# which are necessary for making HTTPS requests to the upstream Gemini API.
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the pre-built binary file from the builder stage.
COPY --from=builder /app/gemini-proxy .

# Expose port 8080 to the outside world. This is the port the server will listen on.
EXPOSE 8080

# Command to run the executable when the container starts.
CMD ["./gemini-proxy"]
