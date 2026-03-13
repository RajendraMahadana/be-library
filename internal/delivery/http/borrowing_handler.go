package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/RajendraMahadana/perpustakaan-clean/internal/dto"
	"github.com/RajendraMahadana/perpustakaan-clean/internal/usecase"
	"github.com/RajendraMahadana/perpustakaan-clean/pkg/middleware"
)

type BorrowingHandler struct {
	borrowingUsecase usecase.BorrowingUsecase
}

func NewBorrowingHandler(borrowingUsecase usecase.BorrowingUsecase) *BorrowingHandler {
	return &BorrowingHandler{
		borrowingUsecase: borrowingUsecase,
	}
}

// PinjamBuku godoc
// @Summary Pinjam buku
// @Tags Borrowings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.PinjamBukuRequest true "Request peminjaman"
// @Success 201 {object} dto.BorrowingResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /borrowings [post]
func (h *BorrowingHandler) PinjamBuku(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req dto.PinjamBukuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.borrowingUsecase.PinjamBuku(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// KembalikanBuku godoc
// @Summary Kembalikan buku
// @Tags Borrowings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.KembalikanBukuRequest true "Request pengembalian"
// @Success 200 {object} dto.BorrowingResponse
// @Failure 400 {object} map[string]interface{}
// @Router /borrowings/return [post]
func (h *BorrowingHandler) KembalikanBuku(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req dto.KembalikanBukuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.borrowingUsecase.KembalikanBuku(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetBorrowingByID godoc
// @Summary Detail peminjaman
// @Tags Borrowings
// @Security BearerAuth
// @Produce json
// @Param id path int true "Borrowing ID"
// @Success 200 {object} dto.BorrowingResponse
// @Failure 404 {object} map[string]interface{}
// @Router /borrowings/{id} [get]
func (h *BorrowingHandler) GetBorrowingByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	role := middleware.GetUserRole(c)
	isAdmin := role == "admin"

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	response, err := h.borrowingUsecase.GetBorrowingByID(uint(id), userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetMyBorrowings godoc
// @Summary Riwayat peminjaman saya
// @Tags Borrowings
// @Security BearerAuth
// @Produce json
// @Param page query int false "Halaman"
// @Param limit query int false "Limit per halaman"
// @Success 200 {object} dto.ListBorrowingResponse
// @Router /borrowings/my [get]
func (h *BorrowingHandler) GetMyBorrowings(c *gin.Context) {
	userID := middleware.GetUserID(c)
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	response, err := h.borrowingUsecase.GetUserBorrowings(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetActiveBorrowings godoc
// @Summary Peminjaman aktif saya
// @Tags Borrowings
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.BorrowingResponse
// @Router /borrowings/active [get]
func (h *BorrowingHandler) GetActiveBorrowings(c *gin.Context) {
	userID := middleware.GetUserID(c)

	responses, err := h.borrowingUsecase.GetActiveBorrowings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses)
}

// GetAllBorrowings godoc
// @Summary Semua peminjaman (Admin only)
// @Tags Admin
// @Security BearerAuth
// @Produce json
// @Param page query int false "Halaman"
// @Param limit query int false "Limit per halaman"
// @Param status query string false "Filter status (borrowed/returned/overdue)"
// @Success 200 {object} dto.ListBorrowingResponse
// @Router /admin/borrowings [get]
func (h *BorrowingHandler) GetAllBorrowings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	response, err := h.borrowingUsecase.GetAllBorrowings(page, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}