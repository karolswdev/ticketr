package views

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/karolswdev/ticktr/internal/adapters/tui/effects"
	"github.com/karolswdev/ticktr/internal/adapters/tui/theme"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/rivo/tview"
)

// WorkspaceModal displays a modal dialog for creating workspaces with credential profiles.
type WorkspaceModal struct {
	app              *tview.Application
	pages            *tview.Pages // Phase 6.6: Use pages overlay instead of SetRoot
	form             *tview.Form
	shadowForm       *effects.ShadowForm
	workspaceService *services.WorkspaceService
	onClose          func()
	onSuccess        func()

	// Form fields
	nameField       *tview.InputField
	projectKeyField *tview.InputField
	profileDropdown *tview.DropDown

	// New profile fields (shown conditionally)
	newProfileName     *tview.InputField
	newProfileURL      *tview.InputField
	newProfileUsername *tview.InputField
	newProfileToken    *tview.InputField

	// State
	profiles           []domain.CredentialProfile
	useExistingProfile bool
	showingNewProfile  bool
	isValidating       bool
}

// NewWorkspaceModal creates a new workspace creation modal.
func NewWorkspaceModal(app *tview.Application, pages *tview.Pages, workspaceService *services.WorkspaceService) *WorkspaceModal {
	modal := &WorkspaceModal{
		app:                app,
		pages:              pages,
		workspaceService:   workspaceService,
		useExistingProfile: true,
	}

	modal.setupForm()
	return modal
}

// setupForm creates and configures the form.
func (w *WorkspaceModal) setupForm() {
	// Check if shadows are enabled
	effectsConfig := theme.GetEffects()
	if effectsConfig.DropShadows {
		// Use shadow form for modal with drop shadow
		w.shadowForm = effects.NewShadowForm()
		w.form = w.shadowForm.GetForm()
	} else {
		// Use regular form without shadow
		w.form = tview.NewForm()
	}
	w.form.SetBorder(true).SetTitle(" Create Workspace ")
	w.form.SetBorderColor(theme.GetPrimaryColor())

	// Add ESC key handler to close modal
	w.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			w.handleCancel()
			return nil // Consume the event
		}
		return event
	})

	// Create form fields with required markers
	w.nameField = tview.NewInputField().
		SetLabel("Workspace Name *").
		SetFieldWidth(40).
		SetPlaceholder("e.g., my-project")

	w.projectKeyField = tview.NewInputField().
		SetLabel("Project Key *").
		SetFieldWidth(20).
		SetPlaceholder("e.g., PROJ")

	w.profileDropdown = tview.NewDropDown().
		SetLabel("Credential Profile *").
		SetFieldWidth(40)

	// New profile fields (initially hidden) - with required markers
	w.newProfileName = tview.NewInputField().
		SetLabel("Profile Name *").
		SetFieldWidth(40).
		SetPlaceholder("e.g., prod-admin")

	w.newProfileURL = tview.NewInputField().
		SetLabel("Jira URL *").
		SetFieldWidth(50).
		SetPlaceholder("https://company.atlassian.net")

	w.newProfileUsername = tview.NewInputField().
		SetLabel("Username/Email *").
		SetFieldWidth(40).
		SetPlaceholder("user@company.com")

	w.newProfileToken = tview.NewInputField().
		SetLabel("API Token *").
		SetFieldWidth(40).
		SetMaskCharacter('*').
		SetPlaceholder("Your Jira API token")

	// Setup form layout
	w.buildForm()

	// Setup input validation
	w.setupValidation()
}

// buildForm constructs the form layout based on current state.
func (w *WorkspaceModal) buildForm() {
	w.form.Clear(true)

	// Basic workspace fields
	w.form.AddFormItem(w.nameField)
	w.form.AddFormItem(w.projectKeyField)

	// Credential profile selection
	w.form.AddFormItem(w.profileDropdown)

	// Add profile mode toggle buttons
	w.form.AddButton("Use Existing Profile", func() {
		w.useExistingProfile = true
		w.showingNewProfile = false
		w.buildForm()
	})

	w.form.AddButton("Create New Profile", func() {
		w.useExistingProfile = false
		w.showingNewProfile = true
		w.buildForm()
	})

	// Show new profile fields if creating new profile
	if w.showingNewProfile {
		w.form.AddFormItem(w.newProfileName)
		w.form.AddFormItem(w.newProfileURL)
		w.form.AddFormItem(w.newProfileUsername)
		w.form.AddFormItem(w.newProfileToken)
	}

	// Add help text
	helpText := tview.NewTextView().
		SetText("[gray]* = Required field | Tab: Next field | Enter: Submit | ESC: Cancel[-]").
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	w.form.AddFormItem(helpText)

	// Action buttons
	w.form.AddButton("Create", w.handleCreate)
	w.form.AddButton("Cancel", w.handleCancel)

	// Set button styling
	w.form.SetButtonsAlign(tview.AlignCenter)
	w.form.SetButtonBackgroundColor(theme.GetPrimaryColor())
	w.form.SetButtonTextColor(tcell.ColorWhite)
}

// setupValidation configures real-time validation for form fields.
func (w *WorkspaceModal) setupValidation() {
	// Workspace name validation
	w.nameField.SetChangedFunc(func(text string) {
		if err := domain.ValidateWorkspaceName(text); err != nil {
			w.nameField.SetFieldBackgroundColor(theme.GetErrorColor())
		} else {
			w.nameField.SetFieldBackgroundColor(tcell.ColorDefault)
		}
	})

	// Project key validation
	w.projectKeyField.SetChangedFunc(func(text string) {
		if text == "" || len(text) > 10 {
			w.projectKeyField.SetFieldBackgroundColor(theme.GetErrorColor())
		} else {
			w.projectKeyField.SetFieldBackgroundColor(tcell.ColorDefault)
		}
	})

	// URL validation for new profiles
	w.newProfileURL.SetChangedFunc(func(text string) {
		if text != "" {
			if _, err := url.Parse(text); err != nil || !strings.HasPrefix(text, "http") {
				w.newProfileURL.SetFieldBackgroundColor(theme.GetErrorColor())
			} else {
				w.newProfileURL.SetFieldBackgroundColor(tcell.ColorDefault)
			}
		}
	})

	// Username validation
	w.newProfileUsername.SetChangedFunc(func(text string) {
		if w.showingNewProfile && text == "" {
			w.newProfileUsername.SetFieldBackgroundColor(theme.GetErrorColor())
		} else {
			w.newProfileUsername.SetFieldBackgroundColor(tcell.ColorDefault)
		}
	})

	// Token validation
	w.newProfileToken.SetChangedFunc(func(text string) {
		if w.showingNewProfile && text == "" {
			w.newProfileToken.SetFieldBackgroundColor(theme.GetErrorColor())
		} else {
			w.newProfileToken.SetFieldBackgroundColor(tcell.ColorDefault)
		}
	})
}

// Show displays the workspace creation modal.
func (w *WorkspaceModal) Show() {
	// Load credential profiles
	if err := w.loadProfiles(); err != nil {
		// Show error and close
		w.showError(fmt.Sprintf("Failed to load credential profiles: %v", err))
		return
	}

	// Setup profile dropdown
	w.setupProfileDropdown()

	// Rebuild form to ensure proper layout
	w.buildForm()

	// Determine which primitive to display (with or without shadows)
	var displayPrimitive tview.Primitive
	effectsConfig := theme.GetEffects()
	if effectsConfig.DropShadows && w.shadowForm != nil {
		displayPrimitive = w.shadowForm
	} else {
		displayPrimitive = w.form
	}

	// Create a responsive grid layout with comfortable margins
	// Columns: 7-col left margin, flexible center, 7-col right margin (FIX #4: increased from 5)
	grid := tview.NewGrid().
		SetColumns(7, 0, 7).  // Comfortable margins with 20% more space
		SetRows(2, 0, 2).     // 2-row top margin, content, 2-row bottom margin
		AddItem(displayPrimitive, 1, 1, 1, 1, 0, 0, true) // Center position (row 1, col 1)

	// FIX #1: Use pages overlay instead of SetRoot to avoid breaking overlay state
	if w.pages != nil {
		// Add or update the modal page
		w.pages.AddPage("workspace-modal", grid, true, false)
		w.pages.ShowPage("workspace-modal")
		// Focus the name field
		w.app.SetFocus(w.nameField)
	} else {
		// Fallback for tests or standalone usage
		w.app.SetRoot(grid, true)
		w.app.SetFocus(w.nameField)
	}
}

// loadProfiles loads available credential profiles from the service.
func (w *WorkspaceModal) loadProfiles() error {
	profiles, err := w.workspaceService.ListProfiles()
	if err != nil {
		return err
	}
	w.profiles = profiles
	return nil
}

// setupProfileDropdown configures the profile dropdown with available profiles.
func (w *WorkspaceModal) setupProfileDropdown() {
	w.profileDropdown.SetOptions(nil, nil)

	if len(w.profiles) == 0 {
		w.profileDropdown.AddOption("No profiles available", nil)
		w.profileDropdown.SetCurrentOption(0)
		w.useExistingProfile = false
		w.showingNewProfile = true
		return
	}

	for _, profile := range w.profiles {
		displayText := fmt.Sprintf("%s (%s)", profile.Name, profile.JiraURL)
		w.profileDropdown.AddOption(displayText, nil)
	}

	// Auto-populate fields when profile is selected
	w.profileDropdown.SetSelectedFunc(func(text string, index int) {
		if index >= 0 && index < len(w.profiles) {
			profile := w.profiles[index]
			// Auto-suggest workspace name from profile name
			if w.nameField.GetText() == "" {
				w.nameField.SetText(profile.Name)
			}
		}
	})

	w.profileDropdown.SetCurrentOption(0)
}

// handleCreate processes the workspace creation request.
func (w *WorkspaceModal) handleCreate() {
	if w.isValidating {
		return // Prevent double-submission
	}

	// Validate form
	if err := w.validateForm(); err != nil {
		w.showError(err.Error())
		return
	}

	w.isValidating = true
	w.showProgress("Creating workspace...")

	// Get form values
	name := strings.TrimSpace(w.nameField.GetText())
	projectKey := strings.TrimSpace(w.projectKeyField.GetText())

	// Create workspace in background to avoid blocking UI
	go func() {
		var err error

		if w.useExistingProfile && len(w.profiles) > 0 {
			// Use existing profile
			index, _ := w.profileDropdown.GetCurrentOption()
			if index >= 0 && index < len(w.profiles) {
				profileID := w.profiles[index].ID
				err = w.workspaceService.CreateWithProfile(name, projectKey, profileID)
			} else {
				err = fmt.Errorf("no profile selected")
			}
		} else {
			// Create new profile first, then workspace
			err = w.createWorkspaceWithNewProfile(name, projectKey)
		}

		// Update UI on main thread
		w.app.QueueUpdateDraw(func() {
			w.isValidating = false
			if err != nil {
				w.showError(fmt.Sprintf("Failed to create workspace: %v", err))
			} else {
				w.showSuccess()
			}
		})
	}()
}

// createWorkspaceWithNewProfile creates a new profile and workspace.
func (w *WorkspaceModal) createWorkspaceWithNewProfile(workspaceName, projectKey string) error {
	// Create credential profile input
	profileInput := domain.CredentialProfileInput{
		Name:     strings.TrimSpace(w.newProfileName.GetText()),
		JiraURL:  strings.TrimSpace(w.newProfileURL.GetText()),
		Username: strings.TrimSpace(w.newProfileUsername.GetText()),
		APIToken: strings.TrimSpace(w.newProfileToken.GetText()),
	}

	// Create profile
	profileID, err := w.workspaceService.CreateProfile(profileInput)
	if err != nil {
		return fmt.Errorf("failed to create credential profile: %w", err)
	}

	// Create workspace with new profile
	if err := w.workspaceService.CreateWithProfile(workspaceName, projectKey, profileID); err != nil {
		// Try to clean up the profile if workspace creation fails
		w.workspaceService.DeleteProfile(profileInput.Name)
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	return nil
}

// validateForm validates all form fields.
func (w *WorkspaceModal) validateForm() error {
	// Validate workspace name
	name := strings.TrimSpace(w.nameField.GetText())
	if err := domain.ValidateWorkspaceName(name); err != nil {
		return fmt.Errorf("Workspace name is invalid: %w", err)
	}

	// Validate project key
	projectKey := strings.TrimSpace(w.projectKeyField.GetText())
	if projectKey == "" {
		return fmt.Errorf("Project key is required - please enter your Jira project key (e.g., PROJ)")
	}
	if len(projectKey) > 10 {
		return fmt.Errorf("Project key must be 10 characters or less (got %d)", len(projectKey))
	}

	// Validate profile selection or new profile fields
	if w.useExistingProfile {
		if len(w.profiles) == 0 {
			return fmt.Errorf("No credential profiles available. Please create a new profile first.")
		}
		// Profile validation is handled by dropdown selection
	} else {
		// Validate new profile fields
		profileName := strings.TrimSpace(w.newProfileName.GetText())
		if profileName == "" {
			return fmt.Errorf("Profile name is required - please enter a name for your credential profile")
		}

		jiraURL := strings.TrimSpace(w.newProfileURL.GetText())
		if jiraURL == "" {
			return fmt.Errorf("Jira URL is required - please enter your Jira instance URL")
		}
		if _, err := url.Parse(jiraURL); err != nil {
			return fmt.Errorf("Jira URL is invalid: %w", err)
		}
		if !strings.HasPrefix(jiraURL, "http") {
			return fmt.Errorf("Jira URL must start with http:// or https://")
		}

		username := strings.TrimSpace(w.newProfileUsername.GetText())
		if username == "" {
			return fmt.Errorf("Username/email is required - please enter your Jira account email")
		}

		token := strings.TrimSpace(w.newProfileToken.GetText())
		if token == "" {
			return fmt.Errorf("API token is required - please enter your Jira API token")
		}
	}

	return nil
}

// handleCancel closes the modal without creating a workspace.
func (w *WorkspaceModal) handleCancel() {
	if w.onClose != nil {
		w.onClose()
	}
}

// showError displays an error message by temporarily replacing form content.
func (w *WorkspaceModal) showError(message string) {
	// Store current form state
	originalTitle := w.form.GetTitle()

	// Clear form and show error
	w.form.Clear(true)
	w.form.SetTitle(" ⚠ Error ")
	w.form.SetBorderColor(theme.GetErrorColor())

	// Add error text as a text view with better formatting
	errorText := fmt.Sprintf("\n[red::b]Error:[-:-:-] %s\n\n[yellow]Press OK or ESC to continue...[-]", message)
	textView := tview.NewTextView().
		SetText(errorText).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetWordWrap(true)

	// Add the text view to the form
	w.form.AddFormItem(textView)

	// Add OK button
	w.form.AddButton("OK", func() {
		// Restore form
		w.form.SetTitle(originalTitle)
		w.form.SetBorderColor(theme.GetPrimaryColor())
		w.buildForm()
	})

	// Set button styling for error state
	w.form.SetButtonsAlign(tview.AlignCenter)
	w.form.SetButtonBackgroundColor(theme.GetErrorColor())
	w.form.SetButtonTextColor(tcell.ColorWhite)
}

// showProgress displays a progress message by updating form title.
func (w *WorkspaceModal) showProgress(message string) {
	originalTitle := w.form.GetTitle()
	w.form.SetTitle(fmt.Sprintf(" %s ", message))

	// Restore after a brief moment (in real implementation, this would be removed when operation completes)
	go func() {
		// This would be called from the async operation completion
		w.app.QueueUpdateDraw(func() {
			w.form.SetTitle(originalTitle)
		})
	}()
}

// showSuccess displays success message and triggers callbacks.
func (w *WorkspaceModal) showSuccess() {
	// Show success state
	w.form.Clear(true)
	w.form.SetTitle(" ✓ Success ")
	w.form.SetBorderColor(theme.GetSuccessColor())

	// Add success text with better formatting
	successText := "\n[green::b]✓ Workspace created successfully![-:-:-]\n\n[white]You can now switch to this workspace and start syncing tickets.[-]"
	textView := tview.NewTextView().
		SetText(successText).
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetWordWrap(true)

	// Add the text view to the form
	w.form.AddFormItem(textView)

	// Add close button
	w.form.AddButton("Close", func() {
		if w.onSuccess != nil {
			w.onSuccess()
		}
		if w.onClose != nil {
			w.onClose()
		}
	})

	// Set button styling for success state
	w.form.SetButtonsAlign(tview.AlignCenter)
	w.form.SetButtonBackgroundColor(theme.GetSuccessColor())
	w.form.SetButtonTextColor(tcell.ColorBlack)
}

// SetOnClose sets the callback for when the modal is closed.
func (w *WorkspaceModal) SetOnClose(callback func()) {
	w.onClose = callback
}

// SetOnSuccess sets the callback for when a workspace is successfully created.
func (w *WorkspaceModal) SetOnSuccess(callback func()) {
	w.onSuccess = callback
}

// Primitive returns the underlying tview primitive.
func (w *WorkspaceModal) Primitive() tview.Primitive {
	// Return shadow form if shadows are enabled, otherwise return regular form
	if w.shadowForm != nil {
		return w.shadowForm
	}
	return w.form
}

// Name returns the view identifier.
func (w *WorkspaceModal) Name() string {
	return "workspace_modal"
}

// OnShow is called when the view becomes active.
func (w *WorkspaceModal) OnShow() {
	// Reset form state
	w.useExistingProfile = true
	w.showingNewProfile = false
	w.isValidating = false

	// Clear all fields
	w.nameField.SetText("")
	w.projectKeyField.SetText("")
	w.newProfileName.SetText("")
	w.newProfileURL.SetText("")
	w.newProfileUsername.SetText("")
	w.newProfileToken.SetText("")

	// Reset field colors
	w.nameField.SetFieldBackgroundColor(tcell.ColorDefault)
	w.projectKeyField.SetFieldBackgroundColor(tcell.ColorDefault)
	w.newProfileURL.SetFieldBackgroundColor(tcell.ColorDefault)
	w.newProfileUsername.SetFieldBackgroundColor(tcell.ColorDefault)
	w.newProfileToken.SetFieldBackgroundColor(tcell.ColorDefault)
}

// OnHide is called when the view is hidden.
func (w *WorkspaceModal) OnHide() {
	// Clear any sensitive data
	w.newProfileToken.SetText("")
}
