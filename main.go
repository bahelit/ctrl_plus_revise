package main

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
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
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
	ollama "github.com/ollama/ollama/api"

	"github.com/bahelit/ctrl_plus_revise/pkg/bytesize"
	dirSize "github.com/bahelit/ctrl_plus_revise/pkg/dir_size"
	"github.com/bahelit/ctrl_plus_revise/version"
)

// TODO: Refactor to remove some global variables
var (
	dockerClient          *docker.Client
	ollamaClient          *ollama.Client
	selectedModel         = Llama3
	selectedPrompt        = CorrectGrammar
	selectedModelBinding  binding.Int
	selectedPromptBinding binding.String
	guiApp                fyne.App
	containerID           string
	ollamaPID             int
	speech                htgotts.Speech
	stopOllamaOnShutDown  = false
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
	go registerHotkeys(sysTray)

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

func handleShutdown(d *docker.Client, p int) {
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

	useDocker := guiApp.Preferences().BoolWithFallback(useDockerKey, false)
	if useDocker {
		dockerClient, err = connectToDocker()
		if err != nil {
			slog.Error("Failed to connect to Docker", "error", err)
		}
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

	if useDocker {
		slog.Info("Starting Ollama container")
		return setupDocker()
	}

	slog.Info("Starting Ollama")
	// Start Ollama locally
	versionCMD := exec.Command("ollama", "--version")
	err = versionCMD.Run()
	if err != nil {
		slog.Error("Can't find Ollama", "error", err)
		return connectedToOllama
	}
	ollamaServe := exec.Command("ollama", "serve")
	err = ollamaServe.Start()
	if err != nil {
		return connectedToOllama
	}
	time.Sleep(5 * time.Second)
	go func() {
		err = ollamaServe.Wait()
		slog.Info("Ollama process exited", "error", err)
	}()
	ollamaPID = ollamaServe.Process.Pid
	slog.Info("Started Ollama", "pid", ollamaPID)

	ollamaClient = connectToOllama()
	connectedToOllama = true
	return connectedToOllama
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

func stopOllama(p int) {
	if p == 0 {
		slog.Error("Ollama PID not found")
		return
	}
	err := syscall.Kill(p, syscall.SIGTERM)
	if err != nil {
		slog.Error("Failed to stop Ollama", "error", err)
	} else {
		return
	}
	// If the above fails, try a more forceful kill
	_ = syscall.Kill(p, syscall.SIGKILL)
}
