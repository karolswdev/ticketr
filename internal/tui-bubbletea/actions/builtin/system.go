package builtin

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/actions"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/actions/predicates"
)

// Message types
type showHelpMsg struct{}

// RegisterSystemActions registers system-level actions
func RegisterSystemActions(reg *actions.Registry) error {
	// Quit
	if err := reg.Register(&actions.Action{
		ID:          "system.quit",
		Name:        "Quit",
		Description: "Quit the application",
		Category:    actions.CategorySystem,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{
			{Key: "q"},
			{Key: "c", Ctrl: true},
		},
		Predicate: predicates.Always(),
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return tea.Quit
		},
		Tags: []string{"quit", "exit", "close"},
		Icon: "🚪",
	}); err != nil {
		return err
	}

	// Show help
	if err := reg.Register(&actions.Action{
		ID:          "system.help",
		Name:        "Show Help",
		Description: "Show help screen",
		Category:    actions.CategorySystem,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{
			{Key: "?"},
		},
		Predicate: predicates.Always(),
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			// Return message to show help
			return func() tea.Msg {
				return showHelpMsg{}
			}
		},
		Tags: []string{"help", "shortcuts", "keys"},
		Icon: "❓",
	}); err != nil {
		return err
	}

	return nil
}
