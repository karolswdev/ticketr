// Package main provides the entry point for the Ticketr Bubbletea TUI proof-of-concept.
//
// Week 2 Day 1 implementation featuring:
// - Real data integration from database
// - Workspace service integration
// - Ticket query service integration
// - Async data loading with loading states
// - Error handling for data failures
// - Real workspace and ticket display
//
// Previous Week 1 features:
// - Three beautiful themes (Default, Dark, Arctic)
// - Responsive FlexBox layout system
// - Enhanced header with live sync status
// - Panel-based UI with focus management
// - Keyboard shortcuts (Tab, 1/2/3, q)
//
// Usage:
//
//	# Normal mode with real data
//	go run cmd/ticketr-tui-poc/main.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/karolswdev/ticktr/internal/adapters/database"
	"github.com/karolswdev/ticktr/internal/adapters/keychain"
	"github.com/karolswdev/ticktr/internal/core/services"
	tuibubbletea "github.com/karolswdev/ticktr/internal/tui-bubbletea"
)

func main() {
	// Initialize PathResolver for database access
	pathResolver, err := services.GetPathResolver()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing path resolver: %v\n", err)
		os.Exit(1)
	}

	// Ensure directories exist
	if err := pathResolver.EnsureDirectories(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directories: %v\n", err)
		os.Exit(1)
	}

	// Initialize database adapter
	dbAdapter, err := database.NewSQLiteAdapter(pathResolver)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer dbAdapter.Close()

	// Initialize credential store (needed for workspace service)
	credStore := keychain.NewKeychainStore()

	// Initialize workspace repository
	workspaceRepo := database.NewWorkspaceRepository(dbAdapter.DB())
	credProfileRepo := database.NewCredentialProfileRepository(dbAdapter.DB())

	// Initialize services
	workspaceService := services.NewWorkspaceService(workspaceRepo, credProfileRepo, credStore)
	ticketQuery := services.NewTicketQueryService(dbAdapter)

	// Create the Bubbletea app with services
	app, err := tuibubbletea.NewApp(workspaceService, ticketQuery)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating TUI app: %v\n", err)
		os.Exit(1)
	}

	// Run the app
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}

	log.Println("Ticketr Bubbletea TUI POC exited successfully")
}
