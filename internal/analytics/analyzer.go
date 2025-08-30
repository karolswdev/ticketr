// Package analytics provides ticket analysis and statistics functionality.
// It processes ticket data to generate insights about project status,
// velocity, and work distribution.
package analytics

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/karolswdev/ticketr/internal/core/domain"
)

// Statistics holds the calculated metrics for a set of tickets
type Statistics struct {
	TotalTickets            int
	TotalTasks              int
	TicketsByType           map[string]int
	TicketsByStatus         map[string]int
	TasksByStatus           map[string]int
	TotalStoryPoints        float64
	TicketsWithJiraID       int
	TasksWithJiraID         int
	AcceptanceCriteriaCount int
}

// Analyzer processes tickets to generate statistics
type Analyzer struct{}

// NewAnalyzer creates a new analyzer instance
//
// Returns:
//   - *Analyzer: A new analyzer ready to process tickets
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// AnalyzeTickets processes a slice of tickets and generates statistics
//
// Parameters:
//   - tickets: The tickets to analyze
//
// Returns:
//   - *Statistics: Calculated statistics for the tickets
func (a *Analyzer) AnalyzeTickets(tickets []domain.Ticket) *Statistics {
	stats := &Statistics{
		TicketsByType:   make(map[string]int),
		TicketsByStatus: make(map[string]int),
		TasksByStatus:   make(map[string]int),
	}

	for _, ticket := range tickets {
		stats.TotalTickets++

		// Count tickets with JIRA IDs
		if ticket.JiraID != "" {
			stats.TicketsWithJiraID++
		}

		// Categorize by type
		ticketType := a.extractType(ticket)
		stats.TicketsByType[ticketType]++

		// Categorize by status
		status := a.extractStatus(ticket)
		stats.TicketsByStatus[status]++

		// Count acceptance criteria
		stats.AcceptanceCriteriaCount += len(ticket.AcceptanceCriteria)

		// Extract story points
		if points := a.extractStoryPoints(ticket); points > 0 {
			stats.TotalStoryPoints += points
		}

		// Process tasks
		for _, task := range ticket.Tasks {
			stats.TotalTasks++

			if task.JiraID != "" {
				stats.TasksWithJiraID++
			}

			taskStatus := a.extractTaskStatus(task)
			stats.TasksByStatus[taskStatus]++

			// Count task acceptance criteria
			stats.AcceptanceCriteriaCount += len(task.AcceptanceCriteria)
		}
	}

	return stats
}

// extractType determines the ticket type from title or custom fields
func (a *Analyzer) extractType(ticket domain.Ticket) string {
	// Check custom fields first
	if ticketType, ok := ticket.CustomFields["Type"]; ok && ticketType != "" {
		return ticketType
	}

	// Try to infer from title
	title := strings.ToLower(ticket.Title)
	if strings.Contains(title, "[bug]") || strings.Contains(title, "bug:") {
		return "Bug"
	}
	if strings.Contains(title, "[epic]") || strings.Contains(title, "epic:") {
		return "Epic"
	}
	if strings.Contains(title, "[feature]") || strings.Contains(title, "feature:") {
		return "Feature"
	}

	// Default to Story
	return "Story"
}

// extractStatus determines the ticket status from title or custom fields
func (a *Analyzer) extractStatus(ticket domain.Ticket) string {
	// Check custom fields first
	if status, ok := ticket.CustomFields["Status"]; ok && status != "" {
		return status
	}

	// Try to infer from title markers
	title := strings.ToLower(ticket.Title)
	if strings.Contains(title, "✅") || strings.Contains(title, "done") {
		return "Done"
	}
	if strings.Contains(title, "🚧") || strings.Contains(title, "in progress") {
		return "In Progress"
	}
	if strings.Contains(title, "🔄") || strings.Contains(title, "review") {
		return "In Review"
	}

	// Check if it has a JIRA ID (likely in progress)
	if ticket.JiraID != "" {
		return "Open"
	}

	return "To Do"
}

// extractTaskStatus determines the task status
func (a *Analyzer) extractTaskStatus(task domain.Task) string {
	// Check custom fields first
	if status, ok := task.CustomFields["Status"]; ok && status != "" {
		return status
	}

	// Try to infer from title
	title := strings.ToLower(task.Title)
	if strings.Contains(title, "✅") || strings.Contains(title, "done") {
		return "Done"
	}
	if strings.Contains(title, "🚧") || strings.Contains(title, "in progress") {
		return "In Progress"
	}

	// Check if it has a JIRA ID
	if task.JiraID != "" {
		return "Open"
	}

	return "To Do"
}

// extractStoryPoints extracts story points from custom fields
func (a *Analyzer) extractStoryPoints(ticket domain.Ticket) float64 {
	if pointsStr, ok := ticket.CustomFields["Story Points"]; ok && pointsStr != "" {
		if points, err := strconv.ParseFloat(pointsStr, 64); err == nil {
			return points
		}
	}

	// Also check for alternative field names
	if pointsStr, ok := ticket.CustomFields["StoryPoints"]; ok && pointsStr != "" {
		if points, err := strconv.ParseFloat(pointsStr, 64); err == nil {
			return points
		}
	}

	return 0
}

// FormatReport generates a formatted text report of the statistics
//
// Parameters:
//   - stats: The statistics to format
//
// Returns:
//   - string: A formatted report suitable for console output
func (a *Analyzer) FormatReport(stats *Statistics) string {
    var b strings.Builder
    writeHeader(&b)
    writeOverall(&b, stats)
    writeSync(&b, stats)
    writeTicketsByType(&b, a, stats)
    writeTicketsByStatus(&b, a, stats)
    writeTasksByStatus(&b, a, stats)
    writeProgress(&b, stats)
    return b.String()
}

func writeHeader(b *strings.Builder) {
    b.WriteString("\n╔══════════════════════════════════════╗\n")
    b.WriteString("║        TICKET ANALYTICS REPORT       ║\n")
    b.WriteString("╚══════════════════════════════════════╝\n\n")
}

func writeOverall(b *strings.Builder, s *Statistics) {
    b.WriteString("📊 Overall Statistics\n")
    b.WriteString("────────────────────\n")
    b.WriteString(fmt.Sprintf("  Total Tickets:      %d\n", s.TotalTickets))
    b.WriteString(fmt.Sprintf("  Total Tasks:        %d\n", s.TotalTasks))
    b.WriteString(fmt.Sprintf("  Total Items:        %d\n", s.TotalTickets+s.TotalTasks))
    if s.TotalStoryPoints > 0 {
        b.WriteString(fmt.Sprintf("  Total Story Points: %.1f\n", s.TotalStoryPoints))
    }
    b.WriteString(fmt.Sprintf("  Acceptance Criteria: %d\n\n", s.AcceptanceCriteriaCount))
}

func writeSync(b *strings.Builder, s *Statistics) {
    b.WriteString("🔄 JIRA Synchronization\n")
    b.WriteString("────────────────────\n")
    tpct := 0; if s.TotalTickets > 0 { tpct = (s.TicketsWithJiraID*100)/s.TotalTickets }
    apct := 0; if s.TotalTasks > 0 { apct = (s.TasksWithJiraID*100)/s.TotalTasks }
    b.WriteString(fmt.Sprintf("  Tickets Synced: %d/%d (%d%%)\n", s.TicketsWithJiraID, s.TotalTickets, tpct))
    b.WriteString(fmt.Sprintf("  Tasks Synced:   %d/%d (%d%%)\n\n", s.TasksWithJiraID, s.TotalTasks, apct))
}

func writeTicketsByType(b *strings.Builder, a *Analyzer, s *Statistics) {
    if len(s.TicketsByType) == 0 { return }
    b.WriteString("📋 Tickets by Type\n")
    b.WriteString("────────────────────\n")
    for t, c := range s.TicketsByType {
        bar := a.makeBar(c, s.TotalTickets, 20)
        b.WriteString(fmt.Sprintf("  %-10s %s %d\n", t+":", bar, c))
    }
    b.WriteString("\n")
}

func writeOrderedBars(b *strings.Builder, a *Analyzer, total int, order []string, m map[string]int) {
    for _, k := range order {
        if c, ok := m[k]; ok {
            bar := a.makeBar(c, total, 20)
            b.WriteString(fmt.Sprintf("  %-12s %s %d\n", k+":", bar, c))
        }
    }
    for k, c := range m {
        found := false
        for _, okk := range order { if k == okk { found = true; break } }
        if !found {
            bar := a.makeBar(c, total, 20)
            b.WriteString(fmt.Sprintf("  %-12s %s %d\n", k+":", bar, c))
        }
    }
    b.WriteString("\n")
}

func writeTicketsByStatus(b *strings.Builder, a *Analyzer, s *Statistics) {
    if len(s.TicketsByStatus) == 0 { return }
    b.WriteString("📈 Tickets by Status\n")
    b.WriteString("────────────────────\n")
    order := []string{"Done", "In Progress", "In Review", "Open", "To Do"}
    writeOrderedBars(b, a, s.TotalTickets, order, s.TicketsByStatus)
}

func writeTasksByStatus(b *strings.Builder, a *Analyzer, s *Statistics) {
    if len(s.TasksByStatus) == 0 || s.TotalTasks == 0 { return }
    b.WriteString("✅ Tasks by Status\n")
    b.WriteString("────────────────────\n")
    order := []string{"Done", "In Progress", "In Review", "Open", "To Do"}
    writeOrderedBars(b, a, s.TotalTasks, order, s.TasksByStatus)
}

func writeProgress(b *strings.Builder, s *Statistics) {
    b.WriteString("🎯 Progress Summary\n")
    b.WriteString("────────────────────\n")
    doneTickets := s.TicketsByStatus["Done"]
    doneTasks := s.TasksByStatus["Done"]
    total := s.TotalTickets + s.TotalTasks
    done := doneTickets + doneTasks
    pct := 0; if total > 0 { pct = (done*100)/total }
    b.WriteString(fmt.Sprintf("  Overall Completion: %d%%\n", pct))
    b.WriteString(fmt.Sprintf("  Items Completed:    %d/%d\n", done, total))
    if s.TotalStoryPoints > 0 && s.TotalTickets > 0 {
        completed := (s.TotalStoryPoints * float64(doneTickets)) / float64(s.TotalTickets)
        b.WriteString(fmt.Sprintf("  Points Completed:   %.1f/%.1f\n", completed, s.TotalStoryPoints))
    }
    b.WriteString("\n────────────────────────────────────────\n")
}

// makeBar creates a simple text progress bar
func (a *Analyzer) makeBar(value, total, width int) string {
	if total == 0 {
		return strings.Repeat("─", width)
	}

	filled := (value * width) / total
	if filled > width {
		filled = width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return bar
}
