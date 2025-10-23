package tree

import (
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// BenchmarkFlattenTickets1000 benchmarks flattening with 1,000 items
func BenchmarkFlattenTickets1000(b *testing.B) {
	tickets := generateBenchmarkTickets(100, 10) // 100 tickets × 10 tasks = 1,000 items
	expandedState := make(map[string]bool)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FlattenTickets(tickets, expandedState)
	}
}

// BenchmarkFlattenTickets10000 benchmarks flattening with 10,000 items
func BenchmarkFlattenTickets10000(b *testing.B) {
	tickets := generateBenchmarkTickets(1000, 10) // 1,000 tickets × 10 tasks = 10,000 items
	expandedState := make(map[string]bool)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FlattenTickets(tickets, expandedState)
	}
}

// BenchmarkFlattenTicketsExpanded1000 benchmarks with all items expanded
func BenchmarkFlattenTicketsExpanded1000(b *testing.B) {
	tickets := generateBenchmarkTickets(100, 10)
	expandedState := make(map[string]bool)
	for i := range tickets {
		expandedState[tickets[i].JiraID] = true
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FlattenTickets(tickets, expandedState)
	}
}

// BenchmarkTreeRebuild benchmarks tree rebuild operations
func BenchmarkTreeRebuild(b *testing.B) {
	tickets := generateBenchmarkTickets(100, 10)
	m := New(80, 24, &theme.DefaultTheme)
	m.SetTickets(tickets)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.rebuildTree()
	}
}

// BenchmarkTreeViewRender benchmarks rendering the tree view
func BenchmarkTreeViewRender(b *testing.B) {
	tickets := generateBenchmarkTickets(100, 10)
	m := New(80, 24, &theme.DefaultTheme)
	m.SetTickets(tickets)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.View()
	}
}

// generateBenchmarkTickets creates test tickets for benchmarking
func generateBenchmarkTickets(numTickets, tasksPerTicket int) []domain.Ticket {
	tickets := make([]domain.Ticket, numTickets)
	for i := 0; i < numTickets; i++ {
		tasks := make([]domain.Task, tasksPerTicket)
		for j := 0; j < tasksPerTicket; j++ {
			tasks[j] = domain.Task{
				JiraID:      "BENCH-TASK",
				Title:       "Benchmark Task for Performance Testing",
				Description: "This is a benchmark task used for performance testing of the tree component rendering and flattening operations.",
				CustomFields: map[string]string{
					"Priority": "Medium",
					"Status":   "In Progress",
				},
			}
		}
		tickets[i] = domain.Ticket{
			JiraID:      "BENCH-TICKET",
			Title:       "Benchmark Ticket for Performance Testing",
			Description: "This is a benchmark ticket used for performance testing.",
			CustomFields: map[string]string{
				"Priority": "High",
				"Status":   "Open",
			},
			Tasks: tasks,
		}
	}
	return tickets
}
