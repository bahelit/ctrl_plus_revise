package chat

import (
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/bahelit/ctrl_plus_revise/internal/gui/loading"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
	"github.com/bahelit/ctrl_plus_revise/internal/store/database"
	"github.com/bahelit/ctrl_plus_revise/internal/store/models/chat"
	ollamaApi "github.com/ollama/ollama/api"
)

func newQuestionContainer(dbClient *database.ChatBot, guiApp fyne.App, tabs *container.AppTabs, ollamaClient *ollamaApi.Client, model ollama.ModelName) *fyne.Container {
	slog.Debug("New Chat")

	chatEntry := &chat.Chat{
		ID:        nil,
		Owner:     DefaultUser,
		Model:     int(model),
		Title:     "Bonkers",
		Questions: []string{},
		Responses: []string{},
	}

	submitText := widget.NewLabel("Press Shift + Enter to submit text.")
	submitText.TextStyle = fyne.TextStyle{Italic: true}

	text := widget.NewMultiLineEntry()
	text.SetMinRowsVisible(3)
	text.PlaceHolder = "Continue your question here, it remembers what is in this chat,\n" +
		"you can ask it to format the response in a certain way,\n" +
		"or to expand on or summarize the response."
	text.OnSubmitted = func(s string) {
		slog.Debug("Question submitted - keyboard shortcut", "text", s)
		err := text.Validate()
		if err != nil {
			slog.Error("Error validating question", "error", err)
			return
		}
		submitNewQuestion(dbClient, guiApp, ollamaClient, text, chatEntry, tabs)
		text.SetText("")
	}
	text.Validator = func(s string) error {
		if len(s) < 10 {
			return fmt.Errorf("question too short")
		}
		if len(s) > 10000000 {
			return fmt.Errorf("question too long, testing is needed before increasing the max length")
		}
		return nil
	}

	submitQuestionsButton := widget.NewButton("Submit Question", func() {
		slog.Debug("Question submitted", "text", text.Text)
		err := text.Validate()
		if err != nil {
			slog.Error("Error validating question", "error", err)
			return
		}
		submitNewQuestion(dbClient, guiApp, ollamaClient, text, chatEntry, tabs)
		text.SetText("")
	})

	questionWindow := container.NewBorder(submitText, submitQuestionsButton, nil, nil, text)
	return questionWindow
}

func submitNewQuestion(dbClient *database.ChatBot, guiApp fyne.App, ollamaClient *ollamaApi.Client, text *widget.Entry, yakityYak *chat.Chat, tabs *container.AppTabs) {
	loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
		"Asking question...")
	loadingScreen.Show()
	// TODO: Pass in the user selected model from dropdown.
	response, err := ollama.AskAI(guiApp, ollamaClient, text.Text)
	if err != nil {
		slog.Error("Failed to ask AI", "error", err)
		loadingScreen.Hide()
		return
	}
	loadingScreen.Hide()

	yakityYak.Context = response.Context
	yakityYak.Questions = []string{}
	yakityYak.Responses = []string{}
	yakityYak.Questions = append(yakityYak.Questions, text.Text)
	yakityYak.Responses = append(yakityYak.Responses, response.Response)

	if dbClient != nil {
		yakityYak.Title = text.Text[:14]
		err = dbClient.SaveChat(yakityYak)
		if err != nil {
			// TODO: Pop-up notification to inform the user their chat isn't being saved.
			slog.Error("Failed to save new chat", "error", err)
		}
	}
	var chitChat chat.Chat
	chitChat = *yakityYak

	newChatTab := createChatEntry(dbClient, guiApp, ollamaClient, chitChat)
	ct := container.NewTabItem(chitChat.Title, newChatTab)
	tabs.Append(ct)
	tabs.Select(ct)
	yakityYak = nil
}

func chatQuestionContainer(dbClient *database.ChatBot, guiApp fyne.App, entries *fyne.Container, scroll *container.Scroll, ollamaClient *ollamaApi.Client, yakity chat.Chat) *fyne.Container {
	slog.Debug("Chatting Question")

	text := widget.NewMultiLineEntry()
	text.SetMinRowsVisible(3)
	text.PlaceHolder = "Continue your question here, it remembers what is in this chat,\n" +
		"you can ask it to format the response in a certain way,\n" +
		"or to expand on or summarize the response."
	text.OnSubmitted = func(s string) {
		slog.Debug("Question submitted - keyboard shortcut", "text", s)
		err := text.Validate()
		if err != nil {
			slog.Error("Error validating question", "error", err)
			return
		}
		submitQuestionToChat(guiApp, err, ollamaClient, dbClient, &yakity, text, entries, s)
		text.SetText("")
		scroll.ScrollToBottom()
	}
	text.Validator = func(s string) error {
		if len(s) < 10 {
			return fmt.Errorf("question too short")
		}
		if len(s) > 10000000 {
			return fmt.Errorf("question too long, testing is needed before increasing the max length")
		}
		return nil
	}

	submitQuestionsButton := widget.NewButton("Submit Question", func() {
		slog.Debug("Question submitted", "text", text.Text)
		err := text.Validate()
		if err != nil {
			slog.Error("Error validating question", "error", err)
			return
		}
		submitQuestionToChat(guiApp, err, ollamaClient, dbClient, &yakity, text, entries, text.Text)
		text.SetText("")
		scroll.ScrollToBottom()
	})

	questionWindow := container.NewBorder(nil, submitQuestionsButton, nil, nil, text)
	return questionWindow
}

func submitQuestionToChat(guiApp fyne.App, err error, ollamaClient *ollamaApi.Client, dbClient *database.ChatBot, yakity *chat.Chat, text *widget.Entry, entries *fyne.Container, questionFromUser string) {
	loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
		"Asking question...")
	loadingScreen.Show()
	response, err := ollama.AskAIWithContext(guiApp, ollamaClient, yakity.Context, questionFromUser)
	if err != nil {
		slog.Error("Failed to ask AI", "error", err)
		loadingScreen.Hide()
		return
	}
	loadingScreen.Hide()
	yakity.Context = response.Context
	if dbClient != nil {
		err = dbClient.UpdateChat(yakity)
		if err != nil {
			slog.Error("Failed to save new chat", "error", err)
		}
	} else {
		slog.Warn("Failed to save new chat", "error", err)
	}
	// TODO: Add tab, add tab close/save buttons, copy button should be with text response
	newEntry := addChatEntry(text.Text, response.Response)
	entries.Add(newEntry)
	yakity.Questions = append(yakity.Questions, text.Text)
	yakity.Responses = append(yakity.Responses, response.Response)
}
