package views

import (
	"fmt"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/rivo/tview"
)

// mockWorkspaceService provides a mock WorkspaceService for benchmarking.
type mockWorkspaceService struct {
	workspace *domain.Workspace
}

func (m *mockWorkspaceService) Current() (*domain.Workspace, error) {
	return m.workspace, nil
}

func (m *mockWorkspaceService) Create(name, jiraURL, projectKey string) error {
	return nil
}

func (m *mockWorkspaceService) List() ([]domain.Workspace, error) {
	return []domain.Workspace{*m.workspace}, nil
}

func (m *mockWorkspaceService) SetDefault(workspaceID string) error {
	return nil
}

func (m *mockWorkspaceService) Get(workspaceID string) (*domain.Workspace, error) {
	return m.workspace, nil
}

func (m *mockWorkspaceService) Delete(workspaceID string) error {
	return nil
}

func (m *mockWorkspaceService) Update(workspaceID, jiraURL, projectKey string) error {
	return nil
}

// mockTicketQueryService provides a mock TicketQueryService for benchmarking.
type mockTicketQueryService struct {
	tickets []domain.Ticket
}

func (m *mockTicketQueryService) ListByWorkspace(workspaceID string) ([]domain.Ticket, error) {
	return m.tickets, nil
}

func (m *mockTicketQueryService) Get(ticketID string) (*domain.Ticket, error) {
	for i := range m.tickets {
		if m.tickets[i].JiraID == ticketID {
			return &m.tickets[i], nil
		}
	}
	return nil, fmt.Errorf("ticket not found")
}

// generateTestTickets creates a slice of test tickets with varying complexity.
func generateTestTickets(count int, tasksPerTicket int) []domain.Ticket {
	tickets := make([]domain.Ticket, count)

	for i := 0; i < count; i++ {
		ticketID := fmt.Sprintf("PROJ-%d", i+1)

		// Generate tasks for this ticket
		tasks := make([]domain.Task, tasksPerTicket)
		for j := 0; j < tasksPerTicket; j++ {
			tasks[j] = domain.Task{
				JiraID: fmt.Sprintf("%s-TASK-%d", ticketID, j+1),
				Title:  fmt.Sprintf("Task %d for ticket %s", j+1, ticketID),
			}
		}

		tickets[i] = domain.Ticket{
			JiraID:      ticketID,
			Title:       fmt.Sprintf("Test Ticket %d - This is a longer title to simulate real ticket titles", i+1),
			Description: fmt.Sprintf("Description for ticket %d with some content", i+1),
			Tasks:       tasks,
			CustomFields: map[string]string{
				"Priority":  "High",
				"Component": fmt.Sprintf("Component-%d", i%5),
				"Sprint":    fmt.Sprintf("Sprint %d", i%10),
			},
		}
	}

	return tickets
}

// BenchmarkTreeRendering100 benchmarks tree rendering with 100 tickets.
func BenchmarkTreeRendering100(b *testing.B) {
	benchmarkTreeRendering(b, 100, 3)
}

// BenchmarkTreeRendering1000 benchmarks tree rendering with 1000 tickets.
func BenchmarkTreeRendering1000(b *testing.B) {
	benchmarkTreeRendering(b, 1000, 3)
}

// BenchmarkTreeRendering10000 benchmarks tree rendering with 10000 tickets.
func BenchmarkTreeRendering10000(b *testing.B) {
	benchmarkTreeRendering(b, 10000, 3)
}

// BenchmarkTreeRenderingWithTasks benchmarks tree rendering with many tasks per ticket.
func BenchmarkTreeRenderingWithTasks(b *testing.B) {
	benchmarkTreeRendering(b, 1000, 10)
}

// benchmarkTreeRendering is the core benchmark function.
func benchmarkTreeRendering(b *testing.B, ticketCount, tasksPerTicket int) {
	// Create test data
	tickets := generateTestTickets(ticketCount, tasksPerTicket)

	// Setup mocks
	workspace := &domain.Workspace{
		ID:   "test-workspace",
		Name: "Test Workspace",
	}

	mockWS := &mockWorkspaceService{workspace: workspace}
	mockTQ := &mockTicketQueryService{tickets: tickets}

	// Create a test app
	app := tview.NewApplication()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Create new view for each iteration to avoid caching effects
		view := NewTicketTreeView(mockWS, mockTQ, app)

		// Build the tree (this is what we're benchmarking)
		view.buildTree(tickets)
	}
}

// BenchmarkTreeRefresh benchmarks the refreshTree operation.
func BenchmarkTreeRefresh(b *testing.B) {
	tickets := generateTestTickets(1000, 3)

	workspace := &domain.Workspace{
		ID:   "test-workspace",
		Name: "Test Workspace",
	}

	mockWS := &mockWorkspaceService{workspace: workspace}
	mockTQ := &mockTicketQueryService{tickets: tickets}
	app := tview.NewApplication()

	view := NewTicketTreeView(mockWS, mockTQ, app)
	view.buildTree(tickets)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Toggle selection to trigger refresh
		view.selectedTickets["PROJ-1"] = !view.selectedTickets["PROJ-1"]
		view.refreshTree()
	}
}

// BenchmarkSelectionOperations benchmarks ticket selection operations.
func BenchmarkSelectionOperations(b *testing.B) {
	tickets := generateTestTickets(1000, 3)

	workspace := &domain.Workspace{
		ID:   "test-workspace",
		Name: "Test Workspace",
	}

	mockWS := &mockWorkspaceService{workspace: workspace}
	mockTQ := &mockTicketQueryService{tickets: tickets}
	app := tview.NewApplication()

	view := NewTicketTreeView(mockWS, mockTQ, app)
	view.buildTree(tickets)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		view.selectAllVisible()
		view.clearSelection()
	}
}

// BenchmarkStringConcatenation benchmarks the string building for ticket nodes.
func BenchmarkStringConcatenation(b *testing.B) {
	ticket := domain.Ticket{
		JiraID: "PROJ-1234",
		Title:  "This is a very long title that might exceed the maximum length and need truncation to fit properly",
	}

	b.Run("CurrentImplementation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			checkbox := "[ ] "
			ticketText := checkbox + ticket.JiraID
			if ticket.Title != "" {
				ticketText += ": " + ticket.Title
			}

			maxLen := 60
			if len(ticketText) > maxLen {
				ticketText = ticketText[:maxLen-3] + "..."
			}
			_ = ticketText
		}
	})

	b.Run("StringBuilderOptimized", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Preallocate with estimated capacity
			var builder []byte
			capacity := len("[ ] ") + len(ticket.JiraID) + 2 + len(ticket.Title)
			if capacity > 60 {
				capacity = 60
			}
			builder = make([]byte, 0, capacity)

			builder = append(builder, "[ ] "...)
			builder = append(builder, ticket.JiraID...)
			if ticket.Title != "" {
				builder = append(builder, ": "...)
				builder = append(builder, ticket.Title...)
			}

			result := string(builder)
			maxLen := 60
			if len(result) > maxLen {
				result = result[:maxLen-3] + "..."
			}
			_ = result
		}
	})
}
