package help

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// HelpModel represents the help screen with keyboard shortcuts.
// Week 2 Day 2: Modal overlay showing all available keyboard shortcuts.
type HelpModel struct {
	viewport viewport.Model
	width    int
	height   int
	visible  bool
	theme    *theme.Theme
	content  string
}

// ShowHelpMsg is sent to show the help modal.
type ShowHelpMsg struct{}

// HideHelpMsg is sent to hide the help modal.
type HideHelpMsg struct{}

// New creates a new help screen model.
func New(width, height int, th *theme.Theme) HelpModel {
	vp := viewport.New(width, height)
	vp.YPosition = 0

	m := HelpModel{
		viewport: vp,
		width:    width,
		height:   height,
		visible:  false,
		theme:    th,
	}

	m.updateContent()
	return m
}

// SetSize updates the dimensions of the help screen.
func (m *HelpModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width
	m.viewport.Height = height
}

// SetTheme updates the theme and regenerates content.
func (m *HelpModel) SetTheme(th *theme.Theme) {
	m.theme = th
	m.updateContent()
}

// Show displays the help modal.
func (m *HelpModel) Show() {
	m.visible = true
	m.updateContent()
	m.viewport.GotoTop()
}

// Hide hides the help modal.
func (m *HelpModel) Hide() {
	m.visible = false
}

// Toggle toggles help visibility.
func (m *HelpModel) Toggle() {
	if m.visible {
		m.Hide()
	} else {
		m.Show()
	}
}

// IsVisible returns whether the help modal is currently visible.
func (m HelpModel) IsVisible() bool {
	return m.visible
}

// Update handles messages for the help screen.
func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case ShowHelpMsg:
		m.Show()
		return m, nil

	case HideHelpMsg:
		m.Hide()
		return m, nil

	case tea.KeyMsg:
		if !m.visible {
			return m, nil
		}

		switch msg.String() {
		case "?", "esc", "q":
			m.Hide()
			return m, nil
		}
	}

	// Update viewport for scrolling (only when visible)
	if !m.visible {
		return m, nil
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the help screen.
func (m HelpModel) View() string {
	if !m.visible {
		return ""
	}

	return m.viewport.View()
}

// Init initializes the help screen.
func (m HelpModel) Init() tea.Cmd {
	return nil
}

// updateContent generates the help content with current theme.
func (m *HelpModel) updateContent() {
	titleStyle := lipgloss.NewStyle().
		Foreground(m.theme.Primary).
		Bold(true).
		Align(lipgloss.Center).
		Width(m.width)

	sectionStyle := lipgloss.NewStyle().
		Foreground(m.theme.Accent).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)

	keyStyle := lipgloss.NewStyle().
		Foreground(m.theme.Success).
		Bold(true).
		Width(15)

	descStyle := lipgloss.NewStyle().
		Foreground(m.theme.Foreground)

	helpStyle := lipgloss.NewStyle().
		Foreground(m.theme.Muted).
		Italic(true).
		MarginTop(1)

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("TICKETR - KEYBOARD SHORTCUTS"))
	b.WriteString("\n\n")

	// Navigation section
	b.WriteString(sectionStyle.Render("NAVIGATION"))
	b.WriteString("\n")
	b.WriteString(keyStyle.Render("Tab") + descStyle.Render("Switch focus between panels"))
	b.WriteString("\n")
	b.WriteString(keyStyle.Render("h") + descStyle.Render("Focus left panel (tree)"))
	b.WriteString("\n")
	b.WriteString(keyStyle.Render("l") + descStyle.Render("Focus right panel (detail)"))
	b.WriteString("\n")
	b.WriteString(keyStyle.Render("↑/↓, j/k") + descStyle.Render("Navigate up/down in lists"))
	b.WriteString("\n")
	b.WriteString(keyStyle.Render("←/→") + descStyle.Render("Collapse/expand tree nodes"))
	b.WriteString("\n")
	b.WriteString(keyStyle.Render("Enter") + descStyle.Render("Select item / show detail"))
	b.WriteString("\n")
	b.WriteString(keyStyle.Render("Esc") + descStyle.Render("Go back / close modal"))
	b.WriteString("\n")

	// Actions section
	b.WriteString(sectionStyle.Render("ACTIONS"))
	b.WriteString("\n")
	b.WriteString(keyStyle.Render("W") + descStyle.Render("Switch workspace"))
	b.WriteString("\n")
	b.WriteString(keyStyle.Render("r") + descStyle.Render("Refresh data"))
	b.WriteString("\n")
	b.WriteString(keyStyle.Render("q, Ctrl+C") + descStyle.Render("Quit application"))
	b.WriteString("\n")

	// Theme section
	b.WriteString(sectionStyle.Render("THEMES"))
	b.WriteString("\n")
	b.WriteString(keyStyle.Render("1") + descStyle.Render("Default theme (Green)"))
	b.WriteString("\n")
	b.WriteString(keyStyle.Render("2") + descStyle.Render("Dark theme (Blue)"))
	b.WriteString("\n")
	b.WriteString(keyStyle.Render("3") + descStyle.Render("Arctic theme (Cyan)"))
	b.WriteString("\n")
	b.WriteString(keyStyle.Render("t") + descStyle.Render("Cycle through themes"))
	b.WriteString("\n")

	// Help section
	b.WriteString(sectionStyle.Render("HELP"))
	b.WriteString("\n")
	b.WriteString(keyStyle.Render("?") + descStyle.Render("Toggle this help screen"))
	b.WriteString("\n")

	// Scrolling help if content is long
	if m.viewport.TotalLineCount() > m.height {
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Use ↑/↓ or j/k to scroll this help. Press ? or Esc to close."))
	} else {
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("Press ? or Esc to close."))
	}

	m.content = b.String()
	m.viewport.SetContent(m.content)
}
