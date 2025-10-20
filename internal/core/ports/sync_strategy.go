package ports

import "github.com/karolswdev/ticktr/internal/core/domain"

// SyncStrategy defines how to resolve conflicts during sync operations.
// It provides methods for determining if a sync is needed and resolving
// conflicts when both local and remote tickets have changed.
//
// Implementations:
//   - LocalWinsStrategy: Always preserves local changes
//   - RemoteWinsStrategy: Always accepts remote changes
//   - ThreeWayMergeStrategy: Merges compatible changes, errors on conflicts
type SyncStrategy interface {
	// ShouldSync determines if a ticket should be synced based on the comparison
	// of local and remote states. It returns true if the remote ticket should
	// be considered for merging or updating the local ticket.
	//
	// Parameters:
	//   - localHash: The hash of the current local ticket state
	//   - remoteHash: The hash of the current remote ticket state
	//   - storedLocalHash: The hash of the last known local state
	//   - storedRemoteHash: The hash of the last known remote state
	//
	// Returns true if sync should proceed (i.e., remote has changed).
	ShouldSync(localHash, remoteHash, storedLocalHash, storedRemoteHash string) bool

	// ResolveConflict merges local and remote tickets according to the strategy's
	// conflict resolution policy. It is called when both local and remote have
	// changed since the last sync.
	//
	// Parameters:
	//   - local: The current local ticket
	//   - remote: The current remote ticket
	//
	// Returns:
	//   - The merged ticket (or one of the originals, depending on strategy)
	//   - An error if the conflict cannot be resolved automatically
	ResolveConflict(local, remote *domain.Ticket) (*domain.Ticket, error)

	// Name returns the human-readable name of the strategy.
	// Used for logging and configuration.
	Name() string
}
