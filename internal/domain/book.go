package domain

import (
	"time"
	"gorm.io/gorm"
)

type Book struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:255;not null;index" json:"title"`
	Author    string         `gorm:"size:255;not null;index" json:"author"`
	ISBN      string         `gorm:"size:100;uniqueIndex" json:"isbn"`
	Year      int            `gorm:"not null" json:"year"`       // 🔴 TAMBAHKAN FIELD YEAR
	Stock     int            `gorm:"default:0" json:"stock"`
	Available int            `gorm:"default:0" json:"available"` // Stok tersedia untuk dipinjam
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi (opsional)
	Borrowings []Borrowing `json:"borrowings,omitempty" gorm:"foreignKey:BookID"`
}

// TableName menentukan nama tabel di database
func (Book) TableName() string {
	return "books"
}

// BeforeCreate hook untuk inisialisasi
func (b *Book) BeforeCreate(tx *gorm.DB) error {
	if b.Available == 0 && b.Stock > 0 {
		b.Available = b.Stock // Set available = stock saat pertama dibuat
	}
	return nil
}

// BeforeUpdate hook untuk memastikan available tidak melebihi stock
func (b *Book) BeforeUpdate(tx *gorm.DB) error {
	if b.Available > b.Stock {
		b.Available = b.Stock
	}
	return nil
}

// KurangiStok mengurangi stok tersedia (dipanggil saat peminjaman)
func (b *Book) KurangiStok() error {
	if b.Available <= 0 {
		return ErrStokHabis
	}
	b.Available--
	return nil
}

// TambahStok menambah stok tersedia (dipanggil saat pengembalian)
func (b *Book) TambahStok() error {
	if b.Available >= b.Stock {
		return ErrStokPenuh
	}
	b.Available++
	return nil
}

// IsAvailable cek apakah buku tersedia
func (b *Book) IsAvailable() bool {
	return b.Available > 0
}