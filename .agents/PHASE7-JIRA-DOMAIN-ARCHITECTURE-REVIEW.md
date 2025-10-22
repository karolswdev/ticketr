# Phase 7: Jira Domain Architecture Review

**Status:** INITIATED
**Priority:** CRITICAL - Architectural Foundation
**Duration:** 3-5 days estimated
**Trigger:** Human mandate for comprehensive domain model review
**Branch:** feature/jira-domain-redesign

---

## CRITICAL MANDATE FROM HUMAN

> "I DEMAND A THOROUGH REVIEW PHASE OF THE CURRENT JIRA SERVICE, JIRA ADAPTER. I WANT YOU TO QUESTION EVERYTHING IN THAT PHASE, THE DESIGN, ARE WE ROBUST ENOUGH? ARE WE PROVIDING THE RIGHT LEVEL OF ABSTRACTION AND DESIGN FOR FLEXIBILITY?"

> "ARE WE CHECKING THIS WITH HEADLESS VERSIONS OF GEMINI-CLI AND CODEX TO FIND OUT OTHER VIEWPOINTS ON THIS? THE DOMAIN MUST BE SOLID."

**Translation:** The human has lost confidence in our Jira integration architecture. We must restore it through rigorous analysis, external validation, and principled redesign.

---

## Executive Summary

### The Problem

After 13 milestones of implementation and reaching Phase 6, the human questions whether our Jira service layer is fundamentally sound. This is not about bugs - it's about architectural integrity.

**Key Concerns:**
1. Is the abstraction level correct?
2. Are we robust enough for production use?
3. Do we have the right flexibility for future requirements?
4. Is the domain model solid and principled?
5. Have we applied industry best practices?

### Why This Review is Critical

**Current State:**
- JiraAdapter: 1,137 lines mixing HTTP, JSON, domain conversion, field mapping
- JiraPort: 42 lines defining 8 methods with varying levels of abstraction
- Domain Models: 21 lines (Ticket, Task) - suspiciously minimal
- Service Layer: Push/Pull services directly couple to adapter implementation
- No clear separation of concerns
- Field mapping hardcoded in adapter
- Error handling inconsistent
- Context support retrofitted, not designed-in

**Risk if We Don't Fix:**
- Technical debt compounds with each feature
- Jira API changes break entire system
- Cannot support multiple issue trackers (GitHub, Linear, etc.)
- Testing becomes harder over time
- Onboarding new developers is painful
- User trust erodes with each integration issue

**Opportunity:**
- Establish architectural patterns for all future adapters
- Create a reference implementation for domain-driven design
- Build confidence through external validation
- Set quality bar for rest of codebase

### What Success Looks Like

**By end of Phase 7:**
1. Domain model validated by DDD principles
2. Clear separation: Domain → Ports → Adapters → Services
3. External validation from industry patterns (Atlassian SDK, GitHub API clients, Linear API)
4. Comprehensive error taxonomy and handling
5. Flexibility to add new issue trackers without touching domain
6. Test coverage that gives confidence
7. Documentation that educates future contributors
8. Human approval of architecture before implementation

---

## Current State Analysis

### Domain Layer (/internal/core/domain/models.go)

**Ticket struct (21 lines total):**
```go
type Ticket struct {
    Title              string
    Description        string
    CustomFields       map[string]string  // ← CONCERN: Stringly-typed
    AcceptanceCriteria []string
    JiraID             string             // ← CONCERN: Jira leaking into domain
    Tasks              []Task
    SourceLine         int
}
```

**Critical Issues:**
1. **Leaky Abstraction:** `JiraID` embeds Jira into domain model
   - Should be: `ExternalID string` or abstracted via ports
   - Why it matters: Prevents supporting GitHub, Linear, etc.
   - Fix difficulty: Medium (requires port interface changes)

2. **Stringly-Typed Custom Fields:** `map[string]string`
   - No type safety: "Story Points" could be "five" instead of 5
   - No validation: Invalid field names silently fail
   - No documentation: What fields are valid?
   - Should be: Typed field objects with validation

3. **Missing Domain Concepts:**
   - No `Priority` type (Low/Medium/High/Critical)
   - No `Status` type (Todo/InProgress/Done)
   - No `IssueType` type (Task/Story/Bug/Epic)
   - No `User` type (Assignee/Reporter)
   - No `Workflow` concept
   - Everything crammed into `CustomFields`

4. **Anemic Domain Model:**
   - No business logic methods
   - No validation
   - No invariants enforced
   - Just a data container

### Port Layer (/internal/core/ports/jira_port.go)

**JiraPort interface (8 methods):**
```go
type JiraPort interface {
    Authenticate() error
    CreateTask(task domain.Task, parentID string) (string, error)
    UpdateTask(task domain.Task) error
    GetProjectIssueTypes() (map[string][]string, error)
    GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error)
    CreateTicket(ticket domain.Ticket) (string, error)
    UpdateTicket(ticket domain.Ticket) error
    SearchTickets(ctx context.Context, projectKey string, jql string, progressCallback JiraProgressCallback) ([]domain.Ticket, error)
}
```

**Critical Issues:**
1. **Inconsistent Abstraction Levels:**
   - High-level: `CreateTicket(ticket)` - domain-centric
   - Low-level: `GetIssueTypeFields(string)` - Jira-specific
   - Mixed: `SearchTickets(...jql string...)` - JQL is Jira query language
   - Should be: Consistent abstraction at domain level

2. **Jira-Specific Naming:**
   - "JiraPort" - should be "IssueTrackerPort" or "TicketSystemPort"
   - Assumes only Jira will ever exist
   - Prevents polymorphism

3. **Return Types:**
   - `map[string]interface{}` - loses type safety
   - `(string, error)` - string is opaque ID
   - Should be: Typed responses, domain objects

4. **Missing Operations:**
   - No batch operations (create multiple tickets)
   - No transactional guarantees
   - No retry/circuit breaker hooks
   - No rate limiting support
   - No caching strategy

### Adapter Layer (/internal/adapters/jira/jira_adapter.go)

**JiraAdapter struct (1,137 lines):**
```go
type JiraAdapter struct {
    baseURL       string
    email         string
    apiKey        string
    projectKey    string
    storyType     string
    subTaskType   string
    client        *http.Client
    fieldMappings map[string]interface{}  // ← CONCERN: Runtime config
}
```

**Critical Issues:**
1. **God Object Anti-Pattern:**
   - Handles HTTP transport
   - Manages authentication
   - Performs JSON marshaling/unmarshaling
   - Maps field names
   - Constructs JQL queries
   - Converts domain objects
   - Handles errors
   - All in one struct

2. **Field Mapping Complexity:**
   - `map[string]interface{}` loses type safety
   - Nested maps for complex fields (Story Points)
   - Runtime configuration, no compile-time checks
   - Hardcoded defaults: `"customfield_10020"` (Sprint)
   - No validation of field IDs
   - No support for field schema discovery

3. **Error Handling:**
   - Generic `fmt.Errorf()` everywhere
   - Lost error context
   - No error taxonomy
   - No retryable vs. fatal distinction
   - HTTP status codes wrapped in strings

4. **Testing Challenges:**
   - HTTP client embedded (hard to mock)
   - No dependency injection
   - 60-second timeout hardcoded
   - Base64 encoding inline
   - Field mappings global

5. **Performance Issues:**
   - Synchronous HTTP calls (fixed with context, but not ideal)
   - No request batching
   - No connection pooling configuration
   - No caching of metadata (issue types, fields)
   - Fetches subtasks serially (N+1 query problem)

6. **Code Organization:**
   - 3 constructors (NewJiraAdapter, NewJiraAdapterWithConfig, NewJiraAdapterFromConfig)
   - Duplicated logic across methods
   - `buildFieldsPayload()` does too much
   - `parseJiraIssue()` and `parseJiraSubtask()` nearly identical

### Service Layer

**PushService (172 lines):**
```go
type PushService struct {
    repository   ports.Repository
    jiraClient   ports.JiraPort      // ← Tight coupling
    stateManager *state.StateManager
}
```

**Critical Issues:**
1. **Direct Adapter Coupling:**
   - Service depends on `ports.JiraPort` directly
   - No abstraction for "any issue tracker"
   - Hard to swap Jira for GitHub

2. **Transaction Semantics:**
   - Updates Jira, then saves file
   - What if file save fails?
   - No rollback mechanism
   - No atomic guarantees

3. **Error Aggregation:**
   - Continues processing after errors
   - Returns aggregate error at end
   - User doesn't know partial state
   - No transaction log

**PullService (269 lines):**
```go
type PullService struct {
    jiraAdapter    ports.JiraPort     // ← Direct dependency
    repository     ports.Repository
    dbRepository   ports.ExtendedRepository
    stateManager   *state.StateManager
    syncStrategy   ports.SyncStrategy
}
```

**Critical Issues:**
1. **Three Repositories:**
   - File repository (backward compat)
   - Database repository (TUI)
   - No clear primary source of truth
   - Conflicting writes possible

2. **Conflict Resolution:**
   - Strategy pattern good
   - But: relies on hash comparison only
   - No semantic merge
   - No user-driven resolution

3. **Progress Reporting:**
   - Callback-based
   - No structured events
   - No cancellation acknowledgment
   - No pause/resume

---

## Architectural Violations

### SOLID Principles Analysis

**Single Responsibility Principle: VIOLATED**
- JiraAdapter does: HTTP, auth, JSON, field mapping, domain conversion
- Should be: Separate transport, auth, serialization, mapping concerns

**Open/Closed Principle: VIOLATED**
- Adding GitHub adapter requires changing ports
- Adding new field type requires changing adapter
- Should be: Extensible without modification

**Liskov Substitution Principle: PARTIAL**
- Can't swap JiraPort for GitHubPort due to Jira specifics
- Should be: Any IssueTrackerPort substitutable

**Interface Segregation Principle: VIOLATED**
- JiraPort has 8 methods, not all clients need all
- Should be: Smaller, focused interfaces

**Dependency Inversion Principle: PARTIAL**
- Services depend on ports (good)
- But ports are Jira-shaped (bad)
- Should be: Ports shaped by domain needs

### Domain-Driven Design Analysis

**Ubiquitous Language: WEAK**
- Domain uses "Ticket", Jira uses "Issue"
- Domain uses "Task", Jira uses "Sub-task"
- Inconsistent terminology
- Should be: Consistent language across layers

**Bounded Contexts: UNCLEAR**
- Is Ticket a domain concept or Jira concept?
- Where does domain end and integration begin?
- Should be: Clear boundaries

**Aggregates: MISSING**
- Ticket contains Tasks, but no aggregate root enforcement
- No consistency boundaries
- No transaction boundaries
- Should be: Ticket is aggregate root over Tasks

**Value Objects: MISSING**
- Custom fields are strings, not value objects
- No immutability
- No validation
- Should be: Priority, Status, etc. as value objects

**Repositories: CONFUSED**
- Three different repositories
- No clear aggregate persistence
- File/DB/Jira all mixed
- Should be: One repository per aggregate

### Hexagonal Architecture Analysis

**Ports & Adapters: PARTIAL**
- Have ports (good)
- Have adapters (good)
- But: Ports are adapter-shaped, not domain-shaped (bad)
- Should be: Ports defined by domain use cases

**Inside vs. Outside: BLURRED**
- Domain has `JiraID` (outside concept inside)
- Services know about Jira (inside knows outside)
- Should be: Strict separation

**Dependency Direction: PARTIAL**
- Adapters depend on ports (good)
- But: Ports leak implementation details (bad)
- Should be: Pure domain-driven ports

---

## Research: Industry Best Practices

### External Validation Sources

We will consult these external references for architectural patterns:

**1. Atlassian Ecosystem:**
- Atlassian Connect framework architecture
- Jira REST API SDK design patterns
- Confluence integration best practices
- Reference: https://developer.atlassian.com/

**2. GitHub API Clients:**
- Octokit (GitHub's official SDK)
- go-github library architecture
- How they handle resources, pagination, rate limiting
- Reference: https://github.com/google/go-github

**3. Linear API:**
- GraphQL-based issue tracking
- Type-safe client generation
- How they model issues vs. Jira
- Reference: https://developers.linear.app/

**4. Domain-Driven Design Literature:**
- Eric Evans: "Domain-Driven Design" patterns
- Vaughn Vernon: "Implementing Domain-Driven Design"
- Martin Fowler: "Patterns of Enterprise Application Architecture"

**5. Go-Specific Patterns:**
- Kat Zien: "How Do You Structure Your Go Apps?"
- Kat Zien: "Practical Domain-Driven Design in Go"
- Robert C. Martin: "Clean Architecture" Go examples

**6. API Client Libraries:**
- Stripe Go SDK (excellent error handling)
- AWS SDK for Go v2 (modular design)
- Shopify Go SDK (resource abstraction)

### Pattern Analysis

**From GitHub Octokit:**
```go
// Resource-based organization
client.Issues.Create(ctx, owner, repo, issue)
client.PullRequests.List(ctx, owner, repo, opts)

// Strongly typed options
type IssueRequest struct {
    Title     *string
    Body      *string
    Assignees *[]string
    Labels    *[]string
}

// Pagination built-in
opts := &github.ListOptions{PerPage: 100}
issues, resp, err := client.Issues.List(ctx, opts)
```

**Key Lessons:**
- Resource-centric API (Issues, PRs, Users)
- Options pattern for flexibility
- Pointer fields for optionality
- Context always first param
- Structured responses with metadata

**From Stripe SDK:**
```go
// Error taxonomy
switch err := err.(type) {
case *stripe.Error:
    // Retryable or not?
    if err.Type == stripe.ErrorTypeAPIConnection {
        // Retry
    }
}

// Idempotency support
params := &stripe.ChargeParams{
    IdempotencyKey: stripe.String(uuid.New().String()),
}
```

**Key Lessons:**
- Typed errors with categories
- Idempotency for retries
- Clear retry semantics
- Network vs. business errors

**From AWS SDK v2:**
```go
// Modular clients
import "github.com/aws/aws-sdk-go-v2/service/s3"
client := s3.NewFromConfig(cfg)

// Configuration separate from client
cfg, err := config.LoadDefaultConfig(ctx)

// Operation-specific types
input := &s3.PutObjectInput{
    Bucket: aws.String("my-bucket"),
    Key:    aws.String("my-key"),
}
```

**Key Lessons:**
- Configuration separate from client
- Service-specific clients (modular)
- Input/Output types per operation
- Middleware for cross-cutting concerns

---

## Proposed Architecture

### Domain Layer Redesign

**Core Domain Models:**

```go
// Value Objects
type IssueID struct {
    value string
    system string // "jira", "github", "linear"
}

type Priority struct {
    level int // 1=Critical, 2=High, 3=Medium, 4=Low
}

type Status struct {
    name string
    category StatusCategory // Todo/InProgress/Done
}

type IssueType struct {
    name string // "Task", "Story", "Bug", "Epic"
}

type User struct {
    ID    string
    Name  string
    Email string
}

// Aggregate Root
type Issue struct {
    id                 IssueID
    title              string
    description        string
    issueType          IssueType
    priority           Priority
    status             Status
    assignee           *User
    reporter           *User
    labels             []string
    customFields       map[string]Field // Typed fields
    acceptanceCriteria []AcceptanceCriterion
    subIssues          []Issue
}

// Custom Field abstraction
type Field interface {
    Name() string
    Value() interface{}
    Validate() error
}

type TextField struct { value string }
type NumberField struct { value float64 }
type DateField struct { value time.Time }
type SelectField struct { value string; options []string }
```

**Benefits:**
- Type safety (can't put string in number field)
- Validation at domain level
- Clear semantics (Issue, not Ticket)
- System-agnostic (works for Jira, GitHub, Linear)
- Rich domain model with business logic

### Port Layer Redesign

**Issue Tracker Port:**

```go
// Main port interface
type IssueTracker interface {
    // Metadata operations
    GetMetadata(ctx context.Context) (*TrackerMetadata, error)

    // Issue operations
    CreateIssue(ctx context.Context, req *CreateIssueRequest) (*Issue, error)
    UpdateIssue(ctx context.Context, id IssueID, req *UpdateIssueRequest) (*Issue, error)
    GetIssue(ctx context.Context, id IssueID) (*Issue, error)

    // Query operations
    QueryIssues(ctx context.Context, query *IssueQuery) (*IssueQueryResult, error)

    // Batch operations
    BatchCreate(ctx context.Context, reqs []*CreateIssueRequest) (*BatchResult, error)
}

// Supporting types
type TrackerMetadata struct {
    Name          string
    IssueTypes    []IssueTypeMetadata
    CustomFields  []CustomFieldMetadata
    Workflows     []WorkflowMetadata
}

type IssueQuery struct {
    ProjectKey string
    IssueTypes []string
    Statuses   []string
    Labels     []string
    Query      string // System-specific query (JQL, GitHub search, etc.)
    Pagination *PaginationOptions
}

type IssueQueryResult struct {
    Issues      []Issue
    Total       int
    HasMore     bool
    NextCursor  string
}
```

**Benefits:**
- System-agnostic interface
- Consistent abstraction level
- Supports batching
- Pagination built-in
- Metadata discovery

**Error Taxonomy:**

```go
// Error types
type ErrorCategory int

const (
    ErrorCategoryAuth ErrorCategory = iota
    ErrorCategoryNetwork
    ErrorCategoryRateLimit
    ErrorCategoryValidation
    ErrorCategoryNotFound
    ErrorCategoryConflict
    ErrorCategoryInternal
)

type TrackerError struct {
    Category   ErrorCategory
    Message    string
    Retryable  bool
    RetryAfter *time.Duration
    Cause      error
}

func (e *TrackerError) Error() string {
    return e.Message
}

func (e *TrackerError) Unwrap() error {
    return e.Cause
}
```

**Benefits:**
- Typed error handling
- Retry semantics clear
- Rate limiting handled
- Error wrapping preserved

### Adapter Layer Redesign

**Jira Adapter Architecture:**

```go
// Separate concerns into components
type JiraAdapter struct {
    transport  *JiraTransport      // HTTP client, auth, retries
    serializer *JiraSerializer     // JSON ↔ Domain
    mapper     *JiraFieldMapper    // Field name ↔ Field ID
    metadata   *JiraMetadataCache  // Cached project/field info
}

// Transport handles HTTP
type JiraTransport struct {
    baseURL    string
    httpClient *http.Client
    auth       AuthProvider
    retrier    RetryPolicy
}

// Serializer handles JSON ↔ Domain
type JiraSerializer struct {
    // Converts between Jira API JSON and domain.Issue
}

// Field mapper handles custom fields
type JiraFieldMapper struct {
    projectKey string
    mappings   map[string]string // "Story Points" → "customfield_10016"
}

// Metadata cache reduces API calls
type JiraMetadataCache struct {
    issueTypes  []IssueTypeMetadata
    customFields []CustomFieldMetadata
    ttl         time.Duration
    lastFetch   time.Time
}
```

**Benefits:**
- Single Responsibility per component
- Easy to test each part
- Easy to mock transport for tests
- Metadata cached, reducing API calls
- Retry logic centralized

### Service Layer Redesign

**Use Case Services:**

```go
// SyncService orchestrates push/pull
type SyncService struct {
    tracker       ports.IssueTracker
    repository    ports.IssueRepository  // Single repository
    stateManager  ports.StateManager
    conflictHandler ports.ConflictHandler
}

func (s *SyncService) PushIssues(ctx context.Context, req *PushRequest) (*PushResult, error) {
    // Use case: Push local issues to remote tracker
    // 1. Load local issues
    // 2. Determine creates vs updates
    // 3. Batch creates
    // 4. Batch updates
    // 5. Handle errors
    // 6. Update state
    // 7. Return result
}

func (s *SyncService) PullIssues(ctx context.Context, req *PullRequest) (*PullResult, error) {
    // Use case: Pull remote issues to local
    // 1. Query remote issues
    // 2. Load local issues
    // 3. Detect conflicts
    // 4. Resolve conflicts (strategy pattern)
    // 5. Merge results
    // 6. Save to repository
    // 7. Update state
    // 8. Return result
}
```

**Benefits:**
- Use case driven (not CRUD)
- Clear transaction boundaries
- Batch operations for performance
- Conflict handling explicit
- Single repository (no file/DB confusion)

---

## Deliverables

### 1. Architecture Decision Records (ADRs)

Create ADRs documenting:
- ADR-001: Issue Tracker Abstraction Layer
- ADR-002: Domain Model Design
- ADR-003: Error Handling Strategy
- ADR-004: Field Mapping Architecture
- ADR-005: Repository Consolidation
- ADR-006: Batch Operations Support
- ADR-007: Conflict Resolution Strategy

**Format:**
```markdown
# ADR-001: Issue Tracker Abstraction Layer

## Status
Proposed / Accepted / Deprecated

## Context
[Why is this decision needed?]

## Decision
[What did we decide?]

## Consequences
[What are the trade-offs?]

## Alternatives Considered
[What else did we evaluate?]
```

### 2. Domain Model Documentation

- Ubiquitous language glossary
- Bounded context map
- Aggregate diagrams
- Value object catalog
- Entity lifecycle documentation

### 3. Port & Adapter Specification

- IssueTracker port interface specification
- Error taxonomy documentation
- Request/Response type catalog
- Adapter implementation guide

### 4. Migration Plan

- Phase 1: Add new domain models (non-breaking)
- Phase 2: Add new ports (parallel to existing)
- Phase 3: Implement new Jira adapter
- Phase 4: Migrate services to new ports
- Phase 5: Remove old implementation
- Phase 6: Performance validation

### 5. Test Strategy

- Unit tests for domain logic
- Contract tests for adapters
- Integration tests for services
- Performance benchmarks
- Migration test plan

### 6. External Validation Report

Document findings from:
- GitHub Octokit patterns
- Stripe SDK error handling
- AWS SDK modular design
- DDD pattern literature
- Go community best practices

**Format:**
```markdown
# External Validation: GitHub Octokit

## Pattern Observed
[What did we find?]

## Applicability to Ticketr
[How does this apply?]

## Recommendation
[Should we adopt? Why/why not?]

## Implementation Notes
[If adopted, how to implement?]
```

---

## Agent Assignments

### Day 1-2: Analysis & Design (Steward + Director)

**Steward (16 hours):**
- Analyze current architecture against SOLID principles
- Analyze current architecture against DDD patterns
- Analyze current architecture against Hexagonal Architecture
- Research GitHub Octokit patterns
- Research Stripe SDK error handling
- Research AWS SDK modular design
- Document findings in validation report
- Create initial ADRs

**Director (4 hours):**
- Coordinate Steward's research
- Review findings
- Identify critical architectural decisions
- Prioritize ADRs
- Plan Day 3-4 design work

**Deliverables:**
- Current state analysis document (this section expanded)
- External validation report (6 patterns analyzed)
- Initial ADR drafts (7 decisions)

### Day 3-4: Design & Validation (Builder + Steward)

**Builder (16 hours):**
- Design new domain model (Issue, IssueID, Priority, Status, etc.)
- Design new port interfaces (IssueTracker, ErrorTaxonomy)
- Design new adapter architecture (components separated)
- Design new service layer (use cases)
- Create prototypes in feature/jira-domain-redesign branch
- Write design documentation

**Steward (8 hours):**
- Review Builder's design
- Challenge design decisions
- Validate against DDD principles
- Validate against external patterns
- Approve or reject designs
- Iterate with Builder

**Deliverables:**
- Domain model specification
- Port & adapter specification
- Service layer specification
- Prototype code (not production, proof of concept)
- Design review meeting notes

### Day 5: Human Approval & Planning (Director)

**Director (8 hours):**
- Compile all findings into review package
- Create executive summary for human
- Prepare design presentation
- Document trade-offs
- Create migration plan
- Estimate implementation effort
- Request human approval

**Human Review:**
- Review architecture proposal
- Review ADRs
- Review external validation
- Approve design OR request changes
- Approve migration plan OR adjust scope

**Deliverables:**
- Architecture review package (all docs)
- Human approval (or change requests)
- Migration plan (approved)
- Implementation estimate (approved)

---

## Success Criteria

### Mandatory (Must Pass)

**Architecture Quality:**
- [ ] SOLID principles: No violations identified
- [ ] DDD patterns: Ubiquitous language documented
- [ ] DDD patterns: Aggregates clearly defined
- [ ] DDD patterns: Value objects for all domain concepts
- [ ] Hexagonal: Domain has zero external dependencies
- [ ] Hexagonal: Ports defined by domain needs, not adapter needs

**External Validation:**
- [ ] 6+ external patterns researched and documented
- [ ] At least 3 patterns adopted in design
- [ ] Validation report shows rationale for all decisions
- [ ] Industry best practices applied

**Documentation:**
- [ ] 7 ADRs created covering all major decisions
- [ ] Domain model fully documented with diagrams
- [ ] Port specifications complete with examples
- [ ] Migration plan approved by human
- [ ] All trade-offs explicitly documented

**Design Quality:**
- [ ] New design supports adding GitHub adapter without changing domain
- [ ] New design supports adding custom fields without changing adapter
- [ ] Error taxonomy covers all failure modes
- [ ] Batch operations reduce API calls by 80%+
- [ ] Test strategy covers domain, ports, adapters, services

### Human Approval Gate

**Required for proceeding to implementation:**
- [ ] Human reviews architecture package
- [ ] Human approves domain model
- [ ] Human approves port design
- [ ] Human approves migration plan
- [ ] Human approves timeline estimate
- [ ] Human confirms confidence restored in Jira integration

**No implementation work begins until human approval.**

---

## Timeline

### Day 1: Research & Analysis
**Morning (4 hours):**
- Steward: Analyze current architecture (SOLID, DDD, Hexagonal)
- Director: Setup feature/jira-domain-redesign branch

**Afternoon (4 hours):**
- Steward: Research GitHub Octokit patterns
- Steward: Research Stripe SDK error handling

**Evening (4 hours):**
- Steward: Research AWS SDK modular design
- Steward: Document findings in validation report

### Day 2: Analysis Completion
**Morning (4 hours):**
- Steward: Complete external validation report
- Steward: Create initial ADR drafts

**Afternoon (4 hours):**
- Director: Review Steward's findings
- Director: Prioritize architectural decisions

**Evening (4 hours):**
- Director + Steward: Alignment meeting (async)
- Steward: Finalize ADRs

### Day 3: Design Phase 1
**Morning (4 hours):**
- Builder: Design new domain model
- Steward: Review domain model design

**Afternoon (4 hours):**
- Builder: Design new port interfaces
- Steward: Review port design

**Evening (4 hours):**
- Builder: Create domain model prototype
- Steward: Create design review notes

### Day 4: Design Phase 2
**Morning (4 hours):**
- Builder: Design new adapter architecture
- Builder: Design new service layer

**Afternoon (4 hours):**
- Builder: Create adapter prototype
- Builder: Create service prototype

**Evening (4 hours):**
- Steward: Final design review
- Steward: Approve or request changes

### Day 5: Human Approval
**Morning (4 hours):**
- Director: Compile architecture review package
- Director: Create executive summary

**Afternoon (4 hours):**
- Director: Prepare human presentation
- Director: Create migration plan

**Evening (4 hours):**
- Human: Review architecture
- Human: Approve or request changes
- Director: Document outcome

---

## Branch Strategy

### Branch: feature/jira-domain-redesign

**Purpose:**
- All architectural design work happens here
- Prototypes created here
- ADRs committed here
- No production code affected

**Workflow:**
1. Branch from: `feature/v3`
2. Work in: `feature/jira-domain-redesign`
3. Merge to: `feature/v3` (after human approval)
4. Then: Implementation in separate feature branches

**Protection:**
- No direct commits to main
- Requires human approval before merge
- All ADRs reviewed before merge
- Design docs finalized before merge

---

## Risk Management

### Risk: Analysis Paralysis

**Symptom:** Research never ends, no decisions made

**Mitigation:**
- Day 1-2: Research only
- Day 3-4: Design only
- Day 5: Decision only
- No extending timeline without human approval

### Risk: Premature Implementation

**Symptom:** Builder starts coding before design approved

**Mitigation:**
- Director enforces: NO production code in Phase 7
- Prototypes only (throwaway code)
- Human approval required before implementation phase

### Risk: External Validation is Superficial

**Symptom:** "Looked at GitHub SDK, seems fine"

**Mitigation:**
- Steward must document: What pattern? Why relevant? How to adopt?
- Director reviews validation depth
- Each pattern requires 1+ page analysis

### Risk: Human Rejects Design

**Symptom:** Day 5 human review rejects architecture

**Mitigation:**
- Director presents trade-offs clearly
- Multiple options presented where uncertain
- Steward documents rationale for all decisions
- If rejected: Extend Phase 7 by 2 days, iterate

---

## Phase 7 Completion Checklist

**Architecture Analysis:**
- [ ] SOLID analysis complete and documented
- [ ] DDD analysis complete and documented
- [ ] Hexagonal analysis complete and documented
- [ ] Current violations cataloged

**External Validation:**
- [ ] GitHub Octokit patterns analyzed
- [ ] Stripe SDK patterns analyzed
- [ ] AWS SDK patterns analyzed
- [ ] 3+ additional sources analyzed
- [ ] Validation report complete (20+ pages)

**Design Artifacts:**
- [ ] 7 ADRs created and reviewed
- [ ] Domain model specification complete
- [ ] Port interface specification complete
- [ ] Adapter architecture specification complete
- [ ] Service layer specification complete
- [ ] Prototype code demonstrates feasibility

**Human Approval:**
- [ ] Architecture review package compiled
- [ ] Human has reviewed all materials
- [ ] Human has approved domain model
- [ ] Human has approved port design
- [ ] Human has approved migration plan
- [ ] Human confidence in Jira integration restored

**Planning:**
- [ ] Migration plan approved (6 phases)
- [ ] Implementation timeline estimated
- [ ] Resource allocation planned
- [ ] Test strategy approved
- [ ] Performance targets defined

**Branch Hygiene:**
- [ ] feature/jira-domain-redesign branch clean
- [ ] All ADRs committed
- [ ] All design docs committed
- [ ] Prototype code committed
- [ ] Ready to merge to feature/v3 (post-approval)

---

## Post-Phase 7: Implementation Plan

**If Human Approves Architecture:**

### Phase 7.1: Domain Model Implementation (Week 1)
- Builder implements new domain models
- Verifier creates domain tests
- Scribe documents domain concepts

### Phase 7.2: Port Layer Implementation (Week 1-2)
- Builder implements new port interfaces
- Builder creates port contract tests
- Scribe documents port specifications

### Phase 7.3: Jira Adapter Rewrite (Week 2-3)
- Builder implements new Jira adapter (components separated)
- Builder implements field mapper
- Builder implements metadata cache
- Verifier creates adapter integration tests

### Phase 7.4: Service Layer Migration (Week 3-4)
- Builder migrates PushService to new ports
- Builder migrates PullService to new ports
- Verifier creates service integration tests
- Verifier runs full regression suite

### Phase 7.5: Legacy Removal (Week 4)
- Builder removes old JiraPort
- Builder removes old JiraAdapter
- Builder removes old domain models
- Scribe updates all documentation

### Phase 7.6: Performance Validation (Week 4)
- Verifier runs performance benchmarks
- Verifier compares old vs new performance
- Verifier validates batch operations
- Steward approves performance results

**Total Implementation Estimate: 4 weeks**

---

## Communication Plan

### Daily Standup (Async)

**Format:**
```markdown
## Phase 7 Day N Standup

**Steward:**
- Yesterday: [What was analyzed/researched]
- Today: [What will be analyzed/researched]
- Blockers: [Any issues]

**Builder:**
- Yesterday: [What was designed/prototyped]
- Today: [What will be designed/prototyped]
- Blockers: [Any issues]

**Director:**
- Yesterday: [What was coordinated]
- Today: [What will be coordinated]
- Decisions needed: [Any human input needed]
```

### Human Check-ins

**Day 1 Evening:** Research progress update
**Day 2 Evening:** Validation report preview
**Day 3 Evening:** Domain model design review
**Day 4 Evening:** Full design preview
**Day 5 Morning:** Final approval request

### Escalation Path

**If stuck:** Director escalates to human within 4 hours
**If behind:** Director proposes timeline adjustment
**If design conflict:** Steward documents options, human decides

---

## Conclusion

Phase 7 is a **mandatory architectural review** before we can confidently build on the Jira integration foundation. The human has demanded we question everything - and we will.

**This is about restoring confidence through rigor.**

**Status:** Ready to initiate
**Next Step:** Director creates feature/jira-domain-redesign branch and assigns Day 1 tasks to Steward

**Created:** 2025-10-21
**Priority:** CRITICAL
**Blocking:** All future Jira integration work, multi-tracker support, architectural integrity

---

## Appendix: Key Questions to Answer

### Domain Model Questions
1. Should we call it "Issue" or "Ticket"? (Jira says Issue, we say Ticket)
2. Should IssueID be in domain or port layer?
3. Should CustomFields be stringly-typed or strongly-typed?
4. Should Status be a value object or enum?
5. Should User be in domain or just a string ID?

### Port Design Questions
1. Should port be called JiraPort or IssueTrackerPort?
2. Should port methods be CRUD or use case driven?
3. Should port handle batching or should service layer?
4. Should port expose JQL or abstract query language?
5. Should port return domain objects or DTOs?

### Adapter Questions
1. Should adapter be one class or multiple components?
2. Should field mapping be configuration or code?
3. Should metadata be cached or fetched per request?
4. Should retry logic be in adapter or transport layer?
5. Should serialization be in adapter or separate component?

### Service Questions
1. Should we have one service or multiple (Push, Pull, Sync)?
2. Should conflict resolution be in service or separate handler?
3. Should we support file + database or just database?
4. Should batch operations be service responsibility?
5. Should transaction boundaries be service level?

**All these questions will be answered by Day 5.**
