package ports

import "github.com/karolswdev/ticktr/internal/core/domain"

// Repository defines the interface for ticket persistence operations
type Repository interface {
	// GetTickets reads and parses tickets from a file
	GetTickets(filepath string) ([]domain.Ticket, error)
	// SaveTickets writes tickets to a file in the custom Markdown format
	SaveTickets(filepath string, tickets []domain.Ticket) error
}