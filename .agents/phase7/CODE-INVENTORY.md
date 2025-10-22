# Phase 7 Code Inventory - Files for Analysis

**Date:** 2025-10-21
**Purpose:** Inventory of all code files Steward must analyze
**Phase:** 7 Day 1 Current State Analysis

---

## Critical Files for SOLID/DDD/Hexagonal Analysis

### 1. Jira Adapter (Primary Analysis Target)

**File:** `/home/karol/dev/private/ticktr/internal/adapters/jira/jira_adapter.go`
**Line Count:** 1,137 lines (as of Phase 6.5)
**Purpose:** Jira REST API adapter implementation

**Known Issues (from Phase 6.5):**
- God object anti-pattern (HTTP, auth, JSON, field mapping, domain conversion all in one)
- No dependency injection (embedded HTTP client)
- Synchronous operations (no context cancellation initially)
- Field mapping complexity (runtime config, no type safety)
- No error taxonomy (generic fmt.Errorf everywhere)

**Analysis Focus:**
- Single Responsibility violations
- Open/Closed violations (hard to extend)
- Component decomposition opportunities
- Testability issues

**Related Test Files:**
- `/home/karol/dev/private/ticktr/internal/adapters/jira/jira_adapter_test.go`
- `/home/karol/dev/private/ticktr/internal/adapters/jira/jira_adapter_dynamic_test.go`
- `/home/karol/dev/private/ticktr/internal/adapters/jira/jira_adapter_error_test.go`

---

### 2. Domain Models (DDD Analysis Target)

**File:** `/home/karol/dev/private/ticktr/internal/core/domain/models.go`
**Purpose:** Core domain entities and value objects

**Known Issues (from Phase 7 spec):**
- Leaky abstractions (JiraID in domain)
- Stringly-typed custom fields (map[string]string)
- Missing value objects (Priority, Status, IssueType are strings or primitives)
- Anemic domain model (no business logic, just data containers)

**Analysis Focus:**
- Ubiquitous language consistency
- Aggregate root identification
- Value object opportunities
- Domain invariants and validation
- External dependencies in domain

**Related Files:**
- `/home/karol/dev/private/ticktr/internal/core/domain/workspace.go`
- `/home/karol/dev/private/ticktr/internal/core/domain/bulk_operation.go`
- `/home/karol/dev/private/ticktr/internal/core/domain/credential_profile.go`
- `/home/karol/dev/private/ticktr/internal/core/domain/jql_alias.go`

---

### 3. Port Interfaces (Hexagonal Analysis Target)

**File:** `/home/karol/dev/private/ticktr/internal/core/ports/jira_port.go`
**Purpose:** Port interface for Jira integration

**Known Issues (from Phase 7 spec):**
- Jira-specific naming (prevents polymorphism)
- Inconsistent abstraction levels (high-level CreateTicket, low-level GetIssueTypeFields)
- Return types lose type safety (map[string]interface{})
- Missing batch operations, retry hooks, caching strategy
- Exposes JQL (Jira Query Language) instead of system-agnostic query

**Analysis Focus:**
- Port ownership (domain needs vs adapter needs)
- Abstraction consistency
- System-agnostic design
- Interface segregation
- Dependency inversion

**Related Port Files:**
- `/home/karol/dev/private/ticktr/internal/core/ports/repository.go`
- `/home/karol/dev/private/ticktr/internal/core/ports/extended_repository.go`
- `/home/karol/dev/private/ticktr/internal/core/ports/workspace_repository.go`
- `/home/karol/dev/private/ticktr/internal/core/ports/bulk_operation_service.go`
- `/home/karol/dev/private/ticktr/internal/core/ports/sync_strategy.go`

---

### 4. Service Layer (Use Case Analysis)

**Files to Analyze:**
```
internal/core/services/
├── push_service.go        (Push use case)
├── pull_service.go        (Pull use case)
├── sync_service.go        (Sync use case)
├── bulk_operation_service.go (Batch operations)
└── workspace_service.go   (Workspace management)
```

**Known Issues (from Phase 6.5):**
- Direct adapter coupling (hard to swap Jira for GitHub)
- No transaction semantics or rollback
- Three repositories (file, database, Jira) with unclear primary source
- Error aggregation loses partial state information
- Synchronous operations in async job queue

**Analysis Focus:**
- Use case driven vs CRUD
- Transaction boundaries
- Repository pattern usage
- Conflict resolution strategy
- Progress reporting chain

---

## Supporting Files for Context

### Configuration & Field Mapping

**Files:**
- `/home/karol/dev/private/ticktr/internal/adapters/jira/config.go` (if exists)
- Field mapping configuration (runtime, external files?)

**Analysis:** How is field mapping configured? Compile-time vs runtime trade-offs?

---

### Repository Implementations

**Files:**
```
internal/adapters/repository/
├── file_repository.go      (Filesystem persistence)
├── database_repository.go  (SQLite persistence)
└── jira_repository.go      (Jira as remote store)
```

**Analysis:** Why three repositories? Which is source of truth? How are conflicts handled?

---

### TUI Integration (For Understanding Service Usage)

**Files:**
```
internal/adapters/tui/
├── app.go                  (Main TUI application)
├── handlers.go             (Event handlers - calls services)
└── views/
    └── workspace_modal.go  (Workspace creation flow)
```

**Analysis:** How does TUI consume services? What abstractions does it depend on?

---

## Analysis Methodology for Steward

### Step 1: Read Phase 6.5 Context
- Understand what failed and why
- Identify architectural debt patterns
- Note human's concerns

### Step 2: Analyze Jira Adapter (1,137 lines)

**SOLID Analysis:**
```
Read: internal/adapters/jira/jira_adapter.go
Identify:
- Single Responsibility violations (multiple concerns in one class?)
- Open/Closed violations (hard to extend? Modification required for new features?)
- Liskov Substitution violations (can't swap adapters?)
- Interface Segregation violations (clients depend on methods they don't use?)
- Dependency Inversion violations (depends on concrete types?)

Document:
- Each violation with specific line ranges
- Code example showing the violation
- Explain WHY it's a violation
- Suggest correct approach
```

**Example Output:**
```markdown
### SOLID Violation: Single Responsibility (Lines 150-350)

**Code Example:**
```go
// Lines 150-350 in jira_adapter.go
func (j *JiraAdapter) CreateTicket(ticket Ticket) (string, error) {
    // HTTP client setup (concern #1)
    client := &http.Client{Timeout: 30 * time.Second}

    // JSON serialization (concern #2)
    body, err := json.Marshal(ticketToJiraPayload(ticket))

    // HTTP request (concern #3)
    req, err := http.NewRequest("POST", j.baseURL+"/issue", bytes.NewBuffer(body))

    // Authentication (concern #4)
    req.Header.Set("Authorization", "Bearer "+j.token)

    // Domain conversion (concern #5)
    jiraIssue := convertToJiraFormat(ticket)

    // ... etc
}
```

**Why This Violates SRP:**
CreateTicket method has 5+ reasons to change:
1. HTTP library changes
2. JSON serialization changes
3. Authentication method changes
4. Domain model changes
5. Jira API format changes

**Correct Approach:**
Separate concerns into focused components:
- JiraTransport (HTTP, auth, retries)
- JiraSerializer (JSON ↔ Domain)
- JiraAdapter (orchestrates components)

Each component has ONE reason to change.
```

### Step 3: Analyze Domain Models

**DDD Analysis:**
```
Read: internal/core/domain/models.go
Identify:
- Ubiquitous language violations (inconsistent terminology?)
- Missing aggregates (no aggregate root enforcement?)
- Missing value objects (primitives/strings where value objects should be?)
- Leaky abstractions (external system concepts in domain?)
- Anemic model (no business logic, just getters/setters?)

Document:
- Each violation with code examples
- Explain correct DDD pattern
- Show what value objects should exist
- Define ubiquitous language terms
```

### Step 4: Analyze Port Interfaces

**Hexagonal Analysis:**
```
Read: internal/core/ports/jira_port.go
Identify:
- Port ownership violations (adapter-shaped ports vs domain-shaped?)
- Abstraction inconsistencies (mixing levels?)
- External dependencies in domain (imports from adapters?)
- Jira-specific concepts in ports (should be system-agnostic)

Document:
- Each violation with code examples
- Show why port is adapter-shaped
- Design system-agnostic alternative
- Explain dependency direction
```

---

## Deliverable Structure for Current State Analysis

```markdown
# Current State Analysis - Phase 7

## Executive Summary (1 page)
- Overall architectural health: POOR / FAIR / GOOD
- Critical violations count: X SOLID, Y DDD, Z Hexagonal
- Root cause summary
- Urgency of redesign

## SOLID Principles Analysis (3-4 pages)

### Single Responsibility Violations
[Code examples from jira_adapter.go, services, etc.]

### Open/Closed Violations
[Examples of hard-coded Jira specifics preventing extension]

### Liskov Substitution Violations
[Examples of non-swappable implementations]

### Interface Segregation Violations
[Examples of fat interfaces]

### Dependency Inversion Violations
[Examples of depending on concrete types]

## DDD Patterns Analysis (3-4 pages)

### Ubiquitous Language Issues
[Terminology inconsistencies]

### Missing Aggregates
[Where aggregate roots should be]

### Missing Value Objects
[Where value objects should replace primitives]

### Leaky Abstractions
[Domain depending on external systems]

### Anemic Domain Model
[Lack of business logic in domain]

## Hexagonal Architecture Analysis (3-4 pages)

### Domain Dependency Violations
[External dependencies in domain layer]

### Port Ownership Issues
[Ports shaped by adapter instead of domain]

### Adapter Coupling
[Examples of adapter-specific concepts leaking]

### Inside/Outside Boundary Violations
[Dependency direction violations]

## Summary of Root Causes (1 page)
- Why did these violations occur?
- What are the consequences?
- What must change?

## Recommendations for Phase 7 Design (1 page)
- Top 3-5 architectural changes needed
- Principles to enforce in redesign
- Success criteria for new design
```

**Total:** 10-15 pages with code examples

---

## Tools for Analysis

**Reading Code:**
- `Read` tool: Read specific files
- `Grep` tool: Search for patterns across codebase
- `Glob` tool: Find files by pattern

**Examples:**
```
# Find all map[string]interface{} usage (stringly-typed)
Grep: pattern="map\[string\]interface\{\}", glob="**/*.go"

# Find all direct imports of jira adapter in domain
Grep: pattern="internal/adapters/jira", path="internal/core/domain"

# Find all TODO comments related to field mapping
Grep: pattern="TODO.*field", glob="**/*.go"
```

---

## Success Criteria for Current State Analysis

**You will have succeeded when:**
- [ ] Every SOLID violation has code example with line numbers
- [ ] Every DDD violation has code example showing the issue
- [ ] Every Hexagonal violation shows dependency direction problem
- [ ] Each violation explains WHY it's a problem
- [ ] Each violation suggests correct approach
- [ ] Document is 10-15 pages
- [ ] Director approves quality and depth
- [ ] Ready to inform external validation research (what to look for)

---

**Inventory Compiled By:** Director Agent
**Date:** 2025-10-21
**For:** Steward Agent Day 1 Analysis
**Phase:** 7 Current State Analysis
