package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/RobiGavranovic/NutritionWebApp/backend/models"
	"github.com/gin-gonic/gin"
)

func GetRandomMeal() (models.Meals, error) {
	// Get random meal by calling MealDB API
	resp, err := http.Get("https://www.themealdb.com/api/json/v1/1/random.php")
	if err != nil {
		log.Println("Error fetching meal:", err)
		return models.Meals{}, err
	}
	defer resp.Body.Close()

	// Read the Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return models.Meals{}, err
	}

	// Transform Data To Struct
	var meals models.Meals
	if err := json.Unmarshal(body, &meals); err != nil {
		log.Println("Error decoding JSON:", err)
		return models.Meals{}, err
	}

	// Get Meal's Allergens
	if len(meals.Meals) > 0 {
		meal := &meals.Meals[0]

		// Get All Ingredients
		ingredients := []string{
			meal.StrIngredient1, meal.StrIngredient2, meal.StrIngredient3,
			meal.StrIngredient4, meal.StrIngredient5, meal.StrIngredient6,
			meal.StrIngredient7, meal.StrIngredient8, meal.StrIngredient9,
			meal.StrIngredient10, meal.StrIngredient11, meal.StrIngredient12,
			meal.StrIngredient13, meal.StrIngredient14, meal.StrIngredient15,
			meal.StrIngredient16, meal.StrIngredient17, meal.StrIngredient18,
			meal.StrIngredient19, meal.StrIngredient20,
		}

		// Get Non-Empty Ingredients
		nonEmptyIngredients := []string{}
		for _, ing := range ingredients {
			if ing != "" {
				nonEmptyIngredients = append(nonEmptyIngredients, ing)
			}
		}

		// Set Allergens
		meal.Allergens = findAllergens(nonEmptyIngredients)

		// Set Intolerances
		meal.Intolerances = findIntolerances(nonEmptyIngredients)
	}

	return meals, err

}

func GetNRandomMeals(c *gin.Context) {
	// Get number of meals from URL
	numOfMealsStr := c.Param("numOfMeals")
	numOfMeals, err := strconv.Atoi(numOfMealsStr)
	if err != nil || numOfMeals <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of meals"})
		return
	}

	meals := models.Meals{
		Meals: []models.Meal{},
	}

	// Map used to track already added meal IDS (avoid duplicates)
	seen := make(map[string]bool)

	// Fetch meals
	maxRetries := 10

	for len(meals.Meals) < numOfMeals && maxRetries > 0 {
		mealWrapper, err := GetRandomMeal()
		if err != nil || len(mealWrapper.Meals) == 0 {
			maxRetries--
			continue
		}

		meal := mealWrapper.Meals[0]
		id := meal.IDMeal

		if !seen[id] {
			meals.Meals = append(meals.Meals, meal)
			seen[id] = true
		} else {
			// Duplicate Found
			maxRetries--
		}
	}

	c.JSON(http.StatusOK, meals)

}

func SearchRecipesByOrigin(c *gin.Context) {
	// Get Origin From URL
	origin := c.Param("origin")
	if origin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Origin parameter is required"})
		return
	}

	// Search All Meals With Origin
	apiURL := fmt.Sprintf("https://www.themealdb.com/api/json/v1/1/filter.php?a=%s", origin)

	resp, err := http.Get(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call external API"})
		return
	}
	defer resp.Body.Close()

	var base struct {
		Meals []struct {
			IDMeal string `json:"idMeal"`
		} `json:"meals"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&base); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode filter API response"})
		return
	}

	if base.Meals == nil {
		c.JSON(http.StatusOK, gin.H{"Meals": nil})
		return
	}

	// Since Origin Search Lacks Data, Do ID Lookup For Each Meal
	detailedMeals := []map[string]interface{}{}
	for _, m := range base.Meals {
		detailResp, err := http.Get("https://www.themealdb.com/api/json/v1/1/lookup.php?i=" + m.IDMeal)
		if err != nil {
			continue
		}
		defer detailResp.Body.Close()

		var detailData struct {
			Meals []map[string]interface{} `json:"meals"`
		}
		if err := json.NewDecoder(detailResp.Body).Decode(&detailData); err != nil {
			continue
		}

		if len(detailData.Meals) > 0 {
			// Change The Structure Names So It Fits The FE Expectations
			meal := detailData.Meals[0]
			normalized := make(map[string]interface{})
			for k, v := range meal {
				// Capitalize first letter
				if len(k) > 0 {
					normalizedKey := strings.ToUpper(string(k[0])) + k[1:]
					normalized[normalizedKey] = v
				}
			}
			detailedMeals = append(detailedMeals, normalized)
		}
	}

	c.JSON(http.StatusOK, gin.H{"meals": detailedMeals})
}

func SearchRecipesByName(c *gin.Context) {
	// Get Name From URL
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name parameter is required"})
		return
	}

	// Search All Meals With Origin
	apiURL := fmt.Sprintf("https://www.themealdb.com/api/json/v1/1/search.php?s=%s", name)

	resp, err := http.Get(apiURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data"})
		return
	}
	defer resp.Body.Close()

	var base struct {
		Meals []struct {
			IDMeal string `json:"idMeal"`
		} `json:"meals"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&base); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
		return
	}

	if base.Meals == nil {
		c.JSON(http.StatusOK, gin.H{"Meals": nil})
		return
	}

	// Since Origin Search Lacks Data, Do ID Lookup For Each Meal
	detailedMeals := []map[string]interface{}{}
	for _, m := range base.Meals {
		detailResp, err := http.Get("https://www.themealdb.com/api/json/v1/1/lookup.php?i=" + m.IDMeal)
		if err != nil {
			continue
		}
		defer detailResp.Body.Close()

		var detailData struct {
			Meals []map[string]interface{} `json:"meals"`
		}
		if err := json.NewDecoder(detailResp.Body).Decode(&detailData); err != nil {
			continue
		}

		if len(detailData.Meals) > 0 {
			// Change The Structure Names So It Fits The FE Expectations
			meal := detailData.Meals[0]
			normalized := make(map[string]interface{})
			for k, v := range meal {
				// Capitalize First Letter
				if len(k) > 0 {
					normalizedKey := strings.ToUpper(string(k[0])) + k[1:]
					normalized[normalizedKey] = v
				}
			}
			detailedMeals = append(detailedMeals, normalized)
		}
	}

	c.JSON(http.StatusOK, gin.H{"meals": detailedMeals})
}

func findAllergens(ingredients []string) []string {
	var Celery = map[string]struct{}{
		"Celeriac":                  {},
		"Celery":                    {},
		"Celery Salt":               {},
		"Vegetable Stock":           {},
		"Vegetable Stock Cube":      {},
		"Chicken Stock":             {},
		"Chicken Stock Cube":        {},
		"Beef Stock":                {},
		"Beef Stock Concentrate":    {},
		"Chicken Stock Concentrate": {},
		"Hot Beef Stock":            {},
		"Bouquet Garni":             {},
	}

	var Gluten = map[string]struct{}{
		"Bread":                   {},
		"Breadcrumbs":             {},
		"Bowtie Pasta":            {},
		"Farfalle":                {},
		"Macaroni":                {},
		"Lasagne Sheets":          {},
		"Rigatoni":                {},
		"Spaghetti":               {},
		"Tagliatelle":             {},
		"Fettuccine":              {},
		"Pappardelle Pasta":       {},
		"Paccheri Pasta":          {},
		"Linguine Pasta":          {},
		"Fideo":                   {},
		"Vermicelli Pasta":        {},
		"Noodles":                 {},
		"Udon Noodles":            {},
		"Pretzels":                {},
		"English Muffins":         {},
		"Muffins":                 {},
		"White Flour":             {},
		"Flour":                   {},
		"Plain Flour":             {},
		"Self-raising Flour":      {},
		"Whole Wheat":             {},
		"Shortcrust Pastry":       {},
		"Filo Pastry":             {},
		"Tortillas":               {},
		"Flour Tortilla":          {},
		"Wonton Skin":             {},
		"Pita Bread":              {},
		"Buns":                    {},
		"Bread Rolls":             {},
		"Potatoe Buns":            {},
		"Sesame Seed Burger Buns": {},
		"Ciabatta":                {},
		"Oatmeal":                 {},
		"Oats":                    {},
		"Couscous":                {},
		"Bulgar Wheat":            {},
		"Semolina":                {},
		"Barley":                  {},
		"Wheat":                   {},
	}

	var Crustaceans = map[string]struct{}{
		"King Prawns":     {},
		"Prawns":          {},
		"Raw King Prawns": {},
		"Tiger Prawns":    {},
		"Baby Squid":      {},
		"Squid":           {},
		"Clams":           {},
		"Oysters":         {},
	}

	var Eggs = map[string]struct{}{
		"Egg Rolls":               {},
		"Egg White":               {},
		"Egg Yolks":               {},
		"Eggs":                    {},
		"Free-range Egg, Beaten":  {},
		"Free-range Eggs, Beaten": {},
		"Flax Eggs":               {},
		"Meringue Nests":          {},
		"Custard":                 {},
		"Custard Powder":          {},
		"Mayonnaise":              {},
		"Egg":                     {},
	}

	var Fish = map[string]struct{}{
		"Salmon":             {},
		"Mackerel":           {},
		"Tuna":               {},
		"White Fish":         {},
		"White Fish Fillets": {},
		"Haddock":            {},
		"Smoked Haddock":     {},
		"Pilchards":          {},
		"Monkfish":           {},
		"Cod":                {},
		"Salt Cod":           {},
		"Red Snapper":        {},
		"Sardines":           {},
		"Anchovy Fillet":     {},
		"Herring":            {},
		"Fish Sauce":         {},
		"Thai Fish Sauce":    {},
		"Fish Stock":         {},
	}

	var Lupin = map[string]struct{}{
		"Bread":                   {},
		"Breadcrumbs":             {},
		"Digestive Biscuits":      {},
		"Flour":                   {},
		"Flour Tortilla":          {},
		"Plain Flour":             {},
		"Self-raising Flour":      {},
		"Shortcrust Pastry":       {},
		"Filo Pastry":             {},
		"Naan Bread":              {},
		"Tortillas":               {},
		"Udon Noodles":            {},
		"Baguette":                {},
		"Crusty Bread":            {},
		"Couscous":                {},
		"Whole Wheat":             {},
		"Wholegrain Bread":        {},
		"English Muffins":         {},
		"Muffins":                 {},
		"White Flour":             {},
		"Bread Rolls":             {},
		"Bun":                     {},
		"Potatoe Buns":            {},
		"Sesame Seed Burger Buns": {},
		"Buns":                    {},
		"Ciabatta":                {},
	}

	var Milk = map[string]struct{}{
		"Butter":                        {},
		"Cheddar Cheese":                {},
		"Cheese":                        {},
		"Cheese Curds":                  {},
		"Chilled Butter":                {},
		"Colby Jack Cheese":             {},
		"Condensed Milk":                {},
		"Cream":                         {},
		"Creme Fraiche":                 {},
		"Cubed Feta Cheese":             {},
		"Double Cream":                  {},
		"Feta":                          {},
		"Full Fat Yogurt":               {},
		"Ghee":                          {},
		"Gouda Cheese":                  {},
		"Greek Yogurt":                  {},
		"Gruyère":                       {},
		"Heavy Cream":                   {},
		"Ice Cream":                     {},
		"Milk":                          {},
		"Monterey Jack Cheese":          {},
		"Mozzarella Balls":              {},
		"Parmesan":                      {},
		"Parmesan Cheese":               {},
		"Parmigiano-reggiano":           {},
		"Pecorino":                      {},
		"Salted Butter":                 {},
		"Semi-skimmed Milk":             {},
		"Shredded Mexican Cheese":       {},
		"Shredded Monterey Jack Cheese": {},
		"Sour Cream":                    {},
		"Whole Milk":                    {},
		"Yogurt":                        {},
		"Cream Cheese":                  {},
		"Clotted Cream":                 {},
		"Fromage Frais":                 {},
		"Stilton Cheese":                {},
		"Mascarpone":                    {},
		"Mozzarella":                    {},
		"Ricotta":                       {},
		"Custard":                       {},
		"Single Cream":                  {},
		"Goats Cheese":                  {},
		"Unsalted Butter":               {},
		"Brie":                          {},
		"Cheese Slices":                 {},
	}

	var Molluscs = map[string]struct{}{
		"Mussels":    {},
		"Baby Squid": {},
		"Squid":      {},
		"Clams":      {},
		"Oysters":    {},
	}

	var Mustard = map[string]struct{}{
		"English Mustard": {},
		"Mustard":         {},
		"Mustard Powder":  {},
		"Mustard Seeds":   {},
		"Dijon Mustard":   {},
	}

	var Nuts = map[string]struct{}{
		"Cashew Nuts":    {},
		"Cashews":        {},
		"Flaked Almonds": {},
		"Ground Almonds": {},
		"Peanut Butter":  {},
		"Peanut Oil":     {},
		"Peanuts":        {},
		"Walnuts":        {},
		"Pecan Nuts":     {},
		"Peanut Cookies": {},
		"Peanut Brittle": {},
		"Hazelnuts":      {},
		"Almonds":        {},
		"Almond Milk":    {},
		"Almond Extract": {},
	}

	var Peanuts = map[string]struct{}{
		"Peanut Butter":  {},
		"Peanut Oil":     {},
		"Peanuts":        {},
		"Peanut Cookies": {},
		"Peanut Brittle": {},
	}

	var SesameSeeds = map[string]struct{}{
		"Sesame Seed":             {},
		"Sesame Seed Oil":         {},
		"Sesame Seed Burger Buns": {},
	}

	var Soya = map[string]struct{}{
		"Dark Soy Sauce": {},
		"Soy Sauce":      {},
		"Soya Milk":      {},
		"Tofu":           {},
	}

	var SulphurDioxide = map[string]struct{}{
		"Brandy":             {},
		"Christmas Pudding":  {},
		"Dried Apricots":     {},
		"Muscovado Sugar":    {},
		"Red Wine":           {},
		"White Wine":         {},
		"Balsamic Vinegar":   {},
		"Red Wine Vinegar":   {},
		"White Wine Vinegar": {},
		"Sherry":             {},
		"Cider":              {},
		"Mars Bar":           {},
		"Pickle Juice":       {},
		"Custard":            {},
		"Dried Fruit":        {},
		"Fruit Mix":          {},
		"Mixed Peel":         {},
	}

	allergenMaps := map[string]map[string]struct{}{
		"Celery":          Celery,
		"Gluten":          Gluten,
		"Crustaceans":     Crustaceans,
		"Eggs":            Eggs,
		"Fish":            Fish,
		"Lupin":           Lupin,
		"Milk":            Milk,
		"Molluscs":        Molluscs,
		"Mustard":         Mustard,
		"Nuts":            Nuts,
		"Peanuts":         Peanuts,
		"Sesame seeds":    SesameSeeds,
		"Soya":            Soya,
		"Sulphur dioxide": SulphurDioxide,
	}

	found := make(map[string]bool)
	var result []string

	for _, ingredient := range ingredients {
		for allergenName, allergenMap := range allergenMaps {
			if _, exists := allergenMap[ingredient]; exists {
				if !found[allergenName] {
					result = append(result, allergenName)
					found[allergenName] = true
				}
			}
		}
	}

	return result
}

func findIntolerances(ingredients []string) []string {
	var Lactose = map[string]struct{}{
		"Butter":                        {},
		"Cheddar Cheese":                {},
		"Cheese":                        {},
		"Cheese Curds":                  {},
		"Chilled Butter":                {},
		"Colby Jack Cheese":             {},
		"Condensed Milk":                {},
		"Cream":                         {},
		"Creme Fraiche":                 {},
		"Cubed Feta Cheese":             {},
		"Double Cream":                  {},
		"Feta":                          {},
		"Full Fat Yogurt":               {},
		"Ghee":                          {},
		"Gouda Cheese":                  {},
		"Greek Yogurt":                  {},
		"Gruyère":                       {},
		"Heavy Cream":                   {},
		"Ice Cream":                     {},
		"Milk":                          {},
		"Monterey Jack Cheese":          {},
		"Mozzarella Balls":              {},
		"Parmesan":                      {},
		"Parmesan Cheese":               {},
		"Parmigiano-reggiano":           {},
		"Pecorino":                      {},
		"Salted Butter":                 {},
		"Semi-skimmed Milk":             {},
		"Shredded Mexican Cheese":       {},
		"Shredded Monterey Jack Cheese": {},
		"Sour Cream":                    {},
		"Whole Milk":                    {},
		"Yogurt":                        {},
		"Cream Cheese":                  {},
		"Clotted Cream":                 {},
		"Fromage Frais":                 {},
		"Stilton Cheese":                {},
		"Mascarpone":                    {},
		"Mozzarella":                    {},
		"Ricotta":                       {},
		"Custard":                       {},
		"Single Cream":                  {},
		"Goats Cheese":                  {},
		"Unsalted Butter":               {},
		"Brie":                          {},
		"Cheese Slices":                 {},
	}

	Gluten := map[string]struct{}{
		"Bread":                   {},
		"Breadcrumbs":             {},
		"Bowtie Pasta":            {},
		"Farfalle":                {},
		"Macaroni":                {},
		"Lasagne Sheets":          {},
		"Rigatoni":                {},
		"Spaghetti":               {},
		"Tagliatelle":             {},
		"Fettuccine":              {},
		"Pappardelle Pasta":       {},
		"Paccheri Pasta":          {},
		"Linguine Pasta":          {},
		"Fideo":                   {},
		"Vermicelli Pasta":        {},
		"Noodles":                 {},
		"Udon Noodles":            {},
		"Pretzels":                {},
		"English Muffins":         {},
		"Muffins":                 {},
		"White Flour":             {},
		"Flour":                   {},
		"Plain Flour":             {},
		"Self-raising Flour":      {},
		"Whole Wheat":             {},
		"Shortcrust Pastry":       {},
		"Filo Pastry":             {},
		"Tortillas":               {},
		"Flour Tortilla":          {},
		"Wonton Skin":             {},
		"Pita Bread":              {},
		"Buns":                    {},
		"Bread Rolls":             {},
		"Potatoe Buns":            {},
		"Sesame Seed Burger Buns": {},
		"Ciabatta":                {},
		"Oatmeal":                 {},
		"Oats":                    {},
		"Couscous":                {},
		"Bulgar Wheat":            {},
		"Semolina":                {},
		"Barley":                  {},
		"Wheat":                   {},
	}

	Histamine := map[string]struct{}{
		"Salmon":                {},
		"Mackerel":              {},
		"Tuna":                  {},
		"Haddock":               {},
		"Smoked Haddock":        {},
		"Smoked Salmon":         {},
		"Anchovy Fillet":        {},
		"Sardines":              {},
		"Pilchards":             {},
		"Squid":                 {},
		"Baby Squid":            {},
		"Mussels":               {},
		"Clams":                 {},
		"Oysters":               {},
		"Monkfish":              {},
		"Fish Sauce":            {},
		"Thai Fish Sauce":       {},
		"Fish Stock":            {},
		"Cheddar Cheese":        {},
		"Cheese":                {},
		"Cheese Curds":          {},
		"Gouda Cheese":          {},
		"Greek Yogurt":          {},
		"Gruyère":               {},
		"Parmesan":              {},
		"Parmesan Cheese":       {},
		"Parmigiano-reggiano":   {},
		"Pecorino":              {},
		"Stilton Cheese":        {},
		"Mascarpone":            {},
		"Mozzarella":            {},
		"Ricotta":               {},
		"Brie":                  {},
		"Goats Cheese":          {},
		"Fromage Frais":         {},
		"Yogurt":                {},
		"Full Fat Yogurt":       {},
		"Cream Cheese":          {},
		"Clotted Cream":         {},
		"Single Cream":          {},
		"Double Cream":          {},
		"Butter":                {},
		"Unsalted Butter":       {},
		"Chilled Butter":        {},
		"Salted Butter":         {},
		"Vinegar":               {},
		"Balsamic Vinegar":      {},
		"Red Wine Vinegar":      {},
		"White Wine Vinegar":    {},
		"Apple Cider Vinegar":   {},
		"Malt Vinegar":          {},
		"Rice Vinegar":          {},
		"White Vinegar":         {},
		"Red Wine":              {},
		"White Wine":            {},
		"Sake":                  {},
		"Sherry":                {},
		"Cooking wine":          {},
		"Rice wine":             {},
		"Dark Rum":              {},
		"Light Rum":             {},
		"Rum":                   {},
		"Grand Marnier":         {},
		"Brandy":                {},
		"Beer":                  {},
		"Stout":                 {},
		"Yeast":                 {},
		"Sauerkraut":            {},
		"Pickle Juice":          {},
		"Dill Pickles":          {},
		"Gherkin Relish":        {},
		"Pickled Onions":        {},
		"Chorizo":               {},
		"Sausages":              {},
		"Polish Sausage":        {},
		"Kielbasa":              {},
		"Parma Ham":             {},
		"Prosciutto":            {},
		"Ham":                   {},
		"Bacon":                 {},
		"Anchovy":               {},
		"Spinach":               {},
		"Tomatoes":              {},
		"Plum Tomatoes":         {},
		"Canned Tomatoes":       {},
		"Chopped Tomatoes":      {},
		"Grape Tomatoes":        {},
		"Baby Plum Tomatoes":    {},
		"Tomato Puree":          {},
		"Tomato Ketchup":        {},
		"Tomato Sauce":          {},
		"Aubergine":             {},
		"Avocado":               {},
		"Chocolate Chips":       {},
		"Milk Chocolate":        {},
		"Dark Chocolate":        {},
		"White Chocolate":       {},
		"Plain Chocolate":       {},
		"Cocoa":                 {},
		"Cacao":                 {},
		"Soy Sauce":             {},
		"Soya Milk":             {},
		"Tofu":                  {},
		"Miso":                  {},
		"Fermented Black Beans": {},
		"Doubanjiang":           {},
		"Tempeh":                {},
		"Chilli":                {},
		"Chili Powder":          {},
		"Chilli Powder":         {},
		"Red Chilli":            {},
		"Red Chilli Powder":     {},
		"Red Chilli Flakes":     {},
		"Harissa Spice":         {},
		"Mustard":               {},
		"Mustard Powder":        {},
		"Mustard Seeds":         {},
		"English Mustard":       {},
		"Dijon Mustard":         {},
		"Eggplant":              {},
		"Egg Plants":            {},
		"Swiss Cheese":          {},
		"Blue Cheese":           {},
	}

	Caffeine := map[string]struct{}{
		"Cacao":                {},
		"Cocoa":                {},
		"Chocolate Chips":      {},
		"Dark Chocolate":       {},
		"Milk Chocolate":       {},
		"White Chocolate":      {},
		"Plain Chocolate":      {},
		"Dark Chocolate Chips": {},
		"Tea":                  {},
		"Coffee":               {},
		"Espresso":             {},
		"Mars Bar":             {},
	}

	Alcohol := map[string]struct{}{
		"Brandy":           {},
		"Dry White Wine":   {},
		"Red Wine":         {},
		"Sake":             {},
		"White Wine":       {},
		"Red Wine Vinegar": {},
		"Rice wine":        {},
		"Cooking wine":     {},
		"Dark Rum":         {},
		"Light Rum":        {},
		"Rum":              {},
		"Sherry":           {},
		"Grand Marnier":    {},
		"Cider":            {},
		"Stout":            {},
	}

	Sulphites := map[string]struct{}{
		"Baking Powder":      {},
		"Balsamic Vinegar":   {},
		"Brandy":             {},
		"Breadcrumbs":        {},
		"Christmas Pudding":  {},
		"Condensed Milk":     {},
		"Digestive Biscuits": {},
		"Dried Apricots":     {},
		"Dried Fruit":        {},
		"Dry White Wine":     {},
		"Glace Cherry":       {},
		"Mars Bar":           {},
		"Pickle Juice":       {},
		"Red Wine":           {},
		"Red Wine Vinegar":   {},
		"Sherry":             {},
		"Stout":              {},
		"Sultanas":           {},
		"Tinned Tomatos":     {},
		"Treacle":            {},
		"White Wine":         {},
		"White Wine Vinegar": {},
		"Wine":               {},
	}

	Salicylates := map[string]struct{}{
		"Avocado":             {},
		"Apple Cider Vinegar": {},
		"Asparagus":           {},
		"Aubergine":           {},
		"Basil":               {},
		"Bay Leaf":            {},
		"Bay Leaves":          {},
		"Balsamic Vinegar":    {},
		"Black Pepper":        {},
		"Broccoli":            {},
		"Capsicum":            {},
		"Chilli":              {},
		"Chili Powder":        {},
		"Chilli Powder":       {},
		"Chocolate Chips":     {},
		"Cocoa":               {},
		"Cucumber":            {},
		"Currants":            {},
		"Dried Apricots":      {},
		"Fennel":              {},
		"Ginger":              {},
		"Green Olives":        {},
		"Honey":               {},
		"Mint":                {},
		"Mushrooms":           {},
		"Olive Oil":           {},
		"Oranges":             {},
		"Oregano":             {},
		"Paprika":             {},
		"Peaches":             {},
		"Plum Tomatoes":       {},
		"Raspberries":         {},
		"Red Chilli":          {},
		"Red Chilli Flakes":   {},
		"Red Pepper":          {},
		"Spinach":             {},
		"Strawberries":        {},
		"Tomato Ketchup":      {},
		"Tomato Puree":        {},
		"Tomatoes":            {},
		"Vanilla":             {},
		"Vinegar":             {},
	}

	intolerancesMaps := map[string]map[string]struct{}{
		"Lactose":     Lactose,
		"Gluten":      Gluten,
		"Histamine":   Histamine,
		"Caffeine":    Caffeine,
		"Alcohol":     Alcohol,
		"Sulphites":   Sulphites,
		"Salicylates": Salicylates,
	}

	found := make(map[string]bool)
	var result []string

	for _, ingredient := range ingredients {
		for allergenName, allergenMap := range intolerancesMaps {
			if _, exists := allergenMap[ingredient]; exists {
				if !found[allergenName] {
					result = append(result, allergenName)
					found[allergenName] = true
				}
			}
		}
	}

	return result

}
