package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// mockJiraAdapter is a test double for the JiraPort interface.
type mockJiraAdapter struct {
	// For controlling behavior
	authenticateError    error
	createTicketError    error
	updateTicketError    error
	searchTicketsError   error
	searchTicketsResults []domain.Ticket

	// For tracking calls
	updateTicketCalls []domain.Ticket
	searchTicketCalls []string

	// Function fields for custom behavior
	updateTicketFunc  func(ticket domain.Ticket) error
	searchTicketsFunc func(projectKey string, jql string) ([]domain.Ticket, error)
}

func newMockJiraAdapter() *mockJiraAdapter {
	return &mockJiraAdapter{
		updateTicketCalls:    make([]domain.Ticket, 0),
		searchTicketCalls:    make([]string, 0),
		searchTicketsResults: make([]domain.Ticket, 0),
	}
}

func (m *mockJiraAdapter) Authenticate() error {
	return m.authenticateError
}

func (m *mockJiraAdapter) CreateTask(task domain.Task, parentID string) (string, error) {
	return "", errors.New("not implemented")
}

func (m *mockJiraAdapter) UpdateTask(task domain.Task) error {
	return errors.New("not implemented")
}

func (m *mockJiraAdapter) GetProjectIssueTypes() (map[string][]string, error) {
	return nil, errors.New("not implemented")
}

func (m *mockJiraAdapter) GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}

func (m *mockJiraAdapter) CreateTicket(ticket domain.Ticket) (string, error) {
	if m.createTicketError != nil {
		return "", m.createTicketError
	}
	return "PROJ-NEW", nil
}

func (m *mockJiraAdapter) UpdateTicket(ticket domain.Ticket) error {
	// Use custom function if provided
	if m.updateTicketFunc != nil {
		return m.updateTicketFunc(ticket)
	}

	m.updateTicketCalls = append(m.updateTicketCalls, ticket)
	if m.updateTicketError != nil {
		return m.updateTicketError
	}
	return nil
}

func (m *mockJiraAdapter) SearchTickets(ctx context.Context, projectKey string, jql string, progressCallback ports.JiraProgressCallback) ([]domain.Ticket, error) {
	// Use custom function if provided
	if m.searchTicketsFunc != nil {
		return m.searchTicketsFunc(projectKey, jql)
	}

	m.searchTicketCalls = append(m.searchTicketCalls, jql)
	if m.searchTicketsError != nil {
		return nil, m.searchTicketsError
	}

	// If JQL contains `key = "TICKET-ID"`, return only that ticket
	// This simulates individual ticket lookups
	if len(jql) > 7 && jql[:6] == `key = ` {
		// Extract ticket ID from JQL (format: key = "PROJ-123")
		ticketID := jql[7 : len(jql)-1] // Remove 'key = "' and trailing '"'
		for _, ticket := range m.searchTicketsResults {
			if ticket.JiraID == ticketID {
				return []domain.Ticket{ticket}, nil
			}
		}
		return []domain.Ticket{}, nil // Ticket not found
	}

	// Otherwise return all results
	return m.searchTicketsResults, nil
}

// Helper to set up mock to return specific tickets
func (m *mockJiraAdapter) withTickets(tickets []domain.Ticket) *mockJiraAdapter {
	m.searchTicketsResults = tickets
	return m
}

// Helper to make update fail on specific ticket IDs
func (m *mockJiraAdapter) failOnTickets(ticketIDs ...string) *mockJiraAdapter {
	failMap := make(map[string]bool)
	for _, id := range ticketIDs {
		failMap[id] = true
	}

	// Override UpdateTicket to fail on specific tickets
	m.updateTicketFunc = func(ticket domain.Ticket) error {
		m.updateTicketCalls = append(m.updateTicketCalls, ticket)
		if failMap[ticket.JiraID] {
			return fmt.Errorf("simulated failure for %s", ticket.JiraID)
		}
		return nil
	}
	return m
}

// TestNewBulkOperationService verifies service construction.
func TestNewBulkOperationService(t *testing.T) {
	adapter := newMockJiraAdapter()
	service := NewBulkOperationService(adapter)

	if service == nil {
		t.Fatal("expected service to be created, got nil")
	}
}

// TestExecuteOperation_ValidatesInput verifies that invalid operations are rejected.
func TestExecuteOperation_ValidatesInput(t *testing.T) {
	tests := []struct {
		name    string
		op      *domain.BulkOperation
		wantErr bool
		tickets []domain.Ticket // Mock tickets to return
	}{
		{
			name: "invalid action",
			op: &domain.BulkOperation{
				Action:    "invalid",
				TicketIDs: []string{"PROJ-1"},
				Changes:   map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name: "no ticket IDs",
			op: &domain.BulkOperation{
				Action:    domain.BulkActionUpdate,
				TicketIDs: []string{},
				Changes:   map[string]interface{}{"status": "Done"},
			},
			wantErr: true,
		},
		{
			name: "update without changes",
			op: &domain.BulkOperation{
				Action:    domain.BulkActionUpdate,
				TicketIDs: []string{"PROJ-1"},
				Changes:   nil,
			},
			wantErr: true,
		},
		{
			name: "valid update operation",
			op: &domain.BulkOperation{
				Action:    domain.BulkActionUpdate,
				TicketIDs: []string{"PROJ-1"},
				Changes:   map[string]interface{}{"status": "Done"},
			},
			wantErr: false,
			tickets: []domain.Ticket{
				{JiraID: "PROJ-1", Title: "Ticket 1", CustomFields: map[string]string{"status": "In Progress"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := newMockJiraAdapter()
			if tt.tickets != nil {
				adapter.withTickets(tt.tickets)
			}
			service := NewBulkOperationService(adapter)
			ctx := context.Background()

			_, err := service.ExecuteOperation(ctx, tt.op, nil)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}

// TestExecuteOperation_SuccessfulUpdate verifies successful bulk update.
func TestExecuteOperation_SuccessfulUpdate(t *testing.T) {
	adapter := newMockJiraAdapter()
	adapter.withTickets([]domain.Ticket{
		{JiraID: "PROJ-1", Title: "Ticket 1", CustomFields: map[string]string{"status": "In Progress"}},
		{JiraID: "PROJ-2", Title: "Ticket 2", CustomFields: map[string]string{"status": "In Progress"}},
		{JiraID: "PROJ-3", Title: "Ticket 3", CustomFields: map[string]string{"status": "In Progress"}},
	})

	service := NewBulkOperationService(adapter)

	op := domain.NewBulkOperation(
		domain.BulkActionUpdate,
		[]string{"PROJ-1", "PROJ-2", "PROJ-3"},
		map[string]interface{}{"status": "Done", "assignee": "john.doe"},
	)

	ctx := context.Background()
	result, err := service.ExecuteOperation(ctx, op, nil)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.SuccessCount != 3 {
		t.Errorf("expected SuccessCount=3, got %d", result.SuccessCount)
	}

	if result.FailureCount != 0 {
		t.Errorf("expected FailureCount=0, got %d", result.FailureCount)
	}

	if len(result.SuccessfulTickets) != 3 {
		t.Errorf("expected 3 successful tickets, got %d", len(result.SuccessfulTickets))
	}

	if len(result.FailedTickets) != 0 {
		t.Errorf("expected 0 failed tickets, got %d", len(result.FailedTickets))
	}

	// Verify UpdateTicket was called 3 times
	if len(adapter.updateTicketCalls) != 3 {
		t.Errorf("expected 3 UpdateTicket calls, got %d", len(adapter.updateTicketCalls))
	}

	// Verify changes were applied
	for _, call := range adapter.updateTicketCalls {
		if call.CustomFields["status"] != "Done" {
			t.Errorf("expected status=Done, got %s", call.CustomFields["status"])
		}
		if call.CustomFields["assignee"] != "john.doe" {
			t.Errorf("expected assignee=john.doe, got %s", call.CustomFields["assignee"])
		}
	}
}

// TestExecuteOperation_PartialFailure verifies handling of partial failures.
func TestExecuteOperation_PartialFailure(t *testing.T) {
	adapter := newMockJiraAdapter()

	// Mock will return all 3 tickets when searched
	tickets := []domain.Ticket{
		{JiraID: "PROJ-1", Title: "Ticket 1", CustomFields: map[string]string{"status": "In Progress"}},
		{JiraID: "PROJ-2", Title: "Ticket 2", CustomFields: map[string]string{"status": "In Progress"}},
		{JiraID: "PROJ-3", Title: "Ticket 3", CustomFields: map[string]string{"status": "In Progress"}},
	}
	adapter.searchTicketsResults = tickets

	// Make update fail for PROJ-2
	updateCalls := 0
	adapter.updateTicketFunc = func(ticket domain.Ticket) error {
		adapter.updateTicketCalls = append(adapter.updateTicketCalls, ticket)
		updateCalls++
		if ticket.JiraID == "PROJ-2" {
			return fmt.Errorf("simulated failure for PROJ-2")
		}
		return nil
	}

	service := NewBulkOperationService(adapter)

	op := domain.NewBulkOperation(
		domain.BulkActionUpdate,
		[]string{"PROJ-1", "PROJ-2", "PROJ-3"},
		map[string]interface{}{"status": "Done"},
	)

	ctx := context.Background()
	result, err := service.ExecuteOperation(ctx, op, nil)

	// Should return error for partial failure
	if err == nil {
		t.Error("expected error for partial failure, got nil")
	}

	if result.SuccessCount != 2 {
		t.Errorf("expected SuccessCount=2, got %d", result.SuccessCount)
	}

	if result.FailureCount != 1 {
		t.Errorf("expected FailureCount=1, got %d", result.FailureCount)
	}

	if len(result.FailedTickets) != 1 {
		t.Errorf("expected 1 failed ticket, got %d", len(result.FailedTickets))
	} else if result.FailedTickets[0] != "PROJ-2" {
		t.Errorf("expected failed ticket PROJ-2, got %s", result.FailedTickets[0])
	}

	// Verify error message is recorded
	if errMsg, ok := result.Errors["PROJ-2"]; !ok {
		t.Error("expected error message for PROJ-2")
	} else if errMsg == "" {
		t.Error("expected non-empty error message")
	}

	// Note: Rollback will be attempted but we can't easily verify it without more sophisticated mocking
}

// TestExecuteOperation_AllFailures verifies handling when all operations fail.
func TestExecuteOperation_AllFailures(t *testing.T) {
	adapter := newMockJiraAdapter()
	adapter.updateTicketError = fmt.Errorf("simulated failure")
	adapter.withTickets([]domain.Ticket{
		{JiraID: "PROJ-1", Title: "Ticket 1", CustomFields: map[string]string{}},
		{JiraID: "PROJ-2", Title: "Ticket 2", CustomFields: map[string]string{}},
	})

	service := NewBulkOperationService(adapter)

	op := domain.NewBulkOperation(
		domain.BulkActionUpdate,
		[]string{"PROJ-1", "PROJ-2"},
		map[string]interface{}{"status": "Done"},
	)

	ctx := context.Background()
	result, err := service.ExecuteOperation(ctx, op, nil)

	if err == nil {
		t.Error("expected error when all operations fail, got nil")
	}

	if result.SuccessCount != 0 {
		t.Errorf("expected SuccessCount=0, got %d", result.SuccessCount)
	}

	if result.FailureCount != 2 {
		t.Errorf("expected FailureCount=2, got %d", result.FailureCount)
	}
}

// TestExecuteOperation_ContextCancellation verifies context cancellation handling.
func TestExecuteOperation_ContextCancellation(t *testing.T) {
	adapter := newMockJiraAdapter()

	// Set up tickets
	tickets := []domain.Ticket{
		{JiraID: "PROJ-1", Title: "Ticket 1", CustomFields: map[string]string{}},
		{JiraID: "PROJ-2", Title: "Ticket 2", CustomFields: map[string]string{}},
		{JiraID: "PROJ-3", Title: "Ticket 3", CustomFields: map[string]string{}},
	}
	adapter.searchTicketsResults = tickets

	// Cancel context after first ticket
	cancelCalls := 0
	adapter.updateTicketFunc = func(ticket domain.Ticket) error {
		adapter.updateTicketCalls = append(adapter.updateTicketCalls, ticket)
		cancelCalls++
		return nil
	}

	service := NewBulkOperationService(adapter)

	op := domain.NewBulkOperation(
		domain.BulkActionUpdate,
		[]string{"PROJ-1", "PROJ-2", "PROJ-3"},
		map[string]interface{}{"status": "Done"},
	)

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	result, err := service.ExecuteOperation(ctx, op, nil)

	// Should return cancellation error early
	if err == nil {
		t.Error("expected error for cancelled context, got nil")
	}

	// Should process at least the first ticket before checking cancellation
	// (cancellation is checked between tickets, not during)
	if result.SuccessCount > 1 {
		t.Errorf("expected at most 1 successful operation before cancellation, got %d", result.SuccessCount)
	}
}

// TestExecuteOperation_SuccessfulMove verifies successful bulk move.
func TestExecuteOperation_SuccessfulMove(t *testing.T) {
	adapter := newMockJiraAdapter()
	adapter.withTickets([]domain.Ticket{
		{JiraID: "PROJ-1", Title: "Ticket 1", CustomFields: map[string]string{"Parent": "PROJ-100"}},
		{JiraID: "PROJ-2", Title: "Ticket 2", CustomFields: map[string]string{"Parent": "PROJ-100"}},
	})

	service := NewBulkOperationService(adapter)

	op := domain.NewBulkOperation(
		domain.BulkActionMove,
		[]string{"PROJ-1", "PROJ-2"},
		map[string]interface{}{"parent": "PROJ-200"},
	)

	ctx := context.Background()
	result, err := service.ExecuteOperation(ctx, op, nil)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.SuccessCount != 2 {
		t.Errorf("expected SuccessCount=2, got %d", result.SuccessCount)
	}

	if result.FailureCount != 0 {
		t.Errorf("expected FailureCount=0, got %d", result.FailureCount)
	}

	// Verify parent was updated
	for _, call := range adapter.updateTicketCalls {
		if call.CustomFields["Parent"] != "PROJ-200" {
			t.Errorf("expected Parent=PROJ-200, got %s", call.CustomFields["Parent"])
		}
	}
}

// TestExecuteOperation_MoveWithoutParent verifies move requires parent field.
func TestExecuteOperation_MoveWithoutParent(t *testing.T) {
	adapter := newMockJiraAdapter()
	service := NewBulkOperationService(adapter)

	op := domain.NewBulkOperation(
		domain.BulkActionMove,
		[]string{"PROJ-1"},
		map[string]interface{}{"status": "Done"}, // Missing "parent" field
	)

	ctx := context.Background()
	_, err := service.ExecuteOperation(ctx, op, nil)

	if err == nil {
		t.Error("expected error for move without parent field, got nil")
	}
}

// TestExecuteOperation_DeleteNotSupported verifies delete returns error.
func TestExecuteOperation_DeleteNotSupported(t *testing.T) {
	adapter := newMockJiraAdapter()
	service := NewBulkOperationService(adapter)

	op := domain.NewBulkOperation(
		domain.BulkActionDelete,
		[]string{"PROJ-1", "PROJ-2"},
		nil,
	)

	ctx := context.Background()
	result, err := service.ExecuteOperation(ctx, op, nil)

	// Should return error indicating delete not supported
	if err == nil {
		t.Error("expected error for unsupported delete operation, got nil")
	}

	// All tickets should fail
	if result.FailureCount != 2 {
		t.Errorf("expected FailureCount=2, got %d", result.FailureCount)
	}
}

// TestExecuteOperation_SearchFailure verifies handling of search failures.
func TestExecuteOperation_SearchFailure(t *testing.T) {
	adapter := newMockJiraAdapter()
	adapter.searchTicketsError = fmt.Errorf("network error")

	service := NewBulkOperationService(adapter)

	op := domain.NewBulkOperation(
		domain.BulkActionUpdate,
		[]string{"PROJ-1"},
		map[string]interface{}{"status": "Done"},
	)

	ctx := context.Background()
	result, err := service.ExecuteOperation(ctx, op, nil)

	// Should handle search failure gracefully
	if err == nil {
		t.Error("expected error when search fails, got nil")
	}

	if result.FailureCount != 1 {
		t.Errorf("expected FailureCount=1, got %d", result.FailureCount)
	}
}

// TestExecuteOperation_TicketNotFound verifies handling when ticket doesn't exist.
func TestExecuteOperation_TicketNotFound(t *testing.T) {
	adapter := newMockJiraAdapter()
	adapter.withTickets([]domain.Ticket{}) // Empty results

	service := NewBulkOperationService(adapter)

	op := domain.NewBulkOperation(
		domain.BulkActionUpdate,
		[]string{"PROJ-999"},
		map[string]interface{}{"status": "Done"},
	)

	ctx := context.Background()
	result, err := service.ExecuteOperation(ctx, op, nil)

	// Should handle missing ticket gracefully
	if err == nil {
		t.Error("expected error when ticket not found, got nil")
	}

	if result.FailureCount != 1 {
		t.Errorf("expected FailureCount=1, got %d", result.FailureCount)
	}

	// Verify error message mentions ticket not found
	if errMsg, ok := result.Errors["PROJ-999"]; !ok {
		t.Error("expected error message for PROJ-999")
	} else if errMsg == "" {
		t.Error("expected non-empty error message")
	}
}

// TestExecuteOperation_LargeBatch verifies handling of large batches.
func TestExecuteOperation_LargeBatch(t *testing.T) {
	adapter := newMockJiraAdapter()

	// Create 50 tickets
	tickets := make([]domain.Ticket, 50)
	ticketIDs := make([]string, 50)
	for i := 0; i < 50; i++ {
		id := fmt.Sprintf("PROJ-%d", i+1)
		tickets[i] = domain.Ticket{
			JiraID:       id,
			Title:        fmt.Sprintf("Ticket %d", i+1),
			CustomFields: map[string]string{"status": "In Progress"},
		}
		ticketIDs[i] = id
	}
	adapter.searchTicketsResults = tickets

	service := NewBulkOperationService(adapter)

	op := domain.NewBulkOperation(
		domain.BulkActionUpdate,
		ticketIDs,
		map[string]interface{}{"status": "Done"},
	)

	ctx := context.Background()
	result, err := service.ExecuteOperation(ctx, op, nil)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.SuccessCount != 50 {
		t.Errorf("expected SuccessCount=50, got %d", result.SuccessCount)
	}

	if result.FailureCount != 0 {
		t.Errorf("expected FailureCount=0, got %d", result.FailureCount)
	}
}

// TestExecuteOperation_ContextTimeout verifies timeout handling.
func TestExecuteOperation_ContextTimeout(t *testing.T) {
	adapter := newMockJiraAdapter()

	tickets := make([]domain.Ticket, 10)
	ticketIDs := make([]string, 10)
	for i := 0; i < 10; i++ {
		id := fmt.Sprintf("PROJ-%d", i+1)
		tickets[i] = domain.Ticket{
			JiraID:       id,
			Title:        fmt.Sprintf("Ticket %d", i+1),
			CustomFields: map[string]string{},
		}
		ticketIDs[i] = id
	}
	adapter.searchTicketsResults = tickets

	// Make SearchTickets slow
	searchCalls := 0
	adapter.searchTicketsFunc = func(projectKey string, jql string) ([]domain.Ticket, error) {
		adapter.searchTicketCalls = append(adapter.searchTicketCalls, jql)
		searchCalls++
		time.Sleep(10 * time.Millisecond) // Simulate slow search
		return adapter.searchTicketsResults, nil
	}

	service := NewBulkOperationService(adapter)

	op := domain.NewBulkOperation(
		domain.BulkActionUpdate,
		ticketIDs,
		map[string]interface{}{"status": "Done"},
	)

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	result, err := service.ExecuteOperation(ctx, op, nil)

	// Should return timeout error
	if err == nil {
		t.Error("expected error for timeout, got nil")
	}

	// Should have partial results
	if result.SuccessCount+result.FailureCount >= 10 {
		t.Errorf("expected operation to be interrupted, but processed all tickets")
	}
}

// TestExecuteOperation_ProgressCallback verifies callback invocation.
func TestExecuteOperation_ProgressCallback(t *testing.T) {
	mockAdapter := &mockJiraAdapter{
		searchTicketsResults: []domain.Ticket{
			{JiraID: "PROJ-1", Title: "Test 1", CustomFields: map[string]string{}},
			{JiraID: "PROJ-2", Title: "Test 2", CustomFields: map[string]string{}},
			{JiraID: "PROJ-3", Title: "Test 3", CustomFields: map[string]string{}},
		},
	}

	service := NewBulkOperationService(mockAdapter)

	op := domain.NewBulkOperation(
		domain.BulkActionUpdate,
		[]string{"PROJ-1", "PROJ-2", "PROJ-3"},
		map[string]interface{}{"status": "Done"},
	)

	// Track callback invocations
	var callbackInvocations []string
	var callbackSuccesses []bool
	var callbackErrors []error

	progress := func(ticketID string, success bool, err error) {
		callbackInvocations = append(callbackInvocations, ticketID)
		callbackSuccesses = append(callbackSuccesses, success)
		callbackErrors = append(callbackErrors, err)
	}

	result, err := service.ExecuteOperation(context.Background(), op, progress)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify callback was invoked for each ticket
	if len(callbackInvocations) != 3 {
		t.Errorf("expected 3 callback invocations, got %d", len(callbackInvocations))
	}

	// Verify all callbacks reported success
	for i, success := range callbackSuccesses {
		if !success {
			t.Errorf("callback[%d] expected success=true, got false", i)
		}
		if callbackErrors[i] != nil {
			t.Errorf("callback[%d] expected nil error, got %v", i, callbackErrors[i])
		}
	}

	// Verify result matches callback invocations
	if result.SuccessCount != len(callbackInvocations) {
		t.Errorf("result.SuccessCount (%d) != callback invocations (%d)", result.SuccessCount, len(callbackInvocations))
	}
}

// TestExecuteOperation_NilProgressCallback verifies nil callback doesn't panic.
func TestExecuteOperation_NilProgressCallback(t *testing.T) {
	mockAdapter := &mockJiraAdapter{
		searchTicketsResults: []domain.Ticket{
			{JiraID: "PROJ-1", Title: "Test 1", CustomFields: map[string]string{}},
		},
	}

	service := NewBulkOperationService(mockAdapter)

	op := domain.NewBulkOperation(
		domain.BulkActionUpdate,
		[]string{"PROJ-1"},
		map[string]interface{}{"status": "Done"},
	)

	// Pass nil callback (should not panic)
	result, err := service.ExecuteOperation(context.Background(), op, nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.SuccessCount != 1 {
		t.Errorf("expected 1 success, got %d", result.SuccessCount)
	}
}

// TestRecordSuccess verifies success tracking.
func TestRecordSuccess(t *testing.T) {
	service := &BulkOperationServiceImpl{}
	result := domain.NewBulkOperationResult()

	service.recordSuccess(result, "PROJ-1")
	service.recordSuccess(result, "PROJ-2")

	if result.SuccessCount != 2 {
		t.Errorf("expected SuccessCount=2, got %d", result.SuccessCount)
	}

	if len(result.SuccessfulTickets) != 2 {
		t.Errorf("expected 2 successful tickets, got %d", len(result.SuccessfulTickets))
	}
}

// TestRecordFailure verifies failure tracking.
func TestRecordFailure(t *testing.T) {
	service := &BulkOperationServiceImpl{}
	result := domain.NewBulkOperationResult()

	err1 := fmt.Errorf("error 1")
	err2 := fmt.Errorf("error 2")

	service.recordFailure(result, "PROJ-1", err1)
	service.recordFailure(result, "PROJ-2", err2)

	if result.FailureCount != 2 {
		t.Errorf("expected FailureCount=2, got %d", result.FailureCount)
	}

	if len(result.FailedTickets) != 2 {
		t.Errorf("expected 2 failed tickets, got %d", len(result.FailedTickets))
	}

	if len(result.Errors) != 2 {
		t.Errorf("expected 2 error messages, got %d", len(result.Errors))
	}
}
