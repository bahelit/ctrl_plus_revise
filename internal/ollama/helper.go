package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

// modelEntry is a minimal model descriptor used when parsing /api/tags.
// It intentionally uses a string for modified_at to avoid time-parse failures
// when the Ollama server returns an empty value for that field.
type modelEntry struct {
	Name string `json:"name"`
}

type tagsResponse struct {
	Models []modelEntry `json:"models"`
}

// ListModels returns the names of all models installed in Ollama.
// It calls /api/tags directly to avoid a time-parsing bug in the upstream
// client library when modified_at is empty.
func ListModels(ctx context.Context, baseURL string) ([]modelEntry, error) {
	url := strings.TrimRight(baseURL, "/") + "/api/tags"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Failed to list models", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	var result tagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("Failed to decode model list", "error", err)
		return nil, err
	}
	return result.Models, nil
}

// FindModel reports whether the named model is installed in Ollama.
func FindModel(ctx context.Context, baseURL, model string) (bool, error) {
	models, err := ListModels(ctx, baseURL)
	if err != nil {
		return false, err
	}
	for _, m := range models {
		if m.Name == model {
			return true, nil
		}
	}
	return false, nil
}
