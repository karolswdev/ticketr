package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TicketTreeView displays a hierarchical tree of tickets (stub implementation).
type TicketTreeView struct {
	tree *tview.TreeView
	root *tview.TreeNode
}

// NewTicketTreeView creates a new ticket tree view.
func NewTicketTreeView() *TicketTreeView {
	root := tview.NewTreeNode("Tickets")
	root.SetColor(tcell.ColorYellow)

	tree := tview.NewTreeView()
	tree.SetRoot(root)
	tree.SetCurrentNode(root)
	tree.SetBorder(true).SetTitle(" Tickets ")

	view := &TicketTreeView{
		tree: tree,
		root: root,
	}

	// Add placeholder nodes
	view.addPlaceholderNodes()

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
	// Future: Load tickets from service
}

// OnHide is called when the view is hidden.
func (v *TicketTreeView) OnHide() {
	// No cleanup needed
}

// addPlaceholderNodes adds example nodes for initial skeleton.
func (v *TicketTreeView) addPlaceholderNodes() {
	// Clear existing children
	v.root.ClearChildren()

	// Add placeholder structure
	epic1 := tview.NewTreeNode("📋 EPIC-001: User Authentication")
	epic1.SetColor(tcell.ColorGreen)

	story1 := tview.NewTreeNode("📝 STORY-101: Login page")
	story1.SetColor(tcell.ColorBlue)

	story2 := tview.NewTreeNode("📝 STORY-102: OAuth integration")
	story2.SetColor(tcell.ColorBlue)

	epic1.AddChild(story1)
	epic1.AddChild(story2)

	epic2 := tview.NewTreeNode("📋 EPIC-002: Dashboard redesign")
	epic2.SetColor(tcell.ColorGreen)

	v.root.AddChild(epic1)
	v.root.AddChild(epic2)

	// Add informational node
	info := tview.NewTreeNode("(Placeholder data - integration pending)")
	info.SetColor(tcell.ColorGray)
	v.root.AddChild(info)
}
