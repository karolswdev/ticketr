package services

import (
	"os"
	"testing"

	"github.com/karolswdev/ticketr/internal/core/domain"
	"github.com/karolswdev/ticketr/internal/state"
)

// TestPullService_ResolvesConflictWithLocalWinsStrategy tests that conflicts are resolved
// by keeping local changes when using the local-wins strategy
func TestPullService_ResolvesConflictWithLocalWinsStrategy(t *testing.T) {
	// Arrange
	mockJira := &mockJiraAdapter{}
	mockRepo := &mockRepository{}
	// Use a unique state file for this test to avoid conflicts
	stateManager := state.NewStateManager("test-local-wins.state")
	// Clean up after test
	defer func() {
		_ = os.Remove("test-local-wins.state")
	}()

	// Create a conflict scenario
	localTicket := domain.Ticket{
		JiraID:      "TEST-123",
		Title:       "Local Version",
		Description: "This is the local version",
	}

	remoteTicket := domain.Ticket{
		JiraID:      "TEST-123",
		Title:       "Remote Version",
		Description: "This is the remote version",
	}

	// Setup mock repository with local ticket
	mockRepo.tickets = []domain.Ticket{localTicket}

	// Setup mock JIRA with remote ticket
	mockJira.searchResult = []domain.Ticket{remoteTicket}

	// Initialize state with different hashes to simulate conflict
	stateManager.SetStoredState("TEST-123", state.TicketState{
		LocalHash:  "old_local_hash",
		RemoteHash: "old_remote_hash",
	})

	pullService := NewPullService(mockJira, mockRepo, stateManager)

	// Act
	result, err := pullService.Pull("test.md", PullOptions{
		ProjectKey: "TEST",
		Strategy:   "local-wins",
	})

	// Assert
	if err != nil {
		t.Fatalf("Expected no error with local-wins strategy, got: %v", err)
	}

	if len(result.Conflicts) != 1 {
		t.Errorf("Expected 1 conflict, got %d", len(result.Conflicts))
	}

	// Verify that local version was kept
	if len(mockRepo.savedTickets) != 1 {
		t.Fatalf("Expected 1 ticket saved, got %d", len(mockRepo.savedTickets))
	}

	savedTicket := mockRepo.savedTickets[0]
	if savedTicket.Title != "Local Version" {
		t.Errorf("Expected local version to be kept, got title: %s", savedTicket.Title)
	}

	// Verify the state was correctly updated
	newState, exists := stateManager.GetStoredState("TEST-123")
	if !exists {
		t.Error("Expected state to be updated")
	}

	// Local hash should be updated
	localHash := stateManager.CalculateHash(localTicket)
	if newState.LocalHash != localHash {
		t.Error("Expected local hash to be updated in state")
	}
}

// TestPullService_ResolvesConflictWithRemoteWinsStrategy tests that conflicts are resolved
// by using remote changes when using the remote-wins strategy
func TestPullService_ResolvesConflictWithRemoteWinsStrategy(t *testing.T) {
	// Arrange
	mockJira := &mockJiraAdapter{}
	mockRepo := &mockRepository{}
	// Use a unique state file for this test to avoid conflicts
	stateManager := state.NewStateManager("test-remote-wins.state")
	// Clean up after test
	defer func() {
		_ = os.Remove("test-remote-wins.state")
	}()

	// Create a conflict scenario
	localTicket := domain.Ticket{
		JiraID:      "TEST-124",
		Title:       "Local Version",
		Description: "This is the local version",
	}

	remoteTicket := domain.Ticket{
		JiraID:      "TEST-124",
		Title:       "Remote Version",
		Description: "This is the remote version",
	}

	// Setup mock repository with local ticket
	mockRepo.tickets = []domain.Ticket{localTicket}

	// Setup mock JIRA with remote ticket
	mockJira.searchResult = []domain.Ticket{remoteTicket}

	// Initialize state with different hashes to simulate conflict
	stateManager.SetStoredState("TEST-124", state.TicketState{
		LocalHash:  "old_local_hash",
		RemoteHash: "old_remote_hash",
	})

	pullService := NewPullService(mockJira, mockRepo, stateManager)

	// Act
	result, err := pullService.Pull("test.md", PullOptions{
		ProjectKey: "TEST",
		Strategy:   "remote-wins",
	})

	// Assert
	if err != nil {
		t.Fatalf("Expected no error with remote-wins strategy, got: %v", err)
	}

	if len(result.Conflicts) != 1 {
		t.Errorf("Expected 1 conflict, got %d", len(result.Conflicts))
	}

	// Verify that remote version was used
	if len(mockRepo.savedTickets) != 1 {
		t.Fatalf("Expected 1 ticket saved, got %d", len(mockRepo.savedTickets))
	}

	savedTicket := mockRepo.savedTickets[0]
	if savedTicket.Title != "Remote Version" {
		t.Errorf("Expected remote version to be used, got title: %s", savedTicket.Title)
	}

	// Verify the state was correctly updated
	newState, exists := stateManager.GetStoredState("TEST-124")
	if !exists {
		t.Error("Expected state to be updated")
	}

	// Both hashes should match the remote version
	remoteHash := stateManager.CalculateHash(remoteTicket)
	if newState.LocalHash != remoteHash || newState.RemoteHash != remoteHash {
		t.Error("Expected both hashes to match remote version in state")
	}
}

// Mock implementations for testing
type mockJiraAdapter struct {
	searchResult []domain.Ticket
}

func (m *mockJiraAdapter) CreateTicket(ticket domain.Ticket) (string, error) {
	return "TEST-999", nil
}

func (m *mockJiraAdapter) UpdateTicket(ticket domain.Ticket) error {
	return nil
}

func (m *mockJiraAdapter) CreateTask(task domain.Task, parentID string) (string, error) {
	return "TEST-1000", nil
}

func (m *mockJiraAdapter) UpdateTask(task domain.Task) error {
	return nil
}

func (m *mockJiraAdapter) SearchTickets(projectKey string, jql string) ([]domain.Ticket, error) {
	return m.searchResult, nil
}

func (m *mockJiraAdapter) GetProjectIssueTypes() (map[string][]string, error) {
	return map[string][]string{"TEST": {"Task", "Sub-task"}}, nil
}

func (m *mockJiraAdapter) GetIssueTypeFields(issueType string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (m *mockJiraAdapter) Authenticate() error {
	return nil
}

type mockRepository struct {
	tickets      []domain.Ticket
	savedTickets []domain.Ticket
}

func (m *mockRepository) GetTickets(filePath string) ([]domain.Ticket, error) {
	return m.tickets, nil
}

func (m *mockRepository) SaveTickets(filePath string, tickets []domain.Ticket) error {
	m.savedTickets = tickets
	return nil
}
