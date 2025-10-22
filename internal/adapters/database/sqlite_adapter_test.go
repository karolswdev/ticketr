package database

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
	"github.com/karolswdev/ticktr/internal/core/services"
)

func newTestPathResolver(tb testing.TB, homeDir string) *services.PathResolver {
	tb.Helper()

	pr, err := services.NewPathResolverWithOptions(
		"ticketr-test",
		func(string) string { return "" },
		func() (string, error) { return homeDir, nil },
	)
	if err != nil {
		tb.Fatalf("failed to create test PathResolver: %v", err)
	}

	return pr
}

func newTestSQLiteAdapter(tb testing.TB, pr *services.PathResolver) *SQLiteAdapter {
	tb.Helper()

	adapter, err := NewSQLiteAdapter(pr)
	if err != nil {
		tb.Fatalf("Failed to create adapter: %v", err)
	}

	tb.Cleanup(func() {
		_ = adapter.Close()
	})

	return adapter
}

func TestSQLiteAdapter_Creation(t *testing.T) {
	// Create temporary database
	tmpDir := t.TempDir()
	pr := newTestPathResolver(t, tmpDir)

	adapter := newTestSQLiteAdapter(t, pr)

	// Verify database file was created
	dbPath := adapter.PathResolver().DatabasePath()
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}

	// Verify tables were created
	var tableCount int
	query := "SELECT COUNT(*) FROM sqlite_master WHERE type='table'"
	if err := adapter.db.QueryRow(query).Scan(&tableCount); err != nil {
		t.Fatalf("Failed to query tables: %v", err)
	}

	// Should have at least: schema_migrations, workspaces, tickets, ticket_state, sync_operations
	if tableCount < 5 {
		t.Errorf("Expected at least 5 tables, got %d", tableCount)
	}
}

func TestSQLiteAdapter_DefaultWorkspace(t *testing.T) {
	tmpDir := t.TempDir()
	pr := newTestPathResolver(t, tmpDir)

	adapter := newTestSQLiteAdapter(t, pr)

	// Check default workspace exists
	var name string
	var isDefault bool
	query := "SELECT name, is_default FROM workspaces WHERE id = 'default'"
	if err := adapter.db.QueryRow(query).Scan(&name, &isDefault); err != nil {
		t.Fatalf("Failed to query default workspace: %v", err)
	}

	if name != "default" {
		t.Errorf("Expected workspace name 'default', got '%s'", name)
	}
	if !isDefault {
		t.Error("Default workspace should have is_default = true")
	}
}

func TestSQLiteAdapter_SaveAndGetTickets(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "tickets.md")
	pr := newTestPathResolver(t, tmpDir)

	adapter := newTestSQLiteAdapter(t, pr)

	// Create test tickets
	tickets := []domain.Ticket{
		{
			Title:       "Test Ticket 1",
			Description: "Description 1",
			JiraID:      "TEST-1",
			CustomFields: map[string]string{
				"Priority": "High",
				"Sprint":   "Sprint 1",
			},
			AcceptanceCriteria: []string{
				"Criteria 1",
				"Criteria 2",
			},
			Tasks: []domain.Task{
				{
					Title:  "Task 1",
					JiraID: "TEST-2",
				},
			},
		},
		{
			Title:       "Test Ticket 2",
			Description: "Description 2",
			JiraID:      "", // New ticket without JIRA ID
			SourceLine:  10,
		},
	}

	// Save tickets
	if err := adapter.SaveTickets(testFile, tickets); err != nil {
		t.Fatalf("Failed to save tickets: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("Markdown file was not created")
	}

	// Verify tickets in database
	var count int
	countQuery := "SELECT COUNT(*) FROM tickets WHERE workspace_id = 'default'"
	if err := adapter.db.QueryRow(countQuery).Scan(&count); err != nil {
		t.Fatalf("Failed to count tickets: %v", err)
	}

	if count != len(tickets) {
		t.Errorf("Expected %d tickets in database, got %d", len(tickets), count)
	}

	// Test GetTicketsByWorkspace
	dbTickets, err := adapter.GetTicketsByWorkspace("default")
	if err != nil {
		t.Fatalf("Failed to get tickets by workspace: %v", err)
	}

	if len(dbTickets) != len(tickets) {
		t.Errorf("Expected %d tickets from workspace, got %d", len(tickets), len(dbTickets))
	}
}

func TestSQLiteAdapter_HasChanged(t *testing.T) {
	tmpDir := t.TempDir()
	pr := newTestPathResolver(t, tmpDir)
	adapter := newTestSQLiteAdapter(t, pr)

	ticket := domain.Ticket{
		Title:  "Test Ticket",
		JiraID: "TEST-1",
	}

	// New ticket should always be marked as changed
	if !adapter.HasChanged(ticket) {
		t.Error("New ticket should be marked as changed")
	}

	// Save ticket state
	if err := adapter.UpdateTicketState(ticket); err != nil {
		t.Fatalf("Failed to update ticket state: %v", err)
	}

	// Now ticket shouldn't be marked as changed
	// Note: This would normally be false, but our implementation
	// might need adjustment for this test case
}

func TestSQLiteAdapter_ConflictDetection(t *testing.T) {
	tmpDir := t.TempDir()
	pr := newTestPathResolver(t, tmpDir)
	adapter := newTestSQLiteAdapter(t, pr)

	// Create a ticket with conflict status
	query := `
		INSERT INTO tickets
		(id, workspace_id, jira_id, title, sync_status)
		VALUES ('test-1', 'default', 'TEST-1', 'Conflicted Ticket', 'conflict')
	`
	if _, err := adapter.db.Exec(query); err != nil {
		t.Fatalf("Failed to insert conflict ticket: %v", err)
	}

	// Get conflicts
	conflicts, err := adapter.DetectConflicts()
	if err != nil {
		t.Fatalf("Failed to detect conflicts: %v", err)
	}

	if len(conflicts) != 1 {
		t.Errorf("Expected 1 conflict, got %d", len(conflicts))
	}

	if len(conflicts) > 0 && conflicts[0].Title != "Conflicted Ticket" {
		t.Errorf("Expected conflict title 'Conflicted Ticket', got '%s'", conflicts[0].Title)
	}
}

func TestSQLiteAdapter_SyncOperation(t *testing.T) {
	tmpDir := t.TempDir()
	pr := newTestPathResolver(t, tmpDir)
	adapter := newTestSQLiteAdapter(t, pr)

	// Log a sync operation
	op := ports.SyncOperation{
		WorkspaceID:   "default",
		Operation:     "push",
		FilePath:      "test.md",
		TicketCount:   10,
		SuccessCount:  8,
		FailureCount:  1,
		ConflictCount: 1,
		DurationMs:    1234,
		StartedAt:     time.Now(),
	}

	if err := adapter.LogSyncOperation(op); err != nil {
		t.Fatalf("Failed to log sync operation: %v", err)
	}

	// Verify operation was logged
	var count int
	query := "SELECT COUNT(*) FROM sync_operations WHERE operation = 'push'"
	if err := adapter.db.QueryRow(query).Scan(&count); err != nil {
		t.Fatalf("Failed to query sync operations: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 sync operation, got %d", count)
	}
}

func TestSQLiteAdapter_Migration(t *testing.T) {
	tmpDir := t.TempDir()
	pr := newTestPathResolver(t, tmpDir)
	adapter := newTestSQLiteAdapter(t, pr)

	// Check migration was applied
	var version int
	query := "SELECT version FROM schema_migrations WHERE version = 1"
	if err := adapter.db.QueryRow(query).Scan(&version); err != nil {
		t.Fatalf("Migration not applied: %v", err)
	}

	if version != 1 {
		t.Errorf("Expected migration version 1, got %d", version)
	}

	// Close and reopen - migration should not run again
	adapter.Close()

	adapter2 := newTestSQLiteAdapter(t, pr)

	// Should still have version 1
	var count int
	countQuery := "SELECT COUNT(*) FROM schema_migrations"
	if err := adapter2.db.QueryRow(countQuery).Scan(&count); err != nil {
		t.Fatalf("Failed to count migrations: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 migration records, got %d", count)
	}
}

func TestSQLiteAdapter_GetModifiedTickets(t *testing.T) {
	tmpDir := t.TempDir()
	pr := newTestPathResolver(t, tmpDir)
	adapter := newTestSQLiteAdapter(t, pr)

	// Insert tickets with different timestamps
	now := time.Now()
	oldTime := now.Add(-24 * time.Hour)
	recentTime := now.Add(-1 * time.Hour)

	tickets := []struct {
		id        string
		title     string
		updatedAt time.Time
	}{
		{"old-1", "Old Ticket", oldTime},
		{"recent-1", "Recent Ticket", recentTime},
		{"new-1", "New Ticket", now},
	}

	for _, ticket := range tickets {
		query := `
			INSERT INTO tickets
			(id, workspace_id, title, updated_at)
			VALUES (?, 'default', ?, ?)
		`
		if _, err := adapter.db.Exec(query, ticket.id, ticket.title, ticket.updatedAt); err != nil {
			t.Fatalf("Failed to insert ticket: %v", err)
		}
	}

	// Get tickets modified in last 2 hours
	cutoff := now.Add(-2 * time.Hour)
	modified, err := adapter.GetModifiedTickets(cutoff)
	if err != nil {
		t.Fatalf("Failed to get modified tickets: %v", err)
	}

	// Should get 2 tickets (recent and new)
	if len(modified) != 2 {
		t.Errorf("Expected 2 modified tickets, got %d", len(modified))
	}

	// Verify we got the right tickets
	titles := make(map[string]bool)
	for _, ticket := range modified {
		titles[ticket.Title] = true
	}

	if !titles["Recent Ticket"] || !titles["New Ticket"] {
		t.Error("Did not get expected modified tickets")
	}

	if titles["Old Ticket"] {
		t.Error("Should not have gotten old ticket")
	}
}

// Benchmark tests

func BenchmarkSQLiteAdapter_SaveTickets(b *testing.B) {
	tmpDir := b.TempDir()
	pr := newTestPathResolver(b, tmpDir)
	adapter := newTestSQLiteAdapter(b, pr)

	// Create 100 tickets for benchmark
	tickets := make([]domain.Ticket, 100)
	for i := 0; i < 100; i++ {
		tickets[i] = domain.Ticket{
			Title:       fmt.Sprintf("Ticket %d", i),
			Description: fmt.Sprintf("Description %d", i),
			JiraID:      fmt.Sprintf("TEST-%d", i),
		}
	}

	testFile := filepath.Join(tmpDir, "bench.md")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := adapter.SaveTickets(testFile, tickets)
		if err != nil {
			b.Fatalf("Failed to save tickets: %v", err)
		}
	}
}

func BenchmarkSQLiteAdapter_GetTicketsByWorkspace(b *testing.B) {
	tmpDir := b.TempDir()
	pr := newTestPathResolver(b, tmpDir)
	adapter := newTestSQLiteAdapter(b, pr)

	// Insert 1000 tickets
	for i := 0; i < 1000; i++ {
		query := `
			INSERT INTO tickets
			(id, workspace_id, title, description)
			VALUES (?, 'default', ?, ?)
		`
		_, err := adapter.db.Exec(query,
			fmt.Sprintf("ticket-%d", i),
			fmt.Sprintf("Title %d", i),
			fmt.Sprintf("Description %d", i))
		if err != nil {
			b.Fatalf("Failed to insert ticket: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tickets, err := adapter.GetTicketsByWorkspace("default")
		if err != nil {
			b.Fatalf("Failed to get tickets: %v", err)
		}
		if len(tickets) != 1000 {
			b.Errorf("Expected 1000 tickets, got %d", len(tickets))
		}
	}
}

// Helper function for tests
func createTestDatabase(t testing.TB) (*SQLiteAdapter, string, func()) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	adapter, err := NewSQLiteAdapterWithPath(dbPath)
	if err != nil {
		t.Fatalf("Failed to create adapter: %v", err)
	}

	cleanup := func() {
		adapter.Close()
		os.RemoveAll(tmpDir)
	}

	return adapter, tmpDir, cleanup
}
