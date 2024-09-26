package shortcuts

import (
	"crypto/sha256"
	"fyne.io/fyne/v2"
	"log/slog"
	"time"

	"github.com/go-vgo/robotgo"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/ollama/ollama/api"
	ollamaApi "github.com/ollama/ollama/api"
	hook "github.com/robotn/gohook"

	"github.com/bahelit/ctrl_plus_revise/internal/config"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/bindings"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/loading"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
	"github.com/bahelit/ctrl_plus_revise/pkg/clipboard"
	"github.com/bahelit/ctrl_plus_revise/pkg/throttle"
)

type keyBoardShortcut struct {
	ModifierKey1 string
	ModifierKey2 *string
	Key          string
}

var (
	Throttle              = throttle.NewThrottle(1)
	Speech                *htgotts.Speech
	LastClipboardContent  [32]byte
	lastKeyPressTime      = time.Now()
	waitBetweenKeyPresses = 1 * time.Second
	keyPressSleep         = 250
	firstRun              = true
	AskKey                = keyBoardShortcut{ModifierKey1: "alt", Key: "a"}
	CtrlReviseKey         = keyBoardShortcut{ModifierKey1: "alt", Key: "c"}
	TranslateKey          = keyBoardShortcut{ModifierKey1: "alt", Key: "t"}

	selectedModel  = ollama.Llama3
	selectedPrompt = ollama.CorrectGrammar

	systemHook chan hook.Event
)

// RegisterHotkeys registers the hotkeys for the application
// TODO: Allow users changes the mappings
func RegisterHotkeys(guiApp fyne.App, ollamaClient *ollamaApi.Client) chan hook.Event {
	hook.Register(hook.KeyDown, GetAskKeys(), func(e hook.Event) {
		slog.Debug("AskKey has been pressed", "event", e)
		if time.Since(lastKeyPressTime) < waitBetweenKeyPresses {
			slog.Info("Ignoring key press", "waitBetweenKeyPresses", waitBetweenKeyPresses)
			lastKeyPressTime = time.Now()
			return
		}
		lastKeyPressTime = time.Now()
		handleAskKeyPressed(guiApp, ollamaClient)
		lastKeyPressTime = time.Now()
	})
	hook.Register(hook.KeyDown, GetCtrlReviseKeys(), func(e hook.Event) {
		slog.Debug("CtrlReviseKey has been pressed", "event", e)
		if time.Since(lastKeyPressTime) < waitBetweenKeyPresses {
			slog.Info("Ignoring key press", "waitBetweenKeyPresses", waitBetweenKeyPresses)
			return
		}
		lastKeyPressTime = time.Now()
		handleUserShortcutKeyPressed(guiApp, ollamaClient)
		lastKeyPressTime = time.Now()
	})
	hook.Register(hook.KeyDown, []string{"p", "alt"}, func(e hook.Event) {
		slog.Debug("cyclePromptKey has been pressed", "event", e)
		if time.Since(lastKeyPressTime) < waitBetweenKeyPresses {
			slog.Info("Ignoring key press", "waitBetweenKeyPresses", waitBetweenKeyPresses)
			lastKeyPressTime = time.Now()
			return
		}
		lastKeyPressTime = time.Now()
		handleCyclePromptKeyPressed(guiApp)
		lastKeyPressTime = time.Now()
	})
	hook.Register(hook.KeyDown, []string{"r", "alt"}, func(e hook.Event) {
		slog.Debug("readTextKey has been pressed", "event", e)
		if time.Since(lastKeyPressTime) < waitBetweenKeyPresses {
			slog.Info("Ignoring key press", "waitBetweenKeyPresses", waitBetweenKeyPresses)
			lastKeyPressTime = time.Now()
			return
		}
		lastKeyPressTime = time.Now()
		handleReadTextPressed(guiApp)
		lastKeyPressTime = time.Now()
	})
	hook.Register(hook.KeyDown, GetTranslateKeys(), func(e hook.Event) {
		slog.Debug("translateTextKey has been pressed", "event", e)
		if time.Since(lastKeyPressTime) < waitBetweenKeyPresses {
			slog.Info("Ignoring key press", "waitBetweenKeyPresses", waitBetweenKeyPresses)
			lastKeyPressTime = time.Now()
			return
		}
		lastKeyPressTime = time.Now()
		handleTranslatePressed(guiApp, ollamaClient)
		lastKeyPressTime = time.Now()
	})

	slog.Info("Registered hotkeys")

	return hook.Start()
}

func StartKeyboardListener(guiApp fyne.App, ollamaClient *ollamaApi.Client) bool {
	systemHook = RegisterHotkeys(guiApp, ollamaClient)
	return <-hook.Process(systemHook)
}

func GetAskKeys() []string {
	if AskKey.ModifierKey2 == nil {
		return []string{AskKey.ModifierKey1, AskKey.Key}
	}
	return []string{AskKey.ModifierKey1, *AskKey.ModifierKey2, AskKey.Key}
}

func GetCtrlReviseKeys() []string {
	if CtrlReviseKey.ModifierKey2 == nil {
		return []string{CtrlReviseKey.ModifierKey1, CtrlReviseKey.Key}
	}
	return []string{CtrlReviseKey.ModifierKey1, *CtrlReviseKey.ModifierKey2, CtrlReviseKey.Key}
}

func GetTranslateKeys() []string {
	if TranslateKey.ModifierKey2 == nil {
		return []string{TranslateKey.ModifierKey1, TranslateKey.Key}
	}
	return []string{TranslateKey.ModifierKey1, *TranslateKey.ModifierKey2, TranslateKey.Key}
}

func handleCyclePromptKeyPressed(guiApp fyne.App) {
	err := Throttle.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer Throttle.Done(err)
	if int(selectedPrompt) == len(ollama.PromptToText)-1 {
		selectedPrompt = ollama.CorrectGrammar
	} else {
		selectedPrompt++
	}
	//err = settings.SelectedPromptBinding.Set(selectedPrompt.String())
	//if err != nil {
	//	slog.Error("Failed to set selectedPromptBinding", "error", err)
	//}
	UpdateDropDownMenus()
	ChangedPromptNotification(guiApp)
	time.Sleep(1 * time.Second)
}

func UpdateDropDownMenus() {
	bindings.AiActionDropdown.SetSelected(selectedPrompt.String())
	bindings.AiModelDropdown.SetSelected(selectedModel.String())
}

func ChangedPromptNotification(guiApp fyne.App) {
	guiApp.Preferences().SetString(config.CurrentPromptKey, selectedPrompt.String())
	guiApp.SendNotification(&fyne.Notification{
		Title:   "AI Action Changed",
		Content: "AI Action has been changed to:\n" + selectedPrompt.String(),
	})
}

func handleReadTextPressed(guiApp fyne.App) {
	err := Throttle.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer Throttle.Done(err)
	speakAIResponse := guiApp.Preferences().BoolWithFallback(config.SpeakAIResponseKey, false)
	if !speakAIResponse {
		slog.Info("Reading text aloud is disabled", "speakAIResponse", speakAIResponse)
		return
	}

	err = copyCommand()
	if err != nil {
		if Speech != nil {
			_ = Speech.Speak("Failed to copy text")
		}
		return
	}

	clip, err := clipboard.ReadAll()
	if err != nil {
		slog.Error("Failed to read clipboard", "error", err)
		if Speech != nil {
			_ = Speech.Speak("Failed to read clipboard")
		}
		return
	}

	_ = Speech.Speak("")
	speakErr := Speech.Speak(clip)
	if speakErr != nil {
		slog.Error("Failed to speak", "error", speakErr)
	}
}

func handleUserShortcutKeyPressed(guiApp fyne.App, ollamaClient *ollamaApi.Client) {
	err := Throttle.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer Throttle.Done(err)

	clip, copiedText := copyTextToClipboard()
	if !copiedText {
		return
	}

	loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
		"Prompt: "+selectedPrompt.String()+"...")
	loadingScreen.Show()

	generated, err := ollama.AskAIWithPromptMsg(guiApp, ollamaClient, selectedPrompt, clip)
	if err != nil {
		// TODO: Implement error handling, tell user to restart ollama, maybe we can restart ollama here?
		slog.Error("Failed to communicate with Ollama", "error", err)
		loadingScreen.Hide()
		return
	}
	loadingScreen.Hide()

	handleGeneratedResponse(guiApp, ollamaClient, clip, &generated)
}

func handleAskKeyPressed(guiApp fyne.App, ollamaClient *ollamaApi.Client) {
	err := Throttle.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer Throttle.Done(err)

	clip, copiedText := copyTextToClipboard()
	if !copiedText {
		return
	}

	loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
		"Asking question")
	loadingScreen.Show()

	generated, err := ollama.AskAI(guiApp, ollamaClient, clip)
	if err != nil {
		slog.Error("Failed to ask AI", "error", err)
		loadingScreen.Hide()
		return
	}
	loadingScreen.Hide()

	handleGeneratedResponse(guiApp, ollamaClient, clip, &generated)
}

func handleTranslatePressed(guiApp fyne.App, ollamaClient *ollamaApi.Client) {
	err := Throttle.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer func() {
		slog.Info("Done translating")
		Throttle.Done(err)
	}()

	clip, copiedText := copyTextToClipboard()
	if !copiedText {
		return
	}

	fromLang := guiApp.Preferences().StringWithFallback(config.CurrentFromLangKey, string(ollama.English))
	toLang := guiApp.Preferences().StringWithFallback(config.CurrentToLangKey, string(ollama.Spanish))

	loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
		"Translating")
	loadingScreen.Show()

	slog.Info("Translating text", "fromLang", fromLang, "toLang", toLang)
	generated, err := ollama.AskAIToTranslate(guiApp, ollamaClient, clip, ollama.Language(fromLang), ollama.Language(toLang))
	if err != nil {
		slog.Error("Failed to ask AI", "error", err)
		loadingScreen.Hide()
		return
	}
	loadingScreen.Hide()

	err = clipboard.WriteAll(generated.Response)
	if err != nil {
		slog.Error("Failed to write to clipboard", "error", err)
		return
	}
	LastClipboardContent = sha256.Sum256([]byte(generated.Response))

	// Send a paste command to the operating system
	replaceText := guiApp.Preferences().BoolWithFallback(config.ReplaceHighlightedText, true)
	if replaceText {
		err = pasteCommand()
		if err != nil {
			return
		}
	}
}

func copyCommand() error {
	robotgo.KeySleep = 100

	if firstRun {
		firstRun = false
		robotgo.Sleep(1)
	}
	robotgo.MilliSleep(keyPressSleep)
	err := robotgo.KeyTap(robotgo.KeyC, robotgo.CmdCtrl())
	if err != nil {
		slog.Error("Failed to send copy command", "error", err)
		return err
	}
	robotgo.MilliSleep(keyPressSleep)
	return nil
}

func pasteCommand() error {
	robotgo.KeySleep = 100

	robotgo.MilliSleep(keyPressSleep)
	err := robotgo.KeyTap(robotgo.KeyV, robotgo.CmdCtrl())
	if err != nil {
		slog.Error("Failed to send paste command", "error", err)
		return err
	}
	robotgo.MilliSleep(keyPressSleep)
	return nil
}

func copyTextToClipboard() (string, bool) {
	err := copyCommand()
	if err != nil {
		return "", false
	}

	clip, err := clipboard.ReadAll()
	if err != nil {
		slog.Error("Failed to read clipboard", "error", err)
		return "", false
	}
	if clip == "" {
		slog.Info("Clipboard is empty, skipping")
		return "", false
	}

	if sha256.Sum256([]byte(clip)) == LastClipboardContent {
		slog.Debug("Clipboard content is the same as last time", "clippy", clip)
		return "", false
	}
	return clip, true
}

func handleGeneratedResponse(guiApp fyne.App, ollamaClient *ollamaApi.Client, question string, response *api.GenerateResponse) {
	slog.Debug("LastClipboardContent", "LastClipboardContent", LastClipboardContent)

	LastClipboardContent = sha256.Sum256([]byte(response.Response))
	err := clipboard.WriteAll(response.Response)
	if err != nil {
		slog.Error("Failed to write to clipboard", "error", err)
		return
	}

	// Send a paste command to the operating system
	replaceText := guiApp.Preferences().BoolWithFallback(config.ReplaceHighlightedText, true)
	if replaceText {
		err = pasteCommand()
		if err != nil {
			return
		}
	}

	//showPopUp := guiApp.Preferences().BoolWithFallback(config.ShowPopUpKey, false)
	//if showPopUp {
	//	clippy.QuestionPopUp(guiApp, ollamaClient, question, response)
	//}
}
