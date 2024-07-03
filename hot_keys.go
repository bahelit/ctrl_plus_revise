package main

import (
	"log/slog"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

const (
	thinkingMsg = "Thinking......"
)

// registerHotkeys registers the hotkeys for the application
// TODO: Allow users changes the mappings
func registerHotkeys(sysTray fyne.Window) {
	hook.Register(hook.KeyDown, []string{robotgo.KeyA, robotgo.Alt, robotgo.Ctrl}, func(e hook.Event) {
		slog.Debug("ctrl-alt-a has been pressed", "event", e)
		handleAskKeyPressed(sysTray)
	})

	hook.Register(hook.KeyDown, []string{robotgo.KeyC, robotgo.Ctrl, robotgo.Shift}, func(e hook.Event) {
		slog.Debug("ctrl-shift-c has been pressed", "event", e)
		handleUserShortcutKeyPressed(sysTray)
	})

	hook.Register(hook.KeyDown, []string{robotgo.Tab, robotgo.Ctrl, robotgo.Alt}, func(e hook.Event) {
		slog.Debug("alt-ctrl-tab has been pressed", "event", e)
		handleCyclePromptKeyPressed()
		slog.Debug("Changed AI action", "prompt", selectedPrompt)
	})

	hook.Register(hook.KeyDown, []string{robotgo.KeyR, robotgo.Ctrl, robotgo.Alt}, func(e hook.Event) {
		slog.Debug("alt-ctrl-r has been pressed", "event", e)
		handleReadTextPressed(sysTray)
		slog.Debug("Reading Text aloud")
	})

	s := hook.Start()
	<-hook.Process(s)
}

func handleCyclePromptKeyPressed() {
	if int(selectedPrompt) == len(promptToText)-1 {
		selectedPrompt = CorrectGrammar
	} else {
		selectedPrompt++
	}
	err := selectedPromptBinding.Set(selectedPrompt.String())
	if err != nil {
		slog.Error("Failed to set selectedPromptBinding", "error", err)

	}
	updateDropDownMenus()
	changedPromptNotification()
	robotgo.Sleep(1)
}

func handleReadTextPressed(sysTray fyne.Window) {
	robotgo.Sleep(1)
	err := robotgo.KeyTap(robotgo.KeyC, robotgo.Ctrl)
	if err != nil {
		slog.Error("Failed to send copy command", "error", err)
	}

	speakAIResponse := guiApp.Preferences().BoolWithFallback(speakAIResponseKey, false)
	if !speakAIResponse {
		slog.Debug("Skipping reading text aloud", "speakAIResponse", speakAIResponse)
		return
	}

	clippy := sysTray.Clipboard().Content()

	_ = speech.Speak("")
	speakErr := speech.Speak(clippy)
	if speakErr != nil {
		slog.Error("Failed to speak", "error", speakErr)
	}
}

func handleUserShortcutKeyPressed(sysTray fyne.Window) {
	robotgo.Sleep(1)
	err := robotgo.KeyTap(robotgo.KeyC, robotgo.Ctrl)
	if err != nil {
		slog.Error("Failed to send copy command", "error", err)
	}

	clippy := sysTray.Clipboard().Content()

	model := guiApp.Preferences().IntWithFallback(currentModelKey, int(Llama3))

	generated, err := askAIWithPromptMsg(ollamaClient, selectedPrompt, ModelName(model), clippy)
	if err != nil {
		// TODO: Implement error handling, tell user to restart ollama, maybe we can restart ollama here?
		slog.Error("Failed to communicate with Ollama", "error", err)
		return
	}

	// Send a paste command to the operating system
	replaceText := guiApp.Preferences().BoolWithFallback(replaceHighlightedText, true)
	if replaceText {
		robotgo.TypeStr(generated.Response)
	} else {
		robotgo.TypeStr(string(clippy) + " " + generated.Response)
	}

	speakResponse := guiApp.Preferences().BoolWithFallback(speakAIResponseKey, false)
	if speakResponse {
		speakErr := speech.Speak(generated.Response)
		if speakErr != nil {
			slog.Error("Failed to speak", "error", speakErr)
		}
	}

	// Copy the generated text to the clipboard
	sysTray.Clipboard().SetContent(generated.Response)
}

func handleAskKeyPressed(sysTray fyne.Window) {
	robotgo.Sleep(1)
	err := robotgo.KeyTap(robotgo.KeyC, robotgo.Ctrl)
	if err != nil {
		slog.Error("Failed to send copy command", "error", err)

	}

	clippy := sysTray.Clipboard().Content()

	robotgo.TypeStr(thinkingMsg)
	generated, err := askAI(ollamaClient, selectedModel, clippy)
	if err != nil {
		slog.Error("Failed to ask AI", "error", err)
	}
	for i := 0; i < len(thinkingMsg); i++ {
		err = robotgo.KeyPress("BackSpace")
		if err != nil {
			slog.Error("Failed to send backspace command", "error", err)
			_ = fallbackPasteCommands
		}
	}

	// Send a paste command to the operating system
	replaceText := guiApp.Preferences().BoolWithFallback(replaceHighlightedText, true)
	if replaceText {
		robotgo.TypeStr(generated.Response)
	} else {
		robotgo.TypeStr(clippy + " " + generated.Response)
	}

	speakResponse := guiApp.Preferences().BoolWithFallback(speakAIResponseKey, false)
	if speakResponse {
		speakErr := speech.Speak(generated.Response)
		if speakErr != nil {
			slog.Error("Failed to speak", "error", speakErr)
		}
	}

	// Copy the generated text to the clipboard
	sysTray.Clipboard().SetContent(generated.Response)
}

func fallbackPasteCommands() error {
	var err error
	if runtime.GOOS == "darwin" { // TODO: Untested
		slog.Debug("Sending paste command to MacOS")
		cmd := exec.Command("cgtool", "createkeyboardevent", "command+v")
		err = cmd.Run()
	} else if runtime.GOOS == "linux" {
		slog.Debug("Sending paste command to Linux X11 with xdotool")
		cmd := exec.Command("xdotool", "key", "BackSpace", "ctrl+v") // Ctrl+V (paste) X11
		err = cmd.Run()
		if err != nil { // TODO: Untested
			slog.Debug("xdotool failed try Wayland equivalent wtype")
			cmd := exec.Command("wtype", "key", "BackSpace", "ctrl+v") // Ctrl+V (paste) Wayland
			err = cmd.Run()
		}
	} else if runtime.GOOS == "windows" { // TODO: Untested
		slog.Debug("Sending paste command to Windows")
		cmd := exec.Command("autohotkey", "sendevent,^v")
		err = cmd.Run()
	}
	if err != nil {
		slog.Error("Failed to send paste command", "error", err)
	}
	return err
}
