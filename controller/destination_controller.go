package controller

import (
	"fmt"
	"kautsar/travel-app-api/entity/web"
	"kautsar/travel-app-api/helper"
	"kautsar/travel-app-api/service"
	"math/rand"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type DestinationController interface {
	Create(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}

type DestinationControllerImpl struct {
	DestinationService service.DestinationService
}

func NewDestinationController(destinationService service.DestinationService) DestinationController {
	return &DestinationControllerImpl{
		DestinationService: destinationService,
	}
}

func (controller *DestinationControllerImpl) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	err := r.ParseMultipartForm(r.ContentLength)
	helper.PanicIfError(err)

	file, header, err := r.FormFile("image_file")
	defer file.Close()
	helper.PanicIfError(err)

	filename := fmt.Sprintf("%v-%v-%s", time.Now().Nanosecond(), rand.Intn(20), header.Filename)

	helper.SaveFile(filename, file)

	request := web.DestinationCreateRequest{
		Title:     r.FormValue("title"),
		Date:      time.Now(),
		Long:      r.FormValue("long"),
		Lat:       r.FormValue("lat"),
		Text:      r.FormValue("text"),
		ImageFile: filename,
	}

	destniationResponse := controller.DestinationService.Create(r.Context(), request)
	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   destniationResponse,
	}
	helper.Response(w, webResponse)
}
