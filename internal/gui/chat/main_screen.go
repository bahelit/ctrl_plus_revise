package chat

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	ollamaApi "github.com/ollama/ollama/api"

	"github.com/bahelit/ctrl_plus_revise/internal/config"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/bindings"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/loading"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/settings"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
)

type Chat struct {
	UUID      uuid.UUID `json:"uuid"`
	Owner     uuid.UUID `json:"owner"`
	Model     int       `json:"model"`
	Context   int       `json:"context"`
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
		screenHeight float32 = 675.0
		screenWidth  float32 = 755.0
		verticalTabs         = container.NewAppTabs()
	)
	w := guiApp.NewWindow("Ctrl+Revise Private Chatbot")
	w.Resize(fyne.NewSize(screenWidth, screenHeight))

	//hello := widget.NewLabel("This conversation is between us 🙈 🙉 🙊")
	//hello.TextStyle = fyne.TextStyle{Bold: true}
	//hello.Alignment = fyne.TextAlignCenter

	chat1 := Chat{
		UUID:      uuid.UUID{},
		Owner:     uuid.UUID{},
		Model:     0,
		Title:     "Bonkers",
		Questions: []string{"What is a fart made of?", "What gases?"},
		Responses: []string{"Stinky gases", "The stinky kind"},
	}
	chat2 := Chat{
		UUID:      uuid.UUID{},
		Owner:     uuid.UUID{},
		Model:     0,
		Title:     "Bonkers",
		Questions: []string{"What is poop made of?", "What kind of leftovers?"},
		Responses: []string{"Leftovers", "From the fridge."},
	}

	mainEntry := createMainEntry(guiApp, verticalTabs, ollamaClient)
	chatLayout := createChatEntry(guiApp, verticalTabs, ollamaClient, &chat1)
	chatLayout2 := createChatEntry(guiApp, verticalTabs, ollamaClient, &chat2)

	startChatTab := container.NewTabItem("Home", mainEntry)
	chatTab0 := container.NewTabItem("Blah Blah Blah", chatLayout)
	chatTab1 := container.NewTabItem("Yakity Yak", chatLayout2)

	verticalTabs.SetItems([]*container.TabItem{startChatTab, chatTab0, chatTab1})
	verticalTabs.SetTabLocation(container.TabLocationLeading)
	tabContainer := container.NewBorder(nil, nil, nil, nil, verticalTabs)
	w.SetContent(tabContainer)
	w.Show()
}

func createMainEntry(guiApp fyne.App, tabs *container.AppTabs, ollamaClient *ollamaApi.Client) *fyne.Container {
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

	questionContainer := mainQuestionContainer(guiApp, tabs, ollamaClient)

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

func createChatEntry(guiApp fyne.App, tabs *container.AppTabs, ollamaClient *ollamaApi.Client, chat *Chat) *fyne.Container {
	chatHeader := widget.NewLabel("Model: " + ollama.ModelName(chat.Model).String())
	entries := container.NewVBox()
	var widgyCard *widget.Card
	for key, questionFromChat := range chat.Questions {
		if chat.Responses[key] != "" {
			widgyCard = addChatEntry(questionFromChat, chat.Responses[key])
		} else {
			widgyCard = addChatEntry(questionFromChat, "")
		}
		entries.Add(widgyCard)
	}

	allChats := container.NewVScroll(entries)
	questionContainer := chatQuestionContainer(guiApp, entries, ollamaClient, chat)

	chatLayout := container.NewBorder(
		chatHeader,
		questionContainer,
		nil,
		nil,
		container.NewVScroll(allChats))
	return chatLayout
}

func addChatEntry(questionFromChat, responseFromAI string) *widget.Card {
	questionLabel := widget.NewLabel("User Question:")
	questionLabel.Alignment = fyne.TextAlignLeading
	questionLabel.Wrapping = fyne.TextWrapWord
	questionLabel.TextStyle = fyne.TextStyle{Bold: true}

	questionText := widget.NewLabel(questionFromChat)
	questionText.Alignment = fyne.TextAlignTrailing
	questionText.Wrapping = fyne.TextWrapWord

	generatedTextLabel := widget.NewLabel("AI Response:")
	generatedTextLabel.Alignment = fyne.TextAlignLeading
	generatedTextLabel.Wrapping = fyne.TextWrapWord
	generatedTextLabel.TextStyle = fyne.TextStyle{Bold: true}

	generatedText := widget.NewRichTextFromMarkdown(responseFromAI)
	generatedText.Wrapping = fyne.TextWrapWord

	chatEntryContainer := container.NewVBox(questionLabel, questionText, generatedTextLabel, generatedText)
	chatLog := widget.NewCard("", "", chatEntryContainer)
	return chatLog
}

func chatButtons(guiApp fyne.App) *fyne.Container {
	buttons := container.NewHBox(
		widget.NewButton("Try Again", func() {
			loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
				"Prompt: ")
			loadingScreen.Show()
			loadingScreen.Hide()
		}),
		widget.NewButton("Make the text more Friendly", func() {
			loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
				"Prompt: ")
			loadingScreen.Show()
			loadingScreen.Hide()
		}),
		widget.NewButton("Make the text more Professional", func() {
			loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
				"Prompt: ")
			loadingScreen.Show()
			loadingScreen.Hide()
		}),
		widget.NewButton("Make the text a Bulleted List", func() {
			loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
				"Prompt: ")
			loadingScreen.Show()
			loadingScreen.Hide()
		}),
		widget.NewButtonWithIcon("Copy generated text to Clipboard", theme.ContentCopyIcon(), func() {
			// TODO: Add clipboard command
			slog.Warn("TO BE Implemented")
		}),
	)
	buttons.Layout = layout.NewAdaptiveGridLayout(3)
	return buttons
}
