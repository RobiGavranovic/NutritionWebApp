package controllers

import (
	"net/http"

	"github.com/RobiGavranovic/NutritionWebApp/backend/initializers"
	"github.com/RobiGavranovic/NutritionWebApp/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func init() {

}

func GetProfileData(c *gin.Context) {
	// Get User ID
	rawID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check User ID
	userID, ok := rawID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

		
	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load profile"})
		return
	}

	// Return User Data
	var profile models.Profile
	profile.Username = user.Username
	profile.Allergens = user.Allergens
	profile.Intolerances = user.Intolerances
	profile.Age = user.Age
	profile.Height = user.Height
	profile.Weight = user.Weight
	profile.DailyCalorieGoal = user.DailyCalorieGoal

	c.JSON(http.StatusOK, profile)
}

func UpdateAllergens(c *gin.Context) {
	// Grab Sent Data
	var req models.UpdateAllergens
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

	// Update Allergens
	if err := initializers.DB.Model(&user).Update("allergens", pq.StringArray(req.Allergens)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update allergens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Allergens updated"})
}

func UpdateIntolarences(c *gin.Context) {
	// Grab Sent Data
	var req models.UpdateIntolarences
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

	// Update Intolerances
	if err := initializers.DB.Model(&user).Update("intolerances", req.Intolerances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update intolerances"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Intolerances updated",
		"user": user})
}

func UpdateUsername(c *gin.Context) {
	// Grab Sent Data
	var req models.UpdateUsername
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

	// Update Username
	if err := initializers.DB.Model(&user).Update("Username", req.Username).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Username updated"})
}

func UpdatePersonalInfo(c *gin.Context) {
	// Grab Sent Data
	var req models.UpdatePersonalInfo
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

	// Update Age, Height, Weight
	if err := initializers.DB.Model(&user).Updates(models.User{
		Age:    req.Age,
		Height: req.Height,
		Weight: req.Weight,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update personal info"})
		return
	}

	// Calculate Daily Calories
	calories := CalculateDailyCalories(req.Age, req.Height, req.Weight, user.Gender)

	// Return OK + Calories
	c.JSON(http.StatusOK, gin.H{"message": "Age, height and weight updated",
		"dailyCalorieGoal": calories})

}

func CalculateDailyCalories(age int, height int, weight int, gender string) (dailyCalories float64) {
	// Daily Calories Based On Gender
	if gender == "male" {
		dailyCalories = 10*float64(weight) + 6.25*float64(height) - 5*float64(age) + 5
	} else if gender == "female" {
		dailyCalories = 10*float64(weight) + 6.25*float64(height) - 5*float64(age) - 161
	} else {
		dailyCalories = 0 // Unknown Gender
	}

	return dailyCalories
}

func UpdateDailyCalorieGoal(c *gin.Context) {
	// Grab Sent Data
	var req models.UpdateDailyCalorieGoal
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

	// Update Daily Calorie Goal And Type
	if err := initializers.DB.Model(&user).Updates(models.User{
		DailyCalorieGoal: req.DailyCalorieGoal,
		DailyGoalType:    req.DailyGoalType,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update daily calorie goal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Daily calorie goal updated"})
}
