# Handover Document: Ticketr v3 Phase 4 Week 12

## Quick Context

You are taking over development of **Ticketr**, a Jira-Markdown synchronization CLI tool, currently in **Phase 4 (TUI Implementation)** of a comprehensive v3 rewrite. Phase 4 Week 11 has just been completed. Your task is to continue with **Week 12: Multi-panel layout and ticket integration**.

## Project Overview

### What is Ticketr?

Ticketr is a command-line tool that synchronizes Jira tickets with local Markdown files, allowing developers to work offline and sync changes bidirectionally. It's built using **hexagonal architecture** (ports and adapters pattern).

**Core Value Proposition:**
- Work on Jira tickets offline using Markdown
- Store tickets in version control
- Bidirectional sync (pull from Jira, push changes back)
- Multi-workspace support (multiple Jira projects)
- Now adding: Terminal User Interface (TUI) for visual interaction

### Technical Stack

- **Language:** Go 1.23+
- **Module:** `github.com/karolswdev/ticktr`
- **Architecture:** Hexagonal (Ports & Adapters)
- **CLI Framework:** Cobra + Viper
- **Database:** SQLite (mattn/go-sqlite3)
- **TUI Framework:** tview + tcell/v2
- **Testing:** Standard Go testing + testify
- **VCS:** Git (branch: `feature/v3`)

### Architecture Pattern

```
┌─────────────────────────────────────┐
│         CLI/TUI Adapters            │  ← External interfaces
├─────────────────────────────────────┤
│      Application Services           │  ← Business logic
│  • WorkspaceService                 │
│  • TicketService                    │
│  • PathResolver                     │
├─────────────────────────────────────┤
│         Core Domain                 │  ← Entities, rules
│  • Ticket  • Workspace              │
├─────────────────────────────────────┤
│    Ports (Interface Definitions)    │  ← Contracts
├─────────────────────────────────────┤
│  Adapters (External Integrations)   │  ← Implementations
│  • SQLite  • Jira  • Filesystem     │
└─────────────────────────────────────┘
```

**Critical Rule:** Adapters (CLI, TUI) depend ONLY on services, NEVER directly on repositories or external systems.

## Current State (What Just Happened)

### Phase 4 Week 11 - TUI Skeleton (COMPLETED)

**Commit:** `326c357` - "feat(tui): Implement Phase 4 Week 11 - TUI adapter skeleton with tview"

**What Was Delivered:**

**11 files, 694 lines of code:**

```
internal/adapters/tui/
├── app.go                  (116 lines) - Main TUI application
├── router.go               (78 lines)  - View navigation
├── keybindings.go          (56 lines)  - Global keyboard handlers
└── views/
    ├── view.go             (19 lines)  - View interface
    ├── workspace_list.go   (105 lines) - Live workspace switcher
    ├── ticket_tree.go      (83 lines)  - Ticket browser (PLACEHOLDER)
    ├── ticket_detail.go    (95 lines)  - Ticket viewer (PLACEHOLDER)
    └── help.go             (84 lines)  - Help screen

cmd/ticketr/
└── tui_command.go          (58 lines)  - `ticketr tui` command
```

**Functionality:**
- ✅ Launch with `ticketr tui`
- ✅ Workspace list shows live data from WorkspaceService
- ✅ Enter key switches workspaces
- ✅ `q` / `Ctrl+C` quits
- ✅ `?` toggles help view
- ⏸️ Ticket tree shows placeholder data (not real tickets)
- ⏸️ Ticket detail shows placeholder content
- ⏸️ Single view at a time (no multi-panel layout yet)

**Key Design Patterns Implemented:**

1. **View Interface Pattern:**
```go
type View interface {
    Name() string
    Primitive() tview.Primitive
    OnShow()  // Called when view becomes active
    OnHide()  // Called when view is hidden
}
```

2. **Router Pattern:**
```go
router.Register(view)  // Add view to router
router.Show("view_name")  // Show view, call lifecycle hooks
router.Current()  // Get active view
```

3. **Dependency Injection:**
```go
func NewTUIApp(
    workspaceService *services.WorkspaceService,
    pathResolver *services.PathResolver,
) (*TUIApp, error)
```

### What Exists in the Codebase

**Services Available (internal/core/services/):**
- ✅ **WorkspaceService** - List(), Current(), Switch(), Create(), Delete()
- ✅ **PathResolver** - DatabasePath(), ConfigDir(), DataDir() (Phase 3)
- ⚠️ **TicketService** - EXISTS but needs verification for TUI integration
- ✅ **PushService** - Sync local changes to Jira
- ✅ **PullService** - Sync Jira changes to local

**Repositories Available (internal/adapters/):**
- ✅ SQLiteAdapter - Database access
- ✅ WorkspaceRepository - Workspace CRUD
- ✅ Jira adapter - Jira API integration
- ✅ Filesystem adapter - Markdown file operations

**Domain Models (internal/core/domain/):**
- ✅ Workspace - Workspace entity
- ✅ Ticket - Ticket entity with fields
- ✅ Config - Configuration structures

**TUI Framework (tview):**
- ✅ Installed: `github.com/rivo/tview v0.42.0`
- ✅ Installed: `github.com/gdamore/tcell/v2 v2.9.0`

## Your Task: Phase 4 Week 12

### Goal

Implement **multi-panel layout** and integrate **real ticket data** into the TUI.

### Roadmap Requirements

From `docs/v3-implementation-roadmap.md` (lines 272-369):

**Week 12 Deliverables:**
1. Multi-panel layout (workspace list + ticket tree side-by-side)
2. Panel focus switching with Tab key
3. Real ticket data integration (replace placeholders)
4. Load tickets on workspace switch
5. Vim-style navigation (j/k for up/down, h/l for collapse/expand)

### Acceptance Criteria

- [ ] TUI shows workspace list and ticket tree simultaneously
- [ ] Tab key switches focus between panels
- [ ] Focused panel has visual indicator (border color/style)
- [ ] Ticket tree loads real tickets from current workspace
- [ ] Changing workspace refreshes ticket tree
- [ ] Vim keybindings work (j/k/h/l)
- [ ] Arrow keys still work (backward compatibility)
- [ ] Loading state shown while fetching tickets
- [ ] Error handling for failed ticket loads
- [ ] All Week 11 functionality still works

### Technical Requirements

**Multi-Panel Layout:**
```go
// Use tview.Flex for side-by-side panels
flex := tview.NewFlex().
    AddItem(workspaceListView.Primitive(), 30, 0, true).  // Left: 30 chars wide
    AddItem(ticketTreeView.Primitive(), 0, 1, false)       // Right: Flex fill

// Focus management
func (t *TUIApp) FocusPanel(panel string) {
    // Set focused primitive
    // Update border colors (focused: green, unfocused: white)
    // Update currentFocus tracker
}
```

**Ticket Integration:**
```go
// In ticket_tree.go - replace OnShow()
func (v *TicketTreeView) OnShow() {
    // Show loading indicator
    v.showLoading()

    // Load tickets from service
    tickets, err := v.ticketService.List(v.currentWorkspace)
    if err != nil {
        v.showError(err)
        return
    }

    // Build tree from real data
    v.buildTree(tickets)
}

// Need to add TicketService to TicketTreeView constructor
func NewTicketTreeView(ticketService *services.TicketService) *TicketTreeView
```

**Panel Focus Switching:**
```go
// In keybindings.go - add Tab handler
case tcell.KeyTab:
    currentFocus := h.app.GetCurrentFocus()
    if currentFocus == "workspace_list" {
        h.app.FocusPanel("ticket_tree")
    } else {
        h.app.FocusPanel("workspace_list")
    }
    return true
```

## Development Approach (Recommended)

### Step 1: Verify Ticket Service Exists

**Action:**
```bash
# Check if TicketService exists
Read internal/core/services/ticket_service.go

# Understand its interface
# Look for: List(), Get(), Create(), Update(), Delete()
```

**If TicketService doesn't have a List() method:**
- Check how tickets are currently loaded (might be in PullService)
- May need to create a ListTickets() method
- Or use repository directly (check ports)

### Step 2: Design Multi-Panel Layout

**Action:**
Create design document or directly implement:
```go
// In app.go - modify setupApp()
func (t *TUIApp) setupApp() {
    // Create main flex container
    mainFlex := tview.NewFlex().
        SetDirection(tview.FlexColumn).
        AddItem(t.workspaceListView.Primitive(), 35, 0, true).
        AddItem(t.ticketTreeView.Primitive(), 0, 1, false)

    // Set as root
    t.app.SetRoot(mainFlex, true)

    // Track focus
    t.currentFocus = "workspace_list"

    // Setup keybindings
    t.setupGlobalKeys()
}
```

### Step 3: Implement Panel Focus Switching

**Files to Modify:**
- `internal/adapters/tui/app.go` - Add focus tracking
- `internal/adapters/tui/keybindings.go` - Add Tab handler

**Pattern:**
```go
// Add to TUIApp struct
type TUIApp struct {
    // ... existing fields
    currentFocus   string
    focusablePanels map[string]tview.Primitive
}

// Add method
func (t *TUIApp) FocusPanel(name string) {
    primitive := t.focusablePanels[name]
    t.app.SetFocus(primitive)
    t.currentFocus = name
    t.updatePanelBorders()
}

func (t *TUIApp) updatePanelBorders() {
    for panelName, _ := range t.focusablePanels {
        if panelName == t.currentFocus {
            // Set green border
        } else {
            // Set gray border
        }
    }
}
```

### Step 4: Integrate Real Ticket Data

**Files to Modify:**
- `internal/adapters/tui/views/ticket_tree.go` - Replace placeholder logic
- `internal/adapters/tui/app.go` - Inject TicketService

**Pattern:**
```go
// Update TicketTreeView struct
type TicketTreeView struct {
    tree          *tview.TreeView
    root          *tview.TreeNode
    ticketService *services.TicketService  // NEW
    workspace     string                    // NEW - current workspace
}

// Update constructor
func NewTicketTreeView(ticketService *services.TicketService) *TicketTreeView {
    // ... existing setup
    view.ticketService = ticketService
    return view
}

// Update OnShow
func (v *TicketTreeView) OnShow() {
    v.loadTickets()
}

func (v *TicketTreeView) loadTickets() {
    // Clear existing
    v.root.ClearChildren()

    // Load from service
    tickets, err := v.ticketService.List()
    if err != nil {
        v.showError(err)
        return
    }

    // Build tree
    for _, ticket := range tickets {
        node := tview.NewTreeNode(ticket.Key + ": " + ticket.Summary)
        v.root.AddChild(node)
    }
}
```

### Step 5: Add Vim-Style Navigation

**Files to Modify:**
- `internal/adapters/tui/views/ticket_tree.go` - Add j/k/h/l handlers

**Pattern:**
```go
// In NewTicketTreeView()
tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
    switch event.Rune() {
    case 'j':
        // Move down (same as Down arrow)
        return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
    case 'k':
        // Move up (same as Up arrow)
        return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
    case 'h':
        // Collapse node (same as Left arrow)
        return tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone)
    case 'l':
        // Expand node (same as Right arrow)
        return tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone)
    }
    return event
})
```

### Step 6: Workspace Switch Triggers Ticket Reload

**Files to Modify:**
- `internal/adapters/tui/views/workspace_list.go` - Add callback

**Pattern:**
```go
// In WorkspaceListView - add callback for workspace change
type WorkspaceListView struct {
    // ... existing
    onWorkspaceChanged func(workspaceName string)
}

func (v *WorkspaceListView) SetWorkspaceChangeHandler(handler func(string)) {
    v.onWorkspaceChanged = handler
}

// In handleInput when Enter is pressed
if v.onSwitch != nil {
    _ = v.onSwitch(workspaces[index].Name)
    v.refresh()

    // NEW: Trigger ticket reload
    if v.onWorkspaceChanged != nil {
        v.onWorkspaceChanged(workspaces[index].Name)
    }
}

// In app.go - wire it up
workspaceList.SetWorkspaceChangeHandler(func(wsName string) {
    ticketTree.SetWorkspace(wsName)
    ticketTree.OnShow() // Reload tickets
})
```

## Development Workflow

### 1. Create Todo List

Use TodoWrite tool to track Week 12 tasks:
```json
[
    {"content": "Verify TicketService interface and capabilities", "status": "pending", "activeForm": "Verifying TicketService..."},
    {"content": "Design multi-panel Flex layout", "status": "pending", "activeForm": "Designing multi-panel layout..."},
    {"content": "Implement panel focus switching with Tab", "status": "pending", "activeForm": "Implementing panel focus switching..."},
    {"content": "Add visual focus indicators (border colors)", "status": "pending", "activeForm": "Adding visual focus indicators..."},
    {"content": "Integrate TicketService with TicketTreeView", "status": "pending", "activeForm": "Integrating TicketService..."},
    {"content": "Replace placeholder data with real tickets", "status": "pending", "activeForm": "Replacing placeholder data..."},
    {"content": "Implement Vim-style navigation (j/k/h/l)", "status": "pending", "activeForm": "Implementing Vim-style navigation..."},
    {"content": "Wire workspace switch to trigger ticket reload", "status": "pending", "activeForm": "Wiring workspace switch..."},
    {"content": "Add loading indicators for async operations", "status": "pending", "activeForm": "Adding loading indicators..."},
    {"content": "Implement error handling and display", "status": "pending", "activeForm": "Implementing error handling..."},
    {"content": "Test multi-panel layout with real data", "status": "pending", "activeForm": "Testing multi-panel layout..."},
    {"content": "Verify all Week 11 functionality still works", "status": "pending", "activeForm": "Verifying Week 11 functionality..."}
]
```

### 2. Implementation Order

**Recommended sequence:**
1. ✅ Create todo list (TodoWrite)
2. ✅ Verify TicketService (Read, Grep)
3. ✅ Implement multi-panel layout (Edit app.go)
4. ✅ Add focus switching (Edit app.go, keybindings.go)
5. ✅ Visual focus indicators (Edit views)
6. ✅ Integrate TicketService (Edit ticket_tree.go)
7. ✅ Test with real data
8. ✅ Add Vim navigation (Edit ticket_tree.go)
9. ✅ Wire workspace reload (Edit workspace_list.go, app.go)
10. ✅ Error handling and loading states
11. ✅ Full integration test
12. ✅ Commit with detailed message

### 3. Agent Delegation Strategy

**Option A: Do It Yourself (Director)**
- Good for: Small changes, single-file edits
- Faster for: Simple pattern replication

**Option B: Delegate to Builder**
- Good for: Multi-file changes, complex logic
- Pattern: Create detailed prompt (see Week 11 example)
- Remember: Builder runs in sandbox - you'll need to materialize files

**Option C: Hybrid**
- Builder designs, Director implements
- Fastest for Week 12 scope

**Recommendation for Week 12:** Hybrid approach
- Ask Builder to design multi-panel layout
- Implement files yourself (faster than materializing)

### 4. Testing Approach

**Build Tests:**
```bash
# After each major change
go build ./cmd/ticketr

# Verify command still works
./ticketr tui --help
```

**Manual Functional Tests:**
```bash
# You'll need workspace setup
./ticketr workspace create test-ws --url https://... --project TEST --username ... --token ...

# Launch TUI
./ticketr tui

# Test:
# - Tab switches focus (borders change color)
# - j/k navigate ticket tree
# - h/l collapse/expand
# - Enter on workspace loads tickets
# - q quits
# - ? shows help (and still works)
```

**Unit Tests (if time permits):**
```bash
# Create ticket_tree_test.go
# Test ticket loading logic
# Mock TicketService
```

### 5. Commit Strategy

**Single commit for Week 12:**
```
feat(tui): Implement Phase 4 Week 12 - Multi-panel layout and ticket integration

Enhance TUI with multi-panel layout showing workspace list and ticket tree
simultaneously. Integrate real ticket data from TicketService.

Architecture:
- Flex layout with workspace list (left) + ticket tree (right)
- Panel focus management with Tab key switching
- Visual focus indicators (green border for active panel)
- TicketService injection into TicketTreeView

Features:
- Multi-panel layout (workspace list + ticket tree side-by-side)
- Tab key switches focus between panels
- Focused panel highlighted with colored border
- Real ticket data from TicketService (no more placeholders)
- Workspace switch triggers automatic ticket reload
- Vim-style navigation: j/k (up/down), h/l (collapse/expand)
- Loading indicators during ticket fetch
- Error display for failed operations
- All Week 11 functionality preserved

Technical Details:
- Modified: app.go, keybindings.go, ticket_tree.go, workspace_list.go
- Dependencies: No new dependencies
- Backward compatible: Arrow keys still work alongside vim keys
- Lines changed: ~200-300 (estimated)

Roadmap: Phase 4 Week 12 complete
Next: Week 13 - Ticket detail editor with validation

Generated with [Claude Code](https://claude.ai/code)
via [Happy](https://happy.engineering)

Co-Authored-By: Claude <noreply@anthropic.com>
Co-Authored-By: Happy <yesreply@happy.engineering>
```

## Critical Patterns to Follow

### 1. Hexagonal Architecture Compliance

**DO:**
```go
// TUI depends on services
type TUIApp struct {
    workspace    *services.WorkspaceService
    ticketService *services.TicketService  // Good
}
```

**DON'T:**
```go
// TUI bypasses services
type TUIApp struct {
    db *sql.DB  // BAD - violates architecture
}
```

### 2. Dependency Injection

**Always inject dependencies via constructor:**
```go
func NewTicketTreeView(ticketService *services.TicketService) *TicketTreeView {
    if ticketService == nil {
        // Validation
    }
    return &TicketTreeView{
        ticketService: ticketService,
    }
}
```

### 3. View Lifecycle

**Always use OnShow/OnHide:**
```go
func (v *TicketTreeView) OnShow() {
    // Refresh data when view becomes active
    v.loadTickets()
}

func (v *TicketTreeView) OnHide() {
    // Cleanup if needed
}
```

### 4. Error Handling

**Don't panic, show errors in UI:**
```go
tickets, err := v.ticketService.List()
if err != nil {
    v.root.ClearChildren()
    errNode := tview.NewTreeNode(fmt.Sprintf("Error: %v", err))
    errNode.SetColor(tcell.ColorRed)
    v.root.AddChild(errNode)
    return
}
```

## Files You'll Likely Modify

### High Probability (90%+)
- `internal/adapters/tui/app.go` - Multi-panel layout, focus management
- `internal/adapters/tui/keybindings.go` - Tab handler
- `internal/adapters/tui/views/ticket_tree.go` - Real data integration
- `internal/adapters/tui/views/workspace_list.go` - Workspace change callback

### Medium Probability (50%)
- `internal/adapters/tui/router.go` - Might need focus-aware routing
- `cmd/ticketr/tui_command.go` - TicketService initialization

### Low Probability (10%)
- `internal/core/services/ticket_service.go` - If List() method missing
- `internal/adapters/tui/views/view.go` - If interface needs extension

### Won't Modify
- Core domain models (no changes needed)
- Database adapters (TUI doesn't touch these)
- Jira adapter (TUI doesn't touch these)

## Common Pitfalls (From Week 11)

### Pitfall 1: Import Path Errors

**Issue:** Wrong module name in imports

**Solution:** Always use `github.com/karolswdev/ticktr` (check go.mod line 1)

### Pitfall 2: Agent Files Don't Persist

**Issue:** Builder agent creates files but they vanish

**Solution:** Manually create files with Write tool based on agent design

### Pitfall 3: Nil Service Panics

**Issue:** Forgot to initialize service before passing to TUI

**Solution:** Always initialize in tui_command.go:
```go
ticketService, err := initTicketService()  // Create this if missing
app, err := tui.NewTUIApp(workspaceService, pathResolver, ticketService)
```

### Pitfall 4: tview Focus Confusion

**Issue:** SetFocus doesn't visually highlight panel

**Solution:** Manually update border colors after SetFocus:
```go
t.app.SetFocus(panel)
panel.SetBorderColor(tcell.ColorGreen)
otherPanel.SetBorderColor(tcell.ColorWhite)
```

## Reference Files

**Must Read Before Starting:**
- `docs/v3-implementation-roadmap.md` (lines 272-369) - Phase 4 full spec
- `docs/PHASE-4-WEEK-11-DEVELOPMENT-PROCESS.md` - Week 11 methodology
- `internal/adapters/tui/app.go` - Current TUI implementation

**Helpful References:**
- `internal/core/services/workspace_service.go` - Service pattern example
- `cmd/ticketr/workspace_commands.go` - Service initialization pattern
- tview documentation: https://github.com/rivo/tview/wiki

## Key Service Methods You'll Need

**WorkspaceService:**
```go
List() ([]*domain.Workspace, error)
Current() (*domain.Workspace, error)
Switch(name string) error
```

**TicketService (verify this exists):**
```go
List() ([]domain.Ticket, error)  // Might need workspace param
Get(key string) (*domain.Ticket, error)
```

**PathResolver:**
```go
DatabasePath() string
DataDir() string
ConfigDir() string
```

## Success Criteria

**You're done with Week 12 when:**
- [ ] `go build ./cmd/ticketr` succeeds
- [ ] `./ticketr tui` shows workspace list + ticket tree side-by-side
- [ ] Tab key switches focus (borders change color)
- [ ] Ticket tree shows real tickets from current workspace
- [ ] Changing workspace reloads tickets automatically
- [ ] j/k/h/l vim keys work for navigation
- [ ] Arrow keys still work (backward compatible)
- [ ] Loading state shows during ticket fetch
- [ ] Errors display clearly if ticket load fails
- [ ] All Week 11 features still work (q, ?, help view)
- [ ] Code committed with descriptive message
- [ ] Todo list all items marked complete

## Your First Actions

**Immediate next steps:**

1. **Create todo list:**
```
Use TodoWrite tool with the 12-item checklist from "Development Workflow" section above
```

2. **Verify TicketService exists:**
```
Read internal/core/services/ticket_service.go
Grep "func.*List" internal/core/services/ticket_service.go
```

3. **If TicketService looks good, start implementation:**
```
Edit internal/adapters/tui/app.go
- Add multi-panel Flex layout
- Add currentFocus tracking
```

4. **If TicketService is missing/incomplete:**
```
Analyze internal/core/services/pull_service.go
- Might have ticket loading logic there
- Decide: extend PullService or create TicketService.List()
```

## Questions to Ask Yourself

Before implementing, answer these:

1. **Does TicketService.List() exist?**
   - If yes: What params does it need? (workspace? filters?)
   - If no: Where is ticket data loaded? (PullService? Repository?)

2. **How should tickets be organized in the tree?**
   - Flat list?
   - Grouped by sprint/epic/status?
   - Hierarchical (epics > stories > subtasks)?

3. **What happens if no workspace is selected?**
   - Show empty tree?
   - Show "No workspace selected" message?
   - Disable ticket panel until workspace chosen?

4. **Should ticket loading be async?**
   - Yes if potentially slow (Jira API call)
   - Can use goroutine + app.QueueUpdateDraw()
   - Consider loading indicator

## Contact/Escalation

If you encounter blockers:

1. **Missing service methods:** Check if functionality exists elsewhere (PullService, repositories)
2. **Architectural questions:** Refer to hexagonal pattern (adapters → services → domain)
3. **tview UI issues:** Check tview wiki and examples on GitHub

## Summary

**Your mission:** Transform the single-view TUI skeleton into a multi-panel interface with real ticket data.

**Your constraints:**
- Follow hexagonal architecture strictly
- Use existing services (don't bypass to DB)
- Maintain all Week 11 functionality
- Keep code clean and testable

**Your tools:**
- TodoWrite (track progress)
- Read/Edit/Write (file operations)
- Bash (build, test, commit)
- Task (optional: delegate to Builder)

**Your success metric:** A working TUI that shows workspaces and tickets side-by-side, with all navigation working.

Good luck! 🚀

---

**Last commit:** `326c357` (Week 11 complete)
**Next commit:** Week 12 multi-panel layout
**Branch:** `feature/v3`
**Roadmap:** Phase 4 (TUI), Week 12 of 16
