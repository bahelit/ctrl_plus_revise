package main

import (
	"log/slog"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"golang.design/x/clipboard"
)

const (
	thinkingMsg = "Thinking......"
)

// registerHotkeys registers the hotkeys for the application
// TODO: Allow users changes the mappings
func registerHotkeys(sysTray fyne.Window) {
	hook.Register(hook.KeyDown, []string{robotgo.KeyA, robotgo.Alt, robotgo.Ctrl}, func(e hook.Event) {
		slog.Info("ctrl-alt-a has been pressed", "event", e)
		handleAskKeyPressed(sysTray)
	})

	hook.Register(hook.KeyDown, []string{robotgo.KeyC, robotgo.Ctrl, robotgo.Shift}, func(e hook.Event) {
		slog.Info("ctrl-shift-c has been pressed", "event", e)
		handleUserShortcutKeyPressed(sysTray)
	})

	// FIXME: This sysTray is not updating the selectedPrompt.
	hook.Register(hook.KeyDown, []string{robotgo.Tab, robotgo.Ctrl, robotgo.Alt}, func(e hook.Event) {
		slog.Info("alt-ctrl-tab has been pressed", "event", e)
		handleCyclePromptKeyPressed()
		slog.Info("Selected prompt", "prompt", selectedPrompt)
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
	changedPromptNotification()
	robotgo.Sleep(1)
}

func handleUserShortcutKeyPressed(sysTray fyne.Window) {
	robotgo.Sleep(1)
	err := robotgo.KeyPress(robotgo.KeyC, robotgo.Ctrl)
	if err != nil {
		slog.Error("Failed to send copy command", "error", err)
	}

	clippy := clipboard.Read(clipboard.FmtText)
	if clippy == nil {
		slog.Info("Clipboard is empty")
		return
	}

	replaceText := guiApp.Preferences().BoolWithFallback(replaceHighlightedText, true)

	generated, err := askAIWithPromptMsg(ollamaClient, selectedPrompt, string(clippy))
	if err != nil {
		// TODO: Implement error handling, tell user to restart ollama, maybe we can restart ollama here?
		slog.Error("Failed to communicate with Ollama", "error", err)
		return
	}

	// Send a paste command to the operating system
	if replaceText {
		robotgo.TypeStr(generated.Response)
	} else {
		robotgo.TypeStr(string(clippy) + " " + generated.Response)
	}

	// Copy the generated text to the clipboard
	sysTray.Clipboard().SetContent(generated.Response)
}

func handleAskKeyPressed(sysTray fyne.Window) {
	robotgo.Sleep(1)
	err := robotgo.KeyPress(robotgo.KeyC, robotgo.Ctrl)
	if err != nil {
		slog.Error("Failed to send copy command", "error", err)

	}

	clippy := clipboard.Read(clipboard.FmtText)
	if clippy == nil {
		slog.Info("Clipboard is empty")
		return
	}

	replaceText := guiApp.Preferences().BoolWithFallback(replaceHighlightedText, true)

	robotgo.TypeStr(thinkingMsg)
	generated, err := askAI(ollamaClient, string(clippy))
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
	if replaceText {
		robotgo.TypeStr(generated.Response)
	} else {
		robotgo.TypeStr(string(clippy) + " " + generated.Response)
	}

	// Copy the generated text to the clipboard
	sysTray.Clipboard().SetContent(generated.Response)
}

func fallbackPasteCommands() error {
	var err error
	if runtime.GOOS == "darwin" { // TODO: Untested
		slog.Info("Sending paste command to MacOS")
		cmd := exec.Command("cgtool", "createkeyboardevent", "command+v")
		err = cmd.Run()
	} else if runtime.GOOS == "linux" {
		slog.Info("Sending paste command to Linux X11 with xdotool")
		cmd := exec.Command("xdotool", "key", "BackSpace", "ctrl+v") // Ctrl+V (paste) X11
		err = cmd.Run()
		if err != nil { // TODO: Untested
			slog.Info("xdotool failed try Wayland equivalent wtype")
			cmd := exec.Command("wtype", "key", "BackSpace", "ctrl+v") // Ctrl+V (paste) Wayland
			err = cmd.Run()
		}
	} else if runtime.GOOS == "windows" { // TODO: Untested
		slog.Info("Sending paste command to Windows")
		cmd := exec.Command("autohotkey", "sendevent,^v")
		err = cmd.Run()
	}
	if err != nil {
		slog.Error("Failed to send paste command", "error", err)
	}
	return err
}
