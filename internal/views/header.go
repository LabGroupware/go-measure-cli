package views

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createHeaderView() *tview.Grid {
	header := tview.NewGrid().
		SetRows(1).
		AddItem(displayCtrCBackAction(), 0, 0, 1, 1, 0, 0, false).
		AddItem(createAppTitle(), 0, 1, 1, 1, 0, 0, false).
		AddItem(displayWidthHeight(), 0, 2, 1, 1, 0, 0, false)

	subscribeWidthResize(func(s vSize) {
		switch s.size {
		case Large:
			header.SetColumns(15, 0, 15)
		case Medium:
			header.SetColumns(10, 0, 12)
		case Small:
			header.SetColumns(7, 0, 12)
		}
	})

	return header
}

func displayWidthHeight() tview.Primitive {
	whDisplay := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetTextStyle(tcell.StyleDefault.Foreground(warningColor))

	var w, h int

	subscribeWidthResize(func(s vSize) {
		w = s.v

		switch s.size {
		case Large:
			whDisplay.SetText(fmt.Sprintf("w:%d / h:%d", w, h))
		case Medium, Small:
			whDisplay.SetText(fmt.Sprintf("w:%d/h:%d", w, h))
		}
	})

	subscribeHeightResize(func(s vSize) {
		h = s.v

		switch s.size {
		case Large:
			whDisplay.SetText(fmt.Sprintf("w:%d / h:%d", w, h))
		case Medium, Small:
			whDisplay.SetText(fmt.Sprintf("w:%d/h:%d", w, h))
		}
	})
	return whDisplay
}

func createAppTitle() tview.Primitive {
	return tview.NewTextView().
		SetText("Go Measure TUI").
		SetTextAlign(tview.AlignCenter)
}

func displayCtrCBackAction() tview.Primitive {
	view := tview.NewTextView().SetTextAlign(tview.AlignLeft)

	subscribeWidthResize(func(s vSize) {
		var str string
		switch s.size {
		case Large:
			str = "  Ctrl+C to exit"
		case Medium:
			str = " Ctrl+C"
		case Small:
			str = " C+C"
		}
		view.SetText(str)
	})

	return view
}
