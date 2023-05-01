package main

import (
	"kautsar/travel-app-api/app"
	"kautsar/travel-app-api/controller"
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
	operatorRepository := repository.NewOperatorRepository()
	operatorService := service.NewOperatorService(operatorRepository, db, validate)
	operatorController := controller.NewOperatorController(operatorService)

	router := httprouter.New()

	router.GET("/api/v1/account", operatorController.FindAll)
	router.POST("/api/v1/account", operatorController.Create)
	router.PATCH("/api/v1/account/:accountId", operatorController.ResetPassword)
	router.DELETE("/api/v1/account/:accountId", operatorController.Destroy)

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
