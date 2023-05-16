package main

import (
	"context"
	"kautsar/travel-app-api/app"
	"kautsar/travel-app-api/exception"
	"kautsar/travel-app-api/helper"
	"kautsar/travel-app-api/middleware"
	"net/http"
	"os"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {

	err := godotenv.Load()
	helper.PanicIfError(err)

	db := app.NewDb()
	validate := validator.New()
	router := httprouter.New()

	defer db.Client().Disconnect(context.Background())

	//user
	app.RegisterUserRoutes(router, db, validate)
	// operator
	app.RegisterAdminRoutes(router, db, validate)
	// destination
	app.RegisterDestinationRoutes(router, db, validate)

	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr:    os.Getenv("IP_ADDRESS"),
		Handler: middleware.NewAuthMiddleware(router),
	}

	err = server.ListenAndServe()
	helper.PanicIfError(err)
}
