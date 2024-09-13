package docker

import (
	"context"
	"log/slog"
	"time"

	docker "github.com/docker/docker/client"
)

var (
	Client      *docker.Client
	ContainerID string
)

func connectToDocker() (*docker.Client, error) {
	cli, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		slog.Error("Failed to create docker client", "error", err)
		return nil, err
	}
	return cli, nil
}

func SetupDocker() (connectedToOllamaContainer bool) {
	connectedToOllamaContainer = false
	dockerClient, err := connectToDocker()
	if err != nil {
		slog.Error("Failed to connect to Docker", "error", err)
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	ping, err := dockerClient.Ping(pingCtx)
	if err != nil {
		slog.Error("Failed to connect to Docker", "error", err)
		pingCancel()
		return connectedToOllamaContainer
	}
	if <-pingCtx.Done(); true {
		slog.Error("Timed out trying to connect to Docker")
	}
	slog.Info("Connected to Docker", "Operating_System", ping.OSType)

	// If we made it hear we can talk to docker but don't have a connection to Ollama
	slog.Info("Starting Ollama container")
	// Check Docker containers
	ContainerID, err = startOllamaContainer(dockerClient)
	if err != nil {
		return connectedToOllamaContainer
	}
	if ContainerID != "" {
		connectedToOllamaContainer = true
		return connectedToOllamaContainer
	}
	return connectedToOllamaContainer
}
