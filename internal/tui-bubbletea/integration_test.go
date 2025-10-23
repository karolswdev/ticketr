package tuibubbletea

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/messages"
)

// TestFullWorkflowAppStartToTicketSelection tests the complete workflow
func TestFullWorkflowAppStartToTicketSelection(t *testing.T) {
	// Create model with nil services (will be updated when we have proper interfaces)
	// For now, test the basic structure
	model := initialModel(nil, nil)

	// Simulate window size message
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 30}
	updatedModel, cmd := model.Update(sizeMsg)
	m := updatedModel.(Model)

	if !m.ready {
		t.Error("Expected model to be ready after WindowSizeMsg")
	}

	// Verify command was returned (should load workspace)
	if cmd == nil {
		t.Error("Expected command to be returned for data loading")
	}
}

// TestWindowResize tests window resize propagation
func TestWindowResize(t *testing.T) {
	// Create model with nil services
	model := initialModel(nil, nil)
	model.ready = true
	model.dataLoaded = true

	// Initial size
	initialWidth := 80
	initialHeight := 24
	msg := tea.WindowSizeMsg{Width: initialWidth, Height: initialHeight}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.width != initialWidth {
		t.Errorf("Expected width %d, got %d", initialWidth, m.width)
	}
	if m.height != initialHeight {
		t.Errorf("Expected height %d, got %d", initialHeight, m.height)
	}

	// Resize
	newWidth := 120
	newHeight := 40
	resizeMsg := tea.WindowSizeMsg{Width: newWidth, Height: newHeight}
	updatedModel, _ = m.Update(resizeMsg)
	m = updatedModel.(Model)

	if m.width != newWidth {
		t.Errorf("Expected width %d, got %d", newWidth, m.width)
	}
	if m.height != newHeight {
		t.Errorf("Expected height %d, got %d", newHeight, m.height)
	}
}

// TestWorkspaceLoading tests workspace loading message handling
func TestWorkspaceLoading(t *testing.T) {
	model := Model{
		ready:   true,
		tickets: []domain.Ticket{},
	}

	// Test successful workspace loading
	workspace := &domain.Workspace{
		ID:         "ws1",
		Name:       "Test Workspace",
		ProjectKey: "TEST",
	}

	msg := messages.CurrentWorkspaceLoadedMsg{
		Workspace: workspace,
		Error:     nil,
	}

	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.currentWorkspace == nil {
		t.Error("Expected workspace to be set")
	}
	if m.currentWorkspace.ID != "ws1" {
		t.Errorf("Expected workspace ID 'ws1', got '%s'", m.currentWorkspace.ID)
	}
	if m.loadingWorkspaces {
		t.Error("Expected loadingWorkspaces to be false after load")
	}
}

// TestTicketsLoading tests tickets loading message handling
func TestTicketsLoading(t *testing.T) {
	// Need to initialize model properly with tree component
	model := initialModel(nil, nil)
	model.ready = true

	// Create test tickets
	tickets := []domain.Ticket{
		{
			JiraID: "TEST-1",
			Title:  "First Ticket",
			Tasks:  []domain.Task{},
		},
		{
			JiraID: "TEST-2",
			Title:  "Second Ticket",
			Tasks: []domain.Task{
				{
					JiraID: "TEST-3",
					Title:  "Subtask",
				},
			},
		},
	}

	msg := messages.TicketsLoadedMsg{
		Tickets: tickets,
		Error:   nil,
	}

	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if len(m.tickets) != 2 {
		t.Errorf("Expected 2 tickets, got %d", len(m.tickets))
	}
	if !m.dataLoaded {
		t.Error("Expected dataLoaded to be true")
	}
	if m.loadingTickets {
		t.Error("Expected loadingTickets to be false")
	}
}

// TestFocusSwitching tests focus management between panels
func TestFocusSwitching(t *testing.T) {
	model := Model{
		ready:   true,
		focused: FocusLeft,
	}

	// Test toggle focus
	model.ToggleFocus()
	if model.focused != FocusRight {
		t.Errorf("Expected focus to be Right, got %d", model.focused)
	}

	model.ToggleFocus()
	if model.focused != FocusLeft {
		t.Errorf("Expected focus to be Left, got %d", model.focused)
	}

	// Test set focus
	model.SetFocus(FocusWorkspace)
	if model.focused != FocusWorkspace {
		t.Errorf("Expected focus to be Workspace, got %d", model.focused)
	}
}

// TestKeyboardNavigation tests keyboard input handling
func TestKeyboardNavigation(t *testing.T) {
	model := Model{
		ready:      true,
		focused:    FocusLeft,
		dataLoaded: true,
	}

	// Test quit key
	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.Update(quitMsg)
	if cmd == nil {
		t.Error("Expected quit command")
	}

	// Test tab key (focus switch)
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := model.Update(tabMsg)
	m := updatedModel.(Model)
	if m.focused != FocusRight {
		t.Error("Expected focus to switch to right panel")
	}

	// Test workspace modal
	wsMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'W'}}
	updatedModel, _ = model.Update(wsMsg)
	m = updatedModel.(Model)
	if !m.showWorkspaceModal {
		t.Error("Expected workspace modal to be shown")
	}
	if m.focused != FocusWorkspace {
		t.Error("Expected focus to be on workspace selector")
	}
}

// TestErrorHandling tests error state handling
func TestErrorHandling(t *testing.T) {
	model := Model{
		ready:   true,
		tickets: []domain.Ticket{},
	}

	// Test workspace loading error
	msg := messages.CurrentWorkspaceLoadedMsg{
		Workspace: nil,
		Error:     errors.New("workspace not found"),
	}

	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.loadError == nil {
		t.Error("Expected loadError to be set")
	}
	if m.loadingWorkspaces {
		t.Error("Expected loadingWorkspaces to be false after error")
	}

	// Test tickets loading error
	ticketMsg := messages.TicketsLoadedMsg{
		Tickets: nil,
		Error:   errors.New("failed to load tickets"),
	}

	updatedModel, _ = model.Update(ticketMsg)
	m = updatedModel.(Model)

	if m.loadError == nil {
		t.Error("Expected loadError to be set")
	}
	if m.loadingTickets {
		t.Error("Expected loadingTickets to be false after error")
	}
}
