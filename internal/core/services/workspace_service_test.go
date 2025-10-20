package services

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// MockWorkspaceRepository implements the WorkspaceRepository interface for testing
type MockWorkspaceRepository struct {
	workspaces map[string]*domain.Workspace
	mu         sync.RWMutex
	createErr  error
	getErr     error
	listErr    error
	updateErr  error
	deleteErr  error
}

func NewMockWorkspaceRepository() *MockWorkspaceRepository {
	return &MockWorkspaceRepository{
		workspaces: make(map[string]*domain.Workspace),
	}
}

func (m *MockWorkspaceRepository) Create(workspace *domain.Workspace) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.createErr != nil {
		return m.createErr
	}

	if _, exists := m.workspaces[workspace.ID]; exists {
		return errors.New("workspace already exists")
	}

	m.workspaces[workspace.ID] = workspace
	return nil
}

func (m *MockWorkspaceRepository) Get(id string) (*domain.Workspace, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.getErr != nil {
		return nil, m.getErr
	}

	ws, exists := m.workspaces[id]
	if !exists {
		return nil, errors.New("workspace not found")
	}
	return ws, nil
}

func (m *MockWorkspaceRepository) GetByName(name string) (*domain.Workspace, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, ws := range m.workspaces {
		if ws.Name == name {
			return ws, nil
		}
	}
	return nil, errors.New("workspace not found")
}

func (m *MockWorkspaceRepository) List() ([]*domain.Workspace, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.listErr != nil {
		return nil, m.listErr
	}

	workspaces := make([]*domain.Workspace, 0, len(m.workspaces))
	for _, ws := range m.workspaces {
		workspaces = append(workspaces, ws)
	}

	// Sort by last_used DESC (most recent first), matching the real repository behavior
	// Use a simple bubble sort for small test datasets
	for i := 0; i < len(workspaces); i++ {
		for j := i + 1; j < len(workspaces); j++ {
			// Sort by last_used DESC, then created_at DESC
			if workspaces[j].LastUsed.After(workspaces[i].LastUsed) ||
				(workspaces[j].LastUsed.Equal(workspaces[i].LastUsed) &&
					workspaces[j].CreatedAt.After(workspaces[i].CreatedAt)) {
				workspaces[i], workspaces[j] = workspaces[j], workspaces[i]
			}
		}
	}

	return workspaces, nil
}

func (m *MockWorkspaceRepository) Update(workspace *domain.Workspace) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.updateErr != nil {
		return m.updateErr
	}

	if _, exists := m.workspaces[workspace.ID]; !exists {
		return errors.New("workspace not found")
	}

	m.workspaces[workspace.ID] = workspace
	return nil
}

func (m *MockWorkspaceRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.deleteErr != nil {
		return m.deleteErr
	}

	if _, exists := m.workspaces[id]; !exists {
		return errors.New("workspace not found")
	}

	delete(m.workspaces, id)
	return nil
}

func (m *MockWorkspaceRepository) SetDefault(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clear all defaults first
	for _, ws := range m.workspaces {
		ws.IsDefault = false
	}

	// Set the new default
	ws, exists := m.workspaces[id]
	if !exists {
		return errors.New("workspace not found")
	}
	ws.IsDefault = true
	return nil
}

func (m *MockWorkspaceRepository) GetDefault() (*domain.Workspace, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, ws := range m.workspaces {
		if ws.IsDefault {
			return ws, nil
		}
	}
	return nil, ports.ErrNoDefaultWorkspace
}

func (m *MockWorkspaceRepository) UpdateLastUsed(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ws, exists := m.workspaces[id]
	if !exists {
		return errors.New("workspace not found")
	}
	ws.LastUsed = time.Now()
	return nil
}

// MockCredentialProfileRepository implements the CredentialProfileRepository interface for testing
type MockCredentialProfileRepository struct {
	profiles       map[string]*domain.CredentialProfile
	workspaceUsage map[string][]string // profileID -> workspaceIDs
	mu             sync.RWMutex
	createErr      error
	getErr         error
	listErr        error
	updateErr      error
	deleteErr      error
}

func NewMockCredentialProfileRepository() *MockCredentialProfileRepository {
	return &MockCredentialProfileRepository{
		profiles:       make(map[string]*domain.CredentialProfile),
		workspaceUsage: make(map[string][]string),
	}
}

func (m *MockCredentialProfileRepository) Create(profile *domain.CredentialProfile) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.createErr != nil {
		return m.createErr
	}

	// Check for duplicate name
	for _, existing := range m.profiles {
		if existing.Name == profile.Name {
			return ports.ErrCredentialProfileExists
		}
	}

	m.profiles[profile.ID] = profile
	return nil
}

func (m *MockCredentialProfileRepository) Get(id string) (*domain.CredentialProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.getErr != nil {
		return nil, m.getErr
	}

	profile, exists := m.profiles[id]
	if !exists {
		return nil, ports.ErrCredentialProfileNotFound
	}
	return profile, nil
}

func (m *MockCredentialProfileRepository) GetByName(name string) (*domain.CredentialProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, profile := range m.profiles {
		if profile.Name == name {
			return profile, nil
		}
	}
	return nil, ports.ErrCredentialProfileNotFound
}

func (m *MockCredentialProfileRepository) List() ([]*domain.CredentialProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.listErr != nil {
		return nil, m.listErr
	}

	profiles := make([]*domain.CredentialProfile, 0, len(m.profiles))
	for _, profile := range m.profiles {
		profiles = append(profiles, profile)
	}

	// Sort by name for consistent testing
	for i := 0; i < len(profiles); i++ {
		for j := i + 1; j < len(profiles); j++ {
			if profiles[i].Name > profiles[j].Name {
				profiles[i], profiles[j] = profiles[j], profiles[i]
			}
		}
	}

	return profiles, nil
}

func (m *MockCredentialProfileRepository) Update(profile *domain.CredentialProfile) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.updateErr != nil {
		return m.updateErr
	}

	if _, exists := m.profiles[profile.ID]; !exists {
		return ports.ErrCredentialProfileNotFound
	}

	// Check for duplicate name (excluding current profile)
	for id, existing := range m.profiles {
		if id != profile.ID && existing.Name == profile.Name {
			return ports.ErrCredentialProfileExists
		}
	}

	m.profiles[profile.ID] = profile
	return nil
}

func (m *MockCredentialProfileRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.deleteErr != nil {
		return m.deleteErr
	}

	if _, exists := m.profiles[id]; !exists {
		return ports.ErrCredentialProfileNotFound
	}

	// Check if profile is in use
	if workspaces, exists := m.workspaceUsage[id]; exists && len(workspaces) > 0 {
		return fmt.Errorf("cannot delete credential profile: it is used by %d workspace(s)", len(workspaces))
	}

	delete(m.profiles, id)
	return nil
}

func (m *MockCredentialProfileRepository) GetWorkspacesUsingProfile(profileID string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	workspaceIDs, exists := m.workspaceUsage[profileID]
	if !exists {
		return []string{}, nil
	}

	// Return a copy to avoid concurrent modification
	result := make([]string, len(workspaceIDs))
	copy(result, workspaceIDs)
	return result, nil
}

// SetWorkspaceUsage is a helper method for testing - sets which workspaces use a profile
func (m *MockCredentialProfileRepository) SetWorkspaceUsage(profileID string, workspaceIDs []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.workspaceUsage[profileID] = workspaceIDs
}

// MockCredentialStore implements the CredentialStore interface for testing
type MockCredentialStore struct {
	credentials map[string]domain.WorkspaceConfig
	mu          sync.RWMutex
	storeErr    error
	retrieveErr error
}

func NewMockCredentialStore() *MockCredentialStore {
	return &MockCredentialStore{
		credentials: make(map[string]domain.WorkspaceConfig),
	}
}

func (m *MockCredentialStore) Store(workspaceID string, config domain.WorkspaceConfig) (domain.CredentialRef, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.storeErr != nil {
		return domain.CredentialRef{}, m.storeErr
	}

	ref := domain.CredentialRef{
		KeychainID: "keychain-" + workspaceID,
		ServiceID:  "ticketr",
	}
	m.credentials[ref.KeychainID] = config
	return ref, nil
}

func (m *MockCredentialStore) Retrieve(ref domain.CredentialRef) (*domain.WorkspaceConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.retrieveErr != nil {
		return nil, m.retrieveErr
	}

	config, exists := m.credentials[ref.KeychainID]
	if !exists {
		return nil, errors.New("credentials not found")
	}
	return &config, nil
}

func (m *MockCredentialStore) Delete(ref domain.CredentialRef) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.credentials, ref.KeychainID)
	return nil
}

func (m *MockCredentialStore) List() ([]domain.CredentialRef, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	refs := make([]domain.CredentialRef, 0, len(m.credentials))
	for keychainID := range m.credentials {
		refs = append(refs, domain.CredentialRef{
			KeychainID: keychainID,
			ServiceID:  "ticketr",
		})
	}
	return refs, nil
}

// TestWorkspaceService_Create tests workspace creation
func TestWorkspaceService_Create(t *testing.T) {
	tests := []struct {
		name      string
		wsName    string
		config    domain.WorkspaceConfig
		wantErr   bool
		errMsg    string
		setupMock func(*MockWorkspaceRepository)
	}{
		{
			name:   "successful creation",
			wsName: "backend",
			config: domain.WorkspaceConfig{
				JiraURL:    "https://company.atlassian.net",
				ProjectKey: "BACK",
				Username:   "user@company.com",
				APIToken:   "token123",
			},
			wantErr: false,
		},
		{
			name:   "empty workspace name",
			wsName: "",
			config: domain.WorkspaceConfig{
				JiraURL:    "https://company.atlassian.net",
				ProjectKey: "BACK",
				Username:   "user@company.com",
				APIToken:   "token123",
			},
			wantErr: true,
			errMsg:  "name cannot be empty",
		},
		{
			name:   "duplicate workspace",
			wsName: "backend",
			config: domain.WorkspaceConfig{
				JiraURL:    "https://company.atlassian.net",
				ProjectKey: "BACK",
				Username:   "user@company.com",
				APIToken:   "token123",
			},
			wantErr: true,
			errMsg:  "already exists",
			setupMock: func(m *MockWorkspaceRepository) {
				_ = m.Create(&domain.Workspace{
					ID:         "existing-id",
					Name:       "backend",
					JiraURL:    "https://company.atlassian.net",
					ProjectKey: "BACK",
				})
			},
		},
		{
			name:   "repository error",
			wsName: "backend",
			config: domain.WorkspaceConfig{
				JiraURL:    "https://company.atlassian.net",
				ProjectKey: "BACK",
				Username:   "user@company.com",
				APIToken:   "token123",
			},
			wantErr: true,
			errMsg:  "database error",
			setupMock: func(m *MockWorkspaceRepository) {
				m.createErr = errors.New("database error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockWorkspaceRepository()
			credStore := NewMockCredentialStore()
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			service := NewWorkspaceService(repo, NewMockCredentialProfileRepository(), credStore)
			err := service.Create(tt.wsName, tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Create() error = %v, want error containing %q", err, tt.errMsg)
				}
			}

			if !tt.wantErr {
				// Verify workspace was created
				ws, getErr := repo.GetByName(tt.wsName)
				if getErr != nil {
					t.Errorf("Expected workspace to be created, got error: %v", getErr)
					return
				}
				if ws.Name != tt.wsName {
					t.Errorf("Expected name %q, got %q", tt.wsName, ws.Name)
				}
				if ws.ID == "" {
					t.Error("Expected non-empty workspace ID")
				}
			}
		})
	}
}

// TestWorkspaceService_Switch tests switching between workspaces
func TestWorkspaceService_Switch(t *testing.T) {
	repo := NewMockWorkspaceRepository()
	service := NewWorkspaceService(repo, NewMockCredentialProfileRepository(), NewMockCredentialStore())

	// Create workspaces
	err := service.Create("backend", domain.WorkspaceConfig{
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		Username:   "user@company.com",
		APIToken:   "token123",
	})
	if err != nil {
		t.Fatalf("Failed to create workspace 1: %v", err)
	}

	err = service.Create("frontend", domain.WorkspaceConfig{
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "FRONT",
		Username:   "user@company.com",
		APIToken:   "token123",
	})
	if err != nil {
		t.Fatalf("Failed to create workspace 2: %v", err)
	}

	// Test switching
	err = service.Switch("backend")
	if err != nil {
		t.Errorf("Switch() error = %v", err)
	}

	current, _ := service.Current()
	if current == nil {
		t.Fatal("Expected current workspace to be set")
	}
	if current.Name != "backend" {
		t.Errorf("Expected current workspace to be 'backend', got %q", current.Name)
	}

	// Switch to another workspace
	err = service.Switch("frontend")
	if err != nil {
		t.Errorf("Switch() error = %v", err)
	}

	current, _ = service.Current()
	if current.Name != "frontend" {
		t.Errorf("Expected current workspace to be 'frontend', got %q", current.Name)
	}

	// Verify LastUsed was updated for frontend workspace
	ws2Updated, _ := repo.GetByName("frontend")
	if ws2Updated.LastUsed.IsZero() {
		t.Error("Expected LastUsed to be updated after switch")
	}

	// Test switching to non-existent workspace
	err = service.Switch("nonexistent")
	if err == nil {
		t.Error("Expected error when switching to non-existent workspace")
	}

	// Verify backend workspace wasn't affected
	ws1Check, _ := repo.GetByName("backend")
	if ws1Check == nil {
		t.Error("Workspace 1 should still exist")
	}
}

// TestWorkspaceService_List tests listing workspaces
func TestWorkspaceService_List(t *testing.T) {
	repo := NewMockWorkspaceRepository()
	service := NewWorkspaceService(repo, NewMockCredentialProfileRepository(), NewMockCredentialStore())

	// Initially empty
	workspaces, err := service.List()
	if err != nil {
		t.Errorf("List() error = %v", err)
	}
	if len(workspaces) != 0 {
		t.Errorf("Expected 0 workspaces, got %d", len(workspaces))
	}

	// Create workspaces
	_ = service.Create("backend", domain.WorkspaceConfig{
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		Username:   "user@company.com",
		APIToken:   "token123",
	})
	_ = service.Create("frontend", domain.WorkspaceConfig{
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "FRONT",
		Username:   "user@company.com",
		APIToken:   "token123",
	})
	_ = service.Create("mobile", domain.WorkspaceConfig{
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "MOB",
		Username:   "user@company.com",
		APIToken:   "token123",
	})

	// List all
	workspaces, err = service.List()
	if err != nil {
		t.Errorf("List() error = %v", err)
	}
	if len(workspaces) != 3 {
		t.Errorf("Expected 3 workspaces, got %d", len(workspaces))
	}

	// Verify workspace names
	names := make(map[string]bool)
	for _, ws := range workspaces {
		names[ws.Name] = true
	}

	expected := []string{"backend", "frontend", "mobile"}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("Expected workspace %q in list", name)
		}
	}
}

// TestWorkspaceService_SetDefault tests setting default workspace
func TestWorkspaceService_SetDefault(t *testing.T) {
	repo := NewMockWorkspaceRepository()
	service := NewWorkspaceService(repo, NewMockCredentialProfileRepository(), NewMockCredentialStore())

	// Create workspaces
	err := service.Create("backend", domain.WorkspaceConfig{
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		Username:   "user@company.com",
		APIToken:   "token123",
	})
	if err != nil {
		t.Fatalf("Failed to create backend workspace: %v", err)
	}

	err = service.Create("frontend", domain.WorkspaceConfig{
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "FRONT",
		Username:   "user@company.com",
		APIToken:   "token123",
	})
	if err != nil {
		t.Fatalf("Failed to create frontend workspace: %v", err)
	}

	// Set backend as default
	err = service.SetDefault("backend")
	if err != nil {
		t.Errorf("SetDefault() error = %v", err)
	}

	// Verify backend is default
	ws1, _ := repo.GetByName("backend")
	if !ws1.IsDefault {
		t.Error("Expected backend to be default")
	}

	// Set frontend as default
	err = service.SetDefault("frontend")
	if err != nil {
		t.Errorf("SetDefault() error = %v", err)
	}

	// Verify frontend is default and backend is not
	ws2, _ := repo.GetByName("frontend")
	if !ws2.IsDefault {
		t.Error("Expected frontend to be default")
	}

	ws1, _ = repo.GetByName("backend")
	if ws1.IsDefault {
		t.Error("Expected backend to not be default")
	}

	// Test setting non-existent workspace as default
	err = service.SetDefault("nonexistent")
	if err == nil {
		t.Error("Expected error when setting non-existent workspace as default")
	}
}

// TestWorkspaceService_Delete tests deleting workspaces
func TestWorkspaceService_Delete(t *testing.T) {
	tests := []struct {
		name      string
		wsName    string
		isDefault bool
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "delete non-default workspace",
			wsName:    "backend",
			isDefault: false,
			wantErr:   false,
		},
		{
			name:      "delete default workspace with another workspace available",
			wsName:    "backend",
			isDefault: true,
			wantErr:   false, // Service automatically sets another workspace as default
		},
		{
			name:    "delete non-existent workspace",
			wsName:  "nonexistent",
			wantErr: true,
			errMsg:  "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockWorkspaceRepository()
			service := NewWorkspaceService(repo, NewMockCredentialProfileRepository(), NewMockCredentialStore())

			// Create workspace if not testing non-existent
			if tt.wsName != "nonexistent" {
				// Create the workspace to delete
				err := service.Create(tt.wsName, domain.WorkspaceConfig{
					JiraURL:    "https://company.atlassian.net",
					ProjectKey: "TEST",
					Username:   "user@company.com",
					APIToken:   "token123",
				})
				if err != nil {
					t.Fatalf("Failed to create workspace: %v", err)
				}

				// Create a second workspace so we can delete the first one
				err = service.Create("other", domain.WorkspaceConfig{
					JiraURL:    "https://company.atlassian.net",
					ProjectKey: "OTHER",
					Username:   "user@company.com",
					APIToken:   "token123",
				})
				if err != nil {
					t.Fatalf("Failed to create second workspace: %v", err)
				}

				if tt.isDefault {
					ws, _ := repo.GetByName(tt.wsName)
					_ = repo.SetDefault(ws.ID)
				}
			}

			err := service.Delete(tt.wsName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("Delete() error = %v, want error containing %q", err, tt.errMsg)
				}
			}

			// Verify deletion
			if !tt.wantErr {
				_, err := repo.GetByName(tt.wsName)
				if err == nil {
					t.Error("Expected workspace to be deleted")
				}
			}
		})
	}
}

// TestWorkspaceService_UpdateConfig tests updating workspace configuration
func TestWorkspaceService_UpdateConfig(t *testing.T) {
	repo := NewMockWorkspaceRepository()
	service := NewWorkspaceService(repo, NewMockCredentialProfileRepository(), NewMockCredentialStore())

	// Create workspace
	err := service.Create("backend", domain.WorkspaceConfig{
		JiraURL:    "https://old.atlassian.net",
		ProjectKey: "OLD",
		Username:   "user@company.com",
		APIToken:   "token123",
	})
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	// Update config
	newConfig := domain.WorkspaceConfig{
		JiraURL:    "https://new.atlassian.net",
		ProjectKey: "NEW",
		Username:   "newuser@company.com",
		APIToken:   "newtoken456",
	}

	err = service.UpdateConfig("backend", newConfig)
	if err != nil {
		t.Errorf("UpdateConfig() error = %v", err)
	}

	// Verify update
	updated, _ := repo.GetByName("backend")
	if updated.JiraURL != "https://new.atlassian.net" {
		t.Errorf("Expected JiraURL to be updated, got %q", updated.JiraURL)
	}
	if updated.ProjectKey != "NEW" {
		t.Errorf("Expected ProjectKey to be updated, got %q", updated.ProjectKey)
	}

	// Test updating non-existent workspace
	err = service.UpdateConfig("nonexistent", newConfig)
	if err == nil {
		t.Error("Expected error when updating non-existent workspace")
	}
}

// TestWorkspaceService_ThreadSafety tests concurrent operations
func TestWorkspaceService_ThreadSafety(t *testing.T) {
	repo := NewMockWorkspaceRepository()
	service := NewWorkspaceService(repo, NewMockCredentialProfileRepository(), NewMockCredentialStore())

	// Create initial workspaces
	for i := 0; i < 5; i++ {
		_ = service.Create(fmt.Sprintf("workspace-%d", i), domain.WorkspaceConfig{
			JiraURL:    "https://company.atlassian.net",
			ProjectKey: fmt.Sprintf("WS%d", i),
			Username:   "user@company.com",
			APIToken:   "token123",
		})
	}

	// Concurrent operations
	var wg sync.WaitGroup
	operations := 50

	// Concurrent switches
	wg.Add(operations)
	for i := 0; i < operations; i++ {
		go func(idx int) {
			defer wg.Done()
			_ = service.Switch(fmt.Sprintf("workspace-%d", idx%5))
		}(i)
	}

	// Concurrent lists
	wg.Add(operations)
	for i := 0; i < operations; i++ {
		go func() {
			defer wg.Done()
			_, _ = service.List()
		}()
	}

	// Concurrent creates
	wg.Add(operations)
	for i := 0; i < operations; i++ {
		go func(idx int) {
			defer wg.Done()
			_ = service.Create(fmt.Sprintf("concurrent-%d", idx), domain.WorkspaceConfig{
				JiraURL:    "https://company.atlassian.net",
				ProjectKey: fmt.Sprintf("C%d", idx),
				Username:   "user@company.com",
				APIToken:   "token123",
			})
		}(i)
	}

	wg.Wait()

	// Verify no data corruption
	workspaces, err := service.List()
	if err != nil {
		t.Errorf("List() after concurrent operations error = %v", err)
	}

	if len(workspaces) < 5 {
		t.Errorf("Expected at least 5 workspaces after concurrent operations, got %d", len(workspaces))
	}
}

// TestWorkspaceService_GetCurrent tests getting current workspace
func TestWorkspaceService_GetCurrent(t *testing.T) {
	repo := NewMockWorkspaceRepository()
	service := NewWorkspaceService(repo, NewMockCredentialProfileRepository(), NewMockCredentialStore())

	// No current workspace initially
	current, _ := service.Current()
	if current != nil {
		t.Error("Expected no current workspace initially")
	}

	// Create and switch to workspace
	_ = service.Create("backend", domain.WorkspaceConfig{
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		Username:   "user@company.com",
		APIToken:   "token123",
	})

	err := service.Switch("backend")
	if err != nil {
		t.Fatalf("Switch() error = %v", err)
	}

	current, _ = service.Current()
	if current == nil {
		t.Fatal("Expected current workspace to be set")
	}
	if current.Name != "backend" {
		t.Errorf("Expected current workspace to be 'backend', got %q", current.Name)
	}
}

// TestWorkspaceService_ErrorConditions tests various error scenarios
func TestWorkspaceService_ErrorConditions(t *testing.T) {
	t.Run("repository errors propagate", func(t *testing.T) {
		repo := NewMockWorkspaceRepository()
		repo.listErr = errors.New("database connection failed")
		service := NewWorkspaceService(repo, NewMockCredentialProfileRepository(), NewMockCredentialStore())

		_, err := service.List()
		if err == nil {
			t.Error("Expected error from repository to propagate")
		}
	})

	t.Run("switching to deleted workspace", func(t *testing.T) {
		repo := NewMockWorkspaceRepository()
		service := NewWorkspaceService(repo, NewMockCredentialProfileRepository(), NewMockCredentialStore())

		err := service.Create("backend", domain.WorkspaceConfig{
			JiraURL:    "https://company.atlassian.net",
			ProjectKey: "BACK",
			Username:   "user@company.com",
			APIToken:   "token123",
		})
		if err != nil {
			t.Fatalf("Failed to create workspace: %v", err)
		}

		_ = service.Switch("backend")
		ws, _ := repo.GetByName("backend")
		_ = repo.Delete(ws.ID)

		// Current should still return the cached workspace
		current, _ := service.Current()
		if current == nil {
			t.Error("Expected current workspace to be cached")
		}

		// But switching should fail
		err = service.Switch("backend")
		if err == nil {
			t.Error("Expected error when switching to deleted workspace")
		}
	})
}

// BenchmarkWorkspaceService_Create benchmarks workspace creation
func BenchmarkWorkspaceService_Create(b *testing.B) {
	repo := NewMockWorkspaceRepository()
	service := NewWorkspaceService(repo, NewMockCredentialProfileRepository(), NewMockCredentialStore())

	config := domain.WorkspaceConfig{
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BENCH",
		Username:   "user@company.com",
		APIToken:   "token123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.Create(fmt.Sprintf("workspace-%d", i), config)
	}
}

// BenchmarkWorkspaceService_Switch benchmarks workspace switching
func BenchmarkWorkspaceService_Switch(b *testing.B) {
	repo := NewMockWorkspaceRepository()
	service := NewWorkspaceService(repo, NewMockCredentialProfileRepository(), NewMockCredentialStore())

	// Create workspaces
	for i := 0; i < 10; i++ {
		_ = service.Create(fmt.Sprintf("workspace-%d", i), domain.WorkspaceConfig{
			JiraURL:    "https://company.atlassian.net",
			ProjectKey: fmt.Sprintf("WS%d", i),
			Username:   "user@company.com",
			APIToken:   "token123",
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.Switch(fmt.Sprintf("workspace-%d", i%10))
	}
}

// BenchmarkWorkspaceService_List benchmarks listing workspaces
func BenchmarkWorkspaceService_List(b *testing.B) {
	repo := NewMockWorkspaceRepository()
	service := NewWorkspaceService(repo, NewMockCredentialProfileRepository(), NewMockCredentialStore())

	// Create 100 workspaces
	for i := 0; i < 100; i++ {
		_ = service.Create(fmt.Sprintf("workspace-%d", i), domain.WorkspaceConfig{
			JiraURL:    "https://company.atlassian.net",
			ProjectKey: fmt.Sprintf("WS%d", i),
			Username:   "user@company.com",
			APIToken:   "token123",
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.List()
	}
}

// TestWorkspaceSwitchPersistence tests that workspace switching persists across service instances
func TestWorkspaceSwitchPersistence(t *testing.T) {
	// Shared repository and credential store (simulates database persistence)
	repo := NewMockWorkspaceRepository()
	credStore := NewMockCredentialStore()

	// Create service instance 1
	service1 := NewWorkspaceService(repo, NewMockCredentialProfileRepository(), credStore)

	// Create two workspaces
	err := service1.Create("workspace-a", domain.WorkspaceConfig{
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "WA",
		Username:   "user@company.com",
		APIToken:   "token123",
	})
	if err != nil {
		t.Fatalf("Failed to create workspace-a: %v", err)
	}

	// Wait a bit to ensure different timestamps
	time.Sleep(10 * time.Millisecond)

	err = service1.Create("workspace-b", domain.WorkspaceConfig{
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "WB",
		Username:   "user@company.com",
		APIToken:   "token123",
	})
	if err != nil {
		t.Fatalf("Failed to create workspace-b: %v", err)
	}

	// Switch to workspace-b
	err = service1.Switch("workspace-b")
	if err != nil {
		t.Fatalf("Failed to switch to workspace-b: %v", err)
	}

	// Verify current workspace in service1
	current, err := service1.Current()
	if err != nil {
		t.Fatalf("Failed to get current workspace: %v", err)
	}
	if current.Name != "workspace-b" {
		t.Errorf("Expected current workspace to be 'workspace-b', got %q", current.Name)
	}

	// Create a new service instance (simulates new command invocation)
	service2 := NewWorkspaceService(repo, NewMockCredentialProfileRepository(), credStore)

	// Verify that Current() returns workspace-b (most recently used)
	current2, err := service2.Current()
	if err != nil {
		t.Fatalf("Failed to get current workspace in service2: %v", err)
	}
	if current2.Name != "workspace-b" {
		t.Errorf("Expected current workspace to persist as 'workspace-b', got %q", current2.Name)
	}

	// Switch to workspace-a in service2
	err = service2.Switch("workspace-a")
	if err != nil {
		t.Fatalf("Failed to switch to workspace-a in service2: %v", err)
	}

	// Create service instance 3 (another new command)
	service3 := NewWorkspaceService(repo, NewMockCredentialProfileRepository(), credStore)

	// Verify that Current() now returns workspace-a
	current3, err := service3.Current()
	if err != nil {
		t.Fatalf("Failed to get current workspace in service3: %v", err)
	}
	if current3.Name != "workspace-a" {
		t.Errorf("Expected current workspace to persist as 'workspace-a', got %q", current3.Name)
	}

	// Verify that the most recently used workspace has the latest timestamp
	wsA, _ := repo.GetByName("workspace-a")
	wsB, _ := repo.GetByName("workspace-b")
	if wsA.LastUsed.Before(wsB.LastUsed) {
		t.Error("Expected workspace-a to have a more recent LastUsed timestamp than workspace-b")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) > 0 && len(s) > 0 && s != substr && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestWorkspaceService_CreateProfile tests creating credential profiles
func TestWorkspaceService_CreateProfile(t *testing.T) {
	tests := []struct {
		name       string
		input      domain.CredentialProfileInput
		wantErr    bool
		errMsg     string
		setupMocks func(*MockCredentialProfileRepository, *MockCredentialStore)
	}{
		{
			name: "successful creation",
			input: domain.CredentialProfileInput{
				Name:     "Production Profile",
				JiraURL:  "https://company.atlassian.net",
				Username: "user@company.com",
				APIToken: "token123",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			input: domain.CredentialProfileInput{
				JiraURL:  "https://company.atlassian.net",
				Username: "user@company.com",
				APIToken: "token123",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "duplicate profile name",
			input: domain.CredentialProfileInput{
				Name:     "Existing Profile",
				JiraURL:  "https://company.atlassian.net",
				Username: "user@company.com",
				APIToken: "token123",
			},
			wantErr: true,
			errMsg:  "already exists",
			setupMocks: func(credRepo *MockCredentialProfileRepository, credStore *MockCredentialStore) {
				profile := &domain.CredentialProfile{
					ID:       "existing-id",
					Name:     "Existing Profile",
					JiraURL:  "https://company.atlassian.net",
					Username: "user@company.com",
				}
				credRepo.Create(profile)
			},
		},
		{
			name: "credential store error",
			input: domain.CredentialProfileInput{
				Name:     "Test Profile",
				JiraURL:  "https://company.atlassian.net",
				Username: "user@company.com",
				APIToken: "token123",
			},
			wantErr: true,
			errMsg:  "failed to store credentials",
			setupMocks: func(credRepo *MockCredentialProfileRepository, credStore *MockCredentialStore) {
				credStore.storeErr = errors.New("keychain error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockWorkspaceRepository()
			credRepo := NewMockCredentialProfileRepository()
			credStore := NewMockCredentialStore()

			if tt.setupMocks != nil {
				tt.setupMocks(credRepo, credStore)
			}

			service := NewWorkspaceService(repo, credRepo, credStore)
			profileID, err := service.CreateProfile(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("CreateProfile() error = %v, want error containing %q", err, tt.errMsg)
				}
			}

			if !tt.wantErr {
				if profileID == "" {
					t.Error("Expected non-empty profile ID")
				}

				// Verify profile was created
				profile, getErr := credRepo.GetByName(tt.input.Name)
				if getErr != nil {
					t.Errorf("Expected profile to be created, got error: %v", getErr)
					return
				}
				if profile.Name != tt.input.Name {
					t.Errorf("Expected name %q, got %q", tt.input.Name, profile.Name)
				}
				if profile.JiraURL != tt.input.JiraURL {
					t.Errorf("Expected JiraURL %q, got %q", tt.input.JiraURL, profile.JiraURL)
				}
			}
		})
	}
}

// TestWorkspaceService_ListProfiles tests listing credential profiles
func TestWorkspaceService_ListProfiles(t *testing.T) {
	repo := NewMockWorkspaceRepository()
	credRepo := NewMockCredentialProfileRepository()
	credStore := NewMockCredentialStore()
	service := NewWorkspaceService(repo, credRepo, credStore)

	// Initially empty
	profiles, err := service.ListProfiles()
	if err != nil {
		t.Errorf("ListProfiles() error = %v", err)
	}
	if len(profiles) != 0 {
		t.Errorf("Expected 0 profiles, got %d", len(profiles))
	}

	// Create profiles
	inputs := []domain.CredentialProfileInput{
		{Name: "Beta Profile", JiraURL: "https://beta.atlassian.net", Username: "beta@company.com", APIToken: "token2"},
		{Name: "Alpha Profile", JiraURL: "https://alpha.atlassian.net", Username: "alpha@company.com", APIToken: "token1"},
		{Name: "Gamma Profile", JiraURL: "https://gamma.atlassian.net", Username: "gamma@company.com", APIToken: "token3"},
	}

	for _, input := range inputs {
		_, err := service.CreateProfile(input)
		if err != nil {
			t.Fatalf("Failed to create profile %s: %v", input.Name, err)
		}
	}

	// List all profiles
	profiles, err = service.ListProfiles()
	if err != nil {
		t.Errorf("ListProfiles() error = %v", err)
	}
	if len(profiles) != 3 {
		t.Errorf("Expected 3 profiles, got %d", len(profiles))
	}

	// Verify profiles are ordered by name
	expectedOrder := []string{"Alpha Profile", "Beta Profile", "Gamma Profile"}
	for i, profile := range profiles {
		if profile.Name != expectedOrder[i] {
			t.Errorf("Expected profile at index %d to be %q, got %q", i, expectedOrder[i], profile.Name)
		}
	}
}

// TestWorkspaceService_CreateWithProfile tests creating workspaces with credential profiles
func TestWorkspaceService_CreateWithProfile(t *testing.T) {
	repo := NewMockWorkspaceRepository()
	credRepo := NewMockCredentialProfileRepository()
	credStore := NewMockCredentialStore()
	service := NewWorkspaceService(repo, credRepo, credStore)

	// Create a credential profile first
	profileID, err := service.CreateProfile(domain.CredentialProfileInput{
		Name:     "Production Profile",
		JiraURL:  "https://company.atlassian.net",
		Username: "user@company.com",
		APIToken: "token123",
	})
	if err != nil {
		t.Fatalf("Failed to create credential profile: %v", err)
	}

	tests := []struct {
		name       string
		wsName     string
		projectKey string
		profileID  string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "successful creation",
			wsName:     "backend",
			projectKey: "BACK",
			profileID:  profileID,
			wantErr:    false,
		},
		{
			name:       "empty workspace name",
			wsName:     "",
			projectKey: "BACK",
			profileID:  profileID,
			wantErr:    true,
			errMsg:     "name cannot be empty",
		},
		{
			name:       "empty project key",
			wsName:     "backend2",
			projectKey: "",
			profileID:  profileID,
			wantErr:    true,
			errMsg:     "project key is required",
		},
		{
			name:       "nonexistent profile",
			wsName:     "backend3",
			projectKey: "BACK",
			profileID:  "nonexistent-profile",
			wantErr:    true,
			errMsg:     "credential profile not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateWithProfile(tt.wsName, tt.projectKey, tt.profileID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateWithProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("CreateWithProfile() error = %v, want error containing %q", err, tt.errMsg)
				}
			}

			if !tt.wantErr {
				// Verify workspace was created
				ws, getErr := repo.GetByName(tt.wsName)
				if getErr != nil {
					t.Errorf("Expected workspace to be created, got error: %v", getErr)
					return
				}
				if ws.Name != tt.wsName {
					t.Errorf("Expected name %q, got %q", tt.wsName, ws.Name)
				}
				if ws.ProjectKey != tt.projectKey {
					t.Errorf("Expected project key %q, got %q", tt.projectKey, ws.ProjectKey)
				}
			}
		})
	}
}

// TestWorkspaceService_GetProfile tests retrieving credential profiles by name
func TestWorkspaceService_GetProfile(t *testing.T) {
	repo := NewMockWorkspaceRepository()
	credRepo := NewMockCredentialProfileRepository()
	credStore := NewMockCredentialStore()
	service := NewWorkspaceService(repo, credRepo, credStore)

	// Create a profile
	input := domain.CredentialProfileInput{
		Name:     "Test Profile",
		JiraURL:  "https://company.atlassian.net",
		Username: "user@company.com",
		APIToken: "token123",
	}
	_, err := service.CreateProfile(input)
	if err != nil {
		t.Fatalf("Failed to create profile: %v", err)
	}

	// Test getting existing profile
	profile, err := service.GetProfile("Test Profile")
	if err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}
	if profile.Name != input.Name {
		t.Errorf("Expected name %q, got %q", input.Name, profile.Name)
	}

	// Test getting non-existent profile
	_, err = service.GetProfile("Nonexistent Profile")
	if err == nil {
		t.Error("Expected error for non-existent profile")
	}
	if !contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
}

// TestWorkspaceService_DeleteProfile tests deleting credential profiles
func TestWorkspaceService_DeleteProfile(t *testing.T) {
	repo := NewMockWorkspaceRepository()
	credRepo := NewMockCredentialProfileRepository()
	credStore := NewMockCredentialStore()
	service := NewWorkspaceService(repo, credRepo, credStore)

	// Create a profile
	profileID, err := service.CreateProfile(domain.CredentialProfileInput{
		Name:     "Test Profile",
		JiraURL:  "https://company.atlassian.net",
		Username: "user@company.com",
		APIToken: "token123",
	})
	if err != nil {
		t.Fatalf("Failed to create profile: %v", err)
	}

	t.Run("successful deletion", func(t *testing.T) {
		err := service.DeleteProfile("Test Profile")
		if err != nil {
			t.Errorf("DeleteProfile() error = %v", err)
		}

		// Verify profile was deleted
		_, err = credRepo.Get(profileID)
		if err == nil {
			t.Error("Expected profile to be deleted")
		}
	})

	t.Run("delete profile in use", func(t *testing.T) {
		// Create a new profile for this test
		profileID, err := service.CreateProfile(domain.CredentialProfileInput{
			Name:     "In Use Profile",
			JiraURL:  "https://company.atlassian.net",
			Username: "user@company.com",
			APIToken: "token123",
		})
		if err != nil {
			t.Fatalf("Failed to create profile: %v", err)
		}

		// Set up mock to indicate profile is in use
		credRepo.SetWorkspaceUsage(profileID, []string{"workspace-1", "workspace-2"})

		err = service.DeleteProfile("In Use Profile")
		if err == nil {
			t.Error("Expected error when deleting profile in use")
		}
		if !contains(err.Error(), "used by") {
			t.Errorf("Expected 'used by' error, got: %v", err)
		}
	})

	t.Run("delete non-existent profile", func(t *testing.T) {
		err := service.DeleteProfile("Nonexistent Profile")
		if err == nil {
			t.Error("Expected error for non-existent profile")
		}
		if !contains(err.Error(), "not found") {
			t.Errorf("Expected 'not found' error, got: %v", err)
		}
	})
}
