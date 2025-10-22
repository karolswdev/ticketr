# Builder Investigation Report: Pull Operation Failure

**Agent:** Builder
**Date:** 2025-10-21
**Session:** Phase 6.5 Critical Fix Session
**Time:** 2.5 hours
**Status:** ROOT CAUSE IDENTIFIED + FIX PLAN READY

---

## Investigation Summary

Investigated critical bug where pull operation fails with truncated error:
```
pull failed: failed to fetch tickets from JIRA: search failed with status
```

**ROOT CAUSE FOUND:** Jira adapter is created once at TUI startup and never updated when user switches workspaces, causing credential mismatch.

---

## Files Examined

1. `/home/karol/dev/private/ticktr/internal/adapters/jira/jira_adapter.go` (line 776)
2. `/home/karol/dev/private/ticktr/internal/core/services/pull_service.go` (lines 91-118)
3. `/home/karol/dev/private/ticktr/internal/adapters/tui/app.go` (lines 246-255, 836-870)
4. `/home/karol/dev/private/ticktr/internal/tui/jobs/pull_job.go` (lines 50-110)
5. `/home/karol/dev/private/ticktr/cmd/ticketr/tui_command.go` (lines 63-85)
6. `/home/karol/dev/private/ticktr/cmd/ticketr/main.go` (lines 131-153)

---

## Key Findings

### Finding 1: Static Jira Adapter

**File:** `cmd/ticketr/tui_command.go:63`

```go
jiraAdapter, err := initJiraAdapter(nil)  // Created ONCE at startup
pullService := services.NewPullServiceWithDB(jiraAdapter, ...)
pushService := services.NewPushService(fileRepo, jiraAdapter, ...)
```

Adapter is created once with initial workspace credentials and never updated.

### Finding 2: No Adapter Recreation on Workspace Switch

**File:** `internal/adapters/tui/app.go:246`

```go
t.workspaceListView.SetWorkspaceChangeHandler(func(workspaceID string) {
    t.ticketTreeView.LoadTickets(workspaceID)
    // MISSING: Recreate jiraAdapter with new workspace credentials
})
```

Workspace change handler does NOT recreate Jira adapter.

### Finding 3: Pull Uses Current Workspace but Old Credentials

**File:** `internal/adapters/tui/app.go:838-858`

```go
func (t *TUIApp) handlePull() {
    ws, _ := t.workspaceService.Current()  // Gets CURRENT workspace
    pullJob := jobs.NewPullJob(t.pullService, filePath, services.PullOptions{
        ProjectKey: ws.ProjectKey,  // NEW project key
    })
    // But t.pullService has OLD Jira adapter with wrong credentials!
}
```

Mismatch between project key (new) and credentials (old).

### Finding 4: Error Message Truncation

**File:** `internal/adapters/jira/jira_adapter.go:776`

```go
if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("search failed with status %d: %s", resp.StatusCode, string(body))
}
```

Error IS formatted correctly with status code and body, but gets truncated in UI display.

---

## Fix Plan

### Fix 1: Add Jira Adapter Recreation Method (CRITICAL)

**File:** `internal/adapters/tui/app.go`

Add method to recreate Jira adapter when workspace changes:

```go
// recreateJiraAdapter creates a new Jira adapter for the current workspace
// and updates all sync services.
func (t *TUIApp) recreateJiraAdapter() error {
    // Get current workspace
    ws, err := t.workspaceService.Current()
    if err != nil {
        return fmt.Errorf("no active workspace: %w", err)
    }

    // Get workspace config
    config, err := t.workspaceService.GetConfig(ws.Name)
    if err != nil {
        return fmt.Errorf("failed to get workspace config: %w", err)
    }

    // Create new Jira adapter with current workspace credentials
    jiraAdapter, err := jira.NewJiraAdapterFromConfig(config, nil)
    if err != nil {
        return fmt.Errorf("failed to create Jira adapter: %w", err)
    }

    // Get file repository and state manager from existing services
    // (We need to recreate services but reuse these dependencies)
    fileRepo := filesystem.NewFileRepository()
    stateManager := state.NewStateManager("")

    // Get database adapter from path resolver
    pathResolver, err := services.GetPathResolver()
    if err != nil {
        return fmt.Errorf("failed to get path resolver: %w", err)
    }
    dbAdapter, err := database.NewSQLiteAdapter(pathResolver)
    if err != nil {
        return fmt.Errorf("failed to get database adapter: %w", err)
    }

    // Recreate services with new adapter
    t.pullService = services.NewPullServiceWithDB(jiraAdapter, dbAdapter, fileRepo, stateManager)
    t.pushService = services.NewPushService(fileRepo, jiraAdapter, stateManager)

    // Update sync coordinator with new services
    t.syncCoordinator = sync.NewSyncCoordinator(
        t.pushService,
        t.pullService,
        t.onSyncStatusChanged,
    )

    return nil
}
```

### Fix 2: Call Recreation on Workspace Change (CRITICAL)

**File:** `internal/adapters/tui/app.go:246`

Update workspace change handler:

```go
t.workspaceListView.SetWorkspaceChangeHandler(func(workspaceID string) {
    // Reload tickets
    if t.ticketTreeView != nil {
        t.ticketTreeView.LoadTickets(workspaceID)
        t.setFocus("ticket_tree")
    }

    // FIX: Recreate Jira adapter with new workspace credentials
    if err := t.recreateJiraAdapter(); err != nil {
        // Show error in status bar
        t.syncStatusView.SetStatus(sync.NewErrorStatus("workspace", err))
    }

    // Auto-close workspace panel after selection
    if t.workspaceSlideOut != nil && t.workspaceSlideOut.IsVisible() {
        t.toggleWorkspacePanel()
    }
})
```

### Fix 3: Show Full Error Messages (HIGH PRIORITY)

**File:** `internal/adapters/tui/views/sync_status.go`

Update error display to show full HTTP status:

```go
// SetStatus updates the sync status display
func (sv *SyncStatusView) SetStatus(status sync.SyncStatus) {
    sv.mu.Lock()
    defer sv.mu.Unlock()

    sv.status = status

    // Extract status code from error message if present
    // Format: "search failed with status 401: {json body}"
    errorMsg := ""
    if status.State == sync.StateError && status.Error != nil {
        errorMsg = status.Error.Error()
        
        // If error contains truncated Jira error, show helpful hint
        if strings.Contains(errorMsg, "search failed with status") {
            // Try to extract status code
            parts := strings.Split(errorMsg, "status")
            if len(parts) > 1 {
                errorMsg = fmt.Sprintf("%s (check Jira credentials and project access)", errorMsg)
            }
        }
    }

    sv.updateDisplay(errorMsg)
}
```

### Fix 4: Add Workspace Validation (MEDIUM PRIORITY)

**File:** `internal/adapters/tui/app.go`

Add validation before pull/push operations:

```go
func (t *TUIApp) handlePull() {
    // Get current workspace
    ws, err := t.workspaceService.Current()
    if err != nil || ws == nil {
        t.syncStatusView.SetStatus(sync.NewErrorStatus("pull", fmt.Errorf("no active workspace")))
        return
    }

    // NEW: Validate workspace has credentials
    config, err := t.workspaceService.GetConfig(ws.Name)
    if err != nil {
        t.syncStatusView.SetStatus(sync.NewErrorStatus("pull", 
            fmt.Errorf("workspace '%s' has no credentials configured", ws.Name)))
        return
    }

    if config.JiraURL == "" || config.Username == "" || config.APIToken == "" {
        t.syncStatusView.SetStatus(sync.NewErrorStatus("pull",
            fmt.Errorf("workspace '%s' is missing Jira credentials", ws.Name)))
        return
    }

    // ... rest of pull logic
}
```

---

## Testing Strategy

### Test Case 1: Workspace Switch Pull
1. Start TUI with workspace A active
2. Pull successfully from workspace A
3. Switch to workspace B
4. Pull from workspace B
5. **Expected:** Pull succeeds with workspace B credentials
6. **Current:** Pull fails with "search failed with status"

### Test Case 2: No Initial Workspace
1. Start TUI with no active workspace
2. Create workspace B
3. Switch to workspace B
4. Pull from workspace B
5. **Expected:** Pull succeeds with workspace B credentials
6. **Current:** Pull fails or uses ENV credentials

### Test Case 3: Multiple Jira Instances
1. Workspace A: companyA.atlassian.net, project KEY1
2. Workspace B: companyB.atlassian.net, project KEY2
3. Switch from A to B
4. Pull from workspace B
5. **Expected:** Connects to companyB.atlassian.net with KEY2
6. **Current:** Tries companyA.atlassian.net with KEY2 (404)

---

## Additional Improvements

### Improvement 1: Lazy Jira Adapter Creation

Instead of creating adapter at TUI startup, create it on-demand for each operation:

```go
func (t *TUIApp) getJiraAdapterForCurrentWorkspace() (ports.JiraPort, error) {
    ws, err := t.workspaceService.Current()
    if err != nil {
        return nil, err
    }

    config, err := t.workspaceService.GetConfig(ws.Name)
    if err != nil {
        return nil, err
    }

    return jira.NewJiraAdapterFromConfig(config, nil)
}
```

### Improvement 2: Service Factory Pattern

Create factory for services that takes workspace config:

```go
type SyncServiceFactory struct {
    workspaceService *services.WorkspaceService
    // ... other dependencies
}

func (f *SyncServiceFactory) CreatePullService(workspaceID string) (*services.PullService, error) {
    // Get workspace config
    // Create Jira adapter
    // Create and return pull service
}
```

### Improvement 3: Better Error Display

Add dedicated error modal for detailed error messages instead of status bar truncation.

---

## Recommended Approach

**IMMEDIATE (Phase 6.5):**
1. Implement Fix 1 (recreateJiraAdapter method)
2. Implement Fix 2 (call on workspace change)
3. Implement Fix 3 (better error display)
4. Test with workspace switching scenarios

**FOLLOW-UP (Phase 7):**
1. Implement Improvement 1 (lazy adapter creation)
2. Refactor to Service Factory pattern
3. Add comprehensive workspace validation
4. Add error details modal

---

## Files to Modify

1. **internal/adapters/tui/app.go** - Add recreateJiraAdapter, update workspace change handler
2. **internal/adapters/tui/views/sync_status.go** - Improve error display
3. **cmd/ticketr/tui_command.go** - (Optional) Document that adapter will be recreated

---

## Risk Assessment

**Risk Level:** LOW
- Changes are isolated to TUI workspace switching logic
- No changes to core services or Jira adapter
- Backward compatible (CLI still works)
- Easy to test and verify

**Breaking Changes:** NONE

**Migration Required:** NONE

---

## Conclusion

This investigation identified a critical architectural flaw where the Jira adapter becomes stale after workspace switching. The fix is straightforward: recreate the adapter when the workspace changes.

**Next Steps:**
1. Implement fixes 1-3 (estimated 1-2 hours)
2. Test workspace switching scenarios
3. Verify error messages are shown correctly
4. Deploy and validate with user

**Estimated Fix Time:** 1-2 hours
**Testing Time:** 30 minutes
**Total:** 1.5-2.5 hours

---

**Agent:** Builder
**Status:** Investigation Complete - Ready for Implementation
**Timestamp:** 2025-10-21
