package services

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

func TestNewPathResolver(t *testing.T) {
	pr, err := NewPathResolver()
	if err != nil {
		t.Fatalf("NewPathResolver() failed: %v", err)
	}

	if pr == nil {
		t.Fatal("NewPathResolver() returned nil")
	}

	if pr.appName != "ticketr" {
		t.Errorf("Expected appName to be 'ticketr', got %s", pr.appName)
	}
}

func TestNewPathResolverWithOptions_EmptyAppName(t *testing.T) {
	_, err := NewPathResolverWithOptions("", os.Getenv, os.UserHomeDir)
	if err == nil {
		t.Fatal("Expected error for empty app name, got nil")
	}
	if !strings.Contains(err.Error(), "app name cannot be empty") {
		t.Errorf("Expected error about empty app name, got: %v", err)
	}
}

func TestNewPathResolverWithOptions_HomeError(t *testing.T) {
	homeGetter := func() (string, error) {
		return "", os.ErrNotExist
	}

	_, err := NewPathResolverWithOptions("test", os.Getenv, homeGetter)
	if err == nil {
		t.Fatal("Expected error when home directory cannot be determined")
	}
	if !strings.Contains(err.Error(), "failed to get home directory") {
		t.Errorf("Expected error about home directory, got: %v", err)
	}
}

func TestPathResolver_UnixPaths(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping Unix path tests on Windows")
	}

	tests := []struct {
		name       string
		envVars    map[string]string
		wantConfig string
		wantData   string
		wantCache  string
	}{
		{
			name:       "default paths without XDG vars",
			envVars:    map[string]string{},
			wantConfig: "/home/testuser/.config/ticketr",
			wantData:   "/home/testuser/.local/share/ticketr",
			wantCache:  "/home/testuser/.cache/ticketr",
		},
		{
			name: "with XDG_CONFIG_HOME",
			envVars: map[string]string{
				"XDG_CONFIG_HOME": "/custom/config",
			},
			wantConfig: "/custom/config/ticketr",
			wantData:   "/home/testuser/.local/share/ticketr",
			wantCache:  "/home/testuser/.cache/ticketr",
		},
		{
			name: "with XDG_DATA_HOME",
			envVars: map[string]string{
				"XDG_DATA_HOME": "/custom/data",
			},
			wantConfig: "/home/testuser/.config/ticketr",
			wantData:   "/custom/data/ticketr",
			wantCache:  "/home/testuser/.cache/ticketr",
		},
		{
			name: "with XDG_CACHE_HOME",
			envVars: map[string]string{
				"XDG_CACHE_HOME": "/custom/cache",
			},
			wantConfig: "/home/testuser/.config/ticketr",
			wantData:   "/home/testuser/.local/share/ticketr",
			wantCache:  "/custom/cache/ticketr",
		},
		{
			name: "with all XDG vars",
			envVars: map[string]string{
				"XDG_CONFIG_HOME": "/xdg/config",
				"XDG_DATA_HOME":   "/xdg/data",
				"XDG_CACHE_HOME":  "/xdg/cache",
			},
			wantConfig: "/xdg/config/ticketr",
			wantData:   "/xdg/data/ticketr",
			wantCache:  "/xdg/cache/ticketr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envGetter := func(key string) string {
				return tt.envVars[key]
			}
			homeGetter := func() (string, error) {
				return "/home/testuser", nil
			}

			pr, err := NewPathResolverWithOptions("ticketr", envGetter, homeGetter)
			if err != nil {
				t.Fatalf("Failed to create PathResolver: %v", err)
			}

			if got := pr.ConfigDir(); got != tt.wantConfig {
				t.Errorf("ConfigDir() = %v, want %v", got, tt.wantConfig)
			}
			if got := pr.DataDir(); got != tt.wantData {
				t.Errorf("DataDir() = %v, want %v", got, tt.wantData)
			}
			if got := pr.CacheDir(); got != tt.wantCache {
				t.Errorf("CacheDir() = %v, want %v", got, tt.wantCache)
			}
		})
	}
}

func TestPathResolver_FilePaths(t *testing.T) {
	envGetter := func(key string) string {
		return ""
	}
	homeGetter := func() (string, error) {
		return "/home/testuser", nil
	}

	pr, err := NewPathResolverWithOptions("ticketr", envGetter, homeGetter)
	if err != nil {
		t.Fatalf("Failed to create PathResolver: %v", err)
	}

	tests := []struct {
		name     string
		method   func(string) string
		filename string
		want     string
	}{
		{
			name:     "ConfigFile with filename",
			method:   pr.ConfigFile,
			filename: "config.yaml",
			want:     filepath.Join(pr.ConfigDir(), "config.yaml"),
		},
		{
			name:     "ConfigFile without filename",
			method:   pr.ConfigFile,
			filename: "",
			want:     pr.ConfigDir(),
		},
		{
			name:     "DataFile with filename",
			method:   pr.DataFile,
			filename: "ticketr.db",
			want:     filepath.Join(pr.DataDir(), "ticketr.db"),
		},
		{
			name:     "DataFile without filename",
			method:   pr.DataFile,
			filename: "",
			want:     pr.DataDir(),
		},
		{
			name:     "CacheFile with filename",
			method:   pr.CacheFile,
			filename: "jira_schema.json",
			want:     filepath.Join(pr.CacheDir(), "jira_schema.json"),
		},
		{
			name:     "CacheFile without filename",
			method:   pr.CacheFile,
			filename: "",
			want:     pr.CacheDir(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.method(tt.filename); got != tt.want {
				t.Errorf("%s(%q) = %v, want %v", tt.name, tt.filename, got, tt.want)
			}
		})
	}
}

func TestPathResolver_SpecificPaths(t *testing.T) {
	envGetter := func(key string) string {
		return ""
	}
	homeGetter := func() (string, error) {
		return "/home/testuser", nil
	}

	pr, err := NewPathResolverWithOptions("ticketr", envGetter, homeGetter)
	if err != nil {
		t.Fatalf("Failed to create PathResolver: %v", err)
	}

	tests := []struct {
		name   string
		method func() string
		want   string
	}{
		{
			name:   "TemplatesDir",
			method: pr.TemplatesDir,
			want:   filepath.Join(pr.DataDir(), "templates"),
		},
		{
			name:   "PluginsDir",
			method: pr.PluginsDir,
			want:   filepath.Join(pr.DataDir(), "plugins"),
		},
		{
			name:   "LogsDir",
			method: pr.LogsDir,
			want:   filepath.Join(pr.CacheDir(), "logs"),
		},
		{
			name:   "DatabasePath",
			method: pr.DatabasePath,
			want:   filepath.Join(pr.DataDir(), "ticketr.db"),
		},
		{
			name:   "ConfigPath",
			method: pr.ConfigPath,
			want:   filepath.Join(pr.ConfigDir(), "config.yaml"),
		},
		{
			name:   "WorkspacesPath",
			method: pr.WorkspacesPath,
			want:   filepath.Join(pr.ConfigDir(), "workspaces.yaml"),
		},
		{
			name:   "JiraCachePath",
			method: pr.JiraCachePath,
			want:   filepath.Join(pr.CacheDir(), "jira_schema.json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.method(); got != tt.want {
				t.Errorf("%s() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestPathResolver_EnsureDirectories(t *testing.T) {
	tempDir := t.TempDir()

	envGetter := func(key string) string {
		switch key {
		case "XDG_CONFIG_HOME":
			return filepath.Join(tempDir, "config")
		case "XDG_DATA_HOME":
			return filepath.Join(tempDir, "data")
		case "XDG_CACHE_HOME":
			return filepath.Join(tempDir, "cache")
		default:
			return ""
		}
	}
	homeGetter := func() (string, error) {
		return tempDir, nil
	}

	pr, err := NewPathResolverWithOptions("ticketr", envGetter, homeGetter)
	if err != nil {
		t.Fatalf("Failed to create PathResolver: %v", err)
	}

	// Ensure directories don't exist initially
	dirs := []string{
		pr.ConfigDir(),
		pr.DataDir(),
		pr.CacheDir(),
		pr.TemplatesDir(),
		pr.PluginsDir(),
		pr.LogsDir(),
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			t.Fatalf("Directory %s should not exist initially", dir)
		}
	}

	// Create directories
	if err := pr.EnsureDirectories(); err != nil {
		t.Fatalf("EnsureDirectories() failed: %v", err)
	}

	// Verify all directories were created
	for _, dir := range dirs {
		info, err := os.Stat(dir)
		if err != nil {
			t.Errorf("Directory %s was not created: %v", dir, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("Path %s is not a directory", dir)
		}
	}

	// Test idempotency - should not fail when directories exist
	if err := pr.EnsureDirectories(); err != nil {
		t.Fatalf("EnsureDirectories() failed on second call: %v", err)
	}
}

func TestPathResolver_Exists(t *testing.T) {
	tempDir := t.TempDir()

	envGetter := func(key string) string { return "" }
	homeGetter := func() (string, error) { return tempDir, nil }

	pr, err := NewPathResolverWithOptions("ticketr", envGetter, homeGetter)
	if err != nil {
		t.Fatalf("Failed to create PathResolver: %v", err)
	}

	// Test non-existent path
	nonExistent := filepath.Join(tempDir, "does-not-exist")
	if pr.Exists(nonExistent) {
		t.Errorf("Exists() returned true for non-existent path")
	}

	// Test existing directory
	if !pr.Exists(tempDir) {
		t.Errorf("Exists() returned false for existing directory")
	}

	// Test existing file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if !pr.Exists(testFile) {
		t.Errorf("Exists() returned false for existing file")
	}
}

func TestPathResolver_IsDirectory(t *testing.T) {
	tempDir := t.TempDir()

	envGetter := func(key string) string { return "" }
	homeGetter := func() (string, error) { return tempDir, nil }

	pr, err := NewPathResolverWithOptions("ticketr", envGetter, homeGetter)
	if err != nil {
		t.Fatalf("Failed to create PathResolver: %v", err)
	}

	// Test non-existent path
	nonExistent := filepath.Join(tempDir, "does-not-exist")
	if pr.IsDirectory(nonExistent) {
		t.Errorf("IsDirectory() returned true for non-existent path")
	}

	// Test directory
	if !pr.IsDirectory(tempDir) {
		t.Errorf("IsDirectory() returned false for directory")
	}

	// Test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if pr.IsDirectory(testFile) {
		t.Errorf("IsDirectory() returned true for file")
	}
}

func TestPathResolver_CleanCache(t *testing.T) {
	tempDir := t.TempDir()

	envGetter := func(key string) string {
		if key == "XDG_CACHE_HOME" {
			return filepath.Join(tempDir, "cache")
		}
		return ""
	}
	homeGetter := func() (string, error) { return tempDir, nil }

	pr, err := NewPathResolverWithOptions("ticketr", envGetter, homeGetter)
	if err != nil {
		t.Fatalf("Failed to create PathResolver: %v", err)
	}

	// Create cache directory with some files
	cacheDir := pr.CacheDir()
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		t.Fatalf("Failed to create cache directory: %v", err)
	}

	// Create test files and directories
	testFile1 := filepath.Join(cacheDir, "test1.json")
	testFile2 := filepath.Join(cacheDir, "test2.log")
	testDir := filepath.Join(cacheDir, "logs")
	testFileInDir := filepath.Join(testDir, "app.log")

	if err := os.WriteFile(testFile1, []byte("test1"), 0644); err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("test2"), 0644); err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	if err := os.WriteFile(testFileInDir, []byte("log"), 0644); err != nil {
		t.Fatalf("Failed to create test file in directory: %v", err)
	}

	// Clean cache
	if err := pr.CleanCache(); err != nil {
		t.Fatalf("CleanCache() failed: %v", err)
	}

	// Verify cache directory still exists but is empty
	if !pr.Exists(cacheDir) {
		t.Errorf("Cache directory was removed, should be preserved")
	}

	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		t.Fatalf("Failed to read cache directory: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("Cache directory is not empty, contains %d entries", len(entries))
	}

	// Test cleaning non-existent cache directory
	if err := os.RemoveAll(cacheDir); err != nil {
		t.Fatalf("Failed to remove cache directory: %v", err)
	}
	if err := pr.CleanCache(); err != nil {
		t.Errorf("CleanCache() failed on non-existent directory: %v", err)
	}
}

func TestPathResolver_Summary(t *testing.T) {
	envGetter := func(key string) string {
		switch key {
		case "XDG_CONFIG_HOME":
			return "/custom/config"
		case "XDG_DATA_HOME":
			return "/custom/data"
		case "XDG_CACHE_HOME":
			return "/custom/cache"
		default:
			return ""
		}
	}
	homeGetter := func() (string, error) {
		return "/home/testuser", nil
	}

	pr, err := NewPathResolverWithOptions("ticketr", envGetter, homeGetter)
	if err != nil {
		t.Fatalf("Failed to create PathResolver: %v", err)
	}

	summary := pr.Summary()

	// Check that summary contains expected paths
	expectedPaths := []string{
		"/custom/config/ticketr",
		"/custom/data/ticketr",
		"/custom/cache/ticketr",
		"ticketr.db",
		"config.yaml",
		"workspaces.yaml",
		"templates",
		"plugins",
		"logs",
		"jira_schema.json",
	}

	for _, expected := range expectedPaths {
		if !strings.Contains(summary, expected) {
			t.Errorf("Summary does not contain expected path: %s", expected)
		}
	}

	// Check that summary has the right format
	if !strings.HasPrefix(summary, "Ticketr Paths:") {
		t.Errorf("Summary does not start with expected header")
	}

	lines := strings.Split(summary, "\n")
	if len(lines) != 11 { // Header + 10 paths
		t.Errorf("Summary should have 11 lines, got %d", len(lines))
	}
}

func TestGetPathResolver_Singleton(t *testing.T) {
	defer ResetPathResolver()

	pr1, err1 := GetPathResolver()
	if err1 != nil {
		t.Fatalf("First GetPathResolver() failed: %v", err1)
	}
	if pr1 == nil {
		t.Fatal("First GetPathResolver() returned nil")
	}

	pr2, err2 := GetPathResolver()
	if err2 != nil {
		t.Fatalf("Second GetPathResolver() failed: %v", err2)
	}
	if pr2 == nil {
		t.Fatal("Second GetPathResolver() returned nil")
	}

	// Verify same instance (pointer equality)
	if pr1 != pr2 {
		t.Error("GetPathResolver() returned different instances - singleton pattern violated")
	}
}

func TestGetPathResolver_Concurrent(t *testing.T) {
	defer ResetPathResolver()

	const goroutines = 100
	instances := make([]*PathResolver, goroutines)
	errors := make([]error, goroutines)

	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Launch concurrent calls to GetPathResolver
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			pr, err := GetPathResolver()
			instances[idx] = pr
			errors[idx] = err
		}(i)
	}

	wg.Wait()

	// Verify all calls succeeded
	for i, err := range errors {
		if err != nil {
			t.Errorf("Goroutine %d: GetPathResolver() failed: %v", i, err)
		}
	}

	// Verify all instances are the same
	firstInstance := instances[0]
	for i, instance := range instances {
		if instance != firstInstance {
			t.Errorf("Goroutine %d: Got different instance (concurrent initialization not properly synchronized)", i)
		}
	}
}

func TestResetPathResolver(t *testing.T) {
	pr1, err := GetPathResolver()
	if err != nil {
		t.Fatalf("GetPathResolver() failed: %v", err)
	}

	ResetPathResolver()

	pr2, err := GetPathResolver()
	if err != nil {
		t.Fatalf("GetPathResolver() after reset failed: %v", err)
	}

	// Should be different instances after reset
	if pr1 == pr2 {
		t.Error("ResetPathResolver() did not clear singleton - instances are the same")
	}
}

func TestGetPathResolver_EnsuresDirectories(t *testing.T) {
	defer ResetPathResolver()

	// GetPathResolver uses actual system paths
	// We verify that it creates directories on first call
	pr, err := GetPathResolver()
	if err != nil {
		t.Fatalf("GetPathResolver() failed: %v", err)
	}

	// Verify that essential directories exist after first call
	dirs := []string{
		pr.ConfigDir(),
		pr.DataDir(),
		pr.CacheDir(),
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				t.Errorf("Directory %s was not created by GetPathResolver()", dir)
			} else {
				t.Errorf("Error checking directory %s: %v", dir, err)
			}
		}
	}
}
