package commands

import (
	"fmt"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry() returned nil")
	}

	if r.commands == nil {
		t.Fatal("commands map is nil")
	}

	if r.Count() != 0 {
		t.Errorf("expected empty registry, got %d commands", r.Count())
	}
}

func TestRegistry_Register(t *testing.T) {
	r := NewRegistry()

	cmd := &Command{
		Name:        "test",
		Description: "Test command",
		Keybinding:  "t",
		Category:    CategorySystem,
		Handler:     func() error { return nil },
	}

	err := r.Register(cmd)
	if err != nil {
		t.Fatalf("Register() failed: %v", err)
	}

	if r.Count() != 1 {
		t.Errorf("expected 1 command, got %d", r.Count())
	}
}

func TestRegistry_Register_NoName(t *testing.T) {
	r := NewRegistry()

	cmd := &Command{
		Description: "Test command",
		Handler:     func() error { return nil },
	}

	err := r.Register(cmd)
	if err == nil {
		t.Error("expected error when registering command without name")
	}
}

func TestRegistry_Register_NoHandler(t *testing.T) {
	r := NewRegistry()

	cmd := &Command{
		Name:        "test",
		Description: "Test command",
	}

	err := r.Register(cmd)
	if err == nil {
		t.Error("expected error when registering command without handler")
	}
}

func TestRegistry_Get(t *testing.T) {
	r := NewRegistry()

	cmd := &Command{
		Name:        "test",
		Description: "Test command",
		Handler:     func() error { return nil },
	}

	r.Register(cmd)

	retrieved, exists := r.Get("test")
	if !exists {
		t.Error("expected command to exist")
	}

	if retrieved.Name != "test" {
		t.Errorf("expected name 'test', got '%s'", retrieved.Name)
	}
}

func TestRegistry_Get_NotFound(t *testing.T) {
	r := NewRegistry()

	_, exists := r.Get("nonexistent")
	if exists {
		t.Error("expected command to not exist")
	}
}

func TestRegistry_All(t *testing.T) {
	r := NewRegistry()

	commands := []*Command{
		{Name: "cmd1", Description: "Command 1", Category: CategorySync, Handler: func() error { return nil }},
		{Name: "cmd2", Description: "Command 2", Category: CategoryView, Handler: func() error { return nil }},
		{Name: "cmd3", Description: "Command 3", Category: CategorySync, Handler: func() error { return nil }},
	}

	for _, cmd := range commands {
		r.Register(cmd)
	}

	all := r.All()
	if len(all) != 3 {
		t.Errorf("expected 3 commands, got %d", len(all))
	}

	// Verify sorted by category then name
	if all[0].Category > all[1].Category {
		t.Error("commands not sorted by category")
	}
}

func TestRegistry_ByCategory(t *testing.T) {
	r := NewRegistry()

	r.Register(&Command{Name: "sync1", Category: CategorySync, Handler: func() error { return nil }})
	r.Register(&Command{Name: "view1", Category: CategoryView, Handler: func() error { return nil }})
	r.Register(&Command{Name: "sync2", Category: CategorySync, Handler: func() error { return nil }})

	byCategory := r.ByCategory()

	if len(byCategory[CategorySync]) != 2 {
		t.Errorf("expected 2 sync commands, got %d", len(byCategory[CategorySync]))
	}

	if len(byCategory[CategoryView]) != 1 {
		t.Errorf("expected 1 view command, got %d", len(byCategory[CategoryView]))
	}

	// Verify sorted within category
	if byCategory[CategorySync][0].Name > byCategory[CategorySync][1].Name {
		t.Error("commands not sorted within category")
	}
}

func TestRegistry_Search(t *testing.T) {
	r := NewRegistry()

	r.Register(&Command{Name: "pull", Description: "Pull tickets from Jira", Handler: func() error { return nil }})
	r.Register(&Command{Name: "push", Description: "Push tickets to Jira", Handler: func() error { return nil }})
	r.Register(&Command{Name: "sync", Description: "Full sync", Handler: func() error { return nil }})

	tests := []struct {
		name          string
		query         string
		expectedCount int
	}{
		{"empty query", "", 3},
		{"match name", "pull", 1},
		{"match description", "Jira", 2},
		{"case insensitive", "PUSH", 1},
		{"no match", "xyz", 0},
		{"partial match", "pu", 2}, // matches "pull" and "push"
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := r.Search(tt.query)
			if len(results) != tt.expectedCount {
				t.Errorf("expected %d results, got %d", tt.expectedCount, len(results))
			}
		})
	}
}

func TestRegistry_Search_PrioritizesNameMatches(t *testing.T) {
	r := NewRegistry()

	r.Register(&Command{Name: "test", Description: "Other description", Handler: func() error { return nil }})
	r.Register(&Command{Name: "other", Description: "Contains test in description", Handler: func() error { return nil }})

	results := r.Search("test")

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	// First result should be the name match
	if results[0].Name != "test" {
		t.Errorf("expected first result to be 'test', got '%s'", results[0].Name)
	}
}

func TestRegistry_Execute(t *testing.T) {
	r := NewRegistry()

	executed := false
	r.Register(&Command{
		Name:    "test",
		Handler: func() error { executed = true; return nil },
	})

	err := r.Execute("test")
	if err != nil {
		t.Errorf("Execute() failed: %v", err)
	}

	if !executed {
		t.Error("command handler was not executed")
	}
}

func TestRegistry_Execute_NotFound(t *testing.T) {
	r := NewRegistry()

	err := r.Execute("nonexistent")
	if err == nil {
		t.Error("expected error when executing non-existent command")
	}
}

func TestRegistry_Execute_HandlerError(t *testing.T) {
	r := NewRegistry()

	expectedError := fmt.Errorf("handler error")
	r.Register(&Command{
		Name:    "test",
		Handler: func() error { return expectedError },
	})

	err := r.Execute("test")
	if err != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}

func TestRegistry_Count(t *testing.T) {
	r := NewRegistry()

	if r.Count() != 0 {
		t.Errorf("expected count 0, got %d", r.Count())
	}

	r.Register(&Command{Name: "cmd1", Handler: func() error { return nil }})
	r.Register(&Command{Name: "cmd2", Handler: func() error { return nil }})

	if r.Count() != 2 {
		t.Errorf("expected count 2, got %d", r.Count())
	}
}
