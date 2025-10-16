# CI/CD Pipeline - Ticketr Continuous Integration

## Overview

Ticketr uses GitHub Actions for continuous integration and quality assurance. The CI pipeline runs automatically on every push and pull request to validate code quality, run tests, and ensure cross-platform compatibility.

**Workflow File:** `.github/workflows/ci.yml`

**Triggers:**
- Push to `main` branch
- Push to `feat/**` branches
- Pull requests targeting `main` branch

---

## Pipeline Jobs

The CI pipeline consists of 5 parallel jobs that validate different aspects of the codebase:

### 1. Build Job

**Purpose:** Validates that the application builds successfully across multiple platforms and Go versions.

**Configuration:**
- **Runs on:** Ubuntu Latest, macOS Latest (matrix)
- **Go Versions:** 1.21, 1.22, 1.23 (matrix)
- **Total combinations:** 6 (2 OS × 3 Go versions)

**Steps:**
1. Checkout code
2. Set up Go with specified version
3. Cache Go modules for faster builds
4. Download dependencies
5. Build application binary
6. Verify build artifact exists

**Success Criteria:**
- Binary builds without errors on all platforms
- Build artifact (`ticketr` or `ticketr.exe`) is created

**Typical Duration:** 2-3 minutes per matrix combination

---

### 2. Test Job

**Purpose:** Runs the full test suite with race detection and generates coverage reports.

**Configuration:**
- **Runs on:** Ubuntu Latest
- **Go Version:** 1.22 (stable)

**Steps:**
1. Checkout code
2. Set up Go 1.22
3. Cache Go modules
4. Download dependencies
5. Run tests with race detector and coverage
6. Upload coverage artifact for downstream jobs

**Commands:**
```bash
go test ./... -v -race -coverprofile=coverage.out -covermode=atomic
```

**Success Criteria:**
- All tests pass (105/106 expected, 3 JIRA integration tests skipped)
- No race conditions detected
- Coverage report generated

**Expected Output:**
```
ok      github.com/karolswdev/ticktr/cmd/ticketr              0.123s  coverage: 65.2% of statements
ok      github.com/karolswdev/ticktr/internal/core/domain     0.089s  coverage: 88.4% of statements
ok      github.com/karolswdev/ticktr/internal/core/services   0.234s  coverage: 91.7% of statements
...
PASS
```

**Typical Duration:** 1-2 minutes

---

### 3. Coverage Job

**Purpose:** Analyzes test coverage and enforces minimum coverage threshold.

**Configuration:**
- **Runs on:** Ubuntu Latest
- **Go Version:** 1.22
- **Depends on:** Test job (needs coverage artifact)

**Steps:**
1. Checkout code
2. Set up Go
3. Download coverage artifact from test job
4. Generate coverage report
5. Check coverage threshold (50% minimum)
6. Upload coverage report as artifact

**Coverage Threshold:** 50%

**Commands:**
```bash
go tool cover -func=coverage.out > coverage.txt
COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
```

**Success Criteria:**
- Total coverage ≥ 50%
- Coverage report generated and uploaded

**If Failed:**
- Review coverage report to identify low-coverage packages
- Add tests for critical uncovered code
- See `docs/qa-checklist.md` for coverage improvement tips

**Artifact Retention:** 30 days (coverage-text-report)

**Typical Duration:** 30 seconds

---

### 4. Lint Job

**Purpose:** Enforces code quality standards through static analysis and formatting checks.

**Configuration:**
- **Runs on:** Ubuntu Latest
- **Go Version:** 1.22

**Steps:**
1. Checkout code
2. Set up Go
3. Cache Go modules
4. Download dependencies
5. Run `go vet` (static analysis)
6. Check code formatting with `gofmt`
7. Install and run `staticcheck`
8. Verify `go.mod` and `go.sum` tidiness

**Checks Performed:**

#### go vet
```bash
go vet ./...
```
- Detects suspicious constructs
- Catches unreachable code
- Identifies common mistakes

#### Code Formatting
```bash
gofmt -l .
```
- Ensures all files follow Go formatting standards
- Fails if any files are not formatted

#### staticcheck
```bash
staticcheck ./...
```
- Advanced static analysis
- Catches bugs, performance issues, and style violations
- More comprehensive than `go vet`

#### Module Tidiness
```bash
go mod tidy
git diff go.mod go.sum
```
- Ensures `go.mod` and `go.sum` are up-to-date
- Prevents dependency drift

**Success Criteria:**
- No `go vet` warnings
- All files properly formatted
- No `staticcheck` issues
- `go.mod` and `go.sum` are tidy

**Typical Duration:** 1-2 minutes

---

### 5. Smoke Tests Job

**Purpose:** Validates end-to-end functionality with real-world usage scenarios.

**Configuration:**
- **Runs on:** Ubuntu Latest
- **Go Version:** 1.22
- **Depends on:** Build job

**Steps:**
1. Checkout code
2. Set up Go
3. Build application binary
4. Make smoke test script executable
5. Run smoke tests
6. Upload smoke test logs (on success or failure)

**Test Script:** `tests/smoke/smoke_test.sh`

**Smoke Test Scenarios:**
- Binary execution (--version, --help)
- Basic command validation
- State file operations
- Parser functionality
- Error handling

**Success Criteria:**
- All smoke tests pass
- No runtime errors or panics
- Logs generated correctly

**Artifact Retention:** 7 days (smoke-test-logs)

**Typical Duration:** 30 seconds

---

## Matrix Strategy

The build job uses a matrix strategy to test across multiple configurations:

```yaml
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest]
    go-version: ['1.21', '1.22', '1.23']
```

**Matrix Combinations:**
1. Ubuntu + Go 1.21
2. Ubuntu + Go 1.22
3. Ubuntu + Go 1.23
4. macOS + Go 1.21
5. macOS + Go 1.22
6. macOS + Go 1.23

**Why Matrix Testing?**
- Ensures cross-platform compatibility (Linux, macOS)
- Validates support for multiple Go versions
- Catches platform-specific issues early
- Provides confidence for diverse user environments

**Note:** Windows is not currently in the matrix but can be added if needed.

---

## How to Read CI Results

### Viewing CI Status

**In GitHub UI:**
1. Navigate to "Actions" tab in repository
2. Click on workflow run to see job details
3. Expand individual jobs to see step-by-step logs

**In Pull Request:**
- CI status appears as checks at bottom of PR
- Green checkmark = all jobs passed
- Red X = at least one job failed
- Yellow circle = jobs still running

**Via GitHub CLI:**
```bash
# Check PR status
gh pr checks

# List recent workflow runs
gh run list --limit 10

# View specific run details
gh run view <run-id>

# Watch live workflow execution
gh run watch
```

---

### Interpreting Job Results

#### All Jobs Green ✅
- Code is ready to merge
- All quality checks passed
- No action required

#### Build Job Failed ❌
**Common Causes:**
- Compilation errors
- Missing dependencies
- Platform-specific code issues

**Fix:**
```bash
# Test locally on same platform
go build ./...

# Check for missing dependencies
go mod tidy
```

#### Test Job Failed ❌
**Common Causes:**
- Failing unit tests
- Race conditions detected
- Test panics

**Fix:**
```bash
# Run tests locally
go test ./... -v -race

# Run specific failing test
go test -v -run TestName ./package
```

#### Coverage Job Failed ❌
**Common Causes:**
- Coverage below 50% threshold
- Missing tests for new code

**Fix:**
```bash
# Generate local coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Add tests for low-coverage packages
```

#### Lint Job Failed ❌
**Common Causes:**
- Code not formatted
- `go vet` warnings
- `staticcheck` issues
- `go.mod` not tidy

**Fix:**
```bash
# Format code
gofmt -w .

# Run vet
go vet ./...

# Install and run staticcheck
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...

# Tidy modules
go mod tidy
```

#### Smoke Tests Failed ❌
**Common Causes:**
- Binary execution errors
- Command-line argument issues
- State file corruption

**Fix:**
```bash
# Build and test locally
go build -o ticketr ./cmd/ticketr
./ticketr --version
./ticketr --help

# Run smoke tests locally
bash tests/smoke/smoke_test.sh
```

---

## Troubleshooting CI Failures

### Strategy 1: Reproduce Locally

**For Build/Test/Lint failures:**
```bash
# Run the exact CI commands locally
go build -v ./cmd/ticketr
go test ./... -v -race -coverprofile=coverage.out -covermode=atomic
go vet ./...
gofmt -l .
staticcheck ./...
go mod tidy
```

### Strategy 2: Check CI Logs

**Where to look:**
1. Click on failed job in GitHub Actions
2. Expand the failing step
3. Read error messages and stack traces
4. Look for "FAIL" or "ERROR" keywords

**Example failure log:**
```
--- FAIL: TestParser_RejectsLegacyStoryFormat (0.00s)
    parser_test.go:113: Expected error message to contain 'Legacy', got: legacy '# STORY:' format detected
FAIL
FAIL    github.com/karolswdev/ticktr/internal/parser  0.123s
```

### Strategy 3: Download Artifacts

**Available artifacts:**
- coverage-report (coverage.out)
- coverage-text-report (coverage.txt)
- smoke-test-logs (test output and execution logs)

**How to download:**
1. Navigate to failed workflow run
2. Scroll to "Artifacts" section
3. Click artifact name to download

### Strategy 4: Run Quality Script Locally

**Comprehensive local validation:**
```bash
# Run all quality checks
bash scripts/quality.sh

# This runs the same checks as CI:
# 1. go vet
# 2. Code formatting
# 3. Build check
# 4. Tests with coverage
# 5. Coverage threshold
# 6. staticcheck
# 7. Module tidiness
```

**If quality.sh passes but CI fails:**
- Check for platform-specific issues (Linux vs macOS)
- Verify Go version matches CI (1.22)
- Ensure all files are committed
- Check for .gitignore exclusions

---

## Running CI Checks Locally

### Prerequisites

**Install required tools:**
```bash
# Install Go 1.22 or later
go version

# Install staticcheck
go install honnef.co/go/tools/cmd/staticcheck@latest

# Add to PATH if needed
export PATH=$PATH:$(go env GOPATH)/bin
```

### Run All Checks

**Option 1: Quality script (recommended)**
```bash
bash scripts/quality.sh
```

**Option 2: Manual commands**
```bash
# Build
go build -v ./cmd/ticketr

# Test
go test ./... -v -race -coverprofile=coverage.out -covermode=atomic

# Coverage
go tool cover -func=coverage.out
COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
echo "Coverage: $COVERAGE%"

# Lint
go vet ./...
gofmt -l .
staticcheck ./...
go mod tidy
git diff go.mod go.sum

# Smoke tests
bash tests/smoke/smoke_test.sh
```

### Run Specific Job Checks

**Build job:**
```bash
go build -v ./cmd/ticketr
```

**Test job:**
```bash
go test ./... -v -race -coverprofile=coverage.out -covermode=atomic
```

**Coverage job:**
```bash
go tool cover -func=coverage.out > coverage.txt
cat coverage.txt
```

**Lint job:**
```bash
go vet ./...
gofmt -l .
staticcheck ./...
go mod tidy && git diff go.mod go.sum
```

**Smoke tests job:**
```bash
go build -o ticketr ./cmd/ticketr
chmod +x tests/smoke/smoke_test.sh
./tests/smoke/smoke_test.sh
```

---

## Badge Integration

Add CI status badge to README.md:

```markdown
![CI](https://github.com/karolswdev/ticktr/actions/workflows/ci.yml/badge.svg)
```

**Badge States:**
- Green "passing" = All jobs passed
- Red "failing" = At least one job failed
- Gray "no status" = No recent runs

---

## CI Best Practices

### Before Pushing

1. **Run quality script:**
   ```bash
   bash scripts/quality.sh
   ```

2. **Ensure tests pass locally:**
   ```bash
   go test ./... -v
   ```

3. **Format code:**
   ```bash
   gofmt -w .
   ```

4. **Commit formatted code:**
   ```bash
   git add .
   git commit -m "feat: your changes"
   ```

### During PR Review

1. **Monitor CI status** in PR checks section
2. **Fix failures immediately** - don't wait for review
3. **Check coverage reports** in artifacts
4. **Respond to reviewer feedback** and re-run CI

### After CI Passes

1. **Review coverage report** to ensure new code is tested
2. **Check smoke test logs** to verify end-to-end functionality
3. **Merge when all checks are green** and approved

---

## Extending the CI Pipeline

### Adding New Jobs

**Example: Security scanning**
```yaml
security:
  name: Security Scan
  runs-on: ubuntu-latest

  steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Run Gosec
      uses: securego/gosec@master
      with:
        args: './...'
```

### Adding Matrix Variants

**Example: Add Windows platform**
```yaml
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest, windows-latest]
    go-version: ['1.21', '1.22', '1.23']
```

### Adding Integration Tests

**Example: JIRA integration tests (requires secrets)**
```yaml
integration-tests:
  name: Integration Tests
  runs-on: ubuntu-latest

  steps:
    - name: Checkout code
      uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.22'

    - name: Run integration tests
      env:
        JIRA_URL: ${{ secrets.JIRA_URL }}
        JIRA_EMAIL: ${{ secrets.JIRA_EMAIL }}
        JIRA_API_KEY: ${{ secrets.JIRA_API_KEY }}
        JIRA_PROJECT_KEY: ${{ secrets.JIRA_PROJECT_KEY }}
      run: go test ./internal/adapters/jira -v
```

**Note:** Requires setting up GitHub Secrets in repository settings.

---

## Related Documentation

- **QA Checklist:** `docs/qa-checklist.md` - Pre-commit, pre-PR, and pre-release checks
- **Contributing Guide:** `CONTRIBUTING.md` - Quality gates and CI/CD section
- **Quality Script:** `scripts/quality.sh` - Local quality validation script
- **Smoke Tests:** `tests/smoke/README.md` - End-to-end smoke test documentation
- **Integration Testing:** `docs/integration-testing-guide.md` - JIRA integration test scenarios
- **Workflow File:** `.github/workflows/ci.yml` - Complete CI configuration

---

## CI Performance Optimization

### Caching Strategy

The CI pipeline uses GitHub Actions cache to speed up builds:

```yaml
- name: Cache Go modules
  uses: actions/cache@v4
  with:
    path: |
      ~/.cache/go-build
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ matrix.go-version }}-${{ hashFiles('**/go.sum') }}
```

**Benefits:**
- Faster dependency downloads (30-60 second savings)
- Reduced network load
- Consistent builds with cached dependencies

**Cache Invalidation:**
- Cache is invalidated when `go.sum` changes
- Separate caches per OS and Go version

### Typical Pipeline Duration

**Full pipeline (all jobs in parallel):**
- **Fastest:** 2-3 minutes (with cache hits)
- **Typical:** 3-4 minutes
- **Slowest:** 5-6 minutes (cold cache, all matrix combinations)

**Individual job durations:**
- Build: 1-3 minutes (per matrix combination)
- Test: 1-2 minutes
- Coverage: 30 seconds
- Lint: 1-2 minutes
- Smoke Tests: 30 seconds

---

**Remember:** CI is your safety net. Green checks don't guarantee perfect code, but red checks definitely indicate problems that need fixing!
