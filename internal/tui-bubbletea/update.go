package tuibubbletea

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/actions"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/commands"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/messages"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/views/cmdpalette"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/views/search"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/views/workspace"
)

// max returns the maximum of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Update is the main update function that handles all messages in the TUI.
// It implements the Bubbletea Update method and is the heart of the TEA pattern:
// messages flow in, state changes are applied, and commands flow out.
//
// This function acts as a message router, dispatching messages to the appropriate
// handlers based on message type. Child components will handle their own updates,
// but global concerns (window sizing, quit, modal overlays) are handled here.
//
// Week 4 Day 1: Now routes to Week 3 components (search, command palette) and
// handles action execution messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Week 4 Day 1: Handle action execution requests from search/palette
	case search.ActionExecuteRequestedMsg:
		actx := m.buildActionContext()
		if msg.Action != nil {
			return m, msg.Action.Execute(actx)
		}
		return m, nil

	case cmdpalette.CommandExecutedMsg:
		actx := m.buildActionContext()
		if msg.Action != nil {
			return m, msg.Action.Execute(actx)
		}
		return m, nil

	// Week 4 Day 1: Handle action messages
	case messages.ToggleHelpMsg:
		if m.helpScreen.IsVisible() {
			m.helpScreen.Hide()
		} else {
			actx := m.buildActionContext()
			m.helpScreen.ShowWithContext(actx)
		}
		return m, nil

	case messages.OpenSearchMsg:
		var cmd tea.Cmd
		m.searchModal, cmd = m.searchModal.Open()
		m.contextManager.Push(actions.ContextSearch)

		// Update action context for search modal
		actx := m.buildActionContext()
		m.searchModal.SetActionContext(actx)
		return m, cmd

	case messages.OpenCommandPaletteMsg:
		var cmd tea.Cmd
		m.cmdPalette, cmd = m.cmdPalette.Open()
		m.contextManager.Push(actions.ContextCommandPalette)

		// Update action context for command palette
		actx := m.buildActionContext()
		m.cmdPalette.SetActionContext(actx)
		return m, cmd

	case search.SearchModalClosedMsg:
		m.contextManager.Pop()
		return m, nil

	case cmdpalette.CommandPaletteClosedMsg:
		m.contextManager.Pop()
		return m, nil

	case messages.OpenWorkspaceSelectorMsg:
		m.showWorkspaceModal = true
		m.focused = FocusWorkspace
		return m, nil

	case messages.SwitchPanelMsg:
		m.ToggleFocus()
		if m.focused == FocusLeft {
			m.ticketTree.Focus()
			m.detailView.Blur()
		} else {
			m.ticketTree.Blur()
			m.detailView.Focus()
		}
		return m, nil

	case messages.FocusLeftMsg:
		if !m.showWorkspaceModal {
			m.SetFocus(FocusLeft)
			m.ticketTree.Focus()
			m.detailView.Blur()
		}
		return m, nil

	case messages.FocusRightMsg:
		if !m.showWorkspaceModal {
			m.SetFocus(FocusRight)
			m.ticketTree.Blur()
			m.detailView.Focus()
		}
		return m, nil

	case messages.CycleThemeMsg:
		m.theme = theme.Next(m.theme)
		m.propagateTheme()
		return m, nil

	case messages.SetThemeMsg:
		m.theme = theme.GetByName(msg.ThemeName)
		m.propagateTheme()
		return m, nil
	// Week 2 Day 1: Handle workspace data loading
	case messages.CurrentWorkspaceLoadedMsg:
		if msg.Error != nil {
			m.loadError = msg.Error
			m.loadingWorkspaces = false
			return m, nil
		}
		m.currentWorkspace = msg.Workspace
		m.loadingWorkspaces = false

		// Load tickets for this workspace
		if msg.Workspace != nil {
			m.loadingTickets = true
			m.loadingSpinner.SetMessage("Loading tickets...")
			return m, commands.LoadTickets(m.ticketQuery, msg.Workspace.ID)
		}
		return m, nil

	case messages.TicketsLoadedMsg:
		if msg.Error != nil {
			m.loadError = msg.Error
			m.loadingTickets = false
			return m, nil
		}
		m.tickets = msg.Tickets
		m.loadingTickets = false
		m.dataLoaded = true

		// Week 2 Days 2-3: Update tree with loaded tickets
		m.ticketTree.SetTickets(msg.Tickets)
		m.ticketTree.Focus() // Focus tree after loading

		return m, nil

	case messages.WorkspacesLoadedMsg:
		if msg.Error != nil {
			m.loadError = msg.Error
			return m, nil
		}
		m.workspaces = msg.Workspaces
		// Update workspace selector with loaded workspaces
		m.workspaceSelector = workspace.New(msg.Workspaces, m.width/2, m.height/2)
		return m, nil

	// Week 3 Day 3: Handle workspace selection and switch
	case workspace.WorkspaceSelectedMsg:
		// Update current workspace
		m.currentWorkspace = &msg.Workspace

		// Clear current tickets and reset tree
		m.tickets = []domain.Ticket{}
		m.ticketTree.SetTickets([]domain.Ticket{})
		m.selectedTicket = nil
		m.detailView.SetTicket(nil)

		// Close modal and return focus to tree
		m.showWorkspaceModal = false
		m.focused = FocusLeft

		// Set loading state
		m.loadingTickets = true
		m.dataLoaded = false
		m.loadingSpinner.SetMessage("Loading tickets for " + msg.Workspace.Name + "...")

		// Reload tickets for the new workspace
		return m, commands.LoadTickets(m.ticketQuery, msg.Workspace.ID)

	case tea.WindowSizeMsg:
		// Handle window resize - this is the FIRST message we receive
		m.width = msg.Width
		m.height = msg.Height

		// Week 2 Day 2: Validate minimum terminal size (80x24)
		if msg.Width < 80 || msg.Height < 24 {
			m.terminalTooSmall = true
			return m, nil
		}
		m.terminalTooSmall = false

		// Update layout with new dimensions
		m.layout.Resize(msg.Width, msg.Height)

		if !m.ready {
			// First window size message - mark as ready and load data
			m.ready = true
			return m, tea.Batch(
				commands.LoadCurrentWorkspace(m.workspaceService),
				commands.LoadWorkspaces(m.workspaceService),
			)
		}

		// Week 2 Days 2-3: Propagate size to child components
		leftWidth, rightWidth, contentHeight := m.layout.GetPanelDimensions()
		m.ticketTree.SetSize(leftWidth-2, contentHeight-2)
		m.detailView.SetSize(rightWidth-2, contentHeight-2)
		m.workspaceSelector.SetSize(m.width/2, m.height/2)

		// Week 2 Day 2: Update help screen size
		// Help modal should be 80% of screen for better readability
		helpWidth := max(m.width*8/10, 60)
		helpHeight := max(m.height*8/10, 20)
		m.helpScreen.SetSize(helpWidth, helpHeight)

		// Week 4 Day 1: Update Week 3 component sizes
		m.searchModal.SetSize(msg.Width, msg.Height)
		m.cmdPalette.SetSize(msg.Width, msg.Height)

		return m, nil

	case tea.KeyMsg:
		// Week 4 Day 1: Route to Week 3 modals first (priority order)
		// 1. Search modal
		if m.searchModal.IsVisible() {
			var cmd tea.Cmd
			m.searchModal, cmd = m.searchModal.Update(msg)
			return m, cmd
		}

		// 2. Command palette
		if m.cmdPalette.IsVisible() {
			var cmd tea.Cmd
			m.cmdPalette, cmd = m.cmdPalette.Update(msg)
			return m, cmd
		}

		// 3. Help screen
		if m.helpScreen.IsVisible() {
			var cmd tea.Cmd
			m.helpScreen, cmd = m.helpScreen.Update(msg)
			return m, cmd
		}

		// 4. Workspace selector modal
		if m.showWorkspaceModal {
			switch msg.String() {
			case "esc", "ctrl+c":
				m.showWorkspaceModal = false
				m.focused = FocusLeft
				return m, nil
			default:
				var cmd tea.Cmd
				m.workspaceSelector, cmd = m.workspaceSelector.Update(msg)
				return m, cmd
			}
		}

		// Week 4 Day 1: Global keybindings (when no modal is active)
		switch msg.String() {
		case "/":
			// Open search modal
			var cmd tea.Cmd
			m.searchModal, cmd = m.searchModal.Open()
			m.contextManager.Push(actions.ContextSearch)
			actx := m.buildActionContext()
			m.searchModal.SetActionContext(actx)
			return m, cmd

		case ":", "ctrl+p":
			// Open command palette
			var cmd tea.Cmd
			m.cmdPalette, cmd = m.cmdPalette.Open()
			m.contextManager.Push(actions.ContextCommandPalette)
			actx := m.buildActionContext()
			m.cmdPalette.SetActionContext(actx)
			return m, cmd

		case "?":
			// Toggle help screen
			if m.helpScreen.IsVisible() {
				m.helpScreen.Hide()
			} else {
				actx := m.buildActionContext()
				m.helpScreen.ShowWithContext(actx)
			}
			return m, nil

		case "ctrl+c", "q":
			// Quit the application
			return m, tea.Quit

		case "W":
			// Open workspace selector
			m.showWorkspaceModal = true
			m.focused = FocusWorkspace
			return m, nil

		// Theme switching
		case "1":
			m.theme = theme.GetByName("Default")
			m.propagateTheme()
			return m, nil
		case "2":
			m.theme = theme.GetByName("Dark")
			m.propagateTheme()
			return m, nil
		case "3":
			m.theme = theme.GetByName("Arctic")
			m.propagateTheme()
			return m, nil
		case "t":
			m.theme = theme.Next(m.theme)
			m.propagateTheme()
			return m, nil

		// Focus switching
		case "tab":
			m.ToggleFocus()
			if m.focused == FocusLeft {
				m.ticketTree.Focus()
				m.detailView.Blur()
			} else {
				m.ticketTree.Blur()
				m.detailView.Focus()
			}
			return m, nil
		case "h":
			if !m.showWorkspaceModal {
				m.SetFocus(FocusLeft)
				m.ticketTree.Focus()
				m.detailView.Blur()
			}
			return m, nil
		case "l":
			if !m.showWorkspaceModal {
				m.SetFocus(FocusRight)
				m.ticketTree.Blur()
				m.detailView.Focus()
			}
			return m, nil
		}

		// Week 2 Days 2-3: Route to focused component
		switch m.focused {
		case FocusLeft:
			// Route to tree component
			var cmd tea.Cmd
			m.ticketTree, cmd = m.ticketTree.Update(msg)

			// Check if selection changed
			selected := m.ticketTree.GetSelected()
			if selected != nil && (m.selectedTicket == nil || m.selectedTicket.JiraID != selected.JiraID) {
				m.SetSelectedTicket(selected)
			}

			return m, cmd

		case FocusRight:
			var cmd tea.Cmd
			m.detailView, cmd = m.detailView.Update(msg)
			return m, cmd
		}

		return m, nil

	// Week 2 Day 2: Forward spinner tick messages to keep animation running
	default:
		var cmd tea.Cmd
		m.loadingSpinner, cmd = m.loadingSpinner.Update(msg)
		return m, cmd
	}

	// TODO(day6): Add custom message handlers
	// case ticketSelectedMsg:
	//     // Handle ticket selection from tree
	//     m.ticketDetail.LoadTicket(msg.ticketID)
	//     m.SetFocus(FocusRight)
	//     return m, loadTicketCmd(msg.ticketID)
	//
	// case workspaceChangedMsg:
	//     // Handle workspace switch
	//     m.currentWorkspace = msg.workspace
	//     return m, loadWorkspaceTicketsCmd(msg.workspace)
	//
	// case syncStartedMsg:
	//     // Handle sync operation start
	//     m.header.SetSyncStatus(SyncInProgress)
	//     return m, nil
	//
	// case syncCompletedMsg:
	//     // Handle sync operation completion
	//     m.header.SetSyncStatus(SyncSuccess)
	//     return m, tea.Batch(
	//         refreshTicketsCmd(),
	//         flashSuccessCmd(),
	//     )
}

// propagateTheme updates theme for all components.
// Week 4 Day 1: Now includes Week 3 components.
func (m *Model) propagateTheme() {
	m.loadingSpinner.SetTheme(m.theme)
	m.helpScreen.SetTheme(m.theme)
	m.ticketTree.SetTheme(m.theme)
	m.searchModal.SetTheme(m.theme)
	m.cmdPalette.SetTheme(m.theme)
}
