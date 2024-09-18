package food

func getMeals() []string {
	meals := []string{
		"Breakfast",
		"Breakfast for dinner",
		"Brunch",
		"Lunch",
		"Dinner",
		"Desert",
		"Snack",
	}
	return meals
}

func getThemes() []string {
	cuisines := []string{
		"American",
		"Brazilian",
		"British",
		"Caribbean",
		"Chinese",
		"French",
		"Greek",
		"Indian",
		"Italian",
		"Japanese",
		"Korean",
		"Mexican",
		"Spanish",
		"Thai",
		"Vietnamese",
	}
	return cuisines
}

func getCookWares() []string {
	cookWares := []string{
		"Air Fryer",
		"Baking Sheet",
		"Barbeque Grill",
		"Blender",
		"Box Grater",
		"Cake Pan",
		"Casserole Dish",
		"Cast Iron Skillet",
		"Chefs Knife",
		"Coffee Maker",
		"Colander",
		"Cooling Rack",
		"Cutting Board",
		"Dutch Oven",
		"Food Processor",
		"Fryer",
		"Garlic Press",
		"Grill Pan",
		"Hand Mixer",
		"Ice Cream Scoop",
		"Immersion Blender",
		"Instant Pot",
		"Kitchen Shears",
		"Ladle",
		"Loaf Pan",
		"Mandoline",
		"Measuring Cups",
		"Measuring Spoons",
		"Meat Thermometer",
		"Microplane",
		"Microwave",
		"Mixing Bowls",
		"Muffin Tin",
		"Non-Stick Pan",
		"Oven",
		"Oven Mitts",
		"Paring Knife",
		"Pastry Brush",
		"Pie Dish",
		"Potato Masher",
		"Pressure Cooker",
		"Rice Cooker",
		"Roasting Pan",
		"Rolling Pin",
		"Salad Spinner",
		"Saucepan",
		"Skillet",
		"Slotted Spoon",
		"Slow Cooker",
		"Smoker",
		"Spatula",
		"Steamer",
		"Stand Mixer",
		"Steamer Basket",
		"Stockpot",
		"Stove",
		"Toaster Oven",
		"Toaster",
		"Tongs",
		"Trivets",
		"Vegetable Peeler",
		"Whisk",
		"Zester",
	}
	return cookWares
}

func getDairy() []string {
	dairy := []string{
		"Almond Milk",
		"American Cheese",
		"Asiago Cheese",
		"Blue Cheese",
		"Brie Cheese",
		"Butter",
		"Buttermilk",
		"Camembert Cheese",
		"Cheddar Cheese",
		"Coconut Milk",
		"Colby Cheese",
		"Colby Jack Cheese",
		"Cottage Cheese",
		"Cream Cheese",
		"Edam Cheese",
		"Eggnog",
		"Feta Cheese",
		"Ghee",
		"Goat Cheese",
		"Gorgonzola Cheese",
		"Gouda Cheese",
		"Greek Yogurt",
		"Gruyere Cheese",
		"Half and Half",
		"Havarti Cheese",
		"Heavy Cream",
		"Jarlsberg Cheese",
		"Kefir",
		"Lactose-Free Milk",
		"Manchego Cheese",
		"Margarine",
		"Mascarpone Cheese",
		"Milk",
		"Monterey Jack Cheese",
		"Mozzarella Cheese",
		"Muenster Cheese",
		"Oat Milk",
		"Paneer Cheese",
		"Parmesan Cheese",
		"Pepper Jack Cheese",
		"Provolone Cheese",
		"Queso Blanco",
		"Queso Fresco",
		"Ricotta Cheese",
		"Smoked Cheese",
		"Sour Cream",
		"Soy Milk",
		"String Cheese",
		"Swiss Cheese",
		"Whipping Cream",
		"Yogurt",
	}
	return dairy
}

func getVegetables() []string {
	vegetables := []string{
		"Acorn Squash",
		"Artichoke",
		"Arugula",
		"Asparagus",
		"Beet",
		"Bell Pepper",
		"Bok Choy",
		"Broccoli",
		"Brussels Sprouts",
		"Butternut Squash",
		"Cabbage",
		"Carrot",
		"Cauliflower",
		"Celery",
		"Collard Greens",
		"Corn",
		"Cucumber",
		"Eggplant",
		"Endive",
		"Escarole",
		"Fennel",
		"Garlic",
		"Green Beans",
		"Jalapeño",
		"Kale",
		"Leek",
		"Lettuce",
		"Mushroom",
		"Okra",
		"Onion",
		"Parsnip",
		"Peas",
		"Potato",
		"Pumpkin",
		"Radicchio",
		"Radish",
		"Rutabaga",
		"Scallion",
		"Shallot",
		"Snow Peas",
		"Spinach",
		"Sweet Potato",
		"Swiss Chard",
		"Tomato",
		"Turnip",
		"Watercress",
		"Zucchini",
	}
	return vegetables
}

func getProteins() []string {
	proteins := []string{
		"Almond Butter",
		"Bacon",
		"Beef Steaks",
		"Black Beans",
		"Canned Chicken",
		"Canned Tuna",
		"Chicken Breasts",
		"Chicken Thighs",
		"Chickpeas",
		"Clams",
		"Cod Fillets",
		"Cottage Cheese",
		"Crab",
		"Deli Sliced Ham",
		"Deli Sliced Roast Beef",
		"Deli Sliced Turkey",
		"Edamame",
		"Eggs",
		"Greek Yogurt",
		"Ground Beef",
		"Ground Pork",
		"Ground Turkey",
		"Ham",
		"Hot Dogs",
		"Hummus",
		"Lamb Chops",
		"Lamb Roast",
		"Lentils",
		"Lobster",
		"Mussels",
		"Peanut Butter",
		"Pork Chops",
		"Quinoa",
		"Salmon Fillets",
		"Sausages",
		"Scallops",
		"Shrimp",
		"Tempeh",
		"Tilapia Fillets",
		"Tofu",
		"Tuna Steaks",
		"Turkey Breast",
		"Whole Chicken",
	}
	return proteins
}

func getFreezerStuffs() []string {
	freezerStuffs := []string{
		"Frozen Berries",
		"Frozen Blueberries",
		"Frozen Broccoli",
		"Frozen Cauliflower Rice",
		"Frozen Chicken Breasts",
		"Frozen Chicken Nuggets",
		"Frozen Cookie Dough",
		"Frozen Corn",
		"Frozen Dumplings",
		"Frozen Edamame",
		"Frozen Fish Sticks",
		"Frozen French Fries",
		"Frozen Garlic Bread",
		"Frozen Hamburger Patties",
		"Frozen Meatballs",
		"Frozen Mixed Vegetables",
		"Frozen Mozzarella Sticks",
		"Frozen Onion Rings",
		"Frozen Pancakes",
		"Frozen Peas",
		"Frozen Pie Crust",
		"Frozen Raspberries",
		"Frozen Sausages",
		"Frozen Shrimp",
		"Frozen Smoothie Packs",
		"Frozen Spinach",
		"Frozen Strawberries",
		"Frozen Tater Tots",
		"Frozen Vegetables Stir Fry",
		"Frozen Waffles",
		"Ice Cream",
	}
	return freezerStuffs
}

func getHerbsSpices() []string {
	herbsAndSpices := []string{
		"Allspice",
		"Anise",
		"Basil",
		"Bay Leaf",
		"Bay Leaves",
		"Black Pepper",
		"Caraway Seeds",
		"Cardamom",
		"Cayenne Pepper",
		"Celery Seeds",
		"Chervil",
		"Chili Powder",
		"Chives",
		"Cilantro",
		"Cinnamon",
		"Cloves",
		"Coriander",
		"Cumin",
		"Curry Powder",
		"Dill",
		"Epazote",
		"Fennel Seeds",
		"Fenugreek Leaves",
		"Fenugreek",
		"Galangal",
		"Garlic Powder",
		"Ginger",
		"Horseradish",
		"Juniper Berries",
		"Kaffir Lime Leaves",
		"Lemongrass",
		"Lovage",
		"Marjoram",
		"Mint",
		"Mustard Seeds",
		"Nutmeg",
		"Onion Powder",
		"Oregano",
		"Paprika",
		"Parsley",
		"Rosemary",
		"Saffron",
		"Sage",
		"Savory",
		"Smoked Paprika",
		"Star Anise",
		"Sumac",
		"Tarragon",
		"Thyme",
		"Turmeric",
		"Vanilla Beans",
		"White Pepper",
		"Za'atar",
	}
	return herbsAndSpices
}

func getPantryStuffs() []string {
	pantryStuffs := []string{
		"Almonds",
		"Baking Powder",
		"Baking Soda",
		"Beef Broth",
		"Black Pepper",
		"Bread",
		"Brown Sugar",
		"Cake Mix",
		"Canned Beans",
		"Canned Chicken",
		"Canned Chili",
		"Canned Corn",
		"Canned Pasta",
		"Canned Tomatoes",
		"Canned Tuna",
		"Cashews",
		"Chicken Broth",
		"Coconut Milk",
		"Condensed Milk",
		"Crackers",
		"Dried Figs",
		"Dried Mango",
		"Evaporated Milk",
		"Flour",
		"Granola",
		"Honey",
		"Hot Sauce",
		"Instant Mashed Potatoes",
		"Instant Noodles",
		"Jelly",
		"Ketchup",
		"Macaroni and Cheese",
		"Maple Syrup",
		"Mayonnaise",
		"Mustard",
		"Oatmeal",
		"Olive Oil",
		"Pancake Mix",
		"Pasta",
		"Peanut Butter",
		"Peanuts",
		"Pecans",
		"Popcorn",
		"Potatoes",
		"Pudding Mix",
		"Quinoa",
		"Rice",
		"Salsa",
		"Salt",
		"Soy Sauce",
		"Sugar",
		"Tomato Sauce",
		"Tortillas",
		"Vegetable Broth",
		"Vegetable Oil",
		"Vinegar",
	}
	return pantryStuffs
}

func getAllergies() []string {
	foodAllergies := []string{
		"Peanut",
		"Tree Nut",
		"Milk",
		"Egg",
		"Wheat",
		"Soy",
		"Fish",
		"Shellfish",
		"Sesame",
		"Mustard",
		"Corn",
		"Gluten",
		"Sulfite",
		"Celery",
		"Lupin",
	}
	return foodAllergies
}