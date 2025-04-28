package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

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

	var meals models.Meals
	if err := json.Unmarshal(body, &meals); err != nil {
		log.Println("Error decoding JSON:", err)
		return models.Meals{}, err
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
	for len(meals.Meals) < numOfMeals {
		mealWrapper, err := GetRandomMeal()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get random meal from MealDB"})
			return
		}

		meal := mealWrapper.Meals[0] //extract meal from meals wrapper

		id := meal.IDMeal

		// Check for duplicates
		if !seen[id] {
			meals.Meals = append(meals.Meals, meal)
			seen[id] = true
		}
	}

	c.JSON(http.StatusOK, meals)

}
