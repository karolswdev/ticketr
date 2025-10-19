package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"
)

// FeatureFlags manages feature toggles for v3.0 migration
type FeatureFlags struct {
	UseSQLite         bool   `mapstructure:"use_sqlite"`
	SQLitePath        string `mapstructure:"sqlite_path"`
	EnableWorkspaces  bool   `mapstructure:"enable_workspaces"`
	EnableTUI         bool   `mapstructure:"enable_tui"`
	AutoMigrate       bool   `mapstructure:"auto_migrate"`
	VerboseLogging    bool   `mapstructure:"verbose_logging"`
}

// DefaultFeatures returns default feature flags for backward compatibility
func DefaultFeatures() *FeatureFlags {
	return &FeatureFlags{
		UseSQLite:        false, // Default to file-based for v2 compatibility
		SQLitePath:       "",    // Will be set to XDG path if empty
		EnableWorkspaces: false, // v3 feature, disabled by default
		EnableTUI:        false, // v3 feature, disabled by default
		AutoMigrate:      false, // Don't auto-migrate without user consent
		VerboseLogging:   false,
	}
}

// LoadFeatures loads feature flags from configuration
func LoadFeatures() (*FeatureFlags, error) {
	features := DefaultFeatures()

	// Check environment variables first (highest priority)
	if envVal := os.Getenv("TICKETR_USE_SQLITE"); envVal != "" {
		if val, err := strconv.ParseBool(envVal); err == nil {
			features.UseSQLite = val
		}
	}

	if envVal := os.Getenv("TICKETR_SQLITE_PATH"); envVal != "" {
		features.SQLitePath = envVal
	}

	if envVal := os.Getenv("TICKETR_ENABLE_WORKSPACES"); envVal != "" {
		if val, err := strconv.ParseBool(envVal); err == nil {
			features.EnableWorkspaces = val
		}
	}

	if envVal := os.Getenv("TICKETR_ENABLE_TUI"); envVal != "" {
		if val, err := strconv.ParseBool(envVal); err == nil {
			features.EnableTUI = val
		}
	}

	if envVal := os.Getenv("TICKETR_AUTO_MIGRATE"); envVal != "" {
		if val, err := strconv.ParseBool(envVal); err == nil {
			features.AutoMigrate = val
		}
	}

	if envVal := os.Getenv("TICKETR_VERBOSE"); envVal != "" {
		if val, err := strconv.ParseBool(envVal); err == nil {
			features.VerboseLogging = val
		}
	}

	// Try to load from config file
	configPath := getConfigPath()
	if configPath != "" {
		viper.SetConfigFile(configPath)
		viper.SetConfigType("yaml")

		if err := viper.ReadInConfig(); err == nil {
			// Config file exists, unmarshal features section
			if err := viper.UnmarshalKey("features", features); err != nil {
				return features, fmt.Errorf("failed to parse features config: %w", err)
			}
		}
	}

	// Set default SQLite path if not specified
	if features.UseSQLite && features.SQLitePath == "" {
		features.SQLitePath = getDefaultSQLitePath()
	}

	// Validate feature combinations
	if err := features.Validate(); err != nil {
		return features, err
	}

	return features, nil
}

// Validate checks for invalid feature flag combinations
func (f *FeatureFlags) Validate() error {
	// Workspaces require SQLite
	if f.EnableWorkspaces && !f.UseSQLite {
		return fmt.Errorf("workspaces feature requires SQLite to be enabled")
	}

	// TUI requires SQLite
	if f.EnableTUI && !f.UseSQLite {
		return fmt.Errorf("TUI feature requires SQLite to be enabled")
	}

	return nil
}

// IsV3Enabled returns true if any v3 feature is enabled
func (f *FeatureFlags) IsV3Enabled() bool {
	return f.UseSQLite || f.EnableWorkspaces || f.EnableTUI
}

// RequiresMigration returns true if migration from v2 is needed
func (f *FeatureFlags) RequiresMigration() bool {
	if !f.UseSQLite {
		return false
	}

	// Check if SQLite database exists
	if f.SQLitePath != "" {
		if _, err := os.Stat(f.SQLitePath); os.IsNotExist(err) {
			// Database doesn't exist, migration might be needed
			return true
		}
	}

	return false
}

// SaveFeatures saves feature flags to configuration file
func SaveFeatures(features *FeatureFlags) error {
	configPath := getConfigPath()
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create a map with keys matching the mapstructure tags
	// to ensure viper writes YAML with snake_case keys that match what LoadFeatures expects
	featuresMap := map[string]interface{}{
		"use_sqlite":        features.UseSQLite,
		"sqlite_path":       features.SQLitePath,
		"enable_workspaces": features.EnableWorkspaces,
		"enable_tui":        features.EnableTUI,
		"auto_migrate":      features.AutoMigrate,
		"verbose_logging":   features.VerboseLogging,
	}

	viper.Set("features", featuresMap)
	return viper.WriteConfigAs(configPath)
}

// EnableV3Features progressively enables v3 features
func EnableV3Features(phase string) (*FeatureFlags, error) {
	features, err := LoadFeatures()
	if err != nil {
		return nil, err
	}

	switch phase {
	case "alpha":
		// Phase 1: SQLite only
		features.UseSQLite = true
		features.AutoMigrate = true

	case "beta":
		// Phase 2: SQLite + Workspaces
		features.UseSQLite = true
		features.EnableWorkspaces = true
		features.AutoMigrate = true

	case "rc":
		// Phase 3: SQLite + Workspaces + TUI
		features.UseSQLite = true
		features.EnableWorkspaces = true
		features.EnableTUI = true
		features.AutoMigrate = true

	case "stable", "v3":
		// All features enabled
		features.UseSQLite = true
		features.EnableWorkspaces = true
		features.EnableTUI = true
		features.AutoMigrate = false // Don't auto-migrate in stable

	default:
		return nil, fmt.Errorf("unknown phase: %s (use: alpha, beta, rc, stable)", phase)
	}

	if err := SaveFeatures(features); err != nil {
		return features, fmt.Errorf("failed to save features: %w", err)
	}

	return features, nil
}

// Helper functions

func getConfigPath() string {
	// Check for explicit config file
	if configFile := os.Getenv("TICKETR_CONFIG"); configFile != "" {
		return configFile
	}

	// Check current directory
	if _, err := os.Stat(".ticketr.yaml"); err == nil {
		return ".ticketr.yaml"
	}

	// Check XDG config directory
	return getDefaultConfigPath()
}

func getDefaultConfigPath() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		configDir = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(configDir, "ticketr", "config.yaml")
}

func getDefaultSQLitePath() string {
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "ticketr.db" // Fallback to current directory
		}
		dataDir = filepath.Join(homeDir, ".local", "share")
	}
	return filepath.Join(dataDir, "ticketr", "ticketr.db")
}

// FeatureStatus provides human-readable feature status
func (f *FeatureFlags) Status() string {
	status := "Feature Flags Status:\n"
	status += "====================\n"
	status += fmt.Sprintf("SQLite Backend:    %s\n", enabledStr(f.UseSQLite))
	if f.UseSQLite {
		status += fmt.Sprintf("  Path: %s\n", f.SQLitePath)
	}
	status += fmt.Sprintf("Workspaces:        %s\n", enabledStr(f.EnableWorkspaces))
	status += fmt.Sprintf("TUI Interface:     %s\n", enabledStr(f.EnableTUI))
	status += fmt.Sprintf("Auto-Migration:    %s\n", enabledStr(f.AutoMigrate))
	status += fmt.Sprintf("Verbose Logging:   %s\n", enabledStr(f.VerboseLogging))

	if f.IsV3Enabled() {
		status += "\nMode: v3.0 (enhanced features enabled)"
	} else {
		status += "\nMode: v2.x (backward compatibility mode)"
	}

	return status
}

func enabledStr(enabled bool) string {
	if enabled {
		return "✓ Enabled"
	}
	return "✗ Disabled"
}