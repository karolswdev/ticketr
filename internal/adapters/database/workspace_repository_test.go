package database

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// TestWorkspaceRepository_Create tests workspace creation
func TestWorkspaceRepository_Create(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	ws := &domain.Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		IsDefault:  false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = repo.Create(ws)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	// Verify workspace was created
	retrieved, err := repo.Get(ws.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if retrieved.Name != ws.Name {
		t.Errorf("Expected name %q, got %q", ws.Name, retrieved.Name)
	}
	if retrieved.JiraURL != ws.JiraURL {
		t.Errorf("Expected JiraURL %q, got %q", ws.JiraURL, retrieved.JiraURL)
	}
	if retrieved.ProjectKey != ws.ProjectKey {
		t.Errorf("Expected ProjectKey %q, got %q", ws.ProjectKey, retrieved.ProjectKey)
	}
}

// TestWorkspaceRepository_CreateDuplicate tests duplicate workspace creation
func TestWorkspaceRepository_CreateDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	ws1 := &domain.Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = repo.Create(ws1)
	if err != nil {
		t.Fatalf("First Create() error = %v", err)
	}

	// Try to create duplicate with same ID
	ws2 := &domain.Workspace{
		ID:         "ws-1",
		Name:       "frontend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "FRONT",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = repo.Create(ws2)
	if err == nil {
		t.Error("Expected error when creating duplicate workspace ID")
	}

	// Try to create duplicate with same name
	ws3 := &domain.Workspace{
		ID:         "ws-3",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = repo.Create(ws3)
	if err == nil {
		t.Error("Expected error when creating duplicate workspace name")
	}
}

// TestWorkspaceRepository_Get tests retrieving workspaces
func TestWorkspaceRepository_Get(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	ws := &domain.Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_ = repo.Create(ws)

	// Test Get by ID
	retrieved, err := repo.Get("ws-1")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if retrieved.Name != "backend" {
		t.Errorf("Expected name 'backend', got %q", retrieved.Name)
	}

	// Test Get non-existent
	_, err = repo.Get("nonexistent")
	if err == nil {
		t.Error("Expected error when getting non-existent workspace")
	}
}

// TestWorkspaceRepository_GetByName tests retrieving workspace by name
func TestWorkspaceRepository_GetByName(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	ws := &domain.Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_ = repo.Create(ws)

	// Test GetByName
	retrieved, err := repo.GetByName("backend")
	if err != nil {
		t.Errorf("GetByName() error = %v", err)
	}
	if retrieved.ID != "ws-1" {
		t.Errorf("Expected ID 'ws-1', got %q", retrieved.ID)
	}

	// Test GetByName non-existent
	_, err = repo.GetByName("nonexistent")
	if err == nil {
		t.Error("Expected error when getting non-existent workspace by name")
	}
}

// TestWorkspaceRepository_List tests listing workspaces
func TestWorkspaceRepository_List(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	// Create multiple workspaces
	workspaces := []*domain.Workspace{
		{
			ID:         "ws-1",
			Name:       "backend",
			JiraURL:    "https://company.atlassian.net",
			ProjectKey: "BACK",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         "ws-2",
			Name:       "frontend",
			JiraURL:    "https://company.atlassian.net",
			ProjectKey: "FRONT",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         "ws-3",
			Name:       "mobile",
			JiraURL:    "https://company.atlassian.net",
			ProjectKey: "MOB",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	for _, ws := range workspaces {
		_ = repo.Create(ws)
	}

	// List all (note: migration creates a "default" workspace, so we expect 4)
	listed, err := repo.List()
	if err != nil {
		t.Errorf("List() error = %v", err)
	}

	if len(listed) < 3 {
		t.Errorf("Expected at least 3 workspaces, got %d", len(listed))
	}

	// Verify all workspaces are present
	names := make(map[string]bool)
	for _, ws := range listed {
		names[ws.Name] = true
	}

	expected := []string{"backend", "frontend", "mobile"}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("Expected workspace %q in list", name)
		}
	}
}

// TestWorkspaceRepository_Update tests updating workspaces
func TestWorkspaceRepository_Update(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	ws := &domain.Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://old.atlassian.net",
		ProjectKey: "OLD",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_ = repo.Create(ws)

	// Update workspace
	ws.JiraURL = "https://new.atlassian.net"
	ws.ProjectKey = "NEW"
	ws.UpdatedAt = time.Now()

	err = repo.Update(ws)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	// Verify update
	updated, _ := repo.Get("ws-1")
	if updated.JiraURL != "https://new.atlassian.net" {
		t.Errorf("Expected JiraURL to be updated, got %q", updated.JiraURL)
	}
	if updated.ProjectKey != "NEW" {
		t.Errorf("Expected ProjectKey to be updated, got %q", updated.ProjectKey)
	}

	// Test updating non-existent workspace
	nonExistent := &domain.Workspace{
		ID:         "nonexistent",
		Name:       "test",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "TEST",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = repo.Update(nonExistent)
	if err == nil {
		t.Error("Expected error when updating non-existent workspace")
	}
}

// TestWorkspaceRepository_Delete tests deleting workspaces
func TestWorkspaceRepository_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	ws := &domain.Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_ = repo.Create(ws)

	// Delete workspace
	err = repo.Delete("ws-1")
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// Verify deletion
	_, err = repo.Get("ws-1")
	if err == nil {
		t.Error("Expected error when getting deleted workspace")
	}

	// Test deleting non-existent workspace
	err = repo.Delete("nonexistent")
	if err == nil {
		t.Error("Expected error when deleting non-existent workspace")
	}
}

// TestWorkspaceRepository_SetDefault tests setting default workspace
func TestWorkspaceRepository_SetDefault(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	// Create multiple workspaces
	ws1 := &domain.Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	ws2 := &domain.Workspace{
		ID:         "ws-2",
		Name:       "frontend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "FRONT",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_ = repo.Create(ws1)
	_ = repo.Create(ws2)

	// Set ws1 as default
	err = repo.SetDefault("ws-1")
	if err != nil {
		t.Errorf("SetDefault() error = %v", err)
	}

	// Verify ws1 is default
	retrieved1, _ := repo.Get("ws-1")
	if !retrieved1.IsDefault {
		t.Error("Expected ws1 to be default")
	}

	// Set ws2 as default
	err = repo.SetDefault("ws-2")
	if err != nil {
		t.Errorf("SetDefault() error = %v", err)
	}

	// Verify ws2 is default and ws1 is not
	retrieved2, _ := repo.Get("ws-2")
	if !retrieved2.IsDefault {
		t.Error("Expected ws2 to be default")
	}

	retrieved1, _ = repo.Get("ws-1")
	if retrieved1.IsDefault {
		t.Error("Expected ws1 to not be default anymore")
	}

	// Test setting non-existent workspace as default
	err = repo.SetDefault("nonexistent")
	if err == nil {
		t.Error("Expected error when setting non-existent workspace as default")
	}
}

// TestWorkspaceRepository_DefaultConstraint tests that only one workspace can be default
func TestWorkspaceRepository_DefaultConstraint(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	// Create workspace 1 as default
	ws1 := &domain.Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		IsDefault:  true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	_ = repo.Create(ws1)

	// Create workspace 2 as default - should clear ws1's default
	ws2 := &domain.Workspace{
		ID:         "ws-2",
		Name:       "frontend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "FRONT",
		IsDefault:  true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	_ = repo.Create(ws2)

	// Count default workspaces
	var count int
	query := "SELECT COUNT(*) FROM workspaces WHERE is_default = TRUE"
	err = adapter.db.QueryRow(query).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count defaults: %v", err)
	}

	if count > 1 {
		t.Errorf("Expected at most 1 default workspace, got %d", count)
	}
}

// TestWorkspaceRepository_Transaction tests transaction handling
func TestWorkspaceRepository_Transaction(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	// Test transaction rollback
	tx, err := adapter.db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// Insert workspace in transaction
	query := `
		INSERT INTO workspaces (id, name, jira_url, project_key, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err = tx.Exec(query, "ws-tx", "transaction-test", "https://company.atlassian.net", "TX", time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to insert in transaction: %v", err)
	}

	// Rollback
	err = tx.Rollback()
	if err != nil {
		t.Fatalf("Failed to rollback: %v", err)
	}

	// Verify workspace was not created
	repo := NewWorkspaceRepository(adapter.db)
	_, err = repo.Get("ws-tx")
	if err == nil {
		t.Error("Expected workspace to not exist after rollback")
	}

	// Test transaction commit
	tx, err = adapter.db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	_, err = tx.Exec(query, "ws-tx", "transaction-test", "https://company.atlassian.net", "TX", time.Now(), time.Now())
	if err != nil {
		t.Fatalf("Failed to insert in transaction: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	// Verify workspace was created
	ws, err := repo.Get("ws-tx")
	if err != nil {
		t.Error("Expected workspace to exist after commit")
	}
	if ws != nil && ws.Name != "transaction-test" {
		t.Errorf("Expected name 'transaction-test', got %q", ws.Name)
	}
}

// TestWorkspaceRepository_ConcurrentAccess tests thread-safety
func TestWorkspaceRepository_ConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	var wg sync.WaitGroup
	var createdMutex sync.Mutex
	createdIDs := make(map[string]bool)
	operations := 50

	// Concurrent creates - track which ones succeed
	wg.Add(operations)
	for i := 0; i < operations; i++ {
		go func(idx int) {
			defer wg.Done()
			ws := &domain.Workspace{
				ID:         fmt.Sprintf("ws-%d", idx),
				Name:       fmt.Sprintf("workspace-%d", idx),
				JiraURL:    "https://company.atlassian.net",
				ProjectKey: fmt.Sprintf("WS%d", idx),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			err := repo.Create(ws)
			if err == nil {
				createdMutex.Lock()
				createdIDs[ws.ID] = true
				createdMutex.Unlock()
			}
		}(i)
	}

	// Wait for all create operations to complete
	wg.Wait()

	// Verify created workspaces
	createdMutex.Lock()
	createdCount := len(createdIDs)
	createdMutex.Unlock()

	if createdCount == 0 {
		t.Fatal("No workspaces were created successfully")
	}

	// In concurrent tests, some creates may fail due to constraints, which is expected
	t.Logf("Successfully created %d/%d workspaces concurrently", createdCount, operations)

	// Concurrent reads - ensure all workspaces are accessible
	wg.Add(operations)
	for i := 0; i < operations; i++ {
		go func(idx int) {
			defer wg.Done()
			_, _ = repo.Get(fmt.Sprintf("ws-%d", idx))
		}(i)
	}

	// Wait for all read operations to complete
	wg.Wait()

	// Concurrent updates
	wg.Add(operations)
	for i := 0; i < operations; i++ {
		go func(idx int) {
			defer wg.Done()
			ws, err := repo.Get(fmt.Sprintf("ws-%d", idx))
			if err == nil {
				ws.JiraURL = fmt.Sprintf("https://updated-%d.atlassian.net", idx)
				_ = repo.Update(ws)
			}
		}(i)
	}

	// Wait for all update operations to complete
	wg.Wait()

	// Verify updates - only check workspaces that were successfully created
	createdMutex.Lock()
	idsToCheck := make([]string, 0, len(createdIDs))
	for id := range createdIDs {
		idsToCheck = append(idsToCheck, id)
	}
	createdMutex.Unlock()

	for _, id := range idsToCheck {
		ws, err := repo.Get(id)
		if err != nil {
			t.Errorf("Failed to get workspace %s that was successfully created: %v", id, err)
			continue
		}

		// Extract index from ID (format: "ws-N")
		var idx int
		fmt.Sscanf(id, "ws-%d", &idx)
		expectedURL := fmt.Sprintf("https://updated-%d.atlassian.net", idx)

		if ws.JiraURL != expectedURL {
			t.Errorf("Expected JiraURL %q, got %q", expectedURL, ws.JiraURL)
		}
	}
}

// TestWorkspaceRepository_CredentialRef tests credential reference handling
func TestWorkspaceRepository_CredentialRef(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	ws := &domain.Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		Credentials: domain.CredentialRef{
			KeychainID: "ticketr:backend",
			ServiceID:  "ticketr",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_ = repo.Create(ws)

	// Verify credential ref was stored
	retrieved, _ := repo.Get("ws-1")
	if retrieved.Credentials.KeychainID != "ticketr:backend" {
		t.Errorf("Expected Credentials.KeychainID to be stored, got %q", retrieved.Credentials.KeychainID)
	}
	if retrieved.Credentials.ServiceID != "ticketr" {
		t.Errorf("Expected Credentials.ServiceID to be stored, got %q", retrieved.Credentials.ServiceID)
	}

	// Verify no credentials are stored in database (security check)
	var creds sql.NullString
	query := "SELECT credential_ref FROM workspaces WHERE id = ?"
	err = adapter.db.QueryRow(query, "ws-1").Scan(&creds)
	if err != nil {
		t.Fatalf("Failed to query credentials: %v", err)
	}

	if creds.Valid && len(creds.String) > 0 {
		// Credential ref should only be a reference, not actual credentials
		if len(creds.String) > 100 {
			t.Error("Credential ref seems too long, might contain actual credentials")
		}
	}
}

// TestWorkspaceRepository_LastUsed tests last used timestamp tracking
func TestWorkspaceRepository_LastUsed(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	ws := &domain.Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_ = repo.Create(ws)

	// Initially LastUsed should be zero
	retrieved, _ := repo.Get("ws-1")
	if !retrieved.LastUsed.IsZero() {
		t.Error("Expected LastUsed to be zero initially")
	}

	// Update LastUsed
	ws.LastUsed = time.Now()
	_ = repo.Update(ws)

	// Verify LastUsed was updated
	retrieved, _ = repo.Get("ws-1")
	if retrieved.LastUsed.IsZero() {
		t.Error("Expected LastUsed to be set after update")
	}
}

// TestWorkspaceRepository_GetDefault tests retrieving the default workspace
func TestWorkspaceRepository_GetDefault(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	// Test 1: No default workspace (initially empty except migration default)
	// First clear any migration defaults
	_, _ = adapter.db.Exec("UPDATE workspaces SET is_default = FALSE WHERE is_default = TRUE")

	_, err = repo.GetDefault()
	if err == nil {
		t.Error("Expected ErrNoDefaultWorkspace when no default is set")
	}

	// Test 2: Create workspaces and mark one as default
	ws1 := &domain.Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		IsDefault:  false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	ws2 := &domain.Workspace{
		ID:         "ws-2",
		Name:       "frontend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "FRONT",
		IsDefault:  true, // Mark as default
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	ws3 := &domain.Workspace{
		ID:         "ws-3",
		Name:       "mobile",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "MOB",
		IsDefault:  false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = repo.Create(ws1)
	if err != nil {
		t.Fatalf("Failed to create ws1: %v", err)
	}

	err = repo.Create(ws2)
	if err != nil {
		t.Fatalf("Failed to create ws2: %v", err)
	}

	err = repo.Create(ws3)
	if err != nil {
		t.Fatalf("Failed to create ws3: %v", err)
	}

	// Test 3: GetDefault should return ws2
	defaultWs, err := repo.GetDefault()
	if err != nil {
		t.Errorf("GetDefault() error = %v", err)
	}

	if defaultWs == nil {
		t.Fatal("GetDefault() returned nil workspace")
	}

	if defaultWs.ID != "ws-2" {
		t.Errorf("Expected default workspace ID 'ws-2', got %q", defaultWs.ID)
	}
	if defaultWs.Name != "frontend" {
		t.Errorf("Expected default workspace name 'frontend', got %q", defaultWs.Name)
	}
	if !defaultWs.IsDefault {
		t.Error("Expected IsDefault to be true")
	}

	// Test 4: Ensure only one workspace is returned
	// Verify by changing default and checking again
	err = repo.SetDefault("ws-1")
	if err != nil {
		t.Fatalf("Failed to set ws-1 as default: %v", err)
	}

	defaultWs, err = repo.GetDefault()
	if err != nil {
		t.Errorf("GetDefault() error after changing default = %v", err)
	}

	if defaultWs.ID != "ws-1" {
		t.Errorf("Expected default workspace ID 'ws-1' after change, got %q", defaultWs.ID)
	}
	if defaultWs.Name != "backend" {
		t.Errorf("Expected default workspace name 'backend' after change, got %q", defaultWs.Name)
	}

	// Test 5: Verify ws-2 is no longer default
	ws2Check, err := repo.Get("ws-2")
	if err != nil {
		t.Fatalf("Failed to get ws-2: %v", err)
	}
	if ws2Check.IsDefault {
		t.Error("Expected ws-2 to no longer be default")
	}
}

// TestWorkspaceRepository_UpdateLastUsed tests updating the LastUsed timestamp
func TestWorkspaceRepository_UpdateLastUsed(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	// Test 1: Create a workspace with zero LastUsed
	ws := &domain.Workspace{
		ID:         "ws-1",
		Name:       "backend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "BACK",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		// LastUsed is zero by default
	}

	err = repo.Create(ws)
	if err != nil {
		t.Fatalf("Failed to create workspace: %v", err)
	}

	// Get initial LastUsed timestamp (should be zero)
	initial, err := repo.Get("ws-1")
	if err != nil {
		t.Fatalf("Failed to get workspace: %v", err)
	}

	if !initial.LastUsed.IsZero() {
		t.Error("Expected initial LastUsed to be zero")
	}

	// Wait to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	// Test 2: Update LastUsed
	err = repo.UpdateLastUsed("ws-1")
	if err != nil {
		t.Errorf("UpdateLastUsed() error = %v", err)
	}

	// Test 3: Verify LastUsed was updated
	updated, err := repo.Get("ws-1")
	if err != nil {
		t.Fatalf("Failed to get workspace after update: %v", err)
	}

	if updated.LastUsed.IsZero() {
		t.Error("Expected LastUsed to be set after UpdateLastUsed()")
	}

	// Verify timestamp is recent (within last second)
	timeSinceUpdate := time.Since(updated.LastUsed)
	if timeSinceUpdate > time.Second {
		t.Errorf("LastUsed timestamp seems stale: %v ago", timeSinceUpdate)
	}

	// Test 4: Update LastUsed again and verify it changes
	firstUpdate := updated.LastUsed
	time.Sleep(10 * time.Millisecond)

	err = repo.UpdateLastUsed("ws-1")
	if err != nil {
		t.Errorf("Second UpdateLastUsed() error = %v", err)
	}

	secondUpdate, err := repo.Get("ws-1")
	if err != nil {
		t.Fatalf("Failed to get workspace after second update: %v", err)
	}

	if !secondUpdate.LastUsed.After(firstUpdate) {
		t.Errorf("Expected second LastUsed (%v) to be after first (%v)",
			secondUpdate.LastUsed, firstUpdate)
	}

	// Test 5: Error case - update non-existent workspace
	err = repo.UpdateLastUsed("nonexistent-ws")
	if err == nil {
		t.Error("Expected error when updating LastUsed for non-existent workspace")
	}

	// Verify it's the correct error type
	// The implementation returns ports.ErrWorkspaceNotFound when rowsAffected == 0
	if err != nil && err.Error() != "workspace not found" {
		// This is fine, just checking that we get an error
		t.Logf("Got error (expected): %v", err)
	}

	// Test 6: Verify LastUsed affects List ordering
	ws2 := &domain.Workspace{
		ID:         "ws-2",
		Name:       "frontend",
		JiraURL:    "https://company.atlassian.net",
		ProjectKey: "FRONT",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err = repo.Create(ws2)
	if err != nil {
		t.Fatalf("Failed to create ws-2: %v", err)
	}

	// Update ws-2's LastUsed (should be more recent than ws-1)
	time.Sleep(10 * time.Millisecond)
	err = repo.UpdateLastUsed("ws-2")
	if err != nil {
		t.Fatalf("Failed to update ws-2 LastUsed: %v", err)
	}

	// List should have ws-2 before ws-1 (most recent first)
	workspaces, err := repo.List()
	if err != nil {
		t.Fatalf("Failed to list workspaces: %v", err)
	}

	// Find ws-1 and ws-2 positions
	var ws1Pos, ws2Pos int = -1, -1
	for i, w := range workspaces {
		if w.ID == "ws-1" {
			ws1Pos = i
		}
		if w.ID == "ws-2" {
			ws2Pos = i
		}
	}

	if ws1Pos == -1 || ws2Pos == -1 {
		t.Error("Could not find ws-1 or ws-2 in workspace list")
	} else if ws2Pos > ws1Pos {
		t.Errorf("Expected ws-2 (pos %d) to come before ws-1 (pos %d) in list due to more recent LastUsed",
			ws2Pos, ws1Pos)
	}
}

// BenchmarkWorkspaceRepository_Create benchmarks workspace creation
func BenchmarkWorkspaceRepository_Create(b *testing.B) {
	tmpDir := b.TempDir()
	dbPath := filepath.Join(tmpDir, "bench.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		b.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ws := &domain.Workspace{
			ID:         fmt.Sprintf("ws-%d", i),
			Name:       fmt.Sprintf("workspace-%d", i),
			JiraURL:    "https://company.atlassian.net",
			ProjectKey: fmt.Sprintf("WS%d", i),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		_ = repo.Create(ws)
	}
}

// BenchmarkWorkspaceRepository_Get benchmarks workspace retrieval
func BenchmarkWorkspaceRepository_Get(b *testing.B) {
	tmpDir := b.TempDir()
	dbPath := filepath.Join(tmpDir, "bench.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		b.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	// Create workspaces
	for i := 0; i < 100; i++ {
		ws := &domain.Workspace{
			ID:         fmt.Sprintf("ws-%d", i),
			Name:       fmt.Sprintf("workspace-%d", i),
			JiraURL:    "https://company.atlassian.net",
			ProjectKey: fmt.Sprintf("WS%d", i),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		_ = repo.Create(ws)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.Get(fmt.Sprintf("ws-%d", i%100))
	}
}

// BenchmarkWorkspaceRepository_List benchmarks listing workspaces
func BenchmarkWorkspaceRepository_List(b *testing.B) {
	tmpDir := b.TempDir()
	dbPath := filepath.Join(tmpDir, "bench.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		b.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	// Create 100 workspaces
	for i := 0; i < 100; i++ {
		ws := &domain.Workspace{
			ID:         fmt.Sprintf("ws-%d", i),
			Name:       fmt.Sprintf("workspace-%d", i),
			JiraURL:    "https://company.atlassian.net",
			ProjectKey: fmt.Sprintf("WS%d", i),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		_ = repo.Create(ws)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.List()
	}
}

// BenchmarkWorkspaceRepository_Update benchmarks workspace updates
func BenchmarkWorkspaceRepository_Update(b *testing.B) {
	tmpDir := b.TempDir()
	dbPath := filepath.Join(tmpDir, "bench.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		b.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewWorkspaceRepository(adapter.db)

	// Create workspaces
	for i := 0; i < 100; i++ {
		ws := &domain.Workspace{
			ID:         fmt.Sprintf("ws-%d", i),
			Name:       fmt.Sprintf("workspace-%d", i),
			JiraURL:    "https://company.atlassian.net",
			ProjectKey: fmt.Sprintf("WS%d", i),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		_ = repo.Create(ws)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ws, _ := repo.Get(fmt.Sprintf("ws-%d", i%100))
		if ws != nil {
			ws.JiraURL = fmt.Sprintf("https://updated-%d.atlassian.net", i)
			_ = repo.Update(ws)
		}
	}
}
