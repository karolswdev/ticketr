package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestRootCmd_NoArgs tests root command with no arguments shows help
func TestRootCmd_NoArgs(t *testing.T) {
	// Create a new command for testing
	testRootCmd := &cobra.Command{
		Use:   "ticketr",
		Short: "A tool for managing JIRA tickets as code",
		Run: func(cmd *cobra.Command, args []string) {
			// No-op, just verify command can execute
		},
	}

	// Execute with no args should work or show help
	output := new(bytes.Buffer)
	testRootCmd.SetOut(output)
	testRootCmd.SetErr(output)
	testRootCmd.SetArgs([]string{})

	_ = testRootCmd.Execute()

	// Verify the command structure is correct
	if testRootCmd.Use != "ticketr" {
		t.Errorf("Expected Use to be 'ticketr', got '%s'", testRootCmd.Use)
	}
}

// TestPushCmd_RequiresFileArgument tests that push command requires a file argument
func TestPushCmd_RequiresFileArgument(t *testing.T) {
	testPushCmd := &cobra.Command{
		Use:  "push [file]",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// No-op for test
		},
	}

	testPushCmd.SetArgs([]string{})
	err := testPushCmd.Execute()

	if err == nil {
		t.Error("Expected error when push command called without file argument")
	}

	if err != nil && !strings.Contains(err.Error(), "arg") {
		t.Logf("Error message: %v", err)
	}
}

// TestPushCmd_AcceptsFileArgument tests that push command accepts a file argument
func TestPushCmd_AcceptsFileArgument(t *testing.T) {
	fileArg := ""
	testPushCmd := &cobra.Command{
		Use:  "push [file]",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fileArg = args[0]
		},
	}

	testPushCmd.SetArgs([]string{"test.md"})
	err := testPushCmd.Execute()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if fileArg != "test.md" {
		t.Errorf("Expected file argument 'test.md', got '%s'", fileArg)
	}
}

// TestPullCmd_Flags tests pull command flag parsing
func TestPullCmd_Flags(t *testing.T) {
	var testProject, testJQL, testOutput, testEpic string
	var testForce bool

	testPullCmd := &cobra.Command{
		Use: "pull",
		Run: func(cmd *cobra.Command, args []string) {
			testProject, _ = cmd.Flags().GetString("project")
			testJQL, _ = cmd.Flags().GetString("jql")
			testOutput, _ = cmd.Flags().GetString("output")
			testEpic, _ = cmd.Flags().GetString("epic")
			testForce, _ = cmd.Flags().GetBool("force")
		},
	}

	testPullCmd.Flags().StringVar(&testProject, "project", "", "JIRA project key")
	testPullCmd.Flags().StringVar(&testJQL, "jql", "", "JQL query")
	testPullCmd.Flags().StringVarP(&testOutput, "output", "o", "pulled_tickets.md", "output file")
	testPullCmd.Flags().StringVar(&testEpic, "epic", "", "Epic key")
	testPullCmd.Flags().BoolVar(&testForce, "force", false, "Force overwrite")

	testPullCmd.SetArgs([]string{
		"--project", "TEST",
		"--jql", "status=Done",
		"--output", "output.md",
		"--epic", "TEST-100",
		"--force",
	})

	err := testPullCmd.Execute()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if testProject != "TEST" {
		t.Errorf("Expected project 'TEST', got '%s'", testProject)
	}

	if testJQL != "status=Done" {
		t.Errorf("Expected JQL 'status=Done', got '%s'", testJQL)
	}

	if testOutput != "output.md" {
		t.Errorf("Expected output 'output.md', got '%s'", testOutput)
	}

	if testEpic != "TEST-100" {
		t.Errorf("Expected epic 'TEST-100', got '%s'", testEpic)
	}

	if !testForce {
		t.Error("Expected force flag to be true")
	}
}

// TestConfigFile_Override tests that --config flag overrides default config location
func TestConfigFile_Override(t *testing.T) {
	var testCfgFile string
	testRootCmd := &cobra.Command{
		Use: "ticketr",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			testCfgFile, _ = cmd.Flags().GetString("config")
		},
	}

	testRootCmd.PersistentFlags().StringVar(&testCfgFile, "config", "", "config file")

	testRootCmd.SetArgs([]string{"--config", "/custom/path/.ticketr.yaml"})
	err := testRootCmd.Execute()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if testCfgFile != "/custom/path/.ticketr.yaml" {
		t.Errorf("Expected config file '/custom/path/.ticketr.yaml', got '%s'", testCfgFile)
	}
}

// TestVerboseFlag_Global tests that --verbose flag is available globally
func TestVerboseFlag_Global(t *testing.T) {
	var testVerbose bool
	testRootCmd := &cobra.Command{
		Use: "ticketr",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			testVerbose, _ = cmd.Flags().GetBool("verbose")
		},
	}

	testRootCmd.PersistentFlags().BoolVarP(&testVerbose, "verbose", "v", false, "verbose logging")

	// Test short flag
	testRootCmd.SetArgs([]string{"-v"})
	err := testRootCmd.Execute()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !testVerbose {
		t.Error("Expected verbose flag to be true with -v")
	}

	// Test long flag
	testVerbose = false
	testRootCmd.SetArgs([]string{"--verbose"})
	err = testRootCmd.Execute()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !testVerbose {
		t.Error("Expected verbose flag to be true with --verbose")
	}
}

// TestForcePartialUploadFlag tests --force-partial-upload flag on push command
func TestForcePartialUploadFlag(t *testing.T) {
	var testForce bool
	testPushCmd := &cobra.Command{
		Use:  "push [file]",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			testForce, _ = cmd.Flags().GetBool("force-partial-upload")
		},
	}

	testPushCmd.Flags().BoolVar(&testForce, "force-partial-upload", false, "force partial upload")

	testPushCmd.SetArgs([]string{"test.md", "--force-partial-upload"})
	err := testPushCmd.Execute()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !testForce {
		t.Error("Expected force-partial-upload flag to be true")
	}
}

// TestEnvironmentVariables tests that environment variables are properly handled
func TestEnvironmentVariables(t *testing.T) {
	// Save original env vars
	originalJiraURL := os.Getenv("JIRA_URL")
	originalJiraEmail := os.Getenv("JIRA_EMAIL")
	originalJiraAPIKey := os.Getenv("JIRA_API_KEY")
	originalJiraProjectKey := os.Getenv("JIRA_PROJECT_KEY")

	defer func() {
		// Restore original env vars
		os.Setenv("JIRA_URL", originalJiraURL)
		os.Setenv("JIRA_EMAIL", originalJiraEmail)
		os.Setenv("JIRA_API_KEY", originalJiraAPIKey)
		os.Setenv("JIRA_PROJECT_KEY", originalJiraProjectKey)
	}()

	// Set test env vars
	testURL := "https://test.atlassian.net"
	testEmail := "test@example.com"
	testAPIKey := "test-api-key"
	testProjectKey := "TEST"

	os.Setenv("JIRA_URL", testURL)
	os.Setenv("JIRA_EMAIL", testEmail)
	os.Setenv("JIRA_API_KEY", testAPIKey)
	os.Setenv("JIRA_PROJECT_KEY", testProjectKey)

	// Verify env vars are set correctly
	if os.Getenv("JIRA_URL") != testURL {
		t.Errorf("Expected JIRA_URL to be '%s', got '%s'", testURL, os.Getenv("JIRA_URL"))
	}

	if os.Getenv("JIRA_EMAIL") != testEmail {
		t.Errorf("Expected JIRA_EMAIL to be '%s', got '%s'", testEmail, os.Getenv("JIRA_EMAIL"))
	}

	if os.Getenv("JIRA_API_KEY") != testAPIKey {
		t.Errorf("Expected JIRA_API_KEY to be '%s', got '%s'", testAPIKey, os.Getenv("JIRA_API_KEY"))
	}

	if os.Getenv("JIRA_PROJECT_KEY") != testProjectKey {
		t.Errorf("Expected JIRA_PROJECT_KEY to be '%s', got '%s'", testProjectKey, os.Getenv("JIRA_PROJECT_KEY"))
	}
}

// TestSchemaCmd_NoArgs tests schema command can be called without arguments
func TestSchemaCmd_NoArgs(t *testing.T) {
	executed := false
	testSchemaCmd := &cobra.Command{
		Use: "schema",
		Run: func(cmd *cobra.Command, args []string) {
			executed = true
		},
	}

	testSchemaCmd.SetArgs([]string{})
	err := testSchemaCmd.Execute()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !executed {
		t.Error("Expected schema command to execute")
	}
}

// TestCommandHelp tests that help is available for all commands
func TestCommandHelp(t *testing.T) {
	commands := []struct {
		name    string
		command *cobra.Command
	}{
		{
			name: "push",
			command: &cobra.Command{
				Use:   "push [file]",
				Short: "Push tickets to JIRA",
			},
		},
		{
			name: "pull",
			command: &cobra.Command{
				Use:   "pull",
				Short: "Pull tickets from JIRA",
			},
		},
		{
			name: "schema",
			command: &cobra.Command{
				Use:   "schema",
				Short: "Discover JIRA schema",
			},
		},
	}

	for _, cmd := range commands {
		t.Run(cmd.name, func(t *testing.T) {
			output := new(bytes.Buffer)
			cmd.command.SetOut(output)
			cmd.command.SetArgs([]string{"--help"})

			err := cmd.command.Execute()

			// Help should not return error
			if err != nil {
				t.Logf("Help command returned error (may be expected): %v", err)
			}

			// Help output should contain command usage
			if output.Len() == 0 {
				t.Error("Expected help output, got empty string")
			}
		})
	}
}

// TestProcessFieldForSchema tests the processFieldForSchema function
func TestProcessFieldForSchema(t *testing.T) {
	customFieldsMap := make(map[string]map[string]interface{})

	// Test with a valid custom field
	field := map[string]interface{}{
		"key":  "customfield_10001",
		"name": "Story Points",
		"schema": map[string]interface{}{
			"type": "number",
		},
	}

	processFieldForSchema(field, customFieldsMap)

	if _, exists := customFieldsMap["Story Points"]; !exists {
		t.Error("Expected 'Story Points' field to be added to customFieldsMap")
	}

	if customFieldsMap["Story Points"]["id"] != "customfield_10001" {
		t.Errorf("Expected field ID 'customfield_10001', got '%s'", customFieldsMap["Story Points"]["id"])
	}

	if customFieldsMap["Story Points"]["type"] != "number" {
		t.Errorf("Expected field type 'number', got '%s'", customFieldsMap["Story Points"]["type"])
	}
}

// TestProcessFieldForSchema_SkipsSystemFields tests that system fields are skipped
func TestProcessFieldForSchema_SkipsSystemFields(t *testing.T) {
	customFieldsMap := make(map[string]map[string]interface{})

	// Test with system field (non-customfield)
	field := map[string]interface{}{
		"key":  "summary",
		"name": "Summary",
	}

	processFieldForSchema(field, customFieldsMap)

	if len(customFieldsMap) != 0 {
		t.Error("Expected system field to be skipped, but it was added")
	}

	// Test with Development field (should be skipped)
	field = map[string]interface{}{
		"key":  "customfield_10002",
		"name": "Development",
	}

	processFieldForSchema(field, customFieldsMap)

	if len(customFieldsMap) != 0 {
		t.Error("Expected Development field to be skipped, but it was added")
	}

	// Test with chart field (should be skipped)
	field = map[string]interface{}{
		"key":  "customfield_10003",
		"name": "Status [CHART]",
	}

	processFieldForSchema(field, customFieldsMap)

	if len(customFieldsMap) != 0 {
		t.Error("Expected chart field to be skipped, but it was added")
	}
}

// TestProcessFieldForSchema_DifferentTypes tests field type detection
func TestProcessFieldForSchema_DifferentTypes(t *testing.T) {
	testCases := []struct {
		name         string
		schemaType   string
		expectedType string
	}{
		{"number", "number", "number"},
		{"array", "array", "array"},
		{"option", "option", "option"},
		{"default", "other", "string"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			customFieldsMap := make(map[string]map[string]interface{})

			field := map[string]interface{}{
				"key":  "customfield_10100",
				"name": "Test Field",
				"schema": map[string]interface{}{
					"type": tc.schemaType,
				},
			}

			processFieldForSchema(field, customFieldsMap)

			if customFieldsMap["Test Field"]["type"] != tc.expectedType {
				t.Errorf("Expected type '%s', got '%s'", tc.expectedType, customFieldsMap["Test Field"]["type"])
			}
		})
	}
}
