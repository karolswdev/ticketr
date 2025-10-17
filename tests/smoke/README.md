# Smoke Test Suite for Ticketr

This directory contains smoke tests for the Ticketr CLI application. Smoke tests verify that critical user-facing functionality works correctly in realistic scenarios.

## Overview

The smoke test suite validates:
- Core CLI commands work without crashes
- File operations are safe and reliable
- State management functions correctly
- Error handling is graceful
- Basic workflows can be completed end-to-end

## Prerequisites

Before running smoke tests, ensure you have:

1. **Go 1.21 or higher** installed
2. **Ticketr binary** built:
   ```bash
   go build -o ticketr ./cmd/ticketr
   ```
3. **Bash shell** (Linux/macOS/WSL)
4. **Standard Unix tools**: grep, find, wc, etc.

## Running Smoke Tests

### Quick Start

From the repository root:

```bash
# Build the application
go build -o ticketr ./cmd/ticketr

# Run smoke tests
./tests/smoke/smoke_test.sh
```

### From the smoke test directory

```bash
cd tests/smoke
./smoke_test.sh
```

### Expected Output

The script will output colored test results:
- Green: Tests that passed
- Red: Tests that failed
- Yellow: Test headers and summary

Example:
```
======================================
TEST 1: Push Dry-Run Validation
======================================
PASS: Push dry-run executed without errors
PASS: Dry-run did not modify source file

======================================
SMOKE TEST SUMMARY
======================================
Tests Run:    12
Tests Passed: 12
Tests Failed: 0

ALL SMOKE TESTS PASSED
```

## Test Scenarios

### Test 1: Push Dry-Run Validation
**Purpose**: Ensure `ticketr push --dry-run` validates tickets without making actual changes.

**What it tests**:
- Dry-run flag functionality
- File integrity during dry-run
- Error handling when Jira config is missing

**Expected behavior**:
- Command executes without crashing
- Source files are not modified
- Graceful error handling for missing config

---

### Test 2: Pull with Missing File (First-Run)
**Purpose**: Verify that `ticketr pull` handles first-run scenarios gracefully.

**What it tests**:
- Behavior when no local files exist
- Error messages for missing configuration
- Graceful failure modes

**Expected behavior**:
- Clear error message about missing config
- No crashes or panics
- Exit code indicates failure

---

### Test 3: State File Creation and Persistence
**Purpose**: Confirm that Ticketr creates and maintains state files correctly.

**What it tests**:
- `.ticketr` directory creation
- State file JSON structure
- Persistence of ticket metadata

**Expected behavior**:
- `.ticketr/state.json` is created
- Valid JSON format
- Ticket tracking data is stored

---

### Test 4: Log File Creation
**Purpose**: Verify that execution logs are created in the correct location.

**What it tests**:
- `.ticketr/logs/` directory creation
- Log file generation
- Log file accessibility

**Expected behavior**:
- Log directory exists
- Log files are created on demand
- Logs are readable

---

### Test 5: Help Command and Basic CLI
**Purpose**: Ensure basic CLI functionality works.

**What it tests**:
- `--help` flag
- Basic command invocation
- User-facing documentation

**Expected behavior**:
- Help text is displayed
- No crashes on basic invocation
- Useful information provided

---

### Test 7: Concurrent File Operations Safety
**Purpose**: Verify that Ticketr handles concurrent operations safely.

**What it tests**:
- Multiple simultaneous file reads
- Data integrity during concurrent access
- No file corruption

**Expected behavior**:
- All files remain intact
- No race conditions
- Concurrent operations complete successfully

## Troubleshooting

### Binary Not Found
**Error**: `ERROR: ticketr binary not found`

**Solution**: Build the binary first:
```bash
go build -o ticketr ./cmd/ticketr
```

### Permission Denied
**Error**: `Permission denied: ./smoke_test.sh`

**Solution**: Make the script executable:
```bash
chmod +x tests/smoke/smoke_test.sh
```

### Tests Fail Due to Missing Dependencies
**Error**: Various failures related to missing tools

**Solution**: Ensure you have standard Unix tools:
```bash
# On Ubuntu/Debian
sudo apt-get install grep findutils coreutils

# On macOS (usually pre-installed)
# No action needed
```

### Temporary Directory Cleanup
The smoke tests create temporary directories in `/tmp/smoke-test-*`. These are automatically cleaned up on exit, but if tests are interrupted:

```bash
rm -rf /tmp/smoke-test-*
```

## Integration with CI/CD

The smoke test suite is integrated into the GitHub Actions CI workflow:

```yaml
- name: Run smoke tests
  run: ./tests/smoke/smoke_test.sh
```

In CI, the tests run after a successful build to ensure the binary works correctly before deployment.

## Writing New Smoke Tests

To add a new smoke test:

1. **Choose a test number** (next available)
2. **Add a test section** using the template:
   ```bash
   log_test "8" "Your Test Name"

   TEST_DIR="$TEST_ROOT/test8"
   mkdir -p "$TEST_DIR"
   cd "$TEST_DIR"

   # Your test logic here

   if [ condition ]; then
       log_pass "Success message"
   else
       log_fail "Failure message"
   fi
   ```
3. **Update this README** with test documentation
4. **Test locally** before committing

### Best Practices

- **Isolate tests**: Each test should use its own temporary directory
- **Clean up**: Use the trap mechanism to ensure cleanup
- **Descriptive messages**: Make pass/fail messages clear
- **No external dependencies**: Tests should work without Jira credentials
- **Idempotent**: Tests should be repeatable without side effects

## Maintenance

These smoke tests should be updated when:
- New major features are added to the CLI
- Critical bugs are fixed (add regression tests)
- Command-line interface changes
- File format or structure changes

## Exit Codes

- `0`: All tests passed
- `1`: One or more tests failed

## Environment Variables

Currently, the smoke tests do not use environment variables. Future enhancements might include:
- `TICKETR_BIN`: Override path to ticketr binary
- `SMOKE_TEST_VERBOSE`: Enable verbose output
- `SMOKE_TEST_KEEP_TEMP`: Preserve temporary directories for debugging

## Related Documentation

- [ROADMAP.md](../../ROADMAP.md) - Milestone 11: Quality Gates & Automation
- [Requirements](../../docs/development/REQUIREMENTS.md) - Functional requirements
- [CI Workflow](../../.github/workflows/ci.yml) - GitHub Actions integration

## Support

For issues with smoke tests:
1. Check the troubleshooting section above
2. Review the test output carefully
3. Run individual test sections manually for debugging
4. Check GitHub Actions logs if CI fails
