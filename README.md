# Gemini Anti-Truncate Go

A Go implementation of the Gemini API proxy service that prevents response truncation. This project is a rewrite of the original Cloudflare Workers-based JavaScript implementation, designed for better performance and easier deployment.

## Features

- Prevents Gemini API response truncation using token injection and retry mechanisms
- Supports both streaming and non-streaming responses
- Handles function calls and structured outputs correctly
- Compatible with the original JavaScript API
- Containerized deployment with Docker
- Comprehensive test suite

## Architecture

The service works by:

1. Intercepting requests to target Gemini models
2. Injecting system instructions to append a finish token to responses
3. Processing responses to detect the finish token
4. Automatically retrying incomplete responses
5. Forwarding complete responses to clients

## Getting Started

### Prerequisites

- Go 1.22 or later
- Docker (for containerized deployment)
- A Gemini API key

### Go Modules

This project uses Go modules for dependency management. The `go.mod` and `go.sum` files are included in the repository. To ensure all dependencies are correctly set up:

```bash
# Initialize and verify modules
go mod tidy

# Download dependencies
go mod download
```

### Building

```bash
# Clone the repository
git clone <repository-url>
cd gemini-anti-truncate-go

# Initialize Go modules (if not already done)
go mod tidy

# Build the binary
go build -o gemini-proxy cmd/gemini-proxy/main.go

# Or use make
make build
```

### Running

```bash
# Set your Gemini API key
export GEMINI_API_KEY=your-api-key

# Run the service
./gemini-proxy

# Or with custom configuration
UPSTREAM_URL_BASE=https://generativelanguage.googleapis.com \
MAX_RETRIES=20 \
DEBUG_MODE=true \
HTTP_PORT=8080 \
./gemini-proxy
```

### Docker

#### 本地构建
```bash
# Build the Docker image
docker build -t gemini-proxy .

# Build with dependency updates
docker build --build-arg DEP_UPDATE=true -t gemini-proxy .

# Build with dependency updates and clean cache
docker build --build-arg DEP_UPDATE=true --build-arg DEP_CLEAN_CACHE=true -t gemini-proxy .

# Run with Docker
docker run -p 8080:8080 -e GEMINI_API_KEY=your-api-key gemini-proxy

# Or with docker-compose
docker-compose up
```

#### Dependency Management in Docker

The Docker build process supports flexible dependency management through build arguments:

- `DEP_UPDATE`: Set to `true` to update dependencies during build
- `DEP_CLEAN_CACHE`: Set to `true` to clean Go module cache after download

This allows for better control over dependency updates and image size optimization.

#### 从 GitHub Container Registry 拉取镜像

本项目已配置自动化 CI/CD，每次推送到 `main` 分支会自动构建并推送 Docker 镜像到 GitHub Container Registry (GHCR)。

**镜像标签**:
- `ghcr.io/your-username/gemini-anti-truncate-go:latest` - 最新版本（main 分支）
- `ghcr.io/your-username/gemini-anti-truncate-go:develop` - 开发版本（develop 分支）
- `ghcr.io/your-username/gemini-anti-truncate-go:v1.0.0` - 版本标签

**拉取镜像**:
```bash
# 拉取最新版本
docker pull ghcr.io/your-username/gemini-anti-truncate-go:latest

# 拉取指定版本
docker pull ghcr.io/your-username/gemini-anti-truncate-go:v1.0.0
```

**运行容器**:
```bash
# 基本运行
docker run -p 8080:8080 \
  -e GEMINI_API_KEY=your_api_key \
  -e UPSTREAM_URL_BASE=https://generativelanguage.googleapis.com \
  ghcr.io/your-username/gemini-anti-truncate-go:latest

# 后台运行
docker run -d -p 8080:8080 \
  --name gemini-proxy \
  -e GEMINI_API_KEY=your_api_key \
  -e UPSTREAM_URL_BASE=https://generativelanguage.googleapis.com \
  ghcr.io/your-username/gemini-anti-truncate-go:latest
```

## Configuration

The service can be configured using environment variables:

- `UPSTREAM_URL_BASE`: The base URL for the Gemini API (default: `https://generativelanguage.googleapis.com`)
- `MAX_RETRIES`: Maximum number of retries for incomplete responses (default: `20`)
- `DEBUG_MODE`: Enable debug logging (default: `false`)
- `HTTP_PORT`: Port to listen on (default: `8080`)
- `GEMINI_API_KEY`: Your Gemini API key (can also be provided in requests)

## API Usage

The service proxies requests to the Gemini API:

```bash
# Non-streaming request
curl -X POST http://localhost:8080/v1beta/models/gemini-1.5-pro-latest:generateContent \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"contents":[{"parts":[{"text":"Hello, world!"}]}]}'

# Streaming request
curl -X POST http://localhost:8080/v1beta/models/gemini-1.5-pro-latest:streamGenerateContent \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"contents":[{"parts":[{"text":"Hello, world!"}]}]}'
```

## Testing

The project includes a comprehensive test suite. See [test/README.md](test/README.md) for detailed information on running tests.

```bash
# Run all tests
make test-all

# Run unit tests
make test

# Run integration tests
make test-integration

# Run tests with coverage
make test-coverage
```

## Dependency Management

This project uses Go modules for dependency management. Dependencies are declared in `go.mod` and their exact versions are tracked in `go.sum`.

### Managing Dependencies

To add a new dependency:
```bash
go get github.com/some/package
```

To update a dependency:
```bash
go get -u github.com/some/package
```

To update all dependencies to their latest minor/patch versions:
```bash
go get -u ./...
```

After updating dependencies, always run:
```bash
go mod tidy
```

### Dependency Conflict Resolution

If you encounter dependency conflicts:

1. **Identify conflicts**: Run `go mod graph` to see the dependency tree
2. **Resolve conflicts**: Use `go mod edit -replace` directive in `go.mod` for temporary overrides
3. **Verify resolution**: Run `go mod tidy` and tests to ensure conflicts are resolved

Example of resolving a conflict:
```bash
# In go.mod, add a replace directive:
replace github.com/conflicting/package => github.com/conflicting/package v1.2.3
```

### Dependency Upgrade Strategy

1. **Regular updates**: Update dependencies regularly to get security fixes and improvements
2. **Major version updates**: Test thoroughly as they may contain breaking changes
3. **Security scanning**: Use tools like `govulncheck` to identify vulnerabilities
4. **Version pinning**: Pin to specific versions in production environments

### Security Scanning

To scan for vulnerabilities in dependencies:
```bash
# Install govulncheck if not already installed
go install golang.org/x/vuln/cmd/govulncheck@latest

# Scan for vulnerabilities
govulncheck ./...
```

Regular security scanning should be part of your CI/CD pipeline to ensure dependencies are secure.

## CI/CD Workflows

本项目配置了两个 GitHub Actions 工作流：

### 1. 测试工作流 (`.github/workflows/test.yml`)
- **触发条件**: 推送到 `main` 或 `develop` 分支，或创建 Pull Request
- **包含测试**:
  - 单元测试
  - 集成测试  
  - 覆盖率测试
  - 基准测试
  - 竞态条件测试
  - Docker 镜像构建测试

### 2. Docker 构建和推送工作流 (`.github/workflows/docker.yml`)
- **触发条件**:
  - 推送到 `main` 或 `develop` 分支
  - 创建 Pull Request
  - 发布 Release
- **包含功能**:
  - 多平台构建 (linux/amd64, linux/arm64)
  - 自动推送镜像到 GitHub Container Registry
  - 自动生成镜像标签
  - 缓存优化
  - Release 发布时自动创建 Release Notes

#### 工作流特性
- **多架构支持**: 支持 AMD64 和 ARM64 架构
- **智能标签**: 根据分支、PR、版本自动生成标签
- **缓存优化**: 使用 GitHub Actions 缓存加速构建
- **安全性**: 使用 GitHub Token 进行安全认证
- **版本管理**: 支持语义化版本标签

#### 镜像推送权限
- **推送权限**: 工作流会自动请求必要的权限
- **认证方式**: 使用 `secrets.GITHUB_TOKEN` 进行认证
- **镜像仓库**: `ghcr.io/your-username/gemini-anti-truncate-go`

## Performance

The Go implementation provides better performance than the original JavaScript version:

- Lower latency due to compiled language
- Better memory management
- Improved concurrency handling
- More efficient JSON processing

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Based on the original Cloudflare Workers implementation
- Uses the Gorilla Mux router for HTTP handling