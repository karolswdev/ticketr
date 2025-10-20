# Bulk Operations API Documentation

**Last Updated:** October 19, 2025 (Week 18)
**Audience:** Developers integrating bulk operations into Ticketr
**Status:** Production-ready

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Domain Model](#domain-model)
3. [Service Interface](#service-interface)
4. [Usage Examples](#usage-examples)
5. [Error Handling](#error-handling)
6. [Testing](#testing)

---

## Architecture Overview

Ticketr's bulk operations feature follows the Hexagonal Architecture pattern, with clear separation between domain logic, service orchestration, and external adapters.

### Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     CLI Layer                                │
│                (cmd/ticketr/bulk_commands.go)                │
│   - Parse flags and arguments                               │
│   - Display progress feedback                               │
│   - Handle user interaction                                 │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                  Service Layer                               │
│        (internal/core/services/bulk_operation_service.go)   │
│   - Validate bulk operations                                │
│   - Execute operations with progress callbacks              │
│   - Handle rollback on partial failure                      │
│   - Coordinate with Jira adapter                            │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                   Ports Layer                                │
│            (internal/core/ports/bulk_operation.go)          │
│   - BulkOperationService interface                          │
│   - BulkOperationProgressCallback type                      │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                  Domain Layer                                │
│          (internal/core/domain/bulk_operation.go)           │
│   - BulkOperation model                                     │
│   - BulkOperationResult model                               │
│   - Validation rules                                        │
│   - JQL injection prevention                                │
└─────────────────────────────────────────────────────────────┘
```

### Data Flow

```
User Input (CLI)
    ↓
Parse ticket IDs and changes
    ↓
Create BulkOperation domain object
    ↓
Validate operation (domain rules)
    ↓
Service.ExecuteOperation(ctx, operation, progressCallback)
    ↓
For each ticket:
    ├─> Fetch current state (for rollback)
    ├─> Apply changes
    ├─> Update via Jira adapter
    ├─> Record success/failure
    └─> Invoke progress callback
    ↓
If partial failure:
    └─> Attempt best-effort rollback
    ↓
Return BulkOperationResult
```

---

## Domain Model

### BulkOperation

**Location:** `internal/core/domain/bulk_operation.go`

Represents a bulk action to be performed on multiple tickets.

```go
type BulkOperation struct {
    // Action specifies the bulk operation to perform (update, move, delete).
    Action BulkOperationAction

    // TicketIDs contains the Jira ticket IDs (e.g., "PROJ-123") to operate on.
    // Must contain between 1 and 100 ticket IDs.
    TicketIDs []string

    // Changes contains the field modifications to apply (for update operations).
    // Key is the field name, value is the new value.
    // Required for "update" action, ignored for other actions.
    Changes map[string]interface{}
}
```

#### BulkOperationAction

```go
type BulkOperationAction string

const (
    // BulkActionUpdate modifies field values on multiple tickets.
    BulkActionUpdate BulkOperationAction = "update"

    // BulkActionMove transitions multiple tickets to a new status.
    BulkActionMove BulkOperationAction = "move"

    // BulkActionDelete removes multiple tickets from the system.
    BulkActionDelete BulkOperationAction = "delete"
)
```

#### Constructor

```go
func NewBulkOperation(
    action BulkOperationAction,
    ticketIDs []string,
    changes map[string]interface{},
) *BulkOperation
```

**Parameters:**
- `action`: The operation type (update, move, delete)
- `ticketIDs`: Slice of Jira ticket IDs
- `changes`: Map of field names to new values (can be `nil` for delete)

**Returns:** Pointer to new `BulkOperation` instance

---

### BulkOperationResult

Represents the outcome of a bulk operation.

```go
type BulkOperationResult struct {
    // SuccessCount is the number of tickets successfully processed.
    SuccessCount int

    // FailureCount is the number of tickets that failed to process.
    FailureCount int

    // Errors contains any errors encountered during processing.
    // Map key is the ticket ID, value is the error message.
    Errors map[string]string

    // SuccessfulTickets contains the IDs of tickets that were processed successfully.
    SuccessfulTickets []string

    // FailedTickets contains the IDs of tickets that failed to process.
    FailedTickets []string
}
```

#### Constructor

```go
func NewBulkOperationResult() *BulkOperationResult
```

**Returns:** Pointer to new `BulkOperationResult` with initialized maps and slices

---

### Validation Rules

#### Ticket ID Validation

All ticket IDs are validated against a strict regex pattern to prevent JQL injection:

```go
var jiraIDRegex = regexp.MustCompile(`^[A-Z]+-\d+$`)
```

**Valid examples:**
- `PROJ-123`
- `BACKEND-42`
- `EPIC-1`

**Invalid examples (rejected):**
- `proj-123` (lowercase)
- `PROJ_123` (underscore)
- `PROJ-123" OR 1=1` (SQL injection attempt)
- `123` (no project key)

#### Operation Validation

The `Validate()` method checks:

1. **Action validity**: Must be `update`, `move`, or `delete`
2. **Ticket count**: Between 1 and 100 tickets
3. **Ticket ID format**: All IDs match Jira pattern
4. **Changes requirement**: Update operations must have at least one change

```go
func (bo *BulkOperation) Validate() error
```

**Returns:** Error if validation fails, `nil` if valid

---

## Service Interface

### BulkOperationService

**Location:** `internal/core/ports/bulk_operation.go`

```go
type BulkOperationService interface {
    ExecuteOperation(
        ctx context.Context,
        op *domain.BulkOperation,
        progress BulkOperationProgressCallback,
    ) (*domain.BulkOperationResult, error)
}
```

#### ExecuteOperation

Performs a bulk operation on multiple tickets with real-time progress tracking.

**Parameters:**
- `ctx`: Context for cancellation and timeout
- `op`: The bulk operation to execute
- `progress`: Callback function invoked after each ticket (can be `nil`)

**Returns:**
- `*domain.BulkOperationResult`: Detailed result including success/failure counts
- `error`: Operation-level error (e.g., invalid operation, all tickets failed)

**Behavior:**
- Validates the operation before execution
- Processes tickets sequentially
- Respects context cancellation
- Invokes progress callback after each ticket
- Attempts best-effort rollback on partial failure (update/move only)
- Returns partial results even if some tickets fail

#### BulkOperationProgressCallback

```go
type BulkOperationProgressCallback func(
    ticketID string,
    success bool,
    err error,
)
```

**Parameters:**
- `ticketID`: The Jira ticket ID being processed
- `success`: `true` if operation succeeded, `false` if failed
- `err`: Error details if `success` is `false`, `nil` otherwise

**Usage:** Invoked after each ticket is processed, allowing real-time progress updates in CLI or TUI.

---

## Usage Examples

### Example 1: Basic Update Operation

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/karolswdev/ticktr/internal/core/domain"
    "github.com/karolswdev/ticktr/internal/core/services"
    "github.com/karolswdev/ticktr/internal/adapters/jira"
)

func main() {
    // 1. Initialize Jira adapter
    jiraAdapter, err := jira.NewJiraAdapter(jiraConfig)
    if err != nil {
        log.Fatalf("Failed to create Jira adapter: %v", err)
    }

    // 2. Create bulk operation service
    service := services.NewBulkOperationService(jiraAdapter)

    // 3. Define bulk operation
    ticketIDs := []string{"PROJ-1", "PROJ-2", "PROJ-3"}
    changes := map[string]interface{}{
        "status": "Done",
        "assignee": "john@example.com",
    }

    operation := domain.NewBulkOperation(
        domain.BulkActionUpdate,
        ticketIDs,
        changes,
    )

    // 4. Define progress callback
    progressCallback := func(ticketID string, success bool, err error) {
        if success {
            fmt.Printf("✓ %s updated\n", ticketID)
        } else {
            fmt.Printf("✗ %s failed: %v\n", ticketID, err)
        }
    }

    // 5. Execute operation
    ctx := context.Background()
    result, err := service.ExecuteOperation(ctx, operation, progressCallback)

    // 6. Handle result
    if err != nil {
        log.Printf("Operation completed with errors: %v", err)
    }

    fmt.Printf("\nSummary:\n")
    fmt.Printf("  Success: %d\n", result.SuccessCount)
    fmt.Printf("  Failure: %d\n", result.FailureCount)

    if result.FailureCount > 0 {
        fmt.Println("\nFailed tickets:")
        for ticketID, errMsg := range result.Errors {
            fmt.Printf("  %s: %s\n", ticketID, errMsg)
        }
    }
}
```

---

### Example 2: Move Operation

```go
package main

import (
    "context"
    "fmt"

    "github.com/karolswdev/ticktr/internal/core/domain"
    "github.com/karolswdev/ticktr/internal/core/services"
)

func moveTicketsToParent(
    service ports.BulkOperationService,
    ticketIDs []string,
    parentID string,
) error {
    // Create move operation
    changes := map[string]interface{}{
        "parent": parentID,
    }

    operation := domain.NewBulkOperation(
        domain.BulkActionMove,
        ticketIDs,
        changes,
    )

    // Validate before execution
    if err := operation.Validate(); err != nil {
        return fmt.Errorf("invalid operation: %w", err)
    }

    // Execute with simple progress callback
    ctx := context.Background()
    result, err := service.ExecuteOperation(ctx, operation, func(id string, ok bool, e error) {
        if ok {
            fmt.Printf("[✓] %s moved to %s\n", id, parentID)
        } else {
            fmt.Printf("[✗] %s failed: %v\n", id, e)
        }
    })

    // Check for complete failure
    if err != nil && result.SuccessCount == 0 {
        return fmt.Errorf("all tickets failed: %w", err)
    }

    return nil
}
```

---

### Example 3: Context Cancellation

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/karolswdev/ticktr/internal/core/domain"
)

func updateWithTimeout(
    service ports.BulkOperationService,
    ticketIDs []string,
    changes map[string]interface{},
    timeout time.Duration,
) (*domain.BulkOperationResult, error) {
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    // Create operation
    operation := domain.NewBulkOperation(
        domain.BulkActionUpdate,
        ticketIDs,
        changes,
    )

    // Execute with timeout
    result, err := service.ExecuteOperation(ctx, operation, nil)

    // Check if operation was cancelled
    if ctx.Err() == context.DeadlineExceeded {
        return result, fmt.Errorf("operation timed out after %v: %w", timeout, err)
    }

    return result, err
}
```

---

### Example 4: TUI Integration

```go
package tui

import (
    "context"
    "fmt"

    "github.com/karolswdev/ticktr/internal/core/domain"
    "github.com/karolswdev/ticktr/internal/core/ports"
)

type BulkOperationModal struct {
    service  ports.BulkOperationService
    progress chan string
}

func (m *BulkOperationModal) executeBulkUpdate(
    ticketIDs []string,
    changes map[string]interface{},
) {
    // Create operation
    operation := domain.NewBulkOperation(
        domain.BulkActionUpdate,
        ticketIDs,
        changes,
    )

    // Progress callback sends updates to TUI channel
    progressCallback := func(ticketID string, success bool, err error) {
        if success {
            m.progress <- fmt.Sprintf("✓ %s", ticketID)
        } else {
            m.progress <- fmt.Sprintf("✗ %s: %v", ticketID, err)
        }
    }

    // Execute in background goroutine
    go func() {
        ctx := context.Background()
        result, err := m.service.ExecuteOperation(ctx, operation, progressCallback)

        // Send summary to TUI
        if err != nil {
            m.progress <- fmt.Sprintf("Completed with %d failures", result.FailureCount)
        } else {
            m.progress <- fmt.Sprintf("All %d tickets updated", result.SuccessCount)
        }

        close(m.progress)
    }()
}
```

---

## Error Handling

### Operation-Level Errors

These errors are returned by `ExecuteOperation` and indicate the entire operation failed:

```go
// Invalid operation
if err := op.Validate(); err != nil {
    return nil, fmt.Errorf("invalid bulk operation: %w", err)
}

// All tickets failed
if result.FailureCount == len(op.TicketIDs) {
    return result, fmt.Errorf("all tickets failed to update")
}

// Partial failure with rollback
if result.FailureCount > 0 && result.SuccessCount > 0 {
    return result, fmt.Errorf("partial failure: %d of %d tickets failed (rollback attempted)",
        result.FailureCount, len(op.TicketIDs))
}
```

### Per-Ticket Errors

Individual ticket failures are captured in `BulkOperationResult.Errors`:

```go
result, err := service.ExecuteOperation(ctx, operation, nil)

// Check individual ticket errors
for ticketID, errMsg := range result.Errors {
    log.Printf("Ticket %s failed: %s", ticketID, errMsg)
}
```

### Error Types

| Error Type | Description | Recovery |
|------------|-------------|----------|
| **Validation Error** | Invalid operation (bad ticket IDs, missing changes) | Fix input and retry |
| **Authentication Error** | Invalid Jira credentials | Update workspace credentials |
| **Not Found Error** | Ticket does not exist | Remove ticket from list |
| **Permission Error** | User lacks access to ticket | Check Jira permissions |
| **Network Error** | Jira API unreachable | Retry after network recovery |
| **Partial Failure** | Some tickets succeeded, some failed | Check errors, retry failed tickets |

### Error Wrapping

Errors are wrapped using Go 1.13+ error wrapping:

```go
if err != nil {
    return fmt.Errorf("failed to update ticket %s: %w", ticketID, err)
}
```

**Unwrap errors:**
```go
import "errors"

if errors.Is(err, ports.ErrWorkspaceNotFound) {
    // Handle workspace not found
}
```

---

## Testing

### Unit Testing with Mock Service

```go
package mypackage_test

import (
    "context"
    "testing"

    "github.com/karolswdev/ticktr/internal/core/domain"
    "github.com/karolswdev/ticktr/internal/core/ports"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockBulkOperationService is a mock implementation for testing
type MockBulkOperationService struct {
    mock.Mock
}

func (m *MockBulkOperationService) ExecuteOperation(
    ctx context.Context,
    op *domain.BulkOperation,
    progress ports.BulkOperationProgressCallback,
) (*domain.BulkOperationResult, error) {
    args := m.Called(ctx, op, progress)
    return args.Get(0).(*domain.BulkOperationResult), args.Error(1)
}

func TestBulkUpdate_Success(t *testing.T) {
    // Create mock service
    mockService := new(MockBulkOperationService)

    // Setup expectations
    expectedResult := &domain.BulkOperationResult{
        SuccessCount: 3,
        FailureCount: 0,
        SuccessfulTickets: []string{"PROJ-1", "PROJ-2", "PROJ-3"},
        Errors: make(map[string]string),
    }

    mockService.On("ExecuteOperation",
        mock.Anything,
        mock.Anything,
        mock.Anything,
    ).Return(expectedResult, nil)

    // Test your code that uses the service
    operation := domain.NewBulkOperation(
        domain.BulkActionUpdate,
        []string{"PROJ-1", "PROJ-2", "PROJ-3"},
        map[string]interface{}{"status": "Done"},
    )

    result, err := mockService.ExecuteOperation(context.Background(), operation, nil)

    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, 3, result.SuccessCount)
    assert.Equal(t, 0, result.FailureCount)

    // Verify mock was called
    mockService.AssertExpectations(t)
}
```

---

### Integration Testing

```go
package services_test

import (
    "context"
    "testing"

    "github.com/karolswdev/ticktr/internal/adapters/jira"
    "github.com/karolswdev/ticktr/internal/core/domain"
    "github.com/karolswdev/ticktr/internal/core/services"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestBulkOperationService_Integration(t *testing.T) {
    // Skip if integration tests disabled
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // Initialize real Jira adapter
    jiraAdapter, err := jira.NewJiraAdapter(testConfig)
    require.NoError(t, err)

    // Create service
    service := services.NewBulkOperationService(jiraAdapter)

    // Create test tickets
    ticketIDs := []string{"PROJ-1", "PROJ-2"}
    changes := map[string]interface{}{
        "status": "In Progress",
    }

    operation := domain.NewBulkOperation(
        domain.BulkActionUpdate,
        ticketIDs,
        changes,
    )

    // Execute operation
    ctx := context.Background()
    result, err := service.ExecuteOperation(ctx, operation, nil)

    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, 2, result.SuccessCount)
    assert.Equal(t, 0, result.FailureCount)

    // Cleanup: Reset ticket states
    defer resetTestTickets(t, jiraAdapter, ticketIDs)
}
```

---

### Test Coverage Expectations

**Minimum Coverage:**
- Domain model (validation): 100%
- Service layer (execution): 85%
- CLI commands (integration): 70%

**Critical Paths:**
- JQL injection prevention: 100% coverage
- Rollback logic: 90% coverage
- Progress callback invocation: 100% coverage

---

## Performance Considerations

### Sequential Processing

Tickets are processed **sequentially** (one at a time) to:
- Simplify rollback logic
- Avoid Jira API rate limiting
- Maintain predictable progress feedback

**Typical Performance:**
- 1 ticket: ~200ms
- 10 tickets: ~2 seconds
- 100 tickets: ~20 seconds

### Future Optimization (v3.2.0)

Parallel processing is planned for v3.2.0:
- Process tickets in batches of 10
- Reduce total execution time by 50-70%
- Maintain rollback capability

---

## Security Considerations

### JQL Injection Prevention

All ticket IDs are validated with strict regex before any Jira API calls:

```go
var jiraIDRegex = regexp.MustCompile(`^[A-Z]+-\d+$`)
```

**Blocked patterns:**
- `PROJ-1" OR 1=1`
- `PROJ-1; DROP TABLE`
- `PROJ-1 UNION SELECT`

**Defense-in-depth:**
1. Regex validation at domain level
2. Parameterized JQL queries in Jira adapter
3. Error on invalid characters

### Credential Security

Bulk operations use workspace credentials stored in OS keychain:
- No credentials in code or logs
- Credentials fetched per operation
- Automatic credential cleanup on workspace deletion

---

## Best Practices for Developers

### 1. Always Validate Operations

```go
operation := domain.NewBulkOperation(action, ticketIDs, changes)

// Validate before execution
if err := operation.Validate(); err != nil {
    return fmt.Errorf("invalid operation: %w", err)
}
```

### 2. Provide Progress Callbacks

```go
// Don't do this (no user feedback)
result, err := service.ExecuteOperation(ctx, operation, nil)

// Do this (real-time feedback)
result, err := service.ExecuteOperation(ctx, operation, progressCallback)
```

### 3. Handle Partial Failures

```go
result, err := service.ExecuteOperation(ctx, operation, progressCallback)

// Check for partial failures
if result.FailureCount > 0 {
    log.Printf("Partial failure: %d of %d failed", result.FailureCount, len(ticketIDs))

    // Retry failed tickets
    retryOperation := domain.NewBulkOperation(action, result.FailedTickets, changes)
    retryResult, retryErr := service.ExecuteOperation(ctx, retryOperation, progressCallback)
}
```

### 4. Use Context for Cancellation

```go
// Support user cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Listen for interrupt signal
go func() {
    <-interruptChan
    cancel()
}()

result, err := service.ExecuteOperation(ctx, operation, progressCallback)
```

---

## Summary

Ticketr's Bulk Operations API provides:

- **Clean domain model** with validation and security
- **Service interface** supporting progress callbacks
- **Best-effort rollback** on partial failures
- **Context support** for cancellation and timeouts
- **Comprehensive error handling** at operation and ticket levels

**Key Files:**
- Domain: `internal/core/domain/bulk_operation.go`
- Service: `internal/core/services/bulk_operation_service.go`
- Interface: `internal/core/ports/bulk_operation.go`
- CLI: `cmd/ticketr/bulk_commands.go`

**Next Steps:**
- Read [Bulk Operations Guide](bulk-operations-guide.md) for user-facing documentation
- Review [ARCHITECTURE.md](ARCHITECTURE.md) for hexagonal architecture overview
- Check [REQUIREMENTS-v2.md](development/REQUIREMENTS.md) for complete feature requirements

---

**Document Version:** 1.0
**Status:** Production-ready
**Feedback:** [GitHub Issues](https://github.com/karolswdev/ticketr/issues)
