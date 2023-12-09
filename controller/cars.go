package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Registration already exists"})
		return
	}

	NewCar := models.Car{Model: input.Model, Registration: input.Registration, Mileage: input.Mileage}
	models.DB.Create(&NewCar)
	UUID := uuid.MustParse(NewCar.ID.String())
	fmt.Println("Car ID:", NewCar.ID)

	c.JSON(http.StatusOK, gin.H{"id": UUID})

}
