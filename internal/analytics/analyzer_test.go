package analytics

import (
	"strings"
	"testing"

	"github.com/karolswdev/ticketr/internal/core/domain"
)

func TestNewAnalyzer(t *testing.T) {
	analyzer := NewAnalyzer()
	if analyzer == nil {
		t.Error("NewAnalyzer() returned nil")
	}
}

func TestAnalyzeTickets(t *testing.T) {
	analyzer := NewAnalyzer()

	tickets := []domain.Ticket{
		{
			Title:       "[Feature] User Authentication",
			JiraID:      "PROJ-123",
			Description: "Implement user authentication system",
			AcceptanceCriteria: []string{
				"Users can login",
				"Passwords are secure",
			},
			CustomFields: map[string]string{
				"Type":         "Feature",
				"Status":       "In Progress",
				"Story Points": "5",
			},
			Tasks: []domain.Task{
				{
					Title:  "✅ Set up database",
					JiraID: "PROJ-124",
					AcceptanceCriteria: []string{
						"Database schema created",
					},
					CustomFields: map[string]string{
						"Status": "Done",
					},
				},
				{
					Title:  "🚧 Implement login API",
					JiraID: "PROJ-125",
					CustomFields: map[string]string{
						"Status": "In Progress",
					},
				},
				{
					Title:        "Create UI components",
					CustomFields: map[string]string{},
				},
			},
		},
		{
			Title:       "[Bug] Login fails with special chars",
			JiraID:      "PROJ-126",
			Description: "Fix login bug",
			CustomFields: map[string]string{
				"Type":   "Bug",
				"Status": "To Do",
			},
			Tasks: []domain.Task{},
		},
		{
			Title:       "✅ [Epic] Cloud Migration",
			Description: "Migrate to cloud",
			CustomFields: map[string]string{
				"Type":         "Epic",
				"Status":       "Done",
				"Story Points": "13",
			},
			Tasks: []domain.Task{},
		},
	}

	stats := analyzer.AnalyzeTickets(tickets)

	// Verify overall statistics
	if stats.TotalTickets != 3 {
		t.Errorf("Expected 3 tickets, got %d", stats.TotalTickets)
	}

	if stats.TotalTasks != 3 {
		t.Errorf("Expected 3 tasks, got %d", stats.TotalTasks)
	}

	if stats.TotalStoryPoints != 18 {
		t.Errorf("Expected 18 story points, got %f", stats.TotalStoryPoints)
	}

	if stats.TicketsWithJiraID != 2 {
		t.Errorf("Expected 2 tickets with JIRA ID, got %d", stats.TicketsWithJiraID)
	}

	if stats.TasksWithJiraID != 2 {
		t.Errorf("Expected 2 tasks with JIRA ID, got %d", stats.TasksWithJiraID)
	}

	if stats.AcceptanceCriteriaCount != 3 {
		t.Errorf("Expected 3 acceptance criteria, got %d", stats.AcceptanceCriteriaCount)
	}

	// Verify tickets by type
	if stats.TicketsByType["Feature"] != 1 {
		t.Errorf("Expected 1 Feature ticket, got %d", stats.TicketsByType["Feature"])
	}

	if stats.TicketsByType["Bug"] != 1 {
		t.Errorf("Expected 1 Bug ticket, got %d", stats.TicketsByType["Bug"])
	}

	if stats.TicketsByType["Epic"] != 1 {
		t.Errorf("Expected 1 Epic ticket, got %d", stats.TicketsByType["Epic"])
	}

	// Verify tickets by status
	if stats.TicketsByStatus["In Progress"] != 1 {
		t.Errorf("Expected 1 ticket In Progress, got %d", stats.TicketsByStatus["In Progress"])
	}

	if stats.TicketsByStatus["To Do"] != 1 {
		t.Errorf("Expected 1 ticket To Do, got %d", stats.TicketsByStatus["To Do"])
	}

	if stats.TicketsByStatus["Done"] != 1 {
		t.Errorf("Expected 1 ticket Done, got %d", stats.TicketsByStatus["Done"])
	}

	// Verify tasks by status
	if stats.TasksByStatus["Done"] != 1 {
		t.Errorf("Expected 1 task Done, got %d", stats.TasksByStatus["Done"])
	}

	if stats.TasksByStatus["In Progress"] != 1 {
		t.Errorf("Expected 1 task In Progress, got %d", stats.TasksByStatus["In Progress"])
	}

	if stats.TasksByStatus["To Do"] != 1 {
		t.Errorf("Expected 1 task To Do, got %d", stats.TasksByStatus["To Do"])
	}
}

func TestExtractType(t *testing.T) {
	analyzer := NewAnalyzer()

	tests := []struct {
		name     string
		ticket   domain.Ticket
		expected string
	}{
		{
			name: "Type from custom field",
			ticket: domain.Ticket{
				Title: "Some ticket",
				CustomFields: map[string]string{
					"Type": "Feature",
				},
			},
			expected: "Feature",
		},
		{
			name: "Bug from title",
			ticket: domain.Ticket{
				Title:        "[Bug] Login issue",
				CustomFields: map[string]string{},
			},
			expected: "Bug",
		},
		{
			name: "Epic from title",
			ticket: domain.Ticket{
				Title:        "Epic: Cloud Migration",
				CustomFields: map[string]string{},
			},
			expected: "Epic",
		},
		{
			name: "Feature from title",
			ticket: domain.Ticket{
				Title:        "[Feature] New dashboard",
				CustomFields: map[string]string{},
			},
			expected: "Feature",
		},
		{
			name: "Default to Story",
			ticket: domain.Ticket{
				Title:        "Regular ticket",
				CustomFields: map[string]string{},
			},
			expected: "Story",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.extractType(tt.ticket)
			if result != tt.expected {
				t.Errorf("Expected type %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestExtractStatus(t *testing.T) {
	analyzer := NewAnalyzer()

	tests := []struct {
		name     string
		ticket   domain.Ticket
		expected string
	}{
		{
			name: "Status from custom field",
			ticket: domain.Ticket{
				Title: "Some ticket",
				CustomFields: map[string]string{
					"Status": "In Review",
				},
			},
			expected: "In Review",
		},
		{
			name: "Done from emoji",
			ticket: domain.Ticket{
				Title:        "✅ Completed task",
				CustomFields: map[string]string{},
			},
			expected: "Done",
		},
		{
			name: "In Progress from emoji",
			ticket: domain.Ticket{
				Title:        "🚧 Working on this",
				CustomFields: map[string]string{},
			},
			expected: "In Progress",
		},
		{
			name: "In Review from emoji",
			ticket: domain.Ticket{
				Title:        "🔄 Ready for review",
				CustomFields: map[string]string{},
			},
			expected: "In Review",
		},
		{
			name: "Open with JIRA ID",
			ticket: domain.Ticket{
				Title:        "Regular ticket",
				JiraID:       "PROJ-123",
				CustomFields: map[string]string{},
			},
			expected: "Open",
		},
		{
			name: "Default to To Do",
			ticket: domain.Ticket{
				Title:        "New ticket",
				CustomFields: map[string]string{},
			},
			expected: "To Do",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.extractStatus(tt.ticket)
			if result != tt.expected {
				t.Errorf("Expected status %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestExtractStoryPoints(t *testing.T) {
	analyzer := NewAnalyzer()

	tests := []struct {
		name     string
		ticket   domain.Ticket
		expected float64
	}{
		{
			name: "Story Points field",
			ticket: domain.Ticket{
				CustomFields: map[string]string{
					"Story Points": "5.5",
				},
			},
			expected: 5.5,
		},
		{
			name: "StoryPoints field (no space)",
			ticket: domain.Ticket{
				CustomFields: map[string]string{
					"StoryPoints": "3",
				},
			},
			expected: 3,
		},
		{
			name: "Invalid points value",
			ticket: domain.Ticket{
				CustomFields: map[string]string{
					"Story Points": "invalid",
				},
			},
			expected: 0,
		},
		{
			name: "No story points",
			ticket: domain.Ticket{
				CustomFields: map[string]string{},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.extractStoryPoints(tt.ticket)
			if result != tt.expected {
				t.Errorf("Expected %f story points, got %f", tt.expected, result)
			}
		})
	}
}

func TestFormatReport(t *testing.T) {
	analyzer := NewAnalyzer()

	stats := &Statistics{
		TotalTickets:            10,
		TotalTasks:              15,
		TotalStoryPoints:        42.5,
		TicketsWithJiraID:       8,
		TasksWithJiraID:         12,
		AcceptanceCriteriaCount: 35,
		TicketsByType: map[string]int{
			"Story":   5,
			"Bug":     3,
			"Feature": 2,
		},
		TicketsByStatus: map[string]int{
			"Done":        4,
			"In Progress": 3,
			"To Do":       3,
		},
		TasksByStatus: map[string]int{
			"Done":        8,
			"In Progress": 4,
			"To Do":       3,
		},
	}

	report := analyzer.FormatReport(stats)

	// Check that report contains expected sections
	if !strings.Contains(report, "TICKET ANALYTICS REPORT") {
		t.Error("Report missing title")
	}

	if !strings.Contains(report, "Overall Statistics") {
		t.Error("Report missing overall statistics")
	}

	if !strings.Contains(report, "Total Tickets:      10") {
		t.Error("Report missing correct ticket count")
	}

	if !strings.Contains(report, "Total Tasks:        15") {
		t.Error("Report missing correct task count")
	}

	if !strings.Contains(report, "Total Story Points: 42.5") {
		t.Error("Report missing story points")
	}

	if !strings.Contains(report, "JIRA Synchronization") {
		t.Error("Report missing JIRA sync section")
	}

	if !strings.Contains(report, "Tickets Synced: 8/10 (80%)") {
		t.Error("Report missing ticket sync percentage")
	}

	if !strings.Contains(report, "Tasks Synced:   12/15 (80%)") {
		t.Error("Report missing task sync percentage")
	}

	if !strings.Contains(report, "Tickets by Type") {
		t.Error("Report missing tickets by type")
	}

	if !strings.Contains(report, "Tickets by Status") {
		t.Error("Report missing tickets by status")
	}

	if !strings.Contains(report, "Tasks by Status") {
		t.Error("Report missing tasks by status")
	}

	if !strings.Contains(report, "Progress Summary") {
		t.Error("Report missing progress summary")
	}

	if !strings.Contains(report, "Overall Completion: 48%") {
		t.Error("Report missing completion percentage")
	}
}

func TestMakeBar(t *testing.T) {
	analyzer := NewAnalyzer()

	tests := []struct {
		name     string
		value    int
		total    int
		width    int
		expected string
	}{
		{
			name:     "Half filled",
			value:    5,
			total:    10,
			width:    10,
			expected: "█████░░░░░",
		},
		{
			name:     "Full bar",
			value:    10,
			total:    10,
			width:    10,
			expected: "██████████",
		},
		{
			name:     "Empty bar",
			value:    0,
			total:    10,
			width:    10,
			expected: "░░░░░░░░░░",
		},
		{
			name:     "Zero total",
			value:    5,
			total:    0,
			width:    10,
			expected: "──────────",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.makeBar(tt.value, tt.total, tt.width)
			if result != tt.expected {
				t.Errorf("Expected bar '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
