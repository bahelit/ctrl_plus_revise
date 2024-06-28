package main

import (
	"crypto/sha256"
	"log/slog"

	hook "github.com/robotn/gohook"
	"golang.design/x/clipboard"
)

func registerHotkeys() {
	hook.Register(hook.KeyDown, []string{"f", "ctrl", "shift"}, func(e hook.Event) {
		slog.Info("ctrl-shift-r has been pressed", "event", e)
		handleHotKeyPressed(MakeItFriendly)
	})
	hook.Register(hook.KeyDown, []string{"g", "ctrl", "shift"}, func(e hook.Event) {
		slog.Info("ctrl-shift-r has been pressed", "event", e)
		handleHotKeyPressed(CorrectGrammar)
	})
	hook.Register(hook.KeyDown, []string{"r", "ctrl", "shift"}, func(e hook.Event) {
		slog.Info("ctrl-shift-r has been pressed", "event", e)
		handleHotKeyPressed(MakeItProfessional)
	})

	s := hook.Start()
	<-hook.Process(s)
}

func handleHotKeyPressed(prompt PromptMsg) {
	clippy := clipboard.Read(clipboard.FmtText)
	if clippy == nil {
		slog.Info("Clipboard is empty")
		return
	}

	generated, err := generateResponseFromOllama(ollamaClient, prompt, string(clippy))
	if err != nil {
		// TODO: Implement error handling, tell user to restart ollama, maybe we can restart ollama here?
		slog.Error("Failed to communicate with Ollama", "error", err)
		return
	}
	clippyPopUp(guiApp, string(clippy), &generated)
	lastRspHash = sha256.New().Sum([]byte(generated.Response))
}
