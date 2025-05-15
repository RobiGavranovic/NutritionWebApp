package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/RobiGavranovic/NutritionWebApp/backend/initializers"
	"github.com/RobiGavranovic/NutritionWebApp/backend/models"
	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	// Get POST Data
	var incomingData models.RegisterUser
	if err := c.ShouldBindJSON(&incomingData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get Google User Data
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+incomingData.TokenResponse.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return
	}

	var userInfo models.GoogleUser
	if err := json.Unmarshal(bodyBytes, &userInfo); err != nil {
		return
	}

	// Create User From Data
	user := models.User{
		GoogleID:     incomingData.TokenResponse.TokenType,
		Email:        userInfo.Email,
		Username:     incomingData.Username,
		Gender:       incomingData.Gender,
		Allergens:    incomingData.Allergens,
		Intolerances: incomingData.Intolerances,
	}

	// Check If User Already Exists
	var existingUser models.User
	if err := initializers.DB.Where("email = ?", userInfo.Email).First(&existingUser).Error; err == nil {
		// User found
		c.JSON(http.StatusConflict, gin.H{
			"error": "User already exists, please login instead",
		})
		return
	}

	// Create User In DB
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.Status(400)
		return
	}

	// Return It
	c.JSON(http.StatusOK, gin.H{
		"message":  "Received successfully",
		"token":    incomingData,
		"userInfo": userInfo,
	})
}

func DeleteUser(c *gin.Context) {
	// Get User Email
	email, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user context"})
		return
	}

	// Find User By Email
	var user models.User
	if err := initializers.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete User Consumptions
	if err := initializers.DB.Where("user_id = ?", user.ID).Delete(&models.Consumption{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user's consumption data"})
		return
	}

	// Delete User Account
	if err := initializers.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user account"})
		return
	}

	// Clear JWT Token Cookie
	c.SetCookie("Authorization", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "User account and data deleted successfully"})
}
