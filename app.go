package main

import (
	"context"
	"log/slog"
	"os/exec"

	ollamaApi "github.com/ollama/ollama/api"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/bahelit/ctrl_plus_revise/internal/config"
	"github.com/bahelit/ctrl_plus_revise/internal/meal"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
	"github.com/bahelit/ctrl_plus_revise/internal/shortcuts"
	"github.com/bahelit/ctrl_plus_revise/internal/store/database"
	chatModel "github.com/bahelit/ctrl_plus_revise/internal/store/models/chat"
)

// App is the main Wails application struct.  All exported methods are bound to
// the JavaScript frontend automatically by Wails.
type App struct {
	ctx          context.Context
	cfg          *config.AppConfig
	ollamaClient *ollamaApi.Client
	db           *database.ChatBot
	shortcutMgr  *shortcuts.Manager
}

// NewApp constructs the App using the provided config.
func NewApp(cfg *config.AppConfig) *App {
	return &App{cfg: cfg}
}

// startup is called by the Wails runtime when the application launches.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Connect to Ollama.
	var connectURL *string
	if a.cfg.UseRemoteOllama && a.cfg.OllamaURL != "" {
		connectURL = &a.cfg.OllamaURL
	}
	client, err := ollama.CheckOllamaConnection(connectURL)
	if err != nil {
		slog.Error("Failed to connect to Ollama on startup", "error", err)
	} else {
		a.ollamaClient = client
	}

	// Initialize SQLite database.
	db, err := database.NewSQLiteDB()
	if err != nil {
		slog.Error("Failed to open database", "error", err)
	} else {
		a.db = db
	}

	// Initialize and start the global keyboard shortcut manager.
	notifyFn := func(title, content string) {
		type notification struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		runtime.EventsEmit(a.ctx, "notification", notification{Title: title, Content: content})
	}
	if a.ollamaClient != nil {
		a.shortcutMgr = shortcuts.NewManager(a.cfg, a.ollamaClient, notifyFn)
		go a.shortcutMgr.Start()
	}
}

// shutdown is called by the Wails runtime when the application is about to
// exit.
func (a *App) shutdown(ctx context.Context) {
	if a.shortcutMgr != nil {
		a.shortcutMgr.Stop()
	}

	if a.cfg.StopOllamaOnShutDown {
		slog.Info("Stopping Ollama")
		cmd := exec.Command("pkill", "ollama")
		if err := cmd.Run(); err != nil {
			slog.Warn("Failed to stop Ollama", "error", err)
		}
	}

	if err := a.cfg.Save(); err != nil {
		slog.Error("Failed to save config", "error", err)
	}
}

// ---------------------------------------------------------------------------
// Settings
// ---------------------------------------------------------------------------

// AppSettings mirrors AppConfig for JSON transport to the frontend.
type AppSettings struct {
	ReplaceHighlightedText bool     `json:"replaceHighlightedText"`
	SpeakAIResponse        bool     `json:"speakAIResponse"`
	ShowPopUp              bool     `json:"showPopUp"`
	ShowStartWindow        bool     `json:"showStartWindow"`
	CurrentPrompt          int      `json:"currentPrompt"`
	CurrentChatModel       string   `json:"currentChatModel"`
	CurrentModel           string   `json:"currentModel"`
	CurrentFromLang        string   `json:"currentFromLang"`
	CurrentToLang          string   `json:"currentToLang"`
	StopOllamaOnShutDown   bool     `json:"stopOllamaOnShutDown"`
	UseRemoteOllama        bool     `json:"useRemoteOllama"`
	OllamaURL              string   `json:"ollamaURL"`
	AskAIShortcut          []string `json:"askAIShortcut"`
	CtrlReviseShortcut     []string `json:"ctrlReviseShortcut"`
	TranslateShortcut      []string `json:"translateShortcut"`
}

// GetSettings returns the current application settings.
func (a *App) GetSettings() AppSettings {
	return AppSettings{
		ReplaceHighlightedText: a.cfg.ReplaceHighlightedText,
		SpeakAIResponse:        a.cfg.SpeakAIResponse,
		ShowPopUp:              a.cfg.ShowPopUp,
		ShowStartWindow:        a.cfg.ShowStartWindow,
		CurrentPrompt:          a.cfg.CurrentPrompt,
		CurrentChatModel:       a.cfg.CurrentChatModel,
		CurrentModel:           a.cfg.CurrentModel,
		CurrentFromLang:        a.cfg.CurrentFromLang,
		CurrentToLang:          a.cfg.CurrentToLang,
		StopOllamaOnShutDown:   a.cfg.StopOllamaOnShutDown,
		UseRemoteOllama:        a.cfg.UseRemoteOllama,
		OllamaURL:              a.cfg.OllamaURL,
		AskAIShortcut:          a.cfg.AskAIShortcut,
		CtrlReviseShortcut:     a.cfg.CtrlReviseShortcut,
		TranslateShortcut:      a.cfg.TranslateShortcut,
	}
}

// SaveSettings persists updated settings from the frontend.
func (a *App) SaveSettings(settings AppSettings) error {
	a.cfg.ReplaceHighlightedText = settings.ReplaceHighlightedText
	a.cfg.SpeakAIResponse = settings.SpeakAIResponse
	a.cfg.ShowPopUp = settings.ShowPopUp
	a.cfg.ShowStartWindow = settings.ShowStartWindow
	a.cfg.CurrentPrompt = settings.CurrentPrompt
	a.cfg.CurrentChatModel = settings.CurrentChatModel
	a.cfg.CurrentModel = settings.CurrentModel
	a.cfg.CurrentFromLang = settings.CurrentFromLang
	a.cfg.CurrentToLang = settings.CurrentToLang
	a.cfg.StopOllamaOnShutDown = settings.StopOllamaOnShutDown
	a.cfg.UseRemoteOllama = settings.UseRemoteOllama
	a.cfg.OllamaURL = settings.OllamaURL
	a.cfg.AskAIShortcut = settings.AskAIShortcut
	a.cfg.CtrlReviseShortcut = settings.CtrlReviseShortcut
	a.cfg.TranslateShortcut = settings.TranslateShortcut
	return a.cfg.Save()
}

// ---------------------------------------------------------------------------
// Models
// ---------------------------------------------------------------------------

// recommendedModels lists models suggested to users who haven't installed any.
var recommendedModels = []string{
	"llama3.2:latest",
	"qwen2.5:3b",
	"gemma3:latest",
}

// ModelInfo describes a single Ollama model.
type ModelInfo struct {
	Name          string `json:"name"`
	IsInstalled   bool   `json:"isInstalled"`
	IsRecommended bool   `json:"isRecommended"`
}

// ollamaBaseURL returns the URL of the Ollama server based on current config.
func (a *App) ollamaBaseURL() string {
	if a.cfg.UseRemoteOllama && a.cfg.OllamaURL != "" {
		return a.cfg.OllamaURL
	}
	return "http://localhost:11434"
}

// GetModels queries Ollama for installed models.
func (a *App) GetModels() []ModelInfo {
	var result []ModelInfo

	if a.ollamaClient != nil {
		models, err := ollama.ListModels(context.Background(), a.ollamaBaseURL())
		if err == nil {
			for _, m := range models {
				isRec := false
				for _, rec := range recommendedModels {
					if rec == m.Name {
						isRec = true
						break
					}
				}
				result = append(result, ModelInfo{
					Name:          m.Name,
					IsInstalled:   true,
					IsRecommended: isRec,
				})
			}
		} else {
			slog.Error("Failed to list models", "error", err)
		}
	}

	return result
}

// ---------------------------------------------------------------------------
// Prompts
// ---------------------------------------------------------------------------

// PromptInfo describes a single prompt option.
type PromptInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetPrompts returns all available AI prompts.
func (a *App) GetPrompts() []PromptInfo {
	last := int(ollama.MakeItAListRedo)
	prompts := make([]PromptInfo, 0, last+1)
	for i := 0; i <= last; i++ {
		prompts = append(prompts, PromptInfo{
			ID:   i,
			Name: ollama.PromptMsg(i).String(),
		})
	}
	return prompts
}

// ---------------------------------------------------------------------------
// AI Ask
// ---------------------------------------------------------------------------

// AskQuestion sends a freeform question to Ollama and returns the response
// text.
func (a *App) AskQuestion(question string) string {
	if a.ollamaClient == nil {
		return "Not connected to Ollama"
	}
	resp, err := ollama.AskAI(a.cfg.CurrentModel, a.ollamaClient, question)
	if err != nil {
		slog.Error("AskQuestion failed", "error", err)
		return "Error: " + err.Error()
	}
	return resp.Response
}

// ---------------------------------------------------------------------------
// AI Revise
// ---------------------------------------------------------------------------

// ReviseText rewrites the given text using the specified prompt.
func (a *App) ReviseText(text string, promptID int) string {
	if a.ollamaClient == nil {
		return "Not connected to Ollama"
	}
	prompt := ollama.PromptMsg(promptID)
	resp, err := ollama.AskAIWithPromptMsg(a.cfg.CurrentModel, a.ollamaClient, prompt, text)
	if err != nil {
		slog.Error("ReviseText failed", "error", err)
		return "Error: " + err.Error()
	}
	return resp.Response
}

// ---------------------------------------------------------------------------
// Translate
// ---------------------------------------------------------------------------

// TranslateText translates the given text between the specified languages.
func (a *App) TranslateText(text, fromLang, toLang string) string {
	if a.ollamaClient == nil {
		return "Not connected to Ollama"
	}
	resp, err := ollama.AskAIToTranslate(a.cfg.CurrentModel, a.ollamaClient, text, ollama.Language(fromLang), ollama.Language(toLang))
	if err != nil {
		slog.Error("TranslateText failed", "error", err)
		return "Error: " + err.Error()
	}
	return resp.Response
}

// GetLanguages returns the list of supported translation languages.
func (a *App) GetLanguages() []string {
	return []string{
		string(ollama.Arabic),
		string(ollama.Chinese),
		string(ollama.English),
		string(ollama.French),
		string(ollama.German),
		string(ollama.Italian),
		string(ollama.Japanese),
		string(ollama.Portuguese),
		string(ollama.Russian),
		string(ollama.Spanish),
		string(ollama.Turkish),
	}
}

// ---------------------------------------------------------------------------
// Chat
// ---------------------------------------------------------------------------

// ChatInfo is the JSON-serialisable representation of a stored chat.
type ChatInfo struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Model     int      `json:"model"`
	Questions []string `json:"questions"`
	Responses []string `json:"responses"`
}

// ChatResponse is returned by SendChatMessage.
type ChatResponse struct {
	Response string `json:"response"`
	ChatID   int64  `json:"chatId"`
	Context  []int  `json:"context"`
}

// GetChats returns all stored chats.
func (a *App) GetChats() []ChatInfo {
	if a.db == nil {
		return nil
	}
	chats, err := a.db.GetAllChats("default")
	if err != nil {
		slog.Error("GetChats failed", "error", err)
		return nil
	}
	result := make([]ChatInfo, 0, len(chats))
	for _, c := range chats {
		var id int64
		if c.ID != nil {
			id = *c.ID
		}
		result = append(result, ChatInfo{
			ID:        id,
			Title:     c.Title,
			Model:     c.Model,
			Questions: c.Questions,
			Responses: c.Responses,
		})
	}
	return result
}

// SendChatMessage sends a chat message to Ollama and stores the result.
// When isNew is true and chatID is 0 a new chat is created; otherwise the
// existing chat is updated.
func (a *App) SendChatMessage(chatID int64, isNew bool, message string, msgContext []int) ChatResponse {
	if a.ollamaClient == nil {
		return ChatResponse{Response: "Not connected to Ollama"}
	}

	resp, err := ollama.AskAiWithStringAndContext(a.cfg.CurrentChatModel, a.ollamaClient, msgContext, message)
	if err != nil {
		slog.Error("SendChatMessage failed", "error", err)
		return ChatResponse{Response: "Error: " + err.Error()}
	}

	if a.db != nil {
		if isNew || chatID == 0 {
			title := message
			if len([]rune(title)) > 14 {
				title = string([]rune(title)[:14])
			}
			entry := &chatModel.Chat{
				Model:     0,
				Owner:     "default",
				Title:     title,
				Questions: []string{message},
				Responses: []string{resp.Response},
				Context:   resp.Context,
			}
			if saveErr := a.db.SaveChat(entry); saveErr != nil {
				slog.Error("Failed to save new chat", "error", saveErr)
			} else if entry.ID != nil {
				chatID = *entry.ID
			}
		} else {
			entry := &chatModel.Chat{
				ID:        &chatID,
				Model:     0,
				Owner:     "default",
				Questions: []string{message},
				Responses: []string{resp.Response},
				Context:   resp.Context,
			}
			if updateErr := a.db.UpdateChat(entry); updateErr != nil {
				slog.Error("Failed to update chat", "error", updateErr)
			}
		}
	}

	return ChatResponse{
		Response: resp.Response,
		ChatID:   chatID,
		Context:  resp.Context,
	}
}

// DeleteChat removes a stored chat by ID.
func (a *App) DeleteChat(id int64) error {
	if a.db == nil {
		return nil
	}
	return a.db.DeleteChat(id)
}

// ---------------------------------------------------------------------------
// Meal Planner
// ---------------------------------------------------------------------------

// MealInfo is the JSON-serialisable meal information sent from the frontend.
type MealInfo struct {
	Consumers string   `json:"consumers"`
	Meal      string   `json:"meal"`
	Theme     string   `json:"theme"`
	Cookware  []string `json:"cookware"`
	Dairy     []string `json:"dairy"`
	Veggies   []string `json:"veggies"`
	Protein   []string `json:"protein"`
	Freezer   []string `json:"freezer"`
	Herbs     []string `json:"herbs"`
	Pantry    []string `json:"pantry"`
	Allergies []string `json:"allergies"`
}

// toMealPkg converts the app-level MealInfo to the meal package type.
func toMealPkg(info MealInfo) meal.MealInfo {
	return meal.MealInfo{
		Consumers: info.Consumers,
		Meal:      info.Meal,
		Theme:     info.Theme,
		Cookware:  info.Cookware,
		Dairy:     info.Dairy,
		Veggies:   info.Veggies,
		Protein:   info.Protein,
		Freezer:   info.Freezer,
		Herbs:     info.Herbs,
		Pantry:    info.Pantry,
		Allergies: info.Allergies,
	}
}

// GenerateMeal asks the AI to suggest a single meal recipe.
func (a *App) GenerateMeal(info MealInfo) string {
	if a.ollamaClient == nil {
		return "Not connected to Ollama"
	}
	prompt := meal.CreateMealPrompt(toMealPkg(info)) + ". format: markdown"
	resp, err := ollama.AskAI(a.cfg.CurrentModel, a.ollamaClient, prompt)
	if err != nil {
		slog.Error("GenerateMeal failed", "error", err)
		return "Error: " + err.Error()
	}
	return resp.Response
}

// GenerateMealPrep asks the AI for a weekly meal prep plan.
func (a *App) GenerateMealPrep(info MealInfo) string {
	if a.ollamaClient == nil {
		return "Not connected to Ollama"
	}
	prompt := meal.CreateMealPrepPrompt(toMealPkg(info)) + ". format: markdown"
	resp, err := ollama.AskAI(a.cfg.CurrentModel, a.ollamaClient, prompt)
	if err != nil {
		slog.Error("GenerateMealPrep failed", "error", err)
		return "Error: " + err.Error()
	}
	return resp.Response
}

// GenerateGroceryList asks the AI for a weekly grocery list.
func (a *App) GenerateGroceryList(info MealInfo) string {
	if a.ollamaClient == nil {
		return "Not connected to Ollama"
	}
	prompt := meal.CreateGroceryListPrompt(toMealPkg(info)) + ". format: markdown"
	resp, err := ollama.AskAI(a.cfg.CurrentModel, a.ollamaClient, prompt)
	if err != nil {
		slog.Error("GenerateGroceryList failed", "error", err)
		return "Error: " + err.Error()
	}
	return resp.Response
}

// GenerateBudgetGroceryList asks the AI for a budget-conscious grocery list.
func (a *App) GenerateBudgetGroceryList(info MealInfo) string {
	if a.ollamaClient == nil {
		return "Not connected to Ollama"
	}
	prompt := meal.CreateBudgetFriendlyGroceryListPrompt(toMealPkg(info)) + ". format: markdown"
	resp, err := ollama.AskAI(a.cfg.CurrentModel, a.ollamaClient, prompt)
	if err != nil {
		slog.Error("GenerateBudgetGroceryList failed", "error", err)
		return "Error: " + err.Error()
	}
	return resp.Response
}

// ---------------------------------------------------------------------------
// Shortcuts
// ---------------------------------------------------------------------------

// ShortcutsConfig holds the keyboard shortcut configuration for the frontend.
type ShortcutsConfig struct {
	AskAI      []string `json:"askAI"`
	CtrlRevise []string `json:"ctrlRevise"`
	Translate  []string `json:"translate"`
}

// GetShortcuts returns the current keyboard shortcut configuration.
func (a *App) GetShortcuts() ShortcutsConfig {
	return ShortcutsConfig{
		AskAI:      a.cfg.AskAIShortcut,
		CtrlRevise: a.cfg.CtrlReviseShortcut,
		Translate:  a.cfg.TranslateShortcut,
	}
}

// SaveShortcuts persists new keyboard shortcut bindings.
func (a *App) SaveShortcuts(s ShortcutsConfig) error {
	a.cfg.AskAIShortcut = s.AskAI
	a.cfg.CtrlReviseShortcut = s.CtrlRevise
	a.cfg.TranslateShortcut = s.Translate
	return a.cfg.Save()
}

// ---------------------------------------------------------------------------
// Ollama management
// ---------------------------------------------------------------------------

// CheckOllama returns true when Ollama is reachable.
func (a *App) CheckOllama() bool {
	var connectURL *string
	if a.cfg.UseRemoteOllama && a.cfg.OllamaURL != "" {
		connectURL = &a.cfg.OllamaURL
	}
	client, err := ollama.CheckOllamaConnection(connectURL)
	if err != nil {
		return false
	}
	a.ollamaClient = client
	return true
}

// TestOllamaURL tests whether a connection to the given URL can be established.
func (a *App) TestOllamaURL(url string) bool {
	_, err := ollama.CheckOllamaConnection(&url)
	return err == nil
}
