package domain

import (
	"testing"
	"time"
)

func TestCredentialProfile_Validate(t *testing.T) {
	validProfile := &CredentialProfile{
		ID:       "test-id",
		Name:     "Test Profile",
		JiraURL:  "https://company.atlassian.net",
		Username: "user@company.com",
		KeychainRef: CredentialRef{
			KeychainID: "keychain-id",
			ServiceID:  "service-id",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name     string
		profile  *CredentialProfile
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "valid profile",
			profile: validProfile,
			wantErr: false,
		},
		{
			name: "empty ID",
			profile: &CredentialProfile{
				Name:     "Test Profile",
				JiraURL:  "https://company.atlassian.net",
				Username: "user@company.com",
				KeychainRef: CredentialRef{
					KeychainID: "keychain-id",
					ServiceID:  "service-id",
				},
			},
			wantErr:  true,
			errorMsg: "credential profile ID cannot be empty",
		},
		{
			name: "empty name",
			profile: &CredentialProfile{
				ID:       "test-id",
				JiraURL:  "https://company.atlassian.net",
				Username: "user@company.com",
				KeychainRef: CredentialRef{
					KeychainID: "keychain-id",
					ServiceID:  "service-id",
				},
			},
			wantErr:  true,
			errorMsg: "name cannot be empty",
		},
		{
			name: "invalid name characters",
			profile: &CredentialProfile{
				ID:       "test-id",
				Name:     "Test@Profile#Invalid",
				JiraURL:  "https://company.atlassian.net",
				Username: "user@company.com",
				KeychainRef: CredentialRef{
					KeychainID: "keychain-id",
					ServiceID:  "service-id",
				},
			},
			wantErr:  true,
			errorMsg: "name must contain only alphanumeric characters, hyphens, underscores, and spaces",
		},
		{
			name: "name too long",
			profile: &CredentialProfile{
				ID:       "test-id",
				Name:     "ThisIsAVeryLongCredentialProfileNameThatExceedsTheMaximumLengthLimit",
				JiraURL:  "https://company.atlassian.net",
				Username: "user@company.com",
				KeychainRef: CredentialRef{
					KeychainID: "keychain-id",
					ServiceID:  "service-id",
				},
			},
			wantErr:  true,
			errorMsg: "name must be 64 characters or less",
		},
		{
			name: "empty jira url",
			profile: &CredentialProfile{
				ID:       "test-id",
				Name:     "Test Profile",
				Username: "user@company.com",
				KeychainRef: CredentialRef{
					KeychainID: "keychain-id",
					ServiceID:  "service-id",
				},
			},
			wantErr:  true,
			errorMsg: "jira_url cannot be empty",
		},
		{
			name: "empty username",
			profile: &CredentialProfile{
				ID:      "test-id",
				Name:    "Test Profile",
				JiraURL: "https://company.atlassian.net",
				KeychainRef: CredentialRef{
					KeychainID: "keychain-id",
					ServiceID:  "service-id",
				},
			},
			wantErr:  true,
			errorMsg: "username cannot be empty",
		},
		{
			name: "empty keychain ID",
			profile: &CredentialProfile{
				ID:       "test-id",
				Name:     "Test Profile",
				JiraURL:  "https://company.atlassian.net",
				Username: "user@company.com",
				KeychainRef: CredentialRef{
					ServiceID: "service-id",
				},
			},
			wantErr:  true,
			errorMsg: "keychain_ref cannot be empty",
		},
		{
			name: "empty service ID",
			profile: &CredentialProfile{
				ID:       "test-id",
				Name:     "Test Profile",
				JiraURL:  "https://company.atlassian.net",
				Username: "user@company.com",
				KeychainRef: CredentialRef{
					KeychainID: "keychain-id",
				},
			},
			wantErr:  true,
			errorMsg: "keychain_ref cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.profile.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("Validate() error = %v, want %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestValidateCredentialProfileInput(t *testing.T) {
	validInput := CredentialProfileInput{
		Name:     "Test Profile",
		JiraURL:  "https://company.atlassian.net",
		Username: "user@company.com",
		APIToken: "test-token",
	}

	tests := []struct {
		name     string
		input    CredentialProfileInput
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "valid input",
			input:   validInput,
			wantErr: false,
		},
		{
			name: "empty name",
			input: CredentialProfileInput{
				JiraURL:  "https://company.atlassian.net",
				Username: "user@company.com",
				APIToken: "test-token",
			},
			wantErr:  true,
			errorMsg: "name is required",
		},
		{
			name: "invalid name characters",
			input: CredentialProfileInput{
				Name:     "Test@Profile#Invalid",
				JiraURL:  "https://company.atlassian.net",
				Username: "user@company.com",
				APIToken: "test-token",
			},
			wantErr:  true,
			errorMsg: "name must contain only alphanumeric characters, hyphens, underscores, and spaces",
		},
		{
			name: "name too long",
			input: CredentialProfileInput{
				Name:     "ThisIsAVeryLongCredentialProfileNameThatExceedsTheMaximumLengthLimit",
				JiraURL:  "https://company.atlassian.net",
				Username: "user@company.com",
				APIToken: "test-token",
			},
			wantErr:  true,
			errorMsg: "name must be 64 characters or less",
		},
		{
			name: "empty jira url",
			input: CredentialProfileInput{
				Name:     "Test Profile",
				Username: "user@company.com",
				APIToken: "test-token",
			},
			wantErr:  true,
			errorMsg: "Jira URL is required",
		},
		{
			name: "empty username",
			input: CredentialProfileInput{
				Name:     "Test Profile",
				JiraURL:  "https://company.atlassian.net",
				APIToken: "test-token",
			},
			wantErr:  true,
			errorMsg: "username is required",
		},
		{
			name: "empty api token",
			input: CredentialProfileInput{
				Name:     "Test Profile",
				JiraURL:  "https://company.atlassian.net",
				Username: "user@company.com",
			},
			wantErr:  true,
			errorMsg: "API token is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCredentialProfileInput(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateCredentialProfileInput() expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("ValidateCredentialProfileInput() error = %v, want %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateCredentialProfileInput() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestValidateCredentialProfileName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "valid name",
			input:   "Test Profile",
			wantErr: false,
		},
		{
			name:    "valid name with hyphens",
			input:   "test-profile",
			wantErr: false,
		},
		{
			name:    "valid name with underscores",
			input:   "test_profile",
			wantErr: false,
		},
		{
			name:    "valid name alphanumeric only",
			input:   "TestProfile123",
			wantErr: false,
		},
		{
			name:     "empty name",
			input:    "",
			wantErr:  true,
			errorMsg: "credential profile name cannot be empty",
		},
		{
			name:     "invalid characters",
			input:    "Test@Profile#Invalid",
			wantErr:  true,
			errorMsg: "credential profile name must contain only alphanumeric characters, hyphens, underscores, and spaces",
		},
		{
			name:     "name too long",
			input:    "ThisIsAVeryLongCredentialProfileNameThatExceedsTheMaximumLengthLimit",
			wantErr:  true,
			errorMsg: "credential profile name must be 64 characters or less",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCredentialProfileName(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateCredentialProfileName() expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("ValidateCredentialProfileName() error = %v, want %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateCredentialProfileName() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestCredentialProfile_Touch(t *testing.T) {
	profile := &CredentialProfile{
		ID:        "test-id",
		Name:      "Test Profile",
		JiraURL:   "https://company.atlassian.net",
		Username:  "user@company.com",
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	oldUpdatedAt := profile.UpdatedAt
	time.Sleep(1 * time.Millisecond) // Ensure time difference

	profile.Touch()

	if !profile.UpdatedAt.After(oldUpdatedAt) {
		t.Errorf("Touch() did not update UpdatedAt field. Old: %v, New: %v", oldUpdatedAt, profile.UpdatedAt)
	}
}
