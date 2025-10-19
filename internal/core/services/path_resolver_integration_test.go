//go:build integration

package services

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestPathResolver_EnsureDirectories_Integration(t *testing.T) {
	// Create a temporary directory for testing
	tempHome := t.TempDir()

	// Override HOME environment variable
	originalHome := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		originalHome = os.Getenv("USERPROFILE")
		os.Setenv("USERPROFILE", tempHome)
	} else {
		os.Setenv("HOME", tempHome)
	}
	t.Cleanup(func() {
		if runtime.GOOS == "windows" {
			os.Setenv("USERPROFILE", originalHome)
		} else {
			os.Setenv("HOME", originalHome)
		}
	})

	tests := []struct {
		name        string
		setupEnv    func()
		cleanupEnv  func()
		checkPaths  func(*testing.T, *PathResolver)
	}{
		{
			name: "create directories with default paths",
			setupEnv: func() {
				// No special environment setup
			},
			cleanupEnv: func() {},
			checkPaths: func(t *testing.T, pr *PathResolver) {
				// Verify directories exist
				configDir := pr.ConfigDir()
				cacheDir := pr.CacheDir()
				dataDir := pr.DataDir()

				if _, err := os.Stat(configDir); os.IsNotExist(err) {
					t.Errorf("Config directory was not created: %s", configDir)
				}
				if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
					t.Errorf("Cache directory was not created: %s", cacheDir)
				}
				if _, err := os.Stat(dataDir); os.IsNotExist(err) {
					t.Errorf("Data directory was not created: %s", dataDir)
				}

				// Verify permissions (Unix-like systems)
				if runtime.GOOS != "windows" {
					info, err := os.Stat(configDir)
					if err != nil {
						t.Fatalf("Failed to stat config dir: %v", err)
					}
					if info.Mode().Perm() != 0755 {
						t.Errorf("Expected config dir permissions 0755, got %v", info.Mode().Perm())
					}
				}
			},
		},
		{
			name: "create directories with XDG environment variables",
			setupEnv: func() {
				if runtime.GOOS != "windows" {
					os.Setenv("XDG_CONFIG_HOME", filepath.Join(tempHome, "custom_config"))
					os.Setenv("XDG_DATA_HOME", filepath.Join(tempHome, "custom_data"))
					os.Setenv("XDG_CACHE_HOME", filepath.Join(tempHome, "custom_cache"))
				}
			},
			cleanupEnv: func() {
				if runtime.GOOS != "windows" {
					os.Unsetenv("XDG_CONFIG_HOME")
					os.Unsetenv("XDG_DATA_HOME")
					os.Unsetenv("XDG_CACHE_HOME")
				}
			},
			checkPaths: func(t *testing.T, pr *PathResolver) {
				if runtime.GOOS == "windows" {
					t.Skip("XDG variables not applicable on Windows")
				}

				configDir := pr.ConfigDir()
				dataDir := pr.DataDir()

				// Verify custom paths are used
				if !filepath.IsAbs(configDir) {
					t.Errorf("Config dir should be absolute: %s", configDir)
				}
				if !filepath.IsAbs(dataDir) {
					t.Errorf("Data dir should be absolute: %s", dataDir)
				}

				// Verify directories exist
				if _, err := os.Stat(configDir); os.IsNotExist(err) {
					t.Errorf("Config directory was not created: %s", configDir)
				}
				if _, err := os.Stat(dataDir); os.IsNotExist(err) {
					t.Errorf("Data directory was not created: %s", dataDir)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			tt.setupEnv()
			t.Cleanup(tt.cleanupEnv)

			// Create PathResolver
			pr, err := NewPathResolver()
			if err != nil {
				t.Fatalf("Failed to create PathResolver: %v", err)
			}

			// Ensure directories are created
			err = pr.EnsureDirectories()
			if err != nil {
				t.Fatalf("EnsureDirectories() failed: %v", err)
			}

			// Run checks
			tt.checkPaths(t, pr)
		})
	}
}

func TestPathResolver_FileOperations_Integration(t *testing.T) {
	// Create a temporary directory for testing
	tempHome := t.TempDir()

	// Override HOME environment variable
	originalHome := os.Getenv("HOME")
	if runtime.GOOS == "windows" {
		originalHome = os.Getenv("USERPROFILE")
		os.Setenv("USERPROFILE", tempHome)
	} else {
		os.Setenv("HOME", tempHome)
	}
	t.Cleanup(func() {
		if runtime.GOOS == "windows" {
			os.Setenv("USERPROFILE", originalHome)
		} else {
			os.Setenv("HOME", originalHome)
		}
	})

	pr, err := NewPathResolver()
	if err != nil {
		t.Fatalf("Failed to create PathResolver: %v", err)
	}

	err = pr.EnsureDirectories()
	if err != nil {
		t.Fatalf("EnsureDirectories() failed: %v", err)
	}

	tests := []struct {
		name      string
		operation func(t *testing.T, pr *PathResolver)
	}{
		{
			name: "write and read config file",
			operation: func(t *testing.T, pr *PathResolver) {
				configPath := pr.ConfigFile("test_config.yaml")
				testData := []byte(`test: data`)

				// Write file
				err := os.WriteFile(configPath, testData, 0644)
				if err != nil {
					t.Fatalf("Failed to write config file: %v", err)
				}

				// Read file
				readData, err := os.ReadFile(configPath)
				if err != nil {
					t.Fatalf("Failed to read config file: %v", err)
				}
				if string(readData) != string(testData) {
					t.Errorf("Data mismatch: got %s, want %s", readData, testData)
				}
			},
		},
		{
			name: "write and read data file",
			operation: func(t *testing.T, pr *PathResolver) {
				dataPath := pr.DataFile("test_data.db")
				testData := []byte("test database content")

				// Write file
				err := os.WriteFile(dataPath, testData, 0644)
				if err != nil {
					t.Fatalf("Failed to write data file: %v", err)
				}

				// Read file
				readData, err := os.ReadFile(dataPath)
				if err != nil {
					t.Fatalf("Failed to read data file: %v", err)
				}
				if string(readData) != string(testData) {
					t.Errorf("Data mismatch: got %s, want %s", readData, testData)
				}
			},
		},
		{
			name: "create database file",
			operation: func(t *testing.T, pr *PathResolver) {
				dbPath := pr.DatabasePath()
				testData := []byte("fake database content")

				// Write file
				err := os.WriteFile(dbPath, testData, 0644)
				if err != nil {
					t.Fatalf("Failed to write database file: %v", err)
				}

				// Verify file exists
				if _, err := os.Stat(dbPath); os.IsNotExist(err) {
					t.Errorf("Database file was not created: %s", dbPath)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.operation(t, pr)
		})
	}
}

func TestPathResolver_XDGCompliance_Integration(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG specification not applicable on Windows")
	}

	// Create a temporary directory for testing
	tempHome := t.TempDir()

	// Override HOME environment variable
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})

	tests := []struct {
		name        string
		xdgConfig   string
		xdgData     string
		xdgCache    string
		verifyPaths func(t *testing.T, pr *PathResolver)
	}{
		{
			name:      "XDG_CONFIG_HOME set",
			xdgConfig: filepath.Join(tempHome, "my_config"),
			xdgData:   "",
			xdgCache:  "",
			verifyPaths: func(t *testing.T, pr *PathResolver) {
				configDir := pr.ConfigDir()
				if !filepath.IsAbs(configDir) {
					t.Errorf("Config dir should be absolute: %s", configDir)
				}
			},
		},
		{
			name:      "XDG_DATA_HOME set",
			xdgConfig: "",
			xdgData:   filepath.Join(tempHome, "my_data"),
			xdgCache:  "",
			verifyPaths: func(t *testing.T, pr *PathResolver) {
				dataDir := pr.DataDir()
				if !filepath.IsAbs(dataDir) {
					t.Errorf("Data dir should be absolute: %s", dataDir)
				}
			},
		},
		{
			name:      "all XDG variables set",
			xdgConfig: filepath.Join(tempHome, "custom_config"),
			xdgData:   filepath.Join(tempHome, "custom_data"),
			xdgCache:  filepath.Join(tempHome, "custom_cache"),
			verifyPaths: func(t *testing.T, pr *PathResolver) {
				configDir := pr.ConfigDir()
				dataDir := pr.DataDir()
				cacheDir := pr.CacheDir()

				// Verify all are absolute
				if !filepath.IsAbs(configDir) {
					t.Errorf("Config dir should be absolute: %s", configDir)
				}
				if !filepath.IsAbs(dataDir) {
					t.Errorf("Data dir should be absolute: %s", dataDir)
				}
				if !filepath.IsAbs(cacheDir) {
					t.Errorf("Cache dir should be absolute: %s", cacheDir)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup XDG environment variables
			if tt.xdgConfig != "" {
				os.Setenv("XDG_CONFIG_HOME", tt.xdgConfig)
				t.Cleanup(func() { os.Unsetenv("XDG_CONFIG_HOME") })
			}
			if tt.xdgData != "" {
				os.Setenv("XDG_DATA_HOME", tt.xdgData)
				t.Cleanup(func() { os.Unsetenv("XDG_DATA_HOME") })
			}
			if tt.xdgCache != "" {
				os.Setenv("XDG_CACHE_HOME", tt.xdgCache)
				t.Cleanup(func() { os.Unsetenv("XDG_CACHE_HOME") })
			}

			// Create PathResolver
			pr, err := NewPathResolver()
			if err != nil {
				t.Fatalf("Failed to create PathResolver: %v", err)
			}

			err = pr.EnsureDirectories()
			if err != nil {
				t.Fatalf("EnsureDirectories() failed: %v", err)
			}

			// Verify paths
			tt.verifyPaths(t, pr)

			// Verify directories actually exist
			if _, err := os.Stat(pr.ConfigDir()); os.IsNotExist(err) {
				t.Errorf("Config directory was not created: %s", pr.ConfigDir())
			}
			if _, err := os.Stat(pr.DataDir()); os.IsNotExist(err) {
				t.Errorf("Data directory was not created: %s", pr.DataDir())
			}
		})
	}
}