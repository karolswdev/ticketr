package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDBForCredentialProfiles(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create schema_migrations table
	_, err = db.Exec(`
		CREATE TABLE schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create schema_migrations table: %v", err)
	}

	// Create credential_profiles table
	_, err = db.Exec(`
		CREATE TABLE credential_profiles (
			id TEXT PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			jira_url TEXT NOT NULL,
			username TEXT NOT NULL,
			keychain_ref TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create credential_profiles table: %v", err)
	}

	// Create workspaces table for foreign key test
	_, err = db.Exec(`
		CREATE TABLE workspaces (
			id TEXT PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			jira_url TEXT,
			project_key TEXT,
			credential_ref TEXT,
			credential_profile_id TEXT REFERENCES credential_profiles(id),
			is_default BOOLEAN DEFAULT FALSE,
			last_used TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create workspaces table: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

func createTestCredentialProfile(id, name string) *domain.CredentialProfile {
	now := time.Now()
	return &domain.CredentialProfile{
		ID:       id,
		Name:     name,
		JiraURL:  "https://company.atlassian.net",
		Username: "user@company.com",
		KeychainRef: domain.CredentialRef{
			KeychainID: "keychain-" + id,
			ServiceID:  "service-" + id,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestCredentialProfileRepository_Create(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)
	profile := createTestCredentialProfile("test-id", "Test Profile")

	err := repo.Create(profile)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Verify profile was created
	retrieved, err := repo.Get("test-id")
	if err != nil {
		t.Fatalf("Get() failed after Create(): %v", err)
	}

	if retrieved.ID != profile.ID {
		t.Errorf("Create() ID mismatch: got %v, want %v", retrieved.ID, profile.ID)
	}
	if retrieved.Name != profile.Name {
		t.Errorf("Create() Name mismatch: got %v, want %v", retrieved.Name, profile.Name)
	}
	if retrieved.JiraURL != profile.JiraURL {
		t.Errorf("Create() JiraURL mismatch: got %v, want %v", retrieved.JiraURL, profile.JiraURL)
	}
	if retrieved.Username != profile.Username {
		t.Errorf("Create() Username mismatch: got %v, want %v", retrieved.Username, profile.Username)
	}
}

func TestCredentialProfileRepository_Create_Duplicate(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)
	profile := createTestCredentialProfile("test-id", "Test Profile")

	// Create first profile
	err := repo.Create(profile)
	if err != nil {
		t.Fatalf("First Create() failed: %v", err)
	}

	// Try to create duplicate
	duplicate := createTestCredentialProfile("test-id-2", "Test Profile")
	err = repo.Create(duplicate)
	if err != ports.ErrCredentialProfileExists {
		t.Errorf("Create() duplicate should return ErrCredentialProfileExists, got: %v", err)
	}
}

func TestCredentialProfileRepository_Create_Nil(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)

	err := repo.Create(nil)
	if err == nil {
		t.Error("Create() with nil profile should return error")
	}
}

func TestCredentialProfileRepository_Get(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)
	profile := createTestCredentialProfile("test-id", "Test Profile")

	// Create profile first
	err := repo.Create(profile)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Get profile
	retrieved, err := repo.Get("test-id")
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if retrieved.ID != profile.ID {
		t.Errorf("Get() ID mismatch: got %v, want %v", retrieved.ID, profile.ID)
	}
}

func TestCredentialProfileRepository_Get_NotFound(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)

	_, err := repo.Get("nonexistent")
	if err != ports.ErrCredentialProfileNotFound {
		t.Errorf("Get() nonexistent should return ErrCredentialProfileNotFound, got: %v", err)
	}
}

func TestCredentialProfileRepository_GetByName(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)
	profile := createTestCredentialProfile("test-id", "Test Profile")

	// Create profile first
	err := repo.Create(profile)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Get profile by name
	retrieved, err := repo.GetByName("Test Profile")
	if err != nil {
		t.Fatalf("GetByName() failed: %v", err)
	}

	if retrieved.Name != profile.Name {
		t.Errorf("GetByName() Name mismatch: got %v, want %v", retrieved.Name, profile.Name)
	}
}

func TestCredentialProfileRepository_GetByName_NotFound(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)

	_, err := repo.GetByName("Nonexistent Profile")
	if err != ports.ErrCredentialProfileNotFound {
		t.Errorf("GetByName() nonexistent should return ErrCredentialProfileNotFound, got: %v", err)
	}
}

func TestCredentialProfileRepository_List(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)

	// Create multiple profiles
	profile1 := createTestCredentialProfile("id-1", "Alpha Profile")
	profile2 := createTestCredentialProfile("id-2", "Beta Profile")
	profile3 := createTestCredentialProfile("id-3", "Gamma Profile")

	for _, profile := range []*domain.CredentialProfile{profile1, profile2, profile3} {
		err := repo.Create(profile)
		if err != nil {
			t.Fatalf("Create() failed: %v", err)
		}
	}

	// List profiles
	profiles, err := repo.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(profiles) != 3 {
		t.Errorf("List() count mismatch: got %d, want 3", len(profiles))
	}

	// Verify ordering (should be alphabetical by name)
	expectedOrder := []string{"Alpha Profile", "Beta Profile", "Gamma Profile"}
	for i, profile := range profiles {
		if profile.Name != expectedOrder[i] {
			t.Errorf("List() ordering incorrect at index %d: got %v, want %v", i, profile.Name, expectedOrder[i])
		}
	}
}

func TestCredentialProfileRepository_List_Empty(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)

	profiles, err := repo.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(profiles) != 0 {
		t.Errorf("List() empty should return 0 profiles, got %d", len(profiles))
	}
}

func TestCredentialProfileRepository_Update(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)
	profile := createTestCredentialProfile("test-id", "Test Profile")

	// Create profile first
	err := repo.Create(profile)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Update profile
	profile.Name = "Updated Profile"
	profile.JiraURL = "https://updated.atlassian.net"
	profile.Username = "updated@company.com"

	err = repo.Update(profile)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	// Verify update
	retrieved, err := repo.Get("test-id")
	if err != nil {
		t.Fatalf("Get() after Update() failed: %v", err)
	}

	if retrieved.Name != "Updated Profile" {
		t.Errorf("Update() Name not updated: got %v, want %v", retrieved.Name, "Updated Profile")
	}
	if retrieved.JiraURL != "https://updated.atlassian.net" {
		t.Errorf("Update() JiraURL not updated: got %v, want %v", retrieved.JiraURL, "https://updated.atlassian.net")
	}
	if retrieved.Username != "updated@company.com" {
		t.Errorf("Update() Username not updated: got %v, want %v", retrieved.Username, "updated@company.com")
	}
}

func TestCredentialProfileRepository_Update_NotFound(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)
	profile := createTestCredentialProfile("nonexistent", "Test Profile")

	err := repo.Update(profile)
	if err != ports.ErrCredentialProfileNotFound {
		t.Errorf("Update() nonexistent should return ErrCredentialProfileNotFound, got: %v", err)
	}
}

func TestCredentialProfileRepository_Update_Nil(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)

	err := repo.Update(nil)
	if err == nil {
		t.Error("Update() with nil profile should return error")
	}
}

func TestCredentialProfileRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)
	profile := createTestCredentialProfile("test-id", "Test Profile")

	// Create profile first
	err := repo.Create(profile)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Delete profile
	err = repo.Delete("test-id")
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// Verify deletion
	_, err = repo.Get("test-id")
	if err != ports.ErrCredentialProfileNotFound {
		t.Errorf("Get() after Delete() should return ErrCredentialProfileNotFound, got: %v", err)
	}
}

func TestCredentialProfileRepository_Delete_NotFound(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)

	err := repo.Delete("nonexistent")
	if err != ports.ErrCredentialProfileNotFound {
		t.Errorf("Delete() nonexistent should return ErrCredentialProfileNotFound, got: %v", err)
	}
}

func TestCredentialProfileRepository_Delete_InUse(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)
	profile := createTestCredentialProfile("test-id", "Test Profile")

	// Create profile first
	err := repo.Create(profile)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create workspace that uses this profile
	_, err = db.Exec(`
		INSERT INTO workspaces (id, name, credential_profile_id)
		VALUES ('workspace-1', 'Test Workspace', 'test-id')
	`)
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	// Try to delete profile that's in use
	err = repo.Delete("test-id")
	if err == nil {
		t.Error("Delete() profile in use should return error")
	}

	// Verify profile still exists
	_, err = repo.Get("test-id")
	if err != nil {
		t.Errorf("Get() after failed Delete() should succeed, got: %v", err)
	}
}

func TestCredentialProfileRepository_GetWorkspacesUsingProfile(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)
	profile := createTestCredentialProfile("test-id", "Test Profile")

	// Create profile first
	err := repo.Create(profile)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Create workspaces that use this profile
	workspaceIDs := []string{"workspace-1", "workspace-2", "workspace-3"}
	for _, wsID := range workspaceIDs {
		_, err = db.Exec(`
			INSERT INTO workspaces (id, name, credential_profile_id)
			VALUES (?, ?, 'test-id')
		`, wsID, "Workspace "+wsID)
		if err != nil {
			t.Fatalf("Failed to create workspace %s: %v", wsID, err)
		}
	}

	// Create workspace that doesn't use this profile
	_, err = db.Exec(`
		INSERT INTO workspaces (id, name)
		VALUES ('workspace-other', 'Other Workspace')
	`)
	if err != nil {
		t.Fatalf("Failed to create other workspace: %v", err)
	}

	// Get workspaces using profile
	usingWorkspaces, err := repo.GetWorkspacesUsingProfile("test-id")
	if err != nil {
		t.Fatalf("GetWorkspacesUsingProfile() failed: %v", err)
	}

	if len(usingWorkspaces) != 3 {
		t.Errorf("GetWorkspacesUsingProfile() count mismatch: got %d, want 3", len(usingWorkspaces))
	}

	// Verify correct workspace IDs are returned
	foundIDs := make(map[string]bool)
	for _, id := range usingWorkspaces {
		foundIDs[id] = true
	}

	for _, expectedID := range workspaceIDs {
		if !foundIDs[expectedID] {
			t.Errorf("GetWorkspacesUsingProfile() missing workspace ID: %s", expectedID)
		}
	}

	if foundIDs["workspace-other"] {
		t.Error("GetWorkspacesUsingProfile() should not include workspace-other")
	}
}

func TestCredentialProfileRepository_GetWorkspacesUsingProfile_Empty(t *testing.T) {
	db, cleanup := setupTestDBForCredentialProfiles(t)
	defer cleanup()

	repo := NewCredentialProfileRepository(db)
	profile := createTestCredentialProfile("test-id", "Test Profile")

	// Create profile first
	err := repo.Create(profile)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Get workspaces using profile (should be empty)
	usingWorkspaces, err := repo.GetWorkspacesUsingProfile("test-id")
	if err != nil {
		t.Fatalf("GetWorkspacesUsingProfile() failed: %v", err)
	}

	if len(usingWorkspaces) != 0 {
		t.Errorf("GetWorkspacesUsingProfile() empty should return 0 workspaces, got %d", len(usingWorkspaces))
	}
}
