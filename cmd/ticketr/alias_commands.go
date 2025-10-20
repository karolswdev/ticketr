package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/karolswdev/ticktr/internal/adapters/database"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/spf13/cobra"
)

var (
	// Alias command flags
	aliasDescription string
	aliasGlobal      bool

	// aliasCmd represents the main alias command
	aliasCmd = &cobra.Command{
		Use:   "alias",
		Short: "Manage JQL query aliases",
		Long: `Manage named JQL query aliases for easier ticket filtering.

Aliases allow you to define reusable JQL queries with memorable names.
They can be workspace-specific or global across all workspaces.

Predefined aliases (cannot be modified):
  - mine:    Tickets assigned to you
  - sprint:  Tickets in active sprints
  - blocked: Blocked tickets or tickets with blocked label

Examples:
  ticketr alias list
  ticketr alias create my-bugs "assignee = currentUser() AND type = Bug"
  ticketr alias create high-priority "priority = High" --description "High priority tickets"
  ticketr pull --alias mine --output my-tickets.md`,
	}

	// aliasListCmd lists all aliases
	aliasListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all JQL aliases",
		Long: `List all available JQL aliases for the current workspace.

Shows both predefined system aliases and user-defined aliases.
User-defined aliases override predefined ones with the same name.`,
		RunE: runAliasList,
	}

	// aliasCreateCmd creates a new alias
	aliasCreateCmd = &cobra.Command{
		Use:   "create <name> <jql>",
		Short: "Create a new JQL alias",
		Long: `Create a new JQL query alias.

The alias name must be alphanumeric with optional hyphens and underscores.
By default, aliases are workspace-specific. Use --global for global aliases.

Examples:
  # Workspace-specific alias
  ticketr alias create my-bugs "assignee = currentUser() AND type = Bug"

  # With description
  ticketr alias create urgent "priority = Highest AND status != Done" --description "Urgent open tickets"

  # Global alias (available in all workspaces)
  ticketr alias create critical "priority = Critical" --global`,
		Args: cobra.ExactArgs(2),
		RunE: runAliasCreate,
	}

	// aliasShowCmd shows a specific alias
	aliasShowCmd = &cobra.Command{
		Use:   "show <name>",
		Short: "Show details of a specific alias",
		Long: `Display the JQL query and details of a specific alias.

Shows the expanded JQL query and metadata for the alias.`,
		Args: cobra.ExactArgs(1),
		RunE: runAliasShow,
	}

	// aliasDeleteCmd deletes an alias
	aliasDeleteCmd = &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a JQL alias",
		Long: `Delete a user-defined JQL alias.

Predefined system aliases (mine, sprint, blocked) cannot be deleted.
Only workspace-specific aliases can be deleted from the current workspace.`,
		Args: cobra.ExactArgs(1),
		RunE: runAliasDelete,
	}

	// aliasUpdateCmd updates an existing alias
	aliasUpdateCmd = &cobra.Command{
		Use:   "update <name> <jql>",
		Short: "Update an existing JQL alias",
		Long: `Update the JQL query of an existing alias.

Predefined system aliases cannot be modified.
You can only update aliases in the current workspace.

Examples:
  ticketr alias update my-bugs "assignee = currentUser() AND type = Bug AND status != Done"
  ticketr alias update urgent "priority = Highest" --description "Updated description"`,
		Args: cobra.ExactArgs(2),
		RunE: runAliasUpdate,
	}
)

func init() {
	// Add subcommands to alias
	aliasCmd.AddCommand(aliasListCmd)
	aliasCmd.AddCommand(aliasCreateCmd)
	aliasCmd.AddCommand(aliasShowCmd)
	aliasCmd.AddCommand(aliasDeleteCmd)
	aliasCmd.AddCommand(aliasUpdateCmd)

	// Flags for create command
	aliasCreateCmd.Flags().StringVar(&aliasDescription, "description", "", "Description of the alias")
	aliasCreateCmd.Flags().BoolVar(&aliasGlobal, "global", false, "Create a global alias (available in all workspaces)")

	// Flags for update command
	aliasUpdateCmd.Flags().StringVar(&aliasDescription, "description", "", "Updated description of the alias")
}

// initAliasService initializes an AliasService instance.
func initAliasService() (*services.AliasService, error) {
	// Get PathResolver singleton
	pathResolver, err := services.GetPathResolver()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize path resolver: %w", err)
	}

	// Initialize database
	adapter, err := database.NewSQLiteAdapter(pathResolver)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create alias repository
	aliasRepo := database.NewAliasRepository(adapter.DB())

	// Create alias service
	svc := services.NewAliasService(aliasRepo)

	return svc, nil
}

// getCurrentWorkspaceID retrieves the current workspace ID.
// Returns empty string for global aliases.
func getCurrentWorkspaceID() (string, error) {
	if aliasGlobal {
		return "", nil
	}

	// Initialize workspace service
	workspaceSvc, err := initWorkspaceService()
	if err != nil {
		// If workspaces are not enabled, use empty string (global)
		return "", nil
	}

	// Get current workspace
	workspace, err := workspaceSvc.Current()
	if err != nil {
		// No workspace configured, use global
		return "", nil
	}

	return workspace.ID, nil
}

// runAliasList handles the alias list command.
func runAliasList(cmd *cobra.Command, args []string) error {
	// Initialize service
	svc, err := initAliasService()
	if err != nil {
		return err
	}

	// Get current workspace ID
	workspaceID, err := getCurrentWorkspaceID()
	if err != nil {
		return fmt.Errorf("failed to get current workspace: %w", err)
	}

	// List aliases
	aliases, err := svc.List(workspaceID)
	if err != nil {
		return fmt.Errorf("failed to list aliases: %w", err)
	}

	if len(aliases) == 0 {
		fmt.Println("No aliases found.")
		fmt.Println("\nCreate an alias with:")
		fmt.Println("  ticketr alias create <name> <jql>")
		return nil
	}

	// Print aliases in table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tTYPE\tDESCRIPTION\tJQL")
	fmt.Fprintln(w, "----\t----\t-----------\t---")

	for _, alias := range aliases {
		aliasType := "user"
		if alias.IsPredefined {
			aliasType = "system"
		} else if alias.WorkspaceID == "" {
			aliasType = "global"
		}

		// Truncate JQL if too long
		jql := alias.JQL
		if len(jql) > 50 {
			jql = jql[:47] + "..."
		}

		// Truncate description if too long
		description := alias.Description
		if len(description) > 30 {
			description = description[:27] + "..."
		}
		if description == "" {
			description = "-"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			alias.Name,
			aliasType,
			description,
			jql,
		)
	}

	w.Flush()

	fmt.Println("\nUse 'ticketr alias show <name>' to see the full JQL query")
	fmt.Println("Use 'ticketr pull --alias <name>' to pull tickets using an alias")

	return nil
}

// runAliasCreate handles the alias create command.
func runAliasCreate(cmd *cobra.Command, args []string) error {
	name := args[0]
	jql := args[1]

	// Initialize service
	svc, err := initAliasService()
	if err != nil {
		return err
	}

	// Get workspace ID
	workspaceID, err := getCurrentWorkspaceID()
	if err != nil {
		return fmt.Errorf("failed to get current workspace: %w", err)
	}

	// Create alias
	if err := svc.Create(name, jql, aliasDescription, workspaceID); err != nil {
		return fmt.Errorf("failed to create alias: %w", err)
	}

	scope := "workspace"
	if workspaceID == "" {
		scope = "global"
	}

	fmt.Printf("\n✓ Alias '%s' created successfully (%s)\n", name, scope)
	fmt.Printf("  JQL: %s\n", jql)
	if aliasDescription != "" {
		fmt.Printf("  Description: %s\n", aliasDescription)
	}
	fmt.Printf("\nUse 'ticketr pull --alias %s' to pull tickets\n", name)

	return nil
}

// runAliasShow handles the alias show command.
func runAliasShow(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Initialize service
	svc, err := initAliasService()
	if err != nil {
		return err
	}

	// Get workspace ID
	workspaceID, err := getCurrentWorkspaceID()
	if err != nil {
		return fmt.Errorf("failed to get current workspace: %w", err)
	}

	// Get alias
	alias, err := svc.Get(name, workspaceID)
	if err != nil {
		return fmt.Errorf("alias '%s' not found: %w", name, err)
	}

	// Display alias details
	fmt.Printf("Alias: %s\n", alias.Name)
	fmt.Printf("Type: ")
	if alias.IsPredefined {
		fmt.Println("system (predefined)")
	} else if alias.WorkspaceID == "" {
		fmt.Println("global")
	} else {
		fmt.Println("workspace")
	}

	if alias.Description != "" {
		fmt.Printf("Description: %s\n", alias.Description)
	}

	fmt.Printf("\nJQL Query:\n  %s\n", alias.JQL)

	// Try to expand the alias to show the full query
	expandedJQL, err := svc.ExpandAlias(name, workspaceID)
	if err == nil && expandedJQL != alias.JQL {
		fmt.Printf("\nExpanded JQL:\n  %s\n", expandedJQL)
	}

	if !alias.IsPredefined {
		fmt.Printf("\nCreated: %s\n", alias.CreatedAt.Format("2006-01-02 15:04:05"))
		if !alias.UpdatedAt.Equal(alias.CreatedAt) {
			fmt.Printf("Updated: %s\n", alias.UpdatedAt.Format("2006-01-02 15:04:05"))
		}
	}

	return nil
}

// runAliasDelete handles the alias delete command.
func runAliasDelete(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Check if this is a predefined alias
	if domain.IsPredefinedAlias(name) {
		return fmt.Errorf("cannot delete predefined alias '%s'", name)
	}

	// Initialize service
	svc, err := initAliasService()
	if err != nil {
		return err
	}

	// Get workspace ID
	workspaceID, err := getCurrentWorkspaceID()
	if err != nil {
		return fmt.Errorf("failed to get current workspace: %w", err)
	}

	// Delete alias
	if err := svc.Delete(name, workspaceID); err != nil {
		return fmt.Errorf("failed to delete alias: %w", err)
	}

	fmt.Printf("Alias '%s' deleted successfully\n", name)

	return nil
}

// runAliasUpdate handles the alias update command.
func runAliasUpdate(cmd *cobra.Command, args []string) error {
	name := args[0]
	jql := args[1]

	// Check if this is a predefined alias
	if domain.IsPredefinedAlias(name) {
		return fmt.Errorf("cannot update predefined alias '%s'", name)
	}

	// Initialize service
	svc, err := initAliasService()
	if err != nil {
		return err
	}

	// Get workspace ID
	workspaceID, err := getCurrentWorkspaceID()
	if err != nil {
		return fmt.Errorf("failed to get current workspace: %w", err)
	}

	// Update alias
	if err := svc.Update(name, jql, aliasDescription, workspaceID); err != nil {
		return fmt.Errorf("failed to update alias: %w", err)
	}

	fmt.Printf("\n✓ Alias '%s' updated successfully\n", name)
	fmt.Printf("  JQL: %s\n", jql)
	if aliasDescription != "" {
		fmt.Printf("  Description: %s\n", aliasDescription)
	}

	return nil
}
