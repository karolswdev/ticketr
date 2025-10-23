# Ticketr TUI Architecture Analysis: Tview to Bubbletea Migration

**Date:** 2025-10-22
**Author:** Claude (Architecture Analysis)
**Purpose:** Comprehensive analysis of current Tview-based TUI for migration planning to Bubbletea

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Current Architecture Overview](#current-architecture-overview)
3. [Component Inventory](#component-inventory)
4. [State Management Analysis](#state-management-analysis)
5. [Rendering Architecture](#rendering-architecture)
6. [Effects System Deep Dive](#effects-system-deep-dive)
7. [Pain Points & Bubbletea Solutions](#pain-points--bubbletea-solutions)
8. [Migration Complexity Assessment](#migration-complexity-assessment)
9. [Recommendations](#recommendations)

---

## Executive Summary

### Current State
Ticketr uses **Tview** (imperative, callback-driven) with a sophisticated custom effects system built on top. The architecture is mature with ~46 Go files across views, widgets, effects, and sync coordination.

### Key Strengths
- Well-separated concerns (views, widgets, effects)
- Robust async job queue with progress tracking
- Sophisticated visual effects (animator, background effects, shimmer)
- Clean theme system with environment variable configuration

### Key Challenges
- Heavy reliance on callbacks and manual state synchronization
- Complex focus management and modal state tracking
- Race conditions requiring extensive mutex protection (serviceMutex, multiple RWMutex)
- Manual UI update coordination via `app.QueueUpdateDraw()`

### Migration Verdict
**Medium-High Complexity**. Bubbletea's declarative model will simplify state management significantly, but the effects system and async coordination require careful translation.

---

## Current Architecture Overview

### High-Level Structure

```
┌─────────────────────────────────────────────────────────────┐
│                         TUIApp                              │
│  (Central coordinator, owns all views & services)           │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
┌───────▼──────┐   ┌─────────▼──────┐   ┌─────────▼─────────┐
│   Views      │   │    Widgets     │   │     Effects       │
│ (Screens)    │   │  (Components)  │   │   (Animations)    │
└──────────────┘   └────────────────┘   └───────────────────┘
│                  │                    │
├─WorkspaceList   ├─ActionBar          ├─Animator
├─TicketTree      ├─CommandPalette     ├─BackgroundAnimator
├─TicketDetail    ├─ProgressBar        ├─Shimmer
├─SearchView      ├─SlideOut           └─ShadowBox
├─SyncStatusView  └─Marquee
├─HelpView
└─Modals
  ├─WorkspaceModal
  ├─BulkOpsModal
  └─ErrorModal
```

### Core Dependencies

```
Tview (github.com/rivo/tview)
  └─ tcell/v2 (terminal cell interface)

Custom Systems:
  ├─ Router (view lifecycle management)
  ├─ JobQueue (async operations)
  ├─ SyncCoordinator (Jira sync orchestration)
  └─ Theme System (colors, effects config)
```

---

## Component Inventory

### 1. Views (8 components)

| Component | File | Tview Primitive | Complexity | Migration Notes |
|-----------|------|----------------|------------|-----------------|
| **WorkspaceListView** | `views/workspace_list.go` | `*tview.List` | Medium | Needs Bubbletea list component |
| **TicketTreeView** | `views/ticket_tree.go` | `*tview.TreeView` | **High** | Complex tree structure, selection state, vim bindings |
| **TicketDetailView** | `views/ticket_detail.go` | `*tview.Flex` + TextView/Form | **High** | Dual-mode (display/edit), form validation |
| **SearchView** | `views/search.go` | `*tview.InputField` + List | Medium | Fuzzy search integration |
| **SyncStatusView** | `views/sync_status.go` | `*tview.TextView` | Low | Simple progress display |
| **HelpView** | `views/help.go` | `*tview.TextView` | Low | Static content |
| **CommandPaletteView** | `views/command.go` | `*tview.Modal` | Medium | Command filtering & execution |
| **WorkspaceModal** | `views/workspace_modal.go` | `*tview.Form` in Modal | Medium | Form validation, profile creation |
| **BulkOperationsModal** | `views/bulk_operations_modal.go` | `*tview.Modal` | Medium | Multi-ticket operations |

**Bubbletea Equivalents:**
- **TreeView**: No built-in equivalent → Need custom implementation or use `github.com/charmbracelet/bubbles/tree`
- **Forms**: No built-in equivalent → Use `github.com/charmbracelet/huh` (form library)
- **Lists**: Use `github.com/charmbracelet/bubbles/list`
- **TextInput**: Use `github.com/charmbracelet/bubbles/textinput`
- **Viewport**: Use `github.com/charmbracelet/bubbles/viewport` (for scrolling)

### 2. Widgets (6 components)

| Widget | File | Purpose | Complexity | Migration Notes |
|--------|------|---------|------------|-----------------|
| **ActionBar** | `widgets/actionbar.go` | Context-aware keybindings display | **High** | Marquee integration, responsive truncation |
| **CommandPalette** | `widgets/palette.go` | Enhanced command search | Medium | Uses command registry |
| **ProgressBar** | `widgets/progressbar.go` | ASCII progress with ETA | Medium | Custom renderer, shimmer integration |
| **SlideOut** | `widgets/slideout.go` | Left-side overlay panel | Medium | Layered rendering (cosmic background) |
| **Marquee** | `widgets/marquee.go` | Text overflow animation | **High** | Complex theatrical animation system |

**Migration Approach:**
- Widgets become **Bubbletea components** with their own `Update()` and `View()` methods
- Most widgets have **isolated state** - easy to encapsulate
- **Marquee** is the most complex - requires animation timing in Update() loop

### 3. Effects System (5 components)

| Effect | File | Purpose | Complexity | Migration Notes |
|--------|------|---------|------------|-----------------|
| **Animator** | `effects/animator.go` | Core animation engine | **High** | Context-based lifecycle, animation registry |
| **BackgroundAnimator** | `effects/background.go` | Ambient particle effects | **Very High** | Custom `tcell.Screen` drawing, particle system |
| **Shimmer** | `effects/shimmer.go` | Progress bar shimmer | Low | Simple character substitution |
| **ShadowBox** | `effects/shadowbox.go` | Drop shadows for modals | Low | Border decoration |
| **Borders** | `effects/borders.go` | Custom border styles | Low | Character mapping |

**Key Challenge:** Tview allows **direct screen manipulation** (`screen.SetContent()`). Bubbletea uses **string-based rendering**. Background effects need complete rearchitecture.

### 4. Infrastructure Components

| Component | File | Purpose | Complexity |
|-----------|------|---------|------------|
| **Router** | `router.go` | View lifecycle management | Low |
| **JobQueue** | `../tui/jobs/queue.go` | Async work orchestration | Medium |
| **SyncCoordinator** | `sync/coordinator.go` | Jira sync operations | Medium |
| **Theme** | `theme/theme.go` | Color & effects config | Low |
| **Commands Registry** | `commands/registry.go` | Command catalog | Low |

---

## State Management Analysis

### Current Approach: Imperative Callbacks

```go
// Example from app.go
type TUIApp struct {
    // State scattered across fields
    currentFocus string
    focusOrder   []string
    inModal      bool
    currentJobID jobs.JobID

    // Mutex for race prevention
    serviceMutex stdSync.RWMutex

    // Callbacks for coordination
    onTicketSelected func(*domain.Ticket)
    onSave           func(*domain.Ticket) error
    onClose          func()
}

// Manually propagate state changes
func (t *TUIApp) setFocus(viewName string) {
    t.currentFocus = viewName
    t.updateFocus() // Manually sync UI
}
```

**Problems:**
1. **State synchronization** is manual (`updateFocus()`, `updateActionBar()`)
2. **Race conditions** require extensive mutex usage
3. **Callback hell** for inter-component communication
4. **Implicit state** - hard to reason about "what changed"

### Bubbletea Approach: Declarative Messages

```go
// Proposed Bubbletea model
type Model struct {
    // All state in one place
    focus        Focus
    modal        *Modal
    activeJob    *Job
    tickets      []Ticket
    selectedIDs  map[string]bool

    // Sub-models
    tree   treeModel
    detail detailModel
    search searchModel
}

// State changes via messages
type focusChangedMsg struct{ newFocus Focus }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case focusChangedMsg:
        m.focus = msg.newFocus
        // Automatic re-render, no manual sync
    }
    return m, nil
}
```

**Benefits:**
1. **Single source of truth** - all state in Model struct
2. **Predictable updates** - only via Update() function
3. **No race conditions** - single-threaded event loop
4. **Easy debugging** - log all messages to see state transitions

---

## Rendering Architecture

### Current: Tview Primitives + Manual Drawing

```go
// Tview approach: Build widget tree
func (t *TUIApp) setupApp() error {
    mainLayout := tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(contentLayout, 0, 1, true).
        AddItem(actionBar.Primitive(), 3, 0, false)

    // Manual focus management
    t.app.SetFocus(t.ticketTreeView.Primitive())

    // Async updates require coordination
    go func() {
        result := fetchData()
        t.app.QueueUpdateDraw(func() {
            t.updateView(result)
        })
    }()
}
```

### Proposed: Bubbletea String Rendering

```go
// Bubbletea approach: Render to string
func (m Model) View() string {
    mainContent := lipgloss.JoinVertical(
        lipgloss.Left,
        m.tree.View(),      // Tree component
        m.detail.View(),    // Detail component
    )

    actionBar := m.actionBar.View()

    return lipgloss.JoinVertical(
        lipgloss.Left,
        mainContent,
        actionBar,
    )
}

// Async updates send messages
func fetchDataCmd() tea.Msg {
    return dataFetchedMsg{data: fetchData()}
}
```

**Key Difference:** Bubbletea is **declarative** - you describe what to show, not how to update it.

---

## Effects System Deep Dive

### 1. Animator System

**Current Implementation:**
- **Context-based lifecycle** - each animation has its own context
- **Registry pattern** - animations tracked by name
- **Frame-rate control** - configurable intervals
- **Graceful shutdown** - wait for all animations

```go
type Animator struct {
    animations map[string]*Animation
    ctx        context.Context
    cancel     context.CancelFunc
    wg         sync.WaitGroup
    enabled    bool // Global kill switch
}

func (a *Animator) Start(name string, interval time.Duration, handler func() bool)
```

**Bubbletea Translation:**
```go
// Convert to tick-based system
type AnimationMsg struct {
    name     string
    frame    int
    complete bool
}

func animationTickCmd(interval time.Duration) tea.Cmd {
    return tea.Tick(interval, func(t time.Time) tea.Msg {
        return AnimationMsg{frame: computeFrame(t)}
    })
}

// In Update()
case AnimationMsg:
    if !m.config.Motion {
        return m, nil // Respect kill switch
    }
    m.animState.frame = msg.frame
    return m, animationTickCmd(interval) // Schedule next tick
```

**Migration Complexity: Medium**
- Replace goroutines with `tea.Tick()` commands
- Animation state moves to model fields
- No more manual `QueueUpdateDraw()` - automatic re-render

### 2. Background Animator (Particles)

**Current Implementation:**
- **Direct screen manipulation** using `tcell.Screen.SetContent()`
- **Particle system** with 100+ particles
- **Custom Draw() method** on overlay primitive

```go
func (bo *BackgroundOverlay) Draw(screen tcell.Screen) {
    bo.Box.Draw(screen)

    for _, p := range bo.animator.particles {
        px := x + p.X
        py := y + p.Y
        style := tcell.StyleDefault.Foreground(p.Color).Dim(true)
        screen.SetContent(px, py, p.Char, nil, style)
    }
}
```

**Bubbletea Challenge:**
- **No direct screen access** - everything is string-based
- Must render particles as **ANSI escape sequences**
- Performance concern: rendering 100+ positioned characters per frame

**Proposed Solutions:**

**Option A: Disable particle effects** (simplest)
- Remove BackgroundAnimator entirely
- Focus on functional UI first
- Add back later if needed

**Option B: Layer-based rendering** (complex but preserves feature)
```go
func (m Model) View() string {
    // Layer 1: Background particles
    bg := renderParticleLayer(m.particles, m.width, m.height)

    // Layer 2: UI content
    content := m.renderMainUI()

    // Composite layers using ANSI positioning
    return compositeAnsiLayers(bg, content)
}

func renderParticleLayer(particles []Particle, w, h int) string {
    // Create blank canvas
    canvas := make([][]rune, h)
    for i := range canvas {
        canvas[i] = make([]rune, w)
        for j := range canvas[i] {
            canvas[i][j] = ' '
        }
    }

    // Place particles
    for _, p := range particles {
        if p.Y < h && p.X < w {
            canvas[p.Y][p.X] = p.Char
        }
    }

    // Render to string with ANSI colors
    return canvasToAnsi(canvas, particles)
}
```

**Migration Complexity: Very High** (if preserving feature)

### 3. Shimmer Effect

**Current:** Character substitution in progress bar string
**Bubbletea:** Same approach works - operate on strings before rendering

**Migration Complexity: Low**

### 4. Marquee Animation

**Current:**
- **Theatrical animations** - slide in, blink, slide out
- **Frame-based timing** - 30 FPS update loop
- **Multi-item queue** with state machine

**Bubbletea Translation:**
```go
type MarqueeModel struct {
    items       []string
    currentIdx  int
    phase       AnimationPhase
    frameCount  int
    blinkOn     bool
}

func (m MarqueeModel) Update(msg tea.Msg) (MarqueeModel, tea.Cmd) {
    switch msg.(type) {
    case tea.KeyMsg:
        return m, nil
    case AnimationTickMsg:
        m.frameCount++
        m.updatePhase() // State machine progression
        return m, tea.Tick(33*time.Millisecond, func(t time.Time) tea.Msg {
            return AnimationTickMsg{}
        })
    }
    return m, nil
}

func (m MarqueeModel) View() string {
    switch m.phase {
    case PhaseSlideIn:
        return m.renderSlideIn()
    case PhaseCenter:
        if m.blinkOn {
            return m.centerText(m.items[m.currentIdx])
        }
        return strings.Repeat(" ", m.width)
    case PhaseSlideOut:
        return m.renderSlideOut()
    }
    return ""
}
```

**Migration Complexity: Medium-High**
- State machine translates cleanly
- Timing via `tea.Tick()` instead of goroutine
- String manipulation logic can be reused

---

## Pain Points & Bubbletea Solutions

### Pain Point 1: Race Conditions

**Current Problem:**
```go
// From app.go - Critical race condition fix
// CRITICAL (Phase 6.5 Emergency Fix): serviceMutex protects service fields
serviceMutex    stdSync.RWMutex
pushService     *services.PushService
pullService     *services.PullService

// Access requires lock
func (t *TUIApp) handlePull() {
    t.serviceMutex.RLock()
    pullService := t.pullService
    t.serviceMutex.RUnlock()

    if pullService == nil {
        // Handle error
        return
    }

    pullService.Pull(...)
}
```

**Bubbletea Solution:**
```go
// No mutex needed - single-threaded
type Model struct {
    pullService *services.PullService
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case pullMsg:
        // Safe to access - no concurrent access
        return m, m.pullService.PullCmd()
}
```

### Pain Point 2: Async Coordination

**Current Problem:**
```go
// Manual coordination of async updates
go func() {
    tickets, err := t.ticketQuery.ListByWorkspace(ws.ID)

    t.app.QueueUpdateDraw(func() {
        t.isLoading = false
        if err != nil {
            t.showError(...)
        } else {
            t.buildTree(tickets)
        }
    })
}()
```

**Bubbletea Solution:**
```go
// Return command, handle result as message
func loadTicketsCmd(workspaceID string) tea.Cmd {
    return func() tea.Msg {
        tickets, err := ticketQuery.ListByWorkspace(workspaceID)
        return ticketsLoadedMsg{tickets, err}
    }
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case ticketsLoadedMsg:
        if msg.err != nil {
            m.error = msg.err
        } else {
            m.tickets = msg.tickets
        }
        return m, nil // Automatic re-render
}
```

### Pain Point 3: Focus Management

**Current Problem:**
```go
// Manual focus tracking and synchronization
currentFocus string
focusOrder   []string

func (t *TUIApp) setFocus(viewName string) {
    t.currentFocus = viewName
    t.updateFocus() // Must manually update all views
}

func (t *TUIApp) updateFocus() {
    t.ticketTreeView.SetFocused(t.currentFocus == "ticket_tree")
    t.ticketDetailView.SetFocused(t.currentFocus == "ticket_detail")
    t.actionBar.SetContext(t.getContextFromFocus())
    // ... more manual updates
}
```

**Bubbletea Solution:**
```go
type Focus int

const (
    FocusTree Focus = iota
    FocusDetail
    FocusWorkspace
)

type Model struct {
    focus Focus
}

func (m Model) View() string {
    // Focus is just data - rendering reacts automatically
    treeStyle := lipgloss.NewStyle()
    if m.focus == FocusTree {
        treeStyle = treeStyle.BorderForeground(lipgloss.Color("green"))
    }

    return treeStyle.Render(m.tree.View())
}
```

### Pain Point 4: Modal State Tracking

**Current Problem:**
```go
// Boolean flags scattered everywhere
inModal bool

func (t *TUIApp) showSearch() {
    t.inModal = true
    t.searchView.OnShow()
    t.app.SetRoot(t.searchView.Primitive(), true)
}

// Easy to get into inconsistent state
```

**Bubbletea Solution:**
```go
type Modal interface {
    Update(tea.Msg) (Modal, tea.Cmd)
    View() string
}

type Model struct {
    modal Modal // nil = no modal
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if m.modal != nil {
        // Modal intercepts input
        newModal, cmd := m.modal.Update(msg)
        m.modal = newModal
        return m, cmd
    }
    // Normal input handling
}

func (m Model) View() string {
    view := m.mainView()
    if m.modal != nil {
        view = lipgloss.Place(
            m.width, m.height,
            lipgloss.Center, lipgloss.Center,
            m.modal.View(),
            lipgloss.WithWhitespaceBackground(lipgloss.Color("0")),
        )
    }
    return view
}
```

---

## Migration Complexity Assessment

### Easy (1-2 days each)

| Component | Effort | Notes |
|-----------|--------|-------|
| HelpView | 4 hours | Static content, no state |
| SyncStatusView | 4 hours | Simple text display with color |
| Theme System | 6 hours | Convert to Lipgloss styles |
| ProgressBar | 8 hours | Reuse string rendering logic |
| Shimmer Effect | 4 hours | String manipulation, works as-is |

### Medium (3-5 days each)

| Component | Effort | Notes |
|-----------|--------|-------|
| WorkspaceListView | 2 days | Use bubbles/list, add handlers |
| SearchView | 2 days | Use bubbles/textinput + list |
| CommandPalette | 2 days | Fuzzy search + bubbles/list |
| WorkspaceModal | 3 days | Use huh for forms |
| BulkOperationsModal | 3 days | Checkbox list + form |
| JobQueue Integration | 3 days | Convert to tea.Cmd pattern |
| SlideOut | 3 days | Layout component, no cosmic bg |

### Hard (1-2 weeks each)

| Component | Effort | Notes |
|-----------|--------|-------|
| TicketTreeView | 1.5 weeks | Custom tree component, vim bindings, selection |
| TicketDetailView | 1 week | Dual-mode (display/edit), scrolling, form validation |
| ActionBar + Marquee | 1.5 weeks | Complex animation system, responsive layout |
| Animator System | 1 week | Convert to tick-based, maintain animation registry |

### Very Hard (2-4 weeks each)

| Component | Effort | Notes |
|-----------|--------|-------|
| BackgroundAnimator | 3 weeks | Particle system rearchitecture, performance tuning |
| Full Integration | 2 weeks | Wire all components, debug state transitions |

### Total Estimated Timeline

**Minimum Viable Product (no effects):** 6-8 weeks
- Core views (tree, detail, list, modals)
- Basic navigation and keyboard shortcuts
- Async job integration
- Theme system

**Full Feature Parity:** 10-14 weeks
- All effects except background particles
- Complete animation system
- All widgets and modals

**Complete Migration (all features):** 14-18 weeks
- Background particle effects
- Full polish and performance optimization

---

## Recommendations

### Phase 1: Foundation (Weeks 1-2)

**Goal:** Establish Bubbletea architecture patterns

1. **Set up project structure**
   ```
   internal/adapters/tui-bubbletea/
   ├── model.go           # Root model
   ├── update.go          # Update logic
   ├── view.go            # View rendering
   ├── messages.go        # Message types
   ├── components/
   │   ├── tree/          # Tree component
   │   ├── detail/        # Detail component
   │   └── list/          # List component
   └── styles/            # Lipgloss styles
   ```

2. **Migrate theme system** → Lipgloss styles
   - Color scheme mapping
   - Border styles
   - Layout constants

3. **Build simple view first** (HelpView or SyncStatusView)
   - Prove out message passing
   - Test async command pattern

### Phase 2: Core Views (Weeks 3-6)

**Goal:** Implement primary user flows

1. **WorkspaceListView** → bubbles/list
2. **TicketTreeView** → Custom tree component
   - Research: Use `github.com/charmbracelet/bubbles/tree` or build custom
   - Implement vim keybindings
   - Selection state management
3. **TicketDetailView** → viewport + forms
   - Display mode: bubbles/viewport
   - Edit mode: huh forms library
4. **Navigation** between views
   - Focus management
   - Tab cycling

### Phase 3: Modals & Widgets (Weeks 7-9)

**Goal:** Secondary features and polish

1. **SearchView** → textinput + list
2. **CommandPalette** → textinput + list + fuzzy search
3. **ActionBar** → Simple version (no marquee initially)
4. **Modals** → huh forms + overlay layout

### Phase 4: Async & Sync (Weeks 10-11)

**Goal:** Jira integration and background jobs

1. **JobQueue** → tea.Cmd pattern
   - Pull operations
   - Push operations
   - Progress reporting
2. **SyncCoordinator** integration
3. **Progress indicators**

### Phase 5: Effects (Weeks 12-14)

**Goal:** Visual polish (optional)

1. **Basic animations** (spinners, progress shimmer)
2. **Marquee** (if action bar needs it)
3. **Animator system** (if needed for future)
4. **Decision point:** Skip background particles or invest 3 weeks

### Parallel Development Strategy

**Keep Tview working** while building Bubbletea in parallel:
```
internal/adapters/
├── tui/              # Current Tview (keep working)
└── tui-bubbletea/    # New implementation
```

**Feature flag for testing:**
```bash
TICKETR_USE_BUBBLETEA=true ./ticketr
```

**Incremental cutover:**
- Ship Bubbletea as experimental in minor release
- Gather feedback for 1-2 sprints
- Make Bubbletea default in next major release
- Deprecate Tview in following release

---

## Appendix: Key Files Reference

### Core Application
- `/home/karol/dev/private/ticktr/internal/adapters/tui/app.go` - Main app coordinator
- `/home/karol/dev/private/ticktr/internal/adapters/tui/router.go` - View routing

### Views
- `/home/karol/dev/private/ticktr/internal/adapters/tui/views/ticket_tree.go` - Tree implementation
- `/home/karol/dev/private/ticktr/internal/adapters/tui/views/ticket_detail.go` - Detail view
- `/home/karol/dev/private/ticktr/internal/adapters/tui/views/workspace_list.go` - Workspace list

### Widgets
- `/home/karol/dev/private/ticktr/internal/adapters/tui/widgets/actionbar.go` - Action bar
- `/home/karol/dev/private/ticktr/internal/adapters/tui/widgets/marquee.go` - Marquee animation
- `/home/karol/dev/private/ticktr/internal/adapters/tui/widgets/progressbar.go` - Progress display

### Effects
- `/home/karol/dev/private/ticktr/internal/adapters/tui/effects/animator.go` - Core animator
- `/home/karol/dev/private/ticktr/internal/adapters/tui/effects/background.go` - Particle system
- `/home/karol/dev/private/ticktr/internal/adapters/tui/effects/shimmer.go` - Shimmer effect

### Infrastructure
- `/home/karol/dev/private/ticktr/internal/adapters/tui/sync/coordinator.go` - Sync orchestration
- `/home/karol/dev/private/ticktr/internal/tui/jobs/queue.go` - Job queue system
- `/home/karol/dev/private/ticktr/internal/adapters/tui/theme/theme.go` - Theme configuration

---

## Conclusion

The current Ticketr TUI is well-architected but suffers from the inherent complexity of imperative, callback-driven UI frameworks. **Bubbletea will simplify** the core application logic significantly through its declarative model, but the sophisticated **effects system requires careful translation**.

**Primary benefits of migration:**
1. ✅ Eliminate race conditions (no more mutex nightmare)
2. ✅ Simplify state management (single source of truth)
3. ✅ Easier testing (pure functions, deterministic updates)
4. ✅ Better long-term maintainability

**Primary challenges:**
1. ⚠️ Tree component - no built-in Bubbletea equivalent
2. ⚠️ Background particle effects - complete rearchitecture needed
3. ⚠️ Marquee animations - complex timing logic
4. ⚠️ Learning curve for team

**Recommended approach:** Incremental migration over 10-14 weeks, with background effects as optional "stretch goal" after core functionality is proven.
