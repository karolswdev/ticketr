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