package shortcuts

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/layout"

	"github.com/bahelit/ctrl_plus_revise/internal/config"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/bindings"
)

func ShowShortcuts(guiApp fyne.App) {
	slog.Debug("Showing Shortcuts")
	shortCuts := guiApp.NewWindow("Ctrl+Revise Keyboard Shortcuts")

	var grid *fyne.Container

	warn := widget.NewIcon(theme.WarningIcon())
	restartToReload := widget.NewLabelWithStyle("Restart application for changes to take effect", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	restartForChanges := container.NewGridWithRows(2, warn, restartToReload)

	askLabel := widget.NewLabel("\nAsk a Question with highlighted text")
	askModDropdown1 := keyboardModifierButtonsDropDown(guiApp, AskQuestion, ModifierKey1)
	askModDropdown2 := keyboardModifierButtonsDropDown(guiApp, AskQuestion, ModifierKey2)
	askKeyDropdown1 := keyboardModifierButtonsDropDown(guiApp, AskQuestion, NormalKey)

	reviseLabel := widget.NewLabel("Selected Revise Action: ")
	label2Binding := widget.NewLabelWithData(bindings.SelectedPromptBinding)
	hBox := container.NewGridWithColumns(2, reviseLabel, label2Binding)
	ctrlReviseModDropdown1 := keyboardModifierButtonsDropDown(guiApp, CtrlRevise, ModifierKey1)
	ctrlReviseModDropdown2 := keyboardModifierButtonsDropDown(guiApp, CtrlRevise, ModifierKey2)
	ctrlReviseKeyDropdown1 := keyboardModifierButtonsDropDown(guiApp, CtrlRevise, NormalKey)

	cyclePromptLabel := widget.NewLabel("\nCycle through the prompt options")
	cyclePromptValue := widget.NewLabel("Alt + P")
	cyclePromptValue.TextStyle = fyne.TextStyle{Bold: true}

	readerLabel := widget.NewLabel("\nRead the highlighted text")
	readerValue := widget.NewLabel("Alt + R")
	readerValue.TextStyle = fyne.TextStyle{Bold: true}

	from, err := bindings.TranslationFromBinding.Get()
	if err != nil {
		slog.Error("Failed to get translationFromBinding", "error", err)
	}
	to, err := bindings.TranslationToBinding.Get()
	if err != nil {
		slog.Error("Failed to get translationToBinding", "error", err)
	}
	slog.Info("Translation languages", "from", from, "to", to)
	translateLabel := widget.NewLabel("\nTranslate the highlighted text, From: " + from + " To: " + to)
	translateKeyModDropdown1 := keyboardModifierButtonsDropDown(guiApp, Translate, ModifierKey1)
	translateKeyModDropdown2 := keyboardModifierButtonsDropDown(guiApp, Translate, ModifierKey2)
	translateKeyDropdown1 := keyboardModifierButtonsDropDown(guiApp, Translate, NormalKey)

	askKeys := container.NewAdaptiveGrid(config.LengthOfKeyBoardShortcuts, askModDropdown1, askModDropdown2, askKeyDropdown1)
	ctrlReviseKeys := container.NewAdaptiveGrid(config.LengthOfKeyBoardShortcuts, ctrlReviseModDropdown1, ctrlReviseModDropdown2, ctrlReviseKeyDropdown1)
	translateKeys := container.NewAdaptiveGrid(config.LengthOfKeyBoardShortcuts, translateKeyModDropdown1, translateKeyModDropdown2, translateKeyDropdown1)

	grid = layout.NewResponsiveLayout(
		restartForChanges,
		hBox, ctrlReviseKeys,
		askLabel, askKeys,
		translateLabel, translateKeys,
		cyclePromptLabel, cyclePromptValue)

	speakAIResponse := guiApp.Preferences().BoolWithFallback(config.SpeakAIResponseKey, false)
	if speakAIResponse {
		grid.Add(readerLabel)
		grid.Add(readerValue)
	}
	shortCuts.SetContent(grid)
	shortCuts.Show()
}
