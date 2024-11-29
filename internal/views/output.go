package views

import (
	"fmt"
	"sync"

	"github.com/rivo/tview"
)

type outputCmdType string

const (
	outputCmdTypeClear outputCmdType = "clear"
	outputCmdTypeWrite outputCmdType = "write"
)

type outputValue struct {
	value any
	cmd   outputCmdType
}

var outputBroadcast = NewBroadcaster[outputValue]()

func createOutputView() *tview.TextView {
	codeView := tview.NewTextView().
		SetWrap(false).
		SetDynamicColors(true)
	codeView.SetBorderPadding(1, 1, 2, 0)

	mu := sync.Mutex{}

	ch := outputBroadcast.subscribe()
	go func(ch chan outputValue) {
		for value := range ch {
			mu.Lock()
			switch value.cmd {
			case outputCmdTypeClear:
				codeView.Clear()
			case outputCmdTypeWrite:
				fmt.Fprint(codeView, value.value)
			}
			mu.Unlock()
		}
	}(ch)

	return codeView
}

func ConsoleOutput(s string) {
	outputBroadcast.broadcast(outputValue{value: s, cmd: outputCmdTypeWrite})
}

func ClearOutput() {
	outputBroadcast.broadcast(outputValue{cmd: outputCmdTypeClear})
}
