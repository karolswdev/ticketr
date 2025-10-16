package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	
	"github.com/karolswdev/ticktr/internal/core/domain"
)

type Parser struct{}

func New() *Parser {
	return &Parser{}
}

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
		line := scanner.Text()

		// Check for legacy # STORY: format
		if strings.HasPrefix(strings.TrimSpace(line), "# STORY:") {
			return nil, fmt.Errorf("Legacy '# STORY:' format detected at line %d. Please migrate to '# TICKET:' format. See REQUIREMENTS-v2.md (PROD-201) or use 'ticketr migrate <file>' command.", lineNum)
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return p.parseLines(lines)
}

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
			
			// Parse ticket sections
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

type multilineResult struct {
	content string
	nextIdx int
}

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

type fieldsResult struct {
	fields  map[string]string
	nextIdx int
}

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

type criteriaResult struct {
	criteria []string
	nextIdx  int
}

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
		
		// Stop at task list item if we're in the parent context and see a dash without further indentation
		// Don't break for acceptance criteria items that are properly indented within their section
		// The AC items within a task should have their - at the start after removing task indent
		// This is only for stopping when we see a new TASK at the parent level
		
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

type tasksResult struct {
	tasks   []domain.Task
	nextIdx int
}

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
			
			// Parse task sections (they should be indented)
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