/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/LabGroupware/go-measure-tui/internal/job"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
	"github.com/LabGroupware/go-measure-tui/internal/views"
	"github.com/LabGroupware/go-measure-tui/internal/ws"
	"github.com/fatih/color"
	"github.com/gdamore/tcell/v2"
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
		wsEventHandler := ws.NewDetailEventResponseMessageHandler()
		wsEventHandler.RegisterHandleFunc(
			job.UserProfileCreatedJobBegin,
			ws.EventHandlerAdapter[ws.EventResponseMessageWithData[job.CreateUserProfileJobSuccessData]]{
				Handler: CreateUserProfileJobBeganHandler{},
			},
		)
		wSock.EventMsgHandler = wsEventHandler.HandleMessage
		done, err := wSock.Connect(container.Config.Web.WebSocket.Url, container.AuthToken.AccessToken)
		if err != nil {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Println(red(fmt.Sprintf("failed to connect to WebSocket server: %v", err)))
			return
		}
		defer wSock.Close()

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		app := tview.NewApplication()

		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyCtrlC:
				interrupt <- os.Interrupt
			}
			return event
		})

		go func() {
			for {
				select {
				case <-done:
					fmt.Println("WebSocket connection closed")
					app.Stop()
					green := color.New(color.FgGreen).SprintFunc()
					fmt.Println(green("Successfully application closed"))
					return
				case <-interrupt:
					fmt.Println("Interrupt signal received")
					if err := wSock.SendCloseMessage(); err != nil {
						red := color.New(color.FgRed).SprintFunc()
						fmt.Println(red(fmt.Sprintf("failed to send close message: %v", err)))
						return
					}
					app.Stop()
					return
				}
			}
		}()

		if err := views.Run(ctx, app, *container); err != nil {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Println(red(fmt.Sprintf("failed to run application: %v", err)))
			container.Logger.Error(ctx, "failed to run application", logger.Value("error", err))
		}

		green := color.New(color.FgGreen).SprintFunc()
		fmt.Println(green("Successfully application closed"))
	},
}

func init() {
	connectCmd.AddCommand(webCmd)
}

type CreateUserProfileJobBeganHandler struct{}

func (h CreateUserProfileJobBeganHandler) Handle(
	ws *ws.WebSocket,
	data ws.EventResponseMessageWithData[job.CreateUserProfileJobSuccessData],
	raw []byte,
) {
	views.ConsoleOutput("Handling CreateUserProfileJobBegan:" + data.Data.JobID)
}
