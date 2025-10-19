package sync

import "time"

// SyncState represents the current state of a sync operation.
type SyncState int

const (
	// StateIdle indicates no sync operation is in progress.
	StateIdle SyncState = iota
	// StateSyncing indicates a sync operation is currently in progress.
	StateSyncing
	// StateSuccess indicates the last sync operation completed successfully.
	StateSuccess
	// StateError indicates the last sync operation failed.
	StateError
)

// String returns a human-readable representation of the sync state.
func (s SyncState) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StateSyncing:
		return "syncing"
	case StateSuccess:
		return "success"
	case StateError:
		return "error"
	default:
		return "unknown"
	}
}

// SyncStatus holds the current sync operation status and metadata.
type SyncStatus struct {
	State     SyncState // Current state of sync operation
	Operation string    // Type of operation: "push", "pull", "sync", "refresh"
	Progress  string    // Human-readable progress message (e.g., "3 created, 5 updated")
	Error     string    // Error message if State is StateError
	Timestamp time.Time // When the status was last updated
}

// NewIdleStatus creates a new SyncStatus in the idle state.
func NewIdleStatus() SyncStatus {
	return SyncStatus{
		State:     StateIdle,
		Operation: "",
		Progress:  "Ready",
		Error:     "",
		Timestamp: time.Now(),
	}
}

// NewSyncingStatus creates a new SyncStatus for an operation in progress.
func NewSyncingStatus(operation string, message string) SyncStatus {
	return SyncStatus{
		State:     StateSyncing,
		Operation: operation,
		Progress:  message,
		Error:     "",
		Timestamp: time.Now(),
	}
}

// NewSuccessStatus creates a new SyncStatus for a successful operation.
func NewSuccessStatus(operation string, message string) SyncStatus {
	return SyncStatus{
		State:     StateSuccess,
		Operation: operation,
		Progress:  message,
		Error:     "",
		Timestamp: time.Now(),
	}
}

// NewErrorStatus creates a new SyncStatus for a failed operation.
func NewErrorStatus(operation string, err error) SyncStatus {
	return SyncStatus{
		State:     StateError,
		Operation: operation,
		Progress:  "",
		Error:     err.Error(),
		Timestamp: time.Now(),
	}
}
