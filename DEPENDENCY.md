# Dependency Management Guide

This document provides detailed information about dependency management strategies, security practices, and upgrade procedures for the Gemini Anti-Truncate Go project.

## Table of Contents

1. [Overview](#overview)
2. [Dependency Management Strategy](#dependency-management-strategy)
3. [Security Practices](#security-practices)
4. [Upgrade Procedures](#upgrade-procedures)
5. [Conflict Resolution](#conflict-resolution)
6. [Docker Dependency Management](#docker-dependency-management)
7. [CI/CD Integration](#cicd-integration)

## Overview

The Gemini Anti-Truncate Go project uses Go modules for dependency management. Dependencies are declared in `go.mod` and their exact versions are tracked in `go.sum`. This document outlines best practices for managing these dependencies securely and efficiently.

## Dependency Management Strategy

### Go Modules

Go modules provide reproducible builds by tracking exact dependency versions. Key principles:

1. **Explicit dependencies**: All dependencies are explicitly declared in `go.mod`
2. **Immutable versions**: Exact versions are recorded in `go.sum`
3. **Semantic versioning**: Dependencies should follow semantic versioning
4. **Minimal dependencies**: Only include necessary dependencies

### Version Pinning

For production environments, pin dependencies to specific versions:
```bash
go get github.com/some/package@v1.2.3
```

### Dependency Graph

To visualize the dependency tree:
```bash
go mod graph
```

To see why a dependency is needed:
```bash
go mod why github.com/some/package
```

## Security Practices

### Vulnerability Scanning

Use `govulncheck` to scan for known vulnerabilities:
```bash
# Install govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest

# Scan for vulnerabilities
govulncheck ./...
```

### Dependency Verification

Regularly verify dependency integrity:
```bash
go mod verify
```

### Security Updates

1. Subscribe to security announcements for critical dependencies
2. Regularly scan for vulnerabilities
3. Update dependencies promptly when security issues are discovered
4. Review security advisories from the Go team

## Upgrade Procedures

### Minor/Patch Updates

For minor and patch version updates:
```bash
# Update all dependencies to latest minor/patch versions
go get -u ./...

# Update a specific dependency
go get -u github.com/some/package

# Clean up unused dependencies
go mod tidy
```

### Major Updates

Major version updates may contain breaking changes:
1. Read release notes carefully
2. Update one dependency at a time
3. Run all tests after each update
4. Verify functionality in a staging environment

### Testing After Updates

Always run the full test suite after dependency updates:
```bash
make test-all
```

## Conflict Resolution

### Identifying Conflicts

To identify dependency conflicts:
```bash
go mod graph | grep conflicting-package
```

### Resolving Conflicts

Use replace directives in `go.mod` for temporary solutions:
```go
replace github.com/conflicting/package => github.com/conflicting/package v1.2.3
```

For permanent solutions:
1. Update all dependent packages to compatible versions
2. Remove replace directives
3. Run `go mod tidy`

### Indirect Dependencies

Minimize indirect dependencies by:
1. Using direct dependencies only when necessary
2. Regularly reviewing and cleaning up dependencies
3. Using tools like `go mod graph` to identify unnecessary indirect dependencies

## Docker Dependency Management

The Docker build process includes flexible dependency management options:

### Build Arguments

- `DEP_UPDATE`: Set to `true` to update dependencies during build
- `DEP_CLEAN_CACHE`: Set to `true` to clean Go module cache after download

### Usage Examples

```bash
# Standard build
docker build -t gemini-proxy .

# Build with dependency updates
docker build --build-arg DEP_UPDATE=true -t gemini-proxy .

# Build with dependency updates and clean cache
docker build --build-arg DEP_UPDATE=true --build-arg DEP_CLEAN_CACHE=true -t gemini-proxy .
```

### Image Size Optimization

To minimize image size:
1. Use multi-stage builds (already implemented)
2. Clean module cache after downloading dependencies
3. Use minimal base images (Alpine Linux)

## CI/CD Integration

### Automated Dependency Updates

Consider implementing automated dependency updates using tools like:
- Dependabot
- Renovate Bot

### Security Scanning in CI/CD

Integrate security scanning into your CI/CD pipeline:
```bash
# In your CI pipeline
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

### Dependency Verification

Add dependency verification to CI/CD:
```bash
# In your CI pipeline
go mod verify
go mod tidy
git diff --exit-code go.mod go.sum
```

## Best Practices

1. **Regular Updates**: Schedule regular dependency updates
2. **Security First**: Prioritize security updates
3. **Testing**: Always test after dependency changes
4. **Documentation**: Document significant dependency changes
5. **Minimize Dependencies**: Regularly review and remove unnecessary dependencies
6. **Version Pinning**: Pin to specific versions in production
7. **Monitoring**: Monitor for new vulnerabilities in dependencies

## Conclusion

Effective dependency management is crucial for maintaining a secure, reliable, and maintainable codebase. By following these guidelines, you can ensure that dependencies are managed properly throughout the project lifecycle.