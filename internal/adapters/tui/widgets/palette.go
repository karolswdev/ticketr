package widgets

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/karolswdev/ticktr/internal/adapters/tui/commands"
	"github.com/rivo/tview"
)

// CommandPalette is an enhanced command palette widget that uses the command registry.
type CommandPalette struct {
	modal       *tview.Flex
	inputField  *tview.InputField
	commandList *tview.List
	statusBar   *tview.TextView

	registry        *commands.Registry
	filteredCmds    []*commands.Command
	onClose         func()
	showKeybindings bool
}

// NewCommandPalette creates a new command palette widget.
func NewCommandPalette(registry *commands.Registry) *CommandPalette {
	cp := &CommandPalette{
		registry:        registry,
		showKeybindings: true,
	}

	// Create input field for fuzzy search
	cp.inputField = tview.NewInputField().
		SetLabel(" ⌘ ").
		SetFieldWidth(0).
		SetPlaceholder("Type to search commands...")

	cp.inputField.SetBorder(true).
		SetTitle(" Command Palette (Ctrl+P or F1) ").
		SetBorderColor(tcell.ColorGreen)

	// Create commands list
	cp.commandList = tview.NewList().
		ShowSecondaryText(true)

	cp.commandList.SetBorder(true).
		SetTitle(" Available Commands ").
		SetBorderColor(tcell.ColorWhite)

	// Create status bar
	cp.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetText("[green]↑↓[white] Navigate  [green]Enter[white] Execute  [red]Esc[white] Cancel  [yellow]Tab[white] Toggle Keybindings")

	cp.statusBar.SetBorder(false)

	// Create modal layout
	cp.modal = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(cp.inputField, 3, 0, true).
		AddItem(cp.commandList, 0, 1, false).
		AddItem(cp.statusBar, 1, 0, false)

	// Set up input field change handler for fuzzy search
	cp.inputField.SetChangedFunc(func(text string) {
		cp.updateCommandList(text)
	})

	// Set up input field key handler
	cp.inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			cp.close()
			return nil
		case tcell.KeyDown:
			// Move focus to command list
			return event
		case tcell.KeyEnter:
			// Execute first filtered command if any
			if len(cp.filteredCmds) > 0 {
				cp.executeCommand(0)
			}
			return nil
		case tcell.KeyTab:
			// Toggle keybinding display
			cp.showKeybindings = !cp.showKeybindings
			cp.updateCommandList(cp.inputField.GetText())
			return nil
		}
		return event
	})

	// Set up command list key handler
	cp.commandList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			cp.close()
			return nil
		case tcell.KeyTab:
			// Toggle keybinding display
			cp.showKeybindings = !cp.showKeybindings
			cp.updateCommandList(cp.inputField.GetText())
			return nil
		}
		return event
	})

	// Set up command selection handler
	cp.commandList.SetSelectedFunc(func(index int, _ string, _ string, _ rune) {
		if index < len(cp.filteredCmds) {
			cp.executeCommand(index)
		}
	})

	return cp
}

// Show displays the command palette.
func (cp *CommandPalette) Show() {
	cp.inputField.SetText("")
	cp.updateCommandList("")
}

// SetOnClose sets the callback for when the palette is closed.
func (cp *CommandPalette) SetOnClose(callback func()) {
	cp.onClose = callback
}

// Primitive returns the underlying tview primitive.
func (cp *CommandPalette) Primitive() tview.Primitive {
	return cp.modal
}

// updateCommandList updates the command list based on search query.
func (cp *CommandPalette) updateCommandList(query string) {
	cp.commandList.Clear()
	cp.filteredCmds = cp.registry.Search(query)

	if len(cp.filteredCmds) == 0 {
		cp.commandList.SetTitle(" No Commands Found ")
		return
	}

	// Group by category for better organization
	categoryMap := make(map[commands.Category][]*commands.Command)
	for _, cmd := range cp.filteredCmds {
		categoryMap[cmd.Category] = append(categoryMap[cmd.Category], cmd)
	}

	// Display commands grouped by category
	categories := []commands.Category{
		commands.CategoryNav,
		commands.CategoryView,
		commands.CategoryEdit,
		commands.CategorySync,
		commands.CategorySystem,
	}

	for _, cat := range categories {
		cmds := categoryMap[cat]
		if len(cmds) == 0 {
			continue
		}

		// Add category header (non-selectable)
		cp.commandList.AddItem(fmt.Sprintf("─── %s ───", cat), "", 0, nil)

		// Add commands in this category
		for _, cmd := range cmds {
			mainText := cmd.Name
			secondaryText := cmd.Description

			// Add keybinding info if enabled
			if cp.showKeybindings && cmd.Keybinding != "" {
				secondaryText = fmt.Sprintf("[%s] %s", cmd.Keybinding, cmd.Description)
			}

			cp.commandList.AddItem(mainText, secondaryText, 0, nil)
		}
	}

	// Update title with count
	title := fmt.Sprintf(" Commands (%d) ", len(cp.filteredCmds))
	if query != "" {
		title = fmt.Sprintf(" Filtered Commands (%d) ", len(cp.filteredCmds))
	}
	cp.commandList.SetTitle(title)
}

// executeCommand executes the command at the given filtered index.
func (cp *CommandPalette) executeCommand(index int) {
	if index < 0 || index >= len(cp.filteredCmds) {
		return
	}

	cmd := cp.filteredCmds[index]

	// Close the palette first
	cp.close()

	// Execute the command
	if cmd.Handler != nil {
		_ = cmd.Handler() // Errors are handled by the command itself
	}
}

// close closes the command palette.
func (cp *CommandPalette) close() {
	if cp.onClose != nil {
		cp.onClose()
	}
}
