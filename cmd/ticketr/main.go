package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/karolswdev/ticktr/internal/adapters/filesystem"
	"github.com/karolswdev/ticktr/internal/adapters/jira"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/karolswdev/ticktr/internal/core/validation"
	"github.com/karolswdev/ticktr/internal/logging"
	"github.com/karolswdev/ticktr/internal/migration"
	"github.com/karolswdev/ticktr/internal/state"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile            string
	verbose            bool
	forcePartialUpload bool
	logger             logging.Logger

	// Pull command flags
	pullProject string
	pullEpic    string
	pullJQL     string
	pullAlias   string
	pullOutput  string
	pullForce   bool

	// Migrate command flags
	writeFlag bool

	rootCmd = &cobra.Command{
		Use:   "ticketr",
		Short: "A tool for managing JIRA tickets as code",
		Long: `Ticketr is a command-line tool that allows you to manage JIRA tickets
using Markdown files stored in version control.`,
	}

	pushCmd = &cobra.Command{
		Use:   "push [file]",
		Short: "Push tickets from Markdown to JIRA",
		Long:  `Read tickets from a Markdown file and create or update them in JIRA.`,
		Args:  cobra.ExactArgs(1),
		Run:   runPush,
	}

	pullCmd = &cobra.Command{
		Use:   "pull",
		Short: "Pull tickets from JIRA to Markdown",
		Long: `Fetch tickets from JIRA and intelligently merge them with your local file.

Detects conflicts when both local and remote tickets have changed. Use --force
to overwrite local changes with remote changes when conflicts occur.`,
		Run: runPull,
	}

	schemaCmd = &cobra.Command{
		Use:   "schema",
		Short: "Discover JIRA schema and generate configuration",
		Long:  `Connect to JIRA and generate field mappings for .ticketr.yaml configuration.`,
		Run:   runSchema,
	}

	migrateCmd = &cobra.Command{
		Use:   "migrate [file]",
		Short: "Convert legacy # STORY: format to # TICKET:",
		Long: `Migrates Markdown files from deprecated # STORY: schema to canonical # TICKET: schema.

By default, runs in dry-run mode showing preview of changes.
Use --write flag to apply changes to files.

Examples:
  ticketr migrate path/to/story.md          # Preview changes
  ticketr migrate path/to/story.md --write  # Apply changes
  ticketr migrate examples/*.md --write     # Batch migration`,
		Args: cobra.MinimumNArgs(1),
		Run:  runMigrate,
	}

	// Legacy commands for backward compatibility
	legacyCmd = &cobra.Command{
		Use:    "legacy",
		Hidden: true,
		Run:    runLegacy,
	}
)

// expandAlias expands a JQL alias name to its full JQL query.
func expandAlias(aliasName string) (string, error) {
	// Initialize alias service
	aliasSvc, err := initAliasService()
	if err != nil {
		// If alias service initialization fails, check if it's a predefined alias
		// and return it directly (for backward compatibility)
		if predefined := domain.GetPredefinedAlias(aliasName); predefined != nil {
			return predefined.JQL, nil
		}
		return "", fmt.Errorf("failed to initialize alias service: %w", err)
	}

	// Get current workspace ID
	var workspaceID string
	workspaceSvc, err := initWorkspaceService()
	if err == nil {
		if workspace, err := workspaceSvc.Current(); err == nil {
			workspaceID = workspace.ID
		}
	}

	// Expand the alias
	jql, err := aliasSvc.ExpandAlias(aliasName, workspaceID)
	if err != nil {
		return "", fmt.Errorf("alias not found or expansion failed: %w", err)
	}

	return jql, nil
}

// initJiraAdapter initializes a Jira adapter using workspace credentials if available,
// otherwise falls back to environment variables for backward compatibility.
func initJiraAdapter(fieldMappings map[string]interface{}) (ports.JiraPort, error) {
	// Try to load workspace configuration first
	workspaceService, err := initWorkspaceService()
	if err == nil {
		// Workspace service available, try to get current workspace
		currentWorkspace, err := workspaceService.Current()
		if err == nil {
			// We have a current workspace, get its config
			config, err := workspaceService.GetConfig(currentWorkspace.Name)
			if err == nil {
				// Successfully retrieved workspace credentials
				return jira.NewJiraAdapterFromConfig(config, fieldMappings)
			}
			// Failed to get credentials, fall through to env vars
			fmt.Fprintf(os.Stderr, "Warning: Failed to retrieve workspace credentials: %v\n", err)
			fmt.Fprintf(os.Stderr, "Falling back to environment variables\n\n")
		}
		// No current workspace, fall through to env vars
	}
	// Workspace system not available or no workspace configured
	// Fall back to environment variables (v2 compatibility)
	return jira.NewJiraAdapterWithConfig(fieldMappings)
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .ticketr.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")

	// Push command flags
	pushCmd.Flags().BoolVar(&forcePartialUpload, "force-partial-upload", false, "continue processing even if some items fail")

	// Pull command flags
	pullCmd.Flags().StringVar(&pullProject, "project", "", "JIRA project key to pull from")
	pullCmd.Flags().StringVar(&pullEpic, "epic", "", "JIRA epic key to pull tickets from")
	pullCmd.Flags().StringVar(&pullJQL, "jql", "", "JQL query to filter tickets")
	pullCmd.Flags().StringVar(&pullAlias, "alias", "", "Use a named JQL alias (e.g., mine, sprint, blocked)")
	pullCmd.Flags().StringVarP(&pullOutput, "output", "o", "pulled_tickets.md", "output file path")
	pullCmd.Flags().BoolVar(&pullForce, "force", false, "Force overwrite local changes with remote changes when conflicts are detected")

	// Migrate command flags
	migrateCmd.Flags().BoolVarP(&writeFlag, "write", "w", false, "Write changes to files")

	// Add commands to root
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(schemaCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(legacyCmd)
	rootCmd.AddCommand(workspaceCmd)
	rootCmd.AddCommand(credentialsCmd)
	rootCmd.AddCommand(bulkCmd)
	rootCmd.AddCommand(templateCmd)
	rootCmd.AddCommand(aliasCmd)

	// Legacy flags for backward compatibility
	rootCmd.PersistentFlags().StringP("file", "f", "", "Path to the input Markdown file (deprecated, use 'push' command)")
	rootCmd.PersistentFlags().Bool("list-issue-types", false, "List available issue types (deprecated)")
	rootCmd.PersistentFlags().String("check-fields", "", "Check required fields for issue type (deprecated)")
	rootCmd.PersistentFlags().MarkHidden("file")
	rootCmd.PersistentFlags().MarkHidden("list-issue-types")
	rootCmd.PersistentFlags().MarkHidden("check-fields")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for config in current directory
		viper.AddConfigPath(".")
		viper.SetConfigName(".ticketr")
		viper.SetConfigType("yaml")
	}

	// Environment variables override config
	viper.SetEnvPrefix("JIRA")
	viper.AutomaticEnv()

	// Read config file if it exists
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			log.Printf("Using config file: %s", viper.ConfigFileUsed())
		}
	}

	// Configure logging
	if verbose {
		log.SetFlags(log.Ltime | log.Lshortfile | log.Lmicroseconds)
		log.Println("Verbose mode enabled")
	} else {
		log.SetFlags(log.Ltime)
	}
}

func runPush(cmd *cobra.Command, args []string) {
	inputFile := args[0]

	if logger != nil {
		logger.Section("PUSH COMMAND")
		logger.Info("Input file: %s", inputFile)
		logger.Info("Force partial upload: %v", forcePartialUpload)
	}

	// Initialize repository
	repo := filesystem.NewFileRepository()

	// Pre-flight validation: Parse tickets first for validation
	tickets, err := repo.GetTickets(inputFile)
	if err != nil {
		fmt.Printf("Error reading tickets from file: %v\n", err)
		os.Exit(1)
	}

	// Initialize validator and run pre-flight validation
	validator := validation.NewValidator()
	validationErrors := validator.ValidateTickets(tickets)
	if len(validationErrors) > 0 {
		if forcePartialUpload {
			// Downgrade to warnings
			fmt.Println("Warning: Validation warnings (processing will continue with --force-partial-upload):")
			for _, vErr := range validationErrors {
				fmt.Printf("  - %s\n", vErr.Error())
			}
			fmt.Printf("\n%d validation warning(s) found. Some items may fail during upload.\n", len(validationErrors))
		} else {
			// Hard fail without force flag
			fmt.Println("Validation errors found:")
			for _, vErr := range validationErrors {
				fmt.Printf("  - %s\n", vErr.Error())
			}
			fmt.Printf("\n%d validation error(s) found. Fix these issues before pushing to JIRA.\n", len(validationErrors))
			fmt.Println("Tip: Use --force-partial-upload to continue despite validation errors.")
			os.Exit(1)
		}
	}

	// Initialize Jira adapter (workspace credentials or environment variables)
	jiraAdapter, err := initJiraAdapter(nil)
	if err != nil {
		fmt.Printf("Error initializing Jira adapter: %v\n", err)
		fmt.Println("\nOption 1: Use workspace credentials (recommended):")
		fmt.Println("  ticketr workspace create <name> --url <jira-url> --project <key> --username <email> --token <api-token>")
		fmt.Println("\nOption 2: Set environment variables (legacy):")
		fmt.Println("  - JIRA_URL")
		fmt.Println("  - JIRA_EMAIL")
		fmt.Println("  - JIRA_API_KEY")
		fmt.Println("  - JIRA_PROJECT_KEY")
		fmt.Println("\nOptional environment variables:")
		fmt.Println("  - JIRA_STORY_TYPE (defaults to 'Task')")
		fmt.Println("  - JIRA_SUBTASK_TYPE (defaults to 'Sub-task')")
		os.Exit(1)
	}

	// Initialize state manager (uses PathResolver for global paths)
	stateManager := state.NewStateManager("")

	// Initialize push service with state management
	service := services.NewPushService(repo, jiraAdapter, stateManager)

	// Process tickets
	options := services.ProcessOptions{
		ForcePartialUpload: forcePartialUpload,
	}

	result, err := service.PushTickets(inputFile, options)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}

	// Print summary
	fmt.Println("\n=== Summary ===")
	if result.TicketsCreated > 0 {
		fmt.Printf("Tickets created: %d\n", result.TicketsCreated)
	}
	if result.TicketsUpdated > 0 {
		fmt.Printf("Tickets updated: %d\n", result.TicketsUpdated)
	}
	if result.TasksCreated > 0 {
		fmt.Printf("Tasks created: %d\n", result.TasksCreated)
	}
	if result.TasksUpdated > 0 {
		fmt.Printf("Tasks updated: %d\n", result.TasksUpdated)
	}

	// Print errors if any
	if len(result.Errors) > 0 {
		fmt.Printf("\n=== Errors (%d) ===\n", len(result.Errors))
		for _, err := range result.Errors {
			fmt.Printf("  - %s\n", err)
		}

		if !forcePartialUpload {
			os.Exit(2)
		}
	}

	// Log execution summary
	if logger != nil {
		logger.Section("EXECUTION SUMMARY")
		logger.Info("Tickets created: %d", result.TicketsCreated)
		logger.Info("Tickets updated: %d", result.TicketsUpdated)
		logger.Info("Tasks created: %d", result.TasksCreated)
		logger.Info("Tasks updated: %d", result.TasksUpdated)

		if len(result.Errors) > 0 {
			logger.Section("ERRORS")
			for _, err := range result.Errors {
				logger.Error("%s", err)
			}
		}
	}

	fmt.Println("\nProcessing complete!")
}

// runPull handles the pull command
func runPull(cmd *cobra.Command, args []string) {
	if logger != nil {
		logger.Section("PULL COMMAND")
		logger.Info("Project: %s", pullProject)
		logger.Info("Epic: %s", pullEpic)
		logger.Info("JQL: %s", pullJQL)
		logger.Info("Alias: %s", pullAlias)
		logger.Info("Output: %s", pullOutput)
		logger.Info("Force: %v", pullForce)
	}

	// Initialize JIRA adapter with field mappings from config
	fieldMappings := viper.GetStringMap("field_mappings")

	// Convert to proper format for adapter
	mappings := make(map[string]interface{})
	for key, value := range fieldMappings {
		mappings[key] = value
	}

	// Initialize Jira adapter (workspace credentials or environment variables)
	jiraAdapter, err := initJiraAdapter(mappings)
	if err != nil {
		fmt.Printf("Error initializing JIRA adapter: %v\n", err)
		fmt.Println("\nOption 1: Use workspace credentials (recommended):")
		fmt.Println("  ticketr workspace create <name> --url <jira-url> --project <key> --username <email> --token <api-token>")
		fmt.Println("\nOption 2: Set environment variables (legacy):")
		fmt.Println("  - JIRA_URL")
		fmt.Println("  - JIRA_EMAIL")
		fmt.Println("  - JIRA_API_KEY")
		fmt.Println("  - JIRA_PROJECT_KEY")
		os.Exit(1)
	}

	// Get project key from flag, workspace, or environment (in that order)
	projectKey := pullProject
	if projectKey == "" {
		// Try to get project key from current workspace
		workspaceService, err := initWorkspaceService()
		if err == nil {
			if currentWorkspace, err := workspaceService.Current(); err == nil {
				if config, err := workspaceService.GetConfig(currentWorkspace.Name); err == nil {
					projectKey = config.ProjectKey
					if verbose {
						log.Printf("Using project key from workspace '%s': %s", currentWorkspace.Name, projectKey)
					}
				}
			}
		}
	}
	if projectKey == "" {
		// Fall back to environment variable
		projectKey = os.Getenv("JIRA_PROJECT_KEY")
	}
	if projectKey == "" {
		fmt.Println("Error: Project key is required. Use --project flag, configure a workspace, or set JIRA_PROJECT_KEY environment variable")
		os.Exit(1)
	}

	// Handle alias expansion if --alias flag is provided
	jql := pullJQL
	if pullAlias != "" {
		// Cannot use both --jql and --alias
		if pullJQL != "" {
			fmt.Println("Error: Cannot use both --jql and --alias flags. Please use one or the other.")
			os.Exit(1)
		}

		// Expand the alias to JQL
		aliasJQL, err := expandAlias(pullAlias)
		if err != nil {
			fmt.Printf("Error expanding alias '%s': %v\n", pullAlias, err)
			fmt.Println("\nAvailable aliases:")
			fmt.Println("  mine    - Tickets assigned to you")
			fmt.Println("  sprint  - Tickets in active sprints")
			fmt.Println("  blocked - Blocked tickets")
			fmt.Println("\nList all aliases with: ticketr alias list")
			fmt.Println("Create an alias with: ticketr alias create <name> <jql>")
			os.Exit(1)
		}
		jql = aliasJQL

		if verbose {
			log.Printf("Expanded alias '%s' to JQL: %s", pullAlias, jql)
		}
	}

	// Construct JQL based on flags
	if pullEpic != "" {
		epicFilter := fmt.Sprintf(`"Epic Link" = "%s"`, pullEpic)
		if jql != "" {
			jql = fmt.Sprintf("%s AND %s", jql, epicFilter)
		} else {
			jql = epicFilter
		}
	}

	// Log the query if verbose
	if verbose {
		log.Printf("Pulling tickets from project: %s", projectKey)
		if jql != "" {
			log.Printf("Using JQL filter: %s", jql)
		}
	}

	// Initialize state manager (uses PathResolver for global paths)
	stateManager := state.NewStateManager("")

	// Initialize file repository
	fileRepo := filesystem.NewFileRepository()

	// Create pull service
	pullService := services.NewPullService(jiraAdapter, fileRepo, stateManager)

	// Create progress callback for user feedback
	progressCallback := func(current, total int, message string) {
		if message != "" {
			fmt.Println(message)
		}
		if total > 0 && current > 0 {
			// Show progress counter for large datasets
			fmt.Printf("\rPulling tickets: %d/%d", current, total)
			if current == total {
				// Complete the line when done
				fmt.Println()
			}
		}
	}

	// Execute pull
	result, err := pullService.Pull(pullOutput, services.PullOptions{
		ProjectKey:       projectKey,
		JQL:              jql,
		EpicKey:          pullEpic,
		Force:            pullForce,
		ProgressCallback: progressCallback,
	})

	// Handle errors and conflicts
	if err != nil {
		if errors.Is(err, services.ErrConflictDetected) {
			fmt.Println("⚠️  Conflict detected! The following tickets have both local and remote changes:")
			for _, ticketID := range result.Conflicts {
				fmt.Printf("  - %s\n", ticketID)
			}
			fmt.Println("\nTo force overwrite local changes with remote changes, use --force flag")
			os.Exit(1)
		}
		fmt.Printf("Error pulling tickets: %v\n", err)
		os.Exit(1)
	}

	// Print summary
	fmt.Printf("\nSuccessfully updated %s\n", pullOutput)
	totalCount := result.TicketsPulled + result.TicketsUpdated + result.TicketsSkipped
	fmt.Printf("  - %d ticket(s) pulled\n", totalCount)
	if len(result.Conflicts) > 0 {
		fmt.Printf("  - %d conflict(s) detected\n", len(result.Conflicts))
	} else {
		fmt.Printf("  - 0 conflicts detected\n")
	}
	if result.TicketsUpdated > 0 {
		fmt.Printf("  - %d ticket(s) updated\n", result.TicketsUpdated)
	}

	// Log execution summary
	if logger != nil {
		logger.Section("EXECUTION SUMMARY")
		logger.Info("Tickets pulled: %d", result.TicketsPulled)
		logger.Info("Tickets updated: %d", result.TicketsUpdated)
		logger.Info("Tickets skipped: %d", result.TicketsSkipped)
		logger.Info("Conflicts: %d", len(result.Conflicts))
	}
}

// runSchema handles the schema discovery command
func runSchema(cmd *cobra.Command, args []string) {
	// Initialize JIRA adapter (workspace credentials or environment variables)
	jiraAdapter, err := initJiraAdapter(nil)
	if err != nil {
		fmt.Printf("Error initializing JIRA adapter: %v\n", err)
		fmt.Println("\nOption 1: Use workspace credentials (recommended):")
		fmt.Println("  ticketr workspace create <name> --url <jira-url> --project <key> --username <email> --token <api-token>")
		fmt.Println("\nOption 2: Set environment variables (legacy):")
		fmt.Println("  - JIRA_URL")
		fmt.Println("  - JIRA_EMAIL")
		fmt.Println("  - JIRA_API_KEY")
		fmt.Println("  - JIRA_PROJECT_KEY")
		os.Exit(1)
	}

	// Get project issue types
	issueTypes, err := jiraAdapter.GetProjectIssueTypes()
	if err != nil {
		fmt.Printf("Error fetching project issue types: %v\n", err)
		os.Exit(1)
	}

	// Start building the YAML output
	fmt.Println("# Generated field mappings for .ticketr.yaml")
	fmt.Println("field_mappings:")

	// Always include standard fields
	fmt.Println("  \"Type\": \"issuetype\"")
	fmt.Println("  \"Project\": \"project\"")
	fmt.Println("  \"Summary\": \"summary\"")
	fmt.Println("  \"Description\": \"description\"")
	fmt.Println("  \"Assignee\": \"assignee\"")
	fmt.Println("  \"Reporter\": \"reporter\"")
	fmt.Println("  \"Priority\": \"priority\"")
	fmt.Println("  \"Labels\": \"labels\"")
	fmt.Println("  \"Components\": \"components\"")
	fmt.Println("  \"Fix Version\": \"fixVersions\"")
	fmt.Println("  \"Sprint\": \"customfield_10020\"  # Common sprint field")

	// Collect custom fields from all issue types
	customFieldsMap := make(map[string]map[string]interface{})

	for projectKey, types := range issueTypes {
		if verbose {
			fmt.Fprintf(os.Stderr, "Processing project: %s\n", projectKey)
		}
		for _, issueType := range types {
			if verbose {
				fmt.Fprintf(os.Stderr, "  Fetching fields for issue type: %s\n", issueType)
			}

			fields, err := jiraAdapter.GetIssueTypeFields(issueType)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Could not fetch fields for %s: %v\n", issueType, err)
				continue
			}

			// Process optional fields (custom fields are usually here)
			if optionalInterface, ok := fields["optional"]; ok {
				if optional, ok := optionalInterface.([]interface{}); ok {
					for _, field := range optional {
						if fieldMap, ok := field.(map[string]interface{}); ok {
							processFieldForSchema(fieldMap, customFieldsMap)
						}
					}
				}
			}
		}
	}

	// Output discovered custom fields
	for fieldName, fieldInfo := range customFieldsMap {
		id := fieldInfo["id"].(string)
		fieldType := fieldInfo["type"].(string)

		// Format based on type
		if fieldType == "string" || fieldType == "option" {
			fmt.Printf("  \"%s\": \"%s\"\n", fieldName, id)
		} else {
			fmt.Printf("  \"%s\":\n", fieldName)
			fmt.Printf("    id: \"%s\"\n", id)
			fmt.Printf("    type: \"%s\"\n", fieldType)
		}
	}

	// Add example sync configuration
	fmt.Println("\n# Example sync configuration")
	fmt.Println("sync:")
	fmt.Println("  pull:")
	fmt.Println("    # Fields to pull from JIRA to Markdown")
	fmt.Println("    fields:")
	fmt.Println("      - \"Story Points\"")
	fmt.Println("      - \"Sprint\"")
	fmt.Println("      - \"Priority\"")
	fmt.Println("  ignored_fields:")
	fmt.Println("    # Fields to never pull")
	fmt.Println("    - \"updated\"")
	fmt.Println("    - \"created\"")
}

// processFieldForSchema extracts relevant field information for schema generation
func processFieldForSchema(field map[string]interface{}, customFieldsMap map[string]map[string]interface{}) {
	key, hasKey := field["key"].(string)
	if !hasKey || !strings.HasPrefix(key, "customfield_") {
		return
	}

	name := ""
	if nameVal, ok := field["name"]; ok {
		name = nameVal.(string)
	}

	if name == "" || name == "Development" || strings.Contains(name, "[CHART]") {
		return // Skip system or chart fields
	}

	// Determine field type
	fieldType := "string" // default
	if schema, ok := field["schema"]; ok {
		if schemaMap, ok := schema.(map[string]interface{}); ok {
			if typeVal, ok := schemaMap["type"]; ok {
				switch typeVal.(string) {
				case "number":
					fieldType = "number"
				case "array":
					fieldType = "array"
				case "option":
					fieldType = "option"
				}
			}
		}
	}

	// Store field info if not already present or if this is a better match
	if _, exists := customFieldsMap[name]; !exists {
		customFieldsMap[name] = map[string]interface{}{
			"id":   key,
			"type": fieldType,
		}
	}
}

// runMigrate handles the migrate command
func runMigrate(cmd *cobra.Command, args []string) {
	// Create migrator with DryRun based on writeFlag
	migrator := &migration.Migrator{
		DryRun: !writeFlag,
	}

	// Track overall success/failure
	hasErrors := false
	totalFiles := 0
	totalChanges := 0

	// Process each file argument (supports glob patterns)
	for _, pattern := range args {
		// Expand glob pattern
		matches, err := filepath.Glob(pattern)
		if err != nil {
			fmt.Printf("Error processing pattern '%s': %v\n", pattern, err)
			hasErrors = true
			continue
		}

		if len(matches) == 0 {
			// No glob match, treat as literal file path
			matches = []string{pattern}
		}

		// Process each matched file
		for _, filePath := range matches {
			totalFiles++

			// Get absolute path for display
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				absPath = filePath
			}

			// Perform migration
			content, changed, err := migrator.MigrateFile(filePath)
			if err != nil {
				fmt.Printf("Error migrating %s: %v\n", absPath, err)
				hasErrors = true
				continue
			}

			if !changed {
				if verbose {
					fmt.Printf("No changes needed: %s\n", absPath)
				}
				continue
			}

			totalChanges++

			// If dry-run, show preview
			if migrator.DryRun {
				originalContent, _ := os.ReadFile(filePath)
				preview := migrator.PreviewDiff(absPath, string(originalContent), content)
				fmt.Println(preview)
			} else {
				// Write the migration
				err = migrator.WriteMigration(filePath, content)
				if err != nil {
					fmt.Printf("Error writing %s: %v\n", absPath, err)
					hasErrors = true
					continue
				}

				// Count how many replacements were made
				changeCount := strings.Count(content, "# TICKET:") - strings.Count(string(content), "# STORY:")
				if changeCount < 0 {
					changeCount = -changeCount
				}

				fmt.Printf("Migrated: %s (%d change(s))\n", absPath, changeCount)
			}
		}
	}

	// Print summary
	if totalFiles == 0 {
		fmt.Println("No files matched the provided pattern(s)")
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("\nProcessed %d file(s), %d file(s) with changes\n", totalFiles, totalChanges)
	}

	if hasErrors {
		os.Exit(1)
	}
}

// runLegacy handles the old command-line interface for backward compatibility
func runLegacy(cmd *cobra.Command, args []string) {
	// Check for legacy flags
	inputFile, _ := cmd.Flags().GetString("file")
	listIssueTypes, _ := cmd.Flags().GetBool("list-issue-types")
	checkFields, _ := cmd.Flags().GetString("check-fields")

	// Initialize Jira adapter for legacy commands (workspace credentials or environment variables)
	jiraAdapter, err := initJiraAdapter(nil)
	if err != nil {
		fmt.Printf("Error initializing Jira adapter: %v\n", err)
		fmt.Println("\nOption 1: Use workspace credentials (recommended):")
		fmt.Println("  ticketr workspace create <name> --url <jira-url> --project <key> --username <email> --token <api-token>")
		fmt.Println("\nOption 2: Set environment variables (legacy):")
		fmt.Println("  - JIRA_URL")
		fmt.Println("  - JIRA_EMAIL")
		fmt.Println("  - JIRA_API_KEY")
		fmt.Println("  - JIRA_PROJECT_KEY")
		os.Exit(1)
	}

	// Handle list-issue-types
	if listIssueTypes {
		fmt.Println("Fetching issue types from JIRA...")
		issueTypesInfo, err := jiraAdapter.GetProjectIssueTypes()
		if err != nil {
			fmt.Printf("Error fetching issue types: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n" + "=" + string(make([]byte, 50)))
		if projectName, ok := issueTypesInfo["project"]; ok && len(projectName) > 0 {
			fmt.Printf("Project: %s", projectName[0])
			if key, ok := issueTypesInfo["key"]; ok && len(key) > 0 {
				fmt.Printf(" (%s)\n", key[0])
			}
		}
		fmt.Println("=" + string(make([]byte, 50)))

		if issueTypes, ok := issueTypesInfo["types"]; ok {
			fmt.Println("\nAvailable Issue Types:")
			for _, issueType := range issueTypes {
				fmt.Printf("  - %s\n", issueType)
			}
		}

		if subtaskTypes, ok := issueTypesInfo["subtasks"]; ok && len(subtaskTypes) > 0 {
			fmt.Println("\nAvailable Subtask Types:")
			for _, subtaskType := range subtaskTypes {
				fmt.Printf("  - %s\n", subtaskType)
			}
		}
		return
	}

	// Handle check-fields
	if checkFields != "" {
		fmt.Printf("Checking fields for issue type: %s\n", checkFields)
		fields, err := jiraAdapter.GetIssueTypeFields(checkFields)
		if err != nil {
			fmt.Printf("Error fetching fields: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n%s Issue Type Fields:\n", checkFields)
		fmt.Println("=" + string(make([]byte, 50)))

		if requiredInterface, ok := fields["required"]; ok {
			if required, ok := requiredInterface.([]interface{}); ok && len(required) > 0 {
				fmt.Println("\nRequired Fields:")
				for _, field := range required {
					if fieldMap, ok := field.(map[string]interface{}); ok {
						printFieldInfo(fieldMap)
					}
				}
			}
		}

		if optionalInterface, ok := fields["optional"]; ok {
			if optional, ok := optionalInterface.([]interface{}); ok && len(optional) > 0 {
				fmt.Println("\nOptional Fields:")
				for _, field := range optional {
					if fieldMap, ok := field.(map[string]interface{}); ok {
						printFieldInfo(fieldMap)
					}
				}
			}
		}
		return
	}

	// Handle file processing (default behavior)
	if inputFile != "" {
		runPush(cmd, []string{inputFile})
		return
	}

	// No valid command provided
	cmd.Help()
}

// printFieldInfo prints formatted field information
func printFieldInfo(field map[string]interface{}) {
	key := field["key"].(string)
	name := ""
	if n, ok := field["name"].(string); ok {
		name = n
	}

	fieldType := ""
	if t, ok := field["type"].(string); ok {
		fieldType = t
		if items, ok := field["items"].(string); ok {
			fieldType = fmt.Sprintf("%s[%s]", fieldType, items)
		}
	}

	fmt.Printf("\n  %s (%s)\n", name, key)
	if fieldType != "" {
		fmt.Printf("    Type: %s\n", fieldType)
	}

	if values, ok := field["allowedValues"].([]string); ok && len(values) > 0 {
		fmt.Printf("    Allowed Values: %s\n", strings.Join(values, ", "))
		if len(values) > 5 {
			fmt.Printf("    (showing first 5 of %d values)\n", len(values))
		}
	}
}

func main() {
	// Initialize logger
	logDir := os.Getenv("TICKETR_LOG_DIR")
	if logDir == "" {
		// Use PathResolver for global logs directory
		pathResolver, err := services.GetPathResolver()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to get path resolver: %v\n", err)
			// Fall back to local path for backward compatibility
			logDir = ".ticketr/logs"
		} else {
			logDir = pathResolver.LogsDir()
		}
	}

	fileLogger, err := logging.NewFileLogger(logging.LogConfig{
		LogDir:  logDir,
		Verbose: verbose,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to initialize file logging: %v\n", err)
		// Continue without file logging
	} else {
		logger = fileLogger
		defer fileLogger.Close()

		// Cleanup old logs (keep last 10)
		if err := logging.Cleanup(logDir, 10); err != nil {
			logger.Warn("Failed to cleanup old logs: %v", err)
		}
	}

	// Check for legacy usage (no subcommand)
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		// If first arg is not a flag and not a known command, assume it's a file (legacy)
		knownCommands := []string{"push", "pull", "schema", "migrate", "workspace", "credentials", "bulk", "template", "v3", "help", "completion"}
		isKnownCommand := false
		for _, cmd := range knownCommands {
			if os.Args[1] == cmd {
				isKnownCommand = true
				break
			}
		}

		if !isKnownCommand {
			// Legacy mode: no subcommand, treat as file argument
			// This maintains backward compatibility
		}
	} else if len(os.Args) == 1 {
		// No arguments at all, show help
		rootCmd.Help()
		return
	} else {
		// Check for legacy flags without subcommand
		hasLegacyFlag := false
		for _, arg := range os.Args[1:] {
			if strings.Contains(arg, "-file") || strings.Contains(arg, "-f=") ||
				strings.Contains(arg, "list-issue-types") || strings.Contains(arg, "check-fields") {
				hasLegacyFlag = true
				break
			}
		}

		if hasLegacyFlag {
			// Use legacy command handler
			runLegacy(rootCmd, os.Args[1:])
			return
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
