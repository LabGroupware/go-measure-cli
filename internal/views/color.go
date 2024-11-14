package views

import (
	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/gdamore/tcell/v2"
)

var (
	primaryColor   tcell.Color
	secondaryColor tcell.Color
	successColor   tcell.Color
	dangerColor    tcell.Color
	warningColor   tcell.Color
	infoColor      tcell.Color
	lightColor     tcell.Color
	darkColor      tcell.Color
	theme          string
	mainColor      tcell.Color
	subColor       tcell.Color
	contractColor  tcell.Color
)

func initColors(ctr app.Container) {
	primaryColor = tcell.GetColor(ctr.Config.View.Color.Primary)
	secondaryColor = tcell.GetColor(ctr.Config.View.Color.Secondary)
	successColor = tcell.GetColor(ctr.Config.View.Color.Success)
	dangerColor = tcell.GetColor(ctr.Config.View.Color.Danger)
	warningColor = tcell.GetColor(ctr.Config.View.Color.Warning)
	infoColor = tcell.GetColor(ctr.Config.View.Color.Info)
	lightColor = tcell.GetColor(ctr.Config.View.Color.Light)
	darkColor = tcell.GetColor(ctr.Config.View.Color.Dark)

	theme = ctr.Config.View.Theme
	switch theme {
	case "light":
		mainColor = lightColor
		subColor = tcell.ColorWhite
		contractColor = darkColor
	case "dark":
		mainColor = darkColor
		subColor = tcell.ColorBlack
		contractColor = lightColor
	default:
		mainColor = lightColor
		subColor = tcell.ColorWhite
		contractColor = darkColor
	}
}
