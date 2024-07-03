package main

import (
	"context"
	"log/slog"

	"fyne.io/fyne/v2"
	"github.com/ollama/ollama/api"
)

//go:generate stringer -linecomment -type=ModelName
type ModelName int

const (
	BashBot      ModelName = iota // bashbot:latest
	CodeLlama                     // codellama:latest
	CodeLlama13b                  // codellama:13b
	Gemma                         // gemma:latest
	Llama3                        // llama3:latest
)

//go:generate stringer -linecomment -type=PromptMsg
type PromptMsg int

const (
	CorrectGrammar PromptMsg = iota // Correct Grammar

	MakeItAList        // Make it a List
	MakeItFriendly     // Make it Friendly
	MakeItFriendlyRedo // Make it Friendly
	MakeItProfessional // Make it Professional
	MakeASummary       // Make a Summary
	MakeExplanation    // Explain it like I'm 5
	MakeExpanded       // Expand on the text
	MakeHeadline       // Make a Headline
)

type promptText struct {
	prompt      string
	promptExtra string
}

var promptToText = map[PromptMsg]promptText{
	CorrectGrammar: {
		prompt:      "Act as a writer. Correct the following block of text by fixing all spelling, grammar, punctuation, and capitalization errors. Provide an error-free version of the original text: ",
		promptExtra: " Return the corrected text without explaining what changed, just provide the corrected text"},
	MakeItFriendly: {
		prompt:      "Give the following text a friendly makeover by injecting a touch of humor, warmth, and approachability: ",
		promptExtra: " You don't have to explain your changes, just make the text more friendly"},
	MakeItFriendlyRedo: {
		prompt:      "Act as a writer. Give the previous text a friendly makeover by injecting a touch of humor, warmth, and approachability. ",
		promptExtra: " You don't have to explain your changes, just make the text more friendly"},
	MakeItAList: {
		prompt:      "Read the following text and create a bulleted list summarizing its main points: ",
		promptExtra: " No need to explain your list, just provide the main points in a list format."},
	MakeItProfessional: {
		prompt:      "Act as a writer. Read the following text carefully and revise it to present a more professional tone, ensuring accurate and proper usage of grammar and punctuation: ",
		promptExtra: " Revised text should be free from errors in spelling, capitalization, punctuation, and grammar, while conveying a polished and professional writing style. Please submit your revised text in a clear and concise format with no explanation."},
	MakeHeadline: {
		prompt:      "Act as a writer. Read the following text carefully and create a concise and attention-grabbing headline that summarizes its main idea or key point: ",
		promptExtra: " Your headline should be no more than 5-7 words, yet effectively capture the essence of the text. Please submit your headline in the format below:\n\n[Headline]"},
	MakeASummary: {
		prompt:      "Act as a writer. Read the following text carefully and create a brief summary that captures the main points and essential information: ",
		promptExtra: "Your summary should be no more than 150-200 words, yet effectively convey the key ideas and takeaways from the original text. Please submit your summary in a clear and concise format with no explanation."}, //nolint:lll - long line
	MakeExplanation: {
		prompt:      "Explain the following block of text in a way that a 5-year-old could understand. Use simple language, relatable examples, and avoid technical jargon: ",
		promptExtra: "Goals: Simplify complex ideas into easy-to-grasp concepts. Use analogies or relatable scenarios to help explain abstract concepts. Make it fun and engaging while still being accurate"},
	MakeExpanded: {
		prompt: "Read the following text carefully and determine its nature: does it appear to be based on factual information or is it fictional in nature?: ",
		promptExtra: " If the text appears to be non-fictional in nature, expand on it by incorporating relevant, accurate, and verifiable information from credible sources. " +
			"Ensure that all additional information is properly sourced and attributed to credible sources.\n\n\nHowever, if the text appears to be fictional in nature, feel free to expand on it by adding to the story, developing characters, or exploring themes. " +
			"Please keep your additions consistent with the tone and style of the original text."},
}

func (prompt PromptMsg) promptToText() string {
	text, ok := promptToText[prompt]
	if !ok {
		slog.Error("Unknown prompt", "prompt", prompt)
		return promptToText[CorrectGrammar].prompt
	}
	return text.prompt
}

func (prompt PromptMsg) promptExtraToText() string {
	text, ok := promptToText[prompt]
	if !ok {
		slog.Error("Unknown prompt", "prompt", prompt)
		return promptToText[CorrectGrammar].promptExtra
	}
	return text.promptExtra
}

func askAIWithPromptMsg(client *api.Client, prompt PromptMsg, model ModelName, inputForPrompt string) (api.GenerateResponse, error) {
	var response api.GenerateResponse
	req := &api.GenerateRequest{
		Model:  model.String(),
		Prompt: prompt.promptToText() + "[" + inputForPrompt + "]" + prompt.promptExtraToText(),
		// set streaming to false
		Stream: new(bool),
	}

	// TODO: implement timeout
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

func askAI(client *api.Client, model ModelName, inputForPrompt string) (api.GenerateResponse, error) {
	var response api.GenerateResponse
	req := &api.GenerateRequest{
		Model: model.String(),
		Prompt: "My name is Ctrl+Revise and I am an AI agent running on a personal computer that doesn't require internet access," +
			"Provide information or answers for the following questions based on my training data. " +
			"If I'm unsure or don't have enough information, please indicate this clearly." +
			inputForPrompt +
			"Note: Please only provide information that you are confident about and acknowledge any limitations or uncertainties in your response. " +
			"Output only the answer and nothing else, do not chat, no preamble, get to the point.",
		// set streaming to false
		Stream: new(bool),
	}

	// TODO: implement timeout
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
		Prompt: prompt.promptToText(),
		// set streaming to false
		Stream:  new(bool),
		Context: msgContext,
	}

	// TODO: implement timeout
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

func pullModel(model ModelName, update bool) error {
	ctx := context.Background()
	req := &api.PullRequest{
		Model: model.String(),
	}

	slog.Info("Pulling model", "model", model.String())
	found, err := findModel(ctx, model)
	if err != nil {
		return err
	}
	if found && !update {
		slog.Info("Model already exists", "model", model)
		return nil
	}

	progressFunc := func(resp api.ProgressResponse) error {
		slog.Info("Progress", "status", resp.Status, "total", resp.Total, "completed", resp.Completed)
		if resp.Total == resp.Completed {
			guiApp.SendNotification(&fyne.Notification{
				Title:   "Model Download Completed",
				Content: "Model " + model.String() + " has been pulled",
			})
		}
		return nil
	}

	err = ollamaClient.Pull(ctx, req, progressFunc)
	if err != nil {
		slog.Error("Failed to pull model", "error", err)
		return err
	}
	return nil
}

func findModel(ctx context.Context, model ModelName) (bool, error) {
	response, err := ollamaClient.List(ctx)
	if err != nil {
		slog.Error("Failed to pull model", "error", err)
		return false, err
	}
	for m := range response.Models {
		slog.Debug("Docker Image", "Name", response.Models[m], "Model", response.Models[m].Model,
			"ParameterSize", response.Models[m].Details.ParameterSize, "Families", response.Models[m].Details.Families)
		if response.Models[m].Name == model.String() {
			slog.Debug("Model found", "model", response.Models[m].Model)
			return true, nil
		}
	}
	return false, nil
}
