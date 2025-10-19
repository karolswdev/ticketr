package views

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/karolswdev/ticktr/internal/adapters/tui/search"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/rivo/tview"
)

// SearchView implements the View interface for search functionality.
type SearchView struct {
	modal       *tview.Flex
	inputField  *tview.InputField
	resultsList *tview.List
	statusBar   *tview.TextView
	app         *tview.Application

	// Data
	allTickets []*domain.Ticket
	matches    []*search.Match

	// Callbacks
	onClose   func()
	onSelect  func(*domain.Ticket)
	lastQuery string
}

// NewSearchView creates a new search view.
func NewSearchView(app *tview.Application) *SearchView {
	v := &SearchView{
		app:        app,
		allTickets: []*domain.Ticket{},
		matches:    []*search.Match{},
	}

	// Create input field for search query
	v.inputField = tview.NewInputField().
		SetLabel("= Search: ").
		SetFieldWidth(0). // Full width
		SetPlaceholder("Type to search... (@user #ID !priority ~sprint /regex/)")

	v.inputField.SetBorder(true).
		SetTitle(" Search Mode ").
		SetBorderColor(tcell.ColorGreen)

	// Create results list
	v.resultsList = tview.NewList().
		ShowSecondaryText(true)

	v.resultsList.SetBorder(true).
		SetTitle(" Results (0) ").
		SetBorderColor(tcell.ColorWhite)

	// Create status bar
	v.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]Filters:[white] @assignee #id !priority ~sprint /regex/ [yellow][white] [green]Enter[white]=Open [yellow][white] [red]Esc[white]=Close")

	v.statusBar.SetBorder(false)

	// Create modal layout
	v.modal = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(v.inputField, 3, 0, true).
		AddItem(v.resultsList, 0, 1, false).
		AddItem(v.statusBar, 1, 0, false)

	// Set up input field handler
	v.inputField.SetChangedFunc(func(text string) {
		v.performSearch(text)
	})

	// Set up input field key handler
	v.inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			if v.onClose != nil {
				v.onClose()
			}
			return nil
		case tcell.KeyDown:
			// Move focus to results
			v.app.SetFocus(v.resultsList)
			return nil
		case tcell.KeyEnter:
			// Select first result if any
			if len(v.matches) > 0 && v.onSelect != nil {
				v.onSelect(v.matches[0].Ticket)
			}
			return nil
		}

		// Handle vim-style navigation in input field
		switch event.Rune() {
		case 'j':
			if len(v.inputField.GetText()) == 0 {
				v.app.SetFocus(v.resultsList)
				return nil
			}
		}

		return event
	})

	// Set up results list key handler
	v.resultsList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			if v.onClose != nil {
				v.onClose()
			}
			return nil
		case tcell.KeyUp:
			// If at top of list, return to input
			if v.resultsList.GetCurrentItem() == 0 {
				v.app.SetFocus(v.inputField)
				return nil
			}
		}

		return event
	})

	// Set up selection handler
	v.resultsList.SetSelectedFunc(func(index int, _ string, _ string, _ rune) {
		if index < len(v.matches) && v.onSelect != nil {
			v.onSelect(v.matches[index].Ticket)
		}
	})

	return v
}

// Name returns the view's identifier.
func (v *SearchView) Name() string {
	return "search"
}

// Primitive returns the view's root primitive.
func (v *SearchView) Primitive() tview.Primitive {
	return v.modal
}

// OnShow is called when the view becomes active.
func (v *SearchView) OnShow() {
	// Clear previous search
	v.inputField.SetText("")
	v.resultsList.Clear()
	v.lastQuery = ""
	v.matches = []*search.Match{}
	v.updateResultsTitle()
}

// OnHide is called when the view is hidden.
func (v *SearchView) OnHide() {
	// Nothing to do
}

// SetTickets sets the tickets to search through.
func (v *SearchView) SetTickets(tickets []*domain.Ticket) {
	v.allTickets = tickets
}

// SetOnClose sets the callback for closing the search view.
func (v *SearchView) SetOnClose(callback func()) {
	v.onClose = callback
}

// SetOnSelect sets the callback for selecting a ticket.
func (v *SearchView) SetOnSelect(callback func(*domain.Ticket)) {
	v.onSelect = callback
}

// performSearch executes the search and updates the results.
func (v *SearchView) performSearch(queryStr string) {
	v.lastQuery = queryStr

	// Clear current results
	v.resultsList.Clear()

	if queryStr == "" {
		v.matches = []*search.Match{}
		v.updateResultsTitle()
		return
	}

	// Parse query
	query, _ := search.ParseQuery(queryStr)

	// Apply filters
	filtered := search.ApplyFilters(v.allTickets, query)

	// Perform fuzzy search on filtered results
	v.matches = search.SearchTickets(filtered, query.Text)

	// Populate results list
	for i, match := range v.matches {
		v.addResult(i, match)
	}

	v.updateResultsTitle()
}

// addResult adds a search result to the list.
func (v *SearchView) addResult(index int, match *search.Match) {
	// Main text: JiraID + Title + Score
	mainText := fmt.Sprintf("[%d%%] %s: %s",
		match.Score,
		match.Ticket.JiraID,
		match.Ticket.Title)

	// Secondary text: where matches were found
	var matchInfo string
	if len(match.MatchedIn) > 0 {
		matchInfo = fmt.Sprintf("Matched in: %s", strings.Join(match.MatchedIn, ", "))
	} else {
		matchInfo = "No specific matches"
	}

	v.resultsList.AddItem(mainText, matchInfo, 0, nil)
}

// updateResultsTitle updates the results panel title with count.
func (v *SearchView) updateResultsTitle() {
	title := fmt.Sprintf(" Results (%d) ", len(v.matches))
	v.resultsList.SetTitle(title)
}

// Note: Focus management between inputField and resultsList is handled
// by the application's SetFocus method, not here.
