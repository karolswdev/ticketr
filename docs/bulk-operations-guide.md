# Bulk Operations Guide

**Last Updated:** October 19, 2025 (Week 18)
**Version:** Ticketr v3.0
**Status:** Production-ready (update/move), Planned (delete in v3.1.0)

---

## Table of Contents

1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [Command Reference](#command-reference)
4. [Examples](#examples)
5. [Progress Feedback](#progress-feedback)
6. [Safety Features](#safety-features)
7. [Troubleshooting](#troubleshooting)
8. [Limitations](#limitations)
9. [Roadmap](#roadmap)

---

## Introduction

Bulk operations allow you to perform actions on multiple Jira tickets at once, saving time and reducing repetitive manual work. Instead of updating tickets one at a time, you can apply changes to dozens or hundreds of tickets with a single command.

### What Are Bulk Operations?

Ticketr v3.0 supports three types of bulk operations:

- **Update**: Modify field values on multiple tickets (e.g., change status, assign tickets, update priority)
- **Move**: Change the parent relationship for multiple tickets (e.g., move sub-tasks to a different epic)
- **Delete**: Remove multiple tickets from Jira (planned for v3.1.0, not yet supported)

### When to Use Bulk Operations

**Use bulk operations when you need to:**
- Update status for all tickets in a sprint
- Reassign multiple tickets to a new owner
- Move sub-tasks to a different parent ticket
- Change priority or labels for related tickets
- Apply the same field changes to many tickets

**When NOT to use bulk operations:**
- Single ticket updates (use `ticketr push` instead)
- Complex per-ticket logic (different changes for each ticket)
- More than 100 tickets at once (current limit)

---

## Getting Started

### Prerequisites

Before using bulk operations, ensure:

1. **Active workspace**: You must have a workspace configured
   ```bash
   ticketr workspace current
   ```

2. **Valid credentials**: Your Jira credentials must be stored in the workspace
   ```bash
   ticketr workspace list
   ```

3. **Ticket IDs ready**: Know the exact Jira ticket IDs (e.g., PROJ-123, PROJ-124)

### Basic Concepts

#### Ticket IDs

Ticket IDs must follow Jira format: `PROJECT-123`
- Project key (uppercase letters)
- Hyphen
- Ticket number (digits)

**Valid examples**: `PROJ-1`, `BACKEND-42`, `EPIC-100`
**Invalid examples**: `proj-1`, `PROJ_1`, `123`, `PROJ-`

#### Field Changes

Field changes specify what to modify on each ticket:
- **Field name**: The Jira field to change (e.g., `status`, `assignee`, `priority`)
- **Field value**: The new value to set

**Format**: `field=value`

**Examples**:
- `status=Done`
- `assignee=john@example.com`
- `priority=High`

---

## Command Reference

### `ticketr bulk update`

Update multiple tickets with the same field changes.

#### Syntax

```bash
ticketr bulk update --ids <ticket-ids> --set <field=value> [--set <field=value>...]
```

#### Parameters

- `--ids`: Comma-separated list of ticket IDs (required)
  - Format: `PROJ-1,PROJ-2,PROJ-3`
  - Spaces allowed: `PROJ-1, PROJ-2, PROJ-3`
  - Minimum: 1 ticket
  - Maximum: 100 tickets

- `--set`: Field changes to apply (required, can be used multiple times)
  - Format: `field=value`
  - Use quotes for values with spaces: `--set priority="High Priority"`

#### Examples

**Update status for multiple tickets:**
```bash
ticketr bulk update --ids PROJ-1,PROJ-2,PROJ-3 --set status=Done
```

**Update multiple fields:**
```bash
ticketr bulk update \
  --ids PROJ-1,PROJ-2 \
  --set status="In Progress" \
  --set assignee=developer@company.com
```

**Update with spaces in values:**
```bash
ticketr bulk update \
  --ids PROJ-1,PROJ-2 \
  --set priority="High Priority" \
  --set labels="urgent, critical"
```

#### Expected Output

```
Updating 3 ticket(s) with changes:
  status = Done

[1/3] PROJ-1: ✓ success
[2/3] PROJ-2: ✓ success
[3/3] PROJ-3: ✓ success

=== Summary ===
✓ 3 ticket(s) updated successfully
```

---

### `ticketr bulk move`

Move multiple tickets to a new parent.

#### Syntax

```bash
ticketr bulk move --ids <ticket-ids> --parent <parent-ticket-id>
```

#### Parameters

- `--ids`: Comma-separated list of ticket IDs to move (required)
- `--parent`: New parent ticket ID (required)

#### Examples

**Move sub-tasks to a new parent:**
```bash
ticketr bulk move --ids PROJ-1,PROJ-2 --parent PROJ-100
```

**Move tickets to a different epic:**
```bash
ticketr bulk move \
  --ids TASK-1,TASK-2,TASK-3 \
  --parent EPIC-42
```

#### Expected Output

```
Moving 3 ticket(s) to parent EPIC-42

[1/3] TASK-1: ✓ success
[2/3] TASK-2: ✓ success
[3/3] TASK-3: ✓ success

=== Summary ===
✓ 3 ticket(s) moved successfully
```

---

### `ticketr bulk delete`

Delete multiple tickets (not supported in v3.0, planned for v3.1.0).

#### Syntax

```bash
ticketr bulk delete --ids <ticket-ids> --confirm
```

#### Parameters

- `--ids`: Comma-separated list of ticket IDs to delete (required)
- `--confirm`: Acknowledge destructive operation (required for safety)

#### Status

Bulk delete is **not supported** in Ticketr v3.0. This feature is planned for v3.1.0.

**Current behavior:**
```bash
ticketr bulk delete --ids PROJ-1,PROJ-2 --confirm
```

**Output:**
```
⚠️  Delete operations are not yet supported in Ticketr v3.0

Ticketr currently does not support bulk delete operations via the Jira API.
This feature is planned for v3.1.0.

For now, please delete tickets individually through Jira:

  • PROJ-1
  • PROJ-2

Alternative options:
  1. Use the Jira web interface for bulk deletion
  2. Wait for v3.1.0 which will include this feature
```

---

## Examples

### Common Use Cases

#### 1. Update Status for Sprint Completion

Mark all tickets in a sprint as "Done":

```bash
ticketr bulk update \
  --ids PROJ-10,PROJ-11,PROJ-12,PROJ-13,PROJ-14 \
  --set status=Done
```

#### 2. Reassign Tickets to New Owner

Transfer ownership when a team member leaves:

```bash
ticketr bulk update \
  --ids PROJ-20,PROJ-21,PROJ-22 \
  --set assignee=newowner@company.com
```

#### 3. Update Priority for Critical Issues

Escalate multiple tickets:

```bash
ticketr bulk update \
  --ids BUG-1,BUG-2,BUG-3 \
  --set priority=Critical
```

#### 4. Move Sub-Tasks to New Epic

Reorganize work under a different epic:

```bash
ticketr bulk move \
  --ids STORY-1,STORY-2,STORY-3,STORY-4 \
  --parent EPIC-5
```

#### 5. Update Multiple Fields at Once

Apply several changes simultaneously:

```bash
ticketr bulk update \
  --ids PROJ-30,PROJ-31,PROJ-32 \
  --set status="In Progress" \
  --set assignee=team@company.com \
  --set priority=High \
  --set labels="sprint-5, backend"
```

### Handling Spaces in Values

When field values contain spaces, use quotes:

```bash
# Single quotes (Unix/Linux/macOS)
ticketr bulk update --ids PROJ-1,PROJ-2 --set priority='High Priority'

# Double quotes (works on all platforms)
ticketr bulk update --ids PROJ-1,PROJ-2 --set priority="High Priority"

# Escape spaces (Unix/Linux/macOS only)
ticketr bulk update --ids PROJ-1,PROJ-2 --set priority=High\ Priority
```

**Recommended**: Always use double quotes for cross-platform compatibility.

---

## Progress Feedback

Bulk operations provide real-time progress indicators as each ticket is processed.

### Progress Format

```
[current/total] TICKET-ID: ✓ success
[current/total] TICKET-ID: ✗ failed (error message)
```

### Example Output

```
Updating 5 ticket(s) with changes:
  status = Done

[1/5] PROJ-1: ✓ success
[2/5] PROJ-2: ✓ success
[3/5] PROJ-3: ✗ failed (ticket not found: PROJ-3)
[4/5] PROJ-4: ✓ success
[5/5] PROJ-5: ✓ success

=== Summary ===
✓ 4 ticket(s) updated successfully
✗ 1 ticket(s) failed

Errors:
  PROJ-3: ticket not found: PROJ-3
```

### Understanding the Summary

- **Success count**: Number of tickets successfully updated
- **Failure count**: Number of tickets that failed
- **Errors section**: Detailed error messages for each failed ticket

### Partial Failures

If some tickets succeed and some fail, the operation continues processing all tickets and returns a **partial success**:

- Successful tickets remain updated
- Failed tickets remain unchanged
- For update/move operations, Ticketr attempts **best-effort rollback** (see Safety Features)

---

## Safety Features

Ticketr v3.0 includes several safety mechanisms to protect against accidental damage.

### 1. JQL Injection Prevention

Ticket IDs are validated using a strict regex pattern to prevent JQL injection attacks:

**Pattern**: `^[A-Z]+-\d+$`

**Blocked examples**:
- `PROJ-1" OR 1=1` (SQL-style injection)
- `PROJ-1; DROP TABLE` (command injection)
- `../../../etc/passwd` (path traversal)

**Result**: Invalid ticket IDs are rejected before any Jira API calls.

### 2. Confirmation Prompts (Delete Only)

Delete operations (when supported in v3.1.0) will require:

1. `--confirm` flag to acknowledge destructive operation
2. Interactive confirmation prompt
3. Display of affected tickets before deletion

### 3. Transaction Rollback Behavior

For update and move operations, Ticketr attempts **best-effort rollback** on partial failures:

**How it works:**
1. Before updating each ticket, Ticketr fetches and stores the original state
2. If a later ticket fails, Ticketr attempts to restore all successfully updated tickets
3. Rollback is **best-effort** (cannot guarantee 100% success)

**Why "best-effort"?**
- Network failures during rollback may prevent restoration
- Concurrent Jira edits may conflict with rollback
- Jira API rate limits may block rollback requests

**Recommendation**: Always verify ticket states after partial failures.

### 4. Maximum Ticket Limit

Bulk operations are limited to **100 tickets** per command to:
- Prevent accidental mass changes
- Ensure reasonable execution time
- Reduce risk of Jira API rate limiting

**Workaround for >100 tickets:**
```bash
# Split into multiple batches
ticketr bulk update --ids PROJ-1,...,PROJ-100 --set status=Done
ticketr bulk update --ids PROJ-101,...,PROJ-200 --set status=Done
```

---

## Troubleshooting

### Common Errors

#### 1. "No workspace selected"

**Symptom:**
```
Error: no workspace selected: use 'ticketr workspace create' or 'ticketr workspace switch'
```

**Solution:**
```bash
# List available workspaces
ticketr workspace list

# Switch to a workspace
ticketr workspace switch backend

# Or create a new workspace
ticketr workspace create my-project \
  --url https://company.atlassian.net \
  --project PROJ \
  --username your.email@company.com \
  --token your-api-token
```

---

#### 2. "Invalid Jira ID format"

**Symptom:**
```
Error: invalid bulk operation: ticket_ids[0]: invalid Jira ID format: proj-1 (expected format: PROJECT-123)
```

**Cause:** Ticket IDs must have uppercase project keys.

**Solution:**
```bash
# Incorrect (lowercase)
ticketr bulk update --ids proj-1,proj-2 --set status=Done

# Correct (uppercase)
ticketr bulk update --ids PROJ-1,PROJ-2 --set status=Done
```

---

#### 3. "Authentication failed"

**Symptom:**
```
Error: failed to authenticate with Jira: 401 Unauthorized
```

**Cause:** Workspace credentials are invalid or expired.

**Solution:**
```bash
# Regenerate API token
# Visit: https://id.atlassian.com/manage-profile/security/api-tokens

# Delete old workspace
ticketr workspace delete my-workspace --force

# Recreate with new token
ticketr workspace create my-workspace \
  --url https://company.atlassian.net \
  --project PROJ \
  --username your.email@company.com \
  --token new-api-token
```

---

#### 4. "Ticket not found"

**Symptom:**
```
[1/3] PROJ-999: ✗ failed (ticket not found: PROJ-999)
```

**Cause:** Ticket does not exist in Jira or you don't have access.

**Solution:**
- Verify ticket ID is correct
- Check you have access to the ticket in Jira
- Ensure ticket is in the current workspace's project

---

#### 5. "Partial failure: rollback attempted"

**Symptom:**
```
Error: bulk update completed with errors: partial failure: 2 of 5 tickets failed (rollback attempted)
```

**Cause:** Some tickets failed to update, and Ticketr attempted to restore original state.

**Solution:**
1. Check the error details for each failed ticket
2. Manually verify ticket states in Jira
3. Re-run the command with only the failed ticket IDs
4. If rollback failed, manually revert changes

---

#### 6. "Field name cannot be empty"

**Symptom:**
```
Error: invalid --set format: =Done (field name cannot be empty)
```

**Cause:** Missing field name before `=`.

**Solution:**
```bash
# Incorrect
ticketr bulk update --ids PROJ-1 --set =Done

# Correct
ticketr bulk update --ids PROJ-1 --set status=Done
```

---

#### 7. "Changes are required for 'update' action"

**Symptom:**
```
Error: invalid bulk operation: changes are required for 'update' action
```

**Cause:** No `--set` flags provided.

**Solution:**
```bash
# Incorrect (no --set flags)
ticketr bulk update --ids PROJ-1,PROJ-2

# Correct (at least one --set flag)
ticketr bulk update --ids PROJ-1,PROJ-2 --set status=Done
```

---

### Debug Mode

Enable verbose logging for troubleshooting:

```bash
# Set log level to debug
export TICKETR_LOG_LEVEL=debug

# Run bulk operation
ticketr bulk update --ids PROJ-1,PROJ-2 --set status=Done --verbose
```

**Log location:** `~/.cache/ticketr/logs/` (Linux), `~/Library/Caches/ticketr/logs/` (macOS), `%TEMP%\ticketr\logs\` (Windows)

---

## Limitations

Ticketr v3.0 bulk operations have the following limitations:

### 1. Maximum 100 Tickets Per Operation

You cannot process more than 100 tickets in a single bulk command.

**Workaround**: Split into multiple batches.

### 2. Delete Operation Not Available

Bulk delete is not supported in v3.0. Planned for v3.1.0.

**Alternative**: Use Jira web interface for bulk deletion.

### 3. Best-Effort Rollback

Rollback on partial failures cannot be guaranteed to succeed. Always verify ticket states after errors.

### 4. Sequential Processing

Tickets are processed one at a time (not in parallel). For 100 tickets, expect 1-2 minutes processing time.

**Future enhancement**: Parallel processing planned for v3.2.0.

### 5. Field Type Validation

Ticketr does not validate field types before sending to Jira. Invalid field values will fail at the Jira API level.

**Example failure:**
```bash
# Story Points expects a number, not text
ticketr bulk update --ids PROJ-1 --set storyPoints="five"
# Result: Jira API error
```

**Recommendation**: Verify field types in Jira before running bulk operations.

### 6. Workspace Requirement

Bulk operations require an active workspace. Environment variables are not supported.

**Why?** Bulk operations are a workspace-scoped feature introduced in v3.0.

---

## Roadmap

### v3.1.0 (Planned Q1 2026)

- **Bulk delete support**: Full implementation with confirmation prompts
- **Dry-run mode**: Preview changes before applying (`--dry-run` flag)
- **Enhanced validation**: Field type checking before API calls

### v3.2.0 (Planned Q2 2026)

- **Parallel processing**: Speed up bulk operations with concurrent requests
- **Batch API support**: Reduce HTTP overhead with Jira batch endpoints
- **Progress bar UI**: Enhanced progress indicators with ETA

### v3.3.0 (Future)

- **JQL-based bulk operations**: Operate on tickets matching JQL queries
  ```bash
  ticketr bulk update --jql "project=PROJ AND status='In Progress'" --set assignee=new@company.com
  ```
- **Undo/redo support**: Reverse bulk operations
- **Audit log**: Track all bulk operations with timestamps

---

## Best Practices

### 1. Start Small

Test bulk operations on a small set of tickets before scaling up:

```bash
# Test with 2-3 tickets first
ticketr bulk update --ids PROJ-1,PROJ-2 --set status=Done

# If successful, scale up
ticketr bulk update --ids PROJ-1,PROJ-2,PROJ-3,...,PROJ-50 --set status=Done
```

### 2. Verify in Jira

Always verify ticket states in Jira after bulk operations, especially after partial failures.

### 3. Use Descriptive Field Names

Use Jira's exact field names for clarity:

```bash
# Good (clear)
ticketr bulk update --ids PROJ-1,PROJ-2 --set assignee=user@company.com

# Avoid (ambiguous)
ticketr bulk update --ids PROJ-1,PROJ-2 --set owner=user@company.com
```

### 4. Keep Tickets Under 100

Split large operations into batches to stay within the 100-ticket limit.

### 5. Monitor Progress Output

Watch the real-time progress indicators to catch errors early:

```
[1/50] PROJ-1: ✓ success
[2/50] PROJ-2: ✗ failed (authentication error)
```

If you see authentication errors early, stop the operation (Ctrl+C) and fix credentials.

---

## Advanced Topics

### Scripting Bulk Operations

Bulk operations can be automated in shell scripts:

```bash
#!/bin/bash
# update-sprint-tickets.sh

SPRINT_TICKETS="PROJ-1,PROJ-2,PROJ-3,PROJ-4,PROJ-5"
NEW_STATUS="Done"

echo "Updating sprint tickets to $NEW_STATUS"
ticketr bulk update --ids "$SPRINT_TICKETS" --set "status=$NEW_STATUS"

if [ $? -eq 0 ]; then
  echo "✓ All tickets updated successfully"
else
  echo "✗ Some tickets failed to update"
  exit 1
fi
```

### CI/CD Integration

Use bulk operations in CI/CD pipelines:

```yaml
# .github/workflows/close-sprint.yml
- name: Close sprint tickets
  run: |
    ticketr workspace switch ci-workspace
    ticketr bulk update --ids ${{ env.SPRINT_TICKETS }} --set status=Done
```

---

## Summary

Ticketr v3.0 bulk operations provide:

- **Update**: Modify fields on multiple tickets
- **Move**: Change parent relationships for multiple tickets
- **Delete**: Planned for v3.1.0 (not yet supported)

**Key benefits:**
- Real-time progress feedback
- JQL injection prevention
- Best-effort rollback on partial failures
- Maximum 100 tickets per operation

**Limitations:**
- Sequential processing (not parallel)
- Best-effort rollback (not guaranteed)
- Delete not supported in v3.0

For more information:
- [README.md](../README.md) - Quick start guide
- [Workspace Management Guide](workspace-management-guide.md) - Workspace setup
- [Bulk Operations API](bulk-operations-api.md) - Developer documentation

---

**Document Version:** 1.0
**Status:** Production-ready
**Feedback:** [GitHub Issues](https://github.com/karolswdev/ticketr/issues)
