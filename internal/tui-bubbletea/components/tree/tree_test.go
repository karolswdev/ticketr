package tree

import (
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// TestFlattenTickets tests the flattening of hierarchical tickets
func TestFlattenTickets(t *testing.T) {
	// Create sample tickets with tasks
	tickets := []domain.Ticket{
		{
			JiraID: "PROJ-1",
			Title:  "Parent Ticket 1",
			Tasks: []domain.Task{
				{
					JiraID: "PROJ-1-1",
					Title:  "Child Task 1",
				},
				{
					JiraID: "PROJ-1-2",
					Title:  "Child Task 2",
				},
			},
		},
		{
			JiraID: "PROJ-2",
			Title:  "Parent Ticket 2",
			Tasks:  []domain.Task{},
		},
	}

	// Test with no expanded state
	expandedState := make(map[string]bool)
	items := FlattenTickets(tickets, expandedState)

	// Should have 2 root items (children not expanded)
	if len(items) != 2 {
		t.Errorf("Expected 2 items with no expansion, got %d", len(items))
	}

	// Check root items
	if items[0].Level != 0 {
		t.Errorf("Expected root item at level 0, got %d", items[0].Level)
	}

	if items[0].HasKids != true {
		t.Error("Expected first item to have children")
	}

	if items[1].HasKids != false {
		t.Error("Expected second item to have no children")
	}
}

// TestFlattenTicketsExpanded tests flattening with expanded items
func TestFlattenTicketsExpanded(t *testing.T) {
	tickets := []domain.Ticket{
		{
			JiraID: "PROJ-1",
			Title:  "Parent Ticket 1",
			Tasks: []domain.Task{
				{
					JiraID: "PROJ-1-1",
					Title:  "Child Task 1",
				},
				{
					JiraID: "PROJ-1-2",
					Title:  "Child Task 2",
				},
			},
		},
	}

	// Expand the first ticket
	expandedState := map[string]bool{
		"PROJ-1": true,
	}

	items := FlattenTickets(tickets, expandedState)

	// Should have 3 items (1 parent + 2 children)
	if len(items) != 3 {
		t.Errorf("Expected 3 items with expansion, got %d", len(items))
	}

	// Check children
	if items[1].Level != 1 {
		t.Errorf("Expected child at level 1, got %d", items[1].Level)
	}

	if items[2].Level != 1 {
		t.Errorf("Expected child at level 1, got %d", items[2].Level)
	}

	if items[1].IsTask != true {
		t.Error("Expected second item to be a task")
	}

	if items[1].Parent == nil {
		t.Error("Expected child to have parent reference")
	}
}

// TestTreeModelBasics tests basic tree model operations
func TestTreeModelBasics(t *testing.T) {
	// Create a tree model
	treeModel := New(80, 24, &theme.DefaultTheme)

	if treeModel.width != 80 {
		t.Errorf("Expected width 80, got %d", treeModel.width)
	}

	if treeModel.height != 24 {
		t.Errorf("Expected height 24, got %d", treeModel.height)
	}

	// Set some tickets
	tickets := []domain.Ticket{
		{
			JiraID: "PROJ-1",
			Title:  "Test Ticket",
			Tasks:  []domain.Task{},
		},
	}

	treeModel.SetTickets(tickets)

	if len(treeModel.items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(treeModel.items))
	}

	if len(treeModel.visibleItems) != 1 {
		t.Errorf("Expected 1 visible item, got %d", len(treeModel.visibleItems))
	}
}

// TestTreeItemFilterValue tests the FilterValue method
func TestTreeItemFilterValue(t *testing.T) {
	// Test with ticket
	ticket := &domain.Ticket{
		Title: "Test Ticket",
	}

	item := TreeItem{
		Ticket: ticket,
	}

	if item.FilterValue() != "Test Ticket" {
		t.Errorf("Expected filter value 'Test Ticket', got '%s'", item.FilterValue())
	}

	// Test with task
	task := &domain.Task{
		Title: "Test Task",
	}

	taskItem := TreeItem{
		TaskData: task,
		IsTask:   true,
	}

	if taskItem.FilterValue() != "Test Task" {
		t.Errorf("Expected filter value 'Test Task', got '%s'", taskItem.FilterValue())
	}
}

// TestPerformanceWithLargeDataset tests tree performance with many tickets
func TestPerformanceWithLargeDataset(t *testing.T) {
	// Create 100 tickets with 10 tasks each (1,000 total items)
	tickets := make([]domain.Ticket, 100)
	for i := 0; i < 100; i++ {
		tasks := make([]domain.Task, 10)
		for j := 0; j < 10; j++ {
			tasks[j] = domain.Task{
				JiraID: "TASK",
				Title:  "Task",
			}
		}
		tickets[i] = domain.Ticket{
			JiraID: "TICKET",
			Title:  "Ticket",
			Tasks:  tasks,
		}
	}

	// Flatten without expansion (should be fast)
	expandedState := make(map[string]bool)
	items := FlattenTickets(tickets, expandedState)

	// Should only have root items
	if len(items) != 100 {
		t.Errorf("Expected 100 items, got %d", len(items))
	}

	// Now expand all (1,100 total items)
	for i := range tickets {
		expandedState[tickets[i].JiraID] = true
	}

	itemsExpanded := FlattenTickets(tickets, expandedState)

	// Should have all items
	if len(itemsExpanded) != 1100 {
		t.Errorf("Expected 1100 items when expanded, got %d", len(itemsExpanded))
	}
}
