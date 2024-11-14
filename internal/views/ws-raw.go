package views

import "github.com/rivo/tview"

func createWsRawSendPage() tview.Primitive {
	form := tview.NewForm().
		AddInputField("URL", "", 0, nil, nil).
		AddInputField("Payload", "", 0, nil, nil).
		AddButton("Send", nil).
		AddButton("Back", nil)

	return form
}
