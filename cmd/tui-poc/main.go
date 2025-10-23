package main

import (
	"fmt"
	"os"

	"github.com/karolswdev/ticktr/internal/adapters/database"
	"github.com/karolswdev/ticktr/internal/adapters/keychain"
	"github.com/karolswdev/ticktr/internal/core/services"
	tuibubbletea "github.com/karolswdev/ticktr/internal/tui-bubbletea"
)

// tui-poc is a proof-of-concept demonstrating the new Bubbletea TUI architecture.
//
// NOTE: This is the old Week 1 POC. For the latest version with real data integration,
// use cmd/ticketr-tui-poc instead.
//
// Week 2 Day 1 Update: Now requires service initialization for data integration.
//
// Run: go run cmd/tui-poc/main.go
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

	// Initialize credential store
	credStore := keychain.NewKeychainStore()

	// Initialize repositories
	workspaceRepo := database.NewWorkspaceRepository(dbAdapter.DB())
	credProfileRepo := database.NewCredentialProfileRepository(dbAdapter.DB())

	// Initialize services
	workspaceService := services.NewWorkspaceService(workspaceRepo, credProfileRepo, credStore)
	ticketQuery := services.NewTicketQueryService(dbAdapter)

	// Create the Bubbletea TUI app with services
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
}
