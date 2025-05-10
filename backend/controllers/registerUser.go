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
