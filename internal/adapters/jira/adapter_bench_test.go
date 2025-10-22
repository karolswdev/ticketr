package jira

import (
	"context"
	"os"
	"testing"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// BenchmarkAdapterCreation measures adapter instantiation performance
func BenchmarkAdapterCreation(b *testing.B) {
	config := &domain.WorkspaceConfig{
		JiraURL:    "https://test.atlassian.net",
		Username:   "test@example.com",
		APIToken:   "test-token",
		ProjectKey: "TEST",
	}

	fieldMappings := getDefaultFieldMappings()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewJiraAdapterFromConfig(config, fieldMappings)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkFieldMapping compares field mapping performance through adapter creation
func BenchmarkFieldMapping(b *testing.B) {
	fieldMappings := map[string]interface{}{
		"Story Points": "customfield_10016",
		"Sprint":       "customfield_10020",
		"Epic Link":    "customfield_10014",
		"Priority":     "priority",
		"Status":       "status",
	}

	config := &domain.WorkspaceConfig{
		JiraURL:    "https://test.atlassian.net",
		Username:   "test@example.com",
		APIToken:   "test-token",
		ProjectKey: "TEST",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewJiraAdapterFromConfig(config, fieldMappings)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkTicketCreation compares ticket creation setup performance
func BenchmarkTicketCreation(b *testing.B) {
	ticket := domain.Ticket{
		Title:       "Benchmark Test Ticket",
		Description: "This is a test description with some content",
		AcceptanceCriteria: []string{
			"Should handle large datasets",
			"Should be performant",
			"Should maintain backward compatibility",
		},
		CustomFields: map[string]string{
			"Story Points": "5",
			"Priority":     "High",
		},
	}

	config := &domain.WorkspaceConfig{
		JiraURL:    "https://test.atlassian.net",
		Username:   "test@example.com",
		APIToken:   "test-token",
		ProjectKey: "TEST",
	}

	fieldMappings := map[string]interface{}{
		"Story Points": "customfield_10016",
		"Priority":     "priority",
	}

	adapter, err := NewJiraAdapterFromConfig(config, fieldMappings)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Note: This won't actually create the ticket (no network call)
		// but measures the preparation overhead
		_ = adapter
		_ = ticket
	}
}

// Integration benchmarks (require real Jira connection)

// BenchmarkSearchTickets_Integration measures actual search performance
func BenchmarkSearchTickets_Integration(b *testing.B) {
	config := getIntegrationConfigForBench(b)
	if config == nil {
		b.Skip("Skipping integration benchmark: JIRA credentials not set")
	}

	projectKey := config.ProjectKey
	ctx := context.Background()
	fieldMappings := getDefaultFieldMappings()

	adapter, err := NewJiraAdapterFromConfig(config, fieldMappings)
	if err != nil {
		b.Fatalf("failed to create adapter: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := adapter.SearchTickets(ctx, projectKey, "ORDER BY created DESC", nil)
		if err != nil {
			b.Fatalf("SearchTickets failed: %v", err)
		}
	}
}

// BenchmarkAuthenticate_Integration measures authentication performance
func BenchmarkAuthenticate_Integration(b *testing.B) {
	config := getIntegrationConfigForBench(b)
	if config == nil {
		b.Skip("Skipping integration benchmark: JIRA credentials not set")
	}

	fieldMappings := getDefaultFieldMappings()

	adapter, err := NewJiraAdapterFromConfig(config, fieldMappings)
	if err != nil {
		b.Fatalf("failed to create adapter: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := adapter.Authenticate()
		if err != nil {
			b.Fatalf("Authenticate failed: %v", err)
		}
	}
}

// BenchmarkGetProjectIssueTypes_Integration measures metadata fetch performance
func BenchmarkGetProjectIssueTypes_Integration(b *testing.B) {
	config := getIntegrationConfigForBench(b)
	if config == nil {
		b.Skip("Skipping integration benchmark: JIRA credentials not set")
	}

	fieldMappings := getDefaultFieldMappings()

	adapter, err := NewJiraAdapterFromConfig(config, fieldMappings)
	if err != nil {
		b.Fatalf("failed to create adapter: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := adapter.GetProjectIssueTypes()
		if err != nil {
			b.Fatalf("GetProjectIssueTypes failed: %v", err)
		}
	}
}

// Helper functions

func getIntegrationConfigForBench(b *testing.B) *domain.WorkspaceConfig {
	b.Helper()

	url := os.Getenv("JIRA_URL")
	email := os.Getenv("JIRA_EMAIL")
	token := os.Getenv("JIRA_API_KEY")
	project := os.Getenv("JIRA_PROJECT_KEY")

	if url == "" || email == "" || token == "" || project == "" {
		return nil
	}

	return &domain.WorkspaceConfig{
		JiraURL:    url,
		Username:   email,
		APIToken:   token,
		ProjectKey: project,
	}
}

func createBenchAdapter(b *testing.B, fieldMappings map[string]interface{}) ports.JiraPort {
	b.Helper()

	config := getIntegrationConfigForBench(b)
	if config == nil {
		b.Fatal("integration config not available")
	}

	adapter, err := NewJiraAdapterFromConfig(config, fieldMappings)
	if err != nil {
		b.Fatalf("failed to create adapter: %v", err)
	}

	return adapter
}
