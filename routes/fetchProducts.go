package routes

import (
	"net/http"
	"price-tracker/models"

	"github.com/gin-gonic/gin"
)

func FetchProducts(context *gin.Context) {
	products := models.GetProducts()
	context.JSON(http.StatusOK, products)
}
