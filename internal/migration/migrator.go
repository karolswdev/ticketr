package migration

import (
	"fmt"
	"os"
	"strings"
)

// Migrator handles conversion from legacy # STORY: to # TICKET: format
type Migrator struct {
	DryRun bool
}

// MigrateFile converts a single file from STORY to TICKET format
// Returns: (modified content, changed bool, error)
func (m *Migrator) MigrateFile(filepath string) (string, bool, error) {
	// Read file content
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", false, fmt.Errorf("failed to read file: %w", err)
	}

	originalContent := string(content)

	// Use strings.ReplaceAll to replace # STORY: with # TICKET:
	modifiedContent := strings.ReplaceAll(originalContent, "# STORY:", "# TICKET:")

	// Check if any changes were made
	changed := modifiedContent != originalContent

	return modifiedContent, changed, nil
}

// PreviewDiff shows what would change (for dry-run mode)
func (m *Migrator) PreviewDiff(filepath string, oldContent, newContent string) string {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("Preview of changes for: %s\n", filepath))

	// Split into lines and show changes
	oldLines := strings.Split(oldContent, "\n")
	newLines := strings.Split(newContent, "\n")

	changeCount := 0
	for i := 0; i < len(oldLines) && i < len(newLines); i++ {
		if oldLines[i] != newLines[i] {
			changeCount++
			result.WriteString(fmt.Sprintf("  Line %d: %s -> %s\n", i+1, oldLines[i], newLines[i]))
		}
	}

	result.WriteString(fmt.Sprintf("\n%d change(s) would be made. Use --write to apply.\n", changeCount))

	return result.String()
}

// WriteMigration writes the migrated content to file
func (m *Migrator) WriteMigration(filepath string, content string) error {
	// Write content to file with proper permissions (0644 = rw-r--r--)
	err := os.WriteFile(filepath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
