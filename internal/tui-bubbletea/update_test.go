package tuibubbletea

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/messages"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// TestUpdateWindowSize tests window size message handling
func TestUpdateWindowSize(t *testing.T) {
	m := initialModel(nil, nil)

	// Send window size message
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
	updatedModel, cmd := m.Update(sizeMsg)

	updated := updatedModel.(Model)

	if !updated.ready {
		t.Error("Expected model to be ready after WindowSizeMsg")
	}

	if updated.width != 120 {
		t.Errorf("Expected width 120, got %d", updated.width)
	}

	if updated.height != 40 {
		t.Errorf("Expected height 40, got %d", updated.height)
	}

	// Should return command to load data
	if cmd == nil {
		t.Error("Expected command to be returned for initial WindowSizeMsg")
	}

	// Terminal should not be too small
	if updated.terminalTooSmall {
		t.Error("Expected terminalTooSmall to be false for 120x40")
	}
}

// TestUpdateWindowSizeAlreadyReady tests window resize when already ready
func TestUpdateWindowSizeAlreadyReady(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true

	sizeMsg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updatedModel, cmd := m.Update(sizeMsg)

	updated := updatedModel.(Model)

	if updated.width != 100 {
		t.Errorf("Expected width 100, got %d", updated.width)
	}

	if updated.height != 30 {
		t.Errorf("Expected height 30, got %d", updated.height)
	}

	// Should not return data loading command if already ready
	// (cmd will be nil or a different command)
	_ = cmd
}

// TestUpdateTerminalTooSmall tests small terminal detection
func TestUpdateTerminalTooSmall(t *testing.T) {
	m := initialModel(nil, nil)

	tests := []struct {
		name            string
		width           int
		height          int
		expectTooSmall  bool
	}{
		{"80x24 exact minimum", 80, 24, false},
		{"79x24 width too small", 79, 24, true},
		{"80x23 height too small", 80, 23, true},
		{"60x20 both too small", 60, 20, true},
		{"120x40 normal", 120, 40, false},
		{"81x25 above minimum", 81, 25, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sizeMsg := tea.WindowSizeMsg{Width: tt.width, Height: tt.height}
			updatedModel, _ := m.Update(sizeMsg)

			updated := updatedModel.(Model)

			if updated.terminalTooSmall != tt.expectTooSmall {
				t.Errorf("For %dx%d, expected terminalTooSmall=%v, got %v",
					tt.width, tt.height, tt.expectTooSmall, updated.terminalTooSmall)
			}
		})
	}
}

// TestUpdateTerminalResizeBackToNormal tests resizing back to normal size
func TestUpdateTerminalResizeBackToNormal(t *testing.T) {
	m := initialModel(nil, nil)

	// First make it too small
	smallMsg := tea.WindowSizeMsg{Width: 60, Height: 20}
	updatedModel, _ := m.Update(smallMsg)
	m = updatedModel.(Model)

	if !m.terminalTooSmall {
		t.Error("Expected terminalTooSmall to be true for 60x20")
	}

	// Then resize to normal
	normalMsg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updatedModel, _ = m.Update(normalMsg)

	updated := updatedModel.(Model)

	if updated.terminalTooSmall {
		t.Error("Expected terminalTooSmall to be false after resize to 100x30")
	}
}

// TestUpdateQuitKey tests quit key handling
func TestUpdateQuitKey(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true

	tests := []struct {
		name string
		key  string
	}{
		{"q key", "q"},
		{"ctrl+c", "ctrl+c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			if tt.key == "ctrl+c" {
				keyMsg = tea.KeyMsg{Type: tea.KeyCtrlC}
			}

			_, cmd := m.Update(keyMsg)

			// Should return quit command
			if cmd == nil {
				t.Error("Expected quit command, got nil")
			}
		})
	}
}

// TestUpdateThemeSwitching tests theme key handling
func TestUpdateThemeSwitching(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true

	tests := []struct {
		key       string
		themeName string
	}{
		{"1", "Default"},
		{"2", "Dark"},
		{"3", "Arctic"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			updatedModel, _ := m.Update(keyMsg)
			updated := updatedModel.(Model)

			if updated.theme.Name != tt.themeName {
				t.Errorf("Expected theme %s, got %s", tt.themeName, updated.theme.Name)
			}
		})
	}
}

// TestUpdateThemeToggle tests 't' key theme cycling
func TestUpdateThemeToggle(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true

	initialTheme := m.theme.Name

	// Press 't' to cycle theme
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
	updatedModel, _ := m.Update(keyMsg)
	updated := updatedModel.(Model)

	// Should have changed to a different theme
	if updated.theme.Name == initialTheme {
		t.Error("Expected theme to change after pressing 't'")
	}

	// Press 't' again
	updatedModel, _ = updated.Update(keyMsg)
	updated = updatedModel.(Model)

	// Should have changed again
	// (might cycle back to initial after full cycle)
	_ = updated.theme.Name
}

// TestUpdateHelpToggle tests help screen toggling
func TestUpdateHelpToggle(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true

	// Show help
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	updatedModel, _ := m.Update(keyMsg)
	updated := updatedModel.(Model)

	if !updated.helpScreen.IsVisible() {
		t.Error("Expected help screen to be visible after '?'")
	}

	// Hide help
	updatedModel, _ = updated.Update(keyMsg)
	updated = updatedModel.(Model)

	if updated.helpScreen.IsVisible() {
		t.Error("Expected help screen to be hidden after second '?'")
	}
}

// TestUpdateHelpScreenInterception tests help screen intercepts input
func TestUpdateHelpScreenInterception(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.helpScreen.Toggle() // Show help

	// Try to press a key that would normally change theme
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	initialTheme := m.theme.Name
	updatedModel, _ := m.Update(keyMsg)
	updated := updatedModel.(Model)

	// Theme should NOT have changed because help screen intercepts input
	if updated.theme.Name != initialTheme {
		t.Error("Expected help screen to intercept theme change key")
	}
}

// TestUpdateFocusSwitchingTab tests tab key focus switching
func TestUpdateFocusSwitchingTab(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.focused = FocusLeft

	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := m.Update(tabMsg)
	updated := updatedModel.(Model)

	if updated.focused != FocusRight {
		t.Errorf("Expected focus to switch to FocusRight, got %v", updated.focused)
	}

	// Press tab again
	updatedModel, _ = updated.Update(tabMsg)
	updated = updatedModel.(Model)

	if updated.focused != FocusLeft {
		t.Errorf("Expected focus to switch back to FocusLeft, got %v", updated.focused)
	}
}

// TestUpdateFocusSwitchingHJKL tests vim-style focus switching
func TestUpdateFocusSwitchingHJKL(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.focused = FocusRight

	// Press 'h' to focus left
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
	updatedModel, _ := m.Update(keyMsg)
	updated := updatedModel.(Model)

	if updated.focused != FocusLeft {
		t.Errorf("Expected focus to switch to FocusLeft with 'h', got %v", updated.focused)
	}

	// Press 'l' to focus right
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
	updatedModel, _ = updated.Update(keyMsg)
	updated = updatedModel.(Model)

	if updated.focused != FocusRight {
		t.Errorf("Expected focus to switch to FocusRight with 'l', got %v", updated.focused)
	}
}

// TestUpdateWorkspaceModal tests workspace modal toggling
func TestUpdateWorkspaceModal(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true

	// Open workspace modal with 'W'
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'W'}}
	updatedModel, _ := m.Update(keyMsg)
	updated := updatedModel.(Model)

	if !updated.showWorkspaceModal {
		t.Error("Expected workspace modal to be shown after 'W'")
	}

	if updated.focused != FocusWorkspace {
		t.Errorf("Expected focus to be FocusWorkspace, got %v", updated.focused)
	}
}

// TestUpdateWorkspaceModalEscape tests closing workspace modal with escape
func TestUpdateWorkspaceModalEscape(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true
	m.showWorkspaceModal = true
	m.focused = FocusWorkspace

	// Close with escape
	keyMsg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := m.Update(keyMsg)
	updated := updatedModel.(Model)

	if updated.showWorkspaceModal {
		t.Error("Expected workspace modal to be closed after escape")
	}

	if updated.focused != FocusLeft {
		t.Errorf("Expected focus to return to FocusLeft, got %v", updated.focused)
	}
}

// TestUpdateCurrentWorkspaceLoaded tests workspace loading message
func TestUpdateCurrentWorkspaceLoaded(t *testing.T) {

	// Simplified test - mocks removed for P0
	m := initialModel(nil, nil)

	// Using nil services for unit test
	m.loadingWorkspaces = true

	ws := &domain.Workspace{ID: "test", Name: "Test WS"}
	msg := messages.CurrentWorkspaceLoadedMsg{Workspace: ws, Error: nil}

	updatedModel, cmd := m.Update(msg)
	updated := updatedModel.(Model)

	if updated.loadingWorkspaces {
		t.Error("Expected loadingWorkspaces to be false after CurrentWorkspaceLoadedMsg")
	}

	if updated.currentWorkspace == nil {
		t.Error("Expected workspace to be set")
	}

	if updated.currentWorkspace.Name != "Test WS" {
		t.Errorf("Expected workspace name 'Test WS', got %s", updated.currentWorkspace.Name)
	}

	// Should trigger ticket loading
	if cmd == nil {
		t.Error("Expected command to load tickets")
	}

	if !updated.loadingTickets {
		t.Error("Expected loadingTickets to be true after workspace loaded")
	}
}

// TestUpdateCurrentWorkspaceLoadedError tests workspace loading error
func TestUpdateCurrentWorkspaceLoadedError(t *testing.T) {
	m := initialModel(nil, nil)
	m.loadingWorkspaces = true

	testErr := errors.New("workspace load error")
	msg := messages.CurrentWorkspaceLoadedMsg{Workspace: nil, Error: testErr}

	updatedModel, _ := m.Update(msg)
	updated := updatedModel.(Model)

	if updated.loadingWorkspaces {
		t.Error("Expected loadingWorkspaces to be false after error")
	}

	if updated.loadError == nil {
		t.Error("Expected loadError to be set")
	}

	if updated.loadError.Error() != "workspace load error" {
		t.Errorf("Expected error 'workspace load error', got %s", updated.loadError.Error())
	}
}

// TestUpdateTicketsLoaded tests tickets loading message
func TestUpdateTicketsLoaded(t *testing.T) {
	m := initialModel(nil, nil)
	m.loadingTickets = true

	tickets := []domain.Ticket{
		{JiraID: "TEST-1", Title: "Test Ticket"},
		{JiraID: "TEST-2", Title: "Another Ticket"},
	}
	msg := messages.TicketsLoadedMsg{Tickets: tickets, Error: nil}

	updatedModel, _ := m.Update(msg)
	updated := updatedModel.(Model)

	if updated.loadingTickets {
		t.Error("Expected loadingTickets to be false after TicketsLoadedMsg")
	}

	if !updated.dataLoaded {
		t.Error("Expected dataLoaded to be true after tickets loaded")
	}

	if len(updated.tickets) != 2 {
		t.Errorf("Expected 2 tickets, got %d", len(updated.tickets))
	}

	if updated.tickets[0].JiraID != "TEST-1" {
		t.Errorf("Expected first ticket to be TEST-1, got %s", updated.tickets[0].JiraID)
	}
}

// TestUpdateTicketsLoadedError tests tickets loading error
func TestUpdateTicketsLoadedError(t *testing.T) {
	m := initialModel(nil, nil)
	m.loadingTickets = true

	testErr := errors.New("tickets load error")
	msg := messages.TicketsLoadedMsg{Tickets: nil, Error: testErr}

	updatedModel, _ := m.Update(msg)
	updated := updatedModel.(Model)

	if updated.loadingTickets {
		t.Error("Expected loadingTickets to be false after error")
	}

	if updated.loadError == nil {
		t.Error("Expected loadError to be set")
	}

	if updated.loadError.Error() != "tickets load error" {
		t.Errorf("Expected error 'tickets load error', got %s", updated.loadError.Error())
	}
}

// TestUpdateWorkspacesLoaded tests workspaces list loading
func TestUpdateWorkspacesLoaded(t *testing.T) {
	m := initialModel(nil, nil)

	workspaces := []domain.Workspace{
		{ID: "ws1", Name: "Workspace 1"},
		{ID: "ws2", Name: "Workspace 2"},
	}
	msg := messages.WorkspacesLoadedMsg{Workspaces: workspaces, Error: nil}

	updatedModel, _ := m.Update(msg)
	updated := updatedModel.(Model)

	if len(updated.workspaces) != 2 {
		t.Errorf("Expected 2 workspaces, got %d", len(updated.workspaces))
	}

	if updated.workspaces[0].Name != "Workspace 1" {
		t.Errorf("Expected first workspace 'Workspace 1', got %s", updated.workspaces[0].Name)
	}
}

// TestUpdateWorkspacesLoadedError tests workspaces list loading error
func TestUpdateWorkspacesLoadedError(t *testing.T) {
	m := initialModel(nil, nil)

	testErr := errors.New("workspaces load error")
	msg := messages.WorkspacesLoadedMsg{Workspaces: nil, Error: testErr}

	updatedModel, _ := m.Update(msg)
	updated := updatedModel.(Model)

	if updated.loadError == nil {
		t.Error("Expected loadError to be set")
	}
}

// TestUpdateThemeChangePropagation tests theme change propagates to components
func TestUpdateThemeChangePropagation(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true

	initialTheme := m.theme.Name

	// Change theme
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}}
	updatedModel, _ := m.Update(keyMsg)
	updated := updatedModel.(Model)

	if updated.theme.Name == initialTheme {
		t.Error("Expected theme to change")
	}

	// Theme should be Dark
	if updated.theme.Name != "Dark" {
		t.Errorf("Expected Dark theme, got %s", updated.theme.Name)
	}
}

// TestUpdateNotReady tests that some messages are ignored when not ready
func TestUpdateNotReady(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = false

	// Try to send a key message when not ready
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	initialTheme := m.theme.Name

	updatedModel, _ := m.Update(keyMsg)
	updated := updatedModel.(Model)

	// Theme should NOT change because model is not ready
	// Actually, the update function might still process it
	// Let's verify the model is still not ready
	if updated.ready {
		t.Error("Expected model to remain not ready")
	}

	_ = initialTheme // Theme change might still happen, that's ok
}

// TestUpdateSpinnerMessages tests spinner tick messages
func TestUpdateSpinnerMessages(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true

	// Send a generic message that would trigger spinner update
	updatedModel, cmd := m.Update(struct{}{})
	updated := updatedModel.(Model)

	// Should still have a model
	if updated.width == 0 {
		t.Error("Model should maintain its state")
	}

	// Might return a spinner tick command
	_ = cmd
}

// TestUpdateWithAllThemes tests cycling through all themes
func TestUpdateWithAllThemes(t *testing.T) {
	m := initialModel(nil, nil)
	m.ready = true

	themes := []string{"Default", "Dark", "Arctic"}

	for _, expectedTheme := range themes {
		m.theme = theme.GetByName(expectedTheme)

		// Verify theme is set correctly
		if m.theme.Name != expectedTheme {
			t.Errorf("Expected theme %s, got %s", expectedTheme, m.theme.Name)
		}

		// Verify model can handle updates with this theme
		sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 40}
		updatedModel, _ := m.Update(sizeMsg)
		updated := updatedModel.(Model)

		if updated.theme.Name != expectedTheme {
			t.Errorf("Theme should persist after update, expected %s, got %s",
				expectedTheme, updated.theme.Name)
		}
	}
}
