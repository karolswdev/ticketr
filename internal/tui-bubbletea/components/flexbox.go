// Package components provides reusable Bubbletea components for the Ticketr TUI.
package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// FlexDirection represents the layout direction
type FlexDirection int

const (
	FlexRow FlexDirection = iota
	FlexColumn
)

// FlexBox is a flexible layout container (like CSS flexbox)
type FlexBox struct {
	Direction FlexDirection
	Children  []string
	Width     int
	Height    int
	Gap       int // Space between items
}

// NewFlexBox creates a new FlexBox container
func NewFlexBox(direction FlexDirection) *FlexBox {
	return &FlexBox{
		Direction: direction,
		Children:  []string{},
		Gap:       0,
	}
}

// AddChild adds a child element to the flexbox
func (f *FlexBox) AddChild(content string) {
	f.Children = append(f.Children, content)
}

// SetSize sets the width and height of the flexbox
func (f *FlexBox) SetSize(width, height int) {
	f.Width = width
	f.Height = height
}

// SetGap sets the gap between children
func (f *FlexBox) SetGap(gap int) {
	f.Gap = gap
}

// Render renders the flexbox layout
func (f *FlexBox) Render() string {
	if len(f.Children) == 0 {
		return ""
	}

	if f.Direction == FlexRow {
		return f.renderRow()
	}
	return f.renderColumn()
}

func (f *FlexBox) renderRow() string {
	// Calculate available width per child
	totalGap := f.Gap * (len(f.Children) - 1)
	availableWidth := f.Width - totalGap
	childWidth := availableWidth / len(f.Children)

	// Render each child with fixed width
	rendered := make([]string, len(f.Children))
	for i, child := range f.Children {
		style := lipgloss.NewStyle().Width(childWidth).Height(f.Height)
		rendered[i] = style.Render(child)
	}

	// Join horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
}

func (f *FlexBox) renderColumn() string {
	// Calculate available height per child
	totalGap := f.Gap * (len(f.Children) - 1)
	availableHeight := f.Height - totalGap
	childHeight := availableHeight / len(f.Children)

	// Render each child with fixed height
	rendered := make([]string, len(f.Children))
	for i, child := range f.Children {
		style := lipgloss.NewStyle().Width(f.Width).Height(childHeight)
		rendered[i] = style.Render(child)
	}

	// Join vertically with gap
	if f.Gap > 0 {
		gap := strings.Repeat("\n", f.Gap)
		return strings.Join(rendered, gap)
	}
	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}
