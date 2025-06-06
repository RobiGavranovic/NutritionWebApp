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

	// Check If User With Email Already Exists
	var existingUser models.User
	if err := initializers.DB.Where("email = ?", userInfo.Email).First(&existingUser).Error; err == nil {
		// User found
		c.JSON(http.StatusConflict, gin.H{
			"error": "User already exists, please login instead",
		})
		return
	}

	// Check if User With Username Already Exists
	var existingByUsername models.User
	if err := initializers.DB.Where("username = ?", incomingData.Username).First(&existingByUsername).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken, please choose another"})
		return
	}

	// Create User In DB
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.Status(400)
		return
	}

	// Find User In DB - We Need ID
	var newUser models.User
	if err := initializers.DB.Where("email = ?", userInfo.Email).First(&newUser).Error; err != nil {
		// User Not Found
		c.JSON(http.StatusConflict, gin.H{
			"error": "Failed to create a user, please try again",
		})
		return
	}

	// Create JWT Token
	token, err := GenerateJWT(newUser.ID, newUser.Email)
	if err != nil || token == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to provide session"})
		return
	}

	// Use This In Prod, We Need HTTP For LocalHost Testing, So Secure = false
	//c.SetCookie("session_token", token, 3600*2, "/", "", true, true)
	c.SetCookie("session_token", token, 3600*2, "/", "localhost", false, true)

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
	if err := initializers.DB.Unscoped().Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user account"})
		return
	}

	// Clear JWT Token Cookie
	c.SetCookie("Authorization", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "User account and data deleted successfully"})
}
