package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// RenderError renders an error message with red styling.
// This is used to display errors from failed data loading operations.
func RenderError(err error) string {
	if err == nil {
		return ""
	}

	errorStyle := lipgloss.NewStyle().
		Foreground(&theme.DefaultTheme.Error).
		Bold(true).
		MarginTop(1).
		MarginLeft(2)

	return errorStyle.Render("Error: " + err.Error())
}

// RenderCenteredError renders a centered error message.
// Useful for full-screen error states.
func RenderCenteredError(err error, width, height int) string {
	if err == nil {
		return ""
	}

	errorStyle := lipgloss.NewStyle().
		Foreground(&theme.DefaultTheme.Error).
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(&theme.DefaultTheme.Error).
		Padding(1, 2)

	content := errorStyle.Render("Error: " + err.Error())

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
