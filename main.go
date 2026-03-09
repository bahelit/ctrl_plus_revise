package main

import (
	"embed"
	"log/slog"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"

	"github.com/bahelit/ctrl_plus_revise/internal/config"
	"github.com/bahelit/ctrl_plus_revise/version"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	slog.Info("Starting Ctrl+Revise", "version", version.Version)

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		cfg = config.DefaultConfig()
	}

	application := NewApp(cfg)

	err = wails.Run(&options.App{
		Title:  "Ctrl+Revise AI Text Generator",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 30, G: 30, B: 30, A: 1},
		OnStartup:        application.startup,
		OnShutdown:       application.shutdown,
		Bind: []interface{}{
			application,
		},
		Linux: &linux.Options{
			ProgramName: "ctrl-plus-revise",
		},
	})
	if err != nil {
		slog.Error("Failed to start", "error", err)
	}
}
