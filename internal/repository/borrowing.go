package repository

import (
	"github.com/RajendraMahadana/perpustakaan-clean/internal/domain"
)

type BorrowingRepository interface {
	Create(borrowing *domain.Borrowing) error
	FindByID(id uint) (*domain.Borrowing, error)
	FindByUserID(userID uint, page, limit int) ([]domain.Borrowing, int64, error)
	FindActiveByUserID(userID uint) ([]domain.Borrowing, error)
	FindAll(page, limit int, status string) ([]domain.Borrowing, int64, error)
	FindOverdue() ([]domain.Borrowing, error)
	Update(borrowing *domain.Borrowing) error
	Delete(id uint) error
	CountActiveByUserID(userID uint) (int64, error)
	CountByBookID(bookID uint) (int64, error)
}