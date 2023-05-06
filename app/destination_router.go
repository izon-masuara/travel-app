package app

import (
	"kautsar/travel-app-api/controller"
	"kautsar/travel-app-api/repository"
	"kautsar/travel-app-api/service"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterDestinationRoutes(router *httprouter.Router, db *mongo.Database, validate *validator.Validate) {
	// destination init
	destinationRepository := repository.NewDestinationRepository()
	destinationService := service.NewDestinationService(destinationRepository, db, validate)
	destinationController := controller.NewDestinationController(destinationService)
	//destination
	router.GET("/api/v1/destination", destinationController.FindAll)
	router.POST("/api/v1/destination", destinationController.Create)
	router.GET("/api/v1/destination/:destinationId", destinationController.FindOne)
	router.PUT("/api/v1/destination/:destinationId", destinationController.Update)
	router.DELETE("/api/v1/destination/:destinationId", destinationController.Destroy)
}
