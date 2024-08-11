package hardware

import (
	"log/slog"

	"github.com/jaypipes/ghw"

	"github.com/bahelit/ctrl_plus_revise/pkg/bytesize"
)

const (
	amdDriver    = "amdgpu"
	nvidiaDriver = "nvidia"
)

// GPU Processing devices used for AI inference
//
//go:generate stringer -linecomment -type=GPU
type GPU int

const (
	AMD    GPU = iota // AMD
	NVIDIA            // Nvidia
	noGPU             // CPU
)

// DetectMemory detects the memory of the system
// TODO: Use to determine if we can run certain AI models
func DetectMemory() {
	ram, err := ghw.Memory()
	if err != nil {
		slog.Info("Error getting Memory info", "error", err)
	}

	slog.Debug("Memory", "ram", ram.String())

	total := bytesize.New(float64(ram.TotalPhysicalBytes))
	usable := bytesize.New(float64(ram.TotalUsableBytes))

	slog.Info("System Memory", "Total", total, "Available", usable)
}

func DetectProcessingDevice() GPU {
	gpu, err := ghw.GPU()
	if err != nil {
		slog.Info("Error getting GPU info", "error", err)
	}

	foundGPU := noGPU

	for _, card := range gpu.GraphicsCards {
		// TODO: test on Nvidia system
		if card != nil && card.DeviceInfo != nil {
			slog.Info("GPU Probe", "Driver", card.DeviceInfo.Driver, "Product", card.DeviceInfo.Product.Name)
		} else {
			slog.Warn("Failed to get GPU info")
			continue
		}
		if card.DeviceInfo.Driver == amdDriver {
			slog.Info("Detected AMD GPU", "Product", card.DeviceInfo.Product.Name)
			foundGPU = AMD
			break
		} else if card.DeviceInfo.Driver == nvidiaDriver {
			slog.Info("Detected Nvidia GPU", "Product", card.DeviceInfo.Product.Name)
			foundGPU = NVIDIA
			break
		}
	}

	switch foundGPU {
	case AMD:
		return AMD
	case NVIDIA:
		return NVIDIA
	default:
		slog.Info("No GPU Detected, using CPU")
		return noGPU
	}
}
