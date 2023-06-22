package app

import (
	"database/sql"
	"kautsar/travel-app-api/controller"
	"kautsar/travel-app-api/repository"
	"kautsar/travel-app-api/service"

	"github.com/julienschmidt/httprouter"
)

func RegisterTiketRoutes(router *httprouter.Router, db *sql.DB) {
	tiketRepository := repository.NewTiketRepository()
	tiketService := service.NewTiketService(tiketRepository, db)
	tiketController := controller.NewTiketController(tiketService)

	router.GET("/api/v1/tiket/:code", tiketController.FindOne)
}
