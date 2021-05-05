package router

import (
	"github.com/alelaca/coffee-ratings-api/cmd/controller"
	"github.com/gin-gonic/gin"
)

func Start() {
	router := gin.Default()

	api := router.Group("/api/coffee-ratings")

	api.GET("/greeting", controller.Greeting)

	router.Run(":9000")
}
