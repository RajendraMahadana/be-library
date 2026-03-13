package dto

import (
	"time"
)

// PinjamBukuRequest DTO untuk request peminjaman
type PinjamBukuRequest struct {
	BookID     uint   `json:"book_id" binding:"required"`
	Duration   int    `json:"duration" binding:"required,min=1,max=30"` // lama pinjam dalam hari
}

// KembalikanBukuRequest DTO untuk pengembalian
type KembalikanBukuRequest struct {
	BorrowingID uint `json:"borrowing_id" binding:"required"`
}

// BorrowingResponse DTO untuk response
type BorrowingResponse struct {
	ID          uint       `json:"id"`
	UserID      uint       `json:"user_id"`
	UserName    string     `json:"user_name,omitempty"`
	BookID      uint       `json:"book_id"`
	BookTitle   string     `json:"book_title,omitempty"`
	BorrowDate  time.Time  `json:"borrow_date"`
	DueDate     time.Time  `json:"due_date"`
	ReturnDate  *time.Time `json:"return_date,omitempty"`
	Status      string     `json:"status"`
	Fine        float64    `json:"fine"`
	DaysLeft    int        `json:"days_left,omitempty"`
	IsOverdue   bool       `json:"is_overdue"`
}

// ListBorrowingResponse untuk response list
type ListBorrowingResponse struct {
	Borrowings []BorrowingResponse `json:"borrowings"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
}