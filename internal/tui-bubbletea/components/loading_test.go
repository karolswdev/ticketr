package components

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

func TestNewLoading(t *testing.T) {
	msg := "Test loading message"
	th := &theme.DefaultTheme

	loading := NewLoading(msg, th)

	if loading.message != msg {
		t.Errorf("Expected message %q, got %q", msg, loading.message)
	}

	if loading.theme != th {
		t.Errorf("Expected theme %v, got %v", th, loading.theme)
	}
}

func TestLoadingSetMessage(t *testing.T) {
	loading := NewLoading("Initial", &theme.DefaultTheme)
	newMsg := "Updated message"

	loading.SetMessage(newMsg)

	if loading.message != newMsg {
		t.Errorf("Expected message %q, got %q", newMsg, loading.message)
	}
}

func TestLoadingSetTheme(t *testing.T) {
	loading := NewLoading("Test", &theme.DefaultTheme)
	newTheme := &theme.DarkTheme

	loading.SetTheme(newTheme)

	if loading.theme != newTheme {
		t.Errorf("Expected theme %v, got %v", newTheme, loading.theme)
	}
}

func TestLoadingView(t *testing.T) {
	msg := "Loading data"
	loading := NewLoading(msg, &theme.DefaultTheme)

	view := loading.View()

	// Should contain the message
	if !strings.Contains(view, msg) {
		t.Errorf("Expected view to contain %q, got %q", msg, view)
	}

	// Should not be empty
	if view == "" {
		t.Error("Expected non-empty view")
	}
}

func TestLoadingUpdate(t *testing.T) {
	loading := NewLoading("Test", &theme.DefaultTheme)

	// Update with a tick message
	msg := tea.Msg(struct{}{})
	updated, cmd := loading.Update(msg)

	// Should return updated model
	if updated.message != loading.message {
		t.Error("Expected message to remain unchanged after update")
	}

	// Should return a command (spinner tick)
	// We can't easily test the cmd itself, but it should not panic
	_ = cmd
}

func TestLoadingInit(t *testing.T) {
	loading := NewLoading("Test", &theme.DefaultTheme)

	cmd := loading.Init()

	// Should return a command
	if cmd == nil {
		t.Error("Expected Init to return a command")
	}
}

func TestLoadingWithDifferentThemes(t *testing.T) {
	msg := "Test"
	themes := []*theme.Theme{&theme.DefaultTheme, &theme.DarkTheme, &theme.ArcticTheme}

	for _, th := range themes {
		loading := NewLoading(msg, th)
		view := loading.View()

		if !strings.Contains(view, msg) {
			t.Errorf("Expected view with theme %s to contain %q", th.Name, msg)
		}
	}
}
