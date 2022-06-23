package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rtanx/golang-restful-api/helper"
	"github.com/rtanx/golang-restful-api/model/domain"
)

type CategoryRepositoryImpl struct {
}

func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryImpl{}
}

func (respository *CategoryRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category {
	SQL := "INSERT INTO category(name) VALUES (?)"
	res, err := tx.ExecContext(ctx, SQL, category.Name)
	helper.PanicfIfErr(err)

	id, err := res.LastInsertId()

	helper.PanicfIfErr(err)

	category.Id = id
	return category

}

func (respository *CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category {
	SQL := "UPDATE category SET name = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, category.Name, category.Id)
	helper.PanicfIfErr(err)
	return category
}

func (respository *CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, category domain.Category) {
	SQL := "DELETE FROM category WHERE id = ?"
	_, err := tx.ExecContext(ctx, SQL, category.Id)
	helper.PanicfIfErr(err)
}

func (respository *CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, categoryId int64) (domain.Category, error) {
	SQL := "SELECT id, name FROM category WHERE id = ?"
	resRows, err := tx.QueryContext(ctx, SQL, categoryId)

	helper.PanicfIfErr(err)
	defer resRows.Close()

	category := domain.Category{}
	if resRows.Next() {
		err = resRows.Scan(&category.Id, &category.Name)
		helper.PanicfIfErr(err)
		return category, nil
	} else {
		return category, errors.New("category is not found")
	}
}

func (respository *CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Category {
	SQL := "SELECT * FROM category"
	resRows, err := tx.QueryContext(ctx, SQL)
	helper.PanicfIfErr(err)
	defer resRows.Close()

	var categories []domain.Category
	for resRows.Next() {
		category := domain.Category{}
		err := resRows.Scan(&category.Id, &category.Name)
		helper.PanicfIfErr(err)
		categories = append(categories, category)
	}
	return categories
}
