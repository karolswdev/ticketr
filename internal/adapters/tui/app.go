package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/karolswdev/ticktr/internal/adapters/tui/views"
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

	// View references for multi-panel layout
	workspaceListView *views.WorkspaceListView
	ticketTreeView    *views.TicketTreeView

	// Focus management
	currentFocus string // "workspace_list" or "ticket_tree"
}

// NewTUIApp creates a new TUI application instance.
func NewTUIApp(
	workspaceService *services.WorkspaceService,
	ticketQuery      *services.TicketQueryService,
	pathResolver     *services.PathResolver,
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

	return &TUIApp{
		app:              tview.NewApplication(),
		router:           NewRouter(),
		workspaceService: workspaceService,
		ticketQuery:      ticketQuery,
		pathResolver:     pathResolver,
		currentFocus:     "workspace_list", // Start with workspace list focused
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

	// Set workspace change callback to reload tickets
	t.workspaceListView.SetWorkspaceChangeHandler(func(workspaceID string) {
		if t.ticketTreeView != nil {
			t.ticketTreeView.LoadTickets(workspaceID)
		}
	})

	// Create help view and register with router
	helpView := views.NewHelpView()
	if err := t.router.Register(helpView); err != nil {
		return fmt.Errorf("failed to register help view: %w", err)
	}

	// Create multi-panel layout
	t.mainLayout = tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(t.workspaceListView.Primitive(), 35, 0, true).  // Left panel, 35 chars fixed
		AddItem(t.ticketTreeView.Primitive(), 0, 1, false)       // Right panel, flex fill

	// Set up global key bindings
	t.app.SetInputCapture(t.globalKeyHandler)

	// Set initial focus
	t.updateFocus()

	// Set root primitive
	t.app.SetRoot(t.mainLayout, true)

	return nil
}

// globalKeyHandler handles global keyboard shortcuts.
func (t *TUIApp) globalKeyHandler(event *tcell.EventKey) *tcell.EventKey {
	// Check if help view is active
	currentView := t.router.Current()
	if currentView != nil && currentView.Name() == "help" {
		// '?' to close help and return to main view
		if event.Rune() == '?' {
			t.app.SetRoot(t.mainLayout, true)
			return nil
		}
	} else {
		// Main view key bindings
		switch event.Key() {
		case tcell.KeyCtrlC:
			t.app.Stop()
			return nil
		case tcell.KeyTab:
			t.toggleFocus()
			return nil
		}

		switch event.Rune() {
		case 'q':
			t.app.Stop()
			return nil
		case '?':
			// Show help view
			_ = t.router.Show("help")
			return nil
		}
	}

	return event
}

// toggleFocus switches focus between workspace list and ticket tree.
func (t *TUIApp) toggleFocus() {
	if t.currentFocus == "workspace_list" {
		t.currentFocus = "ticket_tree"
	} else {
		t.currentFocus = "workspace_list"
	}
	t.updateFocus()
}

// updateFocus applies the current focus state to the UI.
func (t *TUIApp) updateFocus() {
	if t.currentFocus == "workspace_list" {
		// Workspace list focused
		t.workspaceListView.Primitive().(*tview.List).SetBorderColor(tcell.ColorGreen)
		t.ticketTreeView.Primitive().(*tview.TreeView).SetBorderColor(tcell.ColorWhite)
		t.app.SetFocus(t.workspaceListView.Primitive())
	} else {
		// Ticket tree focused
		t.workspaceListView.Primitive().(*tview.List).SetBorderColor(tcell.ColorWhite)
		t.ticketTreeView.Primitive().(*tview.TreeView).SetBorderColor(tcell.ColorGreen)
		t.app.SetFocus(t.ticketTreeView.Primitive())
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
