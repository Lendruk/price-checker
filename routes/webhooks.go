package routes

import (
	"price-tracker/models"

	"github.com/gin-gonic/gin"
)

func RegisterWebhook(context *gin.Context) {
	type RequestBody struct {
		Hook string `json:"hook"`
	}

	var body RequestBody

	context.BindJSON(&body)

	err := models.RegisterWebhook(body.Hook)

	if err != nil {
		context.JSON(500, gin.H{"message": "Error creating webhook"})
	} else {
		context.JSON(201, gin.H{"message": "Webhook created successfully"})
	}
}
