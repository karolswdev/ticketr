package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/karolswdev/ticketr/internal/core/domain"
	"github.com/karolswdev/ticketr/internal/core/services"
	"github.com/karolswdev/ticketr/internal/state"
)

// TestWebhookServer_UpdatesFileOnJiraEvent tests that the webhook server
// correctly processes JIRA webhook events and updates the markdown file
func TestWebhookServer_UpdatesFileOnJiraEvent(t *testing.T) {
	// Arrange
	mockJira := &mockJiraAdapter{}
	mockRepo := &mockRepository{}
	// Use a unique state file for this test
	stateManager := state.NewStateManager("test-webhook.state")
	// Clean up after test
	defer func() { _ = os.Remove("test-webhook.state") }()

	// Set up mock JIRA to return a ticket when searched
	mockJira.searchResult = []domain.Ticket{
		{
			JiraID:      "TEST-100",
			Title:       "Updated from Webhook",
			Description: "This ticket was updated via webhook",
		},
	}

	// Set up mock repository to simulate existing file
	mockRepo.tickets = []domain.Ticket{
		{
			JiraID:      "TEST-100",
			Title:       "Original Title",
			Description: "Original description",
		},
	}

	pullService := services.NewPullService(mockJira, mockRepo, stateManager)
	server := NewServer(pullService, "test.md", "", "TEST")

	// Create a test webhook payload
	payload := JiraWebhookPayload{
		WebhookEvent: "jira:issue_updated",
		Issue: struct {
			ID     string `json:"id"`
			Key    string `json:"key"`
			Fields struct {
				Summary     string `json:"summary"`
				Description string `json:"description"`
				IssueType   struct {
					Name string `json:"name"`
				} `json:"issuetype"`
				Status struct {
					Name string `json:"name"`
				} `json:"status"`
				Project struct {
					Key string `json:"key"`
				} `json:"project"`
			} `json:"fields"`
		}{
			ID:  "10001",
			Key: "TEST-100",
			Fields: struct {
				Summary     string `json:"summary"`
				Description string `json:"description"`
				IssueType   struct {
					Name string `json:"name"`
				} `json:"issuetype"`
				Status struct {
					Name string `json:"name"`
				} `json:"status"`
				Project struct {
					Key string `json:"key"`
				} `json:"project"`
			}{
				Summary:     "Updated from Webhook",
				Description: "This ticket was updated via webhook",
				IssueType: struct {
					Name string `json:"name"`
				}{Name: "Task"},
				Status: struct {
					Name string `json:"name"`
				}{Name: "In Progress"},
				Project: struct {
					Key string `json:"key"`
				}{Key: "TEST"},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	// Act
	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.HandleWebhook(w, req)

	// Assert
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify that the file was updated with webhook data
	if len(mockRepo.savedTickets) != 1 {
		t.Fatalf("Expected 1 ticket saved, got %d", len(mockRepo.savedTickets))
	}

	savedTicket := mockRepo.savedTickets[0]
	if savedTicket.Title != "Updated from Webhook" {
		t.Errorf("Expected ticket to be updated from webhook, got title: %s", savedTicket.Title)
	}

	// Verify state was updated
	newState, exists := stateManager.GetStoredState("TEST-100")
	if !exists {
		t.Error("Expected state to be updated for TEST-100")
	}

	// Both hashes should be equal since we used remote-wins strategy
	if newState.LocalHash != newState.RemoteHash {
		t.Error("Expected local and remote hashes to match after webhook update")
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
