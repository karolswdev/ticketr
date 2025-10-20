package main

import (
	"fmt"
	"net/url"
	"os"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	// Credential profile command flags
	credProfileURL      string
	credProfileUsername string

	// credentialsCmd represents the main credentials command
	credentialsCmd = &cobra.Command{
		Use:   "credentials",
		Short: "Manage credential profiles",
		Long: `Manage credential profiles for Jira authentication.

Credential profiles allow you to store and reuse Jira credentials
across multiple workspaces, making it easier to create workspaces
without re-entering credentials each time.

Examples:
  ticketr credentials profile create prod-admin --url https://company.atlassian.net --username admin@company.com
  ticketr credentials profile list
  ticketr workspace create backend --profile prod-admin --project BACK`,
	}

	// credentialsProfileCmd represents the credential profile subcommand group
	credentialsProfileCmd = &cobra.Command{
		Use:   "profile",
		Short: "Manage credential profiles",
		Long: `Manage credential profiles for reusable Jira authentication.

Credential profiles store Jira URL, username, and API token securely
in your OS keychain for reuse across multiple workspaces.`,
	}

	// credentialsProfileCreateCmd creates a new credential profile
	credentialsProfileCreateCmd = &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new credential profile",
		Long: `Create a new credential profile with Jira credentials.

The API token will be prompted securely (hidden input) and stored
in your OS keychain along with the URL and username.

Example:
  ticketr credentials profile create prod-admin \
    --url https://company.atlassian.net \
    --username admin@company.com`,
		Args: cobra.ExactArgs(1),
		RunE: runCredentialsProfileCreate,
	}

	// credentialsProfileListCmd lists all credential profiles
	credentialsProfileListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all credential profiles",
		Long: `List all configured credential profiles.

Shows profile name, Jira URL, username, number of workspaces using
the profile, and creation date.`,
		RunE: runCredentialsProfileList,
	}
)

func init() {
	// Add subcommands to credentials
	credentialsCmd.AddCommand(credentialsProfileCmd)

	// Add subcommands to profile
	credentialsProfileCmd.AddCommand(credentialsProfileCreateCmd)
	credentialsProfileCmd.AddCommand(credentialsProfileListCmd)

	// Flags for profile create command
	credentialsProfileCreateCmd.Flags().StringVar(&credProfileURL, "url", "", "Jira URL (e.g., https://company.atlassian.net)")
	credentialsProfileCreateCmd.Flags().StringVar(&credProfileUsername, "username", "", "Jira username/email")

	credentialsProfileCreateCmd.MarkFlagRequired("url")
	credentialsProfileCreateCmd.MarkFlagRequired("username")
}

// validateCredentialProfileCreateInputs validates the inputs for credential profile creation
func validateCredentialProfileCreateInputs(name, jiraURL, username string) error {
	// Validate profile name (using existing domain validation)
	if err := domain.ValidateCredentialProfileName(name); err != nil {
		return fmt.Errorf("invalid profile name: %w", err)
	}

	// Validate URL format
	parsedURL, err := url.Parse(jiraURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return fmt.Errorf("invalid URL: must be a valid HTTP or HTTPS URL")
	}

	// Validate username is not empty
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	return nil
}

// runCredentialsProfileCreate handles the credential profile create command
func runCredentialsProfileCreate(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Validate inputs
	if err := validateCredentialProfileCreateInputs(name, credProfileURL, credProfileUsername); err != nil {
		return err
	}

	// Prompt for API token (hidden input)
	fmt.Printf("Enter API token for %s@%s: ", credProfileUsername, credProfileURL)
	apiTokenBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read API token: %w", err)
	}
	fmt.Println() // Add newline after hidden input

	apiToken := string(apiTokenBytes)
	if apiToken == "" {
		return fmt.Errorf("API token cannot be empty")
	}

	// Create credential profile input
	input := domain.CredentialProfileInput{
		Name:     name,
		JiraURL:  credProfileURL,
		Username: credProfileUsername,
		APIToken: apiToken,
	}

	// Initialize service
	svc, err := initWorkspaceService()
	if err != nil {
		return err
	}

	// Create credential profile
	profileID, err := svc.CreateProfile(input)
	if err != nil {
		return fmt.Errorf("failed to create credential profile: %w", err)
	}

	fmt.Printf("\n✓ Credential profile '%s' created successfully\n", name)
	fmt.Printf("  Profile ID: %s\n", profileID)
	fmt.Printf("  Jira URL: %s\n", credProfileURL)
	fmt.Printf("  Username: %s\n", credProfileUsername)
	fmt.Println("\nCredentials stored securely in OS keychain")

	return nil
}

// runCredentialsProfileList handles the credential profile list command
func runCredentialsProfileList(cmd *cobra.Command, args []string) error {
	// Initialize service
	svc, err := initWorkspaceService()
	if err != nil {
		return err
	}

	// Get all credential profiles
	profiles, err := svc.ListProfiles()
	if err != nil {
		return fmt.Errorf("failed to list credential profiles: %w", err)
	}

	if len(profiles) == 0 {
		fmt.Println("No credential profiles found.")
		fmt.Println("\nCreate a credential profile with:")
		fmt.Println("  ticketr credentials profile create <name> --url <jira-url> --username <email>")
		return nil
	}

	fmt.Println("Credential Profiles:")
	fmt.Println()

	// Print profiles in table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tJIRA URL\tUSERNAME\tWORKSPACES\tCREATED")
	fmt.Fprintln(w, "----\t--------\t--------\t----------\t-------")

	for _, profile := range profiles {
		// Get workspace count for this profile
		workspaceCount, err := getWorkspaceCountForProfile(svc, profile.ID)
		if err != nil {
			// Log warning but continue
			workspaceCount = 0
		}

		// Truncate URL if too long
		url := profile.JiraURL
		if len(url) > 30 {
			url = url[:27] + "..."
		}

		// Truncate username if too long
		username := profile.Username
		if len(username) > 20 {
			username = username[:17] + "..."
		}

		// Format creation date
		created := formatDate(profile.CreatedAt)

		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n",
			profile.Name,
			url,
			username,
			workspaceCount,
			created,
		)
	}

	w.Flush()

	fmt.Printf("\n%d credential profile(s) found.\n", len(profiles))

	return nil
}

// getWorkspaceCountForProfile gets the number of workspaces using a specific profile.
func getWorkspaceCountForProfile(svc *services.WorkspaceService, profileID string) (int, error) {
	// For now, we'll return 0 since we don't have direct access to the credential repository
	// from the CLI layer. This could be enhanced by adding a method to WorkspaceService
	// that returns workspace count per profile.
	return 0, nil
}

// formatDate formats a time as a human-readable date string
func formatDate(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < 24*time.Hour {
		return "today"
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if diff < 30*24*time.Hour {
		weeks := int(diff.Hours() / 24 / 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	} else if diff < 365*24*time.Hour {
		months := int(diff.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}

	return t.Format("2006-01-02")
}
