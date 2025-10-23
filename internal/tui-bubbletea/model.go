package tuibubbletea

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/actions"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/components"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/components/help"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/components/tree"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/layout"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/views/cmdpalette"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/views/detail"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/views/search"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/views/workspace"
)

// Focus indicates which panel has focus
type Focus int

const (
	FocusLeft Focus = iota
	FocusRight
	FocusWorkspace
)

// Model represents the root state of the TUI application.
// It contains all child component models and manages the overall application state.
type Model struct {
	// Dimensions
	width  int
	height int

	// Layout manager (Day 2-3: FlexBox integration)
	layout *layout.CompleteLayout

	// Theme (Week 3 Day 1: Proper Elm Architecture - no global state)
	theme *theme.Theme

	// State
	ready   bool  // Whether the initial WindowSizeMsg has been received
	focused Focus // Which panel has focus (Left or Right)

	// Services (Week 2 Day 1: Data integration)
	workspaceService *services.WorkspaceService
	ticketQuery      *services.TicketQueryService

	// Data (Week 2 Day 1: Real data from services)
	currentWorkspace *domain.Workspace
	tickets          []domain.Ticket
	workspaces       []domain.Workspace

	// Loading states (Week 2 Day 1: Async data loading)
	loadingWorkspaces bool
	loadingTickets    bool
	dataLoaded        bool
	loadError         error

	// Day 2-3: Tree component (Week 2 Days 2-3)
	ticketTree tree.TreeModel

	// Day 4-5: View components
	detailView        detail.Model
	workspaceSelector workspace.Model
	selectedTicket    *domain.Ticket

	// Modal state
	showWorkspaceModal bool

	// Terminal size validation (Week 2 Day 2)
	terminalTooSmall bool

	// Loading spinner (Week 2 Day 2)
	loadingSpinner components.LoadingModel

	// Help screen (Week 2 Day 2, Week 4 Day 1: Upgraded to use action registry)
	helpScreen help.HelpModel

	// Week 3 components (Week 4 Day 1: Integration)
	searchModal search.Model
	cmdPalette  cmdpalette.Model

	// Week 3 Day 1: Action system (Week 4 Day 1: Integration)
	actionRegistry *actions.Registry
	contextManager *actions.ContextManager
}

// initialModel creates the initial state of the application.
// This is called once when the program starts, before any messages are processed.
// Week 2 Day 1: Now accepts services for data integration.
// Week 4 Day 1: Now initializes action registry and Week 3 components.
func initialModel(workspaceService *services.WorkspaceService, ticketQuery *services.TicketQueryService) Model {
	// Start with reasonable default dimensions (will be updated on first WindowSizeMsg)
	width, height := 120, 30

	// Calculate panel dimensions for initial component sizing
	leftWidth, rightWidth, contentHeight := layout.NewCompleteLayout(width, height).GetPanelDimensions()

	// Week 4 Day 1: Initialize action system
	actionRegistry := actions.NewRegistry()
	contextManager := actions.NewContextManager(actions.ContextTicketTree)

	// Register all built-in actions
	if err := RegisterBuiltinActions(actionRegistry); err != nil {
		// Log error but continue - we'll handle gracefully
		// In production, this should never fail as it's our built-in actions
		panic("failed to register built-in actions: " + err.Error())
	}

	// Create tree component (Week 2 Days 2-3, Week 3 Day 3: Now theme-aware)
	ticketTree := tree.New(leftWidth, contentHeight, &theme.DefaultTheme)

	// Create loading spinner (Week 2 Day 2)
	loadingSpinner := components.NewLoading("Loading...", &theme.DefaultTheme)

	// Week 4 Day 1: Create help screen with action registry (upgraded from Legacy)
	helpScreen := help.New(width, height, &theme.DefaultTheme, actionRegistry, contextManager)

	// Week 4 Day 1: Create Week 3 components
	searchModal := search.New(actionRegistry, &theme.DefaultTheme)
	cmdPalette := cmdpalette.New(actionRegistry, contextManager, &theme.DefaultTheme)

	return Model{
		width:              width,
		height:             height,
		layout:             layout.NewCompleteLayout(width, height),
		theme:              &theme.DefaultTheme,
		ready:              false,
		focused:            FocusLeft,
		workspaceService:   workspaceService,
		ticketQuery:        ticketQuery,
		loadingWorkspaces:  true,
		loadingTickets:     false,
		dataLoaded:         false,
		tickets:            []domain.Ticket{},
		workspaces:         []domain.Workspace{},
		ticketTree:         ticketTree,
		detailView:         detail.New(rightWidth, contentHeight),
		workspaceSelector:  workspace.New([]domain.Workspace{}, width/2, height/2),
		showWorkspaceModal: false,
		terminalTooSmall:   false,
		loadingSpinner:     loadingSpinner,
		helpScreen:         helpScreen,
		searchModal:        searchModal,
		cmdPalette:         cmdPalette,
		actionRegistry:     actionRegistry,
		contextManager:     contextManager,
	}
}

// Init is called once when the program starts.
// It returns an optional initial command to run.
// Week 2 Day 1: Now loads workspace and ticket data asynchronously.
// Week 2 Days 2-3: Now initializes tree component.
// Week 2 Day 2: Now initializes loading spinner.
// Week 4 Day 1: Now initializes Week 3 components.
func (m Model) Init() tea.Cmd {
	// Initialize all components
	return tea.Batch(
		m.ticketTree.Init(),
		m.detailView.Init(),
		m.workspaceSelector.Init(),
		m.loadingSpinner.Init(),
		m.helpScreen.Init(),
		m.searchModal.Init(),
		m.cmdPalette.Init(),
	)
}

// GetCurrentTheme returns the active theme name
func (m Model) GetCurrentTheme() string {
	return m.theme.Name
}

// SetFocus sets which panel has focus
func (m *Model) SetFocus(f Focus) {
	m.focused = f
}

// ToggleFocus switches focus between left and right panels
func (m *Model) ToggleFocus() {
	if m.focused == FocusLeft {
		m.focused = FocusRight
	} else {
		m.focused = FocusLeft
	}
}

// SetSelectedTicket sets the currently selected ticket and updates the detail view
func (m *Model) SetSelectedTicket(ticket *domain.Ticket) {
	m.selectedTicket = ticket
	m.detailView.SetTicket(ticket)
	m.SetFocus(FocusRight)
}

// buildActionContext creates an ActionContext for executing actions.
// Week 4 Day 1: Provides context to actions for execution and predicate evaluation.
func (m *Model) buildActionContext() *actions.ActionContext {
	// Build list of selected tickets
	var selectedTickets []string
	if m.selectedTicket != nil {
		selectedTickets = []string{m.selectedTicket.JiraID}
	}

	// Build workspace state
	var workspaceState *actions.WorkspaceState
	if m.currentWorkspace != nil {
		workspaceState = &actions.WorkspaceState{
			ID:          m.currentWorkspace.ID,
			Name:        m.currentWorkspace.Name,
			ProfileID:   "", // Domain Workspace doesn't have ProfileID, use empty string
			IsDirty:     false, // TODO: Track dirty state
			TicketCount: len(m.tickets),
		}
	}

	return &actions.ActionContext{
		Context:           m.contextManager.Current(),
		SelectedTickets:   selectedTickets,
		SelectedWorkspace: workspaceState,
		HasUnsavedChanges: false, // TODO: Track unsaved changes
		IsSyncing:         false,  // TODO: Track sync state
		IsOffline:         false,  // TODO: Track offline state
		Width:             m.width,
		Height:            m.height,
		Services:          &actions.ServiceContainer{},
	}
}

// Note: Workspace switching is now handled directly in update.go
// via the workspace.WorkspaceSelectedMsg message (Week 3 Day 3)
// Note: propagateTheme() is defined in update.go (Week 4 Day 1)
