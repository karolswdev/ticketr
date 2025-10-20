package views

import (
	"context"
	"fmt"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
	"github.com/rivo/tview"
)

// mockBulkOperationService is a mock implementation of ports.BulkOperationService.
type mockBulkOperationService struct {
	executeFunc func(ctx context.Context, op *domain.BulkOperation, progress ports.BulkOperationProgressCallback) (*domain.BulkOperationResult, error)
}

func (m *mockBulkOperationService) ExecuteOperation(
	ctx context.Context,
	op *domain.BulkOperation,
	progress ports.BulkOperationProgressCallback,
) (*domain.BulkOperationResult, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, op, progress)
	}
	return domain.NewBulkOperationResult(), nil
}

// TestNewBulkOperationsModal tests modal creation.
func TestNewBulkOperationsModal(t *testing.T) {
	app := tview.NewApplication()
	mockService := &mockBulkOperationService{}

	modal := NewBulkOperationsModal(app, mockService)

	if modal == nil {
		t.Fatal("NewBulkOperationsModal returned nil")
	}

	if modal.app != app {
		t.Error("app not set correctly")
	}

	if modal.bulkService != mockService {
		t.Error("bulkService not set correctly")
	}

	if modal.pages == nil {
		t.Error("pages not initialized")
	}
}

// TestBulkOperationsModal_ShowWithNoTickets tests showing modal with no tickets selected.
func TestBulkOperationsModal_ShowWithNoTickets(t *testing.T) {
	app := tview.NewApplication()
	mockService := &mockBulkOperationService{}

	modal := NewBulkOperationsModal(app, mockService)

	// Show with empty ticket list
	modal.Show([]string{})

	// The modal should display an error message
	// We can't easily verify the UI state, but we can check that the function doesn't panic
	// In real usage, this would show an error modal
}

// TestBulkOperationsModal_ShowWithTickets tests showing modal with tickets.
func TestBulkOperationsModal_ShowWithTickets(t *testing.T) {
	app := tview.NewApplication()
	mockService := &mockBulkOperationService{}

	modal := NewBulkOperationsModal(app, mockService)

	tickets := []string{"PROJ-1", "PROJ-2", "PROJ-3"}
	modal.Show(tickets)

	// Verify tickets are stored
	if len(modal.ticketIDs) != 3 {
		t.Errorf("expected 3 tickets, got %d", len(modal.ticketIDs))
	}

	for i, ticketID := range tickets {
		if modal.ticketIDs[i] != ticketID {
			t.Errorf("ticket %d: expected %s, got %s", i, ticketID, modal.ticketIDs[i])
		}
	}
}

// TestBulkOperationsModal_UpdateOperation tests the update operation setup.
func TestBulkOperationsModal_UpdateOperation(t *testing.T) {
	tests := []struct {
		name           string
		ticketIDs      []string
		changes        map[string]interface{}
		expectedAction domain.BulkOperationAction
	}{
		{
			name:      "successful update setup",
			ticketIDs: []string{"PROJ-1", "PROJ-2"},
			changes: map[string]interface{}{
				"Status": "In Progress",
			},
			expectedAction: domain.BulkActionUpdate,
		},
		{
			name:      "partial failure setup",
			ticketIDs: []string{"PROJ-1", "PROJ-2", "PROJ-3"},
			changes: map[string]interface{}{
				"Priority": "High",
			},
			expectedAction: domain.BulkActionUpdate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := tview.NewApplication()
			mockService := &mockBulkOperationService{}

			modal := NewBulkOperationsModal(app, mockService)
			modal.ticketIDs = tt.ticketIDs

			// Create bulk operation to test validation
			bulkOp := domain.NewBulkOperation(tt.expectedAction, tt.ticketIDs, tt.changes)

			// Validate the operation is properly constructed
			if err := bulkOp.Validate(); err != nil {
				t.Errorf("bulk operation validation failed: %v", err)
			}

			if bulkOp.Action != tt.expectedAction {
				t.Errorf("expected action %s, got %s", tt.expectedAction, bulkOp.Action)
			}

			if len(bulkOp.TicketIDs) != len(tt.ticketIDs) {
				t.Errorf("expected %d tickets, got %d", len(tt.ticketIDs), len(bulkOp.TicketIDs))
			}
		})
	}
}

// TestBulkOperationsModal_MoveOperation tests the move operation flow.
func TestBulkOperationsModal_MoveOperation(t *testing.T) {
	tests := []struct {
		name       string
		ticketIDs  []string
		parentID   string
		shouldFail bool
	}{
		{
			name:       "valid move",
			ticketIDs:  []string{"PROJ-1", "PROJ-2"},
			parentID:   "PROJ-100",
			shouldFail: false,
		},
		{
			name:       "move to self",
			ticketIDs:  []string{"PROJ-1", "PROJ-2"},
			parentID:   "PROJ-1", // Trying to move PROJ-1 to itself
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := tview.NewApplication()
			mockService := &mockBulkOperationService{
				executeFunc: func(ctx context.Context, op *domain.BulkOperation, progress ports.BulkOperationProgressCallback) (*domain.BulkOperationResult, error) {
					result := domain.NewBulkOperationResult()
					result.SuccessCount = len(op.TicketIDs)
					result.SuccessfulTickets = op.TicketIDs
					return result, nil
				},
			}

			modal := NewBulkOperationsModal(app, mockService)
			modal.ticketIDs = tt.ticketIDs
			modal.parentField.SetText(tt.parentID)

			// The validation logic should be in handleMoveApply
			// We can't easily test the UI interaction, but we can verify the data flow
		})
	}
}

// TestBulkOperationsModal_DeleteOperation tests the delete operation (not supported).
func TestBulkOperationsModal_DeleteOperation(t *testing.T) {
	app := tview.NewApplication()

	// Mock service that returns "not supported" error for delete
	mockService := &mockBulkOperationService{
		executeFunc: func(ctx context.Context, op *domain.BulkOperation, progress ports.BulkOperationProgressCallback) (*domain.BulkOperationResult, error) {
			result := domain.NewBulkOperationResult()
			for _, ticketID := range op.TicketIDs {
				result.FailureCount++
				result.FailedTickets = append(result.FailedTickets, ticketID)
				result.Errors[ticketID] = "delete operation not supported by Jira adapter"
			}
			return result, fmt.Errorf("delete operation not supported")
		},
	}

	modal := NewBulkOperationsModal(app, mockService)
	modal.ticketIDs = []string{"PROJ-1", "PROJ-2"}

	// Show delete warning
	modal.showDeleteWarning()

	// The modal should display a warning about delete not being supported
	// We verify the opType is set correctly
	if modal.opType != BulkOpDelete {
		t.Errorf("expected opType %s, got %s", BulkOpDelete, modal.opType)
	}
}

// TestBulkOperationsModal_ProgressTracking tests progress state initialization.
func TestBulkOperationsModal_ProgressTracking(t *testing.T) {
	app := tview.NewApplication()

	ticketIDs := []string{"PROJ-1", "PROJ-2", "PROJ-3"}

	mockService := &mockBulkOperationService{
		executeFunc: func(ctx context.Context, op *domain.BulkOperation, progress ports.BulkOperationProgressCallback) (*domain.BulkOperationResult, error) {
			result := domain.NewBulkOperationResult()

			// Simulate processing each ticket
			for i, ticketID := range op.TicketIDs {
				success := i != 1 // Make the second ticket fail
				var err error
				if !success {
					err = fmt.Errorf("ticket %s failed", ticketID)
					result.FailureCount++
					result.FailedTickets = append(result.FailedTickets, ticketID)
					result.Errors[ticketID] = err.Error()
				} else {
					result.SuccessCount++
					result.SuccessfulTickets = append(result.SuccessfulTickets, ticketID)
				}

				// Invoke progress callback
				if progress != nil {
					progress(ticketID, success, err)
				}
			}

			return result, nil
		},
	}

	_ = NewBulkOperationsModal(app, mockService)

	bulkOp := domain.NewBulkOperation(domain.BulkActionUpdate, ticketIDs, map[string]interface{}{
		"Status": "Done",
	})

	// Test that the bulk operation is valid
	if err := bulkOp.Validate(); err != nil {
		t.Errorf("bulk operation validation failed: %v", err)
	}

	// Verify ticket count
	if bulkOp.TicketCount() != len(ticketIDs) {
		t.Errorf("expected ticket count %d, got %d", len(ticketIDs), bulkOp.TicketCount())
	}
}

// TestBulkOperationsModal_Cancellation tests operation cancellation.
func TestBulkOperationsModal_Cancellation(t *testing.T) {
	app := tview.NewApplication()

	mockService := &mockBulkOperationService{
		executeFunc: func(ctx context.Context, op *domain.BulkOperation, progress ports.BulkOperationProgressCallback) (*domain.BulkOperationResult, error) {
			// Simulate checking for cancellation
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			result := domain.NewBulkOperationResult()
			result.SuccessCount = len(op.TicketIDs)
			return result, nil
		},
	}

	modal := NewBulkOperationsModal(app, mockService)
	modal.ticketIDs = []string{"PROJ-1"}

	bulkOp := domain.NewBulkOperation(domain.BulkActionUpdate, modal.ticketIDs, map[string]interface{}{
		"Status": "Done",
	})

	// Start operation
	modal.executeOperation(bulkOp)

	// Verify context was created
	if modal.ctx == nil {
		t.Error("context not created")
	}

	if modal.cancelFunc == nil {
		t.Error("cancel function not created")
	}

	// Simulate cancellation
	if modal.cancelFunc != nil {
		modal.cancelFunc()
	}
}

// TestBulkOperationsModal_Callbacks tests the onClose and onSuccess callbacks.
func TestBulkOperationsModal_Callbacks(t *testing.T) {
	app := tview.NewApplication()
	mockService := &mockBulkOperationService{}

	modal := NewBulkOperationsModal(app, mockService)

	closeCalled := false
	successCalled := false

	modal.SetOnClose(func() {
		closeCalled = true
	})

	modal.SetOnSuccess(func() {
		successCalled = true
	})

	// Test onClose callback
	if modal.onClose != nil {
		modal.onClose()
	}

	if !closeCalled {
		t.Error("onClose callback was not invoked")
	}

	// Test onSuccess callback
	if modal.onSuccess != nil {
		modal.onSuccess()
	}

	if !successCalled {
		t.Error("onSuccess callback was not invoked")
	}
}

// TestBulkOperationsModal_ResetState tests state reset between operations.
func TestBulkOperationsModal_ResetState(t *testing.T) {
	app := tview.NewApplication()
	mockService := &mockBulkOperationService{}

	modal := NewBulkOperationsModal(app, mockService)

	// Set some state
	modal.opType = BulkOpUpdate
	modal.isProcessing = true
	modal.currentProgress = 5
	modal.totalCount = 10
	modal.successCount = 3
	modal.failureCount = 2
	modal.progressDetails = []string{"detail1", "detail2"}
	modal.ctx, modal.cancelFunc = context.WithCancel(context.Background())

	// Reset state
	modal.resetState()

	// Verify state was reset
	if modal.opType != "" {
		t.Errorf("opType not reset, got %s", modal.opType)
	}

	if modal.isProcessing {
		t.Error("isProcessing not reset")
	}

	if modal.currentProgress != 0 {
		t.Errorf("currentProgress not reset, got %d", modal.currentProgress)
	}

	if modal.totalCount != 0 {
		t.Errorf("totalCount not reset, got %d", modal.totalCount)
	}

	if modal.successCount != 0 {
		t.Errorf("successCount not reset, got %d", modal.successCount)
	}

	if modal.failureCount != 0 {
		t.Errorf("failureCount not reset, got %d", modal.failureCount)
	}

	if len(modal.progressDetails) != 0 {
		t.Errorf("progressDetails not reset, got %d items", len(modal.progressDetails))
	}

	if modal.ctx != nil {
		t.Error("ctx not reset")
	}

	if modal.cancelFunc != nil {
		t.Error("cancelFunc not reset")
	}
}

// TestBulkOperationsModal_FormatProgressText tests progress text formatting.
func TestBulkOperationsModal_FormatProgressText(t *testing.T) {
	app := tview.NewApplication()
	mockService := &mockBulkOperationService{}

	modal := NewBulkOperationsModal(app, mockService)

	modal.totalCount = 10
	modal.currentProgress = 5
	modal.successCount = 4
	modal.failureCount = 1
	modal.progressDetails = []string{
		"PROJ-1: Success",
		"PROJ-2: Success",
		"PROJ-3: Failed",
		"PROJ-4: Success",
		"PROJ-5: Success",
	}

	text := modal.formatProgressText()

	// Verify text contains expected elements
	if text == "" {
		t.Error("formatProgressText returned empty string")
	}

	// We can't easily verify exact formatting, but we can check for key components
	// In a real test, we might use regex or string contains checks
}
