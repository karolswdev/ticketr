# Ticketr Workflow Guide

## Complete End-to-End Workflow

This guide demonstrates a complete ticket lifecycle: create tickets in Markdown → push to JIRA → edit in JIRA → pull changes → review logs.

### Prerequisites

Set up your environment:
```bash
export JIRA_URL="https://yourcompany.atlassian.net"
export JIRA_EMAIL="your.email@company.com"
export JIRA_API_KEY="your_api_token"
export JIRA_PROJECT_KEY="PROJ"
```

### Step 1: Create Tickets in Markdown

Create `my-tickets.md`:
```markdown
# TICKET: Implement user authentication

## Description
Add JWT-based authentication to the API.

## Fields
Type: Story
Sprint: Sprint 23
Story Points: 13
Priority: High

## Acceptance Criteria
- Users can log in with email/password
- JWT tokens expire after 24 hours
- Refresh tokens supported

## Tasks
- Setup authentication middleware
  ## Description
  Create middleware to validate JWT tokens on protected routes.

  ## Fields
  Story Points: 5

  ## Acceptance Criteria
  - Middleware validates JWT signature
  - Returns 401 for invalid tokens

- Implement login endpoint
  ## Description
  Create POST /api/login endpoint with email/password validation.

  ## Fields
  Story Points: 8
  Priority: Critical

  ## Acceptance Criteria
  - Validates email format
  - Returns JWT token on success
```

### Step 2: Push to JIRA (First Time)

```bash
$ ticketr push my-tickets.md

=== Summary ===
Tickets created: 1
Tasks created: 2

Processing complete!
```

After first push, `my-tickets.md` is updated with JIRA IDs:
```markdown
# TICKET: [PROJ-123] Implement user authentication
...
## Tasks
- [PROJ-124] Setup authentication middleware
- [PROJ-125] Implement login endpoint
```

**What happened:**
- Parent ticket created in JIRA (PROJ-123)
- Two subtasks created (PROJ-124, PROJ-125)
- Field inheritance applied: Sprint 23 inherited by both tasks
- Priority override applied: PROJ-125 has "Critical" instead of "High"
- State file `.ticketr.state` created with ticket hashes
- Log file `.ticketr/logs/2025-10-16_14-30-00.log` created

### Step 3: Push Again (State-Based Skipping)

```bash
$ ticketr push my-tickets.md

=== Summary ===
Tickets created: 0
Tickets updated: 0
Tasks created: 0
Tasks updated: 0

Processing complete!
```

**What happened:**
- State manager detected no changes
- All tickets skipped (zero API calls to JIRA)
- This is Milestone 9's state-aware push in action!

### Step 4: Make Local Changes

Edit `my-tickets.md` - add a new task:
```markdown
## Tasks
- [PROJ-124] Setup authentication middleware
- [PROJ-125] Implement login endpoint
- Add password reset endpoint
  ## Description
  Create POST /api/reset-password endpoint.

  ## Fields
  Story Points: 3
```

### Step 5: Push Updated Tickets

```bash
$ ticketr push my-tickets.md

=== Summary ===
Tickets updated: 1
Tasks created: 1

Processing complete!
```

**What happened:**
- Parent ticket PROJ-123 updated (new task added)
- New subtask PROJ-126 created
- Existing tasks (PROJ-124, PROJ-125) skipped (unchanged)
- State file updated with new hashes

### Step 6: Edit in JIRA Web UI

Go to JIRA and:
1. Change PROJ-123 Priority from "High" to "Medium"
2. Add comment to PROJ-124
3. Move PROJ-125 to "In Progress"

### Step 7: Pull Changes from JIRA

```bash
$ ticketr pull --project PROJ --output my-tickets.md

Successfully updated my-tickets.md
  - 1 ticket(s) updated with remote changes
```

**What happened:**
- Pulled latest ticket data from JIRA
- Priority field updated from "High" to "Medium"
- Subtasks pulled with current JIRA values
- Conflict detection ran (none detected, local file matches)

### Step 8: Handle Conflicts

Make a local change:
```bash
# Edit my-tickets.md, change Priority to "Low"
```

Then pull:
```bash
$ ticketr pull --project PROJ --output my-tickets.md

⚠️  Conflict detected! The following tickets have both local and remote changes:
  - PROJ-123

To force overwrite local changes with remote changes, use --force flag
```

**Resolution options:**

Option A: Force overwrite local with remote:
```bash
$ ticketr pull --project PROJ --output my-tickets.md --force
Successfully updated my-tickets.md
  - 1 ticket(s) updated with remote changes
```

Option B: Keep local changes, push to JIRA:
```bash
$ ticketr push my-tickets.md
```

### Step 9: Review Execution Logs

```bash
$ cat .ticketr/logs/2025-10-16_14-30-00.log

[2025-10-16 14:30:00] ========================================
[2025-10-16 14:30:00] PUSH COMMAND
[2025-10-16 14:30:00] ========================================
[2025-10-16 14:30:00] Input file: my-tickets.md
[2025-10-16 14:30:00] Force partial upload: false

[2025-10-16 14:30:02] ========================================
[2025-10-16 14:30:02] EXECUTION SUMMARY
[2025-10-16 14:30:02] ========================================
[2025-10-16 14:30:02] Tickets created: 1
[2025-10-16 14:30:02] Tickets updated: 0
[2025-10-16 14:30:02] Tasks created: 2
[2025-10-16 14:30:02] Tasks updated: 0
```

**Log features:**
- Human-readable format with timestamps
- Sections for commands, operations, summary, errors
- Sensitive data redacted (API keys, emails)
- Automatic rotation (keeps last 10 files)
- Location: `.ticketr/logs/` (configurable via `TICKETR_LOG_DIR`)

### Step 10: Advanced Scenarios

#### Scenario A: Partial Upload with Errors

If some tickets have validation errors:
```bash
$ ticketr push my-tickets.md --force-partial-upload

Warning: Validation warnings (processing will continue with --force-partial-upload):
  - Line 45: Missing required field "Type"

=== Summary ===
Tickets created: 2
Tickets updated: 1

=== Errors (1) ===
  - Failed to create ticket 'Invalid Ticket': missing required field

Processing complete!
```

#### Scenario B: Updating Older Files

If you encounter an older Markdown file that still uses `# STORY:` headings, rename them to `# TICKET:` before running Ticketr. The parser will otherwise stop with an error so no unintended changes occur.

## Key Concepts

### State Management
- `.ticketr.state` tracks ticket hashes
- Unchanged tickets skipped automatically
- Improves performance, reduces API calls
- State file updated after each successful push

### Field Inheritance
- Tasks inherit parent ticket's CustomFields
- Task-specific fields override parent values
- Applied during push operations
- Example: Parent has `Sprint: Sprint 1` → all tasks inherit unless overridden

### Conflict Detection
- Compares local file hash vs remote JIRA state
- Detects simultaneous edits
- Requires explicit --force to overwrite
- Prevents accidental data loss

### Logging
- All operations logged to `.ticketr/logs/`
- Human-readable format for auditing
- Sensitive data automatically redacted
- Helps troubleshoot issues

## Common Workflows

### Daily Development Flow
1. Pull latest from JIRA: `ticketr pull --project PROJ --output tickets.md`
2. Edit tickets locally in Markdown
3. Push changes: `ticketr push tickets.md`
4. Repeat

### Sprint Planning Flow
1. Create sprint tickets in Markdown
2. Bulk push: `ticketr push sprint-23-tickets.md`
3. Review in JIRA web UI
4. Make adjustments in JIRA
5. Pull back: `ticketr pull --project PROJ --jql "sprint='Sprint 23'" --output sprint-23.md`

### Collaboration Flow
1. Alice edits in Markdown, pushes to JIRA
2. Bob edits in JIRA web UI
3. Alice pulls, detects conflict
4. Alice reviews changes, decides to force-pull or re-push
5. Team stays in sync

## Troubleshooting

### Issue: Tickets not skipping despite no changes
**Cause:** State file deleted or corrupted
**Solution:** State will rebuild on next push. First push after state loss will update all tickets.

### Issue: Field inheritance not working
**Cause:** Task fields in Markdown override parent fields
**Solution:** Remove field from task Markdown to inherit from parent.

### Issue: Pull conflicts on every pull
**Cause:** JIRA updates timestamps even without content changes
**Solution:** Use --force if you want to always take remote version.

## Next Steps

- Read [State Management Guide](state-management.md) for details on hashing
- Read [Integration Testing Guide](integration-testing-guide.md) for testing
