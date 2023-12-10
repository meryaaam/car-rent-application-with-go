package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"rentCarTest/controller"
	"rentCarTest/models"
)

func main() {

	router := gin.Default()

	models.Connection()
	models.InitDB()
	router.GET(
		"/cars",
		controller.GetALl,
	)
	router.POST(
		"/cars",
		controller.CreateCar,
	)

	router.POST(
		"/cars/:registration/rentals",
		controller.RentCar,
	)

	router.POST(
		"/cars/:registration/returns",
		controller.ReturnCar,
	)

	port := 9000
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}

}
