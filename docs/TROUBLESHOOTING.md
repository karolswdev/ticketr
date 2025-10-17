# Troubleshooting Guide

Comprehensive solutions to common issues with Ticketr.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Authentication & Connection](#authentication--connection)
- [Field & Schema Issues](#field--schema-issues)
- [State Management](#state-management)
- [Push/Pull Issues](#pushpull-issues)
- [Validation Errors](#validation-errors)
- [Performance Issues](#performance-issues)
- [Logging Issues](#logging-issues)
- [Getting Help](#getting-help)

---

## Installation Issues

### Command Not Found

**Problem:** `ticketr: command not found`

**Solutions:**

1. **Verify installation:**
   ```bash
   which ticketr
   go env GOPATH
   ```

2. **Add Go bin to PATH:**
   ```bash
   export PATH=$PATH:$(go env GOPATH)/bin
   # Add to ~/.bashrc or ~/.zshrc for persistence
   ```

3. **Reinstall:**
   ```bash
   go install github.com/karolswdev/ticketr/cmd/ticketr@latest
   ```

### Permission Denied

**Problem:** `permission denied: ./ticketr`

**Solution:**
```bash
chmod +x ./ticketr
# Or for system-wide:
sudo chmod +x /usr/local/bin/ticketr
```

### Import Errors (Building from Source)

**Problem:** Package import errors when building

**Solution:**
```bash
go mod tidy
go mod download
go build -o ticketr cmd/ticketr/main.go
```

---

## Authentication & Connection

### Failed to Authenticate

**Problem:** `failed to authenticate with JIRA`

**Diagnostic Steps:**

1. **Verify credentials:**
   ```bash
   echo $JIRA_URL        # Should include https://
   echo $JIRA_EMAIL      # Your JIRA email
   # JIRA_API_KEY should be set but NOT echoed for security
   ```

2. **Test authentication:**
   ```bash
   ticketr schema --verbose
   ```

3. **Common causes:**
   - Incorrect JIRA_URL (missing `https://`)
   - Wrong email address
   - Expired or invalid API token
   - No project access permissions

**Solutions:**

- **Generate new API token:**
  1. Visit https://id.atlassian.com/manage-profile/security/api-tokens
  2. Create new token
  3. Update JIRA_API_KEY

- **Verify URL format:**
  ```bash
  export JIRA_URL="https://yourcompany.atlassian.net"  # ✅ Correct
  export JIRA_URL="yourcompany.atlassian.net"          # ❌ Missing https://
  ```

### Network/Timeout Errors

**Problem:** Connection timeouts or network errors

**Solutions:**

1. **Check network connectivity:**
   ```bash
   ping yourcompany.atlassian.net
   curl -I https://yourcompany.atlassian.net
   ```

2. **Check firewall/proxy:**
   - Corporate firewall blocking JIRA API?
   - Proxy configuration needed?

3. **Verify JIRA is accessible:**
   - Try accessing JIRA in browser
   - Check JIRA status page

---

## Field & Schema Issues

### Field Not Found

**Problem:** `field not found: Story Points`

**Root Cause:** Field name doesn't match your JIRA instance

**Solution:**

1. **Discover available fields:**
   ```bash
   ticketr schema
   ticketr schema --verbose  # More details
   ```

2. **Check exact field names (case-sensitive!):**
   ```markdown
   ## Fields
   - Story Points: 5      # ❌ May not exist
   - story points: 5      # ❌ Wrong case
   - Storypoints: 5       # ✅ Check with schema command
   ```

3. **Generate config file:**
   ```bash
   ticketr schema > .ticketr.yaml
   # Edit .ticketr.yaml to map fields
   ```

### Field Appears Blank in JIRA

**Problem:** Field synced but doesn't show in JIRA UI

**Causes:**
- Field not on JIRA screen for that issue type
- Field requires specific permissions
- Field ID is incorrect

**Solutions:**

1. **Verify data is in JIRA:**
   ```bash
   ticketr pull --project PROJ -o verify.md
   # Check if field is in pulled data
   ```

2. **Check JIRA screen configuration:**
   - Project Settings → Issue Types → Screens
   - Ensure field is on the appropriate screen

3. **Verify field permissions:**
   - Some fields require specific user roles

### Custom Field Mapping

**Problem:** Custom fields not syncing correctly

**Solution:**

1. **Create field mapping config:**
   ```bash
   ticketr schema > .ticketr.yaml
   ```

2. **Edit .ticketr.yaml:**
   ```yaml
   field_mappings:
     "Story Points":
       id: "customfield_10010"
       type: "number"
     "Sprint": "customfield_10020"
   ```

3. **Use config:**
   ```bash
   ticketr push tickets.md  # Auto-reads .ticketr.yaml
   ```

---

## State Management

### No Changes Detected (But I Made Changes!)

**Problem:** Ticketr says "no changes" but file was edited

**Cause:** `.ticketr.state` tracking old hash

**Solutions:**

1. **Force re-push by deleting state:**
   ```bash
   rm .ticketr.state
   ticketr push tickets.md
   ```

2. **Make a substantive change:**
   - Edit description or title (not just whitespace)
   - Changes must affect ticket content

3. **Check what's being tracked:**
   ```bash
   cat .ticketr.state  # JSON file with hashes
   ```

### State File Out of Sync

**Problem:** State file inconsistent with JIRA/local file

**Solution:**
```bash
# Nuclear option: reset state
rm .ticketr.state
ticketr pull --project PROJ -o tickets.md  # Re-initialize
```

### Should I Commit .ticketr.state?

**Answer:** **NO**

- Already in `.gitignore` (environment-specific)
- Each dev/CI environment should have own state
- State tracks local + remote changes per environment

---

## Push/Pull Issues

### First Pull Fails

**Problem:** Errors on first `ticketr pull`

**Solutions:**

1. **Ensure output directory exists:**
   ```bash
   mkdir -p output
   ticketr pull --project PROJ -o output/tickets.md
   ```

2. **Delete stale state:**
   ```bash
   rm .ticketr.state
   ticketr pull --project PROJ -o tickets.md
   ```

3. **Test with verbose:**
   ```bash
   ticketr pull --project PROJ -o test.md --verbose
   ```

### Conflict Detected

**Problem:** `⚠️ Conflict detected! Both local and remote changes`

**Understanding:**
- Local file changed since last pull
- JIRA ticket also changed remotely
- Ticketr won't overwrite without confirmation

**Solutions:**

1. **Accept remote changes:**
   ```bash
   ticketr pull --project PROJ --force
   ```

2. **Accept local changes:**
   - Push local changes:
   ```bash
   ticketr push tickets.md
   ```

3. **Manual merge:**
   - Review both versions
   - Edit file manually
   - Push when ready

### Push Fails Silently

**Problem:** Push completes but tickets not in JIRA

**Diagnostic:**

1. **Check logs:**
   ```bash
   cat .ticketr/logs/$(ls -t .ticketr/logs/ | head -1)
   ```

2. **Run with verbose:**
   ```bash
   ticketr push tickets.md --verbose
   ```

3. **Look for validation errors:**
   - Hierarchical violations?
   - Required fields missing?

### Subtasks Not Pulling

**Problem:** Parent tickets pull but subtasks missing

**Solution:**

- Subtasks pull automatically with parents
- If missing, check JIRA:
  ```bash
  ticketr pull --project PROJ --verbose
  # Check logs for subtask fetch errors
  ```

---

## Validation Errors

### Unsupported Story Heading

**Problem:** `Error: '# STORY:' format detected`

**Solution:** Update the heading manually:

```markdown
# STORY: Old Format        # ❌ Unsupported
# TICKET: New Format       # ✅ Supported
```

The parser intentionally blocks `# STORY:` so the fix is to rename the heading before running Ticketr again.

### Hierarchical Validation Error

**Problem:** `A 'Story' cannot be the child of an 'Epic'`

**JIRA Hierarchy Rules:**
- Epic → Story → Sub-task ✅
- Epic → Sub-task ❌
- Story → Sub-task ✅
- Task → Sub-task ✅

**Solution:**

Check your issue types:
```markdown
## Fields
- Type: Story  # If parent is Epic, use Story
- Type: Sub-task  # If parent is Story/Task
```

### Validation Fails But Need to Proceed

**Problem:** Validation errors but want to push valid tickets anyway

**Solution:**
```bash
ticketr push tickets.md --force-partial-upload
```

**Behavior:**
- Validation errors become warnings
- Valid tickets process successfully
- Invalid tickets skipped with errors
- Exit code 0 (partial success)

---

## Performance Issues

### Slow Push Operations

**Problem:** Push takes too long

**Optimizations:**

1. **State file is working?**
   - Only changed tickets should push
   - Check: `cat .ticketr.state`

2. **Break into smaller files:**
   ```bash
   # Instead of one large file:
   ticketr push sprint-23.md
   ticketr push sprint-24.md
   ```

3. **Use selective push:**
   - Only push files that changed

### Large Pull Operations

**Problem:** Pulling hundreds of tickets is slow

**Solutions:**

1. **Use selective queries:**
   ```bash
   # Instead of all project tickets:
   ticketr pull --project PROJ

   # Filter by JQL:
   ticketr pull --jql "sprint='Sprint 23'"
   ticketr pull --jql "updated >= -7d"
   ```

2. **Pull specific epics:**
   ```bash
   ticketr pull --epic PROJ-100
   ```

---

## Logging Issues

### No Log Files Created

**Problem:** `.ticketr/logs/` empty or missing

**Solutions:**

1. **Check directory permissions:**
   ```bash
   ls -la .ticketr/
   mkdir -p .ticketr/logs
   chmod 755 .ticketr/logs
   ```

2. **Check custom log dir:**
   ```bash
   echo $TICKETR_LOG_DIR
   # If set, logs go there instead
   ```

3. **Check disk space:**
   ```bash
   df -h .
   ```

### Sensitive Data in Logs?

**Concern:** API tokens in logs?

**Answer:** Ticketr automatically redacts:
- API tokens
- Email addresses
- Passwords
- Other credentials

**Always double-check before sharing logs!**

### Log Rotation Not Working

**Problem:** Old logs not being deleted

**Check:**
- Write permissions on log directory?
- Disk space available?
- 10+ log files exist?

```bash
ls -l .ticketr/logs/ | wc -l  # Should be ~10
```

---

## Exit Codes Reference

Understanding Ticketr exit codes:

| Exit Code | Meaning | Scenario |
|-----------|---------|----------|
| 0 | Success | All operations succeeded |
| 0 | Partial success | With --force-partial-upload, some items succeeded |
| 1 | Validation failure | Pre-flight validation failed (without --force-partial-upload) |
| 2 | Runtime error | JIRA API errors, network issues (without --force-partial-upload) |

**Examples:**

```bash
# Exit 0: All tickets pushed successfully
ticketr push tickets.md
echo $?  # 0

# Exit 1: Validation error without flag
ticketr push invalid-tickets.md
echo $?  # 1

# Exit 0: Partial success with flag
ticketr push invalid-tickets.md --force-partial-upload
echo $?  # 0 (valid tickets succeeded)

# Exit 2: JIRA API error without flag
ticketr push tickets.md  # Network error
echo $?  # 2
```

---

## Common Error Messages

### "accepts 1 arg(s), received 0"

**Problem:** Missing file argument

**Solution:**
```bash
ticketr push            # ❌ Missing file
ticketr push tickets.md # ✅ Correct
```

### "failed to load local tickets: file not found"

**Problem:** File doesn't exist

**Solution:**
```bash
ls -la tickets.md  # Verify file exists
ticketr push ./path/to/tickets.md  # Use correct path
```

### "rate limit exceeded"

**Problem:** Too many API requests

**Solution:**
- Wait a few minutes
- JIRA rate limits API calls
- Reduce frequency of operations

---

## Diagnostic Commands

### Health Check

```bash
# 1. Verify installation
ticketr --version

# 2. Test authentication
ticketr schema

# 3. Check configuration
echo $JIRA_URL
echo $JIRA_EMAIL
echo $JIRA_PROJECT_KEY

# 4. Verify file syntax
grep -n "#\s*STORY:" tickets.md  # Ensure unsupported headings are updated

# 5. Test with verbose
ticketr push tickets.md --verbose

# 6. Check logs
cat .ticketr/logs/$(ls -t .ticketr/logs/ | head -1)
```

### Environment Debugging

```bash
# Print all Ticketr-related env vars
env | grep JIRA
env | grep TICKETR

# Verify Go installation
go version
go env GOPATH

# Check permissions
ls -la ticketr
ls -la .ticketr/
```

---

## Getting Help

If you're still stuck:

1. **Check Documentation:**
   - [WORKFLOW.md](WORKFLOW.md) - Complete usage guide
   - [State Management](state-management.md) - Understanding .ticketr.state
   - [FAQ](https://github.com/karolswdev/ticktr/wiki/FAQ)

2. **Search Issues:**
   - [GitHub Issues](https://github.com/karolswdev/ticktr/issues)
   - Someone may have had the same problem

3. **Ask Community:**
   - [GitHub Discussions](https://github.com/karolswdev/ticktr/discussions)
   - Describe your problem with:
     - Ticketr version
     - OS and Go version
     - Steps to reproduce
     - Relevant logs (redacted!)

4. **File Bug Report:**
   - [New Issue](https://github.com/karolswdev/ticktr/issues/new/choose)
   - Use bug report template
   - Include diagnostic info

5. **Read the Logs:**
   ```bash
   cat .ticketr/logs/$(ls -t .ticketr/logs/ | head -1)
   ```
   Logs often contain the exact error!

---

## Emergency Fixes

### Nuclear Option: Complete Reset

If all else fails:

```bash
# 1. Backup your markdown files!
cp tickets.md tickets.md.backup

# 2. Delete state
rm .ticketr.state
rm -rf .ticketr/logs

# 3. Re-pull from JIRA
ticketr pull --project PROJ -o fresh.md

# 4. Compare with your local file
diff tickets.md fresh.md

# 5. Manually merge if needed
vim tickets.md

# 6. Re-push
ticketr push tickets.md
```

### Clean Slate

```bash
# Remove all Ticketr artifacts
rm .ticketr.state
rm -rf .ticketr/
rm .ticketr.yaml

# Start fresh
ticketr push tickets.md
```

---

**Still having issues?** See [SUPPORT.md](../SUPPORT.md) for support options.
