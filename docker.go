package main

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	ollamaStartName     = "ollama"
	ollamaContainerName = "/ollama"
	ollamaTagRocm       = "ollama/ollama:rocm"
)

func connectToDocker() (*docker.Client, error) {
	cli, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		slog.Error("Failed to create docker client", "error", err)
		return nil, err
	}
	return cli, nil
}

func findContainer(cli *docker.Client) (containerID string, containerIsRunning bool, err error) {
	ctx := context.Background()
	containerIsRunning = false

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		slog.Error("Failed to list containers", "error", err)
		return "", containerIsRunning, err
	}

	for ctr := range containers {
		slog.Debug("Container found", "id", containers[ctr].ID, "name", containers[ctr].Names)
		if containers[ctr].Names[0] == ollamaContainerName {
			slog.Info("Found Ollama container", "id", containers[ctr].ID, "state", containers[ctr].State)
			if containers[ctr].State == "exited" {
				slog.Info("Ollama container is not running")
				return containers[ctr].ID, containerIsRunning, nil
			}
			slog.Info("Ollama container is running")
			containerIsRunning = true
			return containers[ctr].ID, containerIsRunning, nil
		}
	}
	slog.Info("Ollama container not found")
	return "", containerIsRunning, nil
}

func checkForOllamaImage(cli *docker.Client) (string, error) {
	ctx := context.Background()

	images, err := cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		slog.Error("Failed to list images", "error", err)
		return "", err
	}

	for img := range images {
		slog.Debug("Image found", "id", images[img].ID, "tags", images[img].RepoTags)
		if images[img].RepoTags[0] == "ollama:latest" || images[img].RepoTags[0] == ollamaTagRocm {
			slog.Info("Found Ollama image", "id", images[img].ID)
			return images[img].ID, nil
		}
	}
	slog.Info("Ollama image not found")
	return "", nil
}

func removeContainerImage(cli *docker.Client, id string) {
	ctx := context.Background()

	resp, err := cli.ImageRemove(ctx, id, image.RemoveOptions{})
	if err != nil {
		slog.Error("Failed to remove image", "error", err)
		return
	}

	slog.Info("Removed Ollama image", "id", id, "response", resp)
}

func pullOllamaImage(cli *docker.Client) error {
	ctx := context.Background()

	// TODO: Allow user to specify image AMD or Nvidia images
	reader, err := cli.ImagePull(ctx, ollamaTagRocm, image.PullOptions{})
	if err != nil {
		slog.Error("Failed to pull image", "error", err)
		return err
	}
	defer func() { _ = reader.Close() }()

	// cli.ImagePull is asynchronous.
	// The reader needs to be read completely for the pull operation to complete.
	_, err = io.Copy(io.Discard, reader)
	if err != nil {
		slog.Error("Failed to pull docker image", "error", err)
		return err
	}

	slog.Info("Pulled Ollama image")
	return nil
}

func updateOllamaImage(cli *docker.Client) error {
	imageID, err := checkForOllamaImage(cli)
	if err != nil {
		slog.Error("Failed to check for image", "error", err)
		return err
	}

	if imageID == "" {
		return pullOllamaImage(cli)
	}

	removeContainerImage(cli, imageID)
	return pullOllamaImage(cli)

}

func stopOllamaContainer(cli *docker.Client) {
	ctx := context.Background()

	containerID, running, err := findContainer(cli)
	if err != nil {
		slog.Error("Problem talking to docker")
		return
	}
	if containerID == "" {
		slog.Warn("Ollama container not found!")
		return
	}
	if !running {
		slog.Info("Ollama container is already stopped")
		return
	}

	err = cli.ContainerStop(ctx, containerID, container.StopOptions{})
	if err != nil {
		slog.Error("Failed to stop container", "error", err)
		return
	}

	slog.Info("Stopped Ollama container", "id", containerID)
}

func startOllamaExistingContainer(cli *docker.Client, id string) error {
	ctx := context.Background()

	err := cli.ContainerStart(ctx, id, container.StartOptions{})
	if err != nil {
		slog.Error("Failed to start container", "error", err)
		return err
	}

	slog.Info("Started Ollama container", "id", id)
	return nil
}

// startOllamaFirstTime starts the Ollama container
// cmd: docker run -d --device /dev/kfd --device /dev/dri -v ollama:/root/.ollama -p 11434:11434 --name ollama --restart=always ollama/ollama:rocm
func startOllamaFirstTime(cli *docker.Client) error {
	ctx := context.Background()
	var (
		devices      []container.DeviceMapping
		hostConfig   container.HostConfig
		clientConfig container.Config
	)

	containerPort, err := nat.NewPort("tcp", "11434")
	if err != nil {
		slog.Error("Failed to create port", "error", err)
		return err
	}
	hostBinding := nat.PortBinding{
		HostIP:   "0.0.0.0",
		HostPort: "11434",
	}
	hostConfig.PortBindings = nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}

	// Required for AMD Graphics devices
	device1 := container.DeviceMapping{
		PathOnHost:        "/dev/kfd",
		PathInContainer:   "/dev/kfd",
		CgroupPermissions: "rwm"}
	device2 := container.DeviceMapping{
		PathOnHost:        "/dev/dri",
		PathInContainer:   "/dev/dri",
		CgroupPermissions: "rwm"}
	devices = append(devices, device1, device2)
	hostConfig.Devices = append(hostConfig.Devices, devices...)
	hostConfig.RestartPolicy = container.RestartPolicy{Name: "always"}

	// Required for persistent storage
	hostConfig.Binds = []string{"ollama:/root/.ollama"}
	clientConfig.Volumes = map[string]struct{}{"/root/.ollama": {}}

	clientConfig.Image = ollamaTagRocm

	resp, err := cli.ContainerCreate(
		ctx,
		&clientConfig,
		&hostConfig,
		nil,
		nil,
		ollamaStartName)
	if err != nil {
		slog.Error("Failed to create container", "error", err)
		return err
	}

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		slog.Error("Failed to start container", "error", err)
		return err
	}

	slog.Info("Started Ollama container", "id", resp.ID)
	return nil
}

func startOllamaContainer(cli *docker.Client) (string, error) {
	containerID, running, err := findContainer(cli)
	if err != nil {
		slog.Error("Problem talking to docker")
		os.Exit(1)
	}
	if containerID == "" {
		slog.Warn("Ollama container not found!")
		err = startOllamaFirstTime(cli)
		if err != nil {
			slog.Error("Failed to start container", "error", err)
			return "", err
		}
	}

	if !running {
		err = startOllamaExistingContainer(cli, containerID)
		if err != nil {
			slog.Error("Failed to start image", "error", err)
			os.Exit(1)
		}
	}
	return containerID, nil
}
