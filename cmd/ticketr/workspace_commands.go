package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/karolswdev/ticktr/internal/adapters/database"
	"github.com/karolswdev/ticktr/internal/adapters/keychain"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/spf13/cobra"
)

var (
	// Workspace command flags
	wsURL      string
	wsProject  string
	wsUsername string
	wsToken    string
	wsForce    bool
	wsProfile  string

	// workspaceCmd represents the main workspace command
	workspaceCmd = &cobra.Command{
		Use:   "workspace",
		Short: "Manage Jira project workspaces",
		Long: `Manage multiple Jira project workspaces.

Workspaces allow you to work with multiple Jira projects from a single
installation. Each workspace stores its credentials securely in your
OS keychain.

Examples:
  ticketr workspace create my-project --url https://company.atlassian.net --project PROJ --username user@example.com --token abc123
  ticketr workspace create backend --profile prod-admin --project BACK
  ticketr workspace list
  ticketr workspace switch my-project
  ticketr workspace current
  ticketr workspace delete old-project
  ticketr workspace set-default my-project`,
	}

	// workspaceCreateCmd creates a new workspace
	workspaceCreateCmd = &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new workspace",
		Long: `Create a new workspace with Jira credentials.

You can either provide credentials directly or use an existing credential profile:

Option 1: Direct credentials (all flags required)
  --url, --project, --username, --token

Option 2: Use credential profile (--profile and --project required)
  --profile <profile-name> --project <project-key>

Credentials are stored securely in your OS keychain:
- macOS: Keychain Access
- Windows: Credential Manager
- Linux: Secret Service (GNOME Keyring/KWallet)

Examples:
  # Direct credentials
  ticketr workspace create my-project \
    --url https://company.atlassian.net \
    --project PROJ \
    --username user@example.com \
    --token abc123

  # Using credential profile
  ticketr workspace create backend \
    --profile prod-admin \
    --project BACK`,
		Args: cobra.ExactArgs(1),
		RunE: runWorkspaceCreate,
	}

	// workspaceListCmd lists all workspaces
	workspaceListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all workspaces",
		Long: `List all configured workspaces.

Shows workspace name, project key, Jira URL, default status, and last used time.`,
		RunE: runWorkspaceList,
	}

	// workspaceSwitchCmd switches to a different workspace
	workspaceSwitchCmd = &cobra.Command{
		Use:   "switch <name>",
		Short: "Switch to a different workspace",
		Long: `Switch the current active workspace.

This updates the last used timestamp and makes the workspace
the active one for subsequent operations.`,
		Args: cobra.ExactArgs(1),
		RunE: runWorkspaceSwitch,
	}

	// workspaceCurrentCmd shows the current workspace
	workspaceCurrentCmd = &cobra.Command{
		Use:   "current",
		Short: "Show the current workspace",
		Long: `Display the currently active workspace.

If no workspace is active, shows the default workspace.`,
		RunE: runWorkspaceCurrent,
	}

	// workspaceDeleteCmd deletes a workspace
	workspaceDeleteCmd = &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a workspace",
		Long: `Delete a workspace and its credentials.

This removes the workspace from the database and deletes its
credentials from the OS keychain. This operation cannot be undone.

Use --force to skip the confirmation prompt.`,
		Args: cobra.ExactArgs(1),
		RunE: runWorkspaceDelete,
	}

	// workspaceSetDefaultCmd sets the default workspace
	workspaceSetDefaultCmd = &cobra.Command{
		Use:   "set-default <name>",
		Short: "Set the default workspace",
		Long: `Mark a workspace as the default.

The default workspace is used when no explicit workspace is selected.`,
		Args: cobra.ExactArgs(1),
		RunE: runWorkspaceSetDefault,
	}
)

func init() {
	// Add subcommands to workspace
	workspaceCmd.AddCommand(workspaceCreateCmd)
	workspaceCmd.AddCommand(workspaceListCmd)
	workspaceCmd.AddCommand(workspaceSwitchCmd)
	workspaceCmd.AddCommand(workspaceCurrentCmd)
	workspaceCmd.AddCommand(workspaceDeleteCmd)
	workspaceCmd.AddCommand(workspaceSetDefaultCmd)

	// Flags for create command
	workspaceCreateCmd.Flags().StringVar(&wsURL, "url", "", "Jira URL (e.g., https://company.atlassian.net)")
	workspaceCreateCmd.Flags().StringVar(&wsProject, "project", "", "Jira project key")
	workspaceCreateCmd.Flags().StringVar(&wsUsername, "username", "", "Jira username/email")
	workspaceCreateCmd.Flags().StringVar(&wsToken, "token", "", "Jira API token")
	workspaceCreateCmd.Flags().StringVar(&wsProfile, "profile", "", "Credential profile name to use")

	// Project is always required
	workspaceCreateCmd.MarkFlagRequired("project")

	// Note: We'll handle conditional flag requirements in the RunE function
	// since cobra doesn't support conditional required flags natively

	// Flags for delete command
	workspaceDeleteCmd.Flags().BoolVar(&wsForce, "force", false, "Skip confirmation prompt")
}

// initWorkspaceService creates a new WorkspaceService instance
func initWorkspaceService() (*services.WorkspaceService, error) {
	// Get PathResolver singleton
	pathResolver, err := services.GetPathResolver()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize path resolver: %w", err)
	}

	// Check if migration is needed and require explicit user consent
	legacyDBPath := filepath.Join(".ticketr", "ticketr.db")
	migrationStatePath := pathResolver.DataFile(".migration_complete")

	// Check if legacy data exists
	if _, err := os.Stat(legacyDBPath); err == nil {
		// Legacy data exists - check if migration completed
		if _, err := os.Stat(migrationStatePath); os.IsNotExist(err) {
			// Migration needed but not completed - prompt user
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "⚠️  Ticketr v3.0 Migration Required")
			fmt.Fprintln(os.Stderr, "   Detected legacy v2.x data in .ticketr/")
			fmt.Fprintln(os.Stderr, "   ")
			fmt.Fprintln(os.Stderr, "   v3.0 uses platform-standard global directories:")
			fmt.Fprintln(os.Stderr, "   - Linux:   ~/.local/share/ticketr/")
			fmt.Fprintln(os.Stderr, "   - macOS:   ~/Library/Application Support/ticketr/")
			fmt.Fprintln(os.Stderr, "   - Windows: %LOCALAPPDATA%\\ticketr\\")
			fmt.Fprintln(os.Stderr, "   ")
			fmt.Fprintln(os.Stderr, "   Run 'ticketr migrate-paths' to migrate your data")
			fmt.Fprintln(os.Stderr, "   See docs/v3-MIGRATION-GUIDE.md for details")
			fmt.Fprintln(os.Stderr, "")
			return nil, fmt.Errorf("migration required - run 'ticketr migrate-paths'")
		}
	}

	// Initialize database with PathResolver
	adapter, err := database.NewSQLiteAdapter(pathResolver)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create repositories
	repo := database.NewWorkspaceRepository(adapter.DB())
	credentialRepo := database.NewCredentialProfileRepository(adapter.DB())

	// Try to create keychain store first, fall back to file store if keychain unavailable
	var credStore ports.CredentialStore
	keychainStore := keychain.NewKeychainStore()

	// Test if keychain is accessible by trying to list credentials
	// This will fail if the keyring is locked on Linux
	_, err = keychainStore.List()
	if err != nil {
		// Keychain unavailable (locked or not accessible)
		// Fall back to file-based storage
		fmt.Fprintln(os.Stderr, "⚠️  OS keychain unavailable (keyring may be locked)")
		fmt.Fprintln(os.Stderr, "   Using file-based credential storage: ~/.config/ticketr/credentials/")
		fmt.Fprintln(os.Stderr, "   To use OS keychain instead, unlock your keyring with:")
		fmt.Fprintln(os.Stderr, "     gnome-keyring-daemon --unlock")
		fmt.Fprintln(os.Stderr, "")

		fileStore, err := keychain.NewFileStore()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize credential storage: %w", err)
		}
		credStore = fileStore
	} else {
		// Keychain is accessible, use it
		credStore = keychainStore
	}

	// Create workspace service
	svc := services.NewWorkspaceService(repo, credentialRepo, credStore)

	return svc, nil
}

// runWorkspaceCreate handles the workspace create command
func runWorkspaceCreate(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Validate workspace name
	if err := domain.ValidateWorkspaceName(name); err != nil {
		return fmt.Errorf("invalid workspace name: %w", err)
	}

	// Initialize service
	svc, err := initWorkspaceService()
	if err != nil {
		return err
	}

	// Check if using profile-based creation
	if wsProfile != "" {
		// Profile-based creation
		return createWorkspaceWithProfile(svc, name, wsProfile, wsProject)
	}

	// Direct credential creation - validate all required flags are present
	if wsURL == "" || wsUsername == "" || wsToken == "" {
		return fmt.Errorf("when not using --profile, all of --url, --username, and --token are required")
	}

	// Create workspace config
	config := domain.WorkspaceConfig{
		JiraURL:    wsURL,
		ProjectKey: wsProject,
		Username:   wsUsername,
		APIToken:   wsToken,
	}

	// Validate configuration
	if err := domain.ValidateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Create workspace
	if err := svc.Create(name, config); err != nil {
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	fmt.Printf("\n✓ Workspace '%s' created successfully\n", name)
	fmt.Printf("  Project: %s\n", config.ProjectKey)
	fmt.Printf("  URL: %s\n", config.JiraURL)
	fmt.Printf("  Username: %s\n", config.Username)
	fmt.Println("\nCredentials stored securely in OS keychain")

	// Check if this is the first workspace
	workspaces, _ := svc.List()
	if len(workspaces) == 1 {
		fmt.Println("\nThis is your first workspace and has been set as the default.")
	}

	return nil
}

// createWorkspaceWithProfile creates a workspace using an existing credential profile
func createWorkspaceWithProfile(svc *services.WorkspaceService, name, profileName, projectKey string) error {
	// Get the credential profile to validate it exists and get its ID
	profile, err := svc.GetProfile(profileName)
	if err != nil {
		return fmt.Errorf("failed to get credential profile '%s': %w", profileName, err)
	}

	// Create workspace with profile
	if err := svc.CreateWithProfile(name, projectKey, profile.ID); err != nil {
		return fmt.Errorf("failed to create workspace with profile: %w", err)
	}

	fmt.Printf("\n✓ Workspace '%s' created using profile '%s'\n", name, profileName)
	fmt.Printf("  Project: %s\n", projectKey)
	fmt.Printf("  URL: %s\n", profile.JiraURL)
	fmt.Printf("  Username: %s\n", profile.Username)
	fmt.Println("\nUsing existing credentials from OS keychain")

	// Check if this is the first workspace
	workspaces, _ := svc.List()
	if len(workspaces) == 1 {
		fmt.Println("\nThis is your first workspace and has been set as the default.")
	}

	return nil
}

// runWorkspaceList handles the workspace list command
func runWorkspaceList(cmd *cobra.Command, args []string) error {
	// Initialize service
	svc, err := initWorkspaceService()
	if err != nil {
		return err
	}

	// Get all workspaces
	workspaces, err := svc.List()
	if err != nil {
		return fmt.Errorf("failed to list workspaces: %w", err)
	}

	if len(workspaces) == 0 {
		fmt.Println("No workspaces found.")
		fmt.Println("\nCreate a workspace with:")
		fmt.Println("  ticketr workspace create <name> --url <jira-url> --project <key> --username <email> --token <api-token>")
		return nil
	}

	// Print workspaces in table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tPROJECT\tURL\tDEFAULT\tLAST USED")
	fmt.Fprintln(w, "----\t-------\t---\t-------\t---------")

	for _, ws := range workspaces {
		// Format default indicator
		defaultStr := ""
		if ws.IsDefault {
			defaultStr = "*"
		}

		// Format last used
		lastUsed := "never"
		if !ws.LastUsed.IsZero() {
			lastUsed = formatTimeAgo(ws.LastUsed)
		}

		// Truncate URL if too long
		url := ws.JiraURL
		if len(url) > 40 {
			url = url[:37] + "..."
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			ws.Name,
			ws.ProjectKey,
			url,
			defaultStr,
			lastUsed,
		)
	}

	w.Flush()

	// Print legend
	if hasDefault(workspaces) {
		fmt.Println("\n* = default workspace")
	}

	return nil
}

// runWorkspaceSwitch handles the workspace switch command
func runWorkspaceSwitch(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Initialize service
	svc, err := initWorkspaceService()
	if err != nil {
		return err
	}

	// Switch to workspace
	if err := svc.Switch(name); err != nil {
		return fmt.Errorf("failed to switch workspace: %w", err)
	}

	fmt.Printf("Switched to workspace '%s'\n", name)

	return nil
}

// runWorkspaceCurrent handles the workspace current command
func runWorkspaceCurrent(cmd *cobra.Command, args []string) error {
	// Initialize service
	svc, err := initWorkspaceService()
	if err != nil {
		return err
	}

	// Get current workspace
	ws, err := svc.Current()
	if err != nil {
		fmt.Println("No workspace selected and no default workspace configured.")
		fmt.Println("\nCreate a workspace with:")
		fmt.Println("  ticketr workspace create <name> --url <jira-url> --project <key> --username <email> --token <api-token>")
		return nil
	}

	fmt.Printf("Current workspace: %s\n", ws.Name)
	fmt.Printf("  Project: %s\n", ws.ProjectKey)
	fmt.Printf("  URL: %s\n", ws.JiraURL)

	if ws.IsDefault {
		fmt.Println("  Default: yes")
	}

	if !ws.LastUsed.IsZero() {
		fmt.Printf("  Last used: %s\n", formatTimeAgo(ws.LastUsed))
	}

	return nil
}

// runWorkspaceDelete handles the workspace delete command
func runWorkspaceDelete(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Initialize service
	svc, err := initWorkspaceService()
	if err != nil {
		return err
	}

	// Confirm deletion unless --force is set
	if !wsForce {
		fmt.Printf("Are you sure you want to delete workspace '%s'? [y/N]: ", name)
		var response string
		fmt.Scanln(&response)
		response = strings.ToLower(strings.TrimSpace(response))

		if response != "y" && response != "yes" {
			fmt.Println("Deletion cancelled.")
			return nil
		}
	}

	// Delete workspace
	if err := svc.Delete(name); err != nil {
		return fmt.Errorf("failed to delete workspace: %w", err)
	}

	fmt.Printf("Workspace '%s' deleted successfully\n", name)
	fmt.Println("Credentials removed from OS keychain")

	return nil
}

// runWorkspaceSetDefault handles the workspace set-default command
func runWorkspaceSetDefault(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Initialize service
	svc, err := initWorkspaceService()
	if err != nil {
		return err
	}

	// Set as default
	if err := svc.SetDefault(name); err != nil {
		return fmt.Errorf("failed to set default workspace: %w", err)
	}

	fmt.Printf("Workspace '%s' is now the default\n", name)

	return nil
}

// Helper functions

// formatTimeAgo formats a time as a human-readable "time ago" string
func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if duration < 7*24*time.Hour {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if duration < 30*24*time.Hour {
		weeks := int(duration.Hours() / 24 / 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	} else if duration < 365*24*time.Hour {
		months := int(duration.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}

	years := int(duration.Hours() / 24 / 365)
	if years == 1 {
		return "1 year ago"
	}
	return fmt.Sprintf("%d years ago", years)
}

// hasDefault checks if any workspace is marked as default
func hasDefault(workspaces []domain.Workspace) bool {
	for _, ws := range workspaces {
		if ws.IsDefault {
			return true
		}
	}
	return false
}
