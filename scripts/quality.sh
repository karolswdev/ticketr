#!/bin/bash

# Quality Check Script for Ticketr
# Runs all quality checks including linting, formatting, and tests

set -e  # Exit on error
set -u  # Exit on undefined variable

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Track overall success
CHECKS_PASSED=0
CHECKS_FAILED=0

print_header() {
    echo ""
    echo -e "${BLUE}======================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}======================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
    CHECKS_PASSED=$((CHECKS_PASSED + 1))
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
    CHECKS_FAILED=$((CHECKS_FAILED + 1))
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# Check if we're in the project root
if [ ! -f "go.mod" ]; then
    print_error "Must be run from project root directory"
    exit 1
fi

print_header "TICKETR QUALITY CHECKS"
echo "Starting quality checks for Ticketr project..."

# ============================================================================
# 1. Go Vet
# ============================================================================
print_header "1/6: Running go vet"
if go vet ./... 2>&1; then
    print_success "go vet passed"
else
    print_error "go vet found issues"
    exit 1
fi

# ============================================================================
# 2. Formatting Check
# ============================================================================
print_header "2/6: Checking code formatting"
UNFORMATTED=$(gofmt -l . 2>&1 | grep -v "^#" || true)
if [ -z "$UNFORMATTED" ]; then
    print_success "All files are properly formatted"
else
    print_error "The following files are not formatted:"
    echo "$UNFORMATTED"
    echo ""
    print_info "Run 'gofmt -w .' to format all files"
    exit 1
fi

# ============================================================================
# 3. Build Check
# ============================================================================
print_header "3/6: Building application"
if go build -o /tmp/ticketr-test ./cmd/ticketr 2>&1; then
    print_success "Application builds successfully"
    rm -f /tmp/ticketr-test
else
    print_error "Build failed"
    exit 1
fi

# ============================================================================
# 4. Tests with Coverage
# ============================================================================
print_header "4/6: Running tests with coverage"
go test ./... -cover -coverprofile=coverage.out -covermode=atomic 2>&1 | tee /tmp/test-output.txt
TEST_EXIT_CODE=${PIPESTATUS[0]}

if [ $TEST_EXIT_CODE -ne 0 ]; then
    print_error "Tests failed"
    cat /tmp/test-output.txt
    exit 1
else
    print_success "All tests passed"

    # Display coverage summary
    echo ""
    print_info "Coverage by package:"
    go tool cover -func=coverage.out | grep -v "total:" | awk '{printf "  %-50s %s\n", $1":"$2, $3}'

    echo ""
fi

# ============================================================================
# 5. Coverage Threshold Check
# ============================================================================
print_header "5/6: Checking coverage threshold"
COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
THRESHOLD=50

echo "Total coverage: ${COVERAGE}%"
echo "Threshold: ${THRESHOLD}%"

if command -v bc &> /dev/null; then
    if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
        print_error "Coverage ${COVERAGE}% is below threshold ${THRESHOLD}%"
        exit 1
    else
        print_success "Coverage ${COVERAGE}% meets threshold ${THRESHOLD}%"
    fi
else
    # Fallback for systems without bc
    COVERAGE_INT=${COVERAGE%.*}
    if [ "$COVERAGE_INT" -lt "$THRESHOLD" ]; then
        print_error "Coverage ${COVERAGE}% is below threshold ${THRESHOLD}%"
        exit 1
    else
        print_success "Coverage ${COVERAGE}% meets threshold ${THRESHOLD}%"
    fi
fi

# ============================================================================
# 6. Staticcheck (optional but recommended)
# ============================================================================
print_header "6/6: Running staticcheck"
if ! command -v staticcheck &> /dev/null; then
    print_warning "staticcheck not installed, installing..."
    if go install honnef.co/go/tools/cmd/staticcheck@latest 2>&1; then
        print_success "staticcheck installed successfully"
    else
        print_warning "Failed to install staticcheck, skipping this check"
        print_info "Install manually: go install honnef.co/go/tools/cmd/staticcheck@latest"
        CHECKS_PASSED=$((CHECKS_PASSED + 1))  # Don't fail on missing staticcheck
    fi
fi

if command -v staticcheck &> /dev/null; then
    if staticcheck ./... 2>&1; then
        print_success "staticcheck passed"
    else
        print_error "staticcheck found issues"
        exit 1
    fi
else
    print_warning "staticcheck not available, skipping"
    CHECKS_PASSED=$((CHECKS_PASSED + 1))  # Don't fail on missing staticcheck
fi

# ============================================================================
# 7. Go Module Tidiness Check
# ============================================================================
print_header "BONUS: Checking go.mod tidiness"
cp go.mod go.mod.backup
cp go.sum go.sum.backup

if go mod tidy 2>&1; then
    if diff -q go.mod go.mod.backup > /dev/null && diff -q go.sum go.sum.backup > /dev/null; then
        print_success "go.mod and go.sum are tidy"
        rm go.mod.backup go.sum.backup
    else
        print_error "go.mod or go.sum is not tidy"
        print_info "Run 'go mod tidy' to fix"
        mv go.mod.backup go.mod
        mv go.sum.backup go.sum
        exit 1
    fi
else
    print_error "go mod tidy failed"
    mv go.mod.backup go.mod
    mv go.sum.backup go.sum
    exit 1
fi

# ============================================================================
# Summary
# ============================================================================
print_header "QUALITY CHECK SUMMARY"
echo ""
echo "Checks passed: ${GREEN}${CHECKS_PASSED}${NC}"
echo "Checks failed: ${RED}${CHECKS_FAILED}${NC}"
echo ""

if [ $CHECKS_FAILED -eq 0 ]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}   ALL QUALITY CHECKS PASSED!${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    print_info "Coverage report saved to: coverage.out"
    print_info "View detailed coverage: go tool cover -html=coverage.out"
    echo ""
    exit 0
else
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}   QUALITY CHECKS FAILED${NC}"
    echo -e "${RED}========================================${NC}"
    echo ""
    print_info "Fix the issues above and run this script again"
    echo ""
    exit 1
fi
