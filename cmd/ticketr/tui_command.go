package main

import (
	"fmt"

	"github.com/karolswdev/ticktr/internal/adapters/database"
	"github.com/karolswdev/ticktr/internal/adapters/filesystem"
	"github.com/karolswdev/ticktr/internal/adapters/tui"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/karolswdev/ticktr/internal/state"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the Terminal User Interface",
	Long: `Launch the interactive Terminal User Interface for managing workspaces and tickets.

The TUI provides a visual interface for:
- Switching between workspaces
- Browsing ticket hierarchies
- Viewing and editing ticket details
- Syncing with Jira (push, pull, sync operations)

Keyboard shortcuts:
  q, Ctrl+C  - Quit
  ?          - Show help
  Tab        - Switch between panels
  j/k        - Navigate up/down (vim-style)
  h/l        - Collapse/expand tree nodes
  p          - Push tickets to Jira
  P          - Pull tickets from Jira
  r          - Refresh tickets
  s          - Full sync (pull then push)
`,
	RunE: runTUI,
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

func runTUI(cmd *cobra.Command, args []string) error {
	// Initialize workspace service
	workspaceService, err := initWorkspaceService()
	if err != nil {
		return fmt.Errorf("failed to initialize workspace service: %w", err)
	}

	// Initialize ticket query service
	ticketQueryService, err := initTicketQueryService()
	if err != nil {
		return fmt.Errorf("failed to initialize ticket query service: %w", err)
	}

	// Initialize PathResolver
	pathResolver, err := services.NewPathResolver()
	if err != nil {
		return fmt.Errorf("failed to initialize path resolver: %w", err)
	}

	// Initialize Jira adapter for sync operations (workspace credentials or environment variables)
	jiraAdapter, err := initJiraAdapter(nil)
	if err != nil {
		return fmt.Errorf("failed to initialize Jira adapter: %w\n\nPlease configure credentials using one of these methods:\n\nOption 1: Use workspace credentials (recommended):\n  ticketr workspace create <name> --url <jira-url> --project <key> --username <email> --token <api-token>\n\nOption 2: Set environment variables (legacy):\n  - JIRA_URL\n  - JIRA_EMAIL\n  - JIRA_API_KEY\n  - JIRA_PROJECT_KEY", err)
	}

	// Initialize file repository for sync operations
	fileRepo := filesystem.NewFileRepository()

	// Initialize state manager for sync operations (uses PathResolver for global paths)
	stateManager := state.NewStateManager("")

	// Initialize push service
	pushService := services.NewPushService(fileRepo, jiraAdapter, stateManager)

	// Initialize pull service
	pullService := services.NewPullService(jiraAdapter, fileRepo, stateManager)

	// Initialize bulk operation service (Week 18)
	bulkOperationService := services.NewBulkOperationService(jiraAdapter)

	// Create and run TUI application
	app, err := tui.NewTUIApp(
		workspaceService,
		ticketQueryService,
		pathResolver,
		pushService,
		pullService,
		bulkOperationService,
	)
	if err != nil {
		return fmt.Errorf("failed to create TUI application: %w", err)
	}

	if err := app.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}

// initTicketQueryService initializes the ticket query service with SQLite adapter.
func initTicketQueryService() (*services.TicketQueryService, error) {
	// Get PathResolver singleton
	pathResolver, err := services.GetPathResolver()
	if err != nil {
		return nil, fmt.Errorf("failed to get path resolver: %w", err)
	}

	// Initialize SQLite adapter with PathResolver
	adapter, err := database.NewSQLiteAdapter(pathResolver)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create ticket query service
	return services.NewTicketQueryService(adapter), nil
}
