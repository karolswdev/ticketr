package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/karolswdev/ticktr/internal/adapters/filesystem"
	"github.com/karolswdev/ticktr/internal/adapters/jira"
	"github.com/karolswdev/ticktr/internal/core/services"
)

func main() {
	// Parse command-line arguments
	var inputFile string
	var forcePartialUpload bool
	var verbose bool
	
	flag.StringVar(&inputFile, "file", "", "Path to the input Markdown file")
	flag.StringVar(&inputFile, "f", "", "Path to the input Markdown file (shorthand)")
	flag.BoolVar(&forcePartialUpload, "force-partial-upload", false, "Continue processing even if some items fail")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging output")
	flag.BoolVar(&verbose, "v", false, "Enable verbose logging output (shorthand)")
	flag.Parse()
	
	// Configure logging based on verbose flag
	if verbose {
		log.SetFlags(log.Ltime | log.Lshortfile | log.Lmicroseconds)
		log.Println("Verbose mode enabled")
	} else {
		log.SetFlags(log.Ltime)
	}

	// Check if input file was provided
	if inputFile == "" {
		fmt.Println("Error: Input file path is required")
		fmt.Println("\nUsage:")
		fmt.Println("  ticketr -file <path-to-markdown-file> [options]")
		fmt.Println("  ticketr -f <path-to-markdown-file> [options]")
		fmt.Println("\nOptions:")
		fmt.Println("  --force-partial-upload    Continue processing even if some items fail")
		fmt.Println("  --verbose, -v            Enable verbose logging output")
		fmt.Println("\nExample:")
		fmt.Println("  ticketr -f stories.md --verbose --force-partial-upload")
		os.Exit(1)
	}

	// Check if file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' does not exist\n", inputFile)
		os.Exit(1)
	}

	// Initialize adapters
	fileRepo := filesystem.NewFileRepository()
	
	jiraAdapter, err := jira.NewJiraAdapter()
	if err != nil {
		fmt.Printf("Error initializing Jira adapter: %v\n", err)
		fmt.Println("\nMake sure the following environment variables are set:")
		fmt.Println("  - JIRA_URL")
		fmt.Println("  - JIRA_EMAIL")
		fmt.Println("  - JIRA_API_KEY")
		fmt.Println("  - JIRA_PROJECT_KEY")
		os.Exit(1)
	}

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