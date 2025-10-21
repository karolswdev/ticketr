# TUI Async Job Queue Architecture

**Version:** 1.0
**Date:** 2025-10-20
**Phase:** Phase 6, Week 2 Day 6-7
**Status:** Implementation Complete

## Overview

This document describes the async job queue architecture for Ticketr's TUI. The system provides non-blocking operations with real-time progress updates and cancellation support, ensuring the TUI remains responsive during long-running operations like pulling tickets from Jira.

## Goals

1. **Non-blocking TUI**: Users can navigate views while operations run in background
2. **Real-time Progress**: Show progress updates during long-running operations
3. **Cancellation Support**: ESC/Ctrl+C can cancel active operations
4. **Thread Safety**: Zero race conditions, no goroutine leaks
5. **Simple Integration**: Easy to wrap existing services (PullService, PushService)

## Architecture

### Core Components

```
┌─────────────────────────────────────────────────────────────┐
│                         TUI Layer                            │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │  app.go      │───▶│ JobQueue     │◀───│  PullJob     │  │
│  │              │    │  Manager     │    │  (wrapper)   │  │
│  │ - Submit jobs│    │              │    │              │  │
│  │ - Cancel     │    │ - Workers    │    │ - Pull logic │  │
│  │ - Progress   │    │ - Channels   │    │ - Progress   │  │
│  └──────────────┘    └──────────────┘    └──────────────┘  │
│         │                   │                     │          │
│         │                   │                     │          │
│         ▼                   ▼                     ▼          │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Progress Channel (buffered)              │  │
│  │         JobProgress events flow to UI                 │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                      Core Services                           │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │ PullService  │    │ PushService  │    │  Other       │  │
│  │              │    │              │    │  Services    │  │
│  │ - Pull()     │    │ - Push()     │    │              │  │
│  │ - Callbacks  │    │              │    │              │  │
│  └──────────────┘    └──────────────┘    └──────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Job Interface

The `Job` interface defines the contract for any async operation:

```go
type Job interface {
    Execute(ctx context.Context) error  // Run the job
    ID() JobID                           // Unique identifier
    Progress() <-chan JobProgress        // Progress events channel
}
```

**Design Rationale:**
- `Execute()` receives `context.Context` for cancellation support
- `ID()` enables tracking and cancellation by job ID
- `Progress()` returns read-only channel to prevent external writes

### Job Types

```go
type JobID string  // UUID or sequential ID

type JobStatus string
const (
    JobPending    JobStatus = "pending"    // Queued, not started
    JobRunning    JobStatus = "running"    // Currently executing
    JobCompleted  JobStatus = "completed"  // Finished successfully
    JobFailed     JobStatus = "failed"     // Finished with error
    JobCancelled  JobStatus = "cancelled"  // Cancelled by user
)

type JobProgress struct {
    JobID      JobID      // Which job this update is for
    Current    int        // Current count (e.g., 45 tickets)
    Total      int        // Total count (0 if unknown)
    Percentage float64    // Calculated percentage (0-100)
    Message    string     // Human-readable status
}
```

### JobQueue Manager

The `JobQueue` manages job execution, cancellation, and progress aggregation.

**Key Features:**
- **Worker Pool**: Fixed number of goroutines (1-3 workers initially)
- **Job Channel**: Buffered channel for incoming jobs
- **Progress Aggregation**: Single channel for all job progress events
- **Context Management**: Track and cancel jobs by ID
- **Status Tracking**: Query job status at any time

**API:**
```go
type JobQueue struct {
    workers     int
    jobChan     chan Job
    progressChan chan JobProgress
    contexts    map[JobID]context.CancelFunc  // Protected by mutex
    statuses    map[JobID]JobStatus           // Protected by mutex
    mu          sync.Mutex
    wg          sync.WaitGroup
}

func NewJobQueue(workerCount int) *JobQueue
func (jq *JobQueue) Submit(job Job) JobID
func (jq *JobQueue) Cancel(jobID JobID) error
func (jq *JobQueue) Status(jobID JobID) JobStatus
func (jq *JobQueue) Progress() <-chan JobProgress
func (jq *JobQueue) Shutdown() error
```

### PullJob Implementation

Wraps `PullService.Pull()` to implement the `Job` interface:

```go
type PullJob struct {
    id           JobID
    pullService  *services.PullService
    filePath     string
    options      services.PullOptions
    progressChan chan JobProgress
    result       *services.PullResult
    err          error
    mu           sync.Mutex  // Protects result and err
}
```

**Progress Callback Integration:**
```go
// In Execute():
options.ProgressCallback = func(current, total int, message string) {
    select {
    case progressChan <- JobProgress{
        JobID:      job.id,
        Current:    current,
        Total:      total,
        Percentage: calculatePercentage(current, total),
        Message:    message,
    }:
    case <-ctx.Done():
        return  // Don't block if cancelled
    }
}
```

## Sequence Diagrams

### Job Submission Flow

```
User          TUI App      JobQueue     Worker      PullService
 │               │             │           │              │
 │  Press 'P'    │             │           │              │
 ├──────────────▶│             │           │              │
 │               │  Submit()   │           │              │
 │               ├────────────▶│           │              │
 │               │             │  job      │              │
 │               │             ├──────────▶│              │
 │               │  JobID      │           │  Execute()   │
 │               ◀─────────────┤           ├─────────────▶│
 │               │             │           │              │
 │  UI update    │             │  progress │  callback    │
 ◀───────────────┤◀────────────┤◀──────────┤◀─────────────┤
 │  "Pulling..." │             │           │              │
 │               │             │           │              │
 │  UI update    │             │  progress │  callback    │
 ◀───────────────┤◀────────────┤◀──────────┤◀─────────────┤
 │  "45/120..."  │             │           │              │
 │               │             │           │  complete    │
 │               │             │◀──────────┤◀─────────────┤
 │  UI update    │             │           │              │
 ◀───────────────┤             │           │              │
 │  "Complete"   │             │           │              │
```

### Cancellation Flow

```
User          TUI App      JobQueue     Worker      PullService
 │               │             │           │              │
 │  Operation    │             │           │              │
 │  in progress  │             │           ├─────────────▶│
 │               │             │           │   Pull()     │
 │  Press ESC    │             │           │              │
 ├──────────────▶│             │           │              │
 │               │  Cancel(ID) │           │              │
 │               ├────────────▶│           │              │
 │               │             │  ctx.     │              │
 │               │             │  Cancel() │              │
 │               │             ├──────────▶│              │
 │  UI update    │             │           │  ctx.Done()  │
 ◀───────────────┤             │           ├─────────────▶│
 │  "Cancelling" │             │           │              │
 │               │             │           │  return      │
 │               │             │◀──────────┤◀─────────────┤
 │  UI update    │             │           │              │
 ◀───────────────┤             │           │              │
 │  "Cancelled:  │             │           │              │
 │   45 pulled"  │             │           │              │
```

## Concurrency Patterns

### Worker Pool Pattern

Each worker goroutine pulls jobs from a channel:

```go
func (jq *JobQueue) worker() {
    for job := range jq.jobChan {
        // Create cancellable context
        ctx, cancel := context.WithCancel(context.Background())

        // Register context for cancellation
        jq.mu.Lock()
        jq.contexts[job.ID()] = cancel
        jq.statuses[job.ID()] = JobRunning
        jq.mu.Unlock()

        // Forward progress events
        go jq.forwardProgress(job)

        // Execute job
        err := job.Execute(ctx)

        // Update status
        jq.mu.Lock()
        if ctx.Err() == context.Canceled {
            jq.statuses[job.ID()] = JobCancelled
        } else if err != nil {
            jq.statuses[job.ID()] = JobFailed
        } else {
            jq.statuses[job.ID()] = JobCompleted
        }
        delete(jq.contexts, job.ID())
        jq.mu.Unlock()

        cancel()  // Clean up
    }
}
```

### Progress Forwarding

```go
func (jq *JobQueue) forwardProgress(job Job) {
    for progress := range job.Progress() {
        select {
        case jq.progressChan <- progress:
        default:
            // Drop progress if channel full (non-blocking)
        }
    }
}
```

**Design Decision:** We use `default` case to prevent blocking if the progress channel is full. This ensures that slow UI updates don't block job execution.

### Channel Discipline

**Buffering Strategy:**
- `jobChan`: Buffered (10) - Allows burst job submissions
- `progressChan`: Buffered (100) - Prevents blocking on UI updates
- Individual job progress: Buffered (50) - Smooth progress updates

**Closing Channels:**
- Job channels closed by job when complete
- Progress channel remains open (managed by JobQueue)
- No sends to closed channels (checked via context cancellation)

## Thread Safety Strategy

### Mutex Protection

Protected by `JobQueue.mu`:
- `contexts map[JobID]context.CancelFunc`
- `statuses map[JobID]JobStatus`

**Rationale:** Maps require mutex protection for concurrent access.

### Channel-Based Communication

No mutex needed:
- Job distribution (jobChan)
- Progress events (progressChan)
- Individual job progress

**Rationale:** Channels provide built-in synchronization.

### Lock-Free Reads

Job state (result, err) uses mutex only during writes. Reads after completion are safe without locks.

## Error Handling

### Network Failures

```
┌──────────────────────────────────────────────────────┐
│  Network Error During Pull                           │
├──────────────────────────────────────────────────────┤
│  1. PullService.Pull() returns error                 │
│  2. PullJob.Execute() receives error                 │
│  3. Worker marks job as Failed                       │
│  4. Error stored in job.err field                    │
│  5. TUI displays: "Pull failed: network timeout"     │
│  6. Partial results preserved (tickets pulled so far)│
└──────────────────────────────────────────────────────┘
```

### Cancellation During API Call

**Current Limitation:** PullService.Pull() doesn't support context cancellation internally. Cancellation happens at the goroutine boundary.

**Impact:** If user cancels during Jira API call, the call completes before cancellation takes effect. This typically takes 1-5 seconds.

**Future Enhancement:** Add context support to PullService for immediate cancellation.

### Partial Results

When a job is cancelled or fails, partial results are preserved:
- PullResult.TicketsPulled reflects actual tickets saved
- TicketTreeView refreshes to show new tickets
- User sees accurate count: "Cancelled: 45/120 tickets pulled"

## Performance Considerations

### Worker Count

**Initial Choice:** 1 worker

**Rationale:**
- Most Jira operations are I/O bound (network)
- Single worker provides:
  - Simpler reasoning about operation order
  - Lower resource usage
  - Easier debugging
  - No concurrent API calls (respects Jira rate limits)

**Future Scaling:** Increase to 2-3 workers if:
- Users frequently queue multiple operations
- Operations are proven independent and safe

### Channel Buffer Sizes

| Channel | Buffer | Rationale |
|---------|--------|-----------|
| jobChan | 10 | Rare to queue >10 operations |
| progressChan | 100 | High-frequency updates, prevent blocking |
| Job progress | 50 | Per-job updates, smooth UI |

### Memory Footprint

Estimated per active job:
- Job struct: ~200 bytes
- Context: ~100 bytes
- Progress channel: ~50 * 64 bytes = 3.2 KB
- Total: ~3.5 KB per job

**Acceptable:** Even 100 jobs = 350 KB (negligible)

### Goroutine Limits

- Workers: 1 (configurable up to 3)
- Progress forwarders: 1 per active job
- Maximum concurrent goroutines: ~10 (well within safe limits)

## Integration with Existing TUI

### App Initialization

```go
// In tui.NewTUIApp():
jobQueue := jobs.NewJobQueue(1)  // Single worker
tuiApp.jobQueue = jobQueue
```

### Pull Operation

```go
// In app.handlePull():
pullJob := jobs.NewPullJob(pullService, filePath, options)
jobID := jobQueue.Submit(pullJob)

// Track current job for cancellation
tuiApp.currentJobID = jobID

// Subscribe to progress
go func() {
    for progress := range jobQueue.Progress() {
        tuiApp.app.QueueUpdateDraw(func() {
            tuiApp.syncStatusView.SetProgress(progress)
        })
    }
}()
```

### Cancellation

```go
// In globalKeyHandler():
case tcell.KeyESC:
    if tuiApp.currentJobID != "" {
        jobQueue.Cancel(tuiApp.currentJobID)
        tuiApp.currentJobID = ""
    }
```

### Status Display

Progress updates shown in `syncStatusView`:
- "Connecting to Jira..." (0%, unknown total)
- "Querying project PROJ..." (0%, unknown)
- "Found 120 tickets" (0/120, 0%)
- "Pulling tickets: 45/120 (38%)" (45/120, 38%)
- "Pull complete: 120 tickets" (120/120, 100%)

## Testing Strategy

### Unit Tests

**queue_test.go:**
- TestJobQueue_Submit - Job execution
- TestJobQueue_Cancel - Cancellation works
- TestJobQueue_Progress - Progress events received
- TestJobQueue_MultipleJobs - Queue handles concurrent jobs
- TestJobQueue_Shutdown - Graceful shutdown
- TestJobQueue_StatusTracking - Status transitions correct

**Race Detection:**
```bash
go test -race ./internal/tui/jobs/...
```

### Fake Jobs for Testing

```go
type FakeJob struct {
    id           JobID
    progressChan chan JobProgress
    duration     time.Duration
    shouldFail   bool
}

func (f *FakeJob) Execute(ctx context.Context) error {
    for i := 0; i < 10; i++ {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(f.duration / 10):
            f.progressChan <- JobProgress{
                JobID:   f.id,
                Current: i + 1,
                Total:   10,
                Message: fmt.Sprintf("Step %d/10", i+1),
            }
        }
    }

    if f.shouldFail {
        return errors.New("fake job error")
    }
    return nil
}
```

### Integration Tests

Manual testing scenarios:
1. **Large Pull (500+ tickets)** - TUI remains responsive
2. **Rapid Cancellation** - Press ESC immediately after 'P'
3. **Network Failure** - Disconnect network mid-pull
4. **Multiple Operations** - Queue pull, then push
5. **Ctrl+C** - Graceful shutdown during operation

## Limitations and Future Work

### Current Limitations

1. **No Context in PullService**: Cancellation happens at goroutine level, not within PullService.Pull()
   - **Impact**: Cancellation during Jira API call takes 1-5 seconds
   - **Mitigation**: Document expected delay

2. **Single Operation Type**: Only PullJob implemented
   - **Impact**: Push operations still use old SyncCoordinator
   - **Mitigation**: Incremental migration planned

3. **No Retry Logic**: Failed jobs don't auto-retry
   - **Impact**: User must manually retry
   - **Mitigation**: Clear error messages guide user

4. **No Job History**: Completed jobs not persisted
   - **Impact**: Can't view past operation results after restart
   - **Mitigation**: Status shown until next operation

### Future Enhancements

**Week 3:**
- Add context support to PullService for immediate cancellation
- Implement PushJob wrapper
- Migrate all sync operations to job queue

**Week 4:**
- Add retry with exponential backoff
- Implement job history/logging
- Add operation queue view (see pending/active jobs)

**Phase 7:**
- Concurrent operations (multiple workspaces)
- Job priority levels
- Scheduled/deferred jobs

## Conclusion

This job queue architecture provides a robust foundation for async operations in Ticketr's TUI. The design prioritizes:
- **Simplicity**: Clean interfaces, straightforward implementation
- **Safety**: Thread-safe, no race conditions, no goroutine leaks
- **User Experience**: Responsive UI, real-time progress, cancellation
- **Maintainability**: Clear patterns, comprehensive tests, good documentation

The system is ready for immediate use with PullService and easily extensible for future async operations.

---

**Document Status:** Implementation Complete
**Last Updated:** 2025-10-20
**Author:** Builder Agent
**Reviewers:** Director, Verifier
