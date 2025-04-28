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
		AllowOrigins:     []string{"http://localhost:3000"}, // your frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	router.GET("/getRandomMeals/:numOfMeals", controllers.GetNRandomMeals)
	router.POST("/register", controllers.RegisterUser)

	router.Run(":" + os.Getenv("PORT"))
}
