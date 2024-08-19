package main

import (
	"log/slog"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"

	"github.com/bahelit/ctrl_plus_revise/internal/docker"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
	"github.com/bahelit/ctrl_plus_revise/pkg/bytesize"
	dirsize "github.com/bahelit/ctrl_plus_revise/pkg/dir_size"
)

const lengthOfKeyBoardShortcuts = 3

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
	model := ollama.ModelName(guiApp.Preferences().IntWithFallback(CurrentModelKey, int(ollama.Llama3Dot1)))
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
	useDocker := guiApp.Preferences().BoolWithFallback(UseDockerKey, false)
	if stopOllamaOnShutDown {
		if useDocker {
			docker.StopOllamaContainer()
		} else if p != 0 {
			stopOllama(p)
		}
	} else {
		slog.Info("Leaving Ollama running")
	}
}

func setupServices() bool {
	connectedToOllama := checkOllamaConnection(nil)
	if connectedToOllama {
		return connectedToOllama
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

func setKeyboardShortcuts() {
	ask := guiApp.Preferences().StringListWithFallback(AskAIKeyboardShortcut, getAskKeys())
	askKey.ModifierKey1 = ask[0]
	if len(ask) == lengthOfKeyBoardShortcuts && ask[1] != EmptySelection {
		askKey.ModifierKey2 = &ask[1]
		askKey.Key = ask[2]
	} else {
		askKey.Key = ask[1]
	}

	revise := guiApp.Preferences().StringListWithFallback(CtrlReviseKeyboardShortcut, getAskKeys())
	ctrlReviseKey.ModifierKey1 = revise[0]
	if len(revise) == lengthOfKeyBoardShortcuts && ask[1] != EmptySelection {
		ctrlReviseKey.ModifierKey2 = &revise[1]
		ctrlReviseKey.Key = revise[2]
	} else {
		ctrlReviseKey.Key = revise[1]
	}

	translate := guiApp.Preferences().StringListWithFallback(TranslateKeyboardShortcut, getAskKeys())
	translateKey.ModifierKey1 = translate[0]
	if len(translate) == lengthOfKeyBoardShortcuts && ask[1] != EmptySelection {
		translateKey.ModifierKey2 = &translate[1]
		translateKey.Key = translate[2]
	} else {
		translateKey.Key = translate[1]
	}
}

func startOllama() (connectedToOllama bool) {
	_, err := exec.LookPath("ollama")
	if err != nil {
		slog.Info("Ollama not found", "error", err)
		ollama.InstallOrUpdateOllama()
	}

	ollamaServe := exec.Command("ollama", "serve")
	err = ollamaServe.Start()
	if err != nil {
		return connectedToOllama
	}
	time.Sleep(1 * time.Second)
	go func() {
		err = ollamaServe.Wait()
		if err != nil {
			slog.Error("Ollama process exited", "error", err)
		}
	}()

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

func initSpeech() {
	speakResponse := guiApp.Preferences().BoolWithFallback(SpeakAIResponseKey, false)
	if speakResponse {
		// TODO: the audio files need to be cleaned up periodically.
		speech = &htgotts.Speech{Folder: "audio", Language: voices.English, Handler: &handlers.Native{}}
		dirInfo, _ := dirsize.GetDirInfo(os.DirFS("audio"))
		slog.Info("AI Speech Recordings", "fileCount", dirInfo.FileCount, "size", bytesize.New(float64(dirInfo.TotalSize)))
		_ = speech.Speak("")
	}
}
