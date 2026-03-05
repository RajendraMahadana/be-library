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

func (h *BookHandler) Create(c *gin.Context) {
	var req dto.CreateBookRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := domain.Book{
		Title:  req.Title,
		Author: req.Author,
		ISBN:   req.ISBN,
		Stock:  req.Stock,
	}

	if err := h.usecase.CreateBook(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

}

func (h *BookHandler) GetAll(c *gin.Context) {
	books, err := h.usecase.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.BookResponse

	for _, b := range books {
		response = append(response, dto.BookResponse{
			ID:     b.ID,
			Title:  b.Title,
			Author: b.Author,
			ISBN:   b.ISBN,
			Stock:  b.Stock,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *BookHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	book, err := h.usecase.GetBookByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	response := dto.BookResponse{
		ID:     book.ID,
		Title:  book.Title,
		Author: book.Author,
		ISBN:   book.ISBN,
		Stock:  book.Stock,
	}

	c.JSON(http.StatusOK, response)
}

func (h *BookHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var req dto.UpdateBookRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := domain.Book{
		ID:     uint(id),
		Title:  req.Title,
		Author: req.Author,
		ISBN:   req.ISBN,
		Stock:  req.Stock,
	}

	if err := h.usecase.UpdateBook(&book); err !=nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "buku berhasil diperbarui",
	})
}

func (h *BookHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	if err := h.usecase.DeleteBook(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}
