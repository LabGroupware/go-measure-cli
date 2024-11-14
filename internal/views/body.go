package views

import (
	"github.com/rivo/tview"
)

const (
	homePage      = "home"
	wsRawSendPage = "ws-raw-send"
)

func createBodyView() *tview.Pages {
	pages := tview.NewPages()

	pages.AddPage(homePage, createHomePage(), true, true)
	pages.AddPage(wsRawSendPage, createWsRawSendPage(), true, true)

	subscribeSidebarEvent(func(event sidebarEvent) {
		switch event.Key {
		case sidebarKeyGoToHome:
			pages.SwitchToPage(homePage)
		case sidebarKeyGoToRawWs:
			pages.SwitchToPage(wsRawSendPage)
		}
	})

	pages.SendToFront(homePage)

	return pages
}
