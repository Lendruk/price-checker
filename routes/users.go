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

type AddToWatchListRequestBody struct {
	User    int `json:"user"`
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
	context.BindJSON(&body)

	if models.IsProductInWatchlist(body.User, body.Product) {
		context.JSON(400, gin.H{"message": "Product is already in the watchlist of the user"})
	} else {
		err := models.AddProductToWatchlist(body.User, body.Product)

		if err != nil {
			context.JSON(500, gin.H{"message": "There was an error inserting the product into the watchlist"})
		} else {
			context.JSON(200, gin.H{"message": "Product inserted into watchlist successfully"})
		}
	}
}

func RemoveProductFromWatchlist(context *gin.Context) {

}
