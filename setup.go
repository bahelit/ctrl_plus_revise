package main

import (
	ollamaApi "github.com/ollama/ollama/api"
	"log/slog"

	"fyne.io/fyne/v2"

	"github.com/bahelit/ctrl_plus_revise/internal/config"
	"github.com/bahelit/ctrl_plus_revise/internal/docker"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/settings"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/shortcuts"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
)

func sayHello() {
	speakResponse := guiApp.Preferences().BoolWithFallback(config.SpeakAIResponseKey, false)
	if speakResponse {
		go func() {
			prompt := guiApp.Preferences().StringWithFallback(config.CurrentPromptKey, ollama.CorrectGrammar.String())
			_ = shortcuts.Speech.Speak("")
			speakErr := shortcuts.Speech.Speak("Control Plus Revise is set to: " + prompt)
			if speakErr != nil {
				slog.Error("Failed to speak", "error", speakErr)
			}
		}()
	}
}

func fetchModel(ollamaClient *ollamaApi.Client) {
	// Pull the model on startup, will pull updated model if available
	err := settings.PullModelWrapper(guiApp, ollamaClient, false)
	if err != nil {
		slog.Error("Failed to pull model", "error", err)
		guiApp.SendNotification(&fyne.Notification{
			Title: "Ollama Error",
			Content: "Failed to connect to pull model from Ollama\n" +
				"Check logs for more information\n" +
				"Ctrl+Revise will continue running, but may not function correctly",
		})
	}
}

func handleShutdown(p *int) {
	stopOllamaOnShutDown = guiApp.Preferences().BoolWithFallback(config.StopOllamaOnShutDownKey, true)
	useDocker := guiApp.Preferences().BoolWithFallback(config.UseDockerKey, false)
	if stopOllamaOnShutDown {
		if useDocker {
			docker.StopOllamaContainer()
		} else if p != nil {
			settings.StopOllama(nil)
		}
	} else {
		slog.Info("Leaving Ollama running")
	}
}

func setKeyboardShortcuts() {
	ask := guiApp.Preferences().StringListWithFallback(config.AskAIKeyboardShortcut, shortcuts.GetAskKeys())
	shortcuts.AskKey.ModifierKey1 = ask[0]
	if len(ask) == config.LengthOfKeyBoardShortcuts && ask[1] != shortcuts.EmptySelection {
		shortcuts.AskKey.ModifierKey2 = &ask[1]
		shortcuts.AskKey.Key = ask[2]
	} else {
		shortcuts.AskKey.Key = ask[1]
	}

	revise := guiApp.Preferences().StringListWithFallback(config.CtrlReviseKeyboardShortcut, shortcuts.GetCtrlReviseKeys())
	shortcuts.CtrlReviseKey.ModifierKey1 = revise[0]
	if len(revise) == config.LengthOfKeyBoardShortcuts && ask[1] != shortcuts.EmptySelection {
		shortcuts.CtrlReviseKey.ModifierKey2 = &revise[1]
		shortcuts.CtrlReviseKey.Key = revise[2]
	} else {
		shortcuts.CtrlReviseKey.Key = revise[1]
	}

	translate := guiApp.Preferences().StringListWithFallback(config.TranslateKeyboardShortcut, shortcuts.GetTranslateKeys())
	shortcuts.TranslateKey.ModifierKey1 = translate[0]
	if len(translate) == config.LengthOfKeyBoardShortcuts && ask[1] != shortcuts.EmptySelection {
		shortcuts.TranslateKey.ModifierKey2 = &translate[1]
		shortcuts.TranslateKey.Key = translate[2]
	} else {
		shortcuts.TranslateKey.Key = translate[1]
	}
}
