package main

import (
	"fmt"

	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/karolswdev/ticktr/internal/migration"
	"github.com/spf13/cobra"
)

var (
	// migratePathsCmd triggers manual migration to global paths
	migratePathsCmd = &cobra.Command{
		Use:   "migrate-paths",
		Short: "Migrate from v2.x local paths to v3.x global paths",
		Long: `Migrate data from legacy .ticketr/ directory to XDG-compliant global paths.

This command:
- Creates a backup of your existing database
- Copies all data to platform-standard global directories
- Leaves legacy directory intact with migration notice

Safe to run multiple times (idempotent).`,
		RunE: runMigratePaths,
	}

	// rollbackPathsCmd rollback migration to legacy paths
	rollbackPathsCmd = &cobra.Command{
		Use:   "rollback-paths",
		Short: "Rollback PathResolver migration to legacy local paths",
		Long: `Rollback migration by copying data from XDG global paths back to legacy .ticketr/ directory.

WARNING: This is intended for downgrading to v2.x. You will be prompted for confirmation.`,
		RunE: runRollbackPaths,
	}
)

func init() {
	rootCmd.AddCommand(migratePathsCmd)
	rootCmd.AddCommand(rollbackPathsCmd)
}

// runMigratePaths handles the migrate-paths command
func runMigratePaths(cmd *cobra.Command, args []string) error {
	pr, err := services.GetPathResolver()
	if err != nil {
		return fmt.Errorf("failed to get path resolver: %w", err)
	}

	migrator := migration.NewPathResolverMigrator(pr)
	return migrator.Migrate()
}

// runRollbackPaths handles the rollback-paths command
func runRollbackPaths(cmd *cobra.Command, args []string) error {
	pr, err := services.GetPathResolver()
	if err != nil {
		return fmt.Errorf("failed to get path resolver: %w", err)
	}

	migrator := migration.NewPathResolverMigrator(pr)
	return migrator.Rollback()
}
