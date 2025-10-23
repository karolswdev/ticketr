package help

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/actions"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// HelpModel represents the help screen with keyboard shortcuts.
// Week 3 Day 4: Context-aware help using action registry.
type HelpModel struct {
	viewport   viewport.Model
	width      int
	height     int
	visible    bool
	theme      *theme.Theme
	content    string

	// NEW: Action integration
	registry   *actions.Registry
	contextMgr *actions.ContextManager
	actionCtx  *actions.ActionContext

	// Computed data (refreshed on show)
	sections []HelpSection
}

// HelpSection represents a categorized group of actions.
type HelpSection struct {
	Title   string
	Actions []ActionBinding
}

// ActionBinding represents a single keybinding display.
type ActionBinding struct {
	Keys        []string
	Description string
	ActionID    actions.ActionID
}

// ShowHelpMsg is sent to show the help modal.
type ShowHelpMsg struct{}

// HideHelpMsg is sent to hide the help modal.
type HideHelpMsg struct{}

// New creates a new help screen model with action registry integration.
func New(width, height int, th *theme.Theme, registry *actions.Registry, contextMgr *actions.ContextManager) HelpModel {
	vp := viewport.New(width, height)
	vp.YPosition = 0

	m := HelpModel{
		viewport:   vp,
		width:      width,
		height:     height,
		visible:    false,
		theme:      th,
		registry:   registry,
		contextMgr: contextMgr,
		sections:   []HelpSection{},
	}

	// Generate initial content (fallback mode)
	m.updateContent()
	return m
}

// NewLegacy creates a help screen without action registry (backward compatibility).
func NewLegacy(width, height int, th *theme.Theme) HelpModel {
	vp := viewport.New(width, height)
	vp.YPosition = 0

	m := HelpModel{
		viewport: vp,
		width:    width,
		height:   height,
		visible:  false,
		theme:    th,
		sections: []HelpSection{},
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

// SetActionContext updates the action context for dynamic help.
func (m *HelpModel) SetActionContext(actx *actions.ActionContext) {
	m.actionCtx = actx
	if m.visible {
		m.refreshSections()
		m.updateContent()
	}
}

// Show displays the help modal.
func (m *HelpModel) Show() {
	m.visible = true
	m.refreshSections()
	m.updateContent()
	m.viewport.GotoTop()
}

// ShowWithContext displays the help modal with specific action context.
func (m *HelpModel) ShowWithContext(actx *actions.ActionContext) {
	m.actionCtx = actx
	m.Show()
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

// refreshSections generates help sections from action registry.
func (m *HelpModel) refreshSections() {
	// If no registry, use fallback static content
	if m.registry == nil {
		m.sections = m.generateFallbackSections()
		return
	}

	// Get current context
	currentCtx := actions.ContextGlobal
	if m.contextMgr != nil {
		currentCtx = m.contextMgr.Current()
	}

	// Build action context if not provided
	actx := m.actionCtx
	if actx == nil {
		actx = &actions.ActionContext{
			Context: currentCtx,
			Width:   m.width,
			Height:  m.height,
		}
	}

	// Get all available actions for current context
	availableActions := m.registry.ActionsForContext(currentCtx, actx)

	// Group actions by category
	categoryMap := make(map[actions.ActionCategory][]*actions.Action)
	for _, action := range availableActions {
		// Skip actions that shouldn't show in UI
		if action.ShowInUI != nil && !action.ShowInUI(actx) {
			continue
		}

		category := action.Category
		if category == "" {
			category = "Other"
		}
		categoryMap[category] = append(categoryMap[category], action)
	}

	// Convert to HelpSections with sorted categories
	var sections []HelpSection
	categoryOrder := []actions.ActionCategory{
		actions.CategoryNavigation,
		actions.CategoryView,
		actions.CategoryEdit,
		actions.CategoryWorkspace,
		actions.CategorySync,
		actions.CategoryBulk,
		actions.CategorySystem,
		"Other",
	}

	for _, category := range categoryOrder {
		categoryActions, exists := categoryMap[category]
		if !exists || len(categoryActions) == 0 {
			continue
		}

		// Sort actions by name within category
		sort.Slice(categoryActions, func(i, j int) bool {
			return categoryActions[i].Name < categoryActions[j].Name
		})

		// Build action bindings
		var bindings []ActionBinding
		for _, action := range categoryActions {
			keys := m.formatKeybindings(action.Keybindings)
			if len(keys) == 0 {
				// Action with no keybindings, show just the name
				keys = []string{"-"}
			}

			bindings = append(bindings, ActionBinding{
				Keys:        keys,
				Description: action.Description,
				ActionID:    action.ID,
			})
		}

		sections = append(sections, HelpSection{
			Title:   string(category),
			Actions: bindings,
		})
	}

	m.sections = sections
}

// formatKeybindings converts KeyPattern array to display strings.
func (m *HelpModel) formatKeybindings(patterns []actions.KeyPattern) []string {
	var keys []string
	for _, pattern := range patterns {
		keys = append(keys, m.formatKeyPattern(pattern))
	}
	return keys
}

// formatKeyPattern converts a single KeyPattern to display string.
func (m *HelpModel) formatKeyPattern(pattern actions.KeyPattern) string {
	var parts []string

	if pattern.Ctrl {
		parts = append(parts, "Ctrl")
	}
	if pattern.Alt {
		parts = append(parts, "Alt")
	}
	if pattern.Shift {
		parts = append(parts, "Shift")
	}

	// Format key name
	keyName := pattern.Key
	switch strings.ToLower(keyName) {
	case "enter":
		keyName = "Enter"
	case "esc", "escape":
		keyName = "Esc"
	case "tab":
		keyName = "Tab"
	case "up":
		keyName = "↑"
	case "down":
		keyName = "↓"
	case "left":
		keyName = "←"
	case "right":
		keyName = "→"
	case "space":
		keyName = "Space"
	case "backspace":
		keyName = "Backspace"
	case "delete":
		keyName = "Delete"
	}

	parts = append(parts, keyName)

	if len(parts) > 1 {
		return strings.Join(parts, "+")
	}
	return keyName
}

// generateFallbackSections creates static sections when no registry available.
func (m *HelpModel) generateFallbackSections() []HelpSection {
	return []HelpSection{
		{
			Title: "Navigation",
			Actions: []ActionBinding{
				{Keys: []string{"Tab"}, Description: "Switch focus between panels"},
				{Keys: []string{"h"}, Description: "Focus left panel (tree)"},
				{Keys: []string{"l"}, Description: "Focus right panel (detail)"},
				{Keys: []string{"↑", "k"}, Description: "Navigate up"},
				{Keys: []string{"↓", "j"}, Description: "Navigate down"},
				{Keys: []string{"←"}, Description: "Collapse tree node"},
				{Keys: []string{"→"}, Description: "Expand tree node"},
				{Keys: []string{"Enter"}, Description: "Select item / show detail"},
				{Keys: []string{"Esc"}, Description: "Go back / close modal"},
			},
		},
		{
			Title: "Actions",
			Actions: []ActionBinding{
				{Keys: []string{"W"}, Description: "Switch workspace"},
				{Keys: []string{"r"}, Description: "Refresh data"},
			},
		},
		{
			Title: "Themes",
			Actions: []ActionBinding{
				{Keys: []string{"1"}, Description: "Default theme (Green)"},
				{Keys: []string{"2"}, Description: "Dark theme (Blue)"},
				{Keys: []string{"3"}, Description: "Arctic theme (Cyan)"},
				{Keys: []string{"t"}, Description: "Cycle through themes"},
			},
		},
		{
			Title: "System",
			Actions: []ActionBinding{
				{Keys: []string{"?"}, Description: "Toggle this help screen"},
				{Keys: []string{"q", "Ctrl+C"}, Description: "Quit application"},
			},
		},
	}
}

// updateContent generates the help content with current theme and sections.
func (m *HelpModel) updateContent() {
	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(m.theme.Primary).
		Bold(true).
		Align(lipgloss.Center).
		Width(m.width)

	contextStyle := lipgloss.NewStyle().
		Foreground(m.theme.Muted).
		Italic(true).
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
		Width(18)

	descStyle := lipgloss.NewStyle().
		Foreground(m.theme.Foreground)

	helpStyle := lipgloss.NewStyle().
		Foreground(m.theme.Muted).
		Italic(true).
		MarginTop(1)

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("TICKETR - KEYBOARD SHORTCUTS"))
	b.WriteString("\n")

	// Show current context if available
	if m.contextMgr != nil {
		currentCtx := m.contextMgr.Current()
		contextName := m.formatContextName(currentCtx)
		b.WriteString(contextStyle.Render(fmt.Sprintf("Context: %s", contextName)))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Render each section
	for _, section := range m.sections {
		b.WriteString(sectionStyle.Render(strings.ToUpper(section.Title)))
		b.WriteString("\n")

		for _, binding := range section.Actions {
			// Format keys
			keyStr := strings.Join(binding.Keys, ", ")
			b.WriteString(keyStyle.Render(keyStr) + descStyle.Render(binding.Description))
			b.WriteString("\n")
		}
	}

	// Scrolling help if content is long
	b.WriteString("\n")
	if m.viewport.TotalLineCount() > m.height {
		b.WriteString(helpStyle.Render("Use ↑/↓ or j/k to scroll this help. Press ? or Esc to close."))
	} else {
		b.WriteString(helpStyle.Render("Press ? or Esc to close."))
	}

	m.content = b.String()
	m.viewport.SetContent(m.content)
}

// formatContextName converts context constant to human-readable name.
func (m *HelpModel) formatContextName(ctx actions.Context) string {
	switch ctx {
	case actions.ContextWorkspaceList:
		return "Workspace Selector"
	case actions.ContextTicketTree:
		return "Ticket Tree"
	case actions.ContextTicketDetail:
		return "Ticket Detail"
	case actions.ContextSearch:
		return "Search"
	case actions.ContextCommandPalette:
		return "Command Palette"
	case actions.ContextModal:
		return "Modal"
	case actions.ContextSyncing:
		return "Syncing"
	case actions.ContextHelp:
		return "Help"
	case actions.ContextGlobal:
		return "Global"
	default:
		return string(ctx)
	}
}
