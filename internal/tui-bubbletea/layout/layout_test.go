package layout

import (
	"strings"
	"testing"
)

// TestNewDualPanelLayout tests creating a dual panel layout
func TestNewDualPanelLayout(t *testing.T) {
	width, height := 100, 40
	leftRatio := 0.4

	layout := NewDualPanelLayout(width, height, leftRatio)

	if layout == nil {
		t.Fatal("Expected layout to be created")
	}

	if layout.width != width {
		t.Errorf("Expected width %d, got %d", width, layout.width)
	}

	if layout.height != height {
		t.Errorf("Expected height %d, got %d", height, layout.height)
	}

	if layout.leftRatio != leftRatio {
		t.Errorf("Expected leftRatio %.2f, got %.2f", leftRatio, layout.leftRatio)
	}
}

// TestDualPanelLayoutResize tests resizing a dual panel layout
func TestDualPanelLayoutResize(t *testing.T) {
	layout := NewDualPanelLayout(100, 40, 0.4)

	newWidth, newHeight := 120, 50
	layout.Resize(newWidth, newHeight)

	if layout.width != newWidth {
		t.Errorf("Expected width %d after resize, got %d", newWidth, layout.width)
	}

	if layout.height != newHeight {
		t.Errorf("Expected height %d after resize, got %d", newHeight, layout.height)
	}
}

// TestDualPanelLayoutGetPanelWidths tests getting panel widths
func TestDualPanelLayoutGetPanelWidths(t *testing.T) {
	tests := []struct {
		name       string
		totalWidth int
		leftRatio  float64
		expectLeft int
	}{
		{"40% split", 100, 0.4, 40},
		{"50% split", 100, 0.5, 50},
		{"30% split", 120, 0.3, 36},
		{"60% split", 80, 0.6, 48},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout := NewDualPanelLayout(tt.totalWidth, 40, tt.leftRatio)
			left, right := layout.GetPanelWidths()

			if left != tt.expectLeft {
				t.Errorf("Expected left width %d, got %d", tt.expectLeft, left)
			}

			expectedRight := tt.totalWidth - tt.expectLeft
			if right != expectedRight {
				t.Errorf("Expected right width %d, got %d", expectedRight, right)
			}

			// Verify widths sum to total
			if left+right != tt.totalWidth {
				t.Errorf("Panel widths %d + %d = %d, expected %d",
					left, right, left+right, tt.totalWidth)
			}
		})
	}
}

// TestDualPanelLayoutSetLeftRatio tests setting left ratio
func TestDualPanelLayoutSetLeftRatio(t *testing.T) {
	layout := NewDualPanelLayout(100, 40, 0.4)

	tests := []struct {
		name     string
		ratio    float64
		expected float64
	}{
		{"valid 0.5", 0.5, 0.5},
		{"valid 0.3", 0.3, 0.3},
		{"valid 0.7", 0.7, 0.7},
		{"too low 0.05", 0.05, 0.1},  // Clamped to 0.1
		{"too high 0.95", 0.95, 0.9}, // Clamped to 0.9
		{"negative", -0.5, 0.1},      // Clamped to 0.1
		{"over 1.0", 1.5, 0.9},       // Clamped to 0.9
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout.SetLeftRatio(tt.ratio)
			if layout.leftRatio != tt.expected {
				t.Errorf("Expected ratio %.2f, got %.2f", tt.expected, layout.leftRatio)
			}
		})
	}
}

// TestDualPanelLayoutGetDimensions tests getting dimensions
func TestDualPanelLayoutGetDimensions(t *testing.T) {
	width, height := 120, 50
	layout := NewDualPanelLayout(width, height, 0.4)

	w, h := layout.GetDimensions()

	if w != width {
		t.Errorf("Expected width %d, got %d", width, w)
	}

	if h != height {
		t.Errorf("Expected height %d, got %d", height, h)
	}
}

// TestDualPanelLayoutRender tests rendering two panels
func TestDualPanelLayoutRender(t *testing.T) {
	layout := NewDualPanelLayout(100, 40, 0.4)

	leftContent := "LEFT PANEL"
	rightContent := "RIGHT PANEL"

	rendered := layout.Render(leftContent, rightContent)

	// Should contain both panels
	if !strings.Contains(rendered, leftContent) {
		t.Error("Expected rendered output to contain left panel content")
	}

	if !strings.Contains(rendered, rightContent) {
		t.Error("Expected rendered output to contain right panel content")
	}
}

// TestNewTriSectionLayout tests creating a tri-section layout
func TestNewTriSectionLayout(t *testing.T) {
	width, height := 100, 40
	headerHeight, footerHeight := 3, 3

	layout := NewTriSectionLayout(width, height, headerHeight, footerHeight)

	if layout == nil {
		t.Fatal("Expected layout to be created")
	}

	if layout.width != width {
		t.Errorf("Expected width %d, got %d", width, layout.width)
	}

	if layout.height != height {
		t.Errorf("Expected height %d, got %d", height, layout.height)
	}

	if layout.headerHeight != headerHeight {
		t.Errorf("Expected headerHeight %d, got %d", headerHeight, layout.headerHeight)
	}

	if layout.footerHeight != footerHeight {
		t.Errorf("Expected footerHeight %d, got %d", footerHeight, layout.footerHeight)
	}
}

// TestTriSectionLayoutResize tests resizing a tri-section layout
func TestTriSectionLayoutResize(t *testing.T) {
	layout := NewTriSectionLayout(100, 40, 3, 3)

	newWidth, newHeight := 120, 50
	layout.Resize(newWidth, newHeight)

	if layout.width != newWidth {
		t.Errorf("Expected width %d after resize, got %d", newWidth, layout.width)
	}

	if layout.height != newHeight {
		t.Errorf("Expected height %d after resize, got %d", newHeight, layout.height)
	}
}

// TestTriSectionLayoutGetContentHeight tests getting content height
func TestTriSectionLayoutGetContentHeight(t *testing.T) {
	tests := []struct {
		name          string
		totalHeight   int
		headerHeight  int
		footerHeight  int
		expectedContent int
	}{
		{"normal", 40, 3, 3, 34},
		{"large", 100, 5, 5, 90},
		{"small", 20, 2, 2, 16},
		{"exact fit", 10, 3, 3, 4},
		{"too small", 5, 3, 3, 0}, // Content height clamped to 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout := NewTriSectionLayout(100, tt.totalHeight, tt.headerHeight, tt.footerHeight)
			contentHeight := layout.GetContentHeight()

			if contentHeight != tt.expectedContent {
				t.Errorf("Expected content height %d, got %d", tt.expectedContent, contentHeight)
			}

			// Verify it's never negative
			if contentHeight < 0 {
				t.Error("Content height should never be negative")
			}
		})
	}
}

// TestTriSectionLayoutRender tests rendering three sections
func TestTriSectionLayoutRender(t *testing.T) {
	layout := NewTriSectionLayout(100, 40, 3, 3)

	header := "HEADER"
	content := "CONTENT"
	footer := "FOOTER"

	rendered := layout.Render(header, content, footer)

	// Should contain all sections
	if !strings.Contains(rendered, header) {
		t.Error("Expected rendered output to contain header")
	}

	if !strings.Contains(rendered, content) {
		t.Error("Expected rendered output to contain content")
	}

	if !strings.Contains(rendered, footer) {
		t.Error("Expected rendered output to contain footer")
	}
}

// TestNewCompleteLayout tests creating complete layout
func TestNewCompleteLayout(t *testing.T) {
	width, height := 120, 40

	layout := NewCompleteLayout(width, height)

	if layout == nil {
		t.Fatal("Expected layout to be created")
	}

	if layout.width != width {
		t.Errorf("Expected width %d, got %d", width, layout.width)
	}

	if layout.height != height {
		t.Errorf("Expected height %d, got %d", height, layout.height)
	}

	if layout.triSection == nil {
		t.Error("Expected triSection to be initialized")
	}

	if layout.dualPanel == nil {
		t.Error("Expected dualPanel to be initialized")
	}
}

// TestCompleteLayoutResize tests resizing complete layout
func TestCompleteLayoutResize(t *testing.T) {
	layout := NewCompleteLayout(100, 40)

	newWidth, newHeight := 120, 50
	layout.Resize(newWidth, newHeight)

	if layout.width != newWidth {
		t.Errorf("Expected width %d after resize, got %d", newWidth, layout.width)
	}

	if layout.height != newHeight {
		t.Errorf("Expected height %d after resize, got %d", newHeight, layout.height)
	}
}

// TestCompleteLayoutGetPanelDimensions tests getting panel dimensions
func TestCompleteLayoutGetPanelDimensions(t *testing.T) {
	tests := []struct {
		name         string
		width        int
		height       int
		expectedLeft int
	}{
		{"standard 120x40", 120, 40, 48},  // 40% of 120
		{"small 80x24", 80, 24, 32},       // 40% of 80
		{"large 200x60", 200, 60, 80},     // 40% of 200
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout := NewCompleteLayout(tt.width, tt.height)
			leftWidth, rightWidth, contentHeight := layout.GetPanelDimensions()

			if leftWidth != tt.expectedLeft {
				t.Errorf("Expected left width %d, got %d", tt.expectedLeft, leftWidth)
			}

			expectedRight := tt.width - tt.expectedLeft
			if rightWidth != expectedRight {
				t.Errorf("Expected right width %d, got %d", expectedRight, rightWidth)
			}

			// Verify panels sum to total width
			if leftWidth+rightWidth != tt.width {
				t.Errorf("Panel widths %d + %d = %d, expected %d",
					leftWidth, rightWidth, leftWidth+rightWidth, tt.width)
			}

			// Verify content height accounts for header and footer (3 + 3 = 6)
			expectedContentHeight := tt.height - 6
			if expectedContentHeight < 0 {
				expectedContentHeight = 0
			}
			if contentHeight != expectedContentHeight {
				t.Errorf("Expected content height %d, got %d", expectedContentHeight, contentHeight)
			}
		})
	}
}

// TestCompleteLayoutRender tests rendering complete layout
func TestCompleteLayoutRender(t *testing.T) {
	layout := NewCompleteLayout(120, 40)

	header := "HEADER"
	leftPanel := "LEFT"
	rightPanel := "RIGHT"
	footer := "FOOTER"

	rendered := layout.Render(header, leftPanel, rightPanel, footer)

	// Should contain all sections
	if !strings.Contains(rendered, header) {
		t.Error("Expected rendered output to contain header")
	}

	if !strings.Contains(rendered, leftPanel) {
		t.Error("Expected rendered output to contain left panel")
	}

	if !strings.Contains(rendered, rightPanel) {
		t.Error("Expected rendered output to contain right panel")
	}

	if !strings.Contains(rendered, footer) {
		t.Error("Expected rendered output to contain footer")
	}
}

// TestCompleteLayoutEdgeCases tests edge cases
func TestCompleteLayoutEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"minimum size", 80, 24},
		{"very small", 10, 10},
		{"zero width", 0, 40},
		{"zero height", 100, 0},
		{"large", 1000, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			layout := NewCompleteLayout(tt.width, tt.height)
			if layout == nil {
				t.Error("Expected layout to be created even for edge cases")
			}

			// Should be able to get dimensions
			leftWidth, rightWidth, contentHeight := layout.GetPanelDimensions()

			// Verify no negative dimensions
			if leftWidth < 0 || rightWidth < 0 || contentHeight < 0 {
				t.Error("Expected no negative dimensions")
			}

			// Should be able to render without panicking
			rendered := layout.Render("H", "L", "R", "F")
			if rendered == "" {
				t.Log("Rendered empty string for edge case (acceptable)")
			}
		})
	}
}

// TestLayoutResizePreservesRatios tests that resize maintains aspect ratios
func TestLayoutResizePreservesRatios(t *testing.T) {
	layout := NewCompleteLayout(100, 40)

	// Get initial ratio
	left1, right1, _ := layout.GetPanelDimensions()
	initialRatio := float64(left1) / float64(left1+right1)

	// Resize
	layout.Resize(120, 50)

	// Get new ratio
	left2, right2, _ := layout.GetPanelDimensions()
	newRatio := float64(left2) / float64(left2+right2)

	// Ratios should be approximately equal
	if abs(initialRatio-newRatio) > 0.01 {
		t.Errorf("Expected ratio to be preserved (%.3f vs %.3f)", initialRatio, newRatio)
	}
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// TestDualPanelLayoutEmptyContent tests rendering with empty content
func TestDualPanelLayoutEmptyContent(t *testing.T) {
	layout := NewDualPanelLayout(100, 40, 0.4)

	// Render with empty strings
	rendered := layout.Render("", "")

	// Should not panic, even if empty
	_ = rendered
}

// TestTriSectionLayoutEmptyContent tests rendering with empty content
func TestTriSectionLayoutEmptyContent(t *testing.T) {
	layout := NewTriSectionLayout(100, 40, 3, 3)

	// Render with empty strings
	rendered := layout.Render("", "", "")

	// Should not panic, even if empty
	_ = rendered
}

// TestCompleteLayoutEmptyContent tests rendering with empty content
func TestCompleteLayoutEmptyContent(t *testing.T) {
	layout := NewCompleteLayout(100, 40)

	// Render with empty strings
	rendered := layout.Render("", "", "", "")

	// Should not panic, even if empty
	_ = rendered
}

// TestLayoutDimensionConsistency tests that dimensions are consistent across operations
func TestLayoutDimensionConsistency(t *testing.T) {
	width, height := 120, 40
	layout := NewCompleteLayout(width, height)

	// Verify dimensions match what we set
	if layout.width != width || layout.height != height {
		t.Errorf("Expected dimensions %dx%d, got %dx%d", width, height, layout.width, layout.height)
	}

	// Resize and verify
	newWidth, newHeight := 150, 50
	layout.Resize(newWidth, newHeight)

	if layout.width != newWidth || layout.height != newHeight {
		t.Errorf("Expected dimensions %dx%d after resize, got %dx%d",
			newWidth, newHeight, layout.width, layout.height)
	}
}

// TestPanelWidthDistribution tests various width distributions
func TestPanelWidthDistribution(t *testing.T) {
	ratios := []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9}

	for _, ratio := range ratios {
		layout := NewDualPanelLayout(100, 40, ratio)
		left, right := layout.GetPanelWidths()

		// Verify sum equals total
		if left+right != 100 {
			t.Errorf("For ratio %.1f, widths don't sum to 100: %d + %d = %d",
				ratio, left, right, left+right)
		}

		// Verify ratio is approximately correct
		actualRatio := float64(left) / 100.0
		if abs(actualRatio-ratio) > 0.01 {
			t.Errorf("For ratio %.2f, got %.2f", ratio, actualRatio)
		}
	}
}
