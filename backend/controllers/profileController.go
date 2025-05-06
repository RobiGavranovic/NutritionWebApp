package controllers

import (
	"encoding/json"
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

func UpdateAllergens(c *gin.Context) {
	// Grab Sent Data
	var req models.UpdateAllergens
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Verify Google Token
	googleReq, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	googleReq.Header.Set("Authorization", "Bearer "+req.TokenResponse.AccessToken)

	// Fetch User Data
	client := &http.Client{}
	googleRes, err := client.Do(googleReq)
	if err != nil || googleRes.StatusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired Google token"})
		return
	}
	defer googleRes.Body.Close()

	var userInfo models.GoogleUser
	if err := json.NewDecoder(googleRes.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	// Find User By Email
	var user models.User
	if err := initializers.DB.Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
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

	// Verify Google Token
	googleReq, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	googleReq.Header.Set("Authorization", "Bearer "+req.TokenResponse.AccessToken)

	// Fetch User Data
	client := &http.Client{}
	googleRes, err := client.Do(googleReq)
	if err != nil || googleRes.StatusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired Google token"})
		return
	}
	defer googleRes.Body.Close()

	var userInfo models.GoogleUser
	if err := json.NewDecoder(googleRes.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	// Find User By Email
	var user models.User
	if err := initializers.DB.Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
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

	// Verify Google Token
	googleReq, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	googleReq.Header.Set("Authorization", "Bearer "+req.TokenResponse.AccessToken)

	// Fetch User Data
	client := &http.Client{}
	googleRes, err := client.Do(googleReq)
	if err != nil || googleRes.StatusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired Google token"})
		return
	}
	defer googleRes.Body.Close()

	var userInfo models.GoogleUser
	if err := json.NewDecoder(googleRes.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	// Find User By Email
	var user models.User
	if err := initializers.DB.Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update Username
	if err := initializers.DB.Model(&user).Update("Username", req.Username).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Username updated",
		"user": user})
}
