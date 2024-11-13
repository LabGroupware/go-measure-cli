package views

import (
	"github.com/rivo/tview"
)

// NewMenu creates a new menu view
func NewMenu(navigate func(page string)) *tview.List {
	menu := tview.NewList().
		AddItem("Go to Form", "Open the form view", 'f', func() {
			navigate("form")
		}).
		AddItem("Go to List", "Open the list view", 'l', func() {
			navigate("list")
		}).
		AddItem("Quit", "Exit the application", 'q', func() {
			navigate("")
		})
	return menu
}
