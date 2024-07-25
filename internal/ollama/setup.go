package ollama

import (
	"log/slog"
	"os"

	ollama "github.com/ollama/ollama/api"
)

func ConnectToOllama() *ollama.Client {
	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		slog.Error("Failed to create client", "error", err)
		os.Exit(1) // For now...
	}
	return client
}

func startUpOllamaNative() {
	// TODO: Implement this
}
