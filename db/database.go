package db

import (
	"database/sql"
	"time"

	"github.com/rtanx/golang-restful-api/helper"
)

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/learn_golang_restful_api")
	helper.PanicfIfErr(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(60 * time.Minute)
	db.SetConnMaxLifetime(10 * time.Minute)

	return db
}
