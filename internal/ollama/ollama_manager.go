package ollama

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	ollamaApi "github.com/ollama/ollama/api"

	"github.com/bahelit/ctrl_plus_revise/internal/docker"
)

// CheckOllamaConnection attempts to connect to Ollama at the given URL (or
// the environment-configured URL when ollamaURL is nil) and returns the
// connected client.  It returns (nil, error) on failure without showing any
// GUI notifications.
func CheckOllamaConnection(ollamaURL *string) (*ollamaApi.Client, error) {
	heartBeatCtx, heartBeatCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer heartBeatCancel()

	var ollamaClient *ollamaApi.Client
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

// InstallOrUpdateOllama installs or updates Ollama.  When useDocker is true it
// delegates to the docker package; otherwise it runs the platform-native
// installer.
func InstallOrUpdateOllama(useDocker bool, ollamaClient *ollamaApi.Client) error {
	slog.Info("Installing Ollama")
	if useDocker {
		client, err := CheckOllamaConnection(nil)
		if err == nil && client != nil {
			slog.Info("Docker container update not yet implemented")
			return nil
		}
		docker.SetupDocker()
		return nil
	}

	operatingSystem := runtime.GOOS
	architecture := runtime.GOARCH
	slog.Info("Operating System", "OS", operatingSystem, "Architecture", architecture)

	var err error
	switch operatingSystem {
	case "linux":
		slog.Info("Installing Ollama on Linux")
		err = installOllamaLinux()
	case "darwin":
		slog.Info("Installing Ollama on MacOS")
		err = installOllamaMacOS()
	case "windows":
		slog.Info("Installing Ollama on Windows")
		err = installOllamaWindows()
	default:
		slog.Warn("Operating System not detected, manual installation required")
	}
	if err != nil {
		slog.Error("Failed to install Ollama", "error", err)
		return err
	}
	return nil
}

func installOllamaLinux() error {
	curlCommand := "curl -fsSL https://ollama.com/install.sh | sh"
	cmd := exec.Command("bash", "-c", curlCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install Ollama: %w, output: %s", err, string(output))
	}
	slog.Info("Install Output:\n ", "bash", string(output))
	return nil
}

func installOllamaMacOS() error {
	brewCommand := "brew install ollama && brew services start ollama"
	cmd := exec.Command("bash", "-c", brewCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install Ollama: %w, output: %s", err, string(output))
	}
	slog.Info("Install Output:\n ", "brew", string(output))
	return nil
}

func installOllamaWindows() error {
	ollamaURL := "https://ollama.com/download/OllamaSetup.exe"
	ollamaPath := "OllamaSetup.exe"

	tmpDir := os.TempDir()
	ollamaPath = filepath.Join(tmpDir, ollamaPath)

	slog.Info("Downloading file", "url", ollamaURL, "filepath", ollamaPath)
	out, err := os.Create(ollamaPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(ollamaURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	cmd := exec.Command(ollamaPath)
	err = cmd.Run()
	if err != nil {
		slog.Error("Error installing Ollama", "error", err)
		return err
	}
	return nil
}
