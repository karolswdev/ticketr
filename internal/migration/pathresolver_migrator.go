package migration

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/karolswdev/ticktr/internal/core/services"
)

// PathResolverMigrator handles migration from legacy local paths to XDG-compliant global paths.
type PathResolverMigrator struct {
	pathResolver *services.PathResolver
}

// NewPathResolverMigrator creates a new PathResolverMigrator.
func NewPathResolverMigrator(pr *services.PathResolver) *PathResolverMigrator {
	return &PathResolverMigrator{pathResolver: pr}
}

// Migrate performs automatic migration from legacy paths to new XDG paths.
// Creates backup before migrating. Safe to run multiple times (idempotent).
func (m *PathResolverMigrator) Migrate() error {
	legacyDBPath := filepath.Join(".ticketr", "ticketr.db")

	// Check if migration needed
	if _, err := os.Stat(legacyDBPath); os.IsNotExist(err) {
		return nil // No legacy installation found
	}

	// Check if already migrated
	newDBPath := m.pathResolver.DatabasePath()
	if _, err := os.Stat(newDBPath); err == nil {
		fmt.Println("Migration already complete (new database exists)")
		return nil
	}

	// Create backup
	backupDir := filepath.Join(m.pathResolver.DataDir(), "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	backupPath := filepath.Join(backupDir, fmt.Sprintf("legacy-db-%d.db", time.Now().Unix()))
	if err := copyFile(legacyDBPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup database: %w", err)
	}
	fmt.Printf("✓ Backed up legacy database to: %s\n", backupPath)

	// Migrate database
	if err := m.pathResolver.EnsureDirectories(); err != nil {
		return fmt.Errorf("failed to create new directories: %w", err)
	}

	if err := copyFile(legacyDBPath, newDBPath); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	fmt.Printf("✓ Migrated database to: %s\n", newDBPath)

	// Migrate state file
	legacyStatePath := ".ticketr.state"
	if _, err := os.Stat(legacyStatePath); err == nil {
		newStatePath := filepath.Join(m.pathResolver.DataDir(), "state.json")
		if err := copyFile(legacyStatePath, newStatePath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to migrate state file: %v\n", err)
		} else {
			fmt.Printf("✓ Migrated state file to: %s\n", newStatePath)
		}
	}

	// Leave README in legacy directory
	readmePath := filepath.Join(".ticketr", "MIGRATED-README.txt")
	readme := fmt.Sprintf(`This directory has been migrated to: %s

Your data has been safely copied to the new location.
This directory is no longer used by Ticketr v3.0+.

You can safely delete this directory after verifying your data:
  rm -rf .ticketr

Backup location: %s
`, m.pathResolver.DataDir(), backupPath)

	os.WriteFile(readmePath, []byte(readme), 0644)

	// Create migration completion marker
	migrationStatePath := m.pathResolver.DataFile(".migration_complete")
	if err := os.WriteFile(migrationStatePath, []byte(time.Now().Format(time.RFC3339)), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to create migration state marker: %v\n", err)
	}

	fmt.Println("\n✅ Migration complete! Legacy directory preserved for safety.")
	fmt.Println("   Delete .ticketr/ after verifying your data works correctly.")

	return nil
}

// Rollback copies data from new paths back to legacy paths.
// Requires user confirmation.
func (m *PathResolverMigrator) Rollback() error {
	newDBPath := m.pathResolver.DatabasePath()
	legacyDBPath := filepath.Join(".ticketr", "ticketr.db")

	// Check if migration occurred
	if _, err := os.Stat(newDBPath); os.IsNotExist(err) {
		fmt.Println("No migration to rollback")
		return nil
	}

	// Confirm with user
	fmt.Println("WARNING: This will copy data from global paths back to legacy local paths.")
	fmt.Print("Continue? [y/N]: ")
	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "yes" {
		fmt.Println("Rollback cancelled")
		return nil
	}

	// Ensure legacy directory exists
	if err := os.MkdirAll(".ticketr", 0755); err != nil {
		return fmt.Errorf("failed to create legacy directory: %w", err)
	}

	// Copy database back
	if err := copyFile(newDBPath, legacyDBPath); err != nil {
		return fmt.Errorf("failed to rollback database: %w", err)
	}
	fmt.Printf("✓ Rolled back database to: %s\n", legacyDBPath)

	// Copy state back
	newStatePath := filepath.Join(m.pathResolver.DataDir(), "state.json")
	legacyStatePath := ".ticketr.state"
	if _, err := os.Stat(newStatePath); err == nil {
		if err := copyFile(newStatePath, legacyStatePath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to rollback state file: %v\n", err)
		} else {
			fmt.Printf("✓ Rolled back state file to: %s\n", legacyStatePath)
		}
	}

	// Remove migration completion marker
	migrationStatePath := m.pathResolver.DataFile(".migration_complete")
	if err := os.Remove(migrationStatePath); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Warning: Failed to remove migration state marker: %v\n", err)
	}

	fmt.Println("\n✅ Rollback complete! You can now use Ticketr v2.x")
	return nil
}

// copyFile copies a file from src to dst, creating parent directories as needed.
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}
