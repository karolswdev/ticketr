package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// CredentialProfileRepository implements the ports.CredentialProfileRepository interface
// with SQLite backend. All operations are thread-safe through SQLite's
// internal locking mechanism.
type CredentialProfileRepository struct {
	db *sql.DB
}

// NewCredentialProfileRepository creates a new credential profile repository instance
func NewCredentialProfileRepository(db *sql.DB) *CredentialProfileRepository {
	return &CredentialProfileRepository{db: db}
}

// Create creates a new credential profile in the repository.
// Returns ports.ErrCredentialProfileExists if a profile with the same name already exists.
func (r *CredentialProfileRepository) Create(profile *domain.CredentialProfile) error {
	if profile == nil {
		return fmt.Errorf("credential profile cannot be nil")
	}

	// Marshal KeychainRef to JSON
	keychainJSON, err := json.Marshal(profile.KeychainRef)
	if err != nil {
		return fmt.Errorf("failed to marshal keychain ref: %w", err)
	}

	query := `
		INSERT INTO credential_profiles (id, name, jira_url, username, keychain_ref,
		                                created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.Exec(query,
		profile.ID,
		profile.Name,
		profile.JiraURL,
		profile.Username,
		keychainJSON,
		profile.CreatedAt,
		profile.UpdatedAt,
	)

	if err != nil {
		if isSQLiteConstraintError(err) {
			return ports.ErrCredentialProfileExists
		}
		return fmt.Errorf("failed to create credential profile: %w", err)
	}

	return nil
}

// Get retrieves a credential profile by its unique ID.
// Returns ports.ErrCredentialProfileNotFound if the profile doesn't exist.
func (r *CredentialProfileRepository) Get(id string) (*domain.CredentialProfile, error) {
	query := `
		SELECT id, name, jira_url, username, keychain_ref,
		       created_at, updated_at
		FROM credential_profiles
		WHERE id = ?
	`

	profile := &domain.CredentialProfile{}
	var keychainJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&profile.ID,
		&profile.Name,
		&profile.JiraURL,
		&profile.Username,
		&keychainJSON,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ports.ErrCredentialProfileNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get credential profile: %w", err)
	}

	// Unmarshal KeychainRef
	if len(keychainJSON) > 0 {
		if err := json.Unmarshal(keychainJSON, &profile.KeychainRef); err != nil {
			return nil, fmt.Errorf("failed to unmarshal keychain ref: %w", err)
		}
	}

	return profile, nil
}

// GetByName retrieves a credential profile by its name.
// Returns ports.ErrCredentialProfileNotFound if the profile doesn't exist.
func (r *CredentialProfileRepository) GetByName(name string) (*domain.CredentialProfile, error) {
	query := `
		SELECT id, name, jira_url, username, keychain_ref,
		       created_at, updated_at
		FROM credential_profiles
		WHERE name = ?
	`

	profile := &domain.CredentialProfile{}
	var keychainJSON []byte

	err := r.db.QueryRow(query, name).Scan(
		&profile.ID,
		&profile.Name,
		&profile.JiraURL,
		&profile.Username,
		&keychainJSON,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ports.ErrCredentialProfileNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get credential profile by name: %w", err)
	}

	// Unmarshal KeychainRef
	if len(keychainJSON) > 0 {
		if err := json.Unmarshal(keychainJSON, &profile.KeychainRef); err != nil {
			return nil, fmt.Errorf("failed to unmarshal keychain ref: %w", err)
		}
	}

	return profile, nil
}

// List returns all credential profiles, ordered by name.
func (r *CredentialProfileRepository) List() ([]*domain.CredentialProfile, error) {
	query := `
		SELECT id, name, jira_url, username, keychain_ref,
		       created_at, updated_at
		FROM credential_profiles
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list credential profiles: %w", err)
	}
	defer rows.Close()

	var profiles []*domain.CredentialProfile
	for rows.Next() {
		profile := &domain.CredentialProfile{}
		var keychainJSON []byte

		err := rows.Scan(
			&profile.ID,
			&profile.Name,
			&profile.JiraURL,
			&profile.Username,
			&keychainJSON,
			&profile.CreatedAt,
			&profile.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan credential profile: %w", err)
		}

		// Unmarshal KeychainRef
		if len(keychainJSON) > 0 {
			if err := json.Unmarshal(keychainJSON, &profile.KeychainRef); err != nil {
				return nil, fmt.Errorf("failed to unmarshal keychain ref: %w", err)
			}
		}

		profiles = append(profiles, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating credential profiles: %w", err)
	}

	return profiles, nil
}

// Update updates an existing credential profile.
// Returns ports.ErrCredentialProfileNotFound if the profile doesn't exist.
func (r *CredentialProfileRepository) Update(profile *domain.CredentialProfile) error {
	if profile == nil {
		return fmt.Errorf("credential profile cannot be nil")
	}

	// Marshal KeychainRef to JSON
	keychainJSON, err := json.Marshal(profile.KeychainRef)
	if err != nil {
		return fmt.Errorf("failed to marshal keychain ref: %w", err)
	}

	query := `
		UPDATE credential_profiles
		SET name = ?, jira_url = ?, username = ?, keychain_ref = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query,
		profile.Name,
		profile.JiraURL,
		profile.Username,
		keychainJSON,
		time.Now(),
		profile.ID,
	)

	if err != nil {
		if isSQLiteConstraintError(err) {
			return ports.ErrCredentialProfileExists
		}
		return fmt.Errorf("failed to update credential profile: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ports.ErrCredentialProfileNotFound
	}

	return nil
}

// Delete removes a credential profile from the repository.
// Returns ports.ErrCredentialProfileNotFound if the profile doesn't exist.
// This should also check if any workspaces are using this profile and prevent deletion if so.
func (r *CredentialProfileRepository) Delete(id string) error {
	// First, check if any workspaces are using this profile
	workspaceIDs, err := r.GetWorkspacesUsingProfile(id)
	if err != nil {
		return fmt.Errorf("failed to check workspace usage: %w", err)
	}

	if len(workspaceIDs) > 0 {
		return fmt.Errorf("cannot delete credential profile: it is used by %d workspace(s)", len(workspaceIDs))
	}

	query := `DELETE FROM credential_profiles WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete credential profile: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ports.ErrCredentialProfileNotFound
	}

	return nil
}

// GetWorkspacesUsingProfile returns a list of workspace IDs that are using the specified profile.
// This is used to prevent deletion of profiles that are in use.
func (r *CredentialProfileRepository) GetWorkspacesUsingProfile(profileID string) ([]string, error) {
	query := `
		SELECT w.id
		FROM workspaces w
		JOIN credential_profiles cp ON w.credential_profile_id = cp.id
		WHERE cp.id = ?
	`

	rows, err := r.db.Query(query, profileID)
	if err != nil {
		return nil, fmt.Errorf("failed to query workspaces using profile: %w", err)
	}
	defer rows.Close()

	var workspaceIDs []string
	for rows.Next() {
		var workspaceID string
		if err := rows.Scan(&workspaceID); err != nil {
			return nil, fmt.Errorf("failed to scan workspace ID: %w", err)
		}
		workspaceIDs = append(workspaceIDs, workspaceID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating workspace IDs: %w", err)
	}

	return workspaceIDs, nil
}
