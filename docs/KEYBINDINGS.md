# Keybindings Reference

Complete reference for all keyboard shortcuts in Ticketr TUI.

## Quick Reference

| Key | Action | Description |
|-----|--------|-------------|
| `?` | Show Help | Display help screen |
| `F1` | Command Palette | Open command palette with fuzzy search |
| `Ctrl+P` | Command Palette | Alternative shortcut to open command palette |
| `F2` | Pull/Sync | Pull latest tickets from Jira |
| `F5` | Refresh | Refresh current workspace tickets |
| `F10` | Exit | Quit application |
| `Ctrl+C` | Force Quit | Force quit application (works in any view) |
| `q` | Quit | Quit application (from main view) |

## Global Keybindings

These keybindings work from any view in the TUI.

| Key | Action | Context | Description |
|-----|--------|---------|-------------|
| `?` | Help | All Views | Display comprehensive help screen |
| `F1` | Command Palette | All Views | Open enhanced command palette with fuzzy search and categories |
| `Ctrl+P` | Command Palette | All Views | Alternative shortcut for command palette |
| `F2` | Pull | All Views | Pull latest tickets from Jira asynchronously |
| `F5` | Refresh | All Views | Refresh current workspace ticket list |
| `F10` | Exit | All Views | Quit application gracefully |
| `Ctrl+C` | Force Quit | All Views | Immediately quit application (also cancels active jobs) |
| `Tab` | Next Panel | Main View | Cycle focus forward through panels |
| `Shift+Tab` | Previous Panel | Main View | Cycle focus backward through panels |
| `Esc` | Back/Cancel | Context-Aware | Navigate back or cancel active operation |

## Workspace List View

Active when workspace list panel is focused.

| Key | Action | Description |
|-----|--------|-------------|
| `Enter` | Select Workspace | Switch to selected workspace and load tickets |
| `Tab` | Next Panel | Move focus to ticket tree panel |
| `n` | New Workspace | Open workspace creation modal |
| `?` | Help | Display help screen |
| `q` | Quit | Exit application |
| `Ctrl+C` | Force Quit | Immediately exit application |

## Ticket Tree View

Active when ticket tree panel is focused.

| Key | Action | Description |
|-----|--------|-------------|
| `Enter` | Open Ticket | Open selected ticket in detail view |
| `Space` | Select/Deselect | Toggle ticket selection for bulk operations |
| `Tab` | Next Panel | Move focus to detail panel or back to workspace list |
| `Esc` | Back | Navigate back to workspace list |
| `j` / `↓` | Navigate Down | Move selection down in ticket tree |
| `k` / `↑` | Navigate Up | Move selection up in ticket tree |
| `h` | Collapse | Collapse current ticket node |
| `l` | Expand | Expand current ticket node |
| `b` | Bulk Operations | Open bulk operations menu for selected tickets |
| `/` | Search | Open search modal to find tickets |
| `:` | Commands | Open command palette |
| `?` | Help | Display help screen |
| `a` | Select All | Select all visible tickets |
| `A` | Deselect All | Clear all ticket selections |
| `p` | Push | Push tickets to Jira |
| `P` | Pull | Pull tickets from Jira |
| `r` | Refresh | Refresh ticket list |
| `s` | Sync | Full sync (pull then push) |

## Ticket Detail View

Active when ticket detail panel is focused.

| Key | Action | Description |
|-----|--------|-------------|
| `Esc` | Back | Return focus to ticket tree |
| `Tab` | Next Panel | Cycle to next panel |
| `e` | Edit | Edit current ticket (future feature) |
| `:` | Commands | Open command palette |
| `?` | Help | Display help screen |

## Modal Views

These keybindings apply when a modal dialog is active (search, command palette, bulk operations, workspace creation).

| Key | Action | Description |
|-----|--------|-------------|
| `Esc` | Close Modal | Close modal and return to main view |
| `Enter` | Confirm | Confirm action or submit form |
| `Ctrl+C` | Force Quit | Quit application (does not close modal) |

### Command Palette Modal

| Key | Action | Description |
|-----|--------|-------------|
| `Esc` | Close | Close command palette |
| `↑` / `↓` | Navigate | Navigate through command list |
| `Enter` | Execute | Execute selected command |
| `Tab` | Toggle Keybindings | Show/hide keybinding hints in command descriptions |
| Type text | Fuzzy Search | Filter commands by name or description |

### Search Modal

| Key | Action | Description |
|-----|--------|-------------|
| `Esc` | Close | Close search modal |
| `↑` / `↓` | Navigate | Navigate through search results |
| `Enter` | Select | Select ticket and open in detail view |
| Type text | Search | Filter tickets by ID, title, or description |

## Syncing Context

Active during async operations (pull, push, sync).

| Key | Action | Description |
|-----|--------|-------------|
| `Esc` | Cancel Operation | Cancel current async operation gracefully |
| `Ctrl+C` | Force Quit | Cancel operation and quit application |

## Context-Aware Behavior

Some keybindings change behavior based on current context:

### Esc Key Priority
1. **During Sync**: Cancel active job
2. **In Detail View**: Navigate to ticket tree
3. **In Ticket Tree**: Navigate to workspace list
4. **In Workspace List**: No action

### Tab Key Behavior
Cycles through panels in this order:
1. Workspace List
2. Ticket Tree
3. Ticket Detail
4. (Back to Workspace List)

## Discovering Keybindings

The TUI provides multiple ways to discover available keybindings:

1. **Action Bar**: Always visible at bottom of screen, shows context-aware keybindings
2. **Command Palette**: Press `F1` or `Ctrl+P` to see all commands with descriptions and keybindings
3. **Help Screen**: Press `?` for comprehensive help including all keybindings
4. **Status Messages**: Active operations show available actions (e.g., "Press Esc to cancel")

## Customization

Keybindings are currently not customizable. All bindings follow these conventions:

- **F-keys**: Standard application functions (F1=Help, F2=Sync, F5=Refresh, F10=Exit)
- **Ctrl+Key**: Global shortcuts that work anywhere
- **Single letters**: Context-specific actions in focused view
- **Vim-style**: Navigation uses `h/j/k/l` where applicable
- **Standard conventions**: `Esc` = back/cancel, `Enter` = confirm/select, `Tab` = cycle focus

## Terminal Compatibility

All keybindings work in modern terminal emulators:
- iTerm2 (macOS)
- Terminal.app (macOS)
- Alacritty
- Windows Terminal
- GNOME Terminal
- Konsole

**Note**: Some legacy terminals may not support F-keys or Ctrl combinations. In such cases, use alternative keybindings where available (e.g., `?` instead of `F1` for help).

## Related Documentation

- [TUI Guide](TUI-GUIDE.md) - Comprehensive TUI usage guide
- [Bulk Operations Guide](bulk-operations-guide.md) - Bulk operations workflow
- [Workspace Management Guide](workspace-management-guide.md) - Workspace operations
