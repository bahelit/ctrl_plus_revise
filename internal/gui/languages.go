package gui

import "github.com/bahelit/ctrl_plus_revise/internal/ollama"

var (
	Languages = []string{string(ollama.English),
		string(ollama.Arabic),
		string(ollama.Chinese),
		string(ollama.French),
		string(ollama.German),
		string(ollama.Italian),
		string(ollama.Japanese),
		string(ollama.Portuguese),
		string(ollama.Russian),
		string(ollama.Spanish),
		string(ollama.Turkish)}
)
