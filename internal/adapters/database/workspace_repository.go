package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// WorkspaceRepository implements the ports.WorkspaceRepository interface
// with SQLite backend. All operations are thread-safe through SQLite's
// internal locking mechanism.
type WorkspaceRepository struct {
	db *sql.DB
}

// NewWorkspaceRepository creates a new workspace repository instance
func NewWorkspaceRepository(db *sql.DB) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

// Create creates a new workspace in the repository.
// Returns ports.ErrWorkspaceExists if a workspace with the same name already exists.
func (r *WorkspaceRepository) Create(workspace *domain.Workspace) error {
	if workspace == nil {
		return fmt.Errorf("workspace cannot be nil")
	}

	// If this workspace is being set as default, clear other defaults first
	if workspace.IsDefault {
		if err := r.clearAllDefaults(); err != nil {
			return fmt.Errorf("failed to clear existing defaults: %w", err)
		}
	}

	// Marshal CredentialRef to JSON
	credentialJSON, err := json.Marshal(workspace.Credentials)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	query := `
		INSERT INTO workspaces (id, name, jira_url, project_key, credential_ref,
		                       is_default, last_used, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.Exec(query,
		workspace.ID,
		workspace.Name,
		workspace.JiraURL,
		workspace.ProjectKey,
		credentialJSON,
		workspace.IsDefault,
		nullTimeOrNow(workspace.LastUsed),
		workspace.CreatedAt,
		workspace.UpdatedAt,
	)

	if err != nil {
		if isSQLiteConstraintError(err) {
			return ports.ErrWorkspaceExists
		}
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	return nil
}

// Get retrieves a workspace by its unique ID.
// Returns ports.ErrWorkspaceNotFound if the workspace doesn't exist.
func (r *WorkspaceRepository) Get(id string) (*domain.Workspace, error) {
	query := `
		SELECT id, name, jira_url, project_key, credential_ref,
		       is_default, last_used, created_at, updated_at
		FROM workspaces
		WHERE id = ?
	`

	workspace := &domain.Workspace{}
	var jiraURL, projectKey sql.NullString
	var credentialJSON []byte
	var lastUsed sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&workspace.ID,
		&workspace.Name,
		&jiraURL,
		&projectKey,
		&credentialJSON,
		&workspace.IsDefault,
		&lastUsed,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ports.ErrWorkspaceNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}

	if jiraURL.Valid {
		workspace.JiraURL = jiraURL.String
	}
	if projectKey.Valid {
		workspace.ProjectKey = projectKey.String
	}

	// Unmarshal CredentialRef
	if len(credentialJSON) > 0 {
		if err := json.Unmarshal(credentialJSON, &workspace.Credentials); err != nil {
			return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
		}
	}

	if lastUsed.Valid {
		workspace.LastUsed = lastUsed.Time
	}

	return workspace, nil
}

// GetByName retrieves a workspace by its name.
// Returns ports.ErrWorkspaceNotFound if the workspace doesn't exist.
func (r *WorkspaceRepository) GetByName(name string) (*domain.Workspace, error) {
	query := `
		SELECT id, name, jira_url, project_key, credential_ref,
		       is_default, last_used, created_at, updated_at
		FROM workspaces
		WHERE name = ?
	`

	workspace := &domain.Workspace{}
	var jiraURL, projectKey sql.NullString
	var credentialJSON []byte
	var lastUsed sql.NullTime

	err := r.db.QueryRow(query, name).Scan(
		&workspace.ID,
		&workspace.Name,
		&jiraURL,
		&projectKey,
		&credentialJSON,
		&workspace.IsDefault,
		&lastUsed,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ports.ErrWorkspaceNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace by name: %w", err)
	}

	if jiraURL.Valid {
		workspace.JiraURL = jiraURL.String
	}
	if projectKey.Valid {
		workspace.ProjectKey = projectKey.String
	}

	// Unmarshal CredentialRef
	if len(credentialJSON) > 0 {
		if err := json.Unmarshal(credentialJSON, &workspace.Credentials); err != nil {
			return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
		}
	}

	if lastUsed.Valid {
		workspace.LastUsed = lastUsed.Time
	}

	return workspace, nil
}

// List returns all workspaces, ordered by last used (most recent first).
func (r *WorkspaceRepository) List() ([]*domain.Workspace, error) {
	query := `
		SELECT id, name, jira_url, project_key, credential_ref,
		       is_default, last_used, created_at, updated_at
		FROM workspaces
		ORDER BY last_used DESC, created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list workspaces: %w", err)
	}
	defer rows.Close()

	var workspaces []*domain.Workspace
	for rows.Next() {
		workspace := &domain.Workspace{}
		var jiraURL, projectKey sql.NullString
		var credentialJSON []byte
		var lastUsed sql.NullTime

		err := rows.Scan(
			&workspace.ID,
			&workspace.Name,
			&jiraURL,
			&projectKey,
			&credentialJSON,
			&workspace.IsDefault,
			&lastUsed,
			&workspace.CreatedAt,
			&workspace.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workspace: %w", err)
		}

		if jiraURL.Valid {
			workspace.JiraURL = jiraURL.String
		}
		if projectKey.Valid {
			workspace.ProjectKey = projectKey.String
		}

		// Unmarshal CredentialRef
		if len(credentialJSON) > 0 {
			if err := json.Unmarshal(credentialJSON, &workspace.Credentials); err != nil {
				return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
			}
		}

		if lastUsed.Valid {
			workspace.LastUsed = lastUsed.Time
		}

		workspaces = append(workspaces, workspace)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating workspaces: %w", err)
	}

	return workspaces, nil
}

// Update updates an existing workspace.
// Returns ports.ErrWorkspaceNotFound if the workspace doesn't exist.
func (r *WorkspaceRepository) Update(workspace *domain.Workspace) error {
	if workspace == nil {
		return fmt.Errorf("workspace cannot be nil")
	}

	// If this workspace is being set as default, clear other defaults first
	if workspace.IsDefault {
		if err := r.clearAllDefaults(); err != nil {
			return fmt.Errorf("failed to clear existing defaults: %w", err)
		}
	}

	// Marshal CredentialRef to JSON
	credentialJSON, err := json.Marshal(workspace.Credentials)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	query := `
		UPDATE workspaces
		SET name = ?, jira_url = ?, project_key = ?, credential_ref = ?,
		    is_default = ?, last_used = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query,
		workspace.Name,
		workspace.JiraURL,
		workspace.ProjectKey,
		credentialJSON,
		workspace.IsDefault,
		nullTimeOrNow(workspace.LastUsed),
		time.Now(),
		workspace.ID,
	)

	if err != nil {
		if isSQLiteConstraintError(err) {
			return ports.ErrWorkspaceExists
		}
		return fmt.Errorf("failed to update workspace: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ports.ErrWorkspaceNotFound
	}

	return nil
}

// Delete removes a workspace from the repository.
// Returns ports.ErrWorkspaceNotFound if the workspace doesn't exist.
// Cascade deletion of associated tickets is handled by the database.
func (r *WorkspaceRepository) Delete(id string) error {
	query := `DELETE FROM workspaces WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete workspace: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ports.ErrWorkspaceNotFound
	}

	return nil
}

// SetDefault marks a workspace as the default, clearing the default flag on all others.
// Returns ports.ErrWorkspaceNotFound if the workspace doesn't exist.
func (r *WorkspaceRepository) SetDefault(id string) error {
	// Use a transaction to ensure atomicity
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// First, verify the workspace exists
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM workspaces WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check workspace existence: %w", err)
	}
	if !exists {
		return ports.ErrWorkspaceNotFound
	}

	// Clear all default flags
	_, err = tx.Exec("UPDATE workspaces SET is_default = FALSE WHERE is_default = TRUE")
	if err != nil {
		return fmt.Errorf("failed to clear existing defaults: %w", err)
	}

	// Set the new default
	_, err = tx.Exec("UPDATE workspaces SET is_default = TRUE, updated_at = ? WHERE id = ?",
		time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to set default workspace: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetDefault retrieves the default workspace.
// Returns ports.ErrNoDefaultWorkspace if no default workspace is set.
func (r *WorkspaceRepository) GetDefault() (*domain.Workspace, error) {
	query := `
		SELECT id, name, jira_url, project_key, credential_ref,
		       is_default, last_used, created_at, updated_at
		FROM workspaces
		WHERE is_default = TRUE
		LIMIT 2
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query default workspace: %w", err)
	}
	defer rows.Close()

	var workspaces []*domain.Workspace
	for rows.Next() {
		workspace := &domain.Workspace{}
		var jiraURL, projectKey sql.NullString
		var credentialJSON []byte
		var lastUsed sql.NullTime

		err := rows.Scan(
			&workspace.ID,
			&workspace.Name,
			&jiraURL,
			&projectKey,
			&credentialJSON,
			&workspace.IsDefault,
			&lastUsed,
			&workspace.CreatedAt,
			&workspace.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workspace: %w", err)
		}

		if jiraURL.Valid {
			workspace.JiraURL = jiraURL.String
		}
		if projectKey.Valid {
			workspace.ProjectKey = projectKey.String
		}

		// Unmarshal CredentialRef
		if len(credentialJSON) > 0 {
			if err := json.Unmarshal(credentialJSON, &workspace.Credentials); err != nil {
				return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
			}
		}

		if lastUsed.Valid {
			workspace.LastUsed = lastUsed.Time
		}

		workspaces = append(workspaces, workspace)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating workspaces: %w", err)
	}

	if len(workspaces) == 0 {
		return nil, ports.ErrNoDefaultWorkspace
	}

	if len(workspaces) > 1 {
		return nil, ports.ErrMultipleDefaultWorkspaces
	}

	return workspaces[0], nil
}

// UpdateLastUsed updates the LastUsed timestamp for a workspace.
// This is called automatically when switching to a workspace.
func (r *WorkspaceRepository) UpdateLastUsed(id string) error {
	query := `UPDATE workspaces SET last_used = ?, updated_at = ? WHERE id = ?`

	result, err := r.db.Exec(query, time.Now(), time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update last used: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ports.ErrWorkspaceNotFound
	}

	return nil
}

// clearAllDefaults clears the is_default flag on all workspaces.
// This is a helper method used internally before setting a new default.
func (r *WorkspaceRepository) clearAllDefaults() error {
	_, err := r.db.Exec("UPDATE workspaces SET is_default = FALSE WHERE is_default = TRUE")
	if err != nil {
		return fmt.Errorf("failed to clear defaults: %w", err)
	}
	return nil
}

// Helper functions

// nullTimeOrNow returns a sql.NullTime that is NULL if the time is zero,
// otherwise it returns the time value.
func nullTimeOrNow(t time.Time) interface{} {
	if t.IsZero() {
		return nil
	}
	return t
}

// isSQLiteConstraintError checks if an error is a SQLite constraint violation.
// This typically indicates a unique constraint violation (duplicate name/ID).
func isSQLiteConstraintError(err error) bool {
	if err == nil {
		return false
	}
	// SQLite constraint errors contain "UNIQUE constraint failed" or "constraint failed"
	errMsg := err.Error()
	return contains(errMsg, "UNIQUE constraint failed") ||
		contains(errMsg, "constraint failed")
}

// contains is a simple substring check helper
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
