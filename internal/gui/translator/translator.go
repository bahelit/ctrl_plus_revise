package translator

import (
	"fmt"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/loading"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/settings"
	"log/slog"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	ollamaApi "github.com/ollama/ollama/api"

	"github.com/bahelit/ctrl_plus_revise/internal/config"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/shortcuts"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
	"github.com/bahelit/ctrl_plus_revise/pkg/clipboard"
)

func TranslateText(guiApp fyne.App, ollamaClient *ollamaApi.Client) {
	slog.Debug("Asking Question")
	var (
		screenHeight float32 = 480.0
		screenWidth  float32 = 650.0
	)
	translator := guiApp.NewWindow("Ctrl+Revise Translation Tool")
	translator.Resize(fyne.NewSize(screenWidth, screenHeight))

	var (
		from         = widget.NewMultiLineEntry()
		to           = widget.NewMultiLineEntry()
		fromDropdown = settings.SelectTranslationFromDropDown(guiApp)
		toDropdown   = settings.SelectTranslationToDropDown(guiApp)
	)

	to.Wrapping = fyne.TextWrapWord

	var keyPressDelay *time.Timer
	from.SetPlaceHolder("Enter Text to Translate")
	from.Wrapping = fyne.TextWrapWord
	from.OnChanged = func(t string) {
		if keyPressDelay != nil {
			keyPressDelay.Stop()
		}
		keyPressDelay = time.AfterFunc(time.Duration(1)*time.Second, func() {
			slog.Debug("Translating Text")
			err := from.Validate()
			if err != nil {
				slog.Warn("text validating failed for translation", "error", err)
				return
			}
			handleTranslateRequest(guiApp, ollamaClient, from, to)
			translator.Canvas().Focus(from)
		})
	}
	from.Validator = func(s string) error {
		if len(s) < 8 {
			return fmt.Errorf("text is too short")
		}
		if len(s) > 10000000 {
			return fmt.Errorf("text is too long, testing is needed before increasing the max length")
		}
		return nil
	}

	top := container.NewBorder(nil, nil,
		container.NewHBox(fromDropdown),
		container.NewHBox(toDropdown),
		container.NewHBox(layout.NewSpacer(), layout.NewSpacer()),
	)

	combo := container.NewBorder(top, nil, nil, nil,
		container.NewGridWithColumns(2, from, to),
	)

	translator.SetContent(combo)
	translator.Canvas().Focus(from)
	translator.Show()

}

func handleTranslateRequest(guiApp fyne.App, ollamaClient *ollamaApi.Client, from, to *widget.Entry) {
	err := shortcuts.Throttle.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer func() {
		slog.Debug("Done translating")
		shortcuts.Throttle.Done(err)
	}()

	fromLang := guiApp.Preferences().StringWithFallback(config.CurrentFromLangKey, string(ollama.English))
	toLang := guiApp.Preferences().StringWithFallback(config.CurrentToLangKey, string(ollama.Spanish))

	loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg, "Translating text")
	loadingScreen.Show()

	slog.Debug("Translating text", "fromLang", fromLang, "toLang", toLang)
	generated, err := ollama.AskAIToTranslate(guiApp, ollamaClient, from.Text, ollama.Language(fromLang), ollama.Language(toLang))
	if err != nil {
		slog.Error("Failed to ask AI", "error", err)
		loadingScreen.Hide()
		return
	}
	loadingScreen.Hide()
	to.SetText(generated.Response)
	err = clipboard.WriteAll(generated.Response)
	if err != nil {
		slog.Error("Failed to write to clipboard", "error", err)
		return
	}

	return
}
