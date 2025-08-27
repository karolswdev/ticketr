package services

import (
	"fmt"
	"log"

	"github.com/karolswdev/ticktr/internal/core/ports"
)

// StoryService orchestrates the business logic for processing stories
type StoryService struct {
	repository ports.Repository
	jiraClient ports.JiraPort
}

// NewStoryService creates a new instance of StoryService
func NewStoryService(repository ports.Repository, jiraClient ports.JiraPort) *StoryService {
	return &StoryService{
		repository: repository,
		jiraClient: jiraClient,
	}
}

// ProcessResult holds the results of processing stories and tasks
type ProcessResult struct {
	StoriesCreated int
	StoriesUpdated int
	TasksCreated   int
	TasksUpdated   int
	Errors         []string
}

// ProcessStories reads stories from the repository and creates/updates them in Jira
func (s *StoryService) ProcessStories(filePath string) (*ProcessResult, error) {
	result := &ProcessResult{
		Errors: []string{},
	}

	// Read stories from the file
	stories, err := s.repository.GetStories(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read stories from file: %w", err)
	}

	// Process each story
	for i := range stories {
		story := &stories[i]
		
		// Skip if story already has a Jira ID (update functionality will come in Phase 3)
		if story.JiraID != "" {
			log.Printf("Skipping story '%s' - already has Jira ID: %s\n", story.Title, story.JiraID)
			continue
		}

		// Create the story in Jira
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
		log.Printf("Created story '%s' with Jira ID: %s\n", story.Title, jiraID)

		// Process tasks for this story
		for j := range story.Tasks {
			task := &story.Tasks[j]
			
			// Skip if task already has a Jira ID
			if task.JiraID != "" {
				log.Printf("  Skipping task '%s' - already has Jira ID: %s\n", task.Title, task.JiraID)
				continue
			}

			// Create the task in Jira
			taskJiraID, err := s.jiraClient.CreateTask(*task, jiraID)
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

	// Save the updated stories back to the file
	err = s.repository.SaveStories(filePath, stories)
	if err != nil {
		// This is not critical - we've already created the items in Jira
		log.Printf("Warning: Failed to save updated stories back to file: %v\n", err)
	}

	return result, nil
}