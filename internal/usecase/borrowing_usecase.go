package usecase

import (
	"errors"
	"time"

	"github.com/RajendraMahadana/perpustakaan-clean/internal/domain"
	"github.com/RajendraMahadana/perpustakaan-clean/internal/dto"
	"github.com/RajendraMahadana/perpustakaan-clean/internal/repository"
)

type BorrowingUsecase interface {
	PinjamBuku(userID uint, req dto.PinjamBukuRequest) (*dto.BorrowingResponse, error)
	KembalikanBuku(userID uint, req dto.KembalikanBukuRequest) (*dto.BorrowingResponse, error)
	GetBorrowingByID(id uint, userID uint, isAdmin bool) (*dto.BorrowingResponse, error)
	GetUserBorrowings(userID uint, page, limit int) (*dto.ListBorrowingResponse, error)
	GetAllBorrowings(page, limit int, status string) (*dto.ListBorrowingResponse, error)
	GetActiveBorrowings(userID uint) ([]dto.BorrowingResponse, error)
	CheckOverdue() error
}

type borrowingUsecase struct {
	borrowingRepo repository.BorrowingRepository
	bookRepo      repository.BookRepository
	userRepo      repository.UserRepository
}

func NewBorrowingUsecase(
	borrowingRepo repository.BorrowingRepository,
	bookRepo repository.BookRepository,
	userRepo repository.UserRepository,
) BorrowingUsecase {
	return &borrowingUsecase{
		borrowingRepo: borrowingRepo,
		bookRepo:      bookRepo,
		userRepo:      userRepo,
	}
}

// PinjamBuku proses peminjaman buku
func (u *borrowingUsecase) PinjamBuku(userID uint, req dto.PinjamBukuRequest) (*dto.BorrowingResponse, error) {
	// Validasi user
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	// Validasi buku
	book, err := u.bookRepo.FindByID(req.BookID)
	if err != nil {
		return nil, errors.New("buku tidak ditemukan")
	}

	// Cek stok tersedia
	if book.Available <= 0 {
		return nil, errors.New("stok buku tidak tersedia")
	}

	// Cek apakah user masih punya pinjaman aktif
	activeCount, err := u.borrowingRepo.CountActiveByUserID(userID)
	if err != nil {
		return nil, err
	}
	if activeCount >= 3 { // Maksimal 3 buku per user
		return nil, errors.New("anda sudah mencapai batas maksimal peminjaman (3 buku)")
	}

	// Hitung tanggal jatuh tempo
	borrowDate := time.Now()
	dueDate := borrowDate.AddDate(0, 0, req.Duration)

	// Buat record peminjaman
	borrowing := &domain.Borrowing{
		UserID:     userID,
		BookID:     req.BookID,
		BorrowDate: borrowDate,
		DueDate:    dueDate,
		Status:     "borrowed",
		Fine:       0,
	}

	// Gunakan transaksi database
	err = u.bookRepo.WithTransaction(func(tx interface{}) error {
		// Kurangi stok tersedia
		book.Available--
		if err := u.bookRepo.Update(book); err != nil {
			return err
		}

		// Simpan peminjaman
		return u.borrowingRepo.Create(borrowing)
	})

	if err != nil {
		return nil, err
	}

	return u.toResponse(borrowing, user, book), nil
}

// KembalikanBuku proses pengembalian buku
func (u *borrowingUsecase) KembalikanBuku(userID uint, req dto.KembalikanBukuRequest) (*dto.BorrowingResponse, error) {
	// Cari data peminjaman
	borrowing, err := u.borrowingRepo.FindByID(req.BorrowingID)
	if err != nil {
		return nil, errors.New("data peminjaman tidak ditemukan")
	}

	// Validasi kepemilikan (kecuali admin)
	if borrowing.UserID != userID {
		// Cek apakah user adalah admin (perlu logic untuk cek role)
		user, _ := u.userRepo.FindByID(userID)
		if user.Role != "admin" {
			return nil, errors.New("anda tidak berhak mengembalikan buku ini")
		}
	}

	// Cek status
	if borrowing.Status != "borrowed" && borrowing.Status != "overdue" {
		return nil, errors.New("buku sudah dikembalikan sebelumnya")
	}

	// Proses pengembalian
	now := time.Now()
	borrowing.ReturnDate = &now
	borrowing.Status = "returned"

	// Hitung denda jika terlambat
	if now.After(borrowing.DueDate) {
		daysLate := int(now.Sub(borrowing.DueDate).Hours() / 24)
		borrowing.Fine = float64(daysLate) * 1000 // Rp 1000 per hari
	}

	// Update menggunakan transaksi
	err = u.bookRepo.WithTransaction(func(tx interface{}) error {
		// Update status peminjaman
		if err := u.borrowingRepo.Update(borrowing); err != nil {
			return err
		}

		// Tambah stok tersedia buku
		book, err := u.bookRepo.FindByID(borrowing.BookID)
		if err != nil {
			return err
		}
		book.Available++
		return u.bookRepo.Update(book)
	})

	if err != nil {
		return nil, err
	}

	// Ambil data lengkap untuk response
	user, _ := u.userRepo.FindByID(borrowing.UserID)
	book, _ := u.bookRepo.FindByID(borrowing.BookID)

	return u.toResponse(borrowing, user, book), nil
}

// GetBorrowingByID ambil detail peminjaman
func (u *borrowingUsecase) GetBorrowingByID(id uint, userID uint, isAdmin bool) (*dto.BorrowingResponse, error) {
	borrowing, err := u.borrowingRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("data peminjaman tidak ditemukan")
	}

	// Validasi akses (hanya pemilik atau admin)
	if !isAdmin && borrowing.UserID != userID {
		return nil, errors.New("akses ditolak")
	}

	user, _ := u.userRepo.FindByID(borrowing.UserID)
	book, _ := u.bookRepo.FindByID(borrowing.BookID)

	return u.toResponse(borrowing, user, book), nil
}

// GetUserBorrowings ambil semua peminjaman user
func (u *borrowingUsecase) GetUserBorrowings(userID uint, page, limit int) (*dto.ListBorrowingResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	borrowings, total, err := u.borrowingRepo.FindByUserID(userID, page, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.BorrowingResponse, len(borrowings))
	for i, b := range borrowings {
		responses[i] = *u.toResponse(&b, nil, &b.Book)
	}

	return &dto.ListBorrowingResponse{
		Borrowings: responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
	}, nil
}

// GetAllBorrowings ambil semua peminjaman (admin only)
func (u *borrowingUsecase) GetAllBorrowings(page, limit int, status string) (*dto.ListBorrowingResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	borrowings, total, err := u.borrowingRepo.FindAll(page, limit, status)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.BorrowingResponse, len(borrowings))
	for i, b := range borrowings {
		responses[i] = *u.toResponse(&b, &b.User, &b.Book)
	}

	return &dto.ListBorrowingResponse{
		Borrowings: responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
	}, nil
}

// GetActiveBorrowings ambil peminjaman aktif user
func (u *borrowingUsecase) GetActiveBorrowings(userID uint) ([]dto.BorrowingResponse, error) {
	borrowings, err := u.borrowingRepo.FindActiveByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.BorrowingResponse, len(borrowings))
	for i, b := range borrowings {
		responses[i] = *u.toResponse(&b, nil, &b.Book)
	}

	return responses, nil
}

// CheckOverdue cek dan update status peminjaman yang terlambat
func (u *borrowingUsecase) CheckOverdue() error {
	borrowings, err := u.borrowingRepo.FindOverdue()
	if err != nil {
		return err
	}

	for _, b := range borrowings {
		b.Status = "overdue"
		u.borrowingRepo.Update(&b)
	}

	return nil
}

// Helper: konversi ke response DTO
func (u *borrowingUsecase) toResponse(b *domain.Borrowing, user *domain.User, book *domain.Book) *dto.BorrowingResponse {
	resp := &dto.BorrowingResponse{
		ID:         b.ID,
		UserID:     b.UserID,
		BookID:     b.BookID,
		BorrowDate: b.BorrowDate,
		DueDate:    b.DueDate,
		ReturnDate: b.ReturnDate,
		Status:     b.Status,
		Fine:       b.Fine,
		IsOverdue:  b.IsOverdue(),
	}

	if user != nil {
		resp.UserName = user.Name
	}
	if book != nil {
		resp.BookTitle = book.Title
	}

	// Hitung sisa hari
	if b.ReturnDate == nil {
		daysLeft := int(b.DueDate.Sub(time.Now()).Hours() / 24)
		resp.DaysLeft = daysLeft
	}

	return resp
}
