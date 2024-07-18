package main

import (
	"crypto/sha256"
	"log/slog"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"

	"github.com/bahelit/ctrl_plus_revise/pkg/clipboard"
	"github.com/bahelit/ctrl_plus_revise/pkg/throttle"
)

const (
	thinkingMsg = "Thinking..."
)

var (
	th                    = throttle.NewThrottle(1)
	lastClipboardContent  [32]byte
	lastKeyPressTime      = time.Now()
	waitBetweenKeyPresses = 1 * time.Second
	keyPressSleep         = 250
	firstRun              = true
)

// registerHotkeys registers the hotkeys for the application
// TODO: Allow users changes the mappings
func registerHotkeys() {
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
	slog.Info("Registered hotkeys")

	s := hook.Start()
	defer hook.End()
	<-hook.Process(s)

	// Debug code
	/*for ev := range s {
		slog.Info("Event caught", "event", ev)
	}*/
}

func handleCyclePromptKeyPressed() {
	err := th.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer th.Done(err)
	if int(selectedPrompt) == len(promptToText)-1 {
		selectedPrompt = CorrectGrammar
	} else {
		selectedPrompt++
	}
	err = selectedPromptBinding.Set(selectedPrompt.String())
	if err != nil {
		slog.Error("Failed to set selectedPromptBinding", "error", err)

	}
	updateDropDownMenus()
	changedPromptNotification()
	time.Sleep(1 * time.Second)
}

func handleReadTextPressed() {
	err := th.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer th.Done(err)
	speakAIResponse := guiApp.Preferences().BoolWithFallback(speakAIResponseKey, false)
	if !speakAIResponse {
		slog.Info("Reading text aloud is disabled", "speakAIResponse", speakAIResponse)
		return
	}

	err = copyCommand()
	if err != nil {
		_ = speech.Speak("Failed to copy text")
		return
	}

	clippy, err := clipboard.ReadAll()
	if err != nil {
		slog.Error("Failed to read clipboard", "error", err)
		_ = speech.Speak("Failed to read clipboard")
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

	err = copyCommand()
	if err != nil {
		return
	}

	clippy, err := clipboard.ReadAll()
	if err != nil {
		slog.Error("Failed to read clipboard", "error", err)
		return
	}
	if clippy == "" {
		slog.Info("Clipboard is empty, skipping")
		return
	}
	if sha256.Sum256([]byte(clippy)) == lastClipboardContent {
		slog.Info("Clipboard content is the same as last, skipping")
		return
	}

	model := guiApp.Preferences().IntWithFallback(currentModelKey, int(Llama3))
	loadingScreen := loadingScreenWithMessage(thinkingMsg,
		"Using model: "+ModelName(model).String()+" with prompt: "+selectedPrompt.String()+"...")
	loadingScreen.Show()

	generated, err := askAIWithPromptMsg(ollamaClient, selectedPrompt, ModelName(model), clippy)
	if err != nil {
		// TODO: Implement error handling, tell user to restart ollama, maybe we can restart ollama here?
		slog.Error("Failed to communicate with Ollama", "error", err)
		loadingScreen.Hide()
		return
	}
	loadingScreen.Hide()

	slog.Debug("lastClipboardContent", "lastClipboardContent", lastClipboardContent)

	lastClipboardContent = sha256.Sum256([]byte(generated.Response))
	err = clipboard.WriteAll(generated.Response)
	if err != nil {
		slog.Error("Failed to write to clipboard", "error", err)
		return
	}

	speakResponse := guiApp.Preferences().BoolWithFallback(speakAIResponseKey, false)
	if speakResponse {
		speakErr := speech.Speak(generated.Response)
		if speakErr != nil {
			slog.Error("Failed to speak", "error", speakErr)
		}
	}

	// Send a paste command to the operating system
	replaceText := guiApp.Preferences().BoolWithFallback(replaceHighlightedText, true)
	if replaceText {
		err = pasteCommand()
		if err != nil {
			return
		}
	}
}

func handleAskKeyPressed() {
	err := th.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer th.Done(err)

	err = copyCommand()
	if err != nil {
		_ = speech.Speak("Failed to copy text")
		return
	}

	clippy, err := clipboard.ReadAll()
	if err != nil {
		slog.Error("Failed to read clipboard", "error", err)
		return
	}

	if sha256.Sum256([]byte(clippy)) == lastClipboardContent {
		slog.Debug("Clipboard content is the same as last time", "clippy", clippy)
		return
	}

	model := guiApp.Preferences().IntWithFallback(currentModelKey, int(Llama3))
	loadingScreen := loadingScreenWithMessage(thinkingMsg,
		"Asking question with model: "+ModelName(model).String())
	loadingScreen.Show()

	generated, err := askAI(ollamaClient, ModelName(model), clippy)
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

	// Send a paste command to the operating system
	replaceText := guiApp.Preferences().BoolWithFallback(replaceHighlightedText, true)
	if replaceText {
		err = pasteCommand()
		if err != nil {
			return
		}
	}

	speakResponse := guiApp.Preferences().BoolWithFallback(speakAIResponseKey, false)
	if speakResponse {
		speakErr := speech.Speak(generated.Response)
		if speakErr != nil {
			slog.Error("Failed to speak", "error", speakErr)
		}
	}
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

	err = copyCommand()
	if err != nil {
		_ = speech.Speak("Failed to copy text")
		return
	}

	clippy, err := clipboard.ReadAll()
	if err != nil {
		slog.Error("Failed to read clipboard", "error", err)
		return
	}

	if sha256.Sum256([]byte(clippy)) == lastClipboardContent {
		slog.Info("Clipboard content is the same as last time", "clippy", clippy)
		return
	}

	fromLang := guiApp.Preferences().StringWithFallback(currentFromLangKey, string(English))
	toLang := guiApp.Preferences().StringWithFallback(currentToLangKey, string(Spanish))

	model := guiApp.Preferences().IntWithFallback(currentModelKey, int(Llama3))
	loadingScreen := loadingScreenWithMessage(thinkingMsg,
		"Translating with model: "+ModelName(model).String()+"...")
	loadingScreen.Show()

	slog.Info("Translating text", "fromLang", fromLang, "toLang", toLang)
	generated, err := askAIToTranslate(ollamaClient, ModelName(model), clippy, Language(fromLang), Language(toLang))
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

	speakResponse := guiApp.Preferences().BoolWithFallback(speakAIResponseKey, false)
	if speakResponse {
		speakErr := speech.Speak(generated.Response)
		if speakErr != nil {
			slog.Error("Failed to speak", "error", speakErr)
		}
	}

	// Send a paste command to the operating system
	replaceText := guiApp.Preferences().BoolWithFallback(replaceHighlightedText, true)
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
