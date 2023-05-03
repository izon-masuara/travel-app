package controller

import (
	"encoding/json"
	"kautsar/travel-app-api/entity/web"
	"kautsar/travel-app-api/helper"
	"kautsar/travel-app-api/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type OperatorController interface {
	Create(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	ResetPassword(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	Destroy(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}

type OperatorControllerImpl struct {
	OperatorService service.OperatorService
}

func NewOperatorController(operatorService service.OperatorService) OperatorController {
	return &OperatorControllerImpl{
		OperatorService: operatorService,
	}
}

func (controller *OperatorControllerImpl) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	operatorCreateRequest := web.OperatorCreateRequest{}
	err := decoder.Decode(&operatorCreateRequest)
	helper.PanicIfError(err)

	operatorResponse := controller.OperatorService.Create(r.Context(), operatorCreateRequest)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   operatorResponse,
	}

	helper.Response(w, webResponse)
}

func (controller *OperatorControllerImpl) FindAll(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	operatorResponse := controller.OperatorService.FindAll(r.Context())
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   operatorResponse,
	}

	helper.Response(w, webResponse)
}

func (controller *OperatorControllerImpl) ResetPassword(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	operatorId := params.ByName("accountId")
	operatorResponse := controller.OperatorService.ResetPassword(r.Context(), operatorId)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   operatorResponse,
	}

	helper.Response(w, webResponse)
}

func (controller *OperatorControllerImpl) Destroy(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	operatorId := params.ByName("accountId")
	operatorResponse := controller.OperatorService.Destroy(r.Context(), operatorId)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   operatorResponse,
	}

	helper.Response(w, webResponse)
}
