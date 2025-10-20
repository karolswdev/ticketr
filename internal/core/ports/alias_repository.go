package ports

import (
	"errors"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

var (
	// ErrAliasNotFound is returned when an alias doesn't exist
	ErrAliasNotFound = errors.New("alias not found")

	// ErrAliasExists is returned when attempting to create an alias with a duplicate name
	ErrAliasExists = errors.New("alias already exists")

	// ErrCannotModifyPredefined is returned when attempting to modify or delete a predefined alias
	ErrCannotModifyPredefined = errors.New("cannot modify predefined alias")

	// ErrCircularReference is returned when alias expansion would create a circular reference
	ErrCircularReference = errors.New("circular alias reference detected")
)

// AliasRepository defines the interface for JQL alias persistence operations.
// Implementations must ensure thread-safety for concurrent operations.
type AliasRepository interface {
	// Create creates a new alias in the repository.
	// Returns ErrAliasExists if an alias with the same name already exists for the workspace.
	Create(alias *domain.JQLAlias) error

	// Get retrieves an alias by its unique ID.
	// Returns ErrAliasNotFound if the alias doesn't exist.
	Get(id string) (*domain.JQLAlias, error)

	// GetByName retrieves an alias by its name for a specific workspace.
	// workspaceID can be empty string for global aliases.
	// Returns ErrAliasNotFound if the alias doesn't exist.
	GetByName(name string, workspaceID string) (*domain.JQLAlias, error)

	// List returns all aliases for a specific workspace, including global aliases.
	// workspaceID can be empty string to get only global aliases.
	// Returns aliases ordered by name.
	List(workspaceID string) ([]*domain.JQLAlias, error)

	// ListAll returns all aliases across all workspaces (for admin/debugging purposes).
	ListAll() ([]*domain.JQLAlias, error)

	// Update updates an existing alias.
	// Returns ErrAliasNotFound if the alias doesn't exist.
	// Returns ErrCannotModifyPredefined if attempting to modify a predefined alias.
	Update(alias *domain.JQLAlias) error

	// Delete removes an alias from the repository.
	// Returns ErrAliasNotFound if the alias doesn't exist.
	// Returns ErrCannotModifyPredefined if attempting to delete a predefined alias.
	Delete(id string) error

	// DeleteByName removes an alias by name for a specific workspace.
	// Returns ErrAliasNotFound if the alias doesn't exist.
	// Returns ErrCannotModifyPredefined if attempting to delete a predefined alias.
	DeleteByName(name string, workspaceID string) error
}
