package bindings

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"github.com/bahelit/ctrl_plus_revise/internal/config"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
)

var (
	SelectedModelBinding   = binding.NewInt()
	TranslationFromBinding = binding.NewString()
	TranslationToBinding   = binding.NewString()
	SelectedPromptBinding  = binding.NewString()

	AiActionDropdown *widget.Select
	AiModelDropdown  *widget.Select
)

func SetBindingVariables(guiApp fyne.App) error {
	model := guiApp.Preferences().IntWithFallback(config.CurrentModelKey, int(ollama.Llama3Dot2))
	err := SelectedModelBinding.Set(model)
	if err != nil {
		slog.Error("Failed to set SelectedModelBinding", "error", err)
	}

	from := guiApp.Preferences().StringWithFallback(config.CurrentFromLangKey, string(ollama.English))
	err = TranslationFromBinding.Set(from)
	if err != nil {
		slog.Error("Failed to set SelectedModelBinding", "error", err)
	}

	to := guiApp.Preferences().StringWithFallback(config.CurrentToLangKey, string(ollama.Spanish))
	err = TranslationToBinding.Set(to)
	if err != nil {
		slog.Error("Failed to set SelectedModelBinding", "error", err)
	}

	prompt := guiApp.Preferences().StringWithFallback(config.CurrentPromptKey, ollama.CorrectGrammar.String())
	err = SelectedPromptBinding.Set(prompt)
	if err != nil {
		slog.Error("Failed to set SelectedPromptBinding", "error", err)
	}
	return err
}
