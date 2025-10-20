package domain

import (
	"fmt"
	"strings"
	"testing"
)

func TestBulkOperation_Validate(t *testing.T) {
	tests := []struct {
		name      string
		operation *BulkOperation
		wantErr   bool
		errorMsg  string
	}{
		{
			name: "valid update operation",
			operation: &BulkOperation{
				Action:    BulkActionUpdate,
				TicketIDs: []string{"PROJ-123", "PROJ-124"},
				Changes: map[string]interface{}{
					"status":   "In Progress",
					"assignee": "john.doe@example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "valid move operation",
			operation: &BulkOperation{
				Action:    BulkActionMove,
				TicketIDs: []string{"PROJ-123"},
				Changes:   nil,
			},
			wantErr: false,
		},
		{
			name: "valid delete operation",
			operation: &BulkOperation{
				Action:    BulkActionDelete,
				TicketIDs: []string{"PROJ-123", "PROJ-124", "PROJ-125"},
				Changes:   nil,
			},
			wantErr: false,
		},
		{
			name: "invalid action",
			operation: &BulkOperation{
				Action:    "invalid_action",
				TicketIDs: []string{"PROJ-123"},
				Changes:   nil,
			},
			wantErr:  true,
			errorMsg: "invalid action: must be one of 'update', 'move', or 'delete'",
		},
		{
			name: "empty ticket IDs",
			operation: &BulkOperation{
				Action:    BulkActionDelete,
				TicketIDs: []string{},
				Changes:   nil,
			},
			wantErr:  true,
			errorMsg: "ticket_ids cannot be empty: must contain at least 1 ticket",
		},
		{
			name: "too many ticket IDs",
			operation: &BulkOperation{
				Action:    BulkActionDelete,
				TicketIDs: make101TicketIDs(),
				Changes:   nil,
			},
			wantErr:  true,
			errorMsg: "ticket_ids cannot exceed 100 tickets: found 101",
		},
		{
			name: "empty string in ticket IDs",
			operation: &BulkOperation{
				Action:    BulkActionDelete,
				TicketIDs: []string{"PROJ-123", "", "PROJ-125"},
				Changes:   nil,
			},
			wantErr:  true,
			errorMsg: "ticket_ids[1] cannot be an empty string",
		},
		{
			name: "update action without changes",
			operation: &BulkOperation{
				Action:    BulkActionUpdate,
				TicketIDs: []string{"PROJ-123"},
				Changes:   nil,
			},
			wantErr:  true,
			errorMsg: "changes are required for 'update' action",
		},
		{
			name: "update action with empty changes map",
			operation: &BulkOperation{
				Action:    BulkActionUpdate,
				TicketIDs: []string{"PROJ-123"},
				Changes:   map[string]interface{}{},
			},
			wantErr:  true,
			errorMsg: "changes are required for 'update' action",
		},
		{
			name: "exactly 1 ticket (edge case minimum)",
			operation: &BulkOperation{
				Action:    BulkActionDelete,
				TicketIDs: []string{"PROJ-123"},
				Changes:   nil,
			},
			wantErr: false,
		},
		{
			name: "exactly 100 tickets (edge case maximum)",
			operation: &BulkOperation{
				Action:    BulkActionDelete,
				TicketIDs: make100TicketIDs(),
				Changes:   nil,
			},
			wantErr: false,
		},
		{
			name: "move action with changes (allowed but ignored)",
			operation: &BulkOperation{
				Action:    BulkActionMove,
				TicketIDs: []string{"PROJ-123"},
				Changes: map[string]interface{}{
					"status": "Done",
				},
			},
			wantErr: false,
		},
		{
			name: "delete action with changes (allowed but ignored)",
			operation: &BulkOperation{
				Action:    BulkActionDelete,
				TicketIDs: []string{"PROJ-123"},
				Changes: map[string]interface{}{
					"status": "Done",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestNewBulkOperation(t *testing.T) {
	tests := []struct {
		name      string
		action    BulkOperationAction
		ticketIDs []string
		changes   map[string]interface{}
	}{
		{
			name:      "create update operation",
			action:    BulkActionUpdate,
			ticketIDs: []string{"PROJ-1", "PROJ-2"},
			changes:   map[string]interface{}{"status": "Done"},
		},
		{
			name:      "create move operation",
			action:    BulkActionMove,
			ticketIDs: []string{"PROJ-1"},
			changes:   nil,
		},
		{
			name:      "create delete operation",
			action:    BulkActionDelete,
			ticketIDs: []string{"PROJ-1", "PROJ-2", "PROJ-3"},
			changes:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewBulkOperation(tt.action, tt.ticketIDs, tt.changes)

			if op == nil {
				t.Fatal("NewBulkOperation() returned nil")
			}

			if op.Action != tt.action {
				t.Errorf("Action = %v, want %v", op.Action, tt.action)
			}

			if len(op.TicketIDs) != len(tt.ticketIDs) {
				t.Errorf("TicketIDs length = %d, want %d", len(op.TicketIDs), len(tt.ticketIDs))
			}

			if tt.changes != nil && len(op.Changes) != len(tt.changes) {
				t.Errorf("Changes length = %d, want %d", len(op.Changes), len(tt.changes))
			}
		})
	}
}

func TestIsValidAction(t *testing.T) {
	tests := []struct {
		name   string
		action string
		want   bool
	}{
		{"update action", "update", true},
		{"move action", "move", true},
		{"delete action", "delete", true},
		{"invalid action", "invalid", false},
		{"empty string", "", false},
		{"uppercase UPDATE", "UPDATE", false},
		{"mixed case Update", "Update", false},
		{"create action", "create", false},
		{"archive action", "archive", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidAction(tt.action)
			if got != tt.want {
				t.Errorf("IsValidAction(%q) = %v, want %v", tt.action, got, tt.want)
			}
		})
	}
}

func TestBulkOperation_TicketCount(t *testing.T) {
	tests := []struct {
		name      string
		ticketIDs []string
		want      int
	}{
		{"zero tickets", []string{}, 0},
		{"one ticket", []string{"PROJ-1"}, 1},
		{"multiple tickets", []string{"PROJ-1", "PROJ-2", "PROJ-3"}, 3},
		{"max tickets", make100TicketIDs(), 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &BulkOperation{
				Action:    BulkActionDelete,
				TicketIDs: tt.ticketIDs,
			}

			got := op.TicketCount()
			if got != tt.want {
				t.Errorf("TicketCount() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestNewBulkOperationResult(t *testing.T) {
	result := NewBulkOperationResult()

	if result == nil {
		t.Fatal("NewBulkOperationResult() returned nil")
	}

	if result.SuccessCount != 0 {
		t.Errorf("SuccessCount = %d, want 0", result.SuccessCount)
	}

	if result.FailureCount != 0 {
		t.Errorf("FailureCount = %d, want 0", result.FailureCount)
	}

	if result.Errors == nil {
		t.Error("Errors map should be initialized, got nil")
	}

	if result.SuccessfulTickets == nil {
		t.Error("SuccessfulTickets slice should be initialized, got nil")
	}

	if result.FailedTickets == nil {
		t.Error("FailedTickets slice should be initialized, got nil")
	}

	if len(result.Errors) != 0 {
		t.Errorf("Errors map should be empty, got length %d", len(result.Errors))
	}

	if len(result.SuccessfulTickets) != 0 {
		t.Errorf("SuccessfulTickets should be empty, got length %d", len(result.SuccessfulTickets))
	}

	if len(result.FailedTickets) != 0 {
		t.Errorf("FailedTickets should be empty, got length %d", len(result.FailedTickets))
	}
}

func TestBulkOperation_validateAction(t *testing.T) {
	tests := []struct {
		name    string
		action  BulkOperationAction
		wantErr bool
	}{
		{"valid update", BulkActionUpdate, false},
		{"valid move", BulkActionMove, false},
		{"valid delete", BulkActionDelete, false},
		{"invalid empty", "", true},
		{"invalid custom", "custom", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &BulkOperation{
				Action: tt.action,
			}

			err := op.validateAction()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBulkOperation_validateTicketIDs(t *testing.T) {
	tests := []struct {
		name      string
		ticketIDs []string
		wantErr   bool
		errorMsg  string
	}{
		{
			name:      "valid single ticket",
			ticketIDs: []string{"PROJ-123"},
			wantErr:   false,
		},
		{
			name:      "valid multiple tickets",
			ticketIDs: []string{"PROJ-1", "PROJ-2", "PROJ-3"},
			wantErr:   false,
		},
		{
			name:      "valid 100 tickets",
			ticketIDs: make100TicketIDs(),
			wantErr:   false,
		},
		{
			name:      "empty list",
			ticketIDs: []string{},
			wantErr:   true,
			errorMsg:  "ticket_ids cannot be empty: must contain at least 1 ticket",
		},
		{
			name:      "nil list",
			ticketIDs: nil,
			wantErr:   true,
			errorMsg:  "ticket_ids cannot be empty: must contain at least 1 ticket",
		},
		{
			name:      "too many tickets",
			ticketIDs: make101TicketIDs(),
			wantErr:   true,
			errorMsg:  "ticket_ids cannot exceed 100 tickets: found 101",
		},
		{
			name:      "empty string at start",
			ticketIDs: []string{"", "PROJ-2"},
			wantErr:   true,
			errorMsg:  "ticket_ids[0] cannot be an empty string",
		},
		{
			name:      "empty string in middle",
			ticketIDs: []string{"PROJ-1", "", "PROJ-3"},
			wantErr:   true,
			errorMsg:  "ticket_ids[1] cannot be an empty string",
		},
		{
			name:      "empty string at end",
			ticketIDs: []string{"PROJ-1", "PROJ-2", ""},
			wantErr:   true,
			errorMsg:  "ticket_ids[2] cannot be an empty string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &BulkOperation{
				TicketIDs: tt.ticketIDs,
			}

			err := op.validateTicketIDs()
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateTicketIDs() expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("validateTicketIDs() error = %v, want %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateTicketIDs() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestBulkOperation_validateChanges(t *testing.T) {
	tests := []struct {
		name     string
		action   BulkOperationAction
		changes  map[string]interface{}
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "update with valid changes",
			action:  BulkActionUpdate,
			changes: map[string]interface{}{"status": "Done"},
			wantErr: false,
		},
		{
			name:     "update with nil changes",
			action:   BulkActionUpdate,
			changes:  nil,
			wantErr:  true,
			errorMsg: "changes are required for 'update' action",
		},
		{
			name:     "update with empty changes",
			action:   BulkActionUpdate,
			changes:  map[string]interface{}{},
			wantErr:  true,
			errorMsg: "changes are required for 'update' action",
		},
		{
			name:    "move with nil changes (valid)",
			action:  BulkActionMove,
			changes: nil,
			wantErr: false,
		},
		{
			name:    "move with changes (valid, ignored)",
			action:  BulkActionMove,
			changes: map[string]interface{}{"status": "Done"},
			wantErr: false,
		},
		{
			name:    "delete with nil changes (valid)",
			action:  BulkActionDelete,
			changes: nil,
			wantErr: false,
		},
		{
			name:    "delete with changes (valid, ignored)",
			action:  BulkActionDelete,
			changes: map[string]interface{}{"status": "Done"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &BulkOperation{
				Action:  tt.action,
				Changes: tt.changes,
			}

			err := op.validateChanges()
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateChanges() expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("validateChanges() error = %v, want %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateChanges() unexpected error = %v", err)
				}
			}
		})
	}
}

// Helper functions

// make100TicketIDs generates exactly 100 ticket IDs for edge case testing.
func make100TicketIDs() []string {
	ids := make([]string, 100)
	for i := 0; i < 100; i++ {
		ids[i] = fmt.Sprintf("PROJ-%d", i+1)
	}
	return ids
}

// make101TicketIDs generates 101 ticket IDs to test the upper limit validation.
func make101TicketIDs() []string {
	ids := make([]string, 101)
	for i := 0; i < 101; i++ {
		ids[i] = fmt.Sprintf("PROJ-%d", i+1)
	}
	return ids
}

func TestBulkOperation_ValidateJiraIDFormat(t *testing.T) {
	tests := []struct {
		name      string
		operation *BulkOperation
		wantErr   bool
		errorMsg  string
	}{
		{
			name: "valid Jira ID format",
			operation: &BulkOperation{
				Action:    BulkActionUpdate,
				TicketIDs: []string{"PROJ-123", "ABC-1", "TEAM-9999"},
				Changes:   map[string]interface{}{"status": "Done"},
			},
			wantErr: false,
		},
		{
			name: "JQL injection attempt with quotes",
			operation: &BulkOperation{
				Action:    BulkActionUpdate,
				TicketIDs: []string{`PROJ-1" OR key != "`},
				Changes:   map[string]interface{}{"status": "Done"},
			},
			wantErr:  true,
			errorMsg: "invalid Jira ID format",
		},
		{
			name: "JQL injection with OR operator",
			operation: &BulkOperation{
				Action:    BulkActionUpdate,
				TicketIDs: []string{"PROJ-1 OR 1=1"},
				Changes:   map[string]interface{}{"status": "Done"},
			},
			wantErr:  true,
			errorMsg: "invalid Jira ID format",
		},
		{
			name: "lowercase project key",
			operation: &BulkOperation{
				Action:    BulkActionUpdate,
				TicketIDs: []string{"proj-123"},
				Changes:   map[string]interface{}{"status": "Done"},
			},
			wantErr:  true,
			errorMsg: "invalid Jira ID format",
		},
		{
			name: "missing hyphen",
			operation: &BulkOperation{
				Action:    BulkActionUpdate,
				TicketIDs: []string{"PROJ123"},
				Changes:   map[string]interface{}{"status": "Done"},
			},
			wantErr:  true,
			errorMsg: "invalid Jira ID format",
		},
		{
			name: "special characters",
			operation: &BulkOperation{
				Action:    BulkActionUpdate,
				TicketIDs: []string{"PROJ-123; DROP TABLE tickets;"},
				Changes:   map[string]interface{}{"status": "Done"},
			},
			wantErr:  true,
			errorMsg: "invalid Jira ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errorMsg)
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

// Benchmark tests

func BenchmarkBulkOperation_Validate(b *testing.B) {
	op := &BulkOperation{
		Action:    BulkActionUpdate,
		TicketIDs: []string{"PROJ-1", "PROJ-2", "PROJ-3"},
		Changes:   map[string]interface{}{"status": "Done"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = op.Validate()
	}
}

func BenchmarkBulkOperation_Validate_MaxTickets(b *testing.B) {
	op := &BulkOperation{
		Action:    BulkActionDelete,
		TicketIDs: make100TicketIDs(),
		Changes:   nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = op.Validate()
	}
}

func BenchmarkNewBulkOperation(b *testing.B) {
	ticketIDs := []string{"PROJ-1", "PROJ-2", "PROJ-3"}
	changes := map[string]interface{}{"status": "Done"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBulkOperation(BulkActionUpdate, ticketIDs, changes)
	}
}
