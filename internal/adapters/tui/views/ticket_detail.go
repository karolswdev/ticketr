package views

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/karolswdev/ticktr/internal/adapters/tui/theme"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/rivo/tview"
)

var jiraIDPattern = regexp.MustCompile(`^[A-Z][A-Z0-9]+-\d+$`)

// TicketDetailView displays and edits ticket details
type TicketDetailView struct {
	container      *tview.Flex
	displayView    *tview.TextView
	editForm       *tview.Form
	statusBar      *tview.TextView
	currentTicket  *domain.Ticket
	originalTicket *domain.Ticket // For cancel operation
	editMode       bool
	isDirty        bool
	onSave         func(*domain.Ticket) error
	app            *tview.Application
}

// NewTicketDetailView creates a new ticket detail view
func NewTicketDetailView(app *tview.Application) *TicketDetailView {
	v := &TicketDetailView{
		container:   tview.NewFlex().SetDirection(tview.FlexRow),
		displayView: tview.NewTextView(),
		editForm:    tview.NewForm(),
		statusBar:   tview.NewTextView(),
		app:         app,
	}

	v.setupDisplayView()
	v.setupEditForm()
	v.setupStatusBar()
	v.setupLayout()
	v.setupKeybindings()

	return v
}

func (v *TicketDetailView) setupDisplayView() {
	v.displayView.
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true).
		SetBorder(true).
		SetTitle(" Ticket Detail (read-only) ").
		SetBorderColor(tcell.ColorWhite)

	v.displayView.SetText("[gray]No ticket selected. Select a ticket from the tree view.[white]")
}

func (v *TicketDetailView) setupEditForm() {
	v.editForm.
		SetBorder(true).
		SetTitle(" Edit Ticket ").
		SetBorderColor(tcell.ColorWhite)

	// Form will be populated dynamically when entering edit mode
}

func (v *TicketDetailView) setupStatusBar() {
	v.statusBar.
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetText("[gray]Keys: [green]e[gray] edit | [green]Esc[gray] back | [green]?[gray] help")
}

func (v *TicketDetailView) setupLayout() {
	// Initially show display view
	v.container.
		AddItem(v.displayView, 0, 1, true).
		AddItem(v.statusBar, 1, 0, false)
}

func (v *TicketDetailView) setupKeybindings() {
	// Display mode keybindings
	v.displayView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if v.currentTicket == nil {
			return event
		}

		switch event.Key() {
		case tcell.KeyCtrlF:
			// Page down (full page)
			_, _, _, height := v.displayView.GetInnerRect()
			row, col := v.displayView.GetScrollOffset()
			v.displayView.ScrollTo(row+height, col)
			return nil
		case tcell.KeyCtrlB:
			// Page up (full page)
			_, _, _, height := v.displayView.GetInnerRect()
			row, col := v.displayView.GetScrollOffset()
			newRow := row - height
			if newRow < 0 {
				newRow = 0
			}
			v.displayView.ScrollTo(newRow, col)
			return nil
		case tcell.KeyCtrlD:
			// Half-page down
			_, _, _, height := v.displayView.GetInnerRect()
			row, col := v.displayView.GetScrollOffset()
			v.displayView.ScrollTo(row+height/2, col)
			return nil
		case tcell.KeyCtrlU:
			// Half-page up
			_, _, _, height := v.displayView.GetInnerRect()
			row, col := v.displayView.GetScrollOffset()
			newRow := row - height/2
			if newRow < 0 {
				newRow = 0
			}
			v.displayView.ScrollTo(newRow, col)
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'e', 'E':
				v.enterEditMode()
				return nil
			case 'j':
				// Scroll down
				row, col := v.displayView.GetScrollOffset()
				v.displayView.ScrollTo(row+1, col)
				return nil
			case 'k':
				// Scroll up
				row, col := v.displayView.GetScrollOffset()
				v.displayView.ScrollTo(row-1, col)
				return nil
			case 'g':
				// Go to top
				v.displayView.ScrollToBeginning()
				return nil
			case 'G':
				// Go to bottom
				v.displayView.ScrollToEnd()
				return nil
			}
		}
		return event
	})

	// Edit form keybindings handled via form buttons and Cancel handler
}

// Name implements View interface
func (v *TicketDetailView) Name() string {
	return "ticket_detail"
}

// Primitive implements View interface
func (v *TicketDetailView) Primitive() tview.Primitive {
	return v.container
}

// OnShow implements View interface
func (v *TicketDetailView) OnShow() {
	// Reset to display mode when showing
	if v.editMode {
		v.exitEditMode(false)
	}
}

// OnHide implements View interface
func (v *TicketDetailView) OnHide() {
	// Discard unsaved changes when hiding
	if v.isDirty {
		v.isDirty = false
		v.currentTicket = v.copyTicket(v.originalTicket)
	}
}

// SetTicket updates the displayed ticket
func (v *TicketDetailView) SetTicket(ticket *domain.Ticket) {
	if ticket == nil {
		v.currentTicket = nil
		v.originalTicket = nil
		v.displayView.SetText("[gray]No ticket selected. Select a ticket from the tree view.[white]")
		v.updateStatusBar()
		return
	}

	v.currentTicket = v.copyTicket(ticket)
	v.originalTicket = v.copyTicket(ticket)
	v.isDirty = false
	v.renderDisplayView()
	v.updateStatusBar()
}

// SetOnSave sets the callback for saving ticket changes
func (v *TicketDetailView) SetOnSave(callback func(*domain.Ticket) error) {
	v.onSave = callback
}

func (v *TicketDetailView) renderDisplayView() {
	if v.currentTicket == nil {
		return
	}

	var b strings.Builder

	// Title
	b.WriteString(fmt.Sprintf("[yellow::b]%s[-:-:-]\n\n", v.currentTicket.Title))

	// JiraID
	if v.currentTicket.JiraID != "" {
		b.WriteString(fmt.Sprintf("[gray]Jira ID:[-] [green]%s[-]\n", v.currentTicket.JiraID))
	} else {
		b.WriteString("[gray]Jira ID:[-] [red]Not synced[-]\n")
	}

	// Source Line
	b.WriteString(fmt.Sprintf("[gray]Source Line:[-] %d\n\n", v.currentTicket.SourceLine))

	// Description
	if v.currentTicket.Description != "" {
		b.WriteString("[cyan::b]Description[-:-:-]\n")
		b.WriteString(v.formatMultilineText(v.currentTicket.Description))
		b.WriteString("\n\n")
	}

	// Custom Fields
	if len(v.currentTicket.CustomFields) > 0 {
		b.WriteString("[cyan::b]Custom Fields[-:-:-]\n")
		for key, value := range v.currentTicket.CustomFields {
			b.WriteString(fmt.Sprintf("  [gray]%s:[-] %s\n", key, value))
		}
		b.WriteString("\n")
	}

	// Acceptance Criteria
	if len(v.currentTicket.AcceptanceCriteria) > 0 {
		b.WriteString("[cyan::b]Acceptance Criteria[-:-:-]\n")
		for i, ac := range v.currentTicket.AcceptanceCriteria {
			b.WriteString(fmt.Sprintf("  [gray]%d.[-] %s\n", i+1, ac))
		}
		b.WriteString("\n")
	}

	// Tasks
	if len(v.currentTicket.Tasks) > 0 {
		b.WriteString(fmt.Sprintf("[cyan::b]Tasks[-:-:-] [gray](%d)[-]\n", len(v.currentTicket.Tasks)))
		for i, task := range v.currentTicket.Tasks {
			b.WriteString(fmt.Sprintf("\n  [yellow]Task %d:[-] %s\n", i+1, task.Title))
			if task.JiraID != "" {
				b.WriteString(fmt.Sprintf("  [gray]Jira ID:[-] [green]%s[-]\n", task.JiraID))
			}
			if task.Description != "" {
				b.WriteString(fmt.Sprintf("  [gray]Description:[-] %s\n", v.formatMultilineText(task.Description)))
			}
			if len(task.AcceptanceCriteria) > 0 {
				b.WriteString("  [gray]Acceptance Criteria:[-]\n")
				for j, ac := range task.AcceptanceCriteria {
					b.WriteString(fmt.Sprintf("    %d. %s\n", j+1, ac))
				}
			}
		}
	}

	v.displayView.SetText(b.String())
	v.displayView.ScrollToBeginning()
}

func (v *TicketDetailView) formatMultilineText(text string) string {
	lines := strings.Split(text, "\n")
	var formatted []string
	for _, line := range lines {
		if line == "" {
			formatted = append(formatted, "")
		} else {
			formatted = append(formatted, "  "+line)
		}
	}
	return strings.Join(formatted, "\n")
}

func (v *TicketDetailView) enterEditMode() {
	if v.currentTicket == nil {
		return
	}

	v.editMode = true
	v.populateEditForm()

	// Switch layout to show edit form
	v.container.Clear()
	v.container.
		AddItem(v.editForm, 0, 1, true).
		AddItem(v.statusBar, 1, 0, false)

	v.updateStatusBar()
	v.app.SetFocus(v.editForm)
}

func (v *TicketDetailView) populateEditForm() {
	v.editForm.Clear(true)

	// Title field
	titleField := tview.NewInputField().
		SetLabel("Title").
		SetText(v.currentTicket.Title).
		SetFieldWidth(0).
		SetChangedFunc(func(text string) {
			v.isDirty = true
		})

	// JiraID field
	jiraIDField := tview.NewInputField().
		SetLabel("Jira ID").
		SetText(v.currentTicket.JiraID).
		SetFieldWidth(20).
		SetChangedFunc(func(text string) {
			v.isDirty = true
		})

	// Description field
	descField := tview.NewTextArea().
		SetLabel("Description").
		SetText(v.currentTicket.Description, true).
		SetChangedFunc(func() {
			v.isDirty = true
		})

	v.editForm.
		AddFormItem(titleField).
		AddFormItem(jiraIDField).
		AddFormItem(descField)

	// Custom Fields (editable as key=value pairs)
	customFieldsText := v.formatCustomFieldsForEdit()
	customFieldsArea := tview.NewTextArea().
		SetLabel("Custom Fields (key=value per line)").
		SetText(customFieldsText, true).
		SetChangedFunc(func() {
			v.isDirty = true
		})
	v.editForm.AddFormItem(customFieldsArea)

	// Acceptance Criteria (editable as list)
	acText := strings.Join(v.currentTicket.AcceptanceCriteria, "\n")
	acArea := tview.NewTextArea().
		SetLabel("Acceptance Criteria (one per line)").
		SetText(acText, true).
		SetChangedFunc(func() {
			v.isDirty = true
		})
	v.editForm.AddFormItem(acArea)

	// Buttons
	v.editForm.
		AddButton("Save", v.saveTicket).
		AddButton("Cancel", func() {
			v.exitEditMode(false)
		})

	// Set cancel handler
	v.editForm.SetCancelFunc(func() {
		v.exitEditMode(false)
	})
}

func (v *TicketDetailView) formatCustomFieldsForEdit() string {
	var lines []string
	for key, value := range v.currentTicket.CustomFields {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(lines, "\n")
}

func (v *TicketDetailView) parseCustomFields(text string) map[string]string {
	fields := make(map[string]string)
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" {
				fields[key] = value
			}
		}
	}
	return fields
}

func (v *TicketDetailView) saveTicket() {
	// Extract values from form
	title := v.editForm.GetFormItemByLabel("Title").(*tview.InputField).GetText()
	jiraID := v.editForm.GetFormItemByLabel("Jira ID").(*tview.InputField).GetText()
	description := v.editForm.GetFormItemByLabel("Description").(*tview.TextArea).GetText()
	customFieldsText := v.editForm.GetFormItemByLabel("Custom Fields (key=value per line)").(*tview.TextArea).GetText()
	acText := v.editForm.GetFormItemByLabel("Acceptance Criteria (one per line)").(*tview.TextArea).GetText()

	// Build updated ticket
	updatedTicket := v.copyTicket(v.currentTicket)
	updatedTicket.Title = strings.TrimSpace(title)
	updatedTicket.JiraID = strings.TrimSpace(jiraID)
	updatedTicket.Description = description
	updatedTicket.CustomFields = v.parseCustomFields(customFieldsText)
	updatedTicket.AcceptanceCriteria = v.parseAcceptanceCriteria(acText)

	// Validate
	errors := v.validateTicket(updatedTicket)
	if len(errors) > 0 {
		v.showValidationErrors(errors)
		return
	}

	// Update current ticket
	v.currentTicket = updatedTicket
	v.originalTicket = v.copyTicket(updatedTicket)
	v.isDirty = false

	// Call save callback if provided
	if v.onSave != nil {
		if err := v.onSave(v.currentTicket); err != nil {
			v.showError(fmt.Sprintf("Save failed: %v", err))
			return
		}
	}

	// Exit edit mode and show updated display
	v.exitEditMode(true)
}

func (v *TicketDetailView) parseAcceptanceCriteria(text string) []string {
	var criteria []string
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			criteria = append(criteria, line)
		}
	}
	return criteria
}

func (v *TicketDetailView) validateTicket(ticket *domain.Ticket) []string {
	var errors []string

	// Required field: Title
	if ticket.Title == "" {
		errors = append(errors, "Title is required")
	}

	// JiraID format validation (if provided)
	if ticket.JiraID != "" && !jiraIDPattern.MatchString(ticket.JiraID) {
		errors = append(errors, "Jira ID must match format PROJECT-123 (uppercase letters, dash, numbers)")
	}

	return errors
}

func (v *TicketDetailView) showValidationErrors(errors []string) {
	title := "Validation Errors"
	message := strings.Join(errors, "\n")

	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			v.app.SetRoot(v.container, true)
			v.app.SetFocus(v.editForm)
		})

	modal.SetTitle(title).SetBorder(true)
	v.app.SetRoot(modal, true)
}

func (v *TicketDetailView) showError(message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			v.app.SetRoot(v.container, true)
			v.app.SetFocus(v.editForm)
		})

	modal.SetTitle("Error").SetBorder(true)
	v.app.SetRoot(modal, true)
}

func (v *TicketDetailView) exitEditMode(saved bool) {
	v.editMode = false

	if !saved {
		// Restore original ticket (cancel operation)
		v.currentTicket = v.copyTicket(v.originalTicket)
		v.isDirty = false
	}

	// Switch layout back to display view
	v.container.Clear()
	v.container.
		AddItem(v.displayView, 0, 1, true).
		AddItem(v.statusBar, 1, 0, false)

	v.renderDisplayView()
	v.updateStatusBar()
	v.app.SetFocus(v.displayView)
}

func (v *TicketDetailView) updateStatusBar() {
	if v.currentTicket == nil {
		v.statusBar.SetText("[gray]No ticket selected")
		return
	}

	if v.editMode {
		dirtyIndicator := ""
		if v.isDirty {
			dirtyIndicator = " [red]*[gray] (unsaved changes)"
		}
		v.statusBar.SetText(fmt.Sprintf("[gray]Edit Mode%s | Keys: [green]Ctrl+S[gray] save | [green]Esc[gray] cancel", dirtyIndicator))
	} else {
		v.statusBar.SetText("[gray]Keys: [green]e[gray] edit | [green]j/k[gray] scroll | [green]Ctrl+F/B[gray] page | [green]Ctrl+D/U[gray] half | [green]g/G[gray] top/bottom | [green]Esc[gray] back")
	}
}

func (v *TicketDetailView) copyTicket(ticket *domain.Ticket) *domain.Ticket {
	if ticket == nil {
		return nil
	}

	copy := &domain.Ticket{
		Title:              ticket.Title,
		Description:        ticket.Description,
		JiraID:             ticket.JiraID,
		SourceLine:         ticket.SourceLine,
		CustomFields:       make(map[string]string),
		AcceptanceCriteria: make([]string, len(ticket.AcceptanceCriteria)),
		Tasks:              make([]domain.Task, len(ticket.Tasks)),
	}

	// Deep copy maps and slices
	for k, v := range ticket.CustomFields {
		copy.CustomFields[k] = v
	}
	copy.AcceptanceCriteria = append(copy.AcceptanceCriteria[:0], ticket.AcceptanceCriteria...)
	for i, task := range ticket.Tasks {
		copy.Tasks[i] = domain.Task{
			Title:              task.Title,
			Description:        task.Description,
			JiraID:             task.JiraID,
			SourceLine:         task.SourceLine,
			CustomFields:       make(map[string]string),
			AcceptanceCriteria: make([]string, len(task.AcceptanceCriteria)),
		}
		for k, v := range task.CustomFields {
			copy.Tasks[i].CustomFields[k] = v
		}
		copy.Tasks[i].AcceptanceCriteria = append(copy.Tasks[i].AcceptanceCriteria[:0], task.AcceptanceCriteria...)
	}

	return copy
}

// SetFocused updates border color when focus changes
func (v *TicketDetailView) SetFocused(focused bool) {
	color := theme.GetSecondaryColor()
	if focused {
		color = theme.GetPrimaryColor()
	}

	if v.editMode {
		v.editForm.SetBorderColor(color)
	} else {
		v.displayView.SetBorderColor(color)
	}
}
