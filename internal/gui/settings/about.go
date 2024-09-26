package settings

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/layout"

	"github.com/bahelit/ctrl_plus_revise/version"
)

func ShowAbout(guiApp fyne.App) {
	slog.Debug("Showing about")
	about := guiApp.NewWindow("About Ctrl+Revise!")

	label1 := widget.NewLabel("Version")
	value1 := widget.NewLabel(version.Version)
	value1.TextStyle = fyne.TextStyle{Bold: true}
	label2 := widget.NewLabel("Author/Maintainer")
	value2 := widget.NewLabel("Michael Salmons")
	value2.TextStyle = fyne.TextStyle{Bold: true}
	label3 := widget.NewLabel("Contributors")
	value3 := widget.NewLabel("Coming Soon!")
	value3.TextStyle = fyne.TextStyle{Bold: true}
	grid := layout.NewResponsiveLayout(label1, value1, label2, value2, label3, value3)

	aboutTitle := widget.NewLabel("About Ctrl+Revise!")
	aboutTitle.Alignment = fyne.TextAlignCenter
	aboutTitle.TextStyle = fyne.TextStyle{Bold: true}
	aboutText := widget.NewLabel("Ctrl+Revise is here to help you unleash your inner wordsmith!\n" +
		"This nifty tool uses clever local AI agents to generate text based from any highlighted text.\n\n" +
		"Need some professional flair? Got a friendly tone in mind?\nOr maybe you just want to make sure your writing is grammatically correct?\n" +
		"Simply highlight the text you want to fix or ask about then press keyboard shortcut, and you're good to go!")

	aboutWindow := container.NewVBox(
		aboutTitle,
		aboutText,
		grid,
	)
	about.SetContent(aboutWindow)
	about.Show()
}
