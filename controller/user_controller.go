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
