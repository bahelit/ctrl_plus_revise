package meal

import "fmt"

// MealInfo holds all the ingredients and preferences used to generate meal
// prompts.
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

// joinSlice joins a slice of strings with ", " separator.
func joinSlice(items []string) string {
	result := ""
	for i, item := range items {
		if i == 0 {
			result = item
		} else {
			result = result + ", " + item
		}
	}
	return result
}

// appendIngredients appends ingredient sections that are common across all
// prompt types.
func appendIngredients(recipe string, info MealInfo) string {
	if len(info.Cookware) != 0 {
		recipe = fmt.Sprintf("%s. For tools to work with in the kitchen we have %s", recipe, joinSlice(info.Cookware))
	}
	if len(info.Dairy) != 0 {
		recipe = fmt.Sprintf("%s. For Dairy products we have %s", recipe, joinSlice(info.Dairy))
	}
	if len(info.Veggies) != 0 {
		recipe = fmt.Sprintf("%s. For Vegetables we have %s", recipe, joinSlice(info.Veggies))
	}
	if len(info.Protein) != 0 {
		recipe = fmt.Sprintf("%s. For Proteins we have %s", recipe, joinSlice(info.Protein))
	}
	if len(info.Freezer) != 0 {
		recipe = fmt.Sprintf("%s. In the freezer we have %s", recipe, joinSlice(info.Freezer))
	}
	if len(info.Herbs) != 0 {
		recipe = fmt.Sprintf("%s. For herbs and spices we have %s", recipe, joinSlice(info.Herbs))
	}
	if len(info.Pantry) != 0 {
		recipe = fmt.Sprintf("%s. In the Pantry we have we have %s", recipe, joinSlice(info.Pantry))
	}
	if len(info.Allergies) != 0 {
		recipe = fmt.Sprintf("%s. Please be aware of our food allergy to %s", recipe, joinSlice(info.Allergies))
	}
	return recipe
}

// CreateMealPrompt builds a prompt asking the AI for a single meal recipe.
func CreateMealPrompt(info MealInfo) string {
	recipe := "Create a recipe with step-by-step instructions based on the following items I have available, please include the estimated cook time and cost"

	if info.Consumers == "" || info.Consumers == "1" {
		recipe = fmt.Sprintf("%s. Let's prepare %s", recipe, info.Meal)
	} else {
		recipe = fmt.Sprintf("%s. Let's prepare %s for %s people", recipe, info.Meal, info.Consumers)
	}
	if info.Theme != "" {
		recipe = fmt.Sprintf("%s,  with %s style dish", recipe, info.Theme)
	}

	return appendIngredients(recipe, info)
}

// CreateMealPrepPrompt builds a prompt asking the AI for a weekly meal prep
// plan.
func CreateMealPrepPrompt(info MealInfo) string {
	recipe := "Create a plan for doing food prep for the week in my kitchen with step-by-step instructions based on the following items I have available"

	if info.Consumers == "" || info.Consumers == "1" {
		recipe = fmt.Sprintf("%s. Let's prepare %s", recipe, info.Meal)
	} else {
		recipe = fmt.Sprintf("%s. Let's prepare %s for %s people", recipe, info.Meal, info.Consumers)
	}
	if info.Theme != "" {
		recipe = fmt.Sprintf("%s, with %s cuisine", recipe, info.Theme)
	}

	return appendIngredients(recipe, info)
}

// CreateGroceryListPrompt builds a prompt asking the AI for a weekly grocery
// list.
func CreateGroceryListPrompt(info MealInfo) string {
	recipe := "Create a grocery list that avoids processed foods for the week taking into consideration the items I have in stock"

	if info.Consumers == "" || info.Consumers == "1" {
		recipe = fmt.Sprintf("%s. I'am shopping for just myself", recipe)
	} else {
		recipe = fmt.Sprintf("%s. I'am shopping for %s people", recipe, info.Consumers)
	}
	if info.Theme != "" {
		recipe = fmt.Sprintf("%s, primarily shopping for %s cuisine", recipe, info.Theme)
	}

	return appendIngredients(recipe, info)
}

// CreateBudgetFriendlyGroceryListPrompt builds a prompt asking the AI for a
// budget-conscious grocery list.
func CreateBudgetFriendlyGroceryListPrompt(info MealInfo) string {
	recipe := "I don't have much money. Create a budget friendly grocery list for the week taking into consideration the items I have in stock, " +
		"if I have what I need to survive suggest items that will give me the best value for my money as well as last longest, " +
		"suggest tips and strategies for making the food go further and last longer," +
		"try to remain healthy by avoiding processed foods."

	if info.Consumers == "" || info.Consumers == "1" {
		recipe = fmt.Sprintf("%s. I'am shopping for just myself", recipe)
	} else {
		recipe = fmt.Sprintf("%s. I'am shopping for %s people", recipe, info.Consumers)
	}
	if info.Theme != "" {
		recipe = fmt.Sprintf("%s, primarily shopping for %s style dishes", recipe, info.Theme)
	}

	return appendIngredients(recipe, info)
}
