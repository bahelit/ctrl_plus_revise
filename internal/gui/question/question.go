package question

import (
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/layout"
	ollamaApi "github.com/ollama/ollama/api"

	"github.com/bahelit/ctrl_plus_revise/internal/gui/clippy"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/loading"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
)

func AskQuestionWindow(guiApp fyne.App, ollamaClient *ollamaApi.Client) {
	slog.Debug("Asking Question")
	var (
		screenHeight float32 = 200.0
		screenWidth  float32 = 480.0
	)
	question := guiApp.NewWindow("Ctrl+Revise Questions")
	question.Resize(fyne.NewSize(screenWidth, screenHeight))

	label1 := widget.NewLabel("Ask a Question")
	label1.TextStyle = fyne.TextStyle{Bold: true}
	label2 := widget.NewLabel("Press Shift + Enter to submit your question.")
	label2.TextStyle = fyne.TextStyle{Italic: true}

	text := widget.NewMultiLineEntry()
	text.SetMinRowsVisible(4)
	text.PlaceHolder = "Ask your question here, remember this is an AI and important\n" +
		"questions should be verified with other sources."
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
		clippy.QuestionPopUp(guiApp, ollamaClient, s, &response)
		question.Close()
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
		clippy.QuestionPopUp(guiApp, ollamaClient, text.Text, &response)
		question.Close()
	})

	topText := container.NewHBox(label1, label2)
	questionLayout := layout.NewResponsiveLayout(topText)
	buttonLayout := layout.NewResponsiveLayout(layout.Responsive(submitQuestionsButton))
	questionWindow := container.NewBorder(questionLayout, buttonLayout, nil, nil, text)
	question.SetContent(container.NewVScroll(questionWindow))
	question.Canvas().Focus(text)
	question.Show()
}

func AskQuestionContainer(guiApp fyne.App, tabs *container.AppTabs, ollamaClient *ollamaApi.Client) *fyne.Container {
	slog.Debug("Asking Question in Tab")

	label1 := widget.NewLabel("Ask a Question")
	label1.TextStyle = fyne.TextStyle{Bold: true}
	label2 := widget.NewLabel("Press Shift + Enter to submit your question.")
	label2.TextStyle = fyne.TextStyle{Italic: true}

	text := widget.NewMultiLineEntry()
	text.SetMinRowsVisible(3)
	text.PlaceHolder = "Ask your question here, remember this is an AI and important\n" +
		"questions should be verified with other sources."
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
		clippy.QuestionPopUp(guiApp, ollamaClient, s, &response)
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
		clippy.QuestionTab(guiApp, tabs, ollamaClient, text.Text, &response)
	})

	topText := container.NewHBox(label1, label2)
	questionLayout := layout.NewResponsiveLayout(topText)
	buttonLayout := layout.NewResponsiveLayout(layout.Responsive(submitQuestionsButton))
	questionWindow := container.NewBorder(questionLayout, buttonLayout, nil, nil, text)
	return questionWindow
}
