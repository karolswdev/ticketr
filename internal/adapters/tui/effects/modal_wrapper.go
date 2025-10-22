package effects

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ModalWrapper provides a centered, properly-sized modal container with optional shadow.
// It ensures modals are:
// - Centered both horizontally and vertically
// - Sized appropriately (not too wide, not too narrow)
// - Have proper padding and spacing
// - Display drop shadows when enabled
// - Responsive to terminal size
type ModalWrapper struct {
	*tview.Flex
	content       tview.Primitive
	shadowEnabled bool
	shadowChar    rune
	shadowColor   tcell.Color
	minWidth      int
	maxWidth      int
	minHeight     int
	maxHeight     int
}

// NewModalWrapper creates a new modal wrapper that centers and sizes content appropriately.
func NewModalWrapper(content tview.Primitive) *ModalWrapper {
	wrapper := &ModalWrapper{
		Flex:          tview.NewFlex(),
		content:       content,
		shadowEnabled: true,
		shadowChar:    '▒',
		shadowColor:   tcell.ColorGray,
		minWidth:      40,  // Minimum comfortable width
		maxWidth:      80,  // Maximum width (professional, not too wide)
		minHeight:     10,  // Minimum height
		maxHeight:     40,  // Maximum height
	}

	wrapper.setupLayout()
	return wrapper
}

// setupLayout creates the centered layout structure.
func (mw *ModalWrapper) setupLayout() {
	// Create vertical centering
	mw.SetDirection(tview.FlexRow)

	// Add top spacer (flexible)
	mw.AddItem(nil, 0, 1, false)

	// Create horizontal centering for content row
	horizontalFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Add left spacer (flexible)
	horizontalFlex.AddItem(nil, 0, 1, false)

	// Add content with fixed or proportional size
	// We'll use a fixed size that adapts to terminal size
	if mw.shadowEnabled {
		// Wrap content with shadow
		shadowContent := NewShadowFlex()
		shadowContent.GetFlex().AddItem(mw.content, 0, 1, true)
		horizontalFlex.AddItem(shadowContent, mw.maxWidth, 1, true)
	} else {
		horizontalFlex.AddItem(mw.content, mw.maxWidth, 1, true)
	}

	// Add right spacer (flexible)
	horizontalFlex.AddItem(nil, 0, 1, false)

	// Add the horizontal flex to the main vertical flex
	mw.AddItem(horizontalFlex, 0, 1, true)

	// Add bottom spacer (flexible)
	mw.AddItem(nil, 0, 1, false)
}

// SetShadowEnabled enables or disables drop shadows.
func (mw *ModalWrapper) SetShadowEnabled(enabled bool) *ModalWrapper {
	mw.shadowEnabled = enabled
	mw.setupLayout()
	return mw
}

// SetSizeConstraints sets the size constraints for the modal.
func (mw *ModalWrapper) SetSizeConstraints(minWidth, maxWidth, minHeight, maxHeight int) *ModalWrapper {
	mw.minWidth = minWidth
	mw.maxWidth = maxWidth
	mw.minHeight = minHeight
	mw.maxHeight = maxHeight
	mw.setupLayout()
	return mw
}

// SetMaxWidth sets the maximum width of the modal.
func (mw *ModalWrapper) SetMaxWidth(width int) *ModalWrapper {
	mw.maxWidth = width
	mw.setupLayout()
	return mw
}

// InputHandler passes input to the wrapped content.
func (mw *ModalWrapper) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return mw.Flex.InputHandler()
}

// CenteredModal creates a properly centered and sized modal from a tview.Modal.
// This is a convenience function for simple modals.
func CenteredModal(modal *tview.Modal, shadowEnabled bool) *ModalWrapper {
	wrapper := NewModalWrapper(modal)
	wrapper.SetShadowEnabled(shadowEnabled)
	wrapper.SetMaxWidth(60) // Modals are typically narrower
	return wrapper
}

// CenteredForm creates a properly centered and sized form modal.
// This is a convenience function for form-based modals.
func CenteredForm(form *tview.Form, shadowEnabled bool) *ModalWrapper {
	wrapper := NewModalWrapper(form)
	wrapper.SetShadowEnabled(shadowEnabled)
	wrapper.SetMaxWidth(70) // Forms can be slightly wider
	return wrapper
}

// CenteredPages creates a properly centered and sized pages container.
// This is useful for multi-page modals like bulk operations.
func CenteredPages(pages *tview.Pages, shadowEnabled bool) *ModalWrapper {
	wrapper := NewModalWrapper(pages)
	wrapper.SetShadowEnabled(shadowEnabled)
	wrapper.SetMaxWidth(75) // Multi-page modals can be wider
	return wrapper
}

// ShowModal displays a modal as a page overlay (proper modal behavior).
// This should be preferred over app.SetRoot() for modal dialogs.
func ShowModal(pages *tview.Pages, name string, modal tview.Primitive, shadowEnabled bool) {
	wrapper := NewModalWrapper(modal)
	wrapper.SetShadowEnabled(shadowEnabled)
	pages.AddPage(name, wrapper, true, true)
}

// HideModal removes a modal from the pages stack.
func HideModal(pages *tview.Pages, name string) {
	pages.RemovePage(name)
}

// ModalPage represents a modal that can be shown/hidden in a pages container.
type ModalPage struct {
	pages         *tview.Pages
	name          string
	content       tview.Primitive
	wrapper       *ModalWrapper
	shadowEnabled bool
}

// NewModalPage creates a new modal page manager.
func NewModalPage(pages *tview.Pages, name string, content tview.Primitive) *ModalPage {
	return &ModalPage{
		pages:         pages,
		name:          name,
		content:       content,
		shadowEnabled: true,
	}
}

// Show displays the modal.
func (mp *ModalPage) Show() {
	mp.wrapper = NewModalWrapper(mp.content)
	mp.wrapper.SetShadowEnabled(mp.shadowEnabled)
	mp.pages.AddPage(mp.name, mp.wrapper, true, true)
}

// Hide removes the modal.
func (mp *ModalPage) Hide() {
	mp.pages.RemovePage(mp.name)
}

// SetShadowEnabled enables or disables shadows for this modal.
func (mp *ModalPage) SetShadowEnabled(enabled bool) *ModalPage {
	mp.shadowEnabled = enabled
	if mp.wrapper != nil {
		mp.wrapper.SetShadowEnabled(enabled)
	}
	return mp
}
