package helper

import (
	"encoding/json"
	"kautsar/travel-app-api/entity/web"
	"net/http"
)

func Response(w http.ResponseWriter, operatorResponse interface{}) {
	webResponse := web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   operatorResponse,
	}

	w.Header().Add("Content-Type", "application/json")
	endcoder := json.NewEncoder(w)
	err := endcoder.Encode(webResponse)
	PanicIfError(err)
}
