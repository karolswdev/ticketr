package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/karolswdev/ticktr/internal/adapters/tui/views"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/rivo/tview"
)

// TUIApp represents the Terminal User Interface application.
type TUIApp struct {
	app          *tview.Application
	workspace    *services.WorkspaceService
	pathResolver *services.PathResolver
	router       *Router
	keyHandler   *KeyHandler
	previousView string // Track previous view for help toggle
}

// NewTUIApp creates a new TUI application.
func NewTUIApp(workspaceService *services.WorkspaceService, pathResolver *services.PathResolver) (*TUIApp, error) {
	if workspaceService == nil {
		return nil, fmt.Errorf("workspace service is required")
	}
	if pathResolver == nil {
		return nil, fmt.Errorf("path resolver is required")
	}

	tuiApp := &TUIApp{
		app:          tview.NewApplication(),
		workspace:    workspaceService,
		pathResolver: pathResolver,
		router:       NewRouter(),
	}

	// Initialize key handler
	tuiApp.keyHandler = NewKeyHandler(tuiApp, tuiApp.router)

	// Initialize views
	if err := tuiApp.initializeViews(); err != nil {
		return nil, fmt.Errorf("failed to initialize views: %w", err)
	}

	// Set up the application
	tuiApp.setupApp()

	return tuiApp, nil
}

// initializeViews creates and registers all views.
func (t *TUIApp) initializeViews() error {
	// Create workspace list view
	workspaceList := views.NewWorkspaceListView(t.workspace)
	workspaceList.SetSwitchHandler(func(name string) error {
		return t.workspace.Switch(name)
	})
	if err := t.router.Register(workspaceList); err != nil {
		return err
	}

	// Create ticket tree view
	ticketTree := views.NewTicketTreeView()
	if err := t.router.Register(ticketTree); err != nil {
		return err
	}

	// Create ticket detail view
	ticketDetail := views.NewTicketDetailView()
	if err := t.router.Register(ticketDetail); err != nil {
		return err
	}

	// Create help view
	helpView := views.NewHelpView()
	if err := t.router.Register(helpView); err != nil {
		return err
	}

	return nil
}

// setupApp configures the tview application.
func (t *TUIApp) setupApp() {
	// Set root to pages managed by router
	t.app.SetRoot(t.router.Pages(), true)

	// Set global input capture for keybindings
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Let key handler process global shortcuts
		if t.keyHandler.Handle(event) {
			return nil
		}
		// Pass through to current view
		return event
	})

	// Show initial view (workspace list)
	_ = t.router.Show("workspace_list")
}

// Run starts the TUI application.
func (t *TUIApp) Run() error {
	return t.app.Run()
}

// Stop stops the TUI application.
func (t *TUIApp) Stop() {
	t.app.Stop()
}

// Router returns the view router for testing.
func (t *TUIApp) Router() *Router {
	return t.router
}
