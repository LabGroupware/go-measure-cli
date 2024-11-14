package views

import "github.com/rivo/tview"

const (
	subtitle   = `tview - Rich Widgets for Terminal UIs`
	navigation = `[yellow]Ctrl-N[-]: Next slide    [yellow]Ctrl-P[-]: Previous slide    [yellow]Ctrl-C[-]: Exit`
	mouse      = `(or use your mouse)`
)

func createHomePage() tview.Primitive {
	frame := tview.NewFrame(tview.NewBox()).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(subtitle, true, tview.AlignCenter, contractColor).
		AddText("", true, tview.AlignCenter, contractColor).
		AddText(navigation, true, tview.AlignCenter, successColor).
		AddText("", true, tview.AlignCenter, contractColor).
		AddText(mouse, true, tview.AlignCenter, mainColor)

	return frame
}
