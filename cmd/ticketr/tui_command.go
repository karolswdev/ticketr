package main

import (
	"fmt"

	"github.com/karolswdev/ticktr/internal/adapters/database"
	"github.com/karolswdev/ticktr/internal/adapters/tui"
	"github.com/karolswdev/ticktr/internal/config"
	"github.com/karolswdev/ticktr/internal/core/services"
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
- Syncing with Jira

Keyboard shortcuts:
  q, Ctrl+C  - Quit
  ?          - Show help
  Tab        - Switch between panels
  j/k        - Navigate up/down (vim-style)
  h/l        - Collapse/expand tree nodes
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

	// Create and run TUI application
	app, err := tui.NewTUIApp(workspaceService, ticketQueryService, pathResolver)
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
	// Load feature flags
	features, err := config.LoadFeatures()
	if err != nil {
		return nil, fmt.Errorf("failed to load features: %w", err)
	}

	// Check if workspaces are enabled
	if !features.EnableWorkspaces {
		return nil, fmt.Errorf("workspaces feature is not enabled. Enable with: ticketr v3 enable beta")
	}

	// Initialize SQLite adapter (implements ExtendedRepository)
	adapter, err := database.NewSQLiteAdapter(features.SQLitePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create ticket query service
	return services.NewTicketQueryService(adapter), nil
}
