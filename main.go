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

	router.POST("/products", routes.RegisterProduct)
	router.POST("/products/update", routes.UpdateProducts)
	router.GET("/products", routes.FetchProducts)
	router.Run()
}
