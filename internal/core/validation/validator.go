package validation

import (
	"fmt"

    "github.com/karolswdev/ticketr/internal/core/domain"
)

// ValidationError represents a validation error with context
type ValidationError struct {
	Field   string
	Message string
	Line    int
}

func (v ValidationError) Error() string {
	if v.Line > 0 {
		return fmt.Sprintf("line %d: %s: %s", v.Line, v.Field, v.Message)
	}
	return fmt.Sprintf("%s: %s", v.Field, v.Message)
}

// Validator provides validation services for tickets
type Validator struct {
	hierarchyRules map[string][]string // Maps parent type to allowed child types
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		hierarchyRules: map[string][]string{
			"Epic":    {"Story", "Task", "Bug"},
			"Story":   {"Sub-task", "Task"},
			"Task":    {"Sub-task"},
			"Bug":     {"Sub-task"},
			"Feature": {"Sub-task", "Task"},
		},
	}
}

// ValidateHierarchy validates ticket hierarchy rules
func (v *Validator) ValidateHierarchy(tickets []domain.Ticket) []ValidationError {
	errors := []ValidationError{}
	
	// Build a map of ticket IDs to types
	ticketTypes := make(map[string]string)
	for _, ticket := range tickets {
		if ticket.JiraID != "" {
			if ticketType, exists := ticket.CustomFields["Type"]; exists {
				ticketTypes[ticket.JiraID] = ticketType
			}
		}
	}
	
	// Check each ticket's children
	for _, ticket := range tickets {
		parentType := ticket.CustomFields["Type"]
		if parentType == "" {
			parentType = "Story" // Default type
		}
		
		// Check if parent type has hierarchy rules
		allowedChildTypes, hasRules := v.hierarchyRules[parentType]
		if !hasRules {
			continue // No rules for this parent type
		}
		
		// Validate each child task
		for _, task := range ticket.Tasks {
			childType := task.CustomFields["Type"]
			if childType == "" {
				childType = "Sub-task" // Default child type
			}
			
			// Check if child type is allowed
			allowed := false
			for _, allowedType := range allowedChildTypes {
				if childType == allowedType {
					allowed = true
					break
				}
			}
			
			if !allowed {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("Task '%s'", task.Title),
					Message: fmt.Sprintf("A '%s' cannot be the child of a '%s'", childType, parentType),
					Line:    task.SourceLine,
				})
			}
		}
	}
	
	return errors
}

// ValidateRequiredFields validates that required fields are present
func (v *Validator) ValidateRequiredFields(ticket domain.Ticket, requiredFields []string) []ValidationError {
	errors := []ValidationError{}
	
	// Check title
	if ticket.Title == "" {
		errors = append(errors, ValidationError{
			Field:   "Title",
			Message: "Title is required",
			Line:    ticket.SourceLine,
		})
	}
	
	// Check custom required fields
	for _, field := range requiredFields {
		if value, exists := ticket.CustomFields[field]; !exists || value == "" {
			errors = append(errors, ValidationError{
				Field:   field,
				Message: fmt.Sprintf("Required field '%s' is missing or empty", field),
				Line:    ticket.SourceLine,
			})
		}
	}
	
	return errors
}

// ValidateTickets performs comprehensive validation on tickets
func (v *Validator) ValidateTickets(tickets []domain.Ticket) []ValidationError {
	allErrors := []ValidationError{}
	
	// Validate hierarchy
	hierarchyErrors := v.ValidateHierarchy(tickets)
	allErrors = append(allErrors, hierarchyErrors...)
	
	// Validate each ticket's required fields (basic validation)
	for _, ticket := range tickets {
		// Basic required fields check (title is always required)
		if ticket.Title == "" {
			allErrors = append(allErrors, ValidationError{
				Field:   "Title",
				Message: "Title is required",
				Line:    ticket.SourceLine,
			})
		}
	}
	
	return allErrors
}
