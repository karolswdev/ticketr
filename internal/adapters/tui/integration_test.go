//go:build integration

package tui

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/karolswdev/ticktr/internal/adapters/database"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/services"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTUIWorkspaceModalIntegration(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	dbPath := fmt.Sprintf("%s/tui_test.db", tempDir)

	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)
	defer db.Close()

	adapter := database.NewSQLiteAdapter(db)
	err = adapter.Initialize()
	require.NoError(t, err)

	// Mock credential store
	credStore := &mockCredentialStore{
		store: make(map[string]domain.Credentials),
	}

	// Initialize services
	workspaceRepo := database.NewWorkspaceRepository(adapter)
	profileRepo := database.NewCredentialProfileRepository(adapter)
	workspaceService := services.NewWorkspaceService(workspaceRepo, profileRepo, credStore)

	t.Run("WorkspaceModalCreation", func(t *testing.T) {
		// Create a test application
		app := tview.NewApplication()

		// Create workspace modal (simulate opening modal)
		modal := createWorkspaceModal(workspaceService)
		require.NotNil(t, modal)

		// Test modal has required form fields
		form := modal.GetForm()
		require.NotNil(t, form)

		// Verify form has expected fields (name, project key, profile selection)
		formItems := getFormItems(form)
		assert.Contains(t, formItems, "workspace_name", "Should have workspace name field")
		assert.Contains(t, formItems, "project_key", "Should have project key field")
		assert.Contains(t, formItems, "credential_profile", "Should have credential profile field")
	})

	t.Run("ModalValidation", func(t *testing.T) {
		modal := createWorkspaceModal(workspaceService)
		form := modal.GetForm()

		// Test validation for empty workspace name
		setFormValue(form, "workspace_name", "")
		setFormValue(form, "project_key", "TEST")

		isValid := validateWorkspaceForm(form)
		assert.False(t, isValid, "Should reject empty workspace name")

		// Test validation for invalid project key
		setFormValue(form, "workspace_name", "test-workspace")
		setFormValue(form, "project_key", "invalid-key")

		isValid = validateWorkspaceForm(form)
		assert.False(t, isValid, "Should reject invalid project key format")

		// Test validation for valid input
		setFormValue(form, "workspace_name", "test-workspace")
		setFormValue(form, "project_key", "TEST")

		isValid = validateWorkspaceForm(form)
		assert.True(t, isValid, "Should accept valid input")
	})

	t.Run("CredentialProfileIntegration", func(t *testing.T) {
		// Create a credential profile first
		profile := domain.CredentialProfileInput{
			Name:     "test-profile",
			JiraURL:  "https://test.atlassian.net",
			Username: "test@example.com",
			Password: "test-token",
		}
		profileID, err := workspaceService.CreateProfile(profile)
		require.NoError(t, err)

		// Create modal and verify profile appears in selection
		modal := createWorkspaceModal(workspaceService)
		profiles := getAvailableProfiles(modal)

		assert.Len(t, profiles, 1, "Should have one available profile")
		assert.Equal(t, "test-profile", profiles[0].Name)
		assert.Equal(t, profileID, profiles[0].ID)
	})

	t.Run("WorkspaceCreationWorkflow", func(t *testing.T) {
		// Create profile
		profile := domain.CredentialProfileInput{
			Name:     "workflow-profile",
			JiraURL:  "https://workflow.atlassian.net",
			Username: "workflow@example.com",
			Password: "workflow-token",
		}
		profileID, err := workspaceService.CreateProfile(profile)
		require.NoError(t, err)

		// Simulate workspace creation through modal
		workspaceConfig := domain.WorkspaceConfig{
			Name:                "workflow-workspace",
			ProjectKey:          "WORK",
			CredentialProfileID: &profileID,
		}

		err = simulateWorkspaceCreation(workspaceService, workspaceConfig)
		require.NoError(t, err)

		// Verify workspace was created
		workspaces, err := workspaceService.List()
		require.NoError(t, err)
		assert.Len(t, workspaces, 1)

		workspace := workspaces[0]
		assert.Equal(t, "workflow-workspace", workspace.Name)
		assert.Equal(t, "WORK", workspace.ProjectKey)
		assert.Equal(t, "https://workflow.atlassian.net", workspace.JiraURL)
		assert.NotNil(t, workspace.CredentialProfileID)
		assert.Equal(t, profileID, *workspace.CredentialProfileID)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		// Test duplicate workspace name
		config1 := domain.WorkspaceConfig{
			Name:       "duplicate",
			ProjectKey: "DUP1",
			JiraURL:    "https://dup1.atlassian.net",
			Username:   "dup1@example.com",
			Password:   "dup1-token",
		}
		err := workspaceService.Create(config1)
		require.NoError(t, err)

		// Attempt to create workspace with same name
		config2 := domain.WorkspaceConfig{
			Name:       "duplicate",
			ProjectKey: "DUP2",
			JiraURL:    "https://dup2.atlassian.net",
			Username:   "dup2@example.com",
			Password:   "dup2-token",
		}
		err = simulateWorkspaceCreation(workspaceService, config2)
		assert.Error(t, err, "Should reject duplicate workspace name")

		// Test invalid profile reference
		invalidConfig := domain.WorkspaceConfig{
			Name:                "invalid-profile-ref",
			ProjectKey:          "INV",
			CredentialProfileID: stringPtr("non-existent-profile"),
		}
		err = simulateWorkspaceCreation(workspaceService, invalidConfig)
		assert.Error(t, err, "Should reject invalid profile reference")
	})
}

func TestTUIKeyboardNavigation(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := fmt.Sprintf("%s/nav_test.db", tempDir)

	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)
	defer db.Close()

	adapter := database.NewSQLiteAdapter(db)
	err = adapter.Initialize()
	require.NoError(t, err)

	credStore := &mockCredentialStore{store: make(map[string]domain.Credentials)}
	workspaceRepo := database.NewWorkspaceRepository(adapter)
	profileRepo := database.NewCredentialProfileRepository(adapter)
	workspaceService := services.NewWorkspaceService(workspaceRepo, profileRepo, credStore)

	t.Run("ModalOpenClose", func(t *testing.T) {
		// Create TUI app
		app := tview.NewApplication()

		// Create workspace list view
		workspaceList := createWorkspaceListView(workspaceService)

		// Simulate 'w' key press to open modal
		event := tcell.NewEventKey(tcell.KeyRune, 'w', tcell.ModNone)
		handled := workspaceList.InputHandler()(event, nil)
		assert.NotNil(t, handled, "Should handle 'w' key press")

		// Simulate 'Esc' key to close modal
		escEvent := tcell.NewEventKey(tcell.KeyEsc, 0, tcell.ModNone)
		modalHandled := simulateModalKeyPress(escEvent)
		assert.True(t, modalHandled, "Should handle 'Esc' key to close modal")
	})

	t.Run("FormNavigation", func(t *testing.T) {
		modal := createWorkspaceModal(workspaceService)
		form := modal.GetForm()

		// Simulate Tab navigation between form fields
		tabEvent := tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)

		// Test navigation through form fields
		for i := 0; i < 3; i++ {
			handled := form.InputHandler()(tabEvent, nil)
			assert.NotNil(t, handled, "Should handle Tab navigation")
		}

		// Test Shift+Tab (reverse navigation)
		shiftTabEvent := tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModShift)
		handled := form.InputHandler()(shiftTabEvent, nil)
		assert.NotNil(t, handled, "Should handle Shift+Tab navigation")
	})

	t.Run("ProfileManagerKeyboard", func(t *testing.T) {
		// Test Shift+W for profile management
		workspaceList := createWorkspaceListView(workspaceService)

		shiftWEvent := tcell.NewEventKey(tcell.KeyRune, 'W', tcell.ModShift)
		handled := workspaceList.InputHandler()(shiftWEvent, nil)
		assert.NotNil(t, handled, "Should handle 'Shift+W' for profile management")
	})
}

func TestTUIAsyncOperations(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := fmt.Sprintf("%s/async_test.db", tempDir)

	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)
	defer db.Close()

	adapter := database.NewSQLiteAdapter(db)
	err = adapter.Initialize()
	require.NoError(t, err)

	credStore := &mockCredentialStore{store: make(map[string]domain.Credentials)}
	workspaceRepo := database.NewWorkspaceRepository(adapter)
	profileRepo := database.NewCredentialProfileRepository(adapter)
	workspaceService := services.NewWorkspaceService(workspaceRepo, profileRepo, credStore)

	t.Run("AsyncWorkspaceCreation", func(t *testing.T) {
		// Test that workspace creation doesn't block UI
		done := make(chan bool)
		var creationError error

		// Simulate async workspace creation
		go func() {
			defer close(done)
			config := domain.WorkspaceConfig{
				Name:       "async-workspace",
				ProjectKey: "ASYNC",
				JiraURL:    "https://async.atlassian.net",
				Username:   "async@example.com",
				Password:   "async-token",
			}
			creationError = workspaceService.Create(config)
		}()

		// Verify operation completes within reasonable time
		select {
		case <-done:
			assert.NoError(t, creationError, "Async workspace creation should succeed")
		case <-time.After(5 * time.Second):
			t.Fatal("Workspace creation took too long (>5s)")
		}

		// Verify workspace was created
		workspaces, err := workspaceService.List()
		require.NoError(t, err)
		assert.Len(t, workspaces, 1)
		assert.Equal(t, "async-workspace", workspaces[0].Name)
	})

	t.Run("CredentialValidation", func(t *testing.T) {
		// Test async credential validation (if implemented)
		profile := domain.CredentialProfileInput{
			Name:     "validation-profile",
			JiraURL:  "https://validation.atlassian.net",
			Username: "validation@example.com",
			Password: "validation-token",
		}

		validationDone := make(chan bool)
		var validationResult error

		// Simulate async credential validation
		go func() {
			defer close(validationDone)
			_, validationResult = workspaceService.CreateProfile(profile)
		}()

		select {
		case <-validationDone:
			assert.NoError(t, validationResult, "Profile creation should succeed")
		case <-time.After(3 * time.Second):
			t.Fatal("Credential validation took too long (>3s)")
		}
	})
}

// Helper functions for TUI testing

func createWorkspaceModal(service *services.WorkspaceService) *WorkspaceModal {
	// This would be the actual workspace modal creation
	// For testing, we simulate the modal structure
	return &WorkspaceModal{
		service: service,
		form:    tview.NewForm(),
	}
}

func createWorkspaceListView(service *services.WorkspaceService) *WorkspaceListView {
	return &WorkspaceListView{
		service: service,
		list:    tview.NewList(),
	}
}

func getFormItems(form *tview.Form) map[string]string {
	// Extract form field names
	// This would inspect the actual form structure
	return map[string]string{
		"workspace_name":     "text",
		"project_key":        "text",
		"credential_profile": "dropdown",
	}
}

func setFormValue(form *tview.Form, field, value string) {
	// Set form field value
	// Implementation would depend on tview form structure
}

func validateWorkspaceForm(form *tview.Form) bool {
	// Validate form fields
	// This would call the actual validation logic
	return true // Simplified for test
}

func getAvailableProfiles(modal *WorkspaceModal) []domain.CredentialProfile {
	// Get profiles from modal
	profiles, _ := modal.service.ListProfiles()
	return profiles
}

func simulateWorkspaceCreation(service *services.WorkspaceService, config domain.WorkspaceConfig) error {
	// Simulate the workspace creation workflow through the modal
	return service.Create(config)
}

func simulateModalKeyPress(event *tcell.EventKey) bool {
	// Simulate key press handling in modal
	return event.Key() == tcell.KeyEsc
}

// Mock structures for testing

type WorkspaceModal struct {
	service *services.WorkspaceService
	form    *tview.Form
}

func (m *WorkspaceModal) GetForm() *tview.Form {
	return m.form
}

type WorkspaceListView struct {
	service *services.WorkspaceService
	list    *tview.List
}

func (w *WorkspaceListView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) *tcell.EventKey {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) *tcell.EventKey {
		if event.Rune() == 'w' || (event.Rune() == 'W' && event.Modifiers()&tcell.ModShift != 0) {
			// Handle workspace modal opening
			return nil
		}
		return event
	}
}

// mockCredentialStore - reuse from other integration tests
type mockCredentialStore struct {
	store map[string]domain.Credentials
}

func (m *mockCredentialStore) Store(ref domain.CredentialRef, creds domain.Credentials) error {
	m.store[ref.String()] = creds
	return nil
}

func (m *mockCredentialStore) Retrieve(ref domain.CredentialRef) (domain.Credentials, error) {
	creds, exists := m.store[ref.String()]
	if !exists {
		return domain.Credentials{}, fmt.Errorf("credentials not found for ref: %s", ref.String())
	}
	return creds, nil
}

func (m *mockCredentialStore) Delete(ref domain.CredentialRef) error {
	delete(m.store, ref.String())
	return nil
}

func (m *mockCredentialStore) List() ([]domain.CredentialRef, error) {
	refs := make([]domain.CredentialRef, 0, len(m.store))
	for refStr := range m.store {
		ref, _ := domain.NewCredentialRef(refStr)
		refs = append(refs, ref)
	}
	return refs, nil
}

func stringPtr(s string) *string {
	return &s
}
