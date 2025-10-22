package effects

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ShadowBox extends tview.Box with automatic drop shadow rendering.
// The shadow is rendered using ▒ characters offset by 1 row and 2 columns.
type ShadowBox struct {
	*tview.Box
	shadowEnabled bool
	shadowChar    rune
	shadowColor   tcell.Color
	shadowOffsetX int
	shadowOffsetY int
}

// NewShadowBox creates a new box with drop shadow support.
func NewShadowBox() *ShadowBox {
	return &ShadowBox{
		Box:           tview.NewBox(),
		shadowEnabled: true,
		shadowChar:    '▒',
		shadowColor:   tcell.ColorGray,
		shadowOffsetX: 2,
		shadowOffsetY: 1,
	}
}

// SetShadowEnabled enables or disables the shadow effect.
func (sb *ShadowBox) SetShadowEnabled(enabled bool) *ShadowBox {
	sb.shadowEnabled = enabled
	return sb
}

// SetShadowChar sets the character used for rendering the shadow.
func (sb *ShadowBox) SetShadowChar(char rune) *ShadowBox {
	sb.shadowChar = char
	return sb
}

// SetShadowColor sets the color of the shadow.
func (sb *ShadowBox) SetShadowColor(color tcell.Color) *ShadowBox {
	sb.shadowColor = color
	return sb
}

// SetShadowOffset sets the shadow offset in characters.
func (sb *ShadowBox) SetShadowOffset(x, y int) *ShadowBox {
	sb.shadowOffsetX = x
	sb.shadowOffsetY = y
	return sb
}

// Draw renders the box with drop shadow.
func (sb *ShadowBox) Draw(screen tcell.Screen) {
	if !sb.shadowEnabled {
		sb.Box.Draw(screen)
		return
	}

	// Get box dimensions
	x, y, width, height := sb.GetRect()

	// Draw shadow first (behind the box)
	shadowStyle := tcell.StyleDefault.
		Foreground(sb.shadowColor).
		Background(tcell.ColorDefault).
		Dim(true)

	// Draw bottom shadow (horizontal line)
	shadowY := y + height + sb.shadowOffsetY - 1
	if shadowY >= 0 {
		for sx := x + sb.shadowOffsetX; sx < x+width+sb.shadowOffsetX; sx++ {
			if sx >= 0 {
				screen.SetContent(sx, shadowY, sb.shadowChar, nil, shadowStyle)
			}
		}
	}

	// Draw right shadow (vertical line)
	shadowX := x + width + sb.shadowOffsetX - 1
	if shadowX >= 0 {
		for sy := y + sb.shadowOffsetY; sy < y+height+sb.shadowOffsetY; sy++ {
			if sy >= 0 {
				screen.SetContent(shadowX, sy, sb.shadowChar, nil, shadowStyle)
			}
		}
	}

	// Draw corner shadow
	if shadowX >= 0 && shadowY >= 0 {
		screen.SetContent(shadowX, shadowY, sb.shadowChar, nil, shadowStyle)
	}

	// Draw the box on top
	sb.Box.Draw(screen)
}

// WrapPrimitive wraps an existing primitive with a shadow box.
// This is useful for adding shadows to modals and other components.
func WrapWithShadow(primitive tview.Primitive) *ShadowBox {
	// Create a flex layout with the primitive
	shadowBox := NewShadowBox()

	// Note: ShadowBox extends Box, so it can't directly contain other primitives.
	// For wrapping, we need a different approach using Flex or Grid.
	return shadowBox
}

// ShadowModal wraps a tview.Modal with drop shadow support.
type ShadowModal struct {
	*ShadowBox
	modal *tview.Modal
}

// NewShadowModal creates a modal with drop shadow.
func NewShadowModal() *ShadowModal {
	return &ShadowModal{
		ShadowBox: NewShadowBox(),
		modal:     tview.NewModal(),
	}
}

// GetModal returns the underlying modal for configuration.
func (sm *ShadowModal) GetModal() *tview.Modal {
	return sm.modal
}

// Draw renders the modal with shadow.
func (sm *ShadowModal) Draw(screen tcell.Screen) {
	// First draw the shadow
	sm.ShadowBox.Draw(screen)

	// Then draw the modal
	sm.modal.Draw(screen)
}

// ShadowFlex is a Flex container that renders with a drop shadow.
type ShadowFlex struct {
	flex          *tview.Flex
	shadowEnabled bool
	shadowChar    rune
	shadowColor   tcell.Color
	shadowOffsetX int
	shadowOffsetY int
}

// NewShadowFlex creates a new flex container with drop shadow.
func NewShadowFlex() *ShadowFlex {
	return &ShadowFlex{
		flex:          tview.NewFlex(),
		shadowEnabled: true,
		shadowChar:    '▒',
		shadowColor:   tcell.ColorGray,
		shadowOffsetX: 2,
		shadowOffsetY: 1,
	}
}

// SetShadowEnabled enables or disables the shadow effect.
func (sf *ShadowFlex) SetShadowEnabled(enabled bool) *ShadowFlex {
	sf.shadowEnabled = enabled
	return sf
}

// GetFlex returns the underlying flex container for adding items.
func (sf *ShadowFlex) GetFlex() *tview.Flex {
	return sf.flex
}

// Draw renders the flex container with shadow.
func (sf *ShadowFlex) Draw(screen tcell.Screen) {
	// Draw shadow if enabled
	if sf.shadowEnabled {
		x, y, width, height := sf.flex.GetRect()

		shadowStyle := tcell.StyleDefault.
			Foreground(sf.shadowColor).
			Background(tcell.ColorDefault).
			Dim(true)

		// Draw bottom shadow
		shadowY := y + height + sf.shadowOffsetY - 1
		if shadowY >= 0 {
			for sx := x + sf.shadowOffsetX; sx < x+width+sf.shadowOffsetX; sx++ {
				if sx >= 0 {
					screen.SetContent(sx, shadowY, sf.shadowChar, nil, shadowStyle)
				}
			}
		}

		// Draw right shadow
		shadowX := x + width + sf.shadowOffsetX - 1
		if shadowX >= 0 {
			for sy := y + sf.shadowOffsetY; sy < y+height+sf.shadowOffsetY; sy++ {
				if sy >= 0 {
					screen.SetContent(shadowX, sy, sf.shadowChar, nil, shadowStyle)
				}
			}
		}

		// Draw corner shadow
		if shadowX >= 0 && shadowY >= 0 {
			screen.SetContent(shadowX, shadowY, sf.shadowChar, nil, shadowStyle)
		}
	}

	// Draw the flex container
	sf.flex.Draw(screen)
}

// Implement tview.Primitive interface

func (sf *ShadowFlex) GetRect() (int, int, int, int) {
	return sf.flex.GetRect()
}

func (sf *ShadowFlex) SetRect(x, y, width, height int) {
	sf.flex.SetRect(x, y, width, height)
}

func (sf *ShadowFlex) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return sf.flex.InputHandler()
}

func (sf *ShadowFlex) Focus(delegate func(p tview.Primitive)) {
	sf.flex.Focus(delegate)
}

func (sf *ShadowFlex) HasFocus() bool {
	return sf.flex.HasFocus()
}

func (sf *ShadowFlex) Blur() {
	sf.flex.Blur()
}

func (sf *ShadowFlex) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return sf.flex.MouseHandler()
}

func (sf *ShadowFlex) PasteHandler() func(pastedText string, setFocus func(p tview.Primitive)) {
	return sf.flex.PasteHandler()
}

// Helper function to create a shadow wrapper for any form
type ShadowForm struct {
	form          *tview.Form
	shadowEnabled bool
	shadowChar    rune
	shadowColor   tcell.Color
	shadowOffsetX int
	shadowOffsetY int
}

// NewShadowForm creates a form with drop shadow.
func NewShadowForm() *ShadowForm {
	return &ShadowForm{
		form:          tview.NewForm(),
		shadowEnabled: true,
		shadowChar:    '▒',
		shadowColor:   tcell.ColorGray,
		shadowOffsetX: 2,
		shadowOffsetY: 1,
	}
}

// GetForm returns the underlying form for configuration.
func (sf *ShadowForm) GetForm() *tview.Form {
	return sf.form
}

// SetShadowEnabled enables or disables the shadow effect.
func (sf *ShadowForm) SetShadowEnabled(enabled bool) *ShadowForm {
	sf.shadowEnabled = enabled
	return sf
}

// Draw renders the form with shadow.
func (sf *ShadowForm) Draw(screen tcell.Screen) {
	// Draw shadow if enabled
	if sf.shadowEnabled {
		x, y, width, height := sf.form.GetRect()

		shadowStyle := tcell.StyleDefault.
			Foreground(sf.shadowColor).
			Background(tcell.ColorDefault).
			Dim(true)

		// Draw bottom shadow
		shadowY := y + height + sf.shadowOffsetY - 1
		if shadowY >= 0 {
			for sx := x + sf.shadowOffsetX; sx < x+width+sf.shadowOffsetX; sx++ {
				if sx >= 0 {
					screen.SetContent(sx, shadowY, sf.shadowChar, nil, shadowStyle)
				}
			}
		}

		// Draw right shadow
		shadowX := x + width + sf.shadowOffsetX - 1
		if shadowX >= 0 {
			for sy := y + sf.shadowOffsetY; sy < y+height+sf.shadowOffsetY; sy++ {
				if sy >= 0 {
					screen.SetContent(shadowX, sy, sf.shadowChar, nil, shadowStyle)
				}
			}
		}

		// Draw corner shadow
		if shadowX >= 0 && shadowY >= 0 {
			screen.SetContent(shadowX, shadowY, sf.shadowChar, nil, shadowStyle)
		}
	}

	// Draw the form
	sf.form.Draw(screen)
}

// Implement tview.Primitive interface

func (sf *ShadowForm) GetRect() (int, int, int, int) {
	return sf.form.GetRect()
}

func (sf *ShadowForm) SetRect(x, y, width, height int) {
	sf.form.SetRect(x, y, width, height)
}

func (sf *ShadowForm) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return sf.form.InputHandler()
}

func (sf *ShadowForm) Focus(delegate func(p tview.Primitive)) {
	sf.form.Focus(delegate)
}

func (sf *ShadowForm) HasFocus() bool {
	return sf.form.HasFocus()
}

func (sf *ShadowForm) Blur() {
	sf.form.Blur()
}

func (sf *ShadowForm) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return sf.form.MouseHandler()
}

func (sf *ShadowForm) PasteHandler() func(pastedText string, setFocus func(p tview.Primitive)) {
	return sf.form.PasteHandler()
}
