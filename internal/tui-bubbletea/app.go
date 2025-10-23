package tuibubbletea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/commands"
)

// App represents the root Bubbletea application.
// It encapsulates the tea.Program and provides a clean interface
// for starting the TUI.
type App struct {
	program *tea.Program
}

// NewApp creates a new Bubbletea TUI application.
// It initializes the root model and sets up the Bubbletea program
// with appropriate options for a full-screen terminal UI.
// Week 2 Day 1: Now accepts services for data integration.
func NewApp(workspaceService *services.WorkspaceService, ticketQuery *services.TicketQueryService) (*App, error) {
	if workspaceService == nil {
		return nil, fmt.Errorf("workspace service cannot be nil")
	}
	if ticketQuery == nil {
		return nil, fmt.Errorf("ticket query service cannot be nil")
	}

	// Create the initial model with services
	m := initialModel(workspaceService, ticketQuery)

	// Create the Bubbletea program with full-screen alt screen
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Send initial data loading commands
	p.Send(commands.LoadCurrentWorkspace(workspaceService))
	p.Send(commands.LoadWorkspaces(workspaceService))

	return &App{
		program: p,
	}, nil
}

// Run starts the TUI application and blocks until the user quits.
// Returns an error if the program fails to start or encounters a fatal error.
func (a *App) Run() error {
	if _, err := a.program.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}
	return nil
}
