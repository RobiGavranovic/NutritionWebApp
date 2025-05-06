package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/RobiGavranovic/NutritionWebApp/backend/initializers"
	"github.com/RobiGavranovic/NutritionWebApp/backend/models"
	"github.com/gin-gonic/gin"
)

func LoginUser(c *gin.Context) {
	// Get POST Data
	var googleOAuthData models.TokenResponse
	if err := c.ShouldBindJSON(&googleOAuthData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Verify Google Token
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate user"})
		return
	}
	req.Header.Set("Authorization", "Bearer "+googleOAuthData.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate user"})
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate user"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate user"})
		return
	}

	var userInfo models.GoogleUser
	if err := json.Unmarshal(bodyBytes, &userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate user"})
		return
	}

	// Find User In DB
	var user models.User
	if err := initializers.DB.Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
		// User Not Found
		c.JSON(http.StatusConflict, gin.H{
			"error": "Account with this email address does not exist, please register instead",
		})
		return
	}

	// Create JWT Token
	token, err := GenerateJWT(user.ID, user.Email)
	if err != nil || token == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to provide session"})
		return
	}

	// Use This In Prod, We Need HTTP For LocalHost Testing, So Secure = false
	//c.SetCookie("session_token", token, 3600*2, "/", "", true, true)
	c.SetCookie("session_token", token, 3600*2, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})

}

func LogoutUser(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
