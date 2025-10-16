package services

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/state"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	tickets []domain.Ticket
	savedTickets []domain.Ticket
}

func (m *MockRepository) GetTickets(filepath string) ([]domain.Ticket, error) {
	return m.tickets, nil
}

func (m *MockRepository) SaveTickets(filepath string, tickets []domain.Ticket) error {
	m.savedTickets = tickets
	return nil
}

// MockJiraPort is a mock implementation of the JiraPort interface
type MockJiraPort struct {
	UpdateTicketCalled int
	CreateTicketCalled int
	UpdateTaskCalled   int
	CreateTaskCalled   int
	LastCreatedTask    *domain.Task
	LastUpdatedTask    *domain.Task
}

func (m *MockJiraPort) Authenticate() error {
	return nil
}

func (m *MockJiraPort) CreateTask(task domain.Task, parentID string) (string, error) {
	m.CreateTaskCalled++
	m.LastCreatedTask = &task
	return "MOCK-TASK-123", nil
}

func (m *MockJiraPort) UpdateTask(task domain.Task) error {
	m.UpdateTaskCalled++
	m.LastUpdatedTask = &task
	return nil
}

func (m *MockJiraPort) GetProjectIssueTypes() (map[string][]string, error) {
	return nil, nil
}

func (m *MockJiraPort) GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error) {
	return nil, nil
}

func (m *MockJiraPort) CreateTicket(ticket domain.Ticket) (string, error) {
	m.CreateTicketCalled++
	return fmt.Sprintf("MOCK-%d", m.CreateTicketCalled), nil
}

func (m *MockJiraPort) UpdateTicket(ticket domain.Ticket) error {
	m.UpdateTicketCalled++
	return nil
}

func (m *MockJiraPort) SearchTickets(projectKey string, jql string) ([]domain.Ticket, error) {
	return []domain.Ticket{}, nil
}

func TestPushService_SkipsUnchangedTickets(t *testing.T) {
	// Create a temporary state file
	tmpDir := t.TempDir()
	stateFile := filepath.Join(tmpDir, ".ticketr.state")
	
	// Create state manager
	stateManager := state.NewStateManager(stateFile)
	
	// Create test ticket
	ticket1 := domain.Ticket{
		Title:       "Test Ticket 1",
		Description: "Test Description 1",
		JiraID:      "TICKET-1",
		CustomFields: map[string]string{
			"Priority": "High",
		},
	}
	
	// Pre-populate the state file with the same hash
	stateManager.SetStoredState("TICKET-1", state.TicketState{
		LocalHash:  stateManager.CalculateHash(ticket1),
		RemoteHash: stateManager.CalculateHash(ticket1),
	})
	if err := stateManager.Save(); err != nil {
		t.Fatalf("Failed to save initial state: %v", err)
	}
	
	// Create a second ticket that has changed
	ticket2 := domain.Ticket{
		Title:       "Test Ticket 2",
		Description: "Test Description 2 - Updated",
		JiraID:      "TICKET-2",
		CustomFields: map[string]string{
			"Priority": "Low",
		},
	}
	
	// Pre-populate with a different hash (simulating changed content)
	stateManager.SetStoredState("TICKET-2", state.TicketState{
		LocalHash:  "different-hash",
		RemoteHash: "different-hash",
	})
	if err := stateManager.Save(); err != nil {
		t.Fatalf("Failed to save initial state: %v", err)
	}
	
	// Create mock repository with both tickets
	mockRepo := &MockRepository{
		tickets: []domain.Ticket{ticket1, ticket2},
	}
	
	// Create mock Jira client
	mockJira := &MockJiraPort{}
	
	// Create push service
	pushService := NewPushService(mockRepo, mockJira, stateManager)
	
	// Run push
	result, err := pushService.PushTickets("test.md", ProcessOptions{})
	if err != nil {
		t.Fatalf("PushTickets failed: %v", err)
	}
	
	// Verify that UpdateTicket was only called once (for ticket2)
	if mockJira.UpdateTicketCalled != 1 {
		t.Errorf("Expected UpdateTicket to be called 1 time, got %d", mockJira.UpdateTicketCalled)
	}
	
	// Verify the result
	if result.TicketsUpdated != 1 {
		t.Errorf("Expected 1 ticket updated, got %d", result.TicketsUpdated)
	}
	
	// Verify the state file was updated
	newStateManager := state.NewStateManager(stateFile)
	if err := newStateManager.Load(); err != nil {
		t.Fatalf("Failed to load updated state: %v", err)
	}
	
	// Check that ticket2's hash was updated
	storedState, exists := newStateManager.GetStoredState("TICKET-2")
	if !exists {
		t.Error("TICKET-2 state not found in state")
	}
	expectedHash := newStateManager.CalculateHash(ticket2)
	if storedState.LocalHash != expectedHash || storedState.RemoteHash != expectedHash {
		t.Errorf("TICKET-2 state not updated correctly. Got local=%s remote=%s, expected %s", storedState.LocalHash, storedState.RemoteHash, expectedHash)
	}
	
	// Clean up
	os.Remove(stateFile)
}

func TestPushService_FieldInheritance_ParentFields(t *testing.T) {
	// Test that tasks inherit parent fields when they don't have their own
	tmpDir := t.TempDir()
	stateFile := filepath.Join(tmpDir, ".ticketr.state")

	// Create state manager
	stateManager := state.NewStateManager(stateFile)

	// Create test ticket with CustomFields and a task without fields
	ticket := domain.Ticket{
		Title:       "Parent Ticket",
		Description: "Parent Description",
		JiraID:      "TICKET-100",
		CustomFields: map[string]string{
			"Sprint":   "Sprint 1",
			"Priority": "High",
			"Team":     "Engineering",
		},
		Tasks: []domain.Task{
			{
				Title:       "Child Task",
				Description: "Task Description",
				CustomFields: map[string]string{},
			},
		},
	}

	// Create mock repository
	mockRepo := &MockRepository{
		tickets: []domain.Ticket{ticket},
	}

	// Create mock Jira client
	mockJira := &MockJiraPort{}

	// Create push service
	pushService := NewPushService(mockRepo, mockJira, stateManager)

	// Run push
	result, err := pushService.PushTickets("test.md", ProcessOptions{})
	if err != nil {
		t.Fatalf("PushTickets failed: %v", err)
	}

	// Verify that CreateTask was called
	if mockJira.CreateTaskCalled != 1 {
		t.Errorf("Expected CreateTask to be called 1 time, got %d", mockJira.CreateTaskCalled)
	}

	// Verify that the task received the inherited fields
	if mockJira.LastCreatedTask == nil {
		t.Fatal("LastCreatedTask is nil")
	}

	expectedFields := map[string]string{
		"Sprint":   "Sprint 1",
		"Priority": "High",
		"Team":     "Engineering",
	}

	for key, expectedValue := range expectedFields {
		actualValue, exists := mockJira.LastCreatedTask.CustomFields[key]
		if !exists {
			t.Errorf("Expected field '%s' to be inherited, but it was not present", key)
		} else if actualValue != expectedValue {
			t.Errorf("Expected field '%s' to have value '%s', got '%s'", key, expectedValue, actualValue)
		}
	}

	// Verify the result
	if result.TasksCreated != 1 {
		t.Errorf("Expected 1 task created, got %d", result.TasksCreated)
	}

	// Clean up
	os.Remove(stateFile)
}

func TestPushService_FieldInheritance_TaskOverrides(t *testing.T) {
	// Test that task fields override parent fields
	tmpDir := t.TempDir()
	stateFile := filepath.Join(tmpDir, ".ticketr.state")

	// Create state manager
	stateManager := state.NewStateManager(stateFile)

	// Create test ticket with CustomFields and a task with overriding fields
	ticket := domain.Ticket{
		Title:       "Parent Ticket",
		Description: "Parent Description",
		JiraID:      "TICKET-200",
		CustomFields: map[string]string{
			"Sprint":   "Sprint 1",
			"Priority": "High",
			"Team":     "Engineering",
		},
		Tasks: []domain.Task{
			{
				Title:       "Child Task",
				Description: "Task Description",
				JiraID:      "TASK-200",
				CustomFields: map[string]string{
					"Priority": "Low",      // Override parent
					"Assignee": "John Doe", // New field
				},
			},
		},
	}

	// Create mock repository
	mockRepo := &MockRepository{
		tickets: []domain.Ticket{ticket},
	}

	// Create mock Jira client
	mockJira := &MockJiraPort{}

	// Create push service
	pushService := NewPushService(mockRepo, mockJira, stateManager)

	// Run push
	result, err := pushService.PushTickets("test.md", ProcessOptions{})
	if err != nil {
		t.Fatalf("PushTickets failed: %v", err)
	}

	// Verify that UpdateTask was called (task has JiraID)
	if mockJira.UpdateTaskCalled != 1 {
		t.Errorf("Expected UpdateTask to be called 1 time, got %d", mockJira.UpdateTaskCalled)
	}

	// Verify that the task received the correct merged fields
	if mockJira.LastUpdatedTask == nil {
		t.Fatal("LastUpdatedTask is nil")
	}

	// Check inherited fields
	if val, exists := mockJira.LastUpdatedTask.CustomFields["Sprint"]; !exists || val != "Sprint 1" {
		t.Errorf("Expected 'Sprint' to be inherited as 'Sprint 1', got '%s' (exists: %v)", val, exists)
	}

	if val, exists := mockJira.LastUpdatedTask.CustomFields["Team"]; !exists || val != "Engineering" {
		t.Errorf("Expected 'Team' to be inherited as 'Engineering', got '%s' (exists: %v)", val, exists)
	}

	// Check overridden field
	if val, exists := mockJira.LastUpdatedTask.CustomFields["Priority"]; !exists || val != "Low" {
		t.Errorf("Expected 'Priority' to be overridden to 'Low', got '%s' (exists: %v)", val, exists)
	}

	// Check task-specific field
	if val, exists := mockJira.LastUpdatedTask.CustomFields["Assignee"]; !exists || val != "John Doe" {
		t.Errorf("Expected 'Assignee' to be 'John Doe', got '%s' (exists: %v)", val, exists)
	}

	// Verify the result
	if result.TasksUpdated != 1 {
		t.Errorf("Expected 1 task updated, got %d", result.TasksUpdated)
	}

	// Clean up
	os.Remove(stateFile)
}