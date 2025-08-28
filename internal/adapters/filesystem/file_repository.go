package filesystem

import (
	"bufio"
	"fmt"
	"os"
	
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/parser"
)

// FileRepository implements the Repository port for file-based storage
type FileRepository struct {
	parser *parser.Parser
}

// NewFileRepository creates a new instance of FileRepository
func NewFileRepository() *FileRepository {
	return &FileRepository{
		parser: parser.New(),
	}
}

// GetTickets reads and parses tickets from a file using the new parser
func (r *FileRepository) GetTickets(filepath string) ([]domain.Ticket, error) {
	return r.parser.Parse(filepath)
}

// SaveTickets writes tickets to a file in the new TICKET format
func (r *FileRepository) SaveTickets(filepath string, tickets []domain.Ticket) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i, ticket := range tickets {
		// Write ticket heading with Jira ID if present
		if ticket.JiraID != "" {
			fmt.Fprintf(writer, "# TICKET: [%s] %s\n", ticket.JiraID, ticket.Title)
		} else {
			fmt.Fprintf(writer, "# TICKET: %s\n", ticket.Title)
		}
		fmt.Fprintln(writer)

		// Write description
		if ticket.Description != "" {
			fmt.Fprintln(writer, "## Description")
			fmt.Fprintln(writer, ticket.Description)
			fmt.Fprintln(writer)
		}

		// Write fields
		if len(ticket.CustomFields) > 0 {
			fmt.Fprintln(writer, "## Fields")
			for key, value := range ticket.CustomFields {
				fmt.Fprintf(writer, "%s: %s\n", key, value)
			}
			fmt.Fprintln(writer)
		}

		// Write acceptance criteria
		if len(ticket.AcceptanceCriteria) > 0 {
			fmt.Fprintln(writer, "## Acceptance Criteria")
			for _, ac := range ticket.AcceptanceCriteria {
				fmt.Fprintf(writer, "- %s\n", ac)
			}
			fmt.Fprintln(writer)
		}

		// Write tasks
		if len(ticket.Tasks) > 0 {
			fmt.Fprintln(writer, "## Tasks")
			for _, task := range ticket.Tasks {
				// Write task with Jira ID if present
				if task.JiraID != "" {
					fmt.Fprintf(writer, "- [%s] %s\n", task.JiraID, task.Title)
				} else {
					fmt.Fprintf(writer, "- %s\n", task.Title)
				}

				// Write task description (indented)
				if task.Description != "" {
					fmt.Fprintln(writer, "  ## Description")
					// Indent description lines
					lines := fmt.Sprintf("%s", task.Description)
					fmt.Fprintf(writer, "  %s\n", lines)
					fmt.Fprintln(writer)
				}

				// Write task fields (indented)
				if len(task.CustomFields) > 0 {
					fmt.Fprintln(writer, "  ## Fields")
					for key, value := range task.CustomFields {
						fmt.Fprintf(writer, "  %s: %s\n", key, value)
					}
					fmt.Fprintln(writer)
				}

				// Write task acceptance criteria (indented)
				if len(task.AcceptanceCriteria) > 0 {
					fmt.Fprintln(writer, "  ## Acceptance Criteria")
					for _, ac := range task.AcceptanceCriteria {
						fmt.Fprintf(writer, "  - %s\n", ac)
					}
					fmt.Fprintln(writer)
				}
			}
		}

		// Add spacing between tickets
		if i < len(tickets)-1 {
			fmt.Fprintln(writer)
		}
	}

	return writer.Flush()
}