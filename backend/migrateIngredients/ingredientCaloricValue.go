package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/RobiGavranovic/NutritionWebApp/backend/initializers"
	"github.com/RobiGavranovic/NutritionWebApp/backend/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	// Migrate Ingredients Table
	err := initializers.DB.AutoMigrate(&models.Ingredient{})
	if err != nil {
		log.Fatalf("Failed to run AutoMigrate: %v", err)
	}

	// Open Ingredients.json
	file, err := os.Open("migrateIngredients/ingredients.json")
	if err != nil {
		log.Fatalf("Failed to open JSON file: %v", err)
	}
	defer file.Close()

	// Define Struct Matching The JSON Structure
	type RawIngredient struct {
		ID              string `json:"idIngredient"`
		Name            string `json:"strIngredient"`
		CaloriesPer100g int    `json:"calories100g"`
	}

	var ingredients []RawIngredient

	// Decode JSON Into Struct
	if err := json.NewDecoder(file).Decode(&ingredients); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	// Insert Each Ingredient Into Database
	for _, ing := range ingredients {
		id, err := strconv.Atoi(ing.ID)
		if err != nil {
			log.Printf("Skipping invalid ID %s: %v", ing.ID, err)
			continue
		}

		ingredient := models.Ingredient{
			ID:          uint(id),
			Name:        ing.Name,
			KcalPer100g: float64(ing.CaloriesPer100g),
		}

		if err := initializers.DB.Create(&ingredient).Error; err != nil {
			log.Printf("Failed to insert ingredient %s: %v", ing.Name, err)
		}
	}

	log.Println("Successfully populated ingredients table.")
}
