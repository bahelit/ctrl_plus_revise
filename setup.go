package main

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"github.com/bahelit/ctrl_plus_revise/internal/docker"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
)

func sayHello() {
	speakResponse := guiApp.Preferences().BoolWithFallback(SpeakAIResponseKey, false)
	if speakResponse {
		go func() {
			prompt := guiApp.Preferences().StringWithFallback(CurrentPromptKey, ollama.CorrectGrammar.String())
			_ = speech.Speak("")
			speakErr := speech.Speak("Control Plus Revise is set to: " + prompt)
			if speakErr != nil {
				slog.Error("Failed to speak", "error", speakErr)
			}
		}()
	}
}

func fetchModel() {
	// Pull the model on startup, will pull updated model if available
	model := ollama.ModelName(guiApp.Preferences().IntWithFallback(CurrentModelKey, int(ollama.Llama3)))
	err := PullModelWrapper(model, false)
	if err != nil {
		slog.Error("Failed to pull model", "error", err)
		guiApp.SendNotification(&fyne.Notification{
			Title: "Ollama Error",
			Content: "Failed to connect to pull model from Ollama\n" +
				"Check logs for more information\n" +
				"Ctrl+Revise will continue running, but may not function correctly",
		})
	}
}

func handleShutdown(p int) {
	stopOllamaOnShutDown = guiApp.Preferences().BoolWithFallback(StopOllamaOnShutDownKey, true)
	useDocker := guiApp.Preferences().BoolWithFallback(UseDockerKey, true)
	if stopOllamaOnShutDown {
		if useDocker {
			docker.StopOllamaContainer()
		} else {
			if p != 0 {
				stopOllama(p)
			}
		}
	} else {
		slog.Info("Leaving Ollama container running")
	}
}

func setupServices() bool {
	connectedToOllama := false
	var err error

	heartBeatCtx, heartBeatCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer heartBeatCancel()
	// Start communication with the AI
	ollamaClient = ollama.ConnectToOllama()
	err = ollamaClient.Heartbeat(heartBeatCtx)
	if err == nil {
		slog.Info("Connected to Ollama")
		connectedToOllama = true
		return connectedToOllama
	} else {
		slog.Error("Ollama doesn't appear to be running", "error", err)
		heartBeatCancel()
	}
	if <-heartBeatCtx.Done(); true {
		slog.Error("Ollama heartbeat timed out")
	}

	// Ollama isn't running, should we start it with Docker?
	useDocker := guiApp.Preferences().BoolWithFallback(UseDockerKey, false)
	if useDocker {
		slog.Info("Starting Ollama container")
		return docker.SetupDocker()
	}
	slog.Info("Starting Ollama")
	// Start Ollama locally
	return startOllama()
}

func startOllama() (connectedToOllama bool) {
	versionCMD := exec.Command("ollama", "--version")
	err := versionCMD.Run()
	if err != nil {
		slog.Error("Can't find Ollama", "error", err)
		return connectedToOllama
	}
	ollamaServe := exec.Command("ollama", "serve")
	err = ollamaServe.Start()
	if err != nil {
		return connectedToOllama
	}
	err = ollamaServe.Wait()
	if err != nil {
		slog.Error("Ollama process exited", "error", err)
		return connectedToOllama
	}

	ollamaPID = ollamaServe.Process.Pid
	slog.Info("Started Ollama", "pid", ollamaPID)

	ollamaClient = ollama.ConnectToOllama()
	if ollamaClient != nil {
		connectedToOllama = true
	}

	return connectedToOllama
}

func stopOllama(p int) {
	if p == 0 {
		slog.Error("Ollama PID not found")
		return
	}
	ollamaProcess, err := os.FindProcess(p)
	if err != nil {
		slog.Error("Failed to find process", "error", err)
	}

	err = ollamaProcess.Kill()
	if err != nil {
		slog.Error("Failed to find process", "error", err)
	}
}
