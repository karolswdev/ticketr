package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// WorkspaceService provides business logic for managing workspaces.
// It ensures thread-safe workspace switching and credential management.
type WorkspaceService struct {
	repo          ports.WorkspaceRepository
	credStore     ports.CredentialStore
	currentMutex  sync.RWMutex
	currentCache  *domain.Workspace
}

// NewWorkspaceService creates a new WorkspaceService instance.
func NewWorkspaceService(repo ports.WorkspaceRepository, credStore ports.CredentialStore) *WorkspaceService {
	return &WorkspaceService{
		repo:      repo,
		credStore: credStore,
	}
}

// Create creates a new workspace with the given configuration.
// It validates the workspace name and configuration before creating.
// Credentials are stored securely in the OS keychain.
func (s *WorkspaceService) Create(name string, config domain.WorkspaceConfig) error {
	// Validate workspace name
	if err := domain.ValidateWorkspaceName(name); err != nil {
		return fmt.Errorf("invalid workspace name: %w", err)
	}

	// Validate configuration
	if err := domain.ValidateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Check if workspace already exists
	existing, err := s.repo.GetByName(name)
	if err == nil && existing != nil {
		return ports.ErrWorkspaceExists
	}

	// Generate unique ID
	id := uuid.New().String()

	// Store credentials securely
	credRef, err := s.credStore.Store(id, config)
	if err != nil {
		return fmt.Errorf("failed to store credentials: %w", err)
	}

	// Create workspace
	now := time.Now()
	workspace := &domain.Workspace{
		ID:          id,
		Name:        name,
		JiraURL:     config.JiraURL,
		ProjectKey:  config.ProjectKey,
		Credentials: credRef,
		IsDefault:   false,
		LastUsed:    now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Validate workspace
	if err := workspace.Validate(); err != nil {
		// Clean up credentials if validation fails
		s.credStore.Delete(credRef)
		return fmt.Errorf("workspace validation failed: %w", err)
	}

	// Persist workspace
	if err := s.repo.Create(workspace); err != nil {
		// Clean up credentials if creation fails
		s.credStore.Delete(credRef)
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	// If this is the first workspace, make it the default
	workspaces, _ := s.repo.List()
	if len(workspaces) == 1 {
		s.repo.SetDefault(id)
		workspace.IsDefault = true
	}

	return nil
}

// Switch changes the current workspace to the specified one.
// This is a thread-safe operation that updates the workspace's last used time.
func (s *WorkspaceService) Switch(name string) error {
	// Get workspace by name
	workspace, err := s.repo.GetByName(name)
	if err != nil {
		if err == ports.ErrWorkspaceNotFound {
			return fmt.Errorf("workspace '%s' not found", name)
		}
		return fmt.Errorf("failed to get workspace: %w", err)
	}

	// Update last used timestamp
	if err := s.repo.UpdateLastUsed(workspace.ID); err != nil {
		return fmt.Errorf("failed to update workspace timestamp: %w", err)
	}

	// Update the workspace object
	workspace.Touch()

	// Thread-safe update of current workspace
	s.currentMutex.Lock()
	s.currentCache = workspace
	s.currentMutex.Unlock()

	return nil
}

// List returns all workspaces ordered by last used time.
func (s *WorkspaceService) List() ([]domain.Workspace, error) {
	workspaces, err := s.repo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list workspaces: %w", err)
	}

	// Convert pointer slice to value slice
	result := make([]domain.Workspace, len(workspaces))
	for i, ws := range workspaces {
		result[i] = *ws
	}

	return result, nil
}

// Current returns the currently active workspace.
// It returns the most recently used workspace (based on last_used timestamp).
// If no workspace is available, it returns an error.
func (s *WorkspaceService) Current() (*domain.Workspace, error) {
	// Check cache first (thread-safe read)
	s.currentMutex.RLock()
	cached := s.currentCache
	s.currentMutex.RUnlock()

	if cached != nil {
		return cached, nil
	}

	// Get most recently used workspace (List returns workspaces ordered by last_used DESC)
	workspaces, err := s.repo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to get current workspace: %w", err)
	}

	if len(workspaces) == 0 {
		return nil, fmt.Errorf("no workspaces configured")
	}

	// The first workspace in the list is the most recently used
	workspace := workspaces[0]

	// Update cache (thread-safe write)
	s.currentMutex.Lock()
	s.currentCache = workspace
	s.currentMutex.Unlock()

	return workspace, nil
}

// Delete removes a workspace and its associated credentials.
// It prevents deletion of the default workspace if it's the only workspace.
func (s *WorkspaceService) Delete(name string) error {
	// Get workspace by name
	workspace, err := s.repo.GetByName(name)
	if err != nil {
		if err == ports.ErrWorkspaceNotFound {
			return fmt.Errorf("workspace '%s' not found", name)
		}
		return fmt.Errorf("failed to get workspace: %w", err)
	}

	// Check if this is the only workspace
	workspaces, err := s.repo.List()
	if err != nil {
		return fmt.Errorf("failed to list workspaces: %w", err)
	}

	if len(workspaces) == 1 {
		return fmt.Errorf("cannot delete the only workspace")
	}

	// If deleting the default workspace, clear the cache and set a new default
	if workspace.IsDefault {
		s.currentMutex.Lock()
		s.currentCache = nil
		s.currentMutex.Unlock()

		// Set another workspace as default
		for _, ws := range workspaces {
			if ws.ID != workspace.ID {
				if err := s.repo.SetDefault(ws.ID); err != nil {
					return fmt.Errorf("failed to set new default workspace: %w", err)
				}
				break
			}
		}
	}

	// Clear current cache if deleting the current workspace
	s.currentMutex.Lock()
	if s.currentCache != nil && s.currentCache.ID == workspace.ID {
		s.currentCache = nil
	}
	s.currentMutex.Unlock()

	// Delete credentials
	if err := s.credStore.Delete(workspace.Credentials); err != nil {
		// Log warning but continue with deletion
		// Credentials might already be deleted or keychain unavailable
	}

	// Delete workspace (cascades to tickets)
	if err := s.repo.Delete(workspace.ID); err != nil {
		return fmt.Errorf("failed to delete workspace: %w", err)
	}

	return nil
}

// SetDefault sets the specified workspace as the default.
func (s *WorkspaceService) SetDefault(name string) error {
	// Get workspace by name
	workspace, err := s.repo.GetByName(name)
	if err != nil {
		if err == ports.ErrWorkspaceNotFound {
			return fmt.Errorf("workspace '%s' not found", name)
		}
		return fmt.Errorf("failed to get workspace: %w", err)
	}

	// Set as default in repository
	if err := s.repo.SetDefault(workspace.ID); err != nil {
		return fmt.Errorf("failed to set default workspace: %w", err)
	}

	return nil
}

// GetConfig retrieves the configuration (including credentials) for a workspace.
// This should be used carefully as it exposes sensitive information.
func (s *WorkspaceService) GetConfig(name string) (*domain.WorkspaceConfig, error) {
	// Get workspace by name
	workspace, err := s.repo.GetByName(name)
	if err != nil {
		if err == ports.ErrWorkspaceNotFound {
			return nil, fmt.Errorf("workspace '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to get workspace: %w", err)
	}

	// Retrieve credentials
	config, err := s.credStore.Retrieve(workspace.Credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve credentials: %w", err)
	}

	return config, nil
}

// UpdateConfig updates the configuration for a workspace.
// This updates both the workspace metadata and credentials.
func (s *WorkspaceService) UpdateConfig(name string, config domain.WorkspaceConfig) error {
	// Validate configuration
	if err := domain.ValidateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Get workspace by name
	workspace, err := s.repo.GetByName(name)
	if err != nil {
		if err == ports.ErrWorkspaceNotFound {
			return fmt.Errorf("workspace '%s' not found", name)
		}
		return fmt.Errorf("failed to get workspace: %w", err)
	}

	// Update credentials
	credRef, err := s.credStore.Store(workspace.ID, config)
	if err != nil {
		return fmt.Errorf("failed to update credentials: %w", err)
	}

	// Update workspace
	workspace.JiraURL = config.JiraURL
	workspace.ProjectKey = config.ProjectKey
	workspace.Credentials = credRef
	workspace.UpdatedAt = time.Now()

	if err := s.repo.Update(workspace); err != nil {
		return fmt.Errorf("failed to update workspace: %w", err)
	}

	// Clear cache if updating current workspace
	s.currentMutex.Lock()
	if s.currentCache != nil && s.currentCache.ID == workspace.ID {
		s.currentCache = nil
	}
	s.currentMutex.Unlock()

	return nil
}
