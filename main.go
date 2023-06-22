package main

import (
	"context"
	"fmt"
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

	mongoDb := app.NewMongo()
	mysqlDb := app.NewMysql()
	validate := validator.New()
	router := httprouter.New()

	defer mongoDb.Client().Disconnect(context.Background())
	defer mysqlDb.Close()

	//user
	app.RegisterUserRoutes(router, mongoDb, validate)
	// operator
	app.RegisterAdminRoutes(router, mongoDb, validate)
	// destination
	app.RegisterDestinationRoutes(router, mongoDb, validate)
	// Tiket
	app.RegisterTiketRoutes(router, mysqlDb)

	router.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr:    os.Getenv("PORT"),
		Handler: middleware.NewAuthMiddleware(router),
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
	fmt.Println("app running")
}
