package routes

import (
	"price-tracker/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterUser(context *gin.Context) {
	newUser := models.NewUser()
	context.JSON(201, newUser)
}

func RegisterWebhookUser(context *gin.Context) {
	type RequestBody struct {
		Hook string `json:"hook"`
	}
	var body RequestBody

	context.BindJSON(&body)
	newUser := models.NewUser()
	err := models.AddUserToWebHook(body.Hook, int(newUser.Id))

	if err != nil {
		context.JSON(500, gin.H{"message": "Error creating webhook user"})
	} else {
		context.JSON(201, newUser)
	}
}

type AddToWatchListRequestBody struct {
	Product int `json:"product"`
}

func GetUser(context *gin.Context) {
	userId, success := context.Params.Get("id")

	if success {
		parsedUserId, _ := strconv.ParseInt(userId, 10, 64)
		user, err := models.GetUserById(parsedUserId)

		if err != nil {
			context.JSON(404, gin.H{"Message": "User not found"})
		} else {
			context.JSON(200, user)
		}
	} else {
		context.JSON(404, gin.H{"Message": "User not found"})
	}
}

func AddProductToWatchlist(context *gin.Context) {
	var body AddToWatchListRequestBody
	rawUserId, _ := context.Params.Get("id")
	context.BindJSON(&body)
	parsedUserId, convErr := strconv.ParseInt(rawUserId, 10, 32)

	if convErr != nil {
		context.JSON(400, gin.H{"message": "No valid user id sent in request"})
	} else if models.IsProductInWatchlist(int(parsedUserId), body.Product) {
		context.JSON(400, gin.H{"message": "Product is already in the watchlist of the user"})
	} else {
		err := models.AddProductToWatchlist(int(parsedUserId), body.Product)

		if err != nil {
			context.JSON(500, gin.H{"message": "There was an error inserting the product into the watchlist"})
		} else {
			context.JSON(200, gin.H{"message": "Product inserted into watchlist successfully"})
		}
	}
}

func RemoveProductFromWatchlist(context *gin.Context) {
	user, _ := context.Params.Get("id")
	productSKU, _ := context.Params.Get("product")

	parsedUserId, userConvErr := strconv.ParseInt(user, 10, 32)

	if userConvErr != nil {
		context.JSON(400, gin.H{"message": "No parameters sent in request"})
	}

	err := models.RemoveProductFromWatchlist(int(parsedUserId), productSKU)

	if err != nil {
		context.JSON(500, gin.H{"message": "There was an error removing the product from the watchlist"})
	}

	context.JSON(200, gin.H{"message": "Removed product from watchlist successfully"})
}

func GetUserWatchList(context *gin.Context) {
	user, _ := context.Params.Get("id")

	parsedUserId, userConvErr := strconv.ParseInt(user, 10, 32)

	if userConvErr != nil {
		context.JSON(400, gin.H{"message": "Invalid user id sent"})
	}

	products := models.FetchUserWatchList(int(parsedUserId))

	context.JSON(200, products)
}
