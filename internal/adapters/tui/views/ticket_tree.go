package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/karolswdev/ticktr/internal/adapters/tui/theme"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/rivo/tview"
)

// TicketTreeView displays a hierarchical tree of tickets.
type TicketTreeView struct {
	tree             *tview.TreeView
	root             *tview.TreeNode
	workspace        *services.WorkspaceService
	ticketQuery      *services.TicketQueryService
	onTicketSelected func(*domain.Ticket) // Callback for ticket selection
	app              *tview.Application   // Reference to app for async UI updates
	isLoading        bool                 // Track loading state
	selectedTickets  map[string]bool      // Track selected tickets (ticketID -> selected)
	selectionMode    bool                 // True when any tickets are selected
}

// NewTicketTreeView creates a new ticket tree view.
func NewTicketTreeView(workspace *services.WorkspaceService, ticketQuery *services.TicketQueryService, app *tview.Application) *TicketTreeView {
	root := tview.NewTreeNode("Tickets")
	root.SetColor(tcell.ColorYellow)

	tree := tview.NewTreeView()
	tree.SetRoot(root)
	tree.SetCurrentNode(root)
	tree.SetBorder(true).SetTitle(" Tickets ")

	view := &TicketTreeView{
		tree:            tree,
		root:            root,
		workspace:       workspace,
		ticketQuery:     ticketQuery,
		app:             app,
		isLoading:       false,
		selectedTickets: make(map[string]bool),
		selectionMode:   false,
	}

	// Setup vim-style navigation
	view.setupKeyBindings()

	// Setup selection handler for Enter key
	view.setupSelectionHandler()

	// Load tickets asynchronously on startup
	view.loadInitialTicketsAsync()

	return view
}

// Name returns the view identifier.
func (v *TicketTreeView) Name() string {
	return "ticket_tree"
}

// Primitive returns the tview primitive.
func (v *TicketTreeView) Primitive() tview.Primitive {
	return v.tree
}

// OnShow is called when the view becomes active.
func (v *TicketTreeView) OnShow() {
	// No special behavior needed
}

// OnHide is called when the view is hidden.
func (v *TicketTreeView) OnHide() {
	// No cleanup needed
}

// setupKeyBindings configures vim-style keyboard shortcuts.
func (v *TicketTreeView) setupKeyBindings() {
	v.tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'j':
				// Move down
				return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
			case 'k':
				// Move up
				return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
			case 'h':
				// Collapse node
				node := v.tree.GetCurrentNode()
				if node != nil {
					node.SetExpanded(false)
				}
				return nil
			case 'l':
				// Expand node
				node := v.tree.GetCurrentNode()
				if node != nil {
					node.SetExpanded(true)
				}
				return nil
			case ' ':
				// Toggle ticket selection (Space bar)
				v.toggleCurrentSelection()
				return nil
			case 'a':
				// Select all visible tickets
				v.selectAllVisible()
				return nil
			case 'A':
				// Deselect all tickets
				v.clearSelection()
				return nil
			}
		}
		// Arrow keys still work (backward compatibility)
		return event
	})
}

// setupSelectionHandler configures Enter key to trigger ticket detail view.
func (v *TicketTreeView) setupSelectionHandler() {
	v.tree.SetSelectedFunc(func(node *tview.TreeNode) {
		if node == nil {
			return
		}

		ref := node.GetReference()
		if ref == nil {
			return
		}

		// Check if reference is a ticket (not a task)
		if ticket, ok := ref.(domain.Ticket); ok {
			if v.onTicketSelected != nil {
				v.onTicketSelected(&ticket)
			}
		}
	})
}

// SetOnTicketSelected sets callback for when a ticket is selected.
func (v *TicketTreeView) SetOnTicketSelected(callback func(*domain.Ticket)) {
	v.onTicketSelected = callback
}

// SetFocused updates border color when focus changes.
func (v *TicketTreeView) SetFocused(focused bool) {
	color := theme.GetSecondaryColor()
	if focused {
		color = theme.GetPrimaryColor()
	}
	v.tree.SetBorderColor(color)
}

// loadInitialTicketsAsync loads tickets for the current workspace asynchronously.
func (v *TicketTreeView) loadInitialTicketsAsync() {
	// Get current workspace
	workspace, err := v.workspace.Current()
	if err != nil || workspace == nil {
		v.showMessage("No active workspace", tcell.ColorYellow)
		return
	}

	v.LoadTicketsAsync(workspace.ID)
}

// LoadTicketsAsync loads tickets for the specified workspace asynchronously.
// This is the primary method for loading tickets - it runs in a goroutine to avoid blocking the TUI.
func (v *TicketTreeView) LoadTicketsAsync(workspaceID string) {
	// Prevent multiple concurrent loads
	if v.isLoading {
		return
	}
	v.isLoading = true

	// Show loading state immediately
	v.showMessage("Loading tickets...", tcell.ColorBlue)

	// Fetch tickets in goroutine
	go func() {
		tickets, err := v.ticketQuery.ListByWorkspace(workspaceID)

		// Update UI from goroutine using QueueUpdateDraw
		v.app.QueueUpdateDraw(func() {
			v.isLoading = false

			if err != nil {
				v.showError(fmt.Sprintf("Failed to load tickets: %v", err))
				return
			}

			// Handle empty state
			if len(tickets) == 0 {
				v.showEmptyState()
				return
			}

			// Build tree from tickets
			v.buildTree(tickets)
		})
	}()
}

// LoadTickets loads tickets for the specified workspace (synchronous fallback).
// Deprecated: Use LoadTicketsAsync instead to avoid blocking the TUI.
func (v *TicketTreeView) LoadTickets(workspaceID string) {
	// Just delegate to async version
	v.LoadTicketsAsync(workspaceID)
}

// buildTree constructs the tree structure from tickets.
// Optimized for performance with efficient string building and minimal allocations.
func (v *TicketTreeView) buildTree(tickets []domain.Ticket) {
	// Clear loading message
	v.root.ClearChildren()

	// Pre-allocate children slice to avoid repeated allocations
	children := make([]*tview.TreeNode, 0, len(tickets))

	for i := range tickets {
		ticket := &tickets[i] // Avoid copy

		// Build ticket text efficiently using byte slice
		ticketText := v.formatTicketText(ticket)

		ticketNode := tview.NewTreeNode(ticketText)
		ticketNode.SetColor(tcell.ColorGreen)
		ticketNode.SetReference(*ticket) // Dereference for SetReference

		// Add tasks as children if any
		if len(ticket.Tasks) > 0 {
			// Pre-allocate task nodes
			taskNodes := make([]*tview.TreeNode, 0, len(ticket.Tasks))

			for j := range ticket.Tasks {
				task := &ticket.Tasks[j]
				taskText := v.formatTaskText(task)

				taskNode := tview.NewTreeNode(taskText)
				taskNode.SetColor(tcell.ColorBlue)
				taskNode.SetReference(*task)

				taskNodes = append(taskNodes, taskNode)
			}

			// Add all task nodes at once
			for _, taskNode := range taskNodes {
				ticketNode.AddChild(taskNode)
			}
		}

		children = append(children, ticketNode)
	}

	// Add all children at once
	for _, child := range children {
		v.root.AddChild(child)
	}

	// Expand root by default
	v.root.SetExpanded(true)

	// Update border color based on selection mode
	v.updateSelectionBorder()
}

// formatTicketText builds ticket text efficiently with minimal allocations.
func (v *TicketTreeView) formatTicketText(ticket *domain.Ticket) string {
	const maxLen = 60

	// Calculate capacity (checkbox + ID + separator + title + buffer)
	capacity := 4 + len(ticket.JiraID) + 2
	if ticket.Title != "" {
		capacity += len(ticket.Title)
	} else {
		capacity += 10 // "(No title)"
	}

	// Use byte buffer for efficient string building
	buf := make([]byte, 0, capacity)

	// Add checkbox
	if v.selectedTickets[ticket.JiraID] {
		buf = append(buf, "[x] "...)
	} else {
		buf = append(buf, "[ ] "...)
	}

	// Add ticket ID
	buf = append(buf, ticket.JiraID...)

	// Add title
	if ticket.Title != "" {
		buf = append(buf, ": "...)
		buf = append(buf, ticket.Title...)
	} else {
		buf = append(buf, " (No title)"...)
	}

	// Truncate if needed
	if len(buf) > maxLen {
		buf = buf[:maxLen-3]
		buf = append(buf, "..."...)
	}

	return string(buf)
}

// formatTaskText builds task text efficiently.
func (v *TicketTreeView) formatTaskText(task *domain.Task) string {
	capacity := 6 + len(task.JiraID) + 2 + len(task.Title) // "  └─ " + ID + ": " + title
	if task.JiraID == "" {
		capacity = 10 + len(task.Title) // "  └─ TASK: " + title
	}

	buf := make([]byte, 0, capacity)
	buf = append(buf, "  └─ "...)

	if task.JiraID != "" {
		buf = append(buf, task.JiraID...)
	} else {
		buf = append(buf, "TASK"...)
	}

	if task.Title != "" {
		buf = append(buf, ": "...)
		buf = append(buf, task.Title...)
	}

	return string(buf)
}

// showMessage displays a temporary message in the tree.
func (v *TicketTreeView) showMessage(message string, color tcell.Color) {
	v.root.ClearChildren()
	node := tview.NewTreeNode(message)
	node.SetColor(color)
	node.SetSelectable(false)
	v.root.AddChild(node)
}

// showError displays an error message.
func (v *TicketTreeView) showError(message string) {
	v.root.ClearChildren()
	errorNode := tview.NewTreeNode("✗ " + message)
	errorNode.SetColor(tcell.ColorRed)
	errorNode.SetSelectable(false)
	v.root.AddChild(errorNode)

	hintNode := tview.NewTreeNode("Press Tab to switch to workspace panel")
	hintNode.SetColor(tcell.ColorGray)
	hintNode.SetSelectable(false)
	v.root.AddChild(hintNode)
}

// showEmptyState displays a message when no tickets are found.
func (v *TicketTreeView) showEmptyState() {
	v.root.ClearChildren()

	emptyNode := tview.NewTreeNode("No tickets in this workspace")
	emptyNode.SetColor(tcell.ColorYellow)
	emptyNode.SetSelectable(false)
	v.root.AddChild(emptyNode)

	hintNode := tview.NewTreeNode("Press 'P' to pull tickets from Jira")
	hintNode.SetColor(tcell.ColorBlue)
	hintNode.SetSelectable(false)
	v.root.AddChild(hintNode)
}

// toggleCurrentSelection toggles the selection state of the currently focused ticket.
func (v *TicketTreeView) toggleCurrentSelection() {
	node := v.tree.GetCurrentNode()
	if node == nil || node == v.root {
		return
	}

	ref := node.GetReference()
	if ref == nil {
		return
	}

	// Only toggle tickets, not tasks
	if ticket, ok := ref.(domain.Ticket); ok {
		v.selectedTickets[ticket.JiraID] = !v.selectedTickets[ticket.JiraID]
		v.updateSelectionMode()
		v.refreshTree()
	}
}

// selectAllVisible selects all visible tickets in the tree.
func (v *TicketTreeView) selectAllVisible() {
	children := v.root.GetChildren()
	for _, child := range children {
		ref := child.GetReference()
		if ref == nil {
			continue
		}

		// Only select tickets, not tasks
		if ticket, ok := ref.(domain.Ticket); ok {
			v.selectedTickets[ticket.JiraID] = true
		}
	}

	v.updateSelectionMode()
	v.refreshTree()
}

// clearSelection deselects all tickets.
func (v *TicketTreeView) clearSelection() {
	v.selectedTickets = make(map[string]bool)
	v.updateSelectionMode()
	v.refreshTree()
}

// GetSelectedTickets returns the list of selected ticket IDs.
func (v *TicketTreeView) GetSelectedTickets() []string {
	selected := make([]string, 0, len(v.selectedTickets))
	for ticketID, isSelected := range v.selectedTickets {
		if isSelected {
			selected = append(selected, ticketID)
		}
	}
	return selected
}

// ClearSelection clears all selected tickets (alias for clearSelection for external use).
func (v *TicketTreeView) ClearSelection() {
	v.clearSelection()
}

// updateSelectionMode updates the selection mode flag based on selected tickets.
func (v *TicketTreeView) updateSelectionMode() {
	v.selectionMode = false
	for _, isSelected := range v.selectedTickets {
		if isSelected {
			v.selectionMode = true
			break
		}
	}
}

// updateSelectionBorder updates the border color to indicate selection mode.
func (v *TicketTreeView) updateSelectionBorder() {
	if v.selectionMode {
		// Use info color (teal/blue) to indicate selection mode
		v.tree.SetBorderColor(theme.GetInfoColor())
		v.tree.SetTitle(fmt.Sprintf(" Tickets (%d selected) ", len(v.GetSelectedTickets())))
	} else {
		// Restore default border color based on focus state
		// This will be overridden by SetFocused if needed
		v.tree.SetBorderColor(theme.GetSecondaryColor())
		v.tree.SetTitle(" Tickets ")
	}
}

// refreshTree rebuilds the tree to reflect selection changes.
func (v *TicketTreeView) refreshTree() {
	// Get current workspace
	workspace, err := v.workspace.Current()
	if err != nil || workspace == nil {
		return
	}

	// Reload tickets
	tickets, err := v.ticketQuery.ListByWorkspace(workspace.ID)
	if err != nil || len(tickets) == 0 {
		return
	}

	// Rebuild tree with updated selection state
	v.buildTree(tickets)
}
