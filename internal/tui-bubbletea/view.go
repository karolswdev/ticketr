package tuibubbletea

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/components"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/components/modal"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// View renders the entire TUI to a string.
// This is the V in TEA (The Elm Architecture) - it's a pure function that
// takes the current model state and returns a string representation.
//
// The view is called on every state change and must be fast (<16ms for 60 FPS).
func (m Model) View() string {
	// Week 2 Day 2: Check terminal size first
	if m.terminalTooSmall {
		return m.renderTerminalTooSmallError()
	}

	if !m.ready {
		// Show loading state until we receive the first WindowSizeMsg
		return m.renderLoading()
	}

	// Week 2 Day 1: Show errors if data loading failed
	if m.loadError != nil {
		return components.RenderCenteredError(m.loadError, m.width, m.height)
	}

	// Week 2 Day 1: Show loading state while fetching data
	// Week 2 Day 2: Now uses animated spinner
	if !m.dataLoaded {
		// Center the animated spinner
		spinnerView := m.loadingSpinner.View()
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			spinnerView,
		)
	}

	// Week 2 Days 2-3: FlexBox layout with tree and detail components
	header := m.renderHeader()
	leftPanel := m.renderLeftPanel()
	rightPanel := m.renderRightPanel()
	actionBar := m.renderActionBar()

	// Use CompleteLayout to render everything
	mainView := m.layout.Render(header, leftPanel, rightPanel, actionBar)

	// Week 2 Day 2: Overlay help screen (highest priority)
	// Week 3 Day 3: Now passes theme to modal for consistent styling
	if m.helpScreen.IsVisible() {
		helpContent := m.helpScreen.View()
		return modal.Render(helpContent, m.width, m.height, m.theme)
	}

	// Week 2 Days 4-5: Overlay workspace modal
	// Week 3 Day 3: Now passes theme to modal for consistent styling
	if m.showWorkspaceModal {
		modalContent := m.workspaceSelector.View()
		return modal.Render(modalContent, m.width, m.height, m.theme)
	}

	return mainView
}

// renderLoading shows a centered "Loading..." message.
// This is displayed before we receive the initial WindowSizeMsg.
// Week 2 Day 2: Now uses animated spinner.
func (m Model) renderLoading() string {
	spinnerView := m.loadingSpinner.View()

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		spinnerView,
	)
}

// renderHeader renders the status bar / header
// Week 2 Days 2-3: Now displays workspace name, ticket count, and focus state
func (m Model) renderHeader() string {
	currentTheme := m.theme.Name
	title := "TICKETR v3.1 - Bubbletea POC (Week 2 Days 2-3)"

	// Build status line with real data
	var workspaceName string
	if m.currentWorkspace != nil {
		workspaceName = m.currentWorkspace.Name
	} else {
		workspaceName = "No workspace"
	}

	status := fmt.Sprintf("Workspace: %s | Tickets: %d | Theme: %s | Focus: %s",
		workspaceName, len(m.tickets), currentTheme, m.getFocusName())

	titleStyle := theme.TitleStyle(m.theme)
	headerStyle := theme.HeaderStyle(m.theme)

	headerContent := lipgloss.JoinHorizontal(
		lipgloss.Left,
		titleStyle.Render(title),
		lipgloss.NewStyle().Render(" | "),
		theme.InfoStyle(m.theme).Render(status),
	)

	return headerStyle.Width(m.width).Render(headerContent)
}

// renderLeftPanel renders the left panel content
// Week 2 Days 2-3: Now displays ticket tree component
func (m Model) renderLeftPanel() string {
	leftWidth, _, contentHeight := m.layout.GetPanelDimensions()

	// Determine border style based on focus
	var borderStyle lipgloss.Style
	if m.focused == FocusLeft {
		borderStyle = theme.BorderStyleFocused(m.theme)
	} else {
		borderStyle = theme.BorderStyleBlurred(m.theme)
	}

	// Panel title showing workspace and ticket count
	var titleText string
	if m.currentWorkspace != nil {
		titleText = fmt.Sprintf("📋 %s - %d Tickets", m.currentWorkspace.Name, len(m.tickets))
	} else {
		titleText = fmt.Sprintf("📋 Ticket Tree (%d)", len(m.tickets))
	}
	title := theme.PanelTitleStyle(m.theme).Render(titleText)

	// Week 2 Days 2-3: Render tree component
	treeContent := m.ticketTree.View()

	// Help text for tree navigation
	helpText := theme.HelpTextStyle(m.theme).Render(
		"↑↓/jk: Navigate | →/l: Expand | ←/h: Collapse | Enter: Select | Tab: Switch | W: Workspace",
	)

	// Combine title, tree, and help
	panelContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		treeContent,
		"",
		helpText,
	)

	// Apply border and size
	return borderStyle.
		Width(leftWidth - 2).
		Height(contentHeight - 2).
		Render(panelContent)
}

// renderRightPanel renders the right panel content (ticket detail view)
// Week 2 Days 2-3: Now uses detail view component
func (m Model) renderRightPanel() string {
	_, rightWidth, contentHeight := m.layout.GetPanelDimensions()

	// Determine border style based on focus
	var borderStyle lipgloss.Style
	if m.focused == FocusRight {
		borderStyle = theme.BorderStyleFocused(m.theme)
	} else {
		borderStyle = theme.BorderStyleBlurred(m.theme)
	}

	// Panel title
	title := theme.PanelTitleStyle(m.theme).Render("📄 Ticket Detail")

	// Week 2 Days 2-3: Use detail view component
	detailContent := m.detailView.View()

	// Combine title and content
	panelContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		detailContent,
	)

	// Apply border and size
	return borderStyle.
		Width(rightWidth - 2).
		Height(contentHeight - 2).
		Render(panelContent)
}

// renderActionBar renders the action bar / footer
// Week 2 Days 2-3: Added tree navigation keybindings
func (m Model) renderActionBar() string {
	keybindingStyle := theme.KeybindingStyle(m.theme)
	labelStyle := theme.KeybindingLabelStyle(m.theme)

	actions := []string{
		keybindingStyle.Render("[↑↓/jk]") + " " + labelStyle.Render("Navigate"),
		keybindingStyle.Render("[→/l]") + " " + labelStyle.Render("Expand"),
		keybindingStyle.Render("[←/h]") + " " + labelStyle.Render("Collapse"),
		keybindingStyle.Render("[Enter]") + " " + labelStyle.Render("Select"),
		keybindingStyle.Render("[Tab]") + " " + labelStyle.Render("Switch"),
		keybindingStyle.Render("[W]") + " " + labelStyle.Render("Workspace"),
		keybindingStyle.Render("[?]") + " " + labelStyle.Render("Help"),
		keybindingStyle.Render("[q]") + " " + labelStyle.Render("Quit"),
	}

	actionContent := lipgloss.JoinHorizontal(
		lipgloss.Left,
		actions[0]+"  ",
		actions[1]+"  ",
		actions[2]+"  ",
		actions[3]+"  ",
		actions[4]+"  ",
		actions[5]+"  ",
		actions[6]+"  ",
		actions[7],
	)

	return theme.FooterStyle(m.theme).Width(m.width).Render(actionContent)
}

// getFocusName returns a human-readable name for the current focus
// Week 2 Days 2-3: Added FocusWorkspace case
func (m Model) getFocusName() string {
	switch m.focused {
	case FocusLeft:
		return "Tree"
	case FocusRight:
		return "Detail"
	case FocusWorkspace:
		return "Workspace Selector"
	default:
		return "Unknown"
	}
}

// renderTerminalTooSmallError renders an error message when terminal is too small.
// Week 2 Day 2: Minimum terminal size is 80x24.
func (m Model) renderTerminalTooSmallError() string {
	msg := fmt.Sprintf(
		"Terminal Too Small!\n\n"+
			"Current: %d\u00d7%d\n"+
			"Required: 80\u00d724 minimum\n\n"+
			"Please resize your terminal to continue.",
		m.width, m.height,
	)

	errorStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("#E06C75")).
		Bold(true)

	return errorStyle.Render(msg)
}
