package repository

import (
	"time"

	"gorm.io/gorm"
	"github.com/RajendraMahadana/perpustakaan-clean/internal/domain"
)

type borrowingRepository struct {
	db *gorm.DB
}

func NewBorrowingRepository(db *gorm.DB) *borrowingRepository {
	return &borrowingRepository{db: db}
}

// Create membuat peminjaman baru
func (r *borrowingRepository) Create(borrowing *domain.Borrowing) error {
	return r.db.Create(borrowing).Error
}

// FindByID mencari peminjaman berdasarkan ID
func (r *borrowingRepository) FindByID(id uint) (*domain.Borrowing, error) {
	var borrowing domain.Borrowing
	err := r.db.Preload("User").Preload("Book").
		Where("id = ?", id).
		First(&borrowing).Error
	if err != nil {
		return nil, err
	}
	return &borrowing, nil
}

// FindByUserID mencari peminjaman berdasarkan user ID dengan pagination
func (r *borrowingRepository) FindByUserID(userID uint, page, limit int) ([]domain.Borrowing, int64, error) {
	var borrowings []domain.Borrowing
	var total int64

	query := r.db.Model(&domain.Borrowing{}).Where("user_id = ?", userID)

	// Hitung total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Ambil data dengan pagination
	offset := (page - 1) * limit
	err = query.Preload("Book").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&borrowings).Error

	return borrowings, total, err
}

// FindActiveByUserID mencari peminjaman aktif (belum dikembalikan)
func (r *borrowingRepository) FindActiveByUserID(userID uint) ([]domain.Borrowing, error) {
	var borrowings []domain.Borrowing
	err := r.db.Preload("Book").
		Where("user_id = ? AND status IN (?, ?)", userID, "borrowed", "overdue").
		Order("due_date ASC").
		Find(&borrowings).Error
	return borrowings, err
}

// FindAll mencari semua peminjaman dengan filter status
func (r *borrowingRepository) FindAll(page, limit int, status string) ([]domain.Borrowing, int64, error) {
	var borrowings []domain.Borrowing
	var total int64

	query := r.db.Model(&domain.Borrowing{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Preload("User").Preload("Book").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&borrowings).Error

	return borrowings, total, err
}

// FindOverdue mencari semua peminjaman yang terlambat
func (r *borrowingRepository) FindOverdue() ([]domain.Borrowing, error) {
	var borrowings []domain.Borrowing
	err := r.db.Preload("User").Preload("Book").
		Where("status = ? AND due_date < ?", "borrowed", time.Now()).
		Find(&borrowings).Error
	return borrowings, err
}

// Update mengupdate peminjaman
func (r *borrowingRepository) Update(borrowing *domain.Borrowing) error {
	return r.db.Save(borrowing).Error
}

// Delete menghapus peminjaman
func (r *borrowingRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Borrowing{}, id).Error
}

// CountActiveByUserID menghitung jumlah peminjaman aktif user
func (r *borrowingRepository) CountActiveByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Borrowing{}).
		Where("user_id = ? AND status = ?", userID, "borrowed").
		Count(&count).Error
	return count, err
}

// CountByBookID menghitung jumlah peminjaman buku
func (r *borrowingRepository) CountByBookID(bookID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Borrowing{}).
		Where("book_id = ? AND status = ?", bookID, "borrowed").
		Count(&count).Error
	return count, err
}
