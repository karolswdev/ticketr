package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// Spinner represents an animated loading spinner
type Spinner struct {
	frames []string
	frame  int
}

// NewSpinner creates a new braille spinner
func NewSpinner() *Spinner {
	return &Spinner{
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		frame:  0,
	}
}

// Update advances the spinner to the next frame
func (s *Spinner) Update() {
	s.frame = (s.frame + 1) % len(s.frames)
}

// Render renders the current spinner frame
func (s *Spinner) Render() string {
	th := &theme.DefaultTheme
	style := lipgloss.NewStyle().Foreground(th.Accent)
	return style.Render(s.frames[s.frame])
}
