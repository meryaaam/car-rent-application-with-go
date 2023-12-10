package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"rentCarTest/controller"
	"rentCarTest/models"
	"strings"
	"testing"
)

func TestCars(t *testing.T) {
	// Set up a Gin router
	r := gin.Default()

	// Use an in-memory database for testing
	db, err := setupTestDatabase()
	if err != nil {
		t.Fatal("Failed to set up the test database")
	}

	// Inject the test database into the Gin context
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	// Add the return car endpoint to the router
	r.POST("/cars/:registration/returns", controller.ReturnCar)

	// Test case: Car exists, is marked as rented, and valid JSON input
	t.Run("CarReturnedSuccessfully", func(t *testing.T) {
		// Insert a test car into the database
		testCar := &models.Car{
			Model:        "TestModel",
			Registration: "ABC123",
			Mileage:      5000,
			Status:       "rented",
		}
		db.Create(&testCar)

		// JSON input with kilometers driven
		jsonInput := `{"kilometers_driven": 100}`

		// Create a request to the /cars/:registration/returns endpoint
		req, err := http.NewRequest("POST", "/cars/ABC123/returns", strings.NewReader(jsonInput))
		assert.NoError(t, err)

		// Set the content type to JSON
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder to record the response
		w := httptest.NewRecorder()

		// Perform the request
		r.ServeHTTP(w, req)

		// Check the HTTP status code
		assert.Equal(t, http.StatusOK, w.Code)

		// Check the response body
		assert.JSONEq(t, `{"message": "Car returned successfully"}`, w.Body.String())

		// Check the database to ensure the car status and mileage are updated
		var updatedCar models.Car
		db.Where("registration = ?", "ABC123").First(&updatedCar)
		assert.Equal(t, "available", updatedCar.Status)
		assert.Equal(t, int16(5100), updatedCar.Mileage)
	})

	// Add more test cases as needed
}

// Helper function to set up an in-memory database for testing
func setupTestDatabase() (*gorm.DB, error) {
	dsn := "root:@tcp(localhost:3306)/cars_db?charset=utf8mb4&parseTime=True&loc=Local"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
