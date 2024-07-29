package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ollama/ollama/api"
)

func questionPopUp(a fyne.App, question string, generated *api.GenerateResponse) {
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

	generatedText1 := widget.NewRichTextFromMarkdown(generated.Response)

	vbox := container.NewVScroll(generatedText1)
	vbox.SetMinSize(fyne.NewSize(640, 300))

	buttons := container.NewPadded(container.NewVBox(
		widget.NewButtonWithIcon("Copy generated text to Clipboard", theme.ContentCopyIcon(), func() {
			w.Clipboard().SetContent(generated.Response)
			w.Close()
		})))

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
