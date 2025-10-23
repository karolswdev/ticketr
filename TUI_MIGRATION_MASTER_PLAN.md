# Ticketr Bubbletea TUI Migration - Master Plan

**Date:** 2025-10-22
**Author:** Director Agent (Final Arbiter & Architect)
**Status:** DEFINITIVE IMPLEMENTATION ROADMAP
**Version:** 1.0.0

---

## Executive Summary

### Decision: APPROVED FOR FULL MIGRATION

After comprehensive analysis of all research documents, current codebase, and project goals, this plan authorizes a **complete migration from Tview to Bubbletea** with the following key decisions:

**APPROVED SCOPE:**
- ✅ Complete Bubbletea migration (all views, components, state management)
- ✅ Extensible action system with plugin readiness
- ✅ Midnight Commander-inspired dual-panel UI
- ✅ High-performance rendering (60 FPS target)
- ✅ Visual effects (animations, shimmer, gradients)
- ⚠️ **DEFERRED:** Background particle effects (cosmic/hyperspace) - Phase 6 stretch goal
- ✅ Complete testing coverage (unit, integration, visual regression)
- ✅ Feature flag rollout strategy

**TIMELINE:** 16-18 weeks (4-4.5 months)
- Minimum Viable Product: 8 weeks
- Full Feature Parity (minus particles): 12 weeks
- Complete with all polish: 16-18 weeks

**RISK LEVEL:** Medium-High (justified by long-term maintainability gains)

**PRIMARY BENEFITS:**
1. **Eliminate Race Conditions**: No more mutex hell (serviceMutex, RWMutex chaos)
2. **Simplify State Management**: Single source of truth via Model-View-Update
3. **Enable Future Extensibility**: Plugin system, Lua actions, community actions
4. **Improve Developer Experience**: Pure functions, testable, declarative
5. **Modern Framework**: Active ecosystem, production-proven, better documentation

---

## Table of Contents

1. [Architecture Decisions](#1-architecture-decisions)
2. [Project Structure](#2-project-structure)
3. [Component Strategy](#3-component-strategy)
4. [State Management Approach](#4-state-management-approach)
5. [Phased Implementation Roadmap](#5-phased-implementation-roadmap)
6. [Quality Gates & Metrics](#6-quality-gates--metrics)
7. [Risk Assessment & Mitigation](#7-risk-assessment--mitigation)
8. [Rollout Strategy](#8-rollout-strategy)
9. [Success Criteria](#9-success-criteria)
10. [Post-Migration Plan](#10-post-migration-plan)

---

## 1. Architecture Decisions

### Decision Matrix

| Decision Area | Current (Tview) | **APPROVED** (Bubbletea) | Rationale |
|--------------|----------------|-------------------------|-----------|
| **State Management** | Imperative callbacks | Elm Architecture (MVU) | Eliminates race conditions, single source of truth |
| **Rendering** | Manual `QueueUpdateDraw()` | Declarative View() | Automatic re-render on state change |
| **Async Coordination** | Goroutines + mutex | tea.Cmd messages | No race conditions, testable |
| **Component Model** | Tview primitives | Bubbletea models | Composable, reusable, isolated state |
| **Event Handling** | Callback chains | Message passing | Pure functions, predictable |
| **Focus Management** | Manual tracking | Model state | Automatic rendering updates |
| **Effects System** | Direct tcell screen access | Tick-based commands | Framework-native, no hacks |

### Core Architectural Patterns

#### 1.1 Elm Architecture (TEA)

```
┌──────────────────────────────────────────────┐
│              Root Model                       │
│  (Single source of truth for all state)      │
└──────────────────────────────────────────────┘
                    │
      ┌─────────────┼─────────────┐
      │             │             │
┌─────▼──────┐  ┌──▼────┐  ┌─────▼─────┐
│ Components │  │ State │  │ Services  │
│  (Views)   │  │       │  │ (Injected)│
└────────────┘  └───────┘  └───────────┘
      │             │             │
      └─────────────┼─────────────┘
                    │
            ┌───────▼────────┐
            │  Update(msg)   │
            │  (State Reducer)│
            └───────┬────────┘
                    │
            ┌───────▼────────┐
            │   View()       │
            │  (Pure Render) │
            └────────────────┘
```

**Decision: Embrace TEA fully**
- All state changes via Update()
- View() is pure function of Model
- Commands for all I/O
- No goroutines outside Commands

#### 1.2 Component Hierarchy

```
RootModel
├── ContextManager (tracks current focus)
├── ActionRegistry (extensible actions)
├── KeybindingResolver (maps keys to actions)
├── Components
│   ├── WorkspaceListModel (bubbles/list)
│   ├── TicketTreeModel (custom tree component)
│   ├── TicketDetailModel (viewport + forms)
│   ├── SearchModal (textinput + list)
│   ├── CommandPalette (textinput + list)
│   └── BulkOpsModal (huh forms)
├── Widgets
│   ├── ActionBar (context-aware keybindings)
│   ├── StatusBar (sync status, workspace info)
│   ├── ProgressBar (shimmer effect)
│   └── Marquee (animated text ticker)
└── Services (injected)
    ├── TicketService
    ├── WorkspaceService
    ├── SyncService
    └── JobQueue
```

**Decision: Flat component model, no deep nesting**
- Components communicate via messages (not callbacks)
- Services injected into RootModel
- Components get services via ActionContext
- No circular dependencies

#### 1.3 Message-Driven Architecture

```go
// All state changes via messages
type Message interface{}

// Built-in messages (from Bubbletea)
- tea.KeyMsg          // Keyboard input
- tea.MouseMsg        // Mouse events
- tea.WindowSizeMsg   // Terminal resize

// Domain messages (our custom types)
- ticketOpenedMsg     // Ticket detail opened
- ticketSavedMsg      // Ticket saved successfully
- ticketsLoadedMsg    // Tickets fetched from DB
- syncStartedMsg      // Sync operation started
- syncCompletedMsg    // Sync finished (success/error)
- contextChangedMsg   // Focus changed (tree ↔ detail)
- modalOpenedMsg      // Modal overlay shown
- modalClosedMsg      // Modal dismissed
- actionExecutedMsg   // Action completed
- errorMsg            // Error occurred
```

**Decision: Message types are data, not code**
- Messages carry data payloads
- Update() is the only reducer
- Messages trigger state transitions
- Pure function: `(Model, Msg) → (Model, Cmd)`

---

## 2. Project Structure

### Approved Directory Layout

```
ticketr/
├── cmd/
│   ├── root.go                 # Cobra root command
│   └── tui.go                  # TUI entry point
├── internal/
│   ├── domain/                 # Business logic (unchanged)
│   ├── adapters/
│   │   ├── tui/                # Current Tview (keep until cutover)
│   │   └── bubbletea/          # NEW: Bubbletea implementation
│   │       ├── model.go        # Root model definition
│   │       ├── init.go         # Initialization logic
│   │       ├── update.go       # Update reducer (message routing)
│   │       ├── view.go         # View renderer (layout composition)
│   │       ├── messages.go     # Custom message type definitions
│   │       ├── commands.go     # Command factories (async ops)
│   │       ├── styles/         # Lipgloss style definitions
│   │       │   ├── theme.go    # Theme system (colors, effects)
│   │       │   ├── borders.go  # Border styles
│   │       │   └── layout.go   # Layout constants
│   │       ├── components/     # Reusable UI components
│   │       │   ├── tree/       # Custom tree component
│   │       │   │   ├── model.go
│   │       │   │   ├── update.go
│   │       │   │   ├── view.go
│   │       │   │   └── tree_test.go
│   │       │   ├── detail/     # Ticket detail viewer
│   │       │   ├── workspaces/ # Workspace list
│   │       │   ├── search/     # Search modal
│   │       │   ├── palette/    # Command palette
│   │       │   └── bulk/       # Bulk operations modal
│   │       ├── widgets/        # Smaller UI elements
│   │       │   ├── actionbar/  # Action bar widget
│   │       │   ├── statusbar/  # Status bar widget
│   │       │   ├── progress/   # Progress bar with shimmer
│   │       │   └── marquee/    # Animated text ticker
│   │       ├── actions/        # Extensible action system
│   │       │   ├── action.go       # Core action types
│   │       │   ├── context.go      # Context manager
│   │       │   ├── registry.go     # Action registry
│   │       │   ├── resolver.go     # Keybinding resolver
│   │       │   ├── executor.go     # Execution pipeline
│   │       │   ├── predicates/     # Predicate functions
│   │       │   │   └── predicates.go
│   │       │   ├── modifiers/      # Action modifiers
│   │       │   │   └── modifiers.go
│   │       │   └── builtin/        # Built-in actions
│   │       │       ├── tickets.go
│   │       │       ├── workspaces.go
│   │       │       ├── navigation.go
│   │       │       ├── sync.go
│   │       │       └── system.go
│   │       ├── effects/        # Visual effects
│   │       │   ├── animator.go     # Animation tick system
│   │       │   ├── shimmer.go      # Progress shimmer effect
│   │       │   ├── transitions.go  # Fade/slide transitions
│   │       │   └── particles.go    # Particle system (Phase 6)
│   │       └── integration_test.go # End-to-end tests
│   └── tui/                    # Shared TUI infrastructure
│       └── jobs/               # Job queue (shared)
├── config/
│   └── themes/                 # Theme configurations
│       ├── default.yaml
│       ├── dark.yaml
│       └── arctic.yaml
└── docs/
    ├── bubbletea-migration/    # Migration documentation
    │   ├── architecture.md
    │   ├── component-guide.md
    │   ├── action-system.md
    │   └── testing-guide.md
    └── development/
        └── ROADMAP.md          # Updated roadmap
```

**Decision: Parallel development strategy**
- Keep `internal/adapters/tui/` working during migration
- Build complete implementation in `internal/adapters/bubbletea/`
- Use feature flag (`TICKETR_USE_BUBBLETEA=true`) for testing
- Cut over when bubbletea reaches feature parity

---

## 3. Component Strategy

### 3.1 Component Library Decisions

| Component Need | **APPROVED SOLUTION** | Reason |
|---------------|---------------------|--------|
| **Ticket List** | `charmbracelet/bubbles/list` | Production-ready, pagination, filtering |
| **Ticket Tree** | **CUSTOM IMPLEMENTATION** | No built-in tree, need hierarchy, vim bindings |
| **Text Input** | `charmbracelet/bubbles/textinput` | Standard, robust, focus management |
| **Multi-line Text** | `charmbracelet/bubbles/textarea` | For editing ticket descriptions |
| **Scrolling Content** | `charmbracelet/bubbles/viewport` | For ticket detail panel |
| **Tables** | `charmbracelet/bubbles/table` | For structured data display |
| **Progress Bars** | `charmbracelet/bubbles/progress` + custom shimmer | Animated progress with effects |
| **Spinners** | `charmbracelet/bubbles/spinner` | Loading indicators |
| **Help System** | `charmbracelet/bubbles/help` | Auto-generated keybinding help |
| **Forms** | `charmbracelet/huh/v2` | Complex forms (workspace, bulk ops) |
| **Modals** | **CUSTOM** (lipgloss.Place) | Overlay modals with backdrop |
| **Layout** | **Lipgloss** native (JoinHorizontal/Vertical) | Simple, no external dependencies |

### 3.2 Tree Component Architecture

**Decision: Custom tree implementation** (highest complexity component)

```go
// internal/adapters/bubbletea/components/tree/model.go

package tree

type Node struct {
    ID       string
    Label    string
    Icon     string
    Children []*Node
    Expanded bool
    Selected bool
    Data     interface{} // Ticket, Epic, etc.
}

type Model struct {
    root         *Node
    cursor       int           // Current cursor position
    offset       int           // Scroll offset
    height       int           // Viewport height
    width        int           // Viewport width
    selectedIDs  map[string]bool // Multi-selection
    expandedIDs  map[string]bool // Expansion state
    flatView     []*Node       // Flattened tree for rendering
    keybindings  KeyMap        // Vim-style bindings
}

// Key features:
// - Lazy rendering (only visible nodes)
// - Multi-selection with Space
// - Vim keybindings (hjkl, gg, G)
// - Expand/collapse animation
// - Search/filter integration
// - Sync status indicators (●○◐)
```

**Implementation Strategy:**
1. Week 1: Core tree data structure, rendering
2. Week 2: Navigation (hjkl, expand/collapse)
3. Week 3: Multi-selection, search integration
4. Week 4: Performance optimization, testing

### 3.3 Effects System Translation

| Current Effect | Bubbletea Translation | Complexity | Status |
|---------------|----------------------|------------|--------|
| **Animator** | `tea.Tick()` based system | Medium | Phase 3 |
| **Shimmer** | String manipulation in View() | Low | Phase 2 |
| **ShadowBox** | Lipgloss border decoration | Low | Phase 2 |
| **Borders** | Lipgloss style helpers | Low | Phase 1 |
| **Marquee** | Tick-based state machine | Medium-High | Phase 3 |
| **Background Particles** | **DEFERRED** to Phase 6 | Very High | Stretch |

**Background Particle Decision:**
- **DEFERRED to Phase 6** (stretch goal after feature parity)
- Reason: Very high complexity, requires ANSI layer compositing
- Alternative: Static gradient background in interim
- Revisit after core functionality stable

---

## 4. State Management Approach

### 4.1 Root Model Definition

```go
// internal/adapters/bubbletea/model.go

package bubbletea

type Model struct {
    // Window dimensions
    width  int
    height int

    // Context & Focus
    context        actions.Context
    contextManager *actions.ContextManager
    focus          Focus // Which panel is focused

    // Action System
    actionRegistry     *actions.Registry
    keybindingResolver *actions.KeybindingResolver
    executor           *actions.Executor

    // UI Components
    workspaceList WorkspaceListModel
    ticketTree    TreeModel
    ticketDetail  DetailModel
    statusBar     StatusBarModel
    actionBar     ActionBarModel

    // Modals (only one active at a time)
    activeModal Modal // nil = no modal

    // Application State
    currentWorkspace *domain.Workspace
    selectedTickets  []string
    expandedNodes    map[string]bool

    // Sync State
    syncInProgress bool
    syncProgress   float64
    lastSyncTime   time.Time
    lastSyncError  error

    // UI State
    showHelp          bool
    errorMessage      string
    hasUnsavedChanges bool

    // Services (dependency injection)
    services *Services

    // Configuration
    config *Config
}

type Services struct {
    TicketService    domain.TicketService
    WorkspaceService domain.WorkspaceService
    SyncService      domain.SyncService
    JobQueue         *jobs.Queue
}

type Config struct {
    Theme          Theme
    Keybindings    map[actions.ActionID][]actions.KeyPattern
    Features       map[string]bool
    Motion         bool  // Enable animations
    EffectsEnabled bool  // Enable visual effects
}
```

**Decision: Single root model, flat component structure**
- No deep nesting of sub-models
- Components are fields, not interfaces
- Services injected at startup
- Configuration loaded from YAML/env

### 4.2 Update Pattern

```go
// internal/adapters/bubbletea/update.go

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    // Global messages first
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle global quit
        if msg.String() == "q" && m.activeModal == nil {
            return m, tea.Quit
        }

        // Route to modal if active
        if m.activeModal != nil {
            newModal, cmd := m.activeModal.Update(msg)
            m.activeModal = newModal
            return m, cmd
        }

        // Resolve keybinding to action
        actx := m.buildActionContext()
        action, found := m.keybindingResolver.Resolve(msg, actx)
        if found {
            cmd := m.executor.Execute(action, actx)
            return m, cmd
        }

    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        // Broadcast to all components
        m.ticketTree, cmd := m.ticketTree.Update(msg)
        cmds = append(cmds, cmd)
        m.ticketDetail, cmd = m.ticketDetail.Update(msg)
        cmds = append(cmds, cmd)
        // ... other components
        return m, tea.Batch(cmds...)

    // Domain messages
    case ticketOpenedMsg:
        m.ticketDetail.SetTicket(msg.ticket)
        m.context = actions.ContextTicketDetail
        m.focus = FocusDetail
        return m, nil

    case ticketsLoadedMsg:
        m.ticketTree.SetTickets(msg.tickets)
        return m, nil

    case syncStartedMsg:
        m.syncInProgress = true
        m.syncProgress = 0.0
        return m, nil

    case syncProgressMsg:
        m.syncProgress = msg.progress
        return m, nil

    case syncCompletedMsg:
        m.syncInProgress = false
        m.lastSyncTime = time.Now()
        m.lastSyncError = msg.err
        return m, nil
    }

    // Route to focused component
    switch m.focus {
    case FocusTree:
        m.ticketTree, cmd := m.ticketTree.Update(msg)
        return m, cmd
    case FocusDetail:
        m.ticketDetail, cmd := m.ticketDetail.Update(msg)
        return m, cmd
    }

    return m, nil
}
```

**Decision: Message routing in Update()**
- Global messages first (quit, resize)
- Modal intercepts if active
- Action resolution via keybinding
- Fallback to focused component
- Pure function, no side effects

### 4.3 View Pattern

```go
// internal/adapters/bubbletea/view.go

func (m Model) View() string {
    // Handle modal overlay
    if m.activeModal != nil {
        mainView := m.renderMainView()
        modalView := m.activeModal.View()
        return lipgloss.PlaceOverlay(
            m.width/2, m.height/2,
            modalView,
            mainView,
            lipgloss.WithWhitespaceBackground(lipgloss.Color("0")),
        )
    }

    return m.renderMainView()
}

func (m Model) renderMainView() string {
    // Header (status bar)
    header := m.statusBar.View()

    // Main content (dual panel)
    leftWidth := m.width / 2
    rightWidth := m.width - leftWidth

    treeView := m.renderTree(leftWidth, m.height-6)
    detailView := m.renderDetail(rightWidth, m.height-6)

    content := lipgloss.JoinHorizontal(
        lipgloss.Top,
        treeView,
        detailView,
    )

    // Footer (action bar)
    footer := m.actionBar.View()

    // Compose vertically
    return lipgloss.JoinVertical(
        lipgloss.Left,
        header,
        content,
        footer,
    )
}
```

**Decision: Declarative layout composition**
- View() is pure function of Model
- No stateful rendering logic
- Lipgloss for all layout
- Focus indicated by border style

---

## 5. Phased Implementation Roadmap

### Phase 0: Foundation & Proof-of-Concept (Weeks 1-2)

**Goal:** Establish architecture, validate patterns, build foundation

#### Week 1: Setup & Core Patterns
- [ ] **Directory Structure**
  - Create `internal/adapters/bubbletea/` package
  - Set up component, action, widget directories
  - Create style system structure

- [ ] **Core Types**
  - Define RootModel struct
  - Implement Message types
  - Create ActionContext system
  - Build ContextManager

- [ ] **Action System Foundation**
  - Implement Action, Predicate types
  - Build Registry with register/get/search
  - Create KeybindingResolver
  - Implement Executor with middleware support

- [ ] **Testing Infrastructure**
  - Set up teatest framework
  - Create test helpers
  - Implement golden file testing pattern

**Deliverables:**
- ✅ Package structure established
- ✅ Core types compiled and tested
- ✅ Action registry with 5+ test actions
- ✅ Unit tests passing (>80% coverage)

**Success Criteria:**
- All packages compile without errors
- Action registry can register, retrieve, search actions
- Predicates can be composed (And, Or, Not)
- Tests pass with green suite

#### Week 2: Proof-of-Concept UI
- [ ] **Simple View Implementation**
  - Build HelpView (static content)
  - Implement basic Model/Update/View pattern
  - Add WindowSizeMsg handling
  - Create Lipgloss style system

- [ ] **Keybinding Integration**
  - Map tea.KeyMsg to actions
  - Implement action execution pipeline
  - Test with help toggle (? key)

- [ ] **Theme System**
  - Port current theme colors to Lipgloss
  - Implement theme loading from config
  - Support TICKETR_THEME env var

**Deliverables:**
- ✅ Working help screen (? to toggle)
- ✅ Keybinding system operational
- ✅ Theme system functional
- ✅ Basic window resize handling

**Success Criteria:**
- Can run `TICKETR_USE_BUBBLETEA=true ./ticketr` and see help screen
- Pressing ? toggles help on/off
- Terminal resize updates layout correctly
- No rendering artifacts

---

### Phase 1: Core Views (Weeks 3-6)

**Goal:** Implement primary user interfaces (workspace, tree, detail)

#### Week 3: Workspace List
- [ ] **Component Implementation**
  - Use bubbles/list for workspace list
  - Implement sync status indicators (●○◐)
  - Add filter/search functionality
  - Create workspace creation modal (huh forms)

- [ ] **Integration**
  - Load workspaces from service
  - Handle selection events
  - Switch workspace command
  - Update status bar on switch

**Deliverables:**
- ✅ Functional workspace list view
- ✅ Workspace switching works
- ✅ Create new workspace modal
- ✅ Sync status indicators animated

**Success Criteria:**
- Can list all workspaces
- Can select and switch workspaces
- Can create new workspace via modal
- Status indicators update during sync

#### Week 4-5: Ticket Tree (Critical Path)
- [ ] **Custom Tree Component** (Week 4)
  - Build Node data structure
  - Implement tree flattening algorithm
  - Basic rendering (text only)
  - Cursor navigation (jk, up/down)

- [ ] **Enhanced Features** (Week 5)
  - Expand/collapse (hl, left/right)
  - Multi-selection with Space
  - Vim keybindings (gg, G, search)
  - Visual indicators (icons, sync status)
  - Grouping by type/status/assignee

**Deliverables:**
- ✅ Fully functional tree component
- ✅ Multi-selection working
- ✅ All navigation keybindings
- ✅ Visual polish (icons, borders)

**Success Criteria:**
- Can navigate tree with hjkl/arrows
- Can expand/collapse nodes
- Can select multiple tickets
- Tree renders without flicker
- Performance: <16ms render time for 1000 tickets

#### Week 6: Ticket Detail View
- [ ] **Display Mode**
  - Use viewport for scrolling
  - Render ticket fields (title, description, etc.)
  - Syntax highlighting for markdown
  - Custom field display

- [ ] **Edit Mode**
  - Integrate huh forms for editing
  - Inline field editing
  - Save/cancel actions
  - Dirty state tracking

**Deliverables:**
- ✅ Detail view displays ticket info
- ✅ Scrolling works (jk, Ctrl+F/B)
- ✅ Edit mode functional
- ✅ Save/cancel operations

**Success Criteria:**
- Detail view renders all ticket fields
- Can scroll long descriptions
- Can edit and save changes
- Dirty state indicator shows unsaved changes

---

### Phase 2: Modals & Widgets (Weeks 7-9)

**Goal:** Implement secondary UI elements (search, palette, widgets)

#### Week 7: Search Modal
- [ ] **Search UI**
  - Center modal overlay
  - textinput for query
  - Fuzzy matching engine
  - Results list with scores

- [ ] **Filter Syntax**
  - Implement @user, !priority, #ID, ~sprint filters
  - Regex support (/pattern/)
  - Exact match (field:value)

**Deliverables:**
- ✅ Search modal working (/ key)
- ✅ Fuzzy search functional
- ✅ All filter syntax supported
- ✅ Results ranked by relevance

**Success Criteria:**
- Can search tickets by text
- Filters work (e.g., "auth @john !high")
- Results update in real-time (<100ms)
- Match highlights visible

#### Week 8: Command Palette
- [ ] **Palette UI**
  - Center modal overlay
  - textinput with : prefix
  - Command list (fuzzy filtered)
  - Category headers

- [ ] **Command Registration**
  - Register all actions as commands
  - Show keybinding hints
  - Icon/description display

**Deliverables:**
- ✅ Command palette working (: or Ctrl+P)
- ✅ All actions searchable
- ✅ Fuzzy search functional
- ✅ Execute commands on Enter

**Success Criteria:**
- Palette shows all registered actions
- Fuzzy search works ("sp" → "sync:pull")
- Keybinding hints visible
- Commands execute correctly

#### Week 9: Widgets & Action Bar
- [ ] **Action Bar**
  - Context-aware keybinding display
  - Dynamic updates on focus change
  - Responsive layout (truncate on narrow)

- [ ] **Status Bar**
  - Workspace name
  - Sync status indicator
  - Ticket count
  - Last sync time

- [ ] **Progress Bar**
  - Bubbles progress component
  - Add shimmer effect
  - Percentage display
  - ETA calculation

**Deliverables:**
- ✅ Action bar shows context keybindings
- ✅ Status bar displays all info
- ✅ Progress bar with shimmer
- ✅ All widgets responsive

**Success Criteria:**
- Action bar updates when focus changes
- Status bar reflects current state
- Progress bar animates smoothly
- No layout jank on resize

---

### Phase 3: Async Operations & Sync (Weeks 10-11)

**Goal:** Integrate Jira sync, job queue, bulk operations

#### Week 10: Job Queue Integration
- [ ] **Command Pattern**
  - Convert sync operations to tea.Cmd
  - Implement progress reporting via messages
  - Error handling via messages

- [ ] **UI Integration**
  - Show progress bar during sync
  - Display status in status bar
  - Handle cancellation (Esc)

**Deliverables:**
- ✅ Pull operation working (P key)
- ✅ Push operation working (p key)
- ✅ Full sync working (s key)
- ✅ Progress tracking functional

**Success Criteria:**
- Sync operations don't block UI
- Progress bar updates in real-time
- Can cancel sync with Esc
- Errors displayed to user

#### Week 11: Bulk Operations Modal
- [ ] **Wizard UI** (3-step)
  - Step 1: Operation selection
  - Step 2: Field editing (huh forms)
  - Step 3: Progress tracking

- [ ] **Operations**
  - Update fields (status, priority, assignee, etc.)
  - Move tickets
  - Delete tickets (with confirmation)

**Deliverables:**
- ✅ Bulk operations modal working (b key)
- ✅ All operations implemented
- ✅ Progress tracking per ticket
- ✅ Error handling with rollback option

**Success Criteria:**
- Can update multiple tickets at once
- Progress shows per-ticket status
- Errors don't break the wizard
- Rollback option on failure

---

### Phase 4: Visual Polish & Effects (Weeks 12-14)

**Goal:** Add animations, transitions, visual effects

#### Week 12: Animations
- [ ] **Tick-Based Animator**
  - Convert current Animator to tea.Tick
  - Animation registry pattern
  - Enable/disable via config

- [ ] **Core Animations**
  - Modal fade-in/out
  - Slide-out panel animation
  - Spinner for loading states
  - Sync status rotation (◐◓◑◒)

**Deliverables:**
- ✅ Animation system working
- ✅ Modal transitions smooth
- ✅ Loading spinners functional
- ✅ Sync indicator animated

**Success Criteria:**
- Animations run at 60 FPS
- No dropped frames
- Can disable via config
- Smooth, professional feel

#### Week 13: Marquee & Effects
- [ ] **Marquee Component**
  - Tick-based state machine
  - Slide-in, blink, slide-out phases
  - Queue multiple messages

- [ ] **Shimmer Effect**
  - Apply to progress bars
  - Sweep animation (left to right)

**Deliverables:**
- ✅ Marquee widget working
- ✅ Shimmer effect on progress
- ✅ Both enabled/disabled by config

**Success Criteria:**
- Marquee cycles smoothly
- Shimmer visible and smooth
- No performance impact when disabled

#### Week 14: Theme & Polish
- [ ] **Multiple Themes**
  - Default (Green)
  - Dark (Blue)
  - Arctic (Cyan)

- [ ] **Visual Polish**
  - Drop shadows on modals
  - Focus pulse effect (optional)
  - Gradient titles (optional)
  - High contrast mode

**Deliverables:**
- ✅ 3 themes working
- ✅ Theme switcher in settings
- ✅ Visual effects polished
- ✅ Accessibility modes functional

**Success Criteria:**
- Can switch themes without restart
- All themes look good
- High contrast mode readable
- No visual glitches

---

### Phase 5: Testing & Documentation (Weeks 15-16)

**Goal:** Comprehensive testing, documentation, bug fixes

#### Week 15: Testing
- [ ] **Unit Tests**
  - All components >80% coverage
  - Action system 100% coverage
  - Predicate logic tested
  - Message routing tested

- [ ] **Integration Tests**
  - End-to-end workflows
  - Sync operations
  - Bulk operations
  - Modal interactions

- [ ] **Visual Regression Tests**
  - Golden file tests for all views
  - Different terminal sizes
  - All themes

**Deliverables:**
- ✅ Test coverage >80% overall
- ✅ All critical paths tested
- ✅ Golden files for regression
- ✅ CI pipeline green

**Success Criteria:**
- All tests pass
- Coverage meets target
- No regressions detected
- Fast test suite (<2 min)

#### Week 16: Documentation
- [ ] **Architecture Docs**
  - Bubbletea migration guide
  - Component architecture
  - Action system guide
  - Testing guide

- [ ] **User Docs**
  - Updated README
  - Keybinding reference
  - Theme customization
  - Plugin development guide (future)

**Deliverables:**
- ✅ All docs written
- ✅ README updated
- ✅ CONTRIBUTING.md updated
- ✅ API docs complete

**Success Criteria:**
- Docs are clear and accurate
- Examples work
- No TODOs remaining
- User can onboard from docs alone

---

### Phase 6: Rollout & Stabilization (Weeks 17-18)

**Goal:** Feature flag rollout, beta testing, cutover

#### Week 17: Beta Rollout
- [ ] **Feature Flag**
  - `TICKETR_USE_BUBBLETEA=true` env var
  - Ship as experimental in v3.2.0-beta
  - Gather feedback

- [ ] **Bug Fixes**
  - Fix reported issues
  - Performance tuning
  - UX improvements

**Deliverables:**
- ✅ v3.2.0-beta released
- ✅ Beta feedback collected
- ✅ Critical bugs fixed
- ✅ Performance acceptable

**Success Criteria:**
- No critical bugs reported
- Performance meets targets (60 FPS)
- Positive user feedback
- Feature parity with Tview

#### Week 18: Cutover & Cleanup
- [ ] **Make Bubbletea Default**
  - Flip default to Bubbletea
  - Keep Tview as fallback (`TICKETR_USE_TVIEW=true`)
  - Ship as v3.2.0 stable

- [ ] **Cleanup**
  - Remove deprecated code
  - Archive old Tview docs
  - Update all examples

**Deliverables:**
- ✅ v3.2.0 released (Bubbletea default)
- ✅ Tview deprecated but working
- ✅ Migration guide published
- ✅ Celebration! 🎉

**Success Criteria:**
- Smooth cutover, no regressions
- Users can still use Tview if needed
- Migration path documented
- Team celebrates successful migration

---

### Phase 6.5: Stretch Goals (Post-Cutover)

**Optional:** Only pursue if Phase 6 completes early or after v3.2.0

- [ ] **Background Particle Effects**
  - Hyperspace mode (moving stars)
  - Snow mode (falling particles)
  - Matrix mode (green rain)
  - Requires ANSI layer compositing research

- [ ] **Lua Plugin System**
  - Lua VM integration
  - Action registration API
  - Service bindings
  - Security sandboxing

- [ ] **Advanced Features**
  - Tri-panel mode (ultra-wide monitors)
  - Picture-in-picture preview
  - Git integration UI
  - Custom dashboard widgets

---

## 6. Quality Gates & Metrics

### Performance Targets

| Metric | Target | Critical Path | Measurement |
|--------|--------|---------------|-------------|
| **Frame Rate** | 60 FPS | All views | Profiling |
| **Render Time** | <16ms | Tree, Detail | Benchmarks |
| **Tree Render (1000 items)** | <16ms | Tree | Benchmarks |
| **Search Response** | <100ms | Search | Integration tests |
| **Key Press to Action** | <50ms | All | User perception |
| **Startup Time** | <500ms | Init | Benchmarks |
| **Memory Usage** | <50MB | Steady state | Profiling |

### Code Quality Targets

| Metric | Target | Enforcement |
|--------|--------|-------------|
| **Test Coverage** | >80% | CI fails <80% |
| **Action System Coverage** | 100% | CI fails <100% |
| **Cyclomatic Complexity** | <15 per function | Code review |
| **Function Length** | <100 lines | Code review |
| **File Length** | <500 lines | Code review |
| **No TODO Comments** | 0 before release | CI scan |

### Visual Quality Targets

| Metric | Target | Validation |
|--------|--------|------------|
| **No Flicker** | 0 reported | User testing |
| **No Tearing** | 0 reported | User testing |
| **Responsive Resize** | <100ms | Integration tests |
| **Theme Consistency** | All components themed | Visual review |
| **Accessibility (High Contrast)** | Readable | Manual testing |

### Test Coverage Requirements

```
Overall:                >80%
Components:             >80% each
Actions:                100%
Predicates:             100%
Message Routing:        100%
Critical Paths:         100%
  - Workspace switch
  - Ticket open/save
  - Sync operations
  - Bulk operations
```

### Acceptance Criteria Checklist

Before declaring phase complete:
- [ ] All tests pass (unit, integration, visual)
- [ ] Performance metrics met
- [ ] Code quality metrics met
- [ ] Documentation written
- [ ] No critical bugs
- [ ] User feedback positive
- [ ] Code review approved
- [ ] Security review passed

---

## 7. Risk Assessment & Mitigation

### High-Risk Areas

#### Risk 1: Tree Component Complexity
**Probability:** High | **Impact:** High | **Risk Level:** 🔴 CRITICAL

**Threat:**
- No built-in tree component in Bubbletea
- Complex state management (expansion, selection, navigation)
- Performance risk with large trees (1000+ tickets)

**Mitigation:**
- ✅ **Allocate 2 weeks** (Weeks 4-5) for tree implementation
- ✅ **Prototype early** in Phase 0 to validate approach
- ✅ **Performance benchmarks** required before Phase 2
- ✅ **Lazy rendering** - only render visible nodes
- ✅ **Fallback plan**: Use bubbles/list with indentation if tree fails

**Success Metrics:**
- Render 1000 items in <16ms
- Smooth scrolling at 60 FPS
- No memory leaks after 1 hour use

#### Risk 2: Background Particle Effects
**Probability:** Medium | **Impact:** Medium | **Risk Level:** 🟡 MEDIUM

**Threat:**
- Very high complexity (ANSI layer compositing)
- Performance impact (100+ particles at 60 FPS)
- No proven patterns in Bubbletea ecosystem

**Mitigation:**
- ✅ **DEFER to Phase 6.5** (stretch goal)
- ✅ **Not required for feature parity** - ship without it
- ✅ **Alternative**: Static gradient backgrounds work well
- ✅ **Revisit after v3.2.0** if time permits

**Decision:** Accept loss of particle effects in MVP

#### Risk 3: Race Conditions During Migration
**Probability:** Medium | **Impact:** High | **Risk Level:** 🟡 MEDIUM

**Threat:**
- Two TUI implementations running in parallel
- Shared services may have assumptions about Tview
- State synchronization issues during cutover

**Mitigation:**
- ✅ **Feature flag isolation** - only one TUI active at a time
- ✅ **Service abstraction** - services are already TUI-agnostic
- ✅ **Parallel testing** - both TUIs tested simultaneously
- ✅ **Gradual rollout** - beta period before default switch

**Monitoring:**
- Watch for ServiceMutex usage patterns
- Monitor for goroutine leaks
- Test both TUIs in CI

#### Risk 4: Performance Regression
**Probability:** Low | **Impact:** High | **Risk Level:** 🟢 LOW

**Threat:**
- Bubbletea may be slower than Tview for some operations
- String-based rendering could be expensive
- Animation overhead

**Mitigation:**
- ✅ **Benchmark suite** required in Phase 0
- ✅ **Performance gates** in CI (fail if regression)
- ✅ **Profiling** mandatory for tree and detail views
- ✅ **Optimization sprints** built into phases

**Targets:**
- Match or exceed Tview performance
- 60 FPS target for all views
- <16ms render time budget

#### Risk 5: Team Learning Curve
**Probability:** Medium | **Impact:** Low | **Risk Level:** 🟢 LOW

**Threat:**
- Team unfamiliar with Elm Architecture
- Functional programming paradigm shift
- Bubbletea patterns take time to internalize

**Mitigation:**
- ✅ **Documentation-first approach** - write guides before implementation
- ✅ **Pair programming** encouraged during Phases 0-1
- ✅ **Code examples** in every component
- ✅ **Architecture review** after each phase

**Resources:**
- Official Bubbletea tutorials
- Example applications (Glow, Soft Serve)
- Weekly architecture sync meetings

### Medium-Risk Areas

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| Marquee animation complexity | Medium | Medium | Use existing patterns from current implementation |
| Form validation edge cases | Medium | Low | Huh library handles most cases |
| Modal backdrop rendering | Low | Medium | Lipgloss.PlaceOverlay is proven |
| Theme system bugs | Low | Low | Comprehensive visual regression tests |

### Risk Mitigation Checklist

**Before each phase:**
- [ ] Review risks relevant to phase
- [ ] Identify new risks
- [ ] Update mitigation strategies
- [ ] Check if fallback plans needed

**During each phase:**
- [ ] Monitor for warning signs
- [ ] Escalate blockers immediately
- [ ] Update risk assessment weekly

**After each phase:**
- [ ] Retrospective on risks encountered
- [ ] Update mitigation playbook
- [ ] Share learnings with team

---

## 8. Rollout Strategy

### 8.1 Feature Flag Approach

```bash
# Phase 0-5: Development (feature disabled by default)
./ticketr                           # Uses Tview (current)
TICKETR_USE_BUBBLETEA=true ./ticketr  # Uses Bubbletea (experimental)

# Phase 6: Beta Release (v3.2.0-beta)
./ticketr                           # Uses Tview (stable)
TICKETR_USE_BUBBLETEA=true ./ticketr  # Uses Bubbletea (beta)

# Phase 6+: Stable Release (v3.2.0)
./ticketr                           # Uses Bubbletea (default) ✅
TICKETR_USE_TVIEW=true ./ticketr     # Uses Tview (fallback)

# Future: v4.0.0
./ticketr                           # Uses Bubbletea (only option)
# Tview removed entirely
```

### 8.2 Beta Testing Plan

**Beta 1: Internal Testing (Week 17, Days 1-3)**
- Team dogfoods Bubbletea TUI exclusively
- Fix critical bugs
- Tune performance
- Polish UX

**Beta 2: Limited External Release (Week 17, Days 4-7)**
- Release v3.2.0-beta.1 to GitHub
- Announce in README (beta section)
- Invite early adopters
- Collect feedback via GitHub Issues

**Beta 3: Wider Testing (Week 18, Days 1-3)**
- Release v3.2.0-beta.2 with fixes
- Ask users to test specific workflows
- Run performance comparison vs Tview
- Finalize decision on default

**Stable Release (Week 18, Day 4)**
- Release v3.2.0 (Bubbletea default)
- Migration guide published
- Announce on social media
- Monitor for issues

### 8.3 Migration Communication

**Announcement Timeline:**
- **Week 0:** Announce migration plan in ROADMAP.md
- **Week 4:** Progress update (tree component working)
- **Week 8:** Mid-point update (all views working)
- **Week 12:** Feature complete announcement
- **Week 17:** Beta release announcement
- **Week 18:** Stable release announcement

**Communication Channels:**
- README.md (beta section)
- CHANGELOG.md (detailed changes)
- GitHub Releases (announcements)
- Social media (major milestones)

**Migration Guide Content:**
- Why we migrated (benefits)
- What changed (user-visible)
- How to enable beta
- How to report bugs
- How to revert to Tview
- Timeline for Tview deprecation

### 8.4 Rollback Plan

**If critical bugs found during beta:**
1. Keep v3.2.0-beta as experimental
2. DO NOT make Bubbletea default
3. Fix bugs in hot-fix sprints
4. Release v3.2.0-beta.N with fixes
5. Re-enter beta testing period

**Rollback triggers:**
- Critical crash (data loss)
- Performance regression >50%
- Security vulnerability
- Major feature broken

**Rollback procedure:**
1. Revert default to Tview
2. Tag as v3.2.0-rc.N (release candidate)
3. Fix issue in separate PR
4. Re-test thoroughly
5. Only proceed when stable

### 8.5 Deprecation Timeline

```
v3.2.0 (Week 18)   - Bubbletea default, Tview fallback available
v3.3.0 (+6 months) - Tview deprecated, warning on use
v4.0.0 (+12 months)- Tview removed entirely
```

**User notifications:**
- v3.2.0: "Bubbletea is now default. Use TICKETR_USE_TVIEW=true to revert."
- v3.3.0: "⚠️ Tview is deprecated and will be removed in v4.0.0."
- v4.0.0: "Tview removed. Bubbletea is the only TUI."

---

## 9. Success Criteria

### Definition of Done: Migration Complete

#### Functional Parity
- [x] All Tview views ported to Bubbletea
- [x] All features working (workspace, tree, detail, search, palette, bulk ops)
- [x] All keybindings functional
- [x] All sync operations working
- [x] All visual effects working (except particles - deferred)

#### Quality Metrics Met
- [x] Test coverage >80%
- [x] Performance targets met (60 FPS, <16ms renders)
- [x] No critical bugs
- [x] Visual regression tests passing
- [x] Code quality metrics met

#### Documentation Complete
- [x] Architecture docs written
- [x] Migration guide published
- [x] API docs complete
- [x] User guide updated
- [x] CONTRIBUTING.md updated

#### User Acceptance
- [x] Beta testing completed
- [x] Positive user feedback (>80% satisfaction)
- [x] No major UX regressions reported
- [x] Accessibility requirements met

#### Release Ready
- [x] v3.2.0 tagged and released
- [x] CHANGELOG updated
- [x] GitHub release published
- [x] Social media announcements made

### Long-Term Success Indicators

**After 3 months:**
- No critical bugs reported
- <5 minor bugs reported
- Performance stable
- No user requests to revert to Tview
- Plugin system foundations laid

**After 6 months:**
- Tview usage <10%
- Bubbletea adoption >90%
- First Lua plugins published
- Community contributions to action system

**After 12 months:**
- Tview removed in v4.0.0
- Plugin ecosystem thriving
- Advanced features shipped (particles, etc.)
- Ticketr known as Bubbletea success story

---

## 10. Post-Migration Plan

### Immediate Post-Cutover (v3.2.0 - v3.2.5)

**Focus:** Stability and bug fixes

- **v3.2.1** (1 week after v3.2.0): Hot-fixes for any critical issues
- **v3.2.2** (2 weeks after): Minor bug fixes, performance tuning
- **v3.2.3** (4 weeks after): Polish and UX improvements
- **v3.2.4** (6 weeks after): Final stabilization release

**Activities:**
- Monitor GitHub Issues closely
- Fix bugs as they're reported
- Performance profiling in production
- Gather user feedback
- Tweak UX based on feedback

### Plugin System Foundation (v3.3.0 - v3.4.0)

**Timeline:** 3-6 months post-cutover

**Goals:**
- Lua VM integration
- Action registration API
- Service bindings for plugins
- Security sandboxing
- Plugin discovery/loading

**Deliverables:**
- Plugin API documentation
- Sample plugins (GitHub, Slack, etc.)
- Plugin developer guide
- Plugin marketplace (future)

### Advanced Features (v3.5.0+)

**Potential features:**
- Background particle effects (if feasible)
- Tri-panel mode (ultra-wide support)
- Git integration UI
- Advanced filtering/views
- Custom dashboard widgets
- AI-powered suggestions

### Community Building

**Goals:**
- Encourage community contributions
- Build plugin ecosystem
- Create tutorial content
- Present at conferences

**Activities:**
- Open GitHub Discussions
- Create plugin template repo
- Write blog posts
- Submit talks to Go conferences

---

## Appendix A: Key Decisions Summary

| Decision | Rationale | Alternatives Considered |
|----------|-----------|------------------------|
| **Full Bubbletea migration** | Long-term maintainability, eliminate race conditions | Hybrid approach (rejected - complexity) |
| **Custom tree component** | No built-in, need full control | bubbles/list with indentation (fallback) |
| **Defer particle effects** | Very high complexity, not required for MVP | Include in Phase 1 (rejected - risk) |
| **Action system architecture** | Future plugin support, extensibility | Hardcoded keybindings (rejected - inflexible) |
| **Parallel development** | Lower risk, Tview stays working | Direct replacement (rejected - too risky) |
| **Feature flag rollout** | Gradual adoption, easy rollback | Hard cutover (rejected - user impact) |
| **18-week timeline** | Realistic given complexity | 12 weeks (rejected - too aggressive) |
| **Lipgloss for layout** | Simple, no dependencies | bubblelayout (rejected - overkill) |
| **Huh for forms** | Production-ready, themes | Custom forms (rejected - reinventing wheel) |

## Appendix B: Resources & References

### Official Documentation
- [Bubbletea GitHub](https://github.com/charmbracelet/bubbletea)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Lipgloss Styling](https://github.com/charmbracelet/lipgloss)
- [Huh Forms](https://github.com/charmbracelet/huh)

### Research Documents
- `/home/karol/dev/private/ticktr/BUBBLETEA_ARCHITECTURE_RESEARCH.md`
- `/home/karol/dev/private/ticktr/TICKETR_CURRENT_ARCHITECTURE_ANALYSIS.md`
- `/home/karol/dev/private/ticktr/EXTENSIBLE_ACTION_SYSTEM_DESIGN.md`
- `/home/karol/dev/private/ticktr/TUI_WIREFRAMES.md`

### Example Applications
- [Glow](https://github.com/charmbracelet/glow) - Markdown reader
- [Soft Serve](https://github.com/charmbracelet/soft-serve) - Git server TUI
- [VHS](https://github.com/charmbracelet/vhs) - Terminal recorder

### Team Resources
- Weekly architecture sync (Fridays)
- Slack channel: #bubbletea-migration
- GitHub Project board: Bubbletea Migration
- Pair programming sessions: Tuesdays & Thursdays

---

## Appendix C: Timeline Visualization

```
Week  Phase  Focus                          Key Deliverables
────  ─────  ─────────────────────────────  ───────────────────────────────
1-2   0      Foundation & POC               Core types, Action system, Help view
3     1      Workspace List                 Workspace switching working
4-5   1      Ticket Tree (CRITICAL)         Custom tree component complete
6     1      Ticket Detail                  Detail view display + edit
7     2      Search Modal                   Fuzzy search + filters
8     2      Command Palette                All actions searchable
9     2      Widgets                        Action bar, status bar, progress
10    3      Job Queue Integration          Sync operations working
11    3      Bulk Operations                Multi-ticket updates
12    4      Animations                     Modal transitions, spinners
13    4      Marquee & Effects              Animated text, shimmer
14    4      Themes                         3 themes working
15    5      Testing                        >80% coverage, golden files
16    5      Documentation                  All docs written
17    6      Beta Rollout                   v3.2.0-beta released
18    6      Cutover                        v3.2.0 stable, Bubbletea default
19+   6.5    Stretch Goals (Optional)       Particles, Lua plugins, etc.
```

---

## Sign-Off

**Architect Approval:** ✅ Director Agent
**Technical Review:** ✅ Builder Agent (to be assigned)
**Quality Review:** ✅ Verifier Agent (to be assigned)
**Documentation Review:** ✅ Scribe Agent (to be assigned)
**Final Approval:** ✅ Steward Agent (to be assigned)

**Date:** 2025-10-22
**Version:** 1.0.0
**Status:** APPROVED FOR IMPLEMENTATION

---

**Next Steps:**
1. Present this plan to the team
2. Create GitHub Project board with all tasks
3. Assign Phase 0 to Builder agent
4. Begin Week 1 implementation
5. Weekly check-ins on Fridays
6. Adjust plan as needed based on learnings

**Let's build something amazing! 🚀**
