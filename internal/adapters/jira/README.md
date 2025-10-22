# Jira Adapter

This package provides Jira integration for Ticketr using the ports & adapters architecture pattern.

## Overview

The Jira adapter implements the `ports.JiraPort` interface and provides two implementations controlled by a feature flag system:

- **V1 (jira_adapter.go)**: Custom HTTP client implementation (1,136 lines)
  - **Status:** Deprecated, removal planned for v3.2.0 or v3.3.0

- **V2 (jira_adapter_v2.go)**: Uses `andygrunwald/go-jira` library (757 lines) - **DEFAULT**
  - **Status:** Production default as of v3.1.1
  - **Benefit:** 33% code reduction, battle-tested library (9 years, 868 importers)

**Architecture Decision:** See [docs/adr/001-adopt-go-jira-library.md](../../../docs/adr/001-adopt-go-jira-library.md) for complete rationale, external validation (Gemini + Codex AI architects), and migration strategy.

## Feature Flag System

Starting from version 2.0, the Jira adapter supports multiple implementations through a feature flag system. This allows for easy switching between implementations and safe rollback if issues occur.

### Environment Variable

```bash
TICKETR_JIRA_ADAPTER_VERSION=v2  # Use library-based implementation (default)
TICKETR_JIRA_ADAPTER_VERSION=v1  # Use custom HTTP implementation
```

**Default**: If not set, the system defaults to `v2` (library-based implementation).

### Usage

#### From Code

```go
import (
    "github.com/karolswdev/ticktr/internal/adapters/jira"
    "github.com/karolswdev/ticktr/internal/core/domain"
)

// Using feature flag (recommended)
config := &domain.WorkspaceConfig{
    JiraURL:    "https://your-domain.atlassian.net",
    Username:   "your-email@example.com",
    APIToken:   "your-api-token",
    ProjectKey: "PROJ",
}

adapter, err := jira.NewJiraAdapterFromConfigWithVersion(config, nil)
if err != nil {
    // handle error
}

// The adapter will automatically use v1 or v2 based on TICKETR_JIRA_ADAPTER_VERSION
```

#### Forcing a Specific Version

```go
// Force V1 (custom HTTP)
adapterV1, err := jira.NewJiraAdapterFromConfig(config, nil)

// Force V2 (library-based)
adapterV2, err := jira.NewJiraAdapterV2FromConfig(config, nil)
```

## Implementation Comparison

### V1 (Custom HTTP Client)

**Pros:**
- Zero external dependencies (only Go stdlib)
- Full control over HTTP requests
- Battle-tested in production

**Cons:**
- More code to maintain (1,136 lines)
- Manual handling of Jira API quirks
- More complex error handling

**Use when:**
- You need zero dependencies
- You have custom HTTP requirements
- You're experiencing issues with the library

### V2 (andygrunwald/go-jira Library)

**Pros:**
- Less code to maintain (685 lines - 40% reduction)
- Library handles Jira API details
- Community support and updates
- Battle-tested library (9+ years, 868+ importers)

**Cons:**
- External dependency (12 total including transitive)
- Less control over HTTP layer
- Requires library updates for new Jira features

**Use when:**
- You want less code to maintain
- You trust the library ecosystem
- You want community-driven updates

## Rollback Procedure

If you encounter issues with V2, you can easily rollback to V1 using the feature flag system.

### Option 1: Environment Variable (Recommended)

**Instant rollback with zero downtime:**

```bash
export TICKETR_JIRA_ADAPTER_VERSION=v1
ticketr pull  # Uses V1 automatically, no restart required
```

**Rollback Time:** < 1 minute
**Scope:** Per-user or system-wide
**Risk:** None - V1 code is preserved and fully tested

### Option 2: Configuration Change

```go
// In your application initialization code
os.Setenv("TICKETR_JIRA_ADAPTER_VERSION", "v1")
adapter, _ := jira.NewJiraAdapterFromConfigWithVersion(config, nil)
```

### Rollback Decision Criteria

**Consider rolling back to V1 if:**

1. V2 error rate >2x V1 error rate (sustained over 24 hours)
2. Critical bug affecting >50% of operations
3. Data loss or corruption detected
4. Performance degradation >50% slower than V1

**Check error logs:**
```bash
# Compare V1 vs V2 error rates
grep "\[jira-v1\]" ~/.ticketr/logs/*.log | wc -l
grep "\[jira-v2\]" ~/.ticketr/logs/*.log | wc -l
```

### Verification

Check which version is active:

```go
version := jira.GetAdapterVersion()
fmt.Printf("Using adapter version: %s\n", version)
```

**Version logging in errors:**
- V1 errors contain `[jira-v1]` prefix
- V2 errors contain `[jira-v2]` prefix

This allows easy identification of which adapter version encountered issues in production logs.

## Field Mappings

Both implementations support custom field mappings:

```go
fieldMappings := map[string]interface{}{
    "Story Points": map[string]interface{}{
        "id":   "customfield_10010",
        "type": "number",
    },
    "Sprint": "customfield_10020",
}

adapter, err := jira.NewJiraAdapterV2FromConfig(config, fieldMappings)
```

If no mappings are provided, the adapter uses defaults:

- **Story Points**: customfield_10010 (number)
- **Sprint**: customfield_10020 (string)
- Standard fields: summary, description, issuetype, project, labels, etc.

## API Compatibility

Both implementations provide identical APIs through the `ports.JiraPort` interface:

```go
type JiraPort interface {
    Authenticate() error
    SearchTickets(ctx context.Context, projectKey, jql string, progressCallback JiraProgressCallback) ([]domain.Ticket, error)
    CreateTicket(ticket domain.Ticket) (string, error)
    UpdateTicket(ticket domain.Ticket) error
    CreateTask(task domain.Task, parentID string) (string, error)
    UpdateTask(task domain.Task) error
    GetProjectIssueTypes() (map[string][]string, error)
    GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error)
}
```

**No code changes required** when switching between V1 and V2.

## Testing

### Unit Tests

```bash
# Test all adapter code
go test ./internal/adapters/jira/...

# Test specific version
go test ./internal/adapters/jira/... -run TestJiraAdapterV2

# Run with coverage
go test ./internal/adapters/jira/... -cover
```

### Integration Tests

Integration tests are skipped unless Jira credentials are provided:

```bash
export JIRA_URL="https://your-domain.atlassian.net"
export JIRA_EMAIL="your-email@example.com"
export JIRA_API_KEY="your-api-token"
export JIRA_PROJECT_KEY="PROJ"

go test ./internal/adapters/jira/... -v
```

## Troubleshooting

### Issue: "unknown adapter version" error

**Symptom:**
```
Error: unknown adapter version: xyz (valid: v1, v2)
```

**Cause**: Invalid value for `TICKETR_JIRA_ADAPTER_VERSION` environment variable.

**Solution**:
```bash
# Check current value
echo $TICKETR_JIRA_ADAPTER_VERSION

# Set to valid value
export TICKETR_JIRA_ADAPTER_VERSION=v2  # or v1
```

**Valid values:**
- `v1` - Use custom HTTP client
- `v2` - Use library-based implementation (default)
- Unset - Defaults to `v2`

### Issue: "v2 adapter requires workspace configuration"

**Cause**: V2 adapter doesn't support environment variable initialization.

**Solution**: Use workspace configuration:
```go
config := &domain.WorkspaceConfig{...}
adapter, _ := jira.NewJiraAdapterV2FromConfig(config, nil)
```

### Issue: Errors contain [jira-v1] or [jira-v2] prefix

**This is expected behavior** - Not an issue!

**Purpose**: Version tags in errors help identify which adapter implementation encountered the issue.

**Example:**
```
[jira-v2] failed to search tickets: connection timeout
```

This indicates the V2 (library-based) adapter encountered the error, which is useful for:
- Debugging adapter-specific issues
- Monitoring error rates post-deployment
- Identifying when to rollback to V1

### Issue: V2 behaves differently than V1

**Expected**: Both adapters implement the same `JiraPort` interface and should behave identically.

**Debugging steps:**

1. **Check error logs for differences:**
   ```bash
   grep "\[jira-v2\]" ~/.ticketr/logs/*.log
   ```

2. **Compare against V1:**
   ```bash
   export TICKETR_JIRA_ADAPTER_VERSION=v1
   ticketr pull  # Run same operation with V1
   ```

3. **Report behavioral divergence:**
   - Document exact steps to reproduce
   - Include error messages from both V1 and V2
   - Check if issue exists in upstream library: https://github.com/andygrunwald/go-jira/issues

4. **Temporary workaround:**
   - Rollback to V1 while issue is investigated
   - See "Rollback Procedure" section above

### Issue: Authentication failures

**Check**:
1. Jira URL is correct (no trailing slash added automatically)
2. API token is valid (not password)
3. Username is email address
4. Project key matches exactly

**Test**:
```bash
curl -u "your-email:your-api-token" \
  -H "Content-Type: application/json" \
  https://your-domain.atlassian.net/rest/api/3/myself
```

**Note:** Both V1 and V2 use the same authentication mechanism (Basic Auth with API token), so authentication issues are unlikely to be version-specific.

### Issue: Custom fields not mapping

**Verify field IDs**:
1. Go to Jira → Settings → Issues → Custom Fields
2. Click on the field → Copy field ID from URL
3. Update field mappings accordingly

**Note:** Both V1 and V2 use the same field mapping logic, so this is unlikely to be version-specific.

### Issue: Rate limiting errors

Both implementations respect Jira's rate limits (10,000 requests/hour for Cloud).

**V1**: Manual retry logic recommended
**V2**: Library doesn't handle retries automatically

**Workaround**: Implement exponential backoff in calling code.

### Issue: Performance degradation with V2

**Symptom**: V2 operations noticeably slower than V1

**Diagnosis:**
```bash
# Run benchmarks
go test -bench=. -benchmem ./internal/adapters/jira/

# Look for V1 vs V2 comparison
```

**Acceptance threshold**: V2 should be within 20% of V1 performance

**If degradation >20%:**
1. Check network conditions (library may make different HTTP requests)
2. Review Jira API response times (library may use different endpoints)
3. Profile V2 adapter to identify bottlenecks
4. Consider rollback to V1 while performance is investigated

**Report to maintainers**: Include benchmark results and profiling data

## Migration Guide

### From V1 to V2

**No code changes required** if using the factory pattern:

```bash
# Just set the environment variable
export TICKETR_JIRA_ADAPTER_VERSION=v2
```

If directly instantiating V1:

```diff
- adapter, err := jira.NewJiraAdapterFromConfig(config, fieldMappings)
+ adapter, err := jira.NewJiraAdapterV2FromConfig(config, fieldMappings)
```

### From V2 back to V1

```bash
# Immediate rollback via environment variable
export TICKETR_JIRA_ADAPTER_VERSION=v1
```

## Performance

Both implementations use pagination for large result sets:

- **Page size**: 50 tickets per request
- **Progress callbacks**: Supported in both versions
- **Context cancellation**: Supported in both versions

**V2 advantages**:
- Library handles response parsing efficiently
- Built-in connection pooling

**V1 advantages**:
- Custom timeout configuration
- More granular control over HTTP client

## Dependencies

### V1 Dependencies
- Go standard library only

### V2 Dependencies
```
github.com/andygrunwald/go-jira v1.17.0
├── github.com/fatih/structs v1.1.0
├── github.com/golang-jwt/jwt/v4 v4.5.2
├── github.com/google/go-cmp v0.7.0
├── github.com/google/go-querystring v1.1.0
└── github.com/trivago/tgo v1.0.7
```

**Total**: 12 dependencies (including transitive)

## Security

- Both implementations use Basic Auth with API tokens
- No passwords stored in code or logs
- Credentials passed via workspace configuration
- HTTPS enforced for all Jira API calls

**Vulnerability status** (as of 2025-10-21):
- V1: No dependencies, no CVEs
- V2: No library-specific CVEs (see EXTERNAL-VALIDATION-REPORT.md)

## Contributing

When making changes to the adapter:

1. Update both V1 and V2 if changing the interface
2. Add tests for new functionality
3. Update this README with configuration changes
4. Test both versions before committing

## License

Part of Ticketr project. See LICENSE in repository root.
