# TUI UX Issues - Steward Recommendations Summary

**Quick Reference Guide** | **Date:** October 20, 2025

---

## Critical Issues Summary

| Issue | Priority | Fix Timeline | Effort |
|-------|----------|--------------|--------|
| **Blocking Pull Operations** | P0 (Blocker) | v3.1.1 (2 weeks) | 1 week |
| **Async Progress Indicators** | P1 (Critical) | v3.1.1 (2 weeks) | 2-3 days |
| **Poor Menu Structure** | P1 (Critical) | v3.1.2 (3 weeks) | 4-6 days |

---

## Quick Wins (Can Implement This Week)

### 1. Wire Progress Callbacks (1 day)
**Problem:** Users see static "Pulling tickets..." with no indication of progress.

**Fix:**
```go
// In SyncCoordinator.PullAsync()
opts.ProgressCallback = func(current, total int, message string) {
    t.app.QueueUpdateDraw(func() {
        t.syncStatusView.SetText(fmt.Sprintf(
            "⚡ Pulling: %d/%d tickets | %s",
            current, total, message,
        ))
    })
}
```

**Impact:** Immediate user feedback, reduces "is it frozen?" concerns.

---

### 2. Add Bottom Action Bar (2 days)
**Problem:** Users can't discover available actions without reading docs.

**Fix:**
```
╠═══════════════════════════════════════════════════════════════════╣
║ p:push P:pull s:sync r:refresh b:bulk c:create ?:help :cmd  q:quit║
╚═══════════════════════════════════════════════════════════════════╝
```

**Impact:** Improved discoverability without major refactor.

---

### 3. Enhance Command Palette (1 day)
**Problem:** Command palette (`:`) exists but is hidden.

**Fix:**
- Show keybindings next to command names
- Add command descriptions
- Make Ctrl+P open quick search variant

**Impact:** Better onboarding for new users.

---

## v3.1.1 Patch Release (2 weeks)

### Goal: Eliminate Blocking Behavior

**What's Being Fixed:**
1. **Job Queue System** (4 days)
   - Operations run in background
   - User can navigate while sync runs
   - ESC/Ctrl+C cancels operations

2. **Context Cancellation** (2 days)
   - Refactor service layer to accept `context.Context`
   - Propagate cancellation through Jira adapter

3. **Real-time Progress** (2 days)
   - Wire progress callbacks to status bar
   - Show `45/120 tickets` count
   - Display elapsed time

**Acceptance Criteria:**
- ✅ Pull operation does not block TUI
- ✅ User can cancel long operations
- ✅ Status bar shows real-time progress
- ✅ No regressions in existing functionality

---

## v3.1.2 Enhancement (1 week)

### Goal: Improve Discoverability & Progress Visibility

**What's Being Added:**
1. **Bottom Action Bar** (2 days)
   - Context-aware commands
   - Always visible keybindings

2. **Progress Bar Widget** (2 days)
   - ASCII progress bar: `████████░░░░`
   - Expandable detail panel
   - Per-ticket status

**Acceptance Criteria:**
- ✅ Action bar shows available commands
- ✅ Progress bar displays during operations
- ✅ Detail view shows per-ticket status

---

## v3.2.0 Major Release (4 weeks)

### Goal: Complete TUI Refactor

**What's Being Added:**
1. **Top Menu Bar** (1 week)
   - File | Workspace | Tickets | Sync | View | Help
   - F-key shortcuts (F5=Refresh, F6=Sync)

2. **Background Job Manager** (1 week)
   - Persistent job queue
   - Job history and retry
   - Multi-workspace concurrent sync

3. **Performance Optimization** (1 week)
   - Batch Jira API calls
   - Parallel ticket processing
   - Pull 1000+ tickets in <30 seconds

**Acceptance Criteria:**
- ✅ Menu bar provides full action discovery
- ✅ Background jobs persist across restarts
- ✅ 50%+ speed improvement for large datasets

---

## Architectural Patterns to Adopt

### 1. Async Operations (k9s-style)
```go
type JobQueue struct {
    jobs    []*Job
    updates chan JobUpdate
    mutex   sync.RWMutex
}

func (t *TUIApp) handlePull() {
    ctx, cancel := context.WithCancel(context.Background())
    job, _ := t.asyncPullService.PullAsync(ctx, opts)

    go func() {
        for update := range job.Updates() {
            t.app.QueueUpdateDraw(func() {
                t.syncStatusView.UpdateProgress(update)
            })
        }
    }()
}
```

### 2. Bottom Action Bar (lazygit-style)
```
╠═══════════════════════════════════════════════════════════════════╣
║ p:push P:pull s:sync r:refresh b:bulk c:create ?:help :cmd  q:quit║
╚═══════════════════════════════════════════════════════════════════╝
```

### 3. Real-time Progress (lazydocker-style)
```
╠═══════════════════════════════════════════════════════════════════╣
║ ⚡ Syncing: 45/120 tickets | 2.3 req/s | 00:35 elapsed | ~01:20 ETA║
╚═══════════════════════════════════════════════════════════════════╝
```

---

## Production Readiness

| Version | Status | Use Case |
|---------|--------|----------|
| **v3.1.0 (current)** | ❌ BETA | Not recommended for >100 tickets |
| **v3.1.1 (2 weeks)** | ⚠️ STABLE | Suitable for <1000 tickets |
| **v3.2.0 (6 weeks)** | ✅ PRODUCTION | Enterprise-ready |

---

## Recommended Action

### Immediate (This Sprint)
1. Implement quick wins (progress callbacks + action bar) - **3 days**
2. Begin v3.1.1 implementation (job queue) - **4 days**
3. Tag v3.1.1-beta1 for user testing - **End of week 2**

### Short-term (Next Sprint)
1. Complete v3.1.1 with cancellation support
2. Regression testing + documentation
3. Release v3.1.1 stable
4. Begin v3.1.2 enhancements

### Long-term (v3.2.0)
1. Complete menu bar architecture
2. Background job manager
3. Performance optimization
4. Mark as production-ready

---

## Risk Mitigation

**Feature Flags:**
```yaml
# .config/ticketr/config.yaml
features:
  tui:
    async_operations: true     # v3.1.1
    action_bar: true           # v3.1.2
    menu_bar: false            # v3.2.0 (beta)
```

**Rollback Plan:**
- If v3.1.1 introduces regressions, set `async_operations: false`
- Reverts to synchronous mode (blocking but stable)
- Hotfix release within 48 hours

---

## Next Steps

1. **User:** Review and approve phased implementation plan
2. **Director:** Create v3.1.1 milestone with Builder tasks
3. **Builder:** Implement quick wins (progress + action bar)
4. **Verifier:** Manual testing with 500+ ticket dataset
5. **Scribe:** Update TUI documentation with new features

---

**For detailed analysis, see:** `/home/karol/dev/private/ticktr/docs/STEWARD-UX-ARCHITECTURAL-ASSESSMENT.md`
