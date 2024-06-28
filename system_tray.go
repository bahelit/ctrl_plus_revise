package main

import (
	"fyne.io/fyne/v2/layout"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const (
	actionSelection      = "actionSelection"
	autoRunOnCopy        = "autoRunOnCopy"
	showStartWindow      = "showStartWindow"
	stopOllamaOnShutDown = "stopOllamaOnShutDown"
)

func setupSysTray(guiApp fyne.App) fyne.Window {
	sysTray := guiApp.NewWindow("Ctrl+Revise AI Text Generator")
	sysTray.SetTitle("Ctrl+Revise AI Text Generator")

	// System tray menu
	if desk, ok := guiApp.(desktop.App); ok {
		m := fyne.NewMenu("Ctrl+Revise",
			fyne.NewMenuItem("Settings Window", func() {
				sysTray.Show()
			}),
			fyne.NewMenuItem("Keyboard Shortcuts", func() {
				showShortcuts(guiApp)
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("About", func() {
				showAbout(guiApp)
			}),
		)
		desk.SetSystemTrayMenu(m)
	}

	// System tray window content
	startUpCheckBox := showOnStartUpCheckBox(guiApp)
	runOnCopyCheckBox := autoRunOnCopyCheckBox(guiApp)
	stopOllamaCheckBox := stopOllamaOnShutdownCheckBox(guiApp)
	infoText, welcomeText := mainWindowText()
	hideWindowButton := widget.NewButton("Hide This Window", func() {
		sysTray.Hide()
	})
	keyboardShortcutsButton := widget.NewButton("Show Keyboard Shortcuts", func() {
		showShortcuts(guiApp)
	})
	mainWindow := container.NewVBox(
		welcomeText,
		infoText,
		hideWindowButton,
		keyboardShortcutsButton,
		startUpCheckBox,
		runOnCopyCheckBox,
		stopOllamaCheckBox)

	combo := defaultCopyActionDropDown()
	dropDownMenu := container.NewVBox(
		widget.NewLabel("Choose the response type to use when the clipboard is being monitored:"),
		combo)
	sysTray.SetContent(container.NewVBox(mainWindow, dropDownMenu))

	sysTray.SetCloseIntercept(func() {
		sysTray.Hide()
	})

	return sysTray
}

func mainWindowText() (infoText, welcomeText *widget.Label) {
	welcomeText = widget.NewLabel("Welcome to Ctrl+Revise!")
	welcomeText.Alignment = fyne.TextAlignCenter
	welcomeText.TextStyle = fyne.TextStyle{Bold: true}
	infoText = widget.NewLabel("This window can be closed, the program will keep running in the taskbar")
	return infoText, welcomeText
}

func showOnStartUpCheckBox(guiApp fyne.App) *widget.Check {
	openStartWindow := guiApp.Preferences().BoolWithFallback(showStartWindow, true)
	startUpCheck := widget.NewCheck("Show this window on startup", func(b bool) {
		if b == false {
			slog.Info("Hiding start window")
			guiApp.Preferences().SetBool(showStartWindow, false)
		} else if b == true {
			guiApp.Preferences().SetBool(showStartWindow, true)
			slog.Info("Showing start window")
		}
	})
	startUpCheck.Checked = openStartWindow
	return startUpCheck
}

func autoRunOnCopyCheckBox(guiApp fyne.App) *widget.Check {
	autoGenerate := guiApp.Preferences().BoolWithFallback(autoRunOnCopy, false)
	runOnCopy := widget.NewCheck("Generate anytime text is copied (not saved on restart)", func(b bool) {
		if b == false {
			slog.Info("Stopping clipboard watcher checkbox is off")
			guiApp.Preferences().SetBool(autoRunOnCopy, false)
			close(clippyCheckbox)
		} else if b == true {
			slog.Info("Starting clipboard watcher checkbox is on")
			guiApp.Preferences().SetBool(autoRunOnCopy, true)
			clippyCheckbox = make(chan bool)
			go watchClipboardForChanges(clippyCheckbox)
		}
	})
	runOnCopy.Checked = autoGenerate
	return runOnCopy
}

func stopOllamaOnShutdownCheckBox(guiApp fyne.App) *widget.Check {
	stopOllama := guiApp.Preferences().BoolWithFallback(stopOllamaOnShutDown, false)
	startUpCheck := widget.NewCheck("Stop Ollama on Program Exit", func(b bool) {
		if b == false {
			slog.Info("Hiding start window")
			guiApp.Preferences().SetBool(stopOllamaOnShutDown, false)
		} else if b == true {
			guiApp.Preferences().SetBool(stopOllamaOnShutDown, true)
			slog.Info("Showing start window")
		}
	})
	startUpCheck.Checked = stopOllama
	return startUpCheck
}

func defaultCopyActionDropDown() *widget.Select {
	combo := widget.NewSelect([]string{CorrectGrammar.String(), MakeItProfessional.String(), MakeItFriendly.String()},
		func(value string) {
			switch value {
			case CorrectGrammar.String():
				selectedPrompt = CorrectGrammar
			case MakeItProfessional.String():
				selectedPrompt = MakeItProfessional
			case MakeItFriendly.String():
				selectedPrompt = MakeItFriendly
			default:
				slog.Error("Invalid selection", "value", value)
				selectedPrompt = CorrectGrammar
			}
		})
	combo.SetSelected(CorrectGrammar.String())
	return combo
}

func showAbout(guiApp fyne.App) {
	slog.Info("Showing about")
	about := guiApp.NewWindow("About Ctrl+Revise!")

	label1 := widget.NewLabel("Version")
	value1 := widget.NewLabel(Version)
	value1.TextStyle = fyne.TextStyle{Bold: true}
	label2 := widget.NewLabel("Author/Maintainer")
	value2 := widget.NewLabel("Michael Salmons")
	value2.TextStyle = fyne.TextStyle{Bold: true}
	label3 := widget.NewLabel("Contributors")
	value3 := widget.NewLabel("Your name could be here, Wink Wink.")
	value3.TextStyle = fyne.TextStyle{Bold: true}
	grid := container.New(layout.NewFormLayout(), label1, value1, label2, value2, label3, value3)

	aboutTitle := widget.NewLabel("About Ctrl+Revise!")
	aboutTitle.Alignment = fyne.TextAlignCenter
	aboutTitle.TextStyle = fyne.TextStyle{Bold: true}
	aboutText := widget.NewLabel("Ctrl+Revise is here to help you unleash your inner wordsmith!\nThis nifty tool uses clever local AI agents to generate text based on what you copy and paste.\n\nNeed some professional flair? Got a friendly tone in mind?\nOr maybe you just want to make sure your writing is grammatically correct?\nLlama Launcher's got you covered. Simply copy the generated text right onto your clipboard, and you're good to go!")

	aboutWindow := container.NewVBox(
		aboutTitle,
		aboutText,
		grid,
	)
	about.SetContent(aboutWindow)
	about.Show()
}

func showShortcuts(guiApp fyne.App) {
	slog.Info("Showing Shortcuts")
	shortCuts := guiApp.NewWindow("Ctrl+Revise Shortcuts")

	label1 := widget.NewLabel(CorrectGrammar.String())
	value1 := widget.NewLabel("Ctrl + Shift + G")
	value1.TextStyle = fyne.TextStyle{Bold: true}
	label2 := widget.NewLabel(MakeItProfessional.String())
	value2 := widget.NewLabel("Ctrl + Shift + R")
	value2.TextStyle = fyne.TextStyle{Bold: true}
	label3 := widget.NewLabel(MakeItFriendly.String())
	value3 := widget.NewLabel("Ctrl + Shift + F")
	value3.TextStyle = fyne.TextStyle{Bold: true}
	grid := container.New(layout.NewFormLayout(), label1, value1, label2, value2, label3, value3)
	shortCuts.SetContent(grid)
	shortCuts.Show()
}
