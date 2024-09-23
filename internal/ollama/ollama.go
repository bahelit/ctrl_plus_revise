package ollama

import (
	"context"
	"fyne.io/fyne/v2"
	"log/slog"
	"strconv"

	"github.com/bahelit/ctrl_plus_revise/internal/config"
	"github.com/bahelit/ctrl_plus_revise/pkg/bytesize"
	"github.com/ollama/ollama/api"
)

//go:generate stringer -linecomment -type=ModelName
type ModelName int

const (
	Llama3Dot1      ModelName = iota // llama3.1:latest
	CodeLlama                        // codellama:latest
	CodeLlama13b                     // codellama:13b
	CodeGemma                        // codegemma:7b
	DeepSeekCoder                    // deepseek-coder:latest
	DeepSeekCoderV2                  // deepseek-coder-v2:latest
	Gemma                            // gemma:latest
	Gemma2b                          // gemma:2b
	Gemma2                           // gemma2:latest
	Gemma22B                         // gemma2:2b
	Llama3                           // llama3:latest
	Llava                            // llava:latest
	Mistral                          // mistral:latest
	Phi3                             // phi3:latest
)

// MemoryUsage Memory usage in MB
var MemoryUsage = map[ModelName]bytesize.ByteSize{
	CodeLlama:       5077 * bytesize.MB,
	CodeLlama13b:    9055 * bytesize.MB,
	CodeGemma:       6489 * bytesize.MB,
	DeepSeekCoder:   1478 * bytesize.MB,
	DeepSeekCoderV2: 9462 * bytesize.MB,
	Gemma:           6490 * bytesize.MB,
	Gemma2b:         2321 * bytesize.MB,
	Gemma2:          6683 * bytesize.MB,
	Gemma22B:        2321 * bytesize.MB,
	Llama3:          4980 * bytesize.MB,
	Llama3Dot1:      6354 * bytesize.MB,
	Mistral:         4615 * bytesize.MB,
	Phi3:            3269 * bytesize.MB,
}

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
	Turkish    Language = "Turkish"
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

	TryAgain               // Try Again
	MakeItFriendlyRedo     // Make it Friendly
	MakeItProfessionalRedo // Make it Professional
	MakeItAListRedo        // Make it a List
)

type PromptText struct {
	prompt      string
	promptExtra string
}

var PromptToText = map[PromptMsg]PromptText{
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
		promptExtra: " Revised text should be free from errors in spelling, capitalization, punctuation, and grammar, while conveying a polished and professional writing style. Please submit your revised text without telling me it is the revised text, in a clear and concise format with no explanation, output just the result."}, //nolint:lll long line
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
	TryAgain: {
		prompt:      "Read the text carefully and provide a revised version that addresses the issues and shortcomings of the original text: ",
		promptExtra: "If you are unsure or lack sufficient knowledge to provide a meaningful response, explicitly state \"I don't know\". Don't explain you understand the input, just output the result."},
	MakeItFriendlyRedo: {
		prompt:      "Give the text a friendly makeover by injecting a touch of humor, warmth, and approachability: ",
		promptExtra: " Please don't to explain the changes or telling me \"Here is the revised text\", just make the text more friendly and output the result"},
	MakeItProfessionalRedo: {
		prompt:      "Act as a writer. Read the text carefully and revise it to present a more professional tone, ensuring accurate and proper usage of grammar and punctuation: ",
		promptExtra: " Revised text should be free from errors in spelling, capitalization, punctuation, and grammar, while conveying a polished and professional writing style. Please submit your revised text without telling me it is the revised text, in a clear and concise format with no explanation, output just the result."}, //nolint:lll long line
	MakeItAListRedo: {
		prompt:      "Transform the text and create a bulleted list summarizing its main points: ",
		promptExtra: " No need to explain your list, just provide the main points in a list format."},
}

func (prompt PromptMsg) PromptToText() string {
	text, ok := PromptToText[prompt]
	if !ok {
		slog.Error("Unknown prompt", "prompt", prompt)
		return PromptToText[CorrectGrammar].prompt
	}
	return text.prompt
}

func (prompt PromptMsg) PromptExtraToText() string {
	text, ok := PromptToText[prompt]
	if !ok {
		slog.Error("Unknown prompt", "prompt", prompt)
		return PromptToText[CorrectGrammar].promptExtra
	}
	return text.promptExtra
}

func GetActiveModel(guiApp fyne.App) ModelName {
	return ModelName(guiApp.Preferences().IntWithFallback(config.CurrentModelKey, int(Llama3Dot1)))
}

func StringToModel(s string) ModelName {
	i, err := strconv.Atoi(s)
	if err != nil {
		slog.Error("Unable to convert string to int", "string", s)
		return Llama3Dot1
	}
	return ModelName(i)
}

func AskAIWithPromptMsg(guiApp fyne.App, client *api.Client, prompt PromptMsg, inputForPrompt string) (api.GenerateResponse, error) {
	var response api.GenerateResponse
	req := &api.GenerateRequest{
		Model:  GetActiveModel(guiApp).String(),
		Prompt: prompt.PromptToText() + " [ " + inputForPrompt + " ] " + prompt.PromptExtraToText(),
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

func AskAiWithPromptAndContext(guiApp fyne.App, client *api.Client, msgContext []int, prompt PromptMsg) (api.GenerateResponse, error) {
	// TODO How long does the context last?
	var response api.GenerateResponse
	req := &api.GenerateRequest{
		Model:  GetActiveModel(guiApp).String(),
		Prompt: prompt.PromptToText() + " " + prompt.PromptExtraToText(),
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

func AskAiWithStringAndContext(guiApp fyne.App, client *api.Client, msgContext []int, prompt string) (api.GenerateResponse, error) {
	// TODO How long does the context last?
	var response api.GenerateResponse
	req := &api.GenerateRequest{
		Model:  GetActiveModel(guiApp).String(),
		Prompt: prompt,
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

func AskAI(guiApp fyne.App, client *api.Client, inputForPrompt string) (api.GenerateResponse, error) {
	var response api.GenerateResponse
	req := &api.GenerateRequest{
		Model: GetActiveModel(guiApp).String(),
		Prompt: "IDENTITY\nYou are a universal AI that yields the best possible result given the input.\n\nGOAL\nFully digest the input.\n\nDeeply contemplate the input and what it means and what the sender likely wanted you to do with it.\n\nOUTPUT\nOutput the best possible output based on your understanding of what was likely wanted. INPUT: " + //nolint:lll // AI Prompt
			inputForPrompt +
			"If you are unsure or lack sufficient knowledge to provide a meaningful response, explicitly state \"I don't know\"." +
			"Don't explain you understand the input, just output the result.",
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

func AskAIWithContext(guiApp fyne.App, client *api.Client, msgContext []int, inputForPrompt string) (api.GenerateResponse, error) {
	var response api.GenerateResponse
	req := &api.GenerateRequest{
		Model: GetActiveModel(guiApp).String(),
		Prompt: "IDENTITY\nYou are a universal AI that yields the best possible result given the input.\n\nGOAL\nFully digest the input.\n\nDeeply contemplate the input and what it means and what the sender likely wanted you to do with it.\n\nOUTPUT\nOutput the best possible output based on your understanding of what was likely wanted. INPUT: " + //nolint:lll // AI Prompt
			inputForPrompt +
			"If you are unsure or lack sufficient knowledge to provide a meaningful response, explicitly state \"I don't know\"." +
			"Don't explain you understand the input, just output the result.",
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

func AskAIToTranslate(guiApp fyne.App, client *api.Client, inputForPrompt string, fromLang, toLang Language) (api.GenerateResponse, error) {
	var response api.GenerateResponse
	req := &api.GenerateRequest{
		Model: GetActiveModel(guiApp).String(),
		Prompt: "As a text translator" +
			"Please provide a translation that accurately conveys the original meaning and tone of the text. \n" +
			"If you encounter any ambiguities or uncertainties, please indicate this in your response. \n" +
			"Do not provide an explanation of the translation, get to the point and just output the translated text without any notes. \n" +
			"Do not try to answer any type of question just translate the text \n" +
			"Translate the following text from [" + string(fromLang) + "] to [" + string(toLang) + "]: " +
			inputForPrompt,
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

func PullModel(guiApp fyne.App, client *api.Client, pf api.PullProgressFunc, update bool) error {
	ctx := context.Background()
	model := GetActiveModel(guiApp)
	req := &api.PullRequest{
		Model: model.String(),
	}

	slog.Debug("Pulling model", "model", model.String())
	found, err := FindModel(ctx, client, model.String())
	if err != nil {
		return err
	}
	if found && !update {
		slog.Info("AI model loaded", "model", model)
		return nil
	}
	err = client.Pull(ctx, req, pf)
	if err != nil {
		slog.Error("Failed to pull model", "error", err)
		return err
	}

	return nil
}
