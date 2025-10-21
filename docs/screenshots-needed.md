# Required Screenshots for TUI Documentation

This document lists screenshots needed for comprehensive visual documentation of Phase 6 TUI enhancements.

## Overview

Screenshots should demonstrate the enhanced TUI features including:
- Context-aware action bar
- Enhanced command palette with fuzzy search
- F-key shortcuts in action
- Async operations with progress indicators
- Context switching and discoverability

## Screenshot Requirements

### 1. Action Bar in Different Contexts

**Purpose**: Demonstrate context-aware keybinding hints

**Screenshot: Workspace List Context**
- File: `action-bar-workspace-list.png`
- Show: Bottom action bar displaying workspace list keybindings
- Highlight: `[Enter Select Workspace] [Tab Next Panel] [n New Workspace] [? Help] [q/Ctrl+C Quit]`
- Notes: Ensure workspace list panel is focused (highlighted border)

**Screenshot: Ticket Tree Context**
- File: `action-bar-ticket-tree.png`
- Show: Bottom action bar displaying ticket tree keybindings
- Highlight: `[Enter Open Ticket] [Space Select/Deselect] [j/k Navigate] [b Bulk Ops]`
- Notes: Ensure ticket tree panel is focused, some tickets visible

**Screenshot: Ticket Detail Context**
- File: `action-bar-ticket-detail.png`
- Show: Bottom action bar displaying detail view keybindings
- Highlight: `[Esc Back] [Tab Next Panel] [e Edit] [: Commands]`
- Notes: Ensure detail panel is focused with ticket information visible

**Screenshot: Syncing Context**
- File: `action-bar-syncing.png`
- Show: Bottom action bar during async operation
- Highlight: `[Esc Cancel Operation] [Ctrl+C Quit]`
- Notes: Capture during active pull/push with progress visible

### 2. Command Palette

**Purpose**: Demonstrate fuzzy search and command discovery

**Screenshot: Command Palette Initial View**
- File: `command-palette-initial.png`
- Show: Command palette opened with all commands visible
- Highlight: Title "Command Palette (Ctrl+P or F1)", categories (Navigation, View, Edit, Sync, System)
- Notes: Input field empty, all commands grouped by category

**Screenshot: Command Palette with Search**
- File: `command-palette-search.png`
- Show: Command palette with search query entered (e.g., "pull")
- Highlight: Filtered results showing "pull" command
- Notes: Demonstrate fuzzy matching with query text visible

**Screenshot: Command Palette with Keybindings**
- File: `command-palette-keybindings.png`
- Show: Command palette with keybinding hints displayed
- Highlight: Commands showing "[P or F2] Pull latest tickets from Jira"
- Notes: Press Tab to toggle keybinding display, show this state

**Screenshot: Command Palette - Category Organization**
- File: `command-palette-categories.png`
- Show: Full command palette scrolled to show multiple categories
- Highlight: Category headers (Navigation, Sync, System sections visible)
- Notes: Demonstrate clear categorization and organization

### 3. F-Key Shortcuts in Action

**Purpose**: Show F-key shortcuts being used

**Screenshot: F1 Help/Command Palette**
- File: `fkey-f1-help.png`
- Show: Command palette opened via F1
- Highlight: Cursor or indicator showing F1 was pressed
- Notes: Could overlay "F1 pressed" text or arrow

**Screenshot: F2 Pull/Sync**
- File: `fkey-f2-sync.png`
- Show: Pull operation initiated with F2
- Highlight: Sync status showing "Pull in progress" with progress bar
- Notes: Show active sync triggered by F2

**Screenshot: F5 Refresh**
- File: `fkey-f5-refresh.png`
- Show: Ticket list being refreshed
- Highlight: Status showing "Refreshing..." or "Reloading tickets"
- Notes: Brief moment of refresh state

**Screenshot: F10 Exit Confirmation**
- File: `fkey-f10-exit.png`
- Show: Exit confirmation or graceful shutdown
- Highlight: Application closing cleanly after F10
- Notes: Could show terminal returning to prompt

### 4. Async Operations and Progress

**Purpose**: Demonstrate non-blocking async operations with progress

**Screenshot: Pull in Progress**
- File: `async-pull-progress.png`
- Show: Active pull operation with progress indicators
- Highlight:
  - Spinner animation (⠋ character)
  - Ticket count (e.g., "45/120 tickets")
  - Percentage (e.g., "37%")
  - Time elapsed (e.g., "12s elapsed")
- Notes: Ensure TUI still navigable (show cursor on different panel)

**Screenshot: Pull with Cancellation Hint**
- File: `async-cancel-hint.png`
- Show: Active operation with cancellation instruction
- Highlight: Status message "Press Esc to cancel" or action bar showing cancel option
- Notes: Demonstrate user can cancel at any time

**Screenshot: Pull Completed Successfully**
- File: `async-pull-complete.png`
- Show: Completed pull operation with success message
- Highlight: Green success indicator, final ticket count, "120 tickets loaded"
- Notes: Show ticket tree updated with new tickets

**Screenshot: Pull with Error**
- File: `async-pull-error.png`
- Show: Failed pull operation with error message
- Highlight: Red error indicator, error message in status bar
- Notes: Demonstrate graceful error handling

### 5. Context Switching Demonstration

**Purpose**: Show how interface updates when switching contexts

**Screenshot: Context Switch Sequence 1 - Workspace List**
- File: `context-switch-1-workspace.png`
- Show: Workspace list focused, action bar shows workspace context
- Highlight: Workspace panel border highlighted, appropriate action bar hints
- Notes: First step in Tab cycling sequence

**Screenshot: Context Switch Sequence 2 - Ticket Tree**
- File: `context-switch-2-tickets.png`
- Show: After pressing Tab, ticket tree focused
- Highlight: Ticket tree border highlighted, action bar updated to tree context
- Notes: Second step in Tab cycling sequence

**Screenshot: Context Switch Sequence 3 - Detail View**
- File: `context-switch-3-detail.png`
- Show: After pressing Tab again, detail view focused
- Highlight: Detail panel border highlighted, action bar shows detail context
- Notes: Third step in Tab cycling sequence

### 6. Multi-Select and Bulk Operations

**Purpose**: Demonstrate multi-select UI and bulk operations

**Screenshot: Multi-Select Active**
- File: `multiselect-active.png`
- Show: Multiple tickets selected with checkboxes
- Highlight:
  - `[x]` checkboxes on selected tickets
  - `[ ]` checkboxes on unselected tickets
  - Title showing "(3 selected)"
  - Border color change (teal/blue)
- Notes: Select 3-5 tickets for clear demonstration

**Screenshot: Bulk Operations Menu**
- File: `bulk-operations-menu.png`
- Show: Bulk operations modal opened with selected tickets
- Highlight: Operation choices (Update Fields, Move Tickets, Delete Tickets)
- Notes: Show count of selected tickets in modal title

### 7. Discoverability Features

**Purpose**: Demonstrate how new users discover features

**Screenshot: First Launch - Default View**
- File: `discoverability-first-launch.png`
- Show: TUI immediately after launch
- Highlight: Action bar at bottom with clear hints
- Notes: Clean state, shows immediately visible guidance

**Screenshot: Help Screen**
- File: `discoverability-help.png`
- Show: Help screen opened via ? or F1
- Highlight: Comprehensive keybinding list and descriptions
- Notes: Show help is easily accessible and informative

**Screenshot: Search Modal**
- File: `discoverability-search.png`
- Show: Search modal opened with `/` key
- Highlight: Search input, results, navigation hints at bottom
- Notes: Demonstrate quick ticket finding

## Screenshot Specifications

### Technical Requirements

**Resolution**:
- Minimum: 1920x1080 (for clarity)
- Recommended: 2560x1440 (for detail)

**Format**:
- PNG (lossless, good for text)
- Alternative: High-quality JPG (if file size is concern)

**Color Depth**:
- 24-bit true color preferred
- Minimum: 256 colors

**Terminal Emulator**:
- Recommended: iTerm2 (macOS), Windows Terminal, Alacritty
- Avoid: Legacy terminals with poor Unicode support

**Font**:
- Monospace font with good Unicode support
- Recommended: JetBrains Mono, Fira Code, Cascadia Code
- Size: Large enough to be readable in documentation (14pt+)

### Capture Guidelines

**Window Size**:
- Full TUI visible (all three panels + action bar)
- Avoid excessive whitespace/padding
- Crop to relevant content

**Content**:
- Use realistic ticket data (not "Test Ticket 1", "Test Ticket 2")
- Example: "PROJ-123 Fix authentication bug", "PROJ-124 Add user profile"
- Keep ticket content professional and meaningful

**Annotations**:
- Add callout boxes or arrows for important features (optional)
- Use consistent annotation style across all screenshots
- Don't obscure important UI elements

**Timing**:
- For progress indicators, capture mid-operation (30-60% complete)
- For animations (spinners), ensure active animation is visible
- For errors, use realistic error messages

## Screenshot Organization

Store screenshots in: `/home/karol/dev/private/ticktr/docs/images/tui/`

Naming convention: `{category}-{description}.png`

Example structure:
```
docs/
└── images/
    └── tui/
        ├── action-bar-workspace-list.png
        ├── action-bar-ticket-tree.png
        ├── command-palette-initial.png
        ├── async-pull-progress.png
        └── ...
```

## Usage in Documentation

Once captured, screenshots should be embedded in:

1. **TUI-GUIDE.md**:
   - Interface layout section (context switching screenshots)
   - Action bar section (different contexts)
   - Command palette section (search and categories)
   - Async operations section (progress indicators)

2. **KEYBINDINGS.md**:
   - Quick visual reference for action bar in each context
   - F-key usage examples

3. **README.md**:
   - Hero screenshot showing full TUI with action bar
   - GIF/animation of command palette in action (optional)

4. **Release Notes**:
   - Before/after comparisons (if applicable)
   - Highlight new features visually

## Accessibility Considerations

Ensure screenshots:
- Have sufficient contrast for visibility
- Include descriptive alt text when embedded
- Are supplemented by text descriptions (not screenshots alone)
- Show clear focus indicators (borders, highlights)

## Future Enhancements

Consider creating:
- **Animated GIFs**: Show workflows in action (command palette search, context switching)
- **Video Walkthrough**: 2-3 minute demo of TUI features
- **Interactive Tutorial**: Embedded in TUI for first-time users

## Completion Checklist

- [ ] All 20+ screenshots captured
- [ ] Screenshots stored in `/docs/images/tui/` directory
- [ ] Screenshots referenced in TUI-GUIDE.md
- [ ] Screenshots referenced in KEYBINDINGS.md
- [ ] Hero screenshot added to README.md
- [ ] Alt text provided for all screenshots
- [ ] Screenshots reviewed for quality and clarity
- [ ] File sizes optimized (compress if needed)

## Notes for Photographer

- Use a workspace with realistic ticket data
- Ensure terminal is maximized for clarity
- Capture in high-DPI display if available
- Take multiple shots of each scenario (choose best)
- Review each screenshot for readability before finalizing
- Test screenshots at documentation display size (ensure text readable)

## Contact

For questions about screenshot requirements or to submit captured screenshots:
- File issue: [GitHub Issues](https://github.com/karolswdev/ticketr/issues)
- Discussion: [GitHub Discussions](https://github.com/karolswdev/ticketr/discussions)
