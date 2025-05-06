package main

import (
	"os"

	"github.com/RobiGavranovic/NutritionWebApp/backend/controllers"
	"github.com/RobiGavranovic/NutritionWebApp/backend/initializers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	router := gin.Default()

	// Allow CORS for local development ONLY - DELETE THIS IN FOR PROD
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	router.GET("/getRandomMeals/:numOfMeals", controllers.GetNRandomMeals)
	router.POST("/register", controllers.RegisterUser)
	router.POST("/login", controllers.LoginUser)

	authorized := router.Group("/")
	authorized.Use(controllers.RequireAuth)
	authorized.GET("/profile", controllers.GetProfileData)

	authorized.POST("/logout", controllers.LogoutUser)
	authorized.PUT("/profile/updateAllergens", controllers.UpdateAllergens)
	authorized.PUT("/profile/updateIntolerances", controllers.UpdateIntolarences)
	authorized.PUT("/profile/updateUsername", controllers.UpdateUsername)

	router.Run(":" + os.Getenv("PORT"))
}
