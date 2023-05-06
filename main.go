package main

import (
	"kautsar/travel-app-api/app"
	"kautsar/travel-app-api/exception"
	"kautsar/travel-app-api/helper"
	"kautsar/travel-app-api/middleware"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
)

func main() {

	db := app.NewDb()
	validate := validator.New()
	router := httprouter.New()

	//user
	app.RegisterUserRoutes(router, db, validate)
	// operator
	app.RegisterAdminRoutes(router, db, validate)
	// destination
	app.RegisterDestinationRoutes(router, db, validate)

	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: middleware.NewAuthMiddleware(router),
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
