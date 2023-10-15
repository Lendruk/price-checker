package main

import (
	"fmt"
	"net/http"
	"price-tracker/routes"

	"github.com/gin-gonic/gin"
)

const port = ":8080"

func home(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello World!")
}

func main() {
	fmt.Println("Starting Price Tracker!")
	router := gin.Default()

	router.POST("/product", routes.RegisterProduct)
	router.Run()
}
