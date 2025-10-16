# Ticketr Architecture Overview

**Last Updated:** October 16, 2025 (Milestone 12)
**Status:** Production-ready, actively maintained

## Executive Summary

Ticketr is a command-line tool that bridges Markdown files and Jira, enabling bidirectional synchronization of tickets and tasks. The system is built using Hexagonal Architecture (Ports and Adapters) with a focus on maintainability, testability, and extensibility.

**Key Capabilities:**
- Parse Markdown files into structured ticket data
- Push tickets/tasks to Jira with hierarchical field inheritance
- Pull tickets/tasks from Jira with automatic subtask fetching
- State-aware operations (skip unchanged tickets)
- Conflict detection and resolution
- Dynamic field mapping and schema discovery
- Comprehensive logging and validation

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
│   ├── WORKFLOW.md                   # End-to-end workflow guide
│   ├── ci.md                         # CI/CD pipeline
│   ├── state-management.md           # State tracking details
│   ├── migration-guide.md            # v1 → v2 migration
│   ├── integration-testing-guide.md  # Integration test scenarios
│   ├── qa-checklist.md               # Quality assurance guide
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
├── REQUIREMENTS-v2.md                # Current requirements (canonical)
├── ROADMAP.md                        # Project roadmap and milestones
├── CONTRIBUTING.md                   # Contribution guidelines
├── README.md                         # User documentation
└── ARCHITECTURE.md                   # This file
```

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

**Key Features:**
- Configurable field mappings via `.ticketr.yaml`
- Automatic subtask fetching during pull
- Error handling with Jira-specific messages

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

See [ROADMAP.md](ROADMAP.md) for detailed milestone tracking.

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

See [ROADMAP.md](ROADMAP.md) for planned features. Potential areas:

- Multi-level subtask support
- Advanced conflict resolution (3-way merge)
- Batch operations API (reduce HTTP overhead)
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
- [REQUIREMENTS-v2.md](REQUIREMENTS-v2.md) - Current requirements (canonical)
- [docs/legacy/REQUIREMENTS-v1.md](docs/legacy/REQUIREMENTS-v1.md) - Original v1 requirements (deprecated)

### Historical Context
- [docs/history/](docs/history/) - Phase playbooks from milestones 0-11

---

## Contact & Contribution

For questions, issues, or contributions:
- GitHub Issues: [github.com/karolswdev/ticketr/issues](https://github.com/karolswdev/ticketr/issues)
- Discussions: [github.com/karolswdev/ticketr/discussions](https://github.com/karolswdev/ticketr/discussions)
- Contributing Guide: [CONTRIBUTING.md](CONTRIBUTING.md)

---

**Document Version:** 2.0
**Architecture Pattern:** Hexagonal (Ports & Adapters)
**Go Version:** 1.22+
**Status:** Production-ready
