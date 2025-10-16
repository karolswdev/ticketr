package migration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrator_ConvertsSingleStory(t *testing.T) {
	migrator := &Migrator{DryRun: true}

	// Use testdata/legacy_story/simple_story.md
	testFile := "../../testdata/legacy_story/simple_story.md"

	content, changed, err := migrator.MigrateFile(testFile)
	if err != nil {
		t.Fatalf("MigrateFile failed: %v", err)
	}

	// Verify "# STORY:" becomes "# TICKET:"
	if !strings.Contains(content, "# TICKET:") {
		t.Errorf("Expected migrated content to contain '# TICKET:', got: %s", content)
	}

	if strings.Contains(content, "# STORY:") {
		t.Errorf("Expected migrated content to not contain '# STORY:', but it does")
	}

	// Verify other content unchanged
	if !strings.Contains(content, "USER-123 Simple user story") {
		t.Errorf("Expected content to preserve ticket title")
	}

	if !strings.Contains(content, "As a user, I want to see my profile.") {
		t.Errorf("Expected content to preserve description")
	}

	// Verify changed flag is true
	if !changed {
		t.Errorf("Expected changed=true, got false")
	}
}

func TestMigrator_ConvertsMultipleStories(t *testing.T) {
	migrator := &Migrator{DryRun: true}

	// Use testdata/legacy_story/multiple_stories.md
	testFile := "../../testdata/legacy_story/multiple_stories.md"

	content, changed, err := migrator.MigrateFile(testFile)
	if err != nil {
		t.Fatalf("MigrateFile failed: %v", err)
	}

	// Verify all occurrences converted
	storyCount := strings.Count(content, "# STORY:")
	if storyCount != 0 {
		t.Errorf("Expected 0 occurrences of '# STORY:', found %d", storyCount)
	}

	ticketCount := strings.Count(content, "# TICKET:")
	if ticketCount != 2 {
		t.Errorf("Expected 2 occurrences of '# TICKET:', found %d", ticketCount)
	}

	if !changed {
		t.Errorf("Expected changed=true, got false")
	}
}

func TestMigrator_PreservesFormatting(t *testing.T) {
	migrator := &Migrator{DryRun: true}

	testFile := "../../testdata/legacy_story/simple_story.md"

	// Read original content
	originalBytes, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read original file: %v", err)
	}
	original := string(originalBytes)

	content, _, err := migrator.MigrateFile(testFile)
	if err != nil {
		t.Fatalf("MigrateFile failed: %v", err)
	}

	// Verify whitespace, indentation preserved
	originalLines := strings.Split(original, "\n")
	contentLines := strings.Split(content, "\n")

	if len(originalLines) != len(contentLines) {
		t.Errorf("Expected same number of lines: original=%d, migrated=%d", len(originalLines), len(contentLines))
	}

	// Verify only header line changed
	changedLines := 0
	for i := 0; i < len(originalLines) && i < len(contentLines); i++ {
		if originalLines[i] != contentLines[i] {
			changedLines++
			// The changed line should be the conversion from STORY to TICKET
			if !strings.Contains(originalLines[i], "# STORY:") || !strings.Contains(contentLines[i], "# TICKET:") {
				t.Errorf("Unexpected line change at line %d: '%s' -> '%s'", i+1, originalLines[i], contentLines[i])
			}
		}
	}

	if changedLines != 1 {
		t.Errorf("Expected exactly 1 line to change, got %d", changedLines)
	}
}

func TestMigrator_HandlesNoChangesNeeded(t *testing.T) {
	migrator := &Migrator{DryRun: true}

	// Test with file already using # TICKET:
	testFile := "../../testdata/ticket_simple.md"

	content, changed, err := migrator.MigrateFile(testFile)
	if err != nil {
		t.Fatalf("MigrateFile failed: %v", err)
	}

	// Verify changed=false returned
	if changed {
		t.Errorf("Expected changed=false for file already using # TICKET:, got true")
	}

	// Verify content unchanged
	originalBytes, _ := os.ReadFile(testFile)
	original := string(originalBytes)

	if content != original {
		t.Errorf("Expected content to remain unchanged")
	}
}

func TestMigrator_HandlesFileErrors(t *testing.T) {
	migrator := &Migrator{DryRun: true}

	// Test non-existent file
	_, _, err := migrator.MigrateFile("/nonexistent/path/file.md")
	if err == nil {
		t.Errorf("Expected error for non-existent file, got nil")
	}

	// Verify appropriate error returned
	if !strings.Contains(err.Error(), "failed to read file") {
		t.Errorf("Expected error message to contain 'failed to read file', got: %s", err.Error())
	}
}

func TestMigrator_WriteMigration(t *testing.T) {
	migrator := &Migrator{DryRun: false}

	// Create a temporary file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_migration.md")

	// Write initial content with # STORY:
	initialContent := "# STORY: TEST-001 Test story\n\n## Description\nTest content\n"
	err := os.WriteFile(testFile, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Perform migration
	content, changed, err := migrator.MigrateFile(testFile)
	if err != nil {
		t.Fatalf("MigrateFile failed: %v", err)
	}

	if !changed {
		t.Fatalf("Expected changes, got none")
	}

	// Write the migration
	err = migrator.WriteMigration(testFile, content)
	if err != nil {
		t.Fatalf("WriteMigration failed: %v", err)
	}

	// Verify file was written correctly
	writtenBytes, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	written := string(writtenBytes)
	if !strings.Contains(written, "# TICKET:") {
		t.Errorf("Expected written file to contain '# TICKET:'")
	}

	if strings.Contains(written, "# STORY:") {
		t.Errorf("Expected written file to not contain '# STORY:'")
	}
}

func TestMigrator_PreviewDiff(t *testing.T) {
	migrator := &Migrator{DryRun: true}

	oldContent := "# STORY: TEST-001 Test\n\nContent here\n"
	newContent := "# TICKET: TEST-001 Test\n\nContent here\n"

	preview := migrator.PreviewDiff("test.md", oldContent, newContent)

	// Verify preview contains filename
	if !strings.Contains(preview, "test.md") {
		t.Errorf("Expected preview to contain filename")
	}

	// Verify preview shows line numbers
	if !strings.Contains(preview, "Line 1") {
		t.Errorf("Expected preview to show line numbers")
	}

	// Verify preview shows the change
	if !strings.Contains(preview, "# STORY:") || !strings.Contains(preview, "# TICKET:") {
		t.Errorf("Expected preview to show both old and new content")
	}

	// Verify preview mentions --write flag
	if !strings.Contains(preview, "--write") {
		t.Errorf("Expected preview to mention --write flag")
	}
}
