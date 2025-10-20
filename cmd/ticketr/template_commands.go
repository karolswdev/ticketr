package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/karolswdev/ticktr/internal/adapters/jira"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/karolswdev/ticktr/internal/templates"
	"github.com/spf13/cobra"
)

var (
	// Template command flags
	templateDryRun bool

	// templateCmd represents the main template command
	templateCmd = &cobra.Command{
		Use:   "template",
		Short: "Manage and apply ticket templates",
		Long: `Manage and apply YAML templates for creating ticket hierarchies.

Templates allow you to define reusable ticket structures with variable
substitution, making it easy to create consistent ticket hierarchies
for features, bugs, and other workflows.

Examples:
  ticketr template list
  ticketr template apply feature.yaml
  ticketr template validate bug-investigation.yaml`,
	}

	// templateListCmd lists available templates
	templateListCmd = &cobra.Command{
		Use:   "list",
		Short: "List available templates",
		Long: `List all templates in the templates directory.

Templates are stored in:
- Linux: ~/.local/share/ticketr/templates/
- macOS: ~/Library/Application Support/ticketr/templates/
- Windows: %LOCALAPPDATA%\ticketr\templates\

Create new templates by adding .yaml files to this directory.`,
		RunE: runTemplateList,
	}

	// templateApplyCmd applies a template
	templateApplyCmd = &cobra.Command{
		Use:   "apply <template-file>",
		Short: "Apply a template to create tickets",
		Long: `Apply a template to create a ticket hierarchy in Jira.

The command will:
1. Load and validate the template
2. Extract required variables
3. Prompt for variable values
4. Create tickets in Jira (Epic, Stories, Tasks)
5. Display created ticket IDs

Use --dry-run to validate the template without creating tickets.

Examples:
  ticketr template apply feature.yaml
  ticketr template apply ~/.local/share/ticketr/templates/bug-fix.yaml
  ticketr template apply feature.yaml --dry-run`,
		Args: cobra.ExactArgs(1),
		RunE: runTemplateApply,
	}

	// templateValidateCmd validates a template
	templateValidateCmd = &cobra.Command{
		Use:   "validate <template-file>",
		Short: "Validate a template without applying it",
		Long: `Validate a template file's syntax and structure.

This checks:
- YAML syntax is valid
- Required fields are present
- Template structure is correct
- Variables can be extracted

Does not create any tickets in Jira.

Example:
  ticketr template validate feature.yaml`,
		Args: cobra.ExactArgs(1),
		RunE: runTemplateValidate,
	}
)

func init() {
	// Add subcommands
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateApplyCmd)
	templateCmd.AddCommand(templateValidateCmd)

	// Flags for apply command
	templateApplyCmd.Flags().BoolVar(&templateDryRun, "dry-run", false, "Validate template without creating tickets")
}

// initTemplateService creates a new TemplateService instance
func initTemplateService() (*services.TemplateService, error) {
	// Get PathResolver
	pathResolver, err := services.GetPathResolver()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize path resolver: %w", err)
	}

	// Get current workspace
	workspaceSvc, err := initWorkspaceService()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize workspace service: %w", err)
	}

	workspace, err := workspaceSvc.Current()
	if err != nil {
		return nil, fmt.Errorf("no active workspace - create one with 'ticketr workspace create': %w", err)
	}

	// Get credentials for current workspace
	config, err := workspaceSvc.GetConfig(workspace.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get workspace credentials: %w", err)
	}

	// Create Jira client from workspace config
	jiraClient, err := jira.NewJiraAdapterFromConfig(config, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Jira client: %w", err)
	}

	// Create template service
	svc := services.NewTemplateService(jiraClient, pathResolver)

	return svc, nil
}

// runTemplateList lists available templates
func runTemplateList(cmd *cobra.Command, args []string) error {
	// For list, we don't need Jira client
	pathResolver, err := services.GetPathResolver()
	if err != nil {
		return fmt.Errorf("failed to initialize path resolver: %w", err)
	}

	// Create a minimal service with nil Jira client (not needed for listing)
	svc := services.NewTemplateService(nil, pathResolver)

	templates, err := svc.ListTemplates(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list templates: %w", err)
	}

	if len(templates) == 0 {
		pathResolver, _ := services.GetPathResolver()
		fmt.Println("No templates found.")
		fmt.Printf("\nCreate templates in: %s\n", pathResolver.TemplatesDir())
		fmt.Println("\nExample template structure:")
		fmt.Println(`  name: feature
  structure:
    epic:
      title: "Feature: {{.Name}}"
      description: "Epic for {{.Name}}"
    stories:
      - title: "Frontend: {{.Name}}"
        tasks:
          - "Component implementation"
          - "Unit tests"`)
		return nil
	}

	// Print templates in table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tEPIC\tSTORIES\tTASKS\tDESCRIPTION")
	fmt.Fprintln(w, "----\t----\t-------\t-----\t-----------")

	for _, tmpl := range templates {
		epicIndicator := ""
		if tmpl.HasEpic {
			epicIndicator = "Yes"
		}

		// Truncate description if too long
		desc := tmpl.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}

		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\n",
			tmpl.Name,
			epicIndicator,
			tmpl.StoryCount,
			tmpl.TaskCount,
			desc,
		)
	}

	w.Flush()
	return nil
}

// runTemplateApply applies a template to create tickets
func runTemplateApply(cmd *cobra.Command, args []string) error {
	templatePath := args[0]

	svc, err := initTemplateService()
	if err != nil {
		return err
	}

	// Load template
	fmt.Printf("Loading template: %s\n", templatePath)
	tmpl, err := svc.LoadTemplate(context.Background(), templatePath)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	// Validate template
	if err := svc.ValidateTemplate(context.Background(), tmpl); err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	fmt.Printf("Template '%s' loaded successfully\n\n", tmpl.Name)

	// Show template summary
	if tmpl.HasEpic() {
		fmt.Println("Will create:")
		fmt.Println("- 1 Epic")
		if tmpl.StoryCount() > 0 {
			fmt.Printf("- %d Stories\n", tmpl.StoryCount())
		}
		if tmpl.TotalTaskCount() > 0 {
			fmt.Printf("- %d Tasks\n", tmpl.TotalTaskCount())
		}
	} else {
		fmt.Printf("Will create %d Stories", tmpl.StoryCount())
		if tmpl.TotalTaskCount() > 0 {
			fmt.Printf(" with %d Tasks", tmpl.TotalTaskCount())
		}
		fmt.Println()
	}
	fmt.Println()

	// Extract variables
	vars := templates.ExtractVariables(tmpl)

	if len(vars) > 0 {
		fmt.Printf("Template requires %d variable(s): %v\n\n", len(vars), vars)

		// Prompt for variables
		values := make(map[string]string)
		reader := bufio.NewReader(os.Stdin)

		for _, varName := range vars {
			fmt.Printf("Enter value for {{.%s}}: ", varName)
			value, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			values[varName] = strings.TrimSpace(value)
		}

		// Validate all variables provided
		if err := templates.ValidateVariables(tmpl, values); err != nil {
			return fmt.Errorf("variable validation failed: %w", err)
		}

		fmt.Println()

		// Dry run mode - stop here
		if templateDryRun {
			substituted, err := templates.Substitute(tmpl, values)
			if err != nil {
				return fmt.Errorf("substitution failed: %w", err)
			}

			fmt.Println("DRY RUN - Preview of substituted template:")
			fmt.Println()

			if substituted.Structure.Epic != nil {
				fmt.Printf("Epic: %s\n", substituted.Structure.Epic.Title)
				if substituted.Structure.Epic.Description != "" {
					fmt.Printf("  Description: %s\n", truncateString(substituted.Structure.Epic.Description, 80))
				}
				for _, task := range substituted.Structure.Epic.Tasks {
					fmt.Printf("  - %s\n", task)
				}
				fmt.Println()
			}

			for i, story := range substituted.Structure.Stories {
				fmt.Printf("Story %d: %s\n", i+1, story.Title)
				if story.Description != "" {
					fmt.Printf("  Description: %s\n", truncateString(story.Description, 80))
				}
				for _, task := range story.Tasks {
					fmt.Printf("  - %s\n", task)
				}
				fmt.Println()
			}

			fmt.Println("No tickets created (dry-run mode)")
			return nil
		}

		// Apply template
		fmt.Println("Creating tickets in Jira...")
		result, err := svc.ApplyTemplate(context.Background(), tmpl, values)
		if err != nil {
			return fmt.Errorf("failed to apply template: %w", err)
		}

		// Print results
		fmt.Println()
		fmt.Printf("Successfully created %d ticket(s):\n\n", result.TotalCreated)

		if result.EpicID != "" {
			fmt.Printf("Epic: %s\n", result.EpicID)
		}

		if len(result.StoryIDs) > 0 {
			fmt.Println("Stories:")
			for _, id := range result.StoryIDs {
				fmt.Printf("  - %s\n", id)
			}
		}

		if len(result.TaskIDs) > 0 {
			fmt.Println("Tasks:")
			for _, id := range result.TaskIDs {
				fmt.Printf("  - %s\n", id)
			}
		}
	} else {
		// No variables, apply directly
		if templateDryRun {
			fmt.Println("DRY RUN - Template has no variables")
			fmt.Println("Template is valid and ready to apply")
			return nil
		}

		fmt.Println("Applying template (no variables required)...")
		result, err := svc.ApplyTemplate(context.Background(), tmpl, map[string]string{})
		if err != nil {
			return fmt.Errorf("failed to apply template: %w", err)
		}

		fmt.Printf("\nSuccessfully created %d ticket(s)\n", result.TotalCreated)
	}

	return nil
}

// runTemplateValidate validates a template without applying it
func runTemplateValidate(cmd *cobra.Command, args []string) error {
	templatePath := args[0]

	// For validate, we don't need Jira client
	pathResolver, err := services.GetPathResolver()
	if err != nil {
		return fmt.Errorf("failed to initialize path resolver: %w", err)
	}

	// Create a minimal service with nil Jira client (not needed for validation)
	svc := services.NewTemplateService(nil, pathResolver)

	// Load template
	tmpl, err := svc.LoadTemplate(context.Background(), templatePath)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	// Validate template
	if err := svc.ValidateTemplate(context.Background(), tmpl); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Extract variables
	vars := templates.ExtractVariables(tmpl)

	// Print validation results
	fmt.Printf("Template '%s' is valid\n\n", tmpl.Name)
	fmt.Println("Structure:")
	if tmpl.HasEpic() {
		fmt.Println("- Epic: Yes")
	}
	fmt.Printf("- Stories: %d\n", tmpl.StoryCount())
	fmt.Printf("- Tasks: %d\n", tmpl.TotalTaskCount())

	if len(vars) > 0 {
		fmt.Printf("\nRequired variables (%d):\n", len(vars))
		for _, v := range vars {
			fmt.Printf("  - {{.%s}}\n", v)
		}
	} else {
		fmt.Println("\nNo variables required")
	}

	return nil
}

// Helper functions

// truncateString truncates a string to the specified length
func truncateString(s string, maxLen int) string {
	// Remove newlines for single-line display
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.TrimSpace(s)

	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
