// Package webhook provides HTTP webhook server functionality for receiving JIRA events.
// It processes incoming webhook payloads and triggers local file updates automatically.
package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

    "github.com/karolswdev/ticketr/internal/core/services"
)

// Server handles incoming webhook requests from JIRA
type Server struct {
	pullService *services.PullService
	filePath    string
	secret      string
	projectKey  string
}

// NewServer creates a new webhook server instance
//
// Parameters:
//   - pullService: Service for pulling and merging tickets
//   - filePath: Path to the markdown file to update
//   - secret: Optional webhook secret for validation
//   - projectKey: JIRA project key
//
// Returns:
//   - *Server: Configured webhook server
func NewServer(pullService *services.PullService, filePath string, secret string, projectKey string) *Server {
	return &Server{
		pullService: pullService,
		filePath:    filePath,
		secret:      secret,
		projectKey:  projectKey,
	}
}

// JiraWebhookPayload represents the incoming webhook payload from JIRA
type JiraWebhookPayload struct {
	WebhookEvent string `json:"webhookEvent"`
	Issue        struct {
		ID     string `json:"id"`
		Key    string `json:"key"`
		Fields struct {
			Summary     string `json:"summary"`
			Description string `json:"description"`
			IssueType   struct {
				Name string `json:"name"`
			} `json:"issuetype"`
			Status struct {
				Name string `json:"name"`
			} `json:"status"`
			Project struct {
				Key string `json:"key"`
			} `json:"project"`
		} `json:"fields"`
	} `json:"issue"`
	ChangeLog struct {
		Items []struct {
			Field      string `json:"field"`
			FromString string `json:"fromString"`
			ToString   string `json:"toString"`
		} `json:"items"`
	} `json:"changelog"`
}

// HandleWebhook processes incoming JIRA webhook requests
func (s *Server) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading webhook body: %v", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate webhook signature if secret is configured
	if s.secret != "" {
		signature := r.Header.Get("X-Hub-Signature-256")
		if signature == "" {
			signature = r.Header.Get("X-Atlassian-Webhook-Signature")
		}
		
		if !s.validateSignature(body, signature) {
			log.Printf("Invalid webhook signature")
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}
	}

	// Parse the webhook payload
	var payload JiraWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Error parsing webhook payload: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Check if this webhook is for our project
	if payload.Issue.Fields.Project.Key != s.projectKey {
		log.Printf("Ignoring webhook for project %s (configured for %s)", 
			payload.Issue.Fields.Project.Key, s.projectKey)
		w.WriteHeader(http.StatusOK)
		return
	}

	// Log the event
	log.Printf("Received webhook event: %s for issue %s", 
		payload.WebhookEvent, payload.Issue.Key)

	// Process the webhook based on event type
	switch payload.WebhookEvent {
	case "jira:issue_created", "jira:issue_updated":
		// Trigger a pull to update the local file
		if err := s.updateLocalFile(payload.Issue.Key); err != nil {
			log.Printf("Error updating local file: %v", err)
			// Don't return error to JIRA - log and continue
		}
	case "jira:issue_deleted":
		log.Printf("Issue %s was deleted - manual intervention may be required", 
			payload.Issue.Key)
	}

	// Return success to JIRA
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// validateSignature validates the webhook signature using HMAC-SHA256
func (s *Server) validateSignature(payload []byte, signature string) bool {
	if signature == "" {
		return false
	}

	// Remove "sha256=" prefix if present
	signature = strings.TrimPrefix(signature, "sha256=")

	// Calculate expected signature
	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write(payload)
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	// Compare signatures
	return hmac.Equal([]byte(signature), []byte(expectedSig))
}

// updateLocalFile triggers a pull operation to update the local markdown file
func (s *Server) updateLocalFile(issueKey string) error {
	log.Printf("Updating local file for issue %s", issueKey)

	// Build JQL to fetch just this ticket
	jql := fmt.Sprintf("key = %s", issueKey)

	// Execute pull with remote-wins strategy for webhook updates
	result, err := s.pullService.Pull(s.filePath, services.PullOptions{
		ProjectKey: s.projectKey,
		JQL:        jql,
		Strategy:   "remote-wins", // Always use remote version for webhooks
	})

	if err != nil {
		return fmt.Errorf("failed to pull ticket %s: %w", issueKey, err)
	}

	// Log the result
	if result.TicketsUpdated > 0 {
		log.Printf("Updated ticket %s in %s", issueKey, s.filePath)
	} else if result.TicketsPulled > 0 {
		log.Printf("Added new ticket %s to %s", issueKey, s.filePath)
	}

	return nil
}

// Start starts the webhook server on the specified port
//
// Parameters:
//   - port: The port to listen on
//
// Returns:
//   - error: An error if the server fails to start
func (s *Server) Start(port string) error {
	http.HandleFunc("/webhook", s.HandleWebhook)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Starting webhook server on port %s", port)
	log.Printf("Webhook endpoint: http://localhost:%s/webhook", port)
	log.Printf("Health check: http://localhost:%s/health", port)
	log.Printf("Updating file: %s", s.filePath)
	
	if s.secret != "" {
		log.Printf("Webhook signature validation enabled")
	}

	return http.ListenAndServe(":"+port, nil)
}
