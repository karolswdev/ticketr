package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/karolswdev/ticktr/internal/adapters/jira"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/spf13/cobra"
)

var (
	// Bulk command flags
	bulkIDs     string
	bulkSet     []string
	bulkParent  string
	bulkConfirm bool

	// bulkCmd represents the main bulk command
	bulkCmd = &cobra.Command{
		Use:   "bulk",
		Short: "Perform bulk operations on multiple tickets",
		Long: `Bulk operations allow you to update, move, or delete multiple tickets at once.

This is useful for batch updates across multiple tickets in your Jira project.
All operations provide real-time progress feedback and detailed error reporting.

Available operations:
  update - Update multiple tickets with field changes
  move   - Move multiple tickets to a new parent
  delete - Delete multiple tickets (requires confirmation)

Examples:
  ticketr bulk update --ids PROJ-1,PROJ-2 --set status=Done
  ticketr bulk move --ids PROJ-1,PROJ-2 --parent PROJ-100
  ticketr bulk delete --ids PROJ-1,PROJ-2 --confirm`,
	}

	// bulkUpdateCmd updates multiple tickets
	bulkUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update multiple tickets with field changes",
		Long: `Update multiple tickets with the same field changes.

Use --ids to specify which tickets to update (comma-separated list).
Use --set to specify field changes (can be used multiple times).

Field changes are applied to all specified tickets.
Progress is shown in real-time as each ticket is processed.

Examples:
  # Update status for multiple tickets
  ticketr bulk update --ids PROJ-1,PROJ-2,PROJ-3 --set status=Done

  # Update multiple fields
  ticketr bulk update --ids PROJ-1,PROJ-2 --set status="In Progress" --set assignee=john@example.com

  # Update with spaces in values (use quotes)
  ticketr bulk update --ids PROJ-1,PROJ-2 --set priority="High Priority"`,
		Args: cobra.NoArgs,
		RunE: runBulkUpdate,
	}

	// bulkMoveCmd moves multiple tickets to a new parent
	bulkMoveCmd = &cobra.Command{
		Use:   "move",
		Short: "Move multiple tickets to a new parent",
		Long: `Move multiple tickets under a new parent ticket.

This updates the parent relationship for all specified tickets,
making them children (sub-tasks/sub-issues) of the parent ticket.

Use --ids to specify which tickets to move (comma-separated list).
Use --parent to specify the new parent ticket ID.

Examples:
  # Move tickets under a new parent
  ticketr bulk move --ids PROJ-1,PROJ-2 --parent PROJ-100

  # Move sub-tasks to a different epic
  ticketr bulk move --ids TASK-1,TASK-2,TASK-3 --parent EPIC-42`,
		Args: cobra.NoArgs,
		RunE: runBulkMove,
	}

	// bulkDeleteCmd deletes multiple tickets
	bulkDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete multiple tickets (requires confirmation)",
		Long: `Delete multiple tickets from Jira.

WARNING: This operation cannot be undone. Deleted tickets are permanently removed.

Safety requirements:
  1. Must use --confirm flag
  2. Must confirm deletion at prompt

Use --ids to specify which tickets to delete (comma-separated list).
Use --confirm to acknowledge the destructive operation.

Examples:
  # Delete multiple tickets (requires confirmation)
  ticketr bulk delete --ids PROJ-1,PROJ-2 --confirm

NOTE: Bulk delete is currently not supported in Ticketr v3.0.
      This feature is planned for v3.1.0.`,
		Args: cobra.NoArgs,
		RunE: runBulkDelete,
	}
)

func init() {
	// Add subcommands to bulk
	bulkCmd.AddCommand(bulkUpdateCmd)
	bulkCmd.AddCommand(bulkMoveCmd)
	bulkCmd.AddCommand(bulkDeleteCmd)

	// Flags for update command
	bulkUpdateCmd.Flags().StringVar(&bulkIDs, "ids", "", "Comma-separated list of ticket IDs (e.g., PROJ-1,PROJ-2)")
	bulkUpdateCmd.Flags().StringArrayVar(&bulkSet, "set", []string{}, "Field changes to apply (e.g., --set status=Done --set assignee=user@example.com)")
	bulkUpdateCmd.MarkFlagRequired("ids")
	bulkUpdateCmd.MarkFlagRequired("set")

	// Flags for move command
	bulkMoveCmd.Flags().StringVar(&bulkIDs, "ids", "", "Comma-separated list of ticket IDs (e.g., PROJ-1,PROJ-2)")
	bulkMoveCmd.Flags().StringVar(&bulkParent, "parent", "", "New parent ticket ID (e.g., PROJ-100)")
	bulkMoveCmd.MarkFlagRequired("ids")
	bulkMoveCmd.MarkFlagRequired("parent")

	// Flags for delete command
	bulkDeleteCmd.Flags().StringVar(&bulkIDs, "ids", "", "Comma-separated list of ticket IDs (e.g., PROJ-1,PROJ-2)")
	bulkDeleteCmd.Flags().BoolVar(&bulkConfirm, "confirm", false, "Confirm deletion (required for safety)")
	bulkDeleteCmd.MarkFlagRequired("ids")
	bulkDeleteCmd.MarkFlagRequired("confirm")
}

// parseTicketIDs parses a comma-separated string of ticket IDs into a slice.
// It trims whitespace from each ID to handle "PROJ-1, PROJ-2" style input.
func parseTicketIDs(idsFlag string) []string {
	if idsFlag == "" {
		return []string{}
	}

	ids := strings.Split(idsFlag, ",")
	// Trim whitespace from each ID
	for i, id := range ids {
		ids[i] = strings.TrimSpace(id)
	}

	return ids
}

// parseSetFlags parses --set flags into a map of field changes.
// Expected format: field=value (e.g., "status=Done", "assignee=john@example.com")
func parseSetFlags(setFlags []string) (map[string]interface{}, error) {
	if len(setFlags) == 0 {
		return nil, fmt.Errorf("no field changes specified. Use --set field=value (e.g., --set status=Done)")
	}

	changes := make(map[string]interface{})
	for _, set := range setFlags {
		parts := strings.SplitN(set, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid --set format '%s': expected 'field=value' (e.g., status=Done). Use quotes for values with spaces", set)
		}
		field := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if field == "" {
			return nil, fmt.Errorf("invalid --set format '%s': field name cannot be empty. Example: --set status=Done", set)
		}

		changes[field] = value
	}

	return changes, nil
}

// createProgressCallback creates a progress callback function for bulk operations.
// It displays real-time progress as each ticket is processed.
func createProgressCallback(total int) ports.BulkOperationProgressCallback {
	current := 0
	return func(ticketID string, success bool, err error) {
		current++
		if success {
			fmt.Printf("[%d/%d] %s: ✓ success\n", current, total, ticketID)
		} else {
			fmt.Printf("[%d/%d] %s: ✗ failed (%v)\n", current, total, ticketID, err)
		}
	}
}

// initBulkService initializes a BulkOperationService using workspace credentials.
// It follows the same pattern as workspace commands for authentication.
func initBulkService() (ports.BulkOperationService, error) {
	// 1. Get workspace service
	workspaceSvc, err := initWorkspaceService()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize workspace service: %w", err)
	}

	// 2. Get current workspace
	workspace, err := workspaceSvc.Current()
	if err != nil {
		return nil, fmt.Errorf("no workspace selected: use 'ticketr workspace create' or 'ticketr workspace switch': %w", err)
	}

	// 3. Get credentials
	config, err := workspaceSvc.GetConfig(workspace.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace config: %w", err)
	}

	// 4. Initialize Jira adapter
	jiraAdapter, err := jira.NewJiraAdapterFromConfig(config, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Jira adapter: %w", err)
	}

	// 5. Authenticate with Jira
	if err := jiraAdapter.Authenticate(); err != nil {
		return nil, fmt.Errorf("failed to authenticate with Jira: %w", err)
	}

	// 6. Create bulk operation service
	bulkSvc := services.NewBulkOperationService(jiraAdapter)

	return bulkSvc, nil
}

// runBulkUpdate handles the bulk update command.
func runBulkUpdate(cmd *cobra.Command, args []string) error {
	// Parse ticket IDs
	ticketIDs := parseTicketIDs(bulkIDs)
	if len(ticketIDs) == 0 {
		return fmt.Errorf("no ticket IDs specified. Use --ids with comma-separated ticket IDs (e.g., --ids PROJ-1,PROJ-2,PROJ-3)")
	}

	// Parse field changes
	changes, err := parseSetFlags(bulkSet)
	if err != nil {
		return err
	}

	// Initialize service
	svc, err := initBulkService()
	if err != nil {
		return err
	}

	// Create bulk operation
	op := domain.NewBulkOperation(domain.BulkActionUpdate, ticketIDs, changes)

	// Validate operation
	if err := op.Validate(); err != nil {
		return err
	}

	// Display operation summary
	fmt.Printf("Updating %d ticket(s) with changes:\n", len(ticketIDs))
	for field, value := range changes {
		fmt.Printf("  %s = %v\n", field, value)
	}
	fmt.Println()

	// Execute with progress callback
	ctx := context.Background()
	progressCallback := createProgressCallback(len(ticketIDs))
	result, err := svc.ExecuteOperation(ctx, op, progressCallback)

	// Print summary
	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("✓ %d ticket(s) updated successfully\n", result.SuccessCount)

	if result.FailureCount > 0 {
		fmt.Printf("✗ %d ticket(s) failed\n", result.FailureCount)
		fmt.Println("\nErrors:")
		for ticketID, errMsg := range result.Errors {
			fmt.Printf("  %s: %s\n", ticketID, errMsg)
		}
	}

	// Return error if operation failed
	if err != nil {
		return fmt.Errorf("bulk update completed with errors: %w", err)
	}

	return nil
}

// runBulkMove handles the bulk move command.
func runBulkMove(cmd *cobra.Command, args []string) error {
	// Parse ticket IDs
	ticketIDs := parseTicketIDs(bulkIDs)
	if len(ticketIDs) == 0 {
		return fmt.Errorf("no ticket IDs specified. Use --ids with comma-separated ticket IDs (e.g., --ids PROJ-1,PROJ-2)")
	}

	// Validate parent ID
	if bulkParent == "" {
		return fmt.Errorf("no parent ticket specified. Use --parent with target parent ticket ID (e.g., --parent PROJ-100)")
	}

	// Initialize service
	svc, err := initBulkService()
	if err != nil {
		return err
	}

	// Create changes map with parent field
	changes := map[string]interface{}{
		"parent": bulkParent,
	}

	// Create bulk operation
	op := domain.NewBulkOperation(domain.BulkActionMove, ticketIDs, changes)

	// Validate operation
	if err := op.Validate(); err != nil {
		return err
	}

	// Display operation summary
	fmt.Printf("Moving %d ticket(s) to parent %s\n\n", len(ticketIDs), bulkParent)

	// Execute with progress callback
	ctx := context.Background()
	progressCallback := createProgressCallback(len(ticketIDs))
	result, err := svc.ExecuteOperation(ctx, op, progressCallback)

	// Print summary
	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("✓ %d ticket(s) moved successfully\n", result.SuccessCount)

	if result.FailureCount > 0 {
		fmt.Printf("✗ %d ticket(s) failed\n", result.FailureCount)
		fmt.Println("\nErrors:")
		for ticketID, errMsg := range result.Errors {
			fmt.Printf("  %s: %s\n", ticketID, errMsg)
		}
	}

	// Return error if operation failed
	if err != nil {
		return fmt.Errorf("bulk move completed with errors: %w", err)
	}

	return nil
}

// runBulkDelete handles the bulk delete command.
func runBulkDelete(cmd *cobra.Command, args []string) error {
	// Parse ticket IDs
	ticketIDs := parseTicketIDs(bulkIDs)
	if len(ticketIDs) == 0 {
		return fmt.Errorf("--ids is required and cannot be empty")
	}

	// Check confirm flag (first safety check)
	if !bulkConfirm {
		return fmt.Errorf("delete operations require --confirm flag for safety. Add --confirm to acknowledge this is a destructive operation")
	}

	// Display warning banner
	fmt.Println("⚠️  Delete operations are not yet supported in Ticketr v3.0")
	fmt.Println()
	fmt.Println("Ticketr currently does not support bulk delete operations via the Jira API.")
	fmt.Println("This feature is planned for v3.1.0.")
	fmt.Println()
	fmt.Println("For now, please delete tickets individually through Jira:")

	if len(ticketIDs) <= 10 {
		// Show all tickets if list is reasonable
		fmt.Println()
		for _, ticketID := range ticketIDs {
			fmt.Printf("  • %s\n", ticketID)
		}
	} else {
		// Show first few if list is long
		fmt.Println()
		for i := 0; i < 5; i++ {
			fmt.Printf("  • %s\n", ticketIDs[i])
		}
		fmt.Printf("  ... and %d more\n", len(ticketIDs)-5)
	}

	fmt.Println()
	fmt.Println("Alternative options:")
	fmt.Println("  1. Use the Jira web interface for bulk deletion")
	fmt.Println("  2. Wait for v3.1.0 which will include this feature")
	fmt.Println()

	return fmt.Errorf("bulk delete not supported")
}

// promptConfirmation asks the user to confirm a destructive operation.
// Returns true if the user confirms (types 'y' or 'yes'), false otherwise.
func promptConfirmation(message string) bool {
	fmt.Printf("%s [y/N]: ", message)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
