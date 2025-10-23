package messages

import "github.com/karolswdev/ticktr/internal/core/domain"

// Workspace-related messages for switching and managing workspaces.

// CurrentWorkspaceLoadedMsg is sent when the current workspace is loaded.
type CurrentWorkspaceLoadedMsg struct {
	Workspace *domain.Workspace
	Error     error
}

// WorkspacesLoadedMsg is sent when the list of workspaces is loaded.
type WorkspacesLoadedMsg struct {
	Workspaces []domain.Workspace
	Error      error
}

// WorkspaceChangedMsg is sent when the active workspace changes.
type WorkspaceChangedMsg struct {
	WorkspaceID   string
	WorkspaceName string
}

// WorkspaceCreatedMsg is sent when a new workspace is created.
type WorkspaceCreatedMsg struct {
	WorkspaceID   string
	WorkspaceName string
}

// WorkspaceDeletedMsg is sent when a workspace is deleted.
type WorkspaceDeletedMsg struct {
	WorkspaceID string
}
