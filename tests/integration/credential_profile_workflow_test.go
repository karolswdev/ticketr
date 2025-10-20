//go:build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/karolswdev/ticktr/internal/adapters/database"
	"github.com/karolswdev/ticktr/internal/adapters/keychain"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCredentialProfileWorkflow(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	dbPath := fmt.Sprintf("%s/test.db", tempDir)

	// Initialize database
	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)
	defer db.Close()

	// Apply migrations
	adapter := database.NewSQLiteAdapter(db)
	err = adapter.Initialize()
	require.NoError(t, err)

	// Mock credential store for testing (skip actual keychain)
	credStore := &mockCredentialStore{
		store: make(map[string]domain.Credentials),
	}

	// Initialize services
	workspaceRepo := database.NewWorkspaceRepository(adapter)
	profileRepo := database.NewCredentialProfileRepository(adapter)

	workspaceService := services.NewWorkspaceService(workspaceRepo, profileRepo, credStore)

	t.Run("CreateCredentialProfile", func(t *testing.T) {
		profile := domain.CredentialProfileInput{
			Name:     "test-profile",
			JiraURL:  "https://test.atlassian.net",
			Username: "test@example.com",
			Password: "test-token",
		}

		profileID, err := workspaceService.CreateProfile(profile)
		require.NoError(t, err)
		assert.NotEmpty(t, profileID)

		// Verify profile was stored
		profiles, err := workspaceService.ListProfiles()
		require.NoError(t, err)
		assert.Len(t, profiles, 1)
		assert.Equal(t, "test-profile", profiles[0].Name)
		assert.Equal(t, "https://test.atlassian.net", profiles[0].JiraURL)
		assert.Equal(t, "test@example.com", profiles[0].Username)

		// Verify credentials were stored in keychain
		_, exists := credStore.store[profiles[0].KeychainRef.String()]
		assert.True(t, exists, "Credentials should be stored in keychain")
	})

	t.Run("CreateWorkspaceWithProfile", func(t *testing.T) {
		// First create a profile
		profile := domain.CredentialProfileInput{
			Name:     "company-admin",
			JiraURL:  "https://company.atlassian.net",
			Username: "admin@company.com",
			Password: "admin-token",
		}
		profileID, err := workspaceService.CreateProfile(profile)
		require.NoError(t, err)

		// Create workspace using the profile
		config := domain.WorkspaceConfig{
			Name:                "backend",
			ProjectKey:          "BACK",
			CredentialProfileID: &profileID,
		}

		err = workspaceService.Create(config)
		require.NoError(t, err)

		// Verify workspace was created
		workspaces, err := workspaceService.List()
		require.NoError(t, err)
		assert.Len(t, workspaces, 1)

		workspace := workspaces[0]
		assert.Equal(t, "backend", workspace.Name)
		assert.Equal(t, "BACK", workspace.ProjectKey)
		assert.Equal(t, "https://company.atlassian.net", workspace.JiraURL)
		assert.NotNil(t, workspace.CredentialProfileID)
		assert.Equal(t, profileID, *workspace.CredentialProfileID)
		assert.True(t, workspace.IsDefault, "First workspace should be default")
	})

	t.Run("CreateMultipleWorkspacesFromSameProfile", func(t *testing.T) {
		// Create a profile
		profile := domain.CredentialProfileInput{
			Name:     "multi-project",
			JiraURL:  "https://multi.atlassian.net",
			Username: "multi@company.com",
			Password: "multi-token",
		}
		profileID, err := workspaceService.CreateProfile(profile)
		require.NoError(t, err)

		// Create multiple workspaces
		workspaceConfigs := []domain.WorkspaceConfig{
			{
				Name:                "frontend",
				ProjectKey:          "FRONT",
				CredentialProfileID: &profileID,
			},
			{
				Name:                "mobile",
				ProjectKey:          "MOB",
				CredentialProfileID: &profileID,
			},
			{
				Name:                "devops",
				ProjectKey:          "OPS",
				CredentialProfileID: &profileID,
			},
		}

		for _, config := range workspaceConfigs {
			err := workspaceService.Create(config)
			require.NoError(t, err)
		}

		// Verify all workspaces were created
		workspaces, err := workspaceService.List()
		require.NoError(t, err)

		// Should have 4 total (1 from previous test + 3 new)
		profileWorkspaces := make([]domain.Workspace, 0)
		for _, ws := range workspaces {
			if ws.CredentialProfileID != nil && *ws.CredentialProfileID == profileID {
				profileWorkspaces = append(profileWorkspaces, ws)
			}
		}
		assert.Len(t, profileWorkspaces, 3)

		// Verify all use the same Jira URL from profile
		for _, ws := range profileWorkspaces {
			assert.Equal(t, "https://multi.atlassian.net", ws.JiraURL)
		}
	})

	t.Run("UseWorkspaceCredentials", func(t *testing.T) {
		// Switch to a workspace and verify credentials are loaded
		err := workspaceService.Switch("backend")
		require.NoError(t, err)

		current := workspaceService.Current()
		require.NotNil(t, current)
		assert.Equal(t, "backend", current.Name)

		// Verify credentials can be retrieved
		creds, err := workspaceService.GetCurrentCredentials()
		require.NoError(t, err)
		assert.Equal(t, "admin@company.com", creds.Username)
		assert.Equal(t, "admin-token", creds.Password)
	})

	t.Run("ErrorCases", func(t *testing.T) {
		// Duplicate profile name
		duplicateProfile := domain.CredentialProfileInput{
			Name:     "test-profile", // Already exists
			JiraURL:  "https://duplicate.atlassian.net",
			Username: "duplicate@example.com",
			Password: "duplicate-token",
		}
		_, err := workspaceService.CreateProfile(duplicateProfile)
		assert.Error(t, err, "Should reject duplicate profile name")

		// Workspace with non-existent profile
		invalidConfig := domain.WorkspaceConfig{
			Name:                "invalid",
			ProjectKey:          "INV",
			CredentialProfileID: stringPtr("non-existent-id"),
		}
		err = workspaceService.Create(invalidConfig)
		assert.Error(t, err, "Should reject non-existent profile ID")

		// Empty profile name
		emptyProfile := domain.CredentialProfileInput{
			Name:     "",
			JiraURL:  "https://empty.atlassian.net",
			Username: "empty@example.com",
			Password: "empty-token",
		}
		_, err = workspaceService.CreateProfile(emptyProfile)
		assert.Error(t, err, "Should reject empty profile name")
	})
}

func TestCredentialProfileMigration(t *testing.T) {
	// Test that existing workspaces continue to work when profiles are introduced
	tempDir := t.TempDir()
	dbPath := fmt.Sprintf("%s/test.db", tempDir)

	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)
	defer db.Close()

	adapter := database.NewSQLiteAdapter(db)
	credStore := &mockCredentialStore{store: make(map[string]domain.Credentials)}

	t.Run("LegacyWorkspaceCompatibility", func(t *testing.T) {
		// Apply only initial migrations (without credential profiles)
		err = adapter.ApplyMigration(1)
		require.NoError(t, err)
		err = adapter.ApplyMigration(2)
		require.NoError(t, err)

		// Create workspace using legacy method (direct credentials)
		workspaceRepo := database.NewWorkspaceRepository(adapter)
		profileRepo := database.NewCredentialProfileRepository(adapter)

		workspaceService := services.NewWorkspaceService(workspaceRepo, profileRepo, credStore)

		legacyConfig := domain.WorkspaceConfig{
			Name:       "legacy",
			JiraURL:    "https://legacy.atlassian.net",
			ProjectKey: "LEG",
			Username:   "legacy@example.com",
			Password:   "legacy-token",
		}

		err = workspaceService.Create(legacyConfig)
		require.NoError(t, err)

		// Apply profile migration
		err = adapter.ApplyMigration(3)
		require.NoError(t, err)

		// Verify legacy workspace still works
		workspaces, err := workspaceService.List()
		require.NoError(t, err)
		assert.Len(t, workspaces, 1)

		legacy := workspaces[0]
		assert.Equal(t, "legacy", legacy.Name)
		assert.Nil(t, legacy.CredentialProfileID, "Legacy workspace should not have profile ID")

		// Verify credentials still accessible
		err = workspaceService.Switch("legacy")
		require.NoError(t, err)

		creds, err := workspaceService.GetCurrentCredentials()
		require.NoError(t, err)
		assert.Equal(t, "legacy@example.com", creds.Username)
		assert.Equal(t, "legacy-token", creds.Password)
	})
}

func TestCredentialProfileCleanup(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := fmt.Sprintf("%s/test.db", tempDir)

	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)
	defer db.Close()

	adapter := database.NewSQLiteAdapter(db)
	err = adapter.Initialize()
	require.NoError(t, err)

	credStore := &mockCredentialStore{store: make(map[string]domain.Credentials)}
	workspaceRepo := database.NewWorkspaceRepository(adapter)
	profileRepo := database.NewCredentialProfileRepository(adapter)
	workspaceService := services.NewWorkspaceService(workspaceRepo, profileRepo, credStore)

	t.Run("PreventProfileDeletionWhenInUse", func(t *testing.T) {
		// Create profile
		profile := domain.CredentialProfileInput{
			Name:     "protected-profile",
			JiraURL:  "https://protected.atlassian.net",
			Username: "protected@example.com",
			Password: "protected-token",
		}
		profileID, err := workspaceService.CreateProfile(profile)
		require.NoError(t, err)

		// Create workspace using the profile
		config := domain.WorkspaceConfig{
			Name:                "protected-workspace",
			ProjectKey:          "PROT",
			CredentialProfileID: &profileID,
		}
		err = workspaceService.Create(config)
		require.NoError(t, err)

		// Attempt to delete profile (should fail)
		err = profileRepo.Delete(profileID)
		assert.Error(t, err, "Should not allow deleting profile in use")

		// Delete workspace first
		err = workspaceService.Delete("protected-workspace")
		require.NoError(t, err)

		// Now profile deletion should succeed
		err = profileRepo.Delete(profileID)
		assert.NoError(t, err, "Should allow deleting unused profile")
	})
}

// mockCredentialStore implements ports.CredentialStore for testing
type mockCredentialStore struct {
	store map[string]domain.Credentials
}

func (m *mockCredentialStore) Store(ref domain.CredentialRef, creds domain.Credentials) error {
	m.store[ref.String()] = creds
	return nil
}

func (m *mockCredentialStore) Retrieve(ref domain.CredentialRef) (domain.Credentials, error) {
	creds, exists := m.store[ref.String()]
	if !exists {
		return domain.Credentials{}, fmt.Errorf("credentials not found for ref: %s", ref.String())
	}
	return creds, nil
}

func (m *mockCredentialStore) Delete(ref domain.CredentialRef) error {
	delete(m.store, ref.String())
	return nil
}

func (m *mockCredentialStore) List() ([]domain.CredentialRef, error) {
	refs := make([]domain.CredentialRef, 0, len(m.store))
	for refStr := range m.store {
		ref, _ := domain.NewCredentialRef(refStr)
		refs = append(refs, ref)
	}
	return refs, nil
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}
