// Package views provides the templates for the web application.
package views

import "github.com/rivo/tview"

func CreatePrimitive() (tview.Primitive, error) {
	return tview.NewBox().SetBorder(true).SetTitle("Request"), nil
}
