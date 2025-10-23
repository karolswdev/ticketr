package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// ActionBar represents the bottom action bar with keyboard shortcuts
type ActionBar struct {
	Actions []Action
	Width   int
}

// Action represents a keyboard shortcut
type Action struct {
	Key         string
	Description string
}

// NewActionBar creates a new action bar
func NewActionBar() *ActionBar {
	return &ActionBar{
		Actions: []Action{
			{"F1", "Help"},
			{"F2", "Sync"},
			{"F3", "Workspace"},
			{"F5", "Refresh"},
			{"Tab", "Focus"},
			{"1/2/3", "Theme"},
			{"q", "Quit"},
		},
	}
}

// SetActions sets custom actions
func (a *ActionBar) SetActions(actions []Action) {
	a.Actions = actions
}

// Render renders the action bar
func (a *ActionBar) Render() string {
	th := &theme.DefaultTheme

	// Create action bar style
	barStyle := lipgloss.NewStyle().
		Foreground(th.Foreground).
		Background(th.Background).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(th.Border).
		Padding(0, 1).
		Width(a.Width - 2)

	// Render each action
	keyStyle := lipgloss.NewStyle().
		Foreground(th.Primary).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(th.Secondary)

	var parts []string
	for _, action := range a.Actions {
		key := keyStyle.Render(action.Key)
		desc := descStyle.Render(action.Description)
		parts = append(parts, key+": "+desc)
	}

	content := strings.Join(parts, " | ")

	return barStyle.Render(content)
}
