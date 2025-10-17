# State Management in Ticketr

## Overview

Ticketr uses a `.ticketr.state` file to track ticket content and detect changes between local files and JIRA. This enables intelligent synchronization, preventing redundant updates and detecting conflicts.

## State File Format

The state file is a JSON document located at `.ticketr.state` (default) with the following structure:

```json
{
  "TICKET-123": {
    "local_hash": "abc123...",
    "remote_hash": "def456..."
  },
  "TICKET-124": {
    "local_hash": "ghi789...",
    "remote_hash": "jkl012..."
  }
}
```

**Fields:**
- `local_hash`: SHA256 hash of the ticket content in your local Markdown file
- `remote_hash`: SHA256 hash of the ticket content from JIRA's last known state

## Hash Calculation Algorithm

Ticketr calculates SHA256 hashes of ticket content to detect changes. The hash includes:

1. **Ticket metadata**: JiraID, Title, Description, IssueType, Status
2. **Custom fields**: All custom field key-value pairs (sorted alphabetically)
3. **Tasks**: All child task metadata and custom fields (also sorted)

### Deterministic Hashing (Milestone 4)

**Why determinism matters:**
Go's map iteration order is non-deterministic (random). Without sorting, identical tickets could produce different hashes, causing:
- False positive change detection
- Unnecessary JIRA API calls
- User confusion about "unchanged" tickets being updated

**Solution:**
Before hashing custom fields, we extract map keys to a slice, sort them alphabetically, then iterate in sorted order. This ensures identical content always produces identical hashes.

**Implementation:**
```go
// Extract and sort keys
keys := make([]string, 0, len(customFields))
for key := range customFields {
    keys = append(keys, key)
}
sort.Strings(keys)

// Iterate in sorted order
for _, key := range keys {
    hash.Write([]byte(key))
    hash.Write([]byte(fmt.Sprintf("%v", customFields[key])))
}
```

## Conflict Detection

The state manager tracks three hash values per ticket:

1. **Current local hash**: Calculated from your current Markdown file
2. **Stored local hash**: Hash from last successful sync (in `.ticketr.state`)
3. **Stored remote hash**: Hash from JIRA's last known state (in `.ticketr.state`)

### Conflict States

**No Conflict (Safe to Sync):**
- Local unchanged: `current_local == stored_local`
- Remote unchanged: JIRA hash == `stored_remote`
- Only one side changed

**Conflict Detected:**
- Local changed: `current_local != stored_local`
- Remote changed: JIRA hash != `stored_remote`
- **Both sides modified since last sync**

When a conflict is detected, `ticketr pull` will:
1. Report the conflict with ticket IDs
2. Preserve local changes (unless `--force` is used)
3. Suggest using `--force` to overwrite local with remote

## State File Management

**Location:**
- Default: `.ticketr.state` in the current directory
- Configurable via environment variable (future)

**Git Handling:**
- **Recommendation**: Add `.ticketr.state` to `.gitignore`
- **Rationale**: State is environment-specific and includes local file state

**Cleanup:**
- Delete `.ticketr.state` to reset all tracking (next push/pull treats all as new)
- State file automatically created on first use

## Implementation References

**Code Locations:**
- State manager: `/home/karol/dev/private/ticktr/internal/state/manager.go`
- Hash calculation: `StateManager.CalculateHash()` method (line 80+)
- Conflict detection: `StateManager.CheckForConflict()` method (line 27+)
- Pull service integration: `/home/karol/dev/private/ticktr/internal/core/services/pull_service.go`
- Push service integration: `/home/karol/dev/private/ticktr/internal/core/services/push_service.go`

**Tests:**
- State manager tests: `/home/karol/dev/private/ticktr/internal/state/manager_test.go`
- Determinism tests: `TestStateManager_DeterministicHashing`, etc.

## Troubleshooting

**Issue: Tickets always detected as changed**
- **Cause**: Likely non-deterministic hashing (pre-Milestone 4)
- **Solution**: Delete `.ticketr.state` and re-sync (format stabilized in 1.0)

**Issue: Conflicts on every pull**
- **Cause**: State file out of sync or missing
- **Solution**: Delete `.ticketr.state`, run `ticketr pull`, then resume normal workflow

**Issue: State file growing large**
- **Cause**: Many tickets tracked over time
- **Solution**: Periodically delete `.ticketr.state` (safe to do, just resets tracking)

## Future Enhancements

- Configurable state file location
- State file versioning and migration
- Automatic state cleanup for deleted tickets
- State file compression for large projects
