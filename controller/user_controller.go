package controller

import (
	"encoding/json"
	"kautsar/travel-app-api/entity/web"
	"kautsar/travel-app-api/helper"
	"kautsar/travel-app-api/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UserController interface {
	Login(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FindDestinationByRegion(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FindOneDestinationByRegion(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FindAllRegions(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	FindImage(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}

type UserControllerImpl struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

func (controller *UserControllerImpl) Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	userLoginRequest := web.LoginRequest{}
	err := decoder.Decode(&userLoginRequest)
	helper.PanicIfError(err)

	loginResponse := controller.UserService.Login(r.Context(), userLoginRequest)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   loginResponse,
	}
	helper.Response(w, webResponse)
}

func (controller *UserControllerImpl) FindDestinationByRegion(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("region_name")
	destinationsResponse := controller.UserService.FindDestinationByRegion(r.Context(), id)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   destinationsResponse,
	}
	helper.Response(w, webResponse)
}

func (controller *UserControllerImpl) FindOneDestinationByRegion(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	regionName := params.ByName("region_name")
	destinationId := params.ByName("destination_id")
	destinationsResponse := controller.UserService.FindOneDestinationByRegion(r.Context(), regionName, destinationId)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   destinationsResponse,
	}
	helper.Response(w, webResponse)
}

func (controller *UserControllerImpl) FindAllRegions(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	regions := controller.UserService.FindAllRegions(r.Context())
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   regions,
	}
	helper.Response(w, webResponse)
}

func (controller *UserControllerImpl) FindImage(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	imgName := params.ByName("image_file")
	byteFile := controller.UserService.FindImage(imgName)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(byteFile)
}
