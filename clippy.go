package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ollama/ollama/api"
	"golang.design/x/clipboard"
)

func watchClipboardForChanges(keepListening chan bool) {
	slog.Info("Starting ctrl_plus_revise clipboard listener...")
	changed := clipboard.Watch(context.Background(), clipboard.FmtText)

	var (
		generated   api.GenerateResponse
		clippyHash  []byte
		hashesMatch bool
		err         error
	)
	for clippy := range changed {
		select {
		case <-keepListening:
			slog.Info("Stopping clippy watcher")
			return
		default:
			clippyHash = sha256.New().Sum(clippy)
			hashesMatch = bytes.Equal(lastRspHash, clippyHash)
			if string(clippy) == generated.Response || hashesMatch {
				slog.Info("Skipping duplicate Clipboard")
				continue
			}
			slog.Info("Clipboard data changed", "contents", string(clippy))
			generated, err = generateResponseFromOllama(ollamaClient, selectedPrompt, string(clippy))
			if err != nil {
				// TODO: Implement error handling, tell user to restart ollama, maybe we can restart ollama here?
				slog.Error("Failed to communicate with Ollama", "error", err)
				continue
			}
			lastRspHash = sha256.New().Sum([]byte(generated.Response))
			clippyPopUp(guiApp, string(clippy), &generated)
		}
	}
}

func clippyPopUp(a fyne.App, input string, generated *api.GenerateResponse) {
	w := a.NewWindow("llama listener")
	hello := widget.NewLabel("Glad To Help!")
	hello.TextStyle = fyne.TextStyle{Bold: true}
	hello.Alignment = fyne.TextAlignCenter
	w.SetContent(container.NewVBox(
		hello,
		widget.NewLabel("Original text:\n"+input),
		widget.NewLabel("AI Generated text:\n"+generated.Response),
		widget.NewButton("Copy generated text to Clipboard", func() {
			lastRspHash = sha256.New().Sum([]byte(generated.Response))
			w.Clipboard().SetContent(generated.Response)
			w.Hide()
		}),
		widget.NewButton("Try Again", func() {
			reGenerated, err := reGenerateResponseFromOllama(ollamaClient, generated.Context, TryAgain)
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			lastRspHash = sha256.New().Sum([]byte(reGenerated.Response))
			clippyPopUp(a, input, &reGenerated)
			w.Hide()
		}),
		widget.NewButton("Make the text more Friendly", func() {
			reGenerated, err := reGenerateResponseFromOllama(ollamaClient, generated.Context, MakeItFriendlyRedo)
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			lastRspHash = sha256.New().Sum([]byte(reGenerated.Response))
			clippyPopUp(a, input, &reGenerated)
			w.Hide()
		}),
		widget.NewButton("Make the text a Bulleted List", func() {
			reGenerated, err := reGenerateResponseFromOllama(ollamaClient, generated.Context, MakeItAList)
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			lastRspHash = sha256.New().Sum([]byte(reGenerated.Response))
			clippyPopUp(a, input, &reGenerated)
			w.Hide()
		}),
	))

	// ctrlAltL := &desktop.CustomShortcut{KeyName: fyne.KeyL, Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift}
	// w.Canvas().AddShortcut(ctrlAltL, func(shortcut fyne.Shortcut) {
	//	slog.Info("Ctrl+Alt+L pressed")
	//})

	w.Show()
}
