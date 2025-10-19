package views

import (
	"fmt"
	"strings"

	"github.com/karolswdev/ticktr/internal/adapters/tui/theme"
	"github.com/rivo/tview"
)

// ErrorModal displays error messages in a dismissible modal dialog.
type ErrorModal struct {
	modal   *tview.Modal
	app     *tview.Application
	onClose func()
}

// NewErrorModal creates a new error modal.
func NewErrorModal(app *tview.Application) *ErrorModal {
	modal := tview.NewModal()
	modal.SetBorder(true).SetTitle(" Error ")
	modal.SetBorderColor(theme.GetErrorColor())

	return &ErrorModal{
		modal: modal,
		app:   app,
	}
}

// Show displays the error modal with the given error and optional details.
func (e *ErrorModal) Show(err error, details string) {
	if err == nil {
		return
	}

	// Build error message
	var message strings.Builder
	message.WriteString(fmt.Sprintf("[red]Error:[-] %s\n\n", err.Error()))

	if details != "" {
		message.WriteString(fmt.Sprintf("[gray]Details:[-]\n%s\n\n", details))
	}

	message.WriteString("[yellow]Press OK or Esc to close[-]")

	e.modal.SetText(message.String())
	e.modal.ClearButtons()
	e.modal.AddButtons([]string{"OK"})
	e.modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if e.onClose != nil {
			e.onClose()
		}
	})

	// Show modal
	if e.app != nil {
		e.app.SetRoot(e.modal, true)
	}
}

// ShowSimple displays a simple error message without details.
func (e *ErrorModal) ShowSimple(message string) {
	e.modal.SetText(fmt.Sprintf("[red]%s[-]", message))
	e.modal.ClearButtons()
	e.modal.AddButtons([]string{"OK"})
	e.modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if e.onClose != nil {
			e.onClose()
		}
	})

	if e.app != nil {
		e.app.SetRoot(e.modal, true)
	}
}

// SetOnClose sets the callback for when the modal is dismissed.
func (e *ErrorModal) SetOnClose(callback func()) {
	e.onClose = callback
}

// Primitive returns the underlying tview primitive.
func (e *ErrorModal) Primitive() tview.Primitive {
	return e.modal
}
