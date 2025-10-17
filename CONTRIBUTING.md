# Contributing to Ticketr

Thank you for your interest in contributing to Ticketr! This guide will help you get started.

## Development Setup

### Prerequisites
- Go 1.21+ installed
- JIRA instance for integration testing (optional)
- Git for version control

### Clone and Build
```bash
git clone https://github.com/yourorg/ticktr.git
cd ticktr
go build -o ticketr cmd/ticketr/main.go
```

### Run Tests
```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run specific package
go test ./internal/core/services/... -v

# Run with coverage
go test ./... -cover
```

**Expected output:**
- Total: 69 tests
- Passed: 66
- Skipped: 3 (JIRA integration tests)
- Failed: 0

## Architecture

Ticketr follows **Hexagonal Architecture (Ports & Adapters)**:

```
internal/
├── core/
│   ├── domain/          # Domain models (Ticket, Task)
│   ├── ports/           # Interfaces (Repository, JiraPort)
│   └── services/        # Business logic (PushService, PullService)
├── adapters/
│   ├── jira/            # JIRA API integration
│   └── filesystem/      # Markdown file I/O
├── parser/              # Markdown parser
├── state/               # State management (.ticketr.state)
├── logging/             # File logging
└── migration/           # Legacy format migration
```

**Key principles:**
- Domain logic is independent of external systems
- Adapters implement port interfaces
- Services orchestrate domain operations
- Tests mock ports for isolation

For comprehensive architecture documentation including data flows, design decisions, and component details, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## Testing Guidelines

### Test Organization

Tests are co-located with source files:
- `foo.go` → `foo_test.go`
- Use table-driven tests for multiple scenarios
- Use `t.TempDir()` for file operations
- Mock external dependencies (JIRA, filesystem)

### Test Naming Convention

```go
func TestServiceName_MethodName_Scenario(t *testing.T)
```

Examples:
- `TestPushService_PushTickets_SkipsUnchangedTickets`
- `TestParser_ParseTickets_RejectsLegacyFormat`
- `TestStateManager_CalculateHash_DeterministicOutput`

### Writing Tests

**Good test:**
```go
func TestPushService_FieldInheritance_ParentFields(t *testing.T) {
    // Setup
    mockRepo := &MockRepository{}
    mockJira := &MockJiraPort{}
    stateManager := state.NewStateManager(t.TempDir() + "/.ticketr.state")
    service := services.NewPushService(mockRepo, mockJira, stateManager)

    // Test data
    ticket := domain.Ticket{
        Title: "Parent",
        CustomFields: map[string]string{
            "Sprint": "Sprint 1",
        },
        Tasks: []domain.Task{
            {Title: "Child"},
        },
    }

    // Execute
    _, err := service.PushTickets("test.md", services.ProcessOptions{})

    // Verify
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if mockJira.LastCreatedTask.CustomFields["Sprint"] != "Sprint 1" {
        t.Errorf("expected Sprint inherited, got %v", mockJira.LastCreatedTask.CustomFields)
    }
}
```

### Integration Tests

Integration tests require JIRA credentials:
```bash
export JIRA_URL="https://yourcompany.atlassian.net"
export JIRA_EMAIL="test@example.com"
export JIRA_API_KEY="your_api_token"
export JIRA_PROJECT_KEY="TEST"

go test ./internal/adapters/jira/... -v
```

**Note:** Integration tests are skipped by default when credentials are missing.

## Quality Gates & CI/CD

Ticketr uses automated quality gates to ensure code quality. Before submitting a PR, ensure all checks pass.

### Running Quality Checks Locally

The fastest way to verify your changes:

```bash
bash scripts/quality.sh
```

This script runs all 7 quality checks:
1. go vet (static analysis)
2. Code formatting (gofmt)
3. Build verification
4. Full test suite with coverage
5. Coverage threshold check (50% minimum)
6. staticcheck (advanced static analysis)
7. go.mod tidiness

**Expected output:** All checks passing

### CI/CD Pipeline

Every push to `main`, `feat/**` branches, and all PRs trigger GitHub Actions CI:

| Job | Purpose | Run Time |
|-----|---------|----------|
| **build** | Verify compilation across OS/Go versions | ~2 min |
| **test** | Run test suite with race detector | ~3 min |
| **coverage** | Check coverage threshold (50%+) | ~3 min |
| **lint** | Run go vet, gofmt, staticcheck | ~2 min |
| **smoke-tests** | Execute end-to-end CLI workflows | ~1 min |

**Total pipeline duration:** ~11 minutes

See [`docs/ci.md`](docs/ci.md) for detailed CI documentation.

### Smoke Tests

Run end-to-end CLI workflow tests:

```bash
bash tests/smoke/smoke_test.sh
```

**Scenarios tested:**
1. Legacy file migration
2. Push dry-run validation
3. Pull with missing file
4. State file creation and persistence
5. Log file creation
6. Help command functionality
7. Concurrent file operations safety

**Expected:** 7/7 scenarios passing, 13/13 checks passing

See [`tests/smoke/README.md`](tests/smoke/README.md) for detailed smoke test documentation.

### Quality Checklist

Before creating a PR, verify:

- [ ] All tests pass: `go test ./...`
- [ ] Coverage ≥ 50%: `go test ./... -coverprofile=coverage.out && go tool cover -func coverage.out | tail -1`
- [ ] Code formatted: `gofmt -l .` returns empty
- [ ] No vet issues: `go vet ./...` returns clean
- [ ] Quality script passes: `bash scripts/quality.sh`
- [ ] Smoke tests pass: `bash tests/smoke/smoke_test.sh`
- [ ] Documentation updated (if user-facing changes)
- [ ] ROADMAP updated (if milestone work)

See [`docs/qa-checklist.md`](docs/qa-checklist.md) for complete checklists (pre-commit, pre-PR, pre-release).

## Logging and Debugging

### Execution Logs

Ticketr automatically logs to `.ticketr/logs/`:
- Logs created for every push/pull operation
- Format: `2025-10-16_14-30-00.log`
- Sensitive data automatically redacted
- Kept last 10 files (auto-rotation)

### Viewing Logs

```bash
# Latest log
ls -lt .ticketr/logs/ | head -1

# View specific log
cat .ticketr/logs/2025-10-16_14-30-00.log
```

### Log Format

```
[2025-10-16 14:30:00] ========================================
[2025-10-16 14:30:00] PUSH COMMAND
[2025-10-16 14:30:00] ========================================
[2025-10-16 14:30:00] Input file: my-tickets.md
[2025-10-16 14:30:02] Created ticket 'User Auth' with Jira ID: PROJ-123
```

### Sensitive Data Redaction

The following patterns are automatically redacted in logs:
- API keys (replaced with `[REDACTED]`)
- Email addresses (replaced with `[REDACTED]`)
- Passwords (replaced with `[REDACTED]`)

Implementation: `internal/logging/redactor.go`

## Pull Request Guidelines

### Before Submitting

1. **Run tests:** `go test ./...` (all must pass)
2. **Run linter:** `go vet ./...`
3. **Format code:** `go fmt ./...`
4. **Update tests:** Add tests for new features
5. **Update docs:** Update README if user-facing changes
6. **Update ROADMAP:** Mark milestones complete

### PR Description Template

```markdown
## Summary
Brief description of changes.

## Changes
- Added/Modified/Fixed X
- Updated Y

## Testing
- Added TestServiceName_Scenario
- All 69+ tests passing
- Manual testing: [describe]

## Documentation
- Updated README.md section X
- Added example to docs/

## Related Issues
Fixes #123
```

### Code Review Checklist

Reviewers check for:
- [ ] Tests added/updated for new functionality
- [ ] All tests passing (69+ tests)
- [ ] Code follows existing patterns
- [ ] No sensitive data in code/logs
- [ ] Documentation updated if user-facing
- [ ] Commit messages descriptive

## Documentation Standards

All documentation must follow the [Documentation Style Guide](docs/style-guide.md).

### Quick Reference

**User-Facing Docs (README.md):**
- Clear, concise language
- Include code examples with expected output
- Cross-reference related sections
- Use real-world scenarios

**Developer Docs (docs/):**
- Technical detail appropriate for contributors
- Include implementation notes
- Reference code with file:line format
- Explain design decisions
- Provide troubleshooting guidance

**Code Comments:**
- Explain "why", not "what"
- Use GoDoc format for exported functions
- Keep comments up-to-date with code
- Add TODO comments for future work

**Markdown Formatting:**
- Use kebab-case file names (`migration-guide.md`)
- Always specify language hints in code blocks (```bash, ```go, etc.)
- Use relative links for internal documentation
- Follow heading hierarchy (single H1, proper H2/H3/H4 nesting)
- See [docs/style-guide.md](docs/style-guide.md) for complete standards

## Release Process

(To be defined in Milestone 13)

Planned release workflow:
1. Tag release: `git tag v1.0.0`
2. Generate changelog
3. Build binaries for platforms
4. Publish GitHub release
5. Update documentation

## Getting Help

- **Issues:** Open an issue on GitHub
- **Discussions:** Use GitHub Discussions for questions
- **Slack:** (if available) Join #ticketr channel

## Code of Conduct

(To be added in Milestone 13)

Be respectful, inclusive, and constructive in all interactions.

## License

(Check LICENSE file in repository root)
