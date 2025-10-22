package widgets

import (
	"strings"
	"testing"
)

func TestNewActionBar(t *testing.T) {
	ab := NewActionBar()
	if ab == nil {
		t.Fatal("NewActionBar() returned nil")
	}

	if ab.context != ContextWorkspaceList {
		t.Errorf("expected default context to be ContextWorkspaceList, got %v", ab.context)
	}

	if ab.bindings == nil {
		t.Fatal("bindings map is nil")
	}
}

func TestActionBar_SetContext(t *testing.T) {
	ab := NewActionBar()

	tests := []struct {
		name    string
		context ActionBarContext
	}{
		{"workspace list", ContextWorkspaceList},
		{"ticket tree", ContextTicketTree},
		{"ticket detail", ContextTicketDetail},
		{"modal", ContextModal},
		{"syncing", ContextSyncing},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab.SetContext(tt.context)
			if ab.GetContext() != tt.context {
				t.Errorf("expected context %v, got %v", tt.context, ab.GetContext())
			}
		})
	}
}

func TestActionBar_GetContext(t *testing.T) {
	ab := NewActionBar()
	ab.SetContext(ContextTicketTree)

	if got := ab.GetContext(); got != ContextTicketTree {
		t.Errorf("GetContext() = %v, want %v", got, ContextTicketTree)
	}
}

func TestActionBar_AddBinding(t *testing.T) {
	ab := NewActionBar()
	testContext := ContextWorkspaceList

	// Get initial count
	initialCount := len(ab.bindings[testContext])

	// Add a custom binding
	ab.AddBinding(testContext, "x", "Test Action")

	// Verify binding was added
	if len(ab.bindings[testContext]) != initialCount+1 {
		t.Errorf("expected %d bindings, got %d", initialCount+1, len(ab.bindings[testContext]))
	}

	// Find the added binding
	found := false
	for _, binding := range ab.bindings[testContext] {
		if binding.Key == "x" && binding.Description == "Test Action" {
			found = true
			break
		}
	}

	if !found {
		t.Error("added binding not found in bindings list")
	}
}

func TestActionBar_SetBindings(t *testing.T) {
	ab := NewActionBar()
	testContext := ContextTicketDetail

	customBindings := []KeyBinding{
		{Key: "a", Description: "Action A"},
		{Key: "b", Description: "Action B"},
	}

	ab.SetBindings(testContext, customBindings)

	if len(ab.bindings[testContext]) != 2 {
		t.Errorf("expected 2 bindings, got %d", len(ab.bindings[testContext]))
	}

	if ab.bindings[testContext][0].Key != "a" {
		t.Errorf("expected first binding key 'a', got '%s'", ab.bindings[testContext][0].Key)
	}
}

func TestActionBar_Update(t *testing.T) {
	ab := NewActionBar()

	// Test that update generates text
	ab.SetContext(ContextWorkspaceList)
	text := ab.GetText(false)

	if text == "" {
		t.Error("expected non-empty text after update")
	}

	// Verify text contains keybindings
	if !strings.Contains(text, "Enter") && !strings.Contains(text, "Tab") {
		t.Errorf("expected text to contain keybindings, got: %s", text)
	}
}

func TestActionBar_ContextSwitch(t *testing.T) {
	ab := NewActionBar()
	defer ab.Shutdown() // Ensure cleanup

	// Switch contexts and verify bindings are different
	ab.SetContext(ContextWorkspaceList)
	bindings1 := ab.bindings[ContextWorkspaceList]

	ab.SetContext(ContextTicketTree)
	bindings2 := ab.bindings[ContextTicketTree]

	// Verify different number of bindings (or different keys)
	if len(bindings1) == 0 || len(bindings2) == 0 {
		t.Error("expected both contexts to have bindings")
	}

	// Check that at least one key is different
	keys1 := make(map[string]bool)
	for _, b := range bindings1 {
		keys1[b.Key] = true
	}

	allMatch := true
	for _, b := range bindings2 {
		if !keys1[b.Key] {
			allMatch = false
			break
		}
	}

	// If all keys match AND same length, contexts are the same (unexpected)
	if allMatch && len(bindings1) == len(bindings2) {
		t.Error("expected different keybindings for different contexts")
	}
}

func TestActionBar_EmptyBindings(t *testing.T) {
	ab := NewActionBar()

	// Create a context with no bindings
	customContext := ActionBarContext("custom")
	ab.SetContext(customContext)

	text := ab.GetText(false)
	if text != "" {
		t.Errorf("expected empty text for context with no bindings, got: %s", text)
	}
}
