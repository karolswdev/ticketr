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

#### Required:
- `-file` or `-f`: Path to the input Markdown file containing stories and tasks

#### Optional:
- `--force-partial-upload`: Continue processing even if some items fail. By default, the tool exits with error code 2 when any items fail. With this flag, it will complete all possible operations and exit with code 0.
- `--verbose` or `-v`: Enable verbose logging output for debugging. Shows detailed timestamps and file locations in log messages.

### Example

```bash
# Set environment variables
export JIRA_URL="https://yourcompany.atlassian.net"
export JIRA_EMAIL="your.email@company.com"
export JIRA_API_KEY="your-api-token"
export JIRA_PROJECT_KEY="PROJ"

# Run the tool with default settings
./jira-story-creator -f stories.md

# Run with verbose logging
./jira-story-creator -f stories.md --verbose

# Run with force partial upload (continue on errors)
./jira-story-creator -f stories.md --force-partial-upload

# Run with both options
./jira-story-creator -f stories.md --verbose --force-partial-upload
```

The tool will:
1. Authenticate with Jira using the provided credentials
2. Parse the Markdown file for stories and tasks
3. Create new stories and tasks in Jira for items without existing Jira IDs
4. Update existing stories and tasks in Jira that already have Jira IDs
5. Update the original file with any newly created Jira IDs
6. Display a summary report of all operations

### Update Functionality

The tool supports both creating new items and updating existing ones:

- **Creating New Items**: Stories and tasks without Jira IDs (e.g., `# STORY: New Feature`) will be created in Jira
- **Updating Existing Items**: Stories and tasks with Jira IDs (e.g., `# STORY: [PROJ-123] Updated Feature`) will update the existing Jira issues
- **Adding Tasks to Existing Stories**: You can add new tasks to an existing story by including the story's Jira ID and adding new tasks without IDs

Example workflow:
1. First run: Creates stories and tasks, file is updated with Jira IDs
2. Edit the file: Modify descriptions, titles, or add new tasks
3. Second run: Updates existing items in Jira, creates any new tasks

## Docker Usage

### Building the Docker Image

```bash
# Build the Docker image
docker build -t ticketr:latest .

# Or use docker-compose
docker-compose build
```

### Running with Docker

```bash
# Run with environment variables from host
docker run --rm \
  -e JIRA_URL="$JIRA_URL" \
  -e JIRA_EMAIL="$JIRA_EMAIL" \
  -e JIRA_API_KEY="$JIRA_API_KEY" \
  -e JIRA_PROJECT_KEY="$JIRA_PROJECT_KEY" \
  -v $(pwd)/stories.md:/data/stories.md \
  ticketr:latest -f /data/stories.md

# Or use docker-compose (reads .env automatically)
docker-compose run --rm ticketr -f /data/stories.md
```

### Docker Compose

The included `docker-compose.yml` file provides a convenient way to run the tool:

1. Place your story files in a `stories` directory
2. Ensure your `.env` file contains the required Jira credentials
3. Run: `docker-compose run --rm ticketr`

The Docker image:
- Uses a multi-stage build for minimal size (~15MB)
- Runs as a non-root user for security
- Includes CA certificates for HTTPS connections
- Based on Alpine Linux for a minimal footprint

## Markdown Syntax

The full specification for the Ticktr Markdown Syntax can be found in [STORY-MARKDOWN-SPEC.md](./STORY-MARKDOWN-SPEC.md).

### Including Jira IDs

To update existing Jira issues, include the Jira key in square brackets:

```markdown
# STORY: [PROJ-123] Updated Story Title

## Tasks
- [PROJ-124] Updated task title
- New task to be added to existing story
```