package controller

import (
	"kautsar/travel-app-api/entity/web"
	"kautsar/travel-app-api/helper"
	"kautsar/travel-app-api/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type TiketController interface {
	FindOne(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}

type TiketControllerImpl struct {
	TiketService service.TiketService
}

func NewTiketController(tiketService service.TiketService) TiketController {
	return &TiketControllerImpl{
		TiketService: tiketService,
	}
}

func (controller *TiketControllerImpl) FindOne(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	code := params.ByName("code")
	tiketResponse := controller.TiketService.FindOne(r.Context(), code)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   tiketResponse,
	}
	helper.Response(w, webResponse)
}
