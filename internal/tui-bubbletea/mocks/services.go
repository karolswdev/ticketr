// Package mocks provides mock implementations of services for testing the TUI.
package mocks

import (
	"fmt"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// MockWorkspaceService is a mock implementation of WorkspaceService for testing.
type MockWorkspaceService struct {
	// Configurable behavior
	CurrentWorkspace *domain.Workspace
	Workspaces       []domain.Workspace
	CurrentErr       error
	ListErr          error
}

// NewMockWorkspaceService creates a new mock workspace service with default test data.
func NewMockWorkspaceService() *MockWorkspaceService {
	return &MockWorkspaceService{
		CurrentWorkspace: &domain.Workspace{
			ID:         "ws-test-1",
			Name:       "Test Workspace",
			ProjectKey: "TEST",
			JiraURL:    "https://test.atlassian.net",
		},
		Workspaces: []domain.Workspace{
			{
				ID:         "ws-test-1",
				Name:       "Test Workspace",
				ProjectKey: "TEST",
				JiraURL:    "https://test.atlassian.net",
			},
			{
				ID:         "ws-test-2",
				Name:       "Second Workspace",
				ProjectKey: "PROJ",
				JiraURL:    "https://proj.atlassian.net",
			},
		},
	}
}

// Current returns the current workspace or an error if configured.
func (m *MockWorkspaceService) Current() (*domain.Workspace, error) {
	if m.CurrentErr != nil {
		return nil, m.CurrentErr
	}
	return m.CurrentWorkspace, nil
}

// List returns the list of workspaces or an error if configured.
func (m *MockWorkspaceService) List() ([]domain.Workspace, error) {
	if m.ListErr != nil {
		return nil, m.ListErr
	}
	return m.Workspaces, nil
}

// WithError configures the mock to return errors.
func (m *MockWorkspaceService) WithError(err error) *MockWorkspaceService {
	m.CurrentErr = err
	m.ListErr = err
	return m
}

// WithCurrentError configures the mock to return an error for Current().
func (m *MockWorkspaceService) WithCurrentError(err error) *MockWorkspaceService {
	m.CurrentErr = err
	return m
}

// WithListError configures the mock to return an error for List().
func (m *MockWorkspaceService) WithListError(err error) *MockWorkspaceService {
	m.ListErr = err
	return m
}

// WithWorkspace configures the current workspace.
func (m *MockWorkspaceService) WithWorkspace(ws *domain.Workspace) *MockWorkspaceService {
	m.CurrentWorkspace = ws
	return m
}

// WithWorkspaces configures the workspace list.
func (m *MockWorkspaceService) WithWorkspaces(workspaces []domain.Workspace) *MockWorkspaceService {
	m.Workspaces = workspaces
	return m
}

// MockTicketQueryService is a mock implementation of TicketQueryService for testing.
type MockTicketQueryService struct {
	// Configurable behavior
	Tickets map[string][]domain.Ticket // workspace ID -> tickets
	Err     error
}

// NewMockTicketQueryService creates a new mock ticket query service with default test data.
func NewMockTicketQueryService() *MockTicketQueryService {
	testTickets := []domain.Ticket{
		{
			JiraID:      "TEST-1",
			Title:       "First Test Ticket",
			Description: "This is a test ticket",
			CustomFields: map[string]string{
				"Status":   "To Do",
				"Type":     "Story",
				"Priority": "High",
			},
			AcceptanceCriteria: []string{"AC1: Must work", "AC2: Must be tested"},
			Tasks:              []domain.Task{},
		},
		{
			JiraID:      "TEST-2",
			Title:       "Second Test Ticket",
			Description: "Another test ticket",
			CustomFields: map[string]string{
				"Status":   "In Progress",
				"Type":     "Bug",
				"Priority": "Medium",
			},
			AcceptanceCriteria: []string{"AC1: Fix the bug"},
			Tasks: []domain.Task{
				{
					JiraID:      "TEST-3",
					Title:       "Subtask 1",
					Description: "First subtask",
					CustomFields: map[string]string{
						"Status":   "Done",
						"Type":     "Sub-task",
						"Priority": "Low",
					},
				},
				{
					JiraID:      "TEST-4",
					Title:       "Subtask 2",
					Description: "Second subtask",
					CustomFields: map[string]string{
						"Status":   "To Do",
						"Type":     "Sub-task",
						"Priority": "Medium",
					},
				},
			},
		},
		{
			JiraID:      "TEST-5",
			Title:       "Third Test Ticket",
			Description: "Yet another test ticket",
			CustomFields: map[string]string{
				"Status":   "Done",
				"Type":     "Task",
				"Priority": "Low",
			},
			AcceptanceCriteria: []string{},
			Tasks:              []domain.Task{},
		},
	}

	return &MockTicketQueryService{
		Tickets: map[string][]domain.Ticket{
			"ws-test-1": testTickets,
			"ws-test-2": {
				{
					JiraID:      "PROJ-1",
					Title:       "Project Ticket",
					Description: "A project ticket",
					CustomFields: map[string]string{
						"Status":   "To Do",
						"Type":     "Story",
						"Priority": "High",
					},
					AcceptanceCriteria: []string{},
					Tasks:              []domain.Task{},
				},
			},
		},
	}
}

// ListByWorkspace returns tickets for the given workspace ID or an error if configured.
func (m *MockTicketQueryService) ListByWorkspace(workspaceID string) ([]domain.Ticket, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	if workspaceID == "" {
		return nil, fmt.Errorf("workspace ID cannot be empty")
	}

	tickets, ok := m.Tickets[workspaceID]
	if !ok {
		return []domain.Ticket{}, nil
	}

	return tickets, nil
}

// WithError configures the mock to return an error.
func (m *MockTicketQueryService) WithError(err error) *MockTicketQueryService {
	m.Err = err
	return m
}

// WithTickets configures the tickets for a specific workspace.
func (m *MockTicketQueryService) WithTickets(workspaceID string, tickets []domain.Ticket) *MockTicketQueryService {
	if m.Tickets == nil {
		m.Tickets = make(map[string][]domain.Ticket)
	}
	m.Tickets[workspaceID] = tickets
	return m
}

// WithEmptyTickets configures the mock to return empty tickets for all workspaces.
func (m *MockTicketQueryService) WithEmptyTickets() *MockTicketQueryService {
	m.Tickets = make(map[string][]domain.Ticket)
	return m
}
