package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Logger defines the interface for logging operations
type Logger interface {
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Section(title string)
	Close() error
}

// FileLogger implements Logger with dual console+file output
type FileLogger struct {
	file       *os.File
	multiWrite io.Writer
	verbose    bool
	redactor   *SensitiveRedactor
}

// LogConfig holds logging configuration
type LogConfig struct {
	LogDir  string
	Verbose bool
}

// NewFileLogger creates a new file logger
func NewFileLogger(config LogConfig) (*FileLogger, error) {
	// Default log directory
	if config.LogDir == "" {
		config.LogDir = ".ticketr/logs"
	}

	// Create log directory
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Generate timestamped log file
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logPath := filepath.Join(config.LogDir, fmt.Sprintf("%s.log", timestamp))

	file, err := os.Create(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	// Setup multi-writer (console + file)
	multiWrite := io.MultiWriter(os.Stdout, file)

	logger := &FileLogger{
		file:       file,
		multiWrite: multiWrite,
		verbose:    config.Verbose,
		redactor:   NewSensitiveRedactor(),
	}

	// Write header
	logger.writeHeader(logPath)

	return logger, nil
}

func (l *FileLogger) writeHeader(logPath string) {
	header := fmt.Sprintf("=== Ticketr Execution Log ===\nTime: %s\nLog File: %s\n\n",
		time.Now().Format(time.RFC3339), logPath)
	fmt.Fprint(l.file, header)
}

// Info logs informational messages
func (l *FileLogger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	redacted := l.redactor.Redact(msg)
	timestamp := time.Now().Format("15:04:05")

	// Always write to file
	fmt.Fprintf(l.file, "[%s] INFO: %s\n", timestamp, redacted)

	// Write to console if verbose
	if l.verbose {
		fmt.Fprintf(os.Stdout, "[%s] INFO: %s\n", timestamp, redacted)
	}
}

// Warn logs warning messages
func (l *FileLogger) Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	redacted := l.redactor.Redact(msg)
	timestamp := time.Now().Format("15:04:05")

	output := fmt.Sprintf("[%s] WARN: %s\n", timestamp, redacted)
	// Always write warnings to both console and file
	fmt.Fprint(l.multiWrite, output)
}

// Error logs error messages
func (l *FileLogger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	redacted := l.redactor.Redact(msg)
	timestamp := time.Now().Format("15:04:05")

	output := fmt.Sprintf("[%s] ERROR: %s\n", timestamp, redacted)
	// Always write errors to both console and file
	fmt.Fprint(l.multiWrite, output)
}

// Section writes a section header
func (l *FileLogger) Section(title string) {
	separator := strings.Repeat("=", 50)
	output := fmt.Sprintf("\n%s\n%s\n%s\n", separator, title, separator)
	fmt.Fprint(l.file, output)

	if l.verbose {
		fmt.Fprint(os.Stdout, output)
	}
}

// Close closes the log file
func (l *FileLogger) Close() error {
	if l.file != nil {
		fmt.Fprintf(l.file, "\n=== Execution Complete ===\nTime: %s\n", time.Now().Format(time.RFC3339))
		return l.file.Close()
	}
	return nil
}

// SensitiveRedactor handles sensitive data redaction
type SensitiveRedactor struct {
	patterns []*regexp.Regexp
}

func NewSensitiveRedactor() *SensitiveRedactor {
	return &SensitiveRedactor{
		patterns: []*regexp.Regexp{
			// API keys (various formats)
			regexp.MustCompile(`(?i)(api[_-]?key|token|secret|password)[\s:=]+[^\s\n]+`),
			// Email addresses
			regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
			// URLs with credentials
			regexp.MustCompile(`https?://[^:]+:[^@]+@[^\s]+`),
			// Base64-looking strings (likely credentials)
			regexp.MustCompile(`\b[A-Za-z0-9+/]{40,}={0,2}\b`),
		},
	}
}

func (r *SensitiveRedactor) Redact(text string) string {
	redacted := text
	for _, pattern := range r.patterns {
		redacted = pattern.ReplaceAllString(redacted, "[REDACTED]")
	}
	return redacted
}

// Cleanup removes old log files, keeping only the most recent N files
func Cleanup(logDir string, keepCount int) error {
	if keepCount <= 0 {
		return nil
	}

	entries, err := os.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	// Filter for .log files
	var logFiles []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".log") {
			logFiles = append(logFiles, entry)
		}
	}

	// If we have more than keepCount, remove oldest
	if len(logFiles) > keepCount {
		// Sort by modification time (oldest first)
		// Note: This is a simplified implementation; full implementation would sort by mtime
		toRemove := len(logFiles) - keepCount
		for i := 0; i < toRemove; i++ {
			logPath := filepath.Join(logDir, logFiles[i].Name())
			if err := os.Remove(logPath); err != nil {
				log.Printf("Warning: Failed to remove old log file %s: %v", logPath, err)
			}
		}
	}

	return nil
}
