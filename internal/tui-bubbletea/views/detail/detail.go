package detail

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/core/domain"
)

// Model represents the ticket detail view component.
// It displays full ticket information including title, metadata, description,
// acceptance criteria, and subtasks in a scrollable viewport.
type Model struct {
	viewport viewport.Model
	ticket   *domain.Ticket
	width    int
	height   int
	ready    bool
}

// New creates a new detail view model with the specified dimensions.
func New(width, height int) Model {
	vp := viewport.New(width, height)
	vp.HighPerformanceRendering = false

	return Model{
		viewport: vp,
		width:    width,
		height:   height,
		ready:    true,
	}
}

// SetTicket updates the displayed ticket and refreshes the viewport content.
func (m *Model) SetTicket(ticket *domain.Ticket) {
	m.ticket = ticket

	// Render ticket details into viewport content
	content := m.renderTicketContent()
	m.viewport.SetContent(content)
}

// SetSize updates the viewport dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width
	m.viewport.Height = height

	// Re-render content to fit new dimensions
	if m.ticket != nil {
		content := m.renderTicketContent()
		m.viewport.SetContent(content)
	}
}

// renderTicketContent generates the formatted ticket display content.
func (m Model) renderTicketContent() string {
	if m.ticket == nil {
		return lipgloss.NewStyle().
			Faint(true).
			Padding(1).
			Render("No ticket selected")
	}

	var content strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#61AFEF")).
		PaddingBottom(1)

	content.WriteString(titleStyle.Render(m.ticket.JiraID + ": " + m.ticket.Title))
	content.WriteString("\n\n")

	// Metadata
	if len(m.ticket.CustomFields) > 0 {
		content.WriteString(m.renderMetadata())
		content.WriteString("\n\n")
	}

	// Description
	if m.ticket.Description != "" {
		content.WriteString(m.renderSection("Description", m.ticket.Description))
		content.WriteString("\n\n")
	}

	// Acceptance Criteria
	if len(m.ticket.AcceptanceCriteria) > 0 {
		content.WriteString(m.renderAcceptanceCriteria())
		content.WriteString("\n\n")
	}

	// Subtasks
	if len(m.ticket.Tasks) > 0 {
		content.WriteString(m.renderSubtasks())
		content.WriteString("\n")
	}

	return content.String()
}

// renderMetadata formats the ticket's custom fields as metadata.
func (m Model) renderMetadata() string {
	var parts []string

	// Sort keys for consistent display
	for key, value := range m.ticket.CustomFields {
		parts = append(parts, fmt.Sprintf("%s: %s", key, value))
	}

	metaText := strings.Join(parts, " | ")

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ABB2BF")).
		Render(metaText)
}

// renderSection renders a titled section with content.
func (m Model) renderSection(title, content string) string {
	var result strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#98C379"))

	result.WriteString(headerStyle.Render(title + ":"))
	result.WriteString("\n")
	result.WriteString(content)

	return result.String()
}

// renderAcceptanceCriteria formats the acceptance criteria list.
func (m Model) renderAcceptanceCriteria() string {
	var content strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#98C379"))

	content.WriteString(headerStyle.Render("Acceptance Criteria:"))
	content.WriteString("\n")

	for i, ac := range m.ticket.AcceptanceCriteria {
		content.WriteString(fmt.Sprintf("  %d. %s\n", i+1, ac))
	}

	return strings.TrimRight(content.String(), "\n")
}

// renderSubtasks formats the subtasks list.
func (m Model) renderSubtasks() string {
	var content strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#98C379"))

	content.WriteString(headerStyle.Render(fmt.Sprintf("Subtasks (%d):", len(m.ticket.Tasks))))
	content.WriteString("\n")

	for _, task := range m.ticket.Tasks {
		content.WriteString(fmt.Sprintf("  • %s: %s\n", task.JiraID, task.Title))
	}

	return strings.TrimRight(content.String(), "\n")
}

// Update handles viewport navigation and scrolling.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.viewport.LineUp(1)
		case "down", "j":
			m.viewport.LineDown(1)
		case "pgup", "ctrl+u":
			m.viewport.HalfViewUp()
		case "pgdown", "ctrl+d":
			m.viewport.HalfViewDown()
		case "g":
			m.viewport.GotoTop()
		case "G":
			m.viewport.GotoBottom()
		}

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the detail view.
func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}
	return m.viewport.View()
}

// Init initializes the detail view component.
func (m Model) Init() tea.Cmd {
	return nil
}

// Focus gives focus to the detail view.
func (m *Model) Focus() {
	// Detail view focus state (viewport handles scroll state)
}

// Blur removes focus from the detail view.
func (m *Model) Blur() {
	// Detail view blur state
}
