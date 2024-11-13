// Package app provides the application logic for the TUI application.
package app

import (
	"fmt"

	"github.com/LabGroupware/go-measure-tui/internal/views"
	"github.com/rivo/tview"
)

// App represents the TUI application
type App struct {
	App   *tview.Application
	Pages *tview.Pages
}

// NewApp creates a new TUI application
func NewApp() *App {
	app := tview.NewApplication()
	pages := tview.NewPages()

	// メインメニュー画面を初期化
	menu := views.NewMenu(func(page string) {
		pages.SwitchToPage(page)
	})

	// 各画面を追加
	pages.AddPage("menu", menu, true, true)
	pages.AddPage("form", views.NewForm(), true, false)
	pages.AddPage("list", views.NewList(), true, false)

	return &App{
		App:   app,
		Pages: pages,
	}
}

// Run starts the TUI application
func (a *App) Run() error {
	if err := a.App.SetRoot(a.Pages, true).Run(); err != nil {
		return fmt.Errorf("failed to run application: %w", err)
	}

	return nil
}
