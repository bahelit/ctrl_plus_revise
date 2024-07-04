package main

import (
	"log/slog"
	"os/exec"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"github.com/MarinX/keylogger"
	"github.com/atotto/clipboard"
	"github.com/micmonay/keybd_event"
	hook "github.com/robotn/gohook"
)

const (
	thinkingMsg = "Thinking......"
)

// registerHotkeys registers the hotkeys for the application
// TODO: Allow users changes the mappings
func registerHotkeys(sysTray fyne.Window) {
	keyboard := keylogger.FindKeyboardDevice()
	// init keylogger with keyboard
	k, err := keylogger.New(keyboard)
	if err != nil {
		slog.Error("Failed to create keylogger", "error", err)
		return
	}
	events := k.Read()

	var hitCrl, hitAlt, hitTab, hitR, hitC, hitA bool
	// range of events
	for e := range events {
		switch e.Type {
		// EvKey is used to describe state changes of keyboards, buttons, or other key-like devices.
		// check the input_event.go for more events
		case keylogger.EvKey:
			// if the state of key is pressed
			if e.KeyPress() {
				slog.Debug("[event] press key ", "key", e.KeyString())
				if e.KeyString() == "L_CTRL" {
					hitCrl = true
				}
				if e.KeyString() == "L_ALT" {
					hitAlt = true
				}
				if e.KeyString() == "TAB" {
					hitTab = true
				}
				if e.KeyString() == "R" {
					hitR = true
				}
				if e.KeyString() == "C" {
					hitC = true
				}
				if e.KeyString() == "A" {
					hitA = true
				}
			}
			// if the state of key is released
			if e.KeyRelease() {
				slog.Debug("[event] release key ", "key", e.KeyString())
				if e.KeyString() == "L_CTRL" {
					hitCrl = false
				}
				if e.KeyString() == "L_ALT" {
					hitAlt = false
				}
				if e.KeyString() == "TAB" {
					hitTab = false
				}
				if e.KeyString() == "R" {
					hitR = false
				}
				if e.KeyString() == "C" {
					hitC = false
				}
				if e.KeyString() == "A" {
					hitA = false
				}
			}
			if hitCrl && hitAlt && hitC {
				slog.Debug("Ctrl-Alt-C has been pressed")
				time.Sleep(500 * time.Millisecond)
				hitCrl, hitAlt, hitC = false, false, false
				handleUserShortcutKeyPressed(sysTray)
			}
			if hitCrl && hitAlt && hitA {
				slog.Debug("Ctrl-Alt-A has been pressed")
				time.Sleep(500 * time.Millisecond)
				hitCrl, hitAlt, hitA = false, false, false
				handleAskKeyPressed(sysTray)
			}
			if hitCrl && hitAlt && hitR {
				slog.Debug("Ctrl-Alt-R has been pressed")
				time.Sleep(500 * time.Millisecond)
				hitCrl, hitAlt, hitR = false, false, false
				handleReadTextPressed(sysTray)
			}
			if hitCrl && hitAlt && hitTab {
				slog.Debug("Ctrl-Alt-TAB has been pressed")
				time.Sleep(500 * time.Millisecond)
				hitCrl, hitAlt, hitTab = false, false, false
				handleCyclePromptKeyPressed()
			}
			break
		}
	}

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
	time.Sleep(1 * time.Second)
}

func handleReadTextPressed(sysTray fyne.Window) {
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

func handleUserShortcutKeyPressed(sysTray fyne.Window) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		slog.Error("Failed to create key binding", "error", err)
		panic(err)
	}
	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}
	// Select keys to be pressed
	kb.SetKeys(keybd_event.VK_C)
	// Set shift to be pressed
	kb.HasCTRLR(true)
	// Press the selected keys
	err = kb.Launching()
	if err != nil {
		panic(err)
	}

	clippy, err := clipboard.ReadAll()
	if err != nil {
		slog.Error("Failed to read clipboard", "error", err)
		return
	}
	model := guiApp.Preferences().IntWithFallback(currentModelKey, int(Llama3))

	//kb.Clear()
	//kb.SetKeys(keybd_event.VK_T, keybd_event.VK_H, keybd_event.VK_I, keybd_event.VK_N, keybd_event.VK_K,
	//	keybd_event.VK_I, keybd_event.VK_N, keybd_event.VK_G, keybd_event.VK_DOT, keybd_event.VK_DOT, keybd_event.VK_DOT, keybd_event.VK_DOT)
	//// Press the selected keys
	//err = kb.Launching()
	//if err != nil {
	//	slog.Error("Failed to send thinking message", "error", err)
	//	return
	//}

	generated, err := askAIWithPromptMsg(ollamaClient, selectedPrompt, ModelName(model), string(clippy))
	if err != nil {
		// TODO: Implement error handling, tell user to restart ollama, maybe we can restart ollama here?
		slog.Error("Failed to communicate with Ollama", "error", err)
		return
	}
	//for i := 0; i < len(thinkingMsg); i++ {
	//	kb.Clear()
	//	kb.SetKeys(keybd_event.VK_BACKSPACE)
	//	err = kb.Launching()
	//	if err != nil {
	//		slog.Error("Failed to send backspace command", "error", err)
	//		return
	//	}
	//}

	err = clipboard.WriteAll(generated.Response)
	if err != nil {
		slog.Error("Failed to write to clipboard", "error", err)
		return
	}
	// Send a paste command to the operating system
	replaceText := guiApp.Preferences().BoolWithFallback(replaceHighlightedText, true)
	if replaceText {
		kb.Clear()
		kb.SetKeys(keybd_event.VK_V)
		kb.HasCTRLR(true)
		err = kb.Launching()
		if err != nil {
			slog.Error("Failed to send paste command", "error", err)
		}
	} else {
		kb.Clear()
		kb.SetKeys(keybd_event.VK_RIGHT, keybd_event.VK_SPACE)
		err = kb.Launching()
		if err != nil {
			slog.Error("Failed to send paste command", "error", err)
		}
		kb.Clear()
		kb.SetKeys(keybd_event.VK_V)
		kb.HasCTRLR(true)
		err = kb.Launching()
		if err != nil {
			slog.Error("Failed to send paste command", "error", err)
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

func handleAskKeyPressed(sysTray fyne.Window) {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		slog.Error("Failed to create key binding", "error", err)
		panic(err)
	}
	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}
	// Select keys to be pressed
	kb.SetKeys(keybd_event.VK_C)
	// Set shift to be pressed
	kb.HasCTRLR(true)

	// Press the selected keys
	err = kb.Launching()
	if err != nil {
		panic(err)
	}

	clippy, err := clipboard.ReadAll()
	if err != nil {
		slog.Error("Failed to read clipboard", "error", err)
		return
	}

	//kb.Clear()
	//kb.SetKeys(keybd_event.VK_T, keybd_event.VK_H, keybd_event.VK_I, keybd_event.VK_N, keybd_event.VK_K,
	//	keybd_event.VK_I, keybd_event.VK_N, keybd_event.VK_G, keybd_event.VK_DOT, keybd_event.VK_DOT, keybd_event.VK_DOT, keybd_event.VK_DOT)
	//// Press the selected keys
	//err = kb.Launching()
	//if err != nil {
	//	slog.Error("Failed to send thinking message", "error", err)
	//	return
	//}
	generated, err := askAI(ollamaClient, selectedModel, string(clippy))
	if err != nil {
		slog.Error("Failed to ask AI", "error", err)
	}
	//for i := 0; i < len(thinkingMsg); i++ {
	//	kb.Clear()
	//	kb.SetKeys(keybd_event.VK_BACKSPACE)
	//	err = kb.Launching()
	//	if err != nil {
	//		slog.Error("Failed to send backspace command", "error", err)
	//		return
	//	}
	//}

	err = clipboard.WriteAll(generated.Response)
	if err != nil {
		slog.Error("Failed to write to clipboard", "error", err)
		return
	}
	// Send a paste command to the operating system
	replaceText := guiApp.Preferences().BoolWithFallback(replaceHighlightedText, true)
	if replaceText {
		kb.Clear()
		kb.SetKeys(keybd_event.VK_V)
		kb.HasCTRLR(true)
		err = kb.Launching()
		if err != nil {
			slog.Error("Failed to send paste command", "error", err)
		}
	} else {
		kb.Clear()
		kb.SetKeys(keybd_event.VK_RIGHT, keybd_event.VK_SPACE)
		err = kb.Launching()
		if err != nil {
			slog.Error("Failed to send paste command", "error", err)
		}
		kb.Clear()
		kb.SetKeys(keybd_event.VK_V)
		kb.HasCTRLR(true)
		err = kb.Launching()
		if err != nil {
			slog.Error("Failed to send paste command", "error", err)
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
