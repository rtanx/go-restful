package service

import (
	"context"

	webrequest "github.com/rtanx/golang-restful-api/model/web/request"
	webresponse "github.com/rtanx/golang-restful-api/model/web/response"
)

type CategoryService interface {
	Create(ctx context.Context, request webrequest.CategoryCreateRequest) webresponse.CategoryResponse
	Update(ctx context.Context, request webrequest.CategoryUpdateRequest) webresponse.CategoryResponse
	Delete(ctx context.Context, categoryId int64)
	FindById(ctx context.Context, categoryId int64) webresponse.CategoryResponse
	FindAll(ctx context.Context) []webresponse.CategoryResponse
}
