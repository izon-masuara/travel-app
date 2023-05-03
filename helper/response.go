package helper

import (
	"encoding/json"
	"net/http"
)

func Response(w http.ResponseWriter, response interface{}) {
	w.Header().Add("Content-Type", "application/json")
	endcoder := json.NewEncoder(w)
	err := endcoder.Encode(response)
	PanicIfError(err)
}
