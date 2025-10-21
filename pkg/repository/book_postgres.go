package repository

import (
	"rest/models"
)

type BookPostgres struct {
	db any
}

func NewBookPostgres(db any) *BookPostgres {
	return &BookPostgres{db: db}
}

func (r *BookPostgres) Create(userId int, book models.Book) /*(int, error)*/ {

}

func (r *BookPostgres) GetAll(userId int) /*([]models.Book, error)*/ {

}

func (r *BookPostgres) GetById(userId, bookId int) /*(models.Book, error)*/ {

}

func (r *BookPostgres) Delete(userId, bookId int) error {

	return nil
}

func (r *BookPostgres) Update(userId, listId int, input models.Book) error {
	return nil
}
