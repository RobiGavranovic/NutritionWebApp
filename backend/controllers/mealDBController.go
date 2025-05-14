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
	if origin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name parameter is required"})
		return
	}

	// Search All Meals With Origin
	apiURL := fmt.Sprintf("https://www.themealdb.com/api/json/v1/1/filter.php?s=%s", name)

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
