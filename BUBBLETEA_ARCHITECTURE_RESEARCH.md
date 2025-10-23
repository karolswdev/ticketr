# Bubbletea TUI Architecture Research

## Executive Summary

This document provides comprehensive research on Bubbletea architecture patterns for migrating from Tview to Bubbletea. The goal is to build a world-class TUI inspired by Midnight Commander with high FPS, no rendering artifacts, extensible actions, and potential for visual effects.

**Key Findings:**
- Bubbletea uses The Elm Architecture (Model-View-Update pattern)
- 60 FPS default framerate-based renderer (max 120 FPS)
- Production-ready with hundreds of real-world applications
- Rich ecosystem of component libraries and layout managers
- Strong community support with active development

---

## Table of Contents

1. [Core Architecture](#core-architecture)
2. [State Management](#state-management)
3. [High-Performance Rendering](#high-performance-rendering)
4. [Multi-View & Panel Layouts](#multi-view--panel-layouts)
5. [Event Handling & Message Passing](#event-handling--message-passing)
6. [Component Composition](#component-composition)
7. [Styling with Lipgloss](#styling-with-lipgloss)
8. [Bubbles Component Library](#bubbles-component-library)
9. [Context-Aware Actions](#context-aware-actions)
10. [Modal & Overlay Management](#modal--overlay-management)
11. [Keyboard Navigation](#keyboard-navigation)
12. [Animation & Visual Effects](#animation--visual-effects)
13. [Testing Patterns](#testing-patterns)
14. [Production Applications](#production-applications)
15. [Recommendations](#recommendations)

---

## Core Architecture

### The Elm Architecture Pattern

Bubbletea is based on The Elm Architecture, a purely functional pattern with three core components:

```go
type Model struct {
    // Application state
}

func (m Model) Init() tea.Cmd {
    // Returns initial command for the application to run
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handles incoming events and updates the model
    return m, nil
}

func (m Model) View() string {
    // Renders the UI based on the model
    return "Hello, World!"
}
```

### Key Principles

1. **Immutable State**: The model describes application state; all changes happen in `Update()`
2. **Pure Functions**: Functions cannot change their internal state
3. **Message-Driven**: All interactions are message-driven events
4. **Declarative Rendering**: The `View()` describes the entire UI; Bubbletea handles redrawing

### Critical Rules

- **Any changes to the model should be made in `Update()` and returned immediately** in the first return value
- Straying from this course not only defies the natural order of Bubbletea, it also risks making it slower
- Never use goroutines directly - use Commands for all I/O operations
- Because the view describes the entire UI, you don't have to worry about redrawing logic

---

## State Management

### Model Tree Composition

Build a hierarchical model structure for complex applications:

```go
type RootModel struct {
    currentView tea.Model
    width       int
    height      int
    // Global state
}

type LeftPanelModel struct {
    // Panel-specific state
}

type RightPanelModel struct {
    // Panel-specific state
}
```

### Message Routing Strategies

Handle messages at different levels:

1. **Root Model**: Handles global actions (quit, help, window resize)
2. **Current Model**: Handles primary interactions
3. **Broadcast**: Messages like window resize broadcast to all child models

```go
func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height
        // Broadcast to child models

    case tea.KeyMsg:
        if msg.String() == "q" {
            return m, tea.Quit
        }
        // Route to current view
        return m.currentView.Update(msg)
    }
    return m, nil
}
```

### State Machine Pattern

For complex multi-step processes, use a stage-based state machine:

```go
type Stage struct {
    Name           string
    Action         func() error
    Error          error
    IsComplete     bool
    IsCompleteFunc func() bool
    Reset          func() error
}

type Model struct {
    stages     []Stage
    stageIndex int
}

func (m Model) runStage() tea.Msg {
    if !m.stages[m.stageIndex].IsCompleteFunc() {
        m.stages[m.stageIndex].Error = m.stages[m.stageIndex].Action()
    }
    return stageCompleteMsg{}
}
```

**Benefits:**
- Robust error handling
- Idempotent stage execution
- Support for long-running tasks
- Graceful error recovery
- Skip already completed stages

### Multi-View Screen Switching

Use a root screen model for view transitions:

```go
type RootScreenModel struct {
    model tea.Model  // Current active screen
}

func (m RootScreenModel) SwitchScreen(model tea.Model) (tea.Model, tea.Cmd) {
    m.model = model
    return m.model, m.model.Init()
}
```

**Best Practices:**
- Maintain global state in the root screen
- Use pointer references for flexible screen management
- Initialize screens dynamically via `Init()`
- Each screen implements `Init()`, `Update()`, `View()`

---

## High-Performance Rendering

### Framerate-Based Renderer

Bubbletea uses a `standardRenderer` with performance optimizations:

- **Default FPS**: 60
- **Max FPS**: 120
- **Frame-based updates**: Prevents overloading terminal emulator
- **Smart rendering**: Only redraws when model changes

### Performance Features

1. **High-Performance Rendering Mode**
   - Renderer can exclude ranges of lines for direct writing
   - Useful for scroll-based rendering
   - Deprecated scroll functions: `ScrollUp`, `ScrollDown`, `SyncScrollArea`

2. **Batch Updates**
   ```go
   tea.Batch(cmd1, cmd2, cmd3)  // Execute multiple commands concurrently
   ```

3. **Sequence Commands**
   ```go
   tea.Sequence(cmd1, cmd2, cmd3)  // Execute commands in order
   ```

### Avoiding Rendering Artifacts

**Critical Rules:**
- Never modify model state outside `Update()`
- Return model changes immediately in the first return value
- Keep `Update()` and `View()` fast - offload expensive operations to Commands
- Programs redraw at the frame rate by default

**Common Issues:**
- **Flickering on Windows**: Reported with v0.26.0+ (actively being addressed)
- **Screen tearing**: Typically caused by blocking operations in `Update()` or `View()`
- **Unnecessary redraws**: Bubbletea redraws on every message; consider manual rendering control for static content

### Performance Optimization Tips

1. **Keep Event Loop Fast**
   - Offload expensive operations to `tea.Cmd`
   - Avoid blocking operations in `Update()` and `View()`
   - Use goroutines via Commands for long-running tasks

2. **Viewport High Performance**
   ```go
   vp := viewport.New(width, height)
   vp.HighPerformanceRendering = true  // Bypass normal renderer
   ```
   Note: Now deprecated but useful for content with many ANSI codes

3. **Dynamic Layout Calculations**
   - Use Lipgloss `Height()` and `Width()` methods
   - Avoid hard-coded dimensions
   - Recalculate on `WindowSizeMsg`

---

## Multi-View & Panel Layouts

### Midnight Commander Inspiration

**Dual Panel Design:**
- Vertically split into two panels
- Independent directory navigation
- Tab to switch between panels
- Function keys for common operations
- Command line at bottom
- Menu system accessible via F-keys

### Lipgloss Layout Composition

**Horizontal Layout (Side-by-Side):**
```go
leftPanel := lipgloss.NewStyle().
    Width(m.width / 2).
    Height(m.height).
    Render(m.leftView.View())

rightPanel := lipgloss.NewStyle().
    Width(m.width / 2).
    Height(m.height).
    Render(m.rightView.View())

return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
```

**Vertical Layout (Stacked):**
```go
header := headerStyle.Render(m.headerView())
content := contentStyle.Render(m.contentView())
footer := footerStyle.Render(m.footerView())

return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
```

### Advanced Layout Managers

#### 1. **76creates/stickers** - FlexBox Implementation
```go
import "github.com/76creates/stickers"
```

Features:
- Responsive grid box inspired by CSS flexbox
- Ratios between elements for responsive scaling
- Table components with flexbox behavior

#### 2. **winder/bubblelayout** - Declarative Layout Manager
```go
import "github.com/winder/bubblelayout"
```

Features:
- Declarative API inspired by MiG Layout
- Horizontal and vertical spans
- Grid with components spanning multiple cells
- Translates `tea.WindowSizeMsg` to `bl.BubbleLayoutMsg`
- Components retrieve dimensions using unique ID

Example:
```go
layout := bubblelayout.New(
    bubblelayout.WithColumnTemplate("1fr 2fr 1fr"),
    bubblelayout.WithRowTemplate("auto 1fr auto"),
)
```

#### 3. **treilik/bubbleboxer** - Layout Tree Manager
```go
import "github.com/treilik/bubbleboxer"
```

Features:
- Layout multiple bubbles side-by-side in a layout tree
- Handles complex nested layouts

#### 4. **mieubrisse/teact** - React-like Framework
```go
import "github.com/mieubrisse/teact"
```

Features:
- React-like component/layout framework
- Flexbox layout capabilities
- Component-based architecture

### Responsive Design Pattern

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width, m.height = msg.Width, msg.Height

        // Update child component sizes
        headerHeight := 3
        footerHeight := 3
        contentHeight := m.height - headerHeight - footerHeight

        m.viewport.Width = m.width
        m.viewport.Height = contentHeight

        // Broadcast to all child models
    }
    return m, nil
}
```

**Key Points:**
- `WindowSizeMsg` sent on startup and every terminal resize
- Windows doesn't support resize detection (no SIGWINCH)
- Always handle initial size in `Update()`

### Common Layout Patterns

**Three-Panel Layout (Header/Content/Footer):**
```go
const (
    headerHeight = 3
    footerHeight = 2
)

func (m Model) View() string {
    contentHeight := m.height - headerHeight - footerHeight

    header := headerStyle.Height(headerHeight).Render(m.headerText)
    content := contentStyle.Height(contentHeight).Render(m.content)
    footer := footerStyle.Height(footerHeight).Render(m.footerText)

    return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}
```

**Dual Panel with Splitter:**
```go
func (m Model) View() string {
    leftWidth := m.width / 2 - 1
    rightWidth := m.width / 2

    if m.focusedPanel == leftPanel {
        leftWidth = int(float64(m.width) * 0.6)
        rightWidth = m.width - leftWidth - 1
    }

    left := leftPanelStyle.Width(leftWidth).Render(m.leftPanel.View())
    divider := dividerStyle.Render("│")
    right := rightPanelStyle.Width(rightWidth).Render(m.rightPanel.View())

    return lipgloss.JoinHorizontal(lipgloss.Top, left, divider, right)
}
```

### File Structure for Multi-Panel Apps

Recommended project structure:
```
ticktr/
├── cmd/
│   ├── initialize.go
│   ├── root.go
├── internal/
│   ├── storage/         # Database/persistence
│   ├── tui/             # TUI components
│   │   ├── root.go      # Root model
│   │   ├── panels/      # Panel models
│   │   │   ├── left.go
│   │   │   ├── right.go
│   │   ├── components/  # Reusable UI components
│   │   ├── styles/      # Lipgloss styles
│   │   └── messages/    # Custom message types
├── main.go
```

---

## Event Handling & Message Passing

### Core Message Types

Bubbletea provides built-in message types:

```go
// Keyboard input
type KeyMsg struct {
    Type  KeyType
    Runes []rune
    Alt   bool
}

// Mouse input
type MouseMsg struct {
    X, Y   int
    Type   MouseEventType
    Button MouseButton
}

// Window size changes
type WindowSizeMsg struct {
    Width  int
    Height int
}

// Batch complete
type BatchMsg struct{}
```

### Custom Message Types

Define domain-specific messages:

```go
// Data loaded messages
type dataLoadedMsg struct {
    data []Item
}

type dataLoadErrorMsg struct {
    err error
}

// UI interaction messages
type itemSelectedMsg struct {
    item Item
}

type panelSwitchedMsg struct {
    fromPanel int
    toPanel   int
}

// State transition messages
type stageCompleteMsg struct{}
type screenSwitchMsg struct {
    screen tea.Model
}
```

### Error Handling Pattern

Implement the error interface on messages:

```go
type errMsg struct {
    err error
}

func (e errMsg) Error() string {
    return e.err.Error()
}

// Usage in Update
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case errMsg:
        m.err = msg
        return m, nil
    }
    return m, nil
}
```

### Message Routing in Complex Apps

```go
type RootModel struct {
    activePanel int
    panels      []tea.Model
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab":
            m.activePanel = (m.activePanel + 1) % len(m.panels)
            return m, nil
        case "q", "ctrl+c":
            return m, tea.Quit
        }

    case tea.WindowSizeMsg:
        // Broadcast to all panels
        for i := range m.panels {
            var cmd tea.Cmd
            m.panels[i], cmd = m.panels[i].Update(msg)
            cmds = append(cmds, cmd)
        }
        return m, tea.Batch(cmds...)
    }

    // Route to active panel
    var cmd tea.Cmd
    m.panels[m.activePanel], cmd = m.panels[m.activePanel].Update(msg)
    return m, cmd
}
```

### External Message Injection

Send messages from outside the program loop:

```go
p := tea.NewProgram(initialModel())

// In another goroutine or external context
go func() {
    time.Sleep(5 * time.Second)
    p.Send(refreshDataMsg{})
}()

p.Run()
```

---

## Component Composition

### Reusable Component Pattern

Create self-contained components implementing the three methods:

```go
// Component interface (implicit)
type Component interface {
    Init() tea.Cmd
    Update(tea.Msg) (Component, tea.Cmd)
    View() string
}

// Example: StatusBar component
type StatusBar struct {
    width   int
    message string
    style   lipgloss.Style
}

func NewStatusBar() StatusBar {
    return StatusBar{
        style: lipgloss.NewStyle().
            Background(lipgloss.Color("62")).
            Foreground(lipgloss.Color("230")),
    }
}

func (s StatusBar) Init() tea.Cmd {
    return nil
}

func (s StatusBar) Update(msg tea.Msg) (StatusBar, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        s.width = msg.Width
    case statusMsg:
        s.message = string(msg)
    }
    return s, nil
}

func (s StatusBar) View() string {
    return s.style.Width(s.width).Render(s.message)
}
```

### Composing Multiple Components

```go
type AppModel struct {
    header    HeaderComponent
    leftPanel PanelComponent
    rightPanel PanelComponent
    statusBar StatusBar
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd
    var cmd tea.Cmd

    // Update all components
    m.header, cmd = m.header.Update(msg)
    cmds = append(cmds, cmd)

    m.leftPanel, cmd = m.leftPanel.Update(msg)
    cmds = append(cmds, cmd)

    m.rightPanel, cmd = m.rightPanel.Update(msg)
    cmds = append(cmds, cmd)

    m.statusBar, cmd = m.statusBar.Update(msg)
    cmds = append(cmds, cmd)

    return m, tea.Batch(cmds...)
}

func (m AppModel) View() string {
    return lipgloss.JoinVertical(lipgloss.Left,
        m.header.View(),
        lipgloss.JoinHorizontal(lipgloss.Top,
            m.leftPanel.View(),
            m.rightPanel.View(),
        ),
        m.statusBar.View(),
    )
}
```

### Focus Management Pattern

For components that can receive focus:

```go
type Focusable interface {
    Focus() tea.Cmd
    Blur()
    Focused() bool
}

type TextInputComponent struct {
    input    textinput.Model
    focused  bool
}

func (t *TextInputComponent) Focus() tea.Cmd {
    t.focused = true
    return t.input.Focus()
}

func (t *TextInputComponent) Blur() {
    t.focused = false
    t.input.Blur()
}

func (t TextInputComponent) Focused() bool {
    return t.focused
}
```

### Form with Multiple Inputs

```go
type FormModel struct {
    inputs       []textinput.Model
    focusedInput int
}

func (m FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab":
            m.inputs[m.focusedInput].Blur()
            m.focusedInput = (m.focusedInput + 1) % len(m.inputs)
            return m, m.inputs[m.focusedInput].Focus()

        case "shift+tab":
            m.inputs[m.focusedInput].Blur()
            m.focusedInput--
            if m.focusedInput < 0 {
                m.focusedInput = len(m.inputs) - 1
            }
            return m, m.inputs[m.focusedInput].Focus()
        }
    }

    // Update focused input
    var cmd tea.Cmd
    m.inputs[m.focusedInput], cmd = m.inputs[m.focusedInput].Update(msg)
    return m, cmd
}
```

---

## Styling with Lipgloss

### Core Concepts

Lipgloss takes an expressive, declarative approach to terminal rendering. Users familiar with CSS will feel at home.

**Philosophy:**
- Make assembling TUI views simple and fun
- Focus on building your application, not low-level layout details
- Declarative style definitions

### Basic Styling

```go
var (
    subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
    highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
    special = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
)

var (
    titleStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFFDF5")).
        Background(lipgloss.Color("#25A065")).
        Padding(0, 1)

    listItemStyle = lipgloss.NewStyle().
        PaddingLeft(2).
        Foreground(subtle)

    selectedItemStyle = lipgloss.NewStyle().
        PaddingLeft(1).
        Foreground(highlight).
        Bold(true)

    paginationStyle = lipgloss.NewStyle().
        Foreground(special)
)
```

### Layout Features

**Padding and Margins:**
```go
style := lipgloss.NewStyle().
    Padding(1, 2, 1, 2).      // top, right, bottom, left
    Margin(1)                  // all sides
```

**Text Alignment:**
```go
leftAlign := lipgloss.NewStyle().Align(lipgloss.Left)
centerAlign := lipgloss.NewStyle().Align(lipgloss.Center)
rightAlign := lipgloss.NewStyle().Align(lipgloss.Right)
```

**Dimensions:**
```go
box := lipgloss.NewStyle().
    Width(50).
    Height(10).
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("62"))
```

### Borders

```go
// Border styles
lipgloss.NormalBorder()
lipgloss.RoundedBorder()
lipgloss.BlockBorder()
lipgloss.OuterHalfBlockBorder()
lipgloss.InnerHalfBlockBorder()
lipgloss.ThickBorder()
lipgloss.DoubleBorder()
lipgloss.HiddenBorder()

// Usage
panelStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("62")).
    BorderTop(true).
    BorderBottom(true).
    BorderLeft(true).
    BorderRight(true)
```

### Adaptive Colors

Colors that adapt to terminal background:

```go
var adaptiveColor = lipgloss.AdaptiveColor{
    Light: "#000000",  // Used on light backgrounds
    Dark:  "#FFFFFF",  // Used on dark backgrounds
}

style := lipgloss.NewStyle().Foreground(adaptiveColor)
```

### Compositing (v2 Feature)

Layer-based rendering with positioning:

```go
canvas := lipgloss.NewCanvas(width, height)
canvas.SetString(x, y, "Content", style)
canvas.Render()
```

### Dynamic Sizing

Use methods to get rendered dimensions:

```go
content := "Some text"
rendered := style.Render(content)
width := lipgloss.Width(rendered)
height := lipgloss.Height(rendered)

// Adjust layout based on size
nextElementY := currentY + height
```

### Common Patterns

**Card/Panel:**
```go
cardStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#874BFD")).
    Padding(1, 2).
    Width(40)
```

**Header Bar:**
```go
headerStyle := lipgloss.NewStyle().
    Background(lipgloss.Color("#7D56F4")).
    Foreground(lipgloss.Color("#FFFFFF")).
    Bold(true).
    Padding(0, 1).
    Width(termWidth)
```

**Status Bar:**
```go
statusStyle := lipgloss.NewStyle().
    Background(lipgloss.Color("#3C3C3C")).
    Foreground(lipgloss.Color("#FFFFFF")).
    Padding(0, 1).
    Width(termWidth)
```

**Gradient Progress Bar:**
```go
// From the official examples
func gradientBar(percent float64, width int) string {
    filled := int(float64(width) * percent)
    bar := strings.Repeat("█", filled)
    empty := strings.Repeat("░", width-filled)
    return bar + empty
}
```

---

## Bubbles Component Library

Official component library: `github.com/charmbracelet/bubbles`

### Available Components

#### 1. **Spinner**
Indicates ongoing operations with customizable animation.

```go
import "github.com/charmbracelet/bubbles/spinner"

s := spinner.New()
s.Spinner = spinner.Dot
s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
```

#### 2. **Text Input**
Single-line text input with Unicode, pasting, scrolling.

```go
import "github.com/charmbracelet/bubbles/textinput"

ti := textinput.New()
ti.Placeholder = "Enter text..."
ti.Focus()
ti.CharLimit = 156
ti.Width = 20
```

#### 3. **Text Area**
Multi-line text input with Unicode, pasting, vertical scrolling.

```go
import "github.com/charmbracelet/bubbles/textarea"

ta := textarea.New()
ta.Placeholder = "Enter description..."
ta.Focus()
```

#### 4. **Table**
Display and navigate tabular data with vertical scrolling.

```go
import "github.com/charmbracelet/bubbles/table"

columns := []table.Column{
    {Title: "Rank", Width: 4},
    {Title: "Name", Width: 10},
    {Title: "Score", Width: 5},
}

rows := []table.Row{
    {"1", "Alice", "100"},
    {"2", "Bob", "90"},
}

t := table.New(
    table.WithColumns(columns),
    table.WithRows(rows),
    table.WithFocused(true),
    table.WithHeight(7),
)
```

#### 5. **Progress**
Customizable progress meter with optional animation.

```go
import "github.com/charmbracelet/bubbles/progress"

p := progress.New(progress.WithDefaultGradient())
p.ShowPercentage = true

// Update with value between 0.0 and 1.0
p.SetPercent(0.5)
```

#### 6. **Paginator**
Pagination logic and UI with dot-style or numeric display.

```go
import "github.com/charmbracelet/bubbles/paginator"

p := paginator.New()
p.Type = paginator.Dots
p.PerPage = 10
p.SetTotalPages(100)
```

#### 7. **Viewport**
Vertically scrollable content with pager keybindings and mouse wheel.

```go
import "github.com/charmbracelet/bubbles/viewport"

vp := viewport.New(width, height)
vp.SetContent(longContent)
vp.MouseWheelEnabled = true
vp.MouseWheelDelta = 3
```

**Performance Note:**
- `HighPerformanceRendering` option available (now deprecated)
- Can use significant memory for large content with ANSI codes
- Consider pagination for very large datasets

#### 8. **List**
Comprehensive browsing component with pagination, filtering, help.

```go
import "github.com/charmbracelet/bubbles/list"

type item struct {
    title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

items := []list.Item{
    item{title: "Item 1", desc: "Description 1"},
    item{title: "Item 2", desc: "Description 2"},
}

l := list.New(items, list.NewDefaultDelegate(), 0, 0)
l.Title = "My List"
```

#### 9. **File Picker**
Select files from filesystem with directory navigation.

```go
import "github.com/charmbracelet/bubbles/filepicker"

fp := filepicker.New()
fp.AllowedTypes = []string{".md", ".txt"}
fp.CurrentDirectory, _ = os.UserHomeDir()
```

#### 10. **Timer**
Countdown component with customizable update frequency.

```go
import "github.com/charmbracelet/bubbles/timer"

t := timer.NewWithInterval(5*time.Minute, time.Second)
```

#### 11. **Stopwatch**
Counting up component with customizable format.

```go
import "github.com/charmbracelet/bubbles/stopwatch"

s := stopwatch.NewWithInterval(time.Millisecond * 100)
```

#### 12. **Help**
Automatically generated keybinding help view.

```go
import "github.com/charmbracelet/bubbles/help"
import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
    Up   key.Binding
    Down key.Binding
    Quit key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
    return []key.Binding{k.Up, k.Down, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding{
        {k.Up, k.Down},
        {k.Quit},
    }
}

h := help.New()
// Render: h.View(keys)
```

#### 13. **Key**
Non-visual keybinding management component.

```go
import "github.com/charmbracelet/bubbles/key"

var keyMap = struct {
    Up   key.Binding
    Down key.Binding
}{
    Up: key.NewBinding(
        key.WithKeys("k", "up"),
        key.WithHelp("↑/k", "move up"),
    ),
    Down: key.NewBinding(
        key.WithKeys("j", "down"),
        key.WithHelp("↓/j", "move down"),
    ),
}

// In Update:
switch {
case key.Matches(msg, keyMap.Up):
    // Handle up
case key.Matches(msg, keyMap.Down):
    // Handle down
}
```

### Third-Party Component Libraries

#### **mistakenelf/teacup**
```go
import "github.com/mistakenelf/teacup"
```

Components:
- **Filetree**: Visual file system navigation
- **Code**: Code viewing and display
- **Help**: Contextual help interface
- **Statusbar**: Application status display
- **Markdown**: Markdown document rendering
- **PDF**: PDF file viewer
- **Image**: Image display and preview

Utilities:
- **dirfs**: Filesystem helper functions
- **icons**: File icon rendering

#### **charm-and-friends/additional-bubbles**

Community-maintained components:
- **treilik/bubblelister**: Scrollable list without pagination
- **erikgeiser/promptkit**: Confirmation prompts
- **mritd/bubbles**: Input validation, menu selection, modified progressbar

#### **charmbracelet/huh** - Form Builder
```go
import "github.com/charmbracelet/huh/v2"
```

Build terminal forms and prompts with:
- **Select**: Single-choice dropdown
- **MultiSelect**: Multiple-choice selections
- **Input**: Text input fields
- **Text**: Multi-line text areas
- **Confirm**: Yes/No confirmation prompts

Features:
- Accessible mode for screen readers
- Validation
- Themes (5 predefined + custom)
- Groups (pages) of fields
- Standalone or Bubble Tea integration

```go
form := huh.NewForm(
    huh.NewGroup(
        huh.NewInput().
            Title("Name").
            Placeholder("Enter your name"),

        huh.NewSelect[string]().
            Title("Favorite Color").
            Options(
                huh.NewOption("Red", "red"),
                huh.NewOption("Blue", "blue"),
            ),

        huh.NewConfirm().
            Title("Are you sure?"),
    ),
)
```

---

## Context-Aware Actions

### Keybinding Management

Use the `key` package for context-aware keybindings:

```go
import "github.com/charmbracelet/bubbles/key"

type contextualKeyMap struct {
    // Navigation (always available)
    Up    key.Binding
    Down  key.Binding
    Left  key.Binding
    Right key.Binding

    // Context-specific
    Enter  key.Binding
    Delete key.Binding
    Edit   key.Binding
}

func (k contextualKeyMap) ShortHelp() []key.Binding {
    return []key.Binding{k.Up, k.Down, k.Enter}
}

func (k contextualKeyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding{
        {k.Up, k.Down, k.Left, k.Right},
        {k.Enter, k.Delete, k.Edit},
    }
}
```

### Dynamic Help System

Show different help based on context:

```go
type Model struct {
    mode     string  // "normal", "edit", "delete"
    keys     contextualKeyMap
    help     help.Model
}

func (m Model) View() string {
    var helpView string

    switch m.mode {
    case "normal":
        helpView = m.help.View(m.normalKeys())
    case "edit":
        helpView = m.help.View(m.editKeys())
    case "delete":
        helpView = m.help.View(m.deleteKeys())
    }

    return content + "\n" + helpView
}
```

### Action Registry Pattern

Define actions that are enabled/disabled based on context:

```go
type Action struct {
    Name        string
    Key         key.Binding
    Handler     func(Model) (Model, tea.Cmd)
    IsAvailable func(Model) bool
}

type ActionRegistry struct {
    actions []Action
}

func (r ActionRegistry) HandleKey(m Model, msg tea.KeyMsg) (Model, tea.Cmd, bool) {
    for _, action := range r.actions {
        if !action.IsAvailable(m) {
            continue
        }

        if key.Matches(msg, action.Key) {
            return action.Handler(m), nil, true
        }
    }
    return m, nil, false
}

// Define actions
var actions = ActionRegistry{
    actions: []Action{
        {
            Name: "Delete Item",
            Key: key.NewBinding(
                key.WithKeys("d"),
                key.WithHelp("d", "delete"),
            ),
            Handler: deleteItem,
            IsAvailable: func(m Model) bool {
                return m.selectedItem != nil && !m.selectedItem.Protected
            },
        },
        {
            Name: "Edit Item",
            Key: key.NewBinding(
                key.WithKeys("e"),
                key.WithHelp("e", "edit"),
            ),
            Handler: editItem,
            IsAvailable: func(m Model) bool {
                return m.selectedItem != nil && m.hasPermission("edit")
            },
        },
    },
}
```

### Mode-Based Navigation

Vim-like modes for different interaction contexts:

```go
const (
    modeNormal = iota
    modeInsert
    modeVisual
    modeCommand
)

type Model struct {
    mode int
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch m.mode {
        case modeNormal:
            return m.handleNormalMode(msg)
        case modeInsert:
            return m.handleInsertMode(msg)
        case modeVisual:
            return m.handleVisualMode(msg)
        case modeCommand:
            return m.handleCommandMode(msg)
        }
    }
    return m, nil
}

func (m Model) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "i":
        m.mode = modeInsert
    case "v":
        m.mode = modeVisual
    case ":":
        m.mode = modeCommand
    }
    return m, nil
}
```

---

## Modal & Overlay Management

### bubbletea-overlay Package

Community solution: `github.com/quickphosphat/bubbletea-overlay`

```go
import "github.com/quickphosphat/bubbletea-overlay"

// Create overlay
overlay := overlay.New(
    backgroundModel,
    foregroundModel,
    overlay.WithPosition(overlay.Center, overlay.Center),
)

// Position options
overlay.WithPosition(overlay.Right, overlay.Top)
overlay.WithPosition(overlay.Center, overlay.Center)
overlay.WithPosition(overlay.Left, overlay.Bottom)
```

**How it works:**
- Wraps two `tea.Model`s (background and foreground)
- Composites foreground onto background in `View()`
- Supports positioning on vertical and horizontal axes

### Manual Modal Implementation

Create modals using Lipgloss composition:

```go
type Model struct {
    showModal    bool
    modalContent string
}

func (m Model) View() string {
    mainView := m.renderMainView()

    if !m.showModal {
        return mainView
    }

    // Create modal
    modal := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("62")).
        Padding(1, 2).
        Width(40).
        Render(m.modalContent)

    // Calculate position to center
    modalWidth := lipgloss.Width(modal)
    modalHeight := lipgloss.Height(modal)
    x := (m.width - modalWidth) / 2
    y := (m.height - modalHeight) / 2

    // Overlay modal on main view
    return lipgloss.PlaceOverlay(x, y, modal, mainView)
}
```

### Confirmation Dialog Pattern

```go
type confirmDialog struct {
    message string
    onYes   func() tea.Cmd
    onNo    func() tea.Cmd
}

func (c confirmDialog) Update(msg tea.Msg) (confirmDialog, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "y", "Y":
            return c, c.onYes()
        case "n", "N", "esc":
            return c, c.onNo()
        }
    }
    return c, nil
}

func (c confirmDialog) View() string {
    s := fmt.Sprintf("%s\n\n[Y]es / [N]o", c.message)
    return dialogStyle.Render(s)
}

// Usage
func (m Model) showDeleteConfirm() (Model, tea.Cmd) {
    m.modal = confirmDialog{
        message: "Delete this item?",
        onYes: func() tea.Cmd {
            return deleteItemCmd(m.selectedItem)
        },
        onNo: func() tea.Cmd {
            return hideModalCmd
        },
    }
    m.showModal = true
    return m, nil
}
```

### Modal Stack for Multiple Overlays

```go
type Model struct {
    modalStack []tea.Model
}

func (m *Model) pushModal(modal tea.Model) {
    m.modalStack = append(m.modalStack, modal)
}

func (m *Model) popModal() {
    if len(m.modalStack) > 0 {
        m.modalStack = m.modalStack[:len(m.modalStack)-1]
    }
}

func (m Model) topModal() tea.Model {
    if len(m.modalStack) == 0 {
        return nil
    }
    return m.modalStack[len(m.modalStack)-1]
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Route to top modal if exists
    if modal := m.topModal(); modal != nil {
        updated, cmd := modal.Update(msg)
        m.modalStack[len(m.modalStack)-1] = updated
        return m, cmd
    }

    // Otherwise handle normally
    return m.handleNormalUpdate(msg)
}

func (m Model) View() string {
    view := m.mainView()

    // Render modal stack
    for _, modal := range m.modalStack {
        view = lipgloss.PlaceOverlay(x, y, modal.View(), view)
    }

    return view
}
```

---

## Keyboard Navigation

### Vim-Like Keybindings

```go
import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
    // Movement
    Up    key.Binding
    Down  key.Binding
    Left  key.Binding
    Right key.Binding

    // Vim-style alternatives
    K key.Binding  // up
    J key.Binding  // down
    H key.Binding  // left
    L key.Binding  // right

    // Jump
    GotoTop    key.Binding
    GotoBottom key.Binding

    // Actions
    Select key.Binding
    Back   key.Binding
    Quit   key.Binding
}

var DefaultKeyMap = keyMap{
    Up: key.NewBinding(
        key.WithKeys("up", "k"),
        key.WithHelp("↑/k", "up"),
    ),
    Down: key.NewBinding(
        key.WithKeys("down", "j"),
        key.WithHelp("↓/j", "down"),
    ),
    Left: key.NewBinding(
        key.WithKeys("left", "h"),
        key.WithHelp("←/h", "left"),
    ),
    Right: key.NewBinding(
        key.WithKeys("right", "l"),
        key.WithHelp("→/l", "right"),
    ),
    GotoTop: key.NewBinding(
        key.WithKeys("g", "home"),
        key.WithHelp("g/home", "go to top"),
    ),
    GotoBottom: key.NewBinding(
        key.WithKeys("G", "end"),
        key.WithHelp("G/end", "go to bottom"),
    ),
    Select: key.NewBinding(
        key.WithKeys("enter"),
        key.WithHelp("enter", "select"),
    ),
    Back: key.NewBinding(
        key.WithKeys("esc"),
        key.WithHelp("esc", "back"),
    ),
    Quit: key.NewBinding(
        key.WithKeys("q", "ctrl+c"),
        key.WithHelp("q", "quit"),
    ),
}
```

### Enhanced Keyboard Input

Bubbletea now provides high-fidelity input events:
- Key up events
- Previously unavailable keybindings (e.g., `shift+enter`)
- Extended key combinations

### Tab Navigation Between Panels

```go
type Model struct {
    panels      []tea.Model
    activePanel int
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab":
            m.activePanel = (m.activePanel + 1) % len(m.panels)
            return m, nil

        case "shift+tab":
            m.activePanel--
            if m.activePanel < 0 {
                m.activePanel = len(m.panels) - 1
            }
            return m, nil
        }
    }

    // Route to active panel
    var cmd tea.Cmd
    m.panels[m.activePanel], cmd = m.panels[m.activePanel].Update(msg)
    return m, cmd
}
```

### Quick Jump / Command Palette

```go
type commandPalette struct {
    input     textinput.Model
    commands  []Command
    filtered  []Command
    selected  int
}

type Command struct {
    Name        string
    Description string
    Keybinding  string
    Execute     func() tea.Cmd
}

func (c commandPalette) filterCommands(query string) []Command {
    if query == "" {
        return c.commands
    }

    var filtered []Command
    for _, cmd := range c.commands {
        if strings.Contains(
            strings.ToLower(cmd.Name),
            strings.ToLower(query),
        ) {
            filtered = append(filtered, cmd)
        }
    }
    return filtered
}

func (c commandPalette) Update(msg tea.Msg) (commandPalette, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter":
            if c.selected < len(c.filtered) {
                return c, c.filtered[c.selected].Execute()
            }
        case "up", "ctrl+k":
            if c.selected > 0 {
                c.selected--
            }
        case "down", "ctrl+j":
            if c.selected < len(c.filtered)-1 {
                c.selected++
            }
        }
    }

    var cmd tea.Cmd
    c.input, cmd = c.input.Update(msg)
    c.filtered = c.filterCommands(c.input.Value())

    return c, cmd
}
```

### Focus Ring Pattern

Manage focus across multiple focusable elements:

```go
type FocusRing struct {
    elements []Focusable
    current  int
}

func (f *FocusRing) Next() tea.Cmd {
    f.elements[f.current].Blur()
    f.current = (f.current + 1) % len(f.elements)
    return f.elements[f.current].Focus()
}

func (f *FocusRing) Prev() tea.Cmd {
    f.elements[f.current].Blur()
    f.current--
    if f.current < 0 {
        f.current = len(f.elements) - 1
    }
    return f.elements[f.current].Focus()
}

func (f FocusRing) Current() Focusable {
    return f.elements[f.current]
}
```

---

## Animation & Visual Effects

### Harmonica - Spring Animation Library

```go
import "github.com/charmbracelet/harmonica"
```

Provides smooth, natural motion for animations:

```go
type Model struct {
    spring   harmonica.Spring
    value    float64
    target   float64
}

func (m Model) Init() tea.Cmd {
    m.spring = harmonica.NewSpring(harmonica.FPS(60), 5.0, 0.5)
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case animateMsg:
        m.value, m.spring = m.spring.Update(m.value, m.target, 1.0/60.0)

        if !harmonica.IsFinished(m.value, m.target, m.spring.Velocity()) {
            return m, animateCmd()
        }
    }
    return m, nil
}
```

### Progress Bar Animation

From official examples:

```go
import "github.com/charmbracelet/bubbles/progress"

type Model struct {
    progress progress.Model
    percent  float64
}

func (m Model) Init() tea.Cmd {
    return tea.Batch(
        tickCmd(),  // Start animation loop
    )
}

func tickCmd() tea.Cmd {
    return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case tickMsg:
        if m.percent < 1.0 {
            m.percent += 0.01
            return m, tickCmd()
        }
    }

    return m, nil
}

func (m Model) View() string {
    return m.progress.ViewAs(m.percent)
}
```

### Spinner for Loading States

```go
import "github.com/charmbracelet/bubbles/spinner"

type Model struct {
    spinner spinner.Model
    loading bool
}

func (m Model) Init() tea.Cmd {
    return m.spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case spinner.TickMsg:
        var cmd tea.Cmd
        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd
    }
    return m, nil
}

func (m Model) View() string {
    if m.loading {
        return fmt.Sprintf("%s Loading...", m.spinner.View())
    }
    return m.content
}
```

### Gradient Effects

```go
func gradientText(text string, colors []lipgloss.Color) string {
    var result strings.Builder
    step := len(colors) - 1

    for i, char := range text {
        colorIndex := (i * step) / len(text)
        style := lipgloss.NewStyle().Foreground(colors[colorIndex])
        result.WriteString(style.Render(string(char)))
    }

    return result.String()
}

// Usage
colors := []lipgloss.Color{
    lipgloss.Color("#FF0000"),
    lipgloss.Color("#FF7F00"),
    lipgloss.Color("#FFFF00"),
    lipgloss.Color("#00FF00"),
    lipgloss.Color("#0000FF"),
}
title := gradientText("Ticktr", colors)
```

### Fade In/Out Transitions

```go
type fadeModel struct {
    content string
    opacity float64
    fadingIn bool
}

func (m fadeModel) Update(msg tea.Msg) (fadeModel, tea.Cmd) {
    switch msg.(type) {
    case fadeTickMsg:
        if m.fadingIn {
            m.opacity += 0.1
            if m.opacity >= 1.0 {
                m.opacity = 1.0
                return m, nil
            }
        } else {
            m.opacity -= 0.1
            if m.opacity <= 0.0 {
                m.opacity = 0.0
                return m, nil
            }
        }
        return m, fadeTickCmd()
    }
    return m, nil
}

func (m fadeModel) View() string {
    // Approximate opacity with color lightness
    color := lipgloss.Color(fmt.Sprintf("#%02X%02X%02X",
        int(255*m.opacity),
        int(255*m.opacity),
        int(255*m.opacity),
    ))

    return lipgloss.NewStyle().Foreground(color).Render(m.content)
}
```

### Alternative Screen Buffer Transitions

```go
// Toggle between normal and alternative screen buffer
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+l":
            return m, tea.EnterAltScreen
        case "ctrl+k":
            return m, tea.ExitAltScreen
        }
    }
    return m, nil
}
```

---

## Testing Patterns

### Official teatest Library

```go
import "github.com/charmbracelet/x/exp/teatest"
```

**Note**: Currently experimental (work-in-progress)

#### 1. Full Output Testing

Capture and compare full program output:

```go
func TestFullOutput(t *testing.T) {
    m := initialModel(time.Second)
    tm := teatest.NewTestModel(
        t, m,
        teatest.WithInitialTermSize(300, 100),
    )

    out, err := io.ReadAll(tm.FinalOutput(t))
    require.NoError(t, err)

    teatest.RequireEqualOutput(t, out)
}
```

**Golden File Testing:**
- Use `-update` flag to generate/update golden files
- Stores expected output for comparison
- Useful for regression testing

#### 2. Final Model State Testing

Verify model's final state:

```go
func TestFinalModel(t *testing.T) {
    tm := teatest.NewTestModel(t, initialModel(time.Second))

    fm := tm.FinalModel(t)
    m, ok := fm.(Model)
    require.True(t, ok)

    assert.True(t, m.completed)
    assert.Equal(t, "expected state", m.state)
}
```

#### 3. Intermediate Interaction Testing

Test program behavior during execution:

```go
func TestOutput(t *testing.T) {
    tm := teatest.NewTestModel(
        t,
        initialModel(10*time.Second),
        teatest.WithInitialTermSize(70, 30),
    )

    // Wait for specific output
    teatest.WaitFor(
        t, tm.Output(),
        func(bts []byte) bool {
            return bytes.Contains(bts, []byte("sleeping 8s"))
        },
        teatest.WithCheckInterval(time.Millisecond*100),
        teatest.WithDuration(time.Second*3),
    )

    // Send user input
    tm.Send(tea.KeyMsg{
        Type:  tea.KeyRunes,
        Runes: []rune("q"),
    })

    // Wait for quit
    tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
```

### Catwalk - Data-Driven Unit Testing

```go
import "github.com/knz/catwalk"
```

**Features:**
- Data-driven testing approach
- Verify model state and View as they process messages
- Built on datadriven testing framework

**Test File Format:**
```
# Test case: navigation
type "j"
----
selectedIndex: 1

type "k"
----
selectedIndex: 0

enter ""
----
itemSelected: true
```

**Usage:**
```go
func TestModel(t *testing.T) {
    datadriven.RunTest(t, "testdata/navigation", func(t *testing.T, d *datadriven.TestData) string {
        m := initialModel()

        switch d.Cmd {
        case "type":
            msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(d.Input)}
            m, _ = m.Update(msg)
        case "enter":
            msg := tea.KeyMsg{Type: tea.KeyEnter}
            m, _ = m.Update(msg)
        }

        return fmt.Sprintf("selectedIndex: %d\nitemSelected: %v",
            m.selectedIndex, m.itemSelected)
    })
}
```

### Best Practices

1. **Use Fixed Terminal Sizes**
   ```go
   teatest.WithInitialTermSize(width, height)
   ```

2. **Handle Color Profiles**
   - Tests may fail if color rendering differs
   - Consider using `lipgloss.SetColorProfile()`

3. **Manage Line Endings**
   - Golden files should match platform line endings
   - Consider normalizing in tests

4. **Timeouts and Intervals**
   ```go
   teatest.WithCheckInterval(time.Millisecond*100)
   teatest.WithDuration(time.Second*5)
   teatest.WithFinalTimeout(time.Second*10)
   ```

5. **Important Limitation**
   - Bubbletea can be brittle with fast automated input
   - Integration with unit testing requires care
   - Designed for relatively slow human input

### Unit Testing Without teatest

Test individual methods directly:

```go
func TestModelUpdate(t *testing.T) {
    m := Model{selectedIndex: 0, items: []string{"a", "b", "c"}}

    msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
    updated, _ := m.Update(msg)

    assert.Equal(t, 1, updated.(Model).selectedIndex)
}

func TestModelView(t *testing.T) {
    m := Model{title: "Test", content: "Content"}

    view := m.View()

    assert.Contains(t, view, "Test")
    assert.Contains(t, view, "Content")
}
```

---

## Production Applications

### Charm Official Applications

**Core Tools:**
- **Glow**: Markdown reader for the CLI with pizzazz
- **Soft Serve**: Self-hostable Git server with TUI
- **VHS**: Terminal recorder for creating demos
- **Mods**: AI on the CLI, built for pipelines
- **Gum**: Tool for glamorous shell scripts

### Notable Production Applications

Source: [charm-in-the-wild](https://github.com/charm-and-friends/charm-in-the-wild)

#### AI & Coding Tools

- **Plandex**: Terminal-based AI coding engine for complex tasks
- **ChatGPT-TUI**: ChatGPT client in terminal
- **Gemini-CLI**: Gemini API terminal interactions
- **Aichat**: All-in-one chat and copilot CLI

#### DevOps & Cloud

- **Atmos**: Terraform orchestration (uses multiple Charm libraries)
- **eks-node-viewer**: Kubernetes cluster visualization
- **container-canary**: Container validation tool
- **k9s**: Kubernetes cluster management
- **Kubetui**: Kubernetes cluster explorer

#### Development Tools

- **gh-dash**: GitHub CLI dashboard extension
- **mergestat**: SQL queries for Git repositories
- **Flow**: Customizable interactive task runner
- **chezmoi**: Dotfile manager
- **kubectl-ai**: AI-powered kubectl

#### File Management

- **walk**: Terminal file manager
- **fm**: Terminal file manager
- **superfile**: Pretty fancy and modern file manager

#### Games

- **Gambit**: Terminal chess
- **Typer**: Typing test in terminal
- **Dungeonfs**: FUSE filesystem dungeon crawler
- **CHIP-8**: CHIP-8 emulator

#### Productivity

- **Canard**: Task management
- **Timekeeper**: Time tracking
- **slides**: Terminal-based presentation tool
- **tasktimer**: Task timer

#### Unique Applications

- **termdbms**: Database management system UI
- **ov**: Feature-rich terminal pager
- **WG Commander**: Wireguard TUI
- **radio**: Terminal radio player
- **ascii-movie-player**: ASCII video player

### Architectural Insights from Production Apps

**Common Patterns:**
1. **Modular CLI/TUI Design**: Support both command-line and TUI modes
2. **Configuration Management**: YAML files + CLI flags + runtime config
3. **Separate Concerns**: `/ui` or `/tui` directory for UI components
4. **Integration with Cobra**: For multi-command CLIs
5. **Styling Libraries**: Heavy use of Glamour for markdown, Lipgloss for styling

**File Organization (from Glow):**
```
glow/
├── ui/              # UI-specific components
├── utils/           # Utility functions
├── github.go        # Service integrations
├── gitlab.go
└── config.yml       # Configuration
```

---

## Recommendations

### Architecture Recommendations

#### 1. **Project Structure**

```
ticktr/
├── cmd/
│   ├── root.go              # Cobra root command
│   └── tui.go               # TUI command
├── internal/
│   ├── domain/              # Business logic
│   │   ├── jira/
│   │   └── models/
│   ├── storage/             # Persistence layer
│   └── tui/                 # TUI layer
│       ├── app.go           # Root model
│       ├── styles/          # Lipgloss styles
│       │   ├── colors.go
│       │   └── components.go
│       ├── components/      # Reusable components
│       │   ├── statusbar/
│       │   ├── header/
│       │   └── keybindings/
│       ├── panels/          # Main panels
│       │   ├── left/
│       │   └── right/
│       ├── modals/          # Modal dialogs
│       │   ├── confirm.go
│       │   └── input.go
│       └── messages/        # Custom message types
│           ├── data.go
│           └── ui.go
├── config.yml
└── main.go
```

#### 2. **State Management Strategy**

- **Root Model**: Manages global state, window size, active panel
- **Panel Models**: Independent state for each panel
- **Component Models**: Self-contained components with their own state
- **Message Routing**: Root routes messages to appropriate handlers
- **State Machine**: For complex multi-step workflows

#### 3. **Rendering Strategy**

- **60 FPS Default**: Trust Bubbletea's renderer
- **Dynamic Layouts**: Use Lipgloss `Width()` and `Height()` methods
- **Layout Manager**: Consider `bubblelayout` or `stickers` for complex layouts
- **Responsive Design**: Always handle `WindowSizeMsg`

#### 4. **Component Strategy**

**Official Bubbles:**
- **List**: For file/issue browsing
- **Viewport**: For scrollable content
- **Table**: For tabular data
- **TextInput**: For search/filter
- **Help**: For context-aware keybindings
- **Spinner**: For loading states

**Third-Party:**
- **teacup/filetree**: For file system navigation
- **teacup/code**: For code viewing
- **huh**: For forms and prompts
- **bubbletea-overlay**: For modals

**Custom Components:**
- Dual panel layout
- Action registry
- Command palette
- Status bar with context

#### 5. **Performance Strategy**

- **Keep Update Fast**: Offload I/O to Commands
- **Use tea.Batch**: For multiple concurrent operations
- **Use tea.Sequence**: For ordered operations
- **Viewport Optimization**: Paginate large datasets
- **Avoid Goroutines**: Use Commands exclusively

#### 6. **Keyboard Navigation Strategy**

- **Vim-Like Bindings**: Support both arrow keys and hjkl
- **Tab Navigation**: Between panels
- **Context-Aware Actions**: Different keys in different modes
- **Command Palette**: Quick access to all actions (Ctrl+P)
- **Help System**: Context-sensitive help (?)

#### 7. **Visual Effects Strategy**

- **Spinner**: For async operations
- **Progress Bars**: For long-running tasks
- **Harmonica**: For smooth animations
- **Gradient Effects**: For visual polish
- **Adaptive Colors**: For light/dark terminal support

#### 8. **Testing Strategy**

- **teatest**: For integration testing
- **Golden Files**: For regression testing
- **Unit Tests**: For individual Update/View logic
- **Fixed Terminal Size**: For consistent test output
- **Data-Driven**: Consider catwalk for complex scenarios

### Implementation Roadmap

#### Phase 1: Foundation
1. Set up project structure
2. Implement root model with dual panel layout
3. Integrate Lipgloss styling system
4. Handle WindowSizeMsg for responsive design
5. Basic keyboard navigation (arrows, tab, quit)

#### Phase 2: Core Components
1. Implement left panel (file/issue list)
2. Implement right panel (detail view)
3. Add status bar with context
4. Add header with title and mode
5. Integrate Bubbles list component

#### Phase 3: Interactions
1. Implement focus management between panels
2. Add vim-like navigation keybindings
3. Create action registry system
4. Add context-aware help system
5. Implement Commands for async I/O

#### Phase 4: Advanced Features
1. Modal/overlay system for confirmations
2. Command palette for quick actions
3. Forms using huh for input
4. Loading spinners for async operations
5. Progress bars for long tasks

#### Phase 5: Polish
1. Animations and transitions
2. Visual effects and gradients
3. Themes and color schemes
4. Performance optimization
5. Comprehensive testing

### Key Takeaways

1. **Embrace The Elm Architecture**: Model-View-Update pattern is powerful
2. **Use Commands for I/O**: Never use goroutines directly
3. **Keep Update Fast**: Offload expensive operations
4. **Trust the Renderer**: 60 FPS is smooth
5. **Component Composition**: Build reusable, self-contained components
6. **Context-Aware UX**: Actions and help adapt to current state
7. **Test Thoroughly**: Use teatest and golden files
8. **Study Production Apps**: Glow, Soft Serve, and charm-in-the-wild apps

### Additional Resources

**Official Documentation:**
- [Bubbletea GitHub](https://github.com/charmbracelet/bubbletea)
- [Bubbles GitHub](https://github.com/charmbracelet/bubbles)
- [Lipgloss GitHub](https://github.com/charmbracelet/lipgloss)
- [Huh GitHub](https://github.com/charmbracelet/huh)

**Tutorials & Guides:**
- [Building Bubbletea Programs](https://leg100.github.io/en/posts/building-bubbletea-programs/)
- [Commands in Bubbletea](https://charm.land/blog/commands-in-bubbletea/)
- [State Machine Pattern](https://zackproser.com/blog/bubbletea-state-machine)
- [Multi-View Interfaces](https://shi.foo/weblog/multi-view-interfaces-in-bubble-tea)
- [Writing Bubble Tea Tests](https://carlosbecker.com/posts/teatest/)

**Example Applications:**
- [Bubbletea Examples](https://github.com/charmbracelet/bubbletea/tree/main/examples)
- [Charm in the Wild](https://github.com/charm-and-friends/charm-in-the-wild)
- [Additional Bubbles](https://github.com/charm-and-friends/additional-bubbles)

**Component Libraries:**
- [teacup](https://github.com/mistakenelf/teacup)
- [stickers](https://github.com/76creates/stickers)
- [bubblelayout](https://github.com/winder/bubblelayout)
- [bubbletea-overlay](https://github.com/quickphosphat/bubbletea-overlay)

---

## Conclusion

Bubbletea provides a robust, production-ready framework for building world-class TUIs. Its Elm-inspired architecture, combined with the Charm ecosystem (Lipgloss, Bubbles, Huh), offers all the tools needed to build a Midnight Commander-inspired dual-panel TUI with high performance, extensible actions, and visual polish.

The key to success is:
1. Embracing the Model-View-Update pattern
2. Using Commands for all I/O
3. Building composable, reusable components
4. Leveraging the rich ecosystem of libraries
5. Studying production applications for patterns
6. Testing thoroughly with teatest

With this research as a foundation, the migration from Tview to Bubbletea should result in a more maintainable, performant, and polished TUI application.
