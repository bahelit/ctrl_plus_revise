package main

import (
	"context"
	"log/slog"

	"github.com/ollama/ollama/api"
)

type ModelName string

const (
	Llama3 ModelName = "llama3"
)
const (
	correctGrammar     = "Correct the following text for grammar, punctuation, and spelling errors without explaining what changed, just provide the corrected text: "
	makeItAList        = "Transform the previous text into a well-organized, easy-to-read bulleted list, while maintaining key details and phrases. Ensure the list is concise yet informative, making it suitable for a formal report, presentation, or proposal. I don't need an explanation, just provide the list." //nolint:lll
	makeItFriendly     = "Give the following text a friendly makeover by injecting a touch of humor, warmth, and approachability. You don't have to explain your changes, just make the text more friendly: "
	makeItFriendlyRedo = "Give the previous text a friendly makeover by injecting a touch of humor, warmth, and approachability. You don't have to explain your changes, just make the text more friendly."
	makeItProfessional = "Rephrase the following sentence to make it sound more professional and suitable for a formal business document or presentation, only include the rephrased sentence please don't provide an explanation: "
	tryAgain           = "Please try again, but you don't have to tell me it is a another version."
)

//go:generate stringer -type=PromptMsg
type PromptMsg int

const (
	CorrectGrammar PromptMsg = iota // Correct Grammar
	TryAgain
	MakeItAList
	MakeItFriendly     // Make it Friendly
	MakeItFriendlyRedo // Make it Friendly
	MakeItProfessional // Make it Professional
)

func generateResponseFromOllama(client *api.Client, prompt PromptMsg, inputForPrompt string) (api.GenerateResponse, error) {
	var response api.GenerateResponse
	req := &api.GenerateRequest{
		Model:  string(Llama3),
		Prompt: translatePrompt(prompt) + inputForPrompt,
		// set streaming to false
		Stream: new(bool),
	}

	// TODO implement timeout
	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.
		response = resp
		return nil
	}

	err := client.Generate(ctx, req, respFunc)
	if err != nil {
		slog.Error("Failed to generate", "error", err)
		return api.GenerateResponse{}, err
	}

	return response, nil
}

func reGenerateResponseFromOllama(client *api.Client, msgContext []int, prompt PromptMsg) (api.GenerateResponse, error) {
	// TODO How long does the context last?
	var response api.GenerateResponse
	req := &api.GenerateRequest{
		Model:  string(Llama3),
		Prompt: translatePrompt(prompt),
		// set streaming to false
		Stream:  new(bool),
		Context: msgContext,
	}

	// TODO implement timeout
	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.
		response = resp
		return nil
	}

	err := client.Generate(ctx, req, respFunc)
	if err != nil {
		slog.Error("Failed to generate", "error", err)
		return api.GenerateResponse{}, err
	}

	return response, nil
}

func translatePrompt(prompt PromptMsg) string {
	switch prompt {
	case CorrectGrammar:
		return correctGrammar
	case MakeItFriendly:
		return makeItFriendly
	case MakeItFriendlyRedo:
		return makeItFriendlyRedo
	case MakeItAList:
		return makeItAList
	case MakeItProfessional:
		return makeItProfessional
	case TryAgain:
		return tryAgain
	default:
		slog.Warn("Unknown prompt", "prompt", prompt)
		return tryAgain
	}
}
