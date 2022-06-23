package exception

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rtanx/golang-restful-api/helper"
	webresponse "github.com/rtanx/golang-restful-api/model/web/response"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	switch e := err.(type) {
	case NotFoundError:
		notFoundError(w, r, e)
	case validator.ValidationErrors:
		validationErrors(w, r, e)
	default:
		internalServerError(w, r, e)
	}
}

func validationErrors(w http.ResponseWriter, r *http.Request, err validator.ValidationErrors) {
	status := http.StatusBadRequest

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := webresponse.WebResponse{
		Code:   status,
		Status: "Bad Request",
		Data:   err.Error(),
	}
	helper.WriteToResponseBody(w, resp)
}

func notFoundError(w http.ResponseWriter, r *http.Request, err NotFoundError) {

	status := http.StatusNotFound

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := webresponse.WebResponse{
		Code:   status,
		Status: "Not Found",
		Data:   err.Error,
	}
	helper.WriteToResponseBody(w, resp)
}

func internalServerError(w http.ResponseWriter, r *http.Request, err interface{}) {
	status := http.StatusInternalServerError

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := webresponse.WebResponse{
		Code:   status,
		Status: "Internal Server Error",
		Data:   err,
	}
	helper.WriteToResponseBody(w, resp)
}
