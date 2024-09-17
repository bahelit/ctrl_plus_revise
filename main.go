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
	"fyne.io/x/fyne/theme"
	ollamaApi "github.com/ollama/ollama/api"

	"github.com/bahelit/ctrl_plus_revise/internal/config"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/loading"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/settings"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/shortcuts"
	"github.com/bahelit/ctrl_plus_revise/internal/hardware"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
	"github.com/bahelit/ctrl_plus_revise/version"
)

var (
	guiApp               fyne.App
	stopOllamaOnShutDown = false
)

func main() {
	slog.Info("Starting Ctr+Revise gui Service...", "Version", version.Version, "Compiler", runtime.Version())
	guiApp = app.NewWithID("com.ctrlplusrevise.app")
	guiApp.Settings().SetTheme(theme.AdwaitaTheme())

	var ollamaClient *ollamaApi.Client
	ollamaClient = ollama.CheckOllamaConnection(guiApp, ollamaClient, nil)

	// Prepare the loading screen and system tray
	startupWindow := loading.StartupScreen(guiApp)
	sysTray := SetupSysTray(guiApp, ollamaClient)
	if guiApp.Preferences().BoolWithFallback(config.ShowStartWindowKey, true) {
		slog.Debug("Hiding start window")
		sysTray.Show()
		startupWindow.Show()
	}
	loadIcon(guiApp)

	// Start the services
	hardware.DetectProcessingDevice()
	hardware.DetectMemory()
	go func() {
		if ollamaClient == nil {
			ollamaClient = settings.SetupServices(guiApp, ollamaClient)
			if ollamaClient == nil {
				slog.Error("Failed to connect to Ollama")
				ollama.InstallOrUpdateOllamaWindow(guiApp, ollamaClient)
			}
		}
		fetchModel(ollamaClient)
		time.Sleep(1 * time.Second)
		startupWindow.Close()
	}()

	//sayHello()

	// Listen for global hotkeys
	setKeyboardShortcuts()
	go func() {
		shortcuts.StartKeyboardListener(guiApp, ollamaClient)
	}()

	// Handle shutdown signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-c
		slog.Info("Received shutdown signal")
		handleShutdown(nil)
		os.Exit(0)
	}()
	defer func() {
		slog.Info("Shutting down")
		signal.Stop(c)
		close(c)
		handleShutdown(nil)
	}()

	slog.Info("Ctrl+Revise is ready to help you!")
	// Run the gui event loop
	guiApp.Run()
}
