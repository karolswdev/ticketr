# Development Guide

This guide provides comprehensive instructions for setting up your development environment, running tests, and contributing to Ticketr.

## Prerequisites

### Required Software

- **Go 1.22 or later**: [Download Go](https://go.dev/dl/)
- **Git**: [Download Git](https://git-scm.com/downloads)
- **Make** (optional): For simplified command execution
- **Docker** (optional): For containerized development

### Recommended Tools

- **VS Code** with Go extension
- **GoLand** IDE
- **golangci-lint**: For code quality checks
- **ko**: For building container images

## Setting Up Your Development Environment

### 1. Clone the Repository

```bash
# Clone via HTTPS
git clone https://github.com/karolswdev/ticketr.git

# Or clone via SSH
git clone git@github.com:karolswdev/ticketr.git

cd ticketr
```

### 2. Install Dependencies

```bash
# Download all Go module dependencies
go mod download

# Verify dependencies
go mod verify

# Tidy up dependencies (remove unused)
go mod tidy
```

### 3. Configure Environment Variables

Create a `.env` file in the project root for local development:

```bash
# Copy the example environment file
cp .env.example .env

# Edit with your JIRA credentials
vim .env
```

Example `.env` content:
```env
JIRA_URL=https://yourcompany.atlassian.net
JIRA_EMAIL=your.email@company.com
JIRA_API_KEY=your-api-token-here
JIRA_PROJECT_KEY=PROJ
JIRA_STORY_TYPE=Task
JIRA_SUBTASK_TYPE=Sub-task
```

### 4. Build the Application

```bash
# Build the binary
go build -o ticketr cmd/ticketr/main.go

# Build with optimizations (smaller binary)
go build -ldflags="-s -w" -o ticketr cmd/ticketr/main.go

# Install to GOPATH/bin
go install ./cmd/ticketr
```

## Running the Test Suite

### Unit Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test -v ./internal/parser

# Run tests with coverage
go test -v -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Integration Tests

```bash
# Run integration tests (requires JIRA connection)
go test -tags=integration ./...

# Run with specific timeout
go test -timeout 30s ./...
```

### Race Condition Detection

```bash
# Run tests with race detector
go test -race ./...

## Local Lint and Security Checks

Install and run the same tools used in CI via the Makefile:

```bash
# One-time: install tools into GOPATH/bin
make tools

# Lint (golangci-lint)
make lint

# Vulnerability scan (govulncheck)
make vuln

# Full check (fmt, vet, lint, test)
make check
```

Ensure `$(go env GOPATH)/bin` is on your `PATH` so the tools are found.

### Toolchain

This repository enforces a Go toolchain via `go.mod`:

```
toolchain go1.24.4
```

Most tools (including `go` itself and GitHub Actions `setup-go`) will automatically install/use this version when `go-version-file: go.mod` is configured. If you want to override locally for a single shell, you can use:

```bash
GOTOOLCHAIN=go1.24.4
```

# Note: This makes tests slower but catches concurrency issues
```

### Benchmark Tests

```bash
# Run benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkParser ./internal/parser

# Run benchmarks with memory profiling
go test -bench=. -benchmem ./...
```

### Test Organization

Tests are organized following Go conventions:

- `*_test.go` files contain tests for the corresponding source file
- `testdata/` directories contain test fixtures
- Integration tests use build tags for isolation
- Table-driven tests for comprehensive coverage

Example test structure:
```go
func TestTicketService_ProcessTickets(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    *ProcessResult
        wantErr bool
    }{
        {
            name:  "valid ticket",
            input: "# TICKET: Test ticket",
            want:  &ProcessResult{TicketsCreated: 1},
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Development Workflow

### Branch Strategy

We follow a Git Flow-inspired branching model:

```
main           - Production-ready code
├── develop    - Integration branch for features
├── feat/*     - Feature branches
├── fix/*      - Bug fix branches
├── docs/*     - Documentation updates
└── refactor/* - Code refactoring
```

### Creating a Feature Branch

```bash
# Start from develop branch
git checkout develop
git pull origin develop

# Create feature branch
git checkout -b feat/add-webhook-support

# Work on your feature
# ...

# Push to remote
git push -u origin feat/add-webhook-support
```

### Commit Message Convention

We follow the Conventional Commits specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `chore`: Build process or auxiliary tool changes
- `perf`: Performance improvements

**Examples:**

```bash
# Feature commit
git commit -m "feat(parser): Add support for epic ticket types"

# Bug fix with detailed description
git commit -m "fix(jira): Handle rate limiting with exponential backoff

The JIRA adapter now implements exponential backoff when
encountering rate limit errors (429 status code).

Fixes #123"

# Documentation update
git commit -m "docs(readme): Update installation instructions for Windows"

# Refactoring
git commit -m "refactor(services): Extract validation logic to separate package"
```

### Code Review Process

1. **Create Pull Request**
   - Target the `develop` branch
   - Fill out the PR template
   - Link related issues

2. **Automated Checks**
   - Tests must pass
   - Code coverage maintained
   - Linting checks pass

3. **Review Requirements**
   - At least one approving review
   - All conversations resolved
   - Branch up-to-date with target

4. **Merge Strategy**
   - Squash and merge for features
   - Merge commit for releases

## Code Style Guidelines

### Go Code Style

We follow the official Go style guide with some additions:

1. **Format**: Use `gofmt` (enforced by CI)
2. **Linting**: Pass `golangci-lint` checks
3. **Imports**: Group in stdlib, external, internal order
4. **Comments**: GoDoc for all exported symbols
5. **Error Handling**: Wrap errors with context
6. **Testing**: Minimum 80% code coverage

### Running Code Quality Checks

```bash
# Format code
go fmt ./...

# Run linter (install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
golangci-lint run

# Run go vet
go vet ./...

# Check for security issues
gosec ./...
```

### Editor Configuration

**.vscode/settings.json:**
```json
{
    "go.formatTool": "gofmt",
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "package",
    "go.testOnSave": true,
    "go.coverOnSave": true,
    "go.coverageDecorator": {
        "type": "gutter",
        "coveredHighlightColor": "rgba(64,128,64,0.5)",
        "uncoveredHighlightColor": "rgba(128,64,64,0.5)"
    }
}
```

## Debugging

### Using Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the application
dlv debug cmd/ticketr/main.go -- push testdata/ticket_simple.md

# Debug a test
dlv test ./internal/parser -- -test.run TestParser
```

### VS Code Debug Configuration

**.vscode/launch.json:**
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug Ticketr",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/ticketr",
            "args": ["push", "testdata/ticket_simple.md"],
            "env": {
                "JIRA_URL": "https://test.atlassian.net",
                "JIRA_EMAIL": "test@example.com",
                "JIRA_API_KEY": "test-key",
                "JIRA_PROJECT_KEY": "TEST"
            }
        }
    ]
}
```

### Logging and Debugging Output

```bash
# Enable verbose logging
ticketr push file.md -v

# Set log level via environment
LOG_LEVEL=debug ticketr push file.md

# Enable HTTP debugging for JIRA requests
HTTP_DEBUG=true ticketr push file.md
```

## Performance Profiling

### CPU Profiling

```bash
# Generate CPU profile
go test -cpuprofile=cpu.prof -bench=. ./internal/parser

# Analyze profile
go tool pprof cpu.prof

# Generate flame graph (requires graphviz)
go tool pprof -http=:8080 cpu.prof
```

### Memory Profiling

```bash
# Generate memory profile
go test -memprofile=mem.prof -bench=. ./internal/parser

# Analyze memory usage
go tool pprof mem.prof

# Check for memory leaks
go test -memprofile=mem.prof -memprofilerate=1 ./...
```

## Docker Development

### Building Docker Image

```bash
# Build image
docker build -t ticketr:dev .

# Build with specific Go version
docker build --build-arg GO_VERSION=1.22 -t ticketr:dev .

# Multi-platform build
docker buildx build --platform linux/amd64,linux/arm64 -t ticketr:dev .
```

### Running in Docker

```bash
# Run with environment file
docker run --rm --env-file .env \
  -v $(pwd)/testdata:/data \
  ticketr:dev push /data/ticket_simple.md

# Interactive shell in container
docker run --rm -it --entrypoint /bin/sh ticketr:dev
```

### Docker Compose Development

```yaml
# docker-compose.dev.yml
version: '3.8'

services:
  ticketr:
    build:
      context: .
      target: development
    volumes:
      - .:/workspace
      - go-modules:/go/pkg/mod
    environment:
      - CGO_ENABLED=0
    command: go run cmd/ticketr/main.go

volumes:
  go-modules:
```

## Continuous Integration

### GitHub Actions Workflow

Our CI pipeline runs on every push and pull request:

1. **Lint**: Code style and quality checks
2. **Test**: Unit and integration tests
3. **Build**: Cross-platform binary compilation
4. **Coverage**: Code coverage reporting
5. **Security**: Vulnerability scanning

### Running CI Locally

```bash
# Install act (GitHub Actions locally)
brew install act  # macOS
# or
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash  # Linux

# Run CI workflow locally
act push

# Run specific job
act -j test
```

## Release Process

### Creating a Release

1. **Update Version**
   ```bash
   # Update version in code
   vim internal/version/version.go
   ```

2. **Update Changelog**
   ```bash
   # Update CHANGELOG.md
   vim CHANGELOG.md
   ```

3. **Create Release Branch**
   ```bash
   git checkout -b release/v1.2.0
   ```

4. **Tag Release**
   ```bash
   git tag -a v1.2.0 -m "Release version 1.2.0"
   git push origin v1.2.0
   ```

5. **Build Release Artifacts**
   ```bash
   # Build for multiple platforms
   make release
   ```

### Versioning Strategy

We follow Semantic Versioning (SemVer):

- **MAJOR**: Incompatible API changes
- **MINOR**: Backwards-compatible functionality
- **PATCH**: Backwards-compatible bug fixes

## Troubleshooting

### Common Issues

#### 1. Module Dependencies

```bash
# Clear module cache
go clean -modcache

# Re-download dependencies
go mod download

# Update dependencies
go get -u ./...
```

#### 2. Build Issues

```bash
# Clean build cache
go clean -cache

# Rebuild with verbose output
go build -v -x ./cmd/ticketr
```

#### 3. Test Failures

```bash
# Run specific failing test
go test -v -run TestName ./package

# Skip integration tests
go test -short ./...

# Increase timeout
go test -timeout 60s ./...
```

### Getting Help

- Check the [FAQ](https://github.com/karolswdev/ticketr/wiki/FAQ)
- Search [existing issues](https://github.com/karolswdev/ticketr/issues)
- Join our [Discord server](https://discord.gg/ticketr)
- Create a [new issue](https://github.com/karolswdev/ticketr/issues/new)

## Contributing

Please read our [Contributing Guide](../CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
