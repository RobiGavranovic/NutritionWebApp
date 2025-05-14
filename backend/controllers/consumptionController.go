package controllers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/RobiGavranovic/NutritionWebApp/backend/initializers"
	"github.com/RobiGavranovic/NutritionWebApp/backend/models"
	"github.com/gin-gonic/gin"
)

func Consume(c *gin.Context) {
	// Grab Sent Data
	var req models.ConsumptionRequest
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

	// Get Ingredient From DB
	var ingredient models.Ingredient
	if err := initializers.DB.Where("name = ?", req.Ingredient).First(&ingredient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
		return
	}

	calories := (req.Weight / 100.0) * ingredient.KcalPer100g

	// Create consumption record
	consumption := models.Consumption{
		UserID:       user.ID,
		IngredientID: ingredient.ID,
		Ingredient:   ingredient.Name,
		Weight:       req.Weight,
		Calories:     calories,
	}

	if err := initializers.DB.Create(&consumption).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record consumption"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Consumption recorded",
		"caloriesConsumed": calories,
		"id":               consumption.ID,
	})
}

func GetAllIngredients(c *gin.Context) {
	// Grab Sent Data
	var req models.IngredientsNamesRequest
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

	// Get All Ingredients Names
	var ingredients []models.Ingredient
	if err := initializers.DB.Select("name").Find(&ingredients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load ingredients"})
		return
	}

	// Format For FE
	var response []map[string]string
	for _, ing := range ingredients {
		response = append(response, map[string]string{"name": ing.Name})
	}

	c.JSON(http.StatusOK, response)
}

func DeleteConsumed(c *gin.Context) {
	// Get ID From URL
	IDstr := c.Param("id")
	ID, err := strconv.Atoi(IDstr)
	if err != nil || ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumption ID"})
		return
	}

	// Get User ID from context
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User context not found"})
		return
	}
	userID := userIDRaw.(uint)

	// Find Record By ID And UserID
	var entry models.Consumption
	if err := initializers.DB.Where("id = ? AND user_id = ?", ID, userID).First(&entry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Consumption entry not found or access denied"})
		return
	}

	// Delete Record
	if err := initializers.DB.Delete(&entry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete consumption entry"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Consumption entry deleted successfully"})
}

func GetTodaysConsumption(c *gin.Context) {
	// Get User ID from context
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User context not found"})
		return
	}
	userID := userIDRaw.(uint)

	// Get Current Local Date
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	// Query Today's Consumption
	var entries []models.Consumption
	if err := initializers.DB.Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, todayStart, todayEnd).
		Order("created_at DESC").
		Find(&entries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch consumption history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"consumptions": entries,
	})
}

func GetConsumptionStatistics(c *gin.Context) {
	// Get User ID From Context
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User context not found"})
		return
	}
	userID := userIDRaw.(uint)

	// Parse Range From Query String
	rangeStr := c.Query("range")
	daysRange, err := strconv.Atoi(rangeStr)
	if err != nil || daysRange <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid range"})
		return
	}

	// Define Date Boundaries (Yesterday To (Yesterday - Range + 1))
	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)
	startDate := yesterday.AddDate(0, 0, -daysRange+1)

	// Get User Info
	var user models.User
	if err := initializers.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get All Consumption Entries In The Range
	var consumptions []models.Consumption
	if err := initializers.DB.
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, yesterday).
		Find(&consumptions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch consumptions"})
		return
	}

	// Group Calories By Date
	dailyTotals := make(map[string]float64)
	for _, con := range consumptions {
		dateKey := con.CreatedAt.Format("2006-01-02")
		dailyTotals[dateKey] += con.Calories
	}

	// Evaluate Success/Failure/No Input For Each Day
	success, fail, noInput := 0, 0, 0
	for i := 0; i < daysRange; i++ {
		date := startDate.AddDate(0, 0, i).Format("2006-01-02")
		calories, exists := dailyTotals[date]

		if !exists {
			noInput++
			continue
		}

		switch user.DailyGoalType {
		case "gain":
			if calories >= float64(user.DailyCalorieGoal) {
				success++
			} else {
				fail++
			}
		case "lose":
			if calories <= float64(user.DailyCalorieGoal) {
				success++
			} else {
				fail++
			}
		case "maintain":
			diff := math.Abs(calories - float64(user.DailyCalorieGoal))
			if diff <= 50.0 {
				success++
			} else {
				fail++
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Input your goal calories and goal type."})
			return
		}
	}

	history := GetConsumptionIntakeHistory(daysRange, startDate, yesterday, consumptions)

	c.JSON(http.StatusOK, gin.H{
		"success":       success,
		"fail":          fail,
		"noInput":       noInput,
		"fromDate":      startDate.Format("2006-01-02"),
		"toDate":        yesterday.Format("2006-01-02"),
		"intakeHistory": history,
	})
}

func GetConsumptionIntakeHistory(daysRange int, startDate, endDate time.Time, consumptions []models.Consumption) []map[string]interface{} {
	// Group Entries By Day
	dailyBuckets := make(map[string][]float64)
	for _, con := range consumptions {
		key := con.CreatedAt.Format("2006-01-02")
		dailyBuckets[key] = append(dailyBuckets[key], con.Calories)
	}

	// Helper To Calculate Average Ignoring Empty Days
	average := func(values []float64) float64 {
		if len(values) == 0 {
			return 0
		}
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		return sum / float64(len(values))
	}

	results := []map[string]interface{}{}

	switch {
	case daysRange <= 7:
		// Daily
		for i := 0; i < daysRange; i++ {
			currDate := startDate.AddDate(0, 0, i)
			label := currDate.Format("Jan 2")
			bucket := dailyBuckets[currDate.Format("2006-01-02")]
			avg := average(bucket)
			results = append(results, map[string]interface{}{
				"label":       label,
				"avgCalories": math.Round(avg),
			})
		}

	case daysRange <= 60:
		// Weekly
		weekEnd := endDate
		for weekEnd.After(startDate) {
			weekStart := weekEnd.AddDate(0, 0, -6)
			if weekStart.Before(startDate) {
				weekStart = startDate
			}

			label := fmt.Sprintf("%s - %s", weekStart.Format("Jan 2"), weekEnd.Format("Jan 2"))

			// Sum Calories
			weekValues := []float64{}
			for d := 0; !weekStart.After(weekEnd); d++ {
				key := weekStart.Format("2006-01-02")
				weekValues = append(weekValues, dailyBuckets[key]...)
				weekStart = weekStart.AddDate(0, 0, 1)
			}

			results = append([]map[string]interface{}{
				{
					"label":       label,
					"avgCalories": math.Round(average(weekValues)),
				},
			}, results...)

			// Shift To Previous Week
			weekEnd = weekEnd.AddDate(0, 0, -7)
		}

	default:
		// Monthly
		monthMap := make(map[string][]float64)
		monthOrder := []string{}

		for i := 0; i < daysRange; i++ {
			currDate := startDate.AddDate(0, 0, i)
			monthKey := currDate.Format("Jan 2006")
			dateKey := currDate.Format("2006-01-02")

			if _, exists := monthMap[monthKey]; !exists {
				monthOrder = append(monthOrder, monthKey)
			}
			monthMap[monthKey] = append(monthMap[monthKey], dailyBuckets[dateKey]...)
		}

		for _, month := range monthOrder {
			values := monthMap[month]
			results = append(results, map[string]interface{}{
				"label":       month,
				"avgCalories": math.Round(average(values)),
			})
		}

	}

	return results
}
