package migration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/services"
)

func TestPathResolverMigrator_Migrate(t *testing.T) {
	// Create temp directory with legacy structure
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	legacyDir := filepath.Join(tempDir, ".ticketr")
	if err := os.MkdirAll(legacyDir, 0755); err != nil {
		t.Fatalf("Failed to create legacy directory: %v", err)
	}
	legacyDB := filepath.Join(legacyDir, "ticketr.db")
	if err := os.WriteFile(legacyDB, []byte("test database content"), 0644); err != nil {
		t.Fatalf("Failed to create legacy database: %v", err)
	}

	// Create PathResolver with temp directory
	pr, err := services.NewPathResolverWithOptions("ticketr",
		func(key string) string { return "" },
		func() (string, error) { return filepath.Join(tempDir, "home"), nil })
	if err != nil {
		t.Fatalf("Failed to create PathResolver: %v", err)
	}

	migrator := NewPathResolverMigrator(pr)

	// Run migration
	err = migrator.Migrate()
	if err != nil {
		t.Fatalf("Migrate() failed: %v", err)
	}

	// Verify new database exists
	newDB := pr.DatabasePath()
	if _, err := os.Stat(newDB); os.IsNotExist(err) {
		t.Errorf("New database was not created at: %s", newDB)
	}

	// Verify backup was created
	backupDir := filepath.Join(pr.DataDir(), "backups")
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		t.Errorf("Backup directory was not created at: %s", backupDir)
	}

	// Verify README was created in legacy directory
	readmePath := filepath.Join(legacyDir, "MIGRATED-README.txt")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		t.Errorf("Migration README was not created at: %s", readmePath)
	}

	// Verify content was copied correctly
	content, err := os.ReadFile(newDB)
	if err != nil {
		t.Fatalf("Failed to read new database: %v", err)
	}
	if string(content) != "test database content" {
		t.Errorf("Database content mismatch: got %q, want %q", string(content), "test database content")
	}
}

func TestPathResolverMigrator_Migrate_Idempotent(t *testing.T) {
	// Create temp directory with legacy structure
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	legacyDir := filepath.Join(tempDir, ".ticketr")
	if err := os.MkdirAll(legacyDir, 0755); err != nil {
		t.Fatalf("Failed to create legacy directory: %v", err)
	}
	legacyDB := filepath.Join(legacyDir, "ticketr.db")
	if err := os.WriteFile(legacyDB, []byte("original content"), 0644); err != nil {
		t.Fatalf("Failed to create legacy database: %v", err)
	}

	// Create PathResolver with temp directory
	pr, err := services.NewPathResolverWithOptions("ticketr",
		func(key string) string { return "" },
		func() (string, error) { return filepath.Join(tempDir, "home"), nil })
	if err != nil {
		t.Fatalf("Failed to create PathResolver: %v", err)
	}

	migrator := NewPathResolverMigrator(pr)

	// Run migration first time
	if err := migrator.Migrate(); err != nil {
		t.Fatalf("First Migrate() failed: %v", err)
	}

	// Modify new database to ensure second migration doesn't overwrite
	newDB := pr.DatabasePath()
	if err := os.WriteFile(newDB, []byte("modified content"), 0644); err != nil {
		t.Fatalf("Failed to modify new database: %v", err)
	}

	// Run migration second time
	if err := migrator.Migrate(); err != nil {
		t.Fatalf("Second Migrate() failed: %v", err)
	}

	// Verify database was not overwritten
	content, err := os.ReadFile(newDB)
	if err != nil {
		t.Fatalf("Failed to read new database: %v", err)
	}
	if string(content) != "modified content" {
		t.Errorf("Second migration overwrote database: got %q, want %q", string(content), "modified content")
	}
}

func TestPathResolverMigrator_Migrate_NoLegacy(t *testing.T) {
	// Create temp directory WITHOUT legacy structure
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	// Create PathResolver with temp directory
	pr, err := services.NewPathResolverWithOptions("ticketr",
		func(key string) string { return "" },
		func() (string, error) { return filepath.Join(tempDir, "home"), nil })
	if err != nil {
		t.Fatalf("Failed to create PathResolver: %v", err)
	}

	migrator := NewPathResolverMigrator(pr)

	// Run migration (should be no-op)
	err = migrator.Migrate()
	if err != nil {
		t.Errorf("Migrate() failed when no legacy database exists: %v", err)
	}

	// Verify no database was created (since there's nothing to migrate)
	newDB := pr.DatabasePath()
	if _, err := os.Stat(newDB); err == nil {
		t.Errorf("Database was created even though there was nothing to migrate")
	}
}

func TestPathResolverMigrator_Migrate_StateFile(t *testing.T) {
	// Create temp directory with legacy structure AND state file
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	legacyDir := filepath.Join(tempDir, ".ticketr")
	if err := os.MkdirAll(legacyDir, 0755); err != nil {
		t.Fatalf("Failed to create legacy directory: %v", err)
	}
	legacyDB := filepath.Join(legacyDir, "ticketr.db")
	if err := os.WriteFile(legacyDB, []byte("db content"), 0644); err != nil {
		t.Fatalf("Failed to create legacy database: %v", err)
	}

	// Create legacy state file
	legacyState := ".ticketr.state"
	if err := os.WriteFile(legacyState, []byte("state content"), 0644); err != nil {
		t.Fatalf("Failed to create legacy state file: %v", err)
	}

	// Create PathResolver with temp directory
	pr, err := services.NewPathResolverWithOptions("ticketr",
		func(key string) string { return "" },
		func() (string, error) { return filepath.Join(tempDir, "home"), nil })
	if err != nil {
		t.Fatalf("Failed to create PathResolver: %v", err)
	}

	migrator := NewPathResolverMigrator(pr)

	// Run migration
	if err := migrator.Migrate(); err != nil {
		t.Fatalf("Migrate() failed: %v", err)
	}

	// Verify state file was migrated
	newState := filepath.Join(pr.DataDir(), "state.json")
	if _, err := os.Stat(newState); os.IsNotExist(err) {
		t.Errorf("State file was not migrated to: %s", newState)
	}

	// Verify state content
	content, err := os.ReadFile(newState)
	if err != nil {
		t.Fatalf("Failed to read new state file: %v", err)
	}
	if string(content) != "state content" {
		t.Errorf("State file content mismatch: got %q, want %q", string(content), "state content")
	}
}
