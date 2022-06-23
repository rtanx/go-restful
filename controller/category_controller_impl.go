package controller

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rtanx/golang-restful-api/helper"
	webrequest "github.com/rtanx/golang-restful-api/model/web/request"
	webresponse "github.com/rtanx/golang-restful-api/model/web/response"
	"github.com/rtanx/golang-restful-api/service"
)

type CategoryControllerImpl struct {
	CategoryService service.CategoryService
}

func NewCategoryController(categoryService service.CategoryService) CategoryController {
	return &CategoryControllerImpl{
		CategoryService: categoryService,
	}
}

func (controller *CategoryControllerImpl) Create(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryCreateRequest := webrequest.CategoryCreateRequest{}
	helper.ReadFromRequestBody(request, &categoryCreateRequest)

	categoryResponse := controller.CategoryService.Create(request.Context(), categoryCreateRequest)
	webResponse := webresponse.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *CategoryControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryUpdateRequest := webrequest.CategoryUpdateRequest{}
	helper.ReadFromRequestBody(request, &categoryUpdateRequest)

	id, err := strconv.ParseInt(params.ByName("categoryId"), 10, 64)
	helper.PanicfIfErr(err)

	categoryUpdateRequest.Id = id

	categoryResponse := controller.CategoryService.Update(request.Context(), categoryUpdateRequest)
	webResponse := webresponse.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)

}

func (controller *CategoryControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryId, err := strconv.ParseInt(params.ByName("categoryId"), 10, 64)
	helper.PanicfIfErr(err)

	controller.CategoryService.Delete(request.Context(), categoryId)
	webResponse := webresponse.WebResponse{
		Code:   200,
		Status: "OK",
	}
	helper.WriteToResponseBody(writer, webResponse)

}

func (controller *CategoryControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoryId, err := strconv.ParseInt(params.ByName("categoryId"), 10, 64)
	helper.PanicfIfErr(err)
	categoryResponse := controller.CategoryService.FindById(request.Context(), categoryId)
	webResponse := webresponse.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoryResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)

}

func (controller *CategoryControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	categoriesResponse := controller.CategoryService.FindAll(request.Context())
	webResponse := webresponse.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   categoriesResponse,
	}
	helper.WriteToResponseBody(writer, webResponse)

}
