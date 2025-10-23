package help

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/actions"
	"github.com/karolswdev/ticktr/internal/tui-bubbletea/theme"
)

// Test legacy constructor (backward compatibility)
func TestNewLegacy(t *testing.T) {
	width, height := 80, 24
	th := &theme.DefaultTheme

	help := NewLegacy(width, height, th)

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

	if help.registry != nil {
		t.Error("Expected registry to be nil in legacy mode")
	}

	if help.contextMgr != nil {
		t.Error("Expected contextMgr to be nil in legacy mode")
	}
}

// Test new constructor with action registry
func TestNew(t *testing.T) {
	width, height := 80, 24
	th := &theme.DefaultTheme
	reg := actions.NewRegistry()
	cm := actions.NewContextManager(actions.ContextTicketTree)

	help := New(width, height, th, reg, cm)

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

	if help.registry != reg {
		t.Error("Expected registry to be set")
	}

	if help.contextMgr != cm {
		t.Error("Expected contextMgr to be set")
	}
}

func TestShowHide(t *testing.T) {
	help := NewLegacy(80, 24, &theme.DefaultTheme)

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
	help := NewLegacy(80, 24, &theme.DefaultTheme)

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
	help := NewLegacy(80, 24, &theme.DefaultTheme)

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
	help := NewLegacy(80, 24, &theme.DefaultTheme)

	newTheme := &theme.DarkTheme
	help.SetTheme(newTheme)

	if help.theme != newTheme {
		t.Errorf("Expected theme %v, got %v", newTheme, help.theme)
	}
}

func TestViewWhenHidden(t *testing.T) {
	help := NewLegacy(80, 24, &theme.DefaultTheme)

	view := help.View()

	if view != "" {
		t.Error("Expected empty view when help is hidden")
	}
}

func TestViewWhenVisible_Legacy(t *testing.T) {
	help := NewLegacy(80, 24, &theme.DefaultTheme)
	help.Show()

	view := help.View()

	if view == "" {
		t.Error("Expected non-empty view when help is visible")
	}

	// Should contain keyboard shortcuts
	expectedContent := []string{
		"KEYBOARD SHORTCUTS",
		"NAVIGATION",
		"Tab",
		"Switch focus",
	}

	for _, content := range expectedContent {
		if !strings.Contains(view, content) {
			t.Errorf("Expected view to contain %q", content)
		}
	}
}

func TestUpdateWithKeyMessages(t *testing.T) {
	help := NewLegacy(80, 24, &theme.DefaultTheme)
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
	help := NewLegacy(80, 24, &theme.DefaultTheme)

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
	help := NewLegacy(80, 24, &theme.DefaultTheme)

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
		help := NewLegacy(80, 24, th)
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
	help := NewLegacy(80, 24, &theme.DefaultTheme)

	cmd := help.Init()

	// Init should return nil (no initial command needed)
	if cmd != nil {
		t.Error("Expected Init to return nil")
	}
}

// NEW: Context-aware tests

func TestHelp_ContextAwareness(t *testing.T) {
	reg := actions.NewRegistry()
	cm := actions.NewContextManager(actions.ContextTicketTree)

	// Register test actions for different contexts
	reg.Register(&actions.Action{
		ID:          "tree.nav.down",
		Name:        "Move Down",
		Description: "Move down in tree",
		Category:    actions.CategoryNavigation,
		Contexts:    []actions.Context{actions.ContextTicketTree},
		Keybindings: []actions.KeyPattern{{Key: "j"}, {Key: "down"}},
		Execute:     func(ctx *actions.ActionContext) tea.Cmd { return nil },
	})

	reg.Register(&actions.Action{
		ID:          "detail.edit",
		Name:        "Edit Ticket",
		Description: "Edit ticket details",
		Category:    actions.CategoryEdit,
		Contexts:    []actions.Context{actions.ContextTicketDetail},
		Keybindings: []actions.KeyPattern{{Key: "e"}},
		Execute:     func(ctx *actions.ActionContext) tea.Cmd { return nil },
	})

	help := New(80, 24, &theme.DefaultTheme, reg, cm)

	// Test in TicketTree context
	cm.Switch(actions.ContextTicketTree)
	help.Show()

	view := help.View()
	if !strings.Contains(view, "Move down in tree") {
		t.Error("Expected tree action to be shown in TicketTree context")
	}
	if strings.Contains(view, "Edit ticket details") {
		t.Error("Expected detail action to NOT be shown in TicketTree context")
	}

	// Test in TicketDetail context
	cm.Switch(actions.ContextTicketDetail)
	help.Show()

	view = help.View()
	if !strings.Contains(view, "Edit ticket details") {
		t.Error("Expected detail action to be shown in TicketDetail context")
	}
	if strings.Contains(view, "Move down in tree") {
		t.Error("Expected tree action to NOT be shown in TicketDetail context")
	}
}

func TestHelp_DynamicSections(t *testing.T) {
	reg := actions.NewRegistry()
	cm := actions.NewContextManager(actions.ContextGlobal)

	// Register actions in different categories
	reg.Register(&actions.Action{
		ID:          "nav.up",
		Name:        "Up",
		Description: "Navigate up",
		Category:    actions.CategoryNavigation,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "k"}},
		Execute:     func(ctx *actions.ActionContext) tea.Cmd { return nil },
	})

	reg.Register(&actions.Action{
		ID:          "view.expand",
		Name:        "Expand",
		Description: "Expand view",
		Category:    actions.CategoryView,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "x"}},
		Execute:     func(ctx *actions.ActionContext) tea.Cmd { return nil },
	})

	help := New(80, 24, &theme.DefaultTheme, reg, cm)
	help.Show()

	// Check that sections are properly organized
	if len(help.sections) < 2 {
		t.Errorf("Expected at least 2 sections, got %d", len(help.sections))
	}

	// Find sections by title
	var navSection, viewSection *HelpSection
	for i := range help.sections {
		if help.sections[i].Title == "Navigation" {
			navSection = &help.sections[i]
		}
		if help.sections[i].Title == "View" {
			viewSection = &help.sections[i]
		}
	}

	if navSection == nil {
		t.Error("Expected Navigation section to exist")
	}
	if viewSection == nil {
		t.Error("Expected View section to exist")
	}
}

func TestHelp_CategoryGrouping(t *testing.T) {
	reg := actions.NewRegistry()
	cm := actions.NewContextManager(actions.ContextGlobal)

	// Register multiple actions in same category
	reg.Register(&actions.Action{
		ID:          "nav.up",
		Name:        "Up",
		Description: "Navigate up",
		Category:    actions.CategoryNavigation,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "k"}},
		Execute:     func(ctx *actions.ActionContext) tea.Cmd { return nil },
	})

	reg.Register(&actions.Action{
		ID:          "nav.down",
		Name:        "Down",
		Description: "Navigate down",
		Category:    actions.CategoryNavigation,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "j"}},
		Execute:     func(ctx *actions.ActionContext) tea.Cmd { return nil },
	})

	help := New(80, 24, &theme.DefaultTheme, reg, cm)
	help.Show()

	// Find navigation section
	var navSection *HelpSection
	for i := range help.sections {
		if help.sections[i].Title == "Navigation" {
			navSection = &help.sections[i]
			break
		}
	}

	if navSection == nil {
		t.Fatal("Expected Navigation section to exist")
	}

	// Should have at least our 2 actions
	if len(navSection.Actions) < 2 {
		t.Errorf("Expected at least 2 actions in Navigation section, got %d", len(navSection.Actions))
	}

	// Verify both our test actions are present
	hasUp := false
	hasDown := false
	for _, action := range navSection.Actions {
		if action.Description == "Navigate up" {
			hasUp = true
		}
		if action.Description == "Navigate down" {
			hasDown = true
		}
	}

	if !hasUp {
		t.Error("Expected to find 'Navigate up' action")
	}
	if !hasDown {
		t.Error("Expected to find 'Navigate down' action")
	}
}

func TestHelp_MultipleKeybindings(t *testing.T) {
	reg := actions.NewRegistry()
	cm := actions.NewContextManager(actions.ContextGlobal)

	// Register action with multiple keybindings
	reg.Register(&actions.Action{
		ID:          "nav.down",
		Name:        "Down",
		Description: "Navigate down",
		Category:    actions.CategoryNavigation,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "j"}, {Key: "down"}},
		Execute:     func(ctx *actions.ActionContext) tea.Cmd { return nil },
	})

	help := New(80, 24, &theme.DefaultTheme, reg, cm)
	help.Show()

	view := help.View()

	// Should show both keybindings
	if !strings.Contains(view, "j") {
		t.Error("Expected view to contain 'j' keybinding")
	}
	if !strings.Contains(view, "↓") {
		t.Error("Expected view to contain '↓' keybinding")
	}
}

func TestHelp_EmptyContext(t *testing.T) {
	reg := actions.NewRegistry()
	cm := actions.NewContextManager(actions.ContextGlobal)

	// Register action in specific context
	reg.Register(&actions.Action{
		ID:          "tree.expand",
		Name:        "Expand",
		Description: "Expand node",
		Category:    actions.CategoryView,
		Contexts:    []actions.Context{actions.ContextTicketTree},
		Keybindings: []actions.KeyPattern{{Key: "x"}},
		Execute:     func(ctx *actions.ActionContext) tea.Cmd { return nil },
	})

	help := New(80, 24, &theme.DefaultTheme, reg, cm)

	// Switch to context with no actions
	cm.Switch(actions.ContextSearch)
	help.Show()

	// Should not crash, might show global actions or be empty
	view := help.View()
	if view == "" {
		t.Error("Expected non-empty view even with no context-specific actions")
	}
}

func TestHelp_ThemeAwareness(t *testing.T) {
	reg := actions.NewRegistry()
	cm := actions.NewContextManager(actions.ContextGlobal)

	reg.Register(&actions.Action{
		ID:          "test.action",
		Name:        "Test",
		Description: "Test action",
		Category:    actions.CategorySystem,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "t"}},
		Execute:     func(ctx *actions.ActionContext) tea.Cmd { return nil },
	})

	themes := []*theme.Theme{&theme.DefaultTheme, &theme.DarkTheme, &theme.ArcticTheme}

	for _, th := range themes {
		help := New(80, 24, th, reg, cm)
		help.Show()

		view := help.View()
		if view == "" {
			t.Errorf("Expected non-empty view for theme %s", th.Name)
		}

		// Content should be present regardless of theme
		if !strings.Contains(view, "Test action") {
			t.Errorf("Expected content for theme %s to contain test action", th.Name)
		}
	}
}

func TestHelp_Scrolling(t *testing.T) {
	reg := actions.NewRegistry()
	cm := actions.NewContextManager(actions.ContextGlobal)

	// Register many actions to trigger scrolling
	for i := 0; i < 30; i++ {
		reg.Register(&actions.Action{
			ID:          actions.ActionID("test.action." + string(rune(i))),
			Name:        "Test Action",
			Description: "Test action for scrolling",
			Category:    actions.CategoryNavigation,
			Contexts:    []actions.Context{actions.ContextGlobal},
			Keybindings: []actions.KeyPattern{{Key: string(rune('a' + i))}},
			Execute:     func(ctx *actions.ActionContext) tea.Cmd { return nil },
		})
	}

	help := New(80, 10, &theme.DefaultTheme, reg, cm) // Small height
	help.Show()

	view := help.View()

	// Should mention scrolling
	if !strings.Contains(view, "scroll") {
		t.Error("Expected scrolling help text for large content")
	}
}

func TestShowWithContext(t *testing.T) {
	reg := actions.NewRegistry()
	cm := actions.NewContextManager(actions.ContextGlobal)

	reg.Register(&actions.Action{
		ID:          "test.action",
		Name:        "Test",
		Description: "Test action",
		Category:    actions.CategorySystem,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "t"}},
		Execute:     func(ctx *actions.ActionContext) tea.Cmd { return nil },
		Predicate: func(actx *actions.ActionContext) bool {
			return len(actx.SelectedTickets) > 0
		},
	})

	help := New(80, 24, &theme.DefaultTheme, reg, cm)

	// Show with context that has selection
	actx := &actions.ActionContext{
		Context:         actions.ContextGlobal,
		SelectedTickets: []string{"TICKET-1"},
	}
	help.ShowWithContext(actx)

	if !help.IsVisible() {
		t.Error("Expected help to be visible after ShowWithContext")
	}

	view := help.View()
	if !strings.Contains(view, "Test action") {
		t.Error("Expected action to show when predicate is satisfied")
	}
}

func TestSetActionContext(t *testing.T) {
	reg := actions.NewRegistry()
	cm := actions.NewContextManager(actions.ContextGlobal)

	reg.Register(&actions.Action{
		ID:          "test.action",
		Name:        "Test",
		Description: "Test action",
		Category:    actions.CategorySystem,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{{Key: "t"}},
		Execute:     func(ctx *actions.ActionContext) tea.Cmd { return nil },
	})

	help := New(80, 24, &theme.DefaultTheme, reg, cm)

	actx := &actions.ActionContext{
		Context: actions.ContextGlobal,
		Width:   100,
		Height:  50,
	}

	help.SetActionContext(actx)

	if help.actionCtx != actx {
		t.Error("Expected action context to be set")
	}
}

func TestFormatKeyPattern(t *testing.T) {
	help := NewLegacy(80, 24, &theme.DefaultTheme)

	tests := []struct {
		pattern  actions.KeyPattern
		expected string
	}{
		{actions.KeyPattern{Key: "j"}, "j"},
		{actions.KeyPattern{Key: "enter"}, "Enter"},
		{actions.KeyPattern{Key: "up"}, "↑"},
		{actions.KeyPattern{Key: "down"}, "↓"},
		{actions.KeyPattern{Key: "left"}, "←"},
		{actions.KeyPattern{Key: "right"}, "→"},
		{actions.KeyPattern{Key: "s", Ctrl: true}, "Ctrl+s"},
		{actions.KeyPattern{Key: "c", Ctrl: true, Shift: true}, "Ctrl+Shift+c"},
		{actions.KeyPattern{Key: "f", Alt: true}, "Alt+f"},
		{actions.KeyPattern{Key: "space"}, "Space"},
		{actions.KeyPattern{Key: "tab"}, "Tab"},
		{actions.KeyPattern{Key: "esc"}, "Esc"},
	}

	for _, tt := range tests {
		result := help.formatKeyPattern(tt.pattern)
		if result != tt.expected {
			t.Errorf("formatKeyPattern(%+v) = %q, want %q", tt.pattern, result, tt.expected)
		}
	}
}

func TestFormatContextName(t *testing.T) {
	help := NewLegacy(80, 24, &theme.DefaultTheme)

	tests := []struct {
		context  actions.Context
		expected string
	}{
		{actions.ContextWorkspaceList, "Workspace Selector"},
		{actions.ContextTicketTree, "Ticket Tree"},
		{actions.ContextTicketDetail, "Ticket Detail"},
		{actions.ContextSearch, "Search"},
		{actions.ContextCommandPalette, "Command Palette"},
		{actions.ContextModal, "Modal"},
		{actions.ContextSyncing, "Syncing"},
		{actions.ContextHelp, "Help"},
		{actions.ContextGlobal, "Global"},
		{actions.Context("unknown"), "unknown"},
	}

	for _, tt := range tests {
		result := help.formatContextName(tt.context)
		if result != tt.expected {
			t.Errorf("formatContextName(%q) = %q, want %q", tt.context, result, tt.expected)
		}
	}
}

func TestGenerateFallbackSections(t *testing.T) {
	help := NewLegacy(80, 24, &theme.DefaultTheme)

	sections := help.generateFallbackSections()

	if len(sections) == 0 {
		t.Error("Expected fallback sections to be generated")
	}

	// Should have at least Navigation and System sections
	hasNavigation := false
	hasSystem := false

	for _, section := range sections {
		if section.Title == "Navigation" {
			hasNavigation = true
		}
		if section.Title == "System" {
			hasSystem = true
		}
	}

	if !hasNavigation {
		t.Error("Expected fallback sections to include Navigation")
	}
	if !hasSystem {
		t.Error("Expected fallback sections to include System")
	}
}

func TestActionBindingWithNoKeys(t *testing.T) {
	reg := actions.NewRegistry()
	cm := actions.NewContextManager(actions.ContextGlobal)

	// Register action with no keybindings
	reg.Register(&actions.Action{
		ID:          "test.nokeybing",
		Name:        "No Keys",
		Description: "Action with no keybindings",
		Category:    actions.CategorySystem,
		Contexts:    []actions.Context{actions.ContextGlobal},
		Keybindings: []actions.KeyPattern{},
		Execute:     func(ctx *actions.ActionContext) tea.Cmd { return nil },
	})

	help := New(80, 24, &theme.DefaultTheme, reg, cm)
	help.Show()

	// Should show action with placeholder for keys
	view := help.View()
	if !strings.Contains(view, "Action with no keybindings") {
		t.Error("Expected action with no keybindings to be shown")
	}
}

func TestContextSwitching(t *testing.T) {
	reg := actions.NewRegistry()
	cm := actions.NewContextManager(actions.ContextTicketTree)

	// Register actions for different contexts
	reg.Register(&actions.Action{
		ID:          "tree.action",
		Name:        "Tree Action",
		Description: "Tree-specific action",
		Category:    actions.CategoryNavigation,
		Contexts:    []actions.Context{actions.ContextTicketTree},
		Keybindings: []actions.KeyPattern{{Key: "t"}},
		Execute:     func(ctx *actions.ActionContext) tea.Cmd { return nil },
	})

	reg.Register(&actions.Action{
		ID:          "detail.action",
		Name:        "Detail Action",
		Description: "Detail-specific action",
		Category:    actions.CategoryView,
		Contexts:    []actions.Context{actions.ContextTicketDetail},
		Keybindings: []actions.KeyPattern{{Key: "d"}},
		Execute:     func(ctx *actions.ActionContext) tea.Cmd { return nil },
	})

	help := New(80, 24, &theme.DefaultTheme, reg, cm)

	// Show in tree context
	help.Show()
	view1 := help.View()

	// Switch context
	cm.Switch(actions.ContextTicketDetail)
	help.Show()
	view2 := help.View()

	// Views should be different
	if view1 == view2 {
		t.Error("Expected different help content for different contexts")
	}

	// First view should show tree action
	if !strings.Contains(view1, "Tree-specific") {
		t.Error("Expected tree view to show tree action")
	}

	// Second view should show detail action
	if !strings.Contains(view2, "Detail-specific") {
		t.Error("Expected detail view to show detail action")
	}
}
