package ports

import (
	"errors"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

var (
	// ErrCredentialProfileNotFound is returned when a credential profile doesn't exist
	ErrCredentialProfileNotFound = errors.New("credential profile not found")

	// ErrCredentialProfileExists is returned when attempting to create a credential profile with a duplicate name
	ErrCredentialProfileExists = errors.New("credential profile already exists")
)

// CredentialProfileRepository defines the interface for credential profile persistence operations.
// Implementations must ensure thread-safety for concurrent operations.
type CredentialProfileRepository interface {
	// Create creates a new credential profile in the repository.
	// Returns ErrCredentialProfileExists if a profile with the same name already exists.
	Create(profile *domain.CredentialProfile) error

	// Get retrieves a credential profile by its unique ID.
	// Returns ErrCredentialProfileNotFound if the profile doesn't exist.
	Get(id string) (*domain.CredentialProfile, error)

	// GetByName retrieves a credential profile by its name.
	// Returns ErrCredentialProfileNotFound if the profile doesn't exist.
	GetByName(name string) (*domain.CredentialProfile, error)

	// List returns all credential profiles, ordered by name.
	List() ([]*domain.CredentialProfile, error)

	// Update updates an existing credential profile.
	// Returns ErrCredentialProfileNotFound if the profile doesn't exist.
	Update(profile *domain.CredentialProfile) error

	// Delete removes a credential profile from the repository.
	// Returns ErrCredentialProfileNotFound if the profile doesn't exist.
	// This should also check if any workspaces are using this profile and prevent deletion if so.
	Delete(id string) error

	// GetWorkspacesUsingProfile returns a list of workspace IDs that are using the specified profile.
	// This is used to prevent deletion of profiles that are in use.
	GetWorkspacesUsingProfile(profileID string) ([]string, error)
}
