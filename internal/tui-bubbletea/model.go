package tuibubbletea

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/components"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/components/help"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/components/tree"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/layout"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/views/detail"
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

	// Help screen (Week 2 Day 2)
	helpScreen help.HelpModel

	// TODO(future): Add more child component models:
	// - header          header.Model
	// - actionBar       actionbar.Model
	// - commandPalette  commandpalette.Model
	//
	// TODO(day5): Add action system:
	// - actionRegistry     *actions.Registry
	// - contextManager     *actions.ContextManager
	// - keybindingResolver *actions.KeybindingResolver
}

// initialModel creates the initial state of the application.
// This is called once when the program starts, before any messages are processed.
// Week 2 Day 1: Now accepts services for data integration.
func initialModel(workspaceService *services.WorkspaceService, ticketQuery *services.TicketQueryService) Model {
	// Start with reasonable default dimensions (will be updated on first WindowSizeMsg)
	width, height := 120, 30

	// Calculate panel dimensions for initial component sizing
	leftWidth, rightWidth, contentHeight := layout.NewCompleteLayout(width, height).GetPanelDimensions()

	// Create tree component (Week 2 Days 2-3, Week 3 Day 3: Now theme-aware)
	ticketTree := tree.New(leftWidth, contentHeight, &theme.DefaultTheme)

	// Create loading spinner (Week 2 Day 2)
	loadingSpinner := components.NewLoading("Loading...", &theme.DefaultTheme)

	// Create help screen (Week 2 Day 2, Week 3 Day 4: Legacy mode until action system integrated)
	helpScreen := help.NewLegacy(width, height, &theme.DefaultTheme)

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
	}
}

// Init is called once when the program starts.
// It returns an optional initial command to run.
// Week 2 Day 1: Now loads workspace and ticket data asynchronously.
// Week 2 Days 2-3: Now initializes tree component.
// Week 2 Day 2: Now initializes loading spinner.
func (m Model) Init() tea.Cmd {
	// Initialize tree component and other components
	return tea.Batch(
		m.ticketTree.Init(),
		m.detailView.Init(),
		m.workspaceSelector.Init(),
		m.loadingSpinner.Init(),
		m.helpScreen.Init(),
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

// Note: Workspace switching is now handled directly in update.go
// via the workspace.WorkspaceSelectedMsg message (Week 3 Day 3)
