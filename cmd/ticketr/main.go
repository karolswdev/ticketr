package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/karolswdev/ticktr/internal/adapters/filesystem"
	"github.com/karolswdev/ticktr/internal/adapters/jira"
	"github.com/karolswdev/ticktr/internal/core/services"
)

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
	// Parse command-line arguments
	var inputFile string
	var forcePartialUpload bool
	var verbose bool
	var listIssueTypes bool
	var checkFields string
	
	flag.StringVar(&inputFile, "file", "", "Path to the input Markdown file")
	flag.StringVar(&inputFile, "f", "", "Path to the input Markdown file (shorthand)")
	flag.BoolVar(&forcePartialUpload, "force-partial-upload", false, "Continue processing even if some items fail")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging output")
	flag.BoolVar(&verbose, "v", false, "Enable verbose logging output (shorthand)")
	flag.BoolVar(&listIssueTypes, "list-issue-types", false, "List available issue types for the configured JIRA project")
	flag.StringVar(&checkFields, "check-fields", "", "Check required fields for a specific issue type")
	flag.Parse()
	
	// Configure logging based on verbose flag
	if verbose {
		log.SetFlags(log.Ltime | log.Lshortfile | log.Lmicroseconds)
		log.Println("Verbose mode enabled")
	} else {
		log.SetFlags(log.Ltime)
	}

	// Initialize adapters first (needed for both modes)
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

	// If list-issue-types flag is set, list issue types and exit
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
			} else {
				fmt.Println()
			}
		}
		fmt.Println("=" + string(make([]byte, 50)))
		
		fmt.Println("\nAvailable Issue Types:")
		if types, ok := issueTypesInfo["issueTypes"]; ok {
			for _, issueType := range types {
				fmt.Printf("  - %s\n", issueType)
			}
		}
		
		fmt.Println("\nTo use these issue types, set environment variables:")
		fmt.Println("  export JIRA_STORY_TYPE=\"<issue-type>\"")
		fmt.Println("  export JIRA_SUBTASK_TYPE=\"<subtask-type>\"")
		fmt.Println("\nNote: Subtask types are marked with '(subtask)'")
		fmt.Println("\nTo check required fields for an issue type:")
		fmt.Println("  ticketr --check-fields \"<issue-type>\"")
		os.Exit(0)
	}
	
	// If check-fields flag is set, check fields for the specified issue type
	if checkFields != "" {
		fmt.Printf("Fetching field requirements for issue type '%s'...\n", checkFields)
		fieldsInfo, err := jiraAdapter.GetIssueTypeFields(checkFields)
		if err != nil {
			fmt.Printf("Error fetching fields: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("\n" + "=" + string(make([]byte, 70)))
		fmt.Printf("Field Requirements for Issue Type: %s\n", checkFields)
		fmt.Println("=" + string(make([]byte, 70)))
		
		if fields, ok := fieldsInfo["fields"].([]map[string]interface{}); ok {
			// Separate required and optional fields
			var requiredFields []map[string]interface{}
			var optionalFields []map[string]interface{}
			
			for _, field := range fields {
				if required, ok := field["required"].(bool); ok && required {
					requiredFields = append(requiredFields, field)
				} else {
					optionalFields = append(optionalFields, field)
				}
			}
			
			// Print required fields
			if len(requiredFields) > 0 {
				fmt.Println("\nREQUIRED FIELDS:")
				fmt.Println("-" + string(make([]byte, 69)))
				for _, field := range requiredFields {
					printFieldInfo(field)
				}
			}
			
			// Print optional fields (limited list for brevity)
			if len(optionalFields) > 0 {
				fmt.Println("\nOPTIONAL FIELDS (showing common fields):")
				fmt.Println("-" + string(make([]byte, 69)))
				commonFields := []string{"description", "priority", "labels", "components", "fixVersions", "versions", "assignee", "reporter", "duedate", "timetracking", "environment"}
				for _, field := range optionalFields {
					if key, ok := field["key"].(string); ok {
						for _, common := range commonFields {
							if key == common {
								printFieldInfo(field)
								break
							}
						}
					}
				}
			}
		}
		
		fmt.Println("\n" + "=" + string(make([]byte, 70)))
		fmt.Println("\nNote: The fields 'summary' (Title) and 'description' are handled")
		fmt.Println("automatically by ticketr from your markdown files.")
		os.Exit(0)
	}

	// Check if input file was provided (only needed when not listing issue types)
	if inputFile == "" {
		fmt.Println("Error: Input file path is required")
		fmt.Println("\nUsage:")
		fmt.Println("  ticketr -file <path-to-markdown-file> [options]")
		fmt.Println("  ticketr -f <path-to-markdown-file> [options]")
		fmt.Println("  ticketr --list-issue-types")
		fmt.Println("  ticketr --check-fields <issue-type>")
		fmt.Println("\nOptions:")
		fmt.Println("  --force-partial-upload     Continue processing even if some items fail")
		fmt.Println("  --verbose, -v              Enable verbose logging output")
		fmt.Println("  --list-issue-types         List available issue types for the configured JIRA project")
		fmt.Println("  --check-fields <type>      Check required fields for a specific issue type")
		fmt.Println("\nExamples:")
		fmt.Println("  ticketr -f stories.md --verbose --force-partial-upload")
		fmt.Println("  ticketr --list-issue-types")
		fmt.Println("  ticketr --check-fields \"Task\"")
		fmt.Println("  ticketr --check-fields \"Sub-task\"")
		os.Exit(1)
	}

	// Check if file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' does not exist\n", inputFile)
		os.Exit(1)
	}

	// Initialize file repository
	fileRepo := filesystem.NewFileRepository()

	// Test Jira authentication
	fmt.Println("Authenticating with Jira...")
	if err := jiraAdapter.Authenticate(); err != nil {
		fmt.Printf("Error: Failed to authenticate with Jira: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully authenticated with Jira")

	// Initialize the service
	storyService := services.NewStoryService(fileRepo, jiraAdapter)

	// Process the stories
	fmt.Printf("\nProcessing stories from '%s'...\n", inputFile)
	if forcePartialUpload {
		fmt.Println("Force partial upload mode: Will continue on errors")
	}
	fmt.Println("=" + string(make([]byte, 50)))
	
	result, err := storyService.ProcessStoriesWithOptions(inputFile, services.ProcessOptions{
		ForcePartialUpload: forcePartialUpload,
	})
	if err != nil {
		fmt.Printf("Error processing stories: %v\n", err)
		os.Exit(1)
	}

	// Print summary report
	fmt.Println("\n" + "=" + string(make([]byte, 50)))
	fmt.Println("SUMMARY REPORT")
	fmt.Println("=" + string(make([]byte, 50)))
	fmt.Printf("Stories Created: %d\n", result.StoriesCreated)
	fmt.Printf("Stories Updated: %d\n", result.StoriesUpdated)
	fmt.Printf("Tasks Created:   %d\n", result.TasksCreated)
	fmt.Printf("Tasks Updated:   %d\n", result.TasksUpdated)
	
	if len(result.Errors) > 0 {
		fmt.Printf("\nErrors encountered: %d\n", len(result.Errors))
		for _, errMsg := range result.Errors {
			fmt.Printf("  - %s\n", errMsg)
		}
	} else {
		fmt.Println("\nAll operations completed successfully!")
	}

	// Set appropriate exit code
	if len(result.Errors) > 0 && !forcePartialUpload {
		os.Exit(2) // Partial success - exit with error unless force flag is set
	}
}

func init() {
	// Configure logging format
	log.SetFlags(log.Ltime | log.Lshortfile)
}