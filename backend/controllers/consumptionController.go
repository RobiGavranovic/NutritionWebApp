package controllers

import (
	"net/http"

	"github.com/RobiGavranovic/NutritionWebApp/backend/initializers"
	"github.com/RobiGavranovic/NutritionWebApp/backend/models"
	"github.com/gin-gonic/gin"
)

func Consume(c *gin.Context) {
	// Grab Sent Data
	var req models.ConsumptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Find User By Email
	email, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user context"})
		return
	}

	var user models.User
	if err := initializers.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get Ingredient From DB
	var ingredient models.Ingredient
	if err := initializers.DB.Where("name = ?", req.Ingredient).First(&ingredient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
		return
	}

	calories := (req.Weight / 100.0) * ingredient.KcalPer100g

	// Create consumption record
	consumption := models.Consumption{
		UserID:       user.ID,
		IngredientID: ingredient.ID,
		Ingredient:   ingredient.Name,
		Weight:       req.Weight,
		Calories:     calories,
	}
	if err := initializers.DB.Create(&consumption).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record consumption"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Consumption recorded",
		"caloriesConsumed": calories,
	})
}

func GetAllIngredients(c *gin.Context) {
	// Grab Sent Data
	var req models.IngredientsNamesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Find User By Email
	email, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user context"})
		return
	}

	var user models.User
	if err := initializers.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get All Ingredients Names
	var ingredients []models.Ingredient
	if err := initializers.DB.Select("name").Find(&ingredients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load ingredients"})
		return
	}

	// Format For FE
	var response []map[string]string
	for _, ing := range ingredients {
		response = append(response, map[string]string{"name": ing.Name})
	}

	c.JSON(http.StatusOK, response)
}
