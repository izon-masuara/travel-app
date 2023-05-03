package main

import (
	"kautsar/travel-app-api/app"
	"kautsar/travel-app-api/controller"
	"kautsar/travel-app-api/exception"
	"kautsar/travel-app-api/helper"
	"kautsar/travel-app-api/repository"
	"kautsar/travel-app-api/service"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
)

func main() {

	db := app.NewDb()
	validate := validator.New()

	// operator init
	operatorRepository := repository.NewOperatorRepository()
	operatorService := service.NewOperatorService(operatorRepository, db, validate)
	operatorController := controller.NewOperatorController(operatorService)

	// destination init
	destinationRepository := repository.NewDestinationRepository()
	destinationService := service.NewDestinationService(destinationRepository, db, validate)
	destinationController := controller.NewDestinationController(destinationService)

	router := httprouter.New()

	//account
	router.GET("/api/v1/account", operatorController.FindAll)
	router.POST("/api/v1/account", operatorController.Create)
	router.PATCH("/api/v1/account/:accountId", operatorController.ResetPassword)
	router.DELETE("/api/v1/account/:accountId", operatorController.Destroy)

	//destination
	router.GET("/api/v1/destination", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte("Hello world"))
	})
	router.POST("/api/v1/destination", destinationController.Create)

	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
