package sync

import (
	"fmt"

	"github.com/karolswdev/ticktr/internal/core/services"
)

// StatusCallback is called when sync status changes.
type StatusCallback func(status SyncStatus)

// SyncCoordinator manages async sync operations and provides status updates.
type SyncCoordinator struct {
	pushService    *services.PushService
	pullService    *services.PullService
	onStatusChange StatusCallback
}

// NewSyncCoordinator creates a new sync coordinator with the given services.
func NewSyncCoordinator(
	pushService *services.PushService,
	pullService *services.PullService,
	onStatusChange StatusCallback,
) *SyncCoordinator {
	return &SyncCoordinator{
		pushService:    pushService,
		pullService:    pullService,
		onStatusChange: onStatusChange,
	}
}

// PushAsync performs a push operation asynchronously.
// The operation runs in a goroutine and status updates are sent via the callback.
func (sc *SyncCoordinator) PushAsync(filePath string, options services.ProcessOptions) {
	// Immediately signal that sync is starting
	sc.notifyStatus(NewSyncingStatus("push", "Pushing tickets to Jira..."))

	// Run push in goroutine
	go func() {
		result, err := sc.pushService.PushTickets(filePath, options)
		if err != nil {
			// Operation failed
			sc.notifyStatus(NewErrorStatus("push", err))
			return
		}

		// Operation succeeded - format success message
		msg := sc.formatPushResult(result)
		sc.notifyStatus(NewSuccessStatus("push", msg))
	}()
}

// PullAsync performs a pull operation asynchronously.
// The operation runs in a goroutine and status updates are sent via the callback.
func (sc *SyncCoordinator) PullAsync(filePath string, options services.PullOptions) {
	// Immediately signal that sync is starting
	sc.notifyStatus(NewSyncingStatus("pull", "Pulling tickets from Jira..."))

	// Run pull in goroutine
	go func() {
		result, err := sc.pullService.Pull(filePath, options)
		if err != nil {
			// Operation failed
			sc.notifyStatus(NewErrorStatus("pull", err))
			return
		}

		// Operation succeeded - format success message
		msg := sc.formatPullResult(result)
		sc.notifyStatus(NewSuccessStatus("pull", msg))
	}()
}

// SyncAsync performs a full sync (pull then push) asynchronously.
// The operation runs in a goroutine and status updates are sent via the callback.
func (sc *SyncCoordinator) SyncAsync(filePath string) {
	// Immediately signal that sync is starting
	sc.notifyStatus(NewSyncingStatus("sync", "Syncing with Jira (pull then push)..."))

	// Run sync in goroutine
	go func() {
		// Step 1: Pull from Jira
		pullResult, err := sc.pullService.Pull(filePath, services.PullOptions{})
		if err != nil {
			errMsg := fmt.Errorf("pull failed during sync: %w", err)
			sc.notifyStatus(NewErrorStatus("sync", errMsg))
			return
		}

		// Update status to show pull completed
		sc.notifyStatus(NewSyncingStatus("sync", fmt.Sprintf("Pull complete (%s), now pushing...", sc.formatPullResult(pullResult))))

		// Step 2: Push to Jira
		pushResult, err := sc.pushService.PushTickets(filePath, services.ProcessOptions{})
		if err != nil {
			errMsg := fmt.Errorf("push failed during sync: %w", err)
			sc.notifyStatus(NewErrorStatus("sync", errMsg))
			return
		}

		// Both operations succeeded
		msg := fmt.Sprintf("Pull: %s | Push: %s", sc.formatPullResult(pullResult), sc.formatPushResult(pushResult))
		sc.notifyStatus(NewSuccessStatus("sync", msg))
	}()
}

// notifyStatus sends a status update via the callback if configured.
func (sc *SyncCoordinator) notifyStatus(status SyncStatus) {
	if sc.onStatusChange != nil {
		sc.onStatusChange(status)
	}
}

// formatPushResult creates a human-readable summary of push results.
func (sc *SyncCoordinator) formatPushResult(result *services.ProcessResult) string {
	if result == nil {
		return "No changes"
	}

	parts := []string{}
	if result.TicketsCreated > 0 {
		parts = append(parts, fmt.Sprintf("%d ticket(s) created", result.TicketsCreated))
	}
	if result.TicketsUpdated > 0 {
		parts = append(parts, fmt.Sprintf("%d ticket(s) updated", result.TicketsUpdated))
	}
	if result.TasksCreated > 0 {
		parts = append(parts, fmt.Sprintf("%d task(s) created", result.TasksCreated))
	}
	if result.TasksUpdated > 0 {
		parts = append(parts, fmt.Sprintf("%d task(s) updated", result.TasksUpdated))
	}

	if len(parts) == 0 {
		return "No changes"
	}

	// Join parts with commas
	msg := ""
	for i, part := range parts {
		if i > 0 {
			msg += ", "
		}
		msg += part
	}

	// Add error count if present
	if len(result.Errors) > 0 {
		msg += fmt.Sprintf(" (%d error(s))", len(result.Errors))
	}

	return msg
}

// formatPullResult creates a human-readable summary of pull results.
func (sc *SyncCoordinator) formatPullResult(result *services.PullResult) string {
	if result == nil {
		return "No changes"
	}

	parts := []string{}
	if result.TicketsPulled > 0 {
		parts = append(parts, fmt.Sprintf("%d new ticket(s)", result.TicketsPulled))
	}
	if result.TicketsUpdated > 0 {
		parts = append(parts, fmt.Sprintf("%d ticket(s) updated", result.TicketsUpdated))
	}
	if result.TicketsSkipped > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped", result.TicketsSkipped))
	}

	if len(parts) == 0 {
		return "No changes"
	}

	// Join parts with commas
	msg := ""
	for i, part := range parts {
		if i > 0 {
			msg += ", "
		}
		msg += part
	}

	// Add conflict/error count if present
	if len(result.Conflicts) > 0 {
		msg += fmt.Sprintf(" (%d conflict(s))", len(result.Conflicts))
	}
	if len(result.Errors) > 0 {
		msg += fmt.Sprintf(" (%d error(s))", len(result.Errors))
	}

	return msg
}
