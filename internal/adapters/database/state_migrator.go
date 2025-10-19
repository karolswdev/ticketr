package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/karolswdev/ticktr/internal/parser"
	"github.com/karolswdev/ticktr/internal/state"
)

// StateMigrator handles migration from file-based state to SQLite
type StateMigrator struct {
	db      *sql.DB
	parser  *parser.Parser
	adapter *SQLiteAdapter
	verbose bool
}

// NewStateMigrator creates a new state migrator
func NewStateMigrator(pathResolver *services.PathResolver, verbose bool) (*StateMigrator, error) {
	adapter, err := NewSQLiteAdapter(pathResolver)
	if err != nil {
		return nil, fmt.Errorf("failed to create SQLite adapter: %w", err)
	}

	return &StateMigrator{
		db:      adapter.db,
		parser:  parser.New(),
		adapter: adapter,
		verbose: verbose,
	}, nil
}

// MigrateDirectory migrates all .ticketr.state files in a directory tree
func (m *StateMigrator) MigrateDirectory(rootPath string) (*MigrationReport, error) {
	report := &MigrationReport{
		StartedAt: time.Now(),
		Projects:  make([]ProjectMigration, 0),
	}

	// Find all .ticketr.state files
	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip directories with access errors
		}

		// Look for .ticketr.state files
		if d.Name() == ".ticketr.state" {
			projectDir := filepath.Dir(path)
			if m.verbose {
				fmt.Printf("Found project at: %s\n", projectDir)
			}

			projectMigration := m.migrateProject(projectDir)
			report.Projects = append(report.Projects, projectMigration)
			report.TotalProjects++

			if projectMigration.Success {
				report.SuccessfulProjects++
			} else {
				report.FailedProjects++
			}

			report.TotalTickets += projectMigration.TicketCount
		}

		return nil
	})

	if err != nil {
		return report, fmt.Errorf("failed to walk directory: %w", err)
	}

	report.CompletedAt = time.Now()
	report.DurationMs = int(report.CompletedAt.Sub(report.StartedAt).Milliseconds())

	return report, nil
}

// MigrateProject migrates a single project directory
func (m *StateMigrator) MigrateProject(projectPath string) (*ProjectMigration, error) {
	migration := m.migrateProject(projectPath)
	if !migration.Success {
		return &migration, fmt.Errorf("migration failed: %s", migration.Error)
	}
	return &migration, nil
}

// MigrateStateFile migrates a specific .ticketr.state file
func (m *StateMigrator) MigrateStateFile(stateFilePath string) error {
	// Load the state file
	stateManager := state.NewStateManager(stateFilePath)
	if err := stateManager.Load(); err != nil {
		return fmt.Errorf("failed to load state file: %w", err)
	}

	// Begin transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Migrate state entries
	migrated := 0
	projectDir := filepath.Dir(stateFilePath)

	// Find markdown files to associate with state
	markdownFiles, err := m.findMarkdownFiles(projectDir)
	if err != nil {
		if m.verbose {
			fmt.Printf("Warning: couldn't find markdown files in %s: %v\n", projectDir, err)
		}
	}

	// Parse tickets from markdown files
	var allTickets []domain.Ticket
	for _, mdFile := range markdownFiles {
		tickets, err := m.parser.Parse(mdFile)
		if err != nil {
			if m.verbose {
				fmt.Printf("Warning: couldn't parse %s: %v\n", mdFile, err)
			}
			continue
		}
		allTickets = append(allTickets, tickets...)
	}

	// Migrate tickets to database
	for _, ticket := range allTickets {
		if err := m.migrateTicket(tx, ticket, projectDir); err != nil {
			if m.verbose {
				fmt.Printf("Warning: couldn't migrate ticket %s: %v\n", ticket.JiraID, err)
			}
		} else {
			migrated++
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	if m.verbose {
		fmt.Printf("Migrated %d tickets from %s\n", migrated, stateFilePath)
	}

	return nil
}

// CreateWorkspaceFromProject creates a workspace from a project directory
func (m *StateMigrator) CreateWorkspaceFromProject(projectPath string, workspaceName string) error {
	// Extract project info from environment or config
	envFile := filepath.Join(projectPath, ".env")
	config := m.readProjectConfig(envFile)

	// Create workspace
	query := `
		INSERT OR REPLACE INTO workspaces
		(id, name, jira_url, project_key, is_default, created_at)
		VALUES (?, ?, ?, ?, FALSE, CURRENT_TIMESTAMP)
	`

	workspaceID := fmt.Sprintf("ws-%s-%d", workspaceName, time.Now().Unix())
	_, err := m.db.Exec(query, workspaceID, workspaceName,
		config["JIRA_URL"], config["JIRA_PROJECT_KEY"])

	if err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	if m.verbose {
		fmt.Printf("Created workspace '%s' for project at %s\n", workspaceName, projectPath)
	}

	return nil
}

// Private helper methods

func (m *StateMigrator) migrateProject(projectPath string) ProjectMigration {
	migration := ProjectMigration{
		ProjectPath: projectPath,
		StartedAt:   time.Now(),
	}

	// Get project name from directory
	migration.ProjectName = filepath.Base(projectPath)

	// Check for state file
	stateFile := filepath.Join(projectPath, ".ticketr.state")
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		migration.Error = "No .ticketr.state file found"
		return migration
	}

	// Create workspace
	workspaceName := migration.ProjectName
	if err := m.CreateWorkspaceFromProject(projectPath, workspaceName); err != nil {
		migration.Error = fmt.Sprintf("Failed to create workspace: %v", err)
		return migration
	}
	migration.WorkspaceName = workspaceName

	// Migrate state file
	if err := m.MigrateStateFile(stateFile); err != nil {
		migration.Error = fmt.Sprintf("Failed to migrate state: %v", err)
		return migration
	}

	// Count migrated tickets
	var count int
	query := `
		SELECT COUNT(*) FROM tickets
		WHERE workspace_id IN (
			SELECT id FROM workspaces WHERE name = ?
		)
	`
	m.db.QueryRow(query, workspaceName).Scan(&count)
	migration.TicketCount = count

	migration.Success = true
	migration.CompletedAt = time.Now()

	return migration
}

func (m *StateMigrator) migrateTicket(tx *sql.Tx, ticket domain.Ticket, projectPath string) error {
	// Generate ticket ID
	ticketID := ticket.JiraID
	if ticketID == "" {
		ticketID = fmt.Sprintf("local-%s-%d", projectPath, ticket.SourceLine)
	}

	// Marshal JSON fields
	customFieldsJSON, _ := json.Marshal(ticket.CustomFields)
	acceptanceCriteriaJSON, _ := json.Marshal(ticket.AcceptanceCriteria)
	tasksJSON, _ := json.Marshal(ticket.Tasks)

	// Calculate hash
	stateManager := state.NewStateManager("")
	localHash := stateManager.CalculateHash(ticket)

	query := `
		INSERT OR REPLACE INTO tickets
		(id, workspace_id, jira_id, title, description, custom_fields,
		 acceptance_criteria, tasks, local_hash, remote_hash, sync_status, source_line)
		VALUES (?, 'default', ?, ?, ?, ?, ?, ?, ?, ?, 'synced', ?)
	`

	_, err := tx.Exec(query,
		ticketID, ticket.JiraID, ticket.Title, ticket.Description,
		customFieldsJSON, acceptanceCriteriaJSON, tasksJSON,
		localHash, localHash, ticket.SourceLine)

	return err
}

func (m *StateMigrator) findMarkdownFiles(dir string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) == ".md" {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}

	return files, nil
}

func (m *StateMigrator) readProjectConfig(envFile string) map[string]string {
	config := make(map[string]string)

	// Try to read .env file
	data, err := os.ReadFile(envFile)
	if err != nil {
		return config
	}

	// Simple .env parser (not production-ready, but sufficient for migration)
	lines := string(data)
	for _, line := range filepath.SplitList(lines) {
		if line == "" || line[0] == '#' {
			continue
		}

		parts := filepath.SplitList(line)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			// Remove quotes if present
			if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
				value = value[1 : len(value)-1]
			}
			config[key] = value
		}
	}

	return config
}

// BackupStateFile creates a backup of the state file before migration.
//
// Backup Procedure:
// -----------------
// Backups are stored in the same directory as the original state file with
// the naming convention: <original-filename>.backup-<unix-timestamp>
//
// Example:
//
//	Original: /project/.ticketr.state
//	Backup:   /project/.ticketr.state.backup-1234567890
//
// The backup location is logged to stdout when verbose mode is enabled.
//
// Restore Procedure:
// ------------------
// To restore from a backup:
//  1. Locate the backup file (same directory as .ticketr.state)
//  2. Copy the backup to the original location:
//     cp .ticketr.state.backup-<timestamp> .ticketr.state
//  3. Re-run the migration if needed
//
// Important Notes:
// - Backups are NOT automatically deleted
// - Multiple backups may exist for the same state file
// - Choose the most recent backup based on the timestamp
// - Backups are plain text files and can be inspected before restore
//
// Returns the full path to the created backup file.
func (m *StateMigrator) BackupStateFile(stateFilePath string) (string, error) {
	backupPath := fmt.Sprintf("%s.backup-%d", stateFilePath, time.Now().Unix())

	input, err := os.ReadFile(stateFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read state file: %w", err)
	}

	if err := os.WriteFile(backupPath, input, 0644); err != nil {
		return "", fmt.Errorf("failed to write backup: %w", err)
	}

	// Log backup location for user reference
	if m.verbose {
		fmt.Printf("Backup created: %s\n", backupPath)
		fmt.Printf("To restore: cp %s %s\n", backupPath, stateFilePath)
	}

	return backupPath, nil
}

// MigrationReport represents the results of a migration
type MigrationReport struct {
	StartedAt          time.Time
	CompletedAt        time.Time
	DurationMs         int
	TotalProjects      int
	SuccessfulProjects int
	FailedProjects     int
	TotalTickets       int
	Projects           []ProjectMigration
}

// ProjectMigration represents the migration of a single project
type ProjectMigration struct {
	ProjectPath   string
	ProjectName   string
	WorkspaceName string
	TicketCount   int
	Success       bool
	Error         string
	StartedAt     time.Time
	CompletedAt   time.Time
}

// GenerateReport generates a human-readable migration report
func (r *MigrationReport) GenerateReport() string {
	report := fmt.Sprintf(`
Migration Report
================
Started: %s
Completed: %s
Duration: %dms

Summary:
--------
Total Projects: %d
Successful: %d
Failed: %d
Total Tickets: %d

Projects:
---------
`, r.StartedAt.Format(time.RFC3339),
		r.CompletedAt.Format(time.RFC3339),
		r.DurationMs,
		r.TotalProjects,
		r.SuccessfulProjects,
		r.FailedProjects,
		r.TotalTickets)

	for _, p := range r.Projects {
		status := "✓"
		if !p.Success {
			status = "✗"
		}
		report += fmt.Sprintf("%s %s (%s) - %d tickets",
			status, p.ProjectName, p.ProjectPath, p.TicketCount)
		if p.Error != "" {
			report += fmt.Sprintf(" [Error: %s]", p.Error)
		}
		report += "\n"
	}

	return report
}
