package validation

import (
	"testing"

    "github.com/karolswdev/ticketr/internal/core/domain"
)

// Test Case TC-201.1: TestValidation_Hierarchical
func TestValidation_Hierarchical(t *testing.T) {
	// Arrange: Create a ticket structure that violates a hierarchical rule
	// (e.g., a "Task" cannot be the parent of a "Story")
	tickets := []domain.Ticket{
		{
			Title:       "My Task",
			Description: "This is a task",
			CustomFields: map[string]string{
				"Type": "Task",
			},
			JiraID: "PROJ-100",
			Tasks: []domain.Task{
				{
					Title:       "My Story",
					Description: "This should not be allowed as a child of Task",
					CustomFields: map[string]string{
						"Type": "Story", // Story cannot be a child of Task
					},
					SourceLine: 10,
				},
			},
		},
	}

	// Act: Run it through a new validation service
	validator := NewValidator()
	errors := validator.ValidateHierarchy(tickets)

	// Assert: The service returns a specific validation error
	if len(errors) == 0 {
		t.Error("Expected validation error for invalid hierarchy, but got none")
	}

	// Check for the specific error
	foundError := false
	for _, err := range errors {
		if err.Message == "A 'Story' cannot be the child of a 'Task'" {
			foundError = true
			t.Logf("Successfully caught hierarchy violation: %s", err.Error())
			break
		}
	}

	if !foundError {
		t.Errorf("Expected specific hierarchy error message, got: %v", errors)
	}
}

// Additional test: Valid hierarchy should not produce errors
func TestValidation_ValidHierarchy(t *testing.T) {
	tickets := []domain.Ticket{
		{
			Title: "My Story",
			CustomFields: map[string]string{
				"Type": "Story",
			},
			JiraID: "PROJ-100",
			Tasks: []domain.Task{
				{
					Title: "My Sub-task",
					CustomFields: map[string]string{
						"Type": "Sub-task", // Sub-task is allowed under Story
					},
				},
			},
		},
	}

	validator := NewValidator()
	errors := validator.ValidateHierarchy(tickets)

	if len(errors) > 0 {
		t.Errorf("Valid hierarchy produced errors: %v", errors)
	}
}

// Test Epic with various child types
func TestValidation_EpicHierarchy(t *testing.T) {
	tickets := []domain.Ticket{
		{
			Title: "My Epic",
			CustomFields: map[string]string{
				"Type": "Epic",
			},
			JiraID: "PROJ-100",
			Tasks: []domain.Task{
				{
					Title: "Story under Epic",
					CustomFields: map[string]string{
						"Type": "Story", // Allowed
					},
				},
				{
					Title: "Task under Epic",
					CustomFields: map[string]string{
						"Type": "Task", // Allowed
					},
				},
				{
					Title: "Bug under Epic",
					CustomFields: map[string]string{
						"Type": "Bug", // Allowed
					},
				},
				{
					Title: "Sub-task under Epic",
					CustomFields: map[string]string{
						"Type": "Sub-task", // NOT allowed
					},
					SourceLine: 20,
				},
			},
		},
	}

	validator := NewValidator()
	errors := validator.ValidateHierarchy(tickets)

	if len(errors) != 1 {
		t.Errorf("Expected 1 validation error, got %d: %v", len(errors), errors)
	}

	if len(errors) > 0 && errors[0].Message != "A 'Sub-task' cannot be the child of a 'Epic'" {
		t.Errorf("Unexpected error message: %s", errors[0].Message)
	}
}
