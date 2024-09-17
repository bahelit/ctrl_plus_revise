package main

import (
	"log/slog"
	"net/url"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	layoutv1 "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	ollamaApi "github.com/ollama/ollama/api"

	"github.com/bahelit/ctrl_plus_revise/internal/gui/bindings"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/chat"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/food"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/menu"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/question"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/settings"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/shortcuts"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/translator"
)

const (
	AppTitle      = "Ctrl+Revise AI Text Generator"
	GreetingsText = "Welcome to Ctrl+Revise!"
	TrayMenuTitle = "Ctrl+Revise"
)

// SetupSysTray initializes the system tray for the application
func SetupSysTray(guiApp fyne.App, ollamaClient *ollamaApi.Client) fyne.Window {
	if err := bindings.SetBindingVariables(guiApp); err != nil {
		slog.Error("Failed to set binding variables", "error", err)
		os.Exit(1)
	}

	sysTray := guiApp.NewWindow(AppTitle)
	sysTray.SetTitle(AppTitle)

	sysTray.SetMainMenu(menu.MakeMenu(guiApp, ollamaClient, sysTray))

	setupTrayMenu(guiApp, ollamaClient, sysTray)
	setupTrayWindowContent(guiApp, ollamaClient, sysTray)

	sysTray.SetCloseIntercept(func() {
		sysTray.Hide()
	})
	return sysTray
}

// setupTrayMenu sets up the system tray menu
func setupTrayMenu(guiApp fyne.App, ollamaClient *ollamaApi.Client, sysTray fyne.Window) {
	if desk, ok := guiApp.(desktop.App); ok {
		desk.SetSystemTrayMenu(fyne.NewMenu(TrayMenuTitle,
			fyne.NewMenuItem("Ask a Question", func() { question.AskQuestionWindow(guiApp, ollamaClient) }),
			fyne.NewMenuItem("Meal Planner", func() { food.MealPlanner(guiApp, ollamaClient) }),
			fyne.NewMenuItem("Translate Window", func() { translator.TranslateText(guiApp, ollamaClient) }),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Home Screen", func() { sysTray.Show() }),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Settings", func() { settings.ShowSettings(guiApp, ollamaClient) }),
			fyne.NewMenuItem("Keyboard Shortcuts", func() { shortcuts.ShowShortcuts(guiApp) }),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("About", func() { settings.ShowAbout(guiApp) }),
		))
	}
}

// setupTrayWindowContent sets up the content of the system tray window
func setupTrayWindowContent(guiApp fyne.App, ollamaClient *ollamaApi.Client, sysTray fyne.Window) {
	welcomeText := mainWindowText()

	askQuestionsButton := widget.NewButton("Ask a Question", func() {
		question.AskQuestionWindow(guiApp, ollamaClient)
	})
	chatButton := widget.NewButton("Chat with AI - Beta", func() {
		chat.ConversationManager(guiApp, ollamaClient)
	})
	recipeButton := widget.NewButton("Meal Planner", func() {
		food.MealPlanner(guiApp, ollamaClient)
	})
	translatorButton := widget.NewButton("Translate Text", func() {
		translator.TranslateText(guiApp, ollamaClient)
	})

	buttons := container.NewVBox(
		askQuestionsButton,
		chatButton,
		recipeButton,
		translatorButton,
	)
	footer := footer()
	sysTray.SetContent(container.NewBorder(welcomeText, footer, nil, nil, buttons))
}

func footer() *fyne.Container {
	footer := container.NewHBox(
		layoutv1.NewSpacer(),
		widget.NewHyperlink("Ctrl+Revise", parseURL("https://ctrlplusrevise.com")),
		widget.NewLabel("-"),
		widget.NewHyperlink("Documentation", parseURL("https://ctrlplusrevise.com/docs/tutorials/")),
		widget.NewLabel("-"),
		widget.NewHyperlink("Sponsor", parseURL("https://www.patreon.com/SalmonsStudios")),
		layoutv1.NewSpacer(),
	)
	return footer
}

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

// TODO: Embed icon in binary
func loadIcon(guiApp fyne.App) {
	var (
		icon         fyne.Resource
		errLocation1 error
		errLocation2 error
	)
	icon, errLocation1 = fyne.LoadResourceFromPath("/app/share/icons/hicolor/256x256/apps/com.bahelit.ctrl_plus_revise.png")
	if errLocation1 != nil {
		icon, errLocation2 = fyne.LoadResourceFromPath("images/icon.png")
		if errLocation2 != nil {
			slog.Warn("Failed to load icon", "error", errLocation1)
			slog.Warn("Failed to load icon", "error", errLocation2)
		}
	}
	guiApp.SetIcon(icon)
}

func mainWindowText() *fyne.Container {
	welcomeText := widget.NewLabel(GreetingsText)
	welcomeText.Alignment = fyne.TextAlignCenter
	welcomeText.TextStyle = fyne.TextStyle{Bold: true}

	ctrlReviseKeys := shortcuts.GetCtrlReviseKeys()
	var ctrlReviseString string
	keyLength := len(ctrlReviseKeys)
	for key, value := range ctrlReviseKeys {
		ctrlReviseString += strings.ToUpper(value)
		if key != keyLength-1 {
			ctrlReviseString += " + "
		}
	}
	shortcutText := widget.NewLabel("Pressing \"" + ctrlReviseString + "\" will send the highlighted text to an AI\nthe response is put into the clipboard")
	shortcutText.Alignment = fyne.TextAlignCenter
	shortcutText.TextStyle = fyne.TextStyle{Bold: true}
	closeMeText := widget.NewLabel("This window can be closed, the program will keep running in the taskbar")
	closeMeText.Alignment = fyne.TextAlignCenter
	return container.NewVBox(welcomeText, closeMeText, shortcutText)
}
