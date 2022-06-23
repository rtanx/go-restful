package middleware

import (
	"net/http"

	"github.com/rtanx/golang-restful-api/helper"
	webresponse "github.com/rtanx/golang-restful-api/model/web/response"
)

type AuthMiddleware struct {
	Handler http.Handler
}

func NewAuthMiddleware(handler http.Handler) *AuthMiddleware {
	return &AuthMiddleware{Handler: handler}
}

func (middleware *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-API-KEY") == "RAHASIA" {
		middleware.Handler.ServeHTTP(w, r)

	} else {
		status := http.StatusUnauthorized

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		resp := webresponse.WebResponse{
			Code:   status,
			Status: "UNAUTHORIZED",
		}
		helper.WriteToResponseBody(w, resp)
	}
}
