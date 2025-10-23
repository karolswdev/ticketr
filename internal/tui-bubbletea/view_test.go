package tuibubbletea

import (
	"errors"
	"strings"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// TestViewNotReady tests view rendering before ready
func TestViewNotReady(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = false

	view := m.View()

	if view == "" {
		t.Error("Expected non-empty view even when not ready")
	}

	// Should show loading state
	// The loadingSpinner.View() should be present
	if len(view) < 10 {
		t.Error("Expected loading view to have substantial content")
	}
}

// TestViewTerminalTooSmall tests terminal too small error message
func TestViewTerminalTooSmall(t *testing.T) {
	m := initialModel(nil, nil)
	m.terminalTooSmall = true
	m.width = 60
	m.height = 20

	view := m.View()

	if !strings.Contains(view, "Terminal Too Small") &&
	   !strings.Contains(view, "Too Small") {
		t.Error("Expected 'Terminal Too Small' message in view")
	}

	// Should show current and required dimensions
	if !strings.Contains(view, "60") || !strings.Contains(view, "20") {
		t.Error("Expected current terminal dimensions in error message")
	}

	if !strings.Contains(view, "80") {
		t.Error("Expected minimum width requirement in error message")
	}
}

// TestViewErrorState tests error rendering
func TestViewErrorState(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.loadError = errors.New("test error message")
	m.width = 100
	m.height = 30

	view := m.View()

	// Should contain error information
	if !strings.Contains(view, "Error") && !strings.Contains(view, "error") {
		t.Error("Expected error indicator in view")
	}

	// Note: The exact error message might be rendered by components.RenderCenteredError
	// which may format it differently
	if len(view) < 10 {
		t.Error("Expected error view to have content")
	}
}

// TestViewLoadingState tests loading state rendering
func TestViewLoadingState(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.dataLoaded = false
	m.width = 100
	m.height = 30

	view := m.View()

	// Should show loading spinner
	if len(view) < 10 {
		t.Error("Expected loading view to have content")
	}

	// The view should be centered using lipgloss.Place
	// We can't check exact content without running the full rendering,
	// but we can verify it's not empty
}

// TestViewHelpScreen tests help screen rendering
func TestViewHelpScreen(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.dataLoaded = true
	m.width = 100
	m.height = 30

	// Show help
	m.helpScreen.Toggle()

	view := m.View()

	// Help screen should be visible and contain keyboard shortcuts
	// The exact content depends on the help component, but it should be substantial
	if len(view) < 50 {
		t.Error("Expected help screen to have substantial content")
	}

	// Help screen is rendered as a modal, so it should overlay everything
}

// TestViewWorkspaceModal tests workspace modal rendering
func TestViewWorkspaceModal(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.dataLoaded = true
	m.showWorkspaceModal = true
	m.width = 100
	m.height = 30

	view := m.View()

	// Workspace modal should be visible
	if len(view) < 50 {
		t.Error("Expected workspace modal to have substantial content")
	}
}

// TestViewNormalState tests normal state rendering
func TestViewNormalState(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.dataLoaded = true
	m.width = 120
	m.height = 40
	m.currentWorkspace = &domain.Workspace{
		ID:   "ws-test-1",
		Name: "Test Workspace",
	}
	m.tickets = []domain.Ticket{
		{JiraID: "TEST-1", Title: "Test Ticket"},
	}

	// Update tree with tickets
	m.ticketTree.SetTickets(m.tickets)

	view := m.View()

	// Should contain workspace name
	if !strings.Contains(view, "Test Workspace") {
		t.Error("Expected workspace name in view")
	}

	// Should render substantial content (header + panels + footer)
	if len(view) < 100 {
		t.Errorf("Expected substantial view content, got %d chars", len(view))
	}
}

// TestViewRenderLoading tests the renderLoading method
func TestViewRenderLoading(t *testing.T) {
	m := initialModel(nil, nil)
	m.width = 100
	m.height = 30

	view := m.renderLoading()

	if view == "" {
		t.Error("Expected non-empty loading view")
	}

	// Should contain spinner content
	if len(view) < 10 {
		t.Error("Expected loading view to have content")
	}
}

// TestViewRenderHeader tests header rendering
func TestViewRenderHeader(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.width = 120
	m.currentWorkspace = &domain.Workspace{Name: "Test WS"}
	m.tickets = []domain.Ticket{{JiraID: "TEST-1"}}

	header := m.renderHeader()

	// Should contain workspace name
	if !strings.Contains(header, "Test WS") {
		t.Error("Expected workspace name in header")
	}

	// Should contain ticket count
	if !strings.Contains(header, "1") {
		t.Error("Expected ticket count in header")
	}

	// Should contain theme name
	if !strings.Contains(header, "Default") {
		t.Error("Expected theme name in header")
	}

	// Should contain focus indicator
	if !strings.Contains(header, "Tree") {
		t.Error("Expected focus indicator in header")
	}
}

// TestViewRenderHeaderNoWorkspace tests header with no workspace
func TestViewRenderHeaderNoWorkspace(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.width = 120
	m.currentWorkspace = nil

	header := m.renderHeader()

	if !strings.Contains(header, "No workspace") {
		t.Error("Expected 'No workspace' in header when workspace is nil")
	}
}

// TestViewRenderLeftPanel tests left panel rendering
func TestViewRenderLeftPanel(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.dataLoaded = true
	m.width = 120
	m.height = 40
	m.focused = FocusLeft
	m.currentWorkspace = &domain.Workspace{Name: "Test WS"}
	m.tickets = []domain.Ticket{{JiraID: "TEST-1", Title: "Test"}}

	leftPanel := m.renderLeftPanel()

	// Should have content
	if leftPanel == "" {
		t.Error("Expected non-empty left panel")
	}

	// Should contain workspace name or ticket info
	if !strings.Contains(leftPanel, "Test WS") {
		t.Error("Expected workspace name in left panel")
	}
}

// TestViewRenderRightPanel tests right panel rendering
func TestViewRenderRightPanel(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.dataLoaded = true
	m.width = 120
	m.height = 40
	m.focused = FocusRight

	rightPanel := m.renderRightPanel()

	// Should have content
	if rightPanel == "" {
		t.Error("Expected non-empty right panel")
	}

	// Should contain "Ticket Detail" title
	if !strings.Contains(rightPanel, "Ticket Detail") {
		t.Error("Expected 'Ticket Detail' title in right panel")
	}
}

// TestViewRenderActionBar tests action bar rendering
func TestViewRenderActionBar(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.width = 120

	actionBar := m.renderActionBar()

	// Should contain key bindings
	expectedKeys := []string{"Navigate", "Expand", "Collapse", "Select", "Switch", "Workspace", "Help", "Quit"}

	for _, key := range expectedKeys {
		if !strings.Contains(actionBar, key) {
			t.Errorf("Expected action bar to contain '%s'", key)
		}
	}
}

// TestViewGetFocusName tests focus name retrieval
func TestViewGetFocusName(t *testing.T) {
	m := initialModel(nil, nil)

	tests := []struct {
		focus    Focus
		expected string
	}{
		{FocusLeft, "Tree"},
		{FocusRight, "Detail"},
		{FocusWorkspace, "Workspace Selector"},
	}

	for _, tt := range tests {
		m.focused = tt.focus
		name := m.getFocusName()
		if name != tt.expected {
			t.Errorf("For focus %v, expected '%s', got '%s'", tt.focus, tt.expected, name)
		}
	}
}

// TestViewRenderTerminalTooSmallError tests terminal too small rendering
func TestViewRenderTerminalTooSmallError(t *testing.T) {
	m := initialModel(nil, nil)
	m.width = 70
	m.height = 20

	view := m.renderTerminalTooSmallError()

	// Should contain error message
	if !strings.Contains(view, "Terminal Too Small") &&
	   !strings.Contains(view, "Too Small") {
		t.Error("Expected 'Terminal Too Small' in error view")
	}

	// Should contain current dimensions
	if !strings.Contains(view, "70") || !strings.Contains(view, "20") {
		t.Error("Expected current dimensions in error message")
	}

	// Should contain required dimensions
	if !strings.Contains(view, "80") && !strings.Contains(view, "24") {
		t.Error("Expected required dimensions (80x24) in error message")
	}

	// Should mention resize
	if !strings.Contains(view, "resize") {
		t.Error("Expected 'resize' instruction in error message")
	}
}

// TestViewPriorityOrderTerminalTooSmall tests that terminal too small has highest priority
func TestViewPriorityOrderTerminalTooSmall(t *testing.T) {
	m := initialModel(nil, nil)
	m.terminalTooSmall = true
	m.ready = true
	m.dataLoaded = true
	m.loadError = errors.New("some error")
	m.width = 60
	m.height = 20

	view := m.View()

	// Terminal too small should take precedence over everything
	if !strings.Contains(view, "Terminal Too Small") &&
	   !strings.Contains(view, "Too Small") {
		t.Error("Expected terminal too small error to have highest priority")
	}
}

// TestViewPriorityOrderNotReady tests that not ready takes precedence over errors
func TestViewPriorityOrderNotReady(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = false
	m.loadError = errors.New("some error")
	m.width = 100
	m.height = 30

	view := m.View()

	// Should show loading, not error (not ready takes precedence)
	// Loading spinner should be rendered
	if len(view) == 0 {
		t.Error("Expected view to render loading state")
	}
}

// TestViewPriorityOrderHelpScreen tests that help screen overlays everything
func TestViewPriorityOrderHelpScreen(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.dataLoaded = true
	m.showWorkspaceModal = true
	m.helpScreen.Toggle()
	m.width = 100
	m.height = 30

	view := m.View()

	// Help screen should overlay even the workspace modal
	// We can't easily test the exact content, but we verify it renders
	if len(view) < 50 {
		t.Error("Expected help screen to render with substantial content")
	}
}

// TestViewDifferentFocusStates tests view with different focus states
func TestViewDifferentFocusStates(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.dataLoaded = true
	m.width = 120
	m.height = 40

	focuses := []Focus{FocusLeft, FocusRight, FocusWorkspace}

	for _, focus := range focuses {
		m.focused = focus
		view := m.View()

		// Should render successfully regardless of focus
		if len(view) < 50 {
			t.Errorf("Expected view to render for focus %v, got %d chars", focus, len(view))
		}

		// Should contain focus name
		focusName := m.getFocusName()
		if !strings.Contains(view, focusName) {
			t.Errorf("Expected view to contain focus name '%s' for focus %v", focusName, focus)
		}
	}
}

// TestViewWithSelectedTicket tests view with a selected ticket
func TestViewWithSelectedTicket(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.dataLoaded = true
	m.width = 120
	m.height = 40
	m.selectedTicket = &domain.Ticket{
		JiraID: "TEST-1",
		Title:  "Selected Ticket",
	}
	m.detailView.SetTicket(m.selectedTicket)

	view := m.View()

	// Should render with selected ticket
	if len(view) < 100 {
		t.Error("Expected view to render with selected ticket")
	}
}

// TestViewEmptyTickets tests view with no tickets
func TestViewEmptyTickets(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.dataLoaded = true
	m.width = 120
	m.height = 40
	m.tickets = []domain.Ticket{}
	m.currentWorkspace = &domain.Workspace{Name: "Empty WS"}

	view := m.View()

	// Should still render successfully
	if len(view) < 50 {
		t.Error("Expected view to render even with no tickets")
	}

	// Should show ticket count as 0
	if !strings.Contains(view, "0") {
		t.Error("Expected ticket count 0 in view")
	}
}

// TestViewManyTickets tests view with many tickets
func TestViewManyTickets(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.dataLoaded = true
	m.width = 120
	m.height = 40

	// Create 100 tickets
	tickets := make([]domain.Ticket, 100)
	for i := 0; i < 100; i++ {
		tickets[i] = domain.Ticket{
			JiraID: "TEST-" + string(rune(i)),
			Title:  "Ticket " + string(rune(i)),
		}
	}
	m.tickets = tickets
	m.ticketTree.SetTickets(tickets)

	view := m.View()

	// Should render successfully even with many tickets
	if len(view) < 100 {
		t.Error("Expected view to render with many tickets")
	}

	// Should show ticket count
	if !strings.Contains(view, "100") {
		t.Error("Expected ticket count 100 in view")
	}
}
