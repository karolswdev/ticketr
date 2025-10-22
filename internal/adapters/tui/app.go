package tui

import (
	"fmt"
	"os"
	"os/signal"
	stdSync "sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/karolswdev/ticktr/internal/adapters/database"
	"github.com/karolswdev/ticktr/internal/adapters/filesystem"
	"github.com/karolswdev/ticktr/internal/adapters/jira"
	"github.com/karolswdev/ticktr/internal/adapters/tui/commands"
	"github.com/karolswdev/ticktr/internal/adapters/tui/effects"
	"github.com/karolswdev/ticktr/internal/adapters/tui/sync"
	"github.com/karolswdev/ticktr/internal/adapters/tui/theme"
	"github.com/karolswdev/ticktr/internal/adapters/tui/views"
	"github.com/karolswdev/ticktr/internal/adapters/tui/widgets"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/karolswdev/ticktr/internal/state"
	"github.com/karolswdev/ticktr/internal/tui/jobs"
	"github.com/rivo/tview"
)

// TUIApp represents the main TUI application.
type TUIApp struct {
	app              *tview.Application
	mainLayout       *tview.Flex
	pages            *tview.Pages // For slide-out overlay management
	router           *Router
	workspaceService *services.WorkspaceService
	ticketQuery      *services.TicketQueryService
	pathResolver     *services.PathResolver

	// Sync services (Week 15)
	// CRITICAL (Phase 6.5 Emergency Fix): serviceMutex protects service fields from race conditions.
	// During workspace changes, recreateJiraAdapter() replaces these services.
	// Without this mutex, pull/push jobs could access replaced/stale services causing crashes.
	serviceMutex    stdSync.RWMutex
	pushService     *services.PushService
	pullService     *services.PullService
	syncCoordinator *sync.SyncCoordinator

	// Async job queue (Phase 6, Week 2)
	jobQueue       *jobs.JobQueue
	currentJobID   jobs.JobID // Track active job for cancellation
	currentJobType string     // "pull", "push", "sync"

	// Bulk operations service (Week 18)
	bulkOperationService ports.BulkOperationService

	// Visual effects (Phase 6, Day 12.5)
	animator           *effects.Animator
	backgroundAnimator *effects.BackgroundAnimator

	// View references for 2-panel layout (Phase 6.6)
	workspaceListView *views.WorkspaceListView
	ticketTreeView    *views.TicketTreeView
	ticketDetailView  *views.TicketDetailView
	syncStatusView    *views.SyncStatusView

	// Workspace slide-out (Phase 6.6)
	workspaceSlideOut *widgets.SlideOut

	// Modal views (Week 14)
	searchView          *views.SearchView
	commandView         *views.CommandPaletteView
	workspaceModal      *views.WorkspaceModal
	bulkOperationsModal *views.BulkOperationsModal

	// Phase 6 widgets (Day 8-9)
	actionBar       *widgets.ActionBar
	commandRegistry *commands.Registry
	commandPalette  *widgets.CommandPalette

	// Focus management
	currentFocus string   // "ticket_tree" or "ticket_detail" (workspace removed from default)
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
	bulkOperationService ports.BulkOperationService,
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
	if bulkOperationService == nil {
		return nil, fmt.Errorf("bulk operation service is required")
	}

	app := tview.NewApplication()

	tuiApp := &TUIApp{
		app:                  app,
		router:               NewRouter(),
		workspaceService:     workspaceService,
		ticketQuery:          ticketQuery,
		pathResolver:         pathResolver,
		pushService:          pushService,
		pullService:          pullService,
		bulkOperationService: bulkOperationService,
		currentFocus:         "ticket_tree",
		focusOrder:           []string{"ticket_tree", "ticket_detail"}, // Phase 6.6: workspace removed from default
	}

	// Create sync coordinator with status callback
	tuiApp.syncCoordinator = sync.NewSyncCoordinator(
		pushService,
		pullService,
		tuiApp.onSyncStatusChanged,
	)

	// Create job queue (single worker for now)
	tuiApp.jobQueue = jobs.NewJobQueue(1)

	// Set up signal handler for Ctrl+C
	tuiApp.setupSignalHandler()

	return tuiApp, nil
}

// Run starts the TUI application.
func (t *TUIApp) Run() error {
	if err := t.setupApp(); err != nil {
		return err
	}

	// Initialize and start animator if effects are enabled
	effectsConfig := theme.GetEffects()
	if effectsConfig.Motion {
		t.animator = effects.NewAnimator(t.app)
		// Start the animator to enable visual effects
		// Note: Individual animations are started on-demand by views
	}

	// Initialize background animator if ambient effects are enabled (Phase 6.5)
	fmt.Fprintf(os.Stderr, "[DEBUG] effectsConfig.AmbientEnabled = %v\n", effectsConfig.AmbientEnabled)
	fmt.Fprintf(os.Stderr, "[DEBUG] effectsConfig.AmbientMode = %s\n", effectsConfig.AmbientMode)
	fmt.Fprintf(os.Stderr, "[DEBUG] effectsConfig.AmbientDensity = %f\n", effectsConfig.AmbientDensity)
	fmt.Fprintf(os.Stderr, "[DEBUG] effectsConfig.AmbientSpeed = %d\n", effectsConfig.AmbientSpeed)

	if effectsConfig.AmbientEnabled {
		fmt.Fprintf(os.Stderr, "[DEBUG] Ambient effects ENABLED - creating background animator\n")

		// Map theme ambient mode to background effect type
		var effectType effects.BackgroundEffect
		switch effectsConfig.AmbientMode {
		case "hyperspace":
			effectType = effects.BackgroundHyperspace
			fmt.Fprintf(os.Stderr, "[DEBUG] Effect type: hyperspace\n")
		case "snow":
			effectType = effects.BackgroundSnow
			fmt.Fprintf(os.Stderr, "[DEBUG] Effect type: snow\n")
		default:
			effectType = effects.BackgroundNone
			fmt.Fprintf(os.Stderr, "[DEBUG] Effect type: none (default)\n")
		}

		// Create background animator with theme configuration
		bgConfig := effects.BackgroundConfig{
			Effect:    effectType,
			Density:   effectsConfig.AmbientDensity,
			Speed:     effectsConfig.AmbientSpeed,
			Enabled:   true,
			MaxFPS:    15,
			AutoPause: true,
		}

		fmt.Fprintf(os.Stderr, "[DEBUG] Creating BackgroundAnimator with config: %+v\n", bgConfig)
		t.backgroundAnimator = effects.NewBackgroundAnimator(t.app, bgConfig)

		fmt.Fprintf(os.Stderr, "[DEBUG] Starting BackgroundAnimator...\n")
		t.backgroundAnimator.Start()

		fmt.Fprintf(os.Stderr, "[DEBUG] Getting overlay from animator...\n")
		overlay := t.backgroundAnimator.GetOverlay()
		fmt.Fprintf(os.Stderr, "[DEBUG] Overlay is nil: %v\n", overlay == nil)

		// Add cosmic background to workspace slide-out now that animator exists
		if t.workspaceSlideOut != nil {
			fmt.Fprintf(os.Stderr, "[DEBUG] Setting background overlay on workspace slide-out\n")
			t.workspaceSlideOut.SetBackgroundOverlay(overlay)
		} else {
			fmt.Fprintf(os.Stderr, "[DEBUG] WARNING: workspaceSlideOut is nil!\n")
		}
	} else {
		fmt.Fprintf(os.Stderr, "[DEBUG] Ambient effects DISABLED - skipping background animator\n")
	}

	// Start progress monitoring goroutine
	go t.monitorJobProgress()

	return t.app.Run()
}

// setupSignalHandler sets up Ctrl+C signal handling for graceful shutdown.
func (t *TUIApp) setupSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		<-sigChan
		// Cancel active job if any
		if t.currentJobID != "" {
			t.jobQueue.Cancel(t.currentJobID)
		}
		// Shutdown job queue gracefully
		t.jobQueue.Shutdown()
		// Shutdown animator if initialized
		if t.animator != nil {
			t.animator.Shutdown()
		}
		// Shutdown background animator if initialized (Phase 6.5)
		if t.backgroundAnimator != nil {
			t.backgroundAnimator.Shutdown()
		}
		// Stop the TUI
		t.app.Stop()
	}()
}

// monitorJobProgress listens for progress updates from the job queue and updates the UI.
func (t *TUIApp) monitorJobProgress() {
	for progress := range t.jobQueue.Progress() {
		t.app.QueueUpdateDraw(func() {
			// Update status view with progress bar (Phase 6, Day 10-11)
			status := sync.NewSyncingStatus(t.currentJobType, jobs.FormatProgress(progress))
			t.syncStatusView.SetStatus(status)
			t.syncStatusView.UpdateProgress(progress)
		})
	}
}

// setupApp initializes all views and layouts.
func (t *TUIApp) setupApp() error {
	// Load theme configuration from environment variables
	theme.LoadThemeFromEnv()

	// Apply theme (pass nil to avoid calling Draw() before app is ready)
	theme.Apply(nil)

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

	// Set workspace creation and profile management handlers
	t.workspaceListView.SetCreateWorkspaceHandler(func() {
		t.showWorkspaceModal()
	})
	t.workspaceListView.SetManageProfilesHandler(func() {
		t.showProfileManagement()
	})

	// Create ticket tree view with app reference for async updates
	t.ticketTreeView = views.NewTicketTreeView(t.workspaceService, t.ticketQuery, t.app)

	// Create ticket detail view
	t.ticketDetailView = views.NewTicketDetailView(t.app)

	// Create sync status view (Week 15) with animation support
	t.syncStatusView = views.NewSyncStatusView(t.app)

	// Create action bar widget (Phase 6, Day 8-9)
	t.actionBar = widgets.NewActionBar()
	t.actionBar.SetApp(t.app) // Enable marquee updates
	t.actionBar.SetContext(widgets.ContextWorkspaceList)

	// Create command registry and populate with commands (Phase 6, Day 8-9)
	t.commandRegistry = commands.NewRegistry()
	t.setupCommandRegistry()

	// Create enhanced command palette with registry (Phase 6, Day 8-9)
	t.commandPalette = widgets.NewCommandPalette(t.commandRegistry)
	t.commandPalette.SetOnClose(func() {
		t.inModal = false
		t.app.SetRoot(t.mainLayout, true)
		t.updateFocus()
	})

	// Set workspace change callback to reload tickets
	t.workspaceListView.SetWorkspaceChangeHandler(func(workspaceID string) {
		if t.ticketTreeView != nil {
			t.ticketTreeView.LoadTickets(workspaceID)
			t.setFocus("ticket_tree")
		}

		// FIX (Phase 6.5 Critical): Recreate Jira adapter with new workspace credentials
		// This fixes the bug where adapter was created once at startup and never updated,
		// causing pull/push operations to fail after workspace switch
		if err := t.recreateJiraAdapter(); err != nil {
			// Show error in status bar
			t.syncStatusView.SetStatus(sync.NewErrorStatus("workspace", err))
		}

		// FIX #2: Auto-close workspace panel after selection
		if t.workspaceSlideOut != nil && t.workspaceSlideOut.IsVisible() {
			t.toggleWorkspacePanel()
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

	// Create bulk operations modal (Milestone 18 - Slice 4)
	t.bulkOperationsModal = views.NewBulkOperationsModal(t.app, t.bulkOperationService)
	t.bulkOperationsModal.SetOnClose(func() {
		// Return to main layout
		t.inModal = false
		t.app.SetRoot(t.mainLayout, true)
		t.updateFocus()
	})
	t.bulkOperationsModal.SetOnSuccess(func() {
		// Refresh ticket tree to show updated tickets
		if ws, err := t.workspaceService.Current(); err == nil && ws != nil {
			t.ticketTreeView.LoadTicketsAsync(ws.ID)
		}
		// Clear selection after successful bulk operation
		t.ticketTreeView.ClearSelection()
	})

	// Create right panel (detail view + status bar)
	rightPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.ticketDetailView.Primitive(), 0, 1, false). // Detail takes most space
		AddItem(t.syncStatusView.Primitive(), 3, 0, false)    // Status bar (3 rows fixed)

	// Create main content layout - 2-panel (Phase 6.6: workspace removed)
	contentLayout := t.create2PanelLayout(rightPanel)

	// Create main layout with action bar at bottom (Phase 6, Day 8-9)
	t.mainLayout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(contentLayout, 0, 1, true).           // Content takes most space
		AddItem(t.actionBar.Primitive(), 3, 0, false) // Action bar fixed height

	// Create workspace slide-out (Phase 6.6)
	t.workspaceSlideOut = widgets.NewSlideOut(t.workspaceListView.Primitive(), 35)
	t.workspaceSlideOut.SetOnClose(func() {
		// Return focus to main view
		t.updateFocus()
	})

	// Note: Cosmic background will be added in Run() after animator initialization (Phase 6.5)

	// Create pages for overlay management (Phase 6.6)
	t.pages = tview.NewPages().
		AddPage("main", t.mainLayout, true, true).
		AddPage("workspace-overlay", t.workspaceSlideOut.Primitive(), true, false) // Hidden initially

	// CRITICAL FIX (Phase 6.5): Create workspace modal AFTER pages is initialized
	// Previously, modal was created with nil pages, causing ESC handling to fail
	// and app to become unresponsive when closing modal
	t.workspaceModal = views.NewWorkspaceModal(t.app, t.pages, t.workspaceService)
	t.workspaceModal.SetOnClose(func() {
		// Hide modal page overlay
		t.inModal = false
		if t.pages != nil {
			t.pages.HidePage("workspace-modal")
		}
		// Restore main layout root and focus
		t.app.SetRoot(t.pages, true)
		t.updateFocus()
	})
	t.workspaceModal.SetOnSuccess(func() {
		// Refresh workspace list to show new workspace
		t.workspaceListView.OnShow()
	})

	// Set up global key bindings
	t.app.SetInputCapture(t.globalKeyHandler)

	// Set initial focus
	t.setFocus("ticket_tree")

	// Set root primitive to pages (Phase 6.6)
	t.app.SetRoot(t.pages, true)

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
		// 'q' or F10 to quit from help view
		if event.Rune() == 'q' || event.Key() == tcell.KeyF10 {
			t.app.Stop()
			return nil
		}
	} else if !t.inModal {
		// Main view key bindings (ONLY when main layout is active)
		switch event.Key() {
		case tcell.KeyCtrlC:
			t.app.Stop()
			return nil
		case tcell.KeyCtrlP:
			// Show enhanced command palette (Phase 6, Day 8-9)
			t.showEnhancedCommandPalette()
			return nil
		case tcell.KeyF1:
			// F1: Show help or command palette (Phase 6, Day 8-9)
			t.showEnhancedCommandPalette()
			return nil
		case tcell.KeyF2:
			// F2: Pull/Sync from Jira (Phase 6, Day 8-9)
			t.handlePull()
			return nil
		case tcell.KeyF3:
			// F3: Toggle workspace panel (Phase 6.6)
			t.toggleWorkspacePanel()
			return nil
		case tcell.KeyF5:
			// F5: Refresh view (Phase 6, Day 8-9)
			t.handleRefresh()
			return nil
		case tcell.KeyF10:
			// F10: Exit (Phase 6, Day 8-9)
			t.app.Stop()
			return nil
		case tcell.KeyTab:
			t.cycleFocus()
			return nil
		case tcell.KeyBacktab:
			t.cycleFocusReverse()
			return nil
		case tcell.KeyEsc:
			// Priority 1: Close workspace panel if open (Phase 6.6)
			if t.workspaceSlideOut != nil && t.workspaceSlideOut.IsVisible() {
				t.toggleWorkspacePanel()
				return nil
			}
			// Priority 2: Cancel active job if any
			if t.currentJobID != "" {
				t.handleJobCancellation()
				return nil
			}
			// Priority 3: Context-aware navigation: detail → tree
			if t.currentFocus == "ticket_detail" {
				t.setFocus("ticket_tree")
				return nil
			}
			// From ticket tree, do nothing
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
		case 'W':
			// Toggle workspace panel (Phase 6.6)
			t.toggleWorkspacePanel()
			return nil
		case '/':
			// Show search view (Week 14)
			t.showSearch()
			return nil
		case ':':
			// Show command palette (Week 14)
			t.showCommandPalette()
			return nil
		case 'b':
			// Show bulk operations menu (Week 18)
			t.handleBulkOperations()
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
	// Update border colors for all views using theme colors (Phase 6.6: workspace handled separately)
	t.ticketTreeView.SetFocused(t.currentFocus == "ticket_tree")
	t.ticketDetailView.SetFocused(t.currentFocus == "ticket_detail")

	// Update action bar context based on focus (Phase 6, Day 8-9)
	if t.actionBar != nil {
		var ctx widgets.ActionBarContext
		if t.currentJobID != "" {
			ctx = widgets.ContextSyncing
		} else {
			switch t.currentFocus {
			case "ticket_tree":
				ctx = widgets.ContextTicketTree
			case "ticket_detail":
				ctx = widgets.ContextTicketDetail
			default:
				ctx = widgets.ContextTicketTree
			}
		}
		t.actionBar.SetContext(ctx)
	}

	// Set application focus
	switch t.currentFocus {
	case "ticket_tree":
		t.app.SetFocus(t.ticketTreeView.Primitive())
	case "ticket_detail":
		t.app.SetFocus(t.ticketDetailView.Primitive())
	}
}

// Stop stops the TUI application.
func (t *TUIApp) Stop() {
	// Shutdown animator if initialized
	if t.animator != nil {
		t.animator.Shutdown()
	}
	// Shutdown background animator if initialized (Phase 6.5)
	if t.backgroundAnimator != nil {
		t.backgroundAnimator.Shutdown()
	}
	// Shutdown action bar (cleans up marquee)
	if t.actionBar != nil {
		t.actionBar.Shutdown()
	}
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

// setupCommandRegistry populates the command registry with all available commands (Phase 6, Day 8-9).
func (t *TUIApp) setupCommandRegistry() {
	// Navigation commands
	t.commandRegistry.Register(&commands.Command{
		Name:        "help",
		Description: "Show help screen",
		Keybinding:  "? or F1",
		Category:    commands.CategoryNav,
		Handler: func() error {
			if err := t.router.Show("help"); err == nil {
				t.app.SetRoot(t.router.Pages(), true)
			}
			return nil
		},
	})

	t.commandRegistry.Register(&commands.Command{
		Name:        "workspaces",
		Description: "Toggle workspace panel",
		Keybinding:  "W or F3",
		Category:    commands.CategoryNav,
		Handler: func() error {
			t.toggleWorkspacePanel()
			return nil
		},
	})

	t.commandRegistry.Register(&commands.Command{
		Name:        "search",
		Description: "Search tickets in current workspace",
		Keybinding:  "/",
		Category:    commands.CategoryNav,
		Handler: func() error {
			t.showSearch()
			return nil
		},
	})

	t.commandRegistry.Register(&commands.Command{
		Name:        "command-palette",
		Description: "Show command palette",
		Keybinding:  ": or Ctrl+P or F1",
		Category:    commands.CategoryNav,
		Handler: func() error {
			t.showEnhancedCommandPalette()
			return nil
		},
	})

	// Sync commands
	t.commandRegistry.Register(&commands.Command{
		Name:        "pull",
		Description: "Pull latest tickets from Jira",
		Keybinding:  "P or F2",
		Category:    commands.CategorySync,
		Handler: func() error {
			t.handlePull()
			return nil
		},
	})

	t.commandRegistry.Register(&commands.Command{
		Name:        "push",
		Description: "Push tickets to Jira",
		Keybinding:  "p",
		Category:    commands.CategorySync,
		Handler: func() error {
			t.handlePush()
			return nil
		},
	})

	t.commandRegistry.Register(&commands.Command{
		Name:        "sync",
		Description: "Full sync (pull then push)",
		Keybinding:  "s",
		Category:    commands.CategorySync,
		Handler: func() error {
			t.handleSync()
			return nil
		},
	})

	// View commands
	t.commandRegistry.Register(&commands.Command{
		Name:        "refresh",
		Description: "Refresh current workspace tickets",
		Keybinding:  "r or F5",
		Category:    commands.CategoryView,
		Handler: func() error {
			t.handleRefresh()
			return nil
		},
	})

	// Edit commands
	t.commandRegistry.Register(&commands.Command{
		Name:        "bulk-operations",
		Description: "Perform bulk operations on selected tickets",
		Keybinding:  "b",
		Category:    commands.CategoryEdit,
		Handler: func() error {
			t.handleBulkOperations()
			return nil
		},
	})

	// System commands
	t.commandRegistry.Register(&commands.Command{
		Name:        "quit",
		Description: "Quit application",
		Keybinding:  "q or F10",
		Category:    commands.CategorySystem,
		Handler: func() error {
			t.app.Stop()
			return nil
		},
	})
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

// showEnhancedCommandPalette displays the enhanced command palette modal with registry (Phase 6, Day 8-9).
func (t *TUIApp) showEnhancedCommandPalette() {
	t.inModal = true
	t.actionBar.SetContext(widgets.ContextModal)
	t.commandPalette.Show()
	t.app.SetRoot(t.commandPalette.Primitive(), true)
	t.app.SetFocus(t.commandPalette.Primitive())
}

// onSyncStatusChanged is called when sync status changes (from sync coordinator).
func (t *TUIApp) onSyncStatusChanged(status sync.SyncStatus) {
	// Update UI from goroutine using QueueUpdateDraw
	t.app.QueueUpdateDraw(func() {
		// Update status view
		t.syncStatusView.SetStatus(status)

		// If sync completed successfully, reload tickets asynchronously
		if status.State == sync.StateSuccess {
			if ws, err := t.workspaceService.Current(); err == nil && ws != nil {
				t.ticketTreeView.LoadTicketsAsync(ws.ID)
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

	// CRITICAL (Phase 6.5 Emergency Fix Extension): Acquire read lock to safely access syncCoordinator
	// Prevents race condition where recreateJiraAdapter() replaces syncCoordinator while we're using it
	t.serviceMutex.RLock()
	syncCoordinator := t.syncCoordinator
	t.serviceMutex.RUnlock()

	// Nil check after acquiring lock - service might be replaced during workspace change
	if syncCoordinator == nil {
		t.syncStatusView.SetStatus(sync.NewErrorStatus("push", fmt.Errorf("sync coordinator not initialized")))
		return
	}

	// Start async push
	syncCoordinator.PushAsync(filePath, services.ProcessOptions{})
}

// handlePull initiates an async pull operation using the job queue.
func (t *TUIApp) handlePull() {
	// Get current workspace
	ws, err := t.workspaceService.Current()
	if err != nil || ws == nil {
		t.syncStatusView.SetStatus(sync.NewErrorStatus("pull", fmt.Errorf("no active workspace")))
		return
	}

	// Don't allow multiple concurrent operations
	if t.currentJobID != "" {
		t.syncStatusView.SetStatus(sync.NewErrorStatus("pull", fmt.Errorf("operation already in progress (press ESC to cancel)")))
		return
	}

	// For now, use tickets.md in the current working directory
	// TODO: In future, integrate with workspace file path configuration
	filePath := "tickets.md"

	// CRITICAL (Phase 6.5 Emergency Fix): Acquire read lock to safely access pullService
	// Prevents race condition where recreateJiraAdapter() replaces pullService while we're using it
	t.serviceMutex.RLock()
	pullService := t.pullService
	t.serviceMutex.RUnlock()

	// Nil check after acquiring lock - service might be replaced during workspace change
	if pullService == nil {
		t.syncStatusView.SetStatus(sync.NewErrorStatus("pull", fmt.Errorf("pull service not initialized")))
		return
	}

	// Create pull job with workspace ID for database persistence
	pullJob := jobs.NewPullJob(pullService, filePath, services.PullOptions{
		ProjectKey:  ws.ProjectKey,
		WorkspaceID: ws.ID, // Pass workspace ID for database operations
	})

	// Submit to job queue
	jobID := t.jobQueue.Submit(pullJob)
	t.currentJobID = jobID
	t.currentJobType = "pull"

	// Initial status update
	t.syncStatusView.SetStatus(sync.NewSyncingStatus("pull", "Starting pull operation..."))

	// Monitor job completion in background
	go t.monitorJobCompletion(jobID, pullJob)
}

// monitorJobCompletion watches a job and updates status when complete.
func (t *TUIApp) monitorJobCompletion(jobID jobs.JobID, pullJob *jobs.PullJob) {
	// Poll job status until complete
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		status, exists := t.jobQueue.Status(jobID)
		if !exists {
			return
		}

		if status == jobs.JobCompleted || status == jobs.JobFailed || status == jobs.JobCancelled {
			// Job finished
			t.app.QueueUpdateDraw(func() {
				t.currentJobID = ""
				t.currentJobType = ""

				// Clear progress bar (Phase 6, Day 10-11)
				t.syncStatusView.ClearProgress()

				switch status {
				case jobs.JobCompleted:
					msg := pullJob.FormatResult()
					t.syncStatusView.SetStatus(sync.NewSuccessStatus("pull", msg))
					// Reload tickets
					if ws, err := t.workspaceService.Current(); err == nil && ws != nil {
						t.ticketTreeView.LoadTicketsAsync(ws.ID)
					}
				case jobs.JobFailed:
					errMsg := "Unknown error"
					if err := pullJob.Error(); err != nil {
						errMsg = err.Error()
					}
					t.syncStatusView.SetStatus(sync.NewErrorStatus("pull", fmt.Errorf(errMsg)))
				case jobs.JobCancelled:
					msg := pullJob.FormatResult()
					if msg == "No result available" {
						msg = "Cancelled"
					} else {
						msg = "Cancelled: " + msg
					}
					t.syncStatusView.SetStatus(sync.NewErrorStatus("pull", fmt.Errorf(msg)))
					// Reload tickets to show partial results
					if ws, err := t.workspaceService.Current(); err == nil && ws != nil {
						t.ticketTreeView.LoadTicketsAsync(ws.ID)
					}
				}
			})
			return
		}
	}
}

// handleJobCancellation cancels the currently active job.
func (t *TUIApp) handleJobCancellation() {
	if t.currentJobID == "" {
		return
	}

	err := t.jobQueue.Cancel(t.currentJobID)
	if err != nil {
		t.syncStatusView.SetStatus(sync.NewErrorStatus(t.currentJobType, fmt.Errorf("cancel failed: %v", err)))
		return
	}

	t.syncStatusView.SetStatus(sync.NewSyncingStatus(t.currentJobType, "Cancelling..."))
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

	// CRITICAL (Phase 6.5 Emergency Fix Extension): Acquire read lock to safely access syncCoordinator
	// Prevents race condition where recreateJiraAdapter() replaces syncCoordinator while we're using it
	t.serviceMutex.RLock()
	syncCoordinator := t.syncCoordinator
	t.serviceMutex.RUnlock()

	// Nil check after acquiring lock - service might be replaced during workspace change
	if syncCoordinator == nil {
		t.syncStatusView.SetStatus(sync.NewErrorStatus("sync", fmt.Errorf("sync coordinator not initialized")))
		return
	}

	// Start async sync
	syncCoordinator.SyncAsync(filePath)
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

// create2PanelLayout creates the 2-panel layout (Phase 6.6: workspace moved to slide-out).
// Layout: Tickets (40%) | Detail (60%)
func (t *TUIApp) create2PanelLayout(rightPanel *tview.Flex) *tview.Flex {
	return tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(t.ticketTreeView.Primitive(), 0, 2, true). // 40% (2 of 5 parts), focused by default
		AddItem(rightPanel, 0, 3, false)                   // 60% (3 of 5 parts)
}

// handleRefresh reloads tickets for the current workspace asynchronously.
func (t *TUIApp) handleRefresh() {
	// Get current workspace
	ws, err := t.workspaceService.Current()
	if err != nil || ws == nil {
		t.syncStatusView.SetStatus(sync.NewErrorStatus("refresh", fmt.Errorf("no active workspace")))
		return
	}

	// Update status to show refreshing
	t.syncStatusView.SetStatus(sync.NewSyncingStatus("refresh", "Reloading tickets..."))

	// Reload tickets asynchronously
	// The LoadTicketsAsync method will handle the async loading and UI updates
	// We use a goroutine to also update the status when complete
	go func() {
		tickets, err := t.ticketQuery.ListByWorkspace(ws.ID)

		// Update status based on result
		t.app.QueueUpdateDraw(func() {
			if err != nil {
				t.syncStatusView.SetStatus(sync.NewErrorStatus("refresh", err))
			} else {
				msg := fmt.Sprintf("%d ticket(s) loaded", len(tickets))
				t.syncStatusView.SetStatus(sync.NewSuccessStatus("refresh", msg))
			}
		})
	}()

	// Trigger the tree view to reload (it will show its own loading indicator)
	t.ticketTreeView.LoadTicketsAsync(ws.ID)
}

// showWorkspaceModal displays the workspace creation modal.
func (t *TUIApp) showWorkspaceModal() {
	// BUGFIX (Phase 6.5): Auto-close workspace panel if open to prevent state corruption
	// When modal opens while panel is visible, ESC handler would leave panel in inconsistent state
	// (SlideOut.isVisible=true but page hidden), causing W toggle to work in reverse.
	if t.workspaceSlideOut != nil && t.workspaceSlideOut.IsVisible() {
		t.toggleWorkspacePanel() // Properly sync both isVisible and pages state
	}
	t.inModal = true
	t.workspaceModal.OnShow()
	t.workspaceModal.Show()
}

// showProfileManagement displays the profile management interface.
// For now, this shows the workspace modal in profile creation mode.
// In the future, this could be expanded to a dedicated profile management view.
func (t *TUIApp) showProfileManagement() {
	// For Slice 3, we'll use the workspace modal's profile creation functionality
	// Future enhancement could create a dedicated profile management modal
	t.showWorkspaceModal()
}

// handleBulkOperations opens the bulk operations modal for selected tickets.
func (t *TUIApp) handleBulkOperations() {
	// Get selected tickets from tree view
	selectedTickets := t.ticketTreeView.GetSelectedTickets()

	// Show modal with selected tickets
	t.inModal = true
	t.bulkOperationsModal.Show(selectedTickets)
}

// toggleWorkspacePanel shows or hides the workspace slide-out panel (Phase 6.6).
func (t *TUIApp) toggleWorkspacePanel() {
	if t.workspaceSlideOut == nil {
		return
	}

	if t.workspaceSlideOut.IsVisible() {
		// Hide the panel
		t.workspaceSlideOut.Hide()
		t.pages.HidePage("workspace-overlay")
		// Return focus to main view
		t.updateFocus()
	} else {
		// Show the panel
		t.workspaceListView.OnShow() // Refresh workspace list
		t.workspaceSlideOut.Show()
		t.pages.ShowPage("workspace-overlay")
		// Focus the workspace list
		t.workspaceListView.SetFocused(true)
		t.app.SetFocus(t.workspaceListView.Primitive())
	}
}

// recreateJiraAdapter creates a new Jira adapter for the current workspace
// and updates all sync services (Phase 6.5 Critical Fix).
// This fixes the bug where adapter was created once at startup and never updated
// when switching workspaces, causing credentials to be stale.
func (t *TUIApp) recreateJiraAdapter() error {
	// Get current workspace
	ws, err := t.workspaceService.Current()
	if err != nil {
		return fmt.Errorf("no active workspace: %w", err)
	}

	// Get workspace config
	config, err := t.workspaceService.GetConfig(ws.Name)
	if err != nil {
		return fmt.Errorf("failed to get workspace config: %w", err)
	}

	// Create new Jira adapter with current workspace credentials
	jiraAdapter, err := jira.NewJiraAdapterFromConfig(config, nil)
	if err != nil {
		return fmt.Errorf("failed to create Jira adapter: %w", err)
	}

	// Get file repository and state manager (reuse these dependencies)
	fileRepo := filesystem.NewFileRepository()
	stateManager := state.NewStateManager("")

	// Get database adapter from path resolver
	dbAdapter, err := database.NewSQLiteAdapter(t.pathResolver)
	if err != nil {
		return fmt.Errorf("failed to get database adapter: %w", err)
	}

	// CRITICAL (Phase 6.5 Emergency Fix): Acquire write lock before replacing services
	// This prevents race conditions where handlePull/handlePush access services while we replace them
	t.serviceMutex.Lock()
	defer t.serviceMutex.Unlock()

	// Recreate services with new adapter
	t.pullService = services.NewPullServiceWithDB(jiraAdapter, dbAdapter, fileRepo, stateManager)
	t.pushService = services.NewPushService(fileRepo, jiraAdapter, stateManager)

	// Update sync coordinator with new services
	t.syncCoordinator = sync.NewSyncCoordinator(
		t.pushService,
		t.pullService,
		t.onSyncStatusChanged,
	)

	// Update bulk operation service with new adapter (if used)
	t.bulkOperationService = services.NewBulkOperationService(jiraAdapter)

	return nil
}
