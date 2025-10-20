package services

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// PathResolver handles application directory paths following platform conventions.
// On Linux/macOS, it follows the XDG Base Directory specification.
// On Windows, it uses standard Windows directory conventions.
type PathResolver struct {
	// configHome is the base directory for configuration files
	// XDG_CONFIG_HOME on Unix, %APPDATA% on Windows
	configHome string

	// dataHome is the base directory for data files
	// XDG_DATA_HOME on Unix, %LOCALAPPDATA% on Windows
	dataHome string

	// cacheHome is the base directory for cache files
	// XDG_CACHE_HOME on Unix, %TEMP% on Windows
	cacheHome string

	// homeDir is the user's home directory
	homeDir string

	// appName is the application name used in paths
	appName string

	// envGetter allows for dependency injection of environment variables
	envGetter func(string) string
}

// NewPathResolver creates a new PathResolver with system defaults.
// It automatically detects the operating system and applies appropriate conventions.
func NewPathResolver() (*PathResolver, error) {
	return NewPathResolverWithOptions("ticketr", os.Getenv, os.UserHomeDir)
}

// NewPathResolverWithOptions creates a new PathResolver with custom options.
// This constructor is primarily used for testing to inject dependencies.
func NewPathResolverWithOptions(
	appName string,
	envGetter func(string) string,
	homeGetter func() (string, error),
) (*PathResolver, error) {
	if appName == "" {
		return nil, fmt.Errorf("app name cannot be empty")
	}

	homeDir, err := homeGetter()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	pr := &PathResolver{
		appName:   appName,
		homeDir:   homeDir,
		envGetter: envGetter,
	}

	// Initialize paths based on operating system
	if runtime.GOOS == "windows" {
		pr.initWindowsPaths()
	} else {
		pr.initUnixPaths()
	}

	return pr, nil
}

// initUnixPaths initializes paths following XDG Base Directory specification.
func (pr *PathResolver) initUnixPaths() {
	// Config directory: XDG_CONFIG_HOME or ~/.config
	if configHome := pr.envGetter("XDG_CONFIG_HOME"); configHome != "" {
		pr.configHome = configHome
	} else {
		pr.configHome = filepath.Join(pr.homeDir, ".config")
	}

	// Data directory: XDG_DATA_HOME or ~/.local/share
	if dataHome := pr.envGetter("XDG_DATA_HOME"); dataHome != "" {
		pr.dataHome = dataHome
	} else {
		pr.dataHome = filepath.Join(pr.homeDir, ".local", "share")
	}

	// Cache directory: XDG_CACHE_HOME or ~/.cache
	if cacheHome := pr.envGetter("XDG_CACHE_HOME"); cacheHome != "" {
		pr.cacheHome = cacheHome
	} else {
		pr.cacheHome = filepath.Join(pr.homeDir, ".cache")
	}
}

// initWindowsPaths initializes paths following Windows conventions.
func (pr *PathResolver) initWindowsPaths() {
	// Config directory: %APPDATA%
	if appData := pr.envGetter("APPDATA"); appData != "" {
		pr.configHome = appData
	} else {
		// Fallback to user profile if APPDATA is not set
		pr.configHome = filepath.Join(pr.homeDir, "AppData", "Roaming")
	}

	// Data directory: %LOCALAPPDATA%
	if localAppData := pr.envGetter("LOCALAPPDATA"); localAppData != "" {
		pr.dataHome = localAppData
	} else {
		// Fallback to user profile if LOCALAPPDATA is not set
		pr.dataHome = filepath.Join(pr.homeDir, "AppData", "Local")
	}

	// Cache directory: %TEMP%
	if temp := pr.envGetter("TEMP"); temp != "" {
		pr.cacheHome = temp
	} else if temp := pr.envGetter("TMP"); temp != "" {
		pr.cacheHome = temp
	} else {
		// Fallback to LocalAppData\Temp
		pr.cacheHome = filepath.Join(pr.dataHome, "Temp")
	}
}

// ConfigDir returns the base configuration directory for the application.
// On Unix: ~/.config/ticketr (or $XDG_CONFIG_HOME/ticketr)
// On Windows: %APPDATA%\ticketr
func (pr *PathResolver) ConfigDir() string {
	return filepath.Join(pr.configHome, pr.appName)
}

// DataDir returns the base data directory for the application.
// On Unix: ~/.local/share/ticketr (or $XDG_DATA_HOME/ticketr)
// On Windows: %LOCALAPPDATA%\ticketr
func (pr *PathResolver) DataDir() string {
	return filepath.Join(pr.dataHome, pr.appName)
}

// CacheDir returns the base cache directory for the application.
// On Unix: ~/.cache/ticketr (or $XDG_CACHE_HOME/ticketr)
// On Windows: %TEMP%\ticketr
func (pr *PathResolver) CacheDir() string {
	return filepath.Join(pr.cacheHome, pr.appName)
}

// ConfigFile returns the full path to a configuration file.
func (pr *PathResolver) ConfigFile(filename string) string {
	if filename == "" {
		return pr.ConfigDir()
	}
	return filepath.Join(pr.ConfigDir(), filename)
}

// DataFile returns the full path to a data file.
func (pr *PathResolver) DataFile(filename string) string {
	if filename == "" {
		return pr.DataDir()
	}
	return filepath.Join(pr.DataDir(), filename)
}

// CacheFile returns the full path to a cache file.
func (pr *PathResolver) CacheFile(filename string) string {
	if filename == "" {
		return pr.CacheDir()
	}
	return filepath.Join(pr.CacheDir(), filename)
}

// TemplatesDir returns the directory for user templates.
func (pr *PathResolver) TemplatesDir() string {
	return filepath.Join(pr.DataDir(), "templates")
}

// PluginsDir returns the directory for plugins (future use).
func (pr *PathResolver) PluginsDir() string {
	return filepath.Join(pr.DataDir(), "plugins")
}

// LogsDir returns the directory for log files.
func (pr *PathResolver) LogsDir() string {
	return filepath.Join(pr.CacheDir(), "logs")
}

// DatabasePath returns the full path to the SQLite database file.
func (pr *PathResolver) DatabasePath() string {
	return pr.DataFile("ticketr.db")
}

// ConfigPath returns the full path to the main configuration file.
func (pr *PathResolver) ConfigPath() string {
	return pr.ConfigFile("config.yaml")
}

// WorkspacesPath returns the full path to the workspaces configuration file.
func (pr *PathResolver) WorkspacesPath() string {
	return pr.ConfigFile("workspaces.yaml")
}

// JiraCachePath returns the full path to the Jira schema cache file.
func (pr *PathResolver) JiraCachePath() string {
	return pr.CacheFile("jira_schema.json")
}

// EnsureDirectories creates all required application directories if they don't exist.
// Directories are created with 0755 permissions (subject to umask).
func (pr *PathResolver) EnsureDirectories() error {
	dirs := []string{
		pr.ConfigDir(),
		pr.DataDir(),
		pr.CacheDir(),
		pr.TemplatesDir(),
		pr.PluginsDir(),
		pr.LogsDir(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// EnsureDirectory creates a specific directory if it doesn't exist.
// The directory is created with 0755 permissions (subject to umask).
func (pr *PathResolver) EnsureDirectory(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

// Exists checks if a path exists.
func (pr *PathResolver) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDirectory checks if a path exists and is a directory.
func (pr *PathResolver) IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// CleanCache cleans the cache directory by removing all its contents.
// The cache directory itself is preserved.
func (pr *PathResolver) CleanCache() error {
	cacheDir := pr.CacheDir()
	if !pr.Exists(cacheDir) {
		return nil
	}

	// Read directory contents
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	// Remove each entry
	for _, entry := range entries {
		path := filepath.Join(cacheDir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove %s: %w", path, err)
		}
	}

	return nil
}

// Summary returns a human-readable summary of all paths.
func (pr *PathResolver) Summary() string {
	return fmt.Sprintf(
		"Ticketr Paths:\n"+
			"  Config Dir:     %s\n"+
			"  Data Dir:       %s\n"+
			"  Cache Dir:      %s\n"+
			"  Database:       %s\n"+
			"  Config File:    %s\n"+
			"  Workspaces:     %s\n"+
			"  Templates:      %s\n"+
			"  Plugins:        %s\n"+
			"  Logs:           %s\n"+
			"  Jira Cache:     %s",
		pr.ConfigDir(),
		pr.DataDir(),
		pr.CacheDir(),
		pr.DatabasePath(),
		pr.ConfigPath(),
		pr.WorkspacesPath(),
		pr.TemplatesDir(),
		pr.PluginsDir(),
		pr.LogsDir(),
		pr.JiraCachePath(),
	)
}

// Singleton instance management

var (
	globalPathResolver *PathResolver
	pathResolverOnce   sync.Once
	pathResolverErr    error
)

// GetPathResolver returns the singleton PathResolver instance.
// Safe for concurrent use. Automatically ensures directories exist on first access.
func GetPathResolver() (*PathResolver, error) {
	pathResolverOnce.Do(func() {
		globalPathResolver, pathResolverErr = NewPathResolver()
		if pathResolverErr == nil {
			// Ensure directories exist on first access
			if err := globalPathResolver.EnsureDirectories(); err != nil {
				pathResolverErr = fmt.Errorf("failed to create directories: %w", err)
				globalPathResolver = nil
			}
		}
	})
	return globalPathResolver, pathResolverErr
}

// ResetPathResolver clears the singleton (FOR TESTING ONLY).
// This function is used to isolate test cases that need different PathResolver configurations.
func ResetPathResolver() {
	globalPathResolver = nil
	pathResolverOnce = sync.Once{}
	pathResolverErr = nil
}
