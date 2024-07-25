package ollama

import (
	"context"
	"log/slog"

	ollama "github.com/ollama/ollama/api"
)

func FindModel(ctx context.Context, client *ollama.Client, model string) (bool, error) {
	response, err := client.List(ctx)
	if err != nil {
		slog.Error("Failed to pull model", "error", err)
		return false, err
	}
	for m := range response.Models {
		slog.Debug("Docker Image", "Name", response.Models[m], "Model", response.Models[m].Model,
			"ParameterSize", response.Models[m].Details.ParameterSize, "Families", response.Models[m].Details.Families)
		if response.Models[m].Name == model {
			slog.Debug("Model found", "model", response.Models[m].Model)
			return true, nil
		}
	}
	return false, nil
}
