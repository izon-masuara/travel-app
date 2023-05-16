package helper

import (
	"context"
	"encoding/json"
	"kautsar/travel-app-api/entity/domain"
)

func InterfaceToJsonAuth(ctx context.Context) domain.Auth {
	bit, err := json.Marshal(ctx.Value("auth"))
	PanicIfError(err)
	var auth domain.Auth
	json.Unmarshal(bit, &auth)
	return auth
}
