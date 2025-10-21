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

## Agent Roles & Development Methodology

Ticketr development follows a **6-Agent Methodology** where specialized agents handle different aspects of the development lifecycle. This methodology ensures consistent quality, comprehensive testing, complete documentation, and architectural integrity.

### The 6 Specialized Agents

#### 1. Builder Agent (Feature Developer)
**Expertise:** Go development, hexagonal architecture, testing patterns, clean code practices

**Responsibilities:**
- Implement production-quality code changes
- Write initial unit tests with >80% coverage for critical paths
- Follow hexagonal architecture (ports & adapters)
- Ensure code quality (`gofmt`, `go vet` clean)

**Deliverables:**
- Working code with tests passing
- Implementation summary (files modified, behaviors added)
- Notes for Verifier (areas needing thorough testing)
- Notes for Scribe (documentation updates needed)

See [`.agents/builder.agent.md`](.agents/builder.agent.md) for complete specification.

#### 2. Verifier Agent (Quality & Test Engineer)
**Expertise:** Testing strategies, quality assurance, regression detection, coverage analysis

**Responsibilities:**
- Extend test coverage to meet targets (>80% critical paths, >50% overall)
- Run full test suite (`go test ./...`)
- Run race detector (`go test -race ./...`)
- Validate requirements compliance
- Check for regressions

**Deliverables:**
- Test execution report (pass/fail counts, coverage metrics)
- Requirements validation matrix
- Regression analysis
- Clear recommendation: APPROVE or REQUEST FIXES

See [`.agents/verifier.agent.md`](.agents/verifier.agent.md) for complete specification.

#### 3. Scribe Agent (Documentation Specialist)
**Expertise:** Technical writing, markdown standards, documentation architecture, user experience writing

**Responsibilities:**
- Update README.md for user-facing changes
- Update REQUIREMENTS.md traceability
- Update ROADMAP.md milestone checkboxes
- Create/update feature guides in `docs/`
- Update examples in `examples/`
- Update CHANGELOG.md

**Deliverables:**
- Comprehensive documentation updates
- Accurate examples (tested against actual CLI)
- Cross-references validated (no broken links)
- Spell-checked, well-formatted markdown

See [`.agents/scribe.agent.md`](.agents/scribe.agent.md) for complete specification.

#### 4. Steward Agent (Architect & Final Approver)
**Expertise:** System architecture, security assessment, requirements governance, technical debt management

**Responsibilities:**
- Architecture compliance review (hexagonal boundaries)
- Security assessment (no secrets, credential management)
- Requirements validation (traceability chain)
- Quality assessment (test coverage, test quality)
- Documentation completeness review
- Phase gate GO/NO-GO decisions

**Deliverables:**
- Comprehensive review report (architecture, security, requirements, quality, documentation)
- Final decision: APPROVE / APPROVE WITH CONDITIONS / REJECT
- Remediation plan (if rejected)

See [`.agents/steward.agent.md`](.agents/steward.agent.md) for complete specification.

#### 5. Director Agent (Control Flow Orchestrator)
**Expertise:** Task decomposition, agent coordination, progress tracking, git workflow, quality gate enforcement

**Responsibilities:**
- Break roadmap milestones into atomic tasks
- Delegate to specialized agents (Builder → Verifier → Scribe → Steward)
- Track progress with TodoWrite (one task in_progress at a time)
- Enforce quality gates (never skip Verifier or Scribe)
- Create git commits with proper attribution (Happy + Claude)

**Deliverables:**
- Milestone completion reports
- Git commits with conventional format
- Quality gates validation
- Blocker escalation

See [`.agents/director.agent.md`](.agents/director.agent.md) for complete specification.

#### 6. TUIUX Agent (TUI/UX Expert)
**Expertise:** Terminal UI design, user experience, visual polish, motion design, accessibility

**Responsibilities:**
- TUI visual and experiential polish
- Animation implementation (spinners, pulses, fades)
- Theme system creation
- Accessibility compliance (motion kill switch, graceful degradation)
- Performance optimization (≤3% CPU for animations)

**Deliverables:**
- Visual polish implementation (effects, widgets, themes)
- Performance benchmarks (<3% CPU assertions)
- Accessibility validation
- Demo programs showcasing features

See [`.agents/tuiux.agent.md`](.agents/tuiux.agent.md) for complete specification.

### Standard Workflow

The 6-agent methodology follows a **sequential workflow**:

```
DIRECTOR: Analyze & Plan
    ↓ (create TodoList)
BUILDER: Implement
    ↓ (code + initial tests)
VERIFIER: Validate
    ↓ (full test suite + coverage)
SCRIBE: Document
    ↓ (update all docs)
(STEWARD): Approve (optional for major changes)
    ↓
DIRECTOR: Commit
    ↓ (git commit with attribution)
```

**Key Rules:**
- ✅ Always run agents sequentially (Builder → Verifier → Scribe)
- ✅ Never skip Verifier (even if Builder tests pass)
- ✅ Never skip Scribe (documentation is mandatory)
- ✅ Invoke Steward for phase gates, major changes, releases
- ✅ Create logical git commits with Happy/Claude attribution

### When to Engage Which Agent

| Situation | Agent to Engage | Reason |
|-----------|----------------|--------|
| Implementing a new feature | Builder | Code implementation specialist |
| After code implementation | Verifier | Validate quality, extend tests, check regressions |
| After tests pass | Scribe | Document the feature for users and developers |
| Phase gate or release | Steward | Final architectural and security approval |
| TUI visual polish | TUIUX | Specialized in terminal UI design and UX |
| Milestone orchestration | Director | Coordinates all agents and enforces workflow |

### Quality Gates Enforced

The 6-agent methodology enforces these quality gates:

1. **Code Quality:** `gofmt` clean, `go vet` clean, hexagonal architecture maintained
2. **Test Coverage:** >80% critical paths, >70% service layer, >50% overall
3. **Test Quality:** No regressions, race detector clean, error paths tested
4. **Documentation:** README updated, guides created, requirements traced, roadmap marked
5. **Security:** No secrets in code, credentials in keychain, `.gitignore` proper
6. **Architecture:** Hexagonal boundaries respected, dependency direction correct

### Additional Resources

- **Director's Handbook:** Complete methodology guide (`docs/DIRECTOR-HANDBOOK.md`)
- **Requirements:** Single source of truth (`REQUIREMENTS.md`)
- **Roadmap:** Milestone tracking (`ROADMAP.md`)
- **Architecture:** Hexagonal architecture patterns (`docs/ARCHITECTURE.md`)

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

### Testing Workspace Features

Ticketr v3.0 introduces workspace management with OS keychain integration. Testing workspace features requires special considerations.

#### Unit Tests (Keychain Mocked)

Most workspace tests use mock keychain adapters for fast, reliable unit testing:

```go
func TestWorkspaceService_Create(t *testing.T) {
    // Setup mock credential store
    mockCredStore := &MockCredentialStore{}
    mockRepo := &MockWorkspaceRepository{}

    service := services.NewWorkspaceService(mockRepo, mockCredStore)

    // Test workspace creation with mocked keychain
    err := service.Create("backend", domain.WorkspaceConfig{
        JiraURL:    "https://company.atlassian.net",
        ProjectKey: "BACK",
        Username:   "user@company.com",
        APIToken:   "token123",
    })

    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}
```

**Run unit tests:**
```bash
go test ./internal/core/services/... -v
go test ./internal/adapters/database/... -v
```

#### Integration Tests (Real Keychain)

Integration tests validate actual keychain operations on your OS:

**macOS:**
```bash
# First run may prompt for keychain access
go test ./internal/adapters/keychain/... -v
```

**Windows:**
```bash
# Requires interactive session
go test ./internal/adapters/keychain/... -v
```

**Linux:**
```bash
# Requires keyring daemon (GNOME Keyring or KWallet)
go test ./internal/adapters/keychain/... -v
```

**Expected behavior:**
- First run may prompt for OS permission
- Subsequent runs should pass without prompts
- Tests clean up keychain entries after completion

#### Graceful Test Skipping

Keychain integration tests skip gracefully when keychain is unavailable:

```go
func TestKeychainStore_Integration(t *testing.T) {
    if os.Getenv("SKIP_KEYCHAIN_TESTS") != "" {
        t.Skip("Skipping keychain integration tests (CI environment)")
    }

    // Test real keychain operations
    store := keychain.NewKeychainStore()
    // ...
}
```

**Skip keychain tests in CI:**
```bash
export SKIP_KEYCHAIN_TESTS=1
go test ./...
```

#### Platform-Specific Testing

**macOS:**
- Keychain Access app must be unlocked
- First test run requires user approval
- Check keychain entries: Open Keychain Access → Search "ticketr"

**Windows:**
- Must run in interactive session (not headless CI)
- Check Credential Manager: Control Panel → User Accounts → Credential Manager
- Search for "ticketr" entries

**Linux:**
- Requires `gnome-keyring-daemon` or `kwalletd5` running
- Check with: `secret-tool search service ticketr`
- Start keyring if needed: `gnome-keyring-daemon --start`

**Troubleshooting:**
```bash
# Linux: Check if keyring daemon is running
ps aux | grep keyring

# Linux: Test keyring access
secret-tool store --label='test' service test account test
secret-tool lookup service test account test
secret-tool clear service test account test
```

#### Manual Testing Checklist

Before submitting workspace-related PRs, manually verify:

**Workspace Creation:**
- [ ] Create workspace with valid credentials
- [ ] Verify credentials stored in OS keychain (check Keychain Access/Credential Manager)
- [ ] Verify workspace appears in database
- [ ] Create duplicate workspace (should fail)
- [ ] Create workspace with invalid credentials (should fail)

**Workspace Switching:**
- [ ] Switch between multiple workspaces
- [ ] Verify credentials loaded correctly
- [ ] Verify LastUsed timestamp updated
- [ ] Switch to non-existent workspace (should fail)

**Workspace Deletion:**
- [ ] Delete workspace
- [ ] Verify credentials removed from keychain
- [ ] Verify workspace removed from database
- [ ] Delete default workspace (should reassign default)
- [ ] Delete only workspace (should fail)

**Security:**
- [ ] Verify credentials NOT in database (only CredentialRef)
- [ ] Verify credentials NOT in logs
- [ ] Verify credentials NOT in error messages
- [ ] Verify credentials cleared from memory after use

**Cross-Platform:**
- [ ] Test on macOS (Keychain Access)
- [ ] Test on Windows (Credential Manager)
- [ ] Test on Linux (GNOME Keyring/KWallet)

#### CI/CD Considerations

GitHub Actions and other CI platforms may not support interactive keychain access:

**GitHub Actions Strategy:**
```yaml
- name: Run tests (skip keychain integration)
  run: |
    export SKIP_KEYCHAIN_TESTS=1
    go test ./...
```

**Local Development:**
```bash
# Run all tests including keychain integration
go test ./...

# Run only keychain integration tests
go test ./internal/adapters/keychain/... -v
```

#### Coverage Requirements

Workspace features must meet coverage thresholds:

| Component | Minimum Coverage | Status |
|-----------|------------------|--------|
| Workspace Repository | 80% | ✅ 80.6% |
| Workspace Service | 70% | ✅ 75.2% |
| Keychain Adapter | 70% | ✅ 87.5% |
| CLI Workspace Commands | 60% | ✅ 68.2% |

**Check workspace coverage:**
```bash
go test ./internal/core/services/... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep -i workspace
```

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
