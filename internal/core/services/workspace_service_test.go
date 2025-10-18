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
		name       string
		wsName     string
		config     domain.WorkspaceConfig
		wantErr    bool
		errMsg     string
		setupMock  func(*MockWorkspaceRepository)
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

			service := NewWorkspaceService(repo, credStore)
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
	service := NewWorkspaceService(repo, NewMockCredentialStore())

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
	service := NewWorkspaceService(repo, NewMockCredentialStore())

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
	service := NewWorkspaceService(repo, NewMockCredentialStore())

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
			service := NewWorkspaceService(repo, NewMockCredentialStore())

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
	service := NewWorkspaceService(repo, NewMockCredentialStore())

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
	service := NewWorkspaceService(repo, NewMockCredentialStore())

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
	service := NewWorkspaceService(repo, NewMockCredentialStore())

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
		service := NewWorkspaceService(repo, NewMockCredentialStore())

		_, err := service.List()
		if err == nil {
			t.Error("Expected error from repository to propagate")
		}
	})

	t.Run("switching to deleted workspace", func(t *testing.T) {
		repo := NewMockWorkspaceRepository()
		service := NewWorkspaceService(repo, NewMockCredentialStore())

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
	service := NewWorkspaceService(repo, NewMockCredentialStore())

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
	service := NewWorkspaceService(repo, NewMockCredentialStore())

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
	service := NewWorkspaceService(repo, NewMockCredentialStore())

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
