package tuibubbletea

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// TestInitialModel tests the initialModel function with nil services
func TestInitialModel(t *testing.T) {
	model := initialModel(nil, nil)

	// Verify default theme is set
	if model.theme == nil {
		t.Error("Expected default theme to be set")
	}

	if model.theme.Name != "Default" {
		t.Errorf("Expected theme name 'Default', got '%s'", model.theme.Name)
	}

	// Verify initial focus state
	if model.focused != FocusLeft {
		t.Errorf("Expected focus to be FocusLeft, got %v", model.focused)
	}

	// Verify not ready before WindowSizeMsg
	if model.ready {
		t.Error("Expected model to not be ready before WindowSizeMsg")
	}

	// Verify loading states
	if !model.loadingWorkspaces {
		t.Error("Expected loadingWorkspaces to be true initially")
	}

	if model.loadingTickets {
		t.Error("Expected loadingTickets to be false initially")
	}

	if model.dataLoaded {
		t.Error("Expected dataLoaded to be false initially")
	}

	// Verify components initialized
	if model.ticketTree.View() == "" {
		t.Error("Expected tree component to be initialized")
	}

	// Verify default dimensions
	if model.width == 0 || model.height == 0 {
		t.Error("Expected default dimensions to be set")
	}

	// Verify terminal size validation is off
	if model.terminalTooSmall {
		t.Error("Expected terminalTooSmall to be false initially")
	}

	// Verify modal is closed
	if model.showWorkspaceModal {
		t.Error("Expected showWorkspaceModal to be false initially")
	}

	// Verify empty data slices
	if model.tickets == nil {
		t.Error("Expected tickets slice to be initialized")
	}

	if model.workspaces == nil {
		t.Error("Expected workspaces slice to be initialized")
	}

	// Verify no error state
	if model.loadError != nil {
		t.Error("Expected loadError to be nil initially")
	}
}

// TestInitialModelWithServices tests initialModel accepts services
func TestInitialModelWithServices(t *testing.T) {
	// Test that initialModel accepts non-nil services
	// Note: We can't easily create real service instances in unit tests
	// so we just verify nil handling works correctly
	model := initialModel(nil, nil)

	// Services can be nil for basic model creation
	// The model should still initialize correctly
	if model.theme == nil {
		t.Error("Expected theme to be initialized even with nil services")
	}
}

// TestInitCommand tests the Init command
func TestInitCommand(t *testing.T) {
	model := initialModel(nil, nil)
	cmd := model.Init()

	// Init should return a batch command that initializes all components
	if cmd == nil {
		t.Error("Expected Init to return a command")
	}
}

// TestModelDefaults tests default values
func TestModelDefaults(t *testing.T) {
	m := initialModel(nil, nil)

	tests := []struct {
		name     string
		actual   interface{}
		expected interface{}
		compare  func(a, b interface{}) bool
	}{
		{
			name:     "default theme name",
			actual:   m.theme.Name,
			expected: "Default",
			compare:  func(a, b interface{}) bool { return a.(string) == b.(string) },
		},
		{
			name:     "not ready",
			actual:   m.ready,
			expected: false,
			compare:  func(a, b interface{}) bool { return a.(bool) == b.(bool) },
		},
		{
			name:     "loading workspaces",
			actual:   m.loadingWorkspaces,
			expected: true,
			compare:  func(a, b interface{}) bool { return a.(bool) == b.(bool) },
		},
		{
			name:     "not loading tickets",
			actual:   m.loadingTickets,
			expected: false,
			compare:  func(a, b interface{}) bool { return a.(bool) == b.(bool) },
		},
		{
			name:     "data not loaded",
			actual:   m.dataLoaded,
			expected: false,
			compare:  func(a, b interface{}) bool { return a.(bool) == b.(bool) },
		},
		{
			name:     "no error",
			actual:   m.loadError,
			expected: nil,
			compare:  func(a, b interface{}) bool { return a == nil && b == nil },
		},
		{
			name:     "terminal not too small",
			actual:   m.terminalTooSmall,
			expected: false,
			compare:  func(a, b interface{}) bool { return a.(bool) == b.(bool) },
		},
		{
			name:     "focus on left",
			actual:   m.focused,
			expected: FocusLeft,
			compare:  func(a, b interface{}) bool { return a.(Focus) == b.(Focus) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.compare(tt.actual, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, tt.actual)
			}
		})
	}
}

// TestGetCurrentTheme tests theme retrieval
func TestGetCurrentTheme(t *testing.T) {
	m := initialModel(nil, nil)

	themeName := m.GetCurrentTheme()
	if themeName != "Default" {
		t.Errorf("Expected theme 'Default', got '%s'", themeName)
	}

	// Change theme and verify
	m.theme = theme.GetByName("Dark")
	themeName = m.GetCurrentTheme()
	if themeName != "Dark" {
		t.Errorf("Expected theme 'Dark', got '%s'", themeName)
	}
}

// TestSetFocus tests focus setting
func TestSetFocus(t *testing.T) {
	m := initialModel(nil, nil)

	// Test setting focus to right
	m.SetFocus(FocusRight)
	if m.focused != FocusRight {
		t.Errorf("Expected focus to be FocusRight, got %v", m.focused)
	}

	// Test setting focus to workspace
	m.SetFocus(FocusWorkspace)
	if m.focused != FocusWorkspace {
		t.Errorf("Expected focus to be FocusWorkspace, got %v", m.focused)
	}

	// Test setting focus back to left
	m.SetFocus(FocusLeft)
	if m.focused != FocusLeft {
		t.Errorf("Expected focus to be FocusLeft, got %v", m.focused)
	}
}

// TestToggleFocus tests focus toggling
func TestToggleFocus(t *testing.T) {
	m := initialModel(nil, nil)

	// Start with left focus
	if m.focused != FocusLeft {
		t.Errorf("Expected initial focus to be FocusLeft, got %v", m.focused)
	}

	// Toggle to right
	m.ToggleFocus()
	if m.focused != FocusRight {
		t.Errorf("Expected focus to toggle to FocusRight, got %v", m.focused)
	}

	// Toggle back to left
	m.ToggleFocus()
	if m.focused != FocusLeft {
		t.Errorf("Expected focus to toggle back to FocusLeft, got %v", m.focused)
	}

	// Verify workspace focus doesn't toggle
	m.SetFocus(FocusWorkspace)
	m.ToggleFocus()
	// ToggleFocus only works between Left and Right
	if m.focused != FocusLeft {
		t.Errorf("Expected focus to be FocusLeft after toggle from workspace, got %v", m.focused)
	}
}

// TestSetSelectedTicket tests ticket selection
func TestSetSelectedTicket(t *testing.T) {
	m := initialModel(nil, nil)

	ticket := &domain.Ticket{
		JiraID: "TEST-1",
		Title:  "Test Ticket",
	}

	m.SetSelectedTicket(ticket)

	// Verify ticket is set
	if m.selectedTicket == nil {
		t.Error("Expected selectedTicket to be set")
	}

	if m.selectedTicket.JiraID != "TEST-1" {
		t.Errorf("Expected ticket JiraID 'TEST-1', got '%s'", m.selectedTicket.JiraID)
	}

	// Verify focus switches to right
	if m.focused != FocusRight {
		t.Errorf("Expected focus to be FocusRight after selection, got %v", m.focused)
	}
}

// TestSetSelectedTicketNil tests setting nil ticket
func TestSetSelectedTicketNil(t *testing.T) {
	m := initialModel(nil, nil)

	// First set a ticket
	ticket := &domain.Ticket{JiraID: "TEST-1", Title: "Test"}
	m.SetSelectedTicket(ticket)

	// Now set to nil
	m.SetSelectedTicket(nil)

	if m.selectedTicket != nil {
		t.Error("Expected selectedTicket to be nil")
	}
}

// Note: Workspace switching is now tested in update_test.go
// via workspace.WorkspaceSelectedMsg message handling (Week 3 Day 3)

// TestLayoutInitialization tests layout component initialization
func TestLayoutInitialization(t *testing.T) {
	m := initialModel(nil, nil)

	if m.layout == nil {
		t.Error("Expected layout to be initialized")
	}

	// Verify layout dimensions
	leftWidth, rightWidth, contentHeight := m.layout.GetPanelDimensions()

	if leftWidth == 0 {
		t.Error("Expected leftWidth to be non-zero")
	}

	if rightWidth == 0 {
		t.Error("Expected rightWidth to be non-zero")
	}

	if contentHeight == 0 {
		t.Error("Expected contentHeight to be non-zero")
	}

	// Verify left panel is roughly 40% of width
	totalWidth := leftWidth + rightWidth
	leftRatio := float64(leftWidth) / float64(totalWidth)

	if leftRatio < 0.35 || leftRatio > 0.45 {
		t.Errorf("Expected left panel to be roughly 40%% of width, got %.1f%%", leftRatio*100)
	}
}

// TestComponentsInitialization tests all child components are initialized
func TestComponentsInitialization(t *testing.T) {
	m := initialModel(nil, nil)

	// Test tree component
	if m.ticketTree.View() == "" {
		t.Error("Expected tree component to be initialized")
	}

	// Test detail view
	detailView := m.detailView.View()
	if detailView == "" {
		t.Error("Expected detail view to be initialized")
	}

	// Test workspace selector
	if m.workspaceSelector.View() == "" {
		t.Error("Expected workspace selector to be initialized")
	}

	// Test loading spinner
	if m.loadingSpinner.View() == "" {
		t.Error("Expected loading spinner to be initialized")
	}

	// Test help screen - may return empty view when not visible
	helpView := m.helpScreen.View()
	_ = helpView // Help screen component exists
}

// TestMaxFunction tests the max helper function
func TestMaxFunction(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"a greater", 10, 5, 10},
		{"b greater", 5, 10, 10},
		{"equal", 7, 7, 7},
		{"negative a", -5, 3, 3},
		{"negative b", 3, -5, 3},
		{"both negative", -5, -3, -3},
		{"zero", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := max(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("max(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// TestFocusConstants tests Focus enum values
func TestFocusConstants(t *testing.T) {
	// Verify the constants are distinct
	focuses := []Focus{FocusLeft, FocusRight, FocusWorkspace}

	for i, f1 := range focuses {
		for j, f2 := range focuses {
			if i != j && f1 == f2 {
				t.Errorf("Focus constants at index %d and %d have the same value", i, j)
			}
		}
	}
}

// TestModelWithRealData tests model with actual domain data
func TestModelWithRealData(t *testing.T) {
	m := initialModel(nil, nil)

	// Simulate receiving window size
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
	updatedModel, _ := m.Update(sizeMsg)
	m = updatedModel.(Model)

	// Verify ready state
	if !m.ready {
		t.Error("Expected model to be ready after WindowSizeMsg")
	}

	// Verify dimensions updated
	if m.width != 120 || m.height != 40 {
		t.Errorf("Expected dimensions 120x40, got %dx%d", m.width, m.height)
	}
}
