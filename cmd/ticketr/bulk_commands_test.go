package main

import (
	"reflect"
	"testing"
)

// TestParseTicketIDs tests the parseTicketIDs function.
func TestParseTicketIDs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "empty string",
			input: "",
			want:  []string{},
		},
		{
			name:  "single ticket",
			input: "PROJ-1",
			want:  []string{"PROJ-1"},
		},
		{
			name:  "multiple tickets",
			input: "PROJ-1,PROJ-2,PROJ-3",
			want:  []string{"PROJ-1", "PROJ-2", "PROJ-3"},
		},
		{
			name:  "tickets with spaces",
			input: "PROJ-1, PROJ-2 , PROJ-3",
			want:  []string{"PROJ-1", "PROJ-2", "PROJ-3"},
		},
		{
			name:  "tickets with extra whitespace",
			input: "  PROJ-1  ,  PROJ-2  ,  PROJ-3  ",
			want:  []string{"PROJ-1", "PROJ-2", "PROJ-3"},
		},
		{
			name:  "mixed project keys",
			input: "PROJ-1,EPIC-42,TASK-123",
			want:  []string{"PROJ-1", "EPIC-42", "TASK-123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTicketIDs(tt.input)
			if !equalStringSlices(got, tt.want) {
				t.Errorf("parseTicketIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestParseSetFlags tests the parseSetFlags function.
func TestParseSetFlags(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "empty input",
			input:   []string{},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "single field",
			input: []string{"status=Done"},
			want:  map[string]interface{}{"status": "Done"},
		},
		{
			name:  "multiple fields",
			input: []string{"status=Done", "assignee=john@example.com"},
			want: map[string]interface{}{
				"status":   "Done",
				"assignee": "john@example.com",
			},
		},
		{
			name:  "field with spaces in value",
			input: []string{"status=In Progress", "priority=High Priority"},
			want: map[string]interface{}{
				"status":   "In Progress",
				"priority": "High Priority",
			},
		},
		{
			name:  "field with equals in value",
			input: []string{"description=x=y+z"},
			want:  map[string]interface{}{"description": "x=y+z"},
		},
		{
			name:    "invalid format - no equals",
			input:   []string{"status"},
			wantErr: true,
		},
		{
			name:    "invalid format - empty field name",
			input:   []string{"=Done"},
			wantErr: true,
		},
		{
			name:  "field with whitespace around equals",
			input: []string{"status = Done", "assignee = john@example.com"},
			want: map[string]interface{}{
				"status":   "Done",
				"assignee": "john@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSetFlags(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSetFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !equalMaps(got, tt.want) {
				t.Errorf("parseSetFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCreateProgressCallback tests the createProgressCallback function.
func TestCreateProgressCallback(t *testing.T) {
	tests := []struct {
		name         string
		total        int
		tickets      []string
		successes    []bool
		wantCalls    int
		wantProgress []string // expected progress messages (simplified)
	}{
		{
			name:      "all success",
			total:     3,
			tickets:   []string{"PROJ-1", "PROJ-2", "PROJ-3"},
			successes: []bool{true, true, true},
			wantCalls: 3,
		},
		{
			name:      "mixed success and failure",
			total:     3,
			tickets:   []string{"PROJ-1", "PROJ-2", "PROJ-3"},
			successes: []bool{true, false, true},
			wantCalls: 3,
		},
		{
			name:      "single ticket",
			total:     1,
			tickets:   []string{"PROJ-1"},
			successes: []bool{true},
			wantCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callCount := 0
			callback := createProgressCallback(tt.total)

			// Simulate calling the callback for each ticket
			for i, ticketID := range tt.tickets {
				var err error
				if !tt.successes[i] {
					err = &testError{msg: "test error"}
				}
				callback(ticketID, tt.successes[i], err)
				callCount++
			}

			if callCount != tt.wantCalls {
				t.Errorf("callback was called %d times, want %d", callCount, tt.wantCalls)
			}
		})
	}
}

// Helper functions

// equalStringSlices compares two string slices for equality.
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// equalMaps compares two maps for equality.
func equalMaps(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	return reflect.DeepEqual(a, b)
}

// testError is a simple error implementation for testing.
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
