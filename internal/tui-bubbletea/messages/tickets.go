package messages

import "github.com/karolswdev/ticktr/internal/core/domain"

// Ticket-related messages for loading, selecting, updating tickets.

// TicketSelectedMsg is sent when a ticket is selected in the tree.
type TicketSelectedMsg struct {
	TicketID string
}

// TicketOpenedMsg is sent when a ticket is opened in detail view.
type TicketOpenedMsg struct {
	TicketID string
}

// TicketsLoadedMsg is sent when tickets are loaded from the database.
type TicketsLoadedMsg struct {
	Tickets []domain.Ticket
	Error   error
}

// TicketUpdatedMsg is sent when a ticket is updated.
type TicketUpdatedMsg struct {
	TicketID string
}

// TicketDeletedMsg is sent when a ticket is deleted.
type TicketDeletedMsg struct {
	TicketID string
}

// TicketCreatedMsg is sent when a new ticket is created.
type TicketCreatedMsg struct {
	TicketID string
}

// TreeExpandedMsg is sent when a tree node is expanded.
type TreeExpandedMsg struct {
	TicketID string
}

// TreeCollapsedMsg is sent when a tree node is collapsed.
type TreeCollapsedMsg struct {
	TicketID string
}

// FilterChangedMsg is sent when the tree filter changes.
type FilterChangedMsg struct {
	Query string
}
