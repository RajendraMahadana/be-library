package usecase

import (
	"errors"

	"github.com/RajendraMahadana/perpustakaan-clean/internal/domain"
	"github.com/RajendraMahadana/perpustakaan-clean/internal/repository"
)

type BookUsecase struct {
	repo repository.BookRepository
}

func NewBookUsecase(r repository.BookRepository) *BookUsecase {
	return &BookUsecase{repo: r}
}

func (u *BookUsecase) CreateBook(book *domain.Book) error {
	if book.Title == "" {
		return errors.New("title is required")
	}

	if book.Stock < 0 {
		return errors.New("stock cannot be negative")
	}

	return u.repo.Create(book)
}

func (u *BookUsecase) GetAllBooks() ([]domain.Book, error) {
	return u.repo.FindAll()
}

func (u *BookUsecase) GetBookByID(id uint) (*domain.Book, error) {
	return u.repo.FindByID(id)
}

func (u *BookUsecase) UpdateBook(book *domain.Book) error {
	return u.repo.Update(book)
}

func (u *BookUsecase) DeleteBook(id uint) error {
	return u.repo.Delete(id)
}