package main

import (
	"github.com/RobiGavranovic/NutritionWebApp/backend/controllers"
	"github.com/RobiGavranovic/NutritionWebApp/backend/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	router := gin.Default()

	router.GET("/getRandomMeals/:numOfMeals", controllers.GetNRandomMeals)

	router.Run()
}
