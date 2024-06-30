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

//go:generate stringer -linecomment -type=PromptMsg
type PromptMsg int

const (
	CorrectGrammar PromptMsg = iota // Correct Grammar

	MakeItAList        // Make it a List
	MakeItFriendly     // Make it Friendly
	MakeItFriendlyRedo // Make it Friendly
	MakeItProfessional // Make it Professional
	MakeASummary       // Make a Summary
	MakeExplanation    // Make an Explanation
	MakeExpanded       // Expand on the text
	MakeHeadline       // Make a Headline
)

var promptToText = map[PromptMsg]string{
	CorrectGrammar:     "Act as a writer. Correct the following text for grammar, punctuation, and spelling errors without explaining what changed, just provide the corrected text: ",
	MakeItFriendly:     "Act as a writer. Give the following text a friendly makeover by injecting a touch of humor, warmth, and approachability. You don't have to explain your changes, just make the text more friendly: ",
	MakeItFriendlyRedo: "Act as a writer. Give the previous text a friendly makeover by injecting a touch of humor, warmth, and approachability. You don't have to explain your changes, just make the text more friendly.",
	MakeItAList:        "Act as a writer. Transform the previous text into a well-organized, easy-to-read bulleted list, while maintaining key details and phrases. Ensure the list is concise yet informative, making it suitable for a formal report, presentation, or proposal. Do not provide an explanation, just provide the list.", //nolint:lll - long line
	MakeItProfessional: "Act as a writer. Rephrase the following sentence to make it sound more professional and suitable for a formal business document or presentation, only include the rephrased sentence please don't provide an explanation: ",
	MakeHeadline:       "Act as a writer. Craft a compelling headline that captures the essence of the following block of text. Keep it concise, attention-grabbing, and accurate. Output only the text and nothing else, do not chat, no preamble, get to the point. ",
	MakeASummary:       "Act as a writer. Summarize the following block of text in 50-75 words, focusing on the main ideas and key points. Use your own words to condense the information without sacrificing clarity or accuracy. Output only the text and nothing else, do not chat, no preamble, get to the point. ", //nolint:lll - long line
	MakeExplanation:    "Act as a writer. Explain the following block of text in 100-150 words, providing context, clarifying complex concepts, and making connections to related ideas. Output only the text and nothing else, do not chat, no preamble, get to the point. ",
	MakeExpanded:       "Act as a writer. Expand the text by adding more details while keeping the same meaning. Output only the text and nothing else, do not chat, no preamble, get to the point. ",
}

func (prompt PromptMsg) toText() string {
	text, ok := promptToText[prompt]
	if !ok {
		slog.Error("Unknown prompt", "prompt", prompt)
		return promptToText[CorrectGrammar]
	}
	return text
}

func askAIWithPromptMsg(client *api.Client, prompt PromptMsg, inputForPrompt string) (api.GenerateResponse, error) {
	var response api.GenerateResponse
	req := &api.GenerateRequest{
		Model:  string(Llama3),
		Prompt: prompt.toText() + inputForPrompt,
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

func askAI(client *api.Client, inputForPrompt string) (api.GenerateResponse, error) {
	var response api.GenerateResponse
	req := &api.GenerateRequest{
		Model: string(Llama3),
		Prompt: "My name is Ctrl+Revise and I am an AI agent running on a personal computer that doesn't require internet access," +
			"Provide information or answers for the following questions based on my training data. " +
			"If I'm unsure or don't have enough information, please indicate this clearly." +
			inputForPrompt +
			"Note: Please only provide information that you are confident about and acknowledge any limitations or uncertainties in your response. " +
			"Output only the answer and nothing else, do not chat, no preamble, get to the point.",
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
		Prompt: prompt.toText(),
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
