package main

import (
	"crypto/sha256"
	"log/slog"
	"math/rand"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/bahelit/ctrl_plus_revise/internal/gui"
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

func mealPlanner(guiApp fyne.App) {
	var (
		screenHeight float32 = 550.0
		screenWidth  float32 = 650.0
	)
	mealPlanner := guiApp.NewWindow("Ctrl+Revise Meal Planner")
	mealPlanner.Resize(fyne.NewSize(screenWidth, screenHeight))

	mealPlanner.SetMainMenu(makeMenu(guiApp, mealPlanner))

	topText := widget.NewLabel("Putting an end to the age old question of \"what's for dinner?\"")
	topText.Alignment = fyne.TextAlignCenter

	meal := widget.NewRadioGroup(getMeals(), func(s string) {
		slog.Debug("selected meal", "meal", s)
		guiApp.Preferences().SetString(MealKey, s)
	})
	meal.SetSelected(guiApp.Preferences().String(MealKey))

	consumersValidated := newNumEntry()
	consumersValidated.SetPlaceHolder("How many people are eating?")
	consumersValidated.OnSubmitted = func(s string) {
		guiApp.Preferences().SetString(ConsumersKey, s)
	}
	consumers := guiApp.Preferences().String(ConsumersKey)
	if consumers != "" {
		consumersValidated.SetText(consumers)
	}

	mealAndCount := container.NewVBox(meal, consumersValidated)
	mealCard := widget.NewCard("Meal", "Breakfast or dinner, or how about breakfast for dinner?", mealAndCount)

	flavor := widget.NewRadioGroup(getThemes(), func(s string) {
		slog.Debug("selected cuisine", "flavor", s)
		guiApp.Preferences().SetString(FlavorKey, s)
	})
	flavor.SetSelected(guiApp.Preferences().String(FlavorKey))
	flavorCard := widget.NewCard("Flavor Profile", "What kind of style for our food?", flavor)

	cookWare := widget.NewCheckGroup(getCookWares(), func(strings []string) {
		slog.Debug("selected cookWare", "cookWare", strings)
		guiApp.Preferences().SetStringList(CookwareKey, strings)
	})
	cookWare.SetSelected(guiApp.Preferences().StringList(CookwareKey))
	cookwareCard := widget.NewCard("Cookware", "What do we have cooking with?", cookWare)

	dairyCheck := widget.NewCheckGroup(getDairy(), func(strings []string) {
		slog.Debug("dairy stuff", "dairy", strings)
		guiApp.Preferences().SetStringList(DairyKey, strings)
	})
	dairyCheck.SetSelected(guiApp.Preferences().StringList(DairyKey))
	dairyCard := widget.NewCard("Dairy", "Cool the burn with some dairy", dairyCheck)

	freezerCheck := widget.NewCheckGroup(getFreezerStuffs(), func(strings []string) {
		slog.Debug("stuff in freezer", "stuff in freezer", strings)
		guiApp.Preferences().SetStringList(FreezerKey, strings)
	})
	freezerCheck.SetSelected(guiApp.Preferences().StringList(FreezerKey))
	freezerCard := widget.NewCard("freezer", "What's in the freezer?", freezerCheck)

	herbsCheck := widget.NewCheckGroup(getHerbsSpices(), func(strings []string) {
		slog.Debug("herbs & spices", "herbs", strings)
		guiApp.Preferences().SetStringList(HerbsKey, strings)
	})
	herbsCheck.SetSelected(guiApp.Preferences().StringList(HerbsKey))
	herbsCard := widget.NewCard("Pantry", "Enhance the aroma and flavor.", herbsCheck)

	pantryCheck := widget.NewCheckGroup(getPantryStuffs(), func(strings []string) {
		slog.Debug("stuff in pantry", "pantry", strings)
		guiApp.Preferences().SetStringList(PantryKey, strings)
	})
	pantryCheck.SetSelected(guiApp.Preferences().StringList(PantryKey))
	pantryCard := widget.NewCard("Pantry", "What's in the pantry?", pantryCheck)

	veggieCheck := widget.NewCheckGroup(getVegetables(), func(strings []string) {
		slog.Debug("dont' forget your vegetables", "pantry", strings)
		guiApp.Preferences().SetStringList(VeggiesKey, strings)
	})
	veggieCheck.SetSelected(guiApp.Preferences().StringList(VeggiesKey))
	veggieCard := widget.NewCard("Vegetables", "Don't forget you're vegetables.", veggieCheck)

	proteinCheck := widget.NewCheckGroup(getProteins(), func(strings []string) {
		slog.Debug("Meat or meat alternatives", "pantry", strings)
		guiApp.Preferences().SetStringList(ProteinKey, strings)
	})
	proteinCheck.SetSelected(guiApp.Preferences().StringList(ProteinKey))
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
	tabs := container.NewAppTabs(mealTab, themeTab, cookwareTab, herbsTab, veggiesTab, dairyTab, pantryTab, proteinTab, freezerTab, allergyTab)
	tabs.SetTabLocation(container.TabLocationLeading)

	mealInfo := MealInfo{
		Consumers: consumersValidated.Text,
		Meal:      meal.Selected,
		Theme:     flavor.Selected,
		Cookware:  shuffleStringArray(cookWare.Selected),
		Dairy:     shuffleStringArray(dairyCheck.Selected),
		Freezer:   shuffleStringArray(freezerCheck.Selected),
		Herbs:     shuffleStringArray(herbsCheck.Selected),
		Pantry:    shuffleStringArray(pantryCheck.Selected),
		Veggies:   shuffleStringArray(veggieCheck.Selected),
		Protein:   shuffleStringArray(proteinCheck.Selected),
		Allergies: shuffleStringArray(allergyCheck.Selected),
	}

	suggest := widget.NewButton("Suggest a Single Meal", func() {
		recipe := createMealPrompt(mealInfo)
		slog.Info("Recipe for single meal", "PROMPT", recipe)
		model := guiApp.Preferences().IntWithFallback(CurrentModelKey, int(ollama.Llama3Dot1))
		loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
			"Preparing recipe with model: "+ollama.ModelName(model).String())
		loadingScreen.Show()

		recipe = addMarkdownFomatingToRecipe(recipe)

		generated, err := ollama.AskAI(ollamaClient, ollama.ModelName(model), recipe)
		if err != nil {
			slog.Error("Failed to ask AI", "error", err)
			loadingScreen.Hide()
			return
		}
		loadingScreen.Hide()
		recipePopUp(guiApp, recipe, &generated)
	})
	suggest.Importance = widget.HighImportance
	suggestPrep := widget.NewButton("Prep Multiple Meals", func() {
		recipe := createMealPrepPrompt(mealInfo)
		slog.Info("Meal prep plan", "PROMPT", recipe)
		model := guiApp.Preferences().IntWithFallback(CurrentModelKey, int(ollama.Llama3Dot1))
		loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
			"Creating meal prep with model: "+ollama.ModelName(model).String())
		loadingScreen.Show()

		recipe = addMarkdownFomatingToRecipe(recipe)

		generated, err := ollama.AskAI(ollamaClient, ollama.ModelName(model), recipe)
		if err != nil {
			slog.Error("Failed to ask AI", "error", err)
			loadingScreen.Hide()
			return
		}
		loadingScreen.Hide()
		recipePopUp(guiApp, recipe, &generated)
	})
	suggestPrep.Importance = widget.HighImportance
	groceryList := widget.NewButton("Create a Grocery List", func() {
		recipe := createGroceryListPrompt(mealInfo)
		slog.Info("Grocery List", "PROMPT", recipe)
		model := guiApp.Preferences().IntWithFallback(CurrentModelKey, int(ollama.Llama3Dot1))
		loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
			"Creating grocery list with model: "+ollama.ModelName(model).String())
		loadingScreen.Show()

		recipe = addMarkdownFomatingToRecipe(recipe)

		generated, err := ollama.AskAI(ollamaClient, ollama.ModelName(model), recipe)
		if err != nil {
			slog.Error("Failed to ask AI", "error", err)
			loadingScreen.Hide()
			return
		}
		loadingScreen.Hide()
		recipePopUp(guiApp, recipe, &generated)
	})
	groceryList.Importance = widget.SuccessImportance
	budgetFriendlyGroceryList := widget.NewButton("Create a Budget Friendly Grocery List", func() {
		recipe := createBudgetFriendlyGroceryListPrompt(mealInfo)
		slog.Info("Budget Friendly Shopping plan", "PROMPT", recipe)
		model := guiApp.Preferences().IntWithFallback(CurrentModelKey, int(ollama.Llama3Dot1))
		loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
			"Creating grocery list with model: "+ollama.ModelName(model).String())
		loadingScreen.Show()

		recipe = addMarkdownFomatingToRecipe(recipe)

		generated, err := ollama.AskAI(ollamaClient, ollama.ModelName(model), recipe)
		if err != nil {
			slog.Error("Failed to ask AI", "error", err)
			loadingScreen.Hide()
			return
		}
		loadingScreen.Hide()
		recipePopUp(guiApp, recipe, &generated)
	})
	budgetFriendlyGroceryList.Importance = widget.SuccessImportance

	action := container.NewHBox(suggest, groceryList, suggestPrep, budgetFriendlyGroceryList)
	action.Layout = layout.NewAdaptiveGridLayout(2)

	boarderLayout := container.NewBorder(topText, action, nil, nil, tabs)
	mealPlanner.SetContent(boarderLayout)
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

func addMarkdownFomatingToRecipe(recipe string) string {
	return recipe + ". " + "format: markdown"
}

func recipePopUp(a fyne.App, recipe string, response *ollamaApi.GenerateResponse) {
	w := a.NewWindow("Ctrl+Revise AI Recipes")
	w.Resize(fyne.NewSize(640, 500))

	hello := widget.NewLabel("Let's Get Cooking!")
	hello.TextStyle = fyne.TextStyle{Bold: true}
	hello.Alignment = fyne.TextAlignCenter

	generatedText1 := widget.NewRichTextFromMarkdown(response.Response)
	generatedText1.Wrapping = fyne.TextWrapWord

	vbox := container.NewVScroll(generatedText1)

	model := guiApp.Preferences().IntWithFallback(CurrentModelKey, int(ollama.Llama3Dot1))

	buttons := container.NewPadded(
		widget.NewButton("Something Different", func() {
			loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
				"Using model: "+ollama.ModelName(model).String()+" with prompt: "+selectedPrompt.String()+"...")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithStringAndContext(ollamaClient, ollama.ModelName(model), response.Context,
				"That doesn't sound good, how about something else please.")
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			lastClipboardContent = sha256.Sum256([]byte(response.Response))
			w.Hide()
			loadingScreen.Hide()
			recipePopUp(a, recipe, &reGenerated)
		}),
		widget.NewButton("Pick a different Protein", func() {
			loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
				"Using model: "+ollama.ModelName(model).String()+" with prompt: "+selectedPrompt.String()+"...")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithStringAndContext(ollamaClient, ollama.ModelName(model), response.Context,
				"That doesn't sound good, how about something with a different protein please.")
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			lastClipboardContent = sha256.Sum256([]byte(response.Response))
			w.Hide()
			loadingScreen.Hide()
			recipePopUp(a, recipe, &reGenerated)
		}),
		widget.NewButton("Something Simple and Quick", func() {
			loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
				"Using model: "+ollama.ModelName(model).String()+" with prompt: "+selectedPrompt.String()+"...")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithStringAndContext(ollamaClient, ollama.ModelName(model), response.Context,
				"That doesn't sound good, how about something that's simple and quick to put together please.")
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			lastClipboardContent = sha256.Sum256([]byte(response.Response))
			w.Hide()
			loadingScreen.Hide()
			recipePopUp(a, recipe, &reGenerated)
		}),
		widget.NewButton("Something Healthy", func() {
			loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
				"Using model: "+ollama.ModelName(model).String()+" with prompt: "+selectedPrompt.String()+"...")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithStringAndContext(ollamaClient, ollama.ModelName(model), response.Context,
				"That doesn't sound good, how about something really healthy but still tasty please.")
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			lastClipboardContent = sha256.Sum256([]byte(response.Response))
			w.Hide()
			loadingScreen.Hide()
			recipePopUp(a, recipe, &reGenerated)
		}),
		widget.NewButton("Let's Not Cook", func() {
			loadingScreen := gui.LoadingScreenWithMessage(guiApp, thinkingMsg,
				"Using model: "+ollama.ModelName(model).String()+" with prompt: "+selectedPrompt.String()+"...")
			loadingScreen.Show()
			reGenerated, err := ollama.AskAiWithStringAndContext(ollamaClient, ollama.ModelName(model), response.Context,
				"I don't feel like cooking or heating anything up, how about something cold I could put together please.")
			if err != nil {
				slog.Error("Failed to re-generate", "error", err)
				return
			}
			lastClipboardContent = sha256.Sum256([]byte(response.Response))
			w.Hide()
			loadingScreen.Hide()
			recipePopUp(a, recipe, &reGenerated)
		}),
		widget.NewButtonWithIcon("Copy the Recipe to Clipboard", theme.ContentCopyIcon(), func() {
			w.Clipboard().SetContent(response.Response)
			w.Close()
		}),
	)
	buttons.Layout = layout.NewAdaptiveGridLayout(3)
	center := container.NewVBox(buttons, layout.NewSpacer(), footer())

	grid := container.New(layout.NewAdaptiveGridLayout(1), vbox)

	w.SetContent(container.NewBorder(
		hello,
		center,
		nil,
		nil,
		container.NewVScroll(grid),
	))
	w.Show()
}
