package actions

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// Test helper predicates (to avoid import cycle with predicates package)
func testAlways() PredicateFunc {
	return func(ctx *ActionContext) bool {
		return true
	}
}

func testHasSingleSelection() PredicateFunc {
	return func(ctx *ActionContext) bool {
		return len(ctx.SelectedTickets) == 1
	}
}

func TestRegistryRegister(t *testing.T) {
	reg := NewRegistry()

	action := &Action{
		ID:       "test.action",
		Name:     "Test Action",
		Contexts: []Context{ContextTicketTree},
		Execute: func(ctx *ActionContext) tea.Cmd {
			return nil
		},
	}

	err := reg.Register(action)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check retrieval
	retrieved, exists := reg.Get("test.action")
	if !exists {
		t.Error("Expected action to exist")
	}
	if retrieved.Name != "Test Action" {
		t.Errorf("Expected name 'Test Action', got %s", retrieved.Name)
	}
}

func TestRegistryRegisterNoID(t *testing.T) {
	reg := NewRegistry()

	action := &Action{
		Name: "No ID",
		Execute: func(ctx *ActionContext) tea.Cmd {
			return nil
		},
	}

	err := reg.Register(action)
	if err == nil {
		t.Error("Expected error for action without ID, got nil")
	}
}

func TestRegistryRegisterNoExecute(t *testing.T) {
	reg := NewRegistry()

	action := &Action{
		ID:   "test.no_execute",
		Name: "No Execute",
	}

	err := reg.Register(action)
	if err == nil {
		t.Error("Expected error for action without Execute function, got nil")
	}
}

func TestRegistryDuplicateID(t *testing.T) {
	reg := NewRegistry()

	action1 := &Action{
		ID:      "test.duplicate",
		Name:    "First",
		Execute: func(ctx *ActionContext) tea.Cmd { return nil },
	}

	action2 := &Action{
		ID:      "test.duplicate",
		Name:    "Second",
		Execute: func(ctx *ActionContext) tea.Cmd { return nil },
	}

	err := reg.Register(action1)
	if err != nil {
		t.Fatalf("Expected no error registering first action, got %v", err)
	}

	err = reg.Register(action2)
	if err == nil {
		t.Error("Expected error for duplicate ID, got nil")
	}
}

func TestActionsForContext(t *testing.T) {
	reg := NewRegistry()

	treeAction := &Action{
		ID:        "tree.action",
		Name:      "Tree Action",
		Contexts:  []Context{ContextTicketTree},
		Predicate: testAlways(),
		Execute:   func(ctx *ActionContext) tea.Cmd { return nil },
	}

	detailAction := &Action{
		ID:        "detail.action",
		Name:      "Detail Action",
		Contexts:  []Context{ContextTicketDetail},
		Predicate: testAlways(),
		Execute:   func(ctx *ActionContext) tea.Cmd { return nil },
	}

	globalAction := &Action{
		ID:        "global.action",
		Name:      "Global Action",
		Contexts:  []Context{ContextGlobal},
		Predicate: testAlways(),
		Execute:   func(ctx *ActionContext) tea.Cmd { return nil },
	}

	reg.Register(treeAction)
	reg.Register(detailAction)
	reg.Register(globalAction)

	actx := &ActionContext{}

	// Test tree context
	treeActions := reg.ActionsForContext(ContextTicketTree, actx)
	if len(treeActions) != 2 { // tree + global
		t.Errorf("Expected 2 actions for tree context, got %d", len(treeActions))
	}

	// Test detail context
	detailActions := reg.ActionsForContext(ContextTicketDetail, actx)
	if len(detailActions) != 2 { // detail + global
		t.Errorf("Expected 2 actions for detail context, got %d", len(detailActions))
	}

	// Test workspace context (no specific actions)
	workspaceActions := reg.ActionsForContext(ContextWorkspaceList, actx)
	if len(workspaceActions) != 1 { // only global
		t.Errorf("Expected 1 action for workspace context, got %d", len(workspaceActions))
	}
}

func TestActionsWithPredicates(t *testing.T) {
	reg := NewRegistry()

	// Action with predicate requiring selection
	action := &Action{
		ID:        "test.predicated",
		Name:      "Predicated Action",
		Contexts:  []Context{ContextTicketTree},
		Predicate: testHasSingleSelection(),
		Execute:   func(ctx *ActionContext) tea.Cmd { return nil },
	}

	reg.Register(action)

	// Test without selection
	actx := &ActionContext{
		SelectedTickets: []string{},
	}
	actions := reg.ActionsForContext(ContextTicketTree, actx)
	if len(actions) != 0 {
		t.Errorf("Expected 0 actions without selection, got %d", len(actions))
	}

	// Test with single selection
	actx.SelectedTickets = []string{"TICKET-1"}
	actions = reg.ActionsForContext(ContextTicketTree, actx)
	if len(actions) != 1 {
		t.Errorf("Expected 1 action with single selection, got %d", len(actions))
	}

	// Test with multiple selection
	actx.SelectedTickets = []string{"TICKET-1", "TICKET-2"}
	actions = reg.ActionsForContext(ContextTicketTree, actx)
	if len(actions) != 0 {
		t.Errorf("Expected 0 actions with multiple selection, got %d", len(actions))
	}
}

func TestActionsForKey(t *testing.T) {
	reg := NewRegistry()

	action := &Action{
		ID:       "test.keybind",
		Name:     "Keybind Action",
		Contexts: []Context{ContextTicketTree},
		Keybindings: []KeyPattern{
			{Key: "enter"},
		},
		Predicate: testAlways(),
		Execute:   func(ctx *ActionContext) tea.Cmd { return nil },
	}

	reg.Register(action)

	actx := &ActionContext{}

	// Test matching key
	actions := reg.ActionsForKey("enter", ContextTicketTree, actx)
	if len(actions) != 1 {
		t.Errorf("Expected 1 action for 'enter' key, got %d", len(actions))
	}

	// Test non-matching key
	actions = reg.ActionsForKey("x", ContextTicketTree, actx)
	if len(actions) != 0 {
		t.Errorf("Expected 0 actions for 'x' key, got %d", len(actions))
	}

	// Test wrong context
	actions = reg.ActionsForKey("enter", ContextTicketDetail, actx)
	if len(actions) != 0 {
		t.Errorf("Expected 0 actions in wrong context, got %d", len(actions))
	}
}

func TestSearch(t *testing.T) {
	reg := NewRegistry()

	action1 := &Action{
		ID:          "test.open",
		Name:        "Open Ticket",
		Description: "Open the selected ticket",
		Tags:        []string{"ticket", "open"},
		Predicate:   testAlways(),
		Execute:     func(ctx *ActionContext) tea.Cmd { return nil },
	}

	action2 := &Action{
		ID:          "test.edit",
		Name:        "Edit Ticket",
		Description: "Edit the selected ticket",
		Tags:        []string{"ticket", "edit"},
		Predicate:   testAlways(),
		Execute:     func(ctx *ActionContext) tea.Cmd { return nil },
	}

	action3 := &Action{
		ID:          "test.close",
		Name:        "Close Window",
		Description: "Close the current window",
		Tags:        []string{"window", "close"},
		Predicate:   testAlways(),
		Execute:     func(ctx *ActionContext) tea.Cmd { return nil },
	}

	reg.Register(action1)
	reg.Register(action2)
	reg.Register(action3)

	actx := &ActionContext{}

	// Search for "open"
	results := reg.Search("open", actx)
	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'open', got %d", len(results))
	}
	if len(results) > 0 && results[0].ID != "test.open" {
		t.Errorf("Expected action ID 'test.open', got %s", results[0].ID)
	}

	// Search for "ticket"
	results = reg.Search("ticket", actx)
	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'ticket', got %d", len(results))
	}

	// Search for "close"
	results = reg.Search("close", actx)
	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'close', got %d", len(results))
	}

	// Search for non-existent term
	results = reg.Search("nonexistent", actx)
	if len(results) != 0 {
		t.Errorf("Expected 0 results for 'nonexistent', got %d", len(results))
	}
}

func TestSearchWithPredicates(t *testing.T) {
	reg := NewRegistry()

	action := &Action{
		ID:          "test.predicated_search",
		Name:        "Predicated Action",
		Description: "Action with predicate",
		Predicate:   testHasSingleSelection(),
		Execute:     func(ctx *ActionContext) tea.Cmd { return nil },
	}

	reg.Register(action)

	// Search without selection (predicate fails)
	actx := &ActionContext{
		SelectedTickets: []string{},
	}
	results := reg.Search("predicated", actx)
	if len(results) != 0 {
		t.Errorf("Expected 0 results when predicate fails, got %d", len(results))
	}

	// Search with selection (predicate passes)
	actx.SelectedTickets = []string{"TICKET-1"}
	results = reg.Search("predicated", actx)
	if len(results) != 1 {
		t.Errorf("Expected 1 result when predicate passes, got %d", len(results))
	}
}

func TestUnregister(t *testing.T) {
	reg := NewRegistry()

	action := &Action{
		ID:      "test.unregister",
		Name:    "Unregister Test",
		Execute: func(ctx *ActionContext) tea.Cmd { return nil },
	}

	// Register
	err := reg.Register(action)
	if err != nil {
		t.Fatalf("Expected no error registering, got %v", err)
	}

	// Verify exists
	_, exists := reg.Get("test.unregister")
	if !exists {
		t.Error("Expected action to exist after registration")
	}

	// Unregister
	err = reg.Unregister("test.unregister")
	if err != nil {
		t.Errorf("Expected no error unregistering, got %v", err)
	}

	// Verify doesn't exist
	_, exists = reg.Get("test.unregister")
	if exists {
		t.Error("Expected action to not exist after unregistration")
	}

	// Try to unregister again (should error)
	err = reg.Unregister("test.unregister")
	if err == nil {
		t.Error("Expected error unregistering non-existent action, got nil")
	}
}

func TestAll(t *testing.T) {
	reg := NewRegistry()

	action1 := &Action{
		ID:      "test.all1",
		Name:    "All Test 1",
		Execute: func(ctx *ActionContext) tea.Cmd { return nil },
	}

	action2 := &Action{
		ID:      "test.all2",
		Name:    "All Test 2",
		Execute: func(ctx *ActionContext) tea.Cmd { return nil },
	}

	action3 := &Action{
		ID:      "test.all3",
		Name:    "All Test 3",
		Execute: func(ctx *ActionContext) tea.Cmd { return nil },
	}

	reg.Register(action1)
	reg.Register(action2)
	reg.Register(action3)

	all := reg.All()
	if len(all) != 3 {
		t.Errorf("Expected 3 actions, got %d", len(all))
	}
}

func TestKeyPatternString(t *testing.T) {
	tests := []struct {
		name     string
		pattern  KeyPattern
		expected string
	}{
		{
			name:     "Simple key",
			pattern:  KeyPattern{Key: "a"},
			expected: "a",
		},
		{
			name:     "Ctrl+key",
			pattern:  KeyPattern{Key: "c", Ctrl: true},
			expected: "ctrl+c",
		},
		{
			name:     "Alt+key",
			pattern:  KeyPattern{Key: "x", Alt: true},
			expected: "alt+x",
		},
		{
			name:     "Shift+key",
			pattern:  KeyPattern{Key: "a", Shift: true},
			expected: "shift+a",
		},
		{
			name:     "Ctrl+Alt+key",
			pattern:  KeyPattern{Key: "delete", Ctrl: true, Alt: true},
			expected: "ctrl+alt+delete",
		},
		{
			name:     "All modifiers",
			pattern:  KeyPattern{Key: "f", Ctrl: true, Alt: true, Shift: true},
			expected: "ctrl+alt+shift+f",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pattern.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
