package services

import (
	"fmt"
	"log"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// TicketService orchestrates the business logic for processing tickets
// DEPRECATED: Use PushService for CLI operations (provides state management)
// This service is kept for backward compatibility and testing
type TicketService struct {
	repository ports.Repository
	jiraClient ports.JiraPort
}

// NewTicketService creates a new instance of TicketService
func NewTicketService(repository ports.Repository, jiraClient ports.JiraPort) *TicketService {
	return &TicketService{
		repository: repository,
		jiraClient: jiraClient,
	}
}

// ProcessResult holds the results of processing tickets and tasks
type ProcessResult struct {
	TicketsCreated int
	TicketsUpdated int
	TasksCreated   int
	TasksUpdated   int
	Errors         []string
}

// ProcessOptions contains options for processing tickets
type ProcessOptions struct {
	ForcePartialUpload bool
}

// calculateFinalFields merges parent fields with task fields (task fields override parent fields)
func (s *TicketService) calculateFinalFields(parent domain.Ticket, task domain.Task) map[string]string {
	// Start with parent's fields
	finalFields := make(map[string]string)
	for k, v := range parent.CustomFields {
		finalFields[k] = v
	}

	// Override with task's fields
	for k, v := range task.CustomFields {
		finalFields[k] = v
	}

	return finalFields
}

// ProcessTickets reads tickets from the repository and creates/updates them in Jira
func (s *TicketService) ProcessTickets(filePath string) (*ProcessResult, error) {
	return s.ProcessTicketsWithOptions(filePath, ProcessOptions{})
}

// ProcessTicketsWithOptions reads tickets from the repository and creates/updates them in Jira with options
func (s *TicketService) ProcessTicketsWithOptions(filePath string, options ProcessOptions) (*ProcessResult, error) {
	result := &ProcessResult{
		Errors: []string{},
	}

	// Read tickets from the file
	tickets, err := s.repository.GetTickets(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tickets from file: %w", err)
	}

	// Process each ticket
	for i := range tickets {
		ticket := &tickets[i]

		// Check if ticket needs to be created or updated
		if ticket.JiraID != "" {
			// Update existing ticket in Jira
			err := s.jiraClient.UpdateTicket(*ticket)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to update ticket '%s' (%s): %v", ticket.Title, ticket.JiraID, err)
				result.Errors = append(result.Errors, errMsg)
				log.Println(errMsg)
				continue
			}
			result.TicketsUpdated++
			log.Printf("Updated ticket '%s' with Jira ID: %s\n", ticket.Title, ticket.JiraID)
		} else {
			// Create new ticket in Jira
			jiraID, err := s.jiraClient.CreateTicket(*ticket)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to create ticket '%s': %v", ticket.Title, err)
				result.Errors = append(result.Errors, errMsg)
				log.Println(errMsg)
				continue
			}

			// Update the ticket with the new Jira ID
			ticket.JiraID = jiraID
			result.TicketsCreated++
			log.Printf("Created ticket '%s' with Jira ID: %s\n", ticket.Title, jiraID)
		}

		// Process tasks for this ticket
		for j := range ticket.Tasks {
			task := &ticket.Tasks[j]

			// Calculate final fields for task (inherit from parent + task overrides)
			finalFields := s.calculateFinalFields(*ticket, *task)

			// Create a task with merged fields for Jira operations
			taskWithFields := *task
			taskWithFields.CustomFields = finalFields

			// Check if task needs to be created or updated
			if task.JiraID != "" {
				// Update existing task in Jira
				err := s.jiraClient.UpdateTask(taskWithFields)
				if err != nil {
					errMsg := fmt.Sprintf("  Failed to update task '%s' (%s): %v", task.Title, task.JiraID, err)
					result.Errors = append(result.Errors, errMsg)
					log.Println(errMsg)
					continue
				}
				result.TasksUpdated++
				log.Printf("  Updated task '%s' with Jira ID: %s\n", task.Title, task.JiraID)
			} else {
				// Create new task in Jira (needs parent ticket to exist)
				if ticket.JiraID == "" {
					errMsg := fmt.Sprintf("  Cannot create task '%s' - parent ticket has no Jira ID", task.Title)
					result.Errors = append(result.Errors, errMsg)
					log.Println(errMsg)
					continue
				}

				taskJiraID, err := s.jiraClient.CreateTask(taskWithFields, ticket.JiraID)
				if err != nil {
					errMsg := fmt.Sprintf("  Failed to create task '%s': %v", task.Title, err)
					result.Errors = append(result.Errors, errMsg)
					log.Println(errMsg)
					continue
				}

				// Update the task with the new Jira ID
				task.JiraID = taskJiraID
				result.TasksCreated++
				log.Printf("  Created task '%s' with Jira ID: %s\n", task.Title, taskJiraID)
			}
		}
	}

	// Save the updated tickets back to the file
	err = s.repository.SaveTickets(filePath, tickets)
	if err != nil {
		// This is not critical - we've already created the items in Jira
		log.Printf("Warning: Failed to save updated tickets back to file: %v\n", err)
	}

	return result, nil
}
