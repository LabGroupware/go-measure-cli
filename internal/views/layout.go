package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createGridLayout() *tview.Grid {
	grid := tview.NewGrid().SetBorders(true)

	grid.SetBackgroundColor(mainColor)

	subscribeWidthResize(func(s vSize) {
		v := calculateWidthSize(s.size)
		grid.SetColumns(v, 0)
	})

	subscribeHeightResize(func(s vSize) {
		v := calculateHeightSize(s.size)
		grid.SetRows(1, 0, v)
	})

	header := createHeaderView()
	sidebar := createSidebarView()
	body := createBodyView()
	output := createOutputView()

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlH, tcell.KeyCtrlL:
			if body.HasFocus() {
				requestSetFocus(sidebar)
				sidebar.SetBorder(true)
				body.SetBorder(false)
			} else {
				requestSetFocus(body)
				body.SetBorder(true)
				sidebar.SetBorder(false)
			}
		}
		return event
	})

	grid.AddItem(header, 0, 0, 1, 2, 0, 0, false)
	grid.AddItem(sidebar, 1, 0, 1, 1, 0, 0, true)
	grid.AddItem(body, 1, 1, 1, 1, 0, 0, true)
	grid.AddItem(output, 2, 0, 1, 2, 0, 0, false)

	sidebar.SetBorder(true)
	sidebar.SetBorderColor(successColor)
	body.SetBorderColor(successColor)
	body.SetBorder(false)

	return grid
}

func calculateWidthSize(s size) int {
	switch s {
	case Large:
		return 20
	case Medium:
		return 15
	default:
		return 10
	}
}

func calculateHeightSize(s size) int {
	switch s {
	case Large:
		return 7
	case Medium:
		return 5
	default:
		return 3
	}
}
