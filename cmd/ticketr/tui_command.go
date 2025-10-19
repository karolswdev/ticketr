package main

import (
	"fmt"

	"github.com/karolswdev/ticktr/internal/adapters/tui"
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
  Tab        - Switch panels (future)
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

	// Initialize PathResolver
	pathResolver, err := services.NewPathResolver()
	if err != nil {
		return fmt.Errorf("failed to initialize path resolver: %w", err)
	}

	// Create and run TUI application
	app, err := tui.NewTUIApp(workspaceService, pathResolver)
	if err != nil {
		return fmt.Errorf("failed to create TUI application: %w", err)
	}

	if err := app.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
