// Package utils provides utility functions for the application
package utils

import (
	"github.com/rivo/tview"
)

// SetFocus sets the focus to the specified primitive
func SetFocus(app *tview.Application, primitive tview.Primitive) {
	app.SetFocus(primitive)
}
