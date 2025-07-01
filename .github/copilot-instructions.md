# Omnistrate CTL - Make Commands Guide

This file provides context for GitHub Copilot about the available make commands in this Go project.

## How to Use Make Commands

Make is a build automation tool that reads a Makefile to execute commands. The general syntax is:

```bash
make [target]
```

### Basic Usage Examples:
```bash
make build                    # Run the build target
make unit-test                # Run unit tests
make lint                     # Run code quality checks
```

## Build Commands

### `make build`
Builds the CTL binary for the current OS/architecture. Creates binary in `dist/` directory with platform-specific naming.

Example:
```bash
make build                   # Build for current platform
make GOOS=linux GOARCH=amd64 build  # Cross-compile for Linux AMD64
```

### Platform-specific builds:
- `make ctl-linux-amd64` - Build for Linux AMD64
- `make ctl-linux-arm64` - Build for Linux ARM64  
- `make ctl-darwin-amd64` - Build for macOS AMD64
- `make ctl-darwin-arm64` - Build for macOS ARM64 (Apple Silicon)
- `make ctl-windows-amd64` - Build for Windows AMD64
- `make ctl-windows-arm64` - Build for Windows ARM64

### `make ctl`
Builds binaries for all supported platforms (Linux, macOS, Windows on AMD64 and ARM64).

## Testing Commands

### `make unit-test`
Runs unit tests with coverage reporting. Fails if coverage is below threshold (currently 0%).

Example:
```bash
make unit-test               # Run all unit tests
make ARGS="-v" unit-test     # Run with verbose output
```

### `make smoke-test`
Runs smoke tests. Requires environment variables:
- `TEST_EMAIL` - Test user email
- `TEST_PASSWORD` - Test user password

Example:
```bash
TEST_EMAIL=user@example.com TEST_PASSWORD=secret make smoke-test
```

### `make integration-test`
Runs integration tests. Requires same environment variables as smoke tests.

Example:
```bash
TEST_EMAIL=user@example.com TEST_PASSWORD=secret make integration-test
```

## Code Quality Commands

### `make lint`
Runs golangci-lint checks on all Go code. Install with `make lint-install`.

```bash
make lint-install           # Install golangci-lint first
make lint                   # Run linting
```

## Dependency Management

### `make tidy`
Cleans up Go module dependencies (`go mod tidy`).

### `make download`
Downloads all Go module dependencies.

### `make update-dependencies`
Updates all dependencies to latest versions.

Example:
```bash
make update-dependencies     # Update all dependencies to latest
```

### `make update-omnistrate-dependencies`
Updates only Omnistrate-specific dependencies.

### `make check-dependencies`
Validates that no conflicting Omnistrate dependencies are present.

## Documentation

### `make gen-doc`
Generates CLI documentation in markdown format for mkdocs.

### `make pretty`
Formats code using prettier (runs `npx prettier --write .`).

## Utility Commands

### `make all`
Runs the complete build pipeline: tidy, build, unit-test, lint, check-dependencies, gen-doc, pretty.

```bash
make all                    # Full CI/CD pipeline
```

## Common Workflows

### Development workflow:
```bash
make tidy          # Clean dependencies
make build         # Build for current platform
make unit-test     # Run tests
make lint          # Check code quality
```
### Quick development cycle:
```bash
make tidy build unit-test lint   # Chain multiple targets
```

## Environment Variables

You can override these variables when running make commands:

- `GOOS` - Target operating system (default: current OS)
- `GOARCH` - Target architecture (default: current arch)
- `TAG` - Docker tag (default: latest)
- `GIT_USER` - GitHub username (auto-detected via gh CLI)
- `GIT_TOKEN` - GitHub token (auto-detected via gh CLI)
- `TEST_EMAIL` - Email for test authentication
- `TEST_PASSWORD` - Password for test authentication
- `TESTCOVERAGE_THRESHOLD` - Minimum test coverage percentage (default: 0)
- `DOCKER_PLATFORM` - Docker build platform (default: linux/arm64)

Example usage:
```bash
make GOOS=windows GOARCH=amd64 build
make TESTCOVERAGE_THRESHOLD=75 unit-test
make TAG=dev DOCKER_PLATFORM=linux/amd64 docker
```
