package tuibubbletea

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/actions"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/messages"
)

// RegisterBuiltinActions registers all built-in TUI actions with the action registry.
// Week 4 Day 1: Complete action registration for all core TUI functionality.
func RegisterBuiltinActions(reg *actions.Registry) error {
	// ===== SYSTEM ACTIONS =====

	if err := reg.Register(&actions.Action{
		ID:          "system.quit",
		Name:        "Quit Application",
		Description: "Exit Ticketr",
		Category:    actions.CategorySystem,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{
			{Key: "q"},
			{Ctrl: true, Key: "c"},
		},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return tea.Quit
		},
		Icon: "🚪",
		Tags: []string{"quit", "exit", "close"},
	}); err != nil {
		return err
	}

	if err := reg.Register(&actions.Action{
		ID:          "system.help",
		Name:        "Toggle Help",
		Description: "Show/hide keyboard shortcuts",
		Category:    actions.CategorySystem,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "?"}},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.ToggleHelpMsg{}
			}
		},
		Icon: "❓",
		Tags: []string{"help", "shortcuts", "keybindings"},
	}); err != nil {
		return err
	}

	if err := reg.Register(&actions.Action{
		ID:          "system.search",
		Name:        "Open Search",
		Description: "Search for actions",
		Category:    actions.CategorySystem,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "/"}},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.OpenSearchMsg{}
			}
		},
		Icon: "🔍",
		Tags: []string{"search", "find", "actions"},
	}); err != nil {
		return err
	}

	if err := reg.Register(&actions.Action{
		ID:          "system.command_palette",
		Name:        "Open Command Palette",
		Description: "Quick access to all actions",
		Category:    actions.CategorySystem,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{
			{Ctrl: true, Key: "p"},
			{Key: ":"},
		},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.OpenCommandPaletteMsg{}
			}
		},
		Icon: "🎯",
		Tags: []string{"command", "palette", "actions"},
	}); err != nil {
		return err
	}

	// ===== NAVIGATION ACTIONS =====

	if err := reg.Register(&actions.Action{
		ID:          "nav.move_down",
		Name:        "Move Down",
		Description: "Move selection down",
		Category:    actions.CategoryNavigation,
		Contexts:    []actions.Context{actions.ContextTicketTree},
		Keybindings: []actions.KeyPattern{
			{Key: "j"},
			{Key: "down"},
		},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.MoveCursorDownMsg{}
			}
		},
		Icon: "⬇️",
		Tags: []string{"down", "navigate"},
	}); err != nil {
		return err
	}

	if err := reg.Register(&actions.Action{
		ID:          "nav.move_up",
		Name:        "Move Up",
		Description: "Move selection up",
		Category:    actions.CategoryNavigation,
		Contexts:    []actions.Context{actions.ContextTicketTree},
		Keybindings: []actions.KeyPattern{
			{Key: "k"},
			{Key: "up"},
		},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.MoveCursorUpMsg{}
			}
		},
		Icon: "⬆️",
		Tags: []string{"up", "navigate"},
	}); err != nil {
		return err
	}

	if err := reg.Register(&actions.Action{
		ID:          "nav.expand_node",
		Name:        "Expand Node",
		Description: "Expand selected tree node",
		Category:    actions.CategoryNavigation,
		Contexts:    []actions.Context{actions.ContextTicketTree},
		Keybindings: []actions.KeyPattern{
			{Key: "l"},
			{Key: "right"},
		},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.ExpandNodeMsg{}
			}
		},
		Icon: "➡️",
		Tags: []string{"expand", "open"},
	}); err != nil {
		return err
	}

	if err := reg.Register(&actions.Action{
		ID:          "nav.collapse_node",
		Name:        "Collapse Node",
		Description: "Collapse selected tree node",
		Category:    actions.CategoryNavigation,
		Contexts:    []actions.Context{actions.ContextTicketTree},
		Keybindings: []actions.KeyPattern{
			{Key: "h"},
			{Key: "left"},
		},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.CollapseNodeMsg{}
			}
		},
		Icon: "⬅️",
		Tags: []string{"collapse", "close"},
	}); err != nil {
		return err
	}

	if err := reg.Register(&actions.Action{
		ID:          "nav.switch_panel",
		Name:        "Switch Panel",
		Description: "Switch focus between panels",
		Category:    actions.CategoryNavigation,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "tab"}},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.SwitchPanelMsg{}
			}
		},
		Icon: "↔️",
		Tags: []string{"switch", "focus", "panel"},
	}); err != nil {
		return err
	}

	if err := reg.Register(&actions.Action{
		ID:          "nav.focus_left",
		Name:        "Focus Left Panel",
		Description: "Move focus to left panel (tree)",
		Category:    actions.CategoryNavigation,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "h"}},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.FocusLeftMsg{}
			}
		},
		Icon: "⬅️",
		Tags: []string{"focus", "left", "tree"},
	}); err != nil {
		return err
	}

	if err := reg.Register(&actions.Action{
		ID:          "nav.focus_right",
		Name:        "Focus Right Panel",
		Description: "Move focus to right panel (detail)",
		Category:    actions.CategoryNavigation,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "l"}},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.FocusRightMsg{}
			}
		},
		Icon: "➡️",
		Tags: []string{"focus", "right", "detail"},
	}); err != nil {
		return err
	}

	// ===== VIEW ACTIONS =====

	if err := reg.Register(&actions.Action{
		ID:          "view.select_ticket",
		Name:        "View Ticket",
		Description: "View selected ticket details",
		Category:    actions.CategoryView,
		Contexts:    []actions.Context{actions.ContextTicketTree},
		Keybindings: []actions.KeyPattern{{Key: "enter"}},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.SelectTicketMsg{}
			}
		},
		Icon: "📄",
		Tags: []string{"view", "ticket", "select"},
	}); err != nil {
		return err
	}

	if err := reg.Register(&actions.Action{
		ID:          "view.refresh",
		Name:        "Refresh Data",
		Description: "Reload tickets from workspace",
		Category:    actions.CategoryView,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "r"}},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.RefreshDataMsg{}
			}
		},
		Icon: "🔄",
		Tags: []string{"refresh", "reload"},
	}); err != nil {
		return err
	}

	// ===== WORKSPACE ACTIONS =====

	if err := reg.Register(&actions.Action{
		ID:          "workspace.switch",
		Name:        "Switch Workspace",
		Description: "Open workspace selector",
		Category:    actions.CategoryWorkspace,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Shift: true, Key: "w"}},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.OpenWorkspaceSelectorMsg{}
			}
		},
		Icon: "🗂️",
		Tags: []string{"workspace", "switch"},
	}); err != nil {
		return err
	}

	// ===== THEME ACTIONS =====

	if err := reg.Register(&actions.Action{
		ID:          "theme.cycle",
		Name:        "Cycle Theme",
		Description: "Switch to next theme",
		Category:    actions.CategorySystem,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "t"}},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.CycleThemeMsg{}
			}
		},
		Icon: "🎨",
		Tags: []string{"theme", "color"},
	}); err != nil {
		return err
	}

	if err := reg.Register(&actions.Action{
		ID:          "theme.default",
		Name:        "Default Theme",
		Description: "Switch to Default (Green) theme",
		Category:    actions.CategorySystem,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "1"}},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.SetThemeMsg{ThemeName: "Default"}
			}
		},
		Icon: "🟢",
		Tags: []string{"theme", "default", "green"},
	}); err != nil {
		return err
	}

	if err := reg.Register(&actions.Action{
		ID:          "theme.dark",
		Name:        "Dark Theme",
		Description: "Switch to Dark (Blue) theme",
		Category:    actions.CategorySystem,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "2"}},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.SetThemeMsg{ThemeName: "Dark"}
			}
		},
		Icon: "🔵",
		Tags: []string{"theme", "dark", "blue"},
	}); err != nil {
		return err
	}

	if err := reg.Register(&actions.Action{
		ID:          "theme.arctic",
		Name:        "Arctic Theme",
		Description: "Switch to Arctic (Cyan) theme",
		Category:    actions.CategorySystem,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "3"}},
		Execute: func(ctx *actions.ActionContext) tea.Cmd {
			return func() tea.Msg {
				return messages.SetThemeMsg{ThemeName: "Arctic"}
			}
		},
		Icon: "🔷",
		Tags: []string{"theme", "arctic", "cyan"},
	}); err != nil {
		return err
	}

	return nil
}
