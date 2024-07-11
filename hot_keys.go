package main

import (
	"crypto/sha256"
	"log/slog"
	"time"

	"fyne.io/fyne/v2"
	"github.com/go-vgo/robotgo"
	"github.com/go-vgo/robotgo/clipboard"
	hook "github.com/robotn/gohook"

	"github.com/bahelit/ctrl_plus_revise/pkg/throttle"
)

const (
	thinkingMsg = "Thinking..."
)

var (
	th                    = throttle.NewThrottle(1)
	lastClipboardContent  [32]byte
	lastKeyPressTime      = time.Now()
	timeDiff              time.Duration
	waitBetweenKeyPresses = 125 * time.Millisecond
	keyPressSleep         = 125
)

// registerHotkeys registers the hotkeys for the application
// TODO: Allow users changes the mappings
func registerHotkeys(sysTray fyne.Window) {
	hook.Register(hook.KeyDown, []string{"a", "alt"}, func(e hook.Event) {
		slog.Debug("askKey has been pressed", "event", e)
		handleAskKeyPressed()
	})
	hook.Register(hook.KeyDown, []string{"c", "alt"}, func(e hook.Event) {
		slog.Debug("userShortcutKey has been pressed", "event", e)
		handleUserShortcutKeyPressed()
	})
	hook.Register(hook.KeyDown, []string{"p", "alt"}, func(e hook.Event) {
		slog.Debug("cyclePromptKey has been pressed", "event", e)
		timeDiff = time.Since(lastKeyPressTime)
		if timeDiff < waitBetweenKeyPresses {
			slog.Info("Ignoring key press", "timeDiff", timeDiff, "waitBetweenKeyPresses", waitBetweenKeyPresses)
			lastKeyPressTime = time.Now()
			return
		}
		handleCyclePromptKeyPressed()
		lastKeyPressTime = time.Now()
	})
	hook.Register(hook.KeyDown, []string{"r", "alt"}, func(e hook.Event) {
		slog.Debug("readTextKey has been pressed", "event", e)
		timeDiff = time.Since(lastKeyPressTime)
		if timeDiff < waitBetweenKeyPresses {
			slog.Info("Ignoring key press", "timeDiff", timeDiff, "waitBetweenKeyPresses", waitBetweenKeyPresses)
			lastKeyPressTime = time.Now()
			return
		}
		handleReadTextPressed(sysTray)
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

func handleReadTextPressed(sysTray fyne.Window) {
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

	clippy := sysTray.Clipboard().Content()

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

	robotgo.MilliSleep(keyPressSleep)
	err = robotgo.KeyTap(robotgo.KeyC, robotgo.CmdCtrl())
	if err != nil {
		slog.Error("Failed to send copy command", "error", err)
	}
	robotgo.MilliSleep(keyPressSleep)

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
		"Using model: "+ModelName(model).String()+" with prompt: "+selectedPrompt.String()+"...")
	loadingScreen.Show()

	generated, err := askAIWithPromptMsg(ollamaClient, selectedPrompt, ModelName(model), clippy)
	if err != nil {
		// TODO: Implement error handling, tell user to restart ollama, maybe we can restart ollama here?
		slog.Error("Failed to communicate with Ollama", "error", err)
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

	// Send a paste command to the operating system
	replaceText := guiApp.Preferences().BoolWithFallback(replaceHighlightedText, true)
	if !replaceText {
		_ = robotgo.KeyPress(robotgo.Right, robotgo.Space)
	}
	robotgo.TypeStr(generated.Response)

	speakResponse := guiApp.Preferences().BoolWithFallback(speakAIResponseKey, false)
	if speakResponse {
		speakErr := speech.Speak(generated.Response)
		if speakErr != nil {
			slog.Error("Failed to speak", "error", speakErr)
		}
	}
}

func handleAskKeyPressed() {
	err := th.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer th.Done(err)

	robotgo.MilliSleep(keyPressSleep)
	err = robotgo.KeyTap(robotgo.KeyC, robotgo.CmdCtrl())
	if err != nil {
		slog.Error("Failed to send copy command", "error", err)
	}
	robotgo.MilliSleep(keyPressSleep)

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
		"Asking question with model: "+ModelName(model).String()+"...")
	loadingScreen.Show()

	generated, err := askAI(ollamaClient, ModelName(model), clippy)
	if err != nil {
		slog.Error("Failed to ask AI", "error", err)
	}
	loadingScreen.Hide()

	lastClipboardContent = sha256.Sum256([]byte(generated.Response))
	err = clipboard.WriteAll(generated.Response)
	if err != nil {
		slog.Error("Failed to write to clipboard", "error", err)
		return
	}

	// Send a paste command to the operating system
	replaceText := guiApp.Preferences().BoolWithFallback(replaceHighlightedText, true)
	if !replaceText {
		_ = robotgo.KeyPress(robotgo.Right, robotgo.Space)
	}
	robotgo.TypeStr(generated.Response)

	speakResponse := guiApp.Preferences().BoolWithFallback(speakAIResponseKey, false)
	if speakResponse {
		speakErr := speech.Speak(generated.Response)
		if speakErr != nil {
			slog.Error("Failed to speak", "error", speakErr)
		}
	}
}
