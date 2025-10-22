package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/karolswdev/ticktr/internal/adapters/tui/effects"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	// Create background animator
	bgConfig := effects.BackgroundConfig{
		Effect:    effects.BackgroundHyperspace,
		Density:   0.05, // 5% for visibility
		Speed:     100,
		Enabled:   true,
		MaxFPS:    15,
		AutoPause: false,
	}

	fmt.Fprintf(os.Stderr, "Creating background animator...\n")
	animator := effects.NewBackgroundAnimator(app, bgConfig)

	fmt.Fprintf(os.Stderr, "Starting animator...\n")
	animator.Start()

	fmt.Fprintf(os.Stderr, "Getting overlay...\n")
	overlay := animator.GetOverlay()

	if overlay == nil {
		fmt.Fprintf(os.Stderr, "ERROR: Overlay is nil!\n")
		return
	}

	// Create a simple text view
	textView := tview.NewTextView()
	textView.SetText("Cosmic Background Test\n\nYou should see stars moving from left to right.\n\nPress 'q' to quit.")
	textView.SetBorder(true)
	textView.SetTitle("Test")

	// Create pages with layering
	pages := tview.NewPages()
	pages.AddPage("cosmic", overlay, true, true)
	pages.AddPage("content", textView, true, true)

	// Set input capture
	pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.Stop()
			return nil
		}
		return event
	})

	app.SetRoot(pages, true)

	// Give animator time to spawn particles
	go func() {
		time.Sleep(500 * time.Millisecond)
		app.QueueUpdateDraw(func() {
			fmt.Fprintf(os.Stderr, "Initial draw triggered\n")
		})
	}()

	if err := app.Run(); err != nil {
		panic(err)
	}

	animator.Shutdown()
}
