package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/karolswdev/ticktr/internal/adapters/tui/views"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/rivo/tview"
)

// TUIApp represents the main TUI application.
type TUIApp struct {
	app              *tview.Application
	mainLayout       *tview.Flex
	router           *Router
	workspaceService *services.WorkspaceService
	ticketQuery      *services.TicketQueryService
	pathResolver     *services.PathResolver

	// View references for tri-panel layout
	workspaceListView *views.WorkspaceListView
	ticketTreeView    *views.TicketTreeView
	ticketDetailView  *views.TicketDetailView

	// Modal views (Week 14)
	searchView  *views.SearchView
	commandView *views.CommandPaletteView

	// Focus management
	currentFocus string   // "workspace_list", "ticket_tree", or "ticket_detail"
	focusOrder   []string // Order of focus cycling
}

// NewTUIApp creates a new TUI application instance.
func NewTUIApp(
	workspaceService *services.WorkspaceService,
	ticketQuery *services.TicketQueryService,
	pathResolver *services.PathResolver,
) (*TUIApp, error) {
	if workspaceService == nil {
		return nil, fmt.Errorf("workspace service is required")
	}
	if ticketQuery == nil {
		return nil, fmt.Errorf("ticket query service is required")
	}
	if pathResolver == nil {
		return nil, fmt.Errorf("path resolver is required")
	}

	app := tview.NewApplication()

	return &TUIApp{
		app:              app,
		router:           NewRouter(),
		workspaceService: workspaceService,
		ticketQuery:      ticketQuery,
		pathResolver:     pathResolver,
		currentFocus:     "workspace_list",
		focusOrder:       []string{"workspace_list", "ticket_tree", "ticket_detail"},
	}, nil
}

// Run starts the TUI application.
func (t *TUIApp) Run() error {
	if err := t.setupApp(); err != nil {
		return err
	}
	return t.app.Run()
}

// setupApp initializes all views and layouts.
func (t *TUIApp) setupApp() error {
	// Create workspace list view
	t.workspaceListView = views.NewWorkspaceListView(t.workspaceService)
	t.workspaceListView.SetSwitchHandler(func(name string) error {
		return t.workspaceService.Switch(name)
	})

	// Create ticket tree view
	t.ticketTreeView = views.NewTicketTreeView(t.workspaceService, t.ticketQuery)

	// Create ticket detail view
	t.ticketDetailView = views.NewTicketDetailView(t.app)

	// Set workspace change callback to reload tickets
	t.workspaceListView.SetWorkspaceChangeHandler(func(workspaceID string) {
		if t.ticketTreeView != nil {
			t.ticketTreeView.LoadTickets(workspaceID)
			t.setFocus("ticket_tree")
		}
	})

	// Set ticket selection callback to show detail view
	t.ticketTreeView.SetOnTicketSelected(func(ticket *domain.Ticket) {
		if t.ticketDetailView != nil {
			t.ticketDetailView.SetTicket(ticket)
			t.setFocus("ticket_detail")
		}
	})

	// Set ticket save callback (prepare for future sync)
	t.ticketDetailView.SetOnSave(func(ticket *domain.Ticket) error {
		// TODO: Implement actual sync in future phase
		// For now, just update in-memory state
		return nil
	})

	// Create help view and register with router
	helpView := views.NewHelpView()
	if err := t.router.Register(helpView); err != nil {
		return fmt.Errorf("failed to register help view: %w", err)
	}

	// Create search view (Week 14)
	t.searchView = views.NewSearchView(t.app)
	t.searchView.SetOnClose(func() {
		// Return to main layout
		t.app.SetRoot(t.mainLayout, true)
		t.updateFocus()
	})
	t.searchView.SetOnSelect(func(ticket *domain.Ticket) {
		// Return to main layout and show ticket detail
		t.app.SetRoot(t.mainLayout, true)
		if t.ticketDetailView != nil {
			t.ticketDetailView.SetTicket(ticket)
			t.setFocus("ticket_detail")
		}
	})

	// Create command palette view (Week 14)
	t.commandView = views.NewCommandPaletteView()
	t.commandView.SetOnClose(func() {
		// Return to main layout
		t.app.SetRoot(t.mainLayout, true)
		t.updateFocus()
	})
	// Set up available commands
	t.setupCommands()

	// Create tri-panel layout: workspace (30 fixed) | tree (40%) | detail (60%)
	t.mainLayout = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(t.workspaceListView.Primitive(), 30, 0, true). // Fixed 30 chars
		AddItem(t.ticketTreeView.Primitive(), 0, 2, false).    // 40% (2 of 5 parts)
		AddItem(t.ticketDetailView.Primitive(), 0, 3, false)   // 60% (3 of 5 parts)

	// Set up global key bindings
	t.app.SetInputCapture(t.globalKeyHandler)

	// Set initial focus
	t.setFocus("workspace_list")

	// Set root primitive
	t.app.SetRoot(t.mainLayout, true)

	return nil
}

// globalKeyHandler handles global keyboard shortcuts.
func (t *TUIApp) globalKeyHandler(event *tcell.EventKey) *tcell.EventKey {
	// Check if help view is active
	currentView := t.router.Current()
	if currentView != nil && currentView.Name() == "help" {
		// '?' or Esc to close help and return to main view
		if event.Rune() == '?' || event.Key() == tcell.KeyEsc {
			t.app.SetRoot(t.mainLayout, true)
			t.updateFocus()
			return nil
		}
	} else {
		// Main view key bindings
		switch event.Key() {
		case tcell.KeyCtrlC:
			t.app.Stop()
			return nil
		case tcell.KeyTab:
			t.cycleFocus()
			return nil
		case tcell.KeyBacktab:
			t.cycleFocusReverse()
			return nil
		case tcell.KeyEsc:
			// Context-aware Escape: detail → tree → workspace
			if t.currentFocus == "ticket_detail" {
				t.setFocus("ticket_tree")
				return nil
			}
			if t.currentFocus == "ticket_tree" {
				t.setFocus("workspace_list")
				return nil
			}
			// From workspace list, do nothing
			return event
		}

		switch event.Rune() {
		case 'q':
			t.app.Stop()
			return nil
		case '?':
			// Show help view
			_ = t.router.Show("help")
			return nil
		case '/':
			// Show search view (Week 14)
			t.showSearch()
			return nil
		case ':':
			// Show command palette (Week 14)
			t.showCommandPalette()
			return nil
		}
	}

	return event
}

// cycleFocus moves focus forward through panels.
func (t *TUIApp) cycleFocus() {
	for i, name := range t.focusOrder {
		if name == t.currentFocus {
			nextIndex := (i + 1) % len(t.focusOrder)
			t.setFocus(t.focusOrder[nextIndex])
			return
		}
	}
}

// cycleFocusReverse moves focus backward through panels.
func (t *TUIApp) cycleFocusReverse() {
	for i, name := range t.focusOrder {
		if name == t.currentFocus {
			prevIndex := (i - 1 + len(t.focusOrder)) % len(t.focusOrder)
			t.setFocus(t.focusOrder[prevIndex])
			return
		}
	}
}

// setFocus updates the currently focused panel.
func (t *TUIApp) setFocus(viewName string) {
	t.currentFocus = viewName
	t.updateFocus()
}

// updateFocus applies the current focus state to the UI.
func (t *TUIApp) updateFocus() {
	// Update border colors for all views
	t.workspaceListView.SetFocused(t.currentFocus == "workspace_list")
	t.ticketTreeView.SetFocused(t.currentFocus == "ticket_tree")
	t.ticketDetailView.SetFocused(t.currentFocus == "ticket_detail")

	// Set application focus
	switch t.currentFocus {
	case "workspace_list":
		t.app.SetFocus(t.workspaceListView.Primitive())
	case "ticket_tree":
		t.app.SetFocus(t.ticketTreeView.Primitive())
	case "ticket_detail":
		t.app.SetFocus(t.ticketDetailView.Primitive())
	}
}

// Stop stops the TUI application.
func (t *TUIApp) Stop() {
	t.app.Stop()
}

// Router returns the view router for testing.
func (t *TUIApp) Router() *Router {
	return t.router
}

// setupCommands populates the command palette with available commands.
func (t *TUIApp) setupCommands() {
	commands := []views.Command{
		{
			Name:        "push",
			Description: "Push tickets to Jira (not yet implemented)",
			Action: func() error {
				// TODO: Implement push in future phase
				return nil
			},
		},
		{
			Name:        "pull",
			Description: "Pull tickets from Jira (not yet implemented)",
			Action: func() error {
				// TODO: Implement pull in future phase
				return nil
			},
		},
		{
			Name:        "refresh",
			Description: "Refresh current workspace tickets",
			Action: func() error {
				// Reload tickets for current workspace
				if ws, err := t.workspaceService.Current(); err == nil && ws != nil {
					t.ticketTreeView.LoadTickets(ws.ID)
				}
				return nil
			},
		},
		{
			Name:        "help",
			Description: "Show help",
			Action: func() error {
				return t.router.Show("help")
			},
		},
		{
			Name:        "quit",
			Description: "Quit application",
			Action: func() error {
				t.app.Stop()
				return nil
			},
		},
	}
	t.commandView.SetCommands(commands)
}

// showSearch displays the search modal with current workspace tickets.
func (t *TUIApp) showSearch() {
	// Get tickets from current workspace
	ws, err := t.workspaceService.Current()
	if err == nil && ws != nil {
		tickets, err := t.ticketQuery.ListByWorkspace(ws.ID)
		if err == nil && len(tickets) > 0 {
			// Convert []domain.Ticket to []*domain.Ticket
			ticketPtrs := make([]*domain.Ticket, len(tickets))
			for i := range tickets {
				ticketPtrs[i] = &tickets[i]
			}
			t.searchView.SetTickets(ticketPtrs)
		}
	}

	t.searchView.OnShow()
	t.app.SetRoot(t.searchView.Primitive(), true)
	t.app.SetFocus(t.searchView.Primitive())
}

// showCommandPalette displays the command palette modal.
func (t *TUIApp) showCommandPalette() {
	t.commandView.OnShow()
	t.app.SetRoot(t.commandView.Primitive(), true)
	t.app.SetFocus(t.commandView.Primitive())
}
