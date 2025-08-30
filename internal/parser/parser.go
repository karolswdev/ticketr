// Package parser provides functionality for parsing Markdown files containing ticket definitions.
// It supports the Tickets-as-Code format, allowing teams to define Jira tickets in Markdown
// and sync them with Jira while maintaining version control.
package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	
    "github.com/karolswdev/ticketr/internal/core/domain"
)

// Parser handles the parsing of Markdown files containing ticket definitions.
// It supports hierarchical ticket structures with tasks, descriptions, acceptance criteria,
// and custom fields.
type Parser struct{}

// New creates a new Parser instance.
//
// Returns:
//   - *Parser: A new parser ready to parse Markdown files
func New() *Parser {
	return &Parser{}
}

// Parse reads and parses a Markdown file containing ticket definitions.
// It extracts all tickets and their associated metadata from the file.
//
// Parameters:
//   - filePath: The path to the Markdown file to parse
//
// Returns:
//   - []domain.Ticket: A slice of parsed tickets with their tasks and metadata
//   - error: An error if the file cannot be read or parsed
func (p *Parser) Parse(filePath string) ([]domain.Ticket, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	var lines []string
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		lines = append(lines, scanner.Text())
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	
	return p.parseLines(lines)
}

// parseLines processes an array of lines and extracts ticket definitions.
// It identifies tickets by the "# TICKET:" header pattern and parses their content.
//
// Parameters:
//   - lines: The lines of the Markdown file to parse
//
// Returns:
//   - []domain.Ticket: A slice of parsed tickets
//   - error: An error if parsing fails
func (p *Parser) parseLines(lines []string) ([]domain.Ticket, error) {
	var tickets []domain.Ticket
	ticketRegex := regexp.MustCompile(`^# TICKET:\s*(?:\[([^\]]+)\])?\s*(.+)$`)
	
	for i := 0; i < len(lines); i++ {
		matches := ticketRegex.FindStringSubmatch(lines[i])
		if matches != nil {
			ticket := domain.Ticket{
				JiraID:       matches[1],
				Title:        strings.TrimSpace(matches[2]),
				SourceLine:   i + 1,
				CustomFields: make(map[string]string),
			}
			
			// Parse ticket sections starting from the next line
			// The parseTicketSections function will handle all nested content
			i++
			nextIdx := p.parseTicketSections(&ticket, lines, i, 0)
			
			tickets = append(tickets, ticket)
			
			// Continue from where parseTicketSections left off
			// But subtract 1 because the loop will increment
			i = nextIdx - 1
		}
	}
	
	return tickets, nil
}

// parseTicketSections parses all sections within a ticket definition.
// It handles Description, Fields, Acceptance Criteria, and Tasks sections.
//
// Parameters:
//   - ticket: The ticket object to populate with parsed data
//   - lines: All lines from the Markdown file
//   - startIdx: The starting line index to parse from
//   - indent: The current indentation level (in spaces)
//
// Returns:
//   - int: The next line index after parsing this ticket's content
func (p *Parser) parseTicketSections(ticket *domain.Ticket, lines []string, startIdx int, indent int) int {
	i := startIdx
	indentStr := strings.Repeat(" ", indent)
	
	for i < len(lines) {
		line := lines[i]
		
		// Check if we've reached the next ticket
		if strings.HasPrefix(strings.TrimSpace(line), "# TICKET:") {
			return i
		}
		
		// Check if we've gone back to a lower indent level
		// This indicates we've exited the current ticket's scope
		if indent > 0 && !strings.HasPrefix(line, indentStr) && strings.TrimSpace(line) != "" {
			break
		}
		
		// Remove the expected indentation
		if indent > 0 && strings.HasPrefix(line, indentStr) {
			line = line[indent:]
		}
		
		// Check for section headers
		if strings.HasPrefix(line, "## Description") {
			i++
			desc := p.parseMultilineSection(lines, i, indent)
			ticket.Description = strings.TrimSpace(desc.content)
			i = desc.nextIdx
		} else if strings.HasPrefix(line, "## Fields") {
			i++
			fields := p.parseFieldsSection(lines, i, indent)
			for k, v := range fields.fields {
				ticket.CustomFields[k] = v
			}
			i = fields.nextIdx
		} else if strings.HasPrefix(line, "## Acceptance Criteria") {
			i++
			ac := p.parseAcceptanceCriteria(lines, i, indent)
			ticket.AcceptanceCriteria = ac.criteria
			i = ac.nextIdx
		} else if strings.HasPrefix(line, "## Tasks") {
			i++
			tasks := p.parseTasks(lines, i, indent)
			ticket.Tasks = tasks.tasks
			i = tasks.nextIdx
			// If parseTasks found a next ticket, we should return that index
			if i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), "# TICKET:") {
				return i
			}
		} else {
			i++
		}
	}
	
	return i
}

// multilineResult contains the parsed content and the next line index to process.
type multilineResult struct {
	content string
	nextIdx int
}

// parseMultilineSection parses a multi-line text section, preserving formatting.
// It stops when it encounters a new section header or reaches the end of the current scope.
//
// Parameters:
//   - lines: All lines from the Markdown file
//   - startIdx: The starting line index
//   - baseIndent: The base indentation level to expect
//
// Returns:
//   - multilineResult: The parsed content and next line index
func (p *Parser) parseMultilineSection(lines []string, startIdx int, baseIndent int) multilineResult {
	var content []string
	i := startIdx
	
	for i < len(lines) {
		line := lines[i]
		
		// Check if line starts a new section (## header)
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "##") || strings.HasPrefix(trimmed, "# TICKET:") {
			break
		}
		
		// Check for task list item (starts with -)
		if baseIndent > 0 && strings.TrimSpace(line) != "" && strings.HasPrefix(strings.TrimSpace(line), "-") && !strings.HasPrefix(line, strings.Repeat(" ", baseIndent+2)) {
			break
		}
		
		// Add the line content (removing base indentation if present)
		if baseIndent > 0 && strings.HasPrefix(line, strings.Repeat(" ", baseIndent)) {
			content = append(content, line[baseIndent:])
		} else if baseIndent == 0 {
			content = append(content, line)
		} else if strings.TrimSpace(line) == "" {
			content = append(content, "")
		} else {
			break
		}
		
		i++
	}
	
	return multilineResult{
		content: strings.TrimSpace(strings.Join(content, "\n")),
		nextIdx: i,
	}
}

// fieldsResult contains parsed field key-value pairs and the next line index.
type fieldsResult struct {
	fields  map[string]string
	nextIdx int
}

// parseFieldsSection parses custom field definitions in "key: value" format.
// Fields are used to set Jira custom fields like Story Points, Sprint, etc.
//
// Parameters:
//   - lines: All lines from the Markdown file
//   - startIdx: The starting line index
//   - baseIndent: The base indentation level
//
// Returns:
//   - fieldsResult: A map of field names to values and the next line index
func (p *Parser) parseFieldsSection(lines []string, startIdx int, baseIndent int) fieldsResult {
	fields := make(map[string]string)
	i := startIdx
	fieldRegex := regexp.MustCompile(`^([^:]+):\s*(.*)$`)
	
	for i < len(lines) {
		line := lines[i]
		
		// Remove base indentation
		if baseIndent > 0 && strings.HasPrefix(line, strings.Repeat(" ", baseIndent)) {
			line = line[baseIndent:]
		}
		
		trimmed := strings.TrimSpace(line)
		
		// Stop at next section or end
		if strings.HasPrefix(trimmed, "##") || strings.HasPrefix(trimmed, "# TICKET:") {
			break
		}
		
		// Stop at task list item if we're in a task context
		if baseIndent > 0 && trimmed != "" && strings.HasPrefix(trimmed, "-") {
			break
		}
		
		// Skip comments and empty lines
		if strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, "##") {
			i++
			continue
		}
		
		if trimmed == "" {
			i++
			continue
		}
		
		// Parse field
		if matches := fieldRegex.FindStringSubmatch(trimmed); matches != nil {
			fields[matches[1]] = strings.TrimSpace(matches[2])
		}
		
		i++
	}
	
	return fieldsResult{
		fields:  fields,
		nextIdx: i,
	}
}

// criteriaResult contains parsed acceptance criteria and the next line index.
type criteriaResult struct {
	criteria []string
	nextIdx  int
}

// parseAcceptanceCriteria parses acceptance criteria items from a bulleted list.
// Each criterion should start with a dash (-) and will be added to the ticket.
//
// Parameters:
//   - lines: All lines from the Markdown file
//   - startIdx: The starting line index
//   - baseIndent: The base indentation level
//
// Returns:
//   - criteriaResult: A slice of criteria and the next line index
func (p *Parser) parseAcceptanceCriteria(lines []string, startIdx int, baseIndent int) criteriaResult {
	var criteria []string
	i := startIdx
	
	for i < len(lines) {
		if i >= len(lines) {
			break
		}
		originalLine := lines[i]
		line := originalLine
		
		// Check if line has less indentation than expected (indicates we're back at parent level)
		expectedIndent := strings.Repeat(" ", baseIndent)
		if baseIndent > 0 && !strings.HasPrefix(line, expectedIndent) && strings.TrimSpace(line) != "" {
			// Line has content but doesn't have the required indentation - we've left this section
			break
		}
		
		// Remove base indentation
		if baseIndent > 0 && strings.HasPrefix(line, expectedIndent) {
			line = line[baseIndent:]
		}
		
		trimmed := strings.TrimSpace(line)
		
		// Stop at next section
		if strings.HasPrefix(trimmed, "##") || strings.HasPrefix(trimmed, "# TICKET:") {
			break
		}
		
		// We continue parsing dash-prefixed items as acceptance criteria
		// The parsing stops when we encounter a new section header or
		// return to a lower indentation level
		
		// Parse criteria item
		if strings.HasPrefix(trimmed, "-") {
			criterion := strings.TrimSpace(strings.TrimPrefix(trimmed, "-"))
			if criterion != "" {
				criteria = append(criteria, criterion)
			}
		}
		
		i++
	}
	
	return criteriaResult{
		criteria: criteria,
		nextIdx:  i,
	}
}

// tasksResult contains parsed tasks and the next line index.
type tasksResult struct {
	tasks   []domain.Task
	nextIdx int
}

// parseTasks parses task definitions from a task list.
// Tasks can have their own descriptions, fields, and acceptance criteria.
//
// Parameters:
//   - lines: All lines from the Markdown file
//   - startIdx: The starting line index
//   - baseIndent: The base indentation level
//
// Returns:
//   - tasksResult: A slice of parsed tasks and the next line index
func (p *Parser) parseTasks(lines []string, startIdx int, baseIndent int) tasksResult {
	var tasks []domain.Task
	i := startIdx
	taskRegex := regexp.MustCompile(`^-\s*(?:\[([^\]]+)\])?\s*(.+)$`)
	
	for i < len(lines) {
		line := lines[i]
		
		// Remove base indentation
		if baseIndent > 0 && strings.HasPrefix(line, strings.Repeat(" ", baseIndent)) {
			line = line[baseIndent:]
		}
		
		trimmed := strings.TrimSpace(line)
		
		// Stop at next ticket-level section
		if strings.HasPrefix(trimmed, "## ") && !strings.HasPrefix(line, "  ") {
			break
		}
		if strings.HasPrefix(trimmed, "# TICKET:") {
			break
		}
		
		// Check for task item
		if matches := taskRegex.FindStringSubmatch(trimmed); matches != nil {
			task := domain.Task{
				JiraID:       matches[1],
				Title:        strings.TrimSpace(matches[2]),
				SourceLine:   i + 1,
				CustomFields: make(map[string]string),
			}
			
			// Parse task sections with increased indentation
			// Tasks require 2 additional spaces of indentation for their content
			i++
			i = p.parseTaskSections(&task, lines, i, baseIndent+2)
			
			tasks = append(tasks, task)
			i-- // Adjust because loop will increment
		}
		
		i++
	}
	
	return tasksResult{
		tasks:   tasks,
		nextIdx: i,
	}
}

// parseTaskSections parses all sections within a task definition.
// It handles Description, Fields, and Acceptance Criteria for individual tasks.
//
// Parameters:
//   - task: The task object to populate with parsed data
//   - lines: All lines from the Markdown file
//   - startIdx: The starting line index
//   - indent: The expected indentation level for task content
//
// Returns:
//   - int: The next line index after parsing this task's content
func (p *Parser) parseTaskSections(task *domain.Task, lines []string, startIdx int, indent int) int {
	i := startIdx
	indentStr := strings.Repeat(" ", indent)
	
	for i < len(lines) {
		if i >= len(lines) {
			break
		}
		
		line := lines[i]
		
		// Check if we've reached a new ticket
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# TICKET:") {
			break
		}
		
		// If line doesn't start with expected indent and is not empty, we're done with this task
		if !strings.HasPrefix(line, indentStr) && trimmed != "" {
			// Check if it's a task-level item (starts with -)
			if strings.HasPrefix(trimmed, "-") {
				break
			}
			// Check if it's a ticket-level section
			if strings.HasPrefix(trimmed, "##") && !strings.HasPrefix(line, "  ") {
				break
			}
		}
		
		// Process the line with proper indentation removed
		processLine := line
		if strings.HasPrefix(line, indentStr) {
			processLine = line[indent:]
		}
		
		trimmed = strings.TrimSpace(processLine)
		
		// Parse different sections
		if strings.HasPrefix(trimmed, "## Description") {
			i++
			desc := p.parseMultilineSection(lines, i, indent)
			task.Description = strings.TrimSpace(desc.content)
			i = desc.nextIdx
		} else if strings.HasPrefix(trimmed, "## Fields") {
			i++
			fields := p.parseFieldsSection(lines, i, indent)
			for k, v := range fields.fields {
				task.CustomFields[k] = v
			}
			i = fields.nextIdx
		} else if strings.HasPrefix(trimmed, "## Acceptance Criteria") {
			i++
			ac := p.parseAcceptanceCriteria(lines, i, indent)
			task.AcceptanceCriteria = ac.criteria
			i = ac.nextIdx
		} else {
			i++
		}
	}
	
	return i
}
