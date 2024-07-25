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

// gpu Processing devices used for AI inference
//
//go:generate stringer -linecomment -type=gpu
type gpu int

const (
	amdGPU    gpu = iota // AMD
	nvidiaGPU            // Nvidia
	noGPU                // CPU
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

func detectProcessingDevice() gpu {
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
			foundGPU = amdGPU
			break
		} else if card.DeviceInfo.Driver == nvidiaDriver {
			slog.Info("Detected Nvidia GPU", "Product", card.DeviceInfo.Product.Name)
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
