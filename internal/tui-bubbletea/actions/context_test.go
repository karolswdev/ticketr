package actions

import (
	"testing"
)

func TestContextManagerNew(t *testing.T) {
	cm := NewContextManager(ContextWorkspaceList)

	if cm.Current() != ContextWorkspaceList {
		t.Errorf("Expected current context to be ContextWorkspaceList, got %s", cm.Current())
	}

	if cm.state.Previous != "" {
		t.Errorf("Expected previous context to be empty, got %s", cm.state.Previous)
	}

	if len(cm.state.Stack) != 1 {
		t.Errorf("Expected stack length 1, got %d", len(cm.state.Stack))
	}
}

func TestContextManagerSwitch(t *testing.T) {
	cm := NewContextManager(ContextWorkspaceList)

	// Switch to ticket tree
	cm.Switch(ContextTicketTree)

	if cm.Current() != ContextTicketTree {
		t.Errorf("Expected current context to be ContextTicketTree, got %s", cm.Current())
	}

	if cm.state.Previous != ContextWorkspaceList {
		t.Errorf("Expected previous context to be ContextWorkspaceList, got %s", cm.state.Previous)
	}

	// Switch to ticket detail
	cm.Switch(ContextTicketDetail)

	if cm.Current() != ContextTicketDetail {
		t.Errorf("Expected current context to be ContextTicketDetail, got %s", cm.Current())
	}

	if cm.state.Previous != ContextTicketTree {
		t.Errorf("Expected previous context to be ContextTicketTree, got %s", cm.state.Previous)
	}
}

func TestContextManagerSwitchSameContext(t *testing.T) {
	cm := NewContextManager(ContextWorkspaceList)

	initialStackLen := len(cm.state.Stack)

	// Switch to same context (should be no-op)
	cm.Switch(ContextWorkspaceList)

	if len(cm.state.Stack) != initialStackLen {
		t.Errorf("Expected stack length to remain %d, got %d", initialStackLen, len(cm.state.Stack))
	}
}

func TestContextManagerPushPop(t *testing.T) {
	cm := NewContextManager(ContextTicketTree)

	// Push modal context
	cm.Push(ContextModal)

	if cm.Current() != ContextModal {
		t.Errorf("Expected current context to be ContextModal, got %s", cm.Current())
	}

	if len(cm.state.Stack) != 2 {
		t.Errorf("Expected stack length 2, got %d", len(cm.state.Stack))
	}

	// Pop modal context
	prev := cm.Pop()

	if prev != ContextTicketTree {
		t.Errorf("Expected popped context to be ContextTicketTree, got %s", prev)
	}

	if cm.Current() != ContextTicketTree {
		t.Errorf("Expected current context to be ContextTicketTree after pop, got %s", cm.Current())
	}

	if len(cm.state.Stack) != 1 {
		t.Errorf("Expected stack length 1 after pop, got %d", len(cm.state.Stack))
	}
}

func TestContextManagerPopEmpty(t *testing.T) {
	cm := NewContextManager(ContextWorkspaceList)

	// Try to pop when stack has only one element (should return current)
	result := cm.Pop()

	if result != ContextWorkspaceList {
		t.Errorf("Expected pop to return ContextWorkspaceList, got %s", result)
	}

	if cm.Current() != ContextWorkspaceList {
		t.Errorf("Expected current context to remain ContextWorkspaceList, got %s", cm.Current())
	}

	if len(cm.state.Stack) != 1 {
		t.Errorf("Expected stack length to remain 1, got %d", len(cm.state.Stack))
	}
}

func TestContextManagerOnChange(t *testing.T) {
	cm := NewContextManager(ContextWorkspaceList)

	var oldContext, newContext Context
	callCount := 0

	// Register observer
	cm.OnChange(func(old, new Context) {
		oldContext = old
		newContext = new
		callCount++
	})

	// Switch context
	cm.Switch(ContextTicketTree)

	if callCount != 1 {
		t.Errorf("Expected observer to be called once, got %d times", callCount)
	}

	if oldContext != ContextWorkspaceList {
		t.Errorf("Expected old context to be ContextWorkspaceList, got %s", oldContext)
	}

	if newContext != ContextTicketTree {
		t.Errorf("Expected new context to be ContextTicketTree, got %s", newContext)
	}
}

func TestContextManagerMultipleObservers(t *testing.T) {
	cm := NewContextManager(ContextWorkspaceList)

	call1 := 0
	call2 := 0

	cm.OnChange(func(old, new Context) {
		call1++
	})

	cm.OnChange(func(old, new Context) {
		call2++
	})

	// Switch context
	cm.Switch(ContextTicketTree)

	if call1 != 1 {
		t.Errorf("Expected first observer to be called once, got %d times", call1)
	}

	if call2 != 1 {
		t.Errorf("Expected second observer to be called once, got %d times", call2)
	}
}

func TestContextManagerMetadata(t *testing.T) {
	cm := NewContextManager(ContextWorkspaceList)

	// Set metadata
	cm.SetMetadata("test_key", "test_value")

	// Get metadata
	value, ok := cm.GetMetadata("test_key")
	if !ok {
		t.Error("Expected metadata to exist")
	}

	if value != "test_value" {
		t.Errorf("Expected metadata value 'test_value', got %v", value)
	}

	// Get non-existent metadata
	_, ok = cm.GetMetadata("non_existent")
	if ok {
		t.Error("Expected non-existent metadata to return false")
	}
}

func TestContextStateIsIn(t *testing.T) {
	state := ContextState{
		Current: ContextTicketTree,
	}

	// Test exact match
	if !state.IsIn(ContextTicketTree) {
		t.Error("Expected IsIn to return true for exact match")
	}

	// Test no match
	if state.IsIn(ContextWorkspaceList) {
		t.Error("Expected IsIn to return false for non-matching context")
	}

	// Test multiple contexts with match
	if !state.IsIn(ContextWorkspaceList, ContextTicketTree, ContextTicketDetail) {
		t.Error("Expected IsIn to return true when one of multiple contexts matches")
	}

	// Test global context (always matches)
	if !state.IsIn(ContextGlobal) {
		t.Error("Expected IsIn to return true for ContextGlobal")
	}

	// Test combination with global
	if !state.IsIn(ContextWorkspaceList, ContextGlobal) {
		t.Error("Expected IsIn to return true when ContextGlobal is in the list")
	}
}

func TestContextManagerPushPopSequence(t *testing.T) {
	cm := NewContextManager(ContextWorkspaceList)

	// Push multiple contexts
	cm.Push(ContextTicketTree)
	cm.Push(ContextModal)
	cm.Push(ContextSearch)

	if cm.Current() != ContextSearch {
		t.Errorf("Expected current context to be ContextSearch, got %s", cm.Current())
	}

	if len(cm.state.Stack) != 4 {
		t.Errorf("Expected stack length 4, got %d", len(cm.state.Stack))
	}

	// Pop back through the stack
	cm.Pop()
	if cm.Current() != ContextModal {
		t.Errorf("Expected current context to be ContextModal after first pop, got %s", cm.Current())
	}

	cm.Pop()
	if cm.Current() != ContextTicketTree {
		t.Errorf("Expected current context to be ContextTicketTree after second pop, got %s", cm.Current())
	}

	cm.Pop()
	if cm.Current() != ContextWorkspaceList {
		t.Errorf("Expected current context to be ContextWorkspaceList after third pop, got %s", cm.Current())
	}
}
