package settings

import (
	"github.com/bahelit/ctrl_plus_revise/internal/gui/loading"
	"log/slog"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/layout"
	ollamaApi "github.com/ollama/ollama/api"

	"github.com/bahelit/ctrl_plus_revise/internal/config"
	"github.com/bahelit/ctrl_plus_revise/internal/docker"
	"github.com/bahelit/ctrl_plus_revise/internal/gui"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/bindings"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/shortcuts"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
)

var (
	ollamaPID int
)

func ShowSettings(guiApp fyne.App, ollamaClient *ollamaApi.Client) {
	slog.Debug("Showing settings")
	settingsWindow := guiApp.NewWindow("Ctrl+Revise Settings")

	startUpCheckBox := showOnStartUpCheckBox(guiApp)
	stopOllamaOnShutdownCheckbox := stopOllamaOnShutdownCheckBox(guiApp)
	replaceHighlightedTextCheckBox := replaceHighlightedCheckbox(guiApp)
	speakAIResponseTextCheckBox := speakAIResponseCheckbox(guiApp)
	useDockerTextCheckBox := useDockerCheckBox(guiApp, ollamaClient)
	showPopUpCheckBox := showPopUpCheckbox(guiApp)

	replaceHighlightedTextCheckBox.OnChanged = func(b bool) {
		if b {
			guiApp.Preferences().SetBool(config.ReplaceHighlightedText, true)
			//guiApp.Preferences().SetBool(config.ShowPopUpKey, false)
			showPopUpCheckBox.Checked = false
			showPopUpCheckBox.Refresh()
			slog.Debug("Replace highlighted text is on")
		} else {
			slog.Debug("Replace highlighted text is off")
			guiApp.Preferences().SetBool(config.ReplaceHighlightedText, false)
		}
	}
	showPopUpCheckBox.OnChanged = func(b bool) {
		if b {
			//guiApp.Preferences().SetBool(config.ShowPopUpKey, true)
			guiApp.Preferences().SetBool(config.ReplaceHighlightedText, false)
			replaceHighlightedTextCheckBox.Checked = false
			replaceHighlightedTextCheckBox.Refresh()
			slog.Debug("Show Pop-Up is on")
		} else {
			slog.Debug("Show Pop-Up is off")
			guiApp.Preferences().SetBool(config.ShowPopUpKey, false)
		}
	}

	checkboxLayout := container.NewAdaptiveGrid(2,
		layout.Responsive(speakAIResponseTextCheckBox),
		layout.Responsive(replaceHighlightedTextCheckBox),
		layout.Responsive(useDockerTextCheckBox),
		layout.Responsive(showPopUpCheckBox),
		layout.Responsive(startUpCheckBox),
		layout.Responsive(stopOllamaOnShutdownCheckbox),
	)

	keyboardShortcutsButton := widget.NewButton("Configure Keyboard Shortcuts", func() {
		shortcuts.ShowShortcuts(guiApp)
	})
	configureOllama := widget.NewButton("Configure Ollama", func() {
		ollama.InstallOrUpdateOllamaWindow(guiApp, ollamaClient)
	})
	downloadModel := widget.NewButton("Download/Update Model", func() {
		_ = PullModelWrapper(guiApp, ollamaClient, true)
	})

	buttons := container.NewVBox(
		keyboardShortcutsButton,
		configureOllama,
		downloadModel,
	)

	chooseActionLabel := widget.NewLabel("Choose what the AI should do to the highlighted text:")
	chooseActionLabel.Alignment = fyne.TextAlignTrailing
	bindings.AiActionDropdown = selectCopyActionDropDown(guiApp)
	chooseModelLabel := widget.NewLabel("Choose which AI should respond to the highlighted text:")
	chooseModelLabel.Alignment = fyne.TextAlignTrailing
	bindings.AiModelDropdown = SelectAIModelDropDown(guiApp)
	bindings.AiModelDropdown.OnChanged = func(s string) {
		modelSelected := ollama.StringToModel(s)
		guiApp.Preferences().SetInt(config.CurrentModelKey, int(modelSelected))
	}
	chooseLanguageLabel := widget.NewLabel("Choose the languages for translation")
	chooseLanguageLabel.Alignment = fyne.TextAlignTrailing
	fromLangDropdown := SelectTranslationFromDropDown(guiApp)
	toLangDropdown := SelectTranslationToDropDown(guiApp)
	langDivider := container.NewHBox(
		widget.NewLabel("From: "),
		fromLangDropdown,
		widget.NewLabel("To: "),
		toLangDropdown,
	)
	dropDownMenu := container.NewAdaptiveGrid(2,
		chooseActionLabel,
		bindings.AiActionDropdown,
		chooseModelLabel,
		bindings.AiModelDropdown,
		chooseLanguageLabel,
		langDivider,
	)

	settingsWindow.SetContent(container.NewBorder(buttons, dropDownMenu, nil, nil, checkboxLayout))
	settingsWindow.Show()
}

func PullModelWrapper(guiApp fyne.App, ollamaClient *ollamaApi.Client, update bool) error {
	completed := binding.NewFloat()
	status := binding.NewString()
	progressBar := widget.NewProgressBarWithData(completed)
	startTime := time.Now()

	progressFunc := func(resp ollamaApi.ProgressResponse) error {
		slog.Info("Progress", "status", resp.Status, "total", resp.Total, "completed", resp.Completed)
		err := status.Set(resp.Status)
		if err != nil {
			slog.Error("Failed to set status", "error", err)
		}
		progressBar.Max = float64(resp.Total)
		progressBar.Min = 0.0
		err = completed.Set(float64(resp.Completed))
		if err != nil {
			slog.Error("Failed to set progress", "error", err)
		}
		if resp.Total == resp.Completed {
			slog.Info("Model pulled", "resp", resp)
		}
		return nil
	}

	model := ollama.GetActiveModel(guiApp)
	pulling := loading.LoadingScreenWithProgressAndMessage(guiApp, progressBar, status, "Downloading Model", "Retrieving model: "+model.String())
	pulling.Show()
	defer func() {
		time.Sleep(1 * time.Second)
		pulling.Hide()
	}()

	err := ollama.PullModel(guiApp, ollamaClient, progressFunc, update)
	if err != nil {
		slog.Error("Failed to pull model", "error", err)
		return err
	}
	elapsed := time.Since(startTime)
	if elapsed > 3*time.Second {
		loading.ShowNotification(guiApp, "Model Download Completed", "Model "+model.String()+" has been pulled")
		slog.Info("Model Download Completed", "model", model)
	} else {
		slog.Info("Already have the latest model", "model", model)
	}
	return nil
}

func SelectTranslationFromDropDown(guiApp fyne.App) *widget.Select {
	combo := widget.NewSelect(
		gui.Languages,
		func(value string) {
			guiApp.Preferences().SetString(config.CurrentFromLangKey, value)
			err := bindings.TranslationFromBinding.Set(value)
			if err != nil {
				slog.Error("Failed to set TranslationFromBinding", "error", err)
			}
		})
	language := guiApp.Preferences().StringWithFallback(config.CurrentFromLangKey, string(ollama.English))
	combo.SetSelected(language)

	return combo
}

func SelectTranslationToDropDown(guiApp fyne.App) *widget.Select {
	combo := widget.NewSelect(gui.Languages,
		func(value string) {
			guiApp.Preferences().SetString(config.CurrentToLangKey, value)
			err := bindings.TranslationToBinding.Set(value)
			if err != nil {
				slog.Error("Failed to set TranslationToBinding", "error", err)
			}
		})
	language := guiApp.Preferences().StringWithFallback(config.CurrentToLangKey, string(ollama.Spanish))
	combo.SetSelected(language)

	return combo
}

func showOnStartUpCheckBox(guiApp fyne.App) *widget.Check {
	openStartWindow := guiApp.Preferences().BoolWithFallback(config.ShowStartWindowKey, true)
	startUpCheck := widget.NewCheck("Show this window on startup", func(b bool) {
		if !b {
			slog.Debug("Hiding start window")
			guiApp.Preferences().SetBool(config.ShowStartWindowKey, false)
		} else if b {
			guiApp.Preferences().SetBool(config.ShowStartWindowKey, true)
			slog.Debug("Showing start window")
		}
	})
	startUpCheck.Checked = openStartWindow
	return startUpCheck
}

func stopOllamaOnShutdownCheckBox(guiApp fyne.App) *widget.Check {
	stopOllamaOnShutdown := guiApp.Preferences().BoolWithFallback(config.StopOllamaOnShutDownKey, true)
	stopOllamaCheckbox := widget.NewCheck("Stop Ollama After Exiting", func(b bool) {
		if !b {
			slog.Debug("Leaving ollama running on shutdown")
			guiApp.Preferences().SetBool(config.StopOllamaOnShutDownKey, false)
		} else if b {
			guiApp.Preferences().SetBool(config.StopOllamaOnShutDownKey, true)
			slog.Debug("Stopping ollama on shutdown")
		}
	})
	stopOllamaCheckbox.Checked = stopOllamaOnShutdown
	return stopOllamaCheckbox
}

func replaceHighlightedCheckbox(guiApp fyne.App) *widget.Check {
	replaceText := guiApp.Preferences().BoolWithFallback(config.ReplaceHighlightedText, true)
	runOnCopy := widget.NewCheck("Paste AI Response", func(b bool) {
		if !b {
			slog.Debug("Replace highlighted checkbox is off")
			guiApp.Preferences().SetBool(config.ReplaceHighlightedText, false)
		} else if b {
			slog.Debug("Replace highlighted checkbox is on")
			guiApp.Preferences().SetBool(config.ReplaceHighlightedText, true)
		}
	})
	runOnCopy.Checked = replaceText
	return runOnCopy
}

func showPopUpCheckbox(guiApp fyne.App) *widget.Check {
	//_ = guiApp.Preferences().BoolWithFallback(config.ShowPopUpKey, false)
	popup := widget.NewCheck("Show Revise Window", func(b bool) {
		if !b {
			slog.Debug("Turning off PopUp mode")
			//guiApp.Preferences().SetBool(config.ShowPopUpKey, false)
		} else if b {
			slog.Debug("Turning on PopUp mode")
			//guiApp.Preferences().SetBool(config.ShowPopUpKey, true)
		}
	})
	popup.Disabled()
	popup.Checked = false
	return popup
}

func speakAIResponseCheckbox(guiApp fyne.App) *widget.Check {
	speakAIResponse := guiApp.Preferences().BoolWithFallback(config.SpeakAIResponseKey, false)
	speakAI := widget.NewCheck("Speak AI through speakers", func(b bool) {
		if !b {
			slog.Debug("Turning off Speech mode")
			guiApp.Preferences().SetBool(config.SpeakAIResponseKey, false)
		} else if b {
			slog.Debug("Turning on Speech mode")
			guiApp.Preferences().SetBool(config.SpeakAIResponseKey, true)
		}
	})
	speakAI.Checked = speakAIResponse
	return speakAI
}

func useDockerCheckBox(guiApp fyne.App, ollamaClient *ollamaApi.Client) *widget.Check {
	userDocker := guiApp.Preferences().BoolWithFallback(config.UseDockerKey, false)
	userDockerCheck := widget.NewCheck("Run AI in Docker", func(b bool) {
		if !b {
			slog.Debug("Not using Docker")
			guiApp.Preferences().SetBool(config.UseDockerKey, false)
			docker.StopOllamaContainer()
			gotConnected := SetupServices(guiApp, ollamaClient)
			if gotConnected == nil {
				go func() {
					gotConnected = SetupServices(guiApp, ollamaClient)
					if gotConnected == nil {
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
			guiApp.Preferences().SetBool(config.UseDockerKey, true)
			go func() {
				StopOllama(nil)
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
					guiApp.Preferences().SetBool(config.UseDockerKey, false)
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

func SetupServices(guiApp fyne.App, ollamaClient *ollamaApi.Client) *ollamaApi.Client {
	connectedToOllama := ollama.CheckOllamaConnection(guiApp, ollamaClient, nil)
	if connectedToOllama != nil {
		return connectedToOllama
	}
	// Ollama isn't running, should we start it with Docker?
	useDocker := guiApp.Preferences().BoolWithFallback(config.UseDockerKey, false)
	if useDocker {
		slog.Info("Starting Ollama container")
		return ollama.CheckOllamaConnection(guiApp, ollamaClient, nil)
	}
	slog.Info("Starting Ollama")
	// Start Ollama locally
	if StartOllama(ollamaClient) {
		return ollama.CheckOllamaConnection(guiApp, ollamaClient, nil)
	} else {
		slog.Error("Failed to connect Ollama")
	}
	return nil
}

func StartOllama(ollamaClient *ollamaApi.Client) (connectedToOllama bool) {
	_, err := exec.LookPath("ollama")
	if err != nil {
		slog.Info("Ollama not found", "error", err)
		ollama.InstallOrUpdateOllama()
	}

	ollamaServe := exec.Command("ollama", "serve")
	err = ollamaServe.Start()
	if err != nil {
		return connectedToOllama
	}
	time.Sleep(1 * time.Second)
	go func() {
		err = ollamaServe.Wait()
		if err != nil {
			slog.Error("Ollama process exited", "error", err)
		}
	}()

	ollamaPID = ollamaServe.Process.Pid
	slog.Info("Started Ollama", "pid", ollamaPID)

	ollamaClient = ollama.ConnectToOllama()
	if ollamaClient != nil {
		connectedToOllama = true
	}

	return connectedToOllama
}

func StopOllama(p *int) {
	if p == nil || *p == 0 {
		slog.Error("Ollama PID not found")
		return
	}
	ollamaProcess, err := os.FindProcess(*p)
	if err != nil {
		slog.Error("Failed to find process", "error", err)
	}

	err = ollamaProcess.Kill()
	if err != nil {
		slog.Error("Failed to find process", "error", err)
	}
}

func selectCopyActionDropDown(guiApp fyne.App) *widget.Select {
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
			var selectedPrompt ollama.PromptMsg
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
			guiApp.Preferences().SetString(config.CurrentPromptKey, selectedPrompt.String())
			err := bindings.SelectedPromptBinding.Set(selectedPrompt.String())
			if err != nil {
				slog.Error("Failed to set SelectedPromptBinding", "error", err)
			}
		})
	prompt := guiApp.Preferences().StringWithFallback(config.CurrentPromptKey, ollama.CorrectGrammar.String())
	combo.SetSelected(prompt)

	return combo
}

//nolint:gocyclo // it's a GUI function
func SelectAIModelDropDown(guiApp fyne.App) *widget.Select {
	var (
		llama3Dot2      = "Llama 3.2 - RAM Usage: " + ollama.MemoryUsage[ollama.Llama3Dot2].String() + " (Default)"
		llama3Dot21B    = "Llama 3.2 1B - RAM Usage: " + ollama.MemoryUsage[ollama.Llama3Dot21B].String()
		llama3Dot1      = "Llama 3.1 - RAM Usage: " + ollama.MemoryUsage[ollama.Llama3Dot1].String()
		llamaVision     = "Llama 3.2 Vision - RAM Usage: " + ollama.MemoryUsage[ollama.LlamaVision].String()
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
		mistralNemo     = "Mistral Nemo - RAM Usage: " + ollama.MemoryUsage[ollama.MistralNemo].String()
		nemoTronMini    = "Nemotron Mini - RAM Usage: " + ollama.MemoryUsage[ollama.MistralNemo].String()
		phi3            = "Phi3 - RAM Usage: " + ollama.MemoryUsage[ollama.Phi3].String()
	)
	var itemAndText = map[ollama.ModelName]string{
		ollama.Llama3Dot2:      llama3Dot2,
		ollama.Llama3Dot21B:    llama3Dot21B,
		ollama.Llama3Dot1:      llama3Dot1,
		ollama.LlamaVision:     llamaVision,
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
		ollama.MistralNemo:     mistralNemo,
		ollama.NemoTronMini:    nemoTronMini,
		ollama.Phi3:            phi3,
	}
	combo := widget.NewSelect([]string{
		llama3Dot1,
		llama3Dot2,
		llama3Dot21B,
		llama3,
		llamaVision,
		codeLlama,
		codeLlama13b,
		codeGemma,
		deepSeekCoder,
		deepSeekCoderV2,
		gemma,
		gemma2b,
		gemma2,
		mistral,
		mistralNemo,
		nemoTronMini,
		phi3},
		func(value string) {
			var selectedModel ollama.ModelName
			switch value {
			case llama3Dot1:
				selectedModel = ollama.Llama3Dot1
			case llama3Dot2:
				selectedModel = ollama.Llama3Dot2
			case llama3Dot21B:
				selectedModel = ollama.Llama3Dot21B
			case llamaVision:
				selectedModel = ollama.LlamaVision
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
			case mistralNemo:
				selectedModel = ollama.MistralNemo
			case nemoTronMini:
				selectedModel = ollama.NemoTronMini
			case phi3:
				selectedModel = ollama.Phi3
			default:
				slog.Error("Invalid selection", "value", value)
				selectedModel = ollama.Llama3
			}
			err := bindings.SelectedModelBinding.Set(int(selectedModel))
			if err != nil {
				slog.Error("Failed to set SelectedModelBinding", "error", err)
			}
			slog.Debug("Selected model", "model", selectedModel)
		})
	model := guiApp.Preferences().IntWithFallback(config.CurrentModelKey, int(ollama.Llama3Dot2))
	selection := itemAndText[ollama.ModelName(model)]
	slog.Debug("Selected model", "model", selection)
	combo.SetSelected(selection)

	return combo
}
