package ports

import "github.com/karolswdev/ticktr/internal/core/domain"

// Repository defines the interface for story persistence operations
type Repository interface {
	// GetStories reads and parses stories from a file
	GetStories(filepath string) ([]domain.Story, error)
	// SaveStories writes stories to a file in the custom Markdown format
	SaveStories(filepath string, stories []domain.Story) error
	// GetTickets reads and parses tickets from a file
	GetTickets(filepath string) ([]domain.Ticket, error)
	// SaveTickets writes tickets to a file in the custom Markdown format
	SaveTickets(filepath string, tickets []domain.Ticket) error
}