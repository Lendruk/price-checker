package routes

import (
	"net/http"
	"price-tracker/models"
	"price-tracker/parsers"

	"github.com/gin-gonic/gin"
)

func SearchProducts(context *gin.Context) {
	type RequestBody struct {
		Product string `json:"product"`
	}
	var body RequestBody

	context.BindJSON(&body)

	parsers.RegisterProductByName(body.Product)

	context.JSON(http.StatusOK, gin.H{"product": body.Product})
}

func RegisterWebhookProduct(context *gin.Context) {
	type RequestBody struct {
		User int    `json:"user"`
		Url  string `json:"url"`
		Hook string `json:"hook"`
	}
	var body RequestBody

	context.BindJSON(&body)

	// Search and register products
	product, productError := parsers.RegisterProductFromUrl(body.Url)

	if productError != nil {
		context.JSON(500, gin.H{"message": "Error while registering the product"})
	}

	// Add new product to user watchlist
	if models.IsProductInWatchlist(body.User, product.Id) == false {
		models.AddProductToWatchlist(body.User, product.Id)
	}

	// Add user to webhook user list if not there yet
	models.AddUserToWebHook(body.Hook, body.User)
}
