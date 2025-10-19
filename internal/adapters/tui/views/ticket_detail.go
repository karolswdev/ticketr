package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TicketDetailView displays detailed information about a ticket (stub implementation).
type TicketDetailView struct {
	textView *tview.TextView
}

// NewTicketDetailView creates a new ticket detail view.
func NewTicketDetailView() *TicketDetailView {
	textView := tview.NewTextView()
	textView.SetBorder(true).SetTitle(" Ticket Detail ")
	textView.SetDynamicColors(true)
	textView.SetWordWrap(true)
	textView.SetScrollable(true)

	view := &TicketDetailView{
		textView: textView,
	}

	// Show placeholder content
	view.showPlaceholder()

	// Enable scrolling with arrow keys
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			row, col := textView.GetScrollOffset()
			textView.ScrollTo(row-1, col)
			return nil
		case tcell.KeyDown:
			row, col := textView.GetScrollOffset()
			textView.ScrollTo(row+1, col)
			return nil
		}
		return event
	})

	return view
}

// Name returns the view identifier.
func (v *TicketDetailView) Name() string {
	return "ticket_detail"
}

// Primitive returns the tview primitive.
func (v *TicketDetailView) Primitive() tview.Primitive {
	return v.textView
}

// OnShow is called when the view becomes active.
func (v *TicketDetailView) OnShow() {
	// Future: Load selected ticket details
}

// OnHide is called when the view is hidden.
func (v *TicketDetailView) OnHide() {
	// No cleanup needed
}

// showPlaceholder displays example ticket content.
func (v *TicketDetailView) showPlaceholder() {
	content := `[yellow]STORY-101: Login page[-]

[::b]Status:[::-] In Progress
[::b]Assignee:[::-] john.doe@example.com
[::b]Priority:[::-] High
[::b]Sprint:[::-] Sprint 23

[::u]Description[::-]
Implement a responsive login page with email and password fields.
Support for "Remember Me" checkbox and "Forgot Password" link.

[::u]Acceptance Criteria[::-]
• Form validation for email format
• Password minimum 8 characters
• Error messages display clearly
• Mobile responsive design

[::u]Comments[::-]
[gray]2025-01-15 14:30 - Jane Smith:[-]
Please ensure accessibility standards are met.

[gray]2025-01-16 09:15 - John Doe:[-]
Working on the validation logic now.

[gray](Placeholder data - integration pending)[-]
`
	v.textView.SetText(content)
}
