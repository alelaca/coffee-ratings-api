package router

import (
	"github.com/alelaca/coffee-ratings-api/cmd/controller"
	"github.com/gin-gonic/gin"
)

func Start(ratingsController *controller.Ratings) {
	router := gin.Default()

	api := router.Group("/api/coffee-ratings")

	api.GET("/greeting", controller.Greeting)

	api.POST("/ratings", ratingsController.CreateRating)
	api.GET("/ratings/coffee-types", ratingsController.GetCoffeeTypeList)
	api.GET("/ratings", ratingsController.GetRating)
	api.GET("/recommendation", ratingsController.GetRecommendation)

	router.Run(":9000")
}
