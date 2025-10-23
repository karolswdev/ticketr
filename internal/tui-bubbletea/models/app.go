// Package models contains the main Bubbletea model and application state.
package models

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/components"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

const (
	minWidth  = 80
	minHeight = 24
)

// FocusPanel represents which panel is focused
type FocusPanel int

const (
	FocusLeft FocusPanel = iota
	FocusRight
)

// AppModel is the main Bubbletea model
type AppModel struct {
	// Window dimensions
	width  int
	height int

	// Components
	header    *components.Header
	leftPanel *components.Panel
	rightPanel *components.Panel
	actionBar *components.ActionBar

	// State
	focus       FocusPanel
	loading     bool
	demoMode    bool
	demoCounter int

	// Error handling
	error string
}

// NewAppModel creates a new app model
func NewAppModel(demoMode bool) *AppModel {
	header := components.NewHeader()

	leftPanel := components.NewPanel("Workspace & Tickets")
	leftPanel.SetFocused(true)
	leftPanel.SetHelp("[Shift+F3: Filter] [/: Search]")

	rightPanel := components.NewPanel("Ticket Detail")
	rightPanel.SetHelp("[e: Edit] [c: Comment] [s: Sync]")

	return &AppModel{
		header:     header,
		leftPanel:  leftPanel,
		rightPanel: rightPanel,
		actionBar:  components.NewActionBar(),
		focus:      FocusLeft,
		loading:    true,
		demoMode:   demoMode,
	}
}

// Init implements tea.Model
func (m *AppModel) Init() tea.Cmd {
	if m.demoMode {
		return tea.Batch(
			tickCmd(),
			demoTickCmd(),
		)
	}
	return tea.Batch(
		tickCmd(),
		loadingDoneCmd(),
	)
}

// Update implements tea.Model
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		m.header.Update()
		return m, tickCmd()

	case demoTickMsg:
		if m.demoMode {
			m.demoCounter++
			// Cycle themes every 3 seconds (30 ticks at 100ms)
			if m.demoCounter%30 == 0 {
				theme.Next(&theme.DefaultTheme)
			}
			// Simulate sync progress
			progress := (m.demoCounter % 100)
			m.header.SetSyncing(progress < 80, progress)

			return m, demoTickCmd()
		}
		return m, nil

	case loadingDoneMsg:
		m.loading = false
		return m, nil
	}

	return m, nil
}

// View implements tea.Model
func (m *AppModel) View() string {
	if m.width < minWidth || m.height < minHeight {
		return m.renderTooSmall()
	}

	if m.loading {
		return m.renderLoading()
	}

	return m.renderMain()
}

func (m *AppModel) renderTooSmall() string {
	th := &theme.DefaultTheme

	style := lipgloss.NewStyle().
		Foreground(th.Error).
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(th.Error).
		Padding(1, 2).
		Align(lipgloss.Center)

	content := fmt.Sprintf(
		"Terminal too small!\n\nCurrent: %d×%d\nMinimum: %d×%d\n\nPlease resize your terminal.",
		m.width, m.height,
		minWidth, minHeight,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		style.Render(content),
	)
}

func (m *AppModel) renderLoading() string {
	th := &theme.DefaultTheme

	style := lipgloss.NewStyle().
		Foreground(th.Primary).
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(th.Primary).
		Padding(1, 2)

	spinner := components.NewSpinner()
	content := fmt.Sprintf("\n%s Loading Ticketr TUI...\n", spinner.Render())

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		style.Render(content),
	)
}

func (m *AppModel) renderMain() string {
	// Calculate component heights
	headerHeight := 5
	actionBarHeight := 2
	panelsHeight := m.height - headerHeight - actionBarHeight - 1

	// Set component sizes
	m.header.Width = m.width
	m.actionBar.Width = m.width

	// Update panel content
	m.updatePanelContent()

	// Render header
	m.header.Width = m.width
	header := m.header.Render()

	// Render panels in a row
	panelWidth := m.width / 2
	m.leftPanel.SetSize(panelWidth, panelsHeight)
	m.rightPanel.SetSize(panelWidth, panelsHeight)

	leftContent := m.leftPanel.RenderWithHelp()
	rightContent := m.rightPanel.RenderWithHelp()

	panels := lipgloss.JoinHorizontal(lipgloss.Top, leftContent, rightContent)

	// Render action bar
	actionBar := m.actionBar.Render()

	// Combine all components
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		panels,
		actionBar,
	)
}

func (m *AppModel) updatePanelContent() {
	// Left panel: Workspace & Tickets
	leftContent := m.renderWorkspacePanel()
	m.leftPanel.SetContent(leftContent)

	// Right panel: Ticket Detail
	rightContent := m.renderTicketDetail()
	m.rightPanel.SetContent(rightContent)

	// Update focus
	m.leftPanel.SetFocused(m.focus == FocusLeft)
	m.rightPanel.SetFocused(m.focus == FocusRight)
}

func (m *AppModel) renderWorkspacePanel() string {
	th := &theme.DefaultTheme

	titleStyle := lipgloss.NewStyle().
		Foreground(th.Primary).
		Bold(true)

	itemStyle := lipgloss.NewStyle().
		Foreground(th.Foreground)

	selectedStyle := lipgloss.NewStyle().
		Foreground(th.Accent).
		Bold(true)

	var lines []string
	lines = append(lines, titleStyle.Render("📁 PROJ-123 (My Project)"))
	lines = append(lines, "")
	lines = append(lines, titleStyle.Render("🎫 Tickets (234)"))

	if m.focus == FocusLeft {
		lines = append(lines, selectedStyle.Render(" ▶ 📋 PROJ-1: Setup project"))
	} else {
		lines = append(lines, itemStyle.Render("   📋 PROJ-1: Setup project"))
	}

	lines = append(lines, itemStyle.Render("   🔧 PROJ-2: Fix authentication"))
	lines = append(lines, itemStyle.Render("   ✨ PROJ-3: Add new feature"))
	lines = append(lines, itemStyle.Render("   🐛 PROJ-4: Bug in login flow"))
	lines = append(lines, itemStyle.Render("   📝 PROJ-5: Update documentation"))

	return strings.Join(lines, "\n")
}

func (m *AppModel) renderTicketDetail() string {
	th := &theme.DefaultTheme

	titleStyle := lipgloss.NewStyle().
		Foreground(th.Primary).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(th.Secondary).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(th.Foreground)

	var lines []string
	lines = append(lines, titleStyle.Render("PROJ-2: Fix authentication"))
	lines = append(lines, labelStyle.Render("Type: ")+" "+valueStyle.Render("Bug | ")+labelStyle.Render("Priority: ")+valueStyle.Render("High | ")+labelStyle.Render("Status: ")+valueStyle.Render("Open"))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("Description:"))
	lines = append(lines, valueStyle.Render("Authentication is broken for OAuth users."))
	lines = append(lines, valueStyle.Render("Need to update token refresh logic."))
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("Assignee: ")+valueStyle.Render("John Doe"))
	lines = append(lines, labelStyle.Render("Created: ")+valueStyle.Render("2025-01-20"))
	lines = append(lines, labelStyle.Render("Updated: ")+valueStyle.Render("2 hours ago"))

	return strings.Join(lines, "\n")
}

func (m *AppModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "tab":
		// Switch focus between panels
		if m.focus == FocusLeft {
			m.focus = FocusRight
		} else {
			m.focus = FocusLeft
		}
		return m, nil

	case "1":
		// TODO: Update when theme system is properly integrated; theme.GetByName("Default")
		return m, nil

	case "2":
		// TODO: Update when theme system is properly integrated; theme.GetByName("Dark")
		return m, nil

	case "3":
		// TODO: Update when theme system is properly integrated; theme.GetByName("Arctic")
		return m, nil

	case "?":
		// TODO: Show help modal
		return m, nil
	}

	return m, nil
}

// Messages

type tickMsg time.Time
type demoTickMsg time.Time
type loadingDoneMsg struct{}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func demoTickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return demoTickMsg(t)
	})
}

func loadingDoneCmd() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return loadingDoneMsg{}
	})
}
