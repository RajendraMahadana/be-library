package domain

import (
	"time"
	
	"gorm.io/gorm"
)

type Borrowing struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	BookID     uint           `gorm:"not null;index" json:"book_id"`
	BorrowDate time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"borrow_date"`
	DueDate    time.Time      `gorm:"not null" json:"due_date"`
	ReturnDate *time.Time     `json:"return_date,omitempty"`
	Status     string         `gorm:"size:20;default:borrowed;index" json:"status"` // borrowed, returned, overdue
	Fine       float64        `gorm:"default:0" json:"fine"` // denda jika terlambat
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Book Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
}

// TableName menentukan nama tabel
func (Borrowing) TableName() string {
	return "borrowings"
}

// HitungDenda menghitung denda berdasarkan lama keterlambatan
func (b *Borrowing) HitungDenda(ratePerDay float64) float64 {
	if b.ReturnDate != nil && b.ReturnDate.After(b.DueDate) {
		daysLate := int(b.ReturnDate.Sub(b.DueDate).Hours() / 24)
		if daysLate > 0 {
			return float64(daysLate) * ratePerDay
		}
	}
	return 0
}

// IsOverdue cek apakah sudah melewati batas pengembalian
func (b *Borrowing) IsOverdue() bool {
	if b.ReturnDate != nil {
		return false
	}
	return time.Now().After(b.DueDate)
}