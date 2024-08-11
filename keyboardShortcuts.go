package main

import (
	"log/slog"

	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
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

//nolint:gocyclo // This function is not too complex
//nolint:funlen // This function is not too long
func keyboardModifierButtonsDropDown(action KeyAction, key KeyType) *widget.Select {
	modKeys := modifierKeys()
	if key == ModifierKey2 {
		modKeys = append(modKeys, EmptySelection)
	} else if key == NormalKey {
		modKeys = normalKeys()
	}
	combo := widget.NewSelect(
		modKeys,
		func(value string) {
			setKeyBoardShortcutKey(action, key, value)
		})
	switch key {
	case ModifierKey1:
		switch action {
		case AskQuestion:
			combo.SetSelected(askKey.ModifierKey1)
		case CtrlRevise:
			combo.SetSelected(ctrlReviseKey.ModifierKey1)
		case Translate:
			combo.SetSelected(translateKey.ModifierKey1)
		default:
			slog.Error("Invalid action", "action", action)
		}
	case ModifierKey2:
		switch action {
		case AskQuestion:
			if askKey.ModifierKey2 != nil {
				combo.SetSelected(*askKey.ModifierKey2)
			}
		case CtrlRevise:
			if ctrlReviseKey.ModifierKey2 != nil {
				combo.SetSelected(*ctrlReviseKey.ModifierKey2)
			}
		case Translate:
			if translateKey.ModifierKey2 != nil {
				combo.SetSelected(*translateKey.ModifierKey2)
			}
		default:
			slog.Error("Invalid action", "action", action)
		}
	case NormalKey:
		switch action {
		case AskQuestion:
			combo.SetSelected(askKey.Key)
		case CtrlRevise:
			combo.SetSelected(ctrlReviseKey.Key)
		case Translate:
			combo.SetSelected(translateKey.Key)
		default:
			slog.Error("Invalid action", "action", action)
		}
	default:
		slog.Error("Invalid key type", "key", key)
	}

	return combo
}
