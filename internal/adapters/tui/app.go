package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/karolswdev/ticktr/internal/adapters/tui/sync"
	"github.com/karolswdev/ticktr/internal/adapters/tui/theme"
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

	// Sync services (Week 15)
	pushService     *services.PushService
	pullService     *services.PullService
	syncCoordinator *sync.SyncCoordinator

	// View references for tri-panel layout
	workspaceListView *views.WorkspaceListView
	ticketTreeView    *views.TicketTreeView
	ticketDetailView  *views.TicketDetailView
	syncStatusView    *views.SyncStatusView

	// Modal views (Week 14)
	searchView  *views.SearchView
	commandView *views.CommandPaletteView

	// Focus management
	currentFocus string   // "workspace_list", "ticket_tree", or "ticket_detail"
	focusOrder   []string // Order of focus cycling

	// Modal state tracking
	inModal bool // True when a modal view (search, command palette) is active
}

// NewTUIApp creates a new TUI application instance.
func NewTUIApp(
	workspaceService *services.WorkspaceService,
	ticketQuery *services.TicketQueryService,
	pathResolver *services.PathResolver,
	pushService *services.PushService,
	pullService *services.PullService,
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
	if pushService == nil {
		return nil, fmt.Errorf("push service is required")
	}
	if pullService == nil {
		return nil, fmt.Errorf("pull service is required")
	}

	app := tview.NewApplication()

	tuiApp := &TUIApp{
		app:              app,
		router:           NewRouter(),
		workspaceService: workspaceService,
		ticketQuery:      ticketQuery,
		pathResolver:     pathResolver,
		pushService:      pushService,
		pullService:      pullService,
		currentFocus:     "workspace_list",
		focusOrder:       []string{"workspace_list", "ticket_tree", "ticket_detail"},
	}

	// Create sync coordinator with status callback
	tuiApp.syncCoordinator = sync.NewSyncCoordinator(
		pushService,
		pullService,
		tuiApp.onSyncStatusChanged,
	)

	return tuiApp, nil
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
	// Apply default theme
	theme.Apply(t.app)

	// Check terminal size and show warning if needed
	if err := t.checkTerminalSize(); err != nil {
		return err
	}

	// Create workspace list view
	t.workspaceListView = views.NewWorkspaceListView(t.workspaceService)
	t.workspaceListView.SetSwitchHandler(func(name string) error {
		return t.workspaceService.Switch(name)
	})
	// Load workspaces on startup
	t.workspaceListView.OnShow()

	// Create ticket tree view
	t.ticketTreeView = views.NewTicketTreeView(t.workspaceService, t.ticketQuery)

	// Create ticket detail view
	t.ticketDetailView = views.NewTicketDetailView(t.app)

	// Create sync status view (Week 15)
	t.syncStatusView = views.NewSyncStatusView()

	// Load tickets for current workspace on startup
	if ws, err := t.workspaceService.Current(); err == nil && ws != nil {
		t.ticketTreeView.LoadTickets(ws.ID)
	}

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
		t.inModal = false
		t.app.SetRoot(t.mainLayout, true)
		t.updateFocus()
	})
	t.searchView.SetOnSelect(func(ticket *domain.Ticket) {
		// Return to main layout and show ticket detail
		t.inModal = false
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
		t.inModal = false
		t.app.SetRoot(t.mainLayout, true)
		t.updateFocus()
	})
	// Set up available commands
	t.setupCommands()

	// Create right panel (detail view + status bar)
	rightPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.ticketDetailView.Primitive(), 0, 1, false). // Detail takes most space
		AddItem(t.syncStatusView.Primitive(), 3, 0, false)    // Status bar (3 rows fixed)

	// Create layout based on terminal size
	t.mainLayout = t.createResponsiveLayout(rightPanel)

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
			// Clear router's current view state
			t.router.ClearCurrent()
			// Return to main layout
			t.app.SetRoot(t.mainLayout, true)
			t.updateFocus()
			return nil
		}
	} else if !t.inModal {
		// Main view key bindings (ONLY when main layout is active)
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
			if err := t.router.Show("help"); err == nil {
				// Switch app root to router pages to display help
				t.app.SetRoot(t.router.Pages(), true)
			}
			return nil
		case '/':
			// Show search view (Week 14)
			t.showSearch()
			return nil
		case ':':
			// Show command palette (Week 14)
			t.showCommandPalette()
			return nil
		case 'p':
			// Push tickets to Jira (Week 15)
			t.handlePush()
			return nil
		case 'P':
			// Pull tickets from Jira (Week 15)
			t.handlePull()
			return nil
		case 'r':
			// Refresh current workspace tickets (Week 15)
			t.handleRefresh()
			return nil
		case 's':
			// Full sync: pull then push (Week 15)
			t.handleSync()
			return nil
		}
	} else {
		// In a modal view - only handle Ctrl+C to quit
		if event.Key() == tcell.KeyCtrlC {
			t.app.Stop()
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
	// Update border colors for all views using theme colors
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
			Description: "Push tickets to Jira",
			Action: func() error {
				t.handlePush()
				return nil
			},
		},
		{
			Name:        "pull",
			Description: "Pull tickets from Jira",
			Action: func() error {
				t.handlePull()
				return nil
			},
		},
		{
			Name:        "sync",
			Description: "Full sync (pull then push)",
			Action: func() error {
				t.handleSync()
				return nil
			},
		},
		{
			Name:        "refresh",
			Description: "Refresh current workspace tickets",
			Action: func() error {
				t.handleRefresh()
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

	t.inModal = true
	t.searchView.OnShow()
	t.app.SetRoot(t.searchView.Primitive(), true)
	t.app.SetFocus(t.searchView.Primitive())
}

// showCommandPalette displays the command palette modal.
func (t *TUIApp) showCommandPalette() {
	t.inModal = true
	t.commandView.OnShow()
	t.app.SetRoot(t.commandView.Primitive(), true)
	t.app.SetFocus(t.commandView.Primitive())
}

// onSyncStatusChanged is called when sync status changes (from sync coordinator).
func (t *TUIApp) onSyncStatusChanged(status sync.SyncStatus) {
	// Update UI from goroutine using QueueUpdateDraw
	t.app.QueueUpdateDraw(func() {
		// Update status view
		t.syncStatusView.SetStatus(status)

		// If sync completed successfully, reload tickets
		if status.State == sync.StateSuccess {
			if ws, err := t.workspaceService.Current(); err == nil && ws != nil {
				t.ticketTreeView.LoadTickets(ws.ID)
			}
		}
	})
}

// handlePush initiates an async push operation.
func (t *TUIApp) handlePush() {
	// Get current workspace
	ws, err := t.workspaceService.Current()
	if err != nil || ws == nil {
		// Set error status
		t.syncStatusView.SetStatus(sync.NewErrorStatus("push", fmt.Errorf("no active workspace")))
		return
	}

	// For now, use tickets.md in the current working directory
	// TODO: In future, integrate with workspace file path configuration
	filePath := "tickets.md"

	// Start async push
	t.syncCoordinator.PushAsync(filePath, services.ProcessOptions{})
}

// handlePull initiates an async pull operation.
func (t *TUIApp) handlePull() {
	// Get current workspace
	ws, err := t.workspaceService.Current()
	if err != nil || ws == nil {
		t.syncStatusView.SetStatus(sync.NewErrorStatus("pull", fmt.Errorf("no active workspace")))
		return
	}

	// For now, use tickets.md in the current working directory
	// TODO: In future, integrate with workspace file path configuration
	filePath := "tickets.md"

	// Start async pull with workspace project key
	t.syncCoordinator.PullAsync(filePath, services.PullOptions{
		ProjectKey: ws.ProjectKey,
	})
}

// handleSync initiates an async full sync (pull then push).
func (t *TUIApp) handleSync() {
	// Get current workspace
	ws, err := t.workspaceService.Current()
	if err != nil || ws == nil {
		t.syncStatusView.SetStatus(sync.NewErrorStatus("sync", fmt.Errorf("no active workspace")))
		return
	}

	// For now, use tickets.md in the current working directory
	// TODO: In future, integrate with workspace file path configuration
	filePath := "tickets.md"

	// Start async sync
	t.syncCoordinator.SyncAsync(filePath)
}

// checkTerminalSize checks if terminal is large enough and shows warning if not.
// Note: We'll do a best-effort check during initial layout. Terminal size can't be
// reliably checked before tview initializes the screen.
func (t *TUIApp) checkTerminalSize() error {
	// This will be checked after screen initialization
	// For now, just return nil - we'll handle small terminals gracefully
	// by using a responsive layout
	return nil
}

// createResponsiveLayout creates layout based on terminal width.
// Note: This is a best-effort approach. We use the full layout by default
// and users can resize their terminal to see the compact version in future updates.
func (t *TUIApp) createResponsiveLayout(rightPanel *tview.Flex) *tview.Flex {
	// For now, always create full layout
	// In a future enhancement, we could add a resize handler to switch layouts dynamically
	// tview makes it difficult to query terminal size before initialization

	// TODO: Add dynamic layout switching on terminal resize events
	// This would require hooking into tview's screen resize events

	return t.createFullLayout(rightPanel)
}

// createFullLayout creates the standard tri-panel layout.
func (t *TUIApp) createFullLayout(rightPanel *tview.Flex) *tview.Flex {
	return tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(t.workspaceListView.Primitive(), 30, 0, true). // Fixed 30 chars
		AddItem(t.ticketTreeView.Primitive(), 0, 2, false).    // 40% (2 of 5 parts)
		AddItem(rightPanel, 0, 3, false)                       // 60% (3 of 5 parts)
}

// handleRefresh reloads tickets for the current workspace.
func (t *TUIApp) handleRefresh() {
	// Get current workspace
	ws, err := t.workspaceService.Current()
	if err != nil || ws == nil {
		t.syncStatusView.SetStatus(sync.NewErrorStatus("refresh", fmt.Errorf("no active workspace")))
		return
	}

	// Update status to show refreshing
	t.syncStatusView.SetStatus(sync.NewSyncingStatus("refresh", "Reloading tickets..."))

	// Reload tickets
	t.ticketTreeView.LoadTickets(ws.ID)

	// Update status to success
	t.syncStatusView.SetStatus(sync.NewSuccessStatus("refresh", "Tickets reloaded"))
}
