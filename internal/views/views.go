// Package views provides the templates for the web application.
package views

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var drawBroadcast = NewBroadcaster[struct{}]()

func requestDraw() {
	drawBroadcast.broadcast(struct{}{})
}

var setFocusBroadcast = NewBroadcaster[tview.Primitive]()

func requestSetFocus(p tview.Primitive) {
	setFocusBroadcast.broadcast(p)
}

func Run(ctx context.Context, tapp *tview.Application, ctr app.Container) error {
	// Init Colors
	initColors(ctr)

	t, err := createPrimitive()
	if err != nil {
		return fmt.Errorf("failed to create primitive: %w", err)
	}

	var lastWidth, lastHeight int
	var mu sync.Mutex

	tapp.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		width, height := screen.Size()

		mu.Lock()
		defer mu.Unlock()

		if width != lastWidth {
			lastWidth = width
			broadcastWidthSize(width)
		}

		if height != lastHeight {
			lastHeight = height
			broadcastHeightSize(height)
		}
		// screen.Clear()
		return false
	})

	// Redraw
	ch := drawBroadcast.subscribe()
	go func(ch chan struct{}) {
		for range ch {
			tapp.Draw()
		}
	}(ch)

	// Set Focus
	chFocus := setFocusBroadcast.subscribe()
	go func(ch chan tview.Primitive) {
		for p := range ch {
			tapp.SetFocus(p)
		}
	}(chFocus)

	// Quit
	subscribeSidebarEvent(func(event sidebarEvent) {
		switch event.Key {
		case sidebarKeyQuit:
			tapp.Stop()
		}
	})

	if err := tapp.SetRoot(t, true).SetFocus(t).Run(); err != nil {
		return fmt.Errorf("failed to run application: %w", err)
	}

	return nil
}

func createPrimitive() (tview.Primitive, error) {

	grid := createGridLayout()

	var welcomeMsg = []byte{
		'W', 'e', 'l', 'c', 'o', 'm', 'e', ' ',
		't', 'o', ' ',
		't', 'h', 'e', ' ',
		'T', 'U', 'I', ' ',
		'a', 'p', 'p', 'l', 'i', 'c', 'a', 't', 'i', 'o', 'n',
	}

	go func() {
		time.Sleep(100 * time.Microsecond)
		for _, b := range welcomeMsg {
			time.Sleep(100 * time.Millisecond)
			ConsoleOutput(string(b))
			requestDraw()
		}
		time.Sleep(100 * time.Millisecond)
		ConsoleOutput("\n")
		requestDraw()
	}()

	return grid, nil
}
