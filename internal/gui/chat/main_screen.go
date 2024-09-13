package chat

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/loading"
	"github.com/google/uuid"
	ollamaApi "github.com/ollama/ollama/api"
)

type Chat struct {
	UUID      uuid.UUID `json:"uuid"`
	Owner     uuid.UUID `json:"owner"`
	Title     string    `json:"title"`
	Questions []string  `json:"questions"`
	Responses []string  `json:"responses"`
}

// TODO: Use an LRU cache to limit the items in ChatContents.
var (
	Chats        map[uuid.UUID]string
	ChatContents map[uuid.UUID]Chat
)

func ConversationManager(guiApp fyne.App, ollamaClient *ollamaApi.Client) {
	w := guiApp.NewWindow("Ctrl+Revise Chatbot")
	w.Resize(fyne.NewSize(640, 500))
	hello := widget.NewLabel("Glad to Help!")
	hello.TextStyle = fyne.TextStyle{Bold: true}
	hello.Alignment = fyne.TextAlignCenter

	questionText := widget.NewLabel("Question:")
	questionText.Alignment = fyne.TextAlignLeading
	questionText.Wrapping = fyne.TextWrapWord
	questionText.TextStyle = fyne.TextStyle{Bold: true}
	questionText1 := widget.NewLabel("")
	questionText1.Alignment = fyne.TextAlignLeading
	questionText1.Wrapping = fyne.TextWrapWord

	generatedText := widget.NewLabel("AI Response:")
	generatedText.Alignment = fyne.TextAlignLeading
	generatedText.Wrapping = fyne.TextWrapWord
	generatedText.TextStyle = fyne.TextStyle{Bold: true}

	generatedText1 := widget.NewRichTextFromMarkdown("")
	generatedText1.Wrapping = fyne.TextWrapWord

	vbox := container.NewVScroll(generatedText1)

	buttons := container.NewHBox(
		widget.NewButton("Try Again", func() {
			loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
				"Prompt: ")
			loadingScreen.Show()
			w.Hide()
			loadingScreen.Hide()
		}),
		widget.NewButton("Make the text more Friendly", func() {
			loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
				"Prompt: ")
			loadingScreen.Show()
			w.Hide()
			loadingScreen.Hide()
		}),
		widget.NewButton("Make the text more Professional", func() {
			loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
				"Prompt: ")
			loadingScreen.Show()
			w.Hide()
			loadingScreen.Hide()
		}),
		widget.NewButton("Make the text a Bulleted List", func() {
			loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
				"Prompt: ")
			loadingScreen.Show()
			w.Hide()
			loadingScreen.Hide()
		}),
		widget.NewButtonWithIcon("Copy generated text to Clipboard", theme.ContentCopyIcon(), func() {
			w.Clipboard().SetContent("Easter Egg")
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
