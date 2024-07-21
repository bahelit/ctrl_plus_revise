package main

import (
	"log/slog"

	"github.com/jaypipes/ghw"

	"github.com/bahelit/ctrl_plus_revise/pkg/bytesize"
)

const (
	amdDriver    = "amdgpu"
	nvidiaDriver = "nvidia"
)

// processingDevice Processing devices used for AI inference
//
//go:generate stringer -linecomment -type=processingDevice
type processingDevice int

const (
	amdGPU    processingDevice = iota // AMD
	nvidiaGPU                         // Nvidia
	noGPU                             // CPU
)

// gpuModel GPUs that have been tested
//
//go:generate stringer -linecomment -type=gpuModel
type gpuModel int

const (
	amdProductNavi10 gpuModel = iota // Navi 10 [Radeon RX 5600 OEM/5600 XT / 5700/5700 XT]
	amdProductNavi21                 // Navi 21 [Radeon RX 6800/6800 XT / 6900 XT]
)

// DetectMemory detects the memory of the system
// TODO: Use to determine if we can run certain AI models
func detectMemory() {
	ram, err := ghw.Memory()
	if err != nil {
		slog.Info("Error getting Memory info", "error", err)
	}

	slog.Debug("Memory", "ram", ram.String())

	total := bytesize.New(float64(ram.TotalPhysicalBytes))
	usable := bytesize.New(float64(ram.TotalUsableBytes))

	slog.Info("System Memory", "Total", total, "Available", usable)
}

func detectProcessingDevice() processingDevice {
	gpu, err := ghw.GPU()
	if err != nil {
		slog.Info("Error getting GPU info", "error", err)
	}

	foundGPU := noGPU

	for _, card := range gpu.GraphicsCards {
		// TODO: test on Nvidia system
		if card != nil && card.DeviceInfo != nil {
			slog.Debug("GPU Probe", "Driver", card.DeviceInfo.Driver, "Product", card.DeviceInfo.Product.Name)
		} else {
			slog.Warn("Failed to get GPU info")
			continue
		}
		if card.DeviceInfo.Driver == amdDriver {
			slog.Info("Detected AMD GPU", "Product", card.DeviceInfo.Product.Name)
			foundGPU = amdGPU
			break
		} else if card.DeviceInfo.Driver == nvidiaDriver {
			slog.Info("Detected AMD GPU", "Product", card.DeviceInfo.Product.Name)
			foundGPU = nvidiaGPU
			break
		}
	}

	switch foundGPU {
	case amdGPU:
		return amdGPU
	case nvidiaGPU:
		return nvidiaGPU
	default:
		slog.Info("No GPU Detected, using CPU")
		return noGPU
	}
}
