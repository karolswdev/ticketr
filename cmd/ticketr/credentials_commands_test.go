package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/spf13/cobra"
)

func TestCredentialsProfileCreateValidation(t *testing.T) {
	tests := []struct {
		name        string
		profileName string
		url         string
		username    string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid inputs",
			profileName: "test-profile",
			url:         "https://company.atlassian.net",
			username:    "test@example.com",
			wantErr:     false,
		},
		{
			name:        "invalid URL",
			profileName: "test-profile",
			url:         "not-a-url",
			username:    "test@example.com",
			wantErr:     true,
			errContains: "invalid URL",
		},
		{
			name:        "invalid profile name",
			profileName: "test@profile",
			url:         "https://company.atlassian.net",
			username:    "test@example.com",
			wantErr:     true,
			errContains: "invalid profile name",
		},
		{
			name:        "empty profile name",
			profileName: "",
			url:         "https://company.atlassian.net",
			username:    "test@example.com",
			wantErr:     true,
			errContains: "invalid profile name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test just the validation logic without full command execution
			err := validateCredentialProfileCreateInputs(tt.profileName, tt.url, tt.username)

			// Check error expectation
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error = %v, want error containing %v", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error = %v", err)
				}
			}
		})
	}
}

func TestCredentialsProfileListCommand(t *testing.T) {
	// Test that the list command structure is set up correctly
	t.Run("command_structure", func(t *testing.T) {
		if credentialsProfileListCmd.Use != "list" {
			t.Errorf("credentialsProfileListCmd.Use = %v, want %v", credentialsProfileListCmd.Use, "list")
		}

		if credentialsProfileListCmd.RunE == nil {
			t.Error("credentialsProfileListCmd.RunE should be set")
		}
	})
}

func TestFormatDate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "today",
			input:    now.Add(-1 * time.Hour),
			expected: "today",
		},
		{
			name:     "yesterday",
			input:    now.Add(-25 * time.Hour),
			expected: "1 day ago",
		},
		{
			name:     "two days ago",
			input:    now.Add(-48 * time.Hour),
			expected: "2 days ago",
		},
		{
			name:     "one week ago",
			input:    now.Add(-7 * 24 * time.Hour),
			expected: "1 week ago",
		},
		{
			name:     "two weeks ago",
			input:    now.Add(-14 * 24 * time.Hour),
			expected: "2 weeks ago",
		},
		{
			name:     "one month ago",
			input:    now.Add(-31 * 24 * time.Hour),
			expected: "1 month ago",
		},
		{
			name:     "old date",
			input:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "2020-01-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDate(tt.input)
			if result != tt.expected {
				t.Errorf("formatDate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCredentialsCommandStructure(t *testing.T) {
	// Test that the command structure is properly set up
	if credentialsCmd.Use != "credentials" {
		t.Errorf("credentialsCmd.Use = %v, want %v", credentialsCmd.Use, "credentials")
	}

	if credentialsProfileCmd.Use != "profile" {
		t.Errorf("credentialsProfileCmd.Use = %v, want %v", credentialsProfileCmd.Use, "profile")
	}

	if credentialsProfileCreateCmd.Use != "create <name>" {
		t.Errorf("credentialsProfileCreateCmd.Use = %v, want %v", credentialsProfileCreateCmd.Use, "create <name>")
	}

	if credentialsProfileListCmd.Use != "list" {
		t.Errorf("credentialsProfileListCmd.Use = %v, want %v", credentialsProfileListCmd.Use, "list")
	}

	// Test that required flags are marked
	urlFlag := credentialsProfileCreateCmd.Flags().Lookup("url")
	if urlFlag == nil {
		t.Error("url flag not found")
	}

	usernameFlag := credentialsProfileCreateCmd.Flags().Lookup("username")
	if usernameFlag == nil {
		t.Error("username flag not found")
	}

	// Test that parent-child relationships are set up
	if !hasSubcommand(credentialsCmd, credentialsProfileCmd) {
		t.Error("credentialsProfileCmd is not a subcommand of credentialsCmd")
	}

	if !hasSubcommand(credentialsProfileCmd, credentialsProfileCreateCmd) {
		t.Error("credentialsProfileCreateCmd is not a subcommand of credentialsProfileCmd")
	}

	if !hasSubcommand(credentialsProfileCmd, credentialsProfileListCmd) {
		t.Error("credentialsProfileListCmd is not a subcommand of credentialsProfileCmd")
	}
}

// Helper function to check if a command has a specific subcommand
func hasSubcommand(parent *cobra.Command, child *cobra.Command) bool {
	for _, cmd := range parent.Commands() {
		if cmd == child {
			return true
		}
	}
	return false
}

// Test helper functions for output validation
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestProfileNameValidation(t *testing.T) {
	// Test using domain validation function
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:      "valid name",
			input:     "prod-admin",
			wantError: false,
		},
		{
			name:      "valid with spaces",
			input:     "prod admin",
			wantError: false,
		},
		{
			name:      "valid with underscores",
			input:     "prod_admin",
			wantError: false,
		},
		{
			name:      "empty name",
			input:     "",
			wantError: true,
		},
		{
			name:      "invalid characters",
			input:     "prod@admin",
			wantError: true,
		},
		{
			name:      "too long",
			input:     strings.Repeat("a", 101),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := domain.ValidateCredentialProfileName(tt.input)
			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateCredentialProfileName(%v) expected error but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateCredentialProfileName(%v) unexpected error = %v", tt.input, err)
				}
			}
		})
	}
}
