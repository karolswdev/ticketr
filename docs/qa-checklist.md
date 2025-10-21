# QA Checklist - Ticketr Quality Assurance Guide

## Purpose

This guide provides checklists for maintaining code quality at different stages of the development workflow. Follow these checklists before commits, pull requests, and releases to ensure consistent quality standards.

---

## Pre-Commit Checklist

Run these checks before committing code to ensure basic quality standards are met.

### 1. Code Formatting

**Command:**
```bash
gofmt -l .
```

**Expected Output:**
- Empty output (no files listed)

**If Failed:**
```bash
# Format all Go files
gofmt -w .

# Verify formatting
gofmt -l .
```

**Why:** Ensures consistent code style across the project and prevents formatting-related diffs.

---

### 2. Go Vet (Static Analysis)

**Command:**
```bash
go vet ./...
```

**Expected Output:**
- No errors or warnings

**If Failed:**
- Review reported issues
- Fix suspicious constructs, unreachable code, or common mistakes
- Re-run `go vet ./...` until clean

**Why:** Catches common programming errors that the compiler might miss.

---

### 3. Build Verification

**Command:**
```bash
go build ./...
```

**Expected Output:**
- No errors
- Successful compilation of all packages

**If Failed:**
- Fix compilation errors
- Check import paths and package dependencies
- Run `go mod tidy` if dependencies are missing

**Why:** Ensures code compiles successfully across all packages before committing.

---

### 4. Unit Tests

**Command:**
```bash
go test ./... -v
```

**Expected Output:**
- All tests passing (105/106 expected, 3 skipped JIRA integration tests)
- No test failures or panics

**If Failed:**
- Review test output for specific failures
- Fix failing tests or broken functionality
- Add tests for new features
- Run single test for debugging: `go test -v -run TestName ./path/to/package`

**Why:** Validates that new changes don't break existing functionality.

---

### 5. Quick Coverage Check

**Command:**
```bash
go test ./... -cover
```

**Expected Output:**
- Coverage percentage displayed per package
- Overall coverage above 50% threshold

**If Failed:**
- Add tests for uncovered code paths
- Focus on critical business logic
- Review coverage report: `go tool cover -html=coverage.out`

**Why:** Ensures sufficient test coverage for new code.

---

## Pre-Pull Request Checklist

Run these comprehensive checks before creating a pull request.

### 1. Run Quality Script

**Command:**
```bash
bash scripts/quality.sh
```

**Expected Output:**
- All 6 quality checks passing:
  1. go vet
  2. Code formatting
  3. Build check
  4. Tests with coverage
  5. Coverage threshold (50%)
  6. Staticcheck (if installed)
- Coverage report saved to `coverage.out`

**If Failed:**
- Review specific check that failed
- Follow error messages and fix issues
- Re-run quality script until all checks pass

**Why:** Comprehensive automated quality validation before code review.

---

### 2. Documentation Updates

**Checklist:**
- [ ] Updated README.md if user-facing changes
- [ ] Updated development/REQUIREMENTS.md if requirements changed
- [ ] Updated development/ROADMAP.md milestone status
- [ ] Added inline code comments for complex logic
- [ ] Created/updated docs/ files for new features

**If Needed:**
- Check existing docs style in `docs/` directory
- Follow markdown formatting conventions
- Cross-reference related documentation
- Include code examples for new features

**Why:** Keeps documentation in sync with code changes.

---

### 3. development/ROADMAP.md Update

**Checklist:**
- [ ] Marked completed tasks with checkmarks
- [ ] Updated milestone status if applicable
- [ ] Noted any blockers or dependencies
- [ ] Moved completed milestones to appropriate section

**Why:** Tracks project progress and communicates status to team.

---

### 4. Integration Test Plan (if applicable)

**When Required:**
- Changes to JIRA adapter (`internal/adapters/jira/`)
- Changes to core services (`internal/core/services/`)
- Changes to field inheritance logic

**Checklist:**
- [ ] Reviewed `docs/integration-testing-guide.md`
- [ ] Identified test scenarios affected by changes
- [ ] Prepared test markdown files
- [ ] Documented expected behavior changes

**Why:** Ensures changes work correctly with real JIRA instances.

---

### 5. Git History Cleanup

**Commands:**
```bash
# Check commit history
git log --oneline -5

# Check for uncommitted changes
git status
```

**Checklist:**
- [ ] Commit messages are descriptive
- [ ] No WIP or debug commits in history
- [ ] No untracked files that should be committed
- [ ] No sensitive data in commits (.env, API keys, etc.)

**Why:** Maintains clean git history for easier code review and debugging.

---

## Pre-Release Checklist

Run these checks before creating a release or merging to main branch.

### 1. Full Test Suite

**Command:**
```bash
go test ./... -v -cover -coverprofile=coverage.out
```

**Expected Output:**
- All tests passing (105/106 expected)
- Coverage above 50% threshold
- No skipped tests except JIRA integration tests (3 expected)

**Why:** Final validation that all functionality works correctly.

---

### 2. Build Verification

**Command:**
```bash
# Build binary
go build -o ticketr ./cmd/ticketr

# Test binary
./ticketr --version
./ticketr --help
```

**Expected Output:**
- Binary builds successfully
- Version displays correctly
- Help text shows all commands

**Why:** Ensures distributable binary works as expected.

---

### 3. Integration Testing

**Requirements:**
- Real JIRA instance access
- `.env` file configured
- Test project available

**Commands:**
```bash
# Run integration tests (requires .env)
go test ./internal/adapters/jira -v

# Test push workflow
ticketr push examples/simple-ticket.md

# Test pull workflow
ticketr pull PROJ-123
```

**Expected Outcome:**
- Integration tests pass (not skipped)
- Push creates tickets in JIRA
- Pull retrieves tickets from JIRA
- Field inheritance works correctly

**Reference:** See `docs/integration-testing-guide.md` for detailed scenarios.

**Why:** Validates real-world usage with JIRA API.

---

### 4. CI/CD Pipeline

**Checklist:**
- [ ] All GitHub Actions workflows passing
- [ ] No failing jobs in CI pipeline
- [ ] Coverage reports uploaded successfully
- [ ] Build artifacts generated

**Commands:**
```bash
# Check CI status
gh pr checks

# View workflow runs
gh run list --limit 5
```

**Reference:** See `docs/ci.md` for CI/CD pipeline details.

**Why:** Ensures automated checks pass in clean environment.

---

### 5. Documentation Review

**Checklist:**
- [ ] README.md is up-to-date
- [ ] All docs/ files reflect current functionality
- [ ] Code examples in documentation are tested
- [ ] Migration guides updated if breaking changes
- [ ] CHANGELOG.md updated with release notes

**Why:** Users rely on documentation for onboarding and troubleshooting.

---

### 6. Dependency Audit

**Commands:**
```bash
# Check for outdated dependencies
go list -u -m all

# Verify go.mod/go.sum are tidy
go mod tidy
git diff go.mod go.sum
```

**Expected Output:**
- No uncommitted changes to go.mod/go.sum
- No critical security vulnerabilities

**Why:** Keeps dependencies secure and up-to-date.

---

## Troubleshooting Common Issues

### Issue: "gofmt reports files not formatted"

**Solution:**
```bash
# Format all files
gofmt -w .

# Verify
gofmt -l .
```

---

### Issue: "Tests fail with 'cannot find module'"

**Solution:**
```bash
# Tidy dependencies
go mod tidy

# Download dependencies
go mod download

# Re-run tests
go test ./...
```

---

### Issue: "Coverage below threshold"

**Solution:**
```bash
# Generate coverage report
go test ./... -coverprofile=coverage.out

# View coverage by package
go tool cover -func=coverage.out

# View HTML coverage report
go tool cover -html=coverage.out

# Add tests for low-coverage packages
```

---

### Issue: "Integration tests skipped"

**Cause:** Missing `.env` file or invalid credentials

**Solution:**
```bash
# Copy example file
cp .env.example .env

# Edit with real values
# JIRA_URL=https://yourcompany.atlassian.net
# JIRA_EMAIL=your.email@company.com
# JIRA_API_KEY=your-api-token
# JIRA_PROJECT_KEY=TEST

# Re-run integration tests
go test ./internal/adapters/jira -v
```

---

### Issue: "quality.sh fails but individual commands pass"

**Cause:** Script uses `set -e` and exits on first error

**Solution:**
```bash
# Run script with bash -x for debugging
bash -x scripts/quality.sh

# Run individual checks to identify failure point
go vet ./...
gofmt -l .
go build ./...
go test ./...
```

---

## Quick Reference Commands

```bash
# Pre-commit quick check
gofmt -l . && go vet ./... && go test ./...

# Pre-PR full quality check
bash scripts/quality.sh

# Pre-release comprehensive check
go test ./... -v -cover && go build -o ticketr ./cmd/ticketr && ./ticketr --version

# View coverage report
go tool cover -html=coverage.out

# Run specific test
go test -v -run TestName ./internal/parser

# Run integration tests only
go test ./internal/adapters/jira -v

# Check CI status
gh pr checks
```

---

## Related Documentation

- **Quality Script:** `scripts/quality.sh` - Automated quality checks
- **Contributing Guide:** `CONTRIBUTING.md` - Quality gates and CI/CD section
- **CI/CD Pipeline:** `docs/ci.md` - Continuous integration details
- **Smoke Tests:** `tests/smoke/README.md` - End-to-end smoke test documentation
- **Integration Testing:** `docs/integration-testing-guide.md` - JIRA integration test scenarios
- **State Management:** `docs/state-management.md` - State file operations
- **Migration Guides:** `docs/archive/` - Legacy migration guides (archived)

---

**Remember:** Quality checks are not obstacles - they're safety nets that catch issues before users do. Run them early and often!
