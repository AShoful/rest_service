package service

import (
	"rest/models"
	"rest/pkg/repository"
)

type Service struct {
	Authorization
	Book
}

type Authorization interface {
	CreateUser(user models.User) (int, error)
}

type Book interface {
	Create(userId int, book models.Book) (int, error)
	GetAll(userId int) ([]models.Book, error)
	GetById(userId, bookId int) (models.Book, error)
	Delete(userId, bookId int) error
	Update(userId, bookId int, input models.Book) error
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
