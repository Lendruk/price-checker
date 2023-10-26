package main

import (
	"fmt"
	"price-tracker/cron"
	"price-tracker/routes"

	"github.com/gin-gonic/gin"
)

const port = ":8080"

func main() {
	fmt.Println("Starting Price Tracker!")
	router := gin.Default()

	router.POST("/products/search", routes.SearchProducts)
	router.POST("/products/update", routes.UpdateProducts)
	router.GET("/products", routes.FetchProducts)

	router.POST("/users", routes.RegisterUser)
	router.PUT("/users/:id/watchlist", routes.AddProductToWatchlist)
	router.DELETE("/users/:id/watchlist/:product", routes.RemoveProductFromWatchlist)
	router.GET("/users/:id", routes.GetUser)

	router.POST("/webhooks", routes.RegisterWebhook)
	router.POST("/webhooks/users", routes.RegisterWebhookUser)
	router.POST("/webhooks/products", routes.RegisterWebhookProduct)
	cron.RegisterProductUpdateCronJob()
	router.Run()

}
