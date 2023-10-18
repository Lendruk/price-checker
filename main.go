package main

import (
	"fmt"
	"price-tracker/routes"

	"github.com/gin-gonic/gin"
)

const port = ":8080"

func main() {
	fmt.Println("Starting Price Tracker!")
	router := gin.Default()

	router.POST("/product", routes.RegisterProduct)
	router.Run()
}
