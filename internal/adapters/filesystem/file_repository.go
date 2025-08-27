package filesystem

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/karolswdev/ticktr/internal/core/domain"
)

// FileRepository implements the Repository port for file-based storage
type FileRepository struct{}

// NewFileRepository creates a new instance of FileRepository
func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

// parserState represents the current state of the parser
type parserState int

const (
	stateNone parserState = iota
	stateStoryDescription
	stateStoryAcceptanceCriteria
	stateTasks
	stateTaskDescription
	stateTaskAcceptanceCriteria
)

// GetStories reads and parses stories from a file according to the custom Markdown format
func (r *FileRepository) GetStories(filepath string) ([]domain.Story, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var stories []domain.Story
	var currentStory *domain.Story
	var currentTask *domain.Task
	state := stateNone

	// Regex patterns for parsing
	storyRegex := regexp.MustCompile(`^#\s+STORY:\s*(?:\[([^\]]+)\])?\s*(.+)$`)
	taskRegex := regexp.MustCompile(`^[-*]\s+(?:\[([^\]]+)\])?\s*(.+)$`)
	descriptionRegex := regexp.MustCompile(`^\s*[-*]\s+Description:\s*(.*)$`)
	acHeaderRegex := regexp.MustCompile(`^\s*[-*]\s+Acceptance Criteria:$`)
	acItemRegex := regexp.MustCompile(`^\s*[-*]\s+(.+)$`)

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Handle story separator
		if strings.TrimSpace(line) == "---" {
			if currentStory != nil {
				// Save current task if exists
				if currentTask != nil {
					currentStory.Tasks = append(currentStory.Tasks, *currentTask)
					currentTask = nil
				}
				stories = append(stories, *currentStory)
				currentStory = nil
				state = stateNone
			}
			continue
		}

		// Parse story heading
		if matches := storyRegex.FindStringSubmatch(line); matches != nil {
			// Save previous story if exists
			if currentStory != nil {
				if currentTask != nil {
					currentStory.Tasks = append(currentStory.Tasks, *currentTask)
					currentTask = nil
				}
				stories = append(stories, *currentStory)
			}

			currentStory = &domain.Story{
				JiraID: matches[1],
				Title:  strings.TrimSpace(matches[2]),
			}
			state = stateStoryDescription
			continue
		}

		// Check for malformed story heading  
		if strings.HasPrefix(line, "## STORY:") {
			return nil, fmt.Errorf("malformed story heading at line %d: use '# STORY:' instead of '## STORY:'", lineNum)
		}

		// Parse section headers
		if strings.HasPrefix(line, "## Description") && currentStory != nil {
			state = stateStoryDescription
			continue
		}
		if strings.HasPrefix(line, "## Acceptance Criteria") && currentStory != nil {
			state = stateStoryAcceptanceCriteria
			continue
		}
		if strings.HasPrefix(line, "## Tasks") && currentStory != nil {
			state = stateTasks
			continue
		}

		// Process based on current state
		switch state {
		case stateStoryDescription:
			if currentStory != nil && strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "#") {
				if currentStory.Description != "" {
					currentStory.Description += "\n"
				}
				currentStory.Description += strings.TrimSpace(line)
			}

		case stateStoryAcceptanceCriteria:
			if currentStory != nil && strings.HasPrefix(strings.TrimSpace(line), "-") {
				ac := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "-"))
				if ac != "" {
					currentStory.AcceptanceCriteria = append(currentStory.AcceptanceCriteria, ac)
				}
			}

		case stateTasks:
			// Check for new task
			if matches := taskRegex.FindStringSubmatch(line); matches != nil {
				// Save previous task if exists
				if currentTask != nil && currentStory != nil {
					currentStory.Tasks = append(currentStory.Tasks, *currentTask)
				}

				currentTask = &domain.Task{
					JiraID: matches[1],
					Title:  strings.TrimSpace(matches[2]),
				}
				state = stateTasks
				continue
			}

			// Parse task details (indented content)
			if currentTask != nil {
				if descMatch := descriptionRegex.FindStringSubmatch(line); descMatch != nil {
					currentTask.Description = strings.TrimSpace(descMatch[1])
					state = stateTaskDescription
				} else if acHeaderRegex.MatchString(line) {
					state = stateTaskAcceptanceCriteria
				}
			}

		case stateTaskDescription:
			// Continue reading task description on next lines if they're indented
			if currentTask != nil && strings.HasPrefix(line, "    ") && strings.TrimSpace(line) != "" {
				if currentTask.Description != "" {
					currentTask.Description += " "
				}
				currentTask.Description += strings.TrimSpace(line)
			} else if acHeaderRegex.MatchString(line) {
				state = stateTaskAcceptanceCriteria
			} else if !strings.HasPrefix(line, "  ") {
				// Back to tasks state if we're not indented anymore
				state = stateTasks
				// Reprocess this line in tasks state
				if matches := taskRegex.FindStringSubmatch(line); matches != nil {
					if currentTask != nil && currentStory != nil {
						currentStory.Tasks = append(currentStory.Tasks, *currentTask)
					}
					currentTask = &domain.Task{
						JiraID: matches[1],
						Title:  strings.TrimSpace(matches[2]),
					}
				}
			}

		case stateTaskAcceptanceCriteria:
			if currentTask != nil && strings.HasPrefix(line, "    ") && acItemRegex.MatchString(line) {
				ac := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "-"))
				if ac != "" {
					currentTask.AcceptanceCriteria = append(currentTask.AcceptanceCriteria, ac)
				}
			} else if !strings.HasPrefix(line, "  ") {
				// Back to tasks state if we're not indented anymore
				state = stateTasks
				// Reprocess this line in tasks state
				if matches := taskRegex.FindStringSubmatch(line); matches != nil {
					if currentTask != nil && currentStory != nil {
						currentStory.Tasks = append(currentStory.Tasks, *currentTask)
					}
					currentTask = &domain.Task{
						JiraID: matches[1],
						Title:  strings.TrimSpace(matches[2]),
					}
				}
			}
		}
	}

	// Save last story and task if exists
	if currentTask != nil && currentStory != nil {
		currentStory.Tasks = append(currentStory.Tasks, *currentTask)
	}
	if currentStory != nil {
		stories = append(stories, *currentStory)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return stories, nil
}

// SaveStories writes stories to a file in the custom Markdown format
func (r *FileRepository) SaveStories(filepath string, stories []domain.Story) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i, story := range stories {
		// Write story heading with Jira ID if present
		if story.JiraID != "" {
			fmt.Fprintf(writer, "# STORY: [%s] %s\n", story.JiraID, story.Title)
		} else {
			fmt.Fprintf(writer, "# STORY: %s\n", story.Title)
		}
		fmt.Fprintln(writer)

		// Write description
		if story.Description != "" {
			fmt.Fprintln(writer, "## Description")
			fmt.Fprintln(writer, story.Description)
			fmt.Fprintln(writer)
		}

		// Write acceptance criteria
		if len(story.AcceptanceCriteria) > 0 {
			fmt.Fprintln(writer, "## Acceptance Criteria")
			for _, ac := range story.AcceptanceCriteria {
				fmt.Fprintf(writer, "- %s\n", ac)
			}
			fmt.Fprintln(writer)
		}

		// Write tasks
		if len(story.Tasks) > 0 {
			fmt.Fprintln(writer, "## Tasks")
			for _, task := range story.Tasks {
				// Write task with Jira ID if present
				if task.JiraID != "" {
					fmt.Fprintf(writer, "- [%s] %s\n", task.JiraID, task.Title)
				} else {
					fmt.Fprintf(writer, "- %s\n", task.Title)
				}
				
				// Write task description
				if task.Description != "" {
					fmt.Fprintf(writer, "  - Description: %s\n", task.Description)
				}
				
				// Write task acceptance criteria
				if len(task.AcceptanceCriteria) > 0 {
					fmt.Fprintln(writer, "  - Acceptance Criteria:")
					for _, ac := range task.AcceptanceCriteria {
						fmt.Fprintf(writer, "    - %s\n", ac)
					}
				}
			}
			fmt.Fprintln(writer)
		}

		// Add separator between stories (except for the last one)
		if i < len(stories)-1 {
			fmt.Fprintln(writer, "---")
			fmt.Fprintln(writer)
		}
	}

	return writer.Flush()
}