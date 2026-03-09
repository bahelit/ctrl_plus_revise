package shortcuts

import (
	"crypto/sha256"
	"log/slog"
	"time"

	"github.com/go-vgo/robotgo"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/ollama/ollama/api"
	hook "github.com/robotn/gohook"

	"github.com/bahelit/ctrl_plus_revise/internal/config"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
	"github.com/bahelit/ctrl_plus_revise/pkg/clipboard"
	"github.com/bahelit/ctrl_plus_revise/pkg/throttle"
)

// KeyboardShortcut describes a hotkey combination.
type KeyboardShortcut struct {
	ModifierKey1 string
	ModifierKey2 *string
	Key          string
}

// NotifyFunc is called when a keyboard action completes to surface a
// notification to the user (e.g. via system tray or Wails event).
type NotifyFunc func(title, content string)

const (
	waitBetweenKeyPresses = 1 * time.Second
	keyPressSleep         = 250
)

// Manager owns all global keyboard hook state and callback dependencies.
type Manager struct {
	cfg            *config.AppConfig
	ollamaClient   *api.Client
	notify         NotifyFunc
	speech         *htgotts.Speech
	throttle       *throttle.Throttle
	lastKeyPress   time.Time
	lastClipHash   [32]byte
	firstRun       bool
	systemHook     chan hook.Event
	selectedPrompt ollama.PromptMsg
	selectedModel  ollama.ModelName
}

// NewManager creates a ready-to-use Manager.  Call Start() in a goroutine to
// begin listening for hotkeys.
func NewManager(cfg *config.AppConfig, client *api.Client, notify NotifyFunc) *Manager {
	m := &Manager{
		cfg:            cfg,
		ollamaClient:   client,
		notify:         notify,
		throttle:       throttle.NewThrottle(1),
		lastKeyPress:   time.Now(),
		firstRun:       true,
		selectedPrompt: ollama.CorrectGrammar,
		selectedModel:  ollama.ModelName(cfg.CurrentModel),
	}
	if cfg.SpeakAIResponse {
		m.speech = &htgotts.Speech{Folder: "audio", Language: "en"}
	}
	return m
}

// Start registers all hotkeys and blocks until the hook is stopped.
// It must be called in its own goroutine.
func (m *Manager) Start() {
	m.systemHook = m.registerHotkeys()
	<-hook.Process(m.systemHook)
}

// Stop tears down the global keyboard hook.
func (m *Manager) Stop() {
	hook.End()
}

// GetAskKeys returns the key slice for the Ask AI shortcut.
func (m *Manager) GetAskKeys() []string {
	keys := m.cfg.AskAIShortcut
	if len(keys) == 0 {
		return []string{"alt", "a"}
	}
	return keys
}

// GetCtrlReviseKeys returns the key slice for the Ctrl+Revise shortcut.
func (m *Manager) GetCtrlReviseKeys() []string {
	keys := m.cfg.CtrlReviseShortcut
	if len(keys) == 0 {
		return []string{"alt", "c"}
	}
	return keys
}

// GetTranslateKeys returns the key slice for the Translate shortcut.
func (m *Manager) GetTranslateKeys() []string {
	keys := m.cfg.TranslateShortcut
	if len(keys) == 0 {
		return []string{"alt", "t"}
	}
	return keys
}

// registerHotkeys wires up all global hotkeys and returns the hook channel.
func (m *Manager) registerHotkeys() chan hook.Event {
	hook.Register(hook.KeyDown, m.GetAskKeys(), func(e hook.Event) {
		slog.Debug("AskKey has been pressed", "event", e)
		if time.Since(m.lastKeyPress) < waitBetweenKeyPresses {
			slog.Info("Ignoring key press", "waitBetweenKeyPresses", waitBetweenKeyPresses)
			m.lastKeyPress = time.Now()
			return
		}
		m.lastKeyPress = time.Now()
		m.handleAskKeyPressed()
		m.lastKeyPress = time.Now()
	})

	hook.Register(hook.KeyDown, m.GetCtrlReviseKeys(), func(e hook.Event) {
		slog.Debug("CtrlReviseKey has been pressed", "event", e)
		if time.Since(m.lastKeyPress) < waitBetweenKeyPresses {
			slog.Info("Ignoring key press", "waitBetweenKeyPresses", waitBetweenKeyPresses)
			return
		}
		m.lastKeyPress = time.Now()
		m.handleUserShortcutKeyPressed()
		m.lastKeyPress = time.Now()
	})

	hook.Register(hook.KeyDown, []string{"p", "alt"}, func(e hook.Event) {
		slog.Debug("cyclePromptKey has been pressed", "event", e)
		if time.Since(m.lastKeyPress) < waitBetweenKeyPresses {
			slog.Info("Ignoring key press", "waitBetweenKeyPresses", waitBetweenKeyPresses)
			m.lastKeyPress = time.Now()
			return
		}
		m.lastKeyPress = time.Now()
		m.handleCyclePromptKeyPressed()
		m.lastKeyPress = time.Now()
	})

	hook.Register(hook.KeyDown, []string{"r", "alt"}, func(e hook.Event) {
		slog.Debug("readTextKey has been pressed", "event", e)
		if time.Since(m.lastKeyPress) < waitBetweenKeyPresses {
			slog.Info("Ignoring key press", "waitBetweenKeyPresses", waitBetweenKeyPresses)
			m.lastKeyPress = time.Now()
			return
		}
		m.lastKeyPress = time.Now()
		m.handleReadTextPressed()
		m.lastKeyPress = time.Now()
	})

	hook.Register(hook.KeyDown, m.GetTranslateKeys(), func(e hook.Event) {
		slog.Debug("translateTextKey has been pressed", "event", e)
		if time.Since(m.lastKeyPress) < waitBetweenKeyPresses {
			slog.Info("Ignoring key press", "waitBetweenKeyPresses", waitBetweenKeyPresses)
			m.lastKeyPress = time.Now()
			return
		}
		m.lastKeyPress = time.Now()
		m.handleTranslatePressed()
		m.lastKeyPress = time.Now()
	})

	slog.Info("Registered hotkeys")
	return hook.Start()
}

// handleCyclePromptKeyPressed advances the selected prompt to the next option.
func (m *Manager) handleCyclePromptKeyPressed() {
	err := m.throttle.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer m.throttle.Done(err)

	if int(m.selectedPrompt) == len(ollama.PromptToText)-1 {
		m.selectedPrompt = ollama.CorrectGrammar
	} else {
		m.selectedPrompt++
	}

	if m.notify != nil {
		m.notify("AI Action Changed", "AI Action has been changed to:\n"+m.selectedPrompt.String())
	}
	time.Sleep(1 * time.Second)
}

// handleReadTextPressed reads the currently selected text aloud using TTS.
func (m *Manager) handleReadTextPressed() {
	err := m.throttle.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer m.throttle.Done(err)

	if !m.cfg.SpeakAIResponse {
		slog.Info("Reading text aloud is disabled", "speakAIResponse", m.cfg.SpeakAIResponse)
		return
	}

	if err := copyCommand(m); err != nil {
		if m.speech != nil {
			_ = m.speech.Speak("Failed to copy text")
		}
		return
	}

	clip, err := clipboard.ReadAll()
	if err != nil {
		slog.Error("Failed to read clipboard", "error", err)
		if m.speech != nil {
			_ = m.speech.Speak("Failed to read clipboard")
		}
		return
	}

	if m.speech != nil {
		_ = m.speech.Speak("")
		if speakErr := m.speech.Speak(clip); speakErr != nil {
			slog.Error("Failed to speak", "error", speakErr)
		}
	}
}

// handleUserShortcutKeyPressed copies the selected text, sends it to Ollama
// with the currently selected prompt, and pastes the result back.
func (m *Manager) handleUserShortcutKeyPressed() {
	err := m.throttle.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer m.throttle.Done(err)

	clip, copiedText := copyTextToClipboard(m)
	if !copiedText {
		return
	}

	generated, err := ollama.AskAIWithPromptMsg(m.selectedModel, m.ollamaClient, m.selectedPrompt, clip)
	if err != nil {
		slog.Error("Failed to communicate with Ollama", "error", err)
		return
	}

	m.handleGeneratedResponse(clip, &generated)
}

// handleAskKeyPressed copies the selected text, sends it as a freeform
// question to Ollama, and pastes the result back.
func (m *Manager) handleAskKeyPressed() {
	err := m.throttle.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer m.throttle.Done(err)

	clip, copiedText := copyTextToClipboard(m)
	if !copiedText {
		return
	}

	generated, err := ollama.AskAI(m.selectedModel, m.ollamaClient, clip)
	if err != nil {
		slog.Error("Failed to ask AI", "error", err)
		return
	}

	m.handleGeneratedResponse(clip, &generated)
}

// handleTranslatePressed copies the selected text and translates it using the
// configured language pair.
func (m *Manager) handleTranslatePressed() {
	err := m.throttle.Do()
	if err != nil {
		slog.Error("Failed to create throttle", "error", err)
	}
	defer func() {
		slog.Info("Done translating")
		m.throttle.Done(err)
	}()

	clip, copiedText := copyTextToClipboard(m)
	if !copiedText {
		return
	}

	fromLang := m.cfg.CurrentFromLang
	toLang := m.cfg.CurrentToLang

	slog.Info("Translating text", "fromLang", fromLang, "toLang", toLang)
	generated, err := ollama.AskAIToTranslate(m.selectedModel, m.ollamaClient, clip, ollama.Language(fromLang), ollama.Language(toLang))
	if err != nil {
		slog.Error("Failed to ask AI", "error", err)
		return
	}

	if err = clipboard.WriteAll(generated.Response); err != nil {
		slog.Error("Failed to write to clipboard", "error", err)
		return
	}
	m.lastClipHash = sha256.Sum256([]byte(generated.Response))

	if m.cfg.ReplaceHighlightedText {
		if err = pasteCommand(); err != nil {
			return
		}
	}
}

// copyCommand sends a Ctrl+C (or Cmd+C on macOS) keyboard event.
func copyCommand(m *Manager) error {
	robotgo.KeySleep = 100

	if m.firstRun {
		m.firstRun = false
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

// pasteCommand sends a Ctrl+V (or Cmd+V on macOS) keyboard event.
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

// copyTextToClipboard triggers a copy and reads the clipboard contents.
// Returns ("", false) if the clipboard is empty or unchanged since last call.
func copyTextToClipboard(m *Manager) (string, bool) {
	if err := copyCommand(m); err != nil {
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

	if sha256.Sum256([]byte(clip)) == m.lastClipHash {
		slog.Debug("Clipboard content is the same as last time", "clippy", clip)
		return "", false
	}
	return clip, true
}

// handleGeneratedResponse writes the AI response to the clipboard and
// optionally pastes it into the active application.
func (m *Manager) handleGeneratedResponse(question string, response *api.GenerateResponse) {
	slog.Debug("LastClipboardContent", "LastClipboardContent", m.lastClipHash)

	m.lastClipHash = sha256.Sum256([]byte(response.Response))
	if err := clipboard.WriteAll(response.Response); err != nil {
		slog.Error("Failed to write to clipboard", "error", err)
		return
	}

	if m.cfg.ReplaceHighlightedText {
		if err := pasteCommand(); err != nil {
			return
		}
	}

	if m.cfg.ShowPopUp && m.notify != nil {
		m.notify("AI Response", response.Response)
	}
}
