package views

import (
	"fmt"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/rivo/tview"
)

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

	// Create root node for tree
	root := tview.NewTreeNode("Tickets")
	tree := tview.NewTreeView()
	tree.SetRoot(root)

	// Create minimal view struct for benchmarking buildTree
	view := &TicketTreeView{
		tree:            tree,
		root:            root,
		selectedTickets: make(map[string]bool),
		selectionMode:   false,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Build the tree (this is what we're benchmarking)
		view.buildTree(tickets)
	}
}

// BenchmarkTreeRefresh benchmarks the tree rebuild operation with selections.
func BenchmarkTreeRefresh(b *testing.B) {
	tickets := generateTestTickets(1000, 3)

	// Create root node for tree
	root := tview.NewTreeNode("Tickets")
	tree := tview.NewTreeView()
	tree.SetRoot(root)

	// Create minimal view struct for benchmarking
	view := &TicketTreeView{
		tree:            tree,
		root:            root,
		selectedTickets: make(map[string]bool),
		selectionMode:   false,
	}

	// Build initial tree
	view.buildTree(tickets)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Toggle selection and rebuild tree (simulate refresh without service calls)
		view.selectedTickets["PROJ-1"] = !view.selectedTickets["PROJ-1"]
		view.buildTree(tickets)
	}
}

// BenchmarkSelectionOperations benchmarks ticket selection operations.
func BenchmarkSelectionOperations(b *testing.B) {
	tickets := generateTestTickets(1000, 3)
	root := tview.NewTreeNode("Tickets")
	tree := tview.NewTreeView()
	tree.SetRoot(root)

	view := &TicketTreeView{
		tree:            tree,
		root:            root,
		selectedTickets: make(map[string]bool),
		selectionMode:   false,
	}

	view.buildTree(tickets)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Benchmark selection state changes without calling refreshTree
		for _, child := range root.GetChildren() {
			if ref := child.GetReference(); ref != nil {
				if ticket, ok := ref.(domain.Ticket); ok {
					view.selectedTickets[ticket.JiraID] = true
				}
			}
		}
		// Clear selection
		view.selectedTickets = make(map[string]bool)
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
