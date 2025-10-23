# Extensible Action System Design for Ticketr Bubbletea TUI

**Date:** 2025-10-22
**Author:** Claude (Architecture Design)
**Purpose:** Design a flexible, context-aware action system for ticketr's Bubbletea TUI with Lua plugin readiness

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Design Principles](#design-principles)
3. [Architecture Overview](#architecture-overview)
4. [Core Type Definitions](#core-type-definitions)
5. [Context System](#context-system)
6. [Predicate System](#predicate-system)
7. [Action Registry](#action-registry)
8. [Keybinding System](#keybinding-system)
9. [Action Execution Pipeline](#action-execution-pipeline)
10. [Bubbletea Integration](#bubbletea-integration)
11. [Extension Points](#extension-points)
12. [Code Examples](#code-examples)
13. [Migration Path](#migration-path)
14. [Future: Lua Plugin Architecture](#future-lua-plugin-architecture)

---

## Executive Summary

### The Problem

The current Tview-based TUI has:
- **Hardcoded actions** scattered across views (ticket_tree.go, ticket_detail.go, etc.)
- **Manual keybinding management** in ActionBar with context switching
- **No extensibility** - adding new actions requires modifying core code
- **Callback-based execution** - difficult to test and compose

### The Solution

A **declarative action system** where:
- Actions are **data structures** (not callbacks scattered everywhere)
- Actions have **predicates** (only show when condition is met)
- Actions are **context-aware** (workspace list vs ticket tree vs ticket detail)
- Actions can be **dynamically registered** (by plugins or at runtime)
- Keybindings are **customizable** via configuration
- Actions are **composable** (chains, sequences, conditions)

### Key Benefits

1. ✅ **Midnight Commander-inspired** - F-keys and context menus like mc
2. ✅ **Plugin-ready** - Lua scripts can register custom actions
3. ✅ **Testable** - Pure functions, no side effects in action definitions
4. ✅ **Discoverable** - Command palette shows all available actions
5. ✅ **Maintainable** - Actions defined in one place, not scattered
6. ✅ **Bubbletea-native** - Uses message passing, not callbacks

---

## Design Principles

### 1. Actions are Data, Not Code

```go
// BAD: Current approach (callback-based)
func (t *TicketTreeView) handleKeyPress(key tcell.Key) {
    if key == tcell.KeyEnter {
        t.openTicket() // Direct function call
    }
}

// GOOD: New approach (data-driven)
var OpenTicketAction = Action{
    ID:          "ticket.open",
    Name:        "Open Ticket",
    Description: "Open selected ticket in detail view",
    Keybindings: []KeyPattern{{Key: "enter"}},
    Execute:     openTicketHandler,
}
```

### 2. Context Determines Available Actions

```go
// Workspace list shows different actions than ticket tree
workspaceActions := registry.ActionsForContext(ContextWorkspaceList)
ticketActions := registry.ActionsForContext(ContextTicketTree)
```

### 3. Predicates Control Visibility

```go
// Action only shows when exactly one ticket is selected
action.Predicate = func(ctx *ActionContext) bool {
    return len(ctx.SelectedTickets) == 1
}
```

### 4. Execution Returns Messages (Bubbletea Pattern)

```go
// Execute returns a tea.Cmd (Bubbletea message)
func (a *Action) Execute(ctx *ActionContext) tea.Cmd {
    return func() tea.Msg {
        return ticketOpenedMsg{id: ctx.SelectedTickets[0]}
    }
}
```

### 5. Composition Over Inheritance

```go
// Actions can be composed into sequences, choices, or conditions
SequenceAction(saveAction, closeAction)    // Save then close
ConditionalAction(isDirty, confirmAction)  // Only if dirty
ParallelAction(saveAction, syncAction)     // Save and sync concurrently
```

---

## Architecture Overview

### System Components

```
┌─────────────────────────────────────────────────────────────────┐
│                    Bubbletea Model (Root)                       │
│                                                                 │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────────┐    │
│  │   Context   │  │    Action    │  │    Keybinding      │    │
│  │   Manager   │  │   Registry   │  │     Resolver       │    │
│  └─────────────┘  └──────────────┘  └────────────────────┘    │
│         │                │                      │               │
│         └────────────────┼──────────────────────┘               │
│                          │                                      │
│              ┌───────────▼──────────┐                          │
│              │  Execution Pipeline  │                          │
│              └──────────────────────┘                          │
│                          │                                      │
│         ┌────────────────┼────────────────┐                    │
│         │                │                │                    │
│  ┌──────▼─────┐  ┌──────▼──────┐  ┌─────▼──────┐            │
│  │  Built-in  │  │   Plugin    │  │   User     │            │
│  │  Actions   │  │  Actions    │  │  Actions   │            │
│  └────────────┘  └─────────────┘  └────────────┘            │
└─────────────────────────────────────────────────────────────────┘
                          │
                          ▼
            ┌─────────────────────────┐
            │  Bubbletea Update()     │
            │  (Message Handling)     │
            └─────────────────────────┘
```

### Data Flow

```
User Input (KeyMsg)
        │
        ▼
Keybinding Resolver
        │
        ▼
Find Matching Action (by context + key)
        │
        ▼
Check Predicate (is action available?)
        │
        ▼
Execute Action (returns tea.Cmd)
        │
        ▼
Bubbletea Runtime (dispatches command)
        │
        ▼
Result Message (ticketOpenedMsg, errorMsg, etc.)
        │
        ▼
Model.Update() (handles result)
        │
        ▼
UI Re-renders (Model.View())
```

---

## Core Type Definitions

### 1. Action Definition

```go
// Package: internal/tui/actions

package actions

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/key"
)

// ActionID uniquely identifies an action (plugin-safe namespacing)
type ActionID string

// ActionCategory groups related actions
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

// Action represents a single executable action with metadata
type Action struct {
	// Identity
	ID          ActionID       // e.g., "ticket.open", "workspace.create"
	Name        string         // e.g., "Open Ticket"
	Description string         // e.g., "Open selected ticket in detail view"
	Category    ActionCategory // For grouping in command palette

	// Keybindings (can have multiple for same action)
	Keybindings []KeyPattern // e.g., [{Key: "enter"}, {Key: "o"}]

	// Visibility & Availability
	Contexts  []Context           // Which contexts show this action
	Predicate PredicateFunc       // When is this action available?
	ShowInUI  ShowInUIFunc        // Should this show in command palette/help?

	// Execution
	Execute ExecuteFunc // What to do when action is triggered

	// Metadata (for UI, help, command palette)
	Icon     string            // Optional icon for UI (e.g., "📄", "✓")
	Tags     []string          // For fuzzy search (e.g., ["ticket", "open", "view"])
	Metadata map[string]string // Arbitrary metadata (plugin-extensible)

	// Composition (for complex actions)
	Modifiers []ActionModifier // Pre/post execution hooks
}

// KeyPattern defines a key combination
type KeyPattern struct {
	Key   string // e.g., "enter", "ctrl+s", "F1"
	Alt   bool   // Alt modifier
	Ctrl  bool   // Ctrl modifier
	Shift bool   // Shift modifier

	// Alternative: Use bubbles/key.Binding directly
	Binding key.Binding
}

// PredicateFunc determines if action is available given current context
type PredicateFunc func(ctx *ActionContext) bool

// ShowInUIFunc determines if action should show in UI (command palette, help)
// Allows hiding internal/debugging actions
type ShowInUIFunc func(ctx *ActionContext) bool

// ExecuteFunc performs the action and returns a Bubbletea command
// The command can dispatch messages, perform I/O, etc.
type ExecuteFunc func(ctx *ActionContext) tea.Cmd

// ActionContext provides state to predicates and execute functions
type ActionContext struct {
	// Current view context
	Context Context

	// Selection state
	SelectedTickets   []string           // Ticket IDs
	SelectedWorkspace *WorkspaceState    // Current workspace
	FocusedItem       interface{}        // Currently focused item (generic)

	// Global state (read-only for predicates)
	HasUnsavedChanges bool
	IsSyncing         bool
	IsOffline         bool

	// User preferences (for customization)
	Config *UserConfig

	// Plugin context (for Lua plugins)
	PluginData map[string]interface{}

	// UI state (for rendering)
	Width  int // Terminal width
	Height int // Terminal height

	// Service layer (for execution)
	Services *ServiceContainer // Dependency injection
}

// ActionModifier wraps action execution with pre/post hooks
type ActionModifier interface {
	Before(ctx *ActionContext) error // Called before Execute
	After(ctx *ActionContext, result tea.Cmd, err error) (tea.Cmd, error) // Called after Execute
}
```

### 2. Context Definition

```go
// Context represents where the user is in the application
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

	// Special: matches any context
	ContextGlobal Context = "*"
)

// ContextState tracks current application context
type ContextState struct {
	Current  Context              // Active context
	Previous Context              // For back navigation
	Stack    []Context            // Context history (for nested modals)
	Metadata map[string]interface{} // Context-specific metadata
}

// IsIn checks if current context matches any of the given contexts
func (cs *ContextState) IsIn(contexts ...Context) bool {
	for _, ctx := range contexts {
		if ctx == ContextGlobal || ctx == cs.Current {
			return true
		}
	}
	return false
}
```

### 3. Supporting Types

```go
// WorkspaceState represents current workspace
type WorkspaceState struct {
	ID          string
	Name        string
	ProfileID   string
	IsDirty     bool
	TicketCount int
}

// UserConfig holds user preferences
type UserConfig struct {
	// Keybinding overrides
	Keybindings map[ActionID][]KeyPattern

	// Feature flags
	Features map[string]bool

	// Plugin configuration
	PluginConfig map[string]map[string]interface{}
}

// ServiceContainer provides dependency injection for actions
type ServiceContainer struct {
	TicketService    TicketService
	WorkspaceService WorkspaceService
	SyncService      SyncService
	// ... other services
}
```

---

## Context System

### Context Manager

```go
// Package: internal/tui/actions

// ContextManager tracks and manages application context
type ContextManager struct {
	state    ContextState
	onChange []func(old, new Context) // Observers
}

// NewContextManager creates a new context manager
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

// Switch changes the current context
func (cm *ContextManager) Switch(newContext Context) {
	old := cm.state.Current
	if old == newContext {
		return
	}

	cm.state.Previous = old
	cm.state.Current = newContext
	cm.state.Stack = append(cm.state.Stack, newContext)

	// Notify observers
	for _, fn := range cm.onChange {
		fn(old, newContext)
	}
}

// Push pushes a new context onto the stack (for modals)
func (cm *ContextManager) Push(ctx Context) {
	cm.state.Stack = append(cm.state.Stack, ctx)
	cm.Switch(ctx)
}

// Pop returns to the previous context
func (cm *ContextManager) Pop() Context {
	if len(cm.state.Stack) <= 1 {
		return cm.state.Current
	}

	cm.state.Stack = cm.state.Stack[:len(cm.state.Stack)-1]
	prev := cm.state.Stack[len(cm.state.Stack)-1]
	cm.Switch(prev)
	return prev
}

// OnChange registers a context change observer
func (cm *ContextManager) OnChange(fn func(old, new Context)) {
	cm.onChange = append(cm.onChange, fn)
}

// Current returns the current context
func (cm *ContextManager) Current() Context {
	return cm.state.Current
}

// GetMetadata retrieves context metadata
func (cm *ContextManager) GetMetadata(key string) (interface{}, bool) {
	val, ok := cm.state.Metadata[key]
	return val, ok
}

// SetMetadata stores context metadata
func (cm *ContextManager) SetMetadata(key string, value interface{}) {
	cm.state.Metadata[key] = value
}
```

---

## Predicate System

### Common Predicates

```go
// Package: internal/tui/actions/predicates

package predicates

import "github.com/yourorg/ticketr/internal/tui/actions"

// Always returns a predicate that always returns true
func Always() actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return true
	}
}

// Never returns a predicate that always returns false
func Never() actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return false
	}
}

// HasSelection returns true if at least one ticket is selected
func HasSelection() actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return len(ctx.SelectedTickets) > 0
	}
}

// HasSingleSelection returns true if exactly one ticket is selected
func HasSingleSelection() actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return len(ctx.SelectedTickets) == 1
	}
}

// HasMultipleSelection returns true if more than one ticket is selected
func HasMultipleSelection() actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return len(ctx.SelectedTickets) > 1
	}
}

// IsInWorkspace returns true if a workspace is selected
func IsInWorkspace() actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return ctx.SelectedWorkspace != nil
	}
}

// HasUnsavedChanges returns true if there are unsaved changes
func HasUnsavedChanges() actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return ctx.HasUnsavedChanges
	}
}

// IsOnline returns true if not in offline mode
func IsOnline() actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return !ctx.IsOffline
	}
}

// Not inverts a predicate
func Not(pred actions.PredicateFunc) actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		return !pred(ctx)
	}
}

// And combines predicates with logical AND
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

// Or combines predicates with logical OR
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

// Custom predicate example
func TicketHasStatus(status string) actions.PredicateFunc {
	return func(ctx *actions.ActionContext) bool {
		if len(ctx.SelectedTickets) != 1 {
			return false
		}
		// Get ticket from service
		ticket, err := ctx.Services.TicketService.Get(ctx.SelectedTickets[0])
		if err != nil {
			return false
		}
		return ticket.Status == status
	}
}
```

### Predicate Composition Example

```go
// Only show "Reopen Ticket" action if:
// - Exactly one ticket is selected
// - Ticket is closed
// - User is online
var ReopenTicketAction = Action{
	ID:          "ticket.reopen",
	Name:        "Reopen Ticket",
	Predicate: And(
		HasSingleSelection(),
		TicketHasStatus("Closed"),
		IsOnline(),
	),
	// ...
}
```

---

## Action Registry

### Registry Implementation

```go
// Package: internal/tui/actions

// Registry manages all registered actions
type Registry struct {
	actions   map[ActionID]*Action
	byContext map[Context][]*Action
	byKey     map[string][]*Action // Key pattern string -> actions
	mutex     sync.RWMutex
}

// NewRegistry creates a new action registry
func NewRegistry() *Registry {
	return &Registry{
		actions:   make(map[ActionID]*Action),
		byContext: make(map[Context][]*Action),
		byKey:     make(map[string][]*Action),
	}
}

// Register adds an action to the registry
func (r *Registry) Register(action *Action) error {
	if action.ID == "" {
		return fmt.Errorf("action ID is required")
	}
	if action.Execute == nil {
		return fmt.Errorf("action execute function is required")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check for duplicate ID
	if _, exists := r.actions[action.ID]; exists {
		return fmt.Errorf("action ID already registered: %s", action.ID)
	}

	// Store action
	r.actions[action.ID] = action

	// Index by context
	for _, ctx := range action.Contexts {
		r.byContext[ctx] = append(r.byContext[ctx], action)
	}

	// Index by keybinding
	for _, keyPattern := range action.Keybindings {
		keyStr := keyPattern.String()
		r.byKey[keyStr] = append(r.byKey[keyStr], action)
	}

	return nil
}

// Unregister removes an action from the registry (for plugins)
func (r *Registry) Unregister(id ActionID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	action, exists := r.actions[id]
	if !exists {
		return fmt.Errorf("action not found: %s", id)
	}

	// Remove from main map
	delete(r.actions, id)

	// Remove from context index
	for _, ctx := range action.Contexts {
		r.removeFromSlice(r.byContext[ctx], action)
	}

	// Remove from key index
	for _, keyPattern := range action.Keybindings {
		keyStr := keyPattern.String()
		r.removeFromSlice(r.byKey[keyStr], action)
	}

	return nil
}

// Get retrieves an action by ID
func (r *Registry) Get(id ActionID) (*Action, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	action, exists := r.actions[id]
	return action, exists
}

// ActionsForContext returns all actions available in a context
func (r *Registry) ActionsForContext(ctx Context, actx *ActionContext) []*Action {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Get actions for this context
	contextActions := r.byContext[ctx]
	globalActions := r.byContext[ContextGlobal]

	// Combine and filter by predicate
	var available []*Action
	for _, action := range append(contextActions, globalActions...) {
		if action.Predicate == nil || action.Predicate(actx) {
			available = append(available, action)
		}
	}

	return available
}

// ActionsForKey returns actions bound to a specific key
func (r *Registry) ActionsForKey(key string, ctx Context, actx *ActionContext) []*Action {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Get all actions for this key
	keyActions := r.byKey[key]

	// Filter by context and predicate
	var available []*Action
	for _, action := range keyActions {
		// Check context
		if !r.actionMatchesContext(action, ctx) {
			continue
		}

		// Check predicate
		if action.Predicate != nil && !action.Predicate(actx) {
			continue
		}

		available = append(available, action)
	}

	return available
}

// All returns all registered actions
func (r *Registry) All() []*Action {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	actions := make([]*Action, 0, len(r.actions))
	for _, action := range r.actions {
		actions = append(actions, action)
	}

	return actions
}

// Search performs fuzzy search on actions
func (r *Registry) Search(query string, actx *ActionContext) []*Action {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	query = strings.ToLower(query)
	var results []*Action

	for _, action := range r.actions {
		// Check if action is available
		if action.Predicate != nil && !action.Predicate(actx) {
			continue
		}

		// Check if should show in UI
		if action.ShowInUI != nil && !action.ShowInUI(actx) {
			continue
		}

		// Fuzzy match on name, description, tags
		name := strings.ToLower(action.Name)
		desc := strings.ToLower(action.Description)

		if strings.Contains(name, query) || strings.Contains(desc, query) {
			results = append(results, action)
			continue
		}

		// Check tags
		for _, tag := range action.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				results = append(results, action)
				break
			}
		}
	}

	// Sort by relevance (name matches first)
	sort.Slice(results, func(i, j int) bool {
		iNameMatch := strings.Contains(strings.ToLower(results[i].Name), query)
		jNameMatch := strings.Contains(strings.ToLower(results[j].Name), query)

		if iNameMatch != jNameMatch {
			return iNameMatch
		}

		return results[i].Name < results[j].Name
	})

	return results
}

// Helper: Check if action applies to context
func (r *Registry) actionMatchesContext(action *Action, ctx Context) bool {
	for _, actCtx := range action.Contexts {
		if actCtx == ContextGlobal || actCtx == ctx {
			return true
		}
	}
	return false
}

// Helper: Remove action from slice
func (r *Registry) removeFromSlice(slice []*Action, action *Action) {
	for i, a := range slice {
		if a.ID == action.ID {
			copy(slice[i:], slice[i+1:])
			slice = slice[:len(slice)-1]
			break
		}
	}
}

// KeyPattern.String converts key pattern to string for indexing
func (kp KeyPattern) String() string {
	var parts []string
	if kp.Ctrl {
		parts = append(parts, "ctrl")
	}
	if kp.Alt {
		parts = append(parts, "alt")
	}
	if kp.Shift {
		parts = append(parts, "shift")
	}
	parts = append(parts, strings.ToLower(kp.Key))
	return strings.Join(parts, "+")
}
```

---

## Keybinding System

### Keybinding Resolver

```go
// Package: internal/tui/actions

// KeybindingResolver maps key events to actions
type KeybindingResolver struct {
	registry       *Registry
	contextManager *ContextManager
	userConfig     *UserConfig
}

// NewKeybindingResolver creates a new keybinding resolver
func NewKeybindingResolver(reg *Registry, cm *ContextManager, cfg *UserConfig) *KeybindingResolver {
	return &KeybindingResolver{
		registry:       reg,
		contextManager: cm,
		userConfig:     cfg,
	}
}

// Resolve finds the action for a key event
func (kr *KeybindingResolver) Resolve(msg tea.KeyMsg, actx *ActionContext) (*Action, bool) {
	// Convert tea.KeyMsg to key pattern string
	keyStr := kr.keyMsgToString(msg)

	// Get current context
	ctx := kr.contextManager.Current()

	// Find actions for this key in this context
	actions := kr.registry.ActionsForKey(keyStr, ctx, actx)

	if len(actions) == 0 {
		return nil, false
	}

	// If multiple actions match, prioritize by:
	// 1. Context-specific over global
	// 2. User-configured over default
	// 3. First registered (stable)

	// For now, return first match
	return actions[0], true
}

// GetKeybinding returns the primary keybinding for an action
func (kr *KeybindingResolver) GetKeybinding(actionID ActionID) (KeyPattern, bool) {
	// Check user overrides first
	if overrides, ok := kr.userConfig.Keybindings[actionID]; ok && len(overrides) > 0 {
		return overrides[0], true
	}

	// Get default from action
	action, exists := kr.registry.Get(actionID)
	if !exists || len(action.Keybindings) == 0 {
		return KeyPattern{}, false
	}

	return action.Keybindings[0], true
}

// Helper: Convert tea.KeyMsg to key pattern string
func (kr *KeybindingResolver) keyMsgToString(msg tea.KeyMsg) string {
	var parts []string

	// Handle modifiers
	if msg.Alt {
		parts = append(parts, "alt")
	}

	// Handle key type
	switch msg.Type {
	case tea.KeyCtrlC:
		parts = append(parts, "ctrl", "c")
	case tea.KeyCtrlS:
		parts = append(parts, "ctrl", "s")
	case tea.KeyEnter:
		parts = append(parts, "enter")
	case tea.KeyEsc:
		parts = append(parts, "esc")
	case tea.KeyTab:
		parts = append(parts, "tab")
	case tea.KeySpace:
		parts = append(parts, "space")
	case tea.KeyRunes:
		// Regular character
		if len(msg.Runes) > 0 {
			parts = append(parts, string(msg.Runes[0]))
		}
	default:
		parts = append(parts, msg.String())
	}

	return strings.Join(parts, "+")
}
```

---

## Action Execution Pipeline

### Execution Coordinator

```go
// Package: internal/tui/actions

// Executor handles action execution with middleware support
type Executor struct {
	registry   *Registry
	middleware []ActionModifier
}

// NewExecutor creates a new action executor
func NewExecutor(reg *Registry) *Executor {
	return &Executor{
		registry:   reg,
		middleware: []ActionModifier{},
	}
}

// Use adds middleware to the execution pipeline
func (e *Executor) Use(modifier ActionModifier) {
	e.middleware = append(e.middleware, modifier)
}

// Execute runs an action with middleware
func (e *Executor) Execute(action *Action, ctx *ActionContext) tea.Cmd {
	// Check predicate one more time
	if action.Predicate != nil && !action.Predicate(ctx) {
		return func() tea.Msg {
			return actionErrorMsg{
				actionID: action.ID,
				err:      fmt.Errorf("action not available"),
			}
		}
	}

	// Run before middleware
	for _, mw := range e.middleware {
		if err := mw.Before(ctx); err != nil {
			return func() tea.Msg {
				return actionErrorMsg{actionID: action.ID, err: err}
			}
		}
	}

	// Run action-specific modifiers
	for _, mw := range action.Modifiers {
		if err := mw.Before(ctx); err != nil {
			return func() tea.Msg {
				return actionErrorMsg{actionID: action.ID, err: err}
			}
		}
	}

	// Execute action
	cmd := action.Execute(ctx)

	// Wrap command with after middleware
	return func() tea.Msg {
		// Execute original command
		var msg tea.Msg
		if cmd != nil {
			msg = cmd()
		}

		// Run after middleware (in reverse order)
		var err error
		for i := len(action.Modifiers) - 1; i >= 0; i-- {
			cmd, err = action.Modifiers[i].After(ctx, cmd, err)
			if err != nil {
				return actionErrorMsg{actionID: action.ID, err: err}
			}
		}

		for i := len(e.middleware) - 1; i >= 0; i-- {
			cmd, err = e.middleware[i].After(ctx, cmd, err)
			if err != nil {
				return actionErrorMsg{actionID: action.ID, err: err}
			}
		}

		return msg
	}
}

// actionErrorMsg indicates action execution failed
type actionErrorMsg struct {
	actionID ActionID
	err      error
}
```

### Common Modifiers

```go
// Package: internal/tui/actions/modifiers

// LoggingModifier logs action execution
type LoggingModifier struct {
	logger Logger
}

func (m *LoggingModifier) Before(ctx *ActionContext) error {
	m.logger.Info("executing action", "context", ctx.Context)
	return nil
}

func (m *LoggingModifier) After(ctx *ActionContext, result tea.Cmd, err error) (tea.Cmd, error) {
	if err != nil {
		m.logger.Error("action failed", "context", ctx.Context, "error", err)
	} else {
		m.logger.Info("action completed", "context", ctx.Context)
	}
	return result, err
}

// ConfirmationModifier prompts for confirmation before destructive actions
type ConfirmationModifier struct {
	message string
}

func (m *ConfirmationModifier) Before(ctx *ActionContext) error {
	// In real implementation, would show modal and wait for user input
	// For now, simplified
	return nil
}

func (m *ConfirmationModifier) After(ctx *ActionContext, result tea.Cmd, err error) (tea.Cmd, error) {
	return result, err
}

// MetricsModifier tracks action usage metrics
type MetricsModifier struct {
	metrics MetricsService
}

func (m *MetricsModifier) Before(ctx *ActionContext) error {
	m.metrics.RecordActionStart(ctx.Context)
	return nil
}

func (m *MetricsModifier) After(ctx *ActionContext, result tea.Cmd, err error) (tea.Cmd, error) {
	m.metrics.RecordActionComplete(ctx.Context, err == nil)
	return result, err
}
```

---

## Bubbletea Integration

### Model Integration

```go
// Package: internal/tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourorg/ticketr/internal/tui/actions"
)

// Model is the root Bubbletea model
type Model struct {
	// Action system
	actionRegistry  *actions.Registry
	contextManager  *actions.ContextManager
	keybindingResolver *actions.KeybindingResolver
	executor        *actions.Executor

	// UI components
	workspaceList WorkspaceListModel
	ticketTree    TicketTreeModel
	ticketDetail  TicketDetailModel

	// State
	focus            Focus
	selectedTickets  []string
	currentWorkspace *actions.WorkspaceState

	// Services
	services *actions.ServiceContainer
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Initialize action system
	m.actionRegistry = actions.NewRegistry()
	m.contextManager = actions.NewContextManager(actions.ContextWorkspaceList)
	m.executor = actions.NewExecutor(m.actionRegistry)

	// Register built-in actions
	m.registerBuiltinActions()

	// Create keybinding resolver
	m.keybindingResolver = actions.NewKeybindingResolver(
		m.actionRegistry,
		m.contextManager,
		&actions.UserConfig{}, // Load from config file
	)

	return nil
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Build action context
		actx := m.buildActionContext()

		// Resolve keybinding to action
		action, found := m.keybindingResolver.Resolve(msg, actx)
		if !found {
			// No action for this key - pass to focused component
			return m.updateFocusedComponent(msg)
		}

		// Execute action
		cmd := m.executor.Execute(action, actx)
		return m, cmd

	case ticketOpenedMsg:
		// Handle action result
		m.ticketDetail.SetTicket(msg.ticket)
		m.contextManager.Switch(actions.ContextTicketDetail)
		m.focus = FocusDetail
		return m, nil

	case actionErrorMsg:
		// Handle action error
		m.showError(msg.err)
		return m, nil

	// ... other message handlers
	}

	return m, nil
}

// buildActionContext creates context for action execution
func (m Model) buildActionContext() *actions.ActionContext {
	return &actions.ActionContext{
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

// registerBuiltinActions registers all built-in actions
func (m Model) registerBuiltinActions() {
	// Register actions from different modules
	registerWorkspaceActions(m.actionRegistry)
	registerTicketActions(m.actionRegistry)
	registerNavigationActions(m.actionRegistry)
	registerSyncActions(m.actionRegistry)
	registerSystemActions(m.actionRegistry)
}
```

### Message Types

```go
// Package: internal/tui

// Action result messages
type ticketOpenedMsg struct {
	ticket *domain.Ticket
}

type ticketSavedMsg struct {
	ticket *domain.Ticket
}

type ticketDeletedMsg struct {
	id string
}

type workspaceChangedMsg struct {
	workspace *domain.Workspace
}

type syncStartedMsg struct {
	jobID string
}

type syncCompletedMsg struct {
	jobID string
	err   error
}

type actionErrorMsg struct {
	actionID actions.ActionID
	err      error
}
```

---

## Extension Points

### 1. Custom Actions (User-Defined)

```go
// Users can define custom actions in config file or via API

// Example: Custom action to mark ticket as "In Review"
func registerCustomActions(registry *actions.Registry) {
	registry.Register(&actions.Action{
		ID:          "custom.mark_in_review",
		Name:        "Mark In Review",
		Description: "Mark selected ticket as In Review",
		Category:    actions.CategoryEdit,
		Contexts:    []actions.Context{actions.ContextTicketTree, actions.ContextTicketDetail},
		Keybindings: []actions.KeyPattern{{Key: "r"}},
		Predicate:   predicates.HasSingleSelection(),
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				ticket, err := ctx.Services.TicketService.Get(ctx.SelectedTickets[0])
				if err != nil {
					return actionErrorMsg{err: err}
				}

				ticket.Status = "In Review"
				err = ctx.Services.TicketService.Save(ticket)
				if err != nil {
					return actionErrorMsg{err: err}
				}

				return ticketSavedMsg{ticket: ticket}
			}
		},
	})
}
```

### 2. Plugin Actions (Lua Integration)

```go
// Package: internal/tui/actions/plugins

// PluginAction wraps a Lua-defined action
type PluginAction struct {
	*Action
	luaFunction string // Name of Lua function to call
	plugin      *Plugin
}

// NewPluginAction creates an action from Lua definition
func NewPluginAction(plugin *Plugin, def LuaActionDefinition) (*PluginAction, error) {
	action := &Action{
		ID:          ActionID(def.ID),
		Name:        def.Name,
		Description: def.Description,
		Category:    ActionCategory(def.Category),
		Contexts:    parseContexts(def.Contexts),
		Keybindings: parseKeybindings(def.Keybindings),
		Tags:        def.Tags,
		Execute:     nil, // Set below
	}

	pa := &PluginAction{
		Action:      action,
		luaFunction: def.Function,
		plugin:      plugin,
	}

	// Wrap Lua function call
	action.Execute = pa.executeLua

	return pa, nil
}

// executeLua calls the Lua function
func (pa *PluginAction) executeLua(ctx *ActionContext) tea.Cmd {
	return func() tea.Msg {
		// Call Lua function with context
		result, err := pa.plugin.Call(pa.luaFunction, ctx)
		if err != nil {
			return actionErrorMsg{actionID: pa.ID, err: err}
		}

		// Convert Lua result to message
		return pa.luaResultToMsg(result)
	}
}

// LuaActionDefinition is the Lua table structure
type LuaActionDefinition struct {
	ID          string
	Name        string
	Description string
	Category    string
	Contexts    []string
	Keybindings []string
	Tags        []string
	Function    string // Name of Lua function to execute
}
```

### 3. Action Composition

```go
// Package: internal/tui/actions/compose

// SequenceAction executes actions in sequence
func SequenceAction(actions ...*Action) *Action {
	return &Action{
		ID:          "compose.sequence",
		Name:        "Sequence",
		Description: "Execute actions in sequence",
		Execute: func(ctx *ActionContext) tea.Cmd {
			cmds := make([]tea.Cmd, len(actions))
			for i, action := range actions {
				cmds[i] = action.Execute(ctx)
			}
			return tea.Sequence(cmds...)
		},
	}
}

// ParallelAction executes actions in parallel
func ParallelAction(actions ...*Action) *Action {
	return &Action{
		ID:          "compose.parallel",
		Name:        "Parallel",
		Description: "Execute actions in parallel",
		Execute: func(ctx *ActionContext) tea.Cmd {
			cmds := make([]tea.Cmd, len(actions))
			for i, action := range actions {
				cmds[i] = action.Execute(ctx)
			}
			return tea.Batch(cmds...)
		},
	}
}

// ConditionalAction executes action only if predicate is true
func ConditionalAction(pred PredicateFunc, action *Action) *Action {
	return &Action{
		ID:          action.ID + ".conditional",
		Name:        action.Name,
		Description: action.Description,
		Predicate:   pred,
		Execute:     action.Execute,
	}
}

// ChainAction chains actions where each gets result of previous
func ChainAction(actions ...*Action) *Action {
	return &Action{
		ID:          "compose.chain",
		Name:        "Chain",
		Description: "Chain action execution",
		Execute: func(ctx *ActionContext) tea.Cmd {
			// First action
			if len(actions) == 0 {
				return nil
			}

			return func() tea.Msg {
				var result tea.Msg
				for _, action := range actions {
					cmd := action.Execute(ctx)
					if cmd != nil {
						result = cmd()
					}
					// TODO: Pass result to next action via ctx
				}
				return result
			}
		},
	}
}
```

---

## Code Examples

### Example 1: Registering Built-in Actions

```go
// Package: internal/tui/actions/builtin

package builtin

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourorg/ticketr/internal/tui/actions"
	"github.com/yourorg/ticketr/internal/tui/actions/predicates"
)

// RegisterTicketActions registers all ticket-related actions
func RegisterTicketActions(reg *actions.Registry) {
	// Open ticket
	reg.Register(&actions.Action{
		ID:          "ticket.open",
		Name:        "Open Ticket",
		Description: "Open selected ticket in detail view",
		Category:    actions.CategoryView,
		Contexts: []actions.Context{
			actions.ContextTicketTree,
		},
		Keybindings: []actions.KeyPattern{
			{Key: "enter"},
			{Key: "o"},
		},
		Predicate: predicates.HasSingleSelection(),
		Execute:   executeOpenTicket,
		Tags:      []string{"ticket", "open", "view", "show"},
		Icon:      "📄",
	})

	// Edit ticket
	reg.Register(&actions.Action{
		ID:          "ticket.edit",
		Name:        "Edit Ticket",
		Description: "Edit selected ticket",
		Category:    actions.CategoryEdit,
		Contexts: []actions.Context{
			actions.ContextTicketDetail,
		},
		Keybindings: []actions.KeyPattern{
			{Key: "e"},
		},
		Predicate: predicates.And(
			predicates.HasSingleSelection(),
			predicates.IsOnline(),
		),
		Execute: executeEditTicket,
		Tags:    []string{"ticket", "edit", "modify"},
		Icon:    "✏️",
	})

	// Delete ticket
	reg.Register(&actions.Action{
		ID:          "ticket.delete",
		Name:        "Delete Ticket",
		Description: "Delete selected ticket(s)",
		Category:    actions.CategoryEdit,
		Contexts: []actions.Context{
			actions.ContextTicketTree,
			actions.ContextTicketDetail,
		},
		Keybindings: []actions.KeyPattern{
			{Key: "d"},
			{Key: "delete"},
		},
		Predicate: predicates.And(
			predicates.HasSelection(),
			predicates.IsOnline(),
		),
		Execute: executeDeleteTicket,
		Modifiers: []actions.ActionModifier{
			&ConfirmationModifier{
				message: "Are you sure you want to delete the selected ticket(s)?",
			},
		},
		Tags: []string{"ticket", "delete", "remove"},
		Icon: "🗑️",
	})

	// Toggle selection
	reg.Register(&actions.Action{
		ID:          "ticket.toggle_selection",
		Name:        "Toggle Selection",
		Description: "Select or deselect ticket",
		Category:    actions.CategoryNavigation,
		Contexts: []actions.Context{
			actions.ContextTicketTree,
		},
		Keybindings: []actions.KeyPattern{
			{Key: "space"},
		},
		Predicate: predicates.Always(),
		Execute:   executeToggleSelection,
		Tags:      []string{"ticket", "select", "multi"},
		Icon:      "☑️",
	})

	// Bulk operations
	reg.Register(&actions.Action{
		ID:          "ticket.bulk_operations",
		Name:        "Bulk Operations",
		Description: "Perform bulk operations on selected tickets",
		Category:    actions.CategoryBulk,
		Contexts: []actions.Context{
			actions.ContextTicketTree,
		},
		Keybindings: []actions.KeyPattern{
			{Key: "b"},
		},
		Predicate: predicates.HasMultipleSelection(),
		Execute:   executeBulkOperations,
		Tags:      []string{"ticket", "bulk", "multiple"},
		Icon:      "📦",
	})
}

// Action handlers
func executeOpenTicket(ctx *actions.ActionContext) tea.Cmd {
	return func() tea.Msg {
		if len(ctx.SelectedTickets) == 0 {
			return actionErrorMsg{err: fmt.Errorf("no ticket selected")}
		}

		ticket, err := ctx.Services.TicketService.Get(ctx.SelectedTickets[0])
		if err != nil {
			return actionErrorMsg{err: err}
		}

		return ticketOpenedMsg{ticket: ticket}
	}
}

func executeEditTicket(ctx *actions.ActionContext) tea.Cmd {
	return func() tea.Msg {
		// Switch to edit mode
		return ticketEditModeMsg{}
	}
}

func executeDeleteTicket(ctx *actions.ActionContext) tea.Cmd {
	return func() tea.Msg {
		for _, id := range ctx.SelectedTickets {
			err := ctx.Services.TicketService.Delete(id)
			if err != nil {
				return actionErrorMsg{err: err}
			}
		}

		return ticketsDeletedMsg{ids: ctx.SelectedTickets}
	}
}

func executeToggleSelection(ctx *actions.ActionContext) tea.Cmd {
	return func() tea.Msg {
		// Toggle selection - handled by tree component
		return ticketSelectionToggledMsg{}
	}
}

func executeBulkOperations(ctx *actions.ActionContext) tea.Cmd {
	return func() tea.Msg {
		// Show bulk operations modal
		return showBulkOperationsModalMsg{}
	}
}
```

### Example 2: Context-Dependent Action Filtering

```go
// In the ActionBar component

func (ab ActionBarModel) View() string {
	// Get current context
	ctx := ab.contextManager.Current()

	// Build action context
	actx := ab.buildActionContext()

	// Get available actions for this context
	availableActions := ab.registry.ActionsForContext(ctx, actx)

	// Build keybinding display
	var bindings []string
	for _, action := range availableActions {
		if len(action.Keybindings) > 0 {
			key := action.Keybindings[0].Key
			bindings = append(bindings, fmt.Sprintf("[%s %s]", key, action.Name))
		}
	}

	// Render
	content := strings.Join(bindings, " ")
	return actionBarStyle.Render(content)
}
```

### Example 3: Command Palette Integration

```go
// Package: internal/tui/components

type CommandPaletteModel struct {
	registry    *actions.Registry
	input       textinput.Model
	results     []*actions.Action
	selectedIdx int
}

func (cp CommandPaletteModel) Update(msg tea.Msg) (CommandPaletteModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Execute selected action
			if cp.selectedIdx < len(cp.results) {
				action := cp.results[cp.selectedIdx]
				actx := cp.buildActionContext()
				return cp, action.Execute(actx)
			}

		case "up", "ctrl+k":
			if cp.selectedIdx > 0 {
				cp.selectedIdx--
			}

		case "down", "ctrl+j":
			if cp.selectedIdx < len(cp.results)-1 {
				cp.selectedIdx++
			}
		}
	}

	// Update search input
	var cmd tea.Cmd
	cp.input, cmd = cp.input.Update(msg)

	// Search actions
	query := cp.input.Value()
	actx := cp.buildActionContext()
	cp.results = cp.registry.Search(query, actx)
	cp.selectedIdx = 0

	return cp, cmd
}

func (cp CommandPaletteModel) View() string {
	// Render input
	s := "Search commands:\n"
	s += cp.input.View() + "\n\n"

	// Render results
	if len(cp.results) == 0 {
		s += "No results"
		return s
	}

	for i, action := range cp.results {
		prefix := "  "
		if i == cp.selectedIdx {
			prefix = "> "
		}

		// Show icon, name, keybinding
		keyStr := ""
		if len(action.Keybindings) > 0 {
			keyStr = action.Keybindings[0].Key
		}

		s += fmt.Sprintf("%s%s %s [%s] - %s\n",
			prefix,
			action.Icon,
			action.Name,
			keyStr,
			action.Description,
		)
	}

	return s
}
```

### Example 4: Lua Plugin Action

```lua
-- File: ~/.config/ticketr/plugins/github_link.lua

-- Define action
action = {
    id = "plugin.github_link",
    name = "Open in GitHub",
    description = "Open ticket in GitHub web UI",
    category = "View",
    contexts = {"ticket_tree", "ticket_detail"},
    keybindings = {"g"},
    tags = {"github", "web", "browser"},

    -- Predicate: only show if ticket has GitHub URL
    predicate = function(ctx)
        if #ctx.selected_tickets ~= 1 then
            return false
        end

        local ticket = ctx.services.ticket:get(ctx.selected_tickets[1])
        return ticket.remote_url ~= nil
    end,

    -- Execute: open browser
    execute = function(ctx)
        local ticket = ctx.services.ticket:get(ctx.selected_tickets[1])
        os.execute("xdg-open " .. ticket.remote_url)
        return {
            type = "notification",
            message = "Opened in browser"
        }
    end
}

-- Register action
ticketr.register_action(action)
```

---

## Migration Path

### Phase 1: Foundation (Week 1)

**Goal:** Set up action system infrastructure

1. **Create package structure**
   ```
   internal/tui/actions/
   ├── action.go           # Core types
   ├── context.go          # Context manager
   ├── registry.go         # Action registry
   ├── resolver.go         # Keybinding resolver
   ├── executor.go         # Execution pipeline
   ├── predicates/         # Common predicates
   │   └── predicates.go
   ├── modifiers/          # Common modifiers
   │   └── modifiers.go
   └── builtin/            # Built-in actions
       ├── tickets.go
       ├── workspaces.go
       ├── navigation.go
       └── sync.go
   ```

2. **Implement core types**
   - Action struct
   - ActionContext struct
   - Context types
   - Predicate functions

3. **Implement Registry**
   - Register/Unregister
   - Get/All/Search
   - ActionsForContext
   - ActionsForKey

4. **Write tests**
   - Registry tests
   - Predicate composition tests
   - Context manager tests

### Phase 2: Built-in Actions (Week 2)

**Goal:** Port existing actions from current TUI

1. **Inventory current actions**
   - Map all keybindings from ActionBar
   - Extract handlers from view files
   - Identify contexts

2. **Implement action modules**
   - Ticket actions (open, edit, delete, etc.)
   - Workspace actions (create, select, delete)
   - Navigation actions (tab, back, help)
   - Sync actions (pull, push, status)
   - System actions (quit, help, command palette)

3. **Register in model**
   ```go
   func (m *Model) registerBuiltinActions() {
       builtin.RegisterTicketActions(m.actionRegistry)
       builtin.RegisterWorkspaceActions(m.actionRegistry)
       builtin.RegisterNavigationActions(m.actionRegistry)
       builtin.RegisterSyncActions(m.actionRegistry)
       builtin.RegisterSystemActions(m.actionRegistry)
   }
   ```

### Phase 3: Keybinding Integration (Week 2-3)

**Goal:** Replace manual keybinding handling

1. **Implement KeybindingResolver**
   - Convert tea.KeyMsg to key pattern
   - Resolve to action
   - Priority handling (context-specific > global)

2. **Update Model.Update()**
   ```go
   case tea.KeyMsg:
       actx := m.buildActionContext()
       action, found := m.keybindingResolver.Resolve(msg, actx)
       if found {
           cmd := m.executor.Execute(action, actx)
           return m, cmd
       }
       // Fallback to component handling
   ```

3. **Update ActionBar component**
   - Show actions from registry (not hardcoded)
   - Filter by context and predicates
   - Dynamic updates when context changes

### Phase 4: Command Palette (Week 3)

**Goal:** Implement fuzzy search for all actions

1. **Build CommandPaletteModel**
   - Search input (bubbles/textinput)
   - Results list (bubbles/list)
   - Fuzzy matching on name/description/tags

2. **Integrate with model**
   - Show on Ctrl+P or :
   - Execute selected action
   - Close on Esc

3. **Add search enhancements**
   - Highlight matching text
   - Show category headers
   - Show recent actions

### Phase 5: Testing & Polish (Week 4)

**Goal:** Ensure reliability and UX

1. **Write comprehensive tests**
   - Action execution tests
   - Predicate logic tests
   - Keybinding resolution tests
   - Integration tests

2. **Add middleware**
   - Logging modifier
   - Metrics modifier
   - Confirmation modifier (for destructive actions)

3. **Documentation**
   - API docs for creating actions
   - User guide for keybindings
   - Plugin developer guide

### Phase 6: Extension Points (Week 5-6)

**Goal:** Enable custom actions and plugins

1. **User-defined actions**
   - Load from YAML/JSON config
   - Hot reload on config change

2. **Lua plugin architecture** (future)
   - Lua VM integration
   - Action registration API
   - Context/service bindings

3. **Action marketplace** (future vision)
   - Share actions with community
   - Install from GitHub gists
   - Version management

---

## Future: Lua Plugin Architecture

### Plugin Structure

```
~/.config/ticketr/plugins/
├── github_integration/
│   ├── init.lua          # Entry point
│   ├── actions.lua       # Action definitions
│   └── manifest.json     # Metadata
└── slack_notifications/
    ├── init.lua
    ├── actions.lua
    └── manifest.json
```

### Plugin Manifest

```json
{
  "name": "github_integration",
  "version": "1.0.0",
  "description": "GitHub integration for ticketr",
  "author": "username",
  "license": "MIT",
  "main": "init.lua",
  "dependencies": {
    "ticketr": ">=0.5.0"
  },
  "permissions": [
    "network",
    "exec"
  ]
}
```

### Lua API (Future Design)

```lua
-- Plugin: github_integration/init.lua

-- Register action
ticketr.register_action({
    id = "github.open_pr",
    name = "Open Pull Request",
    description = "Open GitHub PR in browser",
    category = "View",
    contexts = {"ticket_detail"},
    keybindings = {"shift+g"},

    predicate = function(ctx)
        local ticket = ctx.services.ticket:get(ctx.selected_tickets[1])
        return ticket.metadata.github_pr ~= nil
    end,

    execute = function(ctx)
        local ticket = ctx.services.ticket:get(ctx.selected_tickets[1])
        local url = ticket.metadata.github_pr

        -- Call external command
        ticketr.exec("xdg-open " .. url)

        -- Return result
        return {
            type = "notification",
            message = "Opened PR in browser"
        }
    end
})

-- Register custom predicate
ticketr.register_predicate("has_github_pr", function(ctx)
    if #ctx.selected_tickets ~= 1 then
        return false
    end
    local ticket = ctx.services.ticket:get(ctx.selected_tickets[1])
    return ticket.metadata.github_pr ~= nil
end)

-- Register custom keybinding
ticketr.register_keybinding("github.open_pr", "shift+g")

-- Hook into events
ticketr.on_ticket_saved(function(ticket)
    -- Sync to GitHub when ticket is saved
    if ticket.metadata.github_issue then
        github.update_issue(ticket)
    end
end)
```

### Plugin Security

1. **Sandboxing**
   - Restricted Lua environment
   - No direct file access
   - Network access requires permission

2. **Permissions**
   - Declared in manifest
   - User approval required
   - Can be revoked

3. **API surface**
   - Only exposed services/context
   - No access to internal state
   - Versioned API contracts

---

## Appendix: Full Example Action

Here's a complete, production-ready example of the "Open Ticket" action:

```go
// Package: internal/tui/actions/builtin

var OpenTicketAction = &actions.Action{
	ID:          "ticket.open",
	Name:        "Open Ticket",
	Description: "Open the selected ticket in the detail view",
	Category:    actions.CategoryView,

	// Available in ticket tree view only
	Contexts: []actions.Context{
		actions.ContextTicketTree,
	},

	// Multiple keybindings for same action
	Keybindings: []actions.KeyPattern{
		{Key: "enter"},
		{Key: "o"},
	},

	// Only available when exactly one ticket is selected
	Predicate: predicates.HasSingleSelection(),

	// Show in command palette and help
	ShowInUI: func(ctx *actions.ActionContext) bool {
		return true
	},

	// Execute the action
	Execute: func(ctx *actions.ActionContext) tea.Cmd {
		return func() tea.Msg {
			// Validate selection
			if len(ctx.SelectedTickets) != 1 {
				return actionErrorMsg{
					actionID: "ticket.open",
					err:      fmt.Errorf("exactly one ticket must be selected"),
				}
			}

			// Load ticket from service
			ticket, err := ctx.Services.TicketService.Get(ctx.SelectedTickets[0])
			if err != nil {
				return actionErrorMsg{
					actionID: "ticket.open",
					err:      fmt.Errorf("failed to load ticket: %w", err),
				}
			}

			// Return success message
			return ticketOpenedMsg{
				ticket: ticket,
			}
		}
	},

	// Metadata
	Icon: "📄",
	Tags: []string{"ticket", "open", "view", "show", "detail"},
	Metadata: map[string]string{
		"help_url": "https://docs.ticketr.dev/actions/ticket-open",
	},

	// Modifiers
	Modifiers: []actions.ActionModifier{
		&modifiers.LoggingModifier{},
		&modifiers.MetricsModifier{},
	},
}
```

---

## Conclusion

This extensible action system design provides:

1. ✅ **Midnight Commander-inspired** - Context menus and F-keys
2. ✅ **Declarative** - Actions are data, not scattered callbacks
3. ✅ **Testable** - Pure functions, predictable behavior
4. ✅ **Extensible** - Plugins can register custom actions
5. ✅ **Bubbletea-native** - Uses messages, not callbacks
6. ✅ **Discoverable** - Command palette shows all actions
7. ✅ **Customizable** - User can override keybindings

The system is designed to grow with the application, supporting everything from simple keybindings to complex Lua plugins, while maintaining a clean separation between action definitions and their execution.

The migration path is incremental and low-risk, allowing the team to adopt the system piece by piece while the current Tview implementation continues to work.
