package workspace

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/karolswdev/ticktr/internal/core/domain"
)

func TestNew(t *testing.T) {
	workspaces := []domain.Workspace{
		{
			ID:         "ws1",
			Name:       "Workspace 1",
			ProjectKey: "PROJ1",
			JiraURL:    "https://test1.atlassian.net",
		},
		{
			ID:         "ws2",
			Name:       "Workspace 2",
			ProjectKey: "PROJ2",
			JiraURL:    "https://test2.atlassian.net",
		},
	}

	width, height := 60, 20
	m := New(workspaces, width, height)

	if m.width != width {
		t.Errorf("expected width %d, got %d", width, m.width)
	}
	if m.height != height {
		t.Errorf("expected height %d, got %d", height, m.height)
	}

	// Check that list contains the workspaces
	items := m.list.Items()
	if len(items) != len(workspaces) {
		t.Errorf("expected %d items, got %d", len(workspaces), len(items))
	}
}

func TestWorkspaceItem(t *testing.T) {
	ws := domain.Workspace{
		ID:         "test-ws",
		Name:       "Test Workspace",
		ProjectKey: "TEST",
		JiraURL:    "https://test.atlassian.net",
		CreatedAt:  time.Now(),
	}

	item := workspaceItem{workspace: ws}

	if item.Title() != "Test Workspace" {
		t.Errorf("expected title 'Test Workspace', got '%s'", item.Title())
	}

	desc := item.Description()
	if !strings.Contains(desc, "TEST") {
		t.Errorf("description should contain project key, got '%s'", desc)
	}
	if !strings.Contains(desc, "https://test.atlassian.net") {
		t.Errorf("description should contain JIRA URL, got '%s'", desc)
	}

	if item.FilterValue() != "Test Workspace" {
		t.Errorf("expected filter value 'Test Workspace', got '%s'", item.FilterValue())
	}
}

func TestSetOnSelect(t *testing.T) {
	workspaces := []domain.Workspace{
		{
			ID:         "ws1",
			Name:       "Workspace 1",
			ProjectKey: "PROJ1",
			JiraURL:    "https://test1.atlassian.net",
		},
	}

	m := New(workspaces, 60, 20)

	called := false
	var selectedWs domain.Workspace

	m.SetOnSelect(func(ws domain.Workspace) {
		called = true
		selectedWs = ws
	})

	// Simulate selecting the first item
	// Note: In a real scenario, we'd need to navigate and select via Update()
	// For this test, we'll verify the callback is set
	if m.onSelect == nil {
		t.Error("onSelect callback should be set")
	}

	// Call the callback directly to test
	m.onSelect(workspaces[0])

	if !called {
		t.Error("onSelect callback should have been called")
	}
	if selectedWs.ID != "ws1" {
		t.Errorf("expected selected workspace 'ws1', got '%s'", selectedWs.ID)
	}
}

func TestSetSize(t *testing.T) {
	workspaces := []domain.Workspace{
		{
			ID:         "ws1",
			Name:       "Workspace 1",
			ProjectKey: "PROJ1",
		},
	}

	m := New(workspaces, 60, 20)

	newWidth, newHeight := 80, 30
	m.SetSize(newWidth, newHeight)

	if m.width != newWidth {
		t.Errorf("expected width %d, got %d", newWidth, m.width)
	}
	if m.height != newHeight {
		t.Errorf("expected height %d, got %d", newHeight, m.height)
	}
}

func TestUpdate_Enter(t *testing.T) {
	workspaces := []domain.Workspace{
		{
			ID:         "ws1",
			Name:       "Workspace 1",
			ProjectKey: "PROJ1",
		},
	}

	m := New(workspaces, 60, 20)

	called := false
	var selectedWs domain.Workspace

	m.SetOnSelect(func(ws domain.Workspace) {
		called = true
		selectedWs = ws
	})

	// Simulate pressing enter
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := m.Update(enterMsg)
	m = updatedModel

	if !called {
		t.Error("onSelect should be called when enter is pressed")
	}
	if selectedWs.ID != "ws1" {
		t.Errorf("expected selected workspace 'ws1', got '%s'", selectedWs.ID)
	}
}

func TestUpdate_WindowSize(t *testing.T) {
	workspaces := []domain.Workspace{
		{
			ID:         "ws1",
			Name:       "Workspace 1",
			ProjectKey: "PROJ1",
		},
	}

	m := New(workspaces, 60, 20)

	// Simulate window resize
	resizeMsg := tea.WindowSizeMsg{Width: 100, Height: 40}
	updatedModel, _ := m.Update(resizeMsg)
	m = updatedModel

	if m.width != 100 {
		t.Errorf("expected width 100 after resize, got %d", m.width)
	}
	if m.height != 40 {
		t.Errorf("expected height 40 after resize, got %d", m.height)
	}
}

func TestView(t *testing.T) {
	workspaces := []domain.Workspace{
		{
			ID:         "ws1",
			Name:       "Test Workspace",
			ProjectKey: "TEST",
		},
	}

	m := New(workspaces, 60, 20)

	view := m.View()
	if view == "" {
		t.Error("view should not be empty")
	}

	// The view should contain workspace selector elements
	// (exact content depends on the list rendering)
}

func TestEmptyWorkspacesList(t *testing.T) {
	m := New([]domain.Workspace{}, 60, 20)

	items := m.list.Items()
	if len(items) != 0 {
		t.Errorf("expected 0 items for empty workspaces list, got %d", len(items))
	}

	view := m.View()
	if view == "" {
		t.Error("view should not be empty even with no workspaces")
	}
}
