/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/LabGroupware/go-measure-tui/internal/ws"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "A simple TUI for measuring web performance",
	Long: `For measuring web performance, Go Measure TUI is a simple terminal user interface that
allows you to measure the performance of a website by providing the URL and the number of requests
to be made to the website. It then displays the average time taken to make the requests.

This command is used to measure the performance of a website by providing the URL and the number of requests to be made to the website.`,
	Run: func(cmd *cobra.Command, args []string) {
		wSock := ws.NewWebSocket()
		done, err := wSock.Connect(container.Config.Web.WebSocket.Url)
		if err != nil {
			fmt.Printf("failed to connect to WebSocket server: %v\n", err)
			return
		}
		defer wSock.Close()

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		app := tview.NewApplication()

		go func() {
			for {
				select {
				case <-done:
					fmt.Println("WebSocket connection closed")
					app.Stop()
					return
				case <-interrupt:
					fmt.Println("Interrupt signal received")
					app.Stop()
					return
				}
			}
		}()

		<-interrupt

		// var t tview.Primitive
		// if t, err = views.CreatePrimitive(); err != nil {
		// 	fmt.Printf("failed to create primitive: %v\n", err)
		// }

		// if err := app.SetRoot(t, true).SetFocus(t).Run(); err != nil {
		// 	fmt.Printf("failed to run TUI: %v\n", err)
		// }
	},
}

func init() {
	connectCmd.AddCommand(webCmd)
}
