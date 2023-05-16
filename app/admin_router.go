package app

import (
	"kautsar/travel-app-api/controller"
	"kautsar/travel-app-api/repository"
	"kautsar/travel-app-api/service"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterAdminRoutes(router *httprouter.Router, db *mongo.Database, validate *validator.Validate) {

	// operator init
	operatorRepository := repository.NewOperatorRepository()
	operatorService := service.NewOperatorService(operatorRepository, db, validate)
	operatorController := controller.NewOperatorController(operatorService)
	//account

	router.GET("/api/v1/account", operatorController.FindAll)
	router.POST("/api/v1/account", operatorController.Create)
	router.PATCH("/api/v1/account/:accountId", operatorController.ResetPassword)
	router.DELETE("/api/v1/account/:accountId", operatorController.Destroy)
}
