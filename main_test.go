package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"rentCarTest/controller"
	"rentCarTest/models"
	"testing"
)

func setupTestDatabase() (*gorm.DB, error) {
	dsn := "root:@tcp(localhost:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func TestConnection(t *testing.T) {
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

}

func TestGetCars(t *testing.T) {
	models.InitDB()

	// Create a fake Gin context for testing
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/cars", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	controller.GetALl(c)

	// Check the response status
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, w.Code)
	}
}

func TestCreateCar(t *testing.T) {
	models.InitDB()
	// Create a JSON payload for the test request
	newCarInput := controller.CreateCarInput{
		Model:        "Toyota GR 3",
		Registration: "BG0027",
		Mileage:      1005,
	}
	jsonData, err := json.Marshal(newCarInput)
	// Create a new HTTP request
	req, err := http.NewRequest("POST", "/cars", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	w := httptest.NewRecorder()

	// Create a Gin router and set up the route
	r := gin.Default()
	r.POST("/cars", controller.CreateCar)

	// Serve the HTTP request to the response recorder
	r.ServeHTTP(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response JSON to check the content
	// Assuming the response is in the format: {"id": "generated_id"}
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the "id" field is present in the response
	assert.Contains(t, response, "id")
	assert.NotEmpty(t, response["id"])

}
func TestRentCar(t *testing.T) {

	models.InitDB()
	r := gin.Default()

	requestPayload := gin.H{
		"registration": "BG0027",
	}
	jsonPayload, _ := json.Marshal(requestPayload)

	// Create a request with the JSON payload
	req, err := http.NewRequest("POST", "/cars/BG0027/rentals", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	r.POST("/cars/:registration/rentals", controller.RentCar)

	r.ServeHTTP(w, req)

	// Check the response status
	assert.Equal(t, http.StatusOK, w.Code)

	// Check if the car's status has been updated to "rented" in the database
	var updatedCar models.Car
	models.DB.First(&updatedCar, "registration = ?", "BG0027")
	assert.Equal(t, models.Rented, updatedCar.Status)

	// Add more assertions as needed
}
