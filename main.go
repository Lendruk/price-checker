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

	router.POST("/users", routes.RegisterUser)
	router.PUT("/users/:id/watchlist", routes.AddProductToWatchlist)
	router.DELETE("/users/:id/watchlist/:product", routes.RemoveProductFromWatchlist)
	router.GET("/users/:id", routes.GetUser)

	router.Run()
}
