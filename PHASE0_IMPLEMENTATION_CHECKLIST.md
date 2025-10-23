# Phase 0: Foundation & Proof-of-Concept - Implementation Checklist

**Timeline:** Weeks 1-2
**Goal:** Establish architecture, validate patterns, build foundation
**Assignee:** Builder Agent

---

## Week 1: Setup & Core Patterns

### Day 1: Directory Structure ✅
- [ ] Create `internal/adapters/bubbletea/` directory
- [ ] Create subdirectories:
  ```
  bubbletea/
  ├── components/
  │   ├── tree/
  │   ├── detail/
  │   ├── workspaces/
  │   ├── search/
  │   ├── palette/
  │   └── bulk/
  ├── widgets/
  │   ├── actionbar/
  │   ├── statusbar/
  │   ├── progress/
  │   └── marquee/
  ├── actions/
  │   ├── predicates/
  │   ├── modifiers/
  │   └── builtin/
  ├── effects/
  └── styles/
  ```
- [ ] Create initial files:
  - `model.go` - Root model definition
  - `init.go` - Initialization logic
  - `update.go` - Update reducer
  - `view.go` - View renderer
  - `messages.go` - Message type definitions
  - `commands.go` - Command factories

**Validation:** All directories exist, files compile without errors

---

### Day 2: Core Type Definitions ✅

#### File: `internal/adapters/bubbletea/model.go`
```go
package bubbletea

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/yourorg/ticketr/internal/adapters/bubbletea/actions"
    "github.com/yourorg/ticketr/internal/domain"
)

type Model struct {
    // Window dimensions
    width  int
    height int

    // Context & Focus
    context        actions.Context
    contextManager *actions.ContextManager
    focus          Focus

    // Action System
    actionRegistry     *actions.Registry
    keybindingResolver *actions.KeybindingResolver
    executor           *actions.Executor

    // Services (dependency injection)
    services *Services

    // Configuration
    config *Config

    // UI State
    showHelp bool
}

type Services struct {
    TicketService    domain.TicketService
    WorkspaceService domain.WorkspaceService
    SyncService      domain.SyncService
}

type Config struct {
    Theme          string
    MotionEnabled  bool
    EffectsEnabled bool
}

type Focus int

const (
    FocusTree Focus = iota
    FocusDetail
    FocusWorkspace
)
```

**Tasks:**
- [ ] Define `Model` struct
- [ ] Define `Services` struct
- [ ] Define `Config` struct
- [ ] Define `Focus` enum
- [ ] Add comprehensive comments

**Validation:** File compiles, types are complete

---

#### File: `internal/adapters/bubbletea/messages.go`
```go
package bubbletea

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/yourorg/ticketr/internal/domain"
)

// Application lifecycle messages
type appReadyMsg struct{}
type appQuitMsg struct{}

// Ticket messages
type ticketOpenedMsg struct {
    ticket *domain.Ticket
}

type ticketSavedMsg struct {
    ticket *domain.Ticket
}

type ticketsLoadedMsg struct {
    tickets []*domain.Ticket
    err     error
}

// Workspace messages
type workspaceChangedMsg struct {
    workspace *domain.Workspace
}

// Sync messages
type syncStartedMsg struct {
    operation string
}

type syncProgressMsg struct {
    progress float64
}

type syncCompletedMsg struct {
    operation string
    err       error
}

// Modal messages
type modalOpenedMsg struct {
    modal Modal
}

type modalClosedMsg struct{}

// Error messages
type errorMsg struct {
    err error
}
```

**Tasks:**
- [ ] Define all domain message types
- [ ] Add comments for each message type
- [ ] Group by category

**Validation:** All messages compile, comprehensive coverage

---

### Day 3: Action System Foundation ✅

#### File: `internal/adapters/bubbletea/actions/action.go`
```go
package actions

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/key"
)

type ActionID string
type ActionCategory string

const (
    CategoryNavigation ActionCategory = "Navigation"
    CategoryEdit       ActionCategory = "Edit"
    CategoryView       ActionCategory = "View"
    CategorySync       ActionCategory = "Sync"
    CategoryWorkspace  ActionCategory = "Workspace"
    CategoryBulk       ActionCategory = "Bulk Operations"
    CategorySystem     ActionCategory = "System"
)

type Action struct {
    ID          ActionID
    Name        string
    Description string
    Category    ActionCategory
    Keybindings []KeyPattern
    Contexts    []Context
    Predicate   PredicateFunc
    ShowInUI    ShowInUIFunc
    Execute     ExecuteFunc
    Icon        string
    Tags        []string
    Metadata    map[string]string
    Modifiers   []ActionModifier
}

type KeyPattern struct {
    Key     string
    Alt     bool
    Ctrl    bool
    Shift   bool
    Binding key.Binding
}

type PredicateFunc func(ctx *ActionContext) bool
type ShowInUIFunc func(ctx *ActionContext) bool
type ExecuteFunc func(ctx *ActionContext) tea.Cmd

type ActionContext struct {
    Context           Context
    SelectedTickets   []string
    SelectedWorkspace *WorkspaceState
    FocusedItem       interface{}
    HasUnsavedChanges bool
    IsSyncing         bool
    IsOffline         bool
    Config            *UserConfig
    PluginData        map[string]interface{}
    Width             int
    Height            int
    Services          *ServiceContainer
}

type ActionModifier interface {
    Before(ctx *ActionContext) error
    After(ctx *ActionContext, result tea.Cmd, err error) (tea.Cmd, error)
}
```

**Tasks:**
- [ ] Define `Action` struct
- [ ] Define `ActionContext` struct
- [ ] Define function types (PredicateFunc, ExecuteFunc, etc.)
- [ ] Define `ActionModifier` interface
- [ ] Add comprehensive comments

**Validation:** All types compile, design is complete

---

#### File: `internal/adapters/bubbletea/actions/context.go`
```go
package actions

type Context string

const (
    ContextWorkspaceList Context = "workspace_list"
    ContextTicketTree    Context = "ticket_tree"
    ContextTicketDetail  Context = "ticket_detail"
    ContextSearch        Context = "search"
    ContextCommandPalette Context = "command_palette"
    ContextModal         Context = "modal"
    ContextSyncing       Context = "syncing"
    ContextHelp          Context = "help"
    ContextGlobal        Context = "*"
)

type ContextState struct {
    Current  Context
    Previous Context
    Stack    []Context
    Metadata map[string]interface{}
}

type ContextManager struct {
    state    ContextState
    onChange []func(old, new Context)
}

func NewContextManager(initial Context) *ContextManager {
    return &ContextManager{
        state: ContextState{
            Current:  initial,
            Previous: "",
            Stack:    []Context{initial},
            Metadata: make(map[string]interface{}),
        },
        onChange: []func(Context, Context){},
    }
}

func (cm *ContextManager) Switch(newContext Context) {
    old := cm.state.Current
    if old == newContext {
        return
    }
    cm.state.Previous = old
    cm.state.Current = newContext
    cm.state.Stack = append(cm.state.Stack, newContext)
    for _, fn := range cm.onChange {
        fn(old, newContext)
    }
}

func (cm *ContextManager) Current() Context {
    return cm.state.Current
}

// ... additional methods
```

**Tasks:**
- [ ] Define `Context` enum
- [ ] Define `ContextManager` struct
- [ ] Implement Switch() method
- [ ] Implement Push()/Pop() for modals
- [ ] Implement OnChange() observer pattern

**Validation:** Context switching works, tests pass

---

### Day 4: Action Registry ✅

#### File: `internal/adapters/bubbletea/actions/registry.go`
```go
package actions

import (
    "fmt"
    "sync"
)

type Registry struct {
    actions   map[ActionID]*Action
    byContext map[Context][]*Action
    byKey     map[string][]*Action
    mutex     sync.RWMutex
}

func NewRegistry() *Registry {
    return &Registry{
        actions:   make(map[ActionID]*Action),
        byContext: make(map[Context][]*Action),
        byKey:     make(map[string][]*Action),
    }
}

func (r *Registry) Register(action *Action) error {
    if action.ID == "" {
        return fmt.Errorf("action ID is required")
    }
    if action.Execute == nil {
        return fmt.Errorf("action execute function is required")
    }

    r.mutex.Lock()
    defer r.mutex.Unlock()

    if _, exists := r.actions[action.ID]; exists {
        return fmt.Errorf("action ID already registered: %s", action.ID)
    }

    r.actions[action.ID] = action

    for _, ctx := range action.Contexts {
        r.byContext[ctx] = append(r.byContext[ctx], action)
    }

    for _, keyPattern := range action.Keybindings {
        keyStr := keyPattern.String()
        r.byKey[keyStr] = append(r.byKey[keyStr], action)
    }

    return nil
}

func (r *Registry) Get(id ActionID) (*Action, bool) {
    r.mutex.RLock()
    defer r.mutex.RUnlock()
    action, exists := r.actions[id]
    return action, exists
}

func (r *Registry) ActionsForContext(ctx Context, actx *ActionContext) []*Action {
    // Implementation
}

// ... additional methods
```

**Tasks:**
- [ ] Implement `NewRegistry()`
- [ ] Implement `Register()`
- [ ] Implement `Get()`
- [ ] Implement `ActionsForContext()`
- [ ] Implement `ActionsForKey()`
- [ ] Implement `Search()`
- [ ] Add thread-safety (mutex)

**Validation:** Can register/retrieve actions, tests pass

---

### Day 5: Predicates & Testing ✅

#### File: `internal/adapters/bubbletea/actions/predicates/predicates.go`
```go
package predicates

import "github.com/yourorg/ticketr/internal/adapters/bubbletea/actions"

func Always() actions.PredicateFunc {
    return func(ctx *actions.ActionContext) bool {
        return true
    }
}

func Never() actions.PredicateFunc {
    return func(ctx *actions.ActionContext) bool {
        return false
    }
}

func HasSelection() actions.PredicateFunc {
    return func(ctx *actions.ActionContext) bool {
        return len(ctx.SelectedTickets) > 0
    }
}

func HasSingleSelection() actions.PredicateFunc {
    return func(ctx *actions.ActionContext) bool {
        return len(ctx.SelectedTickets) == 1
    }
}

func And(predicates ...actions.PredicateFunc) actions.PredicateFunc {
    return func(ctx *actions.ActionContext) bool {
        for _, pred := range predicates {
            if !pred(ctx) {
                return false
            }
        }
        return true
    }
}

func Or(predicates ...actions.PredicateFunc) actions.PredicateFunc {
    return func(ctx *actions.ActionContext) bool {
        for _, pred := range predicates {
            if pred(ctx) {
                return true
            }
        }
        return false
    }
}

func Not(pred actions.PredicateFunc) actions.PredicateFunc {
    return func(ctx *actions.ActionContext) bool {
        return !pred(ctx)
    }
}
```

**Tasks:**
- [ ] Implement common predicates (Always, Never, HasSelection, etc.)
- [ ] Implement combinator predicates (And, Or, Not)
- [ ] Write comprehensive tests
- [ ] Achieve 100% coverage

**Tests:**
```go
func TestPredicates(t *testing.T) {
    t.Run("Always returns true", func(t *testing.T) {
        pred := Always()
        ctx := &actions.ActionContext{}
        assert.True(t, pred(ctx))
    })

    t.Run("And combines predicates", func(t *testing.T) {
        pred := And(
            HasSelection(),
            HasSingleSelection(),
        )
        ctx := &actions.ActionContext{
            SelectedTickets: []string{"ticket1"},
        }
        assert.True(t, pred(ctx))
    })

    // ... more tests
}
```

**Validation:** All tests pass, 100% coverage

---

## Week 2: Proof-of-Concept UI

### Day 6: Help View Implementation ✅

#### File: `internal/adapters/bubbletea/components/help/model.go`
```go
package help

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type Model struct {
    width   int
    height  int
    content string
    visible bool
}

func New() Model {
    return Model{
        content: generateHelpContent(),
        visible: false,
    }
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    case tea.KeyMsg:
        if msg.String() == "?" || msg.String() == "esc" {
            m.visible = !m.visible
        }
    }
    return m, nil
}

func (m Model) View() string {
    if !m.visible {
        return ""
    }

    style := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("62")).
        Padding(1, 2).
        Width(m.width - 4).
        Height(m.height - 4)

    return lipgloss.Place(
        m.width, m.height,
        lipgloss.Center, lipgloss.Center,
        style.Render(m.content),
    )
}

func generateHelpContent() string {
    return `Ticketr Help

Navigation:
  j/k      - Move down/up
  h/l      - Collapse/expand
  Tab      - Switch panel
  ?        - Toggle help

Actions:
  Enter    - Open ticket
  Space    - Select ticket
  /        - Search
  :        - Command palette
  q        - Quit

Press ? or Esc to close this help.`
}
```

**Tasks:**
- [ ] Create `help/model.go`
- [ ] Implement Model/Update/View pattern
- [ ] Add toggle logic (? key)
- [ ] Style with Lipgloss
- [ ] Test window resize

**Validation:** Help view toggles correctly, looks good

---

### Day 7: Root Model Integration ✅

#### File: `internal/adapters/bubbletea/model.go` (updated)
```go
package bubbletea

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/yourorg/ticketr/internal/adapters/bubbletea/actions"
    "github.com/yourorg/ticketr/internal/adapters/bubbletea/components/help"
)

type Model struct {
    // ... existing fields

    // Components
    help help.Model

    // ... rest of fields
}

func New(services *Services, config *Config) Model {
    return Model{
        services:       services,
        config:         config,
        contextManager: actions.NewContextManager(actions.ContextWorkspaceList),
        actionRegistry: actions.NewRegistry(),
        help:          help.New(),
    }
}

func (m Model) Init() tea.Cmd {
    // Register built-in actions
    m.registerBuiltinActions()

    // Create keybinding resolver
    m.keybindingResolver = actions.NewKeybindingResolver(
        m.actionRegistry,
        m.contextManager,
        &actions.UserConfig{},
    )

    // Create executor
    m.executor = actions.NewExecutor(m.actionRegistry)

    return tea.Batch(
        m.help.Init(),
    )
}
```

**Tasks:**
- [ ] Update `New()` constructor
- [ ] Implement `Init()` method
- [ ] Integrate help component
- [ ] Wire up action system

**Validation:** Model initializes correctly

---

#### File: `internal/adapters/bubbletea/update.go`
```go
package bubbletea

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle global quit
        if msg.String() == "q" || msg.String() == "ctrl+c" {
            return m, tea.Quit
        }

        // Help toggle
        if msg.String() == "?" {
            m.help.visible = !m.help.visible
            return m, nil
        }

    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        m.help, cmd := m.help.Update(msg)
        cmds = append(cmds, cmd)
        return m, tea.Batch(cmds...)
    }

    // Update help component
    var cmd tea.Cmd
    m.help, cmd = m.help.Update(msg)
    cmds = append(cmds, cmd)

    return m, tea.Batch(cmds...)
}
```

**Tasks:**
- [ ] Implement Update() reducer
- [ ] Handle quit (q, Ctrl+C)
- [ ] Handle help toggle (?)
- [ ] Handle WindowSizeMsg
- [ ] Route to help component

**Validation:** Update logic works, no panics

---

#### File: `internal/adapters/bubbletea/view.go`
```go
package bubbletea

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
    // If help is visible, show help overlay
    if m.help.visible {
        return m.help.View()
    }

    // Otherwise show placeholder
    style := lipgloss.NewStyle().
        Foreground(lipgloss.Color("86")).
        Bold(true).
        Padding(1)

    return lipgloss.Place(
        m.width, m.height,
        lipgloss.Center, lipgloss.Center,
        style.Render("Ticketr Bubbletea POC\nPress ? for help\nPress q to quit"),
    )
}
```

**Tasks:**
- [ ] Implement View() renderer
- [ ] Show help overlay when visible
- [ ] Show placeholder otherwise
- [ ] Center content

**Validation:** View renders correctly

---

### Day 8: Entry Point & Testing ✅

#### File: `cmd/tui.go` (updated)
```go
package cmd

import (
    "fmt"
    "os"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/spf13/cobra"

    tuiAdapterBubbletea "github.com/yourorg/ticketr/internal/adapters/bubbletea"
    tuiAdapterTview "github.com/yourorg/ticketr/internal/adapters/tui"
)

var tuiCmd = &cobra.Command{
    Use:   "tui",
    Short: "Start the TUI",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Check feature flag
        if os.Getenv("TICKETR_USE_BUBBLETEA") == "true" {
            return runBubbletea()
        }
        return runTview()
    },
}

func runBubbletea() error {
    // Initialize services
    services := &tuiAdapterBubbletea.Services{
        // ... initialize services
    }

    config := &tuiAdapterBubbletea.Config{
        Theme:          "default",
        MotionEnabled:  true,
        EffectsEnabled: true,
    }

    model := tuiAdapterBubbletea.New(services, config)

    p := tea.NewProgram(
        model,
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),
    )

    if _, err := p.Run(); err != nil {
        return fmt.Errorf("error running Bubbletea TUI: %w", err)
    }

    return nil
}

func runTview() error {
    // Existing Tview implementation
    return tuiAdapterTview.Run()
}
```

**Tasks:**
- [ ] Update `cmd/tui.go` with feature flag
- [ ] Implement `runBubbletea()` function
- [ ] Initialize services (mock for now)
- [ ] Initialize config
- [ ] Create tea.Program with options
- [ ] Run program

**Validation:** Can run with `TICKETR_USE_BUBBLETEA=true ./ticketr`

---

### Day 9: Integration Testing ✅

#### File: `internal/adapters/bubbletea/integration_test.go`
```go
package bubbletea

import (
    "testing"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/x/exp/teatest"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestHelpToggle(t *testing.T) {
    model := New(&Services{}, &Config{})

    tm := teatest.NewTestModel(
        t, model,
        teatest.WithInitialTermSize(120, 30),
    )

    // Send ? key
    tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})

    // Wait for help to appear
    teatest.WaitFor(
        t, tm.Output(),
        func(bts []byte) bool {
            return bytes.Contains(bts, []byte("Ticketr Help"))
        },
        teatest.WithCheckInterval(time.Millisecond*100),
        teatest.WithDuration(time.Second*2),
    )

    // Send ? again to close
    tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})

    // Verify help closed
    tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func TestQuit(t *testing.T) {
    model := New(&Services{}, &Config{})

    tm := teatest.NewTestModel(t, model)

    // Send q key
    tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

    // Program should quit
    tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

func TestWindowResize(t *testing.T) {
    model := New(&Services{}, &Config{})

    // Send resize message
    updated, _ := model.Update(tea.WindowSizeMsg{Width: 160, Height: 40})
    m := updated.(Model)

    assert.Equal(t, 160, m.width)
    assert.Equal(t, 40, m.height)
}
```

**Tasks:**
- [ ] Write help toggle test
- [ ] Write quit test
- [ ] Write window resize test
- [ ] Use teatest for integration tests
- [ ] Achieve >80% coverage for core files

**Validation:** All tests pass, coverage >80%

---

### Day 10: Theme System & Polish ✅

#### File: `internal/adapters/bubbletea/styles/theme.go`
```go
package styles

import (
    "os"

    "github.com/charmbracelet/lipgloss"
)

type Theme struct {
    Name           string
    Primary        lipgloss.Color
    Secondary      lipgloss.Color
    Accent         lipgloss.Color
    Success        lipgloss.Color
    Warning        lipgloss.Color
    Error          lipgloss.Color
    Background     lipgloss.Color
    TextPrimary    lipgloss.Color
    TextMuted      lipgloss.Color
    BorderFocused  lipgloss.Border
    BorderUnfocused lipgloss.Border
}

var DefaultTheme = Theme{
    Name:           "default",
    Primary:        lipgloss.Color("#00FF00"),
    Secondary:      lipgloss.Color("#AAAAAA"),
    Accent:         lipgloss.Color("#00FFFF"),
    Success:        lipgloss.Color("#00FF00"),
    Warning:        lipgloss.Color("#FFFF00"),
    Error:          lipgloss.Color("#FF0000"),
    Background:     lipgloss.Color("#1A1A1A"),
    TextPrimary:    lipgloss.Color("#FFFFFF"),
    TextMuted:      lipgloss.Color("#666666"),
    BorderFocused:  lipgloss.DoubleBorder(),
    BorderUnfocused: lipgloss.RoundedBorder(),
}

func LoadTheme() Theme {
    themeName := os.Getenv("TICKETR_THEME")
    if themeName == "" {
        themeName = "default"
    }

    switch themeName {
    case "dark":
        return DarkTheme
    case "arctic":
        return ArcticTheme
    default:
        return DefaultTheme
    }
}
```

**Tasks:**
- [ ] Define `Theme` struct
- [ ] Define `DefaultTheme`
- [ ] Implement `LoadTheme()` function
- [ ] Support TICKETR_THEME env var
- [ ] Create style helpers

**Validation:** Theme loads correctly, env var works

---

## Phase 0 Acceptance Criteria

### ✅ Functional
- [ ] Can run `TICKETR_USE_BUBBLETEA=true ./ticketr`
- [ ] Help screen appears on `?`
- [ ] Help closes on `?` or `Esc`
- [ ] Program quits on `q` or `Ctrl+C`
- [ ] Window resize updates layout

### ✅ Code Quality
- [ ] All files compile without errors
- [ ] All tests pass
- [ ] Test coverage >80% for Phase 0 code
- [ ] No linter warnings
- [ ] All functions documented

### ✅ Architecture
- [ ] Action registry can register actions
- [ ] Predicates can be composed (And, Or, Not)
- [ ] Context manager tracks state
- [ ] Model-View-Update pattern working
- [ ] Services injected correctly

### ✅ Documentation
- [ ] Code comments on all exported functions
- [ ] README updated with feature flag
- [ ] CONTRIBUTING.md updated with architecture
- [ ] Examples in action system

---

## Blockers & Escalation

### If Stuck
1. Review master plan (`TUI_MIGRATION_MASTER_PLAN.md`)
2. Check decision matrix (`MIGRATION_DECISION_MATRIX.md`)
3. Ask in #bubbletea-migration Slack channel
4. Escalate to Director if blocked >4 hours

### Known Risks
- **Learning curve:** Elm Architecture is new paradigm
  - Mitigation: Pair with someone familiar, read official tutorials
- **Integration complexity:** Wiring up services
  - Mitigation: Start with mock services, add real ones incrementally

---

## Daily Checklist

**Start of Day:**
- [ ] Pull latest from main
- [ ] Review day's tasks
- [ ] Estimate completion time
- [ ] Identify dependencies

**During Day:**
- [ ] Commit frequently (every feature)
- [ ] Write tests alongside code
- [ ] Update this checklist
- [ ] Ask questions early

**End of Day:**
- [ ] Push commits
- [ ] Update status in Slack
- [ ] Plan next day
- [ ] Celebrate progress! 🎉

---

## Success! 🚀

When Phase 0 is complete:
- [ ] All acceptance criteria met
- [ ] Code reviewed by peer
- [ ] Merged to main
- [ ] Phase 1 checklist created
- [ ] Team notified of progress

**Next:** Phase 1 - Core Views (Workspace List, Ticket Tree, Ticket Detail)
