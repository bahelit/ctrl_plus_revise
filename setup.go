package main

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"github.com/docker/docker/client"
	"github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"

	"github.com/bahelit/ctrl_plus_revise/pkg/bytesize"
	"github.com/bahelit/ctrl_plus_revise/pkg/dir_size"
)

func sayHello() {
	speakResponse := guiApp.Preferences().BoolWithFallback(speakAIResponseKey, false)
	if speakResponse {
		go func() {
			prompt := guiApp.Preferences().StringWithFallback(currentPromptKey, CorrectGrammar.String())
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
	model := ModelName(guiApp.Preferences().IntWithFallback(currentModelKey, int(Llama3)))
	err := pullModel(model, false)
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

func handleShutdown(d *client.Client, p int) {
	stopOllamaOnShutDown = guiApp.Preferences().BoolWithFallback(stopOllamaOnShutDownKey, true)
	useDocker := guiApp.Preferences().BoolWithFallback(useDockerKey, true)
	if stopOllamaOnShutDown {
		if useDocker {
			if d != nil {
				stopOllamaContainer(d)
			}
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
	// Start the speech handler
	speech = htgotts.Speech{Folder: "audio", Language: voices.English, Handler: &handlers.Native{}}
	guiApp.Preferences().SetBool(speakAIResponseKey, false)
	speakResponse := guiApp.Preferences().BoolWithFallback(speakAIResponseKey, false)
	if speakResponse {
		// TODO: the audio files need to be cleaned up periodically.
		dirInfo, _ := dirSize.GetDirInfo(os.DirFS("audio"))
		slog.Info("AI Speech Recordings", "fileCount", dirInfo.FileCount, "size", bytesize.New(float64(dirInfo.TotalSize)))
	}

	heartBeatCtx, heartBeatCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer heartBeatCancel()
	// Start communication with the AI
	ollamaClient = connectToOllama()
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

	useDocker := guiApp.Preferences().BoolWithFallback(useDockerKey, false)
	if useDocker {
		dockerClient, err = connectToDocker()
		if err != nil {
			slog.Error("Failed to connect to Docker", "error", err)
		}

		slog.Info("Starting Ollama container")
		return setupDocker()
	}

	slog.Info("Starting Ollama")
	// Start Ollama locally
	return startOllama()
}

func setupDocker() (connectedToOllama bool) {
	connectedToOllama = false
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	ping, err := dockerClient.Ping(pingCtx)
	if err != nil {
		slog.Error("Failed to connect to Docker", "error", err)
		pingCancel()
		return connectedToOllama
	}
	if <-pingCtx.Done(); true {
		slog.Error("Timed out trying to connect to Docker")
	}
	slog.Info("Connected to Docker", "Operating_System", ping.OSType)

	// If we made it hear we can talk to docker but don't have a connection to Ollama
	slog.Info("Starting Ollama container")
	// Check Docker containers
	containerID, err = startOllamaContainer(dockerClient)
	if err != nil {
		return connectedToOllama
	}
	if containerID != "" {
		connectedToOllama = true
		return connectedToOllama
	}
	return connectedToOllama
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

	ollamaClient = connectToOllama()
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
