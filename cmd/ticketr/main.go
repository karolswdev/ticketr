// Package main provides the command-line interface for Ticketr.
// Ticketr is a tool for managing JIRA tickets as code using Markdown files.
// It supports bidirectional synchronization between local Markdown files and JIRA,
// enabling teams to version control their project management workflow.
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/karolswdev/ticktr/internal/adapters/filesystem"
	"github.com/karolswdev/ticktr/internal/adapters/jira"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/karolswdev/ticktr/internal/core/validation"
	"github.com/karolswdev/ticktr/internal/state"
	"github.com/karolswdev/ticktr/internal/renderer"
)

var (
	// cfgFile specifies an alternative configuration file
	cfgFile string
	// verbose enables detailed logging output
	verbose bool
	// forcePartialUpload continues processing even if some items fail
	forcePartialUpload bool
	
	// Pull command flags
	pullProject string // JIRA project key to pull from
	pullEpic    string // Epic key to filter tickets
	pullJQL     string // Custom JQL query for filtering
	pullOutput  string // Output file path for pulled tickets
	
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
		Long:  `Fetch tickets from JIRA and write them to a Markdown file.`,
		Run:   runPull,
	}
	
	schemaCmd = &cobra.Command{
		Use:   "schema",
		Short: "Discover JIRA schema and generate configuration",
		Long:  `Connect to JIRA and generate field mappings for .ticketr.yaml configuration.`,
		Run:   runSchema,
	}
	
	// Legacy commands for backward compatibility
	legacyCmd = &cobra.Command{
		Use:    "legacy",
		Hidden: true,
		Run:    runLegacy,
	}
)

// init initializes all command-line flags and commands.
// It sets up the command hierarchy and registers all available flags.
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
	pullCmd.Flags().StringVarP(&pullOutput, "output", "o", "pulled_tickets.md", "output file path")
	
	// Add commands to root
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(schemaCmd)
	rootCmd.AddCommand(legacyCmd)
	
	// Legacy flags for backward compatibility
	rootCmd.PersistentFlags().StringP("file", "f", "", "Path to the input Markdown file (deprecated, use 'push' command)")
	rootCmd.PersistentFlags().Bool("list-issue-types", false, "List available issue types (deprecated)")
	rootCmd.PersistentFlags().String("check-fields", "", "Check required fields for issue type (deprecated)")
	rootCmd.PersistentFlags().MarkHidden("file")
	rootCmd.PersistentFlags().MarkHidden("list-issue-types")
	rootCmd.PersistentFlags().MarkHidden("check-fields")
}

// initConfig loads configuration from file and environment variables.
// Configuration priority (highest to lowest):
// 1. Environment variables (JIRA_* prefix)
// 2. Configuration file (.ticketr.yaml)
// 3. Default values
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

// runPush handles the push command to sync tickets from Markdown to JIRA.
// It performs the following steps:
// 1. Pre-flight validation of ticket format and hierarchy
// 2. Connects to JIRA and validates credentials
// 3. Creates or updates tickets in JIRA
// 4. Updates the local file with JIRA IDs
//
// Parameters:
//   - cmd: The cobra command that triggered this function
//   - args: Command arguments (expects exactly one: the input file path)
func runPush(cmd *cobra.Command, args []string) {
	inputFile := args[0]
	
	// Initialize repository
	repo := filesystem.NewFileRepository()
	
	// Pre-flight validation: Parse tickets first to catch errors early
	// This prevents partial uploads due to malformed ticket definitions
	tickets, err := repo.GetTickets(inputFile)
	if err != nil {
		fmt.Printf("Error reading tickets from file: %v\n", err)
		os.Exit(1)
	}
	
	// Initialize validator and run pre-flight validation
	validator := validation.NewValidator()
	validationErrors := validator.ValidateTickets(tickets)
	if len(validationErrors) > 0 {
		fmt.Println("Validation errors found:")
		for _, vErr := range validationErrors {
			fmt.Printf("  - %s\n", vErr.Error())
		}
		fmt.Printf("\n%d validation error(s) found. Fix these issues before pushing to JIRA.\n", len(validationErrors))
		os.Exit(1)
	}
	
	// Initialize Jira adapter
	jiraAdapter, err := jira.NewJiraAdapter()
	if err != nil {
		fmt.Printf("Error initializing Jira adapter: %v\n", err)
		fmt.Println("\nMake sure the following environment variables are set:")
		fmt.Println("  - JIRA_URL")
		fmt.Println("  - JIRA_EMAIL")
		fmt.Println("  - JIRA_API_KEY")
		fmt.Println("  - JIRA_PROJECT_KEY")
		fmt.Println("\nOptional environment variables:")
		fmt.Println("  - JIRA_STORY_TYPE (defaults to 'Task')")
		fmt.Println("  - JIRA_SUBTASK_TYPE (defaults to 'Sub-task')")
		os.Exit(1)
	}
	
	// Initialize service
	service := services.NewTicketService(repo, jiraAdapter)
	
	// Process tickets
	options := services.ProcessOptions{
		ForcePartialUpload: forcePartialUpload,
	}
	
	result, err := service.ProcessTicketsWithOptions(inputFile, options)
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
	
	fmt.Println("\nProcessing complete!")
}

// runPull handles the pull command to sync tickets from JIRA to Markdown.
// It supports intelligent conflict detection when pulling into existing files,
// preserving local changes where possible and alerting on conflicts.
//
// Features:
//   - Pull tickets by project, epic, or custom JQL query
//   - Intelligent merge with conflict detection
//   - State tracking to identify changes
//
// Parameters:
//   - cmd: The cobra command that triggered this function
//   - args: Command arguments (none expected)
func runPull(cmd *cobra.Command, args []string) {
	// Initialize JIRA adapter with field mappings from config
	fieldMappings := viper.GetStringMap("field_mappings")
	
	// Convert to proper format for adapter
	mappings := make(map[string]interface{})
	for key, value := range fieldMappings {
		mappings[key] = value
	}
	
	jiraAdapter, err := jira.NewJiraAdapterWithConfig(mappings)
	if err != nil {
		fmt.Printf("Error initializing JIRA adapter: %v\n", err)
		fmt.Println("\nMake sure the following environment variables are set:")
		fmt.Println("  - JIRA_URL")
		fmt.Println("  - JIRA_EMAIL")
		fmt.Println("  - JIRA_API_KEY")
		fmt.Println("  - JIRA_PROJECT_KEY")
		os.Exit(1)
	}
	
	// Get project key from flag or environment
	projectKey := pullProject
	if projectKey == "" {
		projectKey = os.Getenv("JIRA_PROJECT_KEY")
	}
	if projectKey == "" {
		fmt.Println("Error: Project key is required. Use --project flag or set JIRA_PROJECT_KEY environment variable")
		os.Exit(1)
	}
	
	// Construct JQL based on flags
	jql := pullJQL
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
	
	// Check if output file exists to enable conflict detection
	fileRepo := filesystem.NewFileRepository()
	hasExistingFile := false
	var existingTickets []interface{} // We'll need to adapt this based on actual types
	
	if _, err := os.Stat(pullOutput); err == nil {
		hasExistingFile = true
		// Try to parse existing tickets for conflict detection
		if tickets, err := fileRepo.GetTickets(pullOutput); err == nil {
			existingTickets = make([]interface{}, len(tickets))
			for i, t := range tickets {
				existingTickets[i] = t
			}
		}
	}
	
	// Search for tickets from JIRA
	tickets, err := jiraAdapter.SearchTickets(projectKey, jql)
	if err != nil {
		fmt.Printf("Error searching tickets: %v\n", err)
		os.Exit(1)
	}
	
	if len(tickets) == 0 {
		fmt.Println("No tickets found matching the query")
		return
	}
	
	// Use intelligent pull service when merging with existing file
	// This provides conflict detection and safe merging capabilities
	if hasExistingFile && len(existingTickets) > 0 {
		// Initialize state manager for conflict detection
		stateManager := state.NewStateManager(".ticketr.state")
		
		// Create pull service if available
		pullService := services.NewPullService(jiraAdapter, fileRepo, stateManager)
		
		// Execute intelligent pull with conflict detection
		result, err := pullService.Pull(pullOutput, services.PullOptions{
			ProjectKey: projectKey,
			JQL:        jql,
			EpicKey:    pullEpic,
			Force:      false, // Could be a flag in the future
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
		
		// Print summary for intelligent pull
		fmt.Printf("Successfully updated %s\n", pullOutput)
		if result.TicketsPulled > 0 {
			fmt.Printf("  - %d new ticket(s) pulled from JIRA\n", result.TicketsPulled)
		}
		if result.TicketsUpdated > 0 {
			fmt.Printf("  - %d ticket(s) updated with remote changes\n", result.TicketsUpdated)
		}
		if result.TicketsSkipped > 0 {
			fmt.Printf("  - %d ticket(s) skipped (no changes or local changes preserved)\n", result.TicketsSkipped)
		}
		return
	}
	
	// Fallback to simple render-based pull (when no existing file or state)
	// Initialize renderer with field mappings
	ticketRenderer := renderer.NewRenderer(mappings)
	
	// Render tickets to markdown
	markdown := ticketRenderer.RenderMultiple(tickets)
	
	// Write to output file
	err = os.WriteFile(pullOutput, []byte(markdown), 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}
	
	// Print summary for simple pull
	fmt.Printf("Successfully pulled %d ticket(s) to %s\n", len(tickets), pullOutput)
	if verbose {
		fmt.Println("\nTickets pulled:")
		for _, ticket := range tickets {
			fmt.Printf("  - [%s] %s\n", ticket.JiraID, ticket.Title)
		}
	}
}

// runSchema discovers JIRA custom fields and generates configuration.
// It connects to JIRA, fetches all available fields for the project,
// and outputs a YAML configuration file template.
//
// This is essential for projects using custom fields like Story Points,
// Sprint, or any organization-specific fields.
//
// Parameters:
//   - cmd: The cobra command that triggered this function
//   - args: Command arguments (none expected)
func runSchema(cmd *cobra.Command, args []string) {
	// Initialize JIRA adapter
	jiraAdapter, err := jira.NewJiraAdapter()
	if err != nil {
		fmt.Printf("Error initializing JIRA adapter: %v\n", err)
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

// processFieldForSchema extracts relevant field information for schema generation.
// It filters out system fields and identifies custom field types for proper handling.
//
// Parameters:
//   - field: The JIRA field definition map
//   - customFieldsMap: Map to store discovered custom fields
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

// runLegacy handles the old command-line interface for backward compatibility.
// This ensures existing scripts and workflows continue to function while
// encouraging migration to the new command structure.
//
// Deprecated: Users should migrate to the new push/pull/schema commands.
//
// Parameters:
//   - cmd: The cobra command that triggered this function
//   - args: Command arguments
func runLegacy(cmd *cobra.Command, args []string) {
	// Check for legacy flags
	inputFile, _ := cmd.Flags().GetString("file")
	listIssueTypes, _ := cmd.Flags().GetBool("list-issue-types")
	checkFields, _ := cmd.Flags().GetString("check-fields")
	
	// Initialize Jira adapter for legacy commands
	jiraAdapter, err := jira.NewJiraAdapter()
	if err != nil {
		fmt.Printf("Error initializing Jira adapter: %v\n", err)
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

// printFieldInfo prints formatted field information for display.
// It formats JIRA field metadata in a human-readable format.
//
// Parameters:
//   - field: The JIRA field definition to display
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

// main is the entry point for the Ticketr CLI application.
// It handles both the modern command structure (push/pull/schema)
// and legacy command-line interface for backward compatibility.
func main() {
	// Check for legacy usage (no subcommand)
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		// If first arg is not a flag and not a known command, assume it's a file (legacy)
		knownCommands := []string{"push", "pull", "schema", "help", "completion"}
		isKnownCommand := false
		for _, cmd := range knownCommands {
			if os.Args[1] == cmd {
				isKnownCommand = true
				break
			}
		}
		
		if !isKnownCommand && !strings.HasPrefix(os.Args[1], "-") {
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