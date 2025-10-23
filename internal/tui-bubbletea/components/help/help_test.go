package help

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

func TestNew(t *testing.T) {
	width, height := 80, 24
	th := &theme.DefaultTheme

	help := New(width, height, th)

	if help.width != width {
		t.Errorf("Expected width %d, got %d", width, help.width)
	}

	if help.height != height {
		t.Errorf("Expected height %d, got %d", height, help.height)
	}

	if help.visible {
		t.Error("Expected help to be hidden by default")
	}

	if help.theme != th {
		t.Errorf("Expected theme %v, got %v", th, help.theme)
	}
}

func TestShowHide(t *testing.T) {
	help := New(80, 24, &theme.DefaultTheme)

	// Initially hidden
	if help.IsVisible() {
		t.Error("Expected help to be hidden initially")
	}

	// Show
	help.Show()
	if !help.IsVisible() {
		t.Error("Expected help to be visible after Show()")
	}

	// Hide
	help.Hide()
	if help.IsVisible() {
		t.Error("Expected help to be hidden after Hide()")
	}
}

func TestToggle(t *testing.T) {
	help := New(80, 24, &theme.DefaultTheme)

	// Initially hidden, toggle should show
	help.Toggle()
	if !help.IsVisible() {
		t.Error("Expected help to be visible after first toggle")
	}

	// Toggle again should hide
	help.Toggle()
	if help.IsVisible() {
		t.Error("Expected help to be hidden after second toggle")
	}
}

func TestSetSize(t *testing.T) {
	help := New(80, 24, &theme.DefaultTheme)

	newWidth, newHeight := 100, 30
	help.SetSize(newWidth, newHeight)

	if help.width != newWidth {
		t.Errorf("Expected width %d, got %d", newWidth, help.width)
	}

	if help.height != newHeight {
		t.Errorf("Expected height %d, got %d", newHeight, help.height)
	}
}

func TestSetTheme(t *testing.T) {
	help := New(80, 24, &theme.DefaultTheme)

	newTheme := &theme.DarkTheme
	help.SetTheme(newTheme)

	if help.theme != newTheme {
		t.Errorf("Expected theme %v, got %v", newTheme, help.theme)
	}
}

func TestViewWhenHidden(t *testing.T) {
	help := New(80, 24, &theme.DefaultTheme)

	view := help.View()

	if view != "" {
		t.Error("Expected empty view when help is hidden")
	}
}

func TestViewWhenVisible(t *testing.T) {
	help := New(80, 24, &theme.DefaultTheme)
	help.Show()

	view := help.View()

	if view == "" {
		t.Error("Expected non-empty view when help is visible")
	}

	// Should contain keyboard shortcuts
	expectedContent := []string{
		"KEYBOARD SHORTCUTS",
		"NAVIGATION",
		"ACTIONS",
		"THEMES",
		"Tab",
		"Switch focus",
		"Quit",
	}

	for _, content := range expectedContent {
		if !strings.Contains(view, content) {
			t.Errorf("Expected view to contain %q", content)
		}
	}
}

func TestUpdateWithKeyMessages(t *testing.T) {
	help := New(80, 24, &theme.DefaultTheme)
	help.Show()

	testCases := []struct {
		key            string
		shouldBeHidden bool
	}{
		{"?", true},
		{"esc", true},
		{"q", true},
	}

	for _, tc := range testCases {
		help.Show() // Reset to visible
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tc.key)}
		if tc.key == "esc" {
			msg = tea.KeyMsg{Type: tea.KeyEsc}
		}

		help, _ = help.Update(msg)

		if tc.shouldBeHidden && help.IsVisible() {
			t.Errorf("Expected help to be hidden after pressing %q", tc.key)
		}
	}
}

func TestUpdateWithShowHideMessages(t *testing.T) {
	help := New(80, 24, &theme.DefaultTheme)

	// Test ShowHelpMsg
	help, _ = help.Update(ShowHelpMsg{})
	if !help.IsVisible() {
		t.Error("Expected help to be visible after ShowHelpMsg")
	}

	// Test HideHelpMsg
	help, _ = help.Update(HideHelpMsg{})
	if help.IsVisible() {
		t.Error("Expected help to be hidden after HideHelpMsg")
	}
}

func TestUpdateWhenHidden(t *testing.T) {
	help := New(80, 24, &theme.DefaultTheme)

	// Send a key message when hidden - should not crash
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")}
	help, _ = help.Update(msg)

	// Should still be hidden
	if help.IsVisible() {
		t.Error("Expected help to remain hidden")
	}
}

func TestContentWithDifferentThemes(t *testing.T) {
	themes := []*theme.Theme{&theme.DefaultTheme, &theme.DarkTheme, &theme.ArcticTheme}

	for _, th := range themes {
		help := New(80, 24, th)
		help.Show()

		view := help.View()

		if view == "" {
			t.Errorf("Expected non-empty view for theme %s", th.Name)
		}

		// Content should be present regardless of theme
		if !strings.Contains(view, "NAVIGATION") {
			t.Errorf("Expected content for theme %s to contain NAVIGATION", th.Name)
		}
	}
}

func TestInit(t *testing.T) {
	help := New(80, 24, &theme.DefaultTheme)

	cmd := help.Init()

	// Init should return nil (no initial command needed)
	if cmd != nil {
		t.Error("Expected Init to return nil")
	}
}
