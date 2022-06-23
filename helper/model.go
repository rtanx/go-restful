package helper

import (
	"github.com/rtanx/golang-restful-api/model/domain"
	webresponse "github.com/rtanx/golang-restful-api/model/web/response"
)

func ToCategoryResponse(category domain.Category) webresponse.CategoryResponse {
	return webresponse.CategoryResponse{
		Id:   category.Id,
		Name: category.Name,
	}
}

func ToCategoriesResponse(categories []domain.Category) []webresponse.CategoryResponse {
	var categoriesResponse []webresponse.CategoryResponse
	for _, category := range categories {
		categoriesResponse = append(categoriesResponse, ToCategoryResponse(category))
	}
	return categoriesResponse
}
