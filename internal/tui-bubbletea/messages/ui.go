package messages

// UI-related messages for state changes, focus, modals, etc.

// FocusChangedMsg is sent when focus changes between panels.
type FocusChangedMsg struct {
	From Focus
	To   Focus
}

// Focus represents which UI element has focus.
type Focus int

const (
	FocusNone Focus = iota
	FocusTree
	FocusDetail
	FocusSearch
	FocusCommandPalette
	FocusWorkspaceSelector
	FocusModal
)

// Next returns the next focus in tab order.
func (f Focus) Next() Focus {
	switch f {
	case FocusTree:
		return FocusDetail
	case FocusDetail:
		return FocusTree
	default:
		return FocusTree
	}
}

// ModalOpenedMsg is sent when a modal is opened.
type ModalOpenedMsg struct {
	ModalType string // "search", "command", "confirm", "workspace"
}

// ModalClosedMsg is sent when a modal is closed.
type ModalClosedMsg struct {
	ModalType string
}

// ThemeChangedMsg is sent when the theme is changed.
type ThemeChangedMsg struct {
	ThemeName string
}

// ErrorMsg is sent when an error occurs.
type ErrorMsg struct {
	Err error
}

// SuccessMsg is sent when an operation succeeds.
type SuccessMsg struct {
	Message string
}

// StatusMsg is sent to update the status bar.
type StatusMsg struct {
	Message string
}

// SearchModalOpenedMsg is sent when the search modal opens.
type SearchModalOpenedMsg struct{}

// SearchModalClosedMsg is sent when the search modal closes.
type SearchModalClosedMsg struct{}

// ActionExecuteRequestedMsg is sent when an action is selected for execution.
type ActionExecuteRequestedMsg struct {
	ActionID string
	Action   interface{} // *actions.Action (avoiding import cycle)
}

// Week 4 Day 1: Action execution messages

// ToggleHelpMsg toggles help screen visibility
type ToggleHelpMsg struct{}

// OpenSearchMsg opens the search modal
type OpenSearchMsg struct{}

// OpenCommandPaletteMsg opens the command palette
type OpenCommandPaletteMsg struct{}

// OpenWorkspaceSelectorMsg opens the workspace selector
type OpenWorkspaceSelectorMsg struct{}

// Navigation messages
type MoveCursorDownMsg struct{}
type MoveCursorUpMsg struct{}
type ExpandNodeMsg struct{}
type CollapseNodeMsg struct{}
type SwitchPanelMsg struct{}
type FocusLeftMsg struct{}
type FocusRightMsg struct{}

// View messages
type SelectTicketMsg struct{}
type RefreshDataMsg struct{}

// Theme messages
type CycleThemeMsg struct{}

type SetThemeMsg struct {
	ThemeName string
}
