package jira

import (
	"fmt"
	"os"

	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// AdapterVersion represents the Jira adapter version
type AdapterVersion string

const (
	// AdapterV1 uses the custom HTTP implementation
	AdapterV1 AdapterVersion = "v1"
	// AdapterV2 uses the andygrunwald/go-jira library
	AdapterV2 AdapterVersion = "v2"
)

// GetAdapterVersion returns the adapter version from environment variable
// Defaults to v2 if not set
func GetAdapterVersion() AdapterVersion {
	version := os.Getenv("TICKETR_JIRA_ADAPTER_VERSION")
	switch version {
	case "v1":
		return AdapterV1
	case "v2":
		return AdapterV2
	case "":
		// Default to v2 (new implementation)
		return AdapterV2
	default:
		// Unknown version, default to v2
		return AdapterV2
	}
}

// NewJiraAdapterFromConfigWithVersion creates a Jira adapter from workspace config
// using the version specified in the environment variable TICKETR_JIRA_ADAPTER_VERSION
func NewJiraAdapterFromConfigWithVersion(config *domain.WorkspaceConfig, fieldMappings map[string]interface{}) (ports.JiraPort, error) {
	version := GetAdapterVersion()

	switch version {
	case AdapterV1:
		return NewJiraAdapterFromConfig(config, fieldMappings)
	case AdapterV2:
		return NewJiraAdapterV2FromConfig(config, fieldMappings)
	default:
		// This should never happen due to GetAdapterVersion() defaults
		return nil, fmt.Errorf("unsupported adapter version: %s", version)
	}
}

// NewJiraAdapterWithVersion creates a Jira adapter using environment variables
// and the version specified in TICKETR_JIRA_ADAPTER_VERSION
func NewJiraAdapterWithVersion(fieldMappings map[string]interface{}) (ports.JiraPort, error) {
	version := GetAdapterVersion()

	switch version {
	case AdapterV1:
		return NewJiraAdapterWithConfig(fieldMappings)
	case AdapterV2:
		// V2 doesn't support env-based initialization, so we fail gracefully
		return nil, fmt.Errorf("v2 adapter requires workspace configuration, use NewJiraAdapterFromConfigWithVersion instead")
	default:
		return nil, fmt.Errorf("unsupported adapter version: %s", version)
	}
}
