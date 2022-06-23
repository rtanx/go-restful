package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rtanx/golang-restful-api/app"
	"github.com/rtanx/golang-restful-api/controller"
	"github.com/rtanx/golang-restful-api/db"
	"github.com/rtanx/golang-restful-api/helper"
	"github.com/rtanx/golang-restful-api/middleware"
	"github.com/rtanx/golang-restful-api/repository"
	"github.com/rtanx/golang-restful-api/service"
)

const HOST = "localhost"
const PORT = 8080

func main() {
	log.Printf("Starting Application on port :%d", PORT)

	DB := db.NewDB()
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, DB, validate)
	categoryController := controller.NewCategoryController(categoryService)

	router := app.NewRouter(categoryController)

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", HOST, PORT),
		Handler: middleware.NewAuthMiddleware(router),
	}

	err := server.ListenAndServe()
	helper.PanicfIfErr(err)

}
