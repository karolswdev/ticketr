package services

import (
	"fmt"
	"log"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
	"github.com/karolswdev/ticktr/internal/state"
)

// PushService handles pushing tickets to JIRA with state management
type PushService struct {
	repository   ports.Repository
	jiraClient   ports.JiraPort
	stateManager *state.StateManager
}

// NewPushService creates a new instance of PushService
func NewPushService(repository ports.Repository, jiraClient ports.JiraPort, stateManager *state.StateManager) *PushService {
	return &PushService{
		repository:   repository,
		jiraClient:   jiraClient,
		stateManager: stateManager,
	}
}

// calculateFinalFields merges parent fields with task fields (task fields override parent fields)
func (s *PushService) calculateFinalFields(parent domain.Ticket, task domain.Task) map[string]string {
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

// PushTickets processes tickets with state management to avoid redundant updates
func (s *PushService) PushTickets(filePath string, options ProcessOptions) (*ProcessResult, error) {
	result := &ProcessResult{
		Errors: []string{},
	}

	// Load the current state
	if err := s.stateManager.Load(); err != nil {
		log.Printf("Warning: Could not load state file: %v", err)
		// Continue anyway - we'll treat everything as changed
	}

	// Read tickets from the file
	tickets, err := s.repository.GetTickets(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tickets from file: %w", err)
	}

	// Process each ticket
	for i := range tickets {
		ticket := &tickets[i]

		// Check if ticket has changed
		if !s.stateManager.HasChanged(*ticket) {
			log.Printf("Skipping unchanged ticket '%s' (%s)", ticket.Title, ticket.JiraID)
			continue
		}

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
			s.stateManager.UpdateHash(*ticket)
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
			s.stateManager.UpdateHash(*ticket)
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

		// Update the hash after processing tasks too
		s.stateManager.UpdateHash(*ticket)
	}

	// Save the updated tickets back to the file
	err = s.repository.SaveTickets(filePath, tickets)
	if err != nil {
		// This is not critical - we've already created the items in Jira
		log.Printf("Warning: Failed to save updated tickets back to file: %v\n", err)
	}

	// Save the state file
	if err := s.stateManager.Save(); err != nil {
		log.Printf("Warning: Could not save state file: %v", err)
	}

	// Return error if any tickets failed
	if len(result.Errors) > 0 {
		return result, fmt.Errorf("%d ticket(s) failed to process", len(result.Errors))
	}

	return result, nil
}
