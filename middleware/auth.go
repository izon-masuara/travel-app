package middleware

import (
	"context"
	"kautsar/travel-app-api/entity/web"
	"kautsar/travel-app-api/helper"
	"net/http"
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH ,PUT, DELETE, OPTIONS")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	token := r.Header.Get("TOKEN")

	if token == "" {
		middleware.Handler.ServeHTTP(w, r)
		return
	}

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
