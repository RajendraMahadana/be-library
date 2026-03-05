package repository

import "github.com/RajendraMahadana/perpustakaan-clean/internal/domain"

type BookRepository interface {
	Create(book *domain.Book) error
	FindAll() ([]domain.Book, error)
	FindByID(id uint) (*domain.Book, error)
	Update(book *domain.Book) error
	Delete(id uint) error
}
