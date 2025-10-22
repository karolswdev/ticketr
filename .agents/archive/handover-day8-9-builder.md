# Builder Agent Handover: Day 8-9 TUI Menu Structure

**Phase**: 6 - The Enchantment Release
**Milestone**: Day 8-9 TUI Menu Structure
**Status**: Ready for Implementation
**Assigned**: Builder Agent
**Estimated**: 8 hours
**Director**: Active Orchestration

---

## Mission

Transform the Ticketr TUI from functional to discoverable by implementing:
1. Bottom action bar widget (context-aware keybindings display)
2. Enhanced command palette (fuzzy search, descriptions, grouping)
3. F-key shortcuts (F1: Help, F2: Sync, F5: Refresh, F10: Exit)

**Goal**: All TUI actions must be discoverable without external documentation.

---

## Context & Background

### Current State (Day 6-7 Complete)
- ✅ Async job queue implemented (`internal/tui/jobs/`)
- ✅ Non-blocking pull operations functional
- ✅ Progress reporting working
- ✅ ESC/Ctrl+C cancellation operational

### TUI Architecture
**Primary Location**: `/home/karol/dev/private/ticktr/internal/adapters/tui/`

**Current Structure**:
```
internal/adapters/tui/
├── app.go                    # Main TUI application (tview.Application)
├── keybindings.go           # Current keybinding definitions
├── router.go                # View routing logic
├── theme/
│   └── theme.go             # Theme system
├── views/
│   ├── ticket_tree.go       # Main ticket list view
│   ├── ticket_detail.go     # Detail view
│   ├── workspace_list.go    # Workspace selector
│   ├── help.go              # Help screen
│   ├── sync_status.go       # Sync status display
│   └── ...
└── sync/
    └── coordinator.go       # Sync operations coordinator

internal/tui/jobs/           # Async job queue (Day 6-7)
├── job.go                   # Job interface
├── queue.go                 # JobQueue implementation
└── pull_job.go              # PullJob implementation
```

**Entry Point**: `/home/karol/dev/private/ticktr/cmd/ticketr/tui_command.go`
- Lines 84-100: NewTUIApp initialization
- Services wired: workspace, ticket query, push, pull, bulk operations

### Existing Keybindings (from tui_command.go lines 25-34)
```
q, Ctrl+C  - Quit
?          - Show help
Tab        - Switch between panels
j/k        - Navigate up/down (vim-style)
h/l        - Collapse/expand tree nodes
p          - Push tickets to Jira
P          - Pull tickets from Jira
r          - Refresh tickets
s          - Full sync (pull then push)
```

---

## Requirements & Acceptance Criteria

### From REQUIREMENTS-v2.md
- **TUI-001**: Context-aware menu system
- **TUI-002**: Discoverable keybindings
- **TUI-003**: Command palette with fuzzy search
- **UX-001**: No external docs needed for core actions

### From PHASE6-CLEAN-RELEASE.md (lines 519-624)
**Acceptance Criteria**:
- [ ] Bottom action bar visible in all TUI views
- [ ] Action bar shows context-aware keybindings
- [ ] Command palette functional with descriptions
- [ ] F-key shortcuts implemented
- [ ] All actions discoverable without external docs
- [ ] No keybinding conflicts
- [ ] Help comprehensive and accessible

---

## Implementation Tasks

### Task 1: Bottom Action Bar Widget

**Objective**: Full-width action bar (1-2 rows) displaying context-aware keybindings.

**Files to Create**:
- `/home/karol/dev/private/ticktr/internal/adapters/tui/widgets/actionbar.go`

**Implementation Steps**:

1. **Create ActionBar struct** (tview.TextView or custom tview.Box):
   ```go
   type ActionBar struct {
       *tview.TextView
       currentContext ViewContext
   }
   ```

2. **Define ViewContext enum**:
   ```go
   type ViewContext int
   const (
       ContextTicketList ViewContext = iota
       ContextTicketDetail
       ContextWorkspaceList
       ContextHelp
       ContextCommandPalette
   )
   ```

3. **Implement context-aware keybinding display**:
   - TicketList: `[Enter] Open | [Space] Select | [p] Push | [P] Pull | [s] Sync | [?] Help | [q] Quit`
   - TicketDetail: `[Esc] Back | [e] Edit | [d] Delete | [r] Refresh | [?] Help | [q] Quit`
   - WorkspaceList: `[Enter] Switch | [n] New | [d] Delete | [Esc] Back | [q] Quit`

4. **Styling**:
   - Use tview color tags for key highlights (e.g., `[white:blue][F1][-:-]`)
   - Full-width render (fit terminal width)
   - Fixed height (1-2 rows)
   - Position: Bottom of screen (above status bar if exists)

5. **Public Methods**:
   ```go
   func NewActionBar() *ActionBar
   func (ab *ActionBar) SetContext(ctx ViewContext)
   func (ab *ActionBar) Render() string
   ```

**Files to Modify**:
- `/home/karol/dev/private/ticktr/internal/adapters/tui/app.go`
  - Add ActionBar to main layout (tview.Flex)
  - Wire context changes from view router

**Testing**:
- Unit tests: Verify correct keybindings rendered per context
- Visual test: Launch TUI, navigate views, confirm bar updates

---

### Task 2: Command Palette Enhancement

**Objective**: Fuzzy-searchable command palette triggered by Ctrl+P or F1.

**Files to Create/Modify**:
- `/home/karol/dev/private/ticktr/internal/adapters/tui/widgets/palette.go` (new or enhance existing)
- `/home/karol/dev/private/ticktr/internal/adapters/tui/commands/registry.go` (new)

**Implementation Steps**:

1. **Create Command Registry** (`commands/registry.go`):
   ```go
   type Command struct {
       Name        string
       Description string
       Keybinding  string
       Category    string
       Handler     func()
   }

   type CommandRegistry struct {
       commands []Command
   }

   func NewCommandRegistry() *CommandRegistry
   func (cr *CommandRegistry) Register(cmd Command)
   func (cr *CommandRegistry) GetAll() []Command
   func (cr *CommandRegistry) Search(query string) []Command
   ```

2. **Define Commands**:
   - **Sync Category**:
     - "sync" → "Pull latest tickets from Jira" → Keybinding: `s`
     - "push" → "Push local changes to Jira" → Keybinding: `p`
     - "pull" → "Pull tickets from Jira" → Keybinding: `P`
     - "refresh" → "Refresh current view" → Keybinding: `r`
   - **Navigation Category**:
     - "help" → "Show help screen" → Keybinding: `?`
     - "workspace" → "Switch workspace" → Keybinding: (none, palette-only)
   - **View Category**:
     - "ticket-list" → "Go to ticket list" → Keybinding: (internal)
     - "workspace-list" → "Go to workspace list" → Keybinding: (internal)

3. **Implement Command Palette** (`widgets/palette.go`):
   - Use tview.InputField for search input
   - Use tview.List for filtered command results
   - Implement fuzzy matching (simple substring or use a library)
   - Display format: `[Category] Command Name - Description (Keybinding)`
   - Trigger: Capture Ctrl+P or F1 in main event handler

4. **Integration**:
   - Modal overlay (tview.Flex or tview.Modal)
   - Escape to close palette
   - Enter to execute selected command
   - Real-time filtering as user types

**Testing**:
- Unit tests: Command registration, fuzzy search correctness
- Integration test: Trigger palette, search "sync", execute command

---

### Task 3: F-Key Shortcuts

**Objective**: Implement F1-F10 shortcuts for common actions.

**Proposed Mapping**:
- **F1**: Help / Command Palette (toggle)
- **F2**: Sync (full sync)
- **F5**: Refresh current view
- **F10**: Exit application

**Files to Modify**:
- `/home/karol/dev/private/ticktr/internal/adapters/tui/keybindings.go`
  - Add F-key capture in event handler
- `/home/karol/dev/private/ticktr/internal/adapters/tui/app.go`
  - Wire F-key handlers to actions

**Implementation Steps**:

1. **Capture F-key events** in `app.go` or `keybindings.go`:
   ```go
   case tcell.KeyF1:
       // Open command palette or help
   case tcell.KeyF2:
       // Trigger sync
   case tcell.KeyF5:
       // Refresh current view
   case tcell.KeyF10:
       // Exit application
   ```

2. **Ensure consistency**:
   - F1 should toggle command palette (if palette exists)
   - F2 delegates to existing sync handler
   - F5 delegates to existing refresh handler
   - F10 delegates to existing quit handler

3. **Document in ActionBar**:
   - Update ActionBar context-aware display to show F-keys
   - Example: `[F1] Help | [F2] Sync | [F5] Refresh | [F10] Exit`

**Testing**:
- Manual testing: Press each F-key, verify action executes
- Ensure no keybinding conflicts (F-keys should override if collision)

---

## Deliverables

### New Files
1. `/home/karol/dev/private/ticktr/internal/adapters/tui/widgets/actionbar.go`
2. `/home/karol/dev/private/ticktr/internal/adapters/tui/widgets/actionbar_test.go`
3. `/home/karol/dev/private/ticktr/internal/adapters/tui/widgets/palette.go`
4. `/home/karol/dev/private/ticktr/internal/adapters/tui/widgets/palette_test.go`
5. `/home/karol/dev/private/ticktr/internal/adapters/tui/commands/registry.go`
6. `/home/karol/dev/private/ticktr/internal/adapters/tui/commands/registry_test.go`

### Modified Files
1. `/home/karol/dev/private/ticktr/internal/adapters/tui/app.go` (action bar integration, F-keys)
2. `/home/karol/dev/private/ticktr/internal/adapters/tui/keybindings.go` (F-key handlers)
3. `/home/karol/dev/private/ticktr/internal/adapters/tui/router.go` (context tracking for action bar)
4. `/home/karol/dev/private/ticktr/cmd/ticketr/tui_command.go` (update help text with F-keys)

### Documentation Updates (for Scribe)
- `docs/KEYBINDINGS.md` (new file needed)
- `docs/TUI-GUIDE.md` (update or create)
- `README.md` (update TUI section with new keybindings)

---

## Quality Standards

### Code Quality
- ✅ Follow tview patterns (extend tview.Box or tview.Primitive)
- ✅ Idiomatic Go (no tight coupling)
- ✅ Thread-safe if accessing job queue state
- ✅ Proper error handling (graceful degradation)

### Architecture Quality
- ✅ Widgets in `internal/adapters/tui/widgets/`
- ✅ Commands in `internal/adapters/tui/commands/`
- ✅ No business logic in widgets (presentation only)
- ✅ ActionBar receives context, doesn't query it

### Test Quality
- ✅ Unit tests for ActionBar rendering (each context)
- ✅ Unit tests for CommandRegistry (registration, search)
- ✅ Unit tests for fuzzy matching algorithm
- ✅ Integration test: Full TUI launch with new widgets visible

### Coverage Targets
- ActionBar: >70%
- CommandRegistry: >80%
- Palette: >60% (UI components harder to test)

---

## Testing Strategy

### Unit Tests
1. **ActionBar**:
   - Test each ViewContext renders correct keybindings
   - Test width adaptation (if dynamic)
   - Test color/styling (verify ANSI codes if applicable)

2. **CommandRegistry**:
   - Test command registration
   - Test duplicate command handling
   - Test search with various queries
   - Test empty search returns all commands

3. **Palette**:
   - Test fuzzy matching logic
   - Test command filtering
   - Test empty results handling

### Integration Tests
- Launch TUI in test mode (headless if possible)
- Trigger command palette (Ctrl+P)
- Verify command list rendered
- Execute a command, verify handler called

### Manual Testing (for Verifier)
- Launch TUI: `go run ./cmd/ticketr tui`
- Verify action bar visible at bottom
- Navigate between views (ticket list → detail → workspace list)
- Confirm action bar updates per view
- Press F1, verify command palette or help appears
- Press F2, verify sync triggers
- Press F5, verify refresh triggers
- Press F10, verify exit triggers
- Open command palette (Ctrl+P), type "sync", verify filter works
- Select command in palette, press Enter, verify execution

---

## Guardrails

### Never Do
- ❌ Hardcode terminal dimensions (use tview dynamic sizing)
- ❌ Block UI thread with long operations (use job queue from Day 6-7)
- ❌ Put business logic in widgets (delegate to services)
- ❌ Break existing keybindings without migration plan

### Always Do
- ✅ Test rendering in different terminal sizes
- ✅ Ensure F-key shortcuts are optional (fallbacks for terminals without F-key support)
- ✅ Use tview color tags for consistent theming
- ✅ Respect existing theme system (`internal/adapters/tui/theme/`)

---

## Handoff Criteria (Builder → Verifier)

✅ Ready to hand off when:
- [ ] All 6 new files created with unit tests
- [ ] All 4 modified files updated
- [ ] ActionBar renders correctly in all ViewContexts
- [ ] Command palette functional (Ctrl+P trigger, fuzzy search, execution)
- [ ] F-key shortcuts wired (F1, F2, F5, F10)
- [ ] All unit tests passing: `go test ./internal/adapters/tui/...`
- [ ] No keybinding conflicts detected
- [ ] Code formatted: `gofmt -l .` clean
- [ ] Build succeeds: `go build ./cmd/ticketr`
- [ ] Manual smoke test: TUI launches, action bar visible, palette works

❌ NOT ready if:
- Tests failing
- ActionBar not visible in TUI
- Command palette crashes or doesn't filter
- F-keys have no effect
- Existing keybindings broken

---

## Verifier Notes

**Areas Requiring Thorough Validation**:
1. **Keybinding Conflicts**: Test all existing keybindings still work (p, P, s, r, q, ?)
2. **Context Switching**: Rapidly switch views, verify action bar never shows stale context
3. **Terminal Compatibility**: Test in iTerm2, Alacritty, Windows Terminal, basic xterm
4. **F-Key Support**: Test on terminals with/without F-key support
5. **Command Palette Usability**: Verify fuzzy search intuitive (e.g., "syn" matches "sync")
6. **Discoverability**: Can new user navigate TUI using only action bar hints?

**Test Scenarios**:
- Navigate to ticket detail view → Verify action bar shows detail-specific keys
- Open command palette → Type partial command → Verify filtering
- Press F2 during active pull → Verify sync starts (or queues if pull ongoing)
- Resize terminal → Verify action bar width adapts
- Press invalid F-key (F3) → Verify no crash

---

## Scribe Notes

**Documentation Updates Needed**:

1. **Create docs/KEYBINDINGS.md**:
   - Table format: `| Key | Action | Context | Description |`
   - Include all standard keys (j/k, p/P, s, r, q, ?)
   - Include F-keys (F1, F2, F5, F10)
   - Note context-specific keys

2. **Update docs/TUI-GUIDE.md** (or create):
   - Section: "Navigation and Keybindings"
   - Section: "Command Palette"
   - Section: "Action Bar (Context-Aware Help)"
   - Screenshots: Action bar in different views, command palette open

3. **Update README.md**:
   - TUI section: Update keybindings list with F-keys
   - Mention command palette (Ctrl+P)
   - Link to docs/KEYBINDINGS.md

4. **Screenshots** (add to `docs/images/`):
   - `action-bar-ticket-list.png`
   - `action-bar-detail-view.png`
   - `command-palette.png`

---

## Success Checklist

Before reporting to Director:

- [ ] Reviewed PHASE6-CLEAN-RELEASE.md Day 8-9 section (lines 519-624)
- [ ] Code compiles: `go build ./cmd/ticketr`
- [ ] Unit tests passing: `go test ./internal/adapters/tui/...`
- [ ] Coverage measured: `go test -cover ./internal/adapters/tui/widgets/...`
- [ ] Coverage targets met (ActionBar >70%, Registry >80%, Palette >60%)
- [ ] Manual smoke test completed (TUI launch, action bar visible, palette works)
- [ ] No keybinding conflicts detected
- [ ] F-key shortcuts functional
- [ ] Code formatted: `gofmt -l .`
- [ ] Implementation summary prepared (files, line counts, test results)
- [ ] Verifier notes prepared (test scenarios, compatibility checks)
- [ ] Scribe notes prepared (documentation files needed, screenshot requirements)

---

## Cross-References

### Related Phase 6 Work
- **Day 6-7 (Complete)**: Async Job Queue (`internal/tui/jobs/`)
- **Day 10-11 (Next)**: Progress Indicators (will integrate with action bar/status display)
- **Day 12.5 (Future)**: TUI Visual Polish (shadows, animations, theming)

### Related Agent Definitions
- **Verifier Agent**: `.agents/verifier.agent.md` (will test your implementation)
- **Scribe Agent**: `.agents/scribe.agent.md` (will document your changes)
- **TUIUX Agent**: `.agents/tuiux.agent.md` (may enhance visuals in Day 12.5)

### Related Documentation
- **Phase 6 Plan**: `docs/PHASE6-CLEAN-RELEASE.md` (lines 519-624)
- **Architecture**: `docs/ARCHITECTURE.md` (TUI adapter layer)
- **Requirements**: `REQUIREMENTS-v2.md` (TUI-001, TUI-002, TUI-003, UX-001)

---

## Director's Expectations

You are implementing **critical UX improvements** that transform Ticketr from "functional" to "discoverable and delightful." Every keybinding, every hint in the action bar, every command description is an opportunity to respect the user's time and attention.

**The Four Principles of TUI Excellence** (from PHASE6-CLEAN-RELEASE.md):
1. **Subtle Motion is Life** (Day 12.5)
2. **Light, Shadow, and Focus** (Day 12.5)
3. **Atmosphere and Ambient Effects** (Day 12.5, optional)
4. **Small Charms of Quality** ← **YOU ARE HERE** (Day 8-9)

Your action bar and command palette are **small charms of quality**. They show craftsmanship. Make them polished, intuitive, and helpful.

**Quality over speed. Discoverability over feature count.**

---

## Estimated Time Breakdown

- **Task 1 (Action Bar)**: 3 hours (design, implementation, tests, integration)
- **Task 2 (Command Palette)**: 3.5 hours (registry, palette, fuzzy search, tests)
- **Task 3 (F-Key Shortcuts)**: 1 hour (event handling, wiring)
- **Testing & Validation**: 0.5 hours (smoke tests, manual verification)

**Total**: 8 hours

---

## Ready to Begin

This handover document provides all context, requirements, and quality standards needed. You have full autonomy to implement within these guardrails.

**When ready, report back to Director with your standard deliverable format (Implementation Complete).**

Good luck, Builder. Make it discoverable. Make it delightful.

---

**Handover Created**: 2025-10-20
**Phase**: 6 - The Enchantment Release
**Milestone**: Day 8-9 TUI Menu Structure
**Director**: Orchestrating Phase 6 Execution
