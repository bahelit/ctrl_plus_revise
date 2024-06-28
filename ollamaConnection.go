package main

import (
	"log/slog"
	"os"

	"github.com/ollama/ollama/api"
)

func connectToOllama() {
	var err error
	ollamaClient, err = api.ClientFromEnvironment()
	if err != nil {
		slog.Error("Failed to create client", "error", err)
		os.Exit(1) // For now...
	}
}

func startUpOllamaNative() {
	// TODO: Implement this
}
