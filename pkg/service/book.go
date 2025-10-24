package service

import (
	"rest/models"
	"rest/pkg/repository"
)

type BookService struct {
	repo repository.Book
}

func NewBookService(repo repository.Book) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) Create(book models.Book) (uint, error) {
	return s.repo.Create(book)
}

func (s *BookService) GetAll() ([]models.Book, error) {
	return s.repo.GetAll()
}

func (s *BookService) GetById(bookId uint) (models.Book, error) {
	return s.repo.GetById(bookId)
}

func (s *BookService) Delete(userId uint, bookId uint) error {
	return s.repo.Delete(userId, bookId)
}

func (s *BookService) Update(userId, bookId uint, book models.UpdateBook) error {
	return s.repo.Update(userId, bookId, book)
}
