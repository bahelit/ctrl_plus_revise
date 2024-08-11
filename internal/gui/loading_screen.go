package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func ShowNotification(guiApp fyne.App, title, content string) {
	guiApp.SendNotification(&fyne.Notification{
		Title:   title,
		Content: content,
	})
}

func StartupScreen(guiApp fyne.App) fyne.Window {
	startupWindow := guiApp.NewWindow("Starting Control+Revise")
	infinite := widget.NewProgressBarInfinite()
	text := widget.NewLabel("Starting AI services in the background")
	startupWindow.SetContent(container.NewVBox(text, infinite))
	return startupWindow
}

func LoadingScreenWithMessage(guiApp fyne.App, title, msg string) fyne.Window {
	loadingScreen := guiApp.NewWindow(title)
	infinite := widget.NewProgressBarInfinite()
	text := widget.NewLabel(msg)
	loadingScreen.SetContent(container.NewVBox(text, infinite))
	return loadingScreen
}
func LoadingScreenWithProgressAndMessage(guiApp fyne.App, loading *widget.ProgressBar, status binding.String, title, msg string) fyne.Window {
	loadingScreen := guiApp.NewWindow(title)
	loadingScreen.Resize(fyne.NewSize(300, 80))
	text := widget.NewLabel(msg)
	s := widget.NewLabelWithData(status)
	layout := container.NewGridWithColumns(1, text, s, loading)
	loadingScreen.SetContent(layout)
	return loadingScreen
}
