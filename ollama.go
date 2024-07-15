package main

import (
	"context"
	"log/slog"

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

type Language string

const (
	Arabic     Language = "Arabic"
	Chinese    Language = "Chinese"
	English    Language = "English"
	French     Language = "French"
	German     Language = "German"
	Italian    Language = "Italian"
	Japanese   Language = "Japanese"
	Portuguese Language = "Portuguese"
	Russian    Language = "Russian"
	Spanish    Language = "Spanish"
)

//go:generate stringer -linecomment -type=PromptMsg
type PromptMsg int

const (
	CorrectGrammar PromptMsg = iota // Correct Grammar

	MakeItAList        // Make it a List
	MakeItFriendly     // Make it Friendly
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
		prompt:      "IDENTITY and PURPOSE\nYou are a writing expert. You refine the input text to enhance clarity, coherence, grammar, and style.\n\nSteps\nAnalyze the input text for grammatical errors, stylistic inconsistencies, clarity issues, and coherence.\nApply corrections and improvements directly to the text.\nMaintain the original meaning and intent of the user's text, ensuring that the improvements are made within the context of the input language's grammatical norms and stylistic conventions.\nOUTPUT INSTRUCTIONS\nRefined and improved text that has no grammar mistakes.\nReturn in the same language as the input.\nInclude NO additional commentary or explanation in the response.\nINPUT:", //nolint:lll long line
		promptExtra: " Return the corrected text without explaining what changed or telling me \"Here is the revised text\", just provide the corrected text and output just the result"},
	MakeItFriendly: {
		prompt:      "Give the following text a friendly makeover by injecting a touch of humor, warmth, and approachability: ",
		promptExtra: " Please don't to explain the changes or telling me \"Here is the revised text\", just make the text more friendly and output the result"},
	MakeItAList: {
		prompt:      "Read the following text and create a bulleted list summarizing its main points: ",
		promptExtra: " No need to explain your list, just provide the main points in a list format."},
	MakeItProfessional: {
		prompt:      "Act as a writer. Read the following text carefully and revise it to present a more professional tone, ensuring accurate and proper usage of grammar and punctuation: ",
		promptExtra: " Revised text should be free from errors in spelling, capitalization, punctuation, and grammar, while conveying a polished and professional writing style. Please submit your revised text without telling me it is the revised text, in a clear and concise format with no explanation, output just the result."},
	MakeHeadline: {
		prompt:      "Act as a writer. Read the following text carefully and create a concise and attention-grabbing headline that summarizes its main idea or key point: ",
		promptExtra: " Your headline should be no more than 5-7 words, yet effectively capture the essence of the text. Please submit your headline in the format below:\n\n[Headline]"},
	MakeASummary: {
		prompt:      "IDENTITY and PURPOSE\nYou are a summarization system that extracts the most interesting, useful, and surprising aspects of an article.\n\nTake a step back and think step by step about how to achieve the best result possible as defined in the steps below. You have a lot of freedom to make this work well.\n\nOUTPUT SECTIONS\nYou extract a summary of the content in 20 words or less, including who is presenting and the content being discussed into a section called SUMMARY.\n\nYou extract the top 20 ideas from the input in a section called IDEAS:.\n\nYou extract the 10 most insightful and interesting quotes from the input into a section called QUOTES:. Use the exact quote text from the input.\n\nYou extract the 20 most insightful and interesting recommendations that can be collected from the content into a section called RECOMMENDATIONS.\n\nYou combine all understanding of the article into a single, 20-word sentence in a section called ONE SENTENCE SUMMARY:.\n\nOUTPUT INSTRUCTIONS\nYou only output Markdown.\nDo not give warnings or notes; only output the requested sections.\nYou use numbered lists, not bullets.\nDo not repeat ideas, quotes, facts, or resources.\nDo not start items with the same opening words.\nDo not include any commentary or explanation.\n\nINPUT:", //nolint:lll long line
		promptExtra: ""},
	MakeExplanation: {
		prompt:      "Explain the following block of text in a way that a 5-year-old could understand. Use simple language, relatable examples, and avoid technical jargon: ",
		promptExtra: "Goals: Simplify complex ideas into easy-to-grasp concepts. Use analogies or relatable scenarios to help explain abstract concepts. Make it fun and engaging while still being accurate"},
	MakeExpanded: {
		prompt: "Read the following text carefully and determine its nature: does it appear to be based on factual information or is it fictional in nature?: ",
		promptExtra: " If the text appears to be non-fictional in nature, expand on it by incorporating relevant, accurate, and verifiable information from credible sources. " +
			"Ensure that all additional information is properly sourced and attributed to credible sources.\n\n\nHowever, if the text appears to be fictional in nature, feel free to expand on it by adding to the story, developing characters, or exploring themes. " +
			"Please keep your additions consistent with the tone and style of the original text. " +
			"Please do not provide any commentary or explanation, just expand on the text."},
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
		Prompt: prompt.promptToText() + " [ " + inputForPrompt + " ] " + prompt.promptExtraToText(),
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
		Prompt: "IDENTITY\nYou are a universal AI that yields the best possible result given the input.\n\nGOAL\nFully digest the input.\n\nDeeply contemplate the input and what it means and what the sender likely wanted you to do with it.\n\nOUTPUT\nOutput the best possible output based on your understanding of what was likely wanted. INPUT: " + //nolint:lll long line
			inputForPrompt +
			"If you are unsure or lack sufficient knowledge to provide a meaningful response, explicitly state \"I don't know\"." +
			"Don't explain you understand\n" +
			"The output should be in markdown format.",
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

func askAIToTranslate(client *api.Client, model ModelName, inputForPrompt string, fromLang, toLang Language) (api.GenerateResponse, error) {
	var response api.GenerateResponse
	req := &api.GenerateRequest{
		Model: model.String(),
		Prompt: "Translate the following text from [" + string(fromLang) + "] to [" + string(toLang) + "]: " +
			inputForPrompt +
			"Please provide a translation that accurately conveys the original meaning and tone of the text.\n" +
			"If you encounter any ambiguities or uncertainties, please indicate this in your response.\n" +
			"Do not provide an explanation of the translation, get to the point and just output the translated text " +
			"without any notes.",
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

func pullModel(model ModelName, update bool) error {
	ctx := context.Background()
	req := &api.PullRequest{
		Model: model.String(),
	}

	slog.Debug("Pulling model", "model", model.String())
	found, err := findModel(ctx, model)
	if err != nil {
		return err
	}
	if found && !update {
		slog.Info("AI model loaded", "model", model)
		return nil
	}

	progressFunc := func(resp api.ProgressResponse) error {
		slog.Info("Progress", "status", resp.Status, "total", resp.Total, "completed", resp.Completed)
		if resp.Total == resp.Completed {
			slog.Info("Model pulled", "model", model, "resp", resp)
		}
		return nil
	}

	err = ollamaClient.Pull(ctx, req, progressFunc)
	if err != nil {
		slog.Error("Failed to pull model", "error", err)
		return err
	}
	showNotification("Model Download Completed", "Model "+model.String()+" has been pulled")
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
