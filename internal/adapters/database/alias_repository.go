package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// AliasRepository implements ports.AliasRepository using SQLite.
type AliasRepository struct {
	db *sql.DB
}

// NewAliasRepository creates a new AliasRepository instance.
func NewAliasRepository(db *sql.DB) *AliasRepository {
	return &AliasRepository{db: db}
}

// Create creates a new alias in the database.
func (r *AliasRepository) Create(alias *domain.JQLAlias) error {
	// Check for duplicate name in the same workspace
	existing, err := r.GetByName(alias.Name, alias.WorkspaceID)
	if err == nil && existing != nil {
		return ports.ErrAliasExists
	}

	// Insert the alias
	query := `
		INSERT INTO jql_aliases (id, name, jql, description, is_predefined, workspace_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.Exec(query,
		alias.ID,
		alias.Name,
		alias.JQL,
		alias.Description,
		alias.IsPredefined,
		nullString(alias.WorkspaceID),
		alias.CreatedAt,
		alias.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert alias: %w", err)
	}

	return nil
}

// Get retrieves an alias by its ID.
func (r *AliasRepository) Get(id string) (*domain.JQLAlias, error) {
	query := `
		SELECT id, name, jql, description, is_predefined, workspace_id, created_at, updated_at
		FROM jql_aliases
		WHERE id = ?
	`

	var alias domain.JQLAlias
	var workspaceID sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&alias.ID,
		&alias.Name,
		&alias.JQL,
		&alias.Description,
		&alias.IsPredefined,
		&workspaceID,
		&alias.CreatedAt,
		&alias.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ports.ErrAliasNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get alias: %w", err)
	}

	if workspaceID.Valid {
		alias.WorkspaceID = workspaceID.String
	}

	return &alias, nil
}

// GetByName retrieves an alias by its name for a specific workspace.
func (r *AliasRepository) GetByName(name string, workspaceID string) (*domain.JQLAlias, error) {
	query := `
		SELECT id, name, jql, description, is_predefined, workspace_id, created_at, updated_at
		FROM jql_aliases
		WHERE name = ? AND (workspace_id = ? OR workspace_id IS NULL)
		ORDER BY
			CASE WHEN workspace_id IS NOT NULL THEN 0 ELSE 1 END,
			is_predefined ASC
		LIMIT 1
	`

	var alias domain.JQLAlias
	var workspaceIDResult sql.NullString

	err := r.db.QueryRow(query, name, nullString(workspaceID)).Scan(
		&alias.ID,
		&alias.Name,
		&alias.JQL,
		&alias.Description,
		&alias.IsPredefined,
		&workspaceIDResult,
		&alias.CreatedAt,
		&alias.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ports.ErrAliasNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get alias: %w", err)
	}

	if workspaceIDResult.Valid {
		alias.WorkspaceID = workspaceIDResult.String
	}

	return &alias, nil
}

// List returns all aliases for a specific workspace, including global aliases.
func (r *AliasRepository) List(workspaceID string) ([]*domain.JQLAlias, error) {
	query := `
		SELECT id, name, jql, description, is_predefined, workspace_id, created_at, updated_at
		FROM jql_aliases
		WHERE workspace_id = ? OR workspace_id IS NULL
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query, nullString(workspaceID))
	if err != nil {
		return nil, fmt.Errorf("failed to list aliases: %w", err)
	}
	defer rows.Close()

	aliases := []*domain.JQLAlias{}
	for rows.Next() {
		var alias domain.JQLAlias
		var workspaceIDResult sql.NullString

		err := rows.Scan(
			&alias.ID,
			&alias.Name,
			&alias.JQL,
			&alias.Description,
			&alias.IsPredefined,
			&workspaceIDResult,
			&alias.CreatedAt,
			&alias.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alias: %w", err)
		}

		if workspaceIDResult.Valid {
			alias.WorkspaceID = workspaceIDResult.String
		}

		aliases = append(aliases, &alias)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating aliases: %w", err)
	}

	return aliases, nil
}

// ListAll returns all aliases across all workspaces.
func (r *AliasRepository) ListAll() ([]*domain.JQLAlias, error) {
	query := `
		SELECT id, name, jql, description, is_predefined, workspace_id, created_at, updated_at
		FROM jql_aliases
		ORDER BY workspace_id, name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list all aliases: %w", err)
	}
	defer rows.Close()

	aliases := []*domain.JQLAlias{}
	for rows.Next() {
		var alias domain.JQLAlias
		var workspaceID sql.NullString

		err := rows.Scan(
			&alias.ID,
			&alias.Name,
			&alias.JQL,
			&alias.Description,
			&alias.IsPredefined,
			&workspaceID,
			&alias.CreatedAt,
			&alias.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alias: %w", err)
		}

		if workspaceID.Valid {
			alias.WorkspaceID = workspaceID.String
		}

		aliases = append(aliases, &alias)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating aliases: %w", err)
	}

	return aliases, nil
}

// Update updates an existing alias.
func (r *AliasRepository) Update(alias *domain.JQLAlias) error {
	// Check if alias exists
	existing, err := r.Get(alias.ID)
	if err != nil {
		return err
	}

	// Prevent modification of predefined aliases
	if existing.IsPredefined {
		return ports.ErrCannotModifyPredefined
	}

	query := `
		UPDATE jql_aliases
		SET name = ?, jql = ?, description = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query,
		alias.Name,
		alias.JQL,
		alias.Description,
		time.Now(),
		alias.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update alias: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return ports.ErrAliasNotFound
	}

	return nil
}

// Delete removes an alias by ID.
func (r *AliasRepository) Delete(id string) error {
	// Check if alias exists and is not predefined
	existing, err := r.Get(id)
	if err != nil {
		return err
	}

	if existing.IsPredefined {
		return ports.ErrCannotModifyPredefined
	}

	query := `DELETE FROM jql_aliases WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete alias: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return ports.ErrAliasNotFound
	}

	return nil
}

// DeleteByName removes an alias by name for a specific workspace.
func (r *AliasRepository) DeleteByName(name string, workspaceID string) error {
	// Get the alias first to check if it's predefined
	alias, err := r.GetByName(name, workspaceID)
	if err != nil {
		return err
	}

	if alias.IsPredefined {
		return ports.ErrCannotModifyPredefined
	}

	return r.Delete(alias.ID)
}

// nullString converts a string to sql.NullString.
// Empty strings are converted to NULL.
func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// Helper function to convert JSON string to map (if needed in future)
func jsonToMap(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	if jsonStr == "" {
		return result, nil
	}
	err := json.Unmarshal([]byte(jsonStr), &result)
	return result, err
}
