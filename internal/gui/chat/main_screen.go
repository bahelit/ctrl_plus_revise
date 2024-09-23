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
	"github.com/bahelit/ctrl_plus_revise/internal/store/database"
	"github.com/bahelit/ctrl_plus_revise/internal/store/models/chat"
)

// TODO: Use an LRU cache to limit the items in ChatContents.
var (
	Chats        map[uuid.UUID]string
	ChatContents map[uuid.UUID]chat.Chat
)

const (
	DefaultUser = "default"
)

func ConversationManager(guiApp fyne.App, ollamaClient *ollamaApi.Client) {
	var (
		screenHeight float32 = 675.0
		screenWidth  float32 = 755.0
		verticalTabs         = container.NewAppTabs()
		dbClient     *database.ChatBot
		err          error
	)

	dbClient, err = database.NewSQLiteDB()
	if err != nil {
		slog.Error("Can NOT save chats", "err", err.Error())
	}

	w := guiApp.NewWindow("Ctrl+Revise Private Chatbot")
	w.Resize(fyne.NewSize(screenWidth, screenHeight))

	//hello := widget.NewLabel("This conversation is between us 🙈 🙉 🙊")
	//hello.TextStyle = fyne.TextStyle{Bold: true}
	//hello.Alignment = fyne.TextAlignCenter

	mainEntry := createNewChatEntry(dbClient, guiApp, verticalTabs, ollamaClient)
	startChatTab := container.NewTabItem("Home", mainEntry)
	verticalTabs.SetItems([]*container.TabItem{startChatTab})

	if dbClient != nil {
		savedChats, err := dbClient.GetAllChats(DefaultUser)
		if err != nil {
			slog.Error("Can NOT get chat history", "err", err.Error())
		}
		for i := range savedChats {
			sc := createChatEntry(dbClient, guiApp, ollamaClient, *savedChats[i])
			savedChatTab := container.NewTabItem(savedChats[i].Title, sc)
			verticalTabs.Append(savedChatTab)
		}
	} else {
		slog.Warn("Can not access saved chats")
	}

	verticalTabs.SetTabLocation(container.TabLocationLeading)
	tabContainer := container.NewBorder(nil, nil, nil, nil, verticalTabs)
	w.SetContent(tabContainer)
	w.Show()
}

func createNewChatEntry(dbClient *database.ChatBot, guiApp fyne.App, tabs *container.AppTabs, ollamaClient *ollamaApi.Client) *fyne.Container {
	var selectedModel ollama.ModelName
	chatBotSelection := settings.SelectAIModelDropDown(guiApp)
	chatBotSelection.OnChanged = func(s string) {
		selectedModel = ollama.StringToModel(s)
		guiApp.Preferences().SetInt(config.CurrentChatModelKey, int(selectedModel))
	}
	saveDefaultModelButton := widget.NewButton("Set as default", func() {
		bindings.AiModelDropdown.SetSelected(chatBotSelection.Selected)
		guiApp.Preferences().SetInt(config.CurrentChatModelKey, int(selectedModel))
		guiApp.Preferences().SetInt(config.CurrentModelKey, int(selectedModel))
	})
	model := container.NewVBox(chatBotSelection, container.NewCenter(saveDefaultModelButton))

	questionContainer := newQuestionContainer(dbClient, guiApp, tabs, ollamaClient, selectedModel)

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

func createChatEntry(dbClient *database.ChatBot, guiApp fyne.App, ollamaClient *ollamaApi.Client, chatEntry chat.Chat) *fyne.Container {
	chatHeader := widget.NewLabel("Model: " + ollama.ModelName(chatEntry.Model).String())
	entries := container.NewVBox()
	var widgyCard *widget.Card
	for key, questionFromChat := range chatEntry.Questions {
		if chatEntry.Responses[key] != "" {
			widgyCard = addChatEntry(questionFromChat, chatEntry.Responses[key])
		} else {
			widgyCard = addChatEntry(questionFromChat, "")
		}
		entries.Add(widgyCard)
	}

	allChats := container.NewVScroll(entries)
	questionContainer := chatQuestionContainer(dbClient, guiApp, entries, allChats, ollamaClient, chatEntry)

	chatLayout := container.NewBorder(
		chatHeader,
		questionContainer,
		nil,
		nil,
		allChats)
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
