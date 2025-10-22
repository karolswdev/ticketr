# BLOCKER4 Investigation - Document Index

**Issue:** Pull operation completely broken with truncated error message
**Date:** 2025-10-21
**Investigator:** Builder Agent
**Status:** INVESTIGATION COMPLETE - READY FOR IMPLEMENTATION

---

## Quick Summary

Pull operation fails after workspace switching because Jira adapter is created once at startup and never updated with new workspace credentials.

**Root Cause:** Static Jira adapter pattern in TUI initialization
**Impact:** All sync operations (pull/push/sync) broken after workspace change
**Severity:** CRITICAL - Release blocker
**Fix Complexity:** LOW - 2 hours implementation + testing

---

## Investigation Documents

### 1. BLOCKER4-INVESTIGATION-SUMMARY.md
**Purpose:** Executive summary of the issue
**Contents:**
- Error analysis and what user sees vs what actually happens
- Root cause explanation with code flow
- Evidence trail through 6 key files
- Impact assessment (CRITICAL severity)
- Error message location and truncation details

**Read this first for:** High-level understanding of the problem

### 2. BUILDER-INVESTIGATION-REPORT.md
**Purpose:** Complete investigation report with fix plan
**Contents:**
- Detailed findings from all 6 files examined
- Complete fix plan with 4 specific fixes
- Code examples for each fix
- Testing strategy with 3 test cases
- Risk assessment and timeline estimates

**Read this for:** Implementation details and fix plan

### 3. ROOT-CAUSE-ANALYSIS.md
**Purpose:** Deep technical analysis of the architecture flaw
**Contents:**
- Full error flow trace from user action to HTTP request
- Architectural flaw explanation with code evidence
- Three scenarios that trigger the bug
- Why error message gets truncated
- Analysis of workspace configuration issues

**Read this for:** Technical deep dive and architecture understanding

---

## Key Files Referenced

All investigation documents reference these source files:

1. `/home/karol/dev/private/ticktr/internal/adapters/jira/jira_adapter.go:776`
   - Error source: "search failed with status %d: %s"

2. `/home/karol/dev/private/ticktr/internal/core/services/pull_service.go:91-118`
   - Pull service logic and error wrapping

3. `/home/karol/dev/private/ticktr/internal/adapters/tui/app.go:246-255, 836-870`
   - Workspace change handler (missing adapter recreation)
   - Pull handler (uses stale adapter)

4. `/home/karol/dev/private/ticktr/internal/tui/jobs/pull_job.go:50-110`
   - Async job execution wrapper

5. `/home/karol/dev/private/ticktr/cmd/ticketr/tui_command.go:63-85`
   - TUI initialization (creates adapter once)

6. `/home/karol/dev/private/ticktr/cmd/ticketr/main.go:131-153`
   - Jira adapter factory (reads "current" workspace once)

---

## The Bug in One Diagram

```
TUI Startup:
  ├─ Read "current" workspace (e.g., "internal")
  ├─ Create JiraAdapter with "internal" credentials
  ├─ Create PullService(jiraAdapter)
  └─ Store in TUIApp

User switches to "tbct" workspace:
  ├─ UI updates to show "tbct" tickets
  └─ BUG: JiraAdapter still has "internal" credentials

User presses 'P' to pull:
  ├─ handlePull() gets current workspace = "tbct"
  ├─ Passes "tbct" project key to pull options
  └─ But PullService has old JiraAdapter with "internal" credentials
      ↓
  HTTP Request:
  ├─ URL: "internal" workspace Jira URL
  ├─ Token: "internal" workspace API token
  └─ Project: "tbct" ← WRONG! This project doesn't exist on "internal" Jira
      ↓
  Result: HTTP 401 or 404
  Error: "search failed with status" (truncated)
```

---

## Fix Summary

### Required Changes (CRITICAL)

**File: internal/adapters/tui/app.go**

1. Add `recreateJiraAdapter()` method
2. Call it in workspace change handler
3. Improve error display to show full status code

**Estimated Time:** 1-2 hours implementation + 30 min testing

### Testing Checklist

- [ ] Switch from workspace A to B, pull succeeds
- [ ] Start with no workspace, create B, pull succeeds
- [ ] Multiple Jira instances (A and B), switch works
- [ ] Error messages show full HTTP status code
- [ ] Push operation also works after workspace switch
- [ ] Sync operation works after workspace switch

---

## Next Steps

1. **Read:** BLOCKER4-INVESTIGATION-SUMMARY.md (5 minutes)
2. **Review:** BUILDER-INVESTIGATION-REPORT.md fix plan (10 minutes)
3. **Implement:** Fixes 1-3 from report (1-2 hours)
4. **Test:** Workspace switching scenarios (30 minutes)
5. **Verify:** With user UAT (30 minutes)

---

## Questions Answered

**Q: Why is error message truncated?**
A: Error is wrapped 3 times through the stack and UI status bar has character limits. The full error exists but isn't displayed.

**Q: What HTTP status code is being returned?**
A: Likely 401 Unauthorized (wrong credentials) or 404 Not Found (project doesn't exist on that Jira instance).

**Q: Why does this only affect TUI?**
A: CLI creates fresh Jira adapter for each operation. TUI reuses one adapter across operations.

**Q: Does this affect single-workspace users?**
A: Only if they don't have an active workspace at TUI startup, or if their default workspace differs from their working workspace.

**Q: Is this hard to fix?**
A: No - just need to recreate Jira adapter when workspace changes. Low risk, straightforward implementation.

---

**Investigation Status:** COMPLETE
**Fix Status:** READY FOR IMPLEMENTATION
**Agent:** Builder
**Date:** 2025-10-21
