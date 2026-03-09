package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AppConfig holds all application settings, serialized to JSON.
type AppConfig struct {
	ReplaceHighlightedText bool     `json:"replaceHighlightedText"`
	SpeakAIResponse        bool     `json:"speakAIResponse"`
	ShowPopUp              bool     `json:"showPopUp"`
	ShowStartWindow        bool     `json:"showStartWindow"`
	CurrentPrompt          int      `json:"currentPrompt"`
	CurrentChatModel       string   `json:"currentChatModel"`
	CurrentModel           string   `json:"currentModel"`
	CurrentFromLang        string   `json:"currentFromLang"`
	CurrentToLang          string   `json:"currentToLang"`
	StopOllamaOnShutDown   bool     `json:"stopOllamaOnShutDown"`
	UseRemoteOllama        bool     `json:"useRemoteOllama"`
	OllamaURL              string   `json:"ollamaURL"`
	AskAIShortcut          []string `json:"askAIShortcut"`
	CtrlReviseShortcut     []string `json:"ctrlReviseShortcut"`
	TranslateShortcut      []string `json:"translateShortcut"`
}

// DefaultConfig returns an AppConfig populated with sensible defaults.
func DefaultConfig() *AppConfig {
	return &AppConfig{
		ReplaceHighlightedText: true,
		SpeakAIResponse:        false,
		ShowPopUp:              false,
		ShowStartWindow:        true,
		CurrentPrompt:          0,
		CurrentChatModel:       "llama3.2:latest",
		CurrentModel:           "llama3.2:latest",
		CurrentFromLang:        "English",
		CurrentToLang:          "Spanish",
		StopOllamaOnShutDown:   true,
		UseRemoteOllama:        false,
		OllamaURL:              "http://localhost:11434",
		AskAIShortcut:          []string{"alt", "a"},
		CtrlReviseShortcut:     []string{"alt", "c"},
		TranslateShortcut:      []string{"alt", "t"},
	}
}

// configDir returns the directory where the config file is stored.
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "ctrl_plus_revise"), nil
}

// configPath returns the full path to the config file.
func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// LoadConfig reads the config from ~/.config/ctrl_plus_revise/config.json.
// If the file does not exist it returns DefaultConfig with no error.
func LoadConfig() (*AppConfig, error) {
	path, err := configPath()
	if err != nil {
		return DefaultConfig(), err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return DefaultConfig(), nil
	}
	if err != nil {
		return DefaultConfig(), err
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return DefaultConfig(), err
	}
	return cfg, nil
}

// Save writes the config to ~/.config/ctrl_plus_revise/config.json.
func (c *AppConfig) Save() error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	path := filepath.Join(dir, "config.json")
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
