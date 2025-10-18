package ports

import (
	"errors"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

var (
	// ErrWorkspaceNotFound is returned when a workspace doesn't exist
	ErrWorkspaceNotFound = errors.New("workspace not found")

	// ErrWorkspaceExists is returned when attempting to create a workspace with a duplicate name
	ErrWorkspaceExists = errors.New("workspace already exists")

	// ErrNoDefaultWorkspace is returned when no default workspace is set
	ErrNoDefaultWorkspace = errors.New("no default workspace configured")

	// ErrMultipleDefaultWorkspaces is returned when multiple workspaces are marked as default
	ErrMultipleDefaultWorkspaces = errors.New("multiple default workspaces found")
)

// WorkspaceRepository defines the interface for workspace persistence operations.
// Implementations must ensure thread-safety for concurrent operations.
type WorkspaceRepository interface {
	// Create creates a new workspace in the repository.
	// Returns ErrWorkspaceExists if a workspace with the same name already exists.
	Create(workspace *domain.Workspace) error

	// Get retrieves a workspace by its unique ID.
	// Returns ErrWorkspaceNotFound if the workspace doesn't exist.
	Get(id string) (*domain.Workspace, error)

	// GetByName retrieves a workspace by its name.
	// Returns ErrWorkspaceNotFound if the workspace doesn't exist.
	GetByName(name string) (*domain.Workspace, error)

	// List returns all workspaces, ordered by last used (most recent first).
	List() ([]*domain.Workspace, error)

	// Update updates an existing workspace.
	// Returns ErrWorkspaceNotFound if the workspace doesn't exist.
	Update(workspace *domain.Workspace) error

	// Delete removes a workspace from the repository.
	// Returns ErrWorkspaceNotFound if the workspace doesn't exist.
	// Cascade deletion of associated tickets is handled by the database.
	Delete(id string) error

	// SetDefault marks a workspace as the default, clearing the default flag on all others.
	// Returns ErrWorkspaceNotFound if the workspace doesn't exist.
	SetDefault(id string) error

	// GetDefault retrieves the default workspace.
	// Returns ErrNoDefaultWorkspace if no default workspace is set.
	GetDefault() (*domain.Workspace, error)

	// UpdateLastUsed updates the LastUsed timestamp for a workspace.
	// This is called automatically when switching to a workspace.
	UpdateLastUsed(id string) error
}

// CredentialStore defines the interface for secure credential storage.
// Implementations should use the OS keychain/credential manager.
type CredentialStore interface {
	// Store saves credentials securely and returns a reference.
	Store(workspaceID string, config domain.WorkspaceConfig) (domain.CredentialRef, error)

	// Retrieve fetches credentials using a reference.
	Retrieve(ref domain.CredentialRef) (*domain.WorkspaceConfig, error)

	// Delete removes stored credentials.
	Delete(ref domain.CredentialRef) error

	// List returns all credential references for auditing.
	List() ([]domain.CredentialRef, error)
}
