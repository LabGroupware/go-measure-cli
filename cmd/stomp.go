/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// stompCmd represents the stomp command
var stompCmd = &cobra.Command{
	Use:   "stomp",
	Short: "Connect to a remote server using STOMP",
	Long: `Responds to the stomp command. This command is used to connect to a remote server using STOMP.
	For example:
	STOMP connection to a remote server.`,
	Run: func(cmd *cobra.Command, args []string) {
		startTview()
	},
}

func init() {
	connectCmd.AddCommand(stompCmd)
}

func startTview() {
	app := tview.NewApplication()

	// レイアウト
	flex := tview.NewFlex()

	// メニュービュー（左側）
	menu := tview.NewList().
		AddItem("Dashboard", "View system dashboard", 'd', nil).
		AddItem("Form", "Fill in a form", 'f', nil).
		AddItem("Table", "View a table", 't', nil).
		AddItem("Quit", "Exit the application", 'q', func() {
			app.Stop()
		})

	// メインビュー（右側、動的に切り替え）
	mainView := tview.NewTextView().
		SetText("Welcome to the TUI Application!\nSelect an option from the menu.").
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	// フレックスにメニューとメインビューを追加
	flex.AddItem(menu, 20, 1, true). // メニュー：幅 30、サイズ固定
						AddItem(mainView, 0, 3, false) // メインビュー：残りの幅を使用

	// メニューの選択に応じてビューを切り替える
	menu.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		switch mainText {
		case "Dashboard":
			mainView.SetText("[yellow]Dashboard\n[white]This is your system dashboard.\nMore details will be added here.")
		case "Form":
			mainView.SetText("[green]Form\n[white]Navigate to the form view.")
		case "Table":
			mainView.SetText("[blue]Table\n[white]Navigate to the table view.")
		}
	})

	// メニュー項目の選択でビューを切り替える
	menu.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		switch mainText {
		case "Form":
			app.SetRoot(createForm(app, flex), true) // フォームビューに切り替え
		case "Table":
			app.SetRoot(createTable(app, flex), true) // テーブルビューに切り替え
		}
	})

	// 初期ビューを設定してアプリケーションを実行
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

// フォームビューを作成
func createForm(app *tview.Application, previous tview.Primitive) tview.Primitive {
	form := tview.NewForm().
		AddInputField("Name", "", 20, nil, nil).
		AddInputField("Email", "", 30, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Submit", func() {
			app.SetRoot(previous, true) // 元のビューに戻る
		}).
		AddButton("Cancel", func() {
			app.SetRoot(previous, true) // 元のビューに戻る
		})

	form.SetBorder(true).SetTitle("Form").SetTitleAlign(tview.AlignCenter)
	return form
}

// テーブルビューを作成
func createTable(app *tview.Application, previous tview.Primitive) tview.Primitive {
	table := tview.NewTable().
		SetBorders(true)
	lorem := strings.Split("Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.", " ")
	cols, rows := 10, 40
	word := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			color := tcell.ColorWhite
			if c < 1 || r < 1 {
				color = tcell.ColorYellow
			}
			table.SetCell(r, c,
				tview.NewTableCell(lorem[word]).
					SetTextColor(color).
					SetAlign(tview.AlignCenter))
			word = (word + 1) % len(lorem)
		}
	}

	// escape キーで元のビューに戻る
	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.SetRoot(previous, true)
		}
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		table.SetSelectable(false, false)
	})

	table.SetBorder(true).SetTitle("Table").SetTitleAlign(tview.AlignCenter)
	return table
}
