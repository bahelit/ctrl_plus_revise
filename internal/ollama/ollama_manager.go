package ollama

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	ollamaApi "github.com/ollama/ollama/api"

	"github.com/bahelit/ctrl_plus_revise/internal/config"
	"github.com/bahelit/ctrl_plus_revise/internal/docker"
)

func InstallOrUpdateOllamaWindow(guiApp fyne.App, ollamaClient *ollamaApi.Client) {
	slog.Debug("Asking Question")
	var (
		screenHeight float32 = 480.0
		screenWidth  float32 = 650.0
	)

	ollamaManagerWindow := guiApp.NewWindow("Ctrl+Revise Ollama Manager")
	ollamaManagerWindow.Resize(fyne.NewSize(screenWidth, screenHeight))

	ollamaURL := guiApp.Preferences().StringWithFallback(config.OllamaURLKey, "http://localhost:11434")
	urlOverride := os.Getenv("OLLAMA_HOST")
	if urlOverride != "" {
		ollamaURL = urlOverride
	}
	urlNote := widget.NewLabel("Press the \"Enter\" key to test the connection")
	ollamaURLEntry := widget.NewEntry()
	ollamaURLEntry.SetText(ollamaURL)
	ollamaURLEntry.Validator = func(s string) error {
		if s == "" || len(s) < 3 {
			return nil
		}
		_, err := url.Parse(ollamaURLEntry.Text)
		if err != nil {
			slog.Error("Failed to parse URL", "error", err)
			return errors.New("invalid URL")
		}
		return nil
	}
	ollamaURLEntry.OnSubmitted = func(s string) {
		err := ollamaURLEntry.Validate()
		if err != nil {
			w := guiApp.NewWindow("Invalid URL")
			msg := widget.NewLabel("Please enter a valid URL with a port")
			msg.TextStyle = fyne.TextStyle{Bold: true}
			msg.Alignment = fyne.TextAlignCenter
			validationErrMsg := widget.NewLabel("Validation Error: " + err.Error())
			errLayout := container.NewVBox(msg, validationErrMsg)
			w.SetContent(errLayout)
			w.Show()
			time.Sleep(3 * time.Second)
			w.Close()
			slog.Error("Invalid URL", "error", err)
			return
		}
		ollamaClient := CheckOllamaConnection(guiApp, ollamaClient, &s)
		if ollamaClient != nil {
			w := guiApp.NewWindow("Successfully Connected")
			msg := widget.NewLabel("Successfully Connected to Ollama")
			msg.TextStyle = fyne.TextStyle{Bold: true}
			msg.Alignment = fyne.TextAlignCenter
			w.SetContent(msg)
			w.Show()
			time.Sleep(3 * time.Second)
			w.Close()
		} else {
			w := guiApp.NewWindow("Invalid URL")
			msg := widget.NewLabel("Couldn't connect to Ollama - Please check the URL and port are valid")
			msg.TextStyle = fyne.TextStyle{Bold: true}
			msg.Alignment = fyne.TextAlignCenter
			w.SetContent(msg)
			w.Show()
			time.Sleep(3 * time.Second)
			w.Close()
			slog.Error("Invalid URL", "error", err)
		}
	}
	useRemoteOllamaCheckbox := remoteOllamaCheckbox(guiApp)
	useRemoteOllamaCheckbox.OnChanged = func(b bool) {
		if b {
			ollamaURLEntry.Show()
			urlNote.Show()
			guiApp.Preferences().SetBool(config.UseRemoteOllamaKey, true)
			slog.Debug("Show Pop-Up is on")
		} else {
			ollamaURLEntry.Hide()
			urlNote.Hide()
			guiApp.Preferences().SetBool(config.UseRemoteOllamaKey, false)
			slog.Debug("Show Pop-Up is off")
		}
		ollamaURLEntry.Refresh()
	}
	if useRemoteOllamaCheckbox.Checked {
		ollamaURLEntry.Show()
		urlNote.Show()
	} else {
		ollamaURLEntry.Hide()
		urlNote.Hide()
	}
	manualInstall := widget.NewHyperlink(
		"Manual Install Link",
		&url.URL{
			Scheme: "https",
			Host:   "ollama.com",
			Path:   "/download",
		})

	var err error
	fetchOllama := widget.NewButton("Download or Update Ollama", func() {
		err = installOrUpdateOllama(guiApp, ollamaClient)
		if err != nil {
			guiApp.SendNotification(&fyne.Notification{
				Title:   "Ollama Installation Error",
				Content: "Failed to install Ollama."})
			slog.Error("Failed to install Ollama", "error", err)
		} else {
			guiApp.SendNotification(&fyne.Notification{
				Title:   "Ollama Installed",
				Content: "Ollama has been installed."})
			slog.Info("Ollama installed")
		}
	})

	dockerMsg := widget.NewLabel("Docker is enabled")
	useDocker := guiApp.Preferences().BoolWithFallback(config.UseDockerKey, false)
	if useDocker {
		dockerMsg.Hide()
	}

	managerLayout := container.NewVBox(useRemoteOllamaCheckbox, urlNote, ollamaURLEntry, dockerMsg, fetchOllama, manualInstall)

	ollamaManagerWindow.SetContent(managerLayout)
	ollamaManagerWindow.Canvas().Focus(ollamaURLEntry)
	ollamaManagerWindow.Show()

}

func installOrUpdateOllama(guiApp fyne.App, ollamaClient *ollamaApi.Client) error {
	slog.Info("Installing Ollama")
	useDocker := guiApp.Preferences().BoolWithFallback(config.UseDockerKey, false)
	if useDocker {
		ollamaClient = CheckOllamaConnection(guiApp, ollamaClient, nil)
		if ollamaClient != nil {
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
	// Download the file
	slog.Info("Downloading file", "url", ollamaURL, "filepath", ollamaPath)
	// Create the file
	out, err := os.Create(ollamaPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(ollamaURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
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

func remoteOllamaCheckbox(guiApp fyne.App) *widget.Check {
	showPopUp := guiApp.Preferences().BoolWithFallback(config.UseRemoteOllamaKey, false)
	popup := widget.NewCheck("Connect to Ollama server", func(b bool) {
		if !b {
			slog.Debug("Not using remote server connection")
		} else if b {
			slog.Debug("Using remote server connection")
		}
	})
	popup.Checked = showPopUp
	return popup
}

func CheckOllamaConnection(guiApp fyne.App, ollamaClient *ollamaApi.Client, ollamaURL *string) *ollamaApi.Client {
	heartBeatCtx, heartBeatCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer heartBeatCancel()
	// Start communication with the AI
	if ollamaURL != nil {
		ollamaClient = ConnectToOllamaWithURL(*ollamaURL)
	} else {
		ollamaClient = ConnectToOllama()
	}
	err := ollamaClient.Heartbeat(heartBeatCtx)
	if err == nil {
		slog.Info("Connected to Ollama")
		return ollamaClient
	}
	if <-heartBeatCtx.Done(); true {
		guiApp.SendNotification(&fyne.Notification{
			Title:   "Ollama Connection Error",
			Content: "Timed out trying to connect to Ollama."})
		slog.Error("timed out connecting to Ollama")
	}
	return nil
}
