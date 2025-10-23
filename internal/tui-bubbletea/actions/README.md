# Action System

The Action System provides a declarative, extensible framework for defining and managing application actions in the Ticketr Bubbletea TUI.

## Overview

The action system is built on these core principles:

1. **Actions are Data, Not Code** - Actions are declarative data structures
2. **Context Determines Availability** - Actions only show when relevant
3. **Predicates Control Visibility** - Conditions determine when actions are available
4. **Execution Returns Messages** - Bubbletea tea.Cmd pattern
5. **Composition Over Inheritance** - Actions can be composed

## Architecture

```
actions/
├── action.go           # Core types (Action, ActionContext, KeyPattern)
├── context.go          # Context manager
├── registry.go         # Action registry
├── predicates/         # Common predicates
│   └── predicates.go   # Always, Never, HasSelection, etc.
└── builtin/            # Built-in actions
    └── system.go       # System actions (quit, help)
```

## Core Concepts

### Action

An `Action` represents a single executable operation:

```go
action := &Action{
    ID:          "ticket.open",
    Name:        "Open Ticket",
    Description: "Open selected ticket in detail view",
    Category:    CategoryView,
    Contexts:    []Context{ContextTicketTree},
    Keybindings: []KeyPattern{{Key: "enter"}},
    Predicate:   predicates.HasSingleSelection(),
    Execute: func(ctx *ActionContext) tea.Cmd {
        // Return a tea.Cmd
        return func() tea.Msg {
            return ticketOpenedMsg{id: ctx.SelectedTickets[0]}
        }
    },
    Tags: []string{"ticket", "open", "view"},
    Icon: "📄",
}
```

### Context

A `Context` represents where the user is in the application:

- `ContextWorkspaceList` - Workspace selector view
- `ContextTicketTree` - Ticket tree view
- `ContextTicketDetail` - Ticket detail view
- `ContextSearch` - Search modal
- `ContextGlobal` - Matches any context

### Predicates

Predicates determine if an action is available:

```go
// Simple predicate
action.Predicate = predicates.HasSingleSelection()

// Composed predicate
action.Predicate = predicates.And(
    predicates.HasSelection(),
    predicates.IsOnline(),
)

// Complex composition
action.Predicate = predicates.Or(
    predicates.And(
        predicates.HasSingleSelection(),
        predicates.IsOnline(),
    ),
    predicates.HasUnsavedChanges(),
)
```

### Registry

The `Registry` manages all actions:

```go
// Create registry
reg := NewRegistry()

// Register action
reg.Register(action)

// Get actions for context
actions := reg.ActionsForContext(ContextTicketTree, actx)

// Search actions
results := reg.Search("open", actx)

// Get actions for keybinding
actions := reg.ActionsForKey("enter", ContextTicketTree, actx)
```

## Usage Examples

### Registering Built-in Actions

```go
func RegisterTicketActions(reg *Registry) error {
    // Open ticket
    if err := reg.Register(&Action{
        ID:          "ticket.open",
        Name:        "Open Ticket",
        Description: "Open selected ticket in detail view",
        Category:    CategoryView,
        Contexts:    []Context{ContextTicketTree},
        Keybindings: []KeyPattern{{Key: "enter"}, {Key: "o"}},
        Predicate:   predicates.HasSingleSelection(),
        Execute:     executeOpenTicket,
        Tags:        []string{"ticket", "open", "view"},
        Icon:        "📄",
    }); err != nil {
        return err
    }

    return nil
}
```

### Building Action Context

```go
func (m Model) buildActionContext() *ActionContext {
    return &ActionContext{
        Context:           m.contextManager.Current(),
        SelectedTickets:   m.selectedTickets,
        SelectedWorkspace: m.currentWorkspace,
        HasUnsavedChanges: m.ticketDetail.IsDirty(),
        IsSyncing:         m.syncInProgress,
        IsOffline:         m.isOffline,
        Config:            m.userConfig,
        Services:          m.services,
        Width:             m.width,
        Height:            m.height,
    }
}
```

### Executing Actions

```go
// In Model.Update()
case tea.KeyMsg:
    actx := m.buildActionContext()

    // Get available actions for current context
    actions := m.actionRegistry.ActionsForContext(
        m.contextManager.Current(),
        actx,
    )

    // Execute action
    if action != nil {
        cmd := action.Execute(actx)
        return m, cmd
    }
```

## Predicate Reference

### Basic Predicates

- `Always()` - Always returns true
- `Never()` - Always returns false
- `HasSelection()` - At least one ticket selected
- `HasSingleSelection()` - Exactly one ticket selected
- `HasMultipleSelection()` - More than one ticket selected
- `IsInWorkspace()` - A workspace is selected
- `HasUnsavedChanges()` - There are unsaved changes
- `IsOnline()` - Not in offline mode

### Composition

- `Not(pred)` - Inverts predicate
- `And(preds...)` - All predicates must be true
- `Or(preds...)` - At least one predicate must be true

## Context Manager

The `ContextManager` tracks application context:

```go
// Create manager
cm := NewContextManager(ContextWorkspaceList)

// Switch context
cm.Switch(ContextTicketTree)

// Push context (for modals)
cm.Push(ContextModal)

// Pop context
previous := cm.Pop()

// Get current context
current := cm.Current()

// Register observer
cm.OnChange(func(old, new Context) {
    fmt.Printf("Context changed: %s -> %s\n", old, new)
})
```

## Testing

The action system is fully testable:

```go
func TestActionWithPredicate(t *testing.T) {
    reg := NewRegistry()

    action := &Action{
        ID:        "test.action",
        Predicate: predicates.HasSingleSelection(),
        Execute:   func(ctx *ActionContext) tea.Cmd { return nil },
    }

    reg.Register(action)

    // Test without selection
    actx := &ActionContext{SelectedTickets: []string{}}
    actions := reg.ActionsForContext(ContextTicketTree, actx)
    assert.Equal(t, 0, len(actions))

    // Test with selection
    actx.SelectedTickets = []string{"TICKET-1"}
    actions = reg.ActionsForContext(ContextTicketTree, actx)
    assert.Equal(t, 1, len(actions))
}
```

## Integration Points

### Day 2: Search Modal

The search modal will use `Registry.Search()`:

```go
results := registry.Search(query, actx)
// Display results in list
// Execute selected action
```

### Day 3: Command Palette

The command palette will use `Registry.All()`:

```go
allActions := registry.All()
// Filter by query
// Display with categories
// Execute selected action
```

### Future: Lua Plugins

The action system is designed to support plugin-registered actions:

```go
// Plugin can register custom actions
pluginAction := &Action{
    ID:          "plugin.my_action",
    Name:        "My Custom Action",
    Execute:     pluginExecuteFunc,
}
registry.Register(pluginAction)
```

## Performance

The registry is thread-safe and optimized for lookup:

- Actions indexed by ID (O(1) lookup)
- Actions indexed by context (O(1) lookup)
- Actions indexed by keybinding (O(1) lookup)
- Search uses simple substring matching (can be upgraded to fuzzy)

## Design Decisions

### Why Predicates?

Predicates provide:
- **Declarative** - Easy to understand
- **Composable** - Build complex conditions
- **Testable** - Pure functions
- **Reusable** - Share across actions

### Why Registry?

The registry provides:
- **Centralized** - All actions in one place
- **Discoverable** - Easy to list/search
- **Extensible** - Plugins can register
- **Type-safe** - Compile-time checking

### Why Context Manager?

Context tracking enables:
- **Context-aware UI** - Show relevant actions
- **Modal management** - Push/pop contexts
- **State isolation** - Context-specific metadata
- **Observer pattern** - React to changes

## Future Enhancements

Potential future features:

1. **Action Middleware** - Pre/post execution hooks
2. **Action Sequences** - Chain multiple actions
3. **Keybinding Resolver** - Map keys to actions
4. **User-defined Actions** - Load from config
5. **Action History** - Undo/redo support
6. **Lua Plugin Support** - Register actions from Lua

## References

- EXTENSIBLE_ACTION_SYSTEM_DESIGN.md - Full design document
- CLEAN_SLATE_REFACTOR_MASTER_PLAN.md - Week 3 implementation plan
