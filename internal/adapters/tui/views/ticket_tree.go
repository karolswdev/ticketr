package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
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
}

// NewTicketTreeView creates a new ticket tree view.
func NewTicketTreeView(workspace *services.WorkspaceService, ticketQuery *services.TicketQueryService) *TicketTreeView {
	root := tview.NewTreeNode("Tickets")
	root.SetColor(tcell.ColorYellow)

	tree := tview.NewTreeView()
	tree.SetRoot(root)
	tree.SetCurrentNode(root)
	tree.SetBorder(true).SetTitle(" Tickets ")

	view := &TicketTreeView{
		tree:        tree,
		root:        root,
		workspace:   workspace,
		ticketQuery: ticketQuery,
	}

	// Setup vim-style navigation
	view.setupKeyBindings()

	// Setup selection handler for Enter key
	view.setupSelectionHandler()

	// Load initial tickets
	view.loadInitialTickets()

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
	color := tcell.ColorWhite
	if focused {
		color = tcell.ColorGreen
	}
	v.tree.SetBorderColor(color)
}

// loadInitialTickets loads tickets for the current workspace.
func (v *TicketTreeView) loadInitialTickets() {
	// Get current workspace
	workspace, err := v.workspace.Current()
	if err != nil || workspace == nil {
		v.showMessage("No active workspace", tcell.ColorYellow)
		return
	}

	v.LoadTickets(workspace.ID)
}

// LoadTickets loads tickets for the specified workspace.
func (v *TicketTreeView) LoadTickets(workspaceID string) {
	// Clear existing tree
	v.root.ClearChildren()

	// Show loading state
	v.showMessage("Loading tickets...", tcell.ColorBlue)

	// Fetch tickets
	tickets, err := v.ticketQuery.ListByWorkspace(workspaceID)
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
}

// buildTree constructs the tree structure from tickets.
func (v *TicketTreeView) buildTree(tickets []domain.Ticket) {
	// Clear loading message
	v.root.ClearChildren()

	for _, ticket := range tickets {
		// Create ticket node with JiraID and Title
		ticketText := ticket.JiraID
		if ticket.Title != "" {
			ticketText += ": " + ticket.Title
		} else {
			ticketText += " (No title)"
		}

		// Truncate long titles
		maxLen := 60
		if len(ticketText) > maxLen {
			ticketText = ticketText[:maxLen-3] + "..."
		}

		ticketNode := tview.NewTreeNode(ticketText)
		ticketNode.SetColor(tcell.ColorGreen)
		ticketNode.SetReference(ticket)

		// Add tasks as children if any
		for _, task := range ticket.Tasks {
			taskText := task.JiraID
			if taskText == "" {
				taskText = "TASK"
			}
			if task.Title != "" {
				taskText += ": " + task.Title
			}

			taskNode := tview.NewTreeNode("  └─ " + taskText)
			taskNode.SetColor(tcell.ColorBlue)
			taskNode.SetReference(task)

			ticketNode.AddChild(taskNode)
		}

		v.root.AddChild(ticketNode)
	}

	// Expand root by default
	v.root.SetExpanded(true)
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

	emptyNode := tview.NewTreeNode("No tickets found")
	emptyNode.SetColor(tcell.ColorGray)
	emptyNode.SetSelectable(false)
	v.root.AddChild(emptyNode)

	hintNode := tview.NewTreeNode("Run 'ticketr pull' to sync tickets")
	hintNode.SetColor(tcell.ColorYellow)
	hintNode.SetSelectable(false)
	v.root.AddChild(hintNode)
}
