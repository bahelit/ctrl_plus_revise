package docker

import (
	"context"
	"github.com/bahelit/ctrl_plus_revise/internal/hardware"
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
	ollamaTagNvidia     = "ollama/ollama"
	ollamaTagRocm       = "ollama/ollama:rocm"
)

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
	var (
		reader io.ReadCloser
		err    error
	)

	// TODO: Allow user to specify image AMD or Nvidia images
	computeDevice := hardware.DetectProcessingDevice()
	if computeDevice == hardware.AMD {
		reader, err = cli.ImagePull(ctx, ollamaTagRocm, image.PullOptions{})
	} else {
		reader, err = cli.ImagePull(ctx, ollamaTagNvidia, image.PullOptions{})
	}
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

func updateOllamaImage() error {
	dockerClient, err := connectToDocker()
	if err != nil {
		slog.Error("Failed to connect to Docker", "error", err)
		return err
	}
	imageID, err := checkForOllamaImage(dockerClient)
	if err != nil {
		slog.Error("Failed to check for image", "error", err)
		return err
	}

	if imageID == "" {
		return pullOllamaImage(dockerClient)
	}

	removeContainerImage(dockerClient, imageID)
	return pullOllamaImage(dockerClient)

}

func StopOllamaContainer() {
	ctx := context.Background()
	dockerClient, err := connectToDocker()
	if err != nil {
		slog.Warn("Failed to connect to Docker", "error", err)
		return
	}

	containerID, running, err := findContainer(dockerClient)
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

	err = dockerClient.ContainerStop(ctx, containerID, container.StopOptions{})
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
// AMD cmd: docker run -d --device /dev/kfd --device /dev/dri -v ollama:/root/.ollama -p 11434:11434 --name ollama --restart=always ollama/ollama:rocm
// NVIDIA cmd: docker run -d --gpus=all -v ollama:/root/.ollama -p 11434:11434 --restart=always --name ollama ollama/ollama
// CPU cmd: docker run -d -v ollama:/root/.ollama -p 11434:11434 --restart=always --name ollama ollama/ollama
// https://hub.docker.com/r/ollama/ollama
func startOllamaFirstTime(cli *docker.Client) error {
	ctx := context.Background()
	var (
		devices      []container.DeviceMapping
		hostConfig   container.HostConfig
		clientConfig container.Config
	)

	computeDevice := hardware.DetectProcessingDevice()

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

	if computeDevice == hardware.AMD {
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
		clientConfig.Image = ollamaTagRocm
	} else if computeDevice == hardware.NVIDIA {
		err = os.Setenv("NVIDIA_VISIBLE_DEVICES", "all")
		if err != nil {
			slog.Error("Failed to set NVIDIA_VISIBLE_DEVICES", "error", err)
		}
		err = os.Setenv("NVIDIA_DRIVER_CAPABILITIES", "compute,utility")
		if err != nil {
			slog.Error("Failed to set NVIDIA_DRIVER_CAPABILITIES", "error", err)
		}
	}
	clientConfig.Image = ollamaTagNvidia
	hostConfig.RestartPolicy = container.RestartPolicy{Name: "always"}

	// Required for persistent storage
	hostConfig.Binds = []string{"ollama:/root/.ollama"}
	clientConfig.Volumes = map[string]struct{}{"/root/.ollama": {}}

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
	ollamaContainerID, running, err := findContainer(cli)
	if err != nil {
		slog.Error("Problem talking to docker")
		os.Exit(1)
	}
	if ollamaContainerID == "" {
		slog.Warn("Ollama container not found!")
		err = startOllamaFirstTime(cli)
		if err != nil {
			slog.Error("Failed to start container", "error", err)
			return "", err
		}
	}

	if !running {
		err = startOllamaExistingContainer(cli, ollamaContainerID)
		if err != nil {
			slog.Error("Failed to start image", "error", err)
			os.Exit(1)
		}
	}
	return ollamaContainerID, nil
}
