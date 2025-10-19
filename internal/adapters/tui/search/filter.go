package search

import (
	"regexp"
	"strings"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// Query represents a parsed search query with filters and free text
type Query struct {
	Text     string // Free text to fuzzy match
	Assignee string // @username filter
	JiraID   string // #ID filter
	Priority string // !priority filter
	Sprint   string // ~sprint filter
	Regex    string // /pattern/ filter (raw pattern without slashes)
}

// ParseQuery parses a search string into structured query components.
// Supports:
//   - @username - assignee filter
//   - #JIRA-123 - JiraID filter (can be partial like #BACK)
//   - !priority - priority filter
//   - ~sprint - sprint filter
//   - /pattern/ - regex filter
//   - everything else - free text query
//
// Returns a Query struct. Never returns error - malformed input is handled gracefully.
func ParseQuery(input string) (*Query, error) {
	input = strings.TrimSpace(input)

	query := &Query{}

	// Extract regex pattern first (so it doesn't interfere with other parsing)
	regexPattern := regexp.MustCompile(`/([^/]+)/`)
	if match := regexPattern.FindStringSubmatch(input); len(match) > 1 {
		query.Regex = match[1]
		// Remove regex from input for further processing
		input = regexPattern.ReplaceAllString(input, " ")
	}

	// Split into tokens for parsing
	tokens := strings.Fields(input)
	textParts := []string{}

	for _, token := range tokens {
		if len(token) == 0 {
			continue
		}

		switch token[0] {
		case '@':
			// Assignee filter
			if len(token) > 1 {
				query.Assignee = token[1:]
			}
		case '#':
			// JiraID filter
			if len(token) > 1 {
				query.JiraID = strings.ToUpper(token[1:])
			}
		case '!':
			// Priority filter
			if len(token) > 1 {
				query.Priority = token[1:]
			}
		case '~':
			// Sprint filter
			if len(token) > 1 {
				query.Sprint = token[1:]
			}
		case '%':
			// Completion filter - not implemented yet, skip token
			continue
		default:
			// Regular text
			textParts = append(textParts, token)
		}
	}

	// Join remaining text parts
	query.Text = strings.Join(textParts, " ")

	return query, nil
}

// ApplyFilters filters tickets based on parsed query filters.
// Applies filters in order: JiraID, Assignee, Priority, Sprint, Regex.
// Returns filtered slice. Returns empty slice if input is nil.
func ApplyFilters(tickets []*domain.Ticket, query *Query) []*domain.Ticket {
	if tickets == nil {
		return []*domain.Ticket{}
	}

	if query == nil {
		return tickets
	}

	filtered := tickets

	// Apply JiraID filter
	if query.JiraID != "" {
		filtered = filterByJiraID(filtered, query.JiraID)
	}

	// Apply Assignee filter
	if query.Assignee != "" {
		filtered = filterByAssignee(filtered, query.Assignee)
	}

	// Apply Priority filter
	if query.Priority != "" {
		filtered = filterByPriority(filtered, query.Priority)
	}

	// Apply Sprint filter
	if query.Sprint != "" {
		filtered = filterBySprint(filtered, query.Sprint)
	}

	// Apply Regex filter
	if query.Regex != "" {
		filtered = filterByRegex(filtered, query.Regex)
	}

	return filtered
}

// filterByJiraID filters tickets by JiraID (case-insensitive partial match)
func filterByJiraID(tickets []*domain.Ticket, jiraID string) []*domain.Ticket {
	jiraID = strings.ToUpper(jiraID)
	result := make([]*domain.Ticket, 0)

	for _, ticket := range tickets {
		if ticket == nil {
			continue
		}
		if strings.Contains(strings.ToUpper(ticket.JiraID), jiraID) {
			result = append(result, ticket)
		}
	}

	return result
}

// filterByAssignee filters tickets by Assignee custom field (case-insensitive partial match)
func filterByAssignee(tickets []*domain.Ticket, assignee string) []*domain.Ticket {
	assigneeLower := strings.ToLower(assignee)
	result := make([]*domain.Ticket, 0)

	for _, ticket := range tickets {
		if ticket == nil {
			continue
		}

		if ticketAssignee, exists := ticket.CustomFields["Assignee"]; exists {
			if strings.Contains(strings.ToLower(ticketAssignee), assigneeLower) {
				result = append(result, ticket)
			}
		}
	}

	return result
}

// filterByPriority filters tickets by Priority custom field (case-insensitive partial match)
func filterByPriority(tickets []*domain.Ticket, priority string) []*domain.Ticket {
	priorityLower := strings.ToLower(priority)
	result := make([]*domain.Ticket, 0)

	for _, ticket := range tickets {
		if ticket == nil {
			continue
		}

		if ticketPriority, exists := ticket.CustomFields["Priority"]; exists {
			if strings.Contains(strings.ToLower(ticketPriority), priorityLower) {
				result = append(result, ticket)
			}
		}
	}

	return result
}

// filterBySprint filters tickets by Sprint custom field (case-insensitive partial match)
func filterBySprint(tickets []*domain.Ticket, sprint string) []*domain.Ticket {
	sprintLower := strings.ToLower(sprint)
	result := make([]*domain.Ticket, 0)

	for _, ticket := range tickets {
		if ticket == nil {
			continue
		}

		if ticketSprint, exists := ticket.CustomFields["Sprint"]; exists {
			if strings.Contains(strings.ToLower(ticketSprint), sprintLower) {
				result = append(result, ticket)
			}
		}
	}

	return result
}

// filterByRegex filters tickets by regex pattern matching Title and Description.
// Case-insensitive by default. Invalid regex patterns are ignored (returns all tickets).
func filterByRegex(tickets []*domain.Ticket, pattern string) []*domain.Ticket {
	// Compile regex with case-insensitive flag
	re, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		// Invalid regex - return all tickets unchanged
		return tickets
	}

	result := make([]*domain.Ticket, 0)

	for _, ticket := range tickets {
		if ticket == nil {
			continue
		}

		// Match against Title or Description
		if re.MatchString(ticket.Title) || re.MatchString(ticket.Description) {
			result = append(result, ticket)
		}
	}

	return result
}
