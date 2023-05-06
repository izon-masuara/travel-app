package middleware

import (
	"context"
	"kautsar/travel-app-api/entity/web"
	"kautsar/travel-app-api/helper"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{
		Handler: handler,
	}
}

func (middleware *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	if path[3] == "login" {
		middleware.Handler.ServeHTTP(w, r)
		return
	}
	token := r.Header.Get("TOKEN")
	data, err := helper.ValidateToken(token)
	if err != nil {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		webResponse := web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
		}
		helper.Response(w, webResponse)
		return
	}
	ctx := context.WithValue(r.Context(), "auth", data)
	middleware.Handler.ServeHTTP(w, r.WithContext(ctx))
}
