package ports

import "github.com/karolswdev/ticktr/internal/core/domain"

// Repository defines the interface for story persistence operations
type Repository interface {
	// GetStories reads and parses stories from a file
	GetStories(filepath string) ([]domain.Story, error)
	// SaveStories writes stories to a file in the custom Markdown format
	SaveStories(filepath string, stories []domain.Story) error
}