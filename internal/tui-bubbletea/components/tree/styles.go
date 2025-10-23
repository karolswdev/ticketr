package tree

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// TreeStyles defines the styling for tree components.
// This is now theme-aware and accepts a theme parameter.
type TreeStyles struct {
	SelectedItem lipgloss.Style
	NormalItem   lipgloss.Style
	JiraID       lipgloss.Style
	ExpandIcon   lipgloss.Style
	TicketIcon   lipgloss.Style
	TaskIcon     lipgloss.Style
}

// GetTreeStyles returns theme-aware styling for tree components.
// All tree colors now derive from the active theme.
func GetTreeStyles(t *theme.Theme) TreeStyles {
	palette := theme.GetPaletteForTheme(t)

	return TreeStyles{
		SelectedItem: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.Primary)).
			Background(lipgloss.Color(palette.Selection)).
			Bold(true),
		NormalItem: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.Foreground)),
		JiraID: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.Muted)).
			Italic(true),
		ExpandIcon: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.Accent)),
		TicketIcon: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.Primary)),
		TaskIcon: lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.Secondary)),
	}
}
