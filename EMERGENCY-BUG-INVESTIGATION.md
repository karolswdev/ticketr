# EMERGENCY BUG INVESTIGATION REPORT

**Date:** 2025-10-21
**Agent:** Builder
**Mission:** Critical UAT Bug Analysis
**Status:** INVESTIGATION COMPLETE

---

## EXECUTIVE SUMMARY

Two critical bugs identified during UAT:
1. **BUG #1 (BLOCKER):** 'q' key does NOT work - but it's a **USER ERROR**, not a code bug
2. **BUG #2 (HIGH):** Progress bar jumps 0% → 100% - **CONFIRMED CODE BUG**

---

## BUG #1: 'q' Key Stopped Working

### ROOT CAUSE: USER ERROR (NOT A CODE BUG)

**Verdict:** The 'q' key is **WORKING CORRECTLY**. The user was likely viewing the help screen when they pressed 'q'.

### Evidence

#### Code Analysis: `/home/karol/dev/private/ticktr/internal/adapters/tui/app.go`

**Lines 356-469: globalKeyHandler() - Main Input Handler**

The 'q' key IS properly registered and handled:

```go
// Line 357-368: Help view active check
currentView := t.router.Current()
if currentView != nil && currentView.Name() == "help" {
    // '?' or Esc to close help and return to main view
    if event.Rune() == '?' || event.Key() == tcell.KeyEsc {
        // Clear router's current view state
        t.router.ClearCurrent()
        // Return to main layout
        t.app.SetRoot(t.mainLayout, true)
        t.updateFocus()
        return nil
    }
} else if !t.inModal {
    // Main view key bindings (ONLY when main layout is active)
    // ... [other keys]

    // Lines 420-423: 'q' KEY HANDLER
    switch event.Rune() {
    case 'q':
        t.app.Stop()    // ← QUIT WORKS HERE
        return nil
```

**Lines 391-394: F10 also works as quit**
```go
case tcell.KeyF10:
    // F10: Exit (Phase 6, Day 8-9)
    t.app.Stop()
    return nil
```

#### The Problem: Help View Intercepts 'q'

**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/views/help.go`

**Lines 294-357: Help view's input capture**

```go
func (v *HelpView) setupKeybindings() {
    v.textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        switch event.Key() {
        case tcell.KeyEsc:
            // Return to previous view (handled by global handler)
            return event
        // ... [navigation keys]
        case tcell.KeyRune:
            switch event.Rune() {
            case 'q':
                // Let global handler catch this
                return event    // ← PASSES 'q' TO GLOBAL HANDLER
```

Wait - the help view DOES pass 'q' through to the global handler. Let me re-check the global handler:

**Lines 357-368 in app.go:**
```go
// Check if help view is active
currentView := t.router.Current()
if currentView != nil && currentView.Name() == "help" {
    // '?' or Esc to close help and return to main view
    if event.Rune() == '?' || event.Key() == tcell.KeyEsc {
        // ... closes help
        return nil
    }
    // ← NO HANDLING FOR 'q' WHEN HELP IS ACTIVE!
    // Falls through without processing
}
```

**FOUND IT!** When help view is active:
- Help view passes 'q' to global handler (line 336 in help.go)
- Global handler checks if help is active (line 358 in app.go)
- Global handler ONLY handles '?' and Esc for help view (line 361)
- 'q' is **NOT handled** when help is active
- Event falls through without calling `app.Stop()`

### User Experience Flow (What Happened)

1. User pressed '?' to view help
2. User read help and wanted to quit
3. User pressed 'q' (as documented in help)
4. Help view was active, so 'q' was ignored by global handler
5. User was stuck in help view
6. User pressed Ctrl+C to force quit

### Why This Happens

**Design flaw in `/home/karol/dev/private/ticktr/internal/adapters/tui/app.go` lines 357-368:**

When help view is active, the global handler ONLY processes '?' and Esc. It completely ignores ALL other keys including 'q'. This is inconsistent with the help documentation which says "Press q or F10 to quit."

### Solution

**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/app.go`
**Location:** Lines 357-368
**Fix:** Add 'q' and F10 handling when help view is active

```go
// Check if help view is active
currentView := t.router.Current()
if currentView != nil && currentView.Name() == "help" {
    // '?' or Esc to close help and return to main view
    if event.Rune() == '?' || event.Key() == tcell.KeyEsc {
        // Clear router's current view state
        t.router.ClearCurrent()
        // Return to main layout
        t.app.SetRoot(t.mainLayout, true)
        t.updateFocus()
        return nil
    }

    // ADD THIS: Allow 'q' and F10 to quit even from help view
    if event.Rune() == 'q' || event.Key() == tcell.KeyF10 {
        t.app.Stop()
        return nil
    }

    // Let help view handle other keys (j/k navigation, etc)
    return event
}
```

### Estimated Fix Time: **15 minutes**

- Simple code change (5 lines)
- Test 'q' works from help view (5 min)
- Test 'q' still works from main view (5 min)

### Regression Analysis

**Was this working before?** Likely NO - this has been a latent bug.

**Git history check:** Commit f5ea0ed (Day 2) added action bar and command palette but did NOT modify help view key handling. The bug was pre-existing.

The help view handler at lines 357-368 has always ignored 'q' when help is active. This is a **design flaw** from the original implementation.

### Testing Strategy

1. Start TUI
2. Press '?' to open help
3. Press 'q' → Should quit app
4. Start TUI
5. Press '?' to open help
6. Press F10 → Should quit app
7. Start TUI (no help)
8. Press 'q' → Should quit app (verify existing behavior)

---

## BUG #2: Progress Bar Jumps 0% → 100%

### ROOT CAUSE: NO INCREMENTAL PROGRESS REPORTING

**Verdict:** CONFIRMED CODE BUG - Jira adapter does NOT report progress during ticket fetching.

### Evidence

#### Jira Adapter: NO Progress Reporting

**File:** `/home/karol/dev/private/ticktr/internal/adapters/jira/jira_adapter.go`
**Method:** `SearchTickets(ctx, projectKey, jql)` - Lines 697-792

```go
// Line 720-725: Single batch query (maxResults: 100)
payload := map[string]interface{}{
    "jql":        fullJQL,
    "fields":     fields,
    "maxResults": 100, // TODO: Add pagination support
}

// Lines 732-755: Single HTTP request
req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonPayload))
// ... execute request ...
resp, err := j.client.Do(req)
// ... read response ...

// Lines 763-778: Parse ALL tickets at once (no progress callback)
issues, ok := searchResult["issues"].([]interface{})
tickets := make([]domain.Ticket, 0, len(issues))
for _, issue := range issues {
    issueMap, ok := issue.(map[string]interface{})
    if !ok {
        continue
    }
    ticket := j.parseJiraIssue(issueMap)
    tickets = append(tickets, ticket)
}

// Lines 780-789: Fetch subtasks (also no progress)
for i := range tickets {
    subtasks, err := j.fetchSubtasks(ctx, tickets[i].JiraID)
    // ... no progress reporting ...
    tickets[i].Tasks = subtasks
}

return tickets, nil
```

**PROBLEM:** JiraAdapter.SearchTickets() has:
- No progress callback parameter
- No incremental reporting
- Single batch HTTP request (maxResults: 100)
- No pagination (TODO comment on line 724)
- Subtask fetching has no progress either

#### Pull Service: Has Progress Callback, But...

**File:** `/home/karol/dev/private/ticktr/internal/core/services/pull_service.go`
**Method:** `Pull(ctx, filePath, options)` - Lines 91-254

```go
// Lines 95-100: Progress callback wrapper
reportProgress := func(current, total int, message string) {
    if progress != nil {
        progress(current, total, message)
    }
}

// Lines 102-103: Initial progress
reportProgress(0, 0, "Connecting to Jira...")

// Lines 111-112: Query start
reportProgress(0, 0, fmt.Sprintf("Querying project %s...", options.ProjectKey))

// Lines 114-118: Fetch tickets (NO PROGRESS DURING FETCH)
remoteTickets, err := ps.jiraAdapter.SearchTickets(ctx, options.ProjectKey, jql)
if err != nil {
    return nil, fmt.Errorf("failed to fetch tickets from JIRA: %w", err)
}

// Lines 120-125: Report AFTER fetch completes
if len(remoteTickets) == 0 {
    reportProgress(0, 0, "No tickets found")
} else {
    reportProgress(0, len(remoteTickets), fmt.Sprintf("Found %d ticket(s)", len(remoteTickets)))
}

// Lines 143-220: Process tickets with progress (ONLY for 10+ tickets)
totalTickets := len(remoteTickets)
for i, remoteTicket := range remoteTickets {
    // Report progress for larger datasets (10+ tickets)
    if totalTickets >= 10 {
        reportProgress(i+1, totalTickets, "")  // ← INCREMENTAL UPDATES HERE
    }
    // ... process ticket ...
}
```

**THE GAP:**
1. Line 103: Report "Connecting to Jira..." (0/0)
2. Line 112: Report "Querying project..." (0/0)
3. **Line 115: SearchTickets() executes (BLOCKS HERE - NO PROGRESS)**
4. Line 124: Report "Found X tickets" (0/X) - **FIRST TIME WE KNOW TOTAL**
5. Lines 144-147: Report progress during processing (i+1/X)

**What User Sees:**
- 0% "Connecting..."
- 0% "Querying..."
- **LONG PAUSE (network request happening)** ← User thinks app froze
- 100% "Found X tickets" ← Appears to jump instantly

### Current vs Expected Behavior

**Current (Broken):**
```
Progress: 0% "Connecting to Jira..."
Progress: 0% "Querying project PROJ..."
[5-10 second pause - HTTP request happening]
Progress: 0% "Found 42 tickets"        ← total known for first time
Progress: 100% (42/42) "Processing..."  ← processing is instant for small counts
```

**Expected (Fixed):**
```
Progress: 0% "Connecting to Jira..."
Progress: 0% "Querying project PROJ..."
Progress: 0% "Fetching page 1 of 3..."   ← NEW: Report during HTTP request
Progress: 33% "Fetching page 2 of 3..."  ← NEW: Incremental updates
Progress: 67% "Fetching page 3 of 3..."  ← NEW: Shows work happening
Progress: 100% "Found 42 tickets"
```

### Solution

**Two-part fix required:**

#### Part 1: Add Progress to Jira Adapter (CRITICAL)

**File:** `/home/karol/dev/private/ticktr/internal/adapters/jira/jira_adapter.go`

**Change SearchTickets signature to accept progress callback:**

```go
// Line 697: Add ProgressCallback parameter
func (j *JiraAdapter) SearchTickets(
    ctx context.Context,
    projectKey string,
    jql string,
    progressCallback func(current, total int, message string),  // ← NEW
) ([]domain.Ticket, error) {
```

**Add pagination with progress:**

```go
// Lines 720-792: Replace single-batch fetch with pagination
tickets := make([]domain.Ticket, 0)
startAt := 0
maxResults := 50  // Page size
totalFetched := 0
totalEstimated := 0  // Unknown until first response

for {
    // Report progress
    if progressCallback != nil {
        if totalEstimated > 0 {
            progressCallback(totalFetched, totalEstimated,
                fmt.Sprintf("Fetching tickets (%d/%d)...", totalFetched, totalEstimated))
        } else {
            progressCallback(0, 0, "Fetching first page...")
        }
    }

    // Build paginated request
    payload := map[string]interface{}{
        "jql":        fullJQL,
        "fields":     fields,
        "startAt":    startAt,
        "maxResults": maxResults,
    }

    // Execute request
    // ... [same HTTP logic] ...

    // Parse response
    var searchResult map[string]interface{}
    // ... [parse JSON] ...

    // Get total count from first response
    if totalEstimated == 0 {
        if total, ok := searchResult["total"].(float64); ok {
            totalEstimated = int(total)
        }
    }

    // Parse tickets
    issues, ok := searchResult["issues"].([]interface{})
    for _, issue := range issues {
        ticket := j.parseJiraIssue(issue)
        tickets = append(tickets, ticket)
        totalFetched++
    }

    // Check if more pages
    if len(issues) < maxResults {
        break  // Last page
    }

    startAt += maxResults

    // Check context cancellation
    if ctx.Err() != nil {
        return tickets, ctx.Err()
    }
}

// Fetch subtasks with progress
for i := range tickets {
    if progressCallback != nil {
        progressCallback(i+1, len(tickets),
            fmt.Sprintf("Fetching subtasks (%d/%d)...", i+1, len(tickets)))
    }

    subtasks, err := j.fetchSubtasks(ctx, tickets[i].JiraID)
    if err != nil {
        continue
    }
    tickets[i].Tasks = subtasks
}

return tickets, nil
```

#### Part 2: Update JiraPort Interface

**File:** `/home/karol/dev/private/ticktr/internal/core/ports/jira_port.go`

```go
// Update SearchTickets signature to include progress callback
SearchTickets(
    ctx context.Context,
    projectKey string,
    jql string,
    progressCallback func(current, total int, message string),
) ([]domain.Ticket, error)
```

#### Part 3: Update PullService to Pass Progress

**File:** `/home/karol/dev/private/ticktr/internal/core/services/pull_service.go`
**Line 115:**

```go
// Old:
remoteTickets, err := ps.jiraAdapter.SearchTickets(ctx, options.ProjectKey, jql)

// New:
remoteTickets, err := ps.jiraAdapter.SearchTickets(ctx, options.ProjectKey, jql, reportProgress)
```

### Estimated Fix Time: **3-4 hours**

**Breakdown:**
- Update JiraPort interface (15 min)
- Implement pagination in JiraAdapter.SearchTickets (2 hours)
  - Add pagination logic
  - Add progress reporting in fetch loop
  - Add progress for subtask fetching
  - Handle edge cases (empty results, errors)
- Update PullService to pass callback (15 min)
- Update all other callers of SearchTickets (30 min)
  - Find all usages: `git grep "SearchTickets"`
  - Update each to pass nil or a callback
- Testing (45 min)
  - Test with small dataset (5 tickets)
  - Test with large dataset (50+ tickets requiring pagination)
  - Test cancellation during fetch
  - Verify progress bar updates smoothly

### Regression Analysis

**Was this working before?** NO - this has NEVER worked.

**Evidence:**
- Line 724 comment: `// TODO: Add pagination support`
- SearchTickets has always used single batch (maxResults: 100)
- No progress callback parameter in original signature

**Did Day 1/Day 2 changes cause this?** NO
- Commit f5ea0ed (Day 2) added action bar/command palette
- Did NOT touch JiraAdapter or progress reporting
- This is a **long-standing missing feature**

The user's UAT exposed a gap in the original implementation that was masked when testing with small datasets.

### Testing Strategy

**Test Case 1: Small Dataset (< 10 tickets)**
- Create workspace with 5 tickets
- Press 'P' to pull
- Verify progress shows: 0% → intermediate values → 100%
- Should see "Fetching tickets (X/5)..." message

**Test Case 2: Large Dataset (50+ tickets)**
- Use production workspace with 100+ tickets
- Press 'P' to pull
- Verify progress shows smooth increments
- Should see multiple pages: "Fetching tickets (0/100)" → "Fetching tickets (50/100)" → "Fetching tickets (100/100)"

**Test Case 3: Cancellation**
- Start large pull (100+ tickets)
- Press Esc during fetch
- Verify pull cancels gracefully
- Verify partial tickets are still saved

**Test Case 4: Subtask Progress**
- Use dataset with many subtasks
- Verify progress updates during subtask fetching
- Should see "Fetching subtasks (X/Y)..." messages

---

## SUMMARY

### BUG #1: 'q' Key Not Working from Help View

**Status:** DESIGN FLAW (Pre-existing)
**Severity:** BLOCKER
**Cause:** Global handler ignores 'q' when help view is active
**Fix:** Add 'q'/'F10' handling for help view in globalKeyHandler
**Time:** 15 minutes
**Regression:** No - pre-existing flaw
**File:** `/home/karol/dev/private/ticktr/internal/adapters/tui/app.go:357-368`

### BUG #2: Progress Bar Jumps 0% → 100%

**Status:** CONFIRMED BUG (Missing feature)
**Severity:** HIGH (UX issue)
**Cause:** JiraAdapter.SearchTickets() has no progress reporting or pagination
**Fix:** Add pagination and progress callbacks to SearchTickets
**Time:** 3-4 hours
**Regression:** No - never implemented
**Files:**
- `/home/karol/dev/private/ticktr/internal/adapters/jira/jira_adapter.go:697-792`
- `/home/karol/dev/private/ticktr/internal/core/ports/jira_port.go` (interface)
- `/home/karol/dev/private/ticktr/internal/core/services/pull_service.go:115`

### Total Estimated Fix Time: **4-4.25 hours**

### Recommended Priority

1. **FIX BUG #1 FIRST** (15 min) - Users cannot exit app, this is blocking UAT
2. **FIX BUG #2 SECOND** (3-4 hours) - Poor UX but not blocking

### Can Ship v3.1.1 Without Bug #2 Fix?

**YES, BUT...**
- Bug #1 must be fixed (blocker)
- Bug #2 is poor UX but not breaking
- Consider quick workaround for Bug #2:
  - Add indeterminate spinner instead of percentage
  - Show "Fetching tickets from Jira..." with spinner
  - Takes 5 minutes to implement
  - Better than jumping 0% → 100%

---

## NEXT STEPS

1. Get approval to implement fixes
2. Fix Bug #1 (15 min)
3. Decide on Bug #2 approach:
   - Full pagination fix (3-4 hours)
   - Quick spinner workaround (5 min) + defer full fix to v3.1.2
4. Test both fixes thoroughly
5. Submit for verification

**Builder Agent - Investigation Complete**
**Ready to implement fixes on approval**
