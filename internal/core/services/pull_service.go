package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
	"github.com/karolswdev/ticktr/internal/state"
)

var (
	// ErrConflictDetected is returned when a merge conflict is detected
	ErrConflictDetected = errors.New("conflict detected")
)

// ProgressCallback is called to report progress during pull operations
// current: number of items processed so far
// total: total number of items to process (0 if unknown)
// message: human-readable status message
type ProgressCallback func(current, total int, message string)

// PullService handles pulling tickets from JIRA and updating local files
type PullService struct {
	jiraAdapter    ports.JiraPort
	repository     ports.Repository
	dbRepository   ports.ExtendedRepository // Database repository for TUI
	stateManager   *state.StateManager
	syncStrategy   ports.SyncStrategy
}

// NewPullService creates a new pull service instance.
// If syncStrategy is nil, defaults to RemoteWinsStrategy for backward compatibility.
func NewPullService(jiraAdapter ports.JiraPort, repository ports.Repository, stateManager *state.StateManager) *PullService {
	return &PullService{
		jiraAdapter:    jiraAdapter,
		repository:     repository,
		dbRepository:   nil, // No database support in legacy mode
		stateManager:   stateManager,
		syncStrategy:   &RemoteWinsStrategy{}, // Default strategy for backward compatibility
	}
}

// NewPullServiceWithDB creates a new pull service instance with database support.
// This is the constructor used by the TUI to enable database persistence.
func NewPullServiceWithDB(jiraAdapter ports.JiraPort, dbRepository ports.ExtendedRepository, fileRepository ports.Repository, stateManager *state.StateManager) *PullService {
	return &PullService{
		jiraAdapter:    jiraAdapter,
		repository:     fileRepository,
		dbRepository:   dbRepository,
		stateManager:   stateManager,
		syncStrategy:   &RemoteWinsStrategy{}, // Default strategy
	}
}

// NewPullServiceWithStrategy creates a new pull service instance with a custom sync strategy.
func NewPullServiceWithStrategy(jiraAdapter ports.JiraPort, repository ports.Repository, stateManager *state.StateManager, syncStrategy ports.SyncStrategy) *PullService {
	if syncStrategy == nil {
		syncStrategy = &RemoteWinsStrategy{} // Fallback to default
	}
	return &PullService{
		jiraAdapter:    jiraAdapter,
		repository:     repository,
		dbRepository:   nil, // No database support in legacy mode
		stateManager:   stateManager,
		syncStrategy:   syncStrategy,
	}
}

// PullOptions contains options for the pull operation
type PullOptions struct {
	ProjectKey       string
	JQL              string
	EpicKey          string
	Force            bool             // Force overwrite even if conflicts exist
	ProgressCallback ProgressCallback // Optional callback for progress updates
	WorkspaceID      string           // Workspace ID for database operations (TUI mode)
}

// PullResult contains the results of a pull operation
type PullResult struct {
	TicketsPulled  int
	TicketsUpdated int
	TicketsSkipped int
	Conflicts      []string
	Errors         []error
}

// Pull fetches tickets from JIRA and updates the local file with context support for cancellation
func (ps *PullService) Pull(ctx context.Context, filePath string, options PullOptions) (*PullResult, error) {
	result := &PullResult{}
	progress := options.ProgressCallback

	// Helper to safely call progress callback
	reportProgress := func(current, total int, message string) {
		if progress != nil {
			progress(current, total, message)
		}
	}

	// Load current state
	if err := ps.stateManager.Load(); err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	// Build JQL query
	jql := ps.buildJQL(options)

	// Create Jira progress callback that forwards to our progress callback
	jiraProgressCallback := func(current, total int, message string) {
		reportProgress(current, total, message)
	}

	// Fetch tickets from JIRA with context for cancellation and progress reporting
	remoteTickets, err := ps.jiraAdapter.SearchTickets(ctx, options.ProjectKey, jql, jiraProgressCallback)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tickets from JIRA: %w", err)
	}

	// Load local tickets
	localTickets, err := ps.repository.GetTickets(filePath)
	if err != nil && !errors.Is(err, ports.ErrFileNotFound) {
		return nil, fmt.Errorf("failed to load local tickets: %w", err)
	}

	// Create a map of local tickets by JiraID for easier lookup
	localTicketMap := make(map[string]*domain.Ticket)
	for i := range localTickets {
		if localTickets[i].JiraID != "" {
			localTicketMap[localTickets[i].JiraID] = &localTickets[i]
		}
	}

	// Process each remote ticket
	mergedTickets := []domain.Ticket{}
	totalTickets := len(remoteTickets)
	for i, remoteTicket := range remoteTickets {
		// Report progress for larger datasets (10+ tickets)
		if totalTickets >= 10 {
			reportProgress(i+1, totalTickets, "")
		}

		remoteHash := ps.stateManager.CalculateHash(remoteTicket)

		// Check if ticket exists locally
		localTicket, existsLocally := localTicketMap[remoteTicket.JiraID]

		if !existsLocally {
			// New ticket from remote
			mergedTickets = append(mergedTickets, remoteTicket)
			ps.stateManager.UpdateHash(remoteTicket)
			result.TicketsPulled++
		} else {
			// Ticket exists both locally and remotely - check for conflicts
			localHash := ps.stateManager.CalculateHash(*localTicket)
			storedState, hasStoredState := ps.stateManager.GetStoredState(remoteTicket.JiraID)

			if !hasStoredState {
				// No stored state - first time seeing this ticket
				// Take remote version and update state
				mergedTickets = append(mergedTickets, remoteTicket)
				ps.stateManager.UpdateHash(remoteTicket)
				result.TicketsUpdated++
			} else {
				// We have stored state - check for conflicts
				localChanged := localHash != storedState.LocalHash
				remoteChanged := remoteHash != storedState.RemoteHash

				if localChanged && remoteChanged {
					// Conflict detected - use sync strategy to resolve
					result.Conflicts = append(result.Conflicts, remoteTicket.JiraID)

					if options.Force {
						// Force mode - take remote version
						mergedTickets = append(mergedTickets, remoteTicket)
						ps.stateManager.UpdateHash(remoteTicket)
						result.TicketsUpdated++
					} else {
						// Use sync strategy to resolve conflict
						resolvedTicket, err := ps.syncStrategy.ResolveConflict(localTicket, &remoteTicket)
						if err != nil {
							// Strategy failed to resolve - return error
							return result, fmt.Errorf("conflict resolution failed for %s using %s strategy: %w",
								remoteTicket.JiraID, ps.syncStrategy.Name(), err)
						}
						mergedTickets = append(mergedTickets, *resolvedTicket)
						ps.stateManager.UpdateHash(*resolvedTicket)
						result.TicketsUpdated++
					}
				} else if remoteChanged && !localChanged {
					// Only remote changed - safe to update
					mergedTickets = append(mergedTickets, remoteTicket)
					ps.stateManager.SetStoredState(remoteTicket.JiraID, state.TicketState{
						LocalHash:  remoteHash,
						RemoteHash: remoteHash,
					})
					result.TicketsUpdated++
				} else if localChanged && !remoteChanged {
					// Only local changed - keep local version
					mergedTickets = append(mergedTickets, *localTicket)
					ps.stateManager.UpdateLocalHash(*localTicket)
					result.TicketsSkipped++
				} else {
					// No changes - keep as is
					mergedTickets = append(mergedTickets, *localTicket)
					result.TicketsSkipped++
				}
			}

			// Remove from map to track what's left (local-only tickets)
			delete(localTicketMap, remoteTicket.JiraID)
		}
	}

	// Add any remaining local-only tickets
	for _, localTicket := range localTicketMap {
		mergedTickets = append(mergedTickets, *localTicket)
	}

	// Save merged tickets to file (for backward compatibility with CLI)
	if ps.repository != nil && filePath != "" {
		if err := ps.repository.SaveTickets(filePath, mergedTickets); err != nil {
			// Log error but don't fail - database is primary in TUI mode
			reportProgress(0, 0, fmt.Sprintf("Warning: failed to save to file: %v", err))
		}
	}

	// Save tickets to database if database repository is configured (TUI mode)
	if ps.dbRepository != nil && options.WorkspaceID != "" {
		reportProgress(0, len(mergedTickets), "Saving to database...")

		// Save all tickets to database using workspace-specific path
		// The SQLiteAdapter will extract workspace ID and sync to database
		workspacePath := fmt.Sprintf("workspace-%s.md", options.WorkspaceID)
		if err := ps.dbRepository.SaveTickets(workspacePath, mergedTickets); err != nil {
			return nil, fmt.Errorf("failed to save tickets to database: %w", err)
		}

		reportProgress(len(mergedTickets), len(mergedTickets), fmt.Sprintf("Saved %d tickets to database", len(mergedTickets)))
	}

	// Save updated state
	if err := ps.stateManager.Save(); err != nil {
		return nil, fmt.Errorf("failed to save state: %w", err)
	}

	return result, nil
}

// buildJQL constructs the JQL query from options
func (ps *PullService) buildJQL(options PullOptions) string {
	jql := ""

	if options.JQL != "" {
		jql = options.JQL
	}

	if options.EpicKey != "" {
		if jql != "" {
			jql += " AND "
		}
		jql += fmt.Sprintf(`"Epic Link" = %s`, options.EpicKey)
	}

	return jql
}

