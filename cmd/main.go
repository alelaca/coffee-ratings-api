package main

import (
	"github.com/alelaca/coffee-ratings-api/cmd/controller"
	"github.com/alelaca/coffee-ratings-api/cmd/dao"
	"github.com/alelaca/coffee-ratings-api/cmd/router"
	"github.com/alelaca/coffee-ratings-api/cmd/services"
)

func main() {
	ratingsDAO := dao.InitializeLocalRepository()
	ratingsService := services.InitializeRatingsService(ratingsDAO)
	ratingsController := controller.InitializeRatingsController(ratingsService)

	router.Start(ratingsController)
}
