package detail

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/karolswdev/ticktr/internal/core/domain"
)

func TestNew(t *testing.T) {
	width, height := 80, 24
	m := New(width, height)

	if m.width != width {
		t.Errorf("expected width %d, got %d", width, m.width)
	}
	if m.height != height {
		t.Errorf("expected height %d, got %d", height, m.height)
	}
	if !m.ready {
		t.Error("expected model to be ready")
	}
}

func TestSetTicket(t *testing.T) {
	m := New(80, 24)

	ticket := &domain.Ticket{
		JiraID:      "PROJ-123",
		Title:       "Test Ticket",
		Description: "Test description",
		CustomFields: map[string]string{
			"Priority": "High",
			"Status":   "In Progress",
		},
		AcceptanceCriteria: []string{
			"Criterion 1",
			"Criterion 2",
		},
		Tasks: []domain.Task{
			{
				JiraID: "PROJ-124",
				Title:  "Subtask 1",
			},
		},
	}

	m.SetTicket(ticket)

	if m.ticket != ticket {
		t.Error("ticket not set correctly")
	}

	// Verify viewport content is populated
	content := m.viewport.View()
	if !strings.Contains(content, "PROJ-123") {
		t.Error("viewport content should contain ticket ID")
	}
}

func TestSetTicketNil(t *testing.T) {
	m := New(80, 24)
	m.SetTicket(nil)

	content := m.viewport.View()
	if !strings.Contains(content, "No ticket selected") {
		t.Error("expected 'No ticket selected' message for nil ticket")
	}
}

func TestRenderTicketContent(t *testing.T) {
	m := New(80, 24)

	ticket := &domain.Ticket{
		JiraID:      "PROJ-456",
		Title:       "Implement feature X",
		Description: "Detailed description here",
		CustomFields: map[string]string{
			"Priority": "High",
		},
		AcceptanceCriteria: []string{"AC1", "AC2"},
		Tasks: []domain.Task{
			{JiraID: "PROJ-457", Title: "Subtask A"},
		},
	}

	m.SetTicket(ticket)
	content := m.renderTicketContent()

	// Check all sections are present
	if !strings.Contains(content, "PROJ-456") {
		t.Error("content should contain ticket ID")
	}
	if !strings.Contains(content, "Implement feature X") {
		t.Error("content should contain title")
	}
	if !strings.Contains(content, "Priority: High") {
		t.Error("content should contain custom fields")
	}
	if !strings.Contains(content, "Description:") {
		t.Error("content should contain description header")
	}
	if !strings.Contains(content, "Detailed description here") {
		t.Error("content should contain description text")
	}
	if !strings.Contains(content, "Acceptance Criteria:") {
		t.Error("content should contain AC header")
	}
	if !strings.Contains(content, "AC1") {
		t.Error("content should contain AC items")
	}
	if !strings.Contains(content, "Subtasks (1):") {
		t.Error("content should contain subtasks header")
	}
	if !strings.Contains(content, "PROJ-457") {
		t.Error("content should contain subtask ID")
	}
}

func TestUpdate_Navigation(t *testing.T) {
	m := New(80, 24)

	// Set a ticket with enough content to scroll
	ticket := &domain.Ticket{
		JiraID:      "PROJ-789",
		Title:       "Test navigation",
		Description: strings.Repeat("Line\n", 50),
	}
	m.SetTicket(ticket)

	// Test down navigation
	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updatedModel

	// Test up navigation
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updatedModel

	// Test page down
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
	m = updatedModel

	// Test page up
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
	m = updatedModel

	// Test goto top
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	m = updatedModel

	// Test goto bottom
	updatedModel, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	m = updatedModel

	// If we get here without panic, navigation works
}

func TestSetSize(t *testing.T) {
	m := New(80, 24)

	newWidth, newHeight := 100, 30
	m.SetSize(newWidth, newHeight)

	if m.width != newWidth {
		t.Errorf("expected width %d, got %d", newWidth, m.width)
	}
	if m.height != newHeight {
		t.Errorf("expected height %d, got %d", newHeight, m.height)
	}
	if m.viewport.Width != newWidth {
		t.Errorf("expected viewport width %d, got %d", newWidth, m.viewport.Width)
	}
	if m.viewport.Height != newHeight {
		t.Errorf("expected viewport height %d, got %d", newHeight, m.viewport.Height)
	}
}

func TestView(t *testing.T) {
	m := New(80, 24)

	// Test view with no ticket
	view := m.View()
	if view == "" {
		t.Error("view should not be empty")
	}

	// Test view with ticket
	ticket := &domain.Ticket{
		JiraID: "PROJ-999",
		Title:  "Test view",
	}
	m.SetTicket(ticket)

	view = m.View()
	if view == "" {
		t.Error("view should not be empty with ticket")
	}
}
