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
	TotalTickets    int
	TotalTasks      int
	TicketsByType   map[string]int
	TicketsByStatus map[string]int
	TasksByStatus   map[string]int
	TotalStoryPoints float64
	TicketsWithJiraID int
	TasksWithJiraID   int
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
	var report strings.Builder
	
	report.WriteString("\n╔══════════════════════════════════════╗\n")
	report.WriteString("║        TICKET ANALYTICS REPORT       ║\n")
	report.WriteString("╚══════════════════════════════════════╝\n\n")
	
	// Overall Statistics
	report.WriteString("📊 Overall Statistics\n")
	report.WriteString("────────────────────\n")
	report.WriteString(fmt.Sprintf("  Total Tickets:      %d\n", stats.TotalTickets))
	report.WriteString(fmt.Sprintf("  Total Tasks:        %d\n", stats.TotalTasks))
	report.WriteString(fmt.Sprintf("  Total Items:        %d\n", stats.TotalTickets+stats.TotalTasks))
	if stats.TotalStoryPoints > 0 {
		report.WriteString(fmt.Sprintf("  Total Story Points: %.1f\n", stats.TotalStoryPoints))
	}
	report.WriteString(fmt.Sprintf("  Acceptance Criteria: %d\n", stats.AcceptanceCriteriaCount))
	report.WriteString("\n")
	
	// JIRA Sync Status
	report.WriteString("🔄 JIRA Synchronization\n")
	report.WriteString("────────────────────\n")
	ticketSyncPercent := 0
	if stats.TotalTickets > 0 {
		ticketSyncPercent = (stats.TicketsWithJiraID * 100) / stats.TotalTickets
	}
	taskSyncPercent := 0
	if stats.TotalTasks > 0 {
		taskSyncPercent = (stats.TasksWithJiraID * 100) / stats.TotalTasks
	}
	report.WriteString(fmt.Sprintf("  Tickets Synced: %d/%d (%d%%)\n", 
		stats.TicketsWithJiraID, stats.TotalTickets, ticketSyncPercent))
	report.WriteString(fmt.Sprintf("  Tasks Synced:   %d/%d (%d%%)\n", 
		stats.TasksWithJiraID, stats.TotalTasks, taskSyncPercent))
	report.WriteString("\n")
	
	// Tickets by Type
	if len(stats.TicketsByType) > 0 {
		report.WriteString("📋 Tickets by Type\n")
		report.WriteString("────────────────────\n")
		for ticketType, count := range stats.TicketsByType {
			bar := a.makeBar(count, stats.TotalTickets, 20)
			report.WriteString(fmt.Sprintf("  %-10s %s %d\n", ticketType+":", bar, count))
		}
		report.WriteString("\n")
	}
	
	// Tickets by Status
	if len(stats.TicketsByStatus) > 0 {
		report.WriteString("📈 Tickets by Status\n")
		report.WriteString("────────────────────\n")
		// Order: Done, In Progress, In Review, Open, To Do
		statusOrder := []string{"Done", "In Progress", "In Review", "Open", "To Do"}
		for _, status := range statusOrder {
			if count, exists := stats.TicketsByStatus[status]; exists {
				bar := a.makeBar(count, stats.TotalTickets, 20)
				report.WriteString(fmt.Sprintf("  %-12s %s %d\n", status+":", bar, count))
			}
		}
		// Add any other statuses not in the standard order
		for status, count := range stats.TicketsByStatus {
			found := false
			for _, orderedStatus := range statusOrder {
				if status == orderedStatus {
					found = true
					break
				}
			}
			if !found {
				bar := a.makeBar(count, stats.TotalTickets, 20)
				report.WriteString(fmt.Sprintf("  %-12s %s %d\n", status+":", bar, count))
			}
		}
		report.WriteString("\n")
	}
	
	// Tasks by Status
	if len(stats.TasksByStatus) > 0 && stats.TotalTasks > 0 {
		report.WriteString("✅ Tasks by Status\n")
		report.WriteString("────────────────────\n")
		statusOrder := []string{"Done", "In Progress", "In Review", "Open", "To Do"}
		for _, status := range statusOrder {
			if count, exists := stats.TasksByStatus[status]; exists {
				bar := a.makeBar(count, stats.TotalTasks, 20)
				report.WriteString(fmt.Sprintf("  %-12s %s %d\n", status+":", bar, count))
			}
		}
		// Add any other statuses
		for status, count := range stats.TasksByStatus {
			found := false
			for _, orderedStatus := range statusOrder {
				if status == orderedStatus {
					found = true
					break
				}
			}
			if !found {
				bar := a.makeBar(count, stats.TotalTasks, 20)
				report.WriteString(fmt.Sprintf("  %-12s %s %d\n", status+":", bar, count))
			}
		}
		report.WriteString("\n")
	}
	
	// Progress Summary
	report.WriteString("🎯 Progress Summary\n")
	report.WriteString("────────────────────\n")
	
	// Calculate completion percentage
	doneTickets := stats.TicketsByStatus["Done"]
	doneTasks := stats.TasksByStatus["Done"]
	totalItems := stats.TotalTickets + stats.TotalTasks
	doneItems := doneTickets + doneTasks
	
	completionPercent := 0
	if totalItems > 0 {
		completionPercent = (doneItems * 100) / totalItems
	}
	
	report.WriteString(fmt.Sprintf("  Overall Completion: %d%%\n", completionPercent))
	report.WriteString(fmt.Sprintf("  Items Completed:    %d/%d\n", doneItems, totalItems))
	
	if stats.TotalStoryPoints > 0 {
		// Estimate completed story points (simplified - assumes even distribution)
		completedPoints := (stats.TotalStoryPoints * float64(doneTickets)) / float64(stats.TotalTickets)
		report.WriteString(fmt.Sprintf("  Points Completed:   %.1f/%.1f\n", 
			completedPoints, stats.TotalStoryPoints))
	}
	
	report.WriteString("\n────────────────────────────────────────\n")
	
	return report.String()
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
