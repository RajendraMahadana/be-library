package repository

import (
	"github.com/RajendraMahadana/perpustakaan-clean/internal/domain"
	infra "github.com/RajendraMahadana/perpustakaan-clean/internal/infrastructure/database"
	repo "github.com/RajendraMahadana/perpustakaan-clean/internal/repository"
	"gorm.io/gorm"
)

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) repo.BookRepository {
	return &bookRepository{db}
}

func (r *bookRepository) Create(book *domain.Book) error {
	model := infra.BookModel{
		Title:  book.Title,
		Author: book.Author,
		ISBN:   book.ISBN,
		Stock:  book.Stock,
	}

	return r.db.Create(&model).Error
}

func (r *bookRepository) FindAll() ([]domain.Book, error) {
	var models []infra.BookModel
	err := r.db.Find(&models).Error
	if err != nil {
		return nil, err
	}

	var books []domain.Book
	for _, m := range models {
		books = append(books, domain.Book{
			ID:     m.ID,
			Title:  m.Title,
			Author: m.Author,
			ISBN:   m.ISBN,
			Stock:  m.Stock,
		})
	}

	return books, nil
}

func (r *bookRepository) FindByID(id uint) (*domain.Book, error) {
	var model infra.BookModel

	err := r.db.First(&model, id).Error
	if err != nil {
		return nil, err
	}

	book := domain.Book{
		ID:     model.ID,
		Title:  model.Title,
		Author: model.Author,
		ISBN:   model.ISBN,
		Stock:  model.Stock,
	}

	return &book, nil
}

func (r *bookRepository) Update(book *domain.Book) error {
	return r.db.Model(&infra.BookModel{}).
		Where("id = ?", book.ID).
		Updates(map[string]interface{}{
			"title":  book.Title,
			"author": book.Author,
			"isbn":   book.ISBN,
			"stock":  book.Stock,
		}).Error
}

func (r *bookRepository) Delete(id uint) error {
	return r.db.Delete(&infra.BookModel{}, id).Error
}
