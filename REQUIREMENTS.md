# Ticketr Requirements Specification

**Version:** 3.1.1
**Status:** Production-Ready
**Last Updated:** 2025-10-20

## Document Purpose

This document serves as the single, authoritative source of truth for all Ticketr requirements. It consolidates requirements from Phases 1-6 and establishes the foundation for all future development.

**Supersedes:** All previous requirement documents (REQUIREMENTS-v2.md, scattered docs)

---

## Table of Contents

1. [Functional Requirements](#functional-requirements)
   - [Core Ticket Management](#core-ticket-management)
   - [Workspace Management](#workspace-management)
   - [Bulk Operations](#bulk-operations)
   - [JQL Aliases](#jql-aliases)
   - [Sync Strategies](#sync-strategies)
   - [Template System](#template-system)
2. [Non-Functional Requirements](#non-functional-requirements)
   - [Performance](#performance)
   - [Security](#security)
   - [Reliability](#reliability)
   - [Usability](#usability)
3. [TUI/UX Requirements](#tuiux-requirements)
   - [Async Operations](#async-operations)
   - [User Interface](#user-interface)
   - [Visual Polish](#visual-polish)
4. [Integration Requirements](#integration-requirements)
   - [Jira Integration](#jira-integration)
   - [Filesystem Integration](#filesystem-integration)
   - [Database Integration](#database-integration)

---

## Functional Requirements

### Core Ticket Management

#### PROD-001: Markdown-First Ticket Definition
**Priority:** P0 (Critical)
**Status:** Implemented ✅

**Description:**
Users must be able to define tickets using canonical Markdown format with `# TICKET:` headings.

**Acceptance Criteria:**
- Parser recognizes `# TICKET:` as the canonical heading format
- Supports standard sections: Description, Acceptance Criteria, Tasks, Fields
- Hierarchical task nesting with inheritance
- Rejects legacy `# STORY:` format with helpful error message pointing to migration command

**Traceability:**
- README.md lines 271-307
- Implemented in: `internal/parser/parser.go`
- Tests: `internal/parser/parser_test.go` (TC-001 through TC-050)

---

#### PROD-002: Bidirectional Jira Sync
**Priority:** P0 (Critical)
**Status:** Implemented ✅

**Description:**
Users must be able to push tickets to Jira and pull tickets from Jira with full fidelity.

**Acceptance Criteria:**
- **Push**: Creates new Jira tickets from Markdown files
- **Push**: Updates existing Jira tickets when changes detected
- **Pull**: Fetches tickets from Jira with all fields and subtasks
- **Pull**: Writes tickets to Markdown preserving hierarchy
- Round-trip compatibility: Push → Jira → Pull produces identical structure

**Traceability:**
- README.md lines 292-297, 329-366
- Implemented in: `internal/core/services/push_service.go`, `internal/core/services/pull_service.go`
- Tests: `internal/core/services/push_service_test.go`, `internal/core/services/pull_service_test.go`

---

#### PROD-003: Field Inheritance
**Priority:** P0 (Critical)
**Status:** Implemented ✅

**Description:**
Tasks must inherit parent ticket custom fields, with task-specific fields overriding parent values.

**Acceptance Criteria:**
- Tasks inherit all parent CustomFields by default
- Task-level field declarations override parent values for those specific fields
- Inheritance applied before push operations to Jira
- Pull operations preserve inherited vs. explicit field declarations

**Traceability:**
- README.md lines 371-386
- ROADMAP.md Milestone 7
- Implemented in: `internal/core/services/ticket_service.go` (calculateFinalFields)
- Tests: TC-701.1, TC-701.2, TC-701.3, TC-701.4
- Integration test: docs/integration-test-results-milestone-7.md

---

#### PROD-004: Execution Logging
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
All operations must be logged to disk with timestamps, command parameters, execution summary, and errors.

**Acceptance Criteria:**
- Logs written to platform-standard cache directory (`.ticketr/logs/` or equivalent)
- Log filename format: `YYYY-MM-DD_HH-MM-SS.log`
- Sections: Command Parameters, Execution Summary, Errors, Timestamps
- Automatic log rotation (keep last 10 files)
- Sensitive data (API keys, emails, passwords) automatically redacted
- Configurable log directory via `TICKETR_LOG_DIR` environment variable

**Traceability:**
- README.md lines 429-445
- ROADMAP.md Milestone 6
- Implemented in: `internal/logging/`
- Tests: `internal/logging/logger_test.go`

---

#### PROD-005: State-Aware Push
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
Push operations must detect unchanged tickets and skip API calls to improve performance.

**Acceptance Criteria:**
- State file tracks SHA256 hash of each ticket content
- Custom fields sorted alphabetically before hashing (deterministic)
- Unchanged tickets skipped during push (zero API calls)
- Changed tickets identified and pushed to Jira
- State file location: Platform-standard data directory

**Traceability:**
- README.md lines 352-366
- ROADMAP.md Milestones 4, 9
- Implemented in: `internal/state/manager.go`, `internal/core/services/push_service.go`
- Tests: `internal/state/manager_test.go` (determinism tests)

---

#### PROD-006: Conflict Detection
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
Pull operations must detect when both local file and remote Jira ticket have changed since last sync.

**Acceptance Criteria:**
- Track both `local_hash` and `remote_hash` in state file
- Detect conflicts: local_changed AND remote_changed
- Safe merge: Auto-update when only one side changed
- `--force` flag overrides conflict detection and accepts remote changes
- Clear error message when conflict detected, suggesting `--force` or manual merge

**Traceability:**
- README.md lines 390-427
- ROADMAP.md Milestone 2
- Implemented in: `internal/core/services/pull_service.go`
- Tests: `internal/core/services/pull_service_test.go` (TC-303.x)

---

#### PROD-007: Force Partial Upload
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
Push operations with `--force-partial-upload` flag must continue processing valid tickets even when some tickets have validation errors.

**Acceptance Criteria:**
- Validation errors downgraded to warnings when flag set
- Valid tickets pushed successfully
- Invalid tickets skipped with detailed error messages
- Exit code 0 for partial success with flag
- Exit code 1 for validation failure without flag
- Helpful tip shown when validation fails without flag

**Traceability:**
- README.md lines 339-341
- ROADMAP.md Milestone 5
- Implemented in: `cmd/ticketr/main.go` (pre-flight validation)
- Tests: TC-501.1, TC-501.2, TC-501.3, TC-501.4

---

#### PROD-008: Pull with Subtasks
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
Pull operations must fetch all subtasks for parent tickets and represent them hierarchically in Markdown.

**Acceptance Criteria:**
- SearchTickets fetches subtasks via JQL query `parent = "PARENT-KEY"`
- Subtasks converted to domain.Task format
- Field mapping applied to subtasks (reverse mapping from Jira field IDs)
- Description and acceptance criteria extracted from Jira issue body
- Round-trip compatibility: Pull → Markdown → Parse → Push preserves hierarchy

**Traceability:**
- README.md lines 271-307
- ROADMAP.md Milestone 8
- Implemented in: `internal/adapters/jira/jira_adapter.go` (fetchSubtasks, parseJiraSubtask)
- Tests: TC-208.1, TC-208.2, TC-208.3, TC-208.4

---

#### PROD-009: Schema Discovery
**Priority:** P2 (Medium)
**Status:** Implemented ✅

**Description:**
Users must be able to discover available Jira fields and generate field mapping configuration.

**Acceptance Criteria:**
- `ticketr schema` command lists all available Jira fields
- Output includes field name, field ID, field type
- `ticketr schema > .ticketr.yaml` generates configuration file
- Generated config includes field_mappings for common fields

**Traceability:**
- README.md lines 516-533
- Implemented in: `cmd/ticketr/main.go` (schema command)
- Tests: `cmd/ticketr/schema_test.go`

---

#### PROD-010: Migration Tool
**Priority:** P2 (Medium)
**Status:** Implemented ✅

**Description:**
Users must be able to migrate legacy `# STORY:` format files to canonical `# TICKET:` format.

**Acceptance Criteria:**
- `ticketr migrate <file>` converts legacy format to canonical
- Dry-run mode (default) previews changes without writing
- `--write` flag applies changes to file
- Line-by-line error reporting with line numbers
- Helpful error messages when legacy format detected in parser

**Traceability:**
- docs/migration-guide.md
- ROADMAP.md Milestone 1
- Implemented in: `cmd/ticketr/migration_commands.go`, `internal/migration/migrator.go`
- Tests: `internal/migration/migrator_test.go`

---

### Workspace Management

#### USER-001: Multi-Workspace Support
**Priority:** P0 (Critical)
**Status:** Implemented ✅

**Description:**
Users must be able to manage multiple Jira projects from a single Ticketr installation using workspaces.

**Acceptance Criteria:**
- Create workspace with Jira URL, project key, credentials
- List all configured workspaces
- Switch between workspaces
- Set default workspace
- Delete workspace (with confirmation)
- Each workspace has isolated credentials stored in OS keychain

**Traceability:**
- README.md lines 59-135
- docs/workspace-management-guide.md
- Implemented in: `internal/core/services/workspace_service.go`
- Database: `internal/adapters/database/workspace_repository.go`
- Tests: `internal/core/services/workspace_service_test.go`

---

#### USER-002: Credential Profiles
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
Users must be able to create reusable credential profiles for use across multiple workspaces.

**Acceptance Criteria:**
- Create credential profile with URL, username, API token
- List available credential profiles
- Create workspace using existing profile
- Credentials stored securely in OS keychain
- Profiles isolated per user account

**Traceability:**
- README.md lines 62-80
- Implemented in: `cmd/ticketr/credentials_commands.go`
- Tests: `cmd/ticketr/workspace_profile_integration_test.go`

---

#### USER-003: TUI Workspace Creation
**Priority:** P2 (Medium)
**Status:** Implemented ✅

**Description:**
Users must be able to create workspaces interactively via TUI with credential profile support.

**Acceptance Criteria:**
- Press `w` in TUI to open workspace creation modal
- Form fields: name, URL, project key, credentials (manual or profile)
- Validation feedback in real-time
- Creation success message or error

**Traceability:**
- README.md lines 99
- Implemented in: `internal/adapters/tui/views/workspace_create_modal.go`

---

### Bulk Operations

#### PROD-011: Bulk Update
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
Users must be able to update multiple tickets simultaneously with field changes.

**Acceptance Criteria:**
- CLI: `ticketr bulk update --ids X,Y,Z --set field=value`
- TUI: Multi-select tickets (Space, 'a', 'A') → Press 'b' → Update Fields modal
- Support multiple field updates in single operation
- Real-time progress feedback with [X/Y] counters
- Maximum 100 tickets per operation
- Best-effort rollback on partial failures
- JQL injection prevention (strict ticket ID validation)

**Traceability:**
- README.md lines 137-154
- PHASE5-COMPLETE.md Week 18
- docs/bulk-operations-guide.md
- Implemented in: `internal/core/services/bulk_operation_service.go`
- CLI: `cmd/ticketr/bulk_commands.go`
- TUI: `internal/adapters/tui/views/bulk_operations_modal.go`
- Tests: 30 tests dedicated to bulk operations

---

#### PROD-012: Bulk Move
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
Users must be able to move multiple tickets to a new parent simultaneously.

**Acceptance Criteria:**
- CLI: `ticketr bulk move --ids X,Y,Z --parent PARENT-KEY`
- TUI: Multi-select tickets → Press 'b' → Move Tickets modal
- Parent ticket validation
- Real-time progress feedback
- Best-effort rollback on failures

**Traceability:**
- README.md lines 156-166
- docs/bulk-operations-guide.md
- Implemented in: Same as PROD-011

---

#### PROD-013: TUI Multi-Select
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
TUI must support interactive multi-select with visual feedback and bulk operation execution.

**Acceptance Criteria:**
- Space bar toggles selection (shows [x] checkbox)
- 'a' selects all visible tickets
- 'A' (Shift+a) deselects all
- Selected count shown in title: "Tickets (3 selected)"
- Border color changes when tickets selected (teal/blue)
- 'b' key opens bulk operations menu
- Real-time progress during execution
- Cancel button stops operation (Esc also works)

**Traceability:**
- README.md lines 187-226
- PHASE5-COMPLETE.md Week 18
- Implemented in: `internal/adapters/tui/views/bulk_operations_modal.go`
- Tests: 11 TUI bulk operation tests

---

### JQL Aliases

#### USER-004: Predefined Aliases
**Priority:** P2 (Medium)
**Status:** Implemented ✅

**Description:**
System must provide predefined JQL aliases for common queries.

**Acceptance Criteria:**
- `mine`: `assignee = currentUser() AND resolution = Unresolved`
- `sprint`: `sprint in openSprints()`
- `blocked`: `status = Blocked OR labels in (blocked)`
- Predefined aliases available in all workspaces
- Predefined aliases cannot be modified or deleted

**Traceability:**
- README.md lines 456-463
- docs/FEATURES/JQL-ALIASES.md lines 45-89
- Implemented in: `internal/core/services/alias_service.go`
- Tests: `internal/core/services/alias_service_test.go`

---

#### USER-005: Custom Aliases
**Priority:** P2 (Medium)
**Status:** Implemented ✅

**Description:**
Users must be able to create, update, list, show, and delete custom JQL aliases.

**Acceptance Criteria:**
- Create: `ticketr alias create <name> "<jql>"` with optional `--description` and `--global` flags
- List: `ticketr alias list` shows all available aliases
- Show: `ticketr alias show <name>` displays full details and expanded JQL
- Update: `ticketr alias update <name> "<new-jql>"` modifies existing alias
- Delete: `ticketr alias delete <name>` removes user-defined alias
- Alias name validation: alphanumeric, hyphens, underscores, 1-64 chars
- JQL query length limit: 2000 characters

**Traceability:**
- README.md lines 465-513
- docs/FEATURES/JQL-ALIASES.md
- Implemented in: `cmd/ticketr/alias_commands.go`, `internal/core/services/alias_service.go`
- Database: `internal/adapters/database/alias_repository.go`
- Tests: 50+ tests for alias functionality

---

#### USER-006: Recursive Aliases
**Priority:** P2 (Medium)
**Status:** Implemented ✅

**Description:**
Users must be able to reference other aliases within alias definitions using `@` syntax.

**Acceptance Criteria:**
- Syntax: `@alias_name` expands to the referenced alias JQL
- Supports multi-level nesting (aliases referencing aliases referencing aliases)
- Circular reference detection prevents infinite loops
- `ticketr alias show` displays both original and expanded JQL for recursive aliases
- Expansion algorithm: O(n) where n is number of nested references

**Traceability:**
- README.md lines 499-511
- docs/FEATURES/JQL-ALIASES.md lines 194-265
- Implemented in: `internal/core/services/alias_service.go` (ExpandAlias method)
- Tests: Recursive expansion and circular reference detection tests

---

#### USER-007: Alias Pull Integration
**Priority:** P2 (Medium)
**Status:** Implemented ✅

**Description:**
Users must be able to use aliases in pull operations for ticket filtering.

**Acceptance Criteria:**
- `ticketr pull --alias <name> --output <file>` uses alias for JQL query
- Alias and JQL flags are mutually exclusive (cannot use both)
- Verbose mode shows expanded JQL: `--verbose` flag displays full expanded query
- Works with other filters: `--alias mine --epic PROJ-100` combines alias with epic filter

**Traceability:**
- README.md lines 486-496
- docs/FEATURES/JQL-ALIASES.md lines 399-448
- Implemented in: `cmd/ticketr/main.go` (pull command integration)

---

### Sync Strategies

#### PROD-014: Sync Strategy System
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
System must provide configurable conflict resolution strategies for pull operations.

**Acceptance Criteria:**
- Three strategies implemented: LocalWins, RemoteWins (default), ThreeWayMerge
- Hash-based conflict detection (SHA256, immune to clock skew)
- Field-level granularity for merge operations
- Default strategy: RemoteWins (backward compatible with v2.x)

**Traceability:**
- README.md lines 390-428
- PHASE5-COMPLETE.md Week 19 Slice 2
- docs/sync-strategies-guide.md
- Implemented in: `internal/core/services/sync_strategies.go`
- Tests: 64 new tests (93.95% coverage)

---

#### PROD-015: LocalWins Strategy
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
LocalWins strategy must preserve all local changes and ignore remote updates.

**Acceptance Criteria:**
- Local changes always take precedence
- Remote updates ignored completely
- Use case: Offline-first workflows, long-running feature branches
- Trade-off: Remote changes from teammates lost

**Traceability:**
- docs/sync-strategies-guide.md
- Implemented in: `internal/core/services/sync_strategies.go` (LocalWinsStrategy)
- Tests: LocalWins strategy tests

---

#### PROD-016: RemoteWins Strategy
**Priority:** P0 (Critical)
**Status:** Implemented ✅

**Description:**
RemoteWins strategy must accept all remote changes and discard local edits.

**Acceptance Criteria:**
- Remote changes always take precedence
- Local edits discarded if remote changed
- Default strategy (backward compatible)
- Use case: Jira as single source of truth
- Trade-off: Local edits lost if remote changed

**Traceability:**
- docs/sync-strategies-guide.md
- Implemented in: `internal/core/services/sync_strategies.go` (RemoteWinsStrategy)
- Tests: RemoteWins strategy tests (default behavior)

---

#### PROD-017: ThreeWayMerge Strategy
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
ThreeWayMerge strategy must intelligently merge compatible changes and error on conflicts.

**Acceptance Criteria:**
- Auto-merge compatible changes (different fields modified)
- Error on incompatible changes (same field modified differently)
- Field-level granularity (not all-or-nothing)
- Custom field per-key merging
- Task merging by JiraID with recursive conflict detection
- Examples:
  - Compatible: Local description + Remote status = Both preserved
  - Incompatible: Local title="Fix bug", Remote title="Auth improvements" = Error

**Traceability:**
- README.md lines 398-412
- PHASE5-COMPLETE.md Week 19 Slice 2
- docs/sync-strategies-guide.md
- Implemented in: `internal/core/services/sync_strategies.go` (ThreeWayMergeStrategy)
- Tests: 64 comprehensive test scenarios including compatible/incompatible changes

---

### Template System

#### PROD-018: Template Parser
**Priority:** P2 (Medium)
**Status:** Implemented ✅ (Parser only, CLI deferred to v3.1.1)

**Description:**
System must parse YAML-based templates with variable substitution engine.

**Acceptance Criteria:**
- YAML template parsing with schema validation
- Variable substitution: `{{.Name}}`, `{{.Sprint}}`, `{{.Priority}}`
- Support for nested structures (epics, stories, tasks)
- Variable extraction and validation
- Deep copy safety for template reuse
- Example template structure:
  ```yaml
  title: "Feature: {{.Name}}"
  description: |
    As a {{.Actor}}
    I want {{.Goal}}
    So that {{.Benefit}}
  stories:
    - title: "Implement {{.Component}}"
      tasks:
        - "Unit tests"
        - "Integration tests"
  ```

**Traceability:**
- PHASE5-COMPLETE.md Week 19 Slice 1
- Implemented in: `internal/templates/parser.go`
- Tests: 32 parser tests (~85% coverage)

**Note:** CLI integration (`ticketr template apply`, `ticketr template list`, TUI template selector) deferred to v3.1.1 for UX polish.

---

## Non-Functional Requirements

### Performance

#### NFR-001: TUI Responsiveness
**Priority:** P0 (Critical)
**Status:** Implemented ✅

**Target:** TUI must render 1000+ tickets in less than 100ms

**Acceptance Criteria:**
- Tested with datasets of 1000+ tickets
- UI rendering completes in <100ms
- No lag or freezing during navigation
- Smooth scrolling and interaction

**Traceability:**
- PHASE5-COMPLETE.md Week 20 Day 3-5
- Benchmarks: TUI renders 1000+ tickets in ~85ms

---

#### NFR-002: Sync Performance
**Priority:** P1 (High)
**Status:** Implemented ✅

**Target:** Overhead for conflict detection must be minimal (<5ms per ticket)

**Acceptance Criteria:**
- LocalWins overhead: <1ms per ticket
- RemoteWins overhead: <1ms per ticket
- ThreeWayMerge overhead: <2ms per ticket
- State hash calculation: <1ms per ticket

**Traceability:**
- PHASE5-COMPLETE.md Week 19 Slice 2
- Benchmarks documented in PHASE5-COMPLETE.md

---

#### NFR-003: Alias Expansion Performance
**Priority:** P2 (Medium)
**Status:** Implemented ✅

**Target:** Alias expansion must complete in <10ms regardless of nesting depth

**Acceptance Criteria:**
- Single reference: <1ms
- 3-level recursion: ~3ms
- Complex chain (5 levels): <5ms
- No performance degradation with circular reference detection

**Traceability:**
- PHASE5-COMPLETE.md Week 20 Slice 1
- Benchmarks: All targets exceeded

---

#### NFR-004: Database Query Performance
**Priority:** P2 (Medium)
**Status:** Implemented ✅

**Target:** Database operations must complete in milliseconds

**Acceptance Criteria:**
- Alias lookup: <1ms (indexed by name, workspace_id)
- Workspace switch: <10ms
- State update: <5ms
- Indexes on frequently queried columns

**Traceability:**
- docs/ARCHITECTURE.md (Database schema with indexes)
- Implemented in: SQLite with appropriate indexes

---

### Security

#### NFR-005: Credential Security
**Priority:** P0 (Critical)
**Status:** Implemented ✅

**Description:**
Credentials must never be stored in plaintext or logged.

**Acceptance Criteria:**
- Credentials stored in OS-level keychain (macOS Keychain, Windows Credential Manager, Linux Secret Service)
- Database contains only CredentialRef (keychain ID and service ID)
- No credentials in logs or error messages (automatic redaction)
- Credentials cleared from memory after use
- Per-user isolation via OS keychain

**Traceability:**
- README.md lines 92-98, 129-134
- docs/ARCHITECTURE.md lines 356-393
- Implemented in: `internal/adapters/keychain/keychain_store.go`
- Tests: Keychain integration tests

---

#### NFR-006: JQL Injection Prevention
**Priority:** P0 (Critical)
**Status:** Implemented ✅

**Description:**
System must prevent JQL injection attacks via strict input validation.

**Acceptance Criteria:**
- Ticket IDs validated with regex: `^[A-Z]+-\d+$`
- Alias names validated: alphanumeric, hyphens, underscores only
- JQL query length limit: 2000 characters
- No user input directly concatenated into JQL queries

**Traceability:**
- PHASE5-COMPLETE.md Week 18 (Bulk Operations safety features)
- docs/bulk-operations-guide.md
- Implemented in: `internal/core/domain/bulk_operation.go` (validation)

---

#### NFR-007: Sensitive Data Redaction
**Priority:** P0 (Critical)
**Status:** Implemented ✅

**Description:**
Logs must automatically redact sensitive data patterns.

**Acceptance Criteria:**
- API keys redacted: Patterns like `ATATT`, `Bearer`, `token`
- Email addresses redacted
- Password fields redacted
- Redaction applied before writing to log file
- Redaction preserves log readability (shows [REDACTED])

**Traceability:**
- README.md lines 429-445
- Implemented in: `internal/logging/redaction.go`
- Tests: `internal/logging/redaction_test.go`

---

### Reliability

#### NFR-008: Error Handling
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
All errors must be handled gracefully with clear user messages.

**Acceptance Criteria:**
- No panics in production code (recover at top level)
- Error messages include context and suggested fixes
- Exit codes: 0 (success), 1 (validation failure), 2 (runtime error)
- Detailed errors logged, user-friendly messages shown

**Traceability:**
- README.md lines 664-670
- Implemented throughout codebase with consistent error wrapping

---

#### NFR-009: Graceful Degradation
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
System must handle partial failures without corrupting data.

**Acceptance Criteria:**
- Bulk operations: Best-effort rollback on partial failures
- Pull operations: Non-fatal subtask fetch errors return parent tickets
- State file: Atomic writes prevent corruption
- Log file: Rotation failure doesn't block logging

**Traceability:**
- PHASE5-COMPLETE.md (Bulk Operations rollback)
- PHASE5-COMPLETE.md (Pull with subtasks error handling)
- Implemented in: `internal/core/services/bulk_operation_service.go`, `internal/adapters/jira/jira_adapter.go`

---

### Usability

#### NFR-010: Discoverability
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
Users must be able to discover features without reading documentation.

**Acceptance Criteria:**
- `--help` flag for all commands with examples
- TUI: Press `?` for help screen with all keybindings
- Error messages include suggestions for resolution
- Examples directory with common use cases

**Traceability:**
- README.md comprehensive command reference
- examples/ directory with templates
- Implemented in: Cobra command help text, TUI help modal

---

#### NFR-011: Installation Simplicity
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
Installation must be simple with minimal dependencies.

**Acceptance Criteria:**
- Single command: `go install github.com/karolswdev/ticketr/cmd/ticketr@latest`
- Or download binary: No compilation required
- No external dependencies beyond Go standard library and listed packages
- Docker image available: Lightweight (<15MB)

**Traceability:**
- README.md lines 36-45
- Dockerfile with multi-stage build
- GitHub releases with pre-built binaries

---

## TUI/UX Requirements

### Async Operations

#### UX-001: Non-Blocking Operations
**Priority:** P0 (Critical)
**Status:** Planned (Phase 6 Week 2)

**Description:**
Long-running operations must not block the TUI, allowing users to navigate during execution.

**Acceptance Criteria:**
- Job queue system implemented with goroutine worker pool
- Pull operations submitted to job queue (non-blocking)
- User can navigate TUI views while pull is running
- Progress updates received via channel
- Job status visible in UI (spinner/progress bar)

**Traceability:**
- PHASE6-CLEAN-RELEASE.md Day 6-7
- To be implemented in: `internal/tui/jobs/queue.go`

---

#### UX-002: Cancellation Support
**Priority:** P0 (Critical)
**Status:** Planned (Phase 6 Week 2)

**Description:**
Users must be able to cancel long-running operations gracefully.

**Acceptance Criteria:**
- ESC key cancels active job
- Ctrl+C cancels active job and exits cleanly
- Cancel triggers context cancellation
- Partial results handled gracefully
- UI shows "Cancelling..." state during cleanup
- No orphaned goroutines after cancel

**Traceability:**
- PHASE6-CLEAN-RELEASE.md Day 6-7
- To be implemented in: `internal/tui/jobs/queue.go` (context cancellation)

---

#### UX-003: Progress Feedback
**Priority:** P1 (High)
**Status:** Planned (Phase 6 Week 2)

**Description:**
Users must receive real-time progress feedback during long operations.

**Acceptance Criteria:**
- Progress bar showing current/total counts: [=====>    ] 50% (45/120)
- Ticket count display: "45/120 tickets"
- Time elapsed: "Elapsed: 12s"
- ETA calculation: "ETA: 15s" (linear extrapolation)
- Progress bar adapts to terminal width
- Updates smooth (throttled to avoid flicker)

**Traceability:**
- PHASE6-CLEAN-RELEASE.md Day 10-11
- To be implemented in: `internal/tui/widgets/progressbar.go`

---

### User Interface

#### UX-004: Context-Aware Action Bar
**Priority:** P1 (High)
**Status:** Planned (Phase 6 Week 2)

**Description:**
TUI must display context-aware keybindings in a persistent bottom action bar.

**Acceptance Criteria:**
- Action bar visible in all views (full-width, 1-2 rows)
- Shows keybindings relevant to current view
- Examples: List view shows [Enter Open] [Space Select], Detail view shows [Esc Back] [E Edit]
- Updates dynamically when view changes
- Clear, concise labels

**Traceability:**
- PHASE6-CLEAN-RELEASE.md Day 8-9
- To be implemented in: `internal/tui/widgets/actionbar.go`

---

#### UX-005: Command Palette
**Priority:** P2 (Medium)
**Status:** Planned (Phase 6 Week 2)

**Description:**
TUI must provide a searchable command palette for quick access to all actions.

**Acceptance Criteria:**
- Triggered by Ctrl+P or F1
- Fuzzy search commands by name or description
- Shows command description and keybinding
- Commands grouped by category (View, Edit, Sync, etc.)
- Execute command on selection

**Traceability:**
- PHASE6-CLEAN-RELEASE.md Day 8-9
- To be implemented in: `internal/tui/widgets/palette.go`

---

#### UX-006: Keybinding Consistency
**Priority:** P1 (High)
**Status:** Planned (Phase 6 Week 2)

**Description:**
Keybindings must be consistent across all TUI views following standard conventions.

**Acceptance Criteria:**
- Esc: Back/Cancel (universal)
- Enter: Confirm/Open (universal)
- F1: Help/Command Palette
- F2: Sync/Pull
- F5: Refresh view
- No conflicting bindings within same context
- Documented in help screen and keybindings reference

**Traceability:**
- PHASE6-CLEAN-RELEASE.md Day 8-9
- To be implemented: TUI event handler updates

---

### Visual Polish

#### UX-007: Professional Appearance
**Priority:** P2 (Medium)
**Status:** Planned (Phase 6 Day 12.5)

**Description:**
TUI must have polished visual appearance with subtle animations and effects.

**Acceptance Criteria:**
- Active spinner animations (⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏)
- Focus pulse on focused panels
- Modal fade-in effects (100ms dithered: ░→▒→█)
- Drop shadows on modals (▒ offset characters)
- Border styles vary by focus state (double-line focused, single-line unfocused)
- Title gradients on focused panels
- Optional background effects (hyperspace, snow) - default OFF

**Traceability:**
- PHASE6-CLEAN-RELEASE.md Day 12.5 (Director's Cut)
- To be implemented in: `internal/adapters/tui/effects/`

---

#### UX-008: Performance with Effects
**Priority:** P1 (High)
**Status:** Planned (Phase 6 Day 12.5)

**Description:**
Visual effects must not degrade TUI performance.

**Acceptance Criteria:**
- CPU usage minimal with all effects enabled
- UI remains responsive during animations
- Effects system optional (theme-controlled)
- Zero performance impact when effects disabled
- Graceful degradation on limited terminals

**Traceability:**
- PHASE6-CLEAN-RELEASE.md Day 12.5
- Performance target: Same responsiveness with/without effects

---

## Integration Requirements

### Jira Integration

#### INTG-001: Jira API Authentication
**Priority:** P0 (Critical)
**Status:** Implemented ✅

**Description:**
System must authenticate with Jira API using secure credentials.

**Acceptance Criteria:**
- Basic Auth with email and API token
- API token retrieved from OS keychain
- Authentication failure shows helpful error (invalid credentials, expired token)
- Support for Jira Cloud (atlassian.net)

**Traceability:**
- README.md lines 47-54
- Implemented in: `internal/adapters/jira/jira_adapter.go`

---

#### INTG-002: Dynamic Field Mapping
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
System must map human-readable field names to Jira custom field IDs dynamically.

**Acceptance Criteria:**
- `.ticketr.yaml` config file defines field mappings
- Supports field types: string, number, array, date, user
- Reverse mapping for pull operations (Jira field ID → human name)
- Schema discovery generates mapping configuration
- Type conversion (string → number, array, etc.)

**Traceability:**
- README.md lines 516-533
- Implemented in: `internal/adapters/jira/field_mapper.go`
- Tests: `internal/adapters/jira/jira_test.go`

---

#### INTG-003: Issue Type Hierarchy
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
System must respect Jira issue type hierarchy rules.

**Acceptance Criteria:**
- Validate parent-child relationships (e.g., Epic cannot have Sub-tasks)
- Reject invalid hierarchies with clear error messages
- Support common types: Epic, Story, Task, Sub-task, Bug

**Traceability:**
- Implemented in: `internal/core/validation/hierarchy_validator.go`, `internal/adapters/jira/hierarchy.go`
- Tests: Hierarchy validation tests

---

### Filesystem Integration

#### INTG-004: Atomic File Writes
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
File write operations must be atomic to prevent corruption.

**Acceptance Criteria:**
- Write to temporary file first
- Rename to target file (atomic on most filesystems)
- Original file preserved if write fails
- No partial writes visible to user

**Traceability:**
- Implemented in: `internal/adapters/filesystem/file_repository.go`
- Tests: Filesystem adapter tests

---

#### INTG-005: Platform-Standard Paths
**Priority:** P0 (Critical)
**Status:** Implemented ✅

**Description:**
System must use platform-standard directories for config, data, and cache.

**Acceptance Criteria:**
- Linux/Unix: `~/.config/ticketr/`, `~/.local/share/ticketr/`, `~/.cache/ticketr/`
- macOS: `~/Library/Application Support/ticketr/`, `~/Library/Preferences/ticketr/`, `~/Library/Caches/ticketr/`
- Windows: `%LOCALAPPDATA%\ticketr\`, `%APPDATA%\ticketr\`, `%TEMP%\ticketr\`
- Automatic migration from v2.x local paths on first run
- `ticketr migrate-paths` and `ticketr rollback-paths` commands available

**Traceability:**
- README.md lines 228-268
- docs/v3-MIGRATION-GUIDE.md
- Implemented in: `internal/core/services/path_resolver.go`
- Tests: `internal/core/services/path_resolver_test.go`

---

### Database Integration

#### INTG-006: SQLite Backend
**Priority:** P0 (Critical)
**Status:** Implemented ✅

**Description:**
System must use SQLite for persistent storage of workspaces, aliases, and state.

**Acceptance Criteria:**
- Single database file: `ticketr.db` in platform data directory
- Tables: workspaces, jql_aliases (with migrations for future tables)
- Indexes on frequently queried columns: (name, workspace_id)
- Foreign key constraints with CASCADE delete
- Transaction support for multi-step operations

**Traceability:**
- docs/ARCHITECTURE.md lines 315-369
- Implemented in: `internal/adapters/database/`
- Migration system: `internal/adapters/database/migrations.go`

---

#### INTG-007: Database Migrations
**Priority:** P1 (High)
**Status:** Implemented ✅

**Description:**
System must support schema migrations for database evolution.

**Acceptance Criteria:**
- Migration files tracked and versioned
- Automatic migration on startup (checks schema version)
- Backward-compatible migrations (don't break older versions)
- Migration failures logged with rollback instructions

**Traceability:**
- Implemented in: `internal/adapters/database/migrations.go`
- Tests: Migration tests verify schema transitions

---

## Requirements Traceability Matrix

| Requirement ID | Priority | Status | Implementation | Tests | Documentation |
|---|---|---|---|---|---|
| PROD-001 | P0 | ✅ | parser.go | TC-001..TC-050 | README.md |
| PROD-002 | P0 | ✅ | push/pull_service.go | TC-201..TC-250 | README.md |
| PROD-003 | P0 | ✅ | ticket_service.go | TC-701.x | README.md, Milestone 7 |
| PROD-004 | P1 | ✅ | logging/ | logging tests | README.md |
| PROD-005 | P1 | ✅ | state/manager.go | TC-801..TC-850 | README.md |
| PROD-006 | P1 | ✅ | pull_service.go | TC-303.x | README.md |
| PROD-007 | P1 | ✅ | main.go | TC-501.x | README.md |
| PROD-008 | P1 | ✅ | jira_adapter.go | TC-208.x | README.md |
| PROD-009 | P2 | ✅ | main.go (schema) | schema_test.go | README.md |
| PROD-010 | P2 | ✅ | migration_commands.go | migrator_test.go | migration-guide.md |
| USER-001 | P0 | ✅ | workspace_service.go | workspace tests | workspace-management-guide.md |
| USER-002 | P1 | ✅ | credentials_commands.go | integration tests | README.md |
| USER-003 | P2 | ✅ | workspace_create_modal.go | TUI tests | README.md |
| PROD-011 | P1 | ✅ | bulk_operation_service.go | 30 tests | bulk-operations-guide.md |
| PROD-012 | P1 | ✅ | bulk_operation_service.go | bulk tests | bulk-operations-guide.md |
| PROD-013 | P1 | ✅ | bulk_operations_modal.go | 11 TUI tests | README.md |
| USER-004 | P2 | ✅ | alias_service.go | alias tests | JQL-ALIASES.md |
| USER-005 | P2 | ✅ | alias_commands.go | 50+ tests | JQL-ALIASES.md |
| USER-006 | P2 | ✅ | alias_service.go | recursive tests | JQL-ALIASES.md |
| USER-007 | P2 | ✅ | main.go (pull) | pull tests | JQL-ALIASES.md |
| PROD-014 | P1 | ✅ | sync_strategies.go | 64 tests | sync-strategies-guide.md |
| PROD-015 | P1 | ✅ | sync_strategies.go | strategy tests | sync-strategies-guide.md |
| PROD-016 | P0 | ✅ | sync_strategies.go | strategy tests | sync-strategies-guide.md |
| PROD-017 | P1 | ✅ | sync_strategies.go | merge tests | sync-strategies-guide.md |
| PROD-018 | P2 | ✅ | templates/parser.go | 32 tests | PHASE5-COMPLETE.md |
| UX-001 | P0 | 📋 | Phase 6 Week 2 | Planned | PHASE6-CLEAN-RELEASE.md |
| UX-002 | P0 | 📋 | Phase 6 Week 2 | Planned | PHASE6-CLEAN-RELEASE.md |
| UX-003 | P1 | 📋 | Phase 6 Week 2 | Planned | PHASE6-CLEAN-RELEASE.md |
| UX-004 | P1 | 📋 | Phase 6 Week 2 | Planned | PHASE6-CLEAN-RELEASE.md |
| UX-005 | P2 | 📋 | Phase 6 Week 2 | Planned | PHASE6-CLEAN-RELEASE.md |
| UX-006 | P1 | 📋 | Phase 6 Week 2 | Planned | PHASE6-CLEAN-RELEASE.md |
| UX-007 | P2 | 📋 | Phase 6 Day 12.5 | Planned | PHASE6-CLEAN-RELEASE.md |
| UX-008 | P1 | 📋 | Phase 6 Day 12.5 | Planned | PHASE6-CLEAN-RELEASE.md |
| NFR-001 | P0 | ✅ | TUI rendering | TUI benchmarks | PHASE5-COMPLETE.md |
| NFR-002 | P1 | ✅ | sync_strategies.go | strategy tests | PHASE5-COMPLETE.md |
| NFR-003 | P2 | ✅ | alias_service.go | alias tests | PHASE5-COMPLETE.md |
| NFR-004 | P2 | ✅ | SQLite indexes | DB tests | ARCHITECTURE.md |
| NFR-005 | P0 | ✅ | keychain_store.go | keychain tests | ARCHITECTURE.md |
| NFR-006 | P0 | ✅ | bulk_operation.go | validation tests | bulk-operations-guide.md |
| NFR-007 | P0 | ✅ | logging/redaction.go | redaction_test.go | README.md |
| NFR-008 | P1 | ✅ | Error handling | All test suites | README.md |
| NFR-009 | P1 | ✅ | bulk/pull services | rollback tests | PHASE5-COMPLETE.md |
| NFR-010 | P1 | ✅ | Discoverability | Manual testing | README.md |
| NFR-011 | P1 | ✅ | Installation | README.md | Dockerfile |
| INTG-001 | P0 | ✅ | jira_adapter.go | jira auth tests | README.md |
| INTG-002 | P1 | ✅ | field_mapper.go | jira_test.go | README.md |
| INTG-003 | P1 | ✅ | hierarchy_validator.go | hierarchy tests | README.md |
| INTG-004 | P1 | ✅ | file_repository.go | filesystem tests | README.md |
| INTG-005 | P0 | ✅ | path_resolver.go | path_resolver_test.go | v3-MIGRATION-GUIDE.md |
| INTG-006 | P0 | ✅ | database/ | migration tests | ARCHITECTURE.md |
| INTG-007 | P1 | ✅ | migrations.go | migration tests | ARCHITECTURE.md |

---

## Excluded Requirements (Migration/Compatibility)

The following requirements from previous versions are **intentionally excluded** from this specification as they are no longer relevant in v3.1.1:

### v2.x Compatibility
- **EXCLUDED**: Legacy `# STORY:` format support (replaced with migration tool)
- **EXCLUDED**: Backward-compatible state file format (v3 uses new format)
- **EXCLUDED**: v2 CLI flag compatibility (v3 has cleaner flags)

### Migration Features
- **EXCLUDED**: Automatic v2-to-v3 migration on first run (use explicit `migrate-paths` command)
- **EXCLUDED**: v2 fallback behavior (v3 is the only supported version)
- **EXCLUDED**: Feature flags for beta/rc features (v3.1.1 ships everything stable)

### Rationale
Phase 6 removes all migration code and compatibility layers. Ticketr v3.1.1 is a clean, production-ready release with no conditional behaviors. Users upgrading from v2.x should use migration tools explicitly before upgrading.

---

## Document Maintenance

### Update Policy
This document must be updated whenever:
1. New features are added (add new requirement IDs)
2. Requirements change (update acceptance criteria)
3. Features are removed (move to "Excluded Requirements" with rationale)
4. Implementation status changes (update status column in traceability matrix)

### Review Schedule
- **Minor reviews**: Each sprint/phase completion
- **Major reviews**: Each major version release (3.0, 4.0, etc.)
- **Owner**: Builder agent (primary), reviewed by Steward agent (architecture)

### Version History
- **1.0** (2025-10-20): Initial consolidated requirements for v3.1.1
  - Consolidated from: REQUIREMENTS-v2.md, README.md, PHASE5-COMPLETE.md, ROADMAP.md
  - Added: Phase 6 TUI/UX requirements
  - Removed: Migration/compatibility requirements

---

## References

### Primary Documentation
- [README.md](README.md) - User-facing documentation
- [ARCHITECTURE.md](docs/ARCHITECTURE.md) - Technical architecture
- [PHASE5-COMPLETE.md](docs/PHASE5-COMPLETE.md) - Phase 5 feature completion report
- [ROADMAP.md](docs/development/ROADMAP.md) - Development milestones
- [PHASE6-CLEAN-RELEASE.md](docs/PHASE6-CLEAN-RELEASE.md) - Phase 6 execution plan

### Feature Guides
- [bulk-operations-guide.md](docs/bulk-operations-guide.md) - Bulk operations user guide
- [sync-strategies-guide.md](docs/sync-strategies-guide.md) - Sync strategies user guide
- [JQL-ALIASES.md](docs/FEATURES/JQL-ALIASES.md) - JQL aliases comprehensive guide
- [workspace-management-guide.md](docs/workspace-management-guide.md) - Workspace management guide

### Migration Guides
- [migration-guide.md](docs/migration-guide.md) - v1 to v2 migration (legacy format)
- [v3-MIGRATION-GUIDE.md](docs/v3-MIGRATION-GUIDE.md) - v2 to v3 migration (PathResolver)

---

**End of Requirements Specification**

This document is the authoritative source of truth for all Ticketr requirements. Any conflicts between this document and other documentation should be resolved in favor of this specification.
