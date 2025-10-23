package workspace

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/core/domain"
)

// WorkspaceSelectedMsg is sent when a workspace is selected from the selector.
// Week 3 Day 3: Added for proper workspace switching flow.
type WorkspaceSelectedMsg struct {
	Workspace domain.Workspace
}

// workspaceItem wraps a workspace configuration for list display.
type workspaceItem struct {
	workspace domain.Workspace
}

// Title returns the display title for the workspace.
func (i workspaceItem) Title() string {
	return i.workspace.Name
}

// Description returns the display description for the workspace.
func (i workspaceItem) Description() string {
	return i.workspace.ProjectKey + " • " + i.workspace.JiraURL
}

// FilterValue returns the value used for filtering.
func (i workspaceItem) FilterValue() string {
	return i.workspace.Name
}

// Model represents the workspace selector component.
// It displays a filterable list of workspaces and allows selection.
type Model struct {
	list     list.Model
	onSelect func(domain.Workspace)
	width    int
	height   int
}

// New creates a new workspace selector with the specified workspaces.
func New(workspaces []domain.Workspace, width, height int) Model {
	items := make([]list.Item, len(workspaces))
	for i, ws := range workspaces {
		items[i] = workspaceItem{workspace: ws}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#61AFEF")).
		Bold(true)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ABB2BF"))

	l := list.New(items, delegate, width, height)
	l.Title = "Select Workspace"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#98C379")).
		Bold(true).
		Padding(0, 0, 1, 0)

	return Model{
		list:   l,
		width:  width,
		height: height,
	}
}

// SetOnSelect sets the callback function to call when a workspace is selected.
func (m *Model) SetOnSelect(fn func(domain.Workspace)) {
	m.onSelect = fn
}

// SetSize updates the workspace selector dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height)
}

// Update handles workspace selector events.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if item := m.list.SelectedItem(); item != nil {
				ws := item.(workspaceItem).workspace
				if m.onSelect != nil {
					m.onSelect(ws)
				}
				// Week 3 Day 3: Send workspace selected message
				return m, func() tea.Msg {
					return WorkspaceSelectedMsg{Workspace: ws}
				}
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the workspace selector.
func (m Model) View() string {
	return m.list.View()
}

// Init initializes the workspace selector component.
func (m Model) Init() tea.Cmd {
	return nil
}
