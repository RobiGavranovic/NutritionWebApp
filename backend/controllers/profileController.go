package controllers

import (
	"net/http"

	"github.com/RobiGavranovic/NutritionWebApp/backend/initializers"
	"github.com/RobiGavranovic/NutritionWebApp/backend/models"
	"github.com/gin-gonic/gin"
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

	// Get DB Data
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

	c.JSON(http.StatusOK, profile)
}
