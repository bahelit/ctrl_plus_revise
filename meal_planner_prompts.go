package main

import (
	"fmt"
	"log/slog"
)

func createMealPrompt(mealInfo MealInfo) string {
	recipe := "Create a recipe with step-by-step instructions based on the following items I have available, please include the estimated cook time and cost"

	if mealInfo.Consumers == "" || mealInfo.Consumers == "1" {
		recipe = fmt.Sprintf("%s. Let's prepare %s", recipe, mealInfo.Meal)
	} else {
		recipe = fmt.Sprintf("%s. Let's prepare %s for %s people", recipe, mealInfo.Meal, mealInfo.Consumers)
	}
	if mealInfo.Theme != "" {
		recipe = fmt.Sprintf("%s,  with %s cuisine", recipe, mealInfo.Theme)
	}
	if len(mealInfo.Cookware) != 0 {
		var cookware string
		for i, c := range mealInfo.Cookware {
			if i == 0 {
				cookware = cookware + c
			} else {
				cookware = cookware + ", " + c
			}
		}
		recipe = fmt.Sprintf("%s. For tools to work with in the kitchen we have %s", recipe, cookware)
	}
	if len(mealInfo.Dairy) != 0 {
		var dairy string
		for i, d := range mealInfo.Dairy {
			if i == 0 {
				dairy = dairy + d
			} else {
				dairy = dairy + ", " + d
			}
		}
		recipe = fmt.Sprintf("%s. For Dairy products we have %s", recipe, dairy)
	}
	if len(mealInfo.Veggies) != 0 {
		var vegetable string
		for i, v := range mealInfo.Veggies {
			if i == 0 {
				vegetable = vegetable + v
			} else {
				vegetable = vegetable + ", " + v
			}
		}
		recipe = fmt.Sprintf("%s. For Vegetables we have %s", recipe, vegetable)
	}
	if len(mealInfo.Protein) != 0 {
		var protein string
		for i, p := range mealInfo.Protein {
			if i == 0 {
				protein = protein + p
			} else {
				protein = protein + ", " + p
			}
		}
		recipe = fmt.Sprintf("%s. For Proteins we have %s", recipe, protein)
	}
	if len(mealInfo.Freezer) != 0 {
		var freezer string
		for i, f := range mealInfo.Freezer {
			if i == 0 {
				freezer = freezer + f
			} else {
				freezer = freezer + ", " + f
			}
		}
		recipe = fmt.Sprintf("%s. In the freezer we have %s", recipe, freezer)
	}
	if len(mealInfo.Herbs) != 0 {
		var herbsAndSpices string
		for i, h := range mealInfo.Herbs {
			if i == 0 {
				herbsAndSpices = herbsAndSpices + h
			} else {
				herbsAndSpices = herbsAndSpices + ", " + h
			}
		}
		recipe = fmt.Sprintf("%s. For herbs and spices we have %s", recipe, herbsAndSpices)
	}
	if len(mealInfo.Pantry) != 0 {
		var pantry string
		for i, p := range mealInfo.Pantry {
			if i == 0 {
				pantry = pantry + p
			} else {
				pantry = pantry + ", " + p
			}
		}
		recipe = fmt.Sprintf("%s. In the Pantry we have we have %s", recipe, pantry)
	}
	if len(mealInfo.Allergies) != 0 {
		var allergy string
		for i, a := range mealInfo.Allergies {
			if i == 0 {
				allergy = allergy + a
			} else {
				allergy = allergy + ", " + a
			}
		}
		recipe = fmt.Sprintf("%s. Please be aware of our food allergy to %s", recipe, allergy)
	}
	return recipe
}

func createMealPrepPrompt(mealInfo MealInfo) string {
	recipe := "Create a plan for doing food prep for the week in my kitchen with step-by-step instructions based on the following items I have available"

	if mealInfo.Consumers == "" || mealInfo.Consumers == "1" {
		recipe = fmt.Sprintf("%s. Let's prepare %s", recipe, mealInfo.Meal)
	} else {
		recipe = fmt.Sprintf("%s. Let's prepare %s for %s people", recipe, mealInfo.Meal, mealInfo.Consumers)
	}
	if mealInfo.Theme != "" {
		recipe = fmt.Sprintf("%s, with %s cuisine", recipe, mealInfo.Theme)
	}
	if len(mealInfo.Cookware) != 0 {
		var cookware string
		for i, c := range mealInfo.Cookware {
			if i == 0 {
				cookware = cookware + c
			} else {
				cookware = cookware + ", " + c
			}
		}
		recipe = fmt.Sprintf("%s. For tools to work with in the kitchen we have %s", recipe, cookware)
	}
	if len(mealInfo.Dairy) != 0 {
		var dairy string
		for i, d := range mealInfo.Dairy {
			if i == 0 {
				dairy = dairy + d
			} else {
				dairy = dairy + ", " + d
			}
		}
		recipe = fmt.Sprintf("%s. For Dairy products we have %s", recipe, dairy)
	}
	if len(mealInfo.Veggies) != 0 {
		var vegetable string
		for i, v := range mealInfo.Veggies {
			if i == 0 {
				vegetable = vegetable + v
			} else {
				vegetable = vegetable + ", " + v
			}
		}
		recipe = fmt.Sprintf("%s. For Vegetables we have %s", recipe, vegetable)
	}
	if len(mealInfo.Protein) != 0 {
		var protein string
		for i, p := range mealInfo.Protein {
			if i == 0 {
				protein = protein + p
			} else {
				protein = protein + ", " + p
			}
		}
		recipe = fmt.Sprintf("%s. For Proteins we have %s", recipe, protein)
	}
	if len(mealInfo.Freezer) != 0 {
		var freezer string
		for i, f := range mealInfo.Freezer {
			if i == 0 {
				freezer = freezer + f
			} else {
				freezer = freezer + ", " + f
			}
		}
		recipe = fmt.Sprintf("%s. In the freezer we have %s", recipe, freezer)
	}
	if len(mealInfo.Herbs) != 0 {
		var herbsAndSpices string
		for i, h := range mealInfo.Herbs {
			if i == 0 {
				herbsAndSpices = herbsAndSpices + h
			} else {
				herbsAndSpices = herbsAndSpices + ", " + h
			}
		}
		recipe = fmt.Sprintf("%s. For herbs and spices we have %s", recipe, herbsAndSpices)
	}
	if len(mealInfo.Pantry) != 0 {
		var pantry string
		for i, p := range mealInfo.Pantry {
			if i == 0 {
				pantry = pantry + p
			} else {
				pantry = pantry + ", " + p
			}
		}
		recipe = fmt.Sprintf("%s. In the Pantry we have we have %s", recipe, pantry)
	}
	if len(mealInfo.Allergies) != 0 {
		var allergy string
		for i, a := range mealInfo.Allergies {
			if i == 0 {
				allergy = allergy + a
			} else {
				allergy = allergy + ", " + a
			}
		}
		recipe = fmt.Sprintf("%s. Please be aware of our food allergy to %s", recipe, allergy)
	}
	return recipe
}

func createGroceryListPrompt(mealInfo MealInfo) string {
	recipe := "Create a health focused grocery list for the week taking into consideration the items I have in stock"

	if mealInfo.Consumers == "" || mealInfo.Consumers == "1" {
		recipe = fmt.Sprintf("%s. I'am shopping for just myself", recipe)
	} else {
		recipe = fmt.Sprintf("%s. I'am shopping for %s people", recipe, mealInfo.Consumers)
	}
	if mealInfo.Theme != "" {
		recipe = fmt.Sprintf("%s, primarily shopping for %s cuisine", recipe, mealInfo.Theme)
	}
	if len(mealInfo.Cookware) != 0 {
		var cookware string
		for i, c := range mealInfo.Cookware {
			if i == 0 {
				cookware = cookware + c
			} else {
				cookware = cookware + ", " + c
			}
		}
		recipe = fmt.Sprintf("%s. For tools to work with in the kitchen we have %s", recipe, cookware)
	}
	if len(mealInfo.Dairy) != 0 {
		var dairy string
		for i, d := range mealInfo.Dairy {
			if i == 0 {
				dairy = dairy + d
			} else {
				dairy = dairy + ", " + d
			}
		}
		recipe = fmt.Sprintf("%s. For Dairy products we have %s", recipe, dairy)
	}
	if len(mealInfo.Veggies) != 0 {
		var vegetable string
		for i, v := range mealInfo.Veggies {
			if i == 0 {
				vegetable = vegetable + v
			} else {
				vegetable = vegetable + ", " + v
			}
		}
		recipe = fmt.Sprintf("%s. For Vegetables we have %s", recipe, vegetable)
	}
	if len(mealInfo.Protein) != 0 {
		var protein string
		for i, p := range mealInfo.Protein {
			if i == 0 {
				protein = protein + p
			} else {
				protein = protein + ", " + p
			}
		}
		recipe = fmt.Sprintf("%s. For Proteins we have %s", recipe, protein)
	}
	if len(mealInfo.Freezer) != 0 {
		var freezer string
		for i, f := range mealInfo.Freezer {
			if i == 0 {
				freezer = freezer + f
			} else {
				freezer = freezer + ", " + f
			}
		}
		recipe = fmt.Sprintf("%s. In the freezer we have %s", recipe, freezer)
	}
	if len(mealInfo.Herbs) != 0 {
		var herbsAndSpices string
		for i, h := range mealInfo.Herbs {
			if i == 0 {
				herbsAndSpices = herbsAndSpices + h
			} else {
				herbsAndSpices = herbsAndSpices + ", " + h
			}
		}
		recipe = fmt.Sprintf("%s. For herbs and spices we have %s", recipe, herbsAndSpices)
	}
	if len(mealInfo.Pantry) != 0 {
		var pantry string
		for i, p := range mealInfo.Pantry {
			if i == 0 {
				pantry = pantry + p
			} else {
				pantry = pantry + ", " + p
			}
		}
		recipe = fmt.Sprintf("%s. In the Pantry we have we have %s", recipe, pantry)
	}
	if len(mealInfo.Allergies) != 0 {
		var allergy string
		for i, a := range mealInfo.Allergies {
			if i == 0 {
				allergy = allergy + a
			} else {
				allergy = allergy + ", " + a
			}
		}
		recipe = fmt.Sprintf("%s. Please be aware of our food allergy to %s", recipe, allergy)
	}
	return recipe
}

func createSpecialRequestPrompt() string {
	slog.Info("Not Yet Implemented")
	return ""
}
