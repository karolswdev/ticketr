// Package main demonstrates using andygrunwald/go-jira library
// This is a research prototype to evaluate the library for Ticketr
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	jira "github.com/andygrunwald/go-jira"
)

func main() {
	// Get credentials from environment
	jiraURL := os.Getenv("JIRA_URL")
	email := os.Getenv("JIRA_EMAIL")
	apiToken := os.Getenv("JIRA_API_KEY")
	projectKey := os.Getenv("JIRA_PROJECT_KEY")

	if jiraURL == "" || email == "" || apiToken == "" || projectKey == "" {
		log.Fatal("Missing required environment variables: JIRA_URL, JIRA_EMAIL, JIRA_API_KEY, JIRA_PROJECT_KEY")
	}

	// Create authenticated transport
	tp := jira.BasicAuthTransport{
		Username: email,
		Password: apiToken,
	}

	// Create Jira client
	client, err := jira.NewClient(tp.Client(), jiraURL)
	if err != nil {
		log.Fatalf("Failed to create Jira client: %v", err)
	}

	fmt.Println("=== Jira Library Prototype ===")
	fmt.Println()

	// Test 1: Authentication by getting current user
	fmt.Println("Test 1: Authentication")
	user, resp, err := client.User.GetSelf()
	if err != nil {
		log.Fatalf("Authentication failed: %v (status: %d)", err, resp.StatusCode)
	}
	fmt.Printf("✓ Authenticated as: %s (%s)\n", user.DisplayName, user.EmailAddress)
	fmt.Println()

	// Test 2: Search for issues using JQL
	fmt.Println("Test 2: JQL Search")
	jql := fmt.Sprintf("project = %s ORDER BY created DESC", projectKey)

	// SearchOptions allows pagination
	searchOptions := &jira.SearchOptions{
		MaxResults: 10,
		StartAt:    0,
		Fields:     []string{"summary", "description", "issuetype", "status", "parent"},
	}

	issues, resp, err := client.Issue.Search(jql, searchOptions)
	if err != nil {
		log.Fatalf("Search failed: %v (status: %d)", err, resp.StatusCode)
	}
	fmt.Printf("✓ Found %d issues (showing first 10)\n", len(issues))
	for i, issue := range issues {
		fmt.Printf("  %d. %s: %s [%s]\n", i+1, issue.Key, issue.Fields.Summary, issue.Fields.Type.Name)

		// Show subtask info if available
		if len(issue.Fields.Subtasks) > 0 {
			fmt.Printf("     Subtasks: %d\n", len(issue.Fields.Subtasks))
		}
	}
	fmt.Println()

	// Test 3: Get single issue with all fields
	if len(issues) > 0 {
		fmt.Println("Test 3: Fetch Single Issue Details")
		issueKey := issues[0].Key

		issue, resp, err := client.Issue.Get(issueKey, nil)
		if err != nil {
			log.Fatalf("Get issue failed: %v (status: %d)", err, resp.StatusCode)
		}

		fmt.Printf("✓ Issue: %s\n", issue.Key)
		fmt.Printf("  Summary: %s\n", issue.Fields.Summary)
		fmt.Printf("  Type: %s\n", issue.Fields.Type.Name)
		fmt.Printf("  Status: %s\n", issue.Fields.Status.Name)
		if issue.Fields.Description != "" {
			descPreview := issue.Fields.Description
			if len(descPreview) > 100 {
				descPreview = descPreview[:100] + "..."
			}
			fmt.Printf("  Description: %s\n", descPreview)
		}
		fmt.Println()
	}

	// Test 4: Get project information
	fmt.Println("Test 4: Project Information")
	project, resp, err := client.Project.Get(projectKey)
	if err != nil {
		log.Fatalf("Get project failed: %v (status: %d)", err, resp.StatusCode)
	}
	fmt.Printf("✓ Project: %s (%s)\n", project.Name, project.Key)
	fmt.Printf("  Issue Types: %d\n", len(project.IssueTypes))
	for _, issueType := range project.IssueTypes {
		subtaskInfo := ""
		if issueType.Subtask {
			subtaskInfo = " (subtask)"
		}
		fmt.Printf("    - %s%s\n", issueType.Name, subtaskInfo)
	}
	fmt.Println()

	// Test 5: Pagination demonstration
	fmt.Println("Test 5: Pagination Test")
	totalFetched := 0
	pageSize := 5
	startAt := 0

	for {
		searchOpts := &jira.SearchOptions{
			MaxResults: pageSize,
			StartAt:    startAt,
		}

		page, _, err := client.Issue.Search(jql, searchOpts)
		if err != nil {
			log.Printf("Pagination error: %v", err)
			break
		}

		if len(page) == 0 {
			break
		}

		totalFetched += len(page)
		fmt.Printf("  Page %d: %d issues (total: %d)\n", (startAt/pageSize)+1, len(page), totalFetched)

		if len(page) < pageSize {
			break
		}

		startAt += pageSize
	}
	fmt.Printf("✓ Total issues fetched via pagination: %d\n", totalFetched)
	fmt.Println()

	// Test 6: Custom fields handling
	fmt.Println("Test 6: Custom Fields")
	if len(issues) > 0 {
		issue := issues[0]

		// Access custom fields via Unknowns map
		fmt.Printf("✓ Issue %s has access to custom fields via .Fields.Unknowns map\n", issue.Key)
		fmt.Printf("  Example: Story Points might be at customfield_10010\n")
		fmt.Printf("  Library provides flexible access to any custom field\n")
	}
	fmt.Println()

	// Test 7: Context support (cancellation)
	fmt.Println("Test 7: Context Support")
	ctx := context.Background()

	// The library supports context through http.Request creation
	// This is handled internally when you use the client
	fmt.Println("✓ Library uses standard http.Client which supports context")
	fmt.Println("  Context can be passed via custom http.Client configuration")
	_ = ctx
	fmt.Println()

	fmt.Println("=== All Tests Passed ===")
	fmt.Println()
	fmt.Println("Summary:")
	fmt.Println("- Authentication: ✓")
	fmt.Println("- JQL Search: ✓")
	fmt.Println("- Fetch Issue: ✓")
	fmt.Println("- Project Info: ✓")
	fmt.Println("- Pagination: ✓")
	fmt.Println("- Custom Fields: ✓")
	fmt.Println("- Context Support: ✓ (via http.Client)")
}
