package ollama

import (
	"context"
	"fmt"
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

// CheckOllamaConnection attempts to connect to Ollama at the given URL (or
// the environment-configured URL when ollamaURL is nil) and returns the
// connected client.  It returns (nil, error) on failure.
func CheckOllamaConnection(ollamaURL *string) (*ollama.Client, error) {
	heartBeatCtx, heartBeatCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer heartBeatCancel()

	var ollamaClient *ollama.Client
	if ollamaURL != nil {
		ollamaClient = ConnectToOllamaWithURL(*ollamaURL)
	} else {
		ollamaClient = ConnectToOllama()
	}

	err := ollamaClient.Heartbeat(heartBeatCtx)
	if err == nil {
		slog.Info("Connected to Ollama")
		return ollamaClient, nil
	}

	select {
	case <-heartBeatCtx.Done():
		slog.Error("timed out connecting to Ollama")
		return nil, fmt.Errorf("timed out connecting to Ollama: %w", heartBeatCtx.Err())
	default:
		slog.Error("failed to connect to Ollama", "error", err)
		return nil, fmt.Errorf("failed to connect to Ollama: %w", err)
	}
}
