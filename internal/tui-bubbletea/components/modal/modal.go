package modal

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// Render wraps content in a centered modal overlay with a dimmed background.
// The modal is displayed in the center of the screen with a border.
// Now accepts a theme parameter to ensure consistent styling across all themes.
func Render(content string, width, height int, t *theme.Theme) string {
	// Use default theme if none provided
	if t == nil {
		t = &theme.DefaultTheme
	}

	palette := theme.GetPaletteForTheme(t)

	// Calculate modal dimensions (60% of screen size, minimum 40x15)
	modalWidth := max(width*6/10, 40)
	modalHeight := max(height*6/10, 15)

	// Ensure modal fits within screen
	if modalWidth > width-4 {
		modalWidth = width - 4
	}
	if modalHeight > height-4 {
		modalHeight = height - 4
	}

	// Create the modal box with border using theme colors
	modalStyle := lipgloss.NewStyle().
		Width(modalWidth).
		Height(modalHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(palette.BorderFocused)).
		Padding(1, 2).
		Background(lipgloss.Color(palette.Background)).
		Foreground(lipgloss.Color(palette.Foreground))

	modal := modalStyle.Render(content)

	// Center the modal on screen with semi-transparent background
	// Use theme's muted color for backdrop
	positioned := lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceChars("░"),
		lipgloss.WithWhitespaceForeground(lipgloss.Color(palette.Muted)),
	)

	return positioned
}

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
