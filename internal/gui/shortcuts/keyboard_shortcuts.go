package shortcuts

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"

	"github.com/bahelit/ctrl_plus_revise/internal/config"
)

func modifierKeys() []string {
	return []string{
		robotgo.Alt,
		robotgo.Ctrl,
		robotgo.Cmd,
		robotgo.Enter,
		robotgo.Shift,
	}
}

func normalKeys() []string {
	return []string{
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "0",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
		"k", "l", "m", "n", "o", "p", "q", "r", "s", "t",
		"u", "v", "w", "x", "y", "z"}
}

type KeyType int

const (
	ModifierKey1 KeyType = iota
	ModifierKey2
	NormalKey
)

type KeyAction int

const (
	AskQuestion KeyAction = iota
	CtrlRevise
	Translate

	EmptySelection string = "Not Used"
)

//nolint:gocyclo // This function is a switch statement and is not complex
func setKeyBoardShortcutKey(guiApp fyne.App, action KeyAction, key KeyType, value string) {
	var empty string = EmptySelection
	switch action {
	case AskQuestion:
		switch key {
		case ModifierKey1:
			AskKey.ModifierKey1 = value
		case ModifierKey2:
			if value == EmptySelection {
				AskKey.ModifierKey2 = &empty
				break
			}
			AskKey.ModifierKey2 = &value
		case NormalKey:
			AskKey.Key = value
		}
		guiApp.Preferences().SetStringList(config.AskAIKeyboardShortcut, GetAskKeys())
	case CtrlRevise:
		switch key {
		case ModifierKey1:
			CtrlReviseKey.ModifierKey1 = value
		case ModifierKey2:
			if value == EmptySelection {
				CtrlReviseKey.ModifierKey2 = &empty
				break
			}
			CtrlReviseKey.ModifierKey2 = &value
		case NormalKey:
			CtrlReviseKey.Key = value
		}
		guiApp.Preferences().SetStringList(config.CtrlReviseKeyboardShortcut, GetCtrlReviseKeys())
	case Translate:
		switch key {
		case ModifierKey1:
			TranslateKey.ModifierKey1 = value
		case ModifierKey2:
			if value == EmptySelection {
				TranslateKey.ModifierKey2 = &empty
				break
			}
			TranslateKey.ModifierKey2 = &value
		case NormalKey:
			TranslateKey.Key = value
		}
		guiApp.Preferences().SetStringList(config.TranslateKeyboardShortcut, GetTranslateKeys())
	default:
		slog.Error("Unknown action", "action", action)
	}
}

//nolint:gocyclo // This function is not too complex
func keyboardModifierButtonsDropDown(guiApp fyne.App, action KeyAction, key KeyType) *widget.Select {
	modKeys := modifierKeys()
	if key == ModifierKey2 {
		modKeys = append(modKeys, EmptySelection)
	} else if key == NormalKey {
		modKeys = normalKeys()
	}
	combo := widget.NewSelect(
		modKeys,
		func(value string) {
			setKeyBoardShortcutKey(guiApp, action, key, value)
		})
	switch key {
	case ModifierKey1:
		switch action {
		case AskQuestion:
			combo.SetSelected(AskKey.ModifierKey1)
		case CtrlRevise:
			combo.SetSelected(CtrlReviseKey.ModifierKey1)
		case Translate:
			combo.SetSelected(TranslateKey.ModifierKey1)
		default:
			slog.Error("Invalid action", "action", action)
		}
	case ModifierKey2:
		switch action {
		case AskQuestion:
			if AskKey.ModifierKey2 != nil {
				combo.SetSelected(*AskKey.ModifierKey2)
			}
		case CtrlRevise:
			if CtrlReviseKey.ModifierKey2 != nil {
				combo.SetSelected(*CtrlReviseKey.ModifierKey2)
			}
		case Translate:
			if TranslateKey.ModifierKey2 != nil {
				combo.SetSelected(*TranslateKey.ModifierKey2)
			}
		default:
			slog.Error("Invalid action", "action", action)
		}
	case NormalKey:
		switch action {
		case AskQuestion:
			combo.SetSelected(AskKey.Key)
		case CtrlRevise:
			combo.SetSelected(CtrlReviseKey.Key)
		case Translate:
			combo.SetSelected(TranslateKey.Key)
		default:
			slog.Error("Invalid action", "action", action)
		}
	default:
		slog.Error("Invalid key type", "key", key)
	}

	return combo
}
