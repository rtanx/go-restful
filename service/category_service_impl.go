package service

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/rtanx/golang-restful-api/exception"
	"github.com/rtanx/golang-restful-api/helper"
	"github.com/rtanx/golang-restful-api/model/domain"
	webrequest "github.com/rtanx/golang-restful-api/model/web/request"
	webresponse "github.com/rtanx/golang-restful-api/model/web/response"
	"github.com/rtanx/golang-restful-api/repository"
)

type CategoryServiceImpl struct {
	CategoryRepository repository.CategoryRepository
	DB                 *sql.DB
	Validate           *validator.Validate
}

func NewCategoryService(categoryRepository repository.CategoryRepository, DB *sql.DB, validate *validator.Validate) CategoryService {
	return &CategoryServiceImpl{
		CategoryRepository: categoryRepository,
		DB:                 DB,
		Validate:           validate,
	}
}

func (service *CategoryServiceImpl) Create(ctx context.Context, request webrequest.CategoryCreateRequest) webresponse.CategoryResponse {
	err := service.Validate.Struct(request)
	helper.PanicfIfErr(err)

	tx, err := service.DB.Begin()
	helper.PanicfIfErr(err)

	defer helper.CommitOrRollback(tx)

	category := domain.Category{
		Name: request.Name,
	}
	category = service.CategoryRepository.Save(ctx, tx, category)
	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) Update(ctx context.Context, request webrequest.CategoryUpdateRequest) webresponse.CategoryResponse {
	err := service.Validate.Struct(request)
	helper.PanicfIfErr(err)

	tx, err := service.DB.Begin()
	helper.PanicfIfErr(err)

	defer helper.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, request.Id)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	category.Name = request.Name

	category = service.CategoryRepository.Update(ctx, tx, category)
	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) Delete(ctx context.Context, categoryId int64) {
	tx, err := service.DB.Begin()
	helper.PanicfIfErr(err)

	category, err := service.CategoryRepository.FindById(ctx, tx, categoryId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	service.CategoryRepository.Delete(ctx, tx, category)

	defer helper.CommitOrRollback(tx)

	service.CategoryRepository.Delete(ctx, tx, category)
}

func (service *CategoryServiceImpl) FindById(ctx context.Context, categoryId int64) webresponse.CategoryResponse {
	tx, err := service.DB.Begin()
	helper.PanicfIfErr(err)

	defer helper.CommitOrRollback(tx)

	category, err := service.CategoryRepository.FindById(ctx, tx, categoryId)
	if err != nil {
		panic(exception.NewNotFoundError(err.Error()))
	}

	return helper.ToCategoryResponse(category)
}

func (service *CategoryServiceImpl) FindAll(ctx context.Context) []webresponse.CategoryResponse {
	tx, err := service.DB.Begin()
	helper.PanicfIfErr(err)

	defer helper.CommitOrRollback(tx)

	categories := service.CategoryRepository.FindAll(ctx, tx)

	return helper.ToCategoriesResponse(categories)
}
