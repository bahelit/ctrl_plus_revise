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
	ollamaApi "github.com/ollama/ollama/api"
	hook "github.com/robotn/gohook"

	"github.com/bahelit/ctrl_plus_revise/internal/gui"
	"github.com/bahelit/ctrl_plus_revise/internal/hardware"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
	"github.com/bahelit/ctrl_plus_revise/version"
)

// TODO: Refactor to remove some global variables
var (
	ollamaClient           *ollamaApi.Client
	selectedModel          = ollama.Llama3
	selectedPrompt         = ollama.CorrectGrammar
	selectedModelBinding   = binding.NewInt()
	translationFromBinding = binding.NewString()
	translationToBinding   = binding.NewString()
	selectedPromptBinding  = binding.NewString()
	guiApp                 fyne.App
	systemHook             chan hook.Event
	ollamaPID              int
	stopOllamaOnShutDown   = false
)

func main() {
	slog.Info("Starting Ctr+Revise gui Service...", "Version", version.Version, "Compiler", runtime.Version())
	guiApp = app.NewWithID("com.ctrlplusrevise.app")
	guiApp.Settings().SetTheme(theme.AdwaitaTheme())

	// Prepare the loading screen and system tray
	startupWindow := gui.StartupScreen(guiApp)
	sysTray := SetupSysTray(guiApp)
	if guiApp.Preferences().BoolWithFallback(ShowStartWindowKey, true) {
		slog.Debug("Hiding start window")
		sysTray.Show()
		startupWindow.Show()
	}
	LoadIcon(guiApp)

	// Start the services
	hardware.DetectProcessingDevice()
	hardware.DetectMemory()
	go func() {
		connectedToOllama := setupServices()
		if !connectedToOllama {
			slog.Error("Failed to connect to Ollama")
			installOrUpdateOllamaWindow(guiApp)
		}
		fetchModel()
		time.Sleep(1 * time.Second)
		startupWindow.Close()
	}()

	sayHello()

	// Listen for global hotkeys
	setKeyboardShortcuts()
	go func() {
		startKeyboardListener()
	}()
	initSpeech()

	// Handle shutdown signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func(p int) {
		<-c
		slog.Info("Received shutdown signal")
		handleShutdown(p)
		os.Exit(0)
	}(ollamaPID)
	defer func(p int) {
		slog.Info("Shutting down")
		signal.Stop(c)
		close(c)
		handleShutdown(p)
	}(ollamaPID)

	slog.Info("Ctrl+Revise is ready to help you!")
	// Run the gui event loop
	guiApp.Run()
}
