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

// 🔴 UPDATE: Create dengan field Year dan Available
func (r *bookRepository) Create(book *domain.Book) error {
	model := infra.BookModel{
		Title:     book.Title,
		Author:    book.Author,
		ISBN:      book.ISBN,
		Year:      book.Year,      // 🔴 TAMBAHKAN YEAR
		Stock:     book.Stock,
		Available: book.Available,  // 🔴 TAMBAHKAN AVAILABLE
	}

	err := r.db.Create(&model).Error
	if err != nil {
		return err
	}

	// 🔴 Set ID yang tergenerate ke domain book
	book.ID = model.ID
	return nil
}

// 🔴 UPDATE: FindAll dengan field Year dan Available
func (r *bookRepository) FindAll() ([]domain.Book, error) {
	var models []infra.BookModel
	err := r.db.Find(&models).Error
	if err != nil {
		return nil, err
	}

	var books []domain.Book
	for _, m := range models {
		books = append(books, domain.Book{
			ID:        m.ID,
			Title:     m.Title,
			Author:    m.Author,
			ISBN:      m.ISBN,
			Year:      m.Year,      // 🔴 TAMBAHKAN YEAR
			Stock:     m.Stock,
			Available: m.Available,  // 🔴 TAMBAHKAN AVAILABLE
		})
	}

	return books, nil
}

// 🔴 UPDATE: GetBooks dengan field Year dan Available
func (r *bookRepository) GetBooks(page int, limit int, search string) ([]domain.Book, int64, error) {
	var models []infra.BookModel
	var total int64

	query := r.db.Model(&infra.BookModel{})

	if search != "" {
		query = query.Where("title LIKE ? OR author LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)

	offset := (page - 1) * limit

	err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	var books []domain.Book

	for _, m := range models {
		books = append(books, domain.Book{
			ID:        m.ID,
			Title:     m.Title,
			Author:    m.Author,
			ISBN:      m.ISBN,
			Year:      m.Year,
			Stock:     m.Stock,
			Available: m.Available,
		})
	}

	return books, total, nil
}

// 🔴 UPDATE: FindByID dengan field Year dan Available
func (r *bookRepository) FindByID(id uint) (*domain.Book, error) {
	var model infra.BookModel

	err := r.db.First(&model, id).Error
	if err != nil {
		return nil, err
	}

	book := &domain.Book{
		ID:        model.ID,
		Title:     model.Title,
		Author:    model.Author,
		ISBN:      model.ISBN,
		Year:      model.Year,
		Stock:     model.Stock,
		Available: model.Available,
	}

	return book, nil
}

func (r *bookRepository) Update(book *domain.Book) error {
    // Ambil data lama
    var existing infra.BookModel
    if err := r.db.First(&existing, book.ID).Error; err != nil {
        return err
    }

    updates := map[string]interface{}{}
    
    if book.Title != "" {
        updates["title"] = book.Title
    }
    if book.Author != "" {
        updates["author"] = book.Author
    }
    if book.ISBN != "" {
        updates["isbn"] = book.ISBN
    }
    if book.Year != 0 {
        updates["year"] = book.Year
    }
    
    // 🔴 INI YANG PENTING!
    if book.Stock > 0 {
        // Hitung selisih
        stockDiff := book.Stock - existing.Stock
        // Available baru = available lama + selisih
        newAvailable := existing.Available + stockDiff
        
        updates["stock"] = book.Stock
        updates["available"] = newAvailable
    }
    
    // 🔴 KALAU AVAILABLE DIKIRIM MANUAL
    if book.Available > 0 {
        updates["available"] = book.Available
    }

    return r.db.Model(&infra.BookModel{}).
        Where("id = ?", book.ID).
        Updates(updates).Error
}

// WithTransaction menjalankan fungsi dalam transaksi
func (r *bookRepository) WithTransaction(fn func(tx interface{}) error) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// 🔴 TAMBAHKAN: Method untuk UpdateWithTx (untuk transaksi)
func (r *bookRepository) UpdateWithTx(tx *gorm.DB, book *domain.Book) error {
	return tx.Model(&infra.BookModel{}).
		Where("id = ?", book.ID).
		Updates(map[string]interface{}{
			"title":     book.Title,
			"author":    book.Author,
			"isbn":      book.ISBN,
			"year":      book.Year,
			"stock":     book.Stock,
			"available": book.Available,
		}).Error
}

// 🔴 TAMBAHKAN: Method untuk BeginTransaction
func (r *bookRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *bookRepository) Delete(id uint) error {
	return r.db.Delete(&infra.BookModel{}, id).Error
}