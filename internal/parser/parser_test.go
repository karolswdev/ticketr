package parser

import (
	"testing"
)

func TestParser_RecognizesTicketBlock(t *testing.T) {
	parser := New()
	
	tickets, err := parser.Parse("../../testdata/ticket_simple.md")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if len(tickets) != 1 {
		t.Fatalf("Expected 1 ticket, got %d", len(tickets))
	}
	
	ticket := tickets[0]
	
	if ticket.Title != "Create user authentication system" {
		t.Errorf("Expected title 'Create user authentication system', got '%s'", ticket.Title)
	}
	
	expectedDesc := "Implement a complete user authentication system with login, logout, and session management capabilities."
	if ticket.Description != expectedDesc {
		t.Errorf("Expected description '%s', got '%s'", expectedDesc, ticket.Description)
	}
	
	// Check CustomFields
	expectedFields := map[string]string{
		"Type": "Story",
		"Project": "PROJ",
		"Priority": "High",
		"Sprint": "10",
	}
	
	for key, expectedVal := range expectedFields {
		if val, ok := ticket.CustomFields[key]; !ok {
			t.Errorf("Missing field '%s'", key)
		} else if val != expectedVal {
			t.Errorf("Field '%s': expected '%s', got '%s'", key, expectedVal, val)
		}
	}
	
	// Check acceptance criteria
	if len(ticket.AcceptanceCriteria) != 3 {
		t.Errorf("Expected 3 acceptance criteria, got %d", len(ticket.AcceptanceCriteria))
	}
}

func TestParser_ParsesNestedTasks(t *testing.T) {
	parser := New()

	tickets, err := parser.Parse("../../testdata/ticket_with_tasks.md")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(tickets) != 1 {
		t.Fatalf("Expected 1 ticket, got %d", len(tickets))
	}

	ticket := tickets[0]

	if ticket.Title != "Build payment processing system" {
		t.Errorf("Expected title 'Build payment processing system', got '%s'", ticket.Title)
	}

	// Check Tasks
	if len(ticket.Tasks) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(ticket.Tasks))
	}

	// Check first task
	task1 := ticket.Tasks[0]
	if task1.Title != "Create payment gateway interface" {
		t.Errorf("Task 1: Expected title 'Create payment gateway interface', got '%s'", task1.Title)
	}

	if task1.CustomFields["Priority"] != "Low" {
		t.Errorf("Task 1: Expected Priority 'Low', got '%s'", task1.CustomFields["Priority"])
	}

	// Check second task
	task2 := ticket.Tasks[1]
	if task2.Title != "Implement Stripe integration" {
		t.Errorf("Task 2: Expected title 'Implement Stripe integration', got '%s'", task2.Title)
	}

	if task2.CustomFields["Priority"] != "Critical" {
		t.Errorf("Task 2: Expected Priority 'Critical', got '%s'", task2.CustomFields["Priority"])
	}

	if task2.CustomFields["Assignee"] != "John Doe" {
		t.Errorf("Task 2: Expected Assignee 'John Doe', got '%s'", task2.CustomFields["Assignee"])
	}
}

func TestParser_RejectsLegacyStoryFormat(t *testing.T) {
	parser := New()

	// Test that parser rejects # STORY: format
	tickets, err := parser.Parse("../../testdata/legacy_story/simple_story.md")

	// Should return an error
	if err == nil {
		t.Fatal("Expected error for legacy # STORY: format, but got nil")
	}

	// Verify error message contains "Legacy" and "STORY"
	errMsg := err.Error()
	if !contains(errMsg, "Legacy") {
		t.Errorf("Expected error message to contain 'Legacy', got: %s", errMsg)
	}
	if !contains(errMsg, "STORY") {
		t.Errorf("Expected error message to contain 'STORY', got: %s", errMsg)
	}

	// Verify tickets are nil or empty since parsing failed
	if tickets != nil && len(tickets) > 0 {
		t.Errorf("Expected no tickets to be returned on error, got %d tickets", len(tickets))
	}
}

func TestParser_ErrorMessageIncludesMigrationGuidance(t *testing.T) {
	parser := New()

	// Test that error message is helpful
	_, err := parser.Parse("../../testdata/legacy_story/simple_story.md")

	if err == nil {
		t.Fatal("Expected error for legacy # STORY: format, but got nil")
	}

	errMsg := err.Error()

	// Verify error contains "migrate" command reference
	if !contains(errMsg, "migrate") {
		t.Errorf("Expected error message to contain 'migrate', got: %s", errMsg)
	}

	// Verify error contains "REQUIREMENTS-v2.md" reference
	if !contains(errMsg, "REQUIREMENTS-v2.md") {
		t.Errorf("Expected error message to contain 'REQUIREMENTS-v2.md', got: %s", errMsg)
	}

	// Verify error includes line number
	if !contains(errMsg, "line") {
		t.Errorf("Expected error message to contain 'line', got: %s", errMsg)
	}
}

func TestParser_RejectsMultipleLegacyStories(t *testing.T) {
	parser := New()

	// Test with testdata/legacy_story/multiple_stories.md
	_, err := parser.Parse("../../testdata/legacy_story/multiple_stories.md")

	if err == nil {
		t.Fatal("Expected error for legacy # STORY: format with multiple stories, but got nil")
	}

	// Verify parser catches first occurrence (should mention line 1)
	errMsg := err.Error()
	if !contains(errMsg, "line 1") {
		t.Errorf("Expected error message to mention 'line 1' (first occurrence), got: %s", errMsg)
	}
}

// Helper function to check if string contains substring (case-insensitive not needed here)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr))))
}