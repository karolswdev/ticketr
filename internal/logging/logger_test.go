package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileLogger_CreatesLogFile(t *testing.T) {
	tmpDir := t.TempDir()

	config := LogConfig{
		LogDir:  tmpDir,
		Verbose: false,
	}

	logger, err := NewFileLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// Verify log file was created
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read log directory: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 log file, found %d", len(entries))
	}

	if !strings.HasSuffix(entries[0].Name(), ".log") {
		t.Errorf("Expected .log extension, got %s", entries[0].Name())
	}
}

func TestFileLogger_WritesMessages(t *testing.T) {
	tmpDir := t.TempDir()

	config := LogConfig{
		LogDir:  tmpDir,
		Verbose: false,
	}

	logger, err := NewFileLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Write various message types
	logger.Info("Test info message")
	logger.Warn("Test warning message")
	logger.Error("Test error message")
	logger.Section("Test Section")

	logger.Close()

	// Read log file content
	entries, _ := os.ReadDir(tmpDir)
	logPath := filepath.Join(tmpDir, entries[0].Name())
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify messages were written
	if !strings.Contains(logContent, "Test info message") {
		t.Error("Info message not found in log")
	}
	if !strings.Contains(logContent, "Test warning message") {
		t.Error("Warning message not found in log")
	}
	if !strings.Contains(logContent, "Test error message") {
		t.Error("Error message not found in log")
	}
	if !strings.Contains(logContent, "Test Section") {
		t.Error("Section not found in log")
	}
}

func TestSensitiveRedactor_RedactsAPIKeys(t *testing.T) {
	redactor := NewSensitiveRedactor()

	testCases := []struct {
		input         string
		shouldContain string
	}{
		{"API_KEY=secret123", "[REDACTED]"},
		{"token: abc123xyz", "[REDACTED]"},
		{"user@example.com", "[REDACTED]"},
		{"https://user:pass@example.com", "[REDACTED]"},
	}

	for _, tc := range testCases {
		result := redactor.Redact(tc.input)
		if !strings.Contains(result, tc.shouldContain) {
			t.Errorf("Expected redaction in %q, got %q", tc.input, result)
		}
		if strings.Contains(result, "secret") || strings.Contains(result, "pass") {
			t.Errorf("Sensitive data not redacted: %q", result)
		}
	}
}

func TestCleanup_RemovesOldLogs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create 5 dummy log files
	for i := 0; i < 5; i++ {
		logPath := filepath.Join(tmpDir, fmt.Sprintf("2024-01-0%d_10-00-00.log", i+1))
		if err := os.WriteFile(logPath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test log: %v", err)
		}
	}

	// Keep only 3 most recent
	err := Cleanup(tmpDir, 3)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	entries, _ := os.ReadDir(tmpDir)
	if len(entries) != 3 {
		t.Errorf("Expected 3 log files remaining, found %d", len(entries))
	}
}

func TestFileLogger_VerboseMode(t *testing.T) {
	tmpDir := t.TempDir()

	config := LogConfig{
		LogDir:  tmpDir,
		Verbose: true,
	}

	logger, err := NewFileLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	logger.Info("Verbose test message")
	logger.Close()

	// Verify the message was written to file
	entries, _ := os.ReadDir(tmpDir)
	logPath := filepath.Join(tmpDir, entries[0].Name())
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(content), "Verbose test message") {
		t.Error("Verbose message not found in log file")
	}
}

func TestFileLogger_Header(t *testing.T) {
	tmpDir := t.TempDir()

	config := LogConfig{
		LogDir:  tmpDir,
		Verbose: false,
	}

	logger, err := NewFileLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	logger.Close()

	// Read log file content
	entries, _ := os.ReadDir(tmpDir)
	logPath := filepath.Join(tmpDir, entries[0].Name())
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify header was written
	if !strings.Contains(logContent, "=== Ticketr Execution Log ===") {
		t.Error("Header not found in log")
	}
	if !strings.Contains(logContent, "Time:") {
		t.Error("Timestamp not found in header")
	}
	if !strings.Contains(logContent, "Log File:") {
		t.Error("Log file path not found in header")
	}
}

func TestFileLogger_Footer(t *testing.T) {
	tmpDir := t.TempDir()

	config := LogConfig{
		LogDir:  tmpDir,
		Verbose: false,
	}

	logger, err := NewFileLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	logger.Close()

	// Read log file content
	entries, _ := os.ReadDir(tmpDir)
	logPath := filepath.Join(tmpDir, entries[0].Name())
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)

	// Verify footer was written
	if !strings.Contains(logContent, "=== Execution Complete ===") {
		t.Error("Footer not found in log")
	}
}

func TestCleanup_HandlesEmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Cleanup with no files should not error
	err := Cleanup(tmpDir, 5)
	if err != nil {
		t.Fatalf("Cleanup failed on empty directory: %v", err)
	}
}

func TestCleanup_IgnoresNonLogFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a mix of log and non-log files
	os.WriteFile(filepath.Join(tmpDir, "file1.log"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.log"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file3.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("test"), 0644)

	// Keep only 1 log file
	err := Cleanup(tmpDir, 1)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	entries, _ := os.ReadDir(tmpDir)

	// Should have 1 log file + 2 non-log files = 3 total
	if len(entries) != 3 {
		t.Errorf("Expected 3 files remaining, found %d", len(entries))
	}

	// Verify non-log files still exist
	logCount := 0
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".log") {
			logCount++
		}
	}
	if logCount != 1 {
		t.Errorf("Expected 1 log file remaining, found %d", logCount)
	}
}
