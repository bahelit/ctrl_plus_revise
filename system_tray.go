package main

import (
	"fmt"
	"log/slog"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/layout"

	"github.com/bahelit/ctrl_plus_revise/version"
)

const (
	replaceHighlightedText  = "replaceHighlightedText"
	speakAIResponseKey      = "speakAIResponseKey"
	showStartWindowKey      = "showStartWindow"
	firstRunKey             = "firstRun"
	currentPromptKey        = "lastPrompt"
	currentModelKey         = "lastModel"
	stopOllamaOnShutDownKey = "stopOllamaOnShutDown"
	useDockerKey            = "useDocker"
)

var (
	aiActionDropdown *widget.Select
	aiModelDropdown  *widget.Select
)

func setupSysTray(guiApp fyne.App) fyne.Window {
	err := setBindingVariables()
	if err != nil {
		slog.Error("Failed to set binding variables", "error", err)
		os.Exit(1)
	}

	sysTray := guiApp.NewWindow("Ctrl+Revise AI Text Generator")
	sysTray.SetTitle("Ctrl+Revise AI Text Generator")

	// System tray menu
	if desk, ok := guiApp.(desktop.App); ok {
		m := fyne.NewMenu("Ctrl+Revise",
			fyne.NewMenuItem("Ask a Question", func() {
				askQuestion(guiApp)
			}),
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
	replaceHighlightedTextCheckBox := replaceHighlightedCheckbox(guiApp)
	speakAIResponseTextCheckBox := speakAIResponseCheckbox(guiApp)
	useDockerTextCheckBox := useDockerCheckBox(guiApp)
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

	checkboxLayout := container.NewAdaptiveGrid(2,
		layout.Responsive(speakAIResponseTextCheckBox),
		layout.Responsive(replaceHighlightedTextCheckBox),
		layout.Responsive(useDockerTextCheckBox),
		layout.Responsive(startUpCheckBox))

	mainWindow := container.NewVBox(
		welcomeText,
		askQuestionsButton,
		keyboardShortcutsButton,
		hideWindowButton,
		checkboxLayout)

	chooseActionLabel := widget.NewLabel("Choose what the AI should do to the highlighted text:")
	chooseActionLabel.Alignment = fyne.TextAlignTrailing
	aiActionDropdown = selectCopyActionDropDown()

	chooseModelLabel := widget.NewLabel("Choose which AI should respond to the highlighted text:")
	chooseModelLabel.Alignment = fyne.TextAlignTrailing
	aiModelDropdown = selectAIModelDropDown()

	dropDownMenu := container.NewAdaptiveGrid(2,
		chooseActionLabel,
		aiActionDropdown,
		chooseModelLabel,
		aiModelDropdown)
	sysTray.SetContent(container.NewVBox(mainWindow, dropDownMenu))

	sysTray.SetCloseIntercept(func() {
		sysTray.Hide()
	})
	return sysTray
}

func loadIcon() {
	var (
		icon         fyne.Resource
		errLocation1 error
		errLocation2 error
	)
	icon, errLocation1 = fyne.LoadResourceFromPath("/app/share/icons/hicolor/256x256/apps/com.bahelit.ctrl_plus_revise.png")
	if errLocation1 != nil {
		icon, errLocation2 = fyne.LoadResourceFromPath("images/icon.png")
		if errLocation2 != nil {
			slog.Warn("Failed to load icon", "error", errLocation1)
			slog.Warn("Failed to load icon", "error", errLocation2)
		}
	}
	guiApp.SetIcon(icon)
}

func mainWindowText() *fyne.Container {
	welcomeText := widget.NewLabel("Welcome to Ctrl+Revise!")
	welcomeText.Alignment = fyne.TextAlignCenter
	welcomeText.TextStyle = fyne.TextStyle{Bold: true}
	shortcutText := widget.NewLabel("Pressing \"Alt + Ctrl + C\" will replace the highlighted text with an AI generated a response.")
	shortcutText.Alignment = fyne.TextAlignCenter
	shortcutText.TextStyle = fyne.TextStyle{Bold: true}
	closeMeText := widget.NewLabel("This window can be closed, the program will keep running in the taskbar")
	closeMeText.Alignment = fyne.TextAlignCenter
	return container.NewVBox(welcomeText, closeMeText, shortcutText)
}

func showOnStartUpCheckBox(guiApp fyne.App) *widget.Check {
	openStartWindow := guiApp.Preferences().BoolWithFallback(showStartWindowKey, true)
	startUpCheck := widget.NewCheck("Show this window on startup", func(b bool) {
		if b == false {
			slog.Debug("Hiding start window")
			guiApp.Preferences().SetBool(showStartWindowKey, false)
		} else if b == true {
			guiApp.Preferences().SetBool(showStartWindowKey, true)
			slog.Debug("Showing start window")
		}
	})
	startUpCheck.Checked = openStartWindow
	return startUpCheck
}

func replaceHighlightedCheckbox(guiApp fyne.App) *widget.Check {
	replaceText := guiApp.Preferences().BoolWithFallback(replaceHighlightedText, true)
	speakAIResponse := guiApp.Preferences().BoolWithFallback(speakAIResponseKey, false)
	runOnCopy := widget.NewCheck("Replace highlighted text with AI response", func(b bool) {
		if b == false {
			slog.Debug("Replace highlighted checkbox is off")
			guiApp.Preferences().SetBool(replaceHighlightedText, false)
			if speakAIResponse {
				go func() {
					speakErr := speech.Speak("Highlighted text will be appended with an AI response.")
					if speakErr != nil {
						slog.Error("Failed to speak", "error", speakErr)
					}
				}()
			}
		} else if b == true {
			slog.Debug("Replace highlighted checkbox is on")
			guiApp.Preferences().SetBool(replaceHighlightedText, true)
			if speakAIResponse {
				go func() {
					speakErr := speech.Speak("Highlighted text will be replaced with AI response.")
					if speakErr != nil {
						slog.Error("Failed to speak", "error", speakErr)
					}
				}()
			}
		}
	})
	runOnCopy.Checked = replaceText
	return runOnCopy
}

func speakAIResponseCheckbox(guiApp fyne.App) *widget.Check {
	speakAIResponse := guiApp.Preferences().BoolWithFallback(speakAIResponseKey, false)
	speakAI := widget.NewCheck("Speak AI through speakers", func(b bool) {
		if b == false {
			slog.Debug("Turning off Speech mode")
			guiApp.Preferences().SetBool(speakAIResponseKey, false)
			go func() {
				speakErr := speech.Speak("Turning off Speech mode.")
				if speakErr != nil {
					slog.Error("Failed to speak", "error", speakErr)
				}
			}()
		} else if b == true {
			slog.Debug("Turning on Speech mode")
			guiApp.Preferences().SetBool(speakAIResponseKey, true)
			go func() {
				speakErr := speech.Speak("Turning on Speech mode.")
				if speakErr != nil {
					slog.Error("Failed to speak", "error", speakErr)
				}
			}()
		}
	})
	speakAI.Checked = speakAIResponse
	return speakAI
}

func useDockerCheckBox(guiApp fyne.App) *widget.Check {
	userDocker := guiApp.Preferences().BoolWithFallback(useDockerKey, false)
	userDockerCheck := widget.NewCheck("Run AI in Docker", func(b bool) {
		if b == false {
			slog.Debug("Not using Docker")
			guiApp.Preferences().SetBool(useDockerKey, false)
			stopOllamaContainer(dockerClient)
			gotConnected := setupServices()
			if !gotConnected {
				go func() {
					gotConnected := setupServices()
					if !gotConnected {
						slog.Error("Failed to connect to or start Ollama container")
						guiApp.SendNotification(&fyne.Notification{
							Title: "Docker Error",
							Content: "Failed to connect to start Ollama container\n" +
								"Check logs for more information\n" +
								"Ctrl+Revise will shutdown now",
						})
						os.Exit(1)
					} else {
						guiApp.SendNotification(&fyne.Notification{
							Title:   "Connected Docker",
							Content: "Ready to process requests with Docker!",
						})
					}
				}()
			}
		} else if b == true {
			slog.Debug("Using Docker")
			guiApp.Preferences().SetBool(useDockerKey, true)
			go func() {
				stopOllama(ollamaPID)
				gotConnected := setupServices()
				if !gotConnected {
					slog.Error("Failed to connect to Docker or start Ollama container")
					guiApp.SendNotification(&fyne.Notification{
						Title: "Docker Error",
						Content: "Failed to connect to Docker or start Ollama container\n" +
							"Please check your Docker installation and try again\n" +
							"Check logs for more information\n" +
							"Ctrl+Revise will continue to run without Docker",
					})
					guiApp.Preferences().SetBool(useDockerKey, false)
					// TODO restart ollama without docker
				} else {
					guiApp.SendNotification(&fyne.Notification{
						Title:   "Connected Docker",
						Content: "Ready to process requests with Docker!",
					})
				}
			}()
		}
	})
	userDockerCheck.Checked = userDocker
	return userDockerCheck
}

func selectCopyActionDropDown() *widget.Select {
	// FIXME - Not updating when the selection is changed with keyboard shortcut
	combo := widget.NewSelect([]string{
		CorrectGrammar.String(),
		MakeItProfessional.String(),
		MakeItFriendly.String(),
		MakeHeadline.String(),
		MakeASummary.String(),
		MakeExpanded.String(),
		MakeExplanation.String(),
		MakeItAList.String()},
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
			guiApp.Preferences().SetString(currentPromptKey, selectedPrompt.String())
			err := selectedPromptBinding.Set(selectedPrompt.String())
			if err != nil {
				slog.Error("Failed to set selectedPromptBinding", "error", err)
			}
		})
	prompt := guiApp.Preferences().StringWithFallback(currentPromptKey, CorrectGrammar.String())
	combo.SetSelected(prompt)

	return combo
}
func selectAIModelDropDown() *widget.Select {
	// FIXME - Not updating when the selection is changed with keyboard shortcut
	combo := widget.NewSelect([]string{
		BashBot.String(),
		CodeLlama.String(),
		CodeLlama13b.String(),
		Gemma.String(),
		Llama3.String()},
		func(value string) {
			switch value {
			case BashBot.String():
				selectedModel = BashBot
			case CodeLlama.String():
				selectedModel = CodeLlama
			case CodeLlama13b.String():
				selectedModel = CodeLlama13b
			case Gemma.String():
				selectedModel = Gemma
			case Llama3.String():
				selectedModel = Llama3
			default:
				slog.Error("Invalid selection", "value", value)
				selectedModel = Llama3
			}
			guiApp.Preferences().SetInt(currentModelKey, int(selectedModel))
			err := selectedModelBinding.Set(int(selectedModel))
			if err != nil {
				slog.Error("Failed to set selectedModelBinding", "error", err)
			}
		})
	model := guiApp.Preferences().IntWithFallback(currentModelKey, int(Llama3))
	combo.SetSelected(ModelName(model).String())

	return combo
}

func updateDropDownMenus() {
	aiActionDropdown.SetSelected(selectedPrompt.String())
	aiModelDropdown.SetSelected(selectedModel.String())
}

func showAbout(guiApp fyne.App) {
	slog.Debug("Showing about")
	about := guiApp.NewWindow("About Ctrl+Revise!")

	label1 := widget.NewLabel("Version")
	value1 := widget.NewLabel(version.Version)
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
	slog.Debug("Showing Shortcuts")
	shortCuts := guiApp.NewWindow("Ctrl+Revise Shortcuts")

	var grid *fyne.Container

	label1 := widget.NewLabel("Ask a Question with highlighted text")
	value1 := widget.NewLabel("Alt + Ctrl + A")
	value1.TextStyle = fyne.TextStyle{Bold: true}

	label2 := widget.NewLabel("Replace highlighted text with: ")
	label2Binding := widget.NewLabelWithData(selectedPromptBinding)
	value2 := widget.NewLabel("Alt + Ctrl + C")
	value2.TextStyle = fyne.TextStyle{Bold: true}
	hbox := container.NewHBox(label2, label2Binding)

	label3 := widget.NewLabel("Read the highlighted text")
	value3 := widget.NewLabel("Alt + Ctrl + R")
	value3.TextStyle = fyne.TextStyle{Bold: true}
	speakAIResponse := guiApp.Preferences().BoolWithFallback(speakAIResponseKey, false)
	if speakAIResponse {
		grid = layout.NewResponsiveLayout(label1, value1, hbox, value2, label3, value3)
	} else {
		grid = layout.NewResponsiveLayout(label1, value1, hbox, value2)
	}
	shortCuts.SetContent(grid)
	shortCuts.Show()
}

func askQuestion(guiApp fyne.App) {
	slog.Debug("Asking Question")
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
		slog.Debug("Question submitted - keyboard shortcut", "text", s)
		err := text.Validate()
		if err != nil {
			slog.Error("Error validating question", "error", err)
			return
		}
		response, err := askAI(ollamaClient, selectedModel, s)
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
		slog.Debug("Question submitted", "text", text.Text)
		err := text.Validate()
		if err != nil {
			slog.Error("Error validating question", "error", err)
			return
		}
		response, err := askAI(ollamaClient, selectedModel, text.Text)
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
	guiApp.Preferences().SetString(currentPromptKey, selectedPrompt.String())
	guiApp.SendNotification(&fyne.Notification{
		Title:   "AI Action Changed",
		Content: "AI Action has been changed to:\n" + selectedPrompt.String(),
	})
}

func setBindingVariables() error {
	selectedModelBinding = binding.NewInt()
	model := guiApp.Preferences().IntWithFallback(currentModelKey, int(Llama3))
	err := selectedModelBinding.Set(model)
	if err != nil {
		slog.Error("Failed to set selectedModelBinding", "error", err)
	}

	selectedPromptBinding = binding.NewString()
	prompt := guiApp.Preferences().StringWithFallback(currentPromptKey, CorrectGrammar.String())
	err = selectedPromptBinding.Set(prompt)
	if err != nil {
		slog.Error("Failed to set selectedPromptBinding", "error", err)
	}
	return err
}

func showNotification(title, content string) {
	guiApp.SendNotification(&fyne.Notification{
		Title:   title,
		Content: content,
	})
}

func startupScreen() fyne.Window {
	startupWindow := guiApp.NewWindow("Starting Control+Revise")
	infinite := widget.NewProgressBarInfinite()
	text := widget.NewLabel("Starting AI services in the background")
	startupWindow.SetContent(container.NewVBox(text, infinite))
	return startupWindow
}

func loadingScreenWithMessage(title, msg string) fyne.Window {
	loadingScreen := guiApp.NewWindow(title)
	infinite := widget.NewProgressBarInfinite()
	text := widget.NewLabel(msg)
	loadingScreen.SetContent(container.NewVBox(text, infinite))
	return loadingScreen
}
