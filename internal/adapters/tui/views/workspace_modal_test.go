package views

import (
	"fmt"
	"testing"
	"time"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/rivo/tview"
)

// mockWorkspaceRepository implements ports.WorkspaceRepository for testing.
type mockWorkspaceRepository struct {
	workspaces map[string]*domain.Workspace
	nameIndex  map[string]string // name -> id mapping
}

func newMockWorkspaceRepository() *mockWorkspaceRepository {
	return &mockWorkspaceRepository{
		workspaces: make(map[string]*domain.Workspace),
		nameIndex:  make(map[string]string),
	}
}

func (m *mockWorkspaceRepository) Create(workspace *domain.Workspace) error {
	if _, exists := m.nameIndex[workspace.Name]; exists {
		return ports.ErrWorkspaceExists
	}
	m.workspaces[workspace.ID] = workspace
	m.nameIndex[workspace.Name] = workspace.ID
	return nil
}

func (m *mockWorkspaceRepository) GetByName(name string) (*domain.Workspace, error) {
	id, exists := m.nameIndex[name]
	if !exists {
		return nil, ports.ErrWorkspaceNotFound
	}
	return m.workspaces[id], nil
}

func (m *mockWorkspaceRepository) List() ([]*domain.Workspace, error) {
	var workspaces []*domain.Workspace
	for _, ws := range m.workspaces {
		workspaces = append(workspaces, ws)
	}
	return workspaces, nil
}

func (m *mockWorkspaceRepository) Get(id string) (*domain.Workspace, error) {
	ws, exists := m.workspaces[id]
	if !exists {
		return nil, ports.ErrWorkspaceNotFound
	}
	return ws, nil
}

func (m *mockWorkspaceRepository) Update(workspace *domain.Workspace) error {
	m.workspaces[workspace.ID] = workspace
	return nil
}

func (m *mockWorkspaceRepository) Delete(id string) error {
	ws, exists := m.workspaces[id]
	if !exists {
		return ports.ErrWorkspaceNotFound
	}
	delete(m.workspaces, id)
	delete(m.nameIndex, ws.Name)
	return nil
}

func (m *mockWorkspaceRepository) SetDefault(id string) error {
	// Reset all workspaces to not default
	for _, ws := range m.workspaces {
		ws.IsDefault = false
	}
	// Set the specified workspace as default
	if ws, exists := m.workspaces[id]; exists {
		ws.IsDefault = true
	}
	return nil
}

func (m *mockWorkspaceRepository) UpdateLastUsed(id string) error {
	if ws, exists := m.workspaces[id]; exists {
		ws.LastUsed = time.Now()
	}
	return nil
}

func (m *mockWorkspaceRepository) GetDefault() (*domain.Workspace, error) {
	for _, ws := range m.workspaces {
		if ws.IsDefault {
			return ws, nil
		}
	}
	return nil, ports.ErrNoDefaultWorkspace
}

// mockCredentialProfileRepository implements ports.CredentialProfileRepository for testing.
type mockCredentialProfileRepository struct {
	profiles  map[string]*domain.CredentialProfile
	nameIndex map[string]string // name -> id mapping
}

func newMockCredentialProfileRepository() *mockCredentialProfileRepository {
	return &mockCredentialProfileRepository{
		profiles:  make(map[string]*domain.CredentialProfile),
		nameIndex: make(map[string]string),
	}
}

func (m *mockCredentialProfileRepository) Create(profile *domain.CredentialProfile) error {
	if _, exists := m.nameIndex[profile.Name]; exists {
		return ports.ErrCredentialProfileExists
	}
	m.profiles[profile.ID] = profile
	m.nameIndex[profile.Name] = profile.ID
	return nil
}

func (m *mockCredentialProfileRepository) Get(id string) (*domain.CredentialProfile, error) {
	profile, exists := m.profiles[id]
	if !exists {
		return nil, ports.ErrCredentialProfileNotFound
	}
	return profile, nil
}

func (m *mockCredentialProfileRepository) GetByName(name string) (*domain.CredentialProfile, error) {
	id, exists := m.nameIndex[name]
	if !exists {
		return nil, ports.ErrCredentialProfileNotFound
	}
	return m.profiles[id], nil
}

func (m *mockCredentialProfileRepository) List() ([]*domain.CredentialProfile, error) {
	var profiles []*domain.CredentialProfile
	for _, profile := range m.profiles {
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func (m *mockCredentialProfileRepository) Update(profile *domain.CredentialProfile) error {
	m.profiles[profile.ID] = profile
	return nil
}

func (m *mockCredentialProfileRepository) Delete(id string) error {
	profile, exists := m.profiles[id]
	if !exists {
		return ports.ErrCredentialProfileNotFound
	}
	delete(m.profiles, id)
	delete(m.nameIndex, profile.Name)
	return nil
}

func (m *mockCredentialProfileRepository) GetWorkspacesUsingProfile(profileID string) ([]string, error) {
	// For testing, return empty slice (no workspaces using profile)
	return []string{}, nil
}

// mockCredentialStore implements ports.CredentialStore for testing.
type mockCredentialStore struct {
	credentials map[string]domain.WorkspaceConfig
}

func newMockCredentialStore() *mockCredentialStore {
	return &mockCredentialStore{
		credentials: make(map[string]domain.WorkspaceConfig),
	}
}

func (m *mockCredentialStore) Store(id string, config domain.WorkspaceConfig) (domain.CredentialRef, error) {
	ref := domain.CredentialRef{
		KeychainID: "cred_" + id,
		ServiceID:  "ticketr",
	}
	m.credentials[ref.KeychainID] = config
	return ref, nil
}

func (m *mockCredentialStore) Retrieve(ref domain.CredentialRef) (*domain.WorkspaceConfig, error) {
	config, exists := m.credentials[ref.KeychainID]
	if !exists {
		return nil, fmt.Errorf("credential not found")
	}
	return &config, nil
}

func (m *mockCredentialStore) Delete(ref domain.CredentialRef) error {
	delete(m.credentials, ref.KeychainID)
	return nil
}

func (m *mockCredentialStore) List() ([]domain.CredentialRef, error) {
	var refs []domain.CredentialRef
	for id := range m.credentials {
		refs = append(refs, domain.CredentialRef{
			KeychainID: id,
			ServiceID:  "ticketr",
		})
	}
	return refs, nil
}

// setupTestWorkspaceService creates a workspace service with mock dependencies.
func setupTestWorkspaceService() *services.WorkspaceService {
	workspaceRepo := newMockWorkspaceRepository()
	credentialRepo := newMockCredentialProfileRepository()
	credStore := newMockCredentialStore()

	return services.NewWorkspaceService(workspaceRepo, credentialRepo, credStore)
}

func TestWorkspaceModal_Creation(t *testing.T) {
	app := tview.NewApplication()
	workspaceService := setupTestWorkspaceService()

	modal := NewWorkspaceModal(app, workspaceService)

	if modal == nil {
		t.Fatal("Expected modal to be created, got nil")
	}

	if modal.app != app {
		t.Error("Expected modal to store app reference")
	}

	if modal.workspaceService != workspaceService {
		t.Error("Expected modal to store workspace service reference")
	}

	if modal.useExistingProfile != true {
		t.Error("Expected modal to default to using existing profile")
	}
}

func TestWorkspaceModal_LoadProfiles(t *testing.T) {
	app := tview.NewApplication()
	workspaceService := setupTestWorkspaceService()

	// Create test profile
	profileInput := domain.CredentialProfileInput{
		Name:     "test-profile",
		JiraURL:  "https://test.atlassian.net",
		Username: "test@example.com",
		APIToken: "test-token",
	}
	_, err := workspaceService.CreateProfile(profileInput)
	if err != nil {
		t.Fatalf("Failed to create test profile: %v", err)
	}

	modal := NewWorkspaceModal(app, workspaceService)

	// Test loading profiles
	err = modal.loadProfiles()
	if err != nil {
		t.Fatalf("Failed to load profiles: %v", err)
	}

	if len(modal.profiles) != 1 {
		t.Errorf("Expected 1 profile, got %d", len(modal.profiles))
	}

	if modal.profiles[0].Name != "test-profile" {
		t.Errorf("Expected profile name 'test-profile', got '%s'", modal.profiles[0].Name)
	}
}

func TestWorkspaceModal_FormValidation(t *testing.T) {
	app := tview.NewApplication()
	workspaceService := setupTestWorkspaceService()
	modal := NewWorkspaceModal(app, workspaceService)

	tests := []struct {
		name           string
		workspaceName  string
		projectKey     string
		profileMode    bool // true for existing profile, false for new profile
		newProfileData map[string]string
		wantErr        bool
		errContains    string
	}{
		{
			name:          "valid workspace with existing profile",
			workspaceName: "test-workspace",
			projectKey:    "TEST",
			profileMode:   true,
			wantErr:       true, // No profiles available in empty repo
			errContains:   "no credential profiles available",
		},
		{
			name:          "empty workspace name",
			workspaceName: "",
			projectKey:    "TEST",
			profileMode:   true,
			wantErr:       true,
			errContains:   "workspace name",
		},
		{
			name:          "empty project key",
			workspaceName: "test-workspace",
			projectKey:    "",
			profileMode:   true,
			wantErr:       true,
			errContains:   "project key is required",
		},
		{
			name:          "long project key",
			workspaceName: "test-workspace",
			projectKey:    "VERYLONGPROJECTKEY",
			profileMode:   true,
			wantErr:       true,
			errContains:   "project key must be 10 characters or less",
		},
		{
			name:          "valid new profile",
			workspaceName: "test-workspace",
			projectKey:    "TEST",
			profileMode:   false,
			newProfileData: map[string]string{
				"name":     "new-profile",
				"url":      "https://test.atlassian.net",
				"username": "test@example.com",
				"token":    "test-token",
			},
			wantErr: false,
		},
		{
			name:          "new profile missing name",
			workspaceName: "test-workspace",
			projectKey:    "TEST",
			profileMode:   false,
			newProfileData: map[string]string{
				"name":     "",
				"url":      "https://test.atlassian.net",
				"username": "test@example.com",
				"token":    "test-token",
			},
			wantErr:     true,
			errContains: "profile name is required",
		},
		{
			name:          "new profile invalid URL",
			workspaceName: "test-workspace",
			projectKey:    "TEST",
			profileMode:   false,
			newProfileData: map[string]string{
				"name":     "new-profile",
				"url":      "invalid-url",
				"username": "test@example.com",
				"token":    "test-token",
			},
			wantErr:     true,
			errContains: "must start with http",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset modal state
			modal.useExistingProfile = tt.profileMode
			modal.showingNewProfile = !tt.profileMode

			// Set form values
			modal.nameField.SetText(tt.workspaceName)
			modal.projectKeyField.SetText(tt.projectKey)

			if !tt.profileMode && tt.newProfileData != nil {
				modal.newProfileName.SetText(tt.newProfileData["name"])
				modal.newProfileURL.SetText(tt.newProfileData["url"])
				modal.newProfileUsername.SetText(tt.newProfileData["username"])
				modal.newProfileToken.SetText(tt.newProfileData["token"])
			}

			err := modal.validateForm()

			if tt.wantErr {
				if err == nil {
					t.Error("Expected validation error, got nil")
				} else if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no validation error, got: %v", err)
				}
			}
		})
	}
}

func TestWorkspaceModal_CallbacksInvoked(t *testing.T) {
	app := tview.NewApplication()
	workspaceService := setupTestWorkspaceService()
	modal := NewWorkspaceModal(app, workspaceService)

	var closeCalled bool
	var successCalled bool

	modal.SetOnClose(func() {
		closeCalled = true
	})

	modal.SetOnSuccess(func() {
		successCalled = true
	})

	// Test close callback
	modal.handleCancel()
	if !closeCalled {
		t.Error("Expected close callback to be called")
	}

	// Test success callback (showSuccess shows a button that triggers callbacks)
	// We can't easily test the button click in unit tests without complex UI simulation
	// Instead, we'll test the callback setting/getting functionality
	if modal.onSuccess == nil {
		t.Error("Expected success callback to be set")
	}

	if modal.onClose == nil {
		t.Error("Expected close callback to be set")
	}

	// Reset for next test
	closeCalled = false
	successCalled = false

	// Manually trigger callbacks to test they work
	if modal.onSuccess != nil {
		modal.onSuccess()
	}
	if modal.onClose != nil {
		modal.onClose()
	}

	if !successCalled {
		t.Error("Expected success callback to work when called directly")
	}

	if !closeCalled {
		t.Error("Expected close callback to work when called directly")
	}
}

func TestWorkspaceModal_StateManagement(t *testing.T) {
	app := tview.NewApplication()
	workspaceService := setupTestWorkspaceService()
	modal := NewWorkspaceModal(app, workspaceService)

	// Test initial state
	if !modal.useExistingProfile {
		t.Error("Expected modal to start in existing profile mode")
	}

	if modal.showingNewProfile {
		t.Error("Expected modal to not show new profile fields initially")
	}

	if modal.isValidating {
		t.Error("Expected modal to not be in validating state initially")
	}

	// Test OnShow resets state
	modal.useExistingProfile = false
	modal.showingNewProfile = true
	modal.isValidating = true

	modal.OnShow()

	if !modal.useExistingProfile {
		t.Error("Expected OnShow to reset to existing profile mode")
	}

	if modal.showingNewProfile {
		t.Error("Expected OnShow to hide new profile fields")
	}

	if modal.isValidating {
		t.Error("Expected OnShow to reset validating state")
	}
}

func TestWorkspaceModal_FieldClearing(t *testing.T) {
	app := tview.NewApplication()
	workspaceService := setupTestWorkspaceService()
	modal := NewWorkspaceModal(app, workspaceService)

	// Set some field values
	modal.nameField.SetText("test-workspace")
	modal.projectKeyField.SetText("TEST")
	modal.newProfileName.SetText("test-profile")
	modal.newProfileURL.SetText("https://test.atlassian.net")
	modal.newProfileUsername.SetText("test@example.com")
	modal.newProfileToken.SetText("test-token")

	// Call OnShow which should clear fields
	modal.OnShow()

	if modal.nameField.GetText() != "" {
		t.Error("Expected name field to be cleared")
	}

	if modal.projectKeyField.GetText() != "" {
		t.Error("Expected project key field to be cleared")
	}

	if modal.newProfileName.GetText() != "" {
		t.Error("Expected new profile name field to be cleared")
	}

	if modal.newProfileURL.GetText() != "" {
		t.Error("Expected new profile URL field to be cleared")
	}

	if modal.newProfileUsername.GetText() != "" {
		t.Error("Expected new profile username field to be cleared")
	}

	if modal.newProfileToken.GetText() != "" {
		t.Error("Expected new profile token field to be cleared")
	}
}

func TestWorkspaceModal_OnHide(t *testing.T) {
	app := tview.NewApplication()
	workspaceService := setupTestWorkspaceService()
	modal := NewWorkspaceModal(app, workspaceService)

	// Set sensitive data
	modal.newProfileToken.SetText("sensitive-token")

	// Call OnHide
	modal.OnHide()

	// Check that sensitive data is cleared
	if modal.newProfileToken.GetText() != "" {
		t.Error("Expected sensitive token to be cleared on hide")
	}
}

// Helper function to check if a string contains a substring.
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
