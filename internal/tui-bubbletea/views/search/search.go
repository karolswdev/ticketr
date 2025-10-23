package search

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/actions"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/components/modal"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// Model represents the search modal for action search.
// Week 3 Day 2: Fuzzy search modal with action execution.
type Model struct {
	// UI components
	input textinput.Model // Search input field

	// Data
	registry *actions.Registry // Action registry reference
	results  []*actions.Action  // Filtered search results
	actionCtx *actions.ActionContext // Action context for predicates

	// State
	visible       bool         // Is modal open?
	selectedIndex int          // Currently selected result index
	width         int          // Viewport width
	height        int          // Viewport height
	theme         *theme.Theme // Current theme

	// Configuration
	maxResults int // Maximum results to display
}

// SearchModalOpenedMsg is sent when the search modal opens.
type SearchModalOpenedMsg struct{}

// SearchModalClosedMsg is sent when the search modal closes.
type SearchModalClosedMsg struct{}

// ActionExecuteRequestedMsg is sent when an action is selected for execution.
type ActionExecuteRequestedMsg struct {
	ActionID actions.ActionID
	Action   *actions.Action
}

// New creates a new search modal.
func New(registry *actions.Registry, t *theme.Theme) Model {
	input := textinput.New()
	input.Placeholder = "Search actions..."
	input.Focus()
	input.CharLimit = 100

	return Model{
		input:         input,
		registry:      registry,
		results:       []*actions.Action{},
		visible:       false,
		selectedIndex: 0,
		width:         80,
		height:        24,
		theme:         t,
		maxResults:    10,
	}
}

// Init initializes the search modal.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the search modal.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.visible {
			return m, nil
		}

		switch msg.String() {
		case "esc":
			// Close modal
			return m.Close()

		case "enter":
			// Execute selected action
			if len(m.results) > 0 && m.selectedIndex < len(m.results) {
				selectedAction := m.results[m.selectedIndex]
				m, closeCmd := m.Close()
				executeCmd := func() tea.Msg {
					return ActionExecuteRequestedMsg{
						ActionID: selectedAction.ID,
						Action:   selectedAction,
					}
				}
				return m, tea.Batch(closeCmd, executeCmd)
			}
			return m, nil

		case "up", "k":
			// Navigate up
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
			return m, nil

		case "down", "j":
			// Navigate down
			if m.selectedIndex < len(m.results)-1 {
				m.selectedIndex++
			}
			return m, nil

		default:
			// Update input and perform search
			m.input, cmd = m.input.Update(msg)
			m.performSearch()
			m.selectedIndex = 0 // Reset selection when query changes
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

// View renders the search modal.
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	// Build modal content
	var content strings.Builder

	// Title
	palette := theme.GetPaletteForTheme(m.theme)
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(palette.Primary)).
		Bold(true)

	content.WriteString(titleStyle.Render("🔍 Search Actions"))
	content.WriteString("\n\n")

	// Search input
	content.WriteString(m.input.View())
	content.WriteString("\n\n")

	// Results
	if len(m.results) == 0 {
		// Empty state
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.Muted)).
			Italic(true)

		if m.input.Value() == "" {
			content.WriteString(emptyStyle.Render("Type to search for actions..."))
		} else {
			content.WriteString(emptyStyle.Render("No actions found"))
		}
	} else {
		// Render results (max 10)
		displayCount := min(len(m.results), m.maxResults)
		for i := 0; i < displayCount; i++ {
			action := m.results[i]
			m.renderActionItem(&content, action, i == m.selectedIndex)
		}

		// Show count if more results available
		if len(m.results) > m.maxResults {
			moreStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(palette.Muted)).
				Italic(true)
			content.WriteString("\n")
			content.WriteString(moreStyle.Render(fmt.Sprintf("... and %d more results", len(m.results)-m.maxResults)))
		}
	}

	// Help text
	content.WriteString("\n\n")
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(palette.Muted)).
		Italic(true)
	content.WriteString(helpStyle.Render("↑/↓ or j/k: Navigate • Enter: Execute • Esc: Close"))

	// Calculate modal dimensions (40% width, 60% height)
	modalWidth := max(m.width*4/10, 40)
	modalHeight := max(m.height*6/10, 15)

	// Ensure modal fits within screen
	if modalWidth > m.width-4 {
		modalWidth = m.width - 4
	}
	if modalHeight > m.height-4 {
		modalHeight = m.height - 4
	}

	// Wrap in modal overlay
	contentStr := content.String()

	// Trim content if too long for modal
	lines := strings.Split(contentStr, "\n")
	if len(lines) > modalHeight-4 { // Account for padding and border
		lines = lines[:modalHeight-4]
		contentStr = strings.Join(lines, "\n")
	}

	return modal.Render(contentStr, m.width, m.height, m.theme)
}

// Open shows the search modal and focuses the input.
func (m Model) Open() (Model, tea.Cmd) {
	m.visible = true
	m.input.Focus()
	m.input.SetValue("")
	m.results = []*actions.Action{}
	m.selectedIndex = 0

	// Perform initial search with empty query to show all actions
	m.performSearch()

	return m, func() tea.Msg {
		return SearchModalOpenedMsg{}
	}
}

// Close hides the search modal.
func (m Model) Close() (Model, tea.Cmd) {
	m.visible = false
	m.input.Blur()
	m.input.SetValue("")
	m.results = []*actions.Action{}
	m.selectedIndex = 0

	return m, func() tea.Msg {
		return SearchModalClosedMsg{}
	}
}

// IsVisible returns whether the modal is currently visible.
func (m Model) IsVisible() bool {
	return m.visible
}

// SetSize updates the modal dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetTheme updates the theme.
func (m *Model) SetTheme(t *theme.Theme) {
	m.theme = t
}

// SetActionContext updates the action context for predicate evaluation.
func (m *Model) SetActionContext(actx *actions.ActionContext) {
	m.actionCtx = actx
}

// performSearch executes the search query against the action registry.
func (m *Model) performSearch() {
	query := m.input.Value()

	if m.actionCtx == nil {
		// If no action context, create a minimal one
		m.actionCtx = &actions.ActionContext{
			Context: actions.ContextGlobal,
			Width:   m.width,
			Height:  m.height,
		}
	}

	if query == "" {
		// Empty query: show all available actions
		m.results = m.registry.ActionsForContext(actions.ContextGlobal, m.actionCtx)
	} else {
		// Perform fuzzy search
		m.results = m.registry.Search(query, m.actionCtx)
	}

	// Ensure selected index is within bounds
	if m.selectedIndex >= len(m.results) {
		m.selectedIndex = max(0, len(m.results)-1)
	}
}

// renderActionItem renders a single action in the results list.
func (m *Model) renderActionItem(content *strings.Builder, action *actions.Action, selected bool) {
	palette := theme.GetPaletteForTheme(m.theme)

	var style lipgloss.Style
	var prefix string

	if selected {
		// Selected item styling
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.Primary)).
			Background(lipgloss.Color(palette.Selection)).
			Bold(true)
		prefix = "> "
	} else {
		// Normal item styling
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.Foreground))
		prefix = "  "
	}

	// Format: > Icon Name - Description
	icon := action.Icon
	if icon == "" {
		icon = "•"
	}

	descStyle := style.Copy().Foreground(lipgloss.Color(palette.Muted))

	// Build line
	line := fmt.Sprintf("%s%s %s", prefix, icon, action.Name)
	if action.Description != "" {
		line += descStyle.Render(" - " + action.Description)
	}

	content.WriteString(style.Render(line))
	content.WriteString("\n")
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
