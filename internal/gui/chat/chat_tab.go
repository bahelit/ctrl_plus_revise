package chat

import (
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
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

func mainQuestionContainer(guiApp fyne.App, tabs *container.AppTabs, ollamaClient *ollamaApi.Client) *fyne.Container {
	slog.Debug("New Chat")

	chat := &Chat{
		UUID:      uuid.UUID{},
		Owner:     uuid.UUID{},
		Model:     0,
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
		loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
			"Asking question...")
		loadingScreen.Show()
		response, err := ollama.AskAI(guiApp, ollamaClient, s)
		if err != nil {
			slog.Error("Failed to ask AI", "error", err)
			loadingScreen.Hide()
			return
		}
		loadingScreen.Hide()
		// TODO: Add tab, add tab close/save buttons, copy button should be with text response
		chat.Questions = append(chat.Questions, text.Text)
		chat.Responses = append(chat.Responses, response.Response)
		chatTab := createChatEntry(guiApp, tabs, ollamaClient, chat)
		ct := container.NewTabItem(chat.Title, chatTab)
		tabs.Append(ct)
		tabs.Select(ct)
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
		loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
			"Asking question...")
		loadingScreen.Show()
		response, err := ollama.AskAI(guiApp, ollamaClient, text.Text)
		if err != nil {
			slog.Error("Failed to ask AI", "error", err)
			loadingScreen.Hide()
			return
		}
		loadingScreen.Hide()

		chat.Questions = append(chat.Questions, text.Text)
		chat.Responses = append(chat.Responses, response.Response)
		chatTab := createChatEntry(guiApp, tabs, ollamaClient, chat)
		ct := container.NewTabItem(chat.Title, chatTab)
		tabs.Append(ct)
		tabs.Select(ct)
	})

	questionWindow := container.NewBorder(submitText, submitQuestionsButton, nil, nil, text)
	return questionWindow
}

func chatQuestionContainer(guiApp fyne.App, entries *fyne.Container, ollamaClient *ollamaApi.Client, chat *Chat) *fyne.Container {
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
		loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
			"Asking question...")
		loadingScreen.Show()
		response, err := ollama.AskAI(guiApp, ollamaClient, s)
		if err != nil {
			slog.Error("Failed to ask AI", "error", err)
			loadingScreen.Hide()
			return
		}
		loadingScreen.Hide()
		// TODO: Add tab, add tab close/save buttons, copy button should be with text response
		newEntry := addChatEntry(text.Text, response.Response)
		entries.Add(newEntry)
		chat.Questions = append(chat.Questions, text.Text)
		chat.Responses = append(chat.Responses, response.Response)
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
		loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
			"Asking question...")
		loadingScreen.Show()
		response, err := ollama.AskAI(guiApp, ollamaClient, text.Text)
		if err != nil {
			slog.Error("Failed to ask AI", "error", err)
			loadingScreen.Hide()
			return
		}
		loadingScreen.Hide()

		newEntry := addChatEntry(text.Text, response.Response)
		entries.Add(newEntry)
		chat.Questions = append(chat.Questions, text.Text)
		chat.Responses = append(chat.Responses, response.Response)
	})

	questionWindow := container.NewBorder(nil, submitQuestionsButton, nil, nil, text)
	return questionWindow
}
