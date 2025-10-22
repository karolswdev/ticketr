package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// mockJiraAdapter provides a mock Jira adapter for benchmarking.
type mockJiraAdapterForBench struct {
	tickets     map[string]domain.Ticket
	searchDelay int // microseconds
	updateDelay int // microseconds
	callCount   int
	updateCount int
}

func newMockJiraAdapterForBench() *mockJiraAdapterForBench {
	return &mockJiraAdapterForBench{
		tickets: make(map[string]domain.Ticket),
	}
}

func (m *mockJiraAdapterForBench) Authenticate() error {
	return nil
}

func (m *mockJiraAdapterForBench) CreateTask(task domain.Task, parentID string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (m *mockJiraAdapterForBench) UpdateTask(task domain.Task) error {
	return fmt.Errorf("not implemented")
}

func (m *mockJiraAdapterForBench) GetProjectIssueTypes() (map[string][]string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockJiraAdapterForBench) GetIssueTypeFields(issueTypeName string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockJiraAdapterForBench) SearchTickets(ctx context.Context, project, jql string, progressCallback ports.JiraProgressCallback) ([]domain.Ticket, error) {
	m.callCount++

	// Simulate search by parsing JQL for key = "ID" pattern
	if len(jql) > 7 && jql[:6] == `key = ` {
		// Extract ticket ID from JQL (format: key = "PROJ-123")
		ticketID := jql[7 : len(jql)-1] // Remove 'key = "' and trailing '"'
		if ticket, exists := m.tickets[ticketID]; exists {
			return []domain.Ticket{ticket}, nil
		}
		return []domain.Ticket{}, nil
	}

	// Return all tickets if no specific query
	var result []domain.Ticket
	for _, ticket := range m.tickets {
		result = append(result, ticket)
	}
	return result, nil
}

func (m *mockJiraAdapterForBench) UpdateTicket(ticket domain.Ticket) error {
	m.callCount++
	m.updateCount++

	if _, exists := m.tickets[ticket.JiraID]; !exists {
		return fmt.Errorf("ticket not found: %s", ticket.JiraID)
	}

	m.tickets[ticket.JiraID] = ticket
	return nil
}

func (m *mockJiraAdapterForBench) CreateTicket(ticket domain.Ticket) (string, error) {
	m.callCount++
	m.tickets[ticket.JiraID] = ticket
	return ticket.JiraID, nil
}

// generateBulkTestTickets creates test tickets for bulk operations.
func generateBulkTestTickets(count int) map[string]domain.Ticket {
	tickets := make(map[string]domain.Ticket)

	for i := 0; i < count; i++ {
		ticketID := fmt.Sprintf("PROJ-%d", i+1)
		tickets[ticketID] = domain.Ticket{
			JiraID:      ticketID,
			Title:       fmt.Sprintf("Test Ticket %d", i+1),
			Description: fmt.Sprintf("Description for ticket %d", i+1),
			CustomFields: map[string]string{
				"Priority":  "Medium",
				"Component": "Backend",
				"Status":    "Open",
			},
		}
	}

	return tickets
}

// BenchmarkBulkUpdate10 benchmarks bulk update with 10 tickets.
func BenchmarkBulkUpdate10(b *testing.B) {
	benchmarkBulkUpdate(b, 10)
}

// BenchmarkBulkUpdate100 benchmarks bulk update with 100 tickets.
func BenchmarkBulkUpdate100(b *testing.B) {
	benchmarkBulkUpdate(b, 100)
}

// BenchmarkBulkUpdate1000 benchmarks bulk update with 1000 tickets.
func BenchmarkBulkUpdate1000(b *testing.B) {
	benchmarkBulkUpdate(b, 1000)
}

// benchmarkBulkUpdate is the core benchmark function for update operations.
func benchmarkBulkUpdate(b *testing.B, ticketCount int) {
	// Create test data
	testTickets := generateBulkTestTickets(ticketCount)

	// Get ticket IDs
	ticketIDs := make([]string, 0, ticketCount)
	for id := range testTickets {
		ticketIDs = append(ticketIDs, id)
	}

	// Create changes
	changes := map[string]interface{}{
		"Priority": "High",
		"Status":   "In Progress",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		adapter := newMockJiraAdapterForBench()
		adapter.tickets = make(map[string]domain.Ticket)
		for id, ticket := range testTickets {
			adapter.tickets[id] = ticket
		}

		service := NewBulkOperationService(adapter)

		op := &domain.BulkOperation{
			Action:    domain.BulkActionUpdate,
			TicketIDs: ticketIDs,
			Changes:   changes,
		}
		b.StartTimer()

		_, err := service.ExecuteOperation(context.Background(), op, nil)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}

		b.StopTimer()
		// Verify update count
		if adapter.updateCount != ticketCount {
			b.Fatalf("expected %d updates, got %d", ticketCount, adapter.updateCount)
		}
	}

	// Report API call metrics
	adapter := newMockJiraAdapterForBench()
	for id, ticket := range testTickets {
		adapter.tickets[id] = ticket
	}
	service := NewBulkOperationService(adapter)
	op := &domain.BulkOperation{
		Action:    domain.BulkActionUpdate,
		TicketIDs: ticketIDs,
		Changes:   changes,
	}
	service.ExecuteOperation(context.Background(), op, nil)

	b.ReportMetric(float64(adapter.callCount)/float64(ticketCount), "api-calls/ticket")
}

// BenchmarkBulkMove benchmarks bulk move operations.
func BenchmarkBulkMove(b *testing.B) {
	counts := []int{10, 50, 100}

	for _, count := range counts {
		b.Run(fmt.Sprintf("Count%d", count), func(b *testing.B) {
			testTickets := generateBulkTestTickets(count)
			ticketIDs := make([]string, 0, count)
			for id := range testTickets {
				ticketIDs = append(ticketIDs, id)
			}

			changes := map[string]interface{}{
				"parent": "PROJ-PARENT",
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				b.StopTimer()
				adapter := newMockJiraAdapterForBench()
				adapter.tickets = make(map[string]domain.Ticket)
				for id, ticket := range testTickets {
					adapter.tickets[id] = ticket
				}

				service := NewBulkOperationService(adapter)

				op := &domain.BulkOperation{
					Action:    domain.BulkActionMove,
					TicketIDs: ticketIDs,
					Changes:   changes,
				}
				b.StartTimer()

				_, err := service.ExecuteOperation(context.Background(), op, nil)
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// BenchmarkBulkOperationWithProgress benchmarks with progress callbacks.
func BenchmarkBulkOperationWithProgress(b *testing.B) {
	testTickets := generateBulkTestTickets(100)
	ticketIDs := make([]string, 0, 100)
	for id := range testTickets {
		ticketIDs = append(ticketIDs, id)
	}

	changes := map[string]interface{}{
		"Priority": "High",
	}

	progressCount := 0
	progressCallback := func(ticketID string, success bool, err error) {
		progressCount++
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		adapter := newMockJiraAdapterForBench()
		adapter.tickets = make(map[string]domain.Ticket)
		for id, ticket := range testTickets {
			adapter.tickets[id] = ticket
		}

		service := NewBulkOperationService(adapter)

		op := &domain.BulkOperation{
			Action:    domain.BulkActionUpdate,
			TicketIDs: ticketIDs,
			Changes:   changes,
		}

		progressCount = 0
		b.StartTimer()

		_, err := service.ExecuteOperation(context.Background(), op, progressCallback)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}

		b.StopTimer()
		if progressCount != 100 {
			b.Fatalf("expected 100 progress callbacks, got %d", progressCount)
		}
	}
}

// BenchmarkSnapshotCreation benchmarks the snapshot creation for rollback.
func BenchmarkSnapshotCreation(b *testing.B) {
	counts := []int{10, 100, 1000}

	for _, count := range counts {
		b.Run(fmt.Sprintf("Count%d", count), func(b *testing.B) {
			testTickets := generateBulkTestTickets(count)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				snapshots := make([]ticketSnapshot, 0, count)
				for id, ticket := range testTickets {
					snapshots = append(snapshots, ticketSnapshot{
						ticketID: id,
						ticket:   ticket,
					})
				}
				_ = snapshots
			}
		})
	}
}

// BenchmarkResultAggregation benchmarks result aggregation operations.
func BenchmarkResultAggregation(b *testing.B) {
	result := domain.NewBulkOperationResult()

	b.Run("RecordSuccess", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result.SuccessCount++
			result.SuccessfulTickets = append(result.SuccessfulTickets, fmt.Sprintf("PROJ-%d", i))
		}
	})

	b.Run("RecordFailure", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result.FailureCount++
			result.FailedTickets = append(result.FailedTickets, fmt.Sprintf("PROJ-%d", i))
			result.Errors[fmt.Sprintf("PROJ-%d", i)] = "test error"
		}
	})
}

// BenchmarkContextCancellation benchmarks context cancellation handling.
func BenchmarkContextCancellation(b *testing.B) {
	testTickets := generateBulkTestTickets(100)
	ticketIDs := make([]string, 0, 100)
	for id := range testTickets {
		ticketIDs = append(ticketIDs, id)
	}

	changes := map[string]interface{}{
		"Priority": "High",
	}

	b.Run("WithoutCancellation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			adapter := newMockJiraAdapterForBench()
			adapter.tickets = make(map[string]domain.Ticket)
			for id, ticket := range testTickets {
				adapter.tickets[id] = ticket
			}

			service := NewBulkOperationService(adapter)

			op := &domain.BulkOperation{
				Action:    domain.BulkActionUpdate,
				TicketIDs: ticketIDs,
				Changes:   changes,
			}
			b.StartTimer()

			ctx := context.Background()
			_, err := service.ExecuteOperation(ctx, op, nil)
			if err != nil {
				b.Fatalf("unexpected error: %v", err)
			}
		}
	})

	b.Run("WithCancellableContext", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			adapter := newMockJiraAdapterForBench()
			adapter.tickets = make(map[string]domain.Ticket)
			for id, ticket := range testTickets {
				adapter.tickets[id] = ticket
			}

			service := NewBulkOperationService(adapter)

			op := &domain.BulkOperation{
				Action:    domain.BulkActionUpdate,
				TicketIDs: ticketIDs,
				Changes:   changes,
			}
			b.StartTimer()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			_, err := service.ExecuteOperation(ctx, op, nil)
			if err != nil {
				b.Fatalf("unexpected error: %v", err)
			}
		}
	})
}
