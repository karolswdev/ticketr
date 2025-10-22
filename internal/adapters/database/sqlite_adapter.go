package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/karolswdev/ticktr/internal/parser"
	"github.com/karolswdev/ticktr/internal/state"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteAdapter implements the Repository interface with SQLite backend
type SQLiteAdapter struct {
	db           *sql.DB
	dbPath       string
	parser       *parser.Parser
	stateManager *state.StateManager
	pathResolver *services.PathResolver // NEW
	workspaceID  string                 // Default to "default" for backward compatibility
}

// NewSQLiteAdapter creates a new SQLite adapter instance with PathResolver (v3.0+).
// This is the primary constructor for v3.0 and later.
func NewSQLiteAdapter(pathResolver *services.PathResolver) (*SQLiteAdapter, error) {
	if pathResolver == nil {
		// If pathResolver is nil, create one using the singleton
		var err error
		pathResolver, err = services.GetPathResolver()
		if err != nil {
			return nil, fmt.Errorf("failed to get path resolver: %w", err)
		}
	}

	// Ensure database directory exists
	if err := pathResolver.EnsureDirectories(); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Get database path from PathResolver
	dbPath := pathResolver.DatabasePath()

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	adapter := &SQLiteAdapter{
		db:           db,
		dbPath:       dbPath,
		parser:       parser.New(),
		pathResolver: pathResolver,
		workspaceID:  "default", // Default workspace for backward compatibility
	}

	// Run migrations
	if err := adapter.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize state manager
	adapter.stateManager = state.NewStateManager("")

	return adapter, nil
}

// Deprecated: NewSQLiteAdapterWithPath is deprecated. Use NewSQLiteAdapter(pathResolver) instead.
// This function exists for backward compatibility and will be removed in v4.0.
func NewSQLiteAdapterWithPath(dbPath string) (*SQLiteAdapter, error) {
	log.Println("DEPRECATED: NewSQLiteAdapterWithPath() is deprecated, use NewSQLiteAdapter(pathResolver) instead")

	// If no path specified, use PathResolver
	if dbPath == "" {
		pr, err := services.GetPathResolver()
		if err != nil {
			return nil, fmt.Errorf("failed to get path resolver: %w", err)
		}
		return NewSQLiteAdapter(pr)
	}

	// Create temporary PathResolver with custom database path
	// This is a workaround for backward compatibility
	homeDir := filepath.Dir(filepath.Dir(dbPath))
	pr, err := services.NewPathResolverWithOptions("ticketr",
		os.Getenv,
		func() (string, error) { return homeDir, nil })
	if err != nil {
		return nil, fmt.Errorf("failed to create path resolver: %w", err)
	}

	return NewSQLiteAdapter(pr)
}

// PathResolver returns the PathResolver instance (getter)
func (a *SQLiteAdapter) PathResolver() *services.PathResolver {
	return a.pathResolver
}

// Close closes the database connection
func (a *SQLiteAdapter) Close() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

// DB returns the underlying database connection for use with repositories
func (a *SQLiteAdapter) DB() *sql.DB {
	return a.db
}

// GetTickets reads tickets from a file and syncs with database
// This maintains backward compatibility with file-based workflow
func (a *SQLiteAdapter) GetTickets(filepath string) ([]domain.Ticket, error) {
	// Parse tickets from file using existing parser
	tickets, err := a.parser.Parse(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ports.ErrFileNotFound
		}
		return nil, err
	}

	// Sync parsed tickets with database
	for _, ticket := range tickets {
		if err := a.syncTicketToDatabase(ticket, filepath); err != nil {
			return nil, fmt.Errorf("failed to sync ticket to database: %w", err)
		}
	}

	return tickets, nil
}

// SaveTickets writes tickets to a file and syncs with database
func (a *SQLiteAdapter) SaveTickets(filepath string, tickets []domain.Ticket) error {
	// First, save to file to maintain backward compatibility
	if err := a.saveToFile(filepath, tickets); err != nil {
		return err
	}

	// Then sync with database
	for _, ticket := range tickets {
		if err := a.syncTicketToDatabase(ticket, filepath); err != nil {
			return fmt.Errorf("failed to sync ticket to database: %w", err)
		}
	}

	return nil
}

// GetTicketsByWorkspace retrieves all tickets for a workspace from database
func (a *SQLiteAdapter) GetTicketsByWorkspace(workspaceID string) ([]domain.Ticket, error) {
	query := `
		SELECT id, jira_id, title, description, custom_fields,
		       acceptance_criteria, tasks, source_line, updated_at
		FROM tickets
		WHERE workspace_id = ?
		ORDER BY updated_at DESC
	`

	rows, err := a.db.Query(query, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tickets: %w", err)
	}
	defer rows.Close()

	var tickets []domain.Ticket
	for rows.Next() {
		ticket, err := a.scanTicket(rows)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	return tickets, rows.Err()
}

// GetModifiedTickets retrieves tickets modified since a given time
func (a *SQLiteAdapter) GetModifiedTickets(since time.Time) ([]domain.Ticket, error) {
	query := `
		SELECT id, jira_id, title, description, custom_fields,
		       acceptance_criteria, tasks, source_line, updated_at
		FROM tickets
		WHERE workspace_id = ? AND updated_at > ?
		ORDER BY updated_at DESC
	`

	rows, err := a.db.Query(query, a.workspaceID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to query modified tickets: %w", err)
	}
	defer rows.Close()

	var tickets []domain.Ticket
	for rows.Next() {
		ticket, err := a.scanTicket(rows)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	return tickets, rows.Err()
}

// UpdateTicketState updates the sync state of a ticket
func (a *SQLiteAdapter) UpdateTicketState(ticket domain.Ticket) error {
	localHash := a.stateManager.CalculateHash(ticket)

	query := `
		UPDATE tickets
		SET local_hash = ?, sync_status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE jira_id = ? AND workspace_id = ?
	`

	status := "modified"
	if localHash == a.getRemoteHash(ticket.JiraID) {
		status = "synced"
	}

	_, err := a.db.Exec(query, localHash, status, ticket.JiraID, a.workspaceID)
	if err != nil {
		return fmt.Errorf("failed to update ticket state: %w", err)
	}

	// Also update the state table for backward compatibility
	stateQuery := `
		INSERT OR REPLACE INTO ticket_state (ticket_id, local_hash, remote_hash)
		SELECT id, ?, remote_hash FROM tickets
		WHERE jira_id = ? AND workspace_id = ?
	`
	_, err = a.db.Exec(stateQuery, localHash, ticket.JiraID, a.workspaceID)

	return err
}

// HasChanged checks if a ticket has changed since last sync
func (a *SQLiteAdapter) HasChanged(ticket domain.Ticket) bool {
	if ticket.JiraID == "" {
		return true // New tickets always need sync
	}

	var localHash, remoteHash sql.NullString
	query := `
		SELECT local_hash, remote_hash
		FROM tickets
		WHERE jira_id = ? AND workspace_id = ?
	`

	err := a.db.QueryRow(query, ticket.JiraID, a.workspaceID).Scan(&localHash, &remoteHash)
	if err == sql.ErrNoRows {
		return true // Not in database, needs sync
	}
	if err != nil {
		return true // Error, assume needs sync
	}

	currentHash := a.stateManager.CalculateHash(ticket)
	return !localHash.Valid || currentHash != localHash.String
}

// DetectConflicts finds tickets with conflicts between local and remote
func (a *SQLiteAdapter) DetectConflicts() ([]domain.Ticket, error) {
	query := `
		SELECT id, jira_id, title, description, custom_fields,
		       acceptance_criteria, tasks, source_line, updated_at
		FROM tickets
		WHERE workspace_id = ? AND sync_status = 'conflict'
		ORDER BY updated_at DESC
	`

	rows, err := a.db.Query(query, a.workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query conflicts: %w", err)
	}
	defer rows.Close()

	var tickets []domain.Ticket
	for rows.Next() {
		ticket, err := a.scanTicket(rows)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	return tickets, rows.Err()
}

// LogSyncOperation logs a sync operation for audit
func (a *SQLiteAdapter) LogSyncOperation(op ports.SyncOperation) error {
	query := `
		INSERT INTO sync_operations
		(workspace_id, operation, file_path, ticket_count, success_count,
		 failure_count, conflict_count, duration_ms, error_details, started_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var errorJSON []byte
	if op.ErrorDetails != nil {
		errorJSON, _ = json.Marshal(op.ErrorDetails)
	}

	workspaceID := op.WorkspaceID
	if workspaceID == "" {
		workspaceID = a.workspaceID
	}

	_, err := a.db.Exec(query,
		workspaceID, op.Operation, op.FilePath, op.TicketCount,
		op.SuccessCount, op.FailureCount, op.ConflictCount,
		op.DurationMs, errorJSON, op.StartedAt)

	return err
}

// Private helper methods

func (a *SQLiteAdapter) migrate() error {
	// Ensure schema_migrations table exists first
	if err := a.createMigrationTable(); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// Run migrations in order
	migrations := []struct {
		version int
		name    string
		sql     string
	}{
		{1, "001_initial_schema", a.migration001()},
		{3, "003_credential_profiles", a.migration003()},
		{4, "004_jql_aliases", a.migration004()},
	}

	for _, migration := range migrations {
		if err := a.runMigration(migration.version, migration.name, migration.sql); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", migration.name, err)
		}
	}

	return nil
}

func (a *SQLiteAdapter) createMigrationTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := a.db.Exec(query)
	return err
}

func (a *SQLiteAdapter) runMigration(version int, name, sql string) error {
	// Check if migration already applied
	var count int
	err := a.db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Already applied
	}

	// Execute migration
	if _, err := a.db.Exec(sql); err != nil {
		return err
	}

	// Mark as applied
	_, err = a.db.Exec("INSERT INTO schema_migrations (version, name) VALUES (?, ?)", version, name)
	return err
}

func (a *SQLiteAdapter) migration001() string {
	return `
		-- Workspaces
		CREATE TABLE IF NOT EXISTS workspaces (
			id TEXT PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			jira_url TEXT,
			project_key TEXT,
			credential_ref TEXT,
			config JSON,
			is_default BOOLEAN DEFAULT FALSE,
			last_used TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Create default workspace
		INSERT OR IGNORE INTO workspaces (id, name, is_default)
		VALUES ('default', 'default', TRUE);

		-- Create unique index for default workspace
		CREATE UNIQUE INDEX IF NOT EXISTS idx_default_workspace
		ON workspaces(is_default) WHERE is_default = TRUE;

		-- Tickets table with full content storage
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

		-- Indexes for efficient querying
		CREATE INDEX IF NOT EXISTS idx_ticket_workspace ON tickets(workspace_id);
		CREATE INDEX IF NOT EXISTS idx_ticket_jira ON tickets(jira_id);
		CREATE INDEX IF NOT EXISTS idx_ticket_status ON tickets(sync_status);
		CREATE INDEX IF NOT EXISTS idx_ticket_modified ON tickets(updated_at);

		-- Composite indexes for common query patterns
		CREATE INDEX IF NOT EXISTS idx_ticket_workspace_status ON tickets(workspace_id, sync_status);
		CREATE INDEX IF NOT EXISTS idx_ticket_workspace_updated ON tickets(workspace_id, updated_at DESC);

		-- State tracking (replaces .ticketr.state file)
		CREATE TABLE IF NOT EXISTS ticket_state (
			ticket_id TEXT PRIMARY KEY REFERENCES tickets(id) ON DELETE CASCADE,
			local_hash TEXT NOT NULL,
			remote_hash TEXT NOT NULL,
			last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Sync operations log for audit trail
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

		-- Trigger to update updated_at timestamp
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
	`
}

func (a *SQLiteAdapter) migration003() string {
	return `
		-- Create credential_profiles table
		CREATE TABLE IF NOT EXISTS credential_profiles (
			id TEXT PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			jira_url TEXT NOT NULL,
			username TEXT NOT NULL,
			keychain_ref TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		-- Add credential_profile_id foreign key to workspaces table
		-- This allows workspaces to optionally reference a credential profile
		ALTER TABLE workspaces
			ADD COLUMN credential_profile_id TEXT REFERENCES credential_profiles(id);

		-- Create indexes for efficient querying
		CREATE INDEX IF NOT EXISTS idx_credential_profile_name ON credential_profiles(name);
		CREATE INDEX IF NOT EXISTS idx_workspace_profile ON workspaces(credential_profile_id);

		-- Add trigger to update updated_at timestamp for credential_profiles
		CREATE TRIGGER IF NOT EXISTS update_credential_profile_timestamp
		AFTER UPDATE ON credential_profiles
		BEGIN
			UPDATE credential_profiles SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;
	`
}

func (a *SQLiteAdapter) migration004() string {
	return `
		-- JQL aliases table
		CREATE TABLE IF NOT EXISTS jql_aliases (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			jql TEXT NOT NULL,
			description TEXT,
			is_predefined BOOLEAN DEFAULT FALSE,
			workspace_id TEXT REFERENCES workspaces(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(name, workspace_id)
		);

		-- Indexes for efficient querying
		CREATE INDEX IF NOT EXISTS idx_alias_workspace ON jql_aliases(workspace_id);
		CREATE INDEX IF NOT EXISTS idx_alias_name ON jql_aliases(name);
		CREATE INDEX IF NOT EXISTS idx_alias_predefined ON jql_aliases(is_predefined);

		-- Composite index for common query pattern (name + workspace)
		CREATE INDEX IF NOT EXISTS idx_alias_name_workspace ON jql_aliases(name, workspace_id);
		CREATE INDEX IF NOT EXISTS idx_alias_workspace_name ON jql_aliases(workspace_id, name);

		-- Trigger to update updated_at timestamp
		CREATE TRIGGER IF NOT EXISTS update_alias_timestamp
		AFTER UPDATE ON jql_aliases
		BEGIN
			UPDATE jql_aliases SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END;
	`
}

// checkSchemaVersion verifies the database schema version compatibility.
func (a *SQLiteAdapter) checkSchemaVersion() error {
	const currentSchemaVersion = 4

	var version int
	err := a.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check schema version: %w", err)
	}

	if version > currentSchemaVersion {
		return fmt.Errorf("database schema version %d is newer than expected %d, please upgrade Ticketr",
			version, currentSchemaVersion)
	}

	return nil
}

func (a *SQLiteAdapter) syncTicketToDatabase(ticket domain.Ticket, filepath string) error {
	// Extract workspace ID from filepath if it matches pattern "workspace-{id}.md"
	// Otherwise use the default workspace ID
	workspaceID := a.workspaceID
	if len(filepath) > 10 && filepath[:10] == "workspace-" {
		// Extract workspace ID from "workspace-{id}.md" pattern
		endIdx := len(filepath)
		if len(filepath) > 3 && filepath[len(filepath)-3:] == ".md" {
			endIdx = len(filepath) - 3
		}
		workspaceID = filepath[10:endIdx]
	}

	// Generate a unique ID if ticket doesn't have a JIRA ID
	ticketID := ticket.JiraID
	if ticketID == "" {
		ticketID = fmt.Sprintf("local-%s-%d", filepath, ticket.SourceLine)
	}

	// Marshal complex fields to JSON
	customFieldsJSON, _ := json.Marshal(ticket.CustomFields)
	acceptanceCriteriaJSON, _ := json.Marshal(ticket.AcceptanceCriteria)
	tasksJSON, _ := json.Marshal(ticket.Tasks)

	// Calculate hashes
	localHash := a.stateManager.CalculateHash(ticket)

	query := `
		INSERT OR REPLACE INTO tickets
		(id, workspace_id, jira_id, title, description, custom_fields,
		 acceptance_criteria, tasks, local_hash, source_line, sync_status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := a.db.Exec(query,
		ticketID, workspaceID, ticket.JiraID, ticket.Title, ticket.Description,
		customFieldsJSON, acceptanceCriteriaJSON, tasksJSON,
		localHash, ticket.SourceLine, "new")

	return err
}

func (a *SQLiteAdapter) scanTicket(rows *sql.Rows) (domain.Ticket, error) {
	var ticket domain.Ticket
	var id string
	var jiraID, description sql.NullString
	var customFieldsJSON, acceptanceCriteriaJSON, tasksJSON []byte
	var sourceLine sql.NullInt64
	var updatedAt time.Time

	err := rows.Scan(&id, &jiraID, &ticket.Title, &description,
		&customFieldsJSON, &acceptanceCriteriaJSON, &tasksJSON,
		&sourceLine, &updatedAt)
	if err != nil {
		return ticket, fmt.Errorf("failed to scan ticket: %w", err)
	}

	if jiraID.Valid {
		ticket.JiraID = jiraID.String
	}
	if description.Valid {
		ticket.Description = description.String
	}
	if sourceLine.Valid {
		ticket.SourceLine = int(sourceLine.Int64)
	}

	// Unmarshal JSON fields
	if len(customFieldsJSON) > 0 {
		json.Unmarshal(customFieldsJSON, &ticket.CustomFields)
	}
	if len(acceptanceCriteriaJSON) > 0 {
		json.Unmarshal(acceptanceCriteriaJSON, &ticket.AcceptanceCriteria)
	}
	if len(tasksJSON) > 0 {
		json.Unmarshal(tasksJSON, &ticket.Tasks)
	}

	return ticket, nil
}

func (a *SQLiteAdapter) saveToFile(filepath string, tickets []domain.Ticket) error {
	// Use existing file repository for backward compatibility
	fileRepo := &fileRepository{parser: a.parser}
	return fileRepo.SaveTickets(filepath, tickets)
}

func (a *SQLiteAdapter) getRemoteHash(jiraID string) string {
	var remoteHash sql.NullString
	query := `SELECT remote_hash FROM tickets WHERE jira_id = ? AND workspace_id = ?`
	a.db.QueryRow(query, jiraID, a.workspaceID).Scan(&remoteHash)
	if remoteHash.Valid {
		return remoteHash.String
	}
	return ""
}

// Temporary file repository for backward compatibility
type fileRepository struct {
	parser *parser.Parser
}

func (r *fileRepository) SaveTickets(filepath string, tickets []domain.Ticket) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	for i, ticket := range tickets {
		if ticket.JiraID != "" {
			fmt.Fprintf(file, "# TICKET: [%s] %s\n", ticket.JiraID, ticket.Title)
		} else {
			fmt.Fprintf(file, "# TICKET: %s\n", ticket.Title)
		}
		fmt.Fprintln(file)

		if ticket.Description != "" {
			fmt.Fprintln(file, "## Description")
			fmt.Fprintln(file, ticket.Description)
			fmt.Fprintln(file)
		}

		if len(ticket.CustomFields) > 0 {
			fmt.Fprintln(file, "## Fields")
			for key, value := range ticket.CustomFields {
				fmt.Fprintf(file, "%s: %s\n", key, value)
			}
			fmt.Fprintln(file)
		}

		if len(ticket.AcceptanceCriteria) > 0 {
			fmt.Fprintln(file, "## Acceptance Criteria")
			for _, ac := range ticket.AcceptanceCriteria {
				fmt.Fprintf(file, "- %s\n", ac)
			}
			fmt.Fprintln(file)
		}

		if len(ticket.Tasks) > 0 {
			fmt.Fprintln(file, "## Tasks")
			for _, task := range ticket.Tasks {
				if task.JiraID != "" {
					fmt.Fprintf(file, "- [%s] %s\n", task.JiraID, task.Title)
				} else {
					fmt.Fprintf(file, "- %s\n", task.Title)
				}
			}
			fmt.Fprintln(file)
		}

		if i < len(tickets)-1 {
			fmt.Fprintln(file)
		}
	}

	return nil
}

