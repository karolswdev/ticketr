package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/karolswdev/ticktr/internal/adapters/tui/widgets"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Create demo action bar with multiple keybindings
	actionBar := widgets.NewActionBar()
	actionBar.SetApp(app)

	// Set ticket tree context (has many bindings)
	actionBar.SetContext(widgets.ContextTicketTree)

	// Create info box to explain the demo
	infoBox := tview.NewTextView()
	infoBox.SetDynamicColors(true)
	infoBox.SetBorder(true)
	infoBox.SetTitle(" Marquee Animation Demo ")
	infoBox.SetText(`[yellow]Marquee Animation Demo[white]

This demonstrates the new theatrical marquee animation system.

[green]Animation Phases:[white]
1. [cyan]Slide-In:[white] Items slide in from right (with angle brackets ⟨⟩ for visual momentum)
2. [cyan]Center Hold:[white] Items stop in center and blink every 150ms for 1.5 seconds
3. [cyan]Slide-Out:[white] Items slide out to left (angle brackets return)

[green]How to Test:[white]
- Resize your terminal to [yellow]< 80 columns[white] to trigger marquee mode
- Watch the action bar at the bottom animate each keybinding
- Resize to [yellow]> 150 columns[white] to see all bindings at once (no animation)
- Press [red]Ctrl+C[white] to exit

[green]Current Terminal Width:[white] Will update dynamically...
[green]Marquee Active:[white] Will show if animation is running...
`)

	// Create layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(infoBox, 0, 1, false).
		AddItem(actionBar, 3, 0, false)

	// Update width display periodically
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			app.QueueUpdateDraw(func() {
				_, _, width, _ := actionBar.GetInnerRect()
				isMarquee := actionBar.GetContext() == widgets.ContextTicketTree && width < 140

				status := fmt.Sprintf(`[yellow]Marquee Animation Demo[white]

This demonstrates the new theatrical marquee animation system.

[green]Animation Phases:[white]
1. [cyan]Slide-In:[white] Items slide in from right (with angle brackets ⟨⟩ for visual momentum)
2. [cyan]Center Hold:[white] Items stop in center and blink every 150ms for 1.5 seconds
3. [cyan]Slide-Out:[white] Items slide out to left (angle brackets return)

[green]How to Test:[white]
- Resize your terminal to [yellow]< 80 columns[white] to trigger marquee mode
- Watch the action bar at the bottom animate each keybinding
- Resize to [yellow]> 150 columns[white] to see all bindings at once (no animation)
- Press [red]Ctrl+C[white] to exit

[green]Current Terminal Width:[white] %d columns
[green]Marquee Active:[white] %v
`, width, isMarquee)

				infoBox.SetText(status)
			})
		}
	}()

	// Set up input handling
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			app.Stop()
			return nil
		}
		return event
	})

	// Run application
	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	// Cleanup
	actionBar.Shutdown()
}
