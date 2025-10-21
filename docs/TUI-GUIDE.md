# TUI Guide

Complete guide to using the Ticketr Terminal User Interface (TUI).

## Overview

The Ticketr TUI provides a full-featured terminal interface for managing tickets, workspaces, and Jira synchronization. The interface is designed for discoverability, efficiency, and a professional user experience.

**Key Features**:
- Tri-panel layout (workspace list, ticket tree, detail view)
- Context-aware action bar with keybinding hints
- Enhanced command palette with fuzzy search
- Asynchronous operations with real-time progress
- Full keyboard navigation with F-key shortcuts
- Multi-select and bulk operations

## Getting Started

Launch the TUI:

```bash
ticketr tui
```

The TUI will open with three main panels and a bottom action bar.

## Interface Layout

```
┌─────────────────────────────────────────────────────────────────┐
│  Workspace List  │  Ticket Tree        │  Ticket Detail       │
│                  │                     │                      │
│  - Project A     │  [x] PROJ-1 Bug     │  Title: Fix login   │
│  - Project B     │  [ ] PROJ-2 Story   │  Status: In Prog    │
│  - Project C     │      └─ PROJ-3 Task │  Assignee: John     │
│                  │                     │  Description: ...   │
│                  │                     │                      │
│                  │                     │ ─── Sync Status ─── │
│                  │                     │  Ready             │
└─────────────────────────────────────────────────────────────────┘
│ [F1 Help] [F2 Sync] [F5 Refresh] [Space Select] [b Bulk] [q Quit]│
└─────────────────────────────────────────────────────────────────┘
```

### Three-Panel Layout

1. **Workspace List (Left)**: Browse and switch between configured workspaces
2. **Ticket Tree (Center)**: Hierarchical view of tickets with multi-select support
3. **Ticket Detail (Right)**: Detailed view of selected ticket with sync status bar

### Action Bar (Bottom)

The action bar is always visible and shows context-aware keybinding hints:
- Updates automatically based on current panel focus
- Shows relevant shortcuts for current context
- Provides at-a-glance guidance without opening help

## Action Bar

The bottom action bar is your primary discoverability tool. It dynamically updates to show relevant keybindings based on your current view.

### Context-Aware Display

The action bar shows different keybindings depending on which panel is focused:

**Workspace List Context**:
```
[Enter Select Workspace] [Tab Next Panel] [n New Workspace] [? Help] [q/Ctrl+C Quit]
```

**Ticket Tree Context**:
```
[Enter Open Ticket] [Space Select/Deselect] [Tab Next Panel] [Esc Back] [j/k Navigate]
[h/l Collapse/Expand] [b Bulk Ops] [/ Search] [: Commands] [? Help]
```

**Ticket Detail Context**:
```
[Esc Back] [Tab Next Panel] [e Edit] [: Commands] [? Help]
```

**Modal Context** (search, command palette, etc.):
```
[Esc Close] [Enter Confirm]
```

**Syncing Context** (during async operations):
```
[Esc Cancel Operation] [Ctrl+C Quit]
```

### How It Updates

The action bar automatically updates when:
- You change panel focus (Tab, Esc, Enter)
- You open a modal (search, command palette, bulk operations)
- An async operation starts or completes
- You navigate between views

This ensures you always see relevant shortcuts without memorizing commands.

## Command Palette

The enhanced command palette provides quick access to all TUI commands with fuzzy search and categorization.

### Opening the Command Palette

Multiple access paths for maximum discoverability:
- Press `F1` (standard help key)
- Press `Ctrl+P` (command palette convention)
- Press `:` (vim-style command mode)
- From action bar hints

### Using the Command Palette

1. **Open** with `F1`, `Ctrl+P`, or `:`
2. **Type** to fuzzy search commands (searches name and description)
3. **Navigate** with `↑`/`↓` arrow keys
4. **Toggle keybindings** with `Tab` to show/hide key hints
5. **Execute** with `Enter` or `Esc` to cancel

### Command Categories

Commands are organized into five categories:

**Navigation**:
- `help` - Show help screen
- `search` - Search tickets in current workspace
- `command-palette` - Show command palette

**View**:
- `refresh` - Refresh current workspace tickets

**Edit**:
- `bulk-operations` - Perform bulk operations on selected tickets

**Sync**:
- `pull` - Pull latest tickets from Jira
- `push` - Push tickets to Jira
- `sync` - Full sync (pull then push)

**System**:
- `quit` - Quit application

### Fuzzy Search

The search is case-insensitive and matches:
- Command names (e.g., "pull" matches "pull")
- Command descriptions (e.g., "jira" matches "Pull latest tickets from Jira")
- Partial matches (e.g., "ref" matches "refresh")

Results are ranked with name matches first, then description matches.

### Keybinding Display

Press `Tab` while in the command palette to toggle keybinding hints:

**With keybindings shown**:
```
pull
  [P or F2] Pull latest tickets from Jira
```

**With keybindings hidden**:
```
pull
  Pull latest tickets from Jira
```

This helps you learn shortcuts over time without cluttering the interface.

## F-Key Shortcuts

Ticketr uses F-keys for common application-level actions following standard conventions.

| F-Key | Action | Description |
|-------|--------|-------------|
| `F1` | Help/Command Palette | Opens enhanced command palette (standard help key) |
| `F2` | Sync/Pull | Pulls latest tickets from Jira asynchronously |
| `F5` | Refresh | Refreshes current workspace ticket list |
| `F10` | Exit | Quits application gracefully (standard exit key) |

### Why F-Keys?

F-keys provide:
- **Muscle memory**: Standard across many applications (F1=help, F10=exit)
- **Global access**: Work from any panel without mode switching
- **No conflicts**: Don't interfere with text input or navigation
- **Accessibility**: Large, distinct keys easy to find

### Alternative Keybindings

Every F-key has alternative bindings for flexibility:

- **F1**: Also `?` or `Ctrl+P`
- **F2**: Also `P` (when in ticket tree)
- **F5**: Also `r` (when in ticket tree)
- **F10**: Also `q` or `Ctrl+C`

Use whichever feels most natural.

## Discoverability Features

The TUI is designed so new users can discover features without reading documentation.

### 1. Always-Visible Action Bar

The bottom action bar is your primary guide:
- Never hidden or minimized
- Updates context-aware hints automatically
- Shows most common actions for current view
- Reminds you of shortcuts during use

### 2. Multiple Access Paths

Most features have multiple ways to access them:

**Command Palette**:
- `F1` (standard help)
- `Ctrl+P` (command palette convention)
- `:` (vim-style)

**Help Screen**:
- `?` (common help key)
- `F1` (standard help)
- Via command palette

**Sync Operations**:
- `F2` (F-key shortcut)
- `P` (letter key)
- Via command palette

This ensures you can find features even if you don't know the "official" shortcut.

### 3. In-Context Hints

Operations show guidance when relevant:
- During sync: "Press Esc to cancel"
- During bulk operations: "Select tickets with Space, press b for operations"
- In modals: Action bar shows "Esc Close, Enter Confirm"

### 4. Progressive Disclosure

The interface reveals complexity gradually:
- **Beginner**: Follow action bar hints, use command palette
- **Intermediate**: Learn letter shortcuts (p, r, s, b)
- **Advanced**: Use F-keys and vim navigation (hjkl)

You can be productive immediately while learning advanced features over time.

### 5. Consistent Conventions

Keybindings follow standard patterns:
- `Esc` always means "back" or "cancel"
- `Enter` always means "confirm" or "select"
- `Tab` always cycles focus
- `?` always shows help
- `Ctrl+C` always quits (even during operations)

Once you learn these patterns, they apply everywhere.

## Asynchronous Operations

Long-running operations (pull, push, sync) run asynchronously so the TUI remains responsive.

### How Async Works

When you trigger a sync operation:
1. Job is submitted to background queue
2. TUI remains interactive (you can navigate panels)
3. Progress updates appear in sync status bar
4. Operation completes in background
5. Ticket list refreshes automatically on success

### Progress Indicators

During async operations, you'll see:
- **Spinner**: Active animation showing operation in progress
- **Ticket count**: Current/total (e.g., "45/120 tickets")
- **Percentage**: Completion percentage (e.g., "37%")
- **Time elapsed**: How long operation has been running
- **Status message**: Current phase (e.g., "Fetching tickets...", "Processing results...")

Example progress display:
```
Sync Status: Pull in progress ⠋ 45/120 tickets (37%) - 12s elapsed
```

### Cancellation

You can cancel any active operation:
- **Press Esc**: Graceful cancellation, waits for current ticket to finish
- **Press Ctrl+C**: Force cancellation and quit application

Cancellation is immediate and responsive. Partial results are preserved.

### Error Handling

If an operation fails:
- Error message appears in sync status bar
- Partial results are kept (if any)
- You can retry immediately
- Detailed errors are logged

The TUI never crashes or freezes on errors.

## Workspace Management

### Switching Workspaces

1. Focus workspace list panel (press `Esc` until in left panel)
2. Navigate with `↑`/`↓` arrow keys
3. Press `Enter` to switch workspace
4. Ticket tree updates automatically with new workspace tickets

### Creating Workspaces

**From TUI**:
1. Press `n` in workspace list panel
2. Fill in workspace creation form:
   - Workspace name
   - Jira URL
   - Project key
   - Email
   - API token
   - (Or select existing credential profile)
3. Press `Enter` to create
4. New workspace appears in list immediately

**Using Credential Profiles**:
- Select "Use existing profile" option
- Choose from list of saved profiles
- Only enter workspace name and project key
- Credentials are reused from profile

See [Workspace Management Guide](workspace-management-guide.md) for details.

## Ticket Tree Navigation

### Basic Navigation

- `↑`/`↓` or `j`/`k`: Move selection up/down
- `Enter`: Open ticket in detail view
- `Esc`: Return to workspace list
- `Tab`: Move to next panel

### Tree Expansion

- `h`: Collapse current ticket node
- `l`: Expand current ticket node
- Collapsed nodes show `►` indicator
- Expanded nodes show `▼` indicator

### Multi-Select

Select multiple tickets for bulk operations:

1. **Select individual tickets**: Press `Space` on each ticket
   - Selected tickets show `[x]` checkbox
   - Unselected tickets show `[ ]` checkbox
2. **Select all**: Press `a` to select all visible tickets
3. **Deselect all**: Press `A` (Shift+a) to clear selections
4. **Selection count**: Title shows "(3 selected)" when tickets selected
5. **Visual feedback**: Tree border turns teal/blue when selections active

Selection state persists as you navigate the tree.

### Bulk Operations

With tickets selected:

1. Press `b` to open bulk operations menu
2. Choose operation:
   - **Update Fields**: Change Status, Priority, Assignee, or custom fields
   - **Move Tickets**: Move all selected tickets under new parent
   - **Delete Tickets**: Delete selected tickets (with confirmation)
3. Fill in operation parameters
4. Click "Apply" or press `Enter`
5. Watch real-time progress with success/failure indicators
6. Press `Cancel` or `Esc` to abort (partial changes kept)

See [Bulk Operations Guide](bulk-operations-guide.md) for comprehensive documentation.

## Search

Quickly find tickets by ID, title, or description:

1. Press `/` in ticket tree panel
2. Type search query (searches ID, title, description)
3. Navigate results with `↑`/`↓`
4. Press `Enter` to select and open ticket detail
5. Press `Esc` to close search

Search is case-insensitive and matches partial strings.

## Sync Operations

### Pull (F2)

Pull latest tickets from Jira:
- Press `F2` or `P` (in ticket tree)
- Operation runs asynchronously in background
- Progress shown in sync status bar
- Ticket tree refreshes on completion
- Press `Esc` to cancel

### Push (p)

Push local ticket changes to Jira:
- Press `p` in ticket tree
- Analyzes changed tickets
- Pushes updates asynchronously
- Progress shown in sync status bar
- Press `Esc` to cancel

### Sync (s)

Full bidirectional sync (pull then push):
- Press `s` in ticket tree
- Pulls latest changes first
- Then pushes local changes
- Two-phase progress tracking
- Press `Esc` to cancel current phase

### Refresh (F5 / r)

Reload tickets from local database:
- Press `F5` or `r`
- Fast local-only operation
- Useful after bulk operations
- No Jira API calls made

## Tips and Tricks

### Efficient Workflow

1. **Use F-keys for common operations**: F2 (sync), F5 (refresh), F1 (help)
2. **Learn letter shortcuts**: p (push), r (refresh), s (sync), b (bulk ops)
3. **Use Tab to navigate**: Faster than mouse or arrow keys
4. **Keep action bar visible**: Glance at bottom for context-aware hints
5. **Use command palette for discovery**: Press F1 to explore available commands

### Keyboard-Only Navigation

The TUI is fully keyboard-driven:
- Never requires mouse
- All features accessible via keyboard
- Vim-style navigation supported (hjkl)
- Consistent conventions across views

### Learning Shortcuts

1. **Start with action bar**: Use visible hints for first few days
2. **Memorize F-keys**: F1, F2, F5, F10 cover 80% of usage
3. **Practice letter shortcuts**: Add p, r, s, b to muscle memory
4. **Toggle keybindings in palette**: Use Tab to learn command shortcuts
5. **Use help screen**: Press ? for comprehensive reference

### Performance Tips

- **Large workspaces**: Use search (`/`) to find tickets quickly
- **Async operations**: Don't wait for sync to complete, keep working
- **Refresh vs Pull**: Use `r` (refresh) for local updates, `F2` (pull) for Jira sync
- **Bulk operations**: Select multiple tickets and update in one operation

### Troubleshooting

**TUI doesn't render correctly**:
- Ensure terminal supports 256 colors or true color
- Try resizing terminal window
- Check terminal emulator compatibility (see [KEYBINDINGS.md](KEYBINDINGS.md))

**F-keys don't work**:
- Use alternative shortcuts (? for help, P for pull, etc.)
- Check terminal emulator function key settings
- Some terminals require Fn key modifier

**Slow performance with large ticket sets**:
- Use search to filter tickets
- Consider splitting into multiple workspaces
- Refresh locally (r) instead of pulling frequently (F2)

**Sync operation stuck**:
- Press Esc to cancel gracefully
- Press Ctrl+C to force quit if needed
- Check network connection and Jira API status

## Related Documentation

- [KEYBINDINGS.md](KEYBINDINGS.md) - Complete keybinding reference
- [Bulk Operations Guide](bulk-operations-guide.md) - Bulk operations workflow
- [Workspace Management Guide](workspace-management-guide.md) - Workspace setup and management
- [Sync Strategies Guide](sync-strategies-guide.md) - Conflict resolution strategies
- [ARCHITECTURE.md](ARCHITECTURE.md) - TUI technical architecture

## Accessibility

The TUI is designed to be accessible:
- **Keyboard-only**: No mouse required
- **High contrast**: Clear visual hierarchy with borders and colors
- **Context hints**: Always-visible action bar for guidance
- **Screen reader friendly**: Uses standard tview primitives
- **Consistent layout**: Predictable panel positions and behavior

## Future Enhancements

Planned improvements:
- Ticket editing in TUI (currently view-only)
- Custom keybinding configuration
- Theme customization
- Filter panels for advanced JQL queries
- Ticket creation wizard
- Inline ticket updates (edit in tree)

See [ROADMAP.md](../ROADMAP.md) for full feature roadmap.

## Feedback

We welcome feedback on the TUI experience:
- Report bugs: [GitHub Issues](https://github.com/karolswdev/ticketr/issues)
- Request features: [GitHub Discussions](https://github.com/karolswdev/ticketr/discussions)
- Contribute: [CONTRIBUTING.md](../CONTRIBUTING.md)

Your input helps us improve discoverability and user experience.
