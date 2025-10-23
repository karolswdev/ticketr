package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// Panel represents a bordered container with a title
type Panel struct {
	Title    string
	Content  string
	Width    int
	Height   int
	Focused  bool
	ShowHelp bool
	HelpText string
}

// NewPanel creates a new panel
func NewPanel(title string) *Panel {
	return &Panel{
		Title:    title,
		Content:  "",
		Focused:  false,
		ShowHelp: false,
	}
}

// SetContent sets the panel content
func (p *Panel) SetContent(content string) {
	p.Content = content
}

// SetSize sets the panel dimensions
func (p *Panel) SetSize(width, height int) {
	p.Width = width
	p.Height = height
}

// SetFocused sets the focus state
func (p *Panel) SetFocused(focused bool) {
	p.Focused = focused
}

// SetHelp sets the help text
func (p *Panel) SetHelp(help string) {
	p.HelpText = help
	p.ShowHelp = help != ""
}

// Render renders the panel with appropriate border
func (p *Panel) Render() string {
	th := &theme.DefaultTheme

	// Choose border style based on focus
	var borderStyle lipgloss.Border
	var borderColor lipgloss.AdaptiveColor
	var titlePrefix string

	if p.Focused {
		borderStyle = th.BorderFocusStyle
		borderColor = th.BorderFocus
		titlePrefix = "▶ "
	} else {
		borderStyle = th.BorderStyle
		borderColor = th.Border
		titlePrefix = ""
	}

	// Create the style
	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColor).
		Width(p.Width - 2).  // Account for borders
		Height(p.Height - 2). // Account for borders
		Padding(0, 1)

	// Add title if present
	if p.Title != "" {
		fullTitle := titlePrefix + p.Title
		if p.Focused {
			style = style.BorderTop(true).
				BorderTopForeground(th.BorderFocus)
		}

		// Render title in panel header
		titleStyle := lipgloss.NewStyle().
			Foreground(th.Primary).
			Bold(true)

		styledTitle := titleStyle.Render(fullTitle)
		style = style.BorderTop(true)

		// Combine title with border
		rendered := style.Render(p.Content)
		lines := strings.Split(rendered, "\n")
		if len(lines) > 0 {
			// Replace first line with title
			borderChar := "─"
			if p.Focused {
				borderChar = "═"
			}
			titleLine := strings.Replace(lines[0], borderChar, styledTitle, 1)
			lines[0] = titleLine
			return strings.Join(lines, "\n")
		}
	}

	return style.Render(p.Content)
}

// RenderWithHelp renders the panel with help text at the bottom
func (p *Panel) RenderWithHelp() string {
	if !p.ShowHelp {
		return p.Render()
	}

	th := &theme.DefaultTheme

	// Calculate content height (leave room for help)
	contentHeight := p.Height - 3 // Border + help line

	// Render main content
	mainPanel := *p
	mainPanel.Height = contentHeight
	mainPanel.ShowHelp = false
	content := mainPanel.Render()

	// Render help text
	helpStyle := lipgloss.NewStyle().
		Foreground(th.Muted).
		Width(p.Width - 2).
		Align(lipgloss.Center)

	help := helpStyle.Render(p.HelpText)

	return lipgloss.JoinVertical(lipgloss.Left, content, help)
}
