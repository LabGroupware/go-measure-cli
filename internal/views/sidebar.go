package views

import (
	"github.com/rivo/tview"
)

type sidebarKey string

const (
	sidebarKeyGoToHome  = "sidebarKeyGoToHome"
	sidebarKeyGoToRawWs = "sidebarKeyGoToRawWs"

	sidebarKeyQuit = "sidebarKeyQuit"
)

type sidebarEvent struct {
	Key sidebarKey
}

var sidebarBroadcast = NewBroadcaster[sidebarEvent]()

func subscribeSidebarEvent(consumer func(event sidebarEvent)) {
	ch := sidebarBroadcast.subscribe()
	go func(ch chan sidebarEvent) {
		for value := range ch {
			consumer(value)
		}
	}(ch)
}

func createSidebarView() *tview.List {

	var items = []struct {
		Name        string
		Description string
		Key         rune
		Action      func()
	}{
		{
			Name:        "Home",
			Description: "Go Home",
			Key:         'h',
			Action: func() {
				sidebarBroadcast.broadcast(sidebarEvent{Key: sidebarKeyGoToHome})
			},
		},
		{
			Name:        "WS Raw Send",
			Description: "Go to WS Raw Send",
			Key:         'w',
			Action: func() {
				sidebarBroadcast.broadcast(sidebarEvent{Key: sidebarKeyGoToRawWs})
			},
		},
		{
			Name:        "Quit",
			Description: "Press to exit",
			Key:         'q',
			Action: func() {
				sidebarBroadcast.broadcast(sidebarEvent{Key: sidebarKeyQuit})
			},
		},
	}

	list := tview.NewList()
	for _, item := range items {
		item := item
		list.AddItem(item.Name, item.Description, item.Key, item.Action)
	}

	return list
}
