package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/layout"
	"github.com/ollama/ollama/api"

	"github.com/bahelit/ctrl_plus_revise/internal/docker"
	"github.com/bahelit/ctrl_plus_revise/internal/gui"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
	"github.com/bahelit/ctrl_plus_revise/version"
)

const (
	ReplaceHighlightedText  = "ReplaceHighlightedText"
	SpeakAIResponseKey      = "SpeakAIResponseKey"
	ShowStartWindowKey      = "showStartWindow"
	firstRunKey             = "firstRun"
	CurrentPromptKey        = "lastPrompt"
	CurrentModelKey         = "lastModel"
	CurrentFromLangKey      = "fromLang"
	CurrentToLangKey        = "toLang"
	StopOllamaOnShutDownKey = "stopOllamaOnShutDown"
	UseDockerKey            = "useDocker"
)

var (
	aiActionDropdown *widget.Select
	aiModelDropdown  *widget.Select
)

func SetupSysTray(guiApp fyne.App) fyne.Window {
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

	chooseLanguageLabel := widget.NewLabel("Choose the languages for translation")
	chooseLanguageLabel.Alignment = fyne.TextAlignTrailing
	fromLangDropdown := selectTranslationFromDropDown()
	toLangDropdown := selectTranslationToDropDown()

	langDivider := container.NewHBox(
		widget.NewLabel("From: "),
		fromLangDropdown,
		widget.NewLabel("To: "),
		toLangDropdown)

	dropDownMenu := container.NewAdaptiveGrid(2,
		chooseActionLabel,
		aiActionDropdown,
		chooseModelLabel,
		aiModelDropdown,
		chooseLanguageLabel,
		langDivider)
	sysTray.SetContent(container.NewVBox(mainWindow, dropDownMenu))

	sysTray.SetCloseIntercept(func() {
		sysTray.Hide()
	})
	return sysTray
}

func LoadIcon(guiApp fyne.App) {
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
	shortcutText := widget.NewLabel("Pressing \"Alt + C\" will send the highlighted text to an AI\nthe response is put into the clipboard")
	shortcutText.Alignment = fyne.TextAlignCenter
	shortcutText.TextStyle = fyne.TextStyle{Bold: true}
	closeMeText := widget.NewLabel("This window can be closed, the program will keep running in the taskbar")
	closeMeText.Alignment = fyne.TextAlignCenter
	return container.NewVBox(welcomeText, closeMeText, shortcutText)
}

func showOnStartUpCheckBox(guiApp fyne.App) *widget.Check {
	openStartWindow := guiApp.Preferences().BoolWithFallback(ShowStartWindowKey, true)
	startUpCheck := widget.NewCheck("Show this window on startup", func(b bool) {
		if b == false {
			slog.Debug("Hiding start window")
			guiApp.Preferences().SetBool(ShowStartWindowKey, false)
		} else if b == true {
			guiApp.Preferences().SetBool(ShowStartWindowKey, true)
			slog.Debug("Showing start window")
		}
	})
	startUpCheck.Checked = openStartWindow
	return startUpCheck
}

func replaceHighlightedCheckbox(guiApp fyne.App) *widget.Check {
	replaceText := guiApp.Preferences().BoolWithFallback(ReplaceHighlightedText, true)
	speakAIResponse := guiApp.Preferences().BoolWithFallback(SpeakAIResponseKey, false)
	runOnCopy := widget.NewCheck("Auto paste text with AI response", func(b bool) {
		if b == false {
			slog.Debug("Replace highlighted checkbox is off")
			guiApp.Preferences().SetBool(ReplaceHighlightedText, false)
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
			guiApp.Preferences().SetBool(ReplaceHighlightedText, true)
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
	speakAIResponse := guiApp.Preferences().BoolWithFallback(SpeakAIResponseKey, false)
	speakAI := widget.NewCheck("Speak AI through speakers", func(b bool) {
		if b == false {
			slog.Debug("Turning off Speech mode")
			guiApp.Preferences().SetBool(SpeakAIResponseKey, false)
			go func() {
				speakErr := speech.Speak("Turning off Speech mode.")
				if speakErr != nil {
					slog.Error("Failed to speak", "error", speakErr)
				}
			}()
		} else if b == true {
			slog.Debug("Turning on Speech mode")
			guiApp.Preferences().SetBool(SpeakAIResponseKey, true)
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
	userDocker := guiApp.Preferences().BoolWithFallback(UseDockerKey, false)
	userDockerCheck := widget.NewCheck("Run AI in Docker", func(b bool) {
		if b == false {
			slog.Debug("Not using Docker")
			guiApp.Preferences().SetBool(UseDockerKey, false)
			docker.StopOllamaContainer()
			gotConnected := setupServices()
			if !gotConnected {
				go func() {
					gotConnected = setupServices()
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
			guiApp.Preferences().SetBool(UseDockerKey, true)
			go func() {
				stopOllama(ollamaPID)
				gotConnected := docker.SetupDocker()
				if !gotConnected {
					slog.Error("Failed to connect to Docker or start Ollama container")
					guiApp.SendNotification(&fyne.Notification{
						Title: "Docker Error",
						Content: "Failed to connect to Docker or start Ollama container\n" +
							"Please check your Docker installation and try again\n" +
							"Check logs for more information\n" +
							"Ctrl+Revise will continue to run without Docker",
					})
					guiApp.Preferences().SetBool(UseDockerKey, false)
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
	combo := widget.NewSelect([]string{
		ollama.CorrectGrammar.String(),
		ollama.MakeItProfessional.String(),
		ollama.MakeItFriendly.String(),
		ollama.MakeHeadline.String(),
		ollama.MakeASummary.String(),
		ollama.MakeExpanded.String(),
		ollama.MakeExplanation.String(),
		ollama.MakeItAList.String()},
		func(value string) {
			switch value {
			case ollama.CorrectGrammar.String():
				selectedPrompt = ollama.CorrectGrammar
			case ollama.MakeItProfessional.String():
				selectedPrompt = ollama.MakeItProfessional
			case ollama.MakeItFriendly.String():
				selectedPrompt = ollama.MakeItFriendly
			case ollama.MakeHeadline.String():
				selectedPrompt = ollama.MakeHeadline
			case ollama.MakeASummary.String():
				selectedPrompt = ollama.MakeASummary
			case ollama.MakeExpanded.String():
				selectedPrompt = ollama.MakeExpanded
			case ollama.MakeExplanation.String():
				selectedPrompt = ollama.MakeExplanation
			case ollama.MakeItAList.String():
				selectedPrompt = ollama.MakeItAList
			default:
				slog.Error("Invalid selection", "value", value)
				selectedPrompt = ollama.CorrectGrammar
			}
			guiApp.Preferences().SetString(CurrentPromptKey, selectedPrompt.String())
			err := selectedPromptBinding.Set(selectedPrompt.String())
			if err != nil {
				slog.Error("Failed to set selectedPromptBinding", "error", err)
			}
		})
	prompt := guiApp.Preferences().StringWithFallback(CurrentPromptKey, ollama.CorrectGrammar.String())
	combo.SetSelected(prompt)

	return combo
}

func selectAIModelDropDown() *widget.Select {
	var (
		llama3Dot1      = "Llama3.1 - RAM Usage: " + ollama.MemoryUsage[ollama.Llama3Dot1].String() + " (Default)"
		llama3          = "Llama3 - RAM Usage: " + ollama.MemoryUsage[ollama.Llama3].String() + " (Default)"
		codeLlama       = "CodeLlama - RAM Usage: " + ollama.MemoryUsage[ollama.CodeLlama].String()
		codeLlama13b    = "CodeLlama13b - RAM Usage: " + ollama.MemoryUsage[ollama.CodeLlama13b].String()
		codeGemma       = "CodeGemma - RAM Usage: " + ollama.MemoryUsage[ollama.CodeGemma].String()
		deepSeekCoder   = "DeepSeekCoder. - RAM Usage: " + ollama.MemoryUsage[ollama.DeepSeekCoder].String()
		deepSeekCoderV2 = "DeepSeekCoderV2 - RAM Usage: " + ollama.MemoryUsage[ollama.DeepSeekCoderV2].String()
		gemma           = "Gemma - RAM Usage: " + ollama.MemoryUsage[ollama.Gemma].String()
		gemma2b         = "Gemma2b - RAM Usage: " + ollama.MemoryUsage[ollama.Gemma2b].String()
		gemma2          = "Gemma2 - RAM Usage: " + ollama.MemoryUsage[ollama.Gemma2].String()
		mistral         = "Mistral - RAM Usage: " + ollama.MemoryUsage[ollama.Mistral].String()
		phi3            = "Phi3 - RAM Usage: " + ollama.MemoryUsage[ollama.Phi3].String()
	)
	var itemAndText = map[ollama.ModelName]string{
		ollama.Llama3Dot1:      llama3Dot1,
		ollama.Llama3:          llama3,
		ollama.CodeLlama:       codeLlama,
		ollama.CodeLlama13b:    codeLlama13b,
		ollama.CodeGemma:       codeGemma,
		ollama.DeepSeekCoder:   deepSeekCoder,
		ollama.DeepSeekCoderV2: deepSeekCoderV2,
		ollama.Gemma:           gemma,
		ollama.Gemma2b:         gemma2b,
		ollama.Gemma2:          gemma2,
		ollama.Mistral:         mistral,
		ollama.Phi3:            phi3,
	}
	combo := widget.NewSelect([]string{
		llama3Dot1,
		llama3,
		codeLlama,
		codeLlama13b,
		codeGemma,
		deepSeekCoder,
		deepSeekCoderV2,
		gemma,
		gemma2b,
		gemma2,
		mistral,
		phi3},
		func(value string) {
			switch value {
			case llama3Dot1:
				selectedModel = ollama.Llama3Dot1
			case llama3:
				selectedModel = ollama.Llama3
			case codeLlama:
				selectedModel = ollama.CodeLlama
			case codeLlama13b:
				selectedModel = ollama.CodeLlama13b
			case codeGemma:
				selectedModel = ollama.CodeGemma
			case deepSeekCoder:
				selectedModel = ollama.DeepSeekCoder
			case deepSeekCoderV2:
				selectedModel = ollama.DeepSeekCoderV2
			case gemma:
				selectedModel = ollama.Gemma
			case gemma2b:
				selectedModel = ollama.Gemma2b
			case gemma2:
				selectedModel = ollama.Gemma2
			case mistral:
				selectedModel = ollama.Mistral
			case phi3:
				selectedModel = ollama.Phi3
			default:
				slog.Error("Invalid selection", "value", value)
				selectedModel = ollama.Llama3
			}
			guiApp.Preferences().SetInt(CurrentModelKey, int(selectedModel))
			err := selectedModelBinding.Set(int(selectedModel))
			if err != nil {
				slog.Error("Failed to set selectedModelBinding", "error", err)
			}
			slog.Debug("Selected model", "model", selectedModel)
			if ollamaClient != nil {
				_ = PullModelWrapper(selectedModel, false)
			}
		})
	model := guiApp.Preferences().IntWithFallback(CurrentModelKey, int(ollama.Llama3))
	selection := itemAndText[ollama.ModelName(model)]
	slog.Debug("Selected model", "model", selection)
	combo.SetSelected(selection)

	return combo
}

func PullModelWrapper(model ollama.ModelName, update bool) error {
	pulling := gui.LoadingScreenWithMessage(guiApp, "Downloading Model", "fetching model: "+model.String())
	pulling.Show()
	defer func() {
		time.Sleep(500 * time.Millisecond)
		pulling.Hide()
	}()

	progressFunc := func(resp api.ProgressResponse) error {
		slog.Info("Progress", "status", resp.Status, "total", resp.Total, "completed", resp.Completed)
		if resp.Total == resp.Completed {
			slog.Info("Model pulled", "model", model, "resp", resp)
		}
		gui.ShowNotification(guiApp, "Model Download Completed", "Model "+model.String()+" has been pulled")
		return nil
	}

	err := ollama.PullModel(ollamaClient, model, progressFunc, update)
	if err != nil {
		slog.Error("Failed to pull model", "error", err)
		return err
	}
	return nil
}

func selectTranslationFromDropDown() *widget.Select {
	combo := widget.NewSelect([]string{
		string(ollama.English),
		string(ollama.Arabic),
		string(ollama.Chinese),
		string(ollama.French),
		string(ollama.German),
		string(ollama.Italian),
		string(ollama.Japanese),
		string(ollama.Portuguese),
		string(ollama.Russian),
		string(ollama.Spanish)},
		func(value string) {
			guiApp.Preferences().SetString(CurrentFromLangKey, value)
			err := translationFromBinding.Set(selectedModel.String())
			if err != nil {
				slog.Error("Failed to set translationFromBinding", "error", err)
			}
		})
	language := guiApp.Preferences().StringWithFallback(CurrentFromLangKey, string(ollama.English))
	combo.SetSelected(language)

	return combo
}
func selectTranslationToDropDown() *widget.Select {
	combo := widget.NewSelect([]string{
		string(ollama.Spanish),
		string(ollama.Arabic),
		string(ollama.Chinese),
		string(ollama.French),
		string(ollama.German),
		string(ollama.Italian),
		string(ollama.Japanese),
		string(ollama.Portuguese),
		string(ollama.Russian),
		string(ollama.English)},
		func(value string) {
			guiApp.Preferences().SetString(CurrentToLangKey, value)
			err := translationToBinding.Set(selectedModel.String())
			if err != nil {
				slog.Error("Failed to set translationToBinding", "error", err)
			}
		})
	language := guiApp.Preferences().StringWithFallback(CurrentToLangKey, string(ollama.Spanish))
	combo.SetSelected(language)

	return combo
}

func UpdateDropDownMenus() {
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
	value1 := widget.NewLabel("Alt + A")
	value1.TextStyle = fyne.TextStyle{Bold: true}

	label2 := widget.NewLabel("Replace highlighted text with: ")
	label2Binding := widget.NewLabelWithData(selectedPromptBinding)
	value2 := widget.NewLabel("Alt + C")
	value2.TextStyle = fyne.TextStyle{Bold: true}
	hbox := container.NewHBox(label2, label2Binding)

	label3 := widget.NewLabel("Cycle through the prompt options")
	value3 := widget.NewLabel("Alt + P")
	value3.TextStyle = fyne.TextStyle{Bold: true}

	label4 := widget.NewLabel("Read the highlighted text")
	value4 := widget.NewLabel("Alt + R")
	value4.TextStyle = fyne.TextStyle{Bold: true}

	from, err := translationFromBinding.Get()
	if err != nil {
		slog.Error("Failed to get translationFromBinding", "error", err)
	}
	to, err := translationToBinding.Get()
	if err != nil {
		slog.Error("Failed to get translationToBinding", "error", err)
	}
	label5 := widget.NewLabel("Translate the highlighted text, From: " + from + " To: " + to)
	value5 := widget.NewLabel("Alt + T")
	value5.TextStyle = fyne.TextStyle{Bold: true}
	speakAIResponse := guiApp.Preferences().BoolWithFallback(SpeakAIResponseKey, false)
	if speakAIResponse {
		grid = layout.NewResponsiveLayout(label1, value1, hbox, value2, label3, value3, label4, value4, label5, value5)
	} else {
		grid = layout.NewResponsiveLayout(label1, value1, hbox, value2, label4, value4, label5, value5)
	}
	shortCuts.SetContent(grid)
	shortCuts.Show()
}

func askQuestion(guiApp fyne.App) {
	slog.Debug("Asking Question")
	var (
		screenHeight float32 = 200.0
		screenWidth  float32 = 480.0
	)
	question := guiApp.NewWindow("Ctrl+Revise Questions")
	question.Resize(fyne.NewSize(screenWidth, screenHeight))

	label1 := widget.NewLabel("Ask a Question")
	label1.TextStyle = fyne.TextStyle{Bold: true}
	label2 := widget.NewLabel("Press Shift + Enter to submit your question.")
	label2.TextStyle = fyne.TextStyle{Italic: true}

	text := widget.NewMultiLineEntry()
	text.SetMinRowsVisible(4)
	text.PlaceHolder = "Ask your question here, remember this is an AI and important\n" +
		"questions should be verified with other sources."
	text.OnSubmitted = func(s string) {
		slog.Debug("Question submitted - keyboard shortcut", "text", s)
		err := text.Validate()
		if err != nil {
			slog.Error("Error validating question", "error", err)
			return
		}
		loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
			"Asking question with model: "+selectedModel.String()+"...")
		loadingScreen.Show()
		response, err := ollama.AskAI(ollamaClient, selectedModel, s)
		if err != nil {
			slog.Error("Failed to ask AI", "error", err)
			loadingScreen.Hide()
			return
		}
		loadingScreen.Hide()
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
		loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
			"Asking question with model: "+selectedModel.String()+"...")
		loadingScreen.Show()
		response, err := ollama.AskAI(ollamaClient, selectedModel, text.Text)
		if err != nil {
			slog.Error("Failed to ask AI", "error", err)
			loadingScreen.Hide()
			return
		}
		loadingScreen.Hide()
		questionPopUp(guiApp, text.Text, &response)
		question.Close()
	})

	topText := container.NewHBox(label1, label2)
	questionLayout := layout.NewResponsiveLayout(
		layout.Responsive(topText))
	buttonLayout := layout.NewResponsiveLayout(layout.Responsive(submitQuestionsButton))
	questionWindow := container.NewBorder(questionLayout, buttonLayout, nil, nil, text)
	question.SetContent(container.NewVScroll(questionWindow))
	question.Show()
}

func ChangedPromptNotification() {
	guiApp.Preferences().SetString(CurrentPromptKey, selectedPrompt.String())
	guiApp.SendNotification(&fyne.Notification{
		Title:   "AI Action Changed",
		Content: "AI Action has been changed to:\n" + selectedPrompt.String(),
	})
}

func setBindingVariables() error {
	selectedModelBinding = binding.NewInt()
	model := guiApp.Preferences().IntWithFallback(CurrentModelKey, int(ollama.Llama3))
	err := selectedModelBinding.Set(model)
	if err != nil {
		slog.Error("Failed to set selectedModelBinding", "error", err)
	}

	translationFromBinding = binding.NewString()
	from := guiApp.Preferences().StringWithFallback(CurrentFromLangKey, string(ollama.English))
	err = translationFromBinding.Set(from)
	if err != nil {
		slog.Error("Failed to set selectedModelBinding", "error", err)
	}

	translationToBinding = binding.NewString()
	to := guiApp.Preferences().StringWithFallback(CurrentToLangKey, string(ollama.Spanish))
	err = translationToBinding.Set(to)
	if err != nil {
		slog.Error("Failed to set selectedModelBinding", "error", err)
	}

	selectedPromptBinding = binding.NewString()
	prompt := guiApp.Preferences().StringWithFallback(CurrentPromptKey, ollama.CorrectGrammar.String())
	err = selectedPromptBinding.Set(prompt)
	if err != nil {
		slog.Error("Failed to set selectedPromptBinding", "error", err)
	}
	return err
}
