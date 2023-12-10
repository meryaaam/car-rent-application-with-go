package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"rentCarTest/models"
)

func GetALl(c *gin.Context) {
	var cars []models.Car

	if err := models.DB.Find(&cars).Error; err != nil {
		fmt.Println("Error fetching cars:", err)
		fmt.Printf("Context: %+v\n", c.Request)
		fmt.Printf("DB Instance: %+v\n", models.DB)
		fmt.Printf("Error Type: %T\n", err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNoContent, gin.H{
				"message": "No records found",
			})
			return
		}

		// Return an error response to the client
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	// Handle the retrieved data as needed

	c.JSON(http.StatusOK, gin.H{
		"data": cars,
	})

}

type CreateCarInput struct {
	Model        string `json:"model" binding:"required"`
	Registration string `json:"registration" binding:"required"`
	Mileage      int    `json:"mileage" binding:"required"`
}

func registrationExists(registration string) bool {
	var existingCar models.Car
	result := models.DB.Where("registration = ?", registration).First(&existingCar)
	return result.Error == nil
}
func CreateCar(c *gin.Context) {
	var input CreateCarInput
	// Handle the case where the input not valid json input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Handle the case where registration already exists
	if registrationExists(input.Registration) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "the registration number already exists"})
		return
	}

	NewCar := models.Car{Model: input.Model, Registration: input.Registration, Mileage: input.Mileage}
	models.DB.Create(&NewCar)

	c.JSON(http.StatusOK, gin.H{"id": NewCar.ID})

}

func RentCar(c *gin.Context) {

	registration := c.Param("registration")

	// Handle the case where registration already exists
	var car models.Car
	result := models.DB.Where("registration = ?", registration).First(&car)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Car not found"})
		return
	}

	// Check if the car is already rented
	if car.Status == "rented" {
		c.JSON(http.StatusOK, gin.H{"message": "Car is already rented"})
		return
	}

	// Mark the car as rented
	car.Status = "rented"
	models.DB.Save(&car)

	c.JSON(http.StatusOK, gin.H{"message": "Car rented successfully"})
}

func ReturnCar(c *gin.Context) {
	registration := c.Param("registration")

	var Car models.Car
	// Handle the case where registration already exists
	result := models.DB.Where("registration = ?", registration).First(&Car)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
		return
	}
	// Handle error if the car is not marked as rented
	if Car.Status != "rented" {
		c.JSON(http.StatusConflict, gin.H{"error": "Car is not marked as rented"})
		return
	}

	var returnRequest struct {
		KilometersDriven int `json:"kilometers_driven"`
	}

	if err := c.ShouldBindJSON(&returnRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the mileage and status of the car
	Car.Mileage += returnRequest.KilometersDriven
	Car.Status = "available"
	models.DB.Save(&Car)

	c.JSON(http.StatusOK, gin.H{"message": "Car returned successfully"})
}
