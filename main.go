package main

import (
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/x/fyne/theme"
	docker "github.com/docker/docker/client"
	htgotts "github.com/hegedustibor/htgo-tts"
	ollama "github.com/ollama/ollama/api"

	"github.com/bahelit/ctrl_plus_revise/version"
)

// TODO: Refactor to remove some global variables
var (
	dockerClient           *docker.Client
	ollamaClient           *ollama.Client
	selectedModel          = Llama3
	selectedPrompt         = CorrectGrammar
	selectedModelBinding   binding.Int
	translationToBinding   binding.String
	translationFromBinding binding.String
	selectedPromptBinding  binding.String
	guiApp                 fyne.App
	containerID            string
	ollamaPID              int
	speech                 htgotts.Speech
	stopOllamaOnShutDown   = false
)

func main() {
	slog.Info("Starting Ctr+Revise GUI Service...", "Version", version.Version, "Compiler", runtime.Version())
	guiApp = app.NewWithID("com.bahelit.ctrl_plus_revise")
	guiApp.Settings().SetTheme(theme.AdwaitaTheme())

	// Prepare the loading screen and system tray
	loadIcon()
	startupWindow := startupScreen()
	sysTray := setupSysTray(guiApp)
	if guiApp.Preferences().BoolWithFallback(showStartWindowKey, true) {
		slog.Debug("Hiding start window")
		sysTray.Show()
		startupWindow.Show()
	}

	// Start the services
	detectProcessingDevice()
	detectMemory()
	go func() {
		connectedToOllama := setupServices()
		if !connectedToOllama {
			os.Exit(1)
		}
		fetchModel()
		time.Sleep(1 * time.Second)
		startupWindow.Close()
	}()

	sayHello()

	// Listen for global hotkeys
	go registerHotkeys()

	// Handle shutdown signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func(d *docker.Client, p int) {
		<-c
		slog.Info("Received shutdown signal")
		handleShutdown(d, p)
		os.Exit(0)
	}(dockerClient, ollamaPID)
	defer func(d *docker.Client, p int) {
		slog.Info("Shutting down")
		signal.Stop(c)
		close(c)
		handleShutdown(d, p)
	}(dockerClient, ollamaPID)

	slog.Info("Ctrl+Revise is ready to help you!")
	// Run the GUI event loop
	guiApp.Run()
}
