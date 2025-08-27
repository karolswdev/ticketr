package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/karolswdev/ticktr/internal/adapters/filesystem"
	"github.com/karolswdev/ticktr/internal/adapters/jira"
	"github.com/karolswdev/ticktr/internal/core/services"
)

var (
	cfgFile string
	verbose bool
	forcePartialUpload bool
	
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
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Pull command not yet implemented in Phase 1")
		},
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

func init() {
	cobra.OnInitialize(initConfig)
	
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .ticketr.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")
	
	// Push command flags
	pushCmd.Flags().BoolVar(&forcePartialUpload, "force-partial-upload", false, "continue processing even if some items fail")
	
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
	
	// Initialize repository
	repo := filesystem.NewFileRepository()
	
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
	
	// Process stories
	options := services.ProcessOptions{
		ForcePartialUpload: forcePartialUpload,
	}
	
	result, err := service.ProcessStoriesWithOptions(inputFile, options)
	if err != nil {
		fmt.Printf("Error processing file: %v\n", err)
		os.Exit(1)
	}
	
	// Print summary
	fmt.Println("\n=== Summary ===")
	if result.StoriesCreated > 0 || result.TicketsCreated > 0 {
		fmt.Printf("Stories created: %d\n", result.StoriesCreated)
	}
	if result.StoriesUpdated > 0 || result.TicketsUpdated > 0 {
		fmt.Printf("Stories updated: %d\n", result.StoriesUpdated)
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

// runSchema handles the schema discovery command
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

// runLegacy handles the old command-line interface for backward compatibility
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