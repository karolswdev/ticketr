# Testing Ticketr v3 TUI Locally

This guide walks you through compiling, configuring, and testing the v3 TUI implementation (Weeks 11-13) on your local machine.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Compilation](#compilation)
- [Configuration](#configuration)
- [Testing Workflows](#testing-workflows)
- [Verification Checklist](#verification-checklist)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Software
- **Go 1.23+**: Check with `go version`
- **Git**: For cloning and managing the repository
- **Terminal**: Any modern terminal with Unicode support
- **Jira Account**: With API token access

### Operating System Requirements
- **macOS**: Tested on macOS 10.15+
- **Linux**: Tested on Ubuntu 20.04+, Fedora 35+
- **Windows**: WSL2 recommended (native Windows support in development)

### Jira Requirements
- Access to a Jira instance (Cloud or Server)
- API token generated (for Jira Cloud: https://id.atlassian.com/manage-profile/security/api-tokens)
- At least one project with some tickets

---

## Compilation

### 1. Navigate to Repository
```bash
cd /home/karol/dev/code/ticketr
```

### 2. Ensure You're on the Correct Branch
```bash
git branch
# Should show: * feature/v3
# If not:
git checkout feature/v3
```

### 3. Pull Latest Changes
```bash
git pull origin feature/v3
```

### 4. Build the Binary
```bash
go build -o ticketr ./cmd/ticketr
```

**Expected Output:**
- No errors
- Binary created: `./ticketr` (approximately 20-25 MB)

### 5. Verify Build
```bash
./ticketr --version
# Should output version information

./ticketr --help
# Should show available commands including 'v3' and 'workspace'
```

**Success Indicator:** Binary runs without errors and shows help menu.

---

## Configuration

### Step 1: Enable v3 Beta Features

v3 features (workspaces, SQLite, TUI) are behind feature flags. Enable them:

```bash
./ticketr v3 enable beta
```

**Expected Output:**
```
v3.0 features enabled successfully

Enabled features:
✓ SQLite database
✓ Workspace support
✓ TUI interface

Configuration file: ~/.config/ticketr/features.yaml
Database path: ~/.local/share/ticketr/ticketr.db

Next steps:
1. Create a workspace: ticketr workspace create <name>
2. Launch TUI: ticketr tui
```

**What This Does:**
- Creates configuration directory: `~/.config/ticketr/`
- Initializes SQLite database: `~/.local/share/ticketr/ticketr.db`
- Sets feature flags: `enable_workspaces: true`, `enable_tui: true`

### Step 2: Verify Feature Status

```bash
./ticketr v3 status
```

**Expected Output:**
```
Ticketr v3.0 Feature Status

Features:
✓ SQLite Database    Enabled    ~/.local/share/ticketr/ticketr.db
✓ Workspaces         Enabled    No workspaces configured yet
✓ TUI Interface      Enabled    Ready to launch

Database:
  Path: ~/.local/share/ticketr/ticketr.db
  Size: 32 KB
  Schema Version: 1

Configuration:
  Features file: ~/.config/ticketr/features.yaml
  Data directory: ~/.local/share/ticketr/
```

### Step 3: Obtain Jira API Token

**For Jira Cloud:**
1. Go to https://id.atlassian.com/manage-profile/security/api-tokens
2. Click "Create API token"
3. Give it a name (e.g., "Ticketr TUI Testing")
4. Copy the token (you won't see it again!)

**For Jira Server/Data Center:**
- Use your Jira password OR generate a Personal Access Token (PAT) if available
- Consult your Jira admin for authentication method

### Step 4: Create Your First Workspace

Replace placeholders with your actual Jira details:

```bash
./ticketr workspace create my-project \
  --url https://YOUR-DOMAIN.atlassian.net \
  --project YOUR-PROJECT-KEY \
  --username your-email@example.com \
  --token YOUR-API-TOKEN
```

**Example:**
```bash
./ticketr workspace create acme-backend \
  --url https://acme.atlassian.net \
  --project BACK \
  --username john.doe@acme.com \
  --token ATATTxxxxxxxxxxxxxxxxxxxxx
```

**Expected Output:**
```
Workspace 'acme-backend' created successfully
  Project: BACK
  URL: https://acme.atlassian.net
  Username: john.doe@acme.com

Credentials stored securely in OS keychain

This is your first workspace and has been set as the default.
```

**Security Note:**
- Credentials are stored in your OS keychain (macOS Keychain, GNOME Keyring, Windows Credential Manager)
- They are NOT stored in plain text or in the SQLite database
- Only a reference ID is stored in the database

### Step 5: Verify Workspace Created

```bash
./ticketr workspace list
```

**Expected Output:**
```
NAME            PROJECT  URL                              DEFAULT  LAST USED
----            -------  ---                              -------  ---------
acme-backend    BACK     https://acme.atlassian.net       *        never
```

### Step 6: Sync Tickets from Jira (Optional but Recommended)

To populate the database with real ticket data for testing:

```bash
./ticketr pull
```

**Note:** This command may not be fully implemented yet in v3. If it fails, you can manually insert test tickets (see Troubleshooting section).

---

## Testing Workflows

### Workflow 1: Launch TUI and Explore Navigation

#### 1.1 Launch TUI
```bash
./ticketr tui
```

**Expected Display:**
```
┌─ Workspaces ──────────────┐┌─ Ticket Tree ───────────────┐┌─ Ticket Detail ──────────────┐
│ * acme-backend             ││ Tickets                      ││ No ticket selected.           │
│                            ││ Loading tickets...           ││ Select a ticket from the tree │
│                            ││                              ││ view.                         │
│                            ││                              ││                               │
│                            ││                              ││                               │
└────────────────────────────┘└──────────────────────────────┘└───────────────────────────────┘
```

**Tri-Panel Layout:**
- **Left Panel (30 chars):** Workspace list - GREEN border indicates focus
- **Middle Panel (40%):** Ticket tree
- **Right Panel (60%):** Ticket detail view

#### 1.2 Test Global Navigation

| Key       | Action                                    | Expected Result                          |
|-----------|-------------------------------------------|------------------------------------------|
| `Tab`     | Cycle focus forward                       | Border color changes: workspace → tree → detail → workspace |
| `Shift+Tab` | Cycle focus backward                    | Border color changes in reverse          |
| `Esc`     | Go back one panel                         | detail → tree → workspace                |
| `?`       | Show help screen                          | Full-screen help overlay appears         |
| `q`       | Quit application                          | TUI exits cleanly                        |
| `Ctrl+C`  | Quit application                          | TUI exits immediately                    |

#### 1.3 Test Workspace Panel (Left)

**Focus the workspace panel** (should start focused, or press `Tab` until green border on left panel):

| Key       | Action                                    | Expected Result                          |
|-----------|-------------------------------------------|------------------------------------------|
| `j` or `↓` | Move down in workspace list              | Next workspace highlighted               |
| `k` or `↑` | Move up in workspace list                | Previous workspace highlighted           |
| `Enter`   | Select workspace                          | Ticket tree loads tickets from workspace, focus shifts to tree |

**Test:**
1. Press `Enter` on your workspace
2. Verify ticket tree panel shows "Loading tickets..." then displays tickets

#### 1.4 Test Ticket Tree Panel (Middle)

**Focus the ticket tree panel** (press `Tab` once from workspace list):

| Key       | Action                                    | Expected Result                          |
|-----------|-------------------------------------------|------------------------------------------|
| `j` or `↓` | Move down in tree                        | Next ticket/task highlighted             |
| `k` or `↑` | Move up in tree                          | Previous ticket/task highlighted         |
| `h` or `←` | Collapse node                            | If ticket has tasks, they collapse       |
| `l` or `→` | Expand node                              | If ticket has tasks, they expand         |
| `Enter`   | Open ticket detail                        | Detail panel shows ticket information, focus shifts to detail |

**Test:**
1. Navigate to a ticket with tasks using `j`/`k`
2. Press `l` to expand, `h` to collapse
3. Press `Enter` on a ticket to view details

#### 1.5 Test Ticket Detail Panel (Right) - Read-Only Mode

**Focus the detail panel** (press `Enter` on a ticket in tree):

| Key       | Action                                    | Expected Result                          |
|-----------|-------------------------------------------|------------------------------------------|
| `j`       | Scroll down                               | Detail view scrolls down                 |
| `k`       | Scroll up                                 | Detail view scrolls up                   |
| `g`       | Go to top                                 | Scrolls to beginning of ticket           |
| `G`       | Go to bottom (Shift+g)                    | Scrolls to end of ticket                 |
| `e`       | Enter edit mode                           | Form appears for editing                 |
| `Esc`     | Return to ticket tree                     | Focus returns to tree panel              |

**Test:**
1. Select a ticket with `Enter` from tree
2. Verify all ticket fields display:
   - Title (yellow, bold)
   - Jira ID (green if synced, red if not)
   - Source Line
   - Description
   - Custom Fields
   - Acceptance Criteria
   - Tasks (if any)
3. Test scrolling with `j`/`k` and `g`/`G`

### Workflow 2: Edit Ticket and Test Validation

#### 2.1 Enter Edit Mode

1. Focus detail panel (press `Enter` on a ticket)
2. Press `e` key
3. Verify form appears with editable fields

**Expected Display:**
```
┌─ Edit Ticket ─────────────────────────────────────────────┐
│ Title: [Implementation of login feature                 ] │
│                                                            │
│ Jira ID: [BACK-123          ]                              │
│                                                            │
│ Description:                                               │
│ ┌────────────────────────────────────────────────────────┐ │
│ │Implement user authentication with OAuth2               │ │
│ │                                                         │ │
│ └────────────────────────────────────────────────────────┘ │
│                                                            │
│ Custom Fields (key=value per line):                       │
│ ┌────────────────────────────────────────────────────────┐ │
│ │priority=High                                           │ │
│ │sprint=Sprint 5                                         │ │
│ └────────────────────────────────────────────────────────┘ │
│                                                            │
│ [Save]  [Cancel]                                           │
└────────────────────────────────────────────────────────────┘
Edit Mode | Keys: Ctrl+S save | Esc cancel
```

#### 2.2 Test Field Editing

| Key       | Action                                    | Expected Result                          |
|-----------|-------------------------------------------|------------------------------------------|
| `Tab`     | Move to next field                        | Cursor moves to next input               |
| Type text | Edit field content                        | Text appears, dirty indicator (*) shows  |
| `Esc`     | Cancel editing                            | Returns to read-only mode, changes discarded |

**Test:**
1. Modify the Title field
2. Verify status bar shows: `Edit Mode * (unsaved changes)`
3. Press `Esc`
4. Verify changes are discarded and you return to read-only mode

#### 2.3 Test Validation - Empty Title

1. Press `e` to enter edit mode
2. Clear the Title field (delete all text)
3. Navigate to "Save" button and press `Enter`

**Expected Result:**
- Modal dialog appears:
```
┌─ Validation Errors ─┐
│ Title is required    │
│                      │
│       [OK]           │
└──────────────────────┘
```
4. Press `Enter` on OK
5. You remain in edit mode to fix the error

#### 2.4 Test Validation - Invalid Jira ID

1. Press `e` to enter edit mode
2. Change Jira ID to invalid format (e.g., "abc-123" lowercase)
3. Navigate to "Save" button and press `Enter`

**Expected Result:**
- Modal dialog appears:
```
┌─ Validation Errors ──────────────────────────────────────┐
│ Jira ID must match format PROJECT-123 (uppercase         │
│ letters, dash, numbers)                                   │
│                                                           │
│                       [OK]                                │
└───────────────────────────────────────────────────────────┘
```

#### 2.5 Test Successful Save

1. Press `e` to enter edit mode
2. Make a valid change (e.g., add a word to Description)
3. Navigate to "Save" button and press `Enter`

**Expected Result:**
- Edit mode closes
- Returns to read-only display mode
- Updated content is visible
- No dirty indicator

**Note:** Actual sync to Jira is not yet implemented. Save updates in-memory state only.

### Workflow 3: Multi-Workspace Testing (If You Have Multiple Workspaces)

#### 3.1 Create Second Workspace

```bash
./ticketr workspace create acme-frontend \
  --url https://acme.atlassian.net \
  --project FRONT \
  --username john.doe@acme.com \
  --token ATATTxxxxxxxxxxxxxxxxxxxxx
```

#### 3.2 Test Workspace Switching in TUI

1. Launch TUI: `./ticketr tui`
2. Focus workspace list (left panel)
3. Use `j`/`k` to navigate between workspaces
4. Press `Enter` on different workspace
5. Verify ticket tree reloads with tickets from new workspace

### Workflow 4: Help System

1. Launch TUI: `./ticketr tui`
2. Press `?` key

**Expected Result:**
- Help screen appears with comprehensive keybinding documentation
- Covers all three panels and both read/edit modes
- Press `?` or `Esc` to close help

**Verify Help Content Includes:**
- Global navigation (Tab, Shift+Tab, Esc, ?, q)
- Workspace list keybindings (j/k/↓/↑, Enter)
- Ticket tree keybindings (j/k/h/l, arrows, Enter)
- Ticket detail read-only (e, j/k, g/G, Esc)
- Ticket detail edit mode (Tab, Save, Cancel, Esc)
- Field validation rules
- Visual indicators (border colors, sync status)

---

## Verification Checklist

Use this checklist to verify all Week 11-13 features are working:

### Week 11: TUI Skeleton ✓
- [ ] `./ticketr tui` launches without errors
- [ ] TUI displays with proper borders and titles
- [ ] `q` quits cleanly
- [ ] `Ctrl+C` quits immediately
- [ ] `?` shows help view
- [ ] `?` from help view returns to main view
- [ ] Workspace list shows workspaces
- [ ] Workspace names are visible and selectable

### Week 12: Multi-Panel Layout + Real Data ✓
- [ ] Three panels display side-by-side (workspace | tree | detail)
- [ ] Panel proportions: 30 fixed | 40% | 60%
- [ ] `Tab` cycles focus through all three panels
- [ ] `Shift+Tab` cycles focus backward
- [ ] Border colors change on focus (green=focused, white=unfocused)
- [ ] Selecting workspace with `Enter` loads tickets into tree
- [ ] Ticket tree shows real ticket data (not placeholders)
- [ ] Tickets display with JiraID and Title
- [ ] Tasks appear as children under tickets
- [ ] Vim keys work: `j`/`k` navigate, `h`/`l` expand/collapse
- [ ] Arrow keys still work for backward compatibility
- [ ] "No tickets found" message if workspace has no tickets
- [ ] "Loading tickets..." appears during fetch

### Week 13: Ticket Detail Editor ✓
- [ ] Pressing `Enter` on ticket in tree opens detail view
- [ ] Detail panel shows all ticket fields:
  - [ ] Title (yellow, bold)
  - [ ] Jira ID (green if synced, red if not)
  - [ ] Source Line number
  - [ ] Description
  - [ ] Custom Fields (as key-value pairs)
  - [ ] Acceptance Criteria (numbered list)
  - [ ] Tasks (with their own fields)
- [ ] Read-only mode keybindings work:
  - [ ] `j`/`k` scroll up/down
  - [ ] `g` jumps to top
  - [ ] `G` jumps to bottom
  - [ ] `Esc` returns to ticket tree
- [ ] `e` key enters edit mode
- [ ] Edit mode displays form with all fields
- [ ] Edit mode status bar shows key hints
- [ ] Editing any field shows dirty indicator (*)
- [ ] `Tab` navigates between form fields
- [ ] `Esc` cancels and discards changes
- [ ] Empty Title triggers validation error modal
- [ ] Invalid Jira ID (lowercase) triggers validation error modal
- [ ] Valid Jira ID (PROJECT-123) passes validation
- [ ] Save button saves changes and returns to read-only mode
- [ ] Cancel button discards changes and returns to read-only mode
- [ ] Custom fields can be edited as key=value pairs
- [ ] Acceptance criteria can be edited (one per line)

### Focus Management ✓
- [ ] Context-aware Escape works:
  - [ ] From detail → goes to tree
  - [ ] From tree → goes to workspace
  - [ ] From workspace → does nothing (or quits if configured)
- [ ] Focus always visible (green border)
- [ ] Only one panel has green border at a time

### Error Handling ✓
- [ ] No workspace selected shows helpful message
- [ ] Empty workspace shows "No tickets found" with hint
- [ ] Ticket load failure shows error message
- [ ] Validation errors display in modal
- [ ] Modal can be dismissed with Enter/OK

---

## Troubleshooting

### Issue: `./ticketr v3 enable beta` fails

**Symptoms:**
```
Error: failed to create config directory: permission denied
```

**Solution:**
Check directory permissions:
```bash
ls -la ~/.config/
mkdir -p ~/.config/ticketr
chmod 755 ~/.config/ticketr
```

### Issue: `./ticketr workspace create` fails with keychain error

**Symptoms (macOS):**
```
Error: failed to store credentials: keychain access denied
```

**Solution:**
Grant terminal access to Keychain:
1. Open "Keychain Access" app
2. Right-click on "login" keychain → "Change Settings for Keychain 'login'"
3. Ensure "Lock after X minutes of inactivity" is unchecked (or set high value)
4. Try command again

**Symptoms (Linux):**
```
Error: failed to store credentials: no keyring available
```

**Solution:**
Install keyring daemon:
```bash
# Ubuntu/Debian
sudo apt-get install gnome-keyring

# Fedora
sudo dnf install gnome-keyring
```

### Issue: TUI shows "No tickets found" despite having tickets in Jira

**Cause:** Tickets haven't been synced to local database yet.

**Solution:**

**Option 1: Manual ticket insertion (temporary testing workaround)**

If `ticketr pull` is not implemented yet, insert test tickets directly:

```bash
sqlite3 ~/.local/share/ticketr/ticketr.db
```

```sql
-- Find your workspace ID
SELECT id, name FROM workspaces;

-- Insert a test ticket (replace WORKSPACE-ID with actual ID from above)
INSERT INTO tickets (id, workspace_id, jira_id, title, description, content, source_line, created_at, updated_at)
VALUES (
  'test-ticket-1',
  'WORKSPACE-ID-HERE',
  'BACK-123',
  'Test ticket for TUI',
  'This is a test ticket',
  '{}',
  1,
  datetime('now'),
  datetime('now')
);

-- Verify insertion
SELECT jira_id, title FROM tickets;

.exit
```

**Option 2: Wait for `ticketr pull` implementation**

The `pull` command should sync tickets from Jira. Check if it's implemented:
```bash
./ticketr pull --help
```

### Issue: TUI displays garbled text or no borders

**Cause:** Terminal doesn't support Unicode or color.

**Solution:**
1. Use a modern terminal (iTerm2 on macOS, GNOME Terminal on Linux, Windows Terminal on Windows)
2. Ensure `TERM` environment variable is set:
```bash
echo $TERM
# Should be: xterm-256color or similar
export TERM=xterm-256color
```

### Issue: Build fails with "package not found"

**Symptoms:**
```
package github.com/rivo/tview: cannot find package
```

**Solution:**
```bash
go mod download
go mod tidy
go build -o ticketr ./cmd/ticketr
```

### Issue: TUI crashes on launch

**Symptoms:**
```
panic: runtime error: invalid memory address
```

**Solution:**
1. Check database integrity:
```bash
sqlite3 ~/.local/share/ticketr/ticketr.db "PRAGMA integrity_check;"
# Should output: ok
```

2. Check workspace exists:
```bash
./ticketr workspace list
```

3. Re-enable v3 features:
```bash
./ticketr v3 enable beta
```

4. Check logs (if verbose mode exists):
```bash
./ticketr tui --verbose
```

### Issue: Can't see dirty indicator (*) when editing

**Cause:** May need to actually modify a field, not just focus it.

**Solution:**
1. Press `e` in detail view
2. Actually type or delete characters in a field
3. Look at status bar - should show `*` and "unsaved changes"

### Issue: Validation modal appears but I can't dismiss it

**Workaround:**
1. Press `Tab` to focus the OK button
2. Press `Enter` to dismiss
3. Or press `Esc` (may work depending on implementation)

---

## Advanced Testing

### Test Database Directly

Check what's stored in SQLite:

```bash
sqlite3 ~/.local/share/ticketr/ticketr.db

-- List all tables
.tables

-- View workspaces
SELECT * FROM workspaces;

-- View tickets
SELECT jira_id, title FROM tickets LIMIT 10;

-- Check ticket count by workspace
SELECT w.name, COUNT(t.id) as ticket_count
FROM workspaces w
LEFT JOIN tickets t ON t.workspace_id = w.id
GROUP BY w.name;

.exit
```

### Test Performance

```bash
# Time TUI startup
time ./ticketr tui
# Should be < 1 second

# Check database size
ls -lh ~/.local/share/ticketr/ticketr.db
# Should be reasonable (< 10 MB for < 1000 tickets)
```

### Test Memory Usage

```bash
# In another terminal while TUI is running:
ps aux | grep ticketr
# Check RSS column for memory usage
# Should be < 50 MB for typical usage
```

---

## What's NOT Implemented Yet (Known Limitations)

Based on Week 13 completion:

1. **Actual Jira Sync on Save:** Save button updates in-memory state but doesn't sync to Jira
2. **Pull Command:** May not be fully implemented for v3 workspaces
3. **Task Editing:** Can view tasks in detail but can't edit them individually
4. **Search/Filter:** Week 14 feature (not yet implemented)
5. **Command Palette:** Week 14 feature (not yet implemented)
6. **Markdown Rendering:** Descriptions shown as plain text, not rendered markdown
7. **Attachments:** Not displayed in detail view
8. **Comments:** Not displayed in detail view

---

## Success Criteria

You've successfully tested v3 TUI if:

✅ All items in [Verification Checklist](#verification-checklist) pass
✅ No crashes during normal navigation
✅ Validation works as expected (errors shown for invalid data)
✅ Focus management feels intuitive (Tab/Esc navigation)
✅ Ticket data displays correctly (all fields visible)
✅ Edit mode allows changes and validates input
✅ Help system is comprehensive and accessible

---

## Reporting Issues

If you encounter bugs or unexpected behavior:

1. **Note the exact steps to reproduce**
2. **Check if it's a known limitation** (see section above)
3. **Capture any error messages**
4. **Note your environment:**
   - OS and version
   - Go version
   - Terminal emulator
   - Git commit hash: `git log -1 --oneline`

5. **Create an issue** in the repository with:
   - Clear title (e.g., "TUI crashes when editing ticket with empty description")
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details

---

## Next Steps After Testing

If testing is successful:

1. **Continue to Week 14:** Search, filter, and command palette implementation
2. **Performance profiling:** If you notice slowness
3. **Accessibility testing:** Test with screen readers or different terminal sizes
4. **Integration testing:** Test with real Jira instance with complex ticket structures

---

**Document Version:** 1.0
**Created:** 2025-10-18
**For Phase 4 Weeks 11-13** (TUI skeleton + multi-panel + ticket detail editor)
**Last Tested Commit:** 21d853e
