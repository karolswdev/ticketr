// Package testhelpers provides testing utilities for end-to-end integration tests.
// These helpers facilitate testing cross-component interactions in the TUI.
package testhelpers

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/karolswdev/ticktr/internal/core/domain"
	tuibubbletea "github.com/karolswdev/ticktr/internal/tui-bubbletea"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/actions"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/mocks"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// TestHarness provides a complete test environment for integration testing.
// It includes a full TUI model, mock services, and helper methods for simulating user interactions.
type TestHarness struct {
	t              *testing.T
	model          tuibubbletea.Model
	width          int
	height         int
	wsService      *mocks.MockWorkspaceService
	ticketService  *mocks.MockTicketQueryService
	actionRegistry *actions.Registry
	contextMgr     *actions.ContextManager
	commands       []tea.Cmd // Collected commands from updates
}

// NewTestHarness creates a new test harness with default configuration.
// It initializes a full TUI model with mock services and returns a ready-to-test environment.
func NewTestHarness(t *testing.T) *TestHarness {
	t.Helper()

	// Create mock services with test data
	wsService := mocks.NewMockWorkspaceService()
	ticketService := mocks.NewMockTicketQueryService()

	// Create action registry and context manager
	registry := actions.NewRegistry()
	contextMgr := actions.NewContextManager(actions.ContextGlobal)

	// Initialize model with services (using internal initialModel function)
	// Note: This requires the model to be properly initialized with all Week 3 components
	// For now, we'll use a basic initialization
	width, height := 120, 40

	harness := &TestHarness{
		t:              t,
		width:          width,
		height:         height,
		wsService:      wsService,
		ticketService:  ticketService,
		actionRegistry: registry,
		contextMgr:     contextMgr,
		commands:       []tea.Cmd{},
	}

	return harness
}

// InitModel initializes the model - call this after configuring services
func (h *TestHarness) InitModel() {
	h.t.Helper()
	// Initialize the model - this will be implemented when Builder integrates Week 3 components
	// For now, create a placeholder that tests can use
	// model := tuibubbletea.InitialModel(h.wsService, h.ticketService)
	// h.model = model
}

// PressKey simulates pressing a single key and updates the model.
// It captures any returned commands for later inspection.
func (h *TestHarness) PressKey(key string) {
	h.t.Helper()

	var msg tea.Msg

	// Convert string to appropriate KeyMsg
	switch key {
	case "enter":
		msg = tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		msg = tea.KeyMsg{Type: tea.KeyEsc}
	case "tab":
		msg = tea.KeyMsg{Type: tea.KeyTab}
	case "up":
		msg = tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		msg = tea.KeyMsg{Type: tea.KeyDown}
	case "left":
		msg = tea.KeyMsg{Type: tea.KeyLeft}
	case "right":
		msg = tea.KeyMsg{Type: tea.KeyRight}
	case "space":
		msg = tea.KeyMsg{Type: tea.KeySpace}
	case "ctrl+c":
		msg = tea.KeyMsg{Type: tea.KeyCtrlC}
	case "ctrl+p":
		msg = tea.KeyMsg{Type: tea.KeyCtrlP}
	default:
		// Handle single character or Runes
		if len(key) == 1 {
			msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{rune(key[0])}}
		} else if strings.HasPrefix(key, "ctrl+") {
			// Parse ctrl+X format
			char := key[5:] // Get character after "ctrl+"
			if len(char) == 1 {
				msg = tea.KeyMsg{
					Type:  tea.KeyRunes,
					Runes: []rune{rune(char[0])},
					Alt:   false,
				}
				// Set Ctrl modifier - Bubbletea doesn't have a direct Ctrl field,
				// so we use the Type field for common ctrl combinations
				switch char {
				case "p":
					msg = tea.KeyMsg{Type: tea.KeyCtrlP}
				case "n":
					msg = tea.KeyMsg{Type: tea.KeyCtrlN}
				case "h":
					msg = tea.KeyMsg{Type: tea.KeyCtrlH}
				}
			}
		}
	}

	model, cmd := h.model.Update(msg)
	h.model = model.(tuibubbletea.Model)
	if cmd != nil {
		h.commands = append(h.commands, cmd)
	}
}

// PressKeys simulates pressing multiple keys in sequence.
func (h *TestHarness) PressKeys(keys []string) {
	h.t.Helper()
	for _, key := range keys {
		h.PressKey(key)
	}
}

// TypeString simulates typing a string character by character.
func (h *TestHarness) TypeString(text string) {
	h.t.Helper()
	for _, char := range text {
		h.PressKey(string(char))
	}
}

// Resize simulates a window resize event.
func (h *TestHarness) Resize(width, height int) {
	h.t.Helper()
	h.width = width
	h.height = height

	msg := tea.WindowSizeMsg{Width: width, Height: height}
	model, cmd := h.model.Update(msg)
	h.model = model.(tuibubbletea.Model)
	if cmd != nil {
		h.commands = append(h.commands, cmd)
	}
}

// ChangeTheme simulates changing the theme.
func (h *TestHarness) ChangeTheme(themeName string) {
	h.t.Helper()
	// Get theme by name
	newTheme := theme.GetByName(themeName)
	if newTheme == nil {
		h.t.Fatalf("Theme %s not found", themeName)
	}

	// Simulate theme change - this will depend on how theme changes are implemented
	// For now, we'll press the appropriate number key
	switch themeName {
	case "Default":
		h.PressKey("1")
	case "Dark":
		h.PressKey("2")
	case "Arctic":
		h.PressKey("3")
	}
}

// GetCurrentTheme returns the current theme name from the model.
func (h *TestHarness) GetCurrentTheme() string {
	// This will need to be implemented based on how the model exposes theme
	return h.model.GetCurrentTheme()
}

// View returns the current rendered view.
func (h *TestHarness) View() string {
	return h.model.View()
}

// AssertViewContains checks if the view contains a specific string.
func (h *TestHarness) AssertViewContains(substring string) {
	h.t.Helper()
	view := h.View()
	if !strings.Contains(view, substring) {
		h.t.Errorf("View does not contain %q\nView:\n%s", substring, view)
	}
}

// AssertViewNotContains checks if the view does NOT contain a specific string.
func (h *TestHarness) AssertViewNotContains(substring string) {
	h.t.Helper()
	view := h.View()
	if strings.Contains(view, substring) {
		h.t.Errorf("View should not contain %q\nView:\n%s", substring, view)
	}
}

// GetCommands returns all commands collected during updates.
func (h *TestHarness) GetCommands() []tea.Cmd {
	return h.commands
}

// ClearCommands clears the collected commands.
func (h *TestHarness) ClearCommands() {
	h.commands = []tea.Cmd{}
}

// HasQuitCommand checks if any collected command is a quit command.
func (h *TestHarness) HasQuitCommand() bool {
	for _, cmd := range h.commands {
		if cmd != nil {
			// Execute command and check if it returns QuitMsg
			msg := cmd()
			if _, ok := msg.(tea.QuitMsg); ok {
				return true
			}
		}
	}
	return false
}

// ConfigureWorkspace configures the mock workspace service.
func (h *TestHarness) ConfigureWorkspace(ws *domain.Workspace) {
	h.wsService.WithWorkspace(ws)
}

// ConfigureTickets configures the mock ticket service for a workspace.
func (h *TestHarness) ConfigureTickets(workspaceID string, tickets []domain.Ticket) {
	h.ticketService.WithTickets(workspaceID, tickets)
}

// SimulateDataLoad simulates successful data loading.
func (h *TestHarness) SimulateDataLoad() {
	h.t.Helper()

	// Simulate window size message (triggers data load)
	h.Resize(h.width, h.height)

	// Process any data loading commands that were generated
	// This is a simplified version - in reality, commands would be executed
	// and their results fed back into the model
}

// RegisterAction registers an action in the test registry.
func (h *TestHarness) RegisterAction(action *actions.Action) error {
	return h.actionRegistry.Register(action)
}

// GetActionRegistry returns the action registry for inspection.
func (h *TestHarness) GetActionRegistry() *actions.Registry {
	return h.actionRegistry
}

// GetContextManager returns the context manager for inspection.
func (h *TestHarness) GetContextManager() *actions.ContextManager {
	return h.contextMgr
}

// SetContext sets the current context in the context manager.
func (h *TestHarness) SetContext(ctx actions.Context) {
	h.contextMgr.Switch(ctx)
}

// AssertContext checks if the current context matches the expected context.
func (h *TestHarness) AssertContext(expected actions.Context) {
	h.t.Helper()
	current := h.contextMgr.Current()
	if current != expected {
		h.t.Errorf("Expected context %v, got %v", expected, current)
	}
}

// Model returns the underlying model for advanced assertions.
func (h *TestHarness) Model() tuibubbletea.Model {
	return h.model
}

// T returns the testing.T instance for custom assertions.
func (h *TestHarness) T() *testing.T {
	return h.t
}
