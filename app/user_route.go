package app

import (
	"kautsar/travel-app-api/controller"
	"kautsar/travel-app-api/repository"
	"kautsar/travel-app-api/service"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUserRoutes(router *httprouter.Router, db *mongo.Database, validate *validator.Validate) {
	// user init
	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db, validate)
	userController := controller.NewUserController(userService)
	// user
	router.GET("/api/v1/user/:region_name/destinations", userController.FindDestinationByRegion)
	router.GET("/api/v1/user/:region_name/destinations/:destination_id", userController.FindOneDestinationByRegion)
	router.POST("/api/v1/login", userController.Login)
}
