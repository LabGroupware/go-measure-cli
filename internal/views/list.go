package views

import (
	"github.com/rivo/tview"
)

// NewList creates a new list view
func NewList() *tview.List {
	list := tview.NewList().
		AddItem("Item 1", "This is the first item", '1', nil).
		AddItem("Item 2", "This is the second item", '2', nil).
		AddItem("Back", "Return to the menu", 'b', func() {
			// Back to menu
		})
	list.SetBorder(true).SetTitle("List View").SetTitleAlign(tview.AlignLeft)
	return list
}
