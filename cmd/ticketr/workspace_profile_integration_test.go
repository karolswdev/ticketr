//go:build integration

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/karolswdev/ticktr/internal/adapters/database"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/services"
)

func TestWorkspaceCreateWithProfile_Integration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Create temporary directory for test
	tempDir := t.TempDir()

	// Initialize database with temporary file
	dbPath := filepath.Join(tempDir, "test.db")
	adapter, err := database.NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer adapter.Close()

	// Create repositories
	workspaceRepo := database.NewWorkspaceRepository(adapter.DB())
	credentialRepo := database.NewCredentialProfileRepository(adapter.DB())

	// Create mock credential store
	credStore := &MockCredentialStore{
		credentials: make(map[string]domain.WorkspaceConfig),
	}

	// Initialize workspace service
	svc := services.NewWorkspaceService(workspaceRepo, credentialRepo, credStore)

	// Test scenario: Create credential profile and use it for workspace creation
	t.Run("create_profile_and_workspace", func(t *testing.T) {
		// Step 1: Create a credential profile
		profileInput := domain.CredentialProfileInput{
			Name:     "test-profile",
			JiraURL:  "https://test.atlassian.net",
			Username: "test@example.com",
			APIToken: "test-token",
		}

		profileID, err := svc.CreateProfile(profileInput)
		if err != nil {
			t.Fatalf("Failed to create credential profile: %v", err)
		}

		// Step 2: Verify profile was created
		profile, err := svc.GetProfile("test-profile")
		if err != nil {
			t.Fatalf("Failed to get credential profile: %v", err)
		}

		if profile.Name != "test-profile" {
			t.Errorf("Profile name = %v, want %v", profile.Name, "test-profile")
		}

		if profile.JiraURL != "https://test.atlassian.net" {
			t.Errorf("Profile URL = %v, want %v", profile.JiraURL, "https://test.atlassian.net")
		}

		// Step 3: Create workspace using the profile
		err = svc.CreateWithProfile("test-workspace", "TEST", profileID)
		if err != nil {
			t.Fatalf("Failed to create workspace with profile: %v", err)
		}

		// Step 4: Verify workspace was created correctly
		workspaces, err := svc.List()
		if err != nil {
			t.Fatalf("Failed to list workspaces: %v", err)
		}

		if len(workspaces) != 1 {
			t.Fatalf("Expected 1 workspace, got %d", len(workspaces))
		}

		workspace := workspaces[0]
		if workspace.Name != "test-workspace" {
			t.Errorf("Workspace name = %v, want %v", workspace.Name, "test-workspace")
		}

		if workspace.ProjectKey != "TEST" {
			t.Errorf("Workspace project = %v, want %v", workspace.ProjectKey, "TEST")
		}

		if workspace.JiraURL != "https://test.atlassian.net" {
			t.Errorf("Workspace URL = %v, want %v", workspace.JiraURL, "https://test.atlassian.net")
		}

		// Step 5: Verify both workspace and profile share the same credentials
		if workspace.Credentials != profile.KeychainRef {
			t.Errorf("Workspace and profile should share credentials")
		}
	})

	t.Run("workspace_create_with_nonexistent_profile", func(t *testing.T) {
		// Try to create workspace with non-existent profile
		err := svc.CreateWithProfile("test-workspace-2", "TEST2", "non-existent-id")
		if err == nil {
			t.Error("Expected error when using non-existent profile")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("Error should mention profile not found, got: %v", err)
		}
	})

	t.Run("profile_deletion_prevents_when_in_use", func(t *testing.T) {
		// Create profile
		profileInput := domain.CredentialProfileInput{
			Name:     "delete-test-profile",
			JiraURL:  "https://delete.atlassian.net",
			Username: "delete@example.com",
			APIToken: "delete-token",
		}

		profileID, err := svc.CreateProfile(profileInput)
		if err != nil {
			t.Fatalf("Failed to create credential profile: %v", err)
		}

		// Create workspace using the profile
		err = svc.CreateWithProfile("delete-test-workspace", "DEL", profileID)
		if err != nil {
			t.Fatalf("Failed to create workspace with profile: %v", err)
		}

		// Try to delete profile while it's in use
		err = svc.DeleteProfile("delete-test-profile")
		if err == nil {
			t.Error("Expected error when deleting profile in use")
		}

		if !strings.Contains(err.Error(), "used by") {
			t.Errorf("Error should mention profile is in use, got: %v", err)
		}
	})
}

func TestWorkspaceCreateWithProfileCLI_Integration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Create temporary directory for test
	tempDir := t.TempDir()

	// Set up environment to use temp directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	// Create .config directory structure
	configDir := filepath.Join(tempDir, ".config", "ticketr")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Test CLI command integration
	t.Run("cli_profile_creation_flow", func(t *testing.T) {
		// Note: This is a simplified integration test
		// Real CLI testing would require more sophisticated mocking
		// of terminal input and credential stores

		// Test that the CLI command structure is correct
		testArgs := []string{"credentials", "profile", "create", "test-cli-profile",
			"--url", "https://cli.atlassian.net",
			"--username", "cli@example.com"}

		// Verify command parsing (without actual execution)
		cmd := credentialsCmd
		cmd.SetArgs(testArgs[1:]) // Skip "credentials" as it's the root

		// Parse flags
		err := cmd.ParseFlags(testArgs[1:])
		if err != nil {
			t.Errorf("Failed to parse CLI flags: %v", err)
		}

		// Test workspace create with profile flag parsing
		workspaceArgs := []string{"workspace", "create", "test-cli-workspace",
			"--profile", "test-cli-profile",
			"--project", "CLI"}

		workspaceCmd.SetArgs(workspaceArgs[1:])
		err = workspaceCmd.ParseFlags(workspaceArgs[1:])
		if err != nil {
			t.Errorf("Failed to parse workspace CLI flags: %v", err)
		}

		// Verify profile flag was set
		profileFlag := workspaceCreateCmd.Flags().Lookup("profile")
		if profileFlag == nil {
			t.Error("Profile flag not found on workspace create command")
		}
	})
}

// MockCredentialStore for testing
type MockCredentialStore struct {
	credentials map[string]domain.WorkspaceConfig
}

func (m *MockCredentialStore) Store(id string, config domain.WorkspaceConfig) (domain.CredentialRef, error) {
	ref := domain.CredentialRef{
		KeychainID: id,
		ServiceID:  "mock-ticketr",
	}
	m.credentials[id] = config
	return ref, nil
}

func (m *MockCredentialStore) Retrieve(ref domain.CredentialRef) (*domain.WorkspaceConfig, error) {
	config, exists := m.credentials[ref.KeychainID]
	if !exists {
		return nil, os.ErrNotExist
	}
	return &config, nil
}

func (m *MockCredentialStore) Delete(ref domain.CredentialRef) error {
	delete(m.credentials, ref.KeychainID)
	return nil
}

func (m *MockCredentialStore) List() ([]domain.CredentialRef, error) {
	refs := make([]domain.CredentialRef, 0, len(m.credentials))
	for id := range m.credentials {
		refs = append(refs, domain.CredentialRef{
			KeychainID: id,
			ServiceID:  "mock-ticketr",
		})
	}
	return refs, nil
}
