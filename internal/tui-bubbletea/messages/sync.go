package messages

import "time"

// Sync-related messages for Jira synchronization operations.

// SyncStatus represents the current sync state.
type SyncStatus int

const (
	SyncIdle SyncStatus = iota
	SyncInProgress
	SyncSuccess
	SyncError
)

// SyncStartedMsg is sent when a sync operation starts.
type SyncStartedMsg struct {
	Operation string // "pull", "push", "full"
}

// SyncProgressMsg is sent during sync to report progress.
type SyncProgressMsg struct {
	Operation   string
	Percent     float64
	Current     int
	Total       int
	Message     string
	ElapsedTime time.Duration
}

// SyncCompletedMsg is sent when a sync operation completes successfully.
type SyncCompletedMsg struct {
	Operation    string
	TicketsAdded int
	TicketsUpdated int
	TicketsDeleted int
	Duration     time.Duration
}

// SyncErrorMsg is sent when a sync operation fails.
type SyncErrorMsg struct {
	Operation string
	Err       error
}

// SyncCancelledMsg is sent when a sync operation is cancelled by the user.
type SyncCancelledMsg struct {
	Operation string
}
