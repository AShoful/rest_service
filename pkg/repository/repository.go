package repository

import (
	"rest/models"

	"gorm.io/gorm"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) /* return user, err */
}

type Book interface {
	Create(userId int, book models.Book) /*(int, error)*/
	GetAll(userId int)                   /*([]models.Book, error)*/
	GetById(userId, bookId int)          /*(models.Book, error)*/
	Delete(userId, bookId int) error
	Update(userId, bookId int, input models.Book) error
}

type Repository struct {
	Authorization
	Book
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Book:          NewBookPostgres(db),
	}
}
