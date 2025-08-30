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
    if err != nil { return fmt.Errorf("failed to create file: %w", err) }
    defer func() { _ = file.Close() }()

    writer := bufio.NewWriter(file)
    w := func(n int, err error) error { if err != nil { return err }; return nil }

    writeHeader := func(t domain.Ticket) error {
        if t.JiraID != "" {
            return w(fmt.Fprintf(writer, "# TICKET: [%s] %s\n", t.JiraID, t.Title))
        }
        return w(fmt.Fprintf(writer, "# TICKET: %s\n", t.Title))
    }
    writeDescription := func(desc string) error {
        if desc == "" { return nil }
        if err := w(fmt.Fprintln(writer, "## Description")); err != nil { return err }
        if err := w(fmt.Fprintln(writer, desc)); err != nil { return err }
        return w(fmt.Fprintln(writer))
    }
    writeFields := func(fields map[string]string) error {
        if len(fields) == 0 { return nil }
        if err := w(fmt.Fprintln(writer, "## Fields")); err != nil { return err }
        for k, v := range fields {
            if err := w(fmt.Fprintf(writer, "%s: %s\n", k, v)); err != nil { return err }
        }
        return w(fmt.Fprintln(writer))
    }
    writeAcceptance := func(criteria []string) error {
        if len(criteria) == 0 { return nil }
        if err := w(fmt.Fprintln(writer, "## Acceptance Criteria")); err != nil { return err }
        for _, ac := range criteria {
            if err := w(fmt.Fprintf(writer, "- %s\n", ac)); err != nil { return err }
        }
        return w(fmt.Fprintln(writer))
    }
    writeTask := func(t domain.Task) error {
        if t.JiraID != "" {
            if err := w(fmt.Fprintf(writer, "- [%s] %s\n", t.JiraID, t.Title)); err != nil { return err }
        } else {
            if err := w(fmt.Fprintf(writer, "- %s\n", t.Title)); err != nil { return err }
        }
        if t.Description != "" {
            if err := w(fmt.Fprintln(writer, "  ## Description")); err != nil { return err }
            if err := w(fmt.Fprintf(writer, "  %s\n", t.Description)); err != nil { return err }
            if err := w(fmt.Fprintln(writer)); err != nil { return err }
        }
        if len(t.CustomFields) > 0 {
            if err := w(fmt.Fprintln(writer, "  ## Fields")); err != nil { return err }
            for k, v := range t.CustomFields {
                if err := w(fmt.Fprintf(writer, "  %s: %s\n", k, v)); err != nil { return err }
            }
            if err := w(fmt.Fprintln(writer)); err != nil { return err }
        }
        if len(t.AcceptanceCriteria) > 0 {
            if err := w(fmt.Fprintln(writer, "  ## Acceptance Criteria")); err != nil { return err }
            for _, ac := range t.AcceptanceCriteria {
                if err := w(fmt.Fprintf(writer, "  - %s\n", ac)); err != nil { return err }
            }
            if err := w(fmt.Fprintln(writer)); err != nil { return err }
        }
        return nil
    }
    writeTasks := func(tasks []domain.Task) error {
        if len(tasks) == 0 { return nil }
        if err := w(fmt.Fprintln(writer, "## Tasks")); err != nil { return err }
        for _, t := range tasks {
            if err := writeTask(t); err != nil { return err }
        }
        return nil
    }

    for i, ticket := range tickets {
        if err := writeHeader(ticket); err != nil { return err }
        if err := w(fmt.Fprintln(writer)); err != nil { return err }
        if err := writeDescription(ticket.Description); err != nil { return err }
        if err := writeFields(ticket.CustomFields); err != nil { return err }
        if err := writeAcceptance(ticket.AcceptanceCriteria); err != nil { return err }
        if err := writeTasks(ticket.Tasks); err != nil { return err }

        if i < len(tickets)-1 {
            if err := w(fmt.Fprintln(writer)); err != nil { return err }
        }
    }
    return writer.Flush()
}
