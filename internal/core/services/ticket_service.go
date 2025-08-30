// Package services provides the core business logic for Ticketr.
// It orchestrates interactions between the repository, JIRA adapter,
// and other components to implement the Tickets-as-Code workflow.
package services

import (
	"fmt"
	"log"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// TicketService orchestrates the business logic for processing tickets.
// It coordinates between the repository (file system) and JIRA to
// synchronize ticket definitions bidirectionally.
type TicketService struct {
	repository ports.Repository
	jiraClient ports.JiraPort
}

// NewTicketService creates a new instance of TicketService.
//
// Parameters:
//   - repository: The repository for reading/writing ticket files
//   - jiraClient: The JIRA adapter for API interactions
//
// Returns:
//   - *TicketService: A configured ticket service ready for use
func NewTicketService(repository ports.Repository, jiraClient ports.JiraPort) *TicketService {
	return &TicketService{
		repository: repository,
		jiraClient: jiraClient,
	}
}

// ProcessResult holds the results of processing tickets and tasks.
// It provides detailed statistics about the synchronization operation
// and any errors that occurred during processing.
type ProcessResult struct {
	TicketsCreated int      // Number of new tickets created in JIRA
	TicketsUpdated int      // Number of existing tickets updated in JIRA
	TasksCreated   int      // Number of new tasks created in JIRA
	TasksUpdated   int      // Number of existing tasks updated in JIRA
	Errors         []string // List of errors encountered during processing
}

// ProcessOptions contains configuration options for ticket processing.
// These options control how the service handles errors and conflicts.
type ProcessOptions struct {
	ForcePartialUpload bool // Continue processing even if some tickets fail
	DryRun             bool // Validate and show what would be done without making changes
}

// calculateFinalFields merges parent ticket fields with task-specific fields.
// Task fields take precedence over parent fields, allowing tasks to override
// or extend the parent ticket's field values.
//
// Parameters:
//   - parent: The parent ticket containing base field values
//   - task: The task with potential field overrides
//
// Returns:
//   - map[string]string: The merged field map with task overrides applied
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


// ProcessTickets reads tickets from a Markdown file and synchronizes them with JIRA.
// This is a convenience method that uses default processing options.
//
// Parameters:
//   - filePath: Path to the Markdown file containing ticket definitions
//
// Returns:
//   - *ProcessResult: Statistics about the processing operation
//   - error: An error if the file cannot be read or processed
func (s *TicketService) ProcessTickets(filePath string) (*ProcessResult, error) {
	return s.ProcessTicketsWithOptions(filePath, ProcessOptions{})
}

// ProcessTicketsWithOptions reads tickets from a Markdown file and synchronizes them with JIRA.
// It supports advanced options for error handling and partial uploads.
//
// The method performs the following steps:
// 1. Reads and parses tickets from the Markdown file
// 2. Creates new tickets in JIRA (assigns JIRA IDs)
// 3. Updates existing tickets in JIRA
// 4. Processes tasks for each ticket
// 5. Writes updated ticket definitions back to the file
//
// Parameters:
//   - filePath: Path to the Markdown file containing ticket definitions
//   - options: Processing options for error handling
//
// Returns:
//   - *ProcessResult: Detailed statistics and any errors encountered
//   - error: A critical error that prevented processing
func (s *TicketService) ProcessTicketsWithOptions(filePath string, options ProcessOptions) (*ProcessResult, error) {
	result := &ProcessResult{
		Errors: []string{},
	}

	// Read tickets from the file
	tickets, err := s.repository.GetTickets(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tickets from file: %w", err)
	}

	// In dry-run mode, log what would be done
	if options.DryRun {
		log.Println("=== DRY RUN MODE - No changes will be made to JIRA ===")
	}

	// Process each ticket
	for i := range tickets {
		ticket := &tickets[i]
		
		// Determine whether to create or update based on JIRA ID presence
		if ticket.JiraID != "" {
			// Update existing ticket in Jira
			if options.DryRun {
				log.Printf("[DRY RUN] Would update ticket '%s' (JIRA ID: %s)\n", ticket.Title, ticket.JiraID)
				result.TicketsUpdated++
			} else {
				err := s.jiraClient.UpdateTicket(*ticket)
				if err != nil {
					errMsg := fmt.Sprintf("Failed to update ticket '%s' (%s): %v", ticket.Title, ticket.JiraID, err)
					result.Errors = append(result.Errors, errMsg)
					log.Println(errMsg)
					continue
				}
				result.TicketsUpdated++
				log.Printf("Updated ticket '%s' with Jira ID: %s\n", ticket.Title, ticket.JiraID)
			}
		} else {
			// Create new ticket in Jira
			if options.DryRun {
				log.Printf("[DRY RUN] Would create ticket '%s'\n", ticket.Title)
				result.TicketsCreated++
			} else {
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
		}

		// Process tasks for this ticket
		for j := range ticket.Tasks {
			task := &ticket.Tasks[j]
			
			// Determine whether to create or update the task
			if task.JiraID != "" {
				// Update existing task in Jira
				if options.DryRun {
					log.Printf("  [DRY RUN] Would update task '%s' (JIRA ID: %s)\n", task.Title, task.JiraID)
					result.TasksUpdated++
				} else {
					err := s.jiraClient.UpdateTask(*task)
					if err != nil {
						errMsg := fmt.Sprintf("  Failed to update task '%s' (%s): %v", task.Title, task.JiraID, err)
						result.Errors = append(result.Errors, errMsg)
						log.Println(errMsg)
						continue
					}
					result.TasksUpdated++
					log.Printf("  Updated task '%s' with Jira ID: %s\n", task.Title, task.JiraID)
				}
			} else {
				// Tasks require a parent ticket to exist in JIRA
				// Skip task creation if parent has no JIRA ID
				if ticket.JiraID == "" && !options.DryRun {
					errMsg := fmt.Sprintf("  Cannot create task '%s' - parent ticket has no Jira ID", task.Title)
					result.Errors = append(result.Errors, errMsg)
					log.Println(errMsg)
					continue
				}
				
				if options.DryRun {
					log.Printf("  [DRY RUN] Would create task '%s' under parent ticket\n", task.Title)
					result.TasksCreated++
				} else {
					taskJiraID, err := s.jiraClient.CreateTask(*task, ticket.JiraID)
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
	}

	// Write the updated tickets with JIRA IDs back to the file
	// This ensures the file reflects the current state in JIRA
	// Skip file update in dry-run mode
	if !options.DryRun {
		err = s.repository.SaveTickets(filePath, tickets)
		if err != nil {
			// Non-critical error: tickets are already in JIRA but file wasn't updated
			// Users can manually add the JIRA IDs if needed
			log.Printf("Warning: Failed to save updated tickets back to file: %v\n", err)
		}
	} else {
		log.Println("[DRY RUN] File would be updated with JIRA IDs after successful creation/update")
	}

	return result, nil
}

