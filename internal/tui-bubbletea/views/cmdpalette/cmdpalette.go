package cmdpalette

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/actions"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/components/modal"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// Model represents the command palette for quick action access.
// Week 3 Day 3: Enhanced command palette with categories, recent actions, and keybindings.
type Model struct {
	// UI components
	input textinput.Model // Command input field

	// Data
	registry   *actions.Registry       // Action registry reference
	contextMgr *actions.ContextManager // Context manager
	actionCtx  *actions.ActionContext  // Action context for predicates
	results    []ActionItem            // Filtered and sorted results
	recent     []actions.ActionID      // Recent action IDs (max 5)

	// State
	visible       bool                   // Is palette open?
	filterMode    FilterMode             // All, Category, Recent
	selectedCat   actions.ActionCategory // Selected category filter
	selectedIndex int                    // Currently selected result index
	width         int                    // Viewport width
	height        int                    // Viewport height
	theme         *theme.Theme           // Current theme

	// Configuration
	maxResults int // Maximum results to display
	maxRecent  int // Maximum recent actions to track
}

// ActionItem represents an action in the command palette with enhanced metadata.
type ActionItem struct {
	Action      *actions.Action
	Keybindings string // Formatted: "Ctrl+S", "j, ↓"
	Category    string
	IsRecent    bool // Show star indicator
}

// FilterMode determines how actions are filtered.
type FilterMode int

const (
	FilterAll FilterMode = iota
	FilterCategory
	FilterRecent
)

// CommandPaletteOpenedMsg is sent when the command palette opens.
type CommandPaletteOpenedMsg struct{}

// CommandPaletteClosedMsg is sent when the command palette closes.
type CommandPaletteClosedMsg struct{}

// CommandExecutedMsg is sent when an action is executed from the palette.
type CommandExecutedMsg struct {
	ActionID actions.ActionID
	Action   *actions.Action
}

// New creates a new command palette.
func New(registry *actions.Registry, contextMgr *actions.ContextManager, t *theme.Theme) Model {
	input := textinput.New()
	input.Placeholder = "Type to search actions..."
	input.Focus()
	input.CharLimit = 100

	return Model{
		input:         input,
		registry:      registry,
		contextMgr:    contextMgr,
		results:       []ActionItem{},
		recent:        []actions.ActionID{},
		visible:       false,
		filterMode:    FilterAll,
		selectedCat:   "",
		selectedIndex: 0,
		width:         80,
		height:        24,
		theme:         t,
		maxResults:    20,
		maxRecent:     5,
	}
}

// Init initializes the command palette.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the command palette.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.visible {
			return m, nil
		}

		// Handle category filter shortcuts (Ctrl+0-7)
		// Check for Ctrl+0 through Ctrl+7 for category filtering
		if len(msg.String()) >= 6 && strings.HasPrefix(msg.String(), "ctrl+") {
			digit := msg.String()[5:]
			if digit == "0" {
				m.ClearFilter()
				return m, nil
			}
			switch digit {
			case "1":
				m.SetCategoryFilter(actions.CategoryNavigation)
				return m, nil
			case "2":
				m.SetCategoryFilter(actions.CategoryView)
				return m, nil
			case "3":
				m.SetCategoryFilter(actions.CategoryEdit)
				return m, nil
			case "4":
				m.SetCategoryFilter(actions.CategoryWorkspace)
				return m, nil
			case "5":
				m.SetCategoryFilter(actions.CategorySync)
				return m, nil
			case "6":
				m.SetCategoryFilter(actions.CategoryBulk)
				return m, nil
			case "7":
				m.SetCategoryFilter(actions.CategorySystem)
				return m, nil
			}
		}

		switch msg.String() {
		case "esc":
			// Close palette
			return m.Close()

		case "enter":
			// Execute selected action
			if len(m.results) > 0 && m.selectedIndex < len(m.results) {
				selectedItem := m.results[m.selectedIndex]
				m.AddRecent(selectedItem.Action.ID)
				m, closeCmd := m.Close()
				executeCmd := func() tea.Msg {
					return CommandExecutedMsg{
						ActionID: selectedItem.Action.ID,
						Action:   selectedItem.Action,
					}
				}
				return m, tea.Batch(closeCmd, executeCmd)
			}
			return m, nil

		case "up", "k":
			// Navigate up
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
			return m, nil

		case "down", "j":
			// Navigate down
			if m.selectedIndex < len(m.results)-1 {
				m.selectedIndex++
			}
			return m, nil

		default:
			// Update input and perform search
			m.input, cmd = m.input.Update(msg)
			m.performSearch()
			m.selectedIndex = 0 // Reset selection when query changes
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

// View renders the command palette.
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	// Build palette content
	var content strings.Builder

	// Title
	palette := theme.GetPaletteForTheme(m.theme)
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(palette.Primary)).
		Bold(true)

	filterInfo := ""
	switch m.filterMode {
	case FilterCategory:
		filterInfo = fmt.Sprintf(" - %s", m.selectedCat)
	case FilterRecent:
		filterInfo = " - Recent"
	}

	content.WriteString(titleStyle.Render("🎯 Command Palette" + filterInfo))
	content.WriteString("  ")

	// Show keybinding hint
	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(palette.Muted)).
		Italic(true)
	content.WriteString(hintStyle.Render("[Ctrl+P]"))
	content.WriteString("\n\n")

	// Search input
	content.WriteString(m.input.View())
	content.WriteString("\n\n")

	// Results
	if len(m.results) == 0 {
		// Empty state
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.Muted)).
			Italic(true)

		if m.input.Value() == "" {
			content.WriteString(emptyStyle.Render("Type to search for actions or press Ctrl+H for help..."))
		} else {
			content.WriteString(emptyStyle.Render("No actions found (try Ctrl+H for help)"))
		}
	} else {
		m.renderResults(&content, palette)
	}

	// Footer with counts and help
	content.WriteString("\n\n")
	m.renderFooter(&content, palette)

	// Calculate modal dimensions (60% width, 70% height - larger than search modal)
	modalWidth := max(m.width*6/10, 60)
	modalHeight := max(m.height*7/10, 20)

	// Ensure modal fits within screen
	if modalWidth > m.width-4 {
		modalWidth = m.width - 4
	}
	if modalHeight > m.height-4 {
		modalHeight = m.height - 4
	}

	// Wrap in modal overlay
	contentStr := content.String()

	// Trim content if too long for modal
	lines := strings.Split(contentStr, "\n")
	if len(lines) > modalHeight-4 { // Account for padding and border
		lines = lines[:modalHeight-4]
		contentStr = strings.Join(lines, "\n")
	}

	return modal.Render(contentStr, m.width, m.height, m.theme)
}

// Open shows the command palette and focuses the input.
func (m Model) Open() (Model, tea.Cmd) {
	m.visible = true
	m.input.Focus()
	m.input.SetValue("")
	m.filterMode = FilterAll
	m.selectedCat = ""
	m.selectedIndex = 0

	// Perform initial search with empty query to show all actions
	m.performSearch()

	return m, func() tea.Msg {
		return CommandPaletteOpenedMsg{}
	}
}

// Close hides the command palette.
func (m Model) Close() (Model, tea.Cmd) {
	m.visible = false
	m.input.Blur()
	m.input.SetValue("")
	m.results = []ActionItem{}
	m.filterMode = FilterAll
	m.selectedCat = ""
	m.selectedIndex = 0

	return m, func() tea.Msg {
		return CommandPaletteClosedMsg{}
	}
}

// IsVisible returns whether the palette is currently visible.
func (m Model) IsVisible() bool {
	return m.visible
}

// SetSize updates the palette dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetTheme updates the theme.
func (m *Model) SetTheme(t *theme.Theme) {
	m.theme = t
}

// SetActionContext updates the action context for predicate evaluation.
func (m *Model) SetActionContext(actx *actions.ActionContext) {
	m.actionCtx = actx
}

// AddRecent adds an action to the recent list.
func (m *Model) AddRecent(actionID actions.ActionID) {
	// Remove if already in list
	for i, id := range m.recent {
		if id == actionID {
			m.recent = append(m.recent[:i], m.recent[i+1:]...)
			break
		}
	}

	// Add to front
	m.recent = append([]actions.ActionID{actionID}, m.recent...)

	// Keep max recent items
	if len(m.recent) > m.maxRecent {
		m.recent = m.recent[:m.maxRecent]
	}
}

// SetCategoryFilter sets the category filter.
func (m *Model) SetCategoryFilter(category actions.ActionCategory) {
	m.filterMode = FilterCategory
	m.selectedCat = category
	m.performSearch()
	m.selectedIndex = 0
}

// ClearFilter clears the category filter.
func (m *Model) ClearFilter() {
	m.filterMode = FilterAll
	m.selectedCat = ""
	m.performSearch()
	m.selectedIndex = 0
}

// performSearch executes the search query against the action registry.
func (m *Model) performSearch() {
	query := m.input.Value()

	// Build action context if not provided
	if m.actionCtx == nil {
		currentCtx := actions.ContextGlobal
		if m.contextMgr != nil {
			currentCtx = m.contextMgr.Current()
		}
		m.actionCtx = &actions.ActionContext{
			Context: currentCtx,
			Width:   m.width,
			Height:  m.height,
		}
	}

	var rawResults []*actions.Action

	if query == "" {
		// Empty query: show all available actions
		rawResults = m.registry.ActionsForContext(actions.ContextGlobal, m.actionCtx)
	} else {
		// Perform fuzzy search
		rawResults = m.registry.Search(query, m.actionCtx)
	}

	// Apply category filter if active
	if m.filterMode == FilterCategory {
		var filtered []*actions.Action
		for _, action := range rawResults {
			if action.Category == m.selectedCat {
				filtered = append(filtered, action)
			}
		}
		rawResults = filtered
	}

	// Convert to ActionItems with enhanced metadata
	m.results = m.buildActionItems(rawResults)

	// Sort results: recent first, then by relevance
	m.sortResults(query)

	// Ensure selected index is within bounds
	if m.selectedIndex >= len(m.results) {
		m.selectedIndex = max(0, len(m.results)-1)
	}
}

// buildActionItems converts actions to ActionItems with metadata.
func (m *Model) buildActionItems(rawActions []*actions.Action) []ActionItem {
	items := make([]ActionItem, 0, len(rawActions))

	for _, action := range rawActions {
		// Check if this action is in recent list
		isRecent := false
		for _, recentID := range m.recent {
			if recentID == action.ID {
				isRecent = true
				break
			}
		}

		items = append(items, ActionItem{
			Action:      action,
			Keybindings: m.formatKeybindings(action.Keybindings),
			Category:    string(action.Category),
			IsRecent:    isRecent,
		})
	}

	return items
}

// sortResults sorts action items by recent, relevance, and name.
func (m *Model) sortResults(query string) {
	query = strings.ToLower(query)

	sort.Slice(m.results, func(i, j int) bool {
		itemI := m.results[i]
		itemJ := m.results[j]

		// Recent actions come first
		if itemI.IsRecent != itemJ.IsRecent {
			return itemI.IsRecent
		}

		// If query is not empty, sort by relevance
		if query != "" {
			iNameMatch := strings.Contains(strings.ToLower(itemI.Action.Name), query)
			jNameMatch := strings.Contains(strings.ToLower(itemJ.Action.Name), query)

			if iNameMatch != jNameMatch {
				return iNameMatch
			}
		}

		// Sort alphabetically within same relevance
		return itemI.Action.Name < itemJ.Action.Name
	})
}

// formatKeybindings converts keybindings to display string.
func (m *Model) formatKeybindings(patterns []actions.KeyPattern) string {
	if len(patterns) == 0 {
		return ""
	}

	var keys []string
	for _, pattern := range patterns {
		keys = append(keys, m.formatKeyPattern(pattern))
	}

	return strings.Join(keys, ", ")
}

// formatKeyPattern converts a single KeyPattern to display string.
func (m *Model) formatKeyPattern(pattern actions.KeyPattern) string {
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
	}

	parts = append(parts, keyName)

	if len(parts) > 1 {
		return strings.Join(parts, "+")
	}
	return keyName
}

// renderResults renders the action items with categories and keybindings.
func (m *Model) renderResults(content *strings.Builder, palette theme.Palette) {
	displayCount := min(len(m.results), m.maxResults)

	// Group results by category for display
	categoryGroups := m.groupByCategory(m.results[:displayCount])

	// Render each category group
	for i, group := range categoryGroups {
		if i > 0 {
			content.WriteString("\n")
		}

		// Render category header
		m.renderCategoryHeader(content, palette, group.Category, group.IsRecent)

		// Render actions in this category
		for _, item := range group.Items {
			actualIndex := m.findItemIndex(item)
			m.renderActionItem(content, palette, item, actualIndex == m.selectedIndex)
		}
	}

	// Show count if more results available
	if len(m.results) > m.maxResults {
		content.WriteString("\n")
		moreStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.Muted)).
			Italic(true)
		content.WriteString(moreStyle.Render(fmt.Sprintf("... and %d more results", len(m.results)-m.maxResults)))
	}
}

// CategoryGroup represents a group of actions by category.
type CategoryGroup struct {
	Category string
	Items    []ActionItem
	IsRecent bool
}

// groupByCategory groups action items by category, with recent at top.
func (m *Model) groupByCategory(items []ActionItem) []CategoryGroup {
	groups := []CategoryGroup{}
	categoryMap := make(map[string][]ActionItem)
	hasRecent := false

	// Separate recent items
	var recentItems []ActionItem
	var otherItems []ActionItem
	for _, item := range items {
		if item.IsRecent {
			recentItems = append(recentItems, item)
			hasRecent = true
		} else {
			otherItems = append(otherItems, item)
		}
	}

	// Add recent group if exists
	if hasRecent && m.filterMode == FilterAll {
		groups = append(groups, CategoryGroup{
			Category: "RECENT",
			Items:    recentItems,
			IsRecent: true,
		})
	}

	// Group other items by category
	for _, item := range otherItems {
		category := item.Category
		if category == "" {
			category = "Other"
		}
		categoryMap[category] = append(categoryMap[category], item)
	}

	// Add category groups in order
	categoryOrder := []string{
		string(actions.CategoryNavigation),
		string(actions.CategoryView),
		string(actions.CategoryEdit),
		string(actions.CategoryWorkspace),
		string(actions.CategorySync),
		string(actions.CategoryBulk),
		string(actions.CategorySystem),
		"Other",
	}

	for _, category := range categoryOrder {
		if items, exists := categoryMap[category]; exists {
			groups = append(groups, CategoryGroup{
				Category: category,
				Items:    items,
				IsRecent: false,
			})
		}
	}

	return groups
}

// findItemIndex finds the actual index of an item in results.
func (m *Model) findItemIndex(item ActionItem) int {
	for i, result := range m.results {
		if result.Action.ID == item.Action.ID {
			return i
		}
	}
	return -1
}

// renderCategoryHeader renders a category header.
func (m *Model) renderCategoryHeader(content *strings.Builder, palette theme.Palette, category string, isRecent bool) {
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(palette.Muted)).
		Bold(false)

	var header string
	if isRecent {
		header = "⭐ " + category
	} else {
		header = "── " + strings.ToUpper(category) + " ──"
	}

	content.WriteString(headerStyle.Render(header))
	content.WriteString("\n")
}

// renderActionItem renders a single action item.
func (m *Model) renderActionItem(content *strings.Builder, palette theme.Palette, item ActionItem, selected bool) {
	var style lipgloss.Style
	var prefix string

	if selected {
		// Selected item styling
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.Primary)).
			Background(lipgloss.Color(palette.Selection)).
			Bold(true)
		prefix = "> "
	} else {
		// Normal item styling
		style = lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.Foreground))
		prefix = "  "
	}

	// Icon
	icon := item.Action.Icon
	if icon == "" {
		icon = "•"
	}

	// Build main line with name
	line := fmt.Sprintf("%s%s %s", prefix, icon, item.Action.Name)

	// Add keybindings on the right (dimmed)
	if item.Keybindings != "" {
		keyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(palette.Muted))

		// Calculate padding to align keybindings to the right
		// For simplicity, just append with spacing
		line += "  "
		line += keyStyle.Render(item.Keybindings)
	}

	content.WriteString(style.Render(line))
	content.WriteString("\n")

	// Add description on next line if selected
	if selected && item.Action.Description != "" {
		descStyle := style.Copy().
			Foreground(lipgloss.Color(palette.Muted)).
			Italic(true)
		descLine := fmt.Sprintf("    %s", item.Action.Description)
		content.WriteString(descStyle.Render(descLine))
		content.WriteString("\n")
	}
}

// renderFooter renders the footer with counts and help.
func (m *Model) renderFooter(content *strings.Builder, palette theme.Palette) {
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(palette.Muted)).
		Italic(true)

	// Count information
	countInfo := ""
	if len(m.results) > 0 {
		categoryCount := m.countCategories(m.results)
		countInfo = fmt.Sprintf("Showing %d actions in %d categories", len(m.results), categoryCount)
		content.WriteString(helpStyle.Render(countInfo))
		content.WriteString("\n")
	}

	// Help text
	helpText := "↑/↓ or j/k: Navigate  |  Enter: Execute  |  Esc: Close  |  Ctrl+0-7: Filter by category"
	content.WriteString(helpStyle.Render(helpText))
}

// countCategories counts unique categories in results.
func (m *Model) countCategories(items []ActionItem) int {
	categories := make(map[string]bool)
	for _, item := range items {
		category := item.Category
		if category == "" {
			category = "Other"
		}
		categories[category] = true
	}
	return len(categories)
}

// GetRecentActions returns the recent action IDs.
func (m Model) GetRecentActions() []actions.ActionID {
	return m.recent
}

// SetRecentActions sets the recent action IDs (for persistence).
func (m *Model) SetRecentActions(recent []actions.ActionID) {
	if len(recent) > m.maxRecent {
		recent = recent[:m.maxRecent]
	}
	m.recent = recent
}

// min returns the smaller of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
