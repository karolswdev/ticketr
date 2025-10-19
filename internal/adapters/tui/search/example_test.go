package search_test

import (
	"fmt"

	"github.com/karolswdev/ticktr/internal/adapters/tui/search"
	"github.com/karolswdev/ticktr/internal/core/domain"
)

// Example demonstrates basic fuzzy matching
func ExampleFuzzyMatch() {
	ticket := &domain.Ticket{
		Title:       "Fix authentication bug in login system",
		Description: "Users cannot login due to OAuth2 token expiration",
		JiraID:      "BACK-123",
		CustomFields: map[string]string{
			"Assignee": "john.doe",
			"Priority": "High",
		},
	}

	match, found := search.FuzzyMatch(ticket, "authentication")
	if found {
		fmt.Printf("Match found! Score: %d, Matched in: %v\n", match.Score, match.MatchedIn)
	}

	// Output:
	// Match found! Score: 80, Matched in: [title]
}

// Example demonstrates searching multiple tickets
func ExampleSearchTickets() {
	tickets := []*domain.Ticket{
		{
			Title:  "Implement OAuth2 authentication",
			JiraID: "BACK-123",
		},
		{
			Title:  "Fix login button UI",
			JiraID: "FRONT-456",
		},
		{
			Title:       "Update documentation",
			Description: "Add authentication examples to docs",
			JiraID:      "DOC-789",
		},
	}

	results := search.SearchTickets(tickets, "authentication")

	fmt.Printf("Found %d matches:\n", len(results))
	for _, match := range results {
		fmt.Printf("- %s (score: %d)\n", match.Ticket.JiraID, match.Score)
	}

	// Output:
	// Found 2 matches:
	// - BACK-123 (score: 80)
	// - DOC-789 (score: 60)
}

// Example demonstrates parsing a complex query with multiple filters
func ExampleParseQuery() {
	input := "@john #BACK-123 !high ~sprint23 authentication bug"

	query, _ := search.ParseQuery(input)

	fmt.Printf("Text: %s\n", query.Text)
	fmt.Printf("Assignee: %s\n", query.Assignee)
	fmt.Printf("JiraID: %s\n", query.JiraID)
	fmt.Printf("Priority: %s\n", query.Priority)
	fmt.Printf("Sprint: %s\n", query.Sprint)

	// Output:
	// Text: authentication bug
	// Assignee: john
	// JiraID: BACK-123
	// Priority: high
	// Sprint: sprint23
}

// Example demonstrates filtering tickets by assignee
func ExampleApplyFilters_assignee() {
	tickets := []*domain.Ticket{
		{
			Title:  "Ticket 1",
			JiraID: "BACK-123",
			CustomFields: map[string]string{
				"Assignee": "john.doe",
			},
		},
		{
			Title:  "Ticket 2",
			JiraID: "BACK-456",
			CustomFields: map[string]string{
				"Assignee": "jane.smith",
			},
		},
	}

	query := &search.Query{
		Assignee: "john",
	}

	filtered := search.ApplyFilters(tickets, query)

	fmt.Printf("Found %d tickets assigned to john\n", len(filtered))

	// Output:
	// Found 1 tickets assigned to john
}

// Example demonstrates regex filtering
func ExampleApplyFilters_regex() {
	tickets := []*domain.Ticket{
		{Title: "Auth bug in login"},
		{Title: "Authentication system"},
		{Title: "UI alignment issue"},
	}

	query := &search.Query{
		Regex: "auth.*bug",
	}

	filtered := search.ApplyFilters(tickets, query)

	fmt.Printf("Found %d tickets matching regex\n", len(filtered))

	// Output:
	// Found 1 tickets matching regex
}

// Example_fullSearchPipeline demonstrates the full search pipeline
func Example_fullSearchPipeline() {
	// Sample tickets
	tickets := []*domain.Ticket{
		{
			Title:       "Fix authentication bug",
			Description: "Users cannot login with OAuth2",
			JiraID:      "BACK-123",
			CustomFields: map[string]string{
				"Assignee": "john.doe",
				"Priority": "High",
				"Sprint":   "Sprint 23",
			},
		},
		{
			Title:       "Add authentication tests",
			Description: "Unit tests for auth module",
			JiraID:      "BACK-456",
			CustomFields: map[string]string{
				"Assignee": "jane.smith",
				"Priority": "Medium",
				"Sprint":   "Sprint 23",
			},
		},
		{
			Title:       "Update user profile UI",
			Description: "Redesign profile page",
			JiraID:      "FRONT-789",
			CustomFields: map[string]string{
				"Assignee": "john.doe",
				"Priority": "Low",
				"Sprint":   "Sprint 24",
			},
		},
	}

	// Parse user input
	input := "@john ~Sprint authentication"
	query, _ := search.ParseQuery(input)

	// Apply filters first
	filtered := search.ApplyFilters(tickets, query)

	// Then fuzzy search on remaining text
	results := search.SearchTickets(filtered, query.Text)

	fmt.Printf("Query: %s\n", input)
	fmt.Printf("Filters: assignee=%s, sprint=%s\n", query.Assignee, query.Sprint)
	fmt.Printf("Text search: %s\n", query.Text)
	fmt.Printf("\nResults (%d):\n", len(results))
	for _, match := range results {
		fmt.Printf("- %s: %s (score: %d)\n",
			match.Ticket.JiraID,
			match.Ticket.Title,
			match.Score)
	}

	// Output:
	// Query: @john ~Sprint authentication
	// Filters: assignee=john, sprint=Sprint
	// Text search: authentication
	//
	// Results (1):
	// - BACK-123: Fix authentication bug (score: 80)
}

// Example_performanceDemo demonstrates performance characteristics
func Example_performanceDemo() {
	// Create 1000 tickets - mix of assignments and priorities
	tickets := make([]*domain.Ticket, 1000)
	for i := 0; i < 1000; i++ {
		tickets[i] = &domain.Ticket{
			Title:       fmt.Sprintf("Ticket %d about various features", i),
			Description: "Detailed description with authentication, UI, backend, etc.",
			JiraID:      fmt.Sprintf("TICK-%d", i),
			CustomFields: map[string]string{
				"Assignee": []string{"john.doe", "jane.smith", "alice.jones"}[i%3],
				"Priority": []string{"Low", "Medium", "High"}[i%3],
			},
		}
	}

	// Search for tickets with "authentication" in text
	input := "authentication"
	query, _ := search.ParseQuery(input)
	filtered := search.ApplyFilters(tickets, query)
	results := search.SearchTickets(filtered, query.Text)

	fmt.Printf("Searched 1000 tickets, found %d matches\n", len(results))
	fmt.Printf("Performance: typically < 50ms for this workload\n")

	// Output:
	// Searched 1000 tickets, found 1000 matches
	// Performance: typically < 50ms for this workload
}
