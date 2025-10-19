package search

import (
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantText     string
		wantAssignee string
		wantJiraID   string
		wantPriority string
		wantSprint   string
		wantRegex    string
	}{
		{
			name:     "empty query",
			input:    "",
			wantText: "",
		},
		{
			name:     "only text",
			input:    "authentication bug",
			wantText: "authentication bug",
		},
		{
			name:         "assignee filter",
			input:        "@john",
			wantAssignee: "john",
		},
		{
			name:       "JiraID filter",
			input:      "#BACK-123",
			wantJiraID: "BACK-123",
		},
		{
			name:       "partial JiraID filter",
			input:      "#BACK",
			wantJiraID: "BACK",
		},
		{
			name:         "priority filter",
			input:        "!high",
			wantPriority: "high",
		},
		{
			name:       "sprint filter",
			input:      "~sprint23",
			wantSprint: "sprint23",
		},
		{
			name:      "regex filter",
			input:     "/auth.*bug/",
			wantRegex: "auth.*bug",
		},
		{
			name:         "combined filters",
			input:        "@john !high authentication",
			wantText:     "authentication",
			wantAssignee: "john",
			wantPriority: "high",
		},
		{
			name:         "all filters",
			input:        "@john #BACK-123 !high ~sprint23 /auth.*/ login bug",
			wantText:     "login bug",
			wantAssignee: "john",
			wantJiraID:   "BACK-123",
			wantPriority: "high",
			wantSprint:   "sprint23",
			wantRegex:    "auth.*",
		},
		{
			name:     "completion filter ignored",
			input:    "%50 test",
			wantText: "test",
		},
		{
			name:     "malformed filters",
			input:    "@ # ! ~ test",
			wantText: "test",
		},
		{
			name:         "whitespace handling",
			input:        "  @john   !high   test  ",
			wantText:     "test",
			wantAssignee: "john",
			wantPriority: "high",
		},
		{
			name:         "case preservation for assignee/priority/sprint",
			input:        "@JohnDoe !High ~Sprint23",
			wantAssignee: "JohnDoe",
			wantPriority: "High",
			wantSprint:   "Sprint23",
		},
		{
			name:       "JiraID uppercased",
			input:      "#back-123",
			wantJiraID: "BACK-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseQuery(tt.input)

			if err != nil {
				t.Errorf("ParseQuery() unexpected error = %v", err)
				return
			}

			if got.Text != tt.wantText {
				t.Errorf("ParseQuery() Text = %v, want %v", got.Text, tt.wantText)
			}

			if got.Assignee != tt.wantAssignee {
				t.Errorf("ParseQuery() Assignee = %v, want %v", got.Assignee, tt.wantAssignee)
			}

			if got.JiraID != tt.wantJiraID {
				t.Errorf("ParseQuery() JiraID = %v, want %v", got.JiraID, tt.wantJiraID)
			}

			if got.Priority != tt.wantPriority {
				t.Errorf("ParseQuery() Priority = %v, want %v", got.Priority, tt.wantPriority)
			}

			if got.Sprint != tt.wantSprint {
				t.Errorf("ParseQuery() Sprint = %v, want %v", got.Sprint, tt.wantSprint)
			}

			if got.Regex != tt.wantRegex {
				t.Errorf("ParseQuery() Regex = %v, want %v", got.Regex, tt.wantRegex)
			}
		})
	}
}

func TestApplyFilters(t *testing.T) {
	baseTickets := []*domain.Ticket{
		{
			Title:  "Auth bug",
			JiraID: "BACK-123",
			CustomFields: map[string]string{
				"Assignee": "john.doe",
				"Priority": "High",
				"Sprint":   "Sprint 23",
			},
		},
		{
			Title:  "UI enhancement",
			JiraID: "FRONT-456",
			CustomFields: map[string]string{
				"Assignee": "jane.smith",
				"Priority": "Low",
				"Sprint":   "Sprint 24",
			},
		},
		{
			Title:  "Backend refactor",
			JiraID: "BACK-789",
			CustomFields: map[string]string{
				"Assignee": "john.doe",
				"Priority": "Medium",
				"Sprint":   "Sprint 23",
			},
		},
		{
			Title:  "Documentation update",
			JiraID: "DOC-111",
			CustomFields: map[string]string{
				"Assignee": "alice.jones",
				"Priority": "Low",
			},
		},
	}

	tests := []struct {
		name      string
		tickets   []*domain.Ticket
		query     *Query
		wantCount int
		wantIDs   []string
	}{
		{
			name:      "nil tickets",
			tickets:   nil,
			query:     &Query{Assignee: "john"},
			wantCount: 0,
		},
		{
			name:      "nil query",
			tickets:   baseTickets,
			query:     nil,
			wantCount: 4,
		},
		{
			name:    "filter by assignee",
			tickets: baseTickets,
			query: &Query{
				Assignee: "john",
			},
			wantCount: 2,
			wantIDs:   []string{"BACK-123", "BACK-789"},
		},
		{
			name:    "filter by JiraID exact",
			tickets: baseTickets,
			query: &Query{
				JiraID: "BACK-123",
			},
			wantCount: 1,
			wantIDs:   []string{"BACK-123"},
		},
		{
			name:    "filter by JiraID partial",
			tickets: baseTickets,
			query: &Query{
				JiraID: "BACK",
			},
			wantCount: 2,
			wantIDs:   []string{"BACK-123", "BACK-789"},
		},
		{
			name:    "filter by priority",
			tickets: baseTickets,
			query: &Query{
				Priority: "Low",
			},
			wantCount: 2,
			wantIDs:   []string{"FRONT-456", "DOC-111"},
		},
		{
			name:    "filter by sprint",
			tickets: baseTickets,
			query: &Query{
				Sprint: "Sprint 23",
			},
			wantCount: 2,
			wantIDs:   []string{"BACK-123", "BACK-789"},
		},
		{
			name:    "filter by regex",
			tickets: baseTickets,
			query: &Query{
				Regex: "auth.*bug",
			},
			wantCount: 1,
			wantIDs:   []string{"BACK-123"},
		},
		{
			name:    "combined filters - assignee and priority",
			tickets: baseTickets,
			query: &Query{
				Assignee: "john",
				Priority: "High",
			},
			wantCount: 1,
			wantIDs:   []string{"BACK-123"},
		},
		{
			name:    "combined filters - no match",
			tickets: baseTickets,
			query: &Query{
				Assignee: "john",
				Priority: "Low",
			},
			wantCount: 0,
		},
		{
			name:    "all filters combined",
			tickets: baseTickets,
			query: &Query{
				JiraID:   "BACK",
				Assignee: "john",
				Priority: "High",
				Sprint:   "Sprint 23",
			},
			wantCount: 1,
			wantIDs:   []string{"BACK-123"},
		},
		{
			name: "tickets with nil entries",
			tickets: []*domain.Ticket{
				baseTickets[0],
				nil,
				baseTickets[1],
			},
			query: &Query{
				Assignee: "john",
			},
			wantCount: 1,
			wantIDs:   []string{"BACK-123"},
		},
		{
			name:    "case insensitive assignee",
			tickets: baseTickets,
			query: &Query{
				Assignee: "JOHN",
			},
			wantCount: 2,
		},
		{
			name:    "invalid regex - returns all",
			tickets: baseTickets,
			query: &Query{
				Regex: "[invalid(",
			},
			wantCount: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplyFilters(tt.tickets, tt.query)

			if len(result) != tt.wantCount {
				t.Errorf("ApplyFilters() count = %v, want %v", len(result), tt.wantCount)
			}

			if tt.wantIDs != nil {
				gotIDs := make([]string, len(result))
				for i, ticket := range result {
					gotIDs[i] = ticket.JiraID
				}

				if len(gotIDs) != len(tt.wantIDs) {
					t.Errorf("ApplyFilters() IDs = %v, want %v", gotIDs, tt.wantIDs)
					return
				}

				for i, id := range tt.wantIDs {
					found := false
					for _, gotID := range gotIDs {
						if gotID == id {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("ApplyFilters() missing ID %v at position %v", id, i)
					}
				}
			}
		})
	}
}

func TestFilterByJiraID(t *testing.T) {
	tickets := []*domain.Ticket{
		{JiraID: "BACK-123"},
		{JiraID: "BACK-456"},
		{JiraID: "FRONT-123"},
		{JiraID: "back-789"}, // Test case insensitive
	}

	tests := []struct {
		name      string
		jiraID    string
		wantCount int
	}{
		{"exact match", "BACK-123", 1},
		{"partial match", "BACK", 3},
		{"case insensitive", "back-123", 1},
		{"no match", "NOTFOUND", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterByJiraID(tickets, tt.jiraID)
			if len(result) != tt.wantCount {
				t.Errorf("filterByJiraID() count = %v, want %v", len(result), tt.wantCount)
			}
		})
	}
}

func TestFilterByRegex(t *testing.T) {
	tickets := []*domain.Ticket{
		{Title: "Auth bug", Description: "Login fails"},
		{Title: "UI issue", Description: "Button alignment"},
		{Title: "Authentication system", Description: "OAuth2 implementation"},
	}

	tests := []struct {
		name      string
		pattern   string
		wantCount int
	}{
		{"simple pattern", "auth", 2},
		{"regex pattern", "auth.*bug", 1},
		{"case insensitive", "AUTH", 2},
		{"description match", "OAuth2", 1},
		{"invalid regex", "[invalid(", 3}, // Returns all on error
		{"no match", "nonexistent", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterByRegex(tickets, tt.pattern)
			if len(result) != tt.wantCount {
				t.Errorf("filterByRegex() count = %v, want %v", len(result), tt.wantCount)
			}
		})
	}
}

func BenchmarkParseQuery(b *testing.B) {
	input := "@john #BACK-123 !high ~sprint23 /auth.*/ login bug fix"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseQuery(input)
	}
}

func BenchmarkApplyFilters(b *testing.B) {
	// Create 1000 tickets
	tickets := make([]*domain.Ticket, 1000)
	for i := 0; i < 1000; i++ {
		tickets[i] = &domain.Ticket{
			Title:  "Ticket about various topics",
			JiraID: "TICK-" + string(rune(1000+i)),
			CustomFields: map[string]string{
				"Assignee": []string{"john", "jane", "alice"}[i%3],
				"Priority": []string{"Low", "Medium", "High"}[i%3],
				"Sprint":   "Sprint " + string(rune(20+(i%5))),
			},
		}
	}

	query := &Query{
		Assignee: "john",
		Priority: "High",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ApplyFilters(tickets, query)
	}
}

func BenchmarkFullSearch(b *testing.B) {
	// Simulate full search pipeline: parse -> filter -> search
	tickets := make([]*domain.Ticket, 1000)
	for i := 0; i < 1000; i++ {
		tickets[i] = &domain.Ticket{
			Title:       "Implement feature for authentication system",
			Description: "This ticket covers various aspects of the auth system",
			JiraID:      "TICK-" + string(rune(1000+i)),
			CustomFields: map[string]string{
				"Assignee": []string{"john", "jane", "alice"}[i%3],
				"Priority": []string{"Low", "Medium", "High"}[i%3],
				"Sprint":   "Sprint " + string(rune(20+(i%5))),
			},
		}
	}

	input := "@john !high authentication"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query, _ := ParseQuery(input)
		filtered := ApplyFilters(tickets, query)
		SearchTickets(filtered, query.Text)
	}
}
