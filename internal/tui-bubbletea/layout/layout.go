// Package layout provides layout management using Lipgloss
// Pure Lipgloss provides excellent layout control for our dual-panel design.
package layout

import (
	"github.com/charmbracelet/lipgloss"
)

// DualPanelLayout manages a dual-panel layout
// Left panel: 40%, Right panel: 60% (configurable)
type DualPanelLayout struct {
	width  int
	height int

	// Split ratio (0.0 - 1.0)
	leftRatio float64
}

// NewDualPanelLayout creates a new dual-panel layout
// leftRatio: percentage of width for left panel (e.g., 0.4 for 40%)
func NewDualPanelLayout(width, height int, leftRatio float64) *DualPanelLayout {
	return &DualPanelLayout{
		width:     width,
		height:    height,
		leftRatio: leftRatio,
	}
}

// Resize updates the layout dimensions
func (l *DualPanelLayout) Resize(width, height int) {
	l.width = width
	l.height = height
}

// Render renders both panels side-by-side using Lipgloss
func (l *DualPanelLayout) Render(leftContent, rightContent string) string {
	// Render panels side-by-side
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftContent,
		rightContent,
	)
}

// SetLeftRatio updates the split ratio between panels
func (l *DualPanelLayout) SetLeftRatio(ratio float64) {
	if ratio < 0.1 {
		ratio = 0.1
	}
	if ratio > 0.9 {
		ratio = 0.9
	}
	l.leftRatio = ratio
}

// GetDimensions returns the current width and height
func (l *DualPanelLayout) GetDimensions() (width, height int) {
	return l.width, l.height
}

// GetPanelWidths returns the widths of left and right panels
func (l *DualPanelLayout) GetPanelWidths() (left, right int) {
	left = int(float64(l.width) * l.leftRatio)
	right = l.width - left
	return
}

// TriSectionLayout manages a 3-section vertical layout (Header, Content, Footer)
type TriSectionLayout struct {
	width        int
	height       int
	headerHeight int
	footerHeight int
}

// NewTriSectionLayout creates a vertical 3-section layout
// headerHeight and footerHeight are fixed, content is flexible
func NewTriSectionLayout(width, height, headerHeight, footerHeight int) *TriSectionLayout {
	return &TriSectionLayout{
		width:        width,
		height:       height,
		headerHeight: headerHeight,
		footerHeight: footerHeight,
	}
}

// Resize updates the layout dimensions
func (l *TriSectionLayout) Resize(width, height int) {
	l.width = width
	l.height = height
}

// Render renders all sections vertically using Lipgloss
func (l *TriSectionLayout) Render(header, content, footer string) string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		footer,
	)
}

// GetContentHeight returns the available content height
func (l *TriSectionLayout) GetContentHeight() int {
	h := l.height - l.headerHeight - l.footerHeight
	if h < 0 {
		return 0
	}
	return h
}

// CompleteLayout combines TriSection and DualPanel for the full app layout
type CompleteLayout struct {
	triSection *TriSectionLayout
	dualPanel  *DualPanelLayout
	width      int
	height     int
}

// NewCompleteLayout creates the complete app layout
// Header (3 rows), Content (dual panel), Action bar (3 rows)
func NewCompleteLayout(width, height int) *CompleteLayout {
	const (
		headerHeight = 3
		footerHeight = 3
	)

	triSection := NewTriSectionLayout(width, height, headerHeight, footerHeight)
	contentHeight := triSection.GetContentHeight()

	// Create dual panel for content area with 40/60 split
	dualPanel := NewDualPanelLayout(width, contentHeight, 0.4)

	return &CompleteLayout{
		triSection: triSection,
		dualPanel:  dualPanel,
		width:      width,
		height:     height,
	}
}

// Resize handles terminal resize
func (l *CompleteLayout) Resize(width, height int) {
	l.width = width
	l.height = height
	l.triSection.Resize(width, height)

	contentHeight := l.triSection.GetContentHeight()
	l.dualPanel.Resize(width, contentHeight)
}

// Render renders the complete layout
func (l *CompleteLayout) Render(header, leftPanel, rightPanel, footer string) string {
	// Render dual panel content
	content := l.dualPanel.Render(leftPanel, rightPanel)

	// Render complete tri-section layout
	return l.triSection.Render(header, content, footer)
}

// GetPanelDimensions returns the dimensions for left and right panels
func (l *CompleteLayout) GetPanelDimensions() (leftWidth, rightWidth, height int) {
	leftWidth, rightWidth = l.dualPanel.GetPanelWidths()
	height = l.triSection.GetContentHeight()
	return
}
