package services

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

var (
	// ErrConflictUnresolvable is returned when a conflict cannot be automatically resolved
	ErrConflictUnresolvable = errors.New("conflict cannot be automatically resolved")

	// ErrUnknownStrategy is returned when an invalid strategy name is provided
	ErrUnknownStrategy = errors.New("unknown sync strategy")
)

// Strategy name constants
const (
	StrategyLocalWins      = "local-wins"
	StrategyRemoteWins     = "remote-wins"
	StrategyThreeWayMerge  = "three-way-merge"
)

// NewSyncStrategy creates a sync strategy by name.
// Supported strategies: "local-wins", "remote-wins", "three-way-merge"
func NewSyncStrategy(name string) (ports.SyncStrategy, error) {
	switch name {
	case StrategyLocalWins:
		return &LocalWinsStrategy{}, nil
	case StrategyRemoteWins:
		return &RemoteWinsStrategy{}, nil
	case StrategyThreeWayMerge:
		return &ThreeWayMergeStrategy{}, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownStrategy, name)
	}
}

// LocalWinsStrategy always preserves local changes in case of conflict.
// Use this when you want to ensure local edits are never overwritten.
type LocalWinsStrategy struct{}

// Name returns the strategy name.
func (s *LocalWinsStrategy) Name() string {
	return StrategyLocalWins
}

// ShouldSync returns true if the remote has changed.
func (s *LocalWinsStrategy) ShouldSync(localHash, remoteHash, storedLocalHash, storedRemoteHash string) bool {
	return remoteHash != storedRemoteHash
}

// ResolveConflict always returns the local ticket unchanged.
func (s *LocalWinsStrategy) ResolveConflict(local, remote *domain.Ticket) (*domain.Ticket, error) {
	if local == nil {
		return nil, errors.New("local ticket is nil")
	}
	if remote == nil {
		return nil, errors.New("remote ticket is nil")
	}

	// Always return local ticket (make a copy to avoid mutations)
	result := *local
	return &result, nil
}

// RemoteWinsStrategy always accepts remote changes in case of conflict.
// Use this when Jira is the single source of truth.
type RemoteWinsStrategy struct{}

// Name returns the strategy name.
func (s *RemoteWinsStrategy) Name() string {
	return StrategyRemoteWins
}

// ShouldSync returns true if the remote has changed.
func (s *RemoteWinsStrategy) ShouldSync(localHash, remoteHash, storedLocalHash, storedRemoteHash string) bool {
	return remoteHash != storedRemoteHash
}

// ResolveConflict always returns the remote ticket unchanged.
func (s *RemoteWinsStrategy) ResolveConflict(local, remote *domain.Ticket) (*domain.Ticket, error) {
	if local == nil {
		return nil, errors.New("local ticket is nil")
	}
	if remote == nil {
		return nil, errors.New("remote ticket is nil")
	}

	// Always return remote ticket (make a copy to avoid mutations)
	result := *remote
	return &result, nil
}

// ThreeWayMergeStrategy attempts to merge compatible changes from both local and remote.
// It compares field-by-field and merges when different fields have been modified.
// Returns an error when the same field has been modified in both versions.
type ThreeWayMergeStrategy struct{}

// Name returns the strategy name.
func (s *ThreeWayMergeStrategy) Name() string {
	return StrategyThreeWayMerge
}

// ShouldSync returns true if the remote has changed.
func (s *ThreeWayMergeStrategy) ShouldSync(localHash, remoteHash, storedLocalHash, storedRemoteHash string) bool {
	return remoteHash != storedRemoteHash
}

// ResolveConflict performs a three-way merge of local and remote tickets.
// It returns a merged ticket with compatible changes from both sides, or an error
// if incompatible changes are detected (same field modified in both).
func (s *ThreeWayMergeStrategy) ResolveConflict(local, remote *domain.Ticket) (*domain.Ticket, error) {
	if local == nil {
		return nil, errors.New("local ticket is nil")
	}
	if remote == nil {
		return nil, errors.New("remote ticket is nil")
	}

	// Start with a copy of the local ticket
	merged := *local

	// Track conflicts
	conflicts := []string{}

	// Compare and merge Title
	if local.Title != remote.Title {
		// Field differs - check if it's a conflict or a safe merge
		// Since we don't have a base version, we consider any difference a potential conflict
		// However, we can use a simple heuristic: if one is empty, take the non-empty one
		if local.Title == "" {
			merged.Title = remote.Title
		} else if remote.Title == "" {
			merged.Title = local.Title
		} else {
			// Both have values and they differ - this is a conflict
			conflicts = append(conflicts, "Title")
		}
	}

	// Compare and merge Description
	if local.Description != remote.Description {
		if local.Description == "" {
			merged.Description = remote.Description
		} else if remote.Description == "" {
			merged.Description = local.Description
		} else {
			conflicts = append(conflicts, "Description")
		}
	}

	// Compare and merge AcceptanceCriteria
	if !stringSlicesEqual(local.AcceptanceCriteria, remote.AcceptanceCriteria) {
		if len(local.AcceptanceCriteria) == 0 {
			merged.AcceptanceCriteria = remote.AcceptanceCriteria
		} else if len(remote.AcceptanceCriteria) == 0 {
			merged.AcceptanceCriteria = local.AcceptanceCriteria
		} else {
			conflicts = append(conflicts, "AcceptanceCriteria")
		}
	}

	// Compare and merge CustomFields
	if !mapsEqual(local.CustomFields, remote.CustomFields) {
		mergedFields, fieldConflicts := mergeCustomFields(local.CustomFields, remote.CustomFields)
		merged.CustomFields = mergedFields
		if len(fieldConflicts) > 0 {
			for _, field := range fieldConflicts {
				conflicts = append(conflicts, fmt.Sprintf("CustomFields[%s]", field))
			}
		}
	}

	// Compare and merge Tasks
	if !tasksEqual(local.Tasks, remote.Tasks) {
		mergedTasks, taskConflicts := mergeTasks(local.Tasks, remote.Tasks)
		merged.Tasks = mergedTasks
		if len(taskConflicts) > 0 {
			conflicts = append(conflicts, taskConflicts...)
		}
	}

	// If there are any conflicts, return an error
	if len(conflicts) > 0 {
		return nil, fmt.Errorf("%w: fields %v have conflicting changes", ErrConflictUnresolvable, conflicts)
	}

	// Preserve JiraID from remote (it should be the same, but use remote as source of truth)
	merged.JiraID = remote.JiraID

	return &merged, nil
}

// mergeCustomFields merges two custom field maps, detecting conflicts.
// Returns the merged map and a list of conflicting field keys.
func mergeCustomFields(local, remote map[string]string) (map[string]string, []string) {
	merged := make(map[string]string)
	conflicts := []string{}

	// Start with all remote fields
	for k, v := range remote {
		merged[k] = v
	}

	// Add local fields, detecting conflicts
	for k, localValue := range local {
		remoteValue, existsInRemote := remote[k]

		if !existsInRemote {
			// Only in local, add it
			merged[k] = localValue
		} else if localValue != remoteValue {
			// Exists in both with different values - conflict
			conflicts = append(conflicts, k)
			// Keep remote value in merged result (though it will be rejected anyway)
		}
	}

	return merged, conflicts
}

// mergeTasks merges two task slices, detecting conflicts.
// Tasks are matched by JiraID. Returns merged tasks and conflict descriptions.
func mergeTasks(local, remote []domain.Task) ([]domain.Task, []string) {
	// Build maps by JiraID for easier comparison
	localMap := make(map[string]domain.Task)
	remoteMap := make(map[string]domain.Task)

	for _, task := range local {
		if task.JiraID != "" {
			localMap[task.JiraID] = task
		}
	}

	for _, task := range remote {
		if task.JiraID != "" {
			remoteMap[task.JiraID] = task
		}
	}

	conflicts := []string{}
	merged := []domain.Task{}

	// Process all remote tasks
	for jiraID, remoteTask := range remoteMap {
		localTask, existsLocally := localMap[jiraID]

		if !existsLocally {
			// New remote task
			merged = append(merged, remoteTask)
		} else {
			// Task exists in both - check for conflicts
			if tasksConflict(localTask, remoteTask) {
				conflicts = append(conflicts, fmt.Sprintf("Task[%s]", jiraID))
				// Use remote version even though there's a conflict
				merged = append(merged, remoteTask)
			} else {
				// No conflict, use remote
				merged = append(merged, remoteTask)
			}
		}

		// Remove from local map to track local-only tasks
		delete(localMap, jiraID)
	}

	// Add any remaining local-only tasks
	for _, task := range localMap {
		merged = append(merged, task)
	}

	return merged, conflicts
}

// tasksConflict checks if two tasks have conflicting changes.
func tasksConflict(local, remote domain.Task) bool {
	// Check each field for conflicts (non-empty differences)
	if local.Title != remote.Title && local.Title != "" && remote.Title != "" {
		return true
	}
	if local.Description != remote.Description && local.Description != "" && remote.Description != "" {
		return true
	}
	if !stringSlicesEqual(local.AcceptanceCriteria, remote.AcceptanceCriteria) &&
		len(local.AcceptanceCriteria) > 0 && len(remote.AcceptanceCriteria) > 0 {
		return true
	}
	if !mapsEqual(local.CustomFields, remote.CustomFields) {
		_, conflicts := mergeCustomFields(local.CustomFields, remote.CustomFields)
		if len(conflicts) > 0 {
			return true
		}
	}

	return false
}

// tasksEqual checks if two task slices are equal.
func tasksEqual(a, b []domain.Task) bool {
	if len(a) != len(b) {
		return false
	}

	// Build maps by JiraID for comparison
	aMap := make(map[string]domain.Task)
	bMap := make(map[string]domain.Task)

	for _, task := range a {
		if task.JiraID != "" {
			aMap[task.JiraID] = task
		}
	}

	for _, task := range b {
		if task.JiraID != "" {
			bMap[task.JiraID] = task
		}
	}

	// If different number of tasks with JiraIDs, they're different
	if len(aMap) != len(bMap) {
		return false
	}

	// Compare each task
	for jiraID, aTask := range aMap {
		bTask, exists := bMap[jiraID]
		if !exists {
			return false
		}

		if !reflect.DeepEqual(aTask, bTask) {
			return false
		}
	}

	return true
}

// stringSlicesEqual checks if two string slices are equal.
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// mapsEqual checks if two string maps are equal.
func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}
