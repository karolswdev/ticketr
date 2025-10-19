package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Command represents an executable command in the palette.
type Command struct {
	Name        string
	Description string
	Action      func() error
}

// CommandPaletteView implements the View interface for command palette functionality.
type CommandPaletteView struct {
	modal       *tview.Flex
	inputField  *tview.InputField
	commandList *tview.List
	statusBar   *tview.TextView

	// Data
	commands        []Command
	filteredIndices []int

	// Callbacks
	onClose func()
}

// NewCommandPaletteView creates a new command palette view.
func NewCommandPaletteView() *CommandPaletteView {
	v := &CommandPaletteView{
		commands:        []Command{},
		filteredIndices: []int{},
	}

	// Create input field for command filter
	v.inputField = tview.NewInputField().
		SetLabel(": ").
		SetFieldWidth(0). // Full width
		SetPlaceholder("Type command name...")

	v.inputField.SetBorder(true).
		SetTitle(" Command Palette ").
		SetBorderColor(tcell.ColorGreen)

	// Create commands list
	v.commandList = tview.NewList().
		ShowSecondaryText(true)

	v.commandList.SetBorder(true).
		SetTitle(" Available Commands ").
		SetBorderColor(tcell.ColorWhite)

	// Create status bar
	v.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetText("[green]Enter[white]=Execute [yellow][white] [red]Esc[white]=Cancel")

	v.statusBar.SetBorder(false)

	// Create modal layout
	v.modal = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(v.inputField, 3, 0, true).
		AddItem(v.commandList, 0, 1, false).
		AddItem(v.statusBar, 1, 0, false)

	// Set up input field handler
	v.inputField.SetChangedFunc(func(text string) {
		v.filterCommands(text)
	})

	// Set up input field key handler
	v.inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			if v.onClose != nil {
				v.onClose()
			}
			return nil
		case tcell.KeyDown:
			// Focus would move to command list (handled by app)
			return event
		case tcell.KeyEnter:
			// Execute first filtered command if any
			if len(v.filteredIndices) > 0 {
				v.executeCommand(v.filteredIndices[0])
			}
			return nil
		}
		return event
	})

	// Set up command list handler
	v.commandList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			if v.onClose != nil {
				v.onClose()
			}
			return nil
		}
		return event
	})

	// Set up selection handler
	v.commandList.SetSelectedFunc(func(index int, _ string, _ string, _ rune) {
		if index < len(v.filteredIndices) {
			v.executeCommand(v.filteredIndices[index])
		}
	})

	return v
}

// Name returns the view's identifier.
func (v *CommandPaletteView) Name() string {
	return "command"
}

// Primitive returns the view's root primitive.
func (v *CommandPaletteView) Primitive() tview.Primitive {
	return v.modal
}

// OnShow is called when the view becomes active.
func (v *CommandPaletteView) OnShow() {
	// Clear filter and show all commands
	v.inputField.SetText("")
	v.filterCommands("")
}

// OnHide is called when the view is hidden.
func (v *CommandPaletteView) OnHide() {
	// Nothing to do
}

// SetCommands sets the available commands.
func (v *CommandPaletteView) SetCommands(commands []Command) {
	v.commands = commands
	v.filterCommands("")
}

// SetOnClose sets the callback for closing the command palette.
func (v *CommandPaletteView) SetOnClose(callback func()) {
	v.onClose = callback
}

// filterCommands filters the command list based on input.
func (v *CommandPaletteView) filterCommands(filter string) {
	v.commandList.Clear()
	v.filteredIndices = []int{}

	for i, cmd := range v.commands {
		// Simple substring matching (case-insensitive)
		if filter == "" || containsIgnoreCase(cmd.Name, filter) || containsIgnoreCase(cmd.Description, filter) {
			v.filteredIndices = append(v.filteredIndices, i)
			v.commandList.AddItem(cmd.Name, cmd.Description, 0, nil)
		}
	}

	// Update title with count
	title := " Available Commands "
	if filter != "" {
		title = " Filtered Commands "
	}
	v.commandList.SetTitle(title)
}

// executeCommand executes the command at the given index and closes the palette.
func (v *CommandPaletteView) executeCommand(index int) {
	if index >= 0 && index < len(v.commands) {
		cmd := v.commands[index]

		// Close the palette first
		if v.onClose != nil {
			v.onClose()
		}

		// Execute the command
		if cmd.Action != nil {
			_ = cmd.Action() // Errors handled by command itself
		}
	}
}

// containsIgnoreCase performs case-insensitive substring search.
func containsIgnoreCase(s, substr string) bool {
	s, substr = toLower(s), toLower(substr)
	return len(s) >= len(substr) && (s == substr || indexOf(s, substr) >= 0)
}

// toLower converts string to lowercase manually (avoid imports).
func toLower(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + 32
		} else {
			result[i] = r
		}
	}
	return string(result)
}

// indexOf finds the index of substr in s (-1 if not found).
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
