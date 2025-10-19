# Critical Bugs Report - Phase 4 TUI Integration Issues

**Date**: 2025-10-19
**Severity**: CRITICAL
**Status**: Requires immediate Director delegation

---

## 📋 HANDOVER SUMMARY

### Where We Were
**Roadmap Position**: Phase 4 Week 16 COMPLETE → Ready for Phase 5
**Last Milestone**: TUI Final Polish (October 18, 2025)
**Status**: Phase 4 declared production-ready in `PHASE4-WEEK16-COMPLETE.md`

**What Was Just Completed:**
- ✅ 6 weeks of TUI implementation (Weeks 11-16)
- ✅ ~4,400 LOC TUI codebase
- ✅ Page navigation, themes, error handling, help system
- ✅ All Week 16 acceptance criteria met
- ✅ Tests passing, builds successful

**Next Planned Step**: Phase 5 - Backend sync implementation (per `docs/v3-implementation-roadmap.md`)

### What Stopped Us
During **manual end-to-end testing** of the completed Phase 4 TUI, we discovered **4 critical architectural bugs** that prevent the system from functioning correctly:

1. **PathResolver never integrated** - Phase 3 deliverable incomplete
2. **Workspace switching broken** - Core workflow non-functional
3. **TUI ticket loading broken** - Primary TUI feature doesn't work
4. **Pull has no progress** - Poor UX for long operations

**Impact**: Cannot proceed to Phase 5 until these P0/P1 bugs are resolved.

### Project Context
**Main Requirements**: `docs/v3-implementation-roadmap.md` (Phase 1-5 plan)
**Architecture**: Hexagonal (ports/adapters) - documented in `ARCHITECTURE.md`
**Agent Model**: Director orchestration - see `docs/DIRECTOR-ORCHESTRATION-GUIDE.md`
**Current Branch**: `feature/v3`
**Last Clean Commit**: `bd717a1` - "feat(pull): Auto-detect project key from current workspace"

### Deliverables Location
- **Phase Completion Reports**: `PHASE4-WEEK16-COMPLETE.md`, `docs/PHASE-3-COMPLETION-REPORT.md`
- **Implementation Plans**: `docs/PATHRESOLVER-INTEGRATION-PLAN.md`, `docs/v3-implementation-roadmap.md`
- **Requirements**: `REQUIREMENTS-v2.md` (legacy), Phase documents (current)
- **Test Evidence**: Git commits with test runs, completion docs
- **Architecture**: `ARCHITECTURE.md`, `CONTRIBUTING.md`

### Fix Strategy
**Approach**: Director delegates to Builder agents for P0/P1 fixes before returning to Phase 5 roadmap.

**Agent Delegation Plan:**
1. **Workspace Bug (P0)** → Builder agent (1-2 days)
2. **TUI Loading Bug (P0)** → Builder agent (2-3 days)
3. **Pull Progress (P2)** → Builder agent (1 day, low priority)
4. **PathResolver Integration (P1)** → Steward review + Builder (3-4 weeks, planned sprint)

**Success Criteria**: All P0 bugs fixed, TUI fully functional, can demonstrate complete workflow before proceeding to Phase 5.

**Return Path**: After bugs fixed → Commit fixes → Update Phase 4 status → Resume Phase 5 (Backend sync) from roadmap.

### Detailed Delegation Instructions

#### For Bug #2 (Workspace Switching) - P0
**Delegate To**: Builder agent
**Estimated Effort**: 1-2 days
**Files**: `internal/core/services/workspace_service.go:99-120`
**Task**: Fix `Switch()` method to persist workspace selection
**Options**:
- Option A: Make `Switch()` also call `SetDefault()`
- Option B: Store "current_workspace_id" in config file at `~/.config/ticketr/current.yaml`
- Option C: Use most-recently-used workspace as current (check `last_used` timestamp)

**Recommended**: Option C (most-recently-used) - aligns with UX expectations
**Test**: Create workspace → Switch → Exit → Re-run `workspace current` → Should show switched workspace

#### For Bug #3 (TUI Ticket Loading) - P0
**Delegate To**: Builder agent
**Estimated Effort**: 2-3 days
**Files**:
- `internal/adapters/tui/views/ticket_tree.go:45-51` (commented out loading)
- `internal/adapters/tui/app.go:130-135` (commented out async load)
**Task**: Implement async ticket loading with progress indicator
**Requirements**:
1. Uncomment initial load, move to goroutine
2. Show "Loading tickets..." in tree view
3. Wire up 'r' key to reload handler
4. Update tree view when tickets arrive

**Test**: Pull tickets → Open TUI → See tickets appear → Press 'r' → See tickets reload

#### For Bug #4 (Pull Progress) - P2
**Delegate To**: Builder agent
**Estimated Effort**: 1 day
**Files**: `internal/core/services/pull_service.go`
**Task**: Add progress callbacks to pull service
**Requirements**:
1. Add "Connecting to Jira..." message
2. Show "Found N tickets" after query
3. Show progress: "Pulling tickets: 45/101"
4. Show summary with counts

**Test**: Pull large ticket set (100+), verify progress updates appear

#### For Bug #1 (PathResolver Integration) - P1
**Delegate To**: Steward (review) + Builder (implementation)
**Estimated Effort**: 3-4 weeks (full sprint)
**Plan Document**: `docs/PATHRESOLVER-INTEGRATION-PLAN.md`
**Task**: Complete Phase 3 PathResolver integration (deferred work)
**Scope**: This is a separate sprint - NOT part of P0 fixes
**Note**: Can be deferred until after Phase 5, use local paths as workaround for now

### Definition of Done (P0 Bugs Only)

**Bug #2 (Workspace Switching) DONE when:**
- [ ] `workspace switch production` persists across commands
- [ ] `workspace current` shows the switched workspace
- [ ] TUI workspace list shows correct `*` indicator
- [ ] Integration test added: `TestWorkspaceSwitchPersistence`
- [ ] Manual test: Switch → exit → list → verify
- [ ] Commit with message: "fix(workspace): Persist workspace switching across command invocations"

**Bug #3 (TUI Ticket Loading) DONE when:**
- [ ] TUI shows tickets after successful pull
- [ ] Press 'r' reloads tickets from Jira
- [ ] Loading indicator shown during fetch
- [ ] No blocking on startup (async loading)
- [ ] Manual test: Pull 100 tickets → TUI shows all
- [ ] Commit with message: "fix(tui): Implement async ticket loading with reload support"

**Bug #4 (Pull Progress) DONE when:**
- [ ] Pull shows "Connecting..." at start
- [ ] Pull shows "Found N tickets" after query
- [ ] Pull shows progress for large sets (100+)
- [ ] Pull shows final summary with counts
- [ ] Manual test: Pull 100+ tickets, see all messages
- [ ] Commit with message: "feat(pull): Add progress indicators and summary"

### How to Return to Roadmap

**After P0 Bugs Fixed:**

1. **Update Phase 4 Status**
   - Edit `PHASE4-WEEK16-COMPLETE.md`
   - Add section: "## Post-Completion Bugs Fixed"
   - List Bug #2 and #3 fixes with commit hashes
   - Note Bug #1 (PathResolver) deferred to future sprint

2. **Commit Final State**
   ```bash
   git add -A
   git commit -m "fix(phase4): Resolve critical integration bugs discovered in testing

   - Fixed workspace switching persistence (Bug #2)
   - Fixed TUI ticket loading (Bug #3)
   - Added pull progress indicators (Bug #4)

   Phase 4 TUI is now fully functional and production-ready.
   PathResolver integration (Bug #1) deferred to separate sprint.

   Ready to proceed to Phase 5: Backend sync implementation."
   ```

3. **Resume Phase 5**
   - Open `docs/v3-implementation-roadmap.md`
   - Find "Phase 5: Advanced Features (Weeks 17-20)"
   - Begin with Milestone 18: Enhanced Capabilities
   - Follow roadmap implementation steps

4. **Update Tracking**
   - Mark Phase 4 as ✅ COMPLETE (with bug fixes noted)
   - Mark Phase 5 as 🔄 IN PROGRESS
   - Update project board/tracking with Phase 5 tasks

### Key Contacts & Resources

**Documentation:**
- Architecture: `ARCHITECTURE.md`
- Director Guide: `docs/DIRECTOR-ORCHESTRATION-GUIDE.md`
- Roadmap: `docs/v3-implementation-roadmap.md`
- Phase 3 Plan: `docs/PATHRESOLVER-INTEGRATION-PLAN.md`
- Phase 4 Completion: `PHASE4-WEEK16-COMPLETE.md`

**Code References:**
- Workspace Service: `internal/core/services/workspace_service.go`
- TUI App: `internal/adapters/tui/app.go`
- Ticket Tree View: `internal/adapters/tui/views/ticket_tree.go`
- Pull Service: `internal/core/services/pull_service.go`
- PathResolver: `internal/core/services/path_resolver.go`

**Test References:**
- Workspace Tests: `internal/core/services/workspace_service_test.go`
- TUI Tests: `internal/adapters/tui/`
- Integration Tests: `tests/integration/`

---

## Executive Summary

During Phase 4 TUI testing, **four critical architectural bugs** were discovered that prevent the TUI and workspace system from functioning correctly:

1. **PathResolver NOT integrated** - Files scattered in local directories instead of global workspace
2. **Workspace switching doesn't persist** - Switch command has no effect across invocations
3. **TUI ticket reload broken** - Press 'r' doesn't load tickets
4. **Pull command has no progress indicator** - User doesn't know if it's working

---

## Bug #1: PathResolver Implementation NOT Integrated (CRITICAL)

### Summary
PathResolver was implemented and tested in Phase 3 but **NEVER integrated** into the system. All files are still using hardcoded local paths instead of XDG-compliant global directories.

### Current Behavior
```bash
# Files created in project directory (WRONG)
./ticketr pull
ls -la .ticketr.state      # Created HERE
ls -la .ticketr/workspaces.db  # Created HERE
```

### Expected Behavior (Per Phase 3 Design)
```bash
# Linux:
~/.local/share/ticketr/workspaces.db
~/.local/share/ticketr/state/
~/.config/ticketr/config.yaml

# macOS:
~/Library/Application Support/ticketr/workspaces.db
~/Library/Application Support/ticketr/state/

# Windows:
%LOCALAPPDATA%\ticketr\workspaces.db
%APPDATA%\ticketr\config.yaml
```

### Evidence
- ✅ PathResolver implemented: `internal/core/services/path_resolver.go` (290 lines)
- ✅ PathResolver tested: 92.9% coverage, all tests passing
- ✅ Documentation exists: `docs/PATHRESOLVER-INTEGRATION-PLAN.md`
- ❌ **Integration status**: "Phase 3 Complete, Integration Pending"

### Impact
- **Data scattered** across project directories
- **No global state** - cannot switch between projects easily
- **Violates XDG spec** - non-standard directory structure
- **Migration issues** - users have fragmented state files

### Files Affected
```go
// Hardcoded local paths (WRONG):
cmd/ticketr/main.go:248:    stateManager := state.NewStateManager(".ticketr.state")
cmd/ticketr/main.go:389:    stateManager := state.NewStateManager(".ticketr.state")
cmd/ticketr/tui_command.go:73:  stateManager := state.NewStateManager(".ticketr.state")
cmd/ticketr/workspace_commands.go:164: adapter, err := database.NewSQLiteAdapter(features.SQLitePath)
```

### Fix Required
Per `docs/PATHRESOLVER-INTEGRATION-PLAN.md`:

1. Update `SQLiteAdapter` to accept `PathResolver`
2. Update all CLI commands to use PathResolver
3. Update state manager to use PathResolver paths
4. Implement migration from old paths to new paths
5. Update all tests

**Estimated Effort**: 3-4 weeks (per integration plan)

---

## Bug #2: Workspace Switching Doesn't Persist (CRITICAL)

### Summary
`workspace switch` command updates in-memory cache but doesn't persist selection. Every new command invocation falls back to default workspace.

### Current Behavior
```bash
./ticketr workspace switch production
# Switched to workspace 'production'

./ticketr workspace current
# Current workspace: my-project  ← WRONG! Should be 'production'
```

### Root Cause
```go
// internal/core/services/workspace_service.go:99
func (s *WorkspaceService) Switch(name string) error {
    // ... get workspace ...

    // Update in-memory cache (LOST when command exits!)
    s.currentMutex.Lock()
    s.currentCache = workspace
    s.currentMutex.Unlock()

    return nil  // No persistence!
}

// Current() always falls back to GetDefault():
func (s *WorkspaceService) Current() (*domain.Workspace, error) {
    if cached != nil {  // Always nil on new command!
        return cached, nil
    }
    return s.repo.GetDefault()  // Always uses default
}
```

### Impact
- Users cannot switch workspaces effectively
- TUI shows wrong workspace as current
- Confusing UX - switch appears to work but doesn't

### Fix Required
Option A: Make `Switch()` also call `SetDefault()`
Option B: Store "current workspace ID" in config file
Option C: Use most recently used workspace as current

**Estimated Effort**: 1-2 days

---

## Bug #3: TUI Ticket Reload Broken (HIGH)

### Summary
Pressing 'r' in TUI to reload tickets doesn't work. Tickets never appear even after successful pull.

### Current Behavior
```bash
./ticketr pull  # Successfully pulls 101 tickets
./ticketr tui
# Press 'r' to reload tickets
# → Nothing happens, ticket panel stays empty
```

### Root Cause
We commented out initial ticket loading to fix TUI hanging issues:

```go
// internal/adapters/tui/views/ticket_tree.go:45-51
// TEMPORARY FIX: Don't load tickets on startup to avoid blocking
// view.loadInitialTickets()  ← Commented out
view.showEmptyState()
```

But the reload handler ('r' key) is likely also broken or not wired up.

### Evidence
- Pulled 101 tickets to `pulled_tickets.md` (1,024 lines)
- State file `.ticketr.state` has all 101 tickets tracked
- TUI never shows tickets even after manual pull

### Fix Required
1. Check if 'r' key handler is wired up
2. Make ticket loading async (goroutine) to avoid blocking
3. Show loading indicator while fetching
4. Update tree view when tickets loaded

**Estimated Effort**: 2-3 days

---

## Bug #4: Pull Command Has No Progress Indicator (MEDIUM)

### Summary
Pull command provides no feedback while running. User doesn't know if it's working or how many tickets found.

### Current Behavior
```bash
./ticketr pull -v
09:49:52.055640 main.go:182: Verbose mode enabled
09:49:52.062156 main.go:354: Using project key from workspace 'produ': ATL
09:49:52.062171 main.go:382: Pulling tickets from project: ATL
# ... long pause (no feedback) ...
Successfully updated pulled_tickets.md
```

### Expected Behavior
```bash
./ticketr pull
Connecting to Jira (https://terumobct.atlassian.net)...
Querying project ATL...
Found 101 tickets
Pulling tickets: [========================================] 101/101
Successfully updated pulled_tickets.md
  - 101 tickets pulled
  - 0 conflicts detected
```

### Fix Required
1. Add "Connecting..." message
2. Show progress: "Pulling tickets: 45/101"
3. Show final summary with counts
4. Add spinner or progress bar for better UX

**Estimated Effort**: 1 day

---

## Additional Issues Discovered

### TUI Workspace List Display Bug
- Set "production" as default
- TUI shows "* produ" instead of "* production"
- Likely string truncation issue in display code

### Empty Tickets Despite Successful Pull
- Pull succeeded (101 tickets in state file)
- But wrote to `pulled_tickets.md` (not `tickets.md`)
- User expected `tickets.md` based on previous commands
- Minor UX issue, but confusing

---

## Recommended Action Plan

### Phase 1: Critical Workspace Fixes (Week 1)
**Delegate to Builder**
- [ ] Fix workspace switching persistence (Bug #2)
- [ ] Fix TUI workspace display truncation
- [ ] Add integration tests for workspace switching

### Phase 2: TUI Ticket Loading (Week 1-2)
**Delegate to Builder**
- [ ] Fix TUI ticket reload (Bug #3)
- [ ] Make ticket loading async with loading indicator
- [ ] Wire up 'r' key properly
- [ ] Test with large ticket sets

### Phase 3: UX Improvements (Week 2)
**Delegate to Builder**
- [ ] Add pull progress indicator (Bug #4)
- [ ] Show connection status
- [ ] Display ticket counts
- [ ] Add spinner/progress bar

### Phase 4: PathResolver Integration (Weeks 3-4)
**Delegate to Steward + Builder**
- [ ] Review PathResolver integration plan
- [ ] Update SQLiteAdapter signature
- [ ] Update all CLI commands
- [ ] Implement migration tool
- [ ] Update all tests
- [ ] Cross-platform testing

---

## Testing Requirements

### Integration Tests Needed
- [ ] Workspace switching persists across command invocations
- [ ] TUI loads tickets from workspace after pull
- [ ] TUI reload (r) fetches latest tickets
- [ ] Pull shows progress for large ticket sets
- [ ] PathResolver creates correct directories on all platforms

### Manual Testing Checklist
- [ ] Create workspace, switch to it, verify current workspace
- [ ] Pull tickets, open TUI, verify tickets appear
- [ ] Press 'r' in TUI, verify reload works
- [ ] Pull large ticket set, verify progress shown
- [ ] Test on Linux (XDG paths)
- [ ] Test on macOS (Application Support)
- [ ] Test on Windows (AppData)

---

## Priority Assessment

| Bug | Severity | Impact | User Blocking | Priority |
|-----|----------|--------|---------------|----------|
| #1 PathResolver | Critical | High | No (workaround exists) | P1 |
| #2 Workspace Switch | Critical | High | Yes | P0 |
| #3 TUI Reload | High | Medium | Yes | P0 |
| #4 Pull Progress | Medium | Low | No | P2 |

---

## Conclusion

These bugs represent **incomplete integration** of Phase 3 (PathResolver) and Phase 4 (TUI) components. The workspace system and TUI were built but critical glue code was never implemented.

**Immediate Action Required**:
1. Fix workspace switching (P0)
2. Fix TUI ticket reload (P0)
3. Add pull progress indicator (P2)
4. Plan PathResolver integration sprint (P1)

**Recommendation**: Delegate to Director with Builder specialization for implementation.

---

**Report By**: Claude Code
**Reviewed By**: _Pending Director Review_
**Next Step**: Director delegation with task breakdown
