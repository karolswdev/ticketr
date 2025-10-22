package views

import (
	"context"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/karolswdev/ticktr/internal/adapters/tui/theme"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
	"github.com/rivo/tview"
)

// BulkOperationType represents the type of bulk operation.
type BulkOperationType string

const (
	BulkOpUpdate BulkOperationType = "update"
	BulkOpMove   BulkOperationType = "move"
	BulkOpDelete BulkOperationType = "delete"
)

// BulkOperationsModal displays a modal dialog for bulk operations on tickets.
type BulkOperationsModal struct {
	app         *tview.Application
	bulkService ports.BulkOperationService
	onClose     func()
	onSuccess   func()

	// Current operation
	opType       BulkOperationType
	ticketIDs    []string
	ctx          context.Context
	cancelFunc   context.CancelFunc
	isProcessing bool

	// UI components
	pages         *tview.Pages
	menuModal     *tview.Modal
	updateForm    *tview.Form
	moveForm      *tview.Form
	deleteModal   *tview.Modal
	progressModal *tview.Modal
	resultModal   *tview.Modal
	progressText  *tview.TextView

	// Form fields for update
	statusField   *tview.InputField
	priorityField *tview.InputField
	assigneeField *tview.InputField
	customFields  *tview.TextArea

	// Form fields for move
	parentField *tview.InputField

	// Progress tracking
	currentProgress int
	totalCount      int
	successCount    int
	failureCount    int
	progressDetails []string
}

// NewBulkOperationsModal creates a new bulk operations modal.
func NewBulkOperationsModal(
	app *tview.Application,
	bulkService ports.BulkOperationService,
) *BulkOperationsModal {
	modal := &BulkOperationsModal{
		app:         app,
		bulkService: bulkService,
	}

	modal.setupPages()
	return modal
}

// setupPages creates all modal pages.
func (m *BulkOperationsModal) setupPages() {
	m.pages = tview.NewPages()

	// Create menu modal
	m.menuModal = tview.NewModal()
	m.menuModal.SetBorder(true).SetTitle(" Bulk Operations ")
	m.menuModal.SetBorderColor(theme.GetPrimaryColor())
	m.menuModal.AddButtons([]string{"Update Fields", "Move Tickets", "Delete Tickets", "Cancel"})

	// Create update form
	m.setupUpdateForm()

	// Create move form
	m.setupMoveForm()

	// Create delete confirmation modal
	m.setupDeleteModal()

	// Create progress modal
	m.setupProgressModal()

	// Create result modal
	m.setupResultModal()

	// Add pages
	m.pages.AddPage("menu", m.menuModal, true, true)
	m.pages.AddPage("update", m.updateForm, true, false)
	m.pages.AddPage("move", m.moveForm, true, false)
	m.pages.AddPage("delete", m.deleteModal, true, false)
	m.pages.AddPage("progress", m.progressModal, true, false)
	m.pages.AddPage("result", m.resultModal, true, false)
}

// setupUpdateForm creates the bulk update form.
func (m *BulkOperationsModal) setupUpdateForm() {
	m.updateForm = tview.NewForm()
	m.updateForm.SetBorder(true).SetTitle(" Bulk Update ")
	m.updateForm.SetBorderColor(theme.GetPrimaryColor())

	// Create form fields with better placeholders
	m.statusField = tview.NewInputField().
		SetLabel("Status").
		SetFieldWidth(40).
		SetPlaceholder("e.g., In Progress, Done (leave empty to skip)")

	m.priorityField = tview.NewInputField().
		SetLabel("Priority").
		SetFieldWidth(40).
		SetPlaceholder("e.g., High, Medium, Low (leave empty to skip)")

	m.assigneeField = tview.NewInputField().
		SetLabel("Assignee").
		SetFieldWidth(40).
		SetPlaceholder("e.g., user@company.com (leave empty to skip)")

	m.customFields = tview.NewTextArea().
		SetPlaceholder("Custom fields (key=value, one per line)")

	// Add fields to form
	m.updateForm.AddFormItem(m.statusField)
	m.updateForm.AddFormItem(m.priorityField)
	m.updateForm.AddFormItem(m.assigneeField)
	m.updateForm.AddFormItem(m.customFields)

	// Add help text
	helpText := tview.NewTextView().
		SetText("[gray]Fill in at least one field | Tab: Next field | Enter: Apply | ESC: Cancel[-]").
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	m.updateForm.AddFormItem(helpText)

	// Add buttons
	m.updateForm.AddButton("Apply", m.handleUpdateApply)
	m.updateForm.AddButton("Cancel", m.handleCancel)

	m.updateForm.SetButtonsAlign(tview.AlignCenter)
	m.updateForm.SetButtonBackgroundColor(theme.GetPrimaryColor())
	m.updateForm.SetButtonTextColor(tcell.ColorWhite)
}

// setupMoveForm creates the bulk move form.
func (m *BulkOperationsModal) setupMoveForm() {
	m.moveForm = tview.NewForm()
	m.moveForm.SetBorder(true).SetTitle(" Bulk Move ")
	m.moveForm.SetBorderColor(theme.GetPrimaryColor())

	// Create parent field
	m.parentField = tview.NewInputField().
		SetLabel("Parent Ticket ID *").
		SetFieldWidth(40).
		SetPlaceholder("e.g., PROJ-123")

	// Add fields to form
	m.moveForm.AddFormItem(m.parentField)

	// Add help text
	helpText := tview.NewTextView().
		SetText("[gray]* = Required field | Enter: Move | ESC: Cancel[-]").
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	m.moveForm.AddFormItem(helpText)

	// Add buttons
	m.moveForm.AddButton("Move", m.handleMoveApply)
	m.moveForm.AddButton("Cancel", m.handleCancel)

	m.moveForm.SetButtonsAlign(tview.AlignCenter)
	m.moveForm.SetButtonBackgroundColor(theme.GetPrimaryColor())
	m.moveForm.SetButtonTextColor(tcell.ColorWhite)
}

// setupDeleteModal creates the bulk delete confirmation modal.
func (m *BulkOperationsModal) setupDeleteModal() {
	m.deleteModal = tview.NewModal()
	m.deleteModal.SetBorder(true).SetTitle(" ⚠ DANGEROUS OPERATION ")
	m.deleteModal.SetBorderColor(theme.GetErrorColor())
	m.deleteModal.SetBackgroundColor(tcell.ColorDefault)

	m.deleteModal.AddButtons([]string{"Cancel"})
}

// setupProgressModal creates the progress display modal.
func (m *BulkOperationsModal) setupProgressModal() {
	m.progressText = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true)

	m.progressModal = tview.NewModal()
	m.progressModal.SetText("")
	m.progressModal.AddButtons([]string{"Cancel"})
	m.progressModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Cancel" {
			m.handleCancelOperation()
		}
	})

	m.progressModal.SetBorder(true).SetTitle(" Processing... ")
	m.progressModal.SetBorderColor(theme.GetInfoColor())
}

// setupResultModal creates the result display modal.
func (m *BulkOperationsModal) setupResultModal() {
	m.resultModal = tview.NewModal()
	m.resultModal.SetBorder(true).SetTitle(" Operation Complete ")
	m.resultModal.SetBorderColor(theme.GetSuccessColor())
	m.resultModal.AddButtons([]string{"OK"})
}

// Show displays the bulk operations menu.
func (m *BulkOperationsModal) Show(ticketIDs []string) {
	if len(ticketIDs) == 0 {
		// Show error - no tickets selected
		errorModal := tview.NewModal()
		errorModal.SetText("[red]No tickets selected.[-]\n\nPress Space to select tickets in the tree view.")
		errorModal.AddButtons([]string{"OK"})
		errorModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if m.onClose != nil {
				m.onClose()
			}
		})
		m.app.SetRoot(errorModal, true)
		return
	}

	m.ticketIDs = ticketIDs
	m.resetState()

	// Update menu text
	menuText := fmt.Sprintf("Perform bulk operation on %d ticket(s)?\n\nSelected tickets: %s",
		len(ticketIDs),
		strings.Join(ticketIDs, ", "))

	// Truncate if too long
	if len(menuText) > 300 {
		ticketList := strings.Join(ticketIDs[:3], ", ")
		menuText = fmt.Sprintf("Perform bulk operation on %d ticket(s)?\n\nSelected tickets: %s... and %d more",
			len(ticketIDs), ticketList, len(ticketIDs)-3)
	}

	m.menuModal.SetText(menuText)

	// Set up menu button handler
	m.menuModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		switch buttonLabel {
		case "Update Fields":
			m.showUpdateForm()
		case "Move Tickets":
			m.showMoveForm()
		case "Delete Tickets":
			m.showDeleteWarning()
		case "Cancel":
			if m.onClose != nil {
				m.onClose()
			}
		}
	})

	// Show menu
	m.pages.SwitchToPage("menu")
	m.app.SetRoot(m.pages, true)
	m.app.SetFocus(m.menuModal)
}

// showUpdateForm displays the update form with ticket count.
func (m *BulkOperationsModal) showUpdateForm() {
	m.opType = BulkOpUpdate

	// Update form title with ticket count
	m.updateForm.SetTitle(fmt.Sprintf(" Bulk Update (%d tickets) ", len(m.ticketIDs)))

	// Clear previous values
	m.statusField.SetText("")
	m.priorityField.SetText("")
	m.assigneeField.SetText("")
	m.customFields.SetText("", true)

	m.pages.SwitchToPage("update")
	m.app.SetFocus(m.statusField)
}

// showMoveForm displays the move form with ticket count.
func (m *BulkOperationsModal) showMoveForm() {
	m.opType = BulkOpMove

	// Update form title with ticket count
	m.moveForm.SetTitle(fmt.Sprintf(" Bulk Move (%d tickets) ", len(m.ticketIDs)))

	// Clear previous value
	m.parentField.SetText("")

	m.pages.SwitchToPage("move")
	m.app.SetFocus(m.parentField)
}

// showDeleteWarning displays the delete warning modal.
func (m *BulkOperationsModal) showDeleteWarning() {
	m.opType = BulkOpDelete

	// Build ticket list
	ticketList := strings.Join(m.ticketIDs, "\n  • ")
	if len(m.ticketIDs) > 10 {
		ticketList = strings.Join(m.ticketIDs[:10], "\n  • ")
		ticketList += fmt.Sprintf("\n  • ... and %d more", len(m.ticketIDs)-10)
	}

	warningText := fmt.Sprintf(`[red::b]⚠ WARNING: DANGEROUS OPERATION[-:-:-]

[yellow]Bulk delete is not yet supported.[-]

The Jira adapter does not currently implement the DeleteTicket() method.
This feature will be available in v3.1.0.

[white]Tickets you attempted to delete:[-]
  • %s

[cyan]Workaround:[-]
Delete these tickets manually in your Jira web interface.

[gray]This limitation exists to prevent accidental data loss until proper
delete functionality with confirmations is implemented.[-]
`, ticketList)

	m.deleteModal.SetText(warningText)
	m.deleteModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Cancel" {
			m.pages.SwitchToPage("menu")
		}
	})

	m.pages.SwitchToPage("delete")
	m.app.SetFocus(m.deleteModal)
}

// handleUpdateApply processes the bulk update request.
func (m *BulkOperationsModal) handleUpdateApply() {
	if m.isProcessing {
		return
	}

	// Validate that at least one field is set
	changes := make(map[string]interface{})

	status := strings.TrimSpace(m.statusField.GetText())
	if status != "" {
		changes["Status"] = status
	}

	priority := strings.TrimSpace(m.priorityField.GetText())
	if priority != "" {
		changes["Priority"] = priority
	}

	assignee := strings.TrimSpace(m.assigneeField.GetText())
	if assignee != "" {
		changes["Assignee"] = assignee
	}

	// Parse custom fields
	customFieldsText := m.customFields.GetText()
	if customFieldsText != "" {
		lines := strings.Split(customFieldsText, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				m.showError(fmt.Sprintf("Invalid custom field format: %s\nExpected: key=value", line))
				return
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			changes[key] = value
		}
	}

	if len(changes) == 0 {
		m.showError("No fields specified. Please fill in at least one field to update.")
		return
	}

	// Create bulk operation
	bulkOp := domain.NewBulkOperation(domain.BulkActionUpdate, m.ticketIDs, changes)

	// Execute operation
	m.executeOperation(bulkOp)
}

// handleMoveApply processes the bulk move request.
func (m *BulkOperationsModal) handleMoveApply() {
	if m.isProcessing {
		return
	}

	parentID := strings.TrimSpace(m.parentField.GetText())
	if parentID == "" {
		m.showError("Parent ticket ID is required.")
		return
	}

	// Validate parent ID format (basic check)
	if !strings.Contains(parentID, "-") {
		m.showError("Invalid parent ticket ID format. Expected: PROJECT-123")
		return
	}

	// Check that we're not moving a ticket to itself
	for _, ticketID := range m.ticketIDs {
		if ticketID == parentID {
			m.showError(fmt.Sprintf("Cannot move ticket %s to itself.", ticketID))
			return
		}
	}

	// Create bulk operation
	changes := map[string]interface{}{
		"parent": parentID,
	}
	bulkOp := domain.NewBulkOperation(domain.BulkActionMove, m.ticketIDs, changes)

	// Execute operation
	m.executeOperation(bulkOp)
}

// executeOperation executes the bulk operation with progress tracking.
func (m *BulkOperationsModal) executeOperation(bulkOp *domain.BulkOperation) {
	m.isProcessing = true
	m.currentProgress = 0
	m.totalCount = len(m.ticketIDs)
	m.successCount = 0
	m.failureCount = 0
	m.progressDetails = make([]string, 0)

	// Create cancellable context
	m.ctx, m.cancelFunc = context.WithCancel(context.Background())

	// Show progress modal
	m.showProgressModal()

	// Execute in background
	go func() {
		// Progress callback
		progress := func(ticketID string, success bool, err error) {
			m.app.QueueUpdateDraw(func() {
				m.currentProgress++
				if success {
					m.successCount++
					m.progressDetails = append(m.progressDetails, fmt.Sprintf("[green]✓[-] %s: Success", ticketID))
				} else {
					m.failureCount++
					errMsg := "Unknown error"
					if err != nil {
						errMsg = err.Error()
					}
					m.progressDetails = append(m.progressDetails, fmt.Sprintf("[red]✗[-] %s: %s", ticketID, errMsg))
				}
				m.updateProgressDisplay()
			})
		}

		// Execute operation
		result, err := m.bulkService.ExecuteOperation(m.ctx, bulkOp, progress)

		// Update UI on completion
		m.app.QueueUpdateDraw(func() {
			m.isProcessing = false
			m.showResult(result, err)
		})
	}()
}

// showProgressModal displays the progress modal.
func (m *BulkOperationsModal) showProgressModal() {
	m.progressModal.SetText(m.formatProgressText())
	m.progressModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Cancel" {
			m.handleCancelOperation()
		}
	})

	m.pages.SwitchToPage("progress")
	m.app.SetFocus(m.progressModal)
}

// updateProgressDisplay updates the progress modal with current status.
func (m *BulkOperationsModal) updateProgressDisplay() {
	m.progressModal.SetText(m.formatProgressText())
}

// formatProgressText formats the progress text for display.
func (m *BulkOperationsModal) formatProgressText() string {
	progressPercent := 0
	if m.totalCount > 0 {
		progressPercent = (m.currentProgress * 100) / m.totalCount
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Processing %d ticket(s)...\n\n", m.totalCount))
	sb.WriteString(fmt.Sprintf("Progress: [%d/%d] (%d%%)\n", m.currentProgress, m.totalCount, progressPercent))
	sb.WriteString(fmt.Sprintf("[green]Success: %d[-]  [red]Failed: %d[-]\n\n", m.successCount, m.failureCount))

	// Show recent progress (last 10 items)
	sb.WriteString("Recent updates:\n")
	startIdx := 0
	if len(m.progressDetails) > 10 {
		startIdx = len(m.progressDetails) - 10
	}
	for i := startIdx; i < len(m.progressDetails); i++ {
		sb.WriteString(fmt.Sprintf("  %s\n", m.progressDetails[i]))
	}

	return sb.String()
}

// showResult displays the operation result.
func (m *BulkOperationsModal) showResult(result *domain.BulkOperationResult, err error) {
	var resultText string
	var borderColor tcell.Color

	if err != nil {
		// Operation-level error
		borderColor = theme.GetErrorColor()
		resultText = fmt.Sprintf("[red]Operation failed:[-]\n\n%s\n\n", err.Error())

		if result != nil {
			resultText += fmt.Sprintf("Partial results:\n")
			resultText += fmt.Sprintf("[green]Success: %d[-]  [red]Failed: %d[-]", result.SuccessCount, result.FailureCount)

			if result.FailureCount > 0 {
				resultText += "\n\nFailed tickets:\n"
				for _, ticketID := range result.FailedTickets {
					errMsg := result.Errors[ticketID]
					resultText += fmt.Sprintf("  • %s: %s\n", ticketID, errMsg)
				}
			}
		}
	} else if result != nil {
		if result.FailureCount == 0 {
			// Complete success
			borderColor = theme.GetSuccessColor()
			resultText = fmt.Sprintf("[green]✓ Operation completed successfully![-]\n\n")
			resultText += fmt.Sprintf("Processed %d ticket(s)\n", result.SuccessCount)
		} else {
			// Partial success
			borderColor = theme.GetWarningColor()
			resultText = fmt.Sprintf("[yellow]⚠ Operation completed with errors[-]\n\n")
			resultText += fmt.Sprintf("[green]Success: %d[-]  [red]Failed: %d[-]\n\n", result.SuccessCount, result.FailureCount)

			if result.FailureCount > 0 {
				resultText += "Failed tickets:\n"
				maxDisplay := 10
				for i, ticketID := range result.FailedTickets {
					if i >= maxDisplay {
						resultText += fmt.Sprintf("  ... and %d more\n", len(result.FailedTickets)-maxDisplay)
						break
					}
					errMsg := result.Errors[ticketID]
					resultText += fmt.Sprintf("  • %s: %s\n", ticketID, errMsg)
				}
			}
		}
	} else {
		// Unknown state
		borderColor = theme.GetErrorColor()
		resultText = "[red]Unknown error: No result returned[-]"
	}

	m.resultModal.SetText(resultText)
	m.resultModal.SetBorderColor(borderColor)
	m.resultModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "OK" {
			// Call success callback if operation was successful
			if err == nil && result != nil && result.FailureCount == 0 {
				if m.onSuccess != nil {
					m.onSuccess()
				}
			}
			if m.onClose != nil {
				m.onClose()
			}
		}
	})

	m.pages.SwitchToPage("result")
	m.app.SetFocus(m.resultModal)
}

// handleCancelOperation cancels the ongoing operation.
func (m *BulkOperationsModal) handleCancelOperation() {
	if m.cancelFunc != nil {
		m.cancelFunc()
	}

	// Show cancellation message
	m.resultModal.SetText("[yellow]Operation cancelled by user.[-]\n\nPartial changes may have been applied.")
	m.resultModal.SetBorderColor(theme.GetWarningColor())
	m.resultModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if m.onClose != nil {
			m.onClose()
		}
	})

	m.pages.SwitchToPage("result")
	m.app.SetFocus(m.resultModal)
}

// handleCancel closes the modal without executing.
func (m *BulkOperationsModal) handleCancel() {
	if m.onClose != nil {
		m.onClose()
	}
}

// showError displays an error message.
func (m *BulkOperationsModal) showError(message string) {
	errorModal := tview.NewModal()
	errorModal.SetText(fmt.Sprintf("[red::b]Error:[-:-:-]\n\n%s\n\n[yellow]Press OK or ESC to continue...[-]", message))
	errorModal.SetBorder(true).SetTitle(" ⚠ Error ")
	errorModal.SetBorderColor(theme.GetErrorColor())
	errorModal.SetBackgroundColor(tcell.ColorDefault)
	errorModal.AddButtons([]string{"OK"})

	// Set button styling
	errorModal.SetButtonBackgroundColor(theme.GetErrorColor())
	errorModal.SetButtonTextColor(tcell.ColorWhite)

	errorModal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		// Return to current form
		switch m.opType {
		case BulkOpUpdate:
			m.pages.SwitchToPage("update")
			m.app.SetFocus(m.statusField)
		case BulkOpMove:
			m.pages.SwitchToPage("move")
			m.app.SetFocus(m.parentField)
		default:
			m.pages.SwitchToPage("menu")
			m.app.SetFocus(m.menuModal)
		}
	})

	m.pages.AddPage("error", errorModal, true, true)
	m.pages.SwitchToPage("error")
	m.app.SetFocus(errorModal)
}

// resetState resets the modal state for a new operation.
func (m *BulkOperationsModal) resetState() {
	m.opType = ""
	m.isProcessing = false
	m.currentProgress = 0
	m.totalCount = 0
	m.successCount = 0
	m.failureCount = 0
	m.progressDetails = make([]string, 0)

	if m.cancelFunc != nil {
		m.cancelFunc()
		m.cancelFunc = nil
	}
	m.ctx = nil
}

// SetOnClose sets the callback for when the modal is closed.
func (m *BulkOperationsModal) SetOnClose(callback func()) {
	m.onClose = callback
}

// SetOnSuccess sets the callback for when an operation completes successfully.
func (m *BulkOperationsModal) SetOnSuccess(callback func()) {
	m.onSuccess = callback
}

// Primitive returns the underlying tview primitive.
func (m *BulkOperationsModal) Primitive() tview.Primitive {
	return m.pages
}
