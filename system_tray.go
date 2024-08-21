package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/layout"
	"github.com/ollama/ollama/api"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/bahelit/ctrl_plus_revise/internal/docker"
	"github.com/bahelit/ctrl_plus_revise/internal/gui"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
	"github.com/bahelit/ctrl_plus_revise/version"
)

const (
	ReplaceHighlightedText     = "ReplaceHighlightedText"
	SpeakAIResponseKey         = "SpeakAIResponseKey"
	ShowPopUpKey               = "ShowPopUpKey"
	ShowStartWindowKey         = "showStartWindow"
	firstRunKey                = "firstRun"
	CurrentPromptKey           = "lastPrompt"
	CurrentModelKey            = "lastModel"
	CurrentFromLangKey         = "fromLang"
	CurrentToLangKey           = "toLang"
	StopOllamaOnShutDownKey    = "stopOllamaOnShutDown"
	UseRemoteOllamaKey         = "useRemoteOllama"
	OllamaURLKey               = "OllamaURL"
	UseDockerKey               = "useDocker"
	AskAIKeyboardShortcut      = "AskAIKeyboardShortcut"
	CtrlReviseKeyboardShortcut = "CtrlReviseKeyboardShortcut"
	TranslateKeyboardShortcut  = "TranslateKeyboardShortcut"
)

var (
	aiActionDropdown *widget.Select
	aiModelDropdown  *widget.Select
)

const (
	AppTitle      = "Ctrl+Revise AI Text Generator"
	GreetingsText = "Welcome to Ctrl+Revise!"
	TrayMenuTitle = "Ctrl+Revise"
)

// SetupSysTray initializes the system tray for the application
func SetupSysTray(guiApp fyne.App) fyne.Window {
	if err := setBindingVariables(); err != nil {
		slog.Error("Failed to set binding variables", "error", err)
		os.Exit(1)
	}

	sysTray := guiApp.NewWindow(AppTitle)
	sysTray.SetTitle(AppTitle)

	setupTrayMenu(guiApp, sysTray)
	setupTrayWindowContent(guiApp, sysTray)

	sysTray.SetCloseIntercept(func() {
		sysTray.Hide()
	})
	return sysTray
}

// setupTrayMenu sets up the system tray menu
func setupTrayMenu(guiApp fyne.App, sysTray fyne.Window) {
	if desk, ok := guiApp.(desktop.App); ok {
		desk.SetSystemTrayMenu(fyne.NewMenu(TrayMenuTitle,
			fyne.NewMenuItem("Ask a Question", func() { askQuestion(guiApp) }),
			fyne.NewMenuItem("Translate Window", func() { translateText(guiApp) }),
			fyne.NewMenuItem("Keyboard Shortcuts", func() { showShortcuts(guiApp) }),
			fyne.NewMenuItem("Settings Window", func() { sysTray.Show() }),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("About", func() { showAbout(guiApp) }),
		))
	}
}

// setupTrayWindowContent sets up the content of the system tray window
func setupTrayWindowContent(guiApp fyne.App, sysTray fyne.Window) {
	welcomeText := mainWindowText()
	startUpCheckBox := showOnStartUpCheckBox(guiApp)
	stopOllamaOnShutdownCheckbox := stopOllamaOnShutdownCheckBox(guiApp)
	replaceHighlightedTextCheckBox := replaceHighlightedCheckbox(guiApp)
	speakAIResponseTextCheckBox := speakAIResponseCheckbox(guiApp)
	useDockerTextCheckBox := useDockerCheckBox(guiApp)
	showPopUpCheckBox := showPopUpCheckbox(guiApp)

	replaceHighlightedTextCheckBox.OnChanged = func(b bool) {
		if b {
			guiApp.Preferences().SetBool(ReplaceHighlightedText, true)
			guiApp.Preferences().SetBool(ShowPopUpKey, false)
			showPopUpCheckBox.Checked = false
			showPopUpCheckBox.Refresh()
			slog.Debug("Replace highlighted text is on")
		} else {
			slog.Debug("Replace highlighted text is off")
			guiApp.Preferences().SetBool(ReplaceHighlightedText, false)
		}
	}
	showPopUpCheckBox.OnChanged = func(b bool) {
		if b {
			guiApp.Preferences().SetBool(ShowPopUpKey, true)
			guiApp.Preferences().SetBool(ReplaceHighlightedText, false)
			replaceHighlightedTextCheckBox.Checked = false
			replaceHighlightedTextCheckBox.Refresh()
			slog.Debug("Show Pop-Up is on")
		} else {
			slog.Debug("Show Pop-Up is off")
			guiApp.Preferences().SetBool(ShowPopUpKey, false)
		}
	}

	keyboardShortcutsButton := widget.NewButton("Configure Keyboard Shortcuts", func() {
		showShortcuts(guiApp)
	})
	configureOllama := widget.NewButton("Configure Ollama", func() {
		installOrUpdateOllamaWindow(guiApp)
	})
	askQuestionsButton := widget.NewButton("Ask a Question", func() {
		askQuestion(guiApp)
	})
	translatorButton := widget.NewButton("Translate Text", func() {
		translateText(guiApp)
	})

	checkboxLayout := container.NewAdaptiveGrid(2,
		layout.Responsive(speakAIResponseTextCheckBox),
		layout.Responsive(replaceHighlightedTextCheckBox),
		layout.Responsive(useDockerTextCheckBox),
		layout.Responsive(showPopUpCheckBox),
		layout.Responsive(startUpCheckBox),
		layout.Responsive(stopOllamaOnShutdownCheckbox))

	mainWindow := container.NewVBox(
		welcomeText,
		askQuestionsButton,
		translatorButton,
		configureOllama,
		keyboardShortcutsButton,
		checkboxLayout,
	)
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
		toLangDropdown,
	)
	dropDownMenu := container.NewAdaptiveGrid(2,
		chooseActionLabel,
		aiActionDropdown,
		chooseModelLabel,
		aiModelDropdown,
		chooseLanguageLabel,
		langDivider,
	)
	sysTray.SetContent(container.NewBorder(mainWindow, dropDownMenu, nil, nil))
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
	welcomeText := widget.NewLabel(GreetingsText)
	welcomeText.Alignment = fyne.TextAlignCenter
	welcomeText.TextStyle = fyne.TextStyle{Bold: true}

	ctrlReviseKeys := getCtrlReviseKeys()
	var ctrlReviseString string
	keyLength := len(ctrlReviseKeys)
	for key, value := range ctrlReviseKeys {
		ctrlReviseString += strings.ToUpper(value)
		if key != keyLength-1 {
			ctrlReviseString += " + "
		}
	}
	shortcutText := widget.NewLabel("Pressing \"" + ctrlReviseString + "\" will send the highlighted text to an AI\nthe response is put into the clipboard")
	shortcutText.Alignment = fyne.TextAlignCenter
	shortcutText.TextStyle = fyne.TextStyle{Bold: true}
	closeMeText := widget.NewLabel("This window can be closed, the program will keep running in the taskbar")
	closeMeText.Alignment = fyne.TextAlignCenter
	return container.NewVBox(welcomeText, closeMeText, shortcutText)
}

func showOnStartUpCheckBox(guiApp fyne.App) *widget.Check {
	openStartWindow := guiApp.Preferences().BoolWithFallback(ShowStartWindowKey, true)
	startUpCheck := widget.NewCheck("Show this window on startup", func(b bool) {
		if !b {
			slog.Debug("Hiding start window")
			guiApp.Preferences().SetBool(ShowStartWindowKey, false)
		} else if b {
			guiApp.Preferences().SetBool(ShowStartWindowKey, true)
			slog.Debug("Showing start window")
		}
	})
	startUpCheck.Checked = openStartWindow
	return startUpCheck
}

func stopOllamaOnShutdownCheckBox(guiApp fyne.App) *widget.Check {
	stopOllamaOnShutdown := guiApp.Preferences().BoolWithFallback(StopOllamaOnShutDownKey, true)
	stopOllamaCheckbox := widget.NewCheck("Stop Ollama After Exiting", func(b bool) {
		if !b {
			slog.Debug("Leaving ollama running on shutdown")
			guiApp.Preferences().SetBool(StopOllamaOnShutDownKey, false)
		} else if b {
			guiApp.Preferences().SetBool(StopOllamaOnShutDownKey, true)
			slog.Debug("Stopping ollama on shutdown")
		}
	})
	stopOllamaCheckbox.Checked = stopOllamaOnShutdown
	return stopOllamaCheckbox
}

func replaceHighlightedCheckbox(guiApp fyne.App) *widget.Check {
	replaceText := guiApp.Preferences().BoolWithFallback(ReplaceHighlightedText, true)
	runOnCopy := widget.NewCheck("Paste AI Response", func(b bool) {
		if !b {
			slog.Debug("Replace highlighted checkbox is off")
			guiApp.Preferences().SetBool(ReplaceHighlightedText, false)
		} else if b {
			slog.Debug("Replace highlighted checkbox is on")
			guiApp.Preferences().SetBool(ReplaceHighlightedText, true)
		}
	})
	runOnCopy.Checked = replaceText
	return runOnCopy
}

func showPopUpCheckbox(guiApp fyne.App) *widget.Check {
	showPopUp := guiApp.Preferences().BoolWithFallback(ShowPopUpKey, false)
	popup := widget.NewCheck("Show Revise Window", func(b bool) {
		if !b {
			slog.Debug("Turning off PopUp mode")
			guiApp.Preferences().SetBool(ShowPopUpKey, false)
		} else if b {
			slog.Debug("Turning on PopUp mode")
			guiApp.Preferences().SetBool(ShowPopUpKey, true)
		}
	})
	popup.Checked = showPopUp
	return popup
}

func speakAIResponseCheckbox(guiApp fyne.App) *widget.Check {
	speakAIResponse := guiApp.Preferences().BoolWithFallback(SpeakAIResponseKey, false)
	speakAI := widget.NewCheck("Speak AI through speakers", func(b bool) {
		if !b {
			slog.Debug("Turning off Speech mode")
			guiApp.Preferences().SetBool(SpeakAIResponseKey, false)
		} else if b {
			slog.Debug("Turning on Speech mode")
			guiApp.Preferences().SetBool(SpeakAIResponseKey, true)
		}
	})
	speakAI.Checked = speakAIResponse
	return speakAI
}

func useDockerCheckBox(guiApp fyne.App) *widget.Check {
	userDocker := guiApp.Preferences().BoolWithFallback(UseDockerKey, false)
	userDockerCheck := widget.NewCheck("Run AI in Docker", func(b bool) {
		if !b {
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
		} else if b {
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

//nolint:gocyclo // it's a GUI function
func selectAIModelDropDown() *widget.Select {
	var (
		llama3Dot1      = "Llama 3.1 - RAM Usage: " + ollama.MemoryUsage[ollama.Llama3Dot1].String() + " (Default)"
		llama3          = "Llama 3 - RAM Usage: " + ollama.MemoryUsage[ollama.Llama3].String()
		codeLlama       = "CodeLlama - RAM Usage: " + ollama.MemoryUsage[ollama.CodeLlama].String()
		codeLlama13b    = "CodeLlama 13b - RAM Usage: " + ollama.MemoryUsage[ollama.CodeLlama13b].String()
		codeGemma       = "CodeGemma - RAM Usage: " + ollama.MemoryUsage[ollama.CodeGemma].String()
		deepSeekCoder   = "DeepSeekCoder. - RAM Usage: " + ollama.MemoryUsage[ollama.DeepSeekCoder].String()
		deepSeekCoderV2 = "DeepSeekCoderV2 - RAM Usage: " + ollama.MemoryUsage[ollama.DeepSeekCoderV2].String()
		gemma           = "Gemma - RAM Usage: " + ollama.MemoryUsage[ollama.Gemma].String()
		gemma2b         = "Gemma 2b - RAM Usage: " + ollama.MemoryUsage[ollama.Gemma2b].String()
		gemma2          = "Gemma2 - RAM Usage: " + ollama.MemoryUsage[ollama.Gemma2].String()
		gemma22B        = "Gemma2 2B - RAM Usage: " + ollama.MemoryUsage[ollama.Gemma22B].String()
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
		ollama.Gemma22B:        gemma22B,
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
			case gemma22B:
				selectedModel = ollama.Gemma22B
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
	model := guiApp.Preferences().IntWithFallback(CurrentModelKey, int(ollama.Llama3Dot1))
	selection := itemAndText[ollama.ModelName(model)]
	slog.Debug("Selected model", "model", selection)
	combo.SetSelected(selection)

	return combo
}

func PullModelWrapper(model ollama.ModelName, update bool) error {
	completed := binding.NewFloat()
	status := binding.NewString()
	loading := widget.NewProgressBarWithData(completed)
	startTime := time.Now()

	progressFunc := func(resp api.ProgressResponse) error {
		slog.Info("Progress", "status", resp.Status, "total", resp.Total, "completed", resp.Completed)
		err := status.Set(resp.Status)
		if err != nil {
			slog.Error("Failed to set status", "error", err)
		}
		loading.Max = float64(resp.Total)
		loading.Min = 0.0
		err = completed.Set(float64(resp.Completed))
		if err != nil {
			slog.Error("Failed to set progress", "error", err)
		}
		if resp.Total == resp.Completed {
			slog.Info("Model pulled", "model", model, "resp", resp)
		}
		return nil
	}

	pulling := gui.LoadingScreenWithProgressAndMessage(guiApp, loading, status, "Downloading Model", "Retrieving model: "+model.String())
	pulling.Show()
	defer func() {
		time.Sleep(1 * time.Second)
		pulling.Hide()
	}()

	err := ollama.PullModel(ollamaClient, model, progressFunc, update)
	if err != nil {
		slog.Error("Failed to pull model", "error", err)
		return err
	}
	elapsed := time.Since(startTime)
	if elapsed > 3*time.Second {
		gui.ShowNotification(guiApp, "Model Download Completed", "Model "+model.String()+" has been pulled")
		slog.Info("Model Download Completed", "model", model)
	} else {
		slog.Info("Already have the latest model", "model", model)
	}
	return nil
}

func selectTranslationFromDropDown() *widget.Select {
	combo := widget.NewSelect(
		gui.Languages,
		func(value string) {
			guiApp.Preferences().SetString(CurrentFromLangKey, value)
			err := translationFromBinding.Set(value)
			if err != nil {
				slog.Error("Failed to set translationFromBinding", "error", err)
			}
		})
	language := guiApp.Preferences().StringWithFallback(CurrentFromLangKey, string(ollama.English))
	combo.SetSelected(language)

	return combo
}

func selectTranslationToDropDown() *widget.Select {
	combo := widget.NewSelect(gui.Languages,
		func(value string) {
			guiApp.Preferences().SetString(CurrentToLangKey, value)
			err := translationToBinding.Set(value)
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
	value3 := widget.NewLabel("Coming Soon!")
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

//nolint:funlen // it's a GUI function
func showShortcuts(guiApp fyne.App) {
	slog.Debug("Showing Shortcuts")
	shortCuts := guiApp.NewWindow("Ctrl+Revise Keyboard Shortcuts")

	var grid *fyne.Container

	warn := widget.NewIcon(theme.WarningIcon())
	restartToReload := widget.NewLabelWithStyle("Restart application for changes to take effect", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	restartForChanges := container.NewGridWithRows(2, warn, restartToReload)

	askLabel := widget.NewLabel("\nAsk a Question with highlighted text")
	askModDropdown1 := keyboardModifierButtonsDropDown(AskQuestion, ModifierKey1)
	askModDropdown2 := keyboardModifierButtonsDropDown(AskQuestion, ModifierKey2)
	askKeyDropdown1 := keyboardModifierButtonsDropDown(AskQuestion, NormalKey)

	reviseLabel := widget.NewLabel("Selected Revise Action: ")
	label2Binding := widget.NewLabelWithData(selectedPromptBinding)
	hBox := container.NewGridWithColumns(2, reviseLabel, label2Binding)
	ctrlReviseModDropdown1 := keyboardModifierButtonsDropDown(CtrlRevise, ModifierKey1)
	ctrlReviseModDropdown2 := keyboardModifierButtonsDropDown(CtrlRevise, ModifierKey2)
	ctrlReviseKeyDropdown1 := keyboardModifierButtonsDropDown(CtrlRevise, NormalKey)

	cyclePromptLabel := widget.NewLabel("\nCycle through the prompt options")
	cyclePromptValue := widget.NewLabel("Alt + P")
	cyclePromptValue.TextStyle = fyne.TextStyle{Bold: true}

	readerLabel := widget.NewLabel("\nRead the highlighted text")
	readerValue := widget.NewLabel("Alt + R")
	readerValue.TextStyle = fyne.TextStyle{Bold: true}

	from, err := translationFromBinding.Get()
	if err != nil {
		slog.Error("Failed to get translationFromBinding", "error", err)
	}
	to, err := translationToBinding.Get()
	if err != nil {
		slog.Error("Failed to get translationToBinding", "error", err)
	}
	slog.Info("Translation languages", "from", from, "to", to)
	translateLabel := widget.NewLabel("\nTranslate the highlighted text, From: " + from + " To: " + to)
	translateKeyModDropdown1 := keyboardModifierButtonsDropDown(Translate, ModifierKey1)
	translateKeyModDropdown2 := keyboardModifierButtonsDropDown(Translate, ModifierKey2)
	translateKeyDropdown1 := keyboardModifierButtonsDropDown(Translate, NormalKey)

	askKeys := container.NewAdaptiveGrid(lengthOfKeyBoardShortcuts, askModDropdown1, askModDropdown2, askKeyDropdown1)
	ctrlReviseKeys := container.NewAdaptiveGrid(lengthOfKeyBoardShortcuts, ctrlReviseModDropdown1, ctrlReviseModDropdown2, ctrlReviseKeyDropdown1)
	translateKeys := container.NewAdaptiveGrid(lengthOfKeyBoardShortcuts, translateKeyModDropdown1, translateKeyModDropdown2, translateKeyDropdown1)

	grid = layout.NewResponsiveLayout(
		restartForChanges,
		hBox, ctrlReviseKeys,
		askLabel, askKeys,
		translateLabel, translateKeys,
		cyclePromptLabel, cyclePromptValue)

	speakAIResponse := guiApp.Preferences().BoolWithFallback(SpeakAIResponseKey, false)
	if speakAIResponse {
		grid.Add(readerLabel)
		grid.Add(readerValue)
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
	model := guiApp.Preferences().IntWithFallback(CurrentModelKey, int(ollama.Llama3Dot1))
	err := selectedModelBinding.Set(model)
	if err != nil {
		slog.Error("Failed to set selectedModelBinding", "error", err)
	}

	from := guiApp.Preferences().StringWithFallback(CurrentFromLangKey, string(ollama.English))
	err = translationFromBinding.Set(from)
	if err != nil {
		slog.Error("Failed to set selectedModelBinding", "error", err)
	}

	to := guiApp.Preferences().StringWithFallback(CurrentToLangKey, string(ollama.Spanish))
	err = translationToBinding.Set(to)
	if err != nil {
		slog.Error("Failed to set selectedModelBinding", "error", err)
	}

	prompt := guiApp.Preferences().StringWithFallback(CurrentPromptKey, ollama.CorrectGrammar.String())
	err = selectedPromptBinding.Set(prompt)
	if err != nil {
		slog.Error("Failed to set selectedPromptBinding", "error", err)
	}
	return err
}
