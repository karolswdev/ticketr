package predicates

import (
	"testing"

	"github.com/karolswdev/ticktr/internal/tui-bubbletea/actions"
)

func TestAlways(t *testing.T) {
	pred := Always()
	ctx := &actions.ActionContext{}

	if !pred(ctx) {
		t.Error("Always() should always return true")
	}
}

func TestNever(t *testing.T) {
	pred := Never()
	ctx := &actions.ActionContext{}

	if pred(ctx) {
		t.Error("Never() should always return false")
	}
}

func TestHasSelection(t *testing.T) {
	pred := HasSelection()

	tests := []struct {
		name     string
		tickets  []string
		expected bool
	}{
		{"No selection", []string{}, false},
		{"Single selection", []string{"TICKET-1"}, true},
		{"Multiple selection", []string{"TICKET-1", "TICKET-2"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &actions.ActionContext{
				SelectedTickets: tt.tickets,
			}
			result := pred(ctx)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHasSingleSelection(t *testing.T) {
	pred := HasSingleSelection()

	tests := []struct {
		name     string
		tickets  []string
		expected bool
	}{
		{"No selection", []string{}, false},
		{"Single selection", []string{"TICKET-1"}, true},
		{"Multiple selection", []string{"TICKET-1", "TICKET-2"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &actions.ActionContext{
				SelectedTickets: tt.tickets,
			}
			result := pred(ctx)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHasMultipleSelection(t *testing.T) {
	pred := HasMultipleSelection()

	tests := []struct {
		name     string
		tickets  []string
		expected bool
	}{
		{"No selection", []string{}, false},
		{"Single selection", []string{"TICKET-1"}, false},
		{"Multiple selection", []string{"TICKET-1", "TICKET-2"}, true},
		{"Three tickets", []string{"TICKET-1", "TICKET-2", "TICKET-3"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &actions.ActionContext{
				SelectedTickets: tt.tickets,
			}
			result := pred(ctx)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsInWorkspace(t *testing.T) {
	pred := IsInWorkspace()

	tests := []struct {
		name      string
		workspace *actions.WorkspaceState
		expected  bool
	}{
		{"No workspace", nil, false},
		{"With workspace", &actions.WorkspaceState{ID: "ws-1", Name: "Test"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &actions.ActionContext{
				SelectedWorkspace: tt.workspace,
			}
			result := pred(ctx)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHasUnsavedChanges(t *testing.T) {
	pred := HasUnsavedChanges()

	tests := []struct {
		name     string
		dirty    bool
		expected bool
	}{
		{"No unsaved changes", false, false},
		{"Has unsaved changes", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &actions.ActionContext{
				HasUnsavedChanges: tt.dirty,
			}
			result := pred(ctx)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsOnline(t *testing.T) {
	pred := IsOnline()

	tests := []struct {
		name     string
		offline  bool
		expected bool
	}{
		{"Online", false, true},
		{"Offline", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &actions.ActionContext{
				IsOffline: tt.offline,
			}
			result := pred(ctx)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestNot(t *testing.T) {
	// Not(Always) should be Never
	pred := Not(Always())
	ctx := &actions.ActionContext{}

	if pred(ctx) {
		t.Error("Not(Always()) should return false")
	}

	// Not(Never) should be Always
	pred = Not(Never())
	if !pred(ctx) {
		t.Error("Not(Never()) should return true")
	}

	// Not(HasSelection) with no selection
	pred = Not(HasSelection())
	ctx.SelectedTickets = []string{}
	if !pred(ctx) {
		t.Error("Not(HasSelection()) should return true when no selection")
	}

	// Not(HasSelection) with selection
	ctx.SelectedTickets = []string{"TICKET-1"}
	if pred(ctx) {
		t.Error("Not(HasSelection()) should return false when selection exists")
	}
}

func TestAnd(t *testing.T) {
	tests := []struct {
		name       string
		predicates []actions.PredicateFunc
		ctx        *actions.ActionContext
		expected   bool
	}{
		{
			name:       "All true",
			predicates: []actions.PredicateFunc{Always(), Always(), Always()},
			ctx:        &actions.ActionContext{},
			expected:   true,
		},
		{
			name:       "All false",
			predicates: []actions.PredicateFunc{Never(), Never(), Never()},
			ctx:        &actions.ActionContext{},
			expected:   false,
		},
		{
			name:       "Mixed with one false",
			predicates: []actions.PredicateFunc{Always(), Never(), Always()},
			ctx:        &actions.ActionContext{},
			expected:   false,
		},
		{
			name:       "HasSelection AND HasSingleSelection with single",
			predicates: []actions.PredicateFunc{HasSelection(), HasSingleSelection()},
			ctx: &actions.ActionContext{
				SelectedTickets: []string{"TICKET-1"},
			},
			expected: true,
		},
		{
			name:       "HasSelection AND HasSingleSelection with multiple",
			predicates: []actions.PredicateFunc{HasSelection(), HasSingleSelection()},
			ctx: &actions.ActionContext{
				SelectedTickets: []string{"TICKET-1", "TICKET-2"},
			},
			expected: false,
		},
		{
			name:       "Empty predicates list",
			predicates: []actions.PredicateFunc{},
			ctx:        &actions.ActionContext{},
			expected:   true, // Empty AND is true
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pred := And(tt.predicates...)
			result := pred(tt.ctx)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestOr(t *testing.T) {
	tests := []struct {
		name       string
		predicates []actions.PredicateFunc
		ctx        *actions.ActionContext
		expected   bool
	}{
		{
			name:       "All true",
			predicates: []actions.PredicateFunc{Always(), Always(), Always()},
			ctx:        &actions.ActionContext{},
			expected:   true,
		},
		{
			name:       "All false",
			predicates: []actions.PredicateFunc{Never(), Never(), Never()},
			ctx:        &actions.ActionContext{},
			expected:   false,
		},
		{
			name:       "Mixed with one true",
			predicates: []actions.PredicateFunc{Never(), Always(), Never()},
			ctx:        &actions.ActionContext{},
			expected:   true,
		},
		{
			name:       "HasSingleSelection OR HasMultipleSelection with single",
			predicates: []actions.PredicateFunc{HasSingleSelection(), HasMultipleSelection()},
			ctx: &actions.ActionContext{
				SelectedTickets: []string{"TICKET-1"},
			},
			expected: true,
		},
		{
			name:       "HasSingleSelection OR HasMultipleSelection with multiple",
			predicates: []actions.PredicateFunc{HasSingleSelection(), HasMultipleSelection()},
			ctx: &actions.ActionContext{
				SelectedTickets: []string{"TICKET-1", "TICKET-2"},
			},
			expected: true,
		},
		{
			name:       "HasSingleSelection OR HasMultipleSelection with none",
			predicates: []actions.PredicateFunc{HasSingleSelection(), HasMultipleSelection()},
			ctx: &actions.ActionContext{
				SelectedTickets: []string{},
			},
			expected: false,
		},
		{
			name:       "Empty predicates list",
			predicates: []actions.PredicateFunc{},
			ctx:        &actions.ActionContext{},
			expected:   false, // Empty OR is false
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pred := Or(tt.predicates...)
			result := pred(tt.ctx)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestComplexPredicateComposition(t *testing.T) {
	// Test: (HasSelection AND IsOnline) OR HasUnsavedChanges
	pred := Or(
		And(HasSelection(), IsOnline()),
		HasUnsavedChanges(),
	)

	tests := []struct {
		name     string
		ctx      *actions.ActionContext
		expected bool
	}{
		{
			name: "Has selection, online, no unsaved",
			ctx: &actions.ActionContext{
				SelectedTickets:   []string{"TICKET-1"},
				IsOffline:         false,
				HasUnsavedChanges: false,
			},
			expected: true,
		},
		{
			name: "Has selection, offline, no unsaved",
			ctx: &actions.ActionContext{
				SelectedTickets:   []string{"TICKET-1"},
				IsOffline:         true,
				HasUnsavedChanges: false,
			},
			expected: false,
		},
		{
			name: "No selection, offline, has unsaved",
			ctx: &actions.ActionContext{
				SelectedTickets:   []string{},
				IsOffline:         true,
				HasUnsavedChanges: true,
			},
			expected: true,
		},
		{
			name: "No selection, online, no unsaved",
			ctx: &actions.ActionContext{
				SelectedTickets:   []string{},
				IsOffline:         false,
				HasUnsavedChanges: false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pred(tt.ctx)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestNestedNotComposition(t *testing.T) {
	// Test: NOT(HasSelection AND IsOffline)
	// This means: "No selection OR online"
	pred := Not(And(HasSelection(), Not(IsOnline())))

	tests := []struct {
		name     string
		ctx      *actions.ActionContext
		expected bool
	}{
		{
			name: "No selection, offline",
			ctx: &actions.ActionContext{
				SelectedTickets: []string{},
				IsOffline:       true,
			},
			expected: true, // No selection, so outer AND fails, NOT makes it true
		},
		{
			name: "Has selection, online",
			ctx: &actions.ActionContext{
				SelectedTickets: []string{"TICKET-1"},
				IsOffline:       false,
			},
			expected: true, // Online (NOT offline), so inner AND fails, NOT makes it true
		},
		{
			name: "Has selection, offline",
			ctx: &actions.ActionContext{
				SelectedTickets: []string{"TICKET-1"},
				IsOffline:       true,
			},
			expected: false, // Has selection AND offline, NOT makes it false
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pred(tt.ctx)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
