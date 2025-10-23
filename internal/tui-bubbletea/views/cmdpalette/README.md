# Command Palette

The Command Palette provides quick, keyboard-driven access to all available actions in Ticketr with enhanced features like category filtering, recent actions, and keybinding hints.

## Overview

The Command Palette is an enhanced version of the Search Modal (Day 2), offering:

- **Categorized Display**: Actions grouped by category with visual headers
- **Recent Actions**: Quick access to last 5 executed actions (⭐)
- **Keybinding Hints**: See keyboard shortcuts for each action
- **Category Filtering**: Filter by category using Ctrl+0-7
- **Smart Sorting**: Recent actions first, then by relevance
- **Larger Modal**: 60% width, 70% height (vs 40%/60% for search)
- **Theme-Aware**: Supports all 3 themes (Default, Dark, Arctic)

## Architecture

```
cmdpalette/
├── cmdpalette.go       # Main command palette model
├── cmdpalette_test.go  # Comprehensive tests (85%+ coverage)
└── README.md           # This file
```

## Model Structure

```go
type Model struct {
    // UI components
    input textinput.Model       // Command input field

    // Data
    registry   *actions.Registry         // Action registry reference
    contextMgr *actions.ContextManager   // Context manager
    actionCtx  *actions.ActionContext    // Action context for predicates
    results    []ActionItem              // Filtered and sorted results
    recent     []actions.ActionID        // Recent action IDs (max 5)

    // State
    visible       bool                      // Is palette open?
    filterMode    FilterMode                // All, Category, Recent
    selectedCat   actions.ActionCategory    // Selected category filter
    selectedIndex int                       // Currently selected result index
    width         int                       // Viewport width
    height        int                       // Viewport height
    theme         *theme.Theme              // Current theme

    // Configuration
    maxResults int // Maximum results to display (default: 20)
    maxRecent  int // Maximum recent actions to track (default: 5)
}

type ActionItem struct {
    Action      *actions.Action
    Keybindings string // Formatted: "Ctrl+S", "j, ↓"
    Category    string
    IsRecent    bool   // Show star indicator
}

type FilterMode int
const (
    FilterAll FilterMode = iota
    FilterCategory
    FilterRecent
)
```

## Usage Examples

### Creating Command Palette

```go
import (
    "github.com/karolswdev/ticktr/internal/tui-bubbletea/actions"
    "github.com/karolswdev/ticktr/internal/tui-bubbletea/views/cmdpalette"
    "github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// In your root model
type Model struct {
    registry   *actions.Registry
    contextMgr *actions.ContextManager
    cmdPalette cmdpalette.Model
    // ...
}

// Initialize
func initialModel() Model {
    registry := actions.NewRegistry()
    contextMgr := actions.NewContextManager(actions.ContextGlobal)

    // Register your actions
    registerAllActions(registry)

    return Model{
        registry:   registry,
        contextMgr: contextMgr,
        cmdPalette: cmdpalette.New(registry, contextMgr, &theme.DefaultTheme),
    }
}
```

### Opening and Closing

```go
// In Update()
case tea.KeyMsg:
    // Open command palette with Ctrl+P or :
    if (msg.Type == tea.KeyCtrlP || msg.String() == ":") && !m.cmdPalette.IsVisible() {
        m.cmdPalette, cmd = m.cmdPalette.Open()
        return m, cmd
    }

    // Route to command palette when visible
    if m.cmdPalette.IsVisible() {
        m.cmdPalette, cmd = m.cmdPalette.Update(msg)
        return m, cmd
    }
```

### Handling Command Execution

```go
// In Update()
case cmdpalette.CommandExecutedMsg:
    // Action already added to recent by palette

    // Build action context
    actx := m.buildActionContext()

    // Execute the action
    cmd := msg.Action.Execute(actx)

    return m, cmd

case cmdpalette.CommandPaletteClosedMsg:
    // Palette closed, restore previous focus
    return m, nil
```

### Managing Recent Actions

```go
// Add action to recent (done automatically on execution)
m.cmdPalette.AddRecent(actionID)

// Get recent actions (for persistence)
recent := m.cmdPalette.GetRecentActions()
saveToFile(recent)

// Set recent actions (from persistence)
loaded := loadFromFile()
m.cmdPalette.SetRecentActions(loaded)
```

### Category Filtering

```go
// Filter by category (also accessible via Ctrl+1-7)
m.cmdPalette.SetCategoryFilter(actions.CategoryNavigation)

// Clear filter (also accessible via Ctrl+0)
m.cmdPalette.ClearFilter()

// Check current filter mode
if m.cmdPalette.filterMode == cmdpalette.FilterCategory {
    // Category filter is active
}
```

## Keybindings

### Primary Keybindings

| Key | Action |
|-----|--------|
| `Ctrl+P` or `:` | Open command palette |
| `Esc` | Close palette |
| `Enter` | Execute selected action |
| `↑` or `k` | Navigate up |
| `↓` or `j` | Navigate down |
| Type | Filter actions |

### Category Filter Shortcuts

| Key | Category |
|-----|----------|
| `Ctrl+0` | Show all (reset filter) |
| `Ctrl+1` | Navigation only |
| `Ctrl+2` | View only |
| `Ctrl+3` | Edit only |
| `Ctrl+4` | Workspace only |
| `Ctrl+5` | Sync only |
| `Ctrl+6` | Bulk Operations only |
| `Ctrl+7` | System only |

## Visual Design

### Layout

```
┌───────────────────────────────────────────────────┐
│  🎯 Command Palette                    [Ctrl+P]  │
│  ┌─────────────────────────────────────────────┐  │
│  │ query text here_                           │  │
│  └─────────────────────────────────────────────┘  │
│                                                   │
│  ⭐ RECENT                                        │
│  > Edit ticket                       e           │
│    Save changes                      Ctrl+S      │
│                                                   │
│  ── NAVIGATION ──                                 │
│    Move down                         j, ↓        │
│    Move up                           k, ↑        │
│    View detail                       Enter       │
│                                                   │
│  ── EDIT ──                                       │
│    Create ticket                     n           │
│    Delete ticket                     Del         │
│                                                   │
│  Showing 15 actions in 3 categories              │
│  Ctrl+0-7: Filter by category  |  Esc: Close     │
└───────────────────────────────────────────────────┘
```

### Styling Features

- **Selected Item**: Primary color background, bold text
- **Category Headers**: Dimmed, uppercase, with decorative separators
- **Recent Section**: Star icon (⭐) to highlight recent actions
- **Keybindings**: Right-aligned, dimmed to not distract
- **Description**: Shows under selected item in italic, muted color
- **Modal Size**: 60% width, 70% height (larger than search modal)

### Theme Support

The command palette adapts to all 3 themes:

- **Default Theme**: Green accents, terminal-style
- **Dark Theme**: Blue accents, modern professional
- **Arctic Theme**: Cyan accents, cool and crisp

## Feature Comparison: Search Modal vs Command Palette

| Feature | Search Modal | Command Palette |
|---------|--------------|-----------------|
| Keybinding | `/` | `Ctrl+P` or `:` |
| Modal Size | 40% x 60% | 60% x 70% |
| Layout | Simple list | Categorized with headers |
| Keybinding Hints | No | Yes (right side) |
| Recent Actions | No | Yes (top of list, ⭐) |
| Category Filter | No | Yes (Ctrl+0-7) |
| Action Count | Yes | Yes (with category breakdown) |
| Max Results | 10 | 20 |
| Description | In line | Shows under selection |
| Empty State | Simple | With help hint |

## Messages

### Outgoing Messages

```go
// CommandPaletteOpenedMsg is sent when the palette opens
type CommandPaletteOpenedMsg struct{}

// CommandPaletteClosedMsg is sent when the palette closes
type CommandPaletteClosedMsg struct{}

// CommandExecutedMsg is sent when an action is executed
type CommandExecutedMsg struct {
    ActionID actions.ActionID
    Action   *actions.Action
}
```

## Recent Actions Implementation

### How It Works

1. **Tracking**: When an action is executed, `AddRecent()` is called
2. **Storage**: Recent actions stored in `[]actions.ActionID` (max 5)
3. **Ordering**: Most recent first (LIFO order)
4. **Duplicates**: Moving an existing action to the front
5. **Display**: Recent actions marked with `IsRecent: true` flag
6. **Sorting**: Recent actions appear first in results
7. **Grouping**: Recent section shown at top with ⭐ header

### Persistence (Optional)

```go
// Save recent actions to file
func saveRecent(m *cmdpalette.Model) error {
    recent := m.GetRecentActions()
    data, _ := json.Marshal(recent)
    return os.WriteFile("recent.json", data, 0644)
}

// Load recent actions from file
func loadRecent(m *cmdpalette.Model) error {
    data, err := os.ReadFile("recent.json")
    if err != nil {
        return err
    }
    var recent []actions.ActionID
    json.Unmarshal(data, &recent)
    m.SetRecentActions(recent)
    return nil
}
```

## Category Filtering Implementation

### How It Works

1. **Filter Modes**: Three modes - All, Category, Recent
2. **Category Selection**: `SetCategoryFilter(category)` sets filter
3. **Keyboard Shortcuts**: Ctrl+1-7 mapped to categories
4. **Search Integration**: Filter applied after search
5. **Visual Feedback**: Category name shown in title bar
6. **Reset**: Ctrl+0 or `ClearFilter()` resets to All

### Category Order

Actions are grouped in this order:
1. Navigation
2. View
3. Edit
4. Workspace
5. Sync
6. Bulk Operations
7. System
8. Other (fallback)

## Smart Sorting

### Sort Priority

1. **Recent First**: Actions in recent list appear first
2. **Name Match**: Actions with query in name rank higher
3. **Alphabetical**: Within same priority, sort alphabetically
4. **Category Order**: When grouped, follow category order

### Example

Query: "e"
- ⭐ Edit ticket (recent)
- ⭐ Save changes (recent)
- Create ticket (name starts with 'e' in Create)
- Delete ticket
- Refresh Data
- ...

## API Reference

### Core Methods

```go
// Open shows the command palette
func (m Model) Open() (Model, tea.Cmd)

// Close hides the command palette
func (m Model) Close() (Model, tea.Cmd)

// IsVisible returns whether palette is visible
func (m Model) IsVisible() bool

// Update handles messages
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd)

// View renders the palette
func (m Model) View() string
```

### Configuration Methods

```go
// SetSize updates palette dimensions
func (m *Model) SetSize(width, height int)

// SetTheme updates the theme
func (m *Model) SetTheme(t *theme.Theme)

// SetActionContext updates action context
func (m *Model) SetActionContext(actx *actions.ActionContext)
```

### Recent Actions Methods

```go
// AddRecent adds an action to recent list
func (m *Model) AddRecent(actionID actions.ActionID)

// GetRecentActions returns recent action IDs
func (m Model) GetRecentActions() []actions.ActionID

// SetRecentActions sets recent actions (for persistence)
func (m *Model) SetRecentActions(recent []actions.ActionID)
```

### Filtering Methods

```go
// SetCategoryFilter filters by category
func (m *Model) SetCategoryFilter(category actions.ActionCategory)

// ClearFilter clears category filter
func (m *Model) ClearFilter()
```

## Testing

### Test Coverage: 85%+

The command palette has comprehensive test coverage:

```bash
go test ./internal/tui-bubbletea/views/cmdpalette -v -cover
```

### Test Categories

- **Initialization**: New, Init
- **Open/Close**: Open, Close, IsVisible
- **Search**: Empty query, text search, tag search, no results
- **Navigation**: j/k, arrows, bounds checking
- **Execution**: Enter key, action execution, recent tracking
- **Category Filtering**: Ctrl+0-7, SetCategoryFilter, ClearFilter
- **Recent Actions**: AddRecent, GetRecentActions, SetRecentActions, sorting
- **Keybindings**: Display formatting, all key patterns
- **Rendering**: Theme awareness, category headers, empty state
- **Edge Cases**: No results, not visible, window resize

### Running Tests

```bash
# Unit tests
go test ./internal/tui-bubbletea/views/cmdpalette

# With coverage
go test ./internal/tui-bubbletea/views/cmdpalette -cover

# With race detection
go test ./internal/tui-bubbletea/views/cmdpalette -race

# Verbose output
go test ./internal/tui-bubbletea/views/cmdpalette -v
```

## Integration Points

### Week 3 Day 2: Search Modal

The command palette reuses patterns from the Search Modal:
- Modal overlay system
- textinput.Model for search
- Action registry integration
- Theme-aware styling
- Test patterns

### Week 3 Day 4: Context-Aware Help

The command palette complements the help system:
- Help can list all keybindings
- Command palette provides quick access
- Both use action registry
- Ctrl+H opens help from palette

### Week 3 Day 5: Keybinding Resolver

Future integration:
- Resolver will map keys to actions
- Command palette shows resolved keybindings
- User keybinding overrides respected
- Conflict detection and resolution

## Performance

### Optimizations

- **Lazy Search**: Only searches on input change
- **Result Limiting**: Max 20 displayed (configurable)
- **Category Caching**: Category groups computed once per search
- **Theme Styles**: Styles created once, reused
- **String Formatting**: Keybindings formatted once per action

### Benchmarks

Expected performance:
- Open: < 1ms
- Search (100 actions): < 5ms
- Render: < 10ms
- Category filter: < 2ms

## Future Enhancements

### Week 4+

1. **Fuzzy Matching**: Better search algorithm (e.g., fzf-style)
2. **Action Preview**: Show more details in side panel
3. **Action History**: Track usage statistics
4. **Favorites**: Pin frequently used actions
5. **Action Aliases**: Custom short names
6. **Macro Recording**: Combine multiple actions
7. **Plugin Actions**: Register from Lua plugins
8. **Custom Categories**: User-defined grouping
9. **Search Syntax**: `@category:navigation query`
10. **Persistence**: Save recent/favorites to disk

## Troubleshooting

### Command Palette Not Opening

- Check if keybinding is registered (Ctrl+P or :)
- Verify registry is not nil
- Check if another modal is open

### No Actions Showing

- Verify actions are registered in registry
- Check predicate functions (might be filtering out)
- Ensure action context is set correctly

### Keybindings Not Displaying

- Check if action has Keybindings field set
- Verify formatKeyPattern() is working
- Some actions might not have keybindings

### Category Filter Not Working

- Verify action has Category field set
- Check if Ctrl+1-7 keys are being captured
- Try calling SetCategoryFilter() directly

### Recent Actions Not Persisting

- Recent actions are in-memory only by default
- Implement persistence using GetRecentActions/SetRecentActions
- Save to file on quit, load on startup

## Design Decisions

### Why Larger Modal?

The command palette shows more information (categories, keybindings, descriptions), so it needs more space than the search modal.

### Why Recent Actions?

Users often repeat the same actions. Recent actions provide quick access to common operations without searching.

### Why Category Filtering?

When you know the category, filtering reduces cognitive load. Ctrl+1-7 provides fast keyboard access to categories.

### Why Show Keybindings?

Discoverability: Users learn keyboard shortcuts by seeing them in the palette.

### Why Smart Sorting?

Combining recent + relevance + alphabetical provides the best UX: common actions first, then relevant matches, then alphabetical within groups.

## Contributing

When adding features to the command palette:

1. **Maintain Patterns**: Follow Elm Architecture
2. **Test Coverage**: Keep above 85%
3. **Documentation**: Update this README
4. **Theme Support**: Test all 3 themes
5. **Performance**: Profile with large action sets
6. **Accessibility**: Consider keyboard-only users

## References

- **Action System**: `internal/tui-bubbletea/actions/`
- **Search Modal**: `internal/tui-bubbletea/views/search/`
- **Help System**: `internal/tui-bubbletea/components/help/`
- **Theme System**: `internal/tui-bubbletea/theme/`
- **Design Doc**: `docs/EXTENSIBLE_ACTION_SYSTEM_DESIGN.md`
- **Wireframes**: `docs/TUI_WIREFRAMES.md`
