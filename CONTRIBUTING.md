# Contributing to Ticketr

Thank you for your interest in contributing to Ticketr! This document provides guidelines and information for contributors.

## Repository Structure

The Ticketr repository is organized as follows:

```
ticketr/
├── cmd/ticketr/         # CLI entry point and main command
├── internal/            # Internal packages (not for external use)
│   ├── core/           # Core business logic
│   │   ├── domain/     # Domain models and entities
│   │   ├── ports/      # Interface definitions (ports)
│   │   ├── services/   # Core business services
│   │   └── validation/ # Validation logic
│   ├── adapters/       # External integrations (adapters)
│   │   ├── filesystem/ # File I/O operations
│   │   └── jira/       # Jira API client
│   ├── parser/         # Markdown parsing logic
│   ├── renderer/       # Markdown rendering logic
│   └── state/          # State management
├── testdata/           # Test fixtures and sample data
├── examples/           # Example templates and usage
├── docs/              # Project documentation
└── .pm/               # Internal project management artifacts
```

### The `.pm/` Directory

The `.pm/` directory contains internal project tracking and management artifacts that are not part of the public codebase. This includes:
- Evidence of test runs and regression testing
- Internal development notes
- Project phase tracking documents

This directory is gitignored to keep the repository clean and focused on the actual codebase.

## Code Style & Commenting

### GoDoc Comments

All exported types, functions, methods, and variables MUST have clear GoDoc comments. This is a mandatory requirement for code quality.

**Format:**
```go
// TicketService handles the business logic for ticket operations.
// It coordinates between the parser, validator, and Jira adapter.
type TicketService struct {
    // ...
}

// ProcessTickets parses markdown content, validates tickets, and syncs with Jira.
// It returns the updated markdown content with Jira IDs and an error if any step fails.
//
// Parameters:
//   - content: The markdown content containing ticket definitions
//   - options: Configuration options for processing
//
// Returns:
//   - string: Updated markdown content with Jira IDs
//   - error: Any error that occurred during processing
func (s *TicketService) ProcessTickets(content string, options ProcessOptions) (string, error) {
    // ...
}
```

### Inline Comments

Complex or non-obvious logic should be explained with inline comments focusing on the "why" rather than the "what":

```go
// Use exponential backoff to avoid overwhelming the Jira API
// when it's under heavy load
for retries := 0; retries < maxRetries; retries++ {
    // Calculate delay: 2^retries * 100ms with jitter
    delay := time.Duration(1<<retries) * 100 * time.Millisecond
    // ...
}
```

### Variable and Function Names

- Use clear, descriptive names that express intent
- Avoid abbreviations unless they're widely understood (e.g., `ctx` for context)
- Keep functions focused on a single responsibility
- Functions should generally be no longer than 50 lines

## Development Workflow

### Setting Up Your Development Environment

1. **Clone the repository:**
   ```bash
   git clone https://github.com/karolswdev/ticketr.git
   cd ticketr
   ```

2. **Install Go 1.24 or later:**
   Follow the instructions at https://go.dev/doc/install

3. **Install dependencies:**
   ```bash
   go mod download
   ```

4. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your Jira credentials
   ```

### Running Tests

```bash
# Run all tests
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Run tests for a specific package
go test ./internal/parser -v

# Run tests with race detection
go test ./... -race
```

### Building the Project

```bash
# Build the binary
go build -o ticketr cmd/ticketr/main.go

# Build with optimizations
go build -ldflags="-s -w" -o ticketr cmd/ticketr/main.go

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o ticketr-linux cmd/ticketr/main.go
GOOS=darwin GOARCH=amd64 go build -o ticketr-darwin cmd/ticketr/main.go
GOOS=windows GOARCH=amd64 go build -o ticketr.exe cmd/ticketr/main.go
```

## Git Workflow

### Branch Strategy (GitHub Flow)

- `main`: Production-ready code
- Feature branches from `main`: `feat/*`, `fix/*`, `docs/*`, `refactor/*`

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

**Examples:**
```
feat(parser): Add support for epic ticket types

fix(jira): Handle rate limiting with exponential backoff

docs(readme): Update installation instructions for Windows

refactor(services): Extract validation logic to separate package
```

### Pull Request Process

1. Create a feature branch from `main`
2. Make your changes following the code style guidelines
3. Add or update tests as needed
4. Ensure all tests pass
5. Update documentation if applicable
6. Create a pull request with a clear description
7. Address review feedback
8. Keep commits focused; squash if requested

### Community

By participating, you agree to follow the [Code of Conduct](CODE_OF_CONDUCT.md).

## Testing Guidelines

### Test Structure

Tests should follow the Arrange-Act-Assert pattern:

```go
func TestTicketService_ProcessTickets(t *testing.T) {
    // Arrange
    service := NewTicketService(mockParser, mockValidator, mockJira)
    content := "# TICKET: Test ticket"
    
    // Act
    result, err := service.ProcessTickets(content, DefaultOptions)
    
    // Assert
    assert.NoError(t, err)
    assert.Contains(t, result, "[PROJ-123]")
}
```

### Test Coverage

- Aim for at least 80% code coverage
- Focus on testing business logic and edge cases
- Use table-driven tests for multiple scenarios
- Mock external dependencies

## Documentation

### Code Documentation

- All exported symbols must have GoDoc comments
- Complex algorithms should have explanatory comments
- Include examples in GoDoc where helpful

### Project Documentation

- Update README.md for user-facing changes
- Update docs/ for architectural changes
- Keep examples/ up to date with current functionality

## Security

### Handling Credentials

- Never commit credentials or API keys
- Use environment variables for sensitive configuration
- Document required environment variables clearly

### Input Validation

- Always validate user input
- Sanitize data before sending to external services
- Handle errors gracefully without exposing sensitive information

## Getting Help

If you have questions or need help:

1. Check the [documentation](docs/)
2. Search [existing issues](https://github.com/karolswdev/ticketr/issues)
3. Ask in [discussions](https://github.com/karolswdev/ticketr/discussions)
4. Create a new issue with a clear description

## License

By contributing to Ticketr, you agree that your contributions will be licensed under the MIT License.
