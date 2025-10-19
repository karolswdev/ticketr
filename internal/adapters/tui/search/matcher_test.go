package search

import (
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

func TestFuzzyMatch(t *testing.T) {
	tests := []struct {
		name          string
		ticket        *domain.Ticket
		query         string
		wantMatch     bool
		wantMinScore  int
		wantMaxScore  int
		wantMatchedIn []string
	}{
		{
			name:         "nil ticket",
			ticket:       nil,
			query:        "test",
			wantMatch:    false,
			wantMinScore: 0,
			wantMaxScore: 0,
		},
		{
			name: "empty query matches all",
			ticket: &domain.Ticket{
				Title: "Test Ticket",
			},
			query:        "",
			wantMatch:    true,
			wantMinScore: 50,
			wantMaxScore: 50,
		},
		{
			name: "exact title match",
			ticket: &domain.Ticket{
				Title: "authentication bug",
			},
			query:         "authentication bug",
			wantMatch:     true,
			wantMinScore:  100,
			wantMaxScore:  100,
			wantMatchedIn: []string{"title"},
		},
		{
			name: "partial title match",
			ticket: &domain.Ticket{
				Title: "Fix authentication bug in login",
			},
			query:         "authentication",
			wantMatch:     true,
			wantMinScore:  80,
			wantMaxScore:  80,
			wantMatchedIn: []string{"title"},
		},
		{
			name: "exact JiraID match",
			ticket: &domain.Ticket{
				Title:  "Some ticket",
				JiraID: "BACK-123",
			},
			query:         "back-123",
			wantMatch:     true,
			wantMinScore:  90,
			wantMaxScore:  90,
			wantMatchedIn: []string{"jira_id"},
		},
		{
			name: "partial JiraID match",
			ticket: &domain.Ticket{
				Title:  "Some ticket",
				JiraID: "BACK-123",
			},
			query:         "back",
			wantMatch:     true,
			wantMinScore:  70,
			wantMaxScore:  70,
			wantMatchedIn: []string{"jira_id"},
		},
		{
			name: "description match",
			ticket: &domain.Ticket{
				Title:       "Short title",
				Description: "This is a detailed description about authentication",
			},
			query:         "authentication",
			wantMatch:     true,
			wantMinScore:  60,
			wantMaxScore:  60,
			wantMatchedIn: []string{"description"},
		},
		{
			name: "custom field match",
			ticket: &domain.Ticket{
				Title: "Some ticket",
				CustomFields: map[string]string{
					"Assignee": "john.doe",
					"Priority": "High",
				},
			},
			query:         "john",
			wantMatch:     true,
			wantMinScore:  40,
			wantMaxScore:  40,
			wantMatchedIn: []string{"custom_field:Assignee"},
		},
		{
			name: "acceptance criteria match",
			ticket: &domain.Ticket{
				Title: "Some ticket",
				AcceptanceCriteria: []string{
					"User can login",
					"Session persists",
				},
			},
			query:        "session",
			wantMatch:    true,
			wantMinScore: 30,
			wantMaxScore: 30,
		},
		{
			name: "multiple matches - highest score wins",
			ticket: &domain.Ticket{
				Title:       "authentication system",
				Description: "Fix authentication bug",
				JiraID:      "AUTH-123",
			},
			query:        "authentication",
			wantMatch:    true,
			wantMinScore: 80,  // Title partial match
			wantMaxScore: 100, // Could be exact if normalized
		},
		{
			name: "no match",
			ticket: &domain.Ticket{
				Title:       "User interface",
				Description: "Update button colors",
			},
			query:     "authentication",
			wantMatch: false,
		},
		{
			name: "case insensitive",
			ticket: &domain.Ticket{
				Title: "Authentication Bug",
			},
			query:         "AUTHENTICATION",
			wantMatch:     true,
			wantMinScore:  80,
			wantMaxScore:  80,
			wantMatchedIn: []string{"title"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, found := FuzzyMatch(tt.ticket, tt.query)

			if found != tt.wantMatch {
				t.Errorf("FuzzyMatch() found = %v, want %v", found, tt.wantMatch)
			}

			if !found {
				return
			}

			if match.Score < tt.wantMinScore || match.Score > tt.wantMaxScore {
				t.Errorf("FuzzyMatch() score = %v, want between %v and %v",
					match.Score, tt.wantMinScore, tt.wantMaxScore)
			}

			if tt.wantMatchedIn != nil {
				if len(match.MatchedIn) != len(tt.wantMatchedIn) {
					t.Errorf("FuzzyMatch() matchedIn count = %v, want %v",
						len(match.MatchedIn), len(tt.wantMatchedIn))
				}
			}
		})
	}
}

func TestSearchTickets(t *testing.T) {
	tickets := []*domain.Ticket{
		{
			Title:  "Exact Match",
			JiraID: "EXACT-1",
		},
		{
			Title:  "Partial authentication match",
			JiraID: "PART-1",
		},
		{
			Title:       "No match title",
			Description: "But has authentication in description",
			JiraID:      "DESC-1",
		},
		{
			Title:  "Unrelated ticket",
			JiraID: "OTHER-1",
		},
	}

	tests := []struct {
		name        string
		tickets     []*domain.Ticket
		query       string
		wantCount   int
		wantFirstID string
		wantLastID  string
		wantSorted  bool
		wantNonNil  bool
	}{
		{
			name:       "nil tickets",
			tickets:    nil,
			query:      "test",
			wantCount:  0,
			wantNonNil: true,
		},
		{
			name:       "empty tickets",
			tickets:    []*domain.Ticket{},
			query:      "test",
			wantCount:  0,
			wantNonNil: true,
		},
		{
			name:        "search with results",
			tickets:     tickets,
			query:       "authentication",
			wantCount:   2,
			wantFirstID: "PART-1", // Title match (80) > Description match (60)
			wantLastID:  "DESC-1",
			wantSorted:  true,
		},
		{
			name:       "no results",
			tickets:    tickets,
			query:      "nonexistent",
			wantCount:  0,
			wantNonNil: true,
		},
		{
			name:      "empty query matches all",
			tickets:   tickets,
			query:     "",
			wantCount: 4,
		},
		{
			name: "tickets with nil entries",
			tickets: []*domain.Ticket{
				{Title: "Valid", JiraID: "V-1"},
				nil,
				{Title: "Also Valid", JiraID: "V-2"},
			},
			query:     "valid",
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := SearchTickets(tt.tickets, tt.query)

			if tt.wantNonNil && results == nil {
				t.Error("SearchTickets() returned nil, want non-nil slice")
			}

			if len(results) != tt.wantCount {
				t.Errorf("SearchTickets() count = %v, want %v", len(results), tt.wantCount)
			}

			if tt.wantCount > 0 {
				if tt.wantFirstID != "" && results[0].Ticket.JiraID != tt.wantFirstID {
					t.Errorf("SearchTickets() first ID = %v, want %v",
						results[0].Ticket.JiraID, tt.wantFirstID)
				}

				if tt.wantLastID != "" && results[len(results)-1].Ticket.JiraID != tt.wantLastID {
					t.Errorf("SearchTickets() last ID = %v, want %v",
						results[len(results)-1].Ticket.JiraID, tt.wantLastID)
				}
			}

			// Verify sorting (descending by score)
			if tt.wantSorted && len(results) > 1 {
				for i := 0; i < len(results)-1; i++ {
					if results[i].Score < results[i+1].Score {
						t.Errorf("SearchTickets() not sorted: results[%d].Score = %v < results[%d].Score = %v",
							i, results[i].Score, i+1, results[i+1].Score)
					}
				}
			}
		})
	}
}

func BenchmarkFuzzyMatch(b *testing.B) {
	ticket := &domain.Ticket{
		Title:       "Implement authentication system with OAuth2",
		Description: "Create a new authentication system using OAuth2 protocol for secure user login",
		JiraID:      "BACK-123",
		CustomFields: map[string]string{
			"Assignee": "john.doe",
			"Priority": "High",
			"Sprint":   "Sprint 23",
		},
		AcceptanceCriteria: []string{
			"User can login with OAuth2",
			"Tokens are stored securely",
			"Session management works correctly",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FuzzyMatch(ticket, "authentication")
	}
}

func BenchmarkSearchTickets(b *testing.B) {
	// Create 1000 tickets
	tickets := make([]*domain.Ticket, 1000)
	for i := 0; i < 1000; i++ {
		tickets[i] = &domain.Ticket{
			Title:       "Ticket about various topics including auth, ui, backend",
			Description: "This is a detailed description with many words",
			JiraID:      "TICK-" + string(rune(i)),
			CustomFields: map[string]string{
				"Assignee": "user" + string(rune(i%10)),
				"Priority": []string{"Low", "Medium", "High"}[i%3],
			},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SearchTickets(tickets, "auth")
	}
}
