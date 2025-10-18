package ports

import (
	"time"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// SyncOperation represents a synchronization operation for logging
type SyncOperation struct {
	WorkspaceID   string
	Operation     string // push, pull, migrate, conflict_resolve
	FilePath      string
	TicketCount   int
	SuccessCount  int
	FailureCount  int
	ConflictCount int
	DurationMs    int
	ErrorDetails  map[string]interface{}
	StartedAt     time.Time
	CompletedAt   time.Time
}

// ExtendedRepository extends the base Repository interface with v3-specific methods.
// This interface provides enhanced query capabilities for workspace-aware operations,
// conflict detection, and synchronization state management.
//
// Implementations must ensure thread-safety for concurrent operations.
type ExtendedRepository interface {
	// Repository embeds the base v2 interface for backward compatibility
	Repository

	// GetTicketsByWorkspace retrieves all tickets associated with a specific workspace.
	// Returns an empty slice if the workspace has no tickets.
	// Returns an error if the workspace doesn't exist or on database errors.
	GetTicketsByWorkspace(workspaceID string) ([]domain.Ticket, error)

	// GetModifiedTickets retrieves tickets that have been modified since the given timestamp.
	// Useful for incremental synchronization operations.
	// Returns an empty slice if no tickets were modified since the specified time.
	GetModifiedTickets(since time.Time) ([]domain.Ticket, error)

	// UpdateTicketState updates the synchronization state of a ticket.
	// This includes local_hash, remote_hash, and sync_status fields.
	// Returns an error if the ticket doesn't exist or on database errors.
	UpdateTicketState(ticket domain.Ticket) error

	// DetectConflicts identifies tickets with conflicting local and remote states.
	// Returns tickets where sync_status = 'conflict' or where local_hash != remote_hash.
	// Returns an empty slice if no conflicts are detected.
	DetectConflicts() ([]domain.Ticket, error)

	// LogSyncOperation records a synchronization operation for audit and debugging.
	// Operations include push, pull, migrate, and conflict_resolve.
	// Returns an error on database errors.
	LogSyncOperation(op SyncOperation) error
}
