package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// LoadingModel represents an animated loading spinner with message.
// Week 2 Day 2: Migrated to use bubbles spinner for proper animation.
type LoadingModel struct {
	spinner spinner.Model
	message string
	theme   *theme.Theme
}

// NewLoading creates a new loading spinner with default configuration.
func NewLoading(message string, th *theme.Theme) LoadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(th.Primary)

	return LoadingModel{
		spinner: s,
		message: message,
		theme:   th,
	}
}

// SetMessage updates the loading message.
func (m *LoadingModel) SetMessage(message string) {
	m.message = message
}

// SetTheme updates the theme for the spinner.
func (m *LoadingModel) SetTheme(th *theme.Theme) {
	m.theme = th
	m.spinner.Style = lipgloss.NewStyle().Foreground(th.Primary)
}

// Update handles messages for the loading spinner.
func (m LoadingModel) Update(msg tea.Msg) (LoadingModel, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// View renders the loading spinner with message.
func (m LoadingModel) View() string {
	messageStyle := lipgloss.NewStyle().
		Foreground(m.theme.Foreground).
		Bold(true)

	return m.spinner.View() + " " + messageStyle.Render(m.message)
}

// Init initializes the spinner (starts the animation).
func (m LoadingModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// RenderLoading renders a loading message with the given text.
// This is used while asynchronously loading data from services.
// Week 2 Day 2: Now uses bubbles spinner but provides static fallback for compatibility.
func RenderLoading(message string) string {
	style := theme.InfoStyle(&theme.DefaultTheme).
		Bold(true).
		MarginTop(2).
		MarginLeft(2)

	spinner := "◐"
	return style.Render(spinner + " " + message)
}

// RenderCenteredLoading renders a centered loading message.
// Useful for full-screen loading states.
// Week 2 Day 2: Static version for compatibility - use LoadingModel for animated version.
func RenderCenteredLoading(message string, width, height int) string {
	spinner := "◐"
	content := theme.InfoStyle(&theme.DefaultTheme).
		Bold(true).
		Render(spinner + " " + message)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
