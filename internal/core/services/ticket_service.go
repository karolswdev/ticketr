package services

import (
	"fmt"
	"log"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// TicketService orchestrates the business logic for processing tickets
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
	// Legacy compatibility
	StoriesCreated int
	StoriesUpdated int
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

// ProcessStories reads stories from the repository and creates/updates them in Jira (backwards compatibility)
func (s *TicketService) ProcessStories(filePath string) (*ProcessResult, error) {
	return s.ProcessStoriesWithOptions(filePath, ProcessOptions{})
}

// ProcessStoriesWithOptions reads stories from the repository and creates/updates them in Jira with options
func (s *TicketService) ProcessStoriesWithOptions(filePath string, options ProcessOptions) (*ProcessResult, error) {
	result := &ProcessResult{
		Errors: []string{},
	}

	// Read stories from the file (this will internally use the new ticket parser)
	stories, err := s.repository.GetStories(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read stories from file: %w", err)
	}

	// Process each story (which is now backed by a Ticket)
	for i := range stories {
		story := &stories[i]
		
		// Check if story needs to be created or updated
		if story.JiraID != "" {
			// Update existing story in Jira
			err := s.jiraClient.UpdateStory(*story)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to update story '%s' (%s): %v", story.Title, story.JiraID, err)
				result.Errors = append(result.Errors, errMsg)
				log.Println(errMsg)
				continue
			}
			result.StoriesUpdated++
			result.TicketsUpdated++
			log.Printf("Updated story '%s' with Jira ID: %s\n", story.Title, story.JiraID)
		} else {
			// Create new story in Jira
			jiraID, err := s.jiraClient.CreateStory(*story)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to create story '%s': %v", story.Title, err)
				result.Errors = append(result.Errors, errMsg)
				log.Println(errMsg)
				continue
			}
			
			// Update the story with the new Jira ID
			story.JiraID = jiraID
			result.StoriesCreated++
			result.TicketsCreated++
			log.Printf("Created story '%s' with Jira ID: %s\n", story.Title, jiraID)
		}

		// Process tasks for this story
		for j := range story.Tasks {
			task := &story.Tasks[j]
			
			// Check if task needs to be created or updated
			if task.JiraID != "" {
				// Update existing task in Jira
				err := s.jiraClient.UpdateTask(*task)
				if err != nil {
					errMsg := fmt.Sprintf("  Failed to update task '%s' (%s): %v", task.Title, task.JiraID, err)
					result.Errors = append(result.Errors, errMsg)
					log.Println(errMsg)
					continue
				}
				result.TasksUpdated++
				log.Printf("  Updated task '%s' with Jira ID: %s\n", task.Title, task.JiraID)
			} else {
				// Create new task in Jira (needs parent story to exist)
				if story.JiraID == "" {
					errMsg := fmt.Sprintf("  Cannot create task '%s' - parent story has no Jira ID", task.Title)
					result.Errors = append(result.Errors, errMsg)
					log.Println(errMsg)
					continue
				}
				
				taskJiraID, err := s.jiraClient.CreateTask(*task, story.JiraID)
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

	// Save the updated stories back to the file
	err = s.repository.SaveStories(filePath, stories)
	if err != nil {
		// This is not critical - we've already created the items in Jira
		log.Printf("Warning: Failed to save updated stories back to file: %v\n", err)
	}

	return result, nil
}

// ProcessTickets reads tickets from the repository and creates/updates them in Jira
func (s *TicketService) ProcessTickets(filePath string) (*ProcessResult, error) {
	return s.ProcessTicketsWithOptions(filePath, ProcessOptions{})
}

// ProcessTicketsWithOptions reads tickets from the repository and creates/updates them in Jira with options
func (s *TicketService) ProcessTicketsWithOptions(filePath string, options ProcessOptions) (*ProcessResult, error) {
	// For now, use the same implementation as ProcessStories
	// In the future, this will work directly with tickets and custom fields
	return s.ProcessStoriesWithOptions(filePath, options)
}

// Legacy compatibility - kept for backward compatibility
type StoryService = TicketService

// NewStoryService creates a new instance of StoryService (backward compatibility)
func NewStoryService(repository ports.Repository, jiraClient ports.JiraPort) *StoryService {
	return NewTicketService(repository, jiraClient)
}