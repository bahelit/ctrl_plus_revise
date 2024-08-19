package ollama

import (
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

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

func ConnectToOllamaWithURL(rawURL string) *ollama.Client {
	ollamaURL, err := url.Parse(rawURL)
	if err != nil {
		slog.Error("Failed to parse URL", "error", err)
		return nil
	}
	httpClient := http.DefaultClient
	httpClient.Timeout = time.Second * 5

	return ollama.NewClient(ollamaURL, httpClient)
}

func startUpOllamaNative() {
	// TODO: Implement this
}

func InstallOrUpdateOllama() {
	slog.Info("Installing Ollama")

}
