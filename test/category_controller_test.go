package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rtanx/golang-restful-api/app"
	"github.com/rtanx/golang-restful-api/controller"
	"github.com/rtanx/golang-restful-api/helper"
	"github.com/rtanx/golang-restful-api/middleware"
	"github.com/rtanx/golang-restful-api/model/domain"
	"github.com/rtanx/golang-restful-api/repository"
	"github.com/rtanx/golang-restful-api/service"
	"github.com/stretchr/testify/assert"
)

const HOST = "localhost"
const PORT = 8080

func newTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/learn_golang_restful_api_test")
	helper.PanicfIfErr(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(60 * time.Minute)
	db.SetConnMaxLifetime(10 * time.Minute)

	return db
}
func setUpRouter(db *sql.DB) http.Handler {
	log.Println("Starting integration testing ...")

	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)

	router := app.NewRouter(categoryController)

	return middleware.NewAuthMiddleware(router)
}

func truncateCategory(db *sql.DB) {
	db.Exec("TRUNCATE category")
}

func TestCreateCategorySuccess(t *testing.T) {
	url := fmt.Sprintf("http://%s:%d/api/categories", HOST, PORT)

	DB := newTestDB()
	truncateCategory(DB)
	router := setUpRouter(DB)
	body := `{
		"name": "Gadget"
	}`
	requestBody := strings.NewReader(body)
	request := httptest.NewRequest(http.MethodPost, url, requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Result()
	log.Println(resp.Status)
	assert.Equal(t, 200, resp.StatusCode)

	resBodyByte, _ := io.ReadAll(resp.Body)
	var resBody map[string]interface{}
	json.Unmarshal(resBodyByte, &resBody)
	fmt.Println(resBody)

	assert.Equal(t, 200, int(resBody["code"].(float64)))
	assert.Equal(t, "OK", resBody["status"])
	assert.Equal(t, "Gadget", resBody["data"].(map[string]interface{})["name"])

}

func TestCreateCategoryFailed(t *testing.T) {
	url := fmt.Sprintf("http://%s:%d/api/categories", HOST, PORT)

	DB := newTestDB()
	truncateCategory(DB)
	router := setUpRouter(DB)
	body := `{
		"name": ""
	}`
	requestBody := strings.NewReader(body)
	request := httptest.NewRequest(http.MethodPost, url, requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Result()
	log.Println(resp.Status)
	assert.Equal(t, 400, resp.StatusCode)

	resBodyByte, _ := io.ReadAll(resp.Body)
	var resBody map[string]interface{}
	json.Unmarshal(resBodyByte, &resBody)
	fmt.Println(resBody)

	assert.Equal(t, 400, int(resBody["code"].(float64)))
	assert.Equal(t, "Bad Request", resBody["status"])
}

func TestUpdateCategorySuccess(t *testing.T) {
	url := fmt.Sprintf("http://%s:%d/api/categories", HOST, PORT)

	DB := newTestDB()
	truncateCategory(DB)

	tx, _ := DB.Begin()
	cr := repository.NewCategoryRepository()
	c := cr.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	router := setUpRouter(DB)
	body := `{
		"name": "Gadget"
	}`
	requestBody := strings.NewReader(body)
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("%s/%d", url, c.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Result()
	log.Println(resp.Status)
	assert.Equal(t, 200, resp.StatusCode)

	resBodyByte, _ := io.ReadAll(resp.Body)
	var resBody map[string]interface{}
	json.Unmarshal(resBodyByte, &resBody)
	fmt.Println(resBody)

	assert.Equal(t, 200, int(resBody["code"].(float64)))
	assert.Equal(t, "OK", resBody["status"])
	assert.Equal(t, c.Id, int64(resBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "Gadget", resBody["data"].(map[string]interface{})["name"])
}

func TestUpdateCategoryFailed(t *testing.T) {
	url := fmt.Sprintf("http://%s:%d/api/categories", HOST, PORT)

	DB := newTestDB()
	truncateCategory(DB)

	tx, _ := DB.Begin()
	cr := repository.NewCategoryRepository()
	c := cr.Save(context.Background(), tx, domain.Category{
		Name: "",
	})
	tx.Commit()

	router := setUpRouter(DB)
	body := `{
		"name": ""
	}`
	requestBody := strings.NewReader(body)
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("%s/%d", url, c.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Result()
	log.Println(resp.Status)
	assert.Equal(t, 400, resp.StatusCode)

	resBodyByte, _ := io.ReadAll(resp.Body)
	var resBody map[string]interface{}
	json.Unmarshal(resBodyByte, &resBody)
	fmt.Println(resBody)

	assert.Equal(t, 400, int(resBody["code"].(float64)))
	assert.Equal(t, "Bad Request", resBody["status"])
}

func TestGetCategorySuccess(t *testing.T) {
	url := fmt.Sprintf("http://%s:%d/api/categories", HOST, PORT)

	DB := newTestDB()
	truncateCategory(DB)

	tx, _ := DB.Begin()
	cr := repository.NewCategoryRepository()
	c := cr.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	router := setUpRouter(DB)
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/%d", url, c.Id), nil)
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Result()
	log.Println(resp.Status)
	assert.Equal(t, 200, resp.StatusCode)

	resBodyByte, _ := io.ReadAll(resp.Body)
	var resBody map[string]interface{}
	json.Unmarshal(resBodyByte, &resBody)
	fmt.Println(resBody)

	assert.Equal(t, 200, int(resBody["code"].(float64)))
	assert.Equal(t, "OK", resBody["status"])
	assert.Equal(t, c.Id, int64(resBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, c.Name, resBody["data"].(map[string]interface{})["name"])
}

func TestGetCategoryFailed(t *testing.T) {
	url := fmt.Sprintf("http://%s:%d/api/categories", HOST, PORT)

	DB := newTestDB()
	truncateCategory(DB)

	router := setUpRouter(DB)
	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/%d", url, 1000), nil)
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Result()
	log.Println(resp.Status)
	assert.Equal(t, 404, resp.StatusCode)

	resBodyByte, _ := io.ReadAll(resp.Body)
	var resBody map[string]interface{}
	json.Unmarshal(resBodyByte, &resBody)
	fmt.Println(resBody)

	assert.Equal(t, 404, int(resBody["code"].(float64)))
	assert.Equal(t, "Not Found", resBody["status"])
}

func TestDeleteCategorySuccess(t *testing.T) {
	url := fmt.Sprintf("http://%s:%d/api/categories", HOST, PORT)

	DB := newTestDB()
	truncateCategory(DB)

	tx, _ := DB.Begin()
	cr := repository.NewCategoryRepository()
	c := cr.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	tx.Commit()

	router := setUpRouter(DB)
	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d", url, c.Id), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Result()
	log.Println(resp.Status)
	assert.Equal(t, 200, resp.StatusCode)

	resBodyByte, _ := io.ReadAll(resp.Body)
	var resBody map[string]interface{}
	json.Unmarshal(resBodyByte, &resBody)
	fmt.Println(resBody)

	assert.Equal(t, 200, int(resBody["code"].(float64)))
	assert.Equal(t, "OK", resBody["status"])
}
func TestDeleteCategoryFailed(t *testing.T) {
	url := fmt.Sprintf("http://%s:%d/api/categories", HOST, PORT)

	DB := newTestDB()
	truncateCategory(DB)

	router := setUpRouter(DB)
	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d", url, 10000), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Result()
	log.Println(resp.Status)
	assert.Equal(t, 404, resp.StatusCode)

	resBodyByte, _ := io.ReadAll(resp.Body)
	var resBody map[string]interface{}
	json.Unmarshal(resBodyByte, &resBody)
	fmt.Println(resBody)

	assert.Equal(t, 404, int(resBody["code"].(float64)))
	assert.Equal(t, "Not Found", resBody["status"])
}

func TestListCategorySuccess(t *testing.T) {
	url := fmt.Sprintf("http://%s:%d/api/categories", HOST, PORT)

	DB := newTestDB()
	truncateCategory(DB)

	tx, _ := DB.Begin()
	cr := repository.NewCategoryRepository()
	c := cr.Save(context.Background(), tx, domain.Category{
		Name: "Gadget",
	})
	c2 := cr.Save(context.Background(), tx, domain.Category{
		Name: "Computer",
	})
	tx.Commit()

	router := setUpRouter(DB)
	request := httptest.NewRequest(http.MethodGet, url, nil)
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Result()
	log.Println(resp.Status)
	assert.Equal(t, 200, resp.StatusCode)

	resBodyByte, _ := io.ReadAll(resp.Body)
	var resBody map[string]interface{}
	json.Unmarshal(resBodyByte, &resBody)
	fmt.Println(resBody)

	assert.Equal(t, 200, int(resBody["code"].(float64)))
	assert.Equal(t, "OK", resBody["status"])

	var categories []interface{} = resBody["data"].([]interface{})

	assert.Equal(t, c.Id, int64(((categories[0].(map[string]interface{}))["id"]).(float64)))
	assert.Equal(t, c.Name, ((categories[0].(map[string]interface{}))["name"]))
	assert.Equal(t, c2.Id, int64(((categories[1].(map[string]interface{}))["id"]).(float64)))
	assert.Equal(t, c2.Name, ((categories[1].(map[string]interface{}))["name"]))
}

func TestUnauthorized(t *testing.T) {
	url := fmt.Sprintf("http://%s:%d/api/categories", HOST, PORT)

	DB := newTestDB()
	truncateCategory(DB)

	router := setUpRouter(DB)
	request := httptest.NewRequest(http.MethodGet, url, nil)
	request.Header.Add("X-API-KEY", "KEY_SALAH")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	resp := recorder.Result()
	log.Println(resp.Status)
	assert.Equal(t, 401, resp.StatusCode)

	resBodyByte, _ := io.ReadAll(resp.Body)
	var resBody map[string]interface{}
	json.Unmarshal(resBodyByte, &resBody)
	fmt.Println(resBody)

	assert.Equal(t, 401, int(resBody["code"].(float64)))
	assert.Equal(t, "UNAUTHORIZED", resBody["status"])
}
