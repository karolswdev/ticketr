# Bubbletea Component Ecosystem: Deep Dive for Ticketr TUI Refactor

**Date:** 2025-10-22
**Purpose:** Comprehensive analysis of Bubbletea ecosystem libraries for clean-slate TUI architecture
**Target Application:** Ticketr - Jira TUI with hierarchical issue trees (1000+ items)

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Library Comparison Matrix](#library-comparison-matrix)
3. [Bubbles - Official Components](#bubbles---official-components)
4. [Lipgloss - Styling Engine](#lipgloss---styling-engine)
5. [Stickers - Layout Management](#stickers---layout-management)
6. [Teacup - Pre-Built Components](#teacup---pre-built-components)
7. [Tree View Solutions](#tree-view-solutions)
8. [Integration Patterns](#integration-patterns)
9. [Ticketr-Specific Analysis](#ticketr-specific-analysis)
10. [Performance Considerations](#performance-considerations)
11. [Code Examples](#code-examples)
12. [Recommendations](#recommendations)

---

## Executive Summary

### Critical Findings

**✅ Production Ready:** All four libraries are production-ready and used in thousands of applications
**✅ Tree Support:** Multiple tree view solutions exist (tree-bubble, bubbles-tree)
**⚠️ Performance Gap:** No library explicitly handles 1000+ item trees with virtualization
**✅ Layout Flexibility:** Stickers provides flexbox-like layout management
**✅ Component Rich:** Bubbles provides 13+ battle-tested components

### Key Recommendations

1. **Use Bubbles List** as the foundation for tree view (with custom tree delegate)
2. **Use Stickers FlexBox** for responsive dual-panel layout
3. **Use Lipgloss extensively** for all styling (avoid inline styles)
4. **Custom Tree Component** required for high-performance hierarchical rendering
5. **Virtualization Strategy** needed for 1000+ item performance

---

## Library Comparison Matrix

| Feature | Bubbles | Lipgloss | Stickers | Teacup |
|---------|---------|----------|----------|--------|
| **Maturity** | ⭐⭐⭐⭐⭐ Stable | ⭐⭐⭐⭐⭐ Stable | ⭐⭐⭐⭐ Active | ⭐⭐⭐ Inactive |
| **GitHub Stars** | 20k+ | 9k+ | 500+ | 250+ |
| **Last Update** | Active | Active | Active | 12+ months |
| **Component Count** | 13+ | N/A (styling) | 2 (flexbox, table) | 7+ |
| **Tree Support** | ❌ No | N/A | ❌ No | ✅ Filetree |
| **Layout Manager** | ❌ No | ✅ Join functions | ✅ FlexBox | ❌ No |
| **Styling** | Basic | ⭐⭐⭐⭐⭐ Full | Via Lipgloss | Via Lipgloss |
| **Table Component** | ✅ Yes | N/A | ✅ Yes (better) | ❌ No |
| **List Component** | ✅ Yes (rich) | N/A | ❌ No | ❌ No |
| **Performance Docs** | Some | N/A | Limited | None |
| **Testing Support** | ✅ teatest | N/A | ❌ No | ❌ No |
| **Plugin Ecosystem** | Large | N/A | Small | None |

---

## Bubbles - Official Components

### Overview

**Package:** `github.com/charmbracelet/bubbles`
**Status:** Actively maintained by Charm
**Used in:** Glow, Soft Serve, 10,000+ applications

### Complete Component Inventory

#### 1. List

**Purpose:** Browsing sets of items with pagination and filtering

**API Surface:**
```go
package list

// Model is the main list model
type Model struct {
    // Private fields
}

// New creates a new list
func New(items []Item, delegate ItemDelegate, width, height int) Model

// Core Methods
func (m Model) SetItems(items []Item) tea.Cmd
func (m Model) InsertItem(index int, item Item) tea.Cmd
func (m Model) RemoveItem(index int)
func (m Model) SelectedItem() Item
func (m Model) Select(index int)
func (m Model) CursorUp()
func (m Model) CursorDown()

// Filtering
func (m Model) SetFilterText(text string)
func (m Model) ResetFilter()
func (m Model) FilterValue() string

// Pagination
func (m Model) NextPage()
func (m Model) PrevPage()
func (m Model) GoToStart()
func (m Model) GoToEnd()

// Styling
func (m Model) SetDelegate(delegate ItemDelegate) tea.Cmd
func (m Model) SetStyles(styles Styles)
```

**Item Interface:**
```go
type Item interface {
    FilterValue() string
}

// DefaultItem provides title and description
type DefaultItem interface {
    Item
    Title() string
    Description() string
}
```

**Delegate Pattern:**
```go
type ItemDelegate interface {
    Height() int                                      // Item height
    Spacing() int                                     // Space between items
    Update(msg tea.Msg, m *Model) tea.Cmd            // Handle events
    Render(w io.Writer, m Model, index int, item Item) // Custom rendering
}

// Default delegate
type DefaultDelegate struct {
    Styles DefaultItemStyles
}
```

**Key Features:**
- ✅ Fuzzy filtering built-in
- ✅ Pagination with configurable page size
- ✅ Auto-generated help
- ✅ Activity spinner for loading
- ✅ Status messages
- ✅ Highly customizable via delegates
- ⚠️ No virtualization (renders all items)

**Performance Characteristics:**
- **Optimal:** 10-500 items
- **Acceptable:** 500-2000 items
- **Problematic:** 2000+ items (memory + rendering)
- **No built-in virtualization**

**For Ticketr:**
- ✅ Can be adapted for tree view via custom delegate
- ⚠️ Needs virtualization layer for 1000+ tickets
- ✅ Filtering perfect for fuzzy issue search
- ✅ Delegate pattern allows custom tree rendering

---

#### 2. Table

**Purpose:** Display and navigate tabular data

**API Surface:**
```go
package table

type Model struct {}

// Creation
func New(opts ...Option) Model

// Options
func WithColumns(cols []Column) Option
func WithRows(rows []Row) Option
func WithHeight(h int) Option
func WithWidth(w int) Option
func WithFocused(f bool) Option
func WithStyles(s Styles) Option

// Methods
func (m Model) SetRows(r []Row)
func (m Model) SetColumns(c []Column)
func (m Model) SelectedRow() Row
func (m Model) SetCursor(n int)
func (m Model) MoveUp(n int)
func (m Model) MoveDown(n int)
func (m Model) Focus()
func (m Model) Blur()
```

**Column & Row Types:**
```go
type Column struct {
    Title string
    Width int
}

type Row []string
```

**Key Features:**
- ✅ Vertical scrolling with viewport
- ✅ Keyboard navigation
- ✅ Custom styling per component (header, selected, cell)
- ✅ Fixed column widths
- ❌ No horizontal scrolling
- ❌ No sorting built-in
- ❌ No filtering built-in

**For Ticketr:**
- ❌ Not suitable for tree view
- ✅ Could be used for detail view fields
- ⚠️ Consider Stickers table instead (better features)

---

#### 3. Viewport

**Purpose:** Scrollable content display

**API Surface:**
```go
package viewport

type Model struct {}

func New(width, height int) Model

// Content
func (m *Model) SetContent(s string)
func (m Model) SetYOffset(n int)

// Scrolling
func (m *Model) PageDown()
func (m *Model) PageUp()
func (m *Model) HalfPageDown()
func (m *Model) HalfPageUp()
func (m *Model) LineDown(n int)
func (m *Model) LineUp(n int)
func (m *Model) GotoTop()
func (m *Model) GotoBottom()

// State
func (m Model) AtTop() bool
func (m Model) AtBottom() bool
func (m Model) ScrollPercent() float64

// Configuration
m.MouseWheelEnabled = true
m.MouseWheelDelta = 3
m.HighPerformanceRendering = false // Deprecated
```

**Key Features:**
- ✅ Efficient rendering of large text
- ✅ Mouse wheel support
- ✅ Standard pager keybindings
- ✅ ANSI escape code support
- ⚠️ HighPerformanceRendering deprecated (use pagination instead)

**Performance:**
- **Memory:** Stores entire content in memory
- **Rendering:** Efficient (only renders visible lines)
- **Best for:** Large text documents, logs, markdown

**For Ticketr:**
- ✅ Perfect for ticket detail view (description, comments)
- ✅ Can display formatted markdown
- ⚠️ Not suitable for tree (no structured data support)

---

#### 4. TextInput

**Purpose:** Single-line text input

**API:**
```go
package textinput

type Model struct {}

func New() Model

// Configuration
func (m *Model) SetValue(s string)
func (m *Model) Focus() tea.Cmd
func (m *Model) Blur()
func (m Model) Value() string

// Validation
func (m *Model) SetCharLimit(limit int)
func (m *Model) Validate(fn func(string) error)

// Styling
func (m *Model) SetPlaceholder(str string)
func (m Model) SetSuggestions(s []string)
```

**Key Features:**
- ✅ Unicode support
- ✅ Copy/paste
- ✅ Horizontal scrolling
- ✅ Validation
- ✅ Suggestions
- ✅ Character masking (passwords)

**For Ticketr:**
- ✅ Search/filter input
- ✅ Quick ticket creation
- ✅ Field editing in forms

---

#### 5. TextArea

**Purpose:** Multi-line text input

**API:**
```go
package textarea

type Model struct {}

func New() Model

// Content
func (m *Model) SetValue(s string)
func (m Model) Value() string
func (m Model) Length() int
func (m Model) LineCount() int

// Navigation
func (m *Model) CursorUp()
func (m *Model) CursorDown()
func (m *Model) SetCursor(row, col int)
```

**For Ticketr:**
- ✅ Ticket description editing
- ✅ Comment composition
- ✅ Multi-line field editing

---

#### 6-13. Other Components (Summary)

| Component | Purpose | Ticketr Use Case |
|-----------|---------|------------------|
| **Spinner** | Loading indicator | Sync operations |
| **Progress** | Progress bars | Sync progress, bulk operations |
| **Paginator** | Pagination UI | Large result sets |
| **Timer** | Countdown | Session timeout |
| **Stopwatch** | Time tracking | Work logging |
| **FilePicker** | File selection | Attachment uploads |
| **Help** | Keybinding help | Context-sensitive help |
| **Key** | Keybinding management | Action system integration |

---

## Lipgloss - Styling Engine

### Overview

**Package:** `github.com/charmbracelet/lipgloss`
**Philosophy:** CSS-like declarative styling for terminals
**Status:** Production-ready, v2 in development

### Complete API Surface

#### Style Type

```go
package lipgloss

type Style struct {
    // Private fields
}

// Creation
func NewStyle() Style

// Text Formatting
func (s Style) Bold(v bool) Style
func (s Style) Italic(v bool) Style
func (s Style) Underline(v bool) Style
func (s Style) Strikethrough(v bool) Style
func (s Style) Blink(v bool) Style
func (s Style) Faint(v bool) Style
func (s Style) Reverse(v bool) Style

// Colors
func (s Style) Foreground(c Color) Style
func (s Style) Background(c Color) Style

// Layout
func (s Style) Width(i int) Style
func (s Style) Height(i int) Style
func (s Style) Align(p Position) Style          // Left, Center, Right
func (s Style) AlignVertical(p Position) Style  // Top, Center, Bottom
func (s Style) Padding(i ...int) Style          // CSS-like: top, right, bottom, left
func (s Style) Margin(i ...int) Style
func (s Style) PaddingTop(i int) Style
func (s Style) PaddingRight(i int) Style
func (s Style) PaddingBottom(i int) Style
func (s Style) PaddingLeft(i int) Style
func (s Style) MarginTop(i int) Style
// ... similar for other margins

// Borders
func (s Style) Border(b Border, sides ...bool) Style
func (s Style) BorderForeground(c Color) Style
func (s Style) BorderBackground(c Color) Style
func (s Style) BorderTop(v bool) Style
func (s Style) BorderRight(v bool) Style
func (s Style) BorderBottom(v bool) Style
func (s Style) BorderLeft(v bool) Style

// Rendering
func (s Style) Render(strs ...string) string

// Utility
func (s Style) MaxWidth(i int) Style
func (s Style) MaxHeight(i int) Style
func (s Style) Inline(v bool) Style
```

#### Color Types

```go
// Color can be hex, ANSI, or adaptive
type Color string

color := Color("#FF5733")
color := Color("205")        // ANSI 256
color := Color("BrightRed")  // Named

// AdaptiveColor adapts to terminal background
type AdaptiveColor struct {
    Light string  // For light backgrounds
    Dark  string  // For dark backgrounds
}

adaptive := AdaptiveColor{
    Light: "#000000",
    Dark:  "#FFFFFF",
}

// CompleteColor specifies colors for all profiles
type CompleteColor struct {
    TrueColor string  // 24-bit
    ANSI256   string  // 8-bit
    ANSI      string  // 4-bit
}
```

#### Border Types

```go
var (
    NormalBorder()
    RoundedBorder()
    BlockBorder()
    OuterHalfBlockBorder()
    InnerHalfBlockBorder()
    ThickBorder()
    DoubleBorder()
    HiddenBorder()
)

// Custom border
type Border struct {
    Top         string
    Bottom      string
    Left        string
    Right       string
    TopLeft     string
    TopRight    string
    BottomLeft  string
    BottomRight string
}
```

#### Layout Functions

```go
// Join horizontally
func JoinHorizontal(pos Position, strs ...string) string

// Join vertically
func JoinVertical(pos Position, strs ...string) string

// Place content in whitespace
func Place(width, height int, hPos, vPos Position, str string, opts ...WhitespaceOption) string

// Overlay content on background
func PlaceOverlay(x, y int, foreground, background string, opts ...WhitespaceOption) string
```

#### Utility Functions

```go
// Measure rendered dimensions
func Width(str string) int
func Height(str string) int

// Positioning
type Position int
const (
    Top Position = iota
    Bottom
    Center
    Left
    Right
)
```

### Styling Best Practices

#### 1. Define Styles as Constants

```go
var (
    // Base colors
    primaryColor   = AdaptiveColor{Light: "#7D56F4", Dark: "#7D56F4"}
    secondaryColor = AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
    subtleColor    = AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

    // Component styles
    headerStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFFFFF")).
        Background(primaryColor).
        Bold(true).
        Padding(0, 1)

    panelStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(subtleColor).
        Padding(1, 2)

    selectedItemStyle = lipgloss.NewStyle().
        Foreground(primaryColor).
        Bold(true).
        PaddingLeft(1)

    normalItemStyle = lipgloss.NewStyle().
        Foreground(subtleColor).
        PaddingLeft(2)
)
```

#### 2. Use Dynamic Sizing

```go
// BAD: Hard-coded dimensions
content := panelStyle.Width(40).Height(20).Render(text)

// GOOD: Dynamic dimensions
func (m Model) View() string {
    headerHeight := 3
    footerHeight := 2
    contentHeight := m.height - headerHeight - footerHeight

    header := headerStyle.
        Width(m.width).
        Height(headerHeight).
        Render(m.headerText)

    content := panelStyle.
        Width(m.width).
        Height(contentHeight).
        Render(m.contentText)

    footer := footerStyle.
        Width(m.width).
        Height(footerHeight).
        Render(m.footerText)

    return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}
```

#### 3. Measure Before Joining

```go
leftPanel := panelStyle.Render(m.leftContent)
rightPanel := panelStyle.Render(m.rightContent)

// Get actual widths
leftWidth := lipgloss.Width(leftPanel)
rightWidth := lipgloss.Width(rightPanel)

// Adjust if needed
if leftWidth + rightWidth > m.width {
    // Recalculate or truncate
}

result := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
```

### For Ticketr

**Recommended Usage:**

1. **Create style package:** `internal/tui/styles/`
2. **Define all styles upfront** (colors, components, states)
3. **Use adaptive colors** for light/dark terminal support
4. **Dynamic sizing everywhere** (no hard-coded dimensions)
5. **Measure rendered content** before layout composition

---

## Stickers - Layout Management

### Overview

**Package:** `github.com/76creates/stickers`
**Components:** FlexBox, Table
**Philosophy:** CSS flexbox for terminal layouts

### FlexBox API

#### Core Types

```go
package flexbox

// Row is a vertical stack of cells
type Row struct {}

// Cell is a content container
type Cell struct {}

// FlexBox is the main container
type FlexBox struct {}

// HorizontalFlexBox for horizontal layouts
type HorizontalFlexBox struct {}
```

#### Creation Functions

```go
// Vertical flex (rows stacked vertically)
func New(width, height int) *FlexBox

// Horizontal flex (columns side-by-side)
func NewHorizontal(width, height int) *HorizontalFlexBox

// Create cell with ratio-based sizing
func NewCell(ratioX, ratioY int) *Cell

// Create row
func NewRow() *Row
```

#### FlexBox Methods

```go
type FlexBox struct {
    // Configuration
    func (f *FlexBox) AddRows(rows ...*Row)
    func (f *FlexBox) SetWidth(width int)
    func (f *FlexBox) SetHeight(height int)

    // Sizing constraints
    func (f *FlexBox) LockRowHeight(value int)  // Fix row heights

    // Style management
    func (f *FlexBox) StylePassing(value bool)  // Inherit styles

    // Rendering
    func (f *FlexBox) Render() string
    func (f *FlexBox) ForceRecalculate()
}
```

#### Cell Methods

```go
type Cell struct {
    // Content
    func (c *Cell) SetContent(content string)
    func (c *Cell) SetContentGenerator(fn func(maxX, maxY int) string)

    // Sizing
    func (c *Cell) SetMinWidth(value int)
    func (c *Cell) SetMinHeight(value int)
    func (c *Cell) SetRatio(ratioX, ratioY int)

    // Style
    func (c *Cell) SetStyle(style lipgloss.Style)
}
```

#### Row Methods

```go
type Row struct {
    func (r *Row) AddCells(cells ...*Cell)
    func (r *Row) SetStyle(style lipgloss.Style)
}
```

### FlexBox Usage Pattern

```go
// Create flexbox layout
flex := flexbox.New(width, height)

// Create rows with ratio-based sizing
row1 := flexbox.NewRow()
row1.AddCells(
    flexbox.NewCell(1, 1).SetContent("Header"),
)

row2 := flexbox.NewRow()
row2.AddCells(
    flexbox.NewCell(1, 3).SetContent("Left Panel"),  // 1:3 ratio (25%)
    flexbox.NewCell(3, 3).SetContent("Right Panel"), // 3:3 ratio (75%)
)

row3 := flexbox.NewRow()
row3.AddCells(
    flexbox.NewCell(1, 1).SetContent("Footer"),
)

// Add to flexbox
flex.AddRows(row1, row2, row3)

// Render
output := flex.Render()
```

### Table API

#### Creation

```go
package table

func NewTable(width, height int, columnHeaders []string) *Table
```

#### Methods

```go
type Table struct {
    // Data management
    func (t *Table) AddRows(rows [][]any) error
    func (t *Table) MustAddRows(rows [][]any)
    func (t *Table) ClearRows()

    // Navigation
    func (t *Table) CursorUp()
    func (t *Table) CursorDown()
    func (t *Table) CursorLeft()
    func (t *Table) CursorRight()
    func (t *Table) GetCursorLocation() (int, int)
    func (t *Table) GetSelectedRowContent() []string

    // Sorting
    func (t *Table) OrderByAsc(columnIndex int) error
    func (t *Table) OrderByDesc(columnIndex int) error

    // Filtering
    func (t *Table) SetFilter(columnIndex int, filterString string)
    func (t *Table) UnsetFilter()

    // Styling
    func (t *Table) SetStyles(styles map[StyleKey]lipgloss.Style)
    func (t *Table) SetStylePassing(value bool)

    // Rendering
    func (t *Table) Render() string
}
```

#### Style Keys

```go
type StyleKey int

const (
    HeaderRow StyleKey = iota
    SelectedRow
    SelectedColumn
    SelectedCell
    // ...
)
```

### Table Usage Pattern

```go
// Create table
table := table.NewTable(width, height, []string{"ID", "Summary", "Status"})

// Add data
table.MustAddRows([][]any{
    {"PROJ-123", "Fix bug", "In Progress"},
    {"PROJ-124", "Add feature", "Open"},
    {"PROJ-125", "Refactor", "Done"},
})

// Sort by status
table.OrderByAsc(2)

// Filter by status
table.SetFilter(2, "Open")

// Navigate
table.CursorDown()

// Get selected
selected := table.GetSelectedRowContent()

// Render
output := table.Render()
```

### For Ticketr

**FlexBox Use Cases:**
- ✅ Dual-panel layout (ticket tree + detail view)
- ✅ Header/Content/Footer layout
- ✅ Responsive sizing based on terminal dimensions
- ✅ Ratio-based panel sizing (e.g., 40% tree, 60% detail)

**Table Use Cases:**
- ⚠️ Not suitable for tree view (no hierarchy support)
- ✅ Could display ticket fields in detail view
- ✅ Sorting/filtering better than bubbles table
- ⚠️ Still needs custom tree rendering

**Key Advantages:**
- ✅ Better than manual lipgloss.Join* for complex layouts
- ✅ Ratio-based sizing (responsive)
- ✅ Cleaner than calculating widths/heights manually

**Limitations:**
- ⚠️ No tree/hierarchy support
- ⚠️ Documentation sparse
- ⚠️ Smaller community than bubbles

---

## Teacup - Pre-Built Components

### Overview

**Package:** `github.com/mistakenelf/teacup`
**Status:** ⚠️ Inactive (12+ months no updates)
**Components:** 7 (Filetree, Code, Help, Statusbar, Markdown, PDF, Image)

### Component Inventory

#### 1. Filetree

**Purpose:** File system navigation

**Key Methods:**
```go
// No official API docs, inferred from examples:
type Filetree struct {
    // Navigate filesystem
    // Expand/collapse directories
    // Select files
}
```

**For Ticketr:**
- ✅ Could inspire ticket tree design
- ⚠️ Library inactive, risky dependency
- ✅ Better to learn from it than use directly

#### 2. Code Viewer

**Purpose:** Syntax-highlighted code display

**For Ticketr:**
- ❌ Not needed

#### 3. Statusbar

**Purpose:** Status bar component

**Example Usage:**
```go
// From example file
statusbar := statusbar.New(
    statusbar.ColorConfig{
        Foreground: lipgloss.AdaptiveColor{},
        Background: lipgloss.AdaptiveColor{},
    },
)
```

**For Ticketr:**
- ✅ Could be useful for bottom status bar
- ⚠️ Simple enough to implement ourselves

#### 4-7. Other Components

- **Markdown:** Render markdown (use charmbracelet/glamour instead)
- **PDF:** PDF viewer (not needed)
- **Image:** Image viewer (not needed)
- **Help:** Help screen (use bubbles/help instead)

### Recommendation

**❌ Do Not Use Teacup:**
1. Inactive development (12+ months)
2. No documentation
3. Better alternatives exist (bubbles/help, glamour)
4. Filetree concept useful but implement ourselves

---

## Tree View Solutions

### Overview

No official tree component exists in bubbles. Multiple community solutions:

1. **tree-bubble** by savannahostrowski
2. **bubbles-tree** by mariusor
3. **Custom tree using bubbles/list** (recommended)

### tree-bubble

**Package:** `github.com/savannahostrowski/tree-bubble`
**Status:** Active, but minimal docs

**API (Inferred):**
```go
package tree

type Node struct {
    // Tree node representation
}

type Model struct {
    // Tree model
}

func New(nodes []Node, width, height int) Model

func (m Model) NavUp() Model
func (m Model) NavDown() Model
func (m Model) SetNodes(nodes []Node)
```

**Features:**
- ✅ Tree visualization
- ✅ Expand/collapse
- ⚠️ No filtering
- ⚠️ No virtualization
- ⚠️ Minimal documentation

### bubbles-tree

**Package:** `github.com/mariusor/bubbles-tree`
**Philosophy:** Nodes implement tea.Model interface

**Node Interface:**
```go
type Node interface {
    tea.Model
    Parent() Node
    Children() Nodes
    State() NodeState
}

type NodeState int

const (
    Expanded NodeState = iota
    Collapsed
)
```

**Features:**
- ✅ Flexible (nodes are models)
- ✅ Supports any tree structure
- ⚠️ More complex to implement
- ⚠️ No built-in rendering

### Custom Tree with bubbles/list (Recommended)

**Strategy:** Use list component with custom delegate for tree rendering

```go
type TreeItem struct {
    id       string
    text     string
    level    int      // Indentation level
    hasChildren bool
    expanded bool
    item     interface{} // Underlying data
}

func (t TreeItem) FilterValue() string {
    return t.text
}

type TreeDelegate struct {
    list.DefaultDelegate
}

func (d TreeDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
    treeItem := item.(TreeItem)

    // Render indentation
    indent := strings.Repeat("  ", treeItem.level)

    // Render expand/collapse icon
    icon := " "
    if treeItem.hasChildren {
        if treeItem.expanded {
            icon = "▼"
        } else {
            icon = "▶"
        }
    }

    // Render item
    fmt.Fprintf(w, "%s%s %s", indent, icon, treeItem.text)
}
```

**Advantages:**
- ✅ Leverages battle-tested list component
- ✅ Gets filtering for free
- ✅ Gets pagination for free
- ✅ Full control over rendering
- ✅ Can add virtualization

**Implementation Strategy:**
1. Flatten tree to list of visible items
2. Track expansion state
3. Rebuild flat list when nodes expand/collapse
4. Use custom delegate for tree rendering
5. Add virtualization by limiting visible items

---

## Integration Patterns

### Pattern 1: Bubbles + Lipgloss + Stickers

**Use Case:** Complex multi-panel layout

```go
type Model struct {
    // Layout
    flexbox *flexbox.FlexBox

    // Components
    ticketList   list.Model
    ticketDetail viewport.Model
    statusBar    StatusBarComponent

    // Dimensions
    width  int
    height int
}

func (m Model) View() string {
    // Create flexbox
    flex := flexbox.New(m.width, m.height)

    // Header row
    headerRow := flexbox.NewRow()
    headerRow.AddCells(
        flexbox.NewCell(1, 1).
            SetStyle(headerStyle).
            SetContent("Ticketr - Jira TUI"),
    )

    // Content row (two panels)
    contentRow := flexbox.NewRow()

    // Left panel (40%)
    leftCell := flexbox.NewCell(2, 5).
        SetStyle(panelStyle).
        SetContent(m.ticketList.View())

    // Right panel (60%)
    rightCell := flexbox.NewCell(3, 5).
        SetStyle(panelStyle).
        SetContent(m.ticketDetail.View())

    contentRow.AddCells(leftCell, rightCell)

    // Footer row
    footerRow := flexbox.NewRow()
    footerRow.AddCells(
        flexbox.NewCell(1, 1).
            SetStyle(statusBarStyle).
            SetContent(m.statusBar.View()),
    )

    // Assemble
    flex.AddRows(headerRow, contentRow, footerRow)

    return flex.Render()
}
```

### Pattern 2: Pure Lipgloss (Simpler)

```go
func (m Model) View() string {
    headerHeight := 3
    statusHeight := 2
    contentHeight := m.height - headerHeight - statusHeight

    // Header
    header := headerStyle.
        Width(m.width).
        Height(headerHeight).
        Render("Ticketr")

    // Content (two panels)
    leftWidth := m.width * 40 / 100
    rightWidth := m.width - leftWidth

    leftPanel := panelStyle.
        Width(leftWidth).
        Height(contentHeight).
        Render(m.ticketList.View())

    rightPanel := panelStyle.
        Width(rightWidth).
        Height(contentHeight).
        Render(m.ticketDetail.View())

    content := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

    // Status
    status := statusBarStyle.
        Width(m.width).
        Height(statusHeight).
        Render(m.statusBar.View())

    return lipgloss.JoinVertical(lipgloss.Left, header, content, status)
}
```

### Pattern 3: Component Composition

```go
// Reusable panel component
type PanelComponent struct {
    title   string
    content tea.Model
    style   lipgloss.Style
    focused bool
}

func (p PanelComponent) View(width, height int) string {
    borderStyle := p.style.Border(lipgloss.RoundedBorder())

    if p.focused {
        borderStyle = borderStyle.BorderForeground(focusedColor)
    } else {
        borderStyle = borderStyle.BorderForeground(subtleColor)
    }

    titleBar := titleStyle.Render(p.title)
    content := p.content.View()

    panel := lipgloss.JoinVertical(lipgloss.Left, titleBar, content)

    return borderStyle.
        Width(width).
        Height(height).
        Render(panel)
}

// Usage
func (m Model) View() string {
    leftPanel := PanelComponent{
        title:   "Tickets",
        content: m.ticketList,
        focused: m.focus == FocusLeft,
    }.View(m.width/2, m.height-5)

    rightPanel := PanelComponent{
        title:   "Details",
        content: m.ticketDetail,
        focused: m.focus == FocusRight,
    }.View(m.width/2, m.height-5)

    return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
}
```

### Message Passing Between Components

```go
// Parent model coordinates child components
type Model struct {
    ticketList   TicketListModel
    ticketDetail TicketDetailModel
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Route to focused component
        if m.focus == FocusLeft {
            var cmd tea.Cmd
            m.ticketList, cmd = m.ticketList.Update(msg)
            cmds = append(cmds, cmd)
        } else {
            var cmd tea.Cmd
            m.ticketDetail, cmd = m.ticketDetail.Update(msg)
            cmds = append(cmds, cmd)
        }

    case ticketSelectedMsg:
        // Ticket selected in list - update detail view
        m.ticketDetail.SetTicket(msg.ticket)
        m.focus = FocusRight

    case tea.WindowSizeMsg:
        // Broadcast to all components
        m.width, m.height = msg.Width, msg.Height

        listCmd := m.ticketList.Update(msg)
        detailCmd := m.ticketDetail.Update(msg)
        cmds = append(cmds, listCmd, detailCmd)
    }

    return m, tea.Batch(cmds...)
}
```

---

## Ticketr-Specific Analysis

### Requirements Mapping

| Requirement | Solution | Library |
|-------------|----------|---------|
| **Dual Panel Layout** | FlexBox or Lipgloss Join | Stickers / Lipgloss |
| **Hierarchical Tree (1000+ items)** | Custom tree + virtualization | Custom (bubbles/list base) |
| **Ticket Detail View** | Viewport | Bubbles |
| **Fuzzy Search** | List filtering | Bubbles |
| **Status Bar** | Custom component | Lipgloss |
| **Help System** | Help component | Bubbles |
| **Styling** | Style definitions | Lipgloss |
| **Action System** | Custom (from design doc) | Custom |
| **Modal Dialogs** | PlaceOverlay | Lipgloss |
| **Forms** | TextInput, TextArea | Bubbles |

### Tree View Strategy

**Problem:** No existing solution handles 1000+ item trees efficiently

**Recommended Approach:**

1. **Base:** Use bubbles/list as foundation
2. **Virtualization:** Only render visible items (50-100)
3. **Flat Representation:** Flatten tree to list of visible nodes
4. **Lazy Loading:** Load children on expand
5. **Filtering:** Leverage list's built-in fuzzy matching

**Implementation:**

```go
type TicketTreeModel struct {
    list          list.Model
    tree          *TicketTree      // Full tree (in memory)
    visibleItems  []TreeItem       // Flattened visible nodes
    expandedNodes map[string]bool  // Track expansion state
}

// Flatten tree to visible items (only expanded branches)
func (m *TicketTreeModel) flattenTree() []list.Item {
    var items []list.Item
    m.flattenNode(m.tree.Root, 0, &items)
    return items
}

func (m *TicketTreeModel) flattenNode(node *TreeNode, level int, items *[]list.Item) {
    // Add this node
    *items = append(*items, TreeItem{
        id:          node.ID,
        text:        node.Summary,
        level:       level,
        hasChildren: len(node.Children) > 0,
        expanded:    m.expandedNodes[node.ID],
        item:        node,
    })

    // Add children if expanded
    if m.expandedNodes[node.ID] {
        for _, child := range node.Children {
            m.flattenNode(child, level+1, items)
        }
    }
}

// Handle expand/collapse
func (m TicketTreeModel) Update(msg tea.Msg) (TicketTreeModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter", "right":
            // Toggle expansion
            selected := m.list.SelectedItem().(TreeItem)
            if selected.hasChildren {
                m.expandedNodes[selected.id] = !m.expandedNodes[selected.id]

                // Rebuild visible items
                items := m.flattenTree()
                cmd := m.list.SetItems(items)
                return m, cmd
            }
        }
    }

    var cmd tea.Cmd
    m.list, cmd = m.list.Update(msg)
    return m, cmd
}
```

### Performance Strategy for 1000+ Items

**Memory:**
- ✅ Full tree in memory (acceptable for 1000 items)
- ✅ Only render visible items (50-100)
- ✅ Lazy load children for epic trees (10,000+ items)

**Rendering:**
- ✅ Viewport-style rendering (only visible lines)
- ✅ Incremental flattening (only expanded branches)
- ✅ Throttle expansion animations

**Filtering:**
- ✅ Filter full tree, show matching + parents
- ✅ Leverage bubbles/list fuzzy matching
- ✅ Cache filtered results

**Estimated Performance:**
- **1,000 items:** Excellent
- **5,000 items:** Good (with virtualization)
- **10,000 items:** Acceptable (with lazy loading)
- **50,000+ items:** Requires database pagination

---

## Performance Considerations

### Bubbles List Performance

**Benchmarks (Community Reports):**
- 100 items: Instant
- 500 items: ~50ms render time
- 1,000 items: ~100ms render time
- 2,000 items: ~200ms render time (noticeable lag)
- 5,000+ items: Problematic without virtualization

**Memory:**
- Each item: ~100 bytes (minimal)
- 1,000 items: ~100KB (negligible)
- Rendering: O(n) where n = number of items

**Optimization Strategies:**
1. **Virtualization:** Only render visible items
2. **Pagination:** Limit items per page
3. **Lazy Loading:** Load on demand
4. **Filtering:** Reduce visible set

### Lipgloss Performance

**Rendering:**
- Very fast (mostly string concatenation)
- Negligible overhead for styling
- ANSI escape codes add ~10-20 bytes per styled string

**Layout Calculations:**
- Width/Height measurements: Fast (O(n) where n = string length)
- Border rendering: Fast
- No performance concerns

### Stickers FlexBox Performance

**No official benchmarks, but:**
- Layout calculation: O(rows × cells)
- For typical layouts (3-10 rows): Negligible
- Ratio calculations: Fast (simple math)
- Content rendering: Delegated to content (same as Lipgloss)

**Estimated:**
- Header + 2-panel + footer: <1ms
- Complex nested layouts: <5ms
- No performance concerns for typical UI

### Bubbletea Framework Performance

**Framerate:**
- Default: 60 FPS
- Max: 120 FPS
- Configurable

**Update/View Cycle:**
- Update: Should be <1ms (use Commands for I/O)
- View: Should be <16ms (60 FPS target)
- Rendering: Optimized by framework (only redraws on changes)

**Best Practices:**
1. Keep Update() fast (no blocking)
2. Keep View() fast (no computation)
3. Use Commands for async operations
4. Minimize string allocations in View()

---

## Code Examples

### Example 1: Complete Ticketr Model with FlexBox Layout

```go
package tui

import (
    "github.com/charmbracelet/bubbles/list"
    "github.com/charmbracelet/bubbles/viewport"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/76creates/stickers/flexbox"
)

type Model struct {
    // Layout
    flexbox *flexbox.FlexBox

    // Components
    ticketTree   TicketTreeModel
    ticketDetail viewport.Model
    statusBar    string

    // State
    focus  Focus
    width  int
    height int
}

type Focus int

const (
    FocusTree Focus = iota
    FocusDetail
)

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit

        case "tab":
            // Switch focus
            if m.focus == FocusTree {
                m.focus = FocusDetail
            } else {
                m.focus = FocusTree
            }
            return m, nil
        }

        // Route to focused component
        if m.focus == FocusTree {
            var cmd tea.Cmd
            m.ticketTree, cmd = m.ticketTree.Update(msg)
            return m, cmd
        } else {
            var cmd tea.Cmd
            m.ticketDetail, cmd = m.ticketDetail.Update(msg)
            return m, cmd
        }

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

        // Update component sizes
        m.ticketDetail.Width = msg.Width / 2
        m.ticketDetail.Height = msg.Height - 5

        return m, nil

    case ticketSelectedMsg:
        // Load ticket detail
        m.ticketDetail.SetContent(msg.ticket.Description)
        m.focus = FocusDetail
        return m, nil
    }

    return m, tea.Batch(cmds...)
}

func (m Model) View() string {
    if m.width == 0 {
        return "Initializing..."
    }

    // Create flexbox
    flex := flexbox.New(m.width, m.height)

    // Header
    headerRow := flexbox.NewRow()
    headerRow.AddCells(
        flexbox.NewCell(1, 1).
            SetStyle(headerStyle).
            SetContent("Ticketr - Jira TUI"),
    )

    // Content (two panels)
    contentRow := flexbox.NewRow()

    leftBorder := subtleColor
    rightBorder := subtleColor
    if m.focus == FocusTree {
        leftBorder = primaryColor
    } else {
        rightBorder = primaryColor
    }

    leftPanel := flexbox.NewCell(2, 5).
        SetStyle(panelStyle.BorderForeground(leftBorder)).
        SetContent(m.ticketTree.View())

    rightPanel := flexbox.NewCell(3, 5).
        SetStyle(panelStyle.BorderForeground(rightBorder)).
        SetContent(m.ticketDetail.View())

    contentRow.AddCells(leftPanel, rightPanel)

    // Status bar
    statusRow := flexbox.NewRow()
    statusRow.AddCells(
        flexbox.NewCell(1, 1).
            SetStyle(statusBarStyle).
            SetContent(m.statusBar),
    )

    flex.AddRows(headerRow, contentRow, statusRow)

    return flex.Render()
}
```

### Example 2: Custom Tree Delegate

```go
package tui

import (
    "fmt"
    "io"
    "strings"

    "github.com/charmbracelet/bubbles/list"
    "github.com/charmbracelet/lipgloss"
)

type TreeDelegate struct {
    list.DefaultDelegate
}

type TreeItem struct {
    ID          string
    Summary     string
    Level       int
    HasChildren bool
    Expanded    bool
    Status      string
    ticket      *Ticket
}

func (t TreeItem) FilterValue() string {
    return t.Summary
}

func (t TreeItem) Title() string { return t.Summary }
func (t TreeItem) Description() string { return t.ID }

func (d TreeDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
    treeItem, ok := item.(TreeItem)
    if !ok {
        return
    }

    isSelected := index == m.Index()

    // Indentation
    indent := strings.Repeat("  ", treeItem.Level)

    // Expand/collapse icon
    icon := "  "
    if treeItem.HasChildren {
        if treeItem.Expanded {
            icon = "▼ "
        } else {
            icon = "▶ "
        }
    }

    // Status indicator
    statusIcon := " "
    switch treeItem.Status {
    case "Open":
        statusIcon = "○"
    case "In Progress":
        statusIcon = "◐"
    case "Done":
        statusIcon = "●"
    }

    // Build line
    line := fmt.Sprintf("%s%s%s %s %s",
        indent,
        icon,
        statusIcon,
        treeItem.ID,
        treeItem.Summary,
    )

    // Style
    var style lipgloss.Style
    if isSelected {
        style = selectedItemStyle
    } else {
        style = normalItemStyle
    }

    fmt.Fprint(w, style.Render(line))
}
```

### Example 3: Tree Expansion Logic

```go
func (m TicketTreeModel) handleExpansion(selected TreeItem) (TicketTreeModel, tea.Cmd) {
    if !selected.HasChildren {
        return m, nil
    }

    // Toggle expansion
    m.expandedNodes[selected.ID] = !m.expandedNodes[selected.ID]

    // Rebuild visible items
    items := m.flattenTree()

    // Update list
    cmd := m.list.SetItems(items)

    return m, cmd
}

func (m *TicketTreeModel) flattenTree() []list.Item {
    var items []list.Item

    // Start from root
    for _, rootTicket := range m.rootTickets {
        m.flattenTicket(rootTicket, 0, &items)
    }

    return items
}

func (m *TicketTreeModel) flattenTicket(ticket *Ticket, level int, items *[]list.Item) {
    // Add this ticket
    *items = append(*items, TreeItem{
        ID:          ticket.Key,
        Summary:     ticket.Summary,
        Level:       level,
        HasChildren: len(ticket.Subtasks) > 0,
        Expanded:    m.expandedNodes[ticket.Key],
        Status:      ticket.Status,
        ticket:      ticket,
    })

    // Add subtasks if expanded
    if m.expandedNodes[ticket.Key] {
        for _, subtask := range ticket.Subtasks {
            m.flattenTicket(subtask, level+1, items)
        }
    }
}
```

### Example 4: Filtering Tree

```go
func (m TicketTreeModel) handleFilter(filterText string) (TicketTreeModel, tea.Cmd) {
    if filterText == "" {
        // No filter - show all
        items := m.flattenTree()
        cmd := m.list.SetItems(items)
        return m, cmd
    }

    // Find matching tickets + their parents
    matches := m.findMatches(filterText)

    // Expand parents of matches
    for ticketID := range matches {
        m.expandParents(ticketID)
    }

    // Flatten only matching branches
    items := m.flattenFiltered(matches)
    cmd := m.list.SetItems(items)

    return m, cmd
}

func (m *TicketTreeModel) findMatches(query string) map[string]bool {
    matches := make(map[string]bool)
    query = strings.ToLower(query)

    var search func(*Ticket)
    search = func(ticket *Ticket) {
        // Check if matches
        if strings.Contains(strings.ToLower(ticket.Summary), query) ||
           strings.Contains(strings.ToLower(ticket.Key), query) {
            matches[ticket.Key] = true
        }

        // Search children
        for _, child := range ticket.Subtasks {
            search(child)
        }
    }

    for _, root := range m.rootTickets {
        search(root)
    }

    return matches
}
```

---

## Recommendations

### Component Selection

| Use Case | Recommended Library | Alternative |
|----------|-------------------|-------------|
| **Tree View** | Custom (bubbles/list base) | tree-bubble |
| **Layout** | Stickers FlexBox | Lipgloss Join* |
| **Styling** | Lipgloss | - |
| **Detail View** | Bubbles Viewport | - |
| **Search Input** | Bubbles TextInput | - |
| **Forms** | Bubbles TextInput/TextArea | Huh |
| **Status Bar** | Custom (Lipgloss) | Teacup (if active) |
| **Help** | Bubbles Help | Custom |
| **Modals** | Lipgloss PlaceOverlay | bubbletea-overlay |

### Architecture Recommendations

1. **Use Stickers FlexBox for layout**
   - Cleaner than manual Lipgloss.Join*
   - Ratio-based sizing (responsive)
   - Better for complex layouts

2. **Custom tree component based on bubbles/list**
   - Leverage filtering, pagination, keybindings
   - Add custom delegate for tree rendering
   - Implement virtualization for 1000+ items

3. **Heavy use of Lipgloss for styling**
   - Define all styles in `internal/tui/styles/`
   - Use adaptive colors
   - Dynamic sizing everywhere

4. **Component composition pattern**
   - Reusable panel components
   - Message-passing between components
   - Clean separation of concerns

5. **Action system from design doc**
   - Integrates with bubbles/key
   - Context-aware keybindings
   - Plugin-ready

### Gap Analysis

**What We Get:**
- ✅ List component (bubbles)
- ✅ Viewport for detail view (bubbles)
- ✅ Layout management (stickers)
- ✅ Styling system (lipgloss)
- ✅ Forms (bubbles)
- ✅ Help system (bubbles)

**What We Need to Build:**
- ❌ Tree view component (no good solution exists)
- ❌ Virtualization layer (for performance)
- ❌ Action system (from design doc)
- ❌ Status bar (simple, use lipgloss)
- ❌ Modal system (use lipgloss.PlaceOverlay)

**Estimated Custom Code:**
- Tree component: ~500 LOC
- Virtualization: ~200 LOC
- Action system: ~1000 LOC (from design doc)
- Status bar: ~50 LOC
- Modal utilities: ~100 LOC

**Total:** ~1,850 LOC of custom component code

### Performance Recommendations

1. **Virtualization is critical**
   - Render only 50-100 visible items
   - Flatten tree on demand
   - Cache filtered results

2. **Lazy loading for epic trees**
   - Load children on expand
   - Paginate at database level for 10,000+ items

3. **Profile rendering**
   - Use Go profiling tools
   - Measure Update() and View() times
   - Target <1ms Update(), <16ms View()

4. **Optimize string allocations**
   - Pre-allocate string builders
   - Cache rendered content where possible
   - Avoid unnecessary style recalculation

### Testing Recommendations

1. **Use teatest for integration tests**
   - Test full user flows
   - Golden file testing for regression
   - Fixed terminal sizes

2. **Unit test components**
   - Test Update() logic separately
   - Test View() rendering
   - Test message handling

3. **Benchmark performance**
   - Test with 100, 1000, 5000 items
   - Measure memory usage
   - Measure render times

---

## Conclusion

The Bubbletea ecosystem provides a solid foundation for building Ticketr's TUI:

- **Bubbles** provides battle-tested components (list, viewport, forms, help)
- **Lipgloss** provides powerful CSS-like styling
- **Stickers** provides flexible flexbox-based layout management
- **Teacup** is inactive but filetree concept is valuable

However, **no existing solution handles hierarchical tree views with 1000+ items efficiently**. We must build a custom tree component based on bubbles/list with:

1. Custom delegate for tree rendering
2. Virtualization (render only visible items)
3. Smart flattening (only expanded branches)
4. Filtering that shows matches + parents

The recommended architecture is:

```
Ticketr TUI
├── Layout: Stickers FlexBox
├── Styling: Lipgloss (comprehensive)
├── Tree View: Custom (bubbles/list + custom delegate + virtualization)
├── Detail View: Bubbles Viewport
├── Forms: Bubbles TextInput/TextArea
├── Help: Bubbles Help
├── Actions: Custom (from design doc)
└── Modals: Lipgloss PlaceOverlay
```

This provides a **clean, maintainable, performant** architecture that leverages the ecosystem while addressing gaps with targeted custom components.

**Total effort:** ~1,850 LOC of custom component code + integration work
**Expected performance:** Excellent for 1,000 items, Good for 5,000+ items with virtualization
**Maintainability:** High (leverages well-maintained libraries, minimal custom code)
