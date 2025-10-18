package database

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Migration represents a database migration
type Migration struct {
	Version int
	Name    string
	SQL     string
}

// Migrator handles database schema migrations
type Migrator struct {
	db         *sql.DB
	migrations []Migration
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: getEmbeddedMigrations(),
	}
}

// Migrate runs all pending migrations
func (m *Migrator) Migrate() error {
	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Run pending migrations
	for _, migration := range m.migrations {
		if _, exists := applied[migration.Version]; exists {
			continue // Skip already applied migrations
		}

		if err := m.runMigration(migration); err != nil {
			return fmt.Errorf("failed to run migration %d (%s): %w",
				migration.Version, migration.Name, err)
		}

		fmt.Printf("Applied migration %d: %s\n", migration.Version, migration.Name)
	}

	return nil
}

// MigrateFromFile runs migrations from SQL files
func (m *Migrator) MigrateFromFile(migrationsDir string) error {
	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Read migration files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		// If directory doesn't exist, use embedded migrations
		if os.IsNotExist(err) {
			return m.Migrate()
		}
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort files by name (assuming format like 001_initial.sql)
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// Get applied migrations
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Run migrations
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Extract version from filename (e.g., 001_initial.sql -> 1)
		var version int
		var name string
		if _, err := fmt.Sscanf(file.Name(), "%03d_%s", &version, &name); err != nil {
			continue // Skip files that don't match pattern
		}
		name = strings.TrimSuffix(name, ".sql")

		// Skip if already applied
		if _, exists := applied[version]; exists {
			continue
		}

		// Read SQL file
		sqlPath := filepath.Join(migrationsDir, file.Name())
		sqlBytes, err := os.ReadFile(sqlPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		migration := Migration{
			Version: version,
			Name:    name,
			SQL:     string(sqlBytes),
		}

		if err := m.runMigration(migration); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", file.Name(), err)
		}

		fmt.Printf("Applied migration %d: %s\n", version, name)
	}

	return nil
}

// Rollback rolls back the last n migrations
func (m *Migrator) Rollback(n int) error {
	// This is a simplified rollback - in production you'd want down migrations
	return fmt.Errorf("rollback not implemented - restore from backup")
}

// GetVersion returns the current schema version
func (m *Migrator) GetVersion() (int, error) {
	var version int
	query := "SELECT MAX(version) FROM schema_migrations"
	err := m.db.QueryRow(query).Scan(&version)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return version, err
}

// Private methods

func (m *Migrator) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := m.db.Exec(query)
	return err
}

func (m *Migrator) getAppliedMigrations() (map[int]string, error) {
	query := "SELECT version, name FROM schema_migrations"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]string)
	for rows.Next() {
		var version int
		var name string
		if err := rows.Scan(&version, &name); err != nil {
			return nil, err
		}
		applied[version] = name
	}

	return applied, rows.Err()
}

func (m *Migrator) runMigration(migration Migration) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration SQL
	if _, err := tx.Exec(migration.SQL); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// Record migration as applied
	recordQuery := `
		INSERT INTO schema_migrations (version, name)
		VALUES (?, ?)
	`
	if _, err := tx.Exec(recordQuery, migration.Version, migration.Name); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return tx.Commit()
}

// getEmbeddedMigrations returns migrations embedded in the binary
func getEmbeddedMigrations() []Migration {
	return []Migration{
		{
			Version: 1,
			Name:    "initial_schema",
			SQL: `
				-- Workspaces
				CREATE TABLE IF NOT EXISTS workspaces (
					id TEXT PRIMARY KEY,
					name TEXT UNIQUE NOT NULL,
					jira_url TEXT,
					project_key TEXT,
					credential_ref TEXT,
					config JSON,
					is_default BOOLEAN DEFAULT FALSE,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				-- Create default workspace
				INSERT OR IGNORE INTO workspaces (id, name, is_default)
				VALUES ('default', 'default', TRUE);

				-- Create unique index for default workspace
				CREATE UNIQUE INDEX IF NOT EXISTS idx_default_workspace
				ON workspaces(is_default) WHERE is_default = TRUE;

				-- Tickets
				CREATE TABLE IF NOT EXISTS tickets (
					id TEXT PRIMARY KEY,
					workspace_id TEXT NOT NULL DEFAULT 'default' REFERENCES workspaces(id) ON DELETE CASCADE,
					jira_id TEXT,
					title TEXT NOT NULL,
					description TEXT,
					custom_fields JSON,
					acceptance_criteria JSON,
					tasks JSON,
					local_hash TEXT,
					remote_hash TEXT,
					sync_status TEXT DEFAULT 'new' CHECK(sync_status IN ('new', 'synced', 'modified', 'conflict')),
					source_line INTEGER,
					last_synced TIMESTAMP,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					UNIQUE(workspace_id, jira_id)
				);

				-- Indexes
				CREATE INDEX IF NOT EXISTS idx_ticket_workspace ON tickets(workspace_id);
				CREATE INDEX IF NOT EXISTS idx_ticket_jira ON tickets(jira_id);
				CREATE INDEX IF NOT EXISTS idx_ticket_status ON tickets(sync_status);
				CREATE INDEX IF NOT EXISTS idx_ticket_modified ON tickets(updated_at);

				-- State tracking
				CREATE TABLE IF NOT EXISTS ticket_state (
					ticket_id TEXT PRIMARY KEY REFERENCES tickets(id) ON DELETE CASCADE,
					local_hash TEXT NOT NULL,
					remote_hash TEXT NOT NULL,
					last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				-- Sync operations log
				CREATE TABLE IF NOT EXISTS sync_operations (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					workspace_id TEXT REFERENCES workspaces(id),
					operation TEXT NOT NULL CHECK(operation IN ('push', 'pull', 'migrate', 'conflict_resolve')),
					file_path TEXT,
					ticket_count INTEGER,
					success_count INTEGER,
					failure_count INTEGER,
					conflict_count INTEGER,
					duration_ms INTEGER,
					error_details JSON,
					started_at TIMESTAMP,
					completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
				);

				CREATE INDEX IF NOT EXISTS idx_sync_workspace ON sync_operations(workspace_id);
				CREATE INDEX IF NOT EXISTS idx_sync_timestamp ON sync_operations(started_at);

				-- Triggers
				CREATE TRIGGER IF NOT EXISTS update_ticket_timestamp
				AFTER UPDATE ON tickets
				BEGIN
					UPDATE tickets SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
				END;

				CREATE TRIGGER IF NOT EXISTS update_workspace_timestamp
				AFTER UPDATE ON workspaces
				BEGIN
					UPDATE workspaces SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
				END;
			`,
		},
		{
			Version: 2,
			Name:    "workspace_enhancements",
			SQL: `
				-- Migration 002: Workspace enhancements for Phase 2
				-- Placeholder for future workspace schema updates
				-- This migration is reserved for Phase 2 workspace features including:
				--   - Enhanced workspace metadata
				--   - Additional sync state tracking
				--   - Performance optimizations
				-- No schema changes in this placeholder migration.
				SELECT 1;
			`,
		},
	}
}

// MigrationStatus represents the status of a migration
type MigrationStatus struct {
	Version   int
	Name      string
	Applied   bool
	AppliedAt *string
}

// GetStatus returns the status of all migrations
func (m *Migrator) GetStatus() ([]MigrationStatus, error) {
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return nil, err
	}

	var statuses []MigrationStatus
	for _, migration := range m.migrations {
		status := MigrationStatus{
			Version: migration.Version,
			Name:    migration.Name,
			Applied: false,
		}

		if _, exists := applied[migration.Version]; exists {
			status.Applied = true
			// Could fetch applied_at timestamp here if needed
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

// Export exports the current database state to a writer
func (m *Migrator) Export(w io.Writer) error {
	// Dump schema
	schemaQuery := `
		SELECT sql FROM sqlite_master
		WHERE type IN ('table', 'index', 'trigger')
		AND name NOT LIKE 'sqlite_%'
		ORDER BY type DESC, name
	`

	rows, err := m.db.Query(schemaQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Fprintln(w, "-- Ticketr Database Schema Export")
	fmt.Fprintln(w, "-- Generated by Ticketr v3.0")
	fmt.Fprintln(w)

	for rows.Next() {
		var sqlStmt sql.NullString
		if err := rows.Scan(&sqlStmt); err != nil {
			return err
		}
		if sqlStmt.Valid {
			fmt.Fprintf(w, "%s;\n\n", sqlStmt.String)
		}
	}

	return rows.Err()
}
