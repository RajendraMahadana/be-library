package http

import (
	"net/http"
	"strconv"

	"github.com/RajendraMahadana/perpustakaan-clean/internal/domain"
	"github.com/RajendraMahadana/perpustakaan-clean/internal/dto"
	"github.com/RajendraMahadana/perpustakaan-clean/internal/usecase"
	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	usecase *usecase.BookUsecase
}

func NewBookHandler(u *usecase.BookUsecase) *BookHandler {
	return &BookHandler{usecase: u}
}

// 🔴 UPDATE: Create buku dengan field Year
func (h *BookHandler) Create(c *gin.Context) {
	var req dto.CreateBookRequest

	if err := c.ShouldBindJSON(&req); err != nil { // Gunakan ShouldBindJSON, bukan ShouldBindBodyWithJSON
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := domain.Book{
		Title:  req.Title,
		Author: req.Author,
		ISBN:   req.ISBN,
		Year:   req.Year,      // 🔴 TAMBAHKAN FIELD YEAR
		Stock:  req.Stock,
		// Available akan otomatis diisi = Stock melalui BeforeCreate hook di domain
	}

	if err := h.usecase.CreateBook(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Buku berhasil dibuat",
		"data":    book,
	})
}

// 🔴 UPDATE: GetAllBooks dengan response yang lengkap
func (h *BookHandler) GetAllBooks(c *gin.Context) {
	books, err := h.usecase.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.BookResponse

	for _, b := range books {
		response = append(response, dto.BookResponse{
			ID:        b.ID,
			Title:     b.Title,
			Author:    b.Author,
			ISBN:      b.ISBN,
			Year:      b.Year,      // 🔴 TAMBAHKAN FIELD YEAR
			Stock:     b.Stock,
			Available: b.Available,  // 🔴 TAMBAHKAN FIELD AVAILABLE
		})
	}

	c.JSON(http.StatusOK, response)
}

// 🔴 UPDATE: GetBook dengan pagination
func (h *BookHandler) GetBook(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	search := c.Query("search")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	books, total, err := h.usecase.GetBooks(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.BookResponse

	for _, b := range books {
		response = append(response, dto.BookResponse{
			ID:        b.ID,
			Title:     b.Title,
			Author:    b.Author,
			ISBN:      b.ISBN,
			Year:      b.Year,
			Stock:     b.Stock,
			Available: b.Available,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"total": total,
		"data":  response,
	})
}

// 🔴 UPDATE: GetByID dengan response lengkap
func (h *BookHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	book, err := h.usecase.GetBookByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	response := dto.BookResponse{
		ID:        book.ID,
		Title:     book.Title,
		Author:    book.Author,
		ISBN:      book.ISBN,
		Year:      book.Year,
		Stock:     book.Stock,
		Available: book.Available,
	}

	c.JSON(http.StatusOK, response)
}

// 🔴 UPDATE: Update dengan field Year dan Available
func (h *BookHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req dto.UpdateBookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil buku yang ada terlebih dahulu
	existingBook, err := h.usecase.GetBookByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Buku tidak ditemukan"})
		return
	}

	// Update field yang dikirim
	if req.Title != "" {
		existingBook.Title = req.Title
	}
	if req.Author != "" {
		existingBook.Author = req.Author
	}
	if req.ISBN != "" {
		existingBook.ISBN = req.ISBN
	}
	if req.Year != 0 {
		existingBook.Year = req.Year
	}
	if req.Stock != 0 {
		existingBook.Stock = req.Stock
	}
	if req.Available != 0 {
		existingBook.Available = req.Available
	}

	if err := h.usecase.UpdateBook(existingBook); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Buku berhasil diperbarui",
		"data":    existingBook,
	})
}

func (h *BookHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.usecase.DeleteBook(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Buku berhasil dihapus"})
}