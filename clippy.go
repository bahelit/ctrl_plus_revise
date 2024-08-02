package main

import (
	"crypto/sha256"
	"github.com/bahelit/ctrl_plus_revise/internal/gui"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	ollamaApi "github.com/ollama/ollama/api"
)

func questionPopUp(a fyne.App, question string, response *ollamaApi.GenerateResponse) {
	w := a.NewWindow("Ctrl+Revise")
	w.Resize(fyne.NewSize(640, 500))
	hello := widget.NewLabel("Glad to Help!")
	hello.TextStyle = fyne.TextStyle{Bold: true}
	hello.Alignment = fyne.TextAlignCenter

	questionText := widget.NewLabel("Question:")
	questionText.Alignment = fyne.TextAlignLeading
	questionText.Wrapping = fyne.TextWrapWord
	questionText.TextStyle = fyne.TextStyle{Bold: true}
	questionText1 := widget.NewLabel(question)
	questionText1.Alignment = fyne.TextAlignLeading
	questionText1.Wrapping = fyne.TextWrapWord

	generatedText := widget.NewLabel("AI Response:")
	generatedText.Alignment = fyne.TextAlignLeading
	generatedText.Wrapping = fyne.TextWrapWord
	generatedText.TextStyle = fyne.TextStyle{Bold: true}

	generatedText1 := widget.NewRichTextFromMarkdown(response.Response)
	generatedText1.Wrapping = fyne.TextWrapWord

	vbox := container.NewVScroll(generatedText1)

	model := guiApp.Preferences().IntWithFallback(CurrentModelKey, int(ollama.Llama3Dot1))

	buttons := container.NewPadded(container.NewVBox(
		widget.NewButton("Try Again", func() {
			loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
				"Using model: "+ollama.ModelName(model).String()+" with prompt: "+selectedPrompt.String()+"...")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithContext(ollamaClient, ollama.ModelName(model), response.Context, ollama.TryAgain)
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			lastClipboardContent = sha256.Sum256([]byte(response.Response))
			w.Hide()
			loadingScreen.Hide()
			questionPopUp(a, question, &reGenerated)
		}),
		widget.NewButton("Make the text more Friendly", func() {
			loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
				"Using model: "+ollama.ModelName(model).String()+" with prompt: "+selectedPrompt.String()+"...")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithContext(ollamaClient, ollama.ModelName(model), response.Context, ollama.MakeItFriendlyRedo)
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			lastClipboardContent = sha256.Sum256([]byte(response.Response))
			w.Hide()
			loadingScreen.Hide()
			questionPopUp(a, question, &reGenerated)
		}),
		widget.NewButton("Make the text more Professional", func() {
			loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
				"Using model: "+ollama.ModelName(model).String()+" with prompt: "+selectedPrompt.String()+"...")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithContext(ollamaClient, ollama.ModelName(model), response.Context, ollama.MakeItProfessionalRedo)
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			lastClipboardContent = sha256.Sum256([]byte(response.Response))
			w.Hide()
			loadingScreen.Hide()
			questionPopUp(a, question, &reGenerated)
		}),
		widget.NewButton("Make the text a Bulleted List", func() {
			loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
				"Using model: "+ollama.ModelName(model).String()+" with prompt: "+selectedPrompt.String()+"...")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithContext(ollamaClient, ollama.ModelName(model), response.Context, ollama.MakeItAListRedo)
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			lastClipboardContent = sha256.Sum256([]byte(response.Response))
			w.Hide()
			loadingScreen.Hide()
			questionPopUp(a, question, &reGenerated)
		}),
		widget.NewButtonWithIcon("Copy generated text to Clipboard", theme.ContentCopyIcon(), func() {
			w.Clipboard().SetContent(response.Response)
			w.Close()
		}),
	))

	grid := container.New(layout.NewAdaptiveGridLayout(1), vbox)

	questionSection := container.NewVBox(hello, questionText, questionText1, generatedText)
	w.SetContent(container.NewBorder(
		questionSection,
		buttons,
		nil,
		nil,
		container.NewVScroll(grid),
	))
	w.Show()
}
