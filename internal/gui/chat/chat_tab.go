package chat

import (
	"crypto/sha256"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/bahelit/ctrl_plus_revise/internal/gui/loading"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/shortcuts"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
	ollamaApi "github.com/ollama/ollama/api"
)

func chatTab(guiApp fyne.App, tabs *container.AppTabs, ollamaClient *ollamaApi.Client, question string, response *ollamaApi.GenerateResponse) {
	w := guiApp.NewWindow("Ctrl+Revise")
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

	buttons := container.NewHBox(
		widget.NewButton("Try Again", func() {
			loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
				"Trying Again...")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithContext(guiApp, ollamaClient, response.Context, ollama.TryAgain)
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			shortcuts.LastClipboardContent = sha256.Sum256([]byte(response.Response))
			w.Hide()
			loadingScreen.Hide()
			chatTab(guiApp, tabs, ollamaClient, question, &reGenerated)
		}),
		widget.NewButton("Make the text more Friendly", func() {
			loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
				"Making the text more Friendly...")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithContext(guiApp, ollamaClient, response.Context, ollama.MakeItFriendlyRedo)
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			shortcuts.LastClipboardContent = sha256.Sum256([]byte(response.Response))
			w.Hide()
			loadingScreen.Hide()
			chatTab(guiApp, tabs, ollamaClient, question, &reGenerated)
		}),
		widget.NewButton("Make the text more Professional", func() {
			loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
				"Making the text more Professional...")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithContext(guiApp, ollamaClient, response.Context, ollama.MakeItProfessionalRedo)
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			shortcuts.LastClipboardContent = sha256.Sum256([]byte(response.Response))
			w.Hide()
			loadingScreen.Hide()
			chatTab(guiApp, tabs, ollamaClient, question, &reGenerated)
		}),
		widget.NewButton("Make the text a Bulleted List", func() {
			loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
				"Making the text a Bulleted List...")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithContext(guiApp, ollamaClient, response.Context, ollama.MakeItAListRedo)
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			shortcuts.LastClipboardContent = sha256.Sum256([]byte(response.Response))
			w.Hide()
			loadingScreen.Hide()
			chatTab(guiApp, tabs, ollamaClient, question, &reGenerated)
		}),
		widget.NewButtonWithIcon("Copy generated text to Clipboard", theme.ContentCopyIcon(), func() {
			w.Clipboard().SetContent(response.Response)
			w.Close()
		}),
	)
	buttons.Layout = layout.NewAdaptiveGridLayout(3)

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
