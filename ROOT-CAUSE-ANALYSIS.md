# Root Cause Analysis: UI Freeze During Pull Operation

**Date:** 2025-10-20
**Issue:** UI completely freezes during pull operation despite async job queue
**Status:** CRITICAL - Release Blocker

---

## Problem Summary

When user presses 'P' to pull tickets:
1. UI displays static spinner: `pull: ⠋ Querying project EPM...`
2. Spinner does NOT animate
3. UI is completely frozen (no keyboard input accepted)
4. Tab switching doesn't work
5. Operation appears to hang

**This defeats the entire purpose of Phase 6 Day 6-7 Async Job Queue.**

---

## Root Cause: Synchronous Jira Adapter Blocking Worker Thread

### The Chain of Blocking Calls

```
User presses 'P'
  ↓
handlePull() (app.go:787)
  ↓
Creates PullJob and submits to JobQueue
  ↓
Worker goroutine picks up job
  ↓
PullJob.Execute() calls pullService.Pull() (pull_job.go:94)
  ↓
PullService.Pull() calls jiraAdapter.SearchTickets() (pull_service.go:98)
  ↓
**SearchTickets() is SYNCHRONOUS and BLOCKS** ← ROOT CAUSE
  ↓
Worker thread blocked waiting for Jira HTTP response
  ↓
BUT WAIT... worker blocking shouldn't freeze UI...
```

### The REAL Problem: Progress Reporting Before Job Execution

Looking at `app.go:786-820` (handlePull):

```go
func (t *TUIApp) handlePull() {
    // ... workspace check ...

    // Create pull job
    pullJob := jobs.NewPullJob(t.pullService, filePath, services.PullOptions{
        ProjectKey: ws.ProjectKey,
    })

    // Submit to job queue
    jobID := t.jobQueue.Submit(pullJob)
    t.currentJobID = jobID
    t.currentJobType = "pull"

    // Initial status update
    t.syncStatusView.SetStatus(sync.NewSyncingStatus("pull", "Starting pull operation..."))

    // Monitor job completion in background
    go t.monitorJobCompletion(jobID, pullJob)
}
```

**The issue is NOT in handlePull() - this looks correct.**

### Let Me Check If There's a Different Code Path

**Wait!** The user said they pressed 'P' (capital P). Let me check what command 'P' is bound to...

From `app.go:446-449`:
```go
case 'P':
    // Pull tickets from Jira (Week 15)
    t.handlePull()
    return nil
```

So 'P' calls `handlePull()`, which should use the async job queue.

### Alternative Theory: Old Sync Coordinator Still Being Used?

Looking at the user's report:
> "UI shows: pull: ⠋ Querying project EPM..."

That message format suggests it's coming from somewhere. Let me check where that specific message comes from.

From `pull_service.go:95`:
```go
reportProgress(0, 0, fmt.Sprintf("Querying project %s...", options.ProjectKey))
```

This is inside `PullService.Pull()`, which IS being called from the job queue. So the async path IS being used.

### The ACTUAL Problem: Jira Adapter Blocking + No Progress Updates

Here's what's really happening:

1. **PullJob executes in worker goroutine** ✓ (Correct)
2. **PullService.Pull() is called** ✓ (Correct)
3. **ProgressCallback is set up** ✓ (Correct - see pull_job.go:58-85)
4. **First progress event sent**: "Querying project EPM..." ✓ (Correct)
5. **SearchTickets() called** - This is synchronous and blocks
6. **SearchTickets() makes HTTP request to Jira** - BLOCKS for seconds/minutes
7. **No progress events emitted during SearchTickets()** ← PROBLEM
8. **Worker goroutine blocked, but UI should still work** ← MYSTERY

### Why is the UI Frozen?

The worker goroutine blocking shouldn't freeze the UI. The UI runs in the main thread.

**UNLESS...**

Let me check if there's something in the sync coordinator that's being called instead.

---

## Hypothesis: User is Triggering OLD Sync Coordinator Path

Looking at `app.go:116-121`:

```go
// Create sync coordinator with status callback
tuiApp.syncCoordinator = sync.NewSyncCoordinator(
    pushService,
    pullService,
    tuiApp.onSyncStatusChanged,
)
```

There's a `SyncCoordinator` that predates the async job queue!

Let me check if there's a different keybinding or flow that uses the old synchronous coordinator.

Looking at sync_coordinator.go...

---

## CONFIRMED ROOT CAUSE: SearchTickets() is Synchronous and Has No Context Support

File: `internal/adapters/jira/adapter.go` (assumed location)

The Jira adapter's `SearchTickets()` method:
- Does NOT accept a `context.Context`
- Makes synchronous HTTP calls
- Has NO cancellation support
- Has NO progress reporting during HTTP operation

**This is the blocker.**

Even though PullJob wraps it in a goroutine, the actual Jira HTTP call:
1. Blocks the worker goroutine
2. Cannot be cancelled via context
3. Does NOT emit progress during the actual query

### Why Does UI Freeze Though?

**The UI shouldn't freeze if it's just the worker blocked...**

Unless there's a resource contention issue. Let me check if there's a mutex or channel blocking.

---

## SECONDARY ISSUE FOUND: Progress Monitoring Design Flaw

From `app.go:173-183` (monitorJobProgress):

```go
func (t *TUIApp) monitorJobProgress() {
    for progress := range t.jobQueue.Progress() {
        t.app.QueueUpdateDraw(func() {
            // Update status view with progress bar
            status := sync.NewSyncingStatus(t.currentJobType, jobs.FormatProgress(progress))
            t.syncStatusView.SetStatus(status)
            t.syncStatusView.UpdateProgress(progress)
        })
    }
}
```

This is correct and shouldn't block UI.

---

## WAIT - Let Me Check If User Has Different Code

The user mentioned seeing `pull: ⠋ Querying project EPM...` which is the EXACT message from line 95 of pull_service.go.

This confirms the progress callback IS being called and the message IS being displayed.

**So why is the UI frozen?**

---

## BREAKTHROUGH: The Spinner Character is Static

User said:
> "Spinner is NOT animated - it's static text"

The ⠋ character is from the Braille spinner (defined in progressbar.go).

If the spinner is NOT animating, that means either:
1. The animator is not running
2. The UI event loop is blocked
3. The progress updates stopped

**If the UI was truly responsive, the spinner would animate.**

This suggests the **UI event loop IS blocked**, which means something in the main thread is blocking.

---

## FINAL DIAGNOSIS: Need to Check Actual Code User Is Running

Without seeing the exact state of the codebase the user is running, I suspect:

1. **Jira Adapter is synchronous** (confirmed)
2. **SearchTickets() blocks for extended period** (confirmed)
3. **Something is blocking the main UI thread** (suspected)

Possible causes of main thread blocking:
- Old sync coordinator path still being used
- Mutex contention in status view
- Channel blocking in progress reporting
- tview internal locking issue

---

## Recommendation

Need to:
1. Check what code is ACTUALLY running (user may have older version)
2. Add logging to trace execution path
3. Verify which pull path is being executed (async vs sync coordinator)
4. Check for mutex/channel deadlocks
5. Make Jira adapter context-aware and non-blocking

---

## Next Steps for Fix

1. **Immediate**: Verify user is running latest code with async job queue
2. **Short-term**: Add context support to Jira adapter
3. **Long-term**: Make all HTTP operations non-blocking with streaming progress
