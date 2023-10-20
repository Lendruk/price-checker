package routes

import (
	"price-tracker/models"
	"price-tracker/parsers"

	"github.com/gin-gonic/gin"
)

func UpdateProducts(context *gin.Context) {

	// For now all products
	// TODO implement watch flags

	vendorEntires := models.GetAllVendorEntries()

	parsers.UpdateProducts(vendorEntires)
}
