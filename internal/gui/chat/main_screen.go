package chat

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/bahelit/ctrl_plus_revise/internal/config"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/bindings"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/loading"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/question"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/settings"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
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
	var (
		screenHeight float32 = 575.0
		screenWidth  float32 = 655.0
		//tabs         container.AppTabs
	)
	w := guiApp.NewWindow("Ctrl+Revise Chatbot")
	w.Resize(fyne.NewSize(screenWidth, screenHeight))

	hello := widget.NewLabel("Glad to Help!")
	hello.TextStyle = fyne.TextStyle{Bold: true}
	hello.Alignment = fyne.TextAlignCenter

	mainEntry := createMainEntry(guiApp, ollamaClient)
	chatLayout := createChatEntry(guiApp, w)
	chatLayout2 := createChatEntry(guiApp, w)

	startChatTab := container.NewTabItem("Home", mainEntry)
	chatTab := container.NewTabItem("Blah Blah Blah", chatLayout)
	chatTab1 := container.NewTabItem("Yakity Yak", chatLayout2)
	verticalTabs := container.NewAppTabs(startChatTab, chatTab, chatTab1)
	verticalTabs.SetTabLocation(container.TabLocationLeading)
	//verticalTabs.Append(container.NewTabItem("Chat", chatLayout))
	tabContainer := container.NewBorder(hello, nil, nil, nil, verticalTabs)
	w.SetContent(tabContainer)
	w.Show()
}

func createMainEntry(guiApp fyne.App, ollamaClient *ollamaApi.Client) *fyne.Container {
	chatBotSelection := settings.SelectAIModelDropDown(guiApp)
	chatBotSelection.OnChanged = func(s string) {
		modelSelected := ollama.StringToModel(s)
		guiApp.Preferences().SetInt(config.CurrentModelKey, int(modelSelected))
	}
	saveDefaultModelButton := widget.NewButton("Set as default", func() {
		bindings.AiModelDropdown.SetSelected(chatBotSelection.Selected)
		bindings.AiModelDropdown.OnChanged = func(s string) {

			modelSelected := ollama.StringToModel(s)
			guiApp.Preferences().SetInt(config.CurrentModelKey, int(modelSelected))
		}
	})

	model := container.NewVBox(chatBotSelection, container.NewCenter(saveDefaultModelButton))
	questionContainer := question.AskQuestion(guiApp, ollamaClient)

	logo := canvas.NewImageFromResource(guiApp.Icon())
	logo.FillMode = canvas.ImageFillOriginal
	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(192, 192))
	} else {
		logo.SetMinSize(fyne.NewSize(256, 256))
	}

	chatLayout := container.NewBorder(
		model,
		questionContainer,
		nil,
		nil,
		logo)
	return chatLayout
}

func createChatEntry(guiApp fyne.App, w fyne.Window) *fyne.Container {
	questionTextLabel := widget.NewLabel("Question:")
	questionTextLabel.Alignment = fyne.TextAlignLeading
	questionTextLabel.Wrapping = fyne.TextWrapWord
	questionTextLabel.TextStyle = fyne.TextStyle{Bold: true}
	questionText := widget.NewLabel("")
	questionText.Alignment = fyne.TextAlignLeading
	questionText.Wrapping = fyne.TextWrapWord

	generatedTextLabel := widget.NewLabel("AI Response:")
	generatedTextLabel.Alignment = fyne.TextAlignLeading
	generatedTextLabel.Wrapping = fyne.TextWrapWord
	generatedTextLabel.TextStyle = fyne.TextStyle{Bold: true}
	questionSection := container.NewVBox(questionTextLabel, questionText, generatedTextLabel)

	generatedText := widget.NewRichTextFromMarkdown("")
	generatedText.Wrapping = fyne.TextWrapWord
	vbox := container.NewVScroll(generatedText)
	textResponse := container.New(layout.NewAdaptiveGridLayout(1), vbox)

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

	chatLayout := container.NewBorder(
		questionSection,
		buttons,
		nil,
		nil,
		container.NewVScroll(textResponse))
	return chatLayout
}
