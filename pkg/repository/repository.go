package repository

import (
	"rest/models"

	"gorm.io/gorm"
)

type Authorization interface {
	CreateUser(user models.User) (uint, error)
	GetUser(username string) (models.User, error)
}

type Book interface {
	Create(book models.Book) (uint, error)
	GetAll() ([]models.Book, error)
	GetById(bookId uint) (models.Book, error)
	Delete(userId uint, bookId uint) error
	Update(userId, bookId uint, input models.UpdateBook) error
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
