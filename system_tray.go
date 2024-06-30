package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/x/fyne/layout"
	//"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"log/slog"
)

const (
	replaceHighlightedText = "replaceHighlightedText"
	showStartWindow        = "showStartWindow"
	stopOllamaOnShutDown   = "stopOllamaOnShutDown"
)

func setupSysTray(guiApp fyne.App) fyne.Window {
	sysTray := guiApp.NewWindow("Ctrl+Revise AI Text Generator")
	sysTray.SetTitle("Ctrl+Revise AI Text Generator")

	combo := &widget.Select{}
	// System tray menu
	if desk, ok := guiApp.(desktop.App); ok {
		m := fyne.NewMenu("Ctrl+Revise",
			fyne.NewMenuItem("Ask a Question", func() {
				askQuestion(guiApp)
			}),
			fyne.NewMenuItem("Settings Window", func() {
				combo.SetSelected(selectedPrompt.String())
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
	replaceHighlightedTextCheckBox := replaceHighlightedCheckbox(guiApp)
	stopOllamaCheckBox := stopOllamaOnShutdownCheckBox(guiApp)
	welcomeText := mainWindowText()
	hideWindowButton := widget.NewButton("Hide This Window", func() {
		sysTray.Hide()
	})
	keyboardShortcutsButton := widget.NewButton("Show Keyboard Shortcuts", func() {
		showShortcuts(guiApp)
	})
	askQuestionsButton := widget.NewButton("Ask a Question", func() {
		askQuestion(guiApp)
	})
	mainWindow := container.NewVBox(
		welcomeText,
		askQuestionsButton,
		keyboardShortcutsButton,
		hideWindowButton,
		replaceHighlightedTextCheckBox,
		startUpCheckBox,
		stopOllamaCheckBox)

	chooseActionLabel := widget.NewLabel("Choose what the AI should do to the highlighted text:")
	chooseActionLabel.Alignment = fyne.TextAlignCenter
	combo = defaultCopyActionDropDown()

	dropDownMenu := container.NewVBox(
		chooseActionLabel,
		combo)
	sysTray.SetContent(container.NewVBox(mainWindow, dropDownMenu))

	sysTray.SetCloseIntercept(func() {
		sysTray.Hide()
	})

	return sysTray
}

func mainWindowText() *fyne.Container {
	welcomeText := widget.NewLabel("Welcome to Ctrl+Revise!")
	welcomeText.Alignment = fyne.TextAlignCenter
	welcomeText.TextStyle = fyne.TextStyle{Bold: true}
	shortcutText := widget.NewLabel("Pressing \"Ctrl + Shift + C\" will replace the highlighted text with an AI generated a response.")
	shortcutText.Alignment = fyne.TextAlignCenter
	shortcutText.TextStyle = fyne.TextStyle{Bold: true}
	closeMeText := widget.NewLabel("This window can be closed, the program will keep running in the taskbar")
	closeMeText.Alignment = fyne.TextAlignCenter
	return container.NewVBox(welcomeText, closeMeText, shortcutText)
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

func replaceHighlightedCheckbox(guiApp fyne.App) *widget.Check {
	replaceText := guiApp.Preferences().BoolWithFallback(replaceHighlightedText, true)
	runOnCopy := widget.NewCheck("Replace highlighted text with AI response", func(b bool) {
		if b == false {
			slog.Info("Replace highlighted checkbox is off")
			guiApp.Preferences().SetBool(replaceHighlightedText, false)
		} else if b == true {
			slog.Info("Replace highlighted checkbox is on")
			guiApp.Preferences().SetBool(replaceHighlightedText, true)
		}
	})
	runOnCopy.Checked = replaceText
	return runOnCopy
}

func stopOllamaOnShutdownCheckBox(guiApp fyne.App) *widget.Check {
	stopOllama := guiApp.Preferences().BoolWithFallback(stopOllamaOnShutDown, false)
	startUpCheck := widget.NewCheck("Stop AI agent on Program Exit", func(b bool) {
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
	// FIXME - Not updating when the selection is changed with keyboard shortcut
	combo := widget.NewSelect([]string{CorrectGrammar.String(), MakeItProfessional.String(), MakeItFriendly.String(),
		MakeHeadline.String(), MakeASummary.String(), MakeExpanded.String(), MakeExplanation.String(), MakeItAList.String()},
		func(value string) {
			switch value {
			case CorrectGrammar.String():
				selectedPrompt = CorrectGrammar
			case MakeItProfessional.String():
				selectedPrompt = MakeItProfessional
			case MakeItFriendly.String():
				selectedPrompt = MakeItFriendly
			case MakeHeadline.String():
				selectedPrompt = MakeHeadline
			case MakeASummary.String():
				selectedPrompt = MakeASummary
			case MakeExpanded.String():
				selectedPrompt = MakeExpanded
			case MakeExplanation.String():
				selectedPrompt = MakeExplanation
			case MakeItAList.String():
				selectedPrompt = MakeItAList
			default:
				slog.Error("Invalid selection", "value", value)
				selectedPrompt = CorrectGrammar
			}
			err := selectedPromptBinding.Set(selectedPrompt.String())
			if err != nil {
				slog.Error("Failed to set selectedPromptBinding", "error", err)
			}
		})
	combo.SetSelected(selectedPrompt.String())

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
	grid := layout.NewResponsiveLayout(label1, value1, label2, value2, label3, value3)

	aboutTitle := widget.NewLabel("About Ctrl+Revise!")
	aboutTitle.Alignment = fyne.TextAlignCenter
	aboutTitle.TextStyle = fyne.TextStyle{Bold: true}
	aboutText := widget.NewLabel("Ctrl+Revise is here to help you unleash your inner wordsmith!\n" +
		"This nifty tool uses clever local AI agents to generate text based on what you copy and paste.\n\n" +
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

func showShortcuts(guiApp fyne.App) {
	slog.Info("Showing Shortcuts")
	shortCuts := guiApp.NewWindow("Ctrl+Revise Shortcuts")

	label1 := widget.NewLabel("Ask a Question with highlighted text")
	value1 := widget.NewLabel("Alt + Ctrl + A")
	value1.TextStyle = fyne.TextStyle{Bold: true}
	label2 := widget.NewLabel("Replace highlighted text with: ")
	label2Binding := widget.NewLabelWithData(selectedPromptBinding)
	value2 := widget.NewLabel("Ctrl + Shift + C")
	value2.TextStyle = fyne.TextStyle{Bold: true}
	hbox := container.NewHBox(label2, label2Binding)
	grid := layout.NewResponsiveLayout(label1, value1, hbox, value2)
	shortCuts.SetContent(grid)
	shortCuts.Show()
}

func askQuestion(guiApp fyne.App) {
	slog.Info("Asking Question")
	var (
		screenHeight float32 = 180.0
		screenWidth  float32 = 480.0
	)
	question := guiApp.NewWindow("Ctrl+Revise Questions")
	question.Resize(fyne.NewSize(screenWidth, screenHeight))

	label1 := widget.NewLabel("Ask a Question")
	label1.TextStyle = fyne.TextStyle{Bold: true}
	label2 := widget.NewLabel("Press Shift + Enter to submit your question.")
	label2.TextStyle = fyne.TextStyle{Italic: true}

	text := widget.NewMultiLineEntry()
	text.PlaceHolder = "Ask your question here, remember this is an AI and important\n" +
		"questions should be verified with other sources."
	text.OnSubmitted = func(s string) {
		slog.Info("Question submitted - keyboard shortcut", "text", s)
		err := text.Validate()
		if err != nil {
			slog.Error("Error validating question", "error", err)
			return
		}
		response, err := askAI(ollamaClient, s)
		if err != nil {
			slog.Error("Failed to ask AI", "error", err)
		}
		questionPopUp(guiApp, s, &response)
		question.Close()
	}
	text.Validator = func(s string) error {
		if len(s) < 10 {
			return fmt.Errorf("question too short")
		}
		if len(s) > 10000000 {
			return fmt.Errorf("question too long, testing is needed before increasing the max length")
		}
		return nil
	}

	submitQuestionsButton := widget.NewButton("Submit Question", func() {
		slog.Info("Question submitted", "text", text.Text)
		err := text.Validate()
		if err != nil {
			slog.Error("Error validating question", "error", err)
			return
		}
		response, err := askAI(ollamaClient, text.Text)
		if err != nil {
			slog.Error("Failed to ask AI", "error", err)
		}
		questionPopUp(guiApp, text.Text, &response)
		question.Close()
	})

	topText := container.NewHBox(label1, label2)
	questionLayout := layout.NewResponsiveLayout(
		layout.Responsive(topText), // all sizes to 100%
		layout.Responsive(text, 1.0, 1.0))
	buttonLayout := layout.NewResponsiveLayout(layout.Responsive(submitQuestionsButton))
	questionWindow := container.NewVBox(
		questionLayout,
		buttonLayout,
	)
	question.SetContent(questionWindow)
	question.Show()
}

func changedPromptNotification() {
	guiApp.SendNotification(&fyne.Notification{
		Title:   "AI Action Changed",
		Content: "AI Action has been changed to:\n" + selectedPrompt.String(),
	})
}
