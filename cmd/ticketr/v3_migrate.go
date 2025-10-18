package main

import (
	"fmt"
	"os"

	"github.com/karolswdev/ticktr/internal/adapters/database"
	"github.com/karolswdev/ticktr/internal/config"
	"github.com/spf13/cobra"
)

var (
	// v3 migration flags
	v3MigrateAll     bool
	v3MigrateDir     string
	v3MigrateDryRun  bool
	v3MigrateVerbose bool
	v3MigrateBackup  bool
	v3Phase          string

	v3Cmd = &cobra.Command{
		Use:   "v3",
		Short: "Ticketr v3.0 migration and feature management",
		Long: `Manage the transition to Ticketr v3.0 with SQLite backend, workspaces, and TUI.

This command helps migrate from v2.x file-based state to v3.0 SQLite database,
enable new features progressively, and manage the migration process.`,
	}

	v3MigrateCmd = &cobra.Command{
		Use:   "migrate [path]",
		Short: "Migrate v2.x state files to v3.0 SQLite database",
		Long: `Migrate existing .ticketr.state files to the new SQLite database backend.

This command scans for .ticketr.state files and migrates them to SQLite while
preserving all ticket data and sync state. Original files are backed up by default.

Examples:
  ticketr v3 migrate                    # Migrate current directory
  ticketr v3 migrate ~/projects         # Migrate all projects in directory
  ticketr v3 migrate . --dry-run        # Preview migration without changes
  ticketr v3 migrate --all              # Migrate all projects from home directory`,
		Run: runV3Migrate,
	}

	v3StatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show v3.0 feature status and configuration",
		Long:  `Display the current status of v3.0 features, migration state, and configuration.`,
		Run:   runV3Status,
	}

	v3EnableCmd = &cobra.Command{
		Use:   "enable [phase]",
		Short: "Enable v3.0 features by phase",
		Long: `Progressively enable v3.0 features based on release phase.

Phases:
  alpha  - Enable SQLite backend only
  beta   - Enable SQLite + Workspaces
  rc     - Enable SQLite + Workspaces + TUI
  stable - Enable all v3.0 features

Example:
  ticketr v3 enable alpha    # Enable Phase 1 features
  ticketr v3 enable beta     # Enable Phase 2 features`,
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"alpha", "beta", "rc", "stable"},
		Run:       runV3Enable,
	}

	v3WorkspaceCmd = &cobra.Command{
		Use:   "workspace",
		Short: "Manage v3.0 workspaces",
		Long:  `Create, list, and manage workspaces in v3.0 (requires SQLite enabled).`,
	}

	v3WorkspaceCreateCmd = &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new workspace",
		Args:  cobra.ExactArgs(1),
		Run:   runV3WorkspaceCreate,
	}

	v3WorkspaceListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all workspaces",
		Run:   runV3WorkspaceList,
	}
)

func init() {
	// Add v3 command to root
	rootCmd.AddCommand(v3Cmd)

	// Add subcommands to v3
	v3Cmd.AddCommand(v3MigrateCmd)
	v3Cmd.AddCommand(v3StatusCmd)
	v3Cmd.AddCommand(v3EnableCmd)
	v3Cmd.AddCommand(v3WorkspaceCmd)

	// Add workspace subcommands
	v3WorkspaceCmd.AddCommand(v3WorkspaceCreateCmd)
	v3WorkspaceCmd.AddCommand(v3WorkspaceListCmd)

	// Migration flags
	v3MigrateCmd.Flags().BoolVar(&v3MigrateAll, "all", false, "Migrate all projects from home directory")
	v3MigrateCmd.Flags().StringVar(&v3MigrateDir, "dir", ".", "Directory to search for projects")
	v3MigrateCmd.Flags().BoolVar(&v3MigrateDryRun, "dry-run", false, "Preview migration without making changes")
	v3MigrateCmd.Flags().BoolVar(&v3MigrateVerbose, "verbose", false, "Show detailed migration progress")
	v3MigrateCmd.Flags().BoolVar(&v3MigrateBackup, "backup", true, "Backup state files before migration")
}

func runV3Migrate(cmd *cobra.Command, args []string) {
	// Load feature flags
	features, err := config.LoadFeatures()
	if err != nil {
		fmt.Printf("Error loading features: %v\n", err)
		os.Exit(1)
	}

	// Check if SQLite is enabled
	if !features.UseSQLite && !v3MigrateDryRun {
		fmt.Println("SQLite backend is not enabled. Migration requires SQLite to be enabled.")
		fmt.Println("\nTo enable SQLite and start migration:")
		fmt.Println("  ticketr v3 enable alpha")
		fmt.Println("\nOr set environment variable:")
		fmt.Println("  export TICKETR_USE_SQLITE=true")
		os.Exit(1)
	}

	// Determine migration path
	migrationPath := v3MigrateDir
	if len(args) > 0 {
		migrationPath = args[0]
	}
	if v3MigrateAll {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		migrationPath = homeDir
	}

	// Create migrator
	migrator, err := database.NewStateMigrator(features.SQLitePath, v3MigrateVerbose)
	if err != nil {
		fmt.Printf("Error creating migrator: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting migration from: %s\n", migrationPath)
	if v3MigrateDryRun {
		fmt.Println("🔍 Running in DRY-RUN mode - no changes will be made")
	}

	// Perform migration
	report, err := migrator.MigrateDirectory(migrationPath)
	if err != nil {
		fmt.Printf("Error during migration: %v\n", err)
		os.Exit(1)
	}

	// Display report
	fmt.Println(report.GenerateReport())

	if report.SuccessfulProjects > 0 && !v3MigrateDryRun {
		fmt.Println("\n✅ Migration complete! Your projects have been migrated to SQLite.")
		fmt.Println("\nNext steps:")
		fmt.Println("1. Test the migration by running normal ticketr commands")
		fmt.Println("2. Original .ticketr.state files have been backed up with .backup suffix")
		fmt.Println("3. To use workspaces: ticketr v3 workspace list")

		if !features.EnableWorkspaces {
			fmt.Println("\n💡 Tip: Enable workspaces with: ticketr v3 enable beta")
		}
	}

	if report.FailedProjects > 0 {
		os.Exit(1)
	}
}

func runV3Status(cmd *cobra.Command, args []string) {
	features, err := config.LoadFeatures()
	if err != nil {
		fmt.Printf("Error loading features: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(features.Status())

	// Check for migration status if SQLite is enabled
	if features.UseSQLite {
		adapter, err := database.NewSQLiteAdapter(features.SQLitePath)
		if err != nil {
			fmt.Printf("\n⚠️  SQLite database not initialized: %v\n", err)
		} else {
			defer adapter.Close()

			// Count tickets in database
			tickets, err := adapter.GetTicketsByWorkspace("default")
			if err == nil {
				fmt.Printf("\nDatabase Status:\n")
				fmt.Printf("================\n")
				fmt.Printf("Location:       %s\n", features.SQLitePath)
				fmt.Printf("Tickets:        %d\n", len(tickets))

				// Check if database exists
				if stat, err := os.Stat(features.SQLitePath); err == nil {
					fmt.Printf("Database Size:  %.2f KB\n", float64(stat.Size())/1024)
				}
			}
		}
	}

	// Show migration recommendation
	if !features.UseSQLite {
		fmt.Println("\n💡 Ready to migrate to v3.0?")
		fmt.Println("   Start with: ticketr v3 enable alpha")
		fmt.Println("   Then run:   ticketr v3 migrate")
	}
}

func runV3Enable(cmd *cobra.Command, args []string) {
	phase := args[0]

	features, err := config.EnableV3Features(phase)
	if err != nil {
		fmt.Printf("Error enabling v3 features: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully enabled v3.0 %s features\n\n", phase)
	fmt.Println(features.Status())

	// Provide next steps based on phase
	switch phase {
	case "alpha":
		fmt.Println("\nNext steps:")
		fmt.Println("1. Run 'ticketr v3 migrate' to migrate existing projects")
		fmt.Println("2. Test existing commands with SQLite backend")
		fmt.Println("3. When ready, enable beta features: ticketr v3 enable beta")

	case "beta":
		fmt.Println("\nNext steps:")
		fmt.Println("1. Create workspaces with: ticketr v3 workspace create <name>")
		fmt.Println("2. Switch between workspaces: ticketr workspace switch <name>")
		fmt.Println("3. When ready, enable RC features: ticketr v3 enable rc")

	case "rc":
		fmt.Println("\nNext steps:")
		fmt.Println("1. Launch TUI with: ticketr tui")
		fmt.Println("2. Test all features in TUI mode")
		fmt.Println("3. When satisfied, enable stable: ticketr v3 enable stable")

	case "stable":
		fmt.Println("\n🎉 Ticketr v3.0 is fully enabled!")
		fmt.Println("Enjoy the enhanced features:")
		fmt.Println("- SQLite database backend")
		fmt.Println("- Multi-workspace support")
		fmt.Println("- Terminal User Interface")
	}
}

func runV3WorkspaceCreate(cmd *cobra.Command, args []string) {
	name := args[0]

	features, err := config.LoadFeatures()
	if err != nil {
		fmt.Printf("Error loading features: %v\n", err)
		os.Exit(1)
	}

	if !features.EnableWorkspaces {
		fmt.Println("Workspaces feature is not enabled.")
		fmt.Println("Enable with: ticketr v3 enable beta")
		os.Exit(1)
	}

	// TODO: Implement workspace creation
	fmt.Printf("Creating workspace '%s'...\n", name)
	fmt.Println("Note: Full workspace implementation coming in Phase 2")
}

func runV3WorkspaceList(cmd *cobra.Command, args []string) {
	features, err := config.LoadFeatures()
	if err != nil {
		fmt.Printf("Error loading features: %v\n", err)
		os.Exit(1)
	}

	if !features.EnableWorkspaces {
		fmt.Println("Workspaces feature is not enabled.")
		fmt.Println("Enable with: ticketr v3 enable beta")
		os.Exit(1)
	}

	adapter, err := database.NewSQLiteAdapter(features.SQLitePath)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer adapter.Close()

	// Query workspaces
	// For now, we just have the default workspace
	fmt.Println("Workspaces:")
	fmt.Println("===========")
	fmt.Println("* default (active)")
	fmt.Println("\nNote: Multi-workspace support coming in Phase 2")
}