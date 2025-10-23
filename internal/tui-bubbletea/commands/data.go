package commands

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/messages"
)

// LoadCurrentWorkspace loads the current workspace asynchronously.
// This command executes in the background and sends a CurrentWorkspaceLoadedMsg when complete.
func LoadCurrentWorkspace(ws *services.WorkspaceService) tea.Cmd {
	return func() tea.Msg {
		workspace, err := ws.Current()
		return messages.CurrentWorkspaceLoadedMsg{
			Workspace: workspace,
			Error:     err,
		}
	}
}

// LoadTickets loads tickets for a workspace asynchronously.
// This command executes in the background and sends a TicketsLoadedMsg when complete.
func LoadTickets(tq *services.TicketQueryService, workspaceID string) tea.Cmd {
	return func() tea.Msg {
		tickets, err := tq.ListByWorkspace(workspaceID)
		return messages.TicketsLoadedMsg{
			Tickets: tickets,
			Error:   err,
		}
	}
}

// LoadWorkspaces loads all workspaces asynchronously.
// This command executes in the background and sends a WorkspacesLoadedMsg when complete.
func LoadWorkspaces(ws *services.WorkspaceService) tea.Cmd {
	return func() tea.Msg {
		workspaces, err := ws.List()
		return messages.WorkspacesLoadedMsg{
			Workspaces: workspaces,
			Error:      err,
		}
	}
}
