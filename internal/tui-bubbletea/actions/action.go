package actions

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ActionID uniquely identifies an action
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
	ID          ActionID
	Name        string
	Description string
	Category    ActionCategory

	// Keybindings (can have multiple for same action)
	Keybindings []KeyPattern

	// Visibility & Availability
	Contexts  []Context
	Predicate PredicateFunc
	ShowInUI  ShowInUIFunc

	// Execution
	Execute ExecuteFunc

	// Metadata (for UI, help, command palette)
	Icon     string
	Tags     []string
	Metadata map[string]string

	// Composition (for complex actions)
	Modifiers []ActionModifier
}

// KeyPattern defines a key combination
type KeyPattern struct {
	Key   string // e.g., "enter", "ctrl+s", "F1"
	Alt   bool   // Alt modifier
	Ctrl  bool   // Ctrl modifier
	Shift bool   // Shift modifier
}

// String converts key pattern to string for indexing
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
	SelectedTickets   []string
	SelectedWorkspace *WorkspaceState
	FocusedItem       interface{}

	// Global state (read-only for predicates)
	HasUnsavedChanges bool
	IsSyncing         bool
	IsOffline         bool

	// User preferences (for customization)
	Config *UserConfig

	// UI state (for rendering)
	Width  int
	Height int

	// Service layer (for execution)
	Services *ServiceContainer
}

// ActionModifier wraps action execution with pre/post hooks
type ActionModifier interface {
	Before(ctx *ActionContext) error
	After(ctx *ActionContext, result tea.Cmd, err error) (tea.Cmd, error)
}

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
	// Placeholder for now - will integrate with existing services
	// TicketService    TicketService
	// WorkspaceService WorkspaceService
	// SyncService      SyncService
}
