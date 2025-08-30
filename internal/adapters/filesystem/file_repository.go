package filesystem

import (
	"bufio"
	"fmt"
	"os"

	"github.com/karolswdev/ticketr/internal/core/domain"
	"github.com/karolswdev/ticketr/internal/parser"
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
	defer func() { _ = file.Close() }()

	writer := bufio.NewWriter(file)
	w := func(n int, err error) error {
		if err != nil {
			return err
		}
		return nil
	}

	for i, ticket := range tickets {
		// Write ticket heading with Jira ID if present
		if ticket.JiraID != "" {
			if err := w(fmt.Fprintf(writer, "# TICKET: [%s] %s\n", ticket.JiraID, ticket.Title)); err != nil {
				return err
			}
		} else {
			if err := w(fmt.Fprintf(writer, "# TICKET: %s\n", ticket.Title)); err != nil {
				return err
			}
		}
		if err := w(fmt.Fprintln(writer)); err != nil {
			return err
		}

		// Write description
		if ticket.Description != "" {
			if err := w(fmt.Fprintln(writer, "## Description")); err != nil {
				return err
			}
			if err := w(fmt.Fprintln(writer, ticket.Description)); err != nil {
				return err
			}
			if err := w(fmt.Fprintln(writer)); err != nil {
				return err
			}
		}

		// Write fields
		if len(ticket.CustomFields) > 0 {
			if err := w(fmt.Fprintln(writer, "## Fields")); err != nil {
				return err
			}
			for key, value := range ticket.CustomFields {
				if err := w(fmt.Fprintf(writer, "%s: %s\n", key, value)); err != nil {
					return err
				}
			}
			if err := w(fmt.Fprintln(writer)); err != nil {
				return err
			}
		}

		// Write acceptance criteria
		if len(ticket.AcceptanceCriteria) > 0 {
			if err := w(fmt.Fprintln(writer, "## Acceptance Criteria")); err != nil {
				return err
			}
			for _, ac := range ticket.AcceptanceCriteria {
				if err := w(fmt.Fprintf(writer, "- %s\n", ac)); err != nil {
					return err
				}
			}
			if err := w(fmt.Fprintln(writer)); err != nil {
				return err
			}
		}

		// Write tasks
		if len(ticket.Tasks) > 0 {
			if err := w(fmt.Fprintln(writer, "## Tasks")); err != nil {
				return err
			}
			for _, task := range ticket.Tasks {
				// Write task with Jira ID if present
				if task.JiraID != "" {
					if err := w(fmt.Fprintf(writer, "- [%s] %s\n", task.JiraID, task.Title)); err != nil {
						return err
					}
				} else {
					if err := w(fmt.Fprintf(writer, "- %s\n", task.Title)); err != nil {
						return err
					}
				}

				// Write task description (indented)
				if task.Description != "" {
					if err := w(fmt.Fprintln(writer, "  ## Description")); err != nil {
						return err
					}
					// Indent description lines
					if err := w(fmt.Fprintf(writer, "  %s\n", task.Description)); err != nil {
						return err
					}
					if err := w(fmt.Fprintln(writer)); err != nil {
						return err
					}
				}

				// Write task fields (indented)
				if len(task.CustomFields) > 0 {
					if err := w(fmt.Fprintln(writer, "  ## Fields")); err != nil {
						return err
					}
					for key, value := range task.CustomFields {
						if err := w(fmt.Fprintf(writer, "  %s: %s\n", key, value)); err != nil {
							return err
						}
					}
					if err := w(fmt.Fprintln(writer)); err != nil {
						return err
					}
				}

				// Write task acceptance criteria (indented)
				if len(task.AcceptanceCriteria) > 0 {
					if err := w(fmt.Fprintln(writer, "  ## Acceptance Criteria")); err != nil {
						return err
					}
					for _, ac := range task.AcceptanceCriteria {
						if err := w(fmt.Fprintf(writer, "  - %s\n", ac)); err != nil {
							return err
						}
					}
					if err := w(fmt.Fprintln(writer)); err != nil {
						return err
					}
				}
			}
		}

		// Add spacing between tickets
		if i < len(tickets)-1 {
			if err := w(fmt.Fprintln(writer)); err != nil {
				return err
			}
		}
	}

	return writer.Flush()
}
