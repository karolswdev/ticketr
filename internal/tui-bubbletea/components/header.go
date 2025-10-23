package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// Header represents the top application header
type Header struct {
	AppName   string
	Version   string
	Workspace string
	Status    string
	Progress  int // 0-100 for sync progress
	Theme     string
	Width     int
	Syncing   bool
	Spinner   *Spinner
}

// NewHeader creates a new header
func NewHeader() *Header {
	return &Header{
		AppName:   "TICKETR",
		Version:   "v3.2.0-beta (Bubbletea)",
		Workspace: "PROJ-123",
		Status:    "Ready",
		Progress:  0,
		Theme:     "Default",
		Syncing:   false,
		Spinner:   NewSpinner(),
	}
}

// SetSyncing sets the syncing state
func (h *Header) SetSyncing(syncing bool, progress int) {
	h.Syncing = syncing
	h.Progress = progress
	if syncing {
		h.Status = fmt.Sprintf("Syncing (%d%%)", progress)
	} else {
		h.Status = "Ready"
	}
}

// Update updates the header spinner
func (h *Header) Update() {
	if h.Syncing {
		h.Spinner.Update()
	}
}

// Render renders the header
func (h *Header) Render() string {
	th := &theme.DefaultTheme
	h.Theme = th.Name

	// Create header style with borders
	headerStyle := lipgloss.NewStyle().
		Foreground(th.Foreground).
		Background(th.Background).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(th.Primary).
		Padding(0, 1).
		Width(h.Width - 2)

	// Left side: App name + version
	leftStyle := lipgloss.NewStyle().
		Foreground(th.Primary).
		Bold(true)
	left := leftStyle.Render(fmt.Sprintf("🎫 %s", h.AppName))

	versionStyle := lipgloss.NewStyle().
		Foreground(th.Muted)
	version := versionStyle.Render(h.Version)

	// Center: Workspace and status
	workspaceStyle := lipgloss.NewStyle().
		Foreground(th.Accent).
		Bold(true)
	workspace := workspaceStyle.Render(fmt.Sprintf("[Workspace: %s]", h.Workspace))

	statusStyle := lipgloss.NewStyle().
		Foreground(th.Secondary)

	var status string
	if h.Syncing {
		spinner := h.Spinner.Render()
		status = statusStyle.Render(fmt.Sprintf("Status: %s %s", spinner, h.Status))
	} else {
		status = statusStyle.Render(fmt.Sprintf("Status: ✓ %s", h.Status))
	}

	// Right side: Theme
	themeStyle := lipgloss.NewStyle().
		Foreground(th.Muted)
	themeText := themeStyle.Render(fmt.Sprintf("Theme: %s", h.Theme))

	// Build the header content
	line1 := lipgloss.JoinHorizontal(
		lipgloss.Left,
		left,
		"  ",
		version,
		strings.Repeat(" ", max(1, h.Width-lipgloss.Width(left)-lipgloss.Width(version)-lipgloss.Width(workspace)-10)),
		workspace,
	)

	line2 := lipgloss.JoinHorizontal(
		lipgloss.Left,
		status,
		strings.Repeat(" ", max(1, h.Width-lipgloss.Width(status)-lipgloss.Width(themeText)-10)),
		themeText,
	)

	content := lipgloss.JoinVertical(lipgloss.Left, line1, line2)

	return headerStyle.Render(content)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
