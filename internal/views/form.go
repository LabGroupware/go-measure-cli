package views

import (
	"github.com/rivo/tview"
)

// NewForm creates a new form view
func NewForm() *tview.Form {
	form := tview.NewForm().
		AddInputField("Name", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Submit", func() {
			// Submit action
		}).
		AddButton("Back", func() {
			// Back to menu
		})
	form.SetBorder(true).SetTitle("Form View").SetTitleAlign(tview.AlignLeft)
	return form
}
