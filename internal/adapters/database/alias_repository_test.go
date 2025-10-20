package database

import (
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// TestAliasRepository_Create tests alias creation
func TestAliasRepository_Create(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewAliasRepository(adapter.DB())
	// Use "default" workspace which is created by migration
	workspaceID := "default"

	alias := &domain.JQLAlias{
		ID:          "alias-1",
		Name:        "my-bugs",
		JQL:         "assignee = currentUser() AND type = Bug",
		Description: "My bug tickets",
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = repo.Create(alias)
	if err != nil {
		t.Errorf("Create() error = %v", err)
	}

	// Verify alias was created
	retrieved, err := repo.Get(alias.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if retrieved.Name != alias.Name {
		t.Errorf("Expected name %q, got %q", alias.Name, retrieved.Name)
	}
	if retrieved.JQL != alias.JQL {
		t.Errorf("Expected JQL %q, got %q", alias.JQL, retrieved.JQL)
	}
	if retrieved.WorkspaceID != alias.WorkspaceID {
		t.Errorf("Expected WorkspaceID %q, got %q", alias.WorkspaceID, retrieved.WorkspaceID)
	}
}

// TestAliasRepository_CreateDuplicate tests duplicate alias detection
func TestAliasRepository_CreateDuplicate(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewAliasRepository(adapter.DB())
	// Use "default" workspace which is created by migration
	workspaceID := "default"

	alias1 := &domain.JQLAlias{
		ID:          "alias-1",
		Name:        "my-bugs",
		JQL:         "assignee = currentUser() AND type = Bug",
		Description: "My bug tickets",
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = repo.Create(alias1)
	if err != nil {
		t.Fatalf("First Create() error = %v", err)
	}

	// Try to create duplicate with same name in same workspace
	alias2 := &domain.JQLAlias{
		ID:          "alias-2",
		Name:        "my-bugs",
		JQL:         "different JQL",
		Description: "Different description",
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = repo.Create(alias2)
	if err != ports.ErrAliasExists {
		t.Errorf("Expected ErrAliasExists, got %v", err)
	}

	// Create with same name as global should succeed
	alias3 := &domain.JQLAlias{
		ID:          "alias-3",
		Name:        "my-bugs-global",
		JQL:         "assignee = currentUser()",
		Description: "Global alias",
		WorkspaceID: "", // Global
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = repo.Create(alias3)
	if err != nil {
		t.Errorf("Create() global alias error = %v", err)
	}
}

// TestAliasRepository_Get tests retrieving aliases by ID
func TestAliasRepository_Get(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewAliasRepository(adapter.DB())

	alias := &domain.JQLAlias{
		ID:          "alias-1",
		Name:        "test-alias",
		JQL:         "status = Open",
		Description: "Test",
		WorkspaceID: "", // Global alias
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_ = repo.Create(alias)

	// Test Get by ID
	retrieved, err := repo.Get("alias-1")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if retrieved.Name != "test-alias" {
		t.Errorf("Expected name 'test-alias', got %q", retrieved.Name)
	}

	// Test Get non-existent
	_, err = repo.Get("nonexistent")
	if err != ports.ErrAliasNotFound {
		t.Errorf("Expected ErrAliasNotFound, got %v", err)
	}
}

// TestAliasRepository_GetByName tests retrieving aliases by name
func TestAliasRepository_GetByName(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewAliasRepository(adapter.DB())
	// Use "default" workspace which is created by migration
	workspaceID := "default"

	// Create workspace-specific alias
	alias := &domain.JQLAlias{
		ID:          "alias-1",
		Name:        "test-alias",
		JQL:         "status = Open",
		Description: "Test",
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_ = repo.Create(alias)

	// Test GetByName
	retrieved, err := repo.GetByName("test-alias", workspaceID)
	if err != nil {
		t.Errorf("GetByName() error = %v", err)
	}
	if retrieved.ID != "alias-1" {
		t.Errorf("Expected ID 'alias-1', got %q", retrieved.ID)
	}

	// Test GetByName non-existent
	_, err = repo.GetByName("nonexistent", workspaceID)
	if err != ports.ErrAliasNotFound {
		t.Errorf("Expected ErrAliasNotFound, got %v", err)
	}

	// Test GetByName non-existent from same workspace
	_, err = repo.GetByName("nonexistent-alias", workspaceID)
	if err != ports.ErrAliasNotFound {
		t.Errorf("Expected ErrAliasNotFound for non-existent alias, got %v", err)
	}
}

// TestAliasRepository_GetByName_GlobalAlias tests global alias retrieval
func TestAliasRepository_GetByName_GlobalAlias(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewAliasRepository(adapter.DB())

	// Create global alias (empty workspace ID)
	globalAlias := &domain.JQLAlias{
		ID:          "global-1",
		Name:        "global-alias",
		JQL:         "priority = High",
		Description: "Global high priority",
		WorkspaceID: "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_ = repo.Create(globalAlias)

	// Should be accessible from any workspace (including default)
	retrieved1, err := repo.GetByName("global-alias", "default")
	if err != nil {
		t.Errorf("GetByName() from default workspace error = %v", err)
	}
	if retrieved1 == nil || retrieved1.WorkspaceID != "" {
		t.Error("Expected global alias (empty WorkspaceID)")
	}

	retrieved2, err := repo.GetByName("global-alias", "nonexistent-workspace")
	if err != nil {
		t.Errorf("GetByName() from nonexistent workspace error = %v", err)
	}
	if retrieved2 == nil || retrieved2.WorkspaceID != "" {
		t.Error("Expected global alias (empty WorkspaceID)")
	}

	// Create workspace-specific alias with DIFFERENT name
	workspaceAlias := &domain.JQLAlias{
		ID:          "workspace-specific-1",
		Name:        "workspace-specific",
		JQL:         "priority = Highest AND assignee = currentUser()",
		Description: "Workspace specific",
		WorkspaceID: "default",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_ = repo.Create(workspaceAlias)

	// Should get workspace-specific version from default workspace
	retrieved3, err := repo.GetByName("workspace-specific", "default")
	if err != nil {
		t.Errorf("GetByName() error = %v", err)
	}
	if retrieved3.WorkspaceID != "default" {
		t.Errorf("Expected workspace-specific alias, got WorkspaceID %q", retrieved3.WorkspaceID)
	}
	if retrieved3.JQL != "priority = Highest AND assignee = currentUser()" {
		t.Errorf("Expected workspace-specific JQL, got %q", retrieved3.JQL)
	}

	// Workspace-specific should NOT be visible from other workspaces
	_, err = repo.GetByName("workspace-specific", "other-workspace")
	if err != ports.ErrAliasNotFound {
		t.Errorf("Expected ErrAliasNotFound for workspace-specific alias from other workspace, got %v", err)
	}
}

// TestAliasRepository_List tests listing aliases
func TestAliasRepository_List(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewAliasRepository(adapter.DB())
	// Use "default" workspace which is created by migration
	workspaceID := "default"

	// Create global alias
	globalAlias := &domain.JQLAlias{
		ID:          "global-1",
		Name:        "global",
		JQL:         "priority = High",
		Description: "Global",
		WorkspaceID: "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_ = repo.Create(globalAlias)

	// Create workspace-specific aliases
	aliases := []*domain.JQLAlias{
		{
			ID:          "alias-1",
			Name:        "my-bugs",
			JQL:         "assignee = currentUser() AND type = Bug",
			Description: "My bugs",
			WorkspaceID: workspaceID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "alias-2",
			Name:        "high-priority",
			JQL:         "priority = High",
			Description: "High priority",
			WorkspaceID: workspaceID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, alias := range aliases {
		_ = repo.Create(alias)
	}

	// Create global alias (empty workspace)
	otherAlias := &domain.JQLAlias{
		ID:          "other-1",
		Name:        "global-other",
		JQL:         "status = Done",
		Description: "Global alias",
		WorkspaceID: "", // Global
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_ = repo.Create(otherAlias)

	// List for test-workspace
	listed, err := repo.List(workspaceID)
	if err != nil {
		t.Errorf("List() error = %v", err)
	}

	// Should include workspace aliases + global aliases
	if len(listed) < 3 {
		t.Errorf("Expected at least 3 aliases (2 workspace + 1 global), got %d", len(listed))
	}

	// Verify workspace aliases are present
	names := make(map[string]bool)
	for _, alias := range listed {
		names[alias.Name] = true
	}

	expected := []string{"global", "my-bugs", "high-priority"}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("Expected alias %q in list", name)
		}
	}

	// Global aliases should be included
	if !names["global-other"] {
		t.Error("Should include global alias")
	}
}

// TestAliasRepository_Update tests updating aliases
func TestAliasRepository_Update(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewAliasRepository(adapter.DB())

	alias := &domain.JQLAlias{
		ID:          "alias-1",
		Name:        "test-alias",
		JQL:         "status = Open",
		Description: "Old description",
		WorkspaceID: "", // Global for benchmarks/tests
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_ = repo.Create(alias)

	// Update alias
	alias.JQL = "status = In Progress"
	alias.Description = "New description"
	alias.UpdatedAt = time.Now()

	err = repo.Update(alias)
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}

	// Verify update
	updated, _ := repo.Get("alias-1")
	if updated.JQL != "status = In Progress" {
		t.Errorf("Expected JQL to be updated, got %q", updated.JQL)
	}
	if updated.Description != "New description" {
		t.Errorf("Expected Description to be updated, got %q", updated.Description)
	}

	// Test updating non-existent alias
	nonExistent := &domain.JQLAlias{
		ID:          "nonexistent",
		Name:        "test",
		JQL:         "status = Done",
		Description: "Test",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = repo.Update(nonExistent)
	if err != ports.ErrAliasNotFound {
		t.Errorf("Expected ErrAliasNotFound, got %v", err)
	}
}

// TestAliasRepository_Update_PredefinedProtection tests predefined alias protection
func TestAliasRepository_Update_PredefinedProtection(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewAliasRepository(adapter.DB())

	// Create a predefined alias
	predefined := &domain.JQLAlias{
		ID:           "predefined-1",
		Name:         "system-alias",
		JQL:          "assignee = currentUser()",
		Description:  "System alias",
		IsPredefined: true,
		WorkspaceID:  "",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	_ = repo.Create(predefined)

	// Try to update predefined alias
	predefined.JQL = "different JQL"
	err = repo.Update(predefined)
	if err != ports.ErrCannotModifyPredefined {
		t.Errorf("Expected ErrCannotModifyPredefined, got %v", err)
	}
}

// TestAliasRepository_Delete tests deleting aliases
func TestAliasRepository_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewAliasRepository(adapter.DB())

	alias := &domain.JQLAlias{
		ID:          "alias-1",
		Name:        "test-alias",
		JQL:         "status = Open",
		Description: "Test",
		WorkspaceID: "", // Global for benchmarks/tests
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_ = repo.Create(alias)

	// Delete alias
	err = repo.Delete("alias-1")
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	// Verify deletion
	_, err = repo.Get("alias-1")
	if err != ports.ErrAliasNotFound {
		t.Error("Expected alias to be deleted")
	}

	// Test deleting non-existent alias
	err = repo.Delete("nonexistent")
	if err != ports.ErrAliasNotFound {
		t.Errorf("Expected ErrAliasNotFound, got %v", err)
	}
}

// TestAliasRepository_Delete_PredefinedProtection tests predefined alias protection
func TestAliasRepository_Delete_PredefinedProtection(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewAliasRepository(adapter.DB())

	// Create a predefined alias
	predefined := &domain.JQLAlias{
		ID:           "predefined-1",
		Name:         "system-alias",
		JQL:          "assignee = currentUser()",
		Description:  "System alias",
		IsPredefined: true,
		WorkspaceID:  "",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	_ = repo.Create(predefined)

	// Try to delete predefined alias
	err = repo.Delete("predefined-1")
	if err != ports.ErrCannotModifyPredefined {
		t.Errorf("Expected ErrCannotModifyPredefined, got %v", err)
	}

	// Verify alias still exists
	_, err = repo.Get("predefined-1")
	if err != nil {
		t.Error("Predefined alias should not have been deleted")
	}
}

// TestAliasRepository_DeleteByName tests deleting by name
func TestAliasRepository_DeleteByName(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewAliasRepository(adapter.DB())
	// Use "default" workspace which is created by migration
	workspaceID := "default"

	alias := &domain.JQLAlias{
		ID:          "alias-1",
		Name:        "test-alias",
		JQL:         "status = Open",
		Description: "Test",
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_ = repo.Create(alias)

	// Delete by name
	err = repo.DeleteByName("test-alias", workspaceID)
	if err != nil {
		t.Errorf("DeleteByName() error = %v", err)
	}

	// Verify deletion
	_, err = repo.GetByName("test-alias", workspaceID)
	if err != ports.ErrAliasNotFound {
		t.Error("Expected alias to be deleted")
	}

	// Test deleting non-existent from workspace
	err = repo.DeleteByName("nonexistent-alias", workspaceID)
	if err != ports.ErrAliasNotFound {
		t.Errorf("Expected ErrAliasNotFound when deleting non-existent alias, got %v", err)
	}
}

// TestAliasRepository_ConcurrentAccess tests thread-safety
func TestAliasRepository_ConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewAliasRepository(adapter.DB())

	var wg sync.WaitGroup
	operations := 50

	// Concurrent creates
	wg.Add(operations)
	for i := 0; i < operations; i++ {
		go func(idx int) {
			defer wg.Done()
			alias := &domain.JQLAlias{
				ID:          fmt.Sprintf("alias-%d", idx),
				Name:        fmt.Sprintf("alias-%d", idx),
				JQL:         "status = Open",
				Description: "Concurrent test",
				WorkspaceID: "", // Global for benchmarks/tests
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			_ = repo.Create(alias)
		}(i)
	}

	wg.Wait()

	// Verify all aliases were created
	aliases, err := repo.List("test-workspace")
	if err != nil {
		t.Errorf("List() error = %v", err)
	}

	if len(aliases) < operations {
		t.Errorf("Expected at least %d aliases, got %d", operations, len(aliases))
	}

	// Concurrent reads
	wg.Add(operations)
	for i := 0; i < operations; i++ {
		go func(idx int) {
			defer wg.Done()
			_, _ = repo.Get(fmt.Sprintf("alias-%d", idx))
		}(i)
	}

	wg.Wait()

	// Concurrent updates
	wg.Add(operations)
	for i := 0; i < operations; i++ {
		go func(idx int) {
			defer wg.Done()
			alias, err := repo.Get(fmt.Sprintf("alias-%d", idx))
			if err == nil {
				alias.JQL = fmt.Sprintf("status = Updated-%d", idx)
				_ = repo.Update(alias)
			}
		}(i)
	}

	wg.Wait()

	// Verify updates
	for i := 0; i < operations; i++ {
		alias, err := repo.Get(fmt.Sprintf("alias-%d", i))
		if err != nil {
			t.Errorf("Failed to get alias %d: %v", i, err)
			continue
		}
		expectedJQL := fmt.Sprintf("status = Updated-%d", i)
		if alias.JQL != expectedJQL {
			t.Errorf("Expected JQL %q, got %q", expectedJQL, alias.JQL)
		}
	}
}

// BenchmarkAliasRepository_Create benchmarks alias creation
func BenchmarkAliasRepository_Create(b *testing.B) {
	tmpDir := b.TempDir()
	dbPath := filepath.Join(tmpDir, "bench.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		b.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewAliasRepository(adapter.DB())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		alias := &domain.JQLAlias{
			ID:          fmt.Sprintf("alias-%d", i),
			Name:        fmt.Sprintf("alias-%d", i),
			JQL:         "status = Open",
			Description: "Benchmark",
			WorkspaceID: "", // Global for benchmarks/tests
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		_ = repo.Create(alias)
	}
}

// BenchmarkAliasRepository_Get benchmarks alias retrieval
func BenchmarkAliasRepository_Get(b *testing.B) {
	tmpDir := b.TempDir()
	dbPath := filepath.Join(tmpDir, "bench.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		b.Fatalf("Failed to create adapter: %v", err)
	}
	defer adapter.Close()

	repo := NewAliasRepository(adapter.DB())

	// Create aliases
	for i := 0; i < 100; i++ {
		alias := &domain.JQLAlias{
			ID:          fmt.Sprintf("alias-%d", i),
			Name:        fmt.Sprintf("alias-%d", i),
			JQL:         "status = Open",
			Description: "Benchmark",
			WorkspaceID: "", // Global for benchmarks/tests
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		_ = repo.Create(alias)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.Get(fmt.Sprintf("alias-%d", i%100))
	}
}
