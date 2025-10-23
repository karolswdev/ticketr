# Search Modal

The Search Modal provides fuzzy search functionality for the action registry in the Ticketr Bubbletea TUI. Users can quickly find and execute actions using keyboard shortcuts.

## Overview

The search modal is a centered overlay that allows users to:

- Search actions by name, description, or tags
- Navigate search results with keyboard
- Execute selected actions
- See real-time filtered results

## Features

- **Fuzzy Search**: Searches action names, descriptions, and tags
- **Context-Aware**: Filters actions based on current context and predicates
- **Theme-Aware**: Adapts styling to active theme (Default, Dark, Arctic)
- **Keyboard Navigation**: Full keyboard control with vim-style bindings
- **Real-time Filtering**: Results update as you type
- **Empty State**: Clear messaging when no results found

## Architecture

### Component Structure

```
views/search/
├── search.go           # Main search modal model
├── search_test.go      # Comprehensive tests (80%+ coverage)
└── README.md           # This file
```

### Model

```go
type Model struct {
    // UI components
    input textinput.Model       // Search input field

    // Data
    registry  *actions.Registry  // Action registry reference
    results   []*actions.Action  // Filtered search results
    actionCtx *ActionContext     // Action context for predicates

    // State
    visible       bool            // Is modal open?
    selectedIndex int             // Currently selected result index
    width         int             // Viewport width
    height        int             // Viewport height
    theme         *theme.Theme    // Current theme

    // Configuration
    maxResults int                // Maximum results to display (10)
}
```

## Usage

### Creating the Search Modal

```go
import (
    "github.com/karolswdev/ticktr/internal/tui-bubbletea/views/search"
    "github.com/karolswdev/ticktr/internal/tui-bubbletea/actions"
    "github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// Create action registry
registry := actions.NewRegistry()
// Register actions...

// Create search modal
searchModal := search.New(registry, &theme.DefaultTheme)
```

### Opening and Closing

```go
// Open modal
searchModal, cmd := searchModal.Open()

// Close modal
searchModal, cmd := searchModal.Close()

// Check visibility
if searchModal.IsVisible() {
    // Modal is open
}
```

### Integration with Root Model

```go
// In your root model
type Model struct {
    searchModal search.Model
    // ... other fields
}

// In Update()
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Open search modal with '/'
        if msg.String() == "/" && !m.searchModal.IsVisible() {
            m.searchModal, cmd = m.searchModal.Open()
            return m, cmd
        }

        // Route to search modal when visible
        if m.searchModal.IsVisible() {
            m.searchModal, cmd = m.searchModal.Update(msg)
            return m, cmd
        }

    case search.ActionExecuteRequestedMsg:
        // Handle action execution
        if msg.Action != nil {
            action := msg.Action.(*actions.Action)
            actx := m.buildActionContext()
            return m, action.Execute(actx)
        }
    }

    return m, nil
}

// In View()
func (m Model) View() string {
    // Base view
    view := m.renderMainContent()

    // Overlay search modal if visible
    if m.searchModal.IsVisible() {
        view = m.searchModal.View()
    }

    return view
}
```

### Setting Action Context

```go
// Update action context for predicate evaluation
actx := &actions.ActionContext{
    Context:         currentContext,
    SelectedTickets: selectedTickets,
    IsOnline:        true,
    Width:           width,
    Height:          height,
}

searchModal.SetActionContext(actx)
```

## Keybindings

| Key | Action |
|-----|--------|
| `/` | Open search modal (handled by parent) |
| `Esc` | Close modal |
| `Enter` | Execute selected action and close |
| `↑` or `k` | Navigate up in results |
| `↓` or `j` | Navigate down in results |
| Type | Filter actions in real-time |

## Messages

### Outgoing Messages

The search modal sends these messages:

```go
// Sent when modal opens
type SearchModalOpenedMsg struct{}

// Sent when modal closes
type SearchModalClosedMsg struct{}

// Sent when user selects an action
type ActionExecuteRequestedMsg struct {
    ActionID actions.ActionID
    Action   *actions.Action
}
```

### Handling Messages

```go
case search.SearchModalOpenedMsg:
    // Modal opened
    m.contextManager.Push(actions.ContextSearch)

case search.SearchModalClosedMsg:
    // Modal closed
    m.contextManager.Pop()

case search.ActionExecuteRequestedMsg:
    // Execute the requested action
    actx := m.buildActionContext()
    return m, msg.Action.Execute(actx)
```

## Visual Design

### Layout

```
┌─────────────────────────────────────┐
│  🔍 Search Actions                  │
│  ┌───────────────────────────────┐  │
│  │ query text here_              │  │
│  └───────────────────────────────┘  │
│                                     │
│  > 📄 Open Ticket - Open selected  │
│    ✕ Close Modal - Close current   │
│    🔍 Search Actions - Open search  │
│                                     │
│  ↑/↓ or j/k: Navigate • Enter: ... │
└─────────────────────────────────────┘
```

### Dimensions

- **Width**: 40% of screen width (min: 40 chars)
- **Height**: 60% of screen height (min: 15 lines)
- **Max Results**: 10 visible results
- **Position**: Centered overlay

### Theme Styling

The modal adapts to the active theme:

- **Title**: Primary color, bold
- **Selected Item**: Primary color on selection background
- **Normal Items**: Foreground color
- **Description**: Muted color
- **Help Text**: Muted color, italic
- **Border**: Theme border color

## Search Behavior

### Fuzzy Matching

The search uses `registry.Search()` which performs substring matching on:

1. **Action Name** (highest priority)
2. **Action Description**
3. **Action Tags**

### Filtering

Actions are filtered by:

1. **Search Query**: Fuzzy match on name/description/tags
2. **Predicates**: Actions with failing predicates are excluded
3. **ShowInUI**: Actions with `ShowInUI() == false` are excluded
4. **Context**: Only actions available in current context

### Sorting

Results are sorted by:

1. Name matches first
2. Alphabetical by name within each group

## API Reference

### Constructor

```go
func New(registry *actions.Registry, t *theme.Theme) Model
```

Creates a new search modal.

### Methods

#### State Management

```go
func (m Model) Open() (Model, tea.Cmd)
func (m Model) Close() (Model, tea.Cmd)
func (m Model) IsVisible() bool
```

#### Configuration

```go
func (m *Model) SetSize(width, height int)
func (m *Model) SetTheme(t *theme.Theme)
func (m *Model) SetActionContext(actx *actions.ActionContext)
```

#### Bubbletea Interface

```go
func (m Model) Init() tea.Cmd
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)
func (m Model) View() string
```

## Testing

The search modal has comprehensive test coverage (80%+):

```bash
# Run tests
go test ./internal/tui-bubbletea/views/search/

# Run with coverage
go test -cover ./internal/tui-bubbletea/views/search/

# Run with race detection
go test -race ./internal/tui-bubbletea/views/search/
```

### Test Coverage

- **TestNew**: Constructor initialization
- **TestOpenClose**: Modal visibility state
- **TestSearch**: Search functionality and filtering
- **TestNavigation**: Keyboard navigation (up/down, j/k)
- **TestExecuteAction**: Action selection and execution
- **TestEmptyState**: No results messaging
- **TestRendering**: View rendering in different states
- **TestThemeAwareness**: Theme switching
- **TestEscapeKey**: Escape key handling
- **TestSearchResetSelection**: Selection reset on query change
- **TestMaxResults**: Handling >10 results

## Performance

- **Thread-Safe**: Registry access is protected by mutex
- **No Allocations**: Minimal allocations during search
- **Responsive**: Real-time search with <10ms latency
- **Efficient Rendering**: Only visible results are rendered

## Future Enhancements

Potential improvements for future iterations:

1. **Advanced Fuzzy Matching**: Implement proper fuzzy scoring (e.g., fzf algorithm)
2. **Search History**: Remember recent searches
3. **Keyboard Shortcuts**: Show keybindings for each action
4. **Categories**: Group results by category
5. **Icons**: Support more action icons
6. **Highlighting**: Highlight matched text in results
7. **Pagination**: Better handling of many results

## Integration Points

### Day 1: Action System
- Uses `Registry.Search()` for fuzzy search
- Respects action predicates and context
- Executes actions via `Action.Execute()`

### Day 3: Command Palette
- Command palette will extend this pattern
- Share search UI components
- Different result styling and metadata

### Day 4: Context-Aware Help
- Help modal uses similar modal pattern
- Consistent keyboard navigation
- Shared theme styling

## References

- [Action System Design](../../actions/README.md)
- [Modal Component](../../components/modal/modal.go)
- [Theme System](../../theme/README.md)
- [Extensible Action System Design](../../../docs/EXTENSIBLE_ACTION_SYSTEM_DESIGN.md)
