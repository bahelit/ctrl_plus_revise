package main

import (
	"crypto/sha256"
	"log/slog"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
	"github.com/ollama/ollama/api"
	hook "github.com/robotn/gohook"

	"github.com/bahelit/ctrl_plus_revise/internal/gui"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
	"github.com/bahelit/ctrl_plus_revise/pkg/bytesize"
	"github.com/bahelit/ctrl_plus_revise/pkg/clipboard"
	dirsize "github.com/bahelit/ctrl_plus_revise/pkg/dir_size"
	"github.com/bahelit/ctrl_plus_revise/pkg/throttle"
)

const (
	thinkingMsg = "Thinking..."
)

var (
	th                    = throttle.NewThrottle(1)
	speech                *htgotts.Speech
	lastClipboardContent  [32]byte
	lastKeyPressTime      = time.Now()
	waitBetweenKeyPresses = 1 * time.Second
	keyPressSleep         = 250
	firstRun              = true
)

// RegisterHotkeys registers the hotkeys for the application
// TODO: Allow users changes the mappings
func RegisterHotkeys() {
	hook.Register(hook.KeyDown, []string{"a", "alt"}, func(e hook.Event) {
		slog.Debug("askKey has been pressed", "event", e)
		if time.Since(lastKeyPressTime) < waitBetweenKeyPresses {
			slog.Info("Ignoring key press", "waitBetweenKeyPresses", waitBetweenKeyPresses)
			lastKeyPressTime = time.Now()
			return
		}
		lastKeyPressTime = time.Now()
		handleAskKeyPressed()
		lastKeyPressTime = time.Now()
	})
	hook.Register(hook.KeyDown, []string{"c", "alt"}, func(e hook.Event) {
		slog.Debug("userShortcutKey has been pressed", "event", e)
		if time.Since(lastKeyPressTime) < waitBetweenKeyPresses {
			slog.Info("Ignoring key press", "waitBetweenKeyPresses", waitBetweenKeyPresses)
			return
		}
		lastKeyPressTime = time.Now()
		handleUserShortcutKeyPressed()
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
		handleCyclePromptKeyPressed()
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
		handleReadTextPressed()
		lastKeyPressTime = time.Now()
	})
	hook.Register(hook.KeyDown, []string{"t", "alt"}, func(e hook.Event) {
		slog.Debug("translateTextKey has been pressed", "event", e)
		if time.Since(lastKeyPressTime) < waitBetweenKeyPresses {
			slog.Info("Ignoring key press", "waitBetweenKeyPresses", waitBetweenKeyPresses)
			lastKeyPressTime = time.Now()
			return
		}
		lastKeyPressTime = time.Now()
		handleTranslatePressed()
		lastKeyPressTime = time.Now()
	})
	initSpeech()
	slog.Info("Registered hotkeys")

	s := hook.Start()
	defer hook.End()
	<-hook.Process(s)
}

func initSpeech() {
	speakResponse := guiApp.Preferences().BoolWithFallback(SpeakAIResponseKey, false)
	if speakResponse {
		// TODO: the audio files need to be cleaned up periodically.
		speech = &htgotts.Speech{Folder: "audio", Language: voices.English, Handler: &handlers.Native{}}
		dirInfo, _ := dirsize.GetDirInfo(os.DirFS("audio"))
		slog.Info("AI Speech Recordings", "fileCount", dirInfo.FileCount, "size", bytesize.New(float64(dirInfo.TotalSize)))
		_ = speech.Speak("")
	}
}

func handleCyclePromptKeyPressed() {
	err := th.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer th.Done(err)
	if int(selectedPrompt) == len(ollama.PromptToText)-1 {
		selectedPrompt = ollama.CorrectGrammar
	} else {
		selectedPrompt++
	}
	err = selectedPromptBinding.Set(selectedPrompt.String())
	if err != nil {
		slog.Error("Failed to set selectedPromptBinding", "error", err)

	}
	UpdateDropDownMenus()
	ChangedPromptNotification()
	time.Sleep(1 * time.Second)
}

func handleReadTextPressed() {
	err := th.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer th.Done(err)
	speakAIResponse := guiApp.Preferences().BoolWithFallback(SpeakAIResponseKey, false)
	if !speakAIResponse {
		slog.Info("Reading text aloud is disabled", "speakAIResponse", speakAIResponse)
		return
	}

	err = copyCommand()
	if err != nil {
		if speech != nil {
			_ = speech.Speak("Failed to copy text")
		}
		return
	}

	clippy, err := clipboard.ReadAll()
	if err != nil {
		slog.Error("Failed to read clipboard", "error", err)
		if speech != nil {
			_ = speech.Speak("Failed to read clipboard")
		}
		return
	}

	_ = speech.Speak("")
	speakErr := speech.Speak(clippy)
	if speakErr != nil {
		slog.Error("Failed to speak", "error", speakErr)
	}
}

func handleUserShortcutKeyPressed() {
	err := th.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer th.Done(err)

	clippy, copiedText := copyTextToClipboard()
	if !copiedText {
		return
	}

	model := guiApp.Preferences().IntWithFallback(CurrentModelKey, int(ollama.Llama3Dot1))
	loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
		"Using model: "+ollama.ModelName(model).String()+" with prompt: "+selectedPrompt.String()+"...")
	loadingScreen.Show()

	generated, err := ollama.AskAIWithPromptMsg(ollamaClient, ollama.ModelName(model), selectedPrompt, clippy)
	if err != nil {
		// TODO: Implement error handling, tell user to restart ollama, maybe we can restart ollama here?
		slog.Error("Failed to communicate with Ollama", "error", err)
		loadingScreen.Hide()
		return
	}
	loadingScreen.Hide()

	handleGeneratedResponse(clippy, &generated)
}

func handleAskKeyPressed() {
	err := th.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer th.Done(err)

	clippy, copiedText := copyTextToClipboard()
	if !copiedText {
		return
	}

	model := guiApp.Preferences().IntWithFallback(CurrentModelKey, int(ollama.Llama3Dot1))
	loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
		"Asking question with model: "+ollama.ModelName(model).String())
	loadingScreen.Show()

	generated, err := ollama.AskAI(ollamaClient, ollama.ModelName(model), clippy)
	if err != nil {
		slog.Error("Failed to ask AI", "error", err)
		loadingScreen.Hide()
		return
	}
	loadingScreen.Hide()

	handleGeneratedResponse(clippy, &generated)
}

func handleTranslatePressed() {
	err := th.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer func() {
		slog.Info("Done translating")
		th.Done(err)
	}()

	clippy, copiedText := copyTextToClipboard()
	if !copiedText {
		return
	}

	fromLang := guiApp.Preferences().StringWithFallback(CurrentFromLangKey, string(ollama.English))
	toLang := guiApp.Preferences().StringWithFallback(CurrentToLangKey, string(ollama.Spanish))

	model := guiApp.Preferences().IntWithFallback(CurrentModelKey, int(ollama.Llama3Dot1))
	loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
		"Translating with model: "+ollama.ModelName(model).String()+"...")
	loadingScreen.Show()

	slog.Info("Translating text", "fromLang", fromLang, "toLang", toLang)
	generated, err := ollama.AskAIToTranslate(ollamaClient, ollama.ModelName(model), clippy, ollama.Language(fromLang), ollama.Language(toLang))
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
	lastClipboardContent = sha256.Sum256([]byte(generated.Response))

	speakResponse := guiApp.Preferences().BoolWithFallback(SpeakAIResponseKey, false)
	if speakResponse {
		if speech == nil {
			initSpeech()
		}
		if speech != nil {
			speakErr := speech.Speak(generated.Response)
			if speakErr != nil {
				slog.Error("Failed to speak", "error", speakErr)
			}
		}
	}

	// Send a paste command to the operating system
	replaceText := guiApp.Preferences().BoolWithFallback(ReplaceHighlightedText, true)
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

	clippy, err := clipboard.ReadAll()
	if err != nil {
		slog.Error("Failed to read clipboard", "error", err)
		return "", false
	}
	if clippy == "" {
		slog.Info("Clipboard is empty, skipping")
		return "", false
	}

	if sha256.Sum256([]byte(clippy)) == lastClipboardContent {
		slog.Debug("Clipboard content is the same as last time", "clippy", clippy)
		return "", false
	}
	return clippy, true
}

func handleGeneratedResponse(question string, response *api.GenerateResponse) {
	slog.Debug("lastClipboardContent", "lastClipboardContent", lastClipboardContent)

	lastClipboardContent = sha256.Sum256([]byte(response.Response))
	err := clipboard.WriteAll(response.Response)
	if err != nil {
		slog.Error("Failed to write to clipboard", "error", err)
		return
	}

	speakResponse := guiApp.Preferences().BoolWithFallback(SpeakAIResponseKey, false)
	if speakResponse {
		if speech == nil {
			initSpeech()
		}
		if speech != nil {
			speakErr := speech.Speak(response.Response)
			if speakErr != nil {
				slog.Error("Failed to speak", "error", speakErr)
			}
		}
	}

	// Send a paste command to the operating system
	replaceText := guiApp.Preferences().BoolWithFallback(ReplaceHighlightedText, true)
	if replaceText {
		err = pasteCommand()
		if err != nil {
			return
		}
	}

	showPopUp := guiApp.Preferences().BoolWithFallback(ShowPopUpKey, false)
	if showPopUp {
		questionPopUp(guiApp, question, response)
	}
}
