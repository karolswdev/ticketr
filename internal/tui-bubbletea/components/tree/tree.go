package tree

import (
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// TreeItem represents a flat item with hierarchy metadata.
// This allows us to use bubbles/list while maintaining tree structure.
type TreeItem struct {
	Ticket   *domain.Ticket // The underlying ticket or task
	Level    int            // Indentation level (0 = root, 1 = child, etc.)
	HasKids  bool           // Has subtasks?
	Expanded bool           // Is this node expanded?
	Parent   *TreeItem      // Parent reference (nil for root)
	Index    int            // Position in flat list
	Visible  bool           // Visible after expansion/collapse filtering?
	IsTask   bool           // Is this a task (not a ticket)?
	TaskData *domain.Task   // If IsTask is true, this contains the task data
}

// FilterValue implements list.Item interface
func (i TreeItem) FilterValue() string {
	if i.IsTask && i.TaskData != nil {
		return i.TaskData.Title
	}
	if i.Ticket != nil {
		return i.Ticket.Title
	}
	return ""
}

// TreeModel wraps bubbles/list with tree logic.
// It maintains a flat representation of the tree for efficient rendering.
type TreeModel struct {
	list          list.Model           // Bubbles list component
	items         []TreeItem           // All items (flat representation)
	visibleItems  []TreeItem           // Currently visible items (after collapse/expand)
	expandedState map[string]bool      // Track expansion state by ticket ID
	width         int                  // Component width
	height        int                  // Component height
	focused       bool                 // Is this component focused?
	onSelect      func(*domain.Ticket) // Callback when item is selected
	theme         *theme.Theme         // Current theme for styling
}

// New creates a new tree model with the given dimensions and theme.
func New(width, height int, t *theme.Theme) TreeModel {
	if t == nil {
		t = &theme.DefaultTheme
	}

	delegate := newTreeDelegate(t)
	l := list.New([]list.Item{}, delegate, width, height)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)

	return TreeModel{
		list:          l,
		items:         []TreeItem{},
		visibleItems:  []TreeItem{},
		expandedState: make(map[string]bool),
		width:         width,
		height:        height,
		focused:       false,
		theme:         t,
	}
}

// SetTickets sets the tickets to display in the tree.
// This flattens the hierarchical structure and rebuilds the visible items.
func (m *TreeModel) SetTickets(tickets []domain.Ticket) {
	m.items = FlattenTickets(tickets, m.expandedState)
	m.rebuildVisibleItems()
}

// FlattenTickets converts hierarchical tickets to flat tree items.
// It respects the expansion state to determine which items are visible.
func FlattenTickets(tickets []domain.Ticket, expandedState map[string]bool) []TreeItem {
	var items []TreeItem

	for i := range tickets {
		ticket := &tickets[i]

		// Create root level item
		item := TreeItem{
			Ticket:   ticket,
			Level:    0,
			HasKids:  len(ticket.Tasks) > 0,
			Expanded: expandedState[ticket.JiraID],
			Parent:   nil,
			Index:    len(items),
			Visible:  true,
			IsTask:   false,
		}
		items = append(items, item)

		// Add children if expanded
		if item.Expanded {
			items = append(items, flattenChildren(ticket, &items[len(items)-1], 1, expandedState)...)
		}
	}

	// Update indices
	for i := range items {
		items[i].Index = i
	}

	return items
}

// flattenChildren recursively flattens child tasks.
func flattenChildren(ticket *domain.Ticket, parent *TreeItem, level int, expandedState map[string]bool) []TreeItem {
	var items []TreeItem

	for i := range ticket.Tasks {
		task := &ticket.Tasks[i]

		// Create a pseudo-ticket for the task
		taskTicket := &domain.Ticket{
			JiraID:             task.JiraID,
			Title:              task.Title,
			Description:        task.Description,
			CustomFields:       task.CustomFields,
			AcceptanceCriteria: task.AcceptanceCriteria,
			Tasks:              []domain.Task{}, // Tasks don't have subtasks
			SourceLine:         task.SourceLine,
		}

		item := TreeItem{
			Ticket:   taskTicket,
			TaskData: task,
			Level:    level,
			HasKids:  false, // Tasks don't have children
			Expanded: false,
			Parent:   parent,
			Visible:  parent.Expanded,
			IsTask:   true,
		}
		items = append(items, item)
	}

	return items
}

// rebuildVisibleItems updates the visible items list based on expansion state.
func (m *TreeModel) rebuildVisibleItems() {
	m.visibleItems = []TreeItem{}

	for _, item := range m.items {
		if item.Visible {
			m.visibleItems = append(m.visibleItems, item)
		}
	}

	// Convert to list items and update the list
	listItems := make([]list.Item, len(m.visibleItems))
	for i, item := range m.visibleItems {
		listItems[i] = item
	}
	m.list.SetItems(listItems)
}

// ToggleExpand toggles the expansion state of the currently selected item.
func (m *TreeModel) ToggleExpand() tea.Cmd {
	if m.list.SelectedItem() == nil {
		return nil
	}

	item := m.list.SelectedItem().(TreeItem)

	if !item.HasKids {
		return nil // Nothing to expand
	}

	// Toggle expansion state
	item.Expanded = !item.Expanded
	if item.Ticket != nil {
		m.expandedState[item.Ticket.JiraID] = item.Expanded
	}

	// Rebuild the tree with new expansion state
	// We need to get the original tickets and rebuild
	m.rebuildTree()

	return nil
}

// rebuildTree rebuilds the entire tree structure from scratch.
// This is called when expansion state changes.
func (m *TreeModel) rebuildTree() {
	// Extract root tickets from items
	var rootTickets []domain.Ticket
	for _, item := range m.items {
		if item.Level == 0 && item.Ticket != nil && !item.IsTask {
			rootTickets = append(rootTickets, *item.Ticket)
		}
	}

	// Rebuild with current expansion state
	m.items = FlattenTickets(rootTickets, m.expandedState)
	m.rebuildVisibleItems()
}

// GoToParent moves the cursor to the parent of the current item.
func (m *TreeModel) GoToParent() tea.Cmd {
	if m.list.SelectedItem() == nil {
		return nil
	}

	item := m.list.SelectedItem().(TreeItem)
	if item.Parent == nil {
		return nil // Already at root
	}

	// Find parent in visible items
	for i, visItem := range m.visibleItems {
		if visItem.Index == item.Parent.Index {
			m.list.Select(i)
			break
		}
	}

	return nil
}

// SetOnSelect sets the callback function that is called when an item is selected.
func (m *TreeModel) SetOnSelect(fn func(*domain.Ticket)) {
	m.onSelect = fn
}

// GetSelected returns the currently selected ticket (if any).
func (m *TreeModel) GetSelected() *domain.Ticket {
	if m.list.SelectedItem() == nil {
		return nil
	}

	item := m.list.SelectedItem().(TreeItem)
	return item.Ticket
}

// Focus sets focus on the tree component.
func (m *TreeModel) Focus() {
	m.focused = true
}

// Blur removes focus from the tree component.
func (m *TreeModel) Blur() {
	m.focused = false
}

// Focused returns whether the tree component is focused.
func (m TreeModel) Focused() bool {
	return m.focused
}

// SetSize updates the dimensions of the tree component.
func (m *TreeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height)
}

// SetTheme updates the theme for the tree component.
// This rebuilds the delegate with new theme-aware styles.
func (m *TreeModel) SetTheme(t *theme.Theme) {
	if t == nil {
		t = &theme.DefaultTheme
	}
	m.theme = t

	// Rebuild delegate with new theme
	delegate := newTreeDelegate(t)
	m.list.SetDelegate(delegate)
}

// Init initializes the tree component.
func (m TreeModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the tree component.
func (m TreeModel) Update(msg tea.Msg) (TreeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.focused {
			return m, nil
		}

		switch msg.String() {
		case "right", "l":
			// Expand current item
			return m, m.ToggleExpand()

		case "left", "h":
			// Collapse current item OR go to parent
			if m.list.SelectedItem() != nil {
				item := m.list.SelectedItem().(TreeItem)
				if item.Expanded {
					return m, m.ToggleExpand()
				} else {
					return m, m.GoToParent()
				}
			}

		case "enter", "o":
			// Toggle expand/collapse OR trigger selection
			if m.list.SelectedItem() != nil {
				item := m.list.SelectedItem().(TreeItem)
				if item.HasKids {
					return m, m.ToggleExpand()
				} else if m.onSelect != nil && item.Ticket != nil {
					m.onSelect(item.Ticket)
				}
			}

		case "up", "k":
			m.list.CursorUp()

		case "down", "j":
			m.list.CursorDown()

		case "g":
			// Go to top
			m.list.Select(0)

		case "G":
			// Go to bottom
			if len(m.visibleItems) > 0 {
				m.list.Select(len(m.visibleItems) - 1)
			}
		}
	}

	// Delegate remaining messages to bubbles/list
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the tree component.
func (m TreeModel) View() string {
	return m.list.View()
}

// treeDelegate is the custom list.ItemDelegate for rendering tree items.
type treeDelegate struct {
	styles TreeStyles
	theme  *theme.Theme
}

// newTreeDelegate creates a new tree delegate with theme-aware styles.
func newTreeDelegate(t *theme.Theme) treeDelegate {
	if t == nil {
		t = &theme.DefaultTheme
	}
	return treeDelegate{
		styles: GetTreeStyles(t),
		theme:  t,
	}
}

// Height returns the height of each item (1 line per item).
func (d treeDelegate) Height() int {
	return 1
}

// Spacing returns the spacing between items.
func (d treeDelegate) Spacing() int {
	return 0
}

// Update handles updates for the delegate.
func (d treeDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

// Render renders a single tree item.
func (d treeDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	treeItem := item.(TreeItem)
	isSelected := index == m.Index()

	// Build the line content
	var content string

	// Add indentation
	indent := ""
	for i := 0; i < treeItem.Level; i++ {
		indent += "  "
	}
	content += indent

	// Add expand/collapse icon
	if treeItem.HasKids {
		if treeItem.Expanded {
			content += "▼ "
		} else {
			content += "▶ "
		}
	} else {
		content += "  "
	}

	// Add item type icon
	if treeItem.IsTask {
		content += "• "
	} else {
		content += "◆ "
	}

	// Add title
	title := treeItem.Ticket.Title
	if len(title) > 60 {
		title = title[:57] + "..."
	}
	content += title

	// Add JiraID if present
	if treeItem.Ticket.JiraID != "" {
		content += " "
		content += d.styles.JiraID.Render("[" + treeItem.Ticket.JiraID + "]")
	}

	// Apply styling based on selection
	var style lipgloss.Style
	if isSelected {
		style = d.styles.SelectedItem
	} else {
		style = d.styles.NormalItem
	}

	_, _ = io.WriteString(w, style.Render(content))
}
