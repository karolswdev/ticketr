package services

import (
	"fmt"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// TicketRepository defines the minimal interface needed for ticket queries.
type TicketRepository interface {
	GetTicketsByWorkspace(workspaceID string) ([]domain.Ticket, error)
}

// TicketQueryService provides ticket querying operations for the TUI.
// It wraps the repository to maintain hexagonal architecture boundaries.
type TicketQueryService struct {
	repo TicketRepository
}

// NewTicketQueryService creates a new ticket query service.
func NewTicketQueryService(repo TicketRepository) *TicketQueryService {
	return &TicketQueryService{
		repo: repo,
	}
}

// ListByWorkspace retrieves all tickets for a given workspace.
// Returns an empty slice if the workspace has no tickets.
func (s *TicketQueryService) ListByWorkspace(workspaceID string) ([]domain.Ticket, error) {
	if workspaceID == "" {
		return nil, fmt.Errorf("workspace ID cannot be empty")
	}

	tickets, err := s.repo.GetTicketsByWorkspace(workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tickets for workspace %s: %w", workspaceID, err)
	}

	// Return empty slice instead of nil for consistency
	if tickets == nil {
		return []domain.Ticket{}, nil
	}

	return tickets, nil
}
