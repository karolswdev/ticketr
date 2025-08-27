# ticktr

A robust tool for managing user stories and tasks in a custom Markdown format with Jira integration capabilities. Ticktr provides a file-based approach to project management, allowing teams to define and track stories locally before syncing with Jira.

## Project Structure

```
ticketr/
├── cmd/
│   └── jira-story-creator/    # Main application entry point
├── internal/
│   ├── core/                  # Core business logic (Ports and Adapters architecture)
│   │   ├── domain/            # Domain models (Story, Task)
│   │   ├── ports/             # Interface definitions
│   │   └── services/          # Business logic services
│   └── adapters/              # External adapters implementation
│       ├── cli/               # Command-line interface adapter
│       ├── filesystem/        # File system operations adapter
│       └── jira/              # Jira API integration adapter
├── go.mod                     # Go module definition
└── README.md                  # This file
```

### Directory Descriptions

- **cmd/jira-story-creator**: Contains the main application executable code
- **internal/core/domain**: Defines the core domain models including Story and Task structures
- **internal/core/ports**: Contains interface definitions that define contracts for adapters
- **internal/core/services**: Implements the core business logic and orchestration
- **internal/adapters/cli**: Handles command-line interface interactions
- **internal/adapters/filesystem**: Manages file reading/writing operations for story persistence
- **internal/adapters/jira**: Integrates with Jira API for story synchronization

## Configuration

The application requires the following environment variables to connect to Jira:

- `JIRA_URL`: The base URL of your Jira instance (e.g., `https://yourcompany.atlassian.net`)
- `JIRA_EMAIL`: The email address associated with your Jira account
- `JIRA_API_KEY`: Your Jira API token (generate from Jira Account Settings → Security → API tokens)
- `JIRA_PROJECT_KEY`: The key of the Jira project where stories will be created (e.g., `PROJ`)

### Setting Environment Variables

```bash
export JIRA_URL="https://yourcompany.atlassian.net"
export JIRA_EMAIL="your.email@company.com"
export JIRA_API_KEY="your-api-token-here"
export JIRA_PROJECT_KEY="PROJ"
```

## Usage

### Running with Go

```bash
# Build the application
go build -o jira-story-creator cmd/jira-story-creator/main.go

# Run the application
./jira-story-creator -file path/to/stories.md
```

### Command Line Options

- `-file` or `-f`: Path to the input Markdown file containing stories and tasks

### Example

```bash
# Set environment variables
export JIRA_URL="https://yourcompany.atlassian.net"
export JIRA_EMAIL="your.email@company.com"
export JIRA_API_KEY="your-api-token"
export JIRA_PROJECT_KEY="PROJ"

# Run the tool
./jira-story-creator -f stories.md
```

The tool will:
1. Authenticate with Jira using the provided credentials
2. Parse the Markdown file for stories and tasks
3. Create new stories and tasks in Jira for items without existing Jira IDs
4. Update the original file with the newly created Jira IDs
5. Display a summary report of all operations

## Markdown Syntax

The full specification for the Ticktr Markdown Syntax can be found in [STORY-MARKDOWN-SPEC.md](./STORY-MARKDOWN-SPEC.md).