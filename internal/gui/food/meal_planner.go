package food

import (
	"crypto/sha256"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/loading"
	"log/slog"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/bahelit/ctrl_plus_revise/internal/config"
	"github.com/bahelit/ctrl_plus_revise/internal/gui/shortcuts"
	"github.com/bahelit/ctrl_plus_revise/internal/ollama"
	ollamaApi "github.com/ollama/ollama/api"
)

type MealInfo struct {
	Consumers string
	Meal      string
	Theme     string
	Cookware  []string
	Dairy     []string
	Freezer   []string
	Herbs     []string
	Pantry    []string
	Veggies   []string
	Protein   []string
	Allergies []string
}

func MealPlanner(guiApp fyne.App, ollamaClient *ollamaApi.Client) {
	var (
		screenHeight float32 = 575.0
		screenWidth  float32 = 655.0
		tabs         container.AppTabs
	)
	mealPlanner := guiApp.NewWindow("Ctrl+Revise Meal Planner")
	mealPlanner.Resize(fyne.NewSize(screenWidth, screenHeight))

	var mealInfo MealInfo

	topText := widget.NewLabel("Putting an end to the age old question of \"what's for dinner?\"")
	topText.Alignment = fyne.TextAlignCenter

	meal := widget.NewRadioGroup(getMeals(), func(s string) {
		slog.Debug("selected meal", "meal", s)
		guiApp.Preferences().SetString(config.MealKey, s)
	})
	meal.SetSelected(guiApp.Preferences().String(config.MealKey))
	meal.OnChanged = func(s string) {
		guiApp.Preferences().SetString(config.MealKey, s)
		mealInfo.Meal = s
	}

	consumersValidated := newNumEntry()
	consumersValidated.SetPlaceHolder("How many people are eating?")
	consumersValidated.OnSubmitted = func(s string) {
		guiApp.Preferences().SetString(config.ConsumersKey, s)
	}
	consumers := guiApp.Preferences().String(config.ConsumersKey)
	if consumers != "" {
		consumersValidated.SetText(consumers)
	}

	mealAndCount := container.NewVBox(meal, consumersValidated)
	mealCard := widget.NewCard("Meal", "Breakfast or dinner, or how about breakfast for dinner?", mealAndCount)

	flavor := widget.NewRadioGroup(getThemes(), func(s string) {
		slog.Debug("selected cuisine", "flavor", s)
		guiApp.Preferences().SetString(config.FlavorKey, s)
	})
	flavor.SetSelected(guiApp.Preferences().String(config.FlavorKey))
	flavor.OnChanged = func(s string) {
		guiApp.Preferences().SetString(config.FlavorKey, s)
		mealInfo.Theme = s
	}
	flavorCard := widget.NewCard("Flavor Profile", "What kind of style for our food?", flavor)

	cookWare := widget.NewCheckGroup(getCookWares(), func(strings []string) {
		slog.Debug("selected cookWare", "cookWare", strings)
		guiApp.Preferences().SetStringList(config.CookwareKey, strings)
	})
	cookWare.SetSelected(guiApp.Preferences().StringList(config.CookwareKey))
	cookwareCard := widget.NewCard("Cookware", "What do we have cooking with?", cookWare)

	dairyCheck := widget.NewCheckGroup(getDairy(), func(strings []string) {
		slog.Debug("dairy stuff", "dairy", strings)
		guiApp.Preferences().SetStringList(config.DairyKey, strings)
	})
	dairyCheck.SetSelected(guiApp.Preferences().StringList(config.DairyKey))
	dairyCard := widget.NewCard("Dairy", "Cool the burn with some dairy", dairyCheck)

	freezerCheck := widget.NewCheckGroup(getFreezerStuffs(), func(strings []string) {
		slog.Debug("stuff in freezer", "stuff in freezer", strings)
		guiApp.Preferences().SetStringList(config.FreezerKey, strings)
	})
	freezerCheck.SetSelected(guiApp.Preferences().StringList(config.FreezerKey))
	freezerCard := widget.NewCard("freezer", "What's in the freezer?", freezerCheck)

	herbsCheck := widget.NewCheckGroup(getHerbsSpices(), func(strings []string) {
		slog.Debug("herbs & spices", "herbs", strings)
		guiApp.Preferences().SetStringList(config.HerbsKey, strings)
	})
	herbsCheck.SetSelected(guiApp.Preferences().StringList(config.HerbsKey))
	herbsCard := widget.NewCard("Pantry", "Enhance the aroma and flavor.", herbsCheck)

	pantryCheck := widget.NewCheckGroup(getPantryStuffs(), func(strings []string) {
		slog.Debug("stuff in pantry", "pantry", strings)
		guiApp.Preferences().SetStringList(config.PantryKey, strings)
	})
	pantryCheck.SetSelected(guiApp.Preferences().StringList(config.PantryKey))
	pantryCard := widget.NewCard("Pantry", "What's in the pantry?", pantryCheck)

	veggieCheck := widget.NewCheckGroup(getVegetables(), func(strings []string) {
		slog.Debug("dont' forget your vegetables", "pantry", strings)
		guiApp.Preferences().SetStringList(config.VeggiesKey, strings)
	})
	veggieCheck.SetSelected(guiApp.Preferences().StringList(config.VeggiesKey))
	veggieCard := widget.NewCard("Vegetables", "Don't forget you're vegetables.", veggieCheck)

	proteinCheck := widget.NewCheckGroup(getProteins(), func(strings []string) {
		slog.Debug("Meat or meat alternatives", "pantry", strings)
		guiApp.Preferences().SetStringList(config.ProteinKey, strings)
	})
	proteinCheck.SetSelected(guiApp.Preferences().StringList(config.ProteinKey))
	proteinCard := widget.NewCard("Protein", "Meat or meat substitutes", proteinCheck)

	allergyCheck := widget.NewCheckGroup(getAllergies(), func(strings []string) {
		slog.Debug("Food allergies", "allergy", strings)
	})
	allergyCard := widget.NewCard("Allergy", "Got any food allergies?", allergyCheck)

	mealTab := container.NewTabItem("Meal", container.NewVScroll(mealCard))
	themeTab := container.NewTabItem("Theme", container.NewVScroll(flavorCard))
	cookwareTab := container.NewTabItem("Cookwares", container.NewVScroll(cookwareCard))
	dairyTab := container.NewTabItem("Dairy", container.NewVScroll(dairyCard))
	freezerTab := container.NewTabItem("Freezer", container.NewVScroll(freezerCard))
	herbsTab := container.NewTabItem("Herbs & Spices", container.NewVScroll(herbsCard))
	pantryTab := container.NewTabItem("Pantry", container.NewVScroll(pantryCard))
	veggiesTab := container.NewTabItem("Vegetables", container.NewVScroll(veggieCard))
	proteinTab := container.NewTabItem("Protein", container.NewVScroll(proteinCard))
	allergyTab := container.NewTabItem("Allergy", container.NewVScroll(allergyCard))
	verticalTabs := container.NewAppTabs(mealTab, themeTab, cookwareTab, herbsTab, veggiesTab, dairyTab, pantryTab, proteinTab, freezerTab, allergyTab)
	verticalTabs.SetTabLocation(container.TabLocationLeading)

	mealInfo.Consumers = consumersValidated.Text
	mealInfo.Cookware = shuffleStringArray(cookWare.Selected)
	mealInfo.Dairy = shuffleStringArray(dairyCheck.Selected)
	mealInfo.Freezer = shuffleStringArray(freezerCheck.Selected)
	mealInfo.Herbs = shuffleStringArray(herbsCheck.Selected)
	mealInfo.Pantry = shuffleStringArray(pantryCheck.Selected)
	mealInfo.Veggies = shuffleStringArray(veggieCheck.Selected)
	mealInfo.Protein = shuffleStringArray(proteinCheck.Selected)
	mealInfo.Allergies = shuffleStringArray(allergyCheck.Selected)

	suggest := widget.NewButton("Suggest a Single Meal", func() {
		recipe := createMealPrompt(mealInfo)
		slog.Info("Recipe for single meal", "PROMPT", recipe)
		loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
			"Preparing recipe")
		loadingScreen.Show()

		recipe = addMarkdownFormattingToRecipe(recipe)

		generated, err := ollama.AskAI(guiApp, ollamaClient, recipe)
		if err != nil {
			slog.Error("Failed to ask AI", "error", err)
			loadingScreen.Hide()
			return
		}
		loadingScreen.Hide()
		recipePopUp(guiApp, &tabs, ollamaClient, recipe, &generated)
	})
	suggest.Importance = widget.HighImportance
	suggestPrep := widget.NewButton("Prep Multiple Meals", func() {
		recipe := createMealPrepPrompt(mealInfo)
		slog.Info("Meal prep plan", "PROMPT", recipe)
		loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
			"Creating meal prep")
		loadingScreen.Show()

		recipe = addMarkdownFormattingToRecipe(recipe)

		generated, err := ollama.AskAI(guiApp, ollamaClient, recipe)
		if err != nil {
			slog.Error("Failed to ask AI", "error", err)
			loadingScreen.Hide()
			return
		}
		loadingScreen.Hide()
		recipePopUp(guiApp, &tabs, ollamaClient, recipe, &generated)
	})
	suggestPrep.Importance = widget.HighImportance
	groceryList := widget.NewButton("Create a Grocery List", func() {
		recipe := createGroceryListPrompt(mealInfo)
		slog.Info("Grocery List", "PROMPT", recipe)
		loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
			"Creating grocery list")
		loadingScreen.Show()

		recipe = addMarkdownFormattingToRecipe(recipe)

		generated, err := ollama.AskAI(guiApp, ollamaClient, recipe)
		if err != nil {
			slog.Error("Failed to ask AI", "error", err)
			loadingScreen.Hide()
			return
		}
		loadingScreen.Hide()
		recipePopUp(guiApp, &tabs, ollamaClient, recipe, &generated)
	})
	groceryList.Importance = widget.SuccessImportance
	budgetFriendlyGroceryList := widget.NewButton("Create a Budget Friendly Grocery List", func() {
		recipe := createBudgetFriendlyGroceryListPrompt(mealInfo)
		slog.Info("Budget Friendly Shopping plan", "PROMPT", recipe)
		loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
			"Creating budget friendly grocery list")
		loadingScreen.Show()

		recipe = addMarkdownFormattingToRecipe(recipe)

		generated, err := ollama.AskAI(guiApp, ollamaClient, recipe)
		if err != nil {
			slog.Error("Failed to ask AI", "error", err)
			loadingScreen.Hide()
			return
		}
		loadingScreen.Hide()
		recipePopUp(guiApp, &tabs, ollamaClient, recipe, &generated)
	})
	budgetFriendlyGroceryList.Importance = widget.SuccessImportance

	action := container.NewHBox(suggest, groceryList, suggestPrep, budgetFriendlyGroceryList)
	action.Layout = layout.NewAdaptiveGridLayout(2)

	boarderLayout := container.NewBorder(nil, action, nil, nil, verticalTabs)
	mealPlannerTab := container.NewTabItem("Meal Planner", boarderLayout)
	tabs.Append(mealPlannerTab)
	tabContainer := container.NewBorder(topText, nil, nil, nil, &tabs)

	mealPlanner.SetContent(tabContainer)
	mealPlanner.Show()
}

type numEntry struct {
	widget.Entry
}

func (n *numEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

func newNumEntry() *numEntry {
	e := &numEntry{}
	e.ExtendBaseWidget(e)
	e.Validator = validation.NewRegexp(`^\d{1,3}$`, "Must contain a number")
	return e
}

// shuffleStringArray Shuffle the order so the AI isn't recommending the items at the top as often.
func shuffleStringArray(str []string) []string {
	rand.New(rand.NewSource(time.Now().UnixNano())).Shuffle(len(str), func(i, j int) {})
	sort.Slice(str,
		func(i, j int) bool {
			return rand.Int63()&1 == 0 && i < j
		})
	return str
}

func addMarkdownFormattingToRecipe(recipe string) string {
	return recipe + ". " + "format: markdown"
}

func recipePopUp(guiApp fyne.App, tabs *container.AppTabs, ollamaClient *ollamaApi.Client, recipe string, response *ollamaApi.GenerateResponse) {
	generatedText1 := widget.NewRichTextFromMarkdown(response.Response)
	generatedText1.Wrapping = fyne.TextWrapWord

	vbox := container.NewVScroll(generatedText1)
	tabCount := len(tabs.Items)

	buttons := container.NewPadded(
		widget.NewButton("Something Different", func() {
			loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
				"Something Different...")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithStringAndContext(guiApp, ollamaClient, response.Context,
				"That doesn't sound good, how about something else please.")
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			shortcuts.LastClipboardContent = sha256.Sum256([]byte(response.Response))
			loadingScreen.Hide()
			recipePopUp(guiApp, tabs, ollamaClient, recipe, &reGenerated)
		}),
		widget.NewButton("Something Simple and Quick", func() {
			loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
				"Something Simple and Quick...")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithStringAndContext(guiApp, ollamaClient, response.Context,
				"That doesn't sound good, how about something that's simple and quick to put together please.")
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			shortcuts.LastClipboardContent = sha256.Sum256([]byte(response.Response))
			loadingScreen.Hide()
			recipePopUp(guiApp, tabs, ollamaClient, recipe, &reGenerated)
		}),
		widget.NewButton("Something Healthy", func() {
			loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
				"Something Healthy...")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithStringAndContext(guiApp, ollamaClient, response.Context,
				"That doesn't sound good, how about something really healthy but still tasty please.")
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			shortcuts.LastClipboardContent = sha256.Sum256([]byte(response.Response))
			loadingScreen.Hide()
			recipePopUp(guiApp, tabs, ollamaClient, recipe, &reGenerated)
		}),
		widget.NewButton("Let's Not Cook", func() {
			loadingScreen := loading.LoadingScreenWithMessageAddModel(guiApp, loading.ThinkingMsg,
				"Let's Not Cook")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithStringAndContext(guiApp, ollamaClient, response.Context,
				"I don't feel like cooking or heating anything up, how about something cold I could put together please.")
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			shortcuts.LastClipboardContent = sha256.Sum256([]byte(response.Response))
			loadingScreen.Hide()
			recipePopUp(guiApp, tabs, ollamaClient, recipe, &reGenerated)
		}),
		widget.NewButtonWithIcon("Close Recipe", theme.ContentClearIcon(), func() {
			tabs.Remove(tabs.Selected())
		}),
		widget.NewButtonWithIcon("Copy the Recipe to Clipboard", theme.ContentCopyIcon(), func() {
			// TODO: Pass in window?
			//w.Clipboard().SetContent(response.Response)
			//w.Close()
		}),
	)

	buttons.Layout = layout.NewAdaptiveGridLayout(3)
	center := container.NewVBox(buttons, layout.NewSpacer(), footer())
	grid := container.New(layout.NewAdaptiveGridLayout(1), vbox)
	tabSpace := container.NewBorder(
		nil,
		center,
		nil,
		nil,
		container.NewVScroll(grid),
	)
	tabs.Append(container.NewTabItem("Recipe #"+strconv.Itoa(tabCount), tabSpace))
	tabs.SelectIndex(tabCount)
}

func footer() *fyne.Container {
	footer := container.NewHBox(
		layout.NewSpacer(),
		widget.NewHyperlink("Ctrl+Revise", parseURL("https://ctrlplusrevise.com")),
		widget.NewLabel("-"),
		widget.NewHyperlink("Documentation", parseURL("https://ctrlplusrevise.com/docs/tutorials/")),
		widget.NewLabel("-"),
		widget.NewHyperlink("Sponsor", parseURL("https://www.patreon.com/SalmonsStudios")),
		layout.NewSpacer(),
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
