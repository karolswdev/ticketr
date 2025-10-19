package search

import (
	"sort"
	"strings"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// Match represents a search match with relevance scoring
type Match struct {
	Ticket    *domain.Ticket
	Score     int      // 0-100, higher is better
	MatchedIn []string // Where match found: "title", "description", etc.
}

// FuzzyMatch performs fuzzy text matching on a ticket.
// Returns a Match with score and match locations, and a boolean indicating if there was any match.
// An empty query matches all tickets with a score of 50.
func FuzzyMatch(ticket *domain.Ticket, query string) (*Match, bool) {
	if ticket == nil {
		return nil, false
	}

	query = strings.ToLower(strings.TrimSpace(query))

	// Empty query matches everything with neutral score
	if query == "" {
		return &Match{
			Ticket:    ticket,
			Score:     50,
			MatchedIn: []string{},
		}, true
	}

	match := &Match{
		Ticket:    ticket,
		Score:     0,
		MatchedIn: []string{},
	}

	// Title matching (highest weight)
	titleLower := strings.ToLower(ticket.Title)
	if titleLower == query {
		match.Score = 100
		match.MatchedIn = append(match.MatchedIn, "title")
	} else if strings.Contains(titleLower, query) {
		match.Score = max(match.Score, 80)
		match.MatchedIn = append(match.MatchedIn, "title")
	}

	// JiraID matching (very high weight)
	jiraIDLower := strings.ToLower(ticket.JiraID)
	if jiraIDLower == query {
		match.Score = max(match.Score, 90)
		match.MatchedIn = append(match.MatchedIn, "jira_id")
	} else if strings.Contains(jiraIDLower, query) {
		match.Score = max(match.Score, 70)
		match.MatchedIn = append(match.MatchedIn, "jira_id")
	}

	// Description matching (medium weight)
	descLower := strings.ToLower(ticket.Description)
	if strings.Contains(descLower, query) {
		match.Score = max(match.Score, 60)
		match.MatchedIn = append(match.MatchedIn, "description")
	}

	// CustomFields matching (lower weight)
	for key, value := range ticket.CustomFields {
		valueLower := strings.ToLower(value)
		if strings.Contains(valueLower, query) {
			match.Score = max(match.Score, 40)
			match.MatchedIn = append(match.MatchedIn, "custom_field:"+key)
		}
	}

	// AcceptanceCriteria matching (lowest weight)
	for i, criterion := range ticket.AcceptanceCriteria {
		criterionLower := strings.ToLower(criterion)
		if strings.Contains(criterionLower, query) {
			match.Score = max(match.Score, 30)
			match.MatchedIn = append(match.MatchedIn, "acceptance_criteria:"+string(rune(i+1)))
		}
	}

	// If no match found, return false
	if match.Score == 0 {
		return nil, false
	}

	return match, true
}

// SearchTickets searches a list of tickets and returns matches sorted by relevance (descending).
// Returns empty slice if no matches found or input is nil/empty.
func SearchTickets(tickets []*domain.Ticket, query string) []*Match {
	if tickets == nil || len(tickets) == 0 {
		return []*Match{}
	}

	query = strings.TrimSpace(query)

	matches := make([]*Match, 0, len(tickets))

	for _, ticket := range tickets {
		if ticket == nil {
			continue
		}

		match, found := FuzzyMatch(ticket, query)
		if found {
			matches = append(matches, match)
		}
	}

	// Sort by score descending (highest first)
	sort.Slice(matches, func(i, j int) bool {
		// Primary sort: score descending
		if matches[i].Score != matches[j].Score {
			return matches[i].Score > matches[j].Score
		}

		// Secondary sort: title alphabetically (for stable ordering)
		return matches[i].Ticket.Title < matches[j].Ticket.Title
	})

	return matches
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
