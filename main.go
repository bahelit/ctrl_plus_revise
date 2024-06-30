package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/x/fyne/theme"
	"github.com/docker/docker/client"
	"github.com/ollama/ollama/api"
	"golang.design/x/clipboard"
)

var (
	ollamaClient            *api.Client
	selectedPrompt          PromptMsg = CorrectGrammar
	selectedPromptBinding   binding.String
	guiApp                  fyne.App
	containerID             string
	stopContainerOnShutDown bool = false
)

func main() {
	// Start the clipboard listener.
	err := clipboard.Init()
	if err != nil {
		slog.Error("Failed to start clipboard listener!", "details", err.Error())
		os.Exit(1)
	}
	selectedPromptBinding = binding.NewString()
	err = selectedPromptBinding.Set(CorrectGrammar.String())
	if err != nil {
		slog.Error("Failed to set selectedPromptBinding", "error", err)
	}

	// Start the services.
	cli := setupServices()
	if cli == nil {
		os.Exit(1)
	}

	// Start the GUI event loop and system tray.
	slog.Info("Starting ctrl_plus_revise GUI Service...", "version", Version)
	icon, err := fyne.LoadResourceFromPath("images/icon_small.png")
	if err != nil {
		slog.Warn("Failed to load icon", "error", err)
	}
	guiApp = app.NewWithID("Ctrl+Revise")
	guiApp.SetIcon(icon)
	guiApp.Settings().SetTheme(theme.AdwaitaTheme())
	// Prepare the system tray
	sysTray := setupSysTray(guiApp)
	if guiApp.Preferences().BoolWithFallback(showStartWindow, true) {
		slog.Info("Hiding start window")
		sysTray.Show()
	}
	sysTray.SetCloseIntercept(func() {
		sysTray.Hide()
	})

	// Listen for global hotkeys
	go registerHotkeys(sysTray)

	// Handle shutdown signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		slog.Info("Received shutdown signal")
		stopContainerOnShutDown = guiApp.Preferences().BoolWithFallback(stopOllamaOnShutDown, false)
		if stopContainerOnShutDown {
			slog.Info("Stopping Ollama container")
			stopOllamaContainer(cli, containerID)
		} else {
			slog.Info("Leaving Ollama container running")
		}
		os.Exit(0)
	}()

	// Run the GUI event loop.
	guiApp.Run()
}

func setupServices() *client.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to Docker
	cli, err := connectToDocker()
	if err != nil {
		return nil
	}
	ping, err := cli.Ping(ctx)
	if err != nil {
		slog.Error("Failed to connect to Docker", "error", err)
		return nil
	}
	slog.Info("Connected to Docker", "Operating_System", ping.OSType)

	// Start communication with the AI
	connectToOllama()
	err = ollamaClient.Heartbeat(ctx)
	if err == nil {
		slog.Info("Connected to Ollama")
		return cli
	}
	if <-ctx.Done(); true {
		slog.Info("Failed to connect to Ollama, starting docker container", "error", err)
	}

	// TODO: Support native Ollama without Docker
	// If we made it hear we can talk to docker but don't have a connection to Ollama
	slog.Info("Starting Ollama container")
	// Check Docker containers
	containerID, err = startOllamaContainer(cli)
	if err != nil {
		return nil
	}

	return cli
}
