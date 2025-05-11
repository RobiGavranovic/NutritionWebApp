package main

import (
	"log"

	"github.com/RobiGavranovic/NutritionWebApp/backend/initializers"
	"github.com/RobiGavranovic/NutritionWebApp/backend/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	err := initializers.DB.AutoMigrate(&models.Consumption{})
	if err != nil {
		log.Fatalf("Failed to run AutoMigrate: %v", err)
	} else {
		log.Fatalf("Successfully Migrated Consumption Table")
	}
}
