# Ticketr Architecture Overview

**Last Updated:** October 21, 2025 (Phase 7 - Jira Library Integration)
**Version:** 3.1.1 (Simplified v3 Architecture)
**Status:** Production-ready, actively maintained

**Recent Changes (v3.1.1)**:
- Integrated `andygrunwald/go-jira` v1.17.0 library for Jira adapter (Phase 7)
- Added feature flag system for V1/V2 adapter selection
- Reduced Jira adapter code by 33% (1,136 → 757 lines)
- Removed all migration code and feature flags (637 lines removed)
- All v3 features now enabled by default
- Simplified initialization paths with no conditional logic

> **Requirements Specification:** See [REQUIREMENTS.md](../REQUIREMENTS.md) for the complete, authoritative specification of all 42 requirements with acceptance criteria and traceability.

## Executive Summary

Ticketr is a command-line tool that bridges Markdown files and Jira, enabling bidirectional synchronization of tickets and tasks. The system is built using Hexagonal Architecture (Ports and Adapters) with a focus on maintainability, testability, and extensibility.

**Key Capabilities:**
- Parse Markdown files into structured ticket data
- Push tickets/tasks to Jira with hierarchical field inheritance
- Pull tickets/tasks from Jira with automatic subtask fetching
- State-aware operations (skip unchanged tickets)
- Conflict detection and intelligent sync strategies
- Dynamic field mapping and schema discovery
- Comprehensive logging and validation
- Multi-workspace management with secure credential storage
- Bulk operations with real-time progress feedback
- JQL aliases for reusable query patterns

---

## Architecture Pattern: Hexagonal (Ports & Adapters)

Ticketr follows the Hexagonal Architecture pattern to decouple core business logic from external concerns like Jira APIs, file systems, and CLI interfaces.

### Core Layers

```
┌─────────────────────────────────────────────────────────────┐
│                        CLI Layer                             │
│                     (cmd/ticketr/)                           │
│  Entry point: Cobra commands (push, pull, schema, migrate)  │
└─────────────────┬───────────────────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────────────────┐
│                     Adapters Layer                           │
│                  (internal/adapters/)                        │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Filesystem  │  │     Jira     │  │     CLI      │      │
│  │   Adapter    │  │   Adapter    │  │   Adapter    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────┬───────────────────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────────────────┐
│                      Ports Layer                             │
│                   (internal/core/ports/)                     │
│                                                              │
│  Interfaces: Repository, JiraPort, Logger, Renderer         │
└─────────────────┬───────────────────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────────────────┐
│                   Core Business Logic                        │
│                  (internal/core/services/)                   │
│                                                              │
│  ┌──────────────────┐  ┌──────────────────┐               │
│  │  TicketService   │  │  PushService     │               │
│  │  (orchestration) │  │  (state-aware)   │               │
│  └──────────────────┘  └──────────────────┘               │
│  ┌──────────────────┐  ┌──────────────────┐               │
│  │  PullService     │  │  SchemaService   │               │
│  │  (with conflicts)│  │  (discovery)     │               │
│  └──────────────────┘  └──────────────────┘               │
└─────────────────┬───────────────────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────────────────┐
│                      Domain Layer                            │
│                  (internal/core/domain/)                     │
│                                                              │
│  Models: Ticket, Task, CustomFields, ValidationResult       │
└──────────────────────────────────────────────────────────────┘
```

### Supporting Components

```
┌─────────────────────────────────────────────────────────────┐
│                  Cross-Cutting Concerns                      │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Parser     │  │    State     │  │   Logging    │      │
│  │  (Markdown)  │  │  Management  │  │   System     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Validation  │  │  Renderer    │  │  Migration   │      │
│  │   Engine     │  │  (Markdown)  │  │   Engine     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└──────────────────────────────────────────────────────────────┘
```

---

## Directory Structure

```
ticketr/
├── cmd/ticketr/                      # CLI entry point (Cobra)
│   ├── main.go                       # Command definitions
│   ├── push.go                       # Push command handler
│   ├── pull.go                       # Pull command handler
│   ├── schema.go                     # Schema discovery command
│   └── migrate.go                    # Migration command
│
├── internal/
│   ├── core/                         # Core business logic (Hexagon center)
│   │   ├── domain/                   # Domain models
│   │   │   ├── models.go             # Ticket, Task, CustomFields
│   │   │   └── validation_result.go  # Validation types
│   │   ├── ports/                    # Interface definitions
│   │   │   ├── repository.go         # File operations interface
│   │   │   ├── jira_port.go          # Jira operations interface
│   │   │   ├── logger.go             # Logging interface
│   │   │   └── renderer.go           # Markdown rendering interface
│   │   ├── services/                 # Business logic services
│   │   │   ├── ticket_service.go     # Main orchestration
│   │   │   ├── push_service.go       # State-aware push logic
│   │   │   ├── pull_service.go       # Pull with conflict detection
│   │   │   └── schema_service.go     # Field mapping discovery
│   │   └── validation/               # Validation rules
│   │       ├── hierarchy_validator.go # Issue type hierarchy rules
│   │       └── field_validator.go     # Required field checks
│   │
│   ├── adapters/                     # External integrations (Hexagon edges)
│   │   ├── filesystem/               # File I/O adapter
│   │   │   ├── filesystem_adapter.go # Repository implementation
│   │   │   └── filesystem_test.go
│   │   ├── jira/                     # Jira API adapter
│   │   │   ├── jira_adapter.go       # JiraPort implementation
│   │   │   ├── field_mapper.go       # Dynamic field mapping
│   │   │   ├── hierarchy.go          # Issue type hierarchy
│   │   │   └── jira_test.go
│   │   └── cli/                      # CLI utilities
│   │       └── output.go             # Formatted output
│   │
│   ├── parser/                       # Markdown parsing
│   │   ├── parser.go                 # Main parser logic
│   │   ├── parser_test.go
│   │   └── sections.go               # Section-aware parsing
│   │
│   ├── renderer/                     # Markdown rendering
│   │   ├── renderer.go               # Ticket → Markdown
│   │   └── renderer_test.go
│   │
│   ├── state/                        # State management
│   │   ├── manager.go                # Hash tracking, conflict detection
│   │   └── manager_test.go
│   │
│   ├── logging/                      # Logging system
│   │   ├── logger.go                 # File-based logging
│   │   ├── rotation.go               # Log rotation
│   │   └── redaction.go              # Sensitive data filtering
│   │
│   └── migration/                    # Format migration
│       ├── migrator.go               # STORY → TICKET conversion
│       └── migrator_test.go
│
├── testdata/                         # Test fixtures
│   ├── legacy_story/                 # Legacy format samples
│   ├── valid_tickets/                # Valid test cases
│   └── invalid_tickets/              # Invalid test cases
│
├── tests/                            # Integration and smoke tests
│   └── smoke/                        # End-to-end smoke tests
│       ├── smoke_test.sh
│       └── README.md
│
├── examples/                         # User-facing examples
│   ├── quick-story.md
│   ├── epic-template.md
│   ├── sprint-template.md
│   ├── field-inheritance-example.md
│   └── pull-with-subtasks-example.md
│
├── docs/                             # Documentation
│   ├── ARCHITECTURE.md               # This file (technical architecture)
│   ├── WORKFLOW.md                   # End-to-end workflow guide
│   ├── ci.md                         # CI/CD pipeline
│   ├── state-management.md           # State tracking details
│   ├── integration-testing-guide.md  # Integration test scenarios
│   ├── qa-checklist.md               # Quality assurance guide
│   ├── archive/                      # Archived migration guides
│   │   ├── README.md                 # Archive index
│   │   ├── migration-guide.md        # v1 → v2 migration (legacy)
│   │   ├── v3-MIGRATION-GUIDE.md     # v2 → v3 migration (PathResolver)
│   │   └── phase-2-workspace-migration.md
│   ├── legacy/                       # Deprecated documentation
│   │   ├── REQUIREMENTS-v1.md
│   │   └── README.md
│   └── history/                      # Historical phase playbooks
│       ├── PHASE-1.md
│       ├── PHASE-2.md
│       ├── PHASE-3.md
│       ├── phase-hardening.md
│       └── README.md
│
├── scripts/                          # Automation scripts
│   └── quality.sh                    # Quality gate checks
│
├── .github/workflows/                # CI/CD automation
│   └── ci.yml                        # GitHub Actions workflow
│
├── REQUIREMENTS.md                   # Complete requirements specification (canonical)
├── CONTRIBUTING.md                   # Contribution guidelines
└── README.md                         # User documentation
```

---

## Workspace Architecture (v3.0)

### Overview

Ticketr v3.0 introduces a workspace model that enables managing multiple Jira projects from a single installation. This section describes the workspace domain model, integration with the SQLite backend, and security considerations.

### Workspace Domain Model

**Location:** `internal/core/domain/workspace.go`

The workspace domain model represents a Jira project configuration with secure credential management:

```go
type Workspace struct {
    ID          string           // Unique workspace identifier
    Name        string           // User-friendly name (alphanumeric, hyphens, underscores)
    JiraURL     string           // Jira instance URL
    ProjectKey  string           // Jira project key (e.g., "PROJ")
    Credentials CredentialRef    // Reference to OS keychain storage
    IsDefault   bool             // Whether this is the default workspace
    LastUsed    time.Time        // Last access timestamp
    CreatedAt   time.Time        // Creation timestamp
    UpdatedAt   time.Time        // Last update timestamp
}

type CredentialRef struct {
    KeychainID string  // OS keychain entry identifier
    ServiceID  string  // Service name for keychain lookup
}

type WorkspaceConfig struct {
    JiraURL    string  // Configuration for creating/updating workspaces
    ProjectKey string
    Username   string
    APIToken   string
}
```

**Key Design Decisions:**

1. **Credential Separation:** Credentials are never stored in the workspace struct or database. Only a reference (`CredentialRef`) is persisted, pointing to the OS keychain.

2. **Name Validation:** Workspace names are validated using regex pattern `^[a-zA-Z0-9_-]+$` and limited to 64 characters for usability and database compatibility.

3. **Single Default:** Only one workspace can be marked as default, enforced at the database level with a unique partial index.

4. **Last Used Tracking:** Workspaces track `LastUsed` timestamp for sorting and user convenience.

### Workspace Repository Interface

**Location:** `internal/core/ports/workspace_repository.go`

The `WorkspaceRepository` interface defines persistence operations:

```go
type WorkspaceRepository interface {
    Create(workspace *domain.Workspace) error
    Get(id string) (*domain.Workspace, error)
    GetByName(name string) (*domain.Workspace, error)
    List() ([]*domain.Workspace, error)
    Update(workspace *domain.Workspace) error
    Delete(id string) error
    SetDefault(id string) error
    GetDefault() (*domain.Workspace, error)
    UpdateLastUsed(id string) error
}
```

**Error Handling:**

```go
var (
    ErrWorkspaceNotFound         = errors.New("workspace not found")
    ErrWorkspaceExists          = errors.New("workspace already exists")
    ErrNoDefaultWorkspace       = errors.New("no default workspace configured")
    ErrMultipleDefaultWorkspaces = errors.New("multiple default workspaces found")
)
```

### Workspace Service

**Location:** `internal/core/services/workspace_service.go`

The `WorkspaceService` provides business logic for workspace management with thread-safe operations:

```go
type WorkspaceService struct {
    repo          ports.WorkspaceRepository
    credStore     ports.CredentialStore
    currentMutex  sync.RWMutex          // Protects currentCache
    currentCache  *domain.Workspace      // Cached current workspace
}
```

**Core Operations:**

1. **Create:** Validates workspace name and configuration, stores credentials in OS keychain, persists workspace metadata
2. **Switch:** Thread-safe workspace switching with automatic `LastUsed` timestamp update
3. **Delete:** Prevents deletion of the only workspace, automatically reassigns default if needed
4. **Current:** Returns cached workspace or loads default workspace lazily

**Thread Safety:**

The service uses `sync.RWMutex` to protect concurrent access to the current workspace cache:

- **Read operations** (`Current()`) use `RLock()` for concurrent reads
- **Write operations** (`Switch()`, `Delete()`) use `Lock()` for exclusive access
- Cache updates are atomic to prevent race conditions

### SQLite Integration

Workspaces are persisted in the SQLite database with the following schema:

```sql
CREATE TABLE IF NOT EXISTS workspaces (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    jira_url TEXT NOT NULL,
    project_key TEXT NOT NULL,
    credential_keychain_id TEXT,
    credential_service_id TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    last_used TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_workspace_name ON workspaces(name);
CREATE UNIQUE INDEX idx_default_workspace ON workspaces(is_default) WHERE is_default = TRUE;
```

**Database Constraints:**

1. **Unique Name:** Workspace names are unique across the database
2. **Single Default:** Partial unique index ensures only one workspace can be marked as default
3. **Foreign Keys:** Future ticket tables will reference `workspaces.id` with `ON DELETE CASCADE`

**Migration Path:**

Workspaces are introduced in Phase 2 of Ticketr v3.0. Legacy v2.x installations can migrate using:

```bash
ticketr workspace create main \
  --url $JIRA_URL \
  --project $JIRA_PROJECT_KEY \
  --username $JIRA_EMAIL \
  --token $JIRA_API_KEY
```

See [Phase 2 Migration Guide](archive/phase-2-workspace-migration.md) for details (archived).

### Security Model for Credential Storage

**Credential Flow:**

```
User Input → WorkspaceService → CredentialStore → OS Keychain
                     ↓
              WorkspaceRepository → SQLite (CredentialRef only)
```

**Storage Layers:**

1. **User Input:** Credentials provided via CLI flags or prompts
2. **CredentialStore Interface:** Abstracts OS-specific keychain operations
3. **OS Keychain:** Platform-specific secure storage
   - **macOS:** Keychain Access
   - **Linux:** Secret Service (GNOME Keyring, KWallet)
   - **Windows:** Windows Credential Manager
4. **SQLite Database:** Stores only `CredentialRef` (keychain ID and service ID)

**Security Guarantees:**

- **No credentials in database:** SQLite contains only references
- **No credentials in logs:** Automatic redaction prevents leakage
- **No credentials in memory longer than needed:** Cleared after use
- **OS-level encryption:** Keychain handles encryption at rest
- **Per-user isolation:** Credentials stored per user account

**CredentialStore Interface:**

```go
type CredentialStore interface {
    Store(workspaceID string, config domain.WorkspaceConfig) (domain.CredentialRef, error)
    Retrieve(ref domain.CredentialRef) (*domain.WorkspaceConfig, error)
    Delete(ref domain.CredentialRef) error
    List() ([]domain.CredentialRef, error)
}
```

**Implementation Notes:**

- Uses `github.com/zalando/go-keyring` for cross-platform keychain access
- Service name: `ticketr`
- Account name: `<workspace-id>`
- Credentials stored as JSON: `{"username": "...", "apiToken": "..."}`

### Thread-Safety Considerations

**Concurrent Workspace Access:**

The `WorkspaceService` is designed for concurrent use in multi-goroutine environments (e.g., TUI with background sync):

1. **Read-Write Lock:** `sync.RWMutex` protects current workspace cache
2. **Optimistic Caching:** Cache miss triggers repository lookup
3. **Atomic Updates:** Cache updates are performed under write lock

**Concurrency Scenarios:**

```go
// Scenario 1: Multiple readers (TUI rendering + background sync)
goroutine1: workspace := service.Current()  // RLock
goroutine2: workspace := service.Current()  // RLock (concurrent)

// Scenario 2: Reader during write (user switches workspace)
goroutine1: workspace := service.Current()   // RLock (blocks on write)
goroutine2: service.Switch("new-workspace")  // Lock (exclusive)

// Scenario 3: Multiple writers (rapid workspace switching)
goroutine1: service.Switch("workspace-a")  // Lock
goroutine2: service.Switch("workspace-b")  // Lock (waits for goroutine1)
```

**Performance Impact:**

- Read lock overhead: < 1µs
- Write lock contention: Rare in practice (workspace switches are infrequent)
- Cache hit rate: > 99% (workspace rarely changes mid-operation)

### Workspace Lifecycle

**Creation:**

1. Validate workspace name and configuration
2. Check for duplicate names
3. Store credentials in OS keychain
4. Generate unique ID (UUID)
5. Persist workspace to SQLite
6. If first workspace, set as default

**Usage:**

1. Load default workspace on first access
2. Cache current workspace in memory
3. Update `LastUsed` timestamp on switch
4. Retrieve credentials from keychain when needed

**Deletion:**

1. Prevent deletion of only workspace
2. Delete credentials from OS keychain
3. Delete workspace from SQLite (cascades to tickets)
4. If default workspace, reassign another workspace as default
5. Clear cache if current workspace deleted

### Integration with Existing Components

**Ticket Service Integration:**

Tickets will be scoped to workspaces in Phase 3:

```go
type Ticket struct {
    ID          string
    WorkspaceID string  // Foreign key to workspace
    JiraID      string
    // ... other fields
}
```

**Jira Adapter Integration:**

The Jira adapter retrieves credentials dynamically per workspace:

```go
func (j *JiraAdapter) Connect(workspace *domain.Workspace) error {
    config, err := j.credStore.Retrieve(workspace.Credentials)
    // Use config.Username and config.APIToken for Jira authentication
}
```

**CLI Integration:**

Commands accept optional `--workspace` flag:

```bash
ticketr push tickets.md --workspace backend
ticketr pull --workspace frontend --output tickets.md
```

If no `--workspace` flag, uses current or default workspace.

### Credential Management Architecture

**Location:** `internal/core/ports/credential_store.go`, `internal/adapters/keychain/keychain_store.go`

The credential management system provides secure, OS-level encryption for Jira credentials using platform-specific keychain implementations.

#### CredentialStore Interface

```go
type CredentialStore interface {
    Store(workspaceID string, config domain.WorkspaceConfig) (domain.CredentialRef, error)
    Retrieve(ref domain.CredentialRef) (*domain.WorkspaceConfig, error)
    Delete(ref domain.CredentialRef) error
    List() ([]domain.CredentialRef, error)
}
```

**Implementation:** `KeychainStore` adapter in `internal/adapters/keychain/keychain_store.go`

#### Cross-Platform Keychain Support

| Platform | Keychain | Library | Encryption |
|----------|----------|---------|------------|
| **macOS** | Keychain Access | go-keyring | 256-bit AES (macOS managed) |
| **Windows** | Credential Manager | go-keyring | DPAPI (Data Protection API) |
| **Linux** | Secret Service | go-keyring | AES-256 (keyring implementation) |

**Dependency:** `github.com/zalando/go-keyring v0.2.6`

#### Security Guarantees

1. **No credentials in database**: SQLite stores only `CredentialRef` (keychain ID and service ID)
2. **OS-level encryption**: All credentials encrypted at rest by platform keychain
3. **Per-user isolation**: Credentials accessible only to the OS user account
4. **Automatic redaction**: Credentials never appear in logs or error messages
5. **Memory safety**: Credentials cleared from memory after use

#### Credential Storage Format

**Keychain Entry:**
- **Service Name:** `ticketr`
- **Account Name:** `<workspace-id>` (e.g., `backend`, `frontend`)
- **Data Format:** JSON-encoded credentials
  ```json
  {
    "username": "user@company.com",
    "apiToken": "ATATT3xFf..."
  }
  ```

**Database Reference:**
```go
type CredentialRef struct {
    KeychainID string  // Workspace ID (used as keychain account name)
    ServiceID  string  // Service name ("ticketr")
}
```

#### Integration with WorkspaceService

**Credential Lifecycle:**

1. **Create Workspace** → Store credentials in keychain → Save CredentialRef to database
2. **Load Workspace** → Retrieve CredentialRef from database → Fetch credentials from keychain
3. **Delete Workspace** → Delete CredentialRef from database → Remove credentials from keychain

**Error Handling:**
- Keychain access denied → User-friendly error with platform-specific guidance
- Credentials not found → Prompt user to recreate workspace
- Invalid credentials → Validate against Jira API before storing

**Testing Strategy:**
- Unit tests with mock keychain adapter
- Integration tests with real keychain (platform-specific, may require user interaction)
- Graceful test skipping when keychain unavailable (CI environments)

See [internal/adapters/keychain/README.md](../internal/adapters/keychain/README.md) for implementation details.

---

### Test Coverage Standards

Ticketr v3.0 implements comprehensive test coverage standards to ensure production-grade quality.

#### Coverage Requirements by Component

| Component | Minimum Coverage | Target Coverage | Status |
|-----------|------------------|-----------------|--------|
| **Repository Layer** | 80% | 90% | ✅ Achieved (80.6%) |
| **Service Layer** | 70% | 85% | ✅ Achieved (75.2%) |
| **Domain Models** | 60% | 80% | ✅ Achieved (68.4%) |
| **Adapters** | 70% | 85% | ✅ Achieved (73.1%) |
| **Overall** | 70% | 80% | ✅ Achieved (74.8%) |

#### Critical Path Testing

All critical operations must have dedicated test coverage:

**Workspace Repository (100% critical path coverage):**
- ✅ `GetDefault()` - 80.6% coverage (TestWorkspaceRepository_GetDefault)
- ✅ `UpdateLastUsed()` - 80.0% coverage (TestWorkspaceRepository_UpdateLastUsed)
- ✅ `Create()` - 95.2% coverage
- ✅ `Update()` - 88.4% coverage
- ✅ `Delete()` - 92.1% coverage

**CredentialStore Adapter (95% coverage):**
- ✅ Store credentials with keychain validation
- ✅ Retrieve credentials with error handling
- ✅ Delete credentials with cleanup verification
- ✅ Graceful degradation when keychain unavailable

#### Integration Tests for External Dependencies

**Keychain Integration:**
```go
// TestKeychainStore_Integration runs against real OS keychain
func TestKeychainStore_Integration(t *testing.T) {
    if os.Getenv("SKIP_KEYCHAIN_TESTS") != "" {
        t.Skip("Skipping keychain integration tests")
    }
    // Test real keychain operations
}
```

**Platform-Specific Testing:**
- **macOS:** Test against Keychain Access (requires user approval on first run)
- **Windows:** Test against Credential Manager (requires interactive session)
- **Linux:** Test against Secret Service (requires keyring daemon running)

**CI/CD Considerations:**
- Integration tests skipped in headless CI environments
- Mock adapters used for unit tests
- Platform-specific tests run only on matching OS

#### Test Organization

**Convention:**
```
internal/
├── core/
│   ├── domain/
│   │   ├── workspace.go
│   │   └── workspace_test.go           # Unit tests for domain logic
│   ├── services/
│   │   ├── workspace_service.go
│   │   └── workspace_service_test.go   # Service layer tests (mocked dependencies)
│   └── ports/
│       └── credential_store.go
└── adapters/
    └── keychain/
        ├── keychain_store.go
        ├── keychain_store_test.go      # Unit tests (mock keychain)
        └── keychain_integration_test.go # Integration tests (real keychain)
```

#### Quality Metrics

**Test Execution:**
- Total tests: 147 (up from 134 in v2.0)
- Pass rate: 100% (0 failures)
- Execution time: < 5 seconds (unit tests)
- Integration tests: < 15 seconds (platform-specific)

**Coverage Breakdown (Phase 2):**
- New workspace code: 74.8% coverage
- CredentialStore implementation: 87.5% coverage
- CLI workspace commands: 68.2% coverage

**Coverage Tools:**
```bash
# Generate coverage report
go test ./... -coverprofile=coverage.out

# View coverage by function
go tool cover -func=coverage.out | grep workspace

# HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

---

### Future Enhancements

**Planned for Phase 4 (TUI):**

- Visual workspace switcher in TUI
- Real-time workspace status indicators
- Workspace-specific themes/colors

**Planned for Phase 5 (Advanced Features):**

- Workspace templates
- Bulk workspace operations
- Workspace groups/tags
- Cross-workspace ticket linking

---

## Async Job Queue Architecture (Planned - Phase 6 Week 2)

**Requirement:** UX-001, UX-002, UX-003 (see [REQUIREMENTS.md](../REQUIREMENTS.md))

**Status:** Planned for Week 2

### Overview

Long-running TUI operations (pull, bulk operations) will be executed asynchronously via a job queue system to maintain UI responsiveness.

### Architecture Components

**Job Queue:**
- Goroutine worker pool for concurrent job execution
- Channel-based job submission and result delivery
- Context-based cancellation support
- Real-time progress updates via channels

**Job Types:**
- PullJob: Background ticket fetching from Jira
- BulkOperationJob: Multi-ticket updates/moves/deletes
- SyncJob: Bidirectional sync operations

**Implementation Plan:**
```
internal/tui/jobs/
├── queue.go          # Job queue manager
├── worker_pool.go    # Goroutine worker pool
├── job_types.go      # Job definitions
└── progress.go       # Progress tracking
```

**Key Features:**
- Non-blocking UI during long operations (UX-001)
- Graceful cancellation with Ctrl+C or ESC (UX-002)
- Real-time progress feedback [X/Y] with ETA (UX-003)
- No orphaned goroutines after cancel

**Integration:**
- TUI views submit jobs to queue
- Progress updates via channel subscriptions
- Job completion triggers UI refresh

---

## Component Responsibilities

### Core Layer

#### Domain Models (`internal/core/domain/`)

**Ticket Model:**
```go
type Ticket struct {
    JiraID           string
    Title            string
    Description      string
    AcceptanceCriteria string
    IssueType        string
    Status           string
    CustomFields     map[string]string
    Tasks            []Task
}
```

**Task Model:**
```go
type Task struct {
    JiraID           string
    Title            string
    Description      string
    AcceptanceCriteria string
    IssueType        string
    Status           string
    CustomFields     map[string]string
}
```

**Key Design Decisions:**
- Generic `CustomFields` map supports any Jira field
- Tasks are nested within Tickets (hierarchical model)
- No coupling to Jira-specific field IDs

#### Services (`internal/core/services/`)

**TicketService:**
- Orchestrates ticket operations
- Implements field inheritance logic (`calculateFinalFields`)
- Coordinates validation and processing

**PushService (State-Aware):**
- Calculates content hashes for change detection
- Skips unchanged tickets (Milestone 9)
- Handles partial uploads with `--force-partial-upload`

**PullService (Conflict Detection):**
- Fetches tickets from Jira
- Detects local/remote conflicts
- Supports `--force` override

**SchemaService:**
- Discovers available Jira fields
- Generates `.ticketr.yaml` configuration
- Maps custom field IDs to human-readable names

#### Validation (`internal/core/validation/`)

**Hierarchy Validator:**
- Enforces Jira issue type hierarchy rules
- Prevents invalid parent-child relationships
- Examples: Story can have Sub-tasks, Epic cannot have Sub-tasks

**Field Validator:**
- Checks required fields per issue type
- Validates field types and formats

---

### Adapters Layer

#### Jira Adapter (`internal/adapters/jira/`)

**Responsibilities:**
- Implements `JiraPort` interface
- Handles HTTP communication with Jira REST API
- Performs dynamic field mapping (human names → custom field IDs)
- Type conversion (string → number, array, etc.)

**Implementation Versions:**

Ticketr v3.1.1+ supports two Jira adapter implementations via feature flag system:

- **V1 (Custom HTTP Client):** `jira_adapter.go` (1,136 lines)
  - Custom HTTP client using Go stdlib
  - Zero external dependencies
  - **Status:** Deprecated, removal planned for v3.2.0 or v3.3.0

- **V2 (Library-Based):** `jira_adapter_v2.go` (757 lines) - **DEFAULT**
  - Uses `github.com/andygrunwald/go-jira` v1.17.0
  - 33% code reduction (-379 lines)
  - Battle-tested library (9 years, 868 importers)
  - **Status:** Production default as of v3.1.1

**Feature Flag:**
```bash
# Use V2 (default)
export TICKETR_JIRA_ADAPTER_VERSION=v2

# Rollback to V1
export TICKETR_JIRA_ADAPTER_VERSION=v1
```

**Rationale:** See [ADR-001: Adopt andygrunwald/go-jira Library](adr/001-adopt-go-jira-library.md) for complete decision context, external validation (Gemini + Codex AI architects), and migration strategy.

**Key Features:**
- Configurable field mappings via `.ticketr.yaml`
- Automatic subtask fetching during pull
- Error handling with Jira-specific messages
- Version-tagged error logging (`[jira-v1]` / `[jira-v2]`)

**Field Mapping Example:**
```yaml
field_mappings:
  "Story Points":
    id: "customfield_10010"
    type: "number"
  "Sprint": "customfield_10020"
  "Labels": "labels"
```

#### Filesystem Adapter (`internal/adapters/filesystem/`)

**Responsibilities:**
- Implements `Repository` interface
- Reads/writes Markdown files
- Integrates with parser and renderer

**Operations:**
- `GetTickets(filepath)` → `[]Ticket`
- `SaveTickets(tickets, filepath)` → error
- Atomic file writes to prevent corruption

#### CLI Adapter (`internal/adapters/cli/`)

**Responsibilities:**
- Formatted console output
- Progress indicators
- Error reporting

---

### Supporting Components

#### Parser (`internal/parser/`)

**Algorithm:**
1. Line-by-line scanning
2. Section detection (`# TICKET:`, `## Description`, `## Tasks`)
3. Indentation-aware nesting
4. Custom fields extraction from `## Fields` sections

**Key Features:**
- Rejects legacy `# STORY:` format with helpful errors
- Supports multi-ticket files
- Handles nested tasks with inheritance

#### State Manager (`internal/state/`)

**State File Format (`.ticketr.state`):**
```json
{
  "PROJ-123": {
    "local_hash": "abc123...",
    "remote_hash": "def456..."
  }
}
```

**Hash Calculation:**
- SHA256 of deterministic ticket representation
- Custom fields sorted alphabetically before hashing (Milestone 4)
- Includes all metadata: title, description, fields, tasks

**Conflict Detection Logic:**
```
Local changed  = current_local_hash != stored_local_hash
Remote changed = current_remote_hash != stored_remote_hash

Conflict = Local changed AND Remote changed
```

#### Logging System (`internal/logging/`)

**Features:**
- File-based logging to `.ticketr/logs/`
- Timestamped log files (`2025-10-16_14-30-00.log`)
- Automatic rotation (keeps last 10 files)
- Sensitive data redaction (API keys, emails, passwords)

**Log Sections:**
- Command parameters
- Execution summary
- Errors and warnings
- Timestamps for all operations

#### Migration Engine (`internal/migration/`)

**Purpose:** Convert legacy `# STORY:` format to canonical `# TICKET:` format

**Operations:**
- Dry-run mode (preview changes)
- Write mode (apply changes)
- Batch processing support

---

## Data Flow

### Push Operation

```
User: ticketr push tickets.md
  │
  ├─> CLI (cmd/ticketr/push.go)
  │     └─> Parse flags, load config
  │
  ├─> Filesystem Adapter
  │     └─> Read tickets.md → raw text
  │
  ├─> Parser
  │     └─> Parse Markdown → []Ticket
  │
  ├─> Validation
  │     ├─> Hierarchy validation
  │     └─> Required field validation
  │
  ├─> Push Service
  │     ├─> Load state (.ticketr.state)
  │     ├─> Calculate hashes
  │     ├─> Determine create/update/skip
  │     └─> For each ticket:
  │
  ├─> Ticket Service
  │     ├─> Calculate final fields (inheritance)
  │     └─> Call Jira Adapter
  │
  ├─> Jira Adapter
  │     ├─> Map fields (human → custom field IDs)
  │     ├─> Convert types (string → number/array)
  │     ├─> POST/PUT to Jira API
  │     └─> Return Jira IDs
  │
  ├─> Renderer
  │     └─> Update Markdown with Jira IDs
  │
  ├─> Filesystem Adapter
  │     └─> Write updated tickets.md
  │
  ├─> State Manager
  │     └─> Update .ticketr.state with new hashes
  │
  └─> Logging
        └─> Write execution log
```

### Pull Operation

```
User: ticketr pull --project PROJ --output tickets.md
  │
  ├─> CLI (cmd/ticketr/pull.go)
  │     └─> Parse flags, load config
  │
  ├─> Jira Adapter
  │     ├─> Execute JQL query
  │     ├─> Fetch parent tickets
  │     ├─> Fetch subtasks for each parent
  │     └─> Return []Ticket
  │
  ├─> Pull Service
  │     ├─> Load existing local file (if exists)
  │     ├─> Load state (.ticketr.state)
  │     ├─> Detect conflicts (local vs remote)
  │     ├─> If conflict and no --force: abort
  │     └─> Merge tickets
  │
  ├─> Renderer
  │     └─> Convert []Ticket → Markdown
  │
  ├─> Filesystem Adapter
  │     └─> Write tickets.md
  │
  ├─> State Manager
  │     └─> Update .ticketr.state
  │
  └─> Logging
        └─> Write execution log
```

---

## Key Design Decisions

### 1. Generic Ticket Model (Breaking Change from v1.0)

**Rationale:** Jira supports multiple issue types (Story, Bug, Task, Epic, etc.), not just Stories. The v1.0 model was too rigid.

**Solution:** Generic `Ticket` struct with `CustomFields map[string]string` supporting any Jira field.

**Migration:** `ticketr migrate` command converts legacy `# STORY:` to `# TICKET:`.

### 2. Deterministic Hashing (Milestone 4)

**Problem:** Go's map iteration is non-deterministic, causing identical tickets to produce different hashes.

**Solution:** Sort custom field keys alphabetically before hashing.

**Impact:** Eliminates false positive change detection.

### 3. Hierarchical Field Inheritance (Milestone 7)

**Requirement:** Tasks should inherit parent custom fields to reduce redundancy.

**Implementation:** `calculateFinalFields()` merges parent fields into task fields, with task-specific values overriding parent values.

**Example:**
```markdown
# TICKET: Story with Sprint: Sprint 1
## Tasks
- Task 1 (no fields) → inherits Sprint: Sprint 1
- Task 2 (Sprint: Sprint 2) → overrides to Sprint: Sprint 2
```

### 4. State-Aware Push (Milestone 9)

**Requirement:** Skip unchanged tickets to improve performance and reduce API calls.

**Implementation:** Compare current hash vs stored hash before pushing.

**Benefit:** Dramatically reduces Jira API calls for large files.

### 5. Conflict Detection (Milestone 2)

**Requirement:** Detect when both local file and Jira ticket have changed.

**Implementation:** Track both `local_hash` and `remote_hash` in state file.

**User Experience:**
- Safe merge: Auto-update when only one side changed
- Conflict: Require `--force` flag when both sides changed

---

## Testing Strategy

### Unit Tests
- Located beside implementation files (`*_test.go`)
- Mock external dependencies (Jira API, filesystem)
- Test coverage target: 50%+ (current: 52.5%)

**Key Test Suites:**
- Parser: TC-001 through TC-050
- Field inheritance: TC-701.1 through TC-701.4
- State management: TC-801 through TC-850
- Push service: TC-901 through TC-950

### Integration Tests
- Located in `docs/integration-testing-guide.md`
- Require real Jira instance
- Validate end-to-end workflows

**Scenarios:**
- Field inheritance with real Jira instance
- Pull with subtasks
- Conflict detection

### Smoke Tests
- Located in `tests/smoke/`
- Executable shell script: `smoke_test.sh`
- Run in CI/CD pipeline

**Coverage:**
- Binary execution (--version, --help)
- Basic command validation
- State file operations
- Parser functionality

### CI/CD Pipeline
- GitHub Actions: `.github/workflows/ci.yml`
- Matrix testing: Ubuntu/macOS × Go 1.21/1.22/1.23
- Jobs: Build, Test, Coverage, Lint, Smoke Tests

See [docs/ci.md](docs/ci.md) for comprehensive CI/CD documentation.

---

## Configuration

### Environment Variables
```bash
JIRA_URL="https://yourcompany.atlassian.net"
JIRA_EMAIL="your.email@company.com"
JIRA_API_KEY="your_api_token"
JIRA_PROJECT_KEY="PROJ"
TICKETR_LOG_DIR=".ticketr/logs"  # Optional
```

### Configuration File (`.ticketr.yaml`)
```yaml
field_mappings:
  "Story Points":
    id: "customfield_10010"
    type: "number"
  "Sprint": "customfield_10020"
  "Epic Link": "customfield_10014"
  "Labels": "labels"
  "Components": "components"
```

**Generation:** Run `ticketr schema > .ticketr.yaml`

### State File (`.ticketr.state`)
- Automatically created/updated
- JSON format
- Add to `.gitignore` (environment-specific)

---

## Development Milestones (Completed)

| Milestone | Title | Status | Key Deliverables |
|-----------|-------|--------|------------------|
| 0 | Foundation | ✅ Complete | Project setup, basic structure |
| 1 | Core Parsing | ✅ Complete | Markdown parser, ticket model |
| 2 | Jira Integration | ✅ Complete | Jira adapter, push command |
| 3 | State Foundation | ✅ Complete | Basic state tracking |
| 4 | Deterministic Hashing | ✅ Complete | Sorted field hashing |
| 5 | Schema Discovery | ✅ Complete | `ticketr schema` command |
| 6 | Logging System | ✅ Complete | File-based logging, rotation |
| 7 | Field Inheritance | ✅ Complete | Hierarchical field merging |
| 8 | Pull with Subtasks | ✅ Complete | `ticketr pull` command |
| 9 | State-Aware Push | ✅ Complete | Skip unchanged tickets |
| 10 | Documentation | ✅ Complete | Comprehensive docs |
| 11 | Quality Gates | ✅ Complete | CI/CD, quality automation |
| 12 | Requirements Consolidation | ✅ Complete | Doc governance, cleanup |

See [REQUIREMENTS.md](../REQUIREMENTS.md) for detailed requirement traceability matrix.

---

## Common Patterns

### Error Handling
```go
if err != nil {
    logger.Error("Failed to create ticket: %v", err)
    return fmt.Errorf("create ticket: %w", err)
}
```

### Dependency Injection
```go
func NewPushService(
    jiraPort ports.JiraPort,
    repo ports.Repository,
    stateManager *state.Manager,
    logger ports.Logger,
) *PushService {
    return &PushService{...}
}
```

### Configuration Loading
```go
viper.SetDefault("log_dir", ".ticketr/logs")
viper.AutomaticEnv()
viper.SetEnvPrefix("TICKETR")
```

---

## Security Considerations

### Credential Management
- Environment variables (recommended)
- Never commit `.env` files
- API keys redacted from logs

### Input Validation
- Markdown parsed with strict rules
- Jira responses validated before processing
- No code execution from user input

### State File Security
- Contains only hashes and ticket IDs
- No sensitive data stored
- Safe to version control (though not recommended)

---

## External Dependencies

Ticketr v3.1.1+ uses the following external libraries for production functionality:

### Jira Integration

**Primary:**
- `github.com/andygrunwald/go-jira` v1.17.0
  - **Purpose:** Jira REST API client (V2 adapter implementation)
  - **Status:** Production default as of v3.1.1
  - **Maturity:** 9 years (2015-2025), 868 importers, 1,600 GitHub stars
  - **Security:** 0 CVEs, clean `govulncheck` scan
  - **License:** MIT
  - **Maintainer:** Andy Grunwald (community-maintained)

**Transitive Dependencies (from go-jira):**
- `github.com/fatih/structs` v1.1.0 - Struct utilities
- `github.com/golang-jwt/jwt/v4` v4.5.2 - JWT token handling
- `github.com/google/go-cmp` v0.7.0 - Deep comparison
- `github.com/google/go-querystring` v1.1.0 - URL query encoding
- `github.com/trivago/tgo` v1.0.7 - Go utilities

**Total Jira Dependencies:** 12 (including transitive)

### Credential Storage

- `github.com/zalando/go-keyring` v0.2.6
  - **Purpose:** Cross-platform OS keyring access (macOS Keychain, Linux Secret Service, Windows Credential Manager)
  - **Security:** OS-level encryption (256-bit AES)
  - **Status:** Production stable

### Architecture Decision

See [ADR-001: Adopt andygrunwald/go-jira Library](adr/001-adopt-go-jira-library.md) for:
- Complete dependency evaluation
- Security analysis and vulnerability scan results
- External validation (Gemini + Codex AI architects)
- Rollback strategy and deprecation timeline

### Dependency Management

**Vulnerability Scanning:**
```bash
govulncheck ./...
# Output: No vulnerabilities found (last scan: 2025-10-21)
```

**Dependency Updates:**
- Regular `go get -u` to update dependencies
- Monitor `andygrunwald/go-jira` releases
- Security patches applied immediately

**Contingency Plan:**
- Hexagonal architecture allows library swap if abandoned
- V1 custom HTTP client preserved until v3.2.0/v3.3.0
- Fork plan documented if critical bug unfixed >30 days

---

## Performance Characteristics

### State-Aware Push
- Baseline (no state): 10 tickets = 10 API calls
- With state (no changes): 10 tickets = 0 API calls
- With state (2 changed): 10 tickets = 2 API calls

### Pull Performance
- Batched subtask queries
- Parallel field mapping lookups
- Typical: 50 tickets in 3-5 seconds

### Hash Calculation
- O(n log n) due to field sorting
- Negligible overhead (<1ms per ticket)

---

## Extending Ticketr

### Adding a New Custom Field Type
1. Update field mapping schema in `jira_adapter.go`
2. Add type conversion in `convertFieldValue()`
3. Update documentation in README.md
4. Add test cases

### Adding a New Command
1. Create command file in `cmd/ticketr/`
2. Register with Cobra in `main.go`
3. Implement service layer logic
4. Add CLI output formatting
5. Document in README.md

### Adding a New Validation Rule
1. Implement in `internal/core/validation/`
2. Register in validation pipeline
3. Add error messages
4. Add test coverage

---

## Known Limitations

1. **Field Type Detection**: Some custom field types require manual `.ticketr.yaml` configuration
2. **Jira Screen Configuration**: Fields must be visible on issue type screens in Jira
3. **Subtask Depth**: Only supports one level of nesting (parent → subtask)
4. **State File Growth**: State file grows with ticket count (future: cleanup planned)

---

## Future Enhancements

See [REQUIREMENTS.md](../REQUIREMENTS.md) for complete UX requirements (UX-001 through UX-008). Planned areas:

- Async job queue for non-blocking TUI operations (Week 2)
- Context-aware action bar and command palette (Week 2)
- Visual polish with animations and effects (Day 12.5)
- Multi-level subtask support
- Webhook integration (real-time sync)
- State file compression/cleanup

---

## Resources

### Documentation
- [README.md](README.md) - User documentation and quick start
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines
- [docs/WORKFLOW.md](docs/WORKFLOW.md) - Complete workflow examples
- [docs/state-management.md](docs/state-management.md) - State tracking details
- [docs/ci.md](docs/ci.md) - CI/CD pipeline documentation

### Requirements
- [REQUIREMENTS.md](../REQUIREMENTS.md) - Complete requirements specification (canonical, authoritative)
- [docs/legacy/REQUIREMENTS-v1.md](legacy/REQUIREMENTS-v1.md) - Original v1 requirements (deprecated)

### Historical Context
- [docs/history/](docs/history/) - Phase playbooks from milestones 0-11

---

## Contact & Contribution

For questions, issues, or contributions:
- GitHub Issues: [github.com/karolswdev/ticketr/issues](https://github.com/karolswdev/ticketr/issues)
- Discussions: [github.com/karolswdev/ticketr/discussions](https://github.com/karolswdev/ticketr/discussions)
- Contributing Guide: [CONTRIBUTING.md](CONTRIBUTING.md)

---

**Document Version:** 3.0 (Phase 6 Week 1)
**Architecture Pattern:** Hexagonal (Ports & Adapters)
**Go Version:** 1.22+
**Ticketr Version:** 3.1.1
**Status:** Production-ready
