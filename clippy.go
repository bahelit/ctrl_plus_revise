package main

import (
	"fyne.io/fyne/v2/layout"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ollama/ollama/api"
)

func clippyPopUp(a fyne.App, input string, generated *api.GenerateResponse) {
	w := a.NewWindow("Ctrl+Revise")
	hello := widget.NewLabel("Glad To Help!")
	hello.TextStyle = fyne.TextStyle{Bold: true}
	hello.Alignment = fyne.TextAlignCenter

	originalText := widget.NewLabel("Original text:\n" + input)
	originalText.Alignment = fyne.TextAlignLeading

	generatedText := widget.NewLabel("AI Generated text:\n" + generated.Response)
	generatedText.Alignment = fyne.TextAlignLeading

	buttons := container.NewVBox(
		widget.NewButton("Copy generated text to Clipboard", func() {
			w.Clipboard().SetContent(generated.Response)
			w.Close()
		}),
		widget.NewButton("Make the text more Friendly", func() {
			reGenerated, err := reGenerateResponseFromOllama(ollamaClient, generated.Context, MakeItFriendlyRedo)
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			clippyPopUp(a, input, &reGenerated)
			w.Close()
		}),
		widget.NewButton("Make the text a Bulleted List", func() {
			reGenerated, err := reGenerateResponseFromOllama(ollamaClient, generated.Context, MakeItAList)
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			clippyPopUp(a, input, &reGenerated)
			w.Close()
		}))

	w.SetContent(container.NewVBox(
		hello,
		originalText,
		generatedText,
		buttons,
	))
	w.Show()
}

func questionPopUp(a fyne.App, question string, generated *api.GenerateResponse) {
	w := a.NewWindow("Ctrl+Revise")
	w.Resize(fyne.NewSize(640, 400))
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
	generatedText1 := widget.NewLabel(generated.Response)
	generatedText1.Alignment = fyne.TextAlignLeading
	generatedText1.Wrapping = fyne.TextWrapWord

	vbox := container.NewVScroll(generatedText1)
	vbox.SetMinSize(fyne.NewSize(630, 250))

	buttons := container.NewPadded(container.NewVBox(
		widget.NewButton("Copy generated text to Clipboard", func() {
			w.Clipboard().SetContent(generated.Response)
			w.Close()
		}),
		widget.NewButton("Make the text a Bulleted List", func() {
			reGenerated, err := reGenerateResponseFromOllama(ollamaClient, generated.Context, MakeItAList)
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			clippyPopUp(a, question, &reGenerated)
			w.Close()
		})))

	grid := container.New(layout.NewGridLayout(1), vbox)

	w.SetContent(container.NewVBox(
		hello,
		questionText,
		questionText1,
		generatedText,
		grid,
		buttons,
	))
	w.Show()
}
