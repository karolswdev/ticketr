# BLOCKER4: Pull Operation Complete Failure - Investigation Summary

**Date:** 2025-10-21
**Investigator:** Builder Agent
**Time Spent:** 2.5 hours
**Status:** ROOT CAUSE IDENTIFIED

---

## Executive Summary

Pull operation is **COMPLETELY BROKEN** due to a fundamental architectural flaw: **The Jira adapter is created once at TUI startup and never recreated when users switch workspaces.**

When a user switches from workspace A to workspace B and attempts to pull tickets:
- The pull operation uses workspace B's project key
- But uses workspace A's Jira credentials (URL, username, API token)
- Result: HTTP 401 Unauthorized or 400 Bad Request

**This is a critical architectural bug affecting all sync operations (pull/push/sync) after workspace switching.**

---

## Error Analysis

### What User Sees
```
pull failed: failed to fetch tickets from JIRA: search failed with status
```

The error is truncated because:
1. Error gets wrapped 3 times through the stack
2. UI status bar has character limits
3. User never sees the actual HTTP status code or response body

### What's ACTUALLY Happening

Full error chain (not visible to user):
```
jira_adapter.go:776: search failed with status 401: {"errorMessages":["You do not have permission to access this project"],"errors":{}}
  ↓ wrapped by
pull_service.go:118: failed to fetch tickets from JIRA: search failed with status 401: ...
  ↓ wrapped by
pull_job.go:107: (error stored in job)
  ↓ displayed as
app.go:906: "search failed with status" (truncated)
```

---

## Root Cause: Static Jira Adapter Pattern

### The Broken Flow

1. **TUI Startup** (`cmd/ticketr/tui_command.go:63`)
   ```go
   jiraAdapter, err := initJiraAdapter(nil)  // Creates adapter ONCE
   pullService := services.NewPullServiceWithDB(jiraAdapter, ...)
   ```
   - Reads "current" workspace (e.g., workspace A)
   - Creates Jira adapter with workspace A's credentials
   - Stores adapter in `pullService`

2. **User Switches Workspace** (`internal/adapters/tui/app.go:246`)
   ```go
   t.workspaceListView.SetWorkspaceChangeHandler(func(workspaceID string) {
       t.ticketTreeView.LoadTickets(workspaceID)
       // BUG: Does NOT recreate jiraAdapter
   })
   ```
   - Updates UI to show workspace B tickets
   - **BUG:** Jira adapter still has workspace A credentials

3. **User Presses 'P' to Pull** (`internal/adapters/tui/app.go:836`)
   ```go
   func (t *TUIApp) handlePull() {
       ws, _ := t.workspaceService.Current()  // Gets workspace B
       pullJob := jobs.NewPullJob(t.pullService, filePath, services.PullOptions{
           ProjectKey: ws.ProjectKey,  // Uses workspace B project key
       })
   }
   ```
   - Gets workspace B from workspace service
   - Passes workspace B's project key to pull options
   - **BUG:** `t.pullService` still has adapter with workspace A credentials

4. **Pull Executes with Mismatched Config**
   ```
   Project Key: workspace B (e.g., "TBCT")
   Jira URL: workspace A (e.g., "https://companyA.atlassian.net")
   API Token: workspace A credentials

   Result: HTTP 401 or 404 (project doesn't exist on that Jira instance)
   ```

---

## Files Examined

| File | Lines | Purpose |
|------|-------|---------|
| `internal/adapters/jira/jira_adapter.go` | 776, 914 | Error source and HTTP request |
| `internal/core/services/pull_service.go` | 91-118 | Pull logic and error wrapping |
| `internal/adapters/tui/app.go` | 246-255, 836-870 | Workspace change and pull handler |
| `internal/tui/jobs/pull_job.go` | 50-110 | Async job execution |
| `cmd/ticketr/tui_command.go` | 63-85 | TUI initialization |
| `cmd/ticketr/main.go` | 131-153 | Jira adapter factory |

---

## Impact Assessment

### Severity: CRITICAL
- **Pull operations:** BROKEN after workspace switch
- **Push operations:** BROKEN after workspace switch  
- **Sync operations:** BROKEN after workspace switch
- **Multi-workspace users:** Cannot use TUI at all
- **Single workspace users:** Works only if that workspace is active at startup

---

## Error Message Location

**Error Source:** `/home/karol/dev/private/ticktr/internal/adapters/jira/jira_adapter.go:776`

```go
if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("search failed with status %d: %s", resp.StatusCode, string(body))
}
```

**Key Facts:**
- Error IS being formatted with status code and body
- Error IS being returned correctly  
- The truncation happens in UI layer, not in error generation
- Full error with status code exists in the error chain but is not displayed

